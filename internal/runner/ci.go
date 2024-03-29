package runner

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/dropseed/deps/internal/billing"
	"github.com/dropseed/deps/internal/ci"
	"github.com/dropseed/deps/internal/config"
	"github.com/dropseed/deps/internal/git"
	"github.com/dropseed/deps/internal/hooks"
	"github.com/dropseed/deps/internal/output"
	"github.com/dropseed/deps/internal/pullrequest"
	"github.com/dropseed/deps/internal/schemaext"
	"github.com/dropseed/deps/pkg/schema"
)

type updateResult struct {
	update *Update
	err    error
}

func CI(autoconfigure bool, types []string, paths []string) error {

	api, err := billing.NewAPI()
	if err != nil {
		return err
	}

	if err := api.Validate(); err != nil {
		return err
	}

	if git.IsDirty() {
		fmt.Println(git.Status())
		return errors.New("git status must be clean to run deps ci")
	}

	repo := pullrequest.NewRepo()
	if repo == nil {
		return errors.New("Repo not found or not supported")
	}
	ciProvider := ci.NewCIProvider()

	if err := repo.CheckRequirements(); err != nil {
		return err
	}

	if autoconfigure {
		ci.BaseAutoconfigure()

		if err := ciProvider.Autoconfigure(); err != nil {
			return err
		}

		repo.Autoconfigure()
	}

	output.Debug("Fetching all branches so we can check for existing updates")
	git.Fetch()

	startingBranch := getCurrentBranch(ciProvider)

	git.Checkout(startingBranch)

	successfulUpdates := []*updateResult{}
	failedUpdates := []*updateResult{}

	if isDepsBranch := git.IsDepsBranch(startingBranch); isDepsBranch {
		return errors.New("You cannot run deps ci on a deps branch")
	}

	cfg, err := config.FindOrInfer()
	if err != nil {
		return err
	}

	allUpdates, err := collectUpdates(cfg, types, paths)
	if err != nil {
		return err
	}

	output.Debug("%d collected updates", len(allUpdates))

	newUpdates, outdatedUpdates, existingUpdates, err := organizeUpdates(allUpdates)
	if err != nil {
		return err
	}

	// Find existing updates where merge_base enabled
	// and a merge is available, and move them to outdated updates
	for _, update := range existingUpdates {
		if update.mergeUpdatesEnabled() {
			git.Checkout(update.branch)
			mergeAvailable := git.MergeAvailable(startingBranch)
			git.Checkout("-")

			if mergeAvailable {
				outdatedUpdates.addUpdate(update)
				existingUpdates.removeUpdate(update)
			}
		}
	}

	output.Event("%d new updates", len(newUpdates))
	output.Event("%d outdated updates", len(outdatedUpdates))
	output.Event("%d existing updates", len(existingUpdates))

	if len(types) == 0 && len(paths) == 0 {
		inapplicableBranches := getInapplicableBranches(allUpdates)
		output.Event("%d inapplicable branches", len(inapplicableBranches))

		for _, branch := range inapplicableBranches {
			// On GitHub at least, deleting these also closes the PR (so also works for no-PR scenario)
			// at some point we could add a helpful comment but that would require
			// implementing all steps for GitHub/GitLab/Bitbucket
			if err := git.DeleteRemoteBranch(branch); err != nil {
				output.Error("Failed to delete inapplicable branch %s\n%s", branch, err)
			} else {
				output.Event("Deleted inapplicable branch %s", branch)
			}
		}
	} else {
		output.Debug("Deleting inapplicable branches not enabled when filtering types or paths")
	}

	// TODO this is also because collectors may have done some crap and not cleaned up
	if git.IsDirty() {
		output.Event("Temporarily saving your uncommitted changes in a git stash")
		stashed := git.Stash(fmt.Sprintf("Deps save before update"))

		// Stash pop needs to happen last (so be added first)
		defer func() {
			if stashed {
				output.Event("Putting original uncommitted changes back")
				if err := git.StashPop(); err != nil {
					output.Error("Error putting stash back: %v", err)
				}
			}
		}()
	}

	for _, update := range newUpdates {
		output.StartSection("(New) %s", update.title)
		if err := runUpdate(update, startingBranch, update.branch, false); err != nil {
			failedUpdates = append(failedUpdates, &updateResult{
				update: update,
				err:    err,
			})
			output.Error("Update failed: %v", err)
		} else {
			successfulUpdates = append(successfulUpdates, &updateResult{
				update: update,
				err:    err,
			})
			output.Success("Update succeeded: %v", update.title)
		}
		output.EndSection()
	}

	for _, update := range outdatedUpdates {
		output.StartSection("(Outdated) %s", update.title)
		// TODO if update.branch already exists, maybe base could be
		// determined from what it originally branched off of?
		if err := runUpdate(update, startingBranch, update.branch, true); err != nil {
			failedUpdates = append(failedUpdates, &updateResult{
				update: update,
				err:    err,
			})
			output.Error("Update failed: %v", err)
		} else {
			successfulUpdates = append(successfulUpdates, &updateResult{
				update: update,
				err:    err,
			})
			output.Success("Update succeeded: %v", update.title)
		}
		output.EndSection()
	}

	if len(successfulUpdates) > 0 {
		output.Success("%d updates made successfully!", len(successfulUpdates))
		for _, ue := range successfulUpdates {
			output.Success("- [%s] %s", ue.update.id, ue.update.title)
		}
	}

	if len(failedUpdates) > 0 {
		output.Error("There were %d errors making the updates", len(failedUpdates))
		for _, ue := range failedUpdates {
			output.Error("- [%s] %s\n  %v", ue.update.id, ue.update.title, ue.err)
		}
	}

	if len(successfulUpdates) > 0 {
		if err := api.IncrementUsage(len(successfulUpdates)); err != nil {
			return err
		}
	}

	if len(failedUpdates) > 0 {
		return fmt.Errorf("%d errors", len(failedUpdates))
	}

	return nil
}

