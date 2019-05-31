# deps [![GitHub release](https://img.shields.io/github/release/dependencies-io/deps.svg)](https://github.com/dependencies-io/deps/releases) [![Build Status](https://travis-ci.org/dependencies-io/deps.svg?branch=master)](https://travis-ci.org/dependencies-io/deps)

A Go application that can be installed both inside dependencies.io components, and on
your local machine.

## Installing on your machine

```sh
$ curl https://www.dependencies.io/install.sh | bash -s 2.5.0-beta.1
```

Or, you can download binaries manually from the [releases on
GitHub](https://github.com/dependencies-io/deps/releases).

## Usage in a dependencies.io component

Add to your Dockerfile.

```dockerfile
ENV DEPS_VERSION=2.5.0-beta
RUN curl https://www.dependencies.io/install.sh | bash -s -- -b /usr/local/bin $DEPS_VERSION
```

### Available commands

#### For collectors

`deps component collect <JSON file path>` - will report the contents of the JSON file back to dependencies.io

`deps component hook <hook name>` - run a named hook

#### For actors

`deps component branch` - creates and checks out a new branch for this update, using the
`JOB_ID` as a unique identifier

`deps component commit -m "message" <paths>` - adds and commits files to git, similar to
using `git` manually but automatically runs commit-related hooks (see below)

`deps component pullrequest <JSON file path>` - creates a pull request (or merge request)
on the host for the repo, using the contents of the JSON file for generating PR
content

`deps component hook <hook name>` - run a named hook

### Hooks

Hooks provide dependencies.io users a way of injecting their own commands and
scripts into the update process, making updates more flexible.

Hooks can be set by the user in the `settings` section of their config (see below).

#### For collectors

`before_update` - each collector should implement this (via `deps hook before_update`) before collection happens, so that any bootstrapping steps are run before dependencies are collected similarly to how they would be run before a pull request is made

#### For actors

`before_branch` - runs in `deps branch`, before the branch is actually created

`after_branch` - runs in `deps branch`, after the branch has been created

`before_update` - runs in `deps branch`, since that is usually the first step in
creating an update for a dependency

`before_commit` - runs in `deps commit`, before the commit is made

`after_commit` - runs in `deps commit`, after the commit is made

`after_update` - runs in `deps pullrequest`, before the pull request is created
since that is usually the last operation in an update

`before_pullrequest` - runs in `deps pullrequest`, before the pull request is created

`after_pullrequest` - runs in `deps pullrequest`, after the pull request has been created

### dependencies.yml

Any [dependencies-io](https://www.dependencies.io) component using this will have these settings available, so they should be added to the README.

```yaml
settings:
  branch_prefix: our-prefix/

  commit_message_prefix: "(chore) "

  pullrequest_notes: Notes that will be inserted at the top of the PR body.

  # automatically close outdated open PRs (works with GitHub only)
  related_pr_behavior: close

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

  # hooks - always a list of commands
  before_update:
  - ./scripts/bootstrap

  before_commit:
  - ./scripts/generate_file
  - git add examplefile
```

## Environment variables

### Always required

- `GIT_HOST` - "github" or "gitlab" ("test" for testing/development)
- `GIT_BRANCH` - the default branch on the repo (usually "master")
- `JOB_ID` - a unique identifier for this "job", in production it is a UUID4
- `DEPENDENCIES_ENV` - should be "production" for PRs to actually be created

### Always optional

- `SETTING_PULLREQUEST_NOTES` - user-supplied content to insert at the top of the PR body
- `SETTING_RELATED_PR_BEHAVIOR` - can be "close" to automatically close outdated open PRs
- `SETTING_{{HOOK_NAME}}` - JSON encoded list of commands

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

---

<p align="center">
    <a href="https://www.dependencies.io">
        <img src="https://www.dependencies.io/permanent/github-readme-logotype.png" />
    </a>
</p>
