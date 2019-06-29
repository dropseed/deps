package github

import (
	"context"

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

// func (pr *PullRequest) commentOnIssue(number int, comment string) error {
// 	ghc, ctx, err := pr.getGitHubClient()
// 	if err != nil {
// 		return err
// 	}

// 	ic := github.IssueComment{Body: &comment}
// 	fmt.Printf("Commenting on PR #%v\n", number)
// 	_, _, err = ghc.Issues.CreateComment(ctx, pr.RepoOwnerName, pr.RepoName, number, &ic)
// 	return err
// }