func getCurrentBranch(ci ci.CIProvider) string {
	branch := git.CurrentRef()

	// CI environments may be checking out a specific ref,
	// so use the variables they provide to see if we get a different branch name

	if b := ci.Branch(); b != "" {
		branch = b
	}

	if branch == "" {
		panic(errors.New("Unable to determine base branch"))
	}

	if branch == "HEAD" {
		panic(errors.New("Unable to determine base branch, only got HEAD"))
	}

	return branch
}

func (update *Update) mergeUpdatesEnabled() bool {
	branchUpdatesSetting := update.dependencyConfig.GetSettingForSchema("branch_updates", update.dependencies)
	if branchUpdatesSetting == nil {
		return false
	}
	branchUpdates := branchUpdatesSetting.(string)
	branchUpdates = strings.ToLower(branchUpdates)
	// merge_if_failing would be a potential improvement here,
	// but CI dependent because we'd need to know the current status for a PR
	// rebase may be another option
	return branchUpdates == "merge"
}

func runUpdate(update *Update, base, head string, existingUpdate bool) error {
	branchUpdated := false

	if existingUpdate {
		// go straight to it
		git.Checkout(head)

		if update.mergeUpdatesEnabled() {
			if git.MergeWouldConflict(base) {
				// Fine to be quiet on this because conflicts show up in host UI
				output.Event("Merge with %s has a conflict, so skipping automatic merge", base)
			} else {
				output.Event("Merging %s into existing update", base)
				if git.Merge(base) {
					branchUpdated = true
				}
			}
		}
	} else {
		// create a branch for it
		git.Checkout(base)
		git.Branch(head)
	}

	defer func() {
		// There should only be uncommitted changes if we're bailing early
		git.ResetAndClean()
		git.CheckoutLast()
	}()

	outputDeps, err := update.runner.Act(update.dependencies)
	if err != nil {
		return err
	}

	pr, err := pullrequest.NewPullrequest(base, head, outputDeps, update.dependencyConfig)
	if err != nil {
		return err
	}

	if !git.IsDirty() && !branchUpdated {
		if existingUpdate {
			output.Event("No new changes to commit")
			return nil
		}

		return errors.New("Update didn't generate any changes to commit")
	}

	if git.IsDirty() {
		if err := hooks.RunPullrequestHook(pr, "before_commit"); err != nil {
			return err
		}

		templateString := "{{.SubjectAndBody}}"
		if templateSetting := pr.GetSetting("commit_message_template"); templateSetting != nil {
			templateString = templateSetting.(string)
		}
		commitMessage, err := renderCommitMessage(outputDeps, templateString)
		if err != nil {
			return err
		}

		git.Add()
		git.Commit(commitMessage)
	}

	git.PushBranch(head)

	if pr != nil {
		output.Debug("Waiting a second for the push to be processed by the host")
		time.Sleep(2 * time.Second)

		if err := pr.CreateOrUpdate(); err != nil {
			output.Error("Creating or updating pull request failed. Deleting branch.")
			if deleteErr := git.DeleteRemoteBranch(head); deleteErr != nil {
				return deleteErr
			}
			return err
		}
	}

	return nil
}

func renderCommitMessage(deps *schema.Dependencies, templateString string) (string, error) {
	tmpl, err := template.New("commitmessage").Parse(templateString)
	if err != nil {
		return "", err
	}

	subject := schemaext.TitleForDeps(deps)
	subjectAndBody := subject

	// Extra explanation not needed in single manifest scenario
	body := ""
	schemaBody := schemaext.DescriptionItemsForDeps(deps)
	if len(deps.Lockfiles) > 0 || len(strings.Split(schemaBody, "\n")) > 1 {
		body = schemaBody
		subjectAndBody = subjectAndBody + "\n\n" + body
	}

	vars := struct {
		Subject        string
		Body           string
		SubjectAndBody string
	}{
		Subject:        subject,
		Body:           body,
		SubjectAndBody: subjectAndBody,
	}

	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, vars); err != nil {
		return "", err
	}

	message := strings.TrimSpace(buf.String())

	if message == "" {
		return "", errors.New("commit message can not be empty")
	}

	return message, nil
}

func getInapplicableBranches(updates Updates) []string {
	inapplicableBranches := []string{}
	for _, branch := range git.GetDepsBranches() {
		branchFound := false
		for _, update := range updates {
			if strings.HasPrefix(branch, update.branchPrefix()) {
				branchFound = true
				break
			}
		}
		if !branchFound {
			inapplicableBranches = append(inapplicableBranches, branch)
		}
	}
	return inapplicableBranches
}
