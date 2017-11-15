package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func (pr *PullRequest) getGitHubClient() (*github.Client, context.Context, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: pr.APIToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return client, ctx, nil
}

// func (pr *PullRequest) getGitHubRepo() (*github.Repository, error) {
// 	client, ctx, err := pr.getGitHubClient()
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	repo, _, err := client.Repositories.Get(ctx, pr.OwnerName, pr.Name)
// 	return repo, err
// }

func (pr *PullRequest) getRelatedPR() (*github.Issue, error) {
	client, ctx, err := pr.getGitHubClient()
	if err != nil {
		return nil, err
	}

	relatedPRTitleSearch, err := pr.Pullrequest.Config.RelatedPRTitleSearchFromConfig()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(
		"\"%v\" in:title author:%v is:pr is:open repo:%v created:<%v",
		strings.Replace(relatedPRTitleSearch, "\"", "\\\"", -1),
		pr.Config.EnvSettings.RelatedPRAuthor,
		pr.RepoFullName,
		pr.CreatedAt,
	)
	fmt.Printf("Searching for: %v", query)
	issuesResult, _, err := client.Search.Issues(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	if total := issuesResult.GetTotal(); total > 1 {
		return nil, fmt.Errorf("%v issues were found so we quit just to be safe", total)
	} else if total < 1 {
		return nil, nil
	}

	return &issuesResult.Issues[0], nil
}

func (pr *PullRequest) closePR(number int) error {
	ghc, ctx, err := pr.getGitHubClient()
	if err != nil {
		return err
	}
	newState := "closed"
	ir := github.IssueRequest{State: &newState}
	fmt.Printf("Closing related PR #%v\n", number)
	_, _, err = ghc.Issues.Edit(ctx, pr.RepoOwnerName, pr.RepoName, number, &ir)
	if err != nil {
		return err
	}

	fmt.Printf("Deleting branch from PR #%v\n", number)
	ghPR, _, err := ghc.PullRequests.Get(ctx, pr.RepoOwnerName, pr.RepoName, number)
	if err != nil {
		return err
	}
	_, err = ghc.Git.DeleteRef(ctx, pr.RepoOwnerName, pr.RepoName, "heads/"+ghPR.Head.GetRef())
	return err
}

func (pr *PullRequest) commentOnIssue(number int, comment string) error {
	ghc, ctx, err := pr.getGitHubClient()
	if err != nil {
		return err
	}

	ic := github.IssueComment{Body: &comment}
	fmt.Printf("Commenting on PR #%v\n", number)
	_, _, err = ghc.Issues.CreateComment(ctx, pr.RepoOwnerName, pr.RepoName, number, &ic)
	return err
}
