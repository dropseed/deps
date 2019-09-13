package runner

import (
	"errors"
	"fmt"

	"github.com/dropseed/deps/internal/config"
	"github.com/dropseed/deps/internal/git"
	"github.com/dropseed/deps/internal/output"

	"github.com/manifoldco/promptui"
)

// Local runs a full interactive update process
func Local() error {
	if git.HasStagedChanges() {
		return errors.New("You can't have staged changes while running this command. Please commit or unstage them.")
	}

	if git.IsDirty() {
		output.Warning("You have uncommitted changes! We are going to stage them so that we can tell the difference between changes. They will be unstaged when this command exits.\n")
		git.Add()
	}

	cfg, err := config.FindOrInfer()
	if err != nil {
		return err
	}

	allUpdates, err := collectUpdates(cfg, []string{})
	if err != nil {
		return err
	}

	if git.HasStagedChanges() || git.IsDirty() {
		output.Debug("Restoring the state of your repo before updates were collected")
		git.Checkout(".")
		git.Unstage()
	}

	newUpdates, _, _, err := organizeUpdates(allUpdates)
	if err != nil {
		return err
	}

	if err := newUpdates.prompt(); err != nil {
		return err
	}

	return nil
}

func (updates Updates) prompt() error {
	for {
		refs := map[int]string{}
		items := []string{}

		updateIndex := 0
		for _, update := range updates {
			if !update.completed {
				items = append(items, update.title)
				refs[updateIndex] = update.id
				updateIndex++
			}
		}

		if len(items) < 1 {
			// No updates left
			break
		}

		items = append(items, "Skip")

		prompt := promptui.Select{
			Label: fmt.Sprintf("Choose an update to make"),
			Items: items,
		}

		println()
		i, _, err := prompt.Run()
		if err != nil {
			return err
		}

		if i+1 == len(items) {
			// Chose to skip
			break
		}

		update := updates[refs[i]]
		if _, err := update.runner.Act(update.dependencies); err != nil {
			return err
		}
		update.completed = true
	}

	return nil
}
