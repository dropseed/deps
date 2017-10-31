# pullrequest [![GitHub release](https://img.shields.io/github/release/dependencies-io/pullrequest.svg)](https://github.com/dependencies-io/pullrequest/releases) [![Build Status](https://travis-ci.org/dependencies-io/pullrequest.svg?branch=master)](https://travis-ci.org/dependencies-io/pullrequest) [![license](https://img.shields.io/github/license/dependencies-io/pullrequest.svg)](https://github.com/dependencies-io/pullrequest/blob/master/LICENSE)

A simple Go application that can send pull requests to the different Git hosts.
Can be easily installed inside Docker containers of any kind, so that they don't
have to implement this functionality themselves.

This was built for being run within the [dependencies.io](https://www.dependencies.io) platform, though it could probably be used however. The key though is the environment variables, which match up with what [dependencies.io](https://www.dependencies.io) uses.

## Usage in a dependencies.io actor

Add to your Dockerfile.

```sh
# add the pullrequest utility to easily create pull requests on different git hosts
WORKDIR /usr/src/actor
ENV PULLREQUEST_VERSION=0.2.1
RUN wget https://github.com/dependencies-io/pullrequest/releases/download/${PULLREQUEST_VERSION}/pullrequest_${PULLREQUEST_VERSION}_linux_amd64.tar.gz && \
    mkdir pullrequest && \
    tar -zxvf pullrequest_${PULLREQUEST_VERSION}_linux_amd64.tar.gz -C pullrequest && \
    ln -s /usr/src/actor/pullrequest/pullrequest /usr/local/bin/pullrequest
```

Do your `git commit` and `git push`, so that the branch exists on the repo.

Then use `pullrequest`!

```sh
pullrequest --branch="branch-name" --title="PR title" --body="PR body"
```

### dependencies.yml

Any actor using this will have these settings available, so they should be added to the README.

```yaml
settings:
  pullrequest_notes: Notes that will be inserted at the top of the PR body.

  # github options
  github_labels:  # list of label names
  - bug
  github_assignees:  # list of usernames
  - davegaeddert
  github_milestone: 3  # milestone number
  github_base_branch: develop  # branch to make PR against (if something other than your default branch)

  # gitlab options
  gitlab_assignee_id: 1  # assignee user ID
  gitlab_labels:  # labels for MR as a list of strings
  - dependencies
  - update
  gitlab_milestone_id: 1  # the ID of a milestone
  gitlab_target_project_id: 1  # The target project (numeric id)
  gitlab_remove_source_branch: true  # flag indicating if a merge request should remove the source branch when merging
  gitlab_target_branch: develop  # branch to make PR against (if something other than your default branch)
```

## Environment variables

### Always required

- `GIT_HOST` - "github" or "gitlab"
- `GIT_BRANCH` - the default branch on the repo (usually "master")
- `DEPENDENCIES_ENV` - should be "production" for PRs to actually be created

### Always optional

- `SETTING_PULLREQUEST_NOTES` - user-supplied content to insert at the top of the PR body

### GitHub

##### Required

- `GITHUB_REPO_FULL_NAME`
- `GITHUB_API_TOKEN`

##### Optional

- `SETTING_GITHUB_LABELS` - JSON encoded list of strings
- `SETTING_GITHUB_ASSIGNEES` - JSON encoded list of strings
- `SETTING_GITHUB_MILESTONE` - int of milestone number

### GitLab

#### Required

- `GITLAB_API_URL` - url to the project in the GitLab API
- `GITLAB_API_TOKEN`

#### Optional

- `SETTING_GITLAB_ASSIGNEE_ID` - int
- `SETTING_GITLAB_LABELS` - JSON encoded list of strings
- `SETTING_GITLAB_MILESTONE_ID` - work in progress
- `SETTING_GITLAB_TARGET_PROJECT_ID` - int
- `SETTING_GITLAB_REMOVE_SOURCE_BRANCH` - JSON encoded bool

## Development

```sh
go build && env $(cat .env | xargs) ./pullrequest --branch="redux-3.7.2-11.1.0" --title=test --body="Testing it out"
```

---

[![dependencies.io](https://www.dependencies.io/permanent/github-readme-logotype.png)](https://www.dependencies.io)
