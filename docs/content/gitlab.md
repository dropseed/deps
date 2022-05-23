---
title: "GitLab"
description: "Examples and instructions for setting up deps in GitLab"
---

# GitLab

Deps can run directly in GitLab CI using a [pipeline schedule](#pipeline-schedule).

If your CI is separated into language-specific containers,
you can use the `--type` option to tell deps which updates to run.

A minimal example of `.gitlab-ci.yml` would look like this:

```yaml
# .gitlab-ci.yml
deps-js:
  image: "node:latest"
  only: [schedules]
  script:
    - curl -sSL https://deps.app/install.sh | bash -s -- -b $HOME/bin
    - yarn install
    - $HOME/bin/deps ci --type js

deps-python:
  image: "python:3"
  only: [schedules]
  script:
    - curl -sSL https://deps.app/install.sh | bash -s -- -b $HOME/bin
    - pipenv sync --dev
    - $HOME/bin/deps ci --type python
```

## Pipeline schedule

Configure a pipeline schedule for running deps.
It will need two variables:

1. A `DEPS_TOKEN` from [3.dependencies.io](https://3.dependencies.io).
1. A `DEPS_GITLAB_TOKEN` that is a [personal access token](#personal-access-token).

![GitLab pipeline schedule for deps](/assets/img/screenshots/gitlab-ci-pipeline-schedule.png)


## Personal access token

The standard way to give deps write-access to your repo is with a *personal access token*.
You can use your personal account to do this, or a "bot" account that your team has.

1. Log in with the account you want deps to use (this will be the author of deps pull requests)
1. Give it access to the repo you're setting up
1. [Generate a new token](https://gitlab.com/profile/personal_access_tokens) with the `write_repository` and `api` scopes
    [![GitLab personal access token settings for deps](/assets/img/screenshots/gitlab-personal-access-token.png)](/assets/img/screenshots/gitlab-personal-access-token.png)
1. Set the `DEPS_GITLAB_TOKEN` environment variable in your CI


## Merge request settings

When working with a GitLab repo,
there are a few settings you can use to determine what your MRs look like.

```yaml
# deps.yml
version: 3
dependencies:
- type: python
  settings:
    gitlab_target_branch: "dev"
    gitlab_labels: ["dependencies"]
    gitlab_assignee_id: 1
    gitlab_assignee_ids: [1, 2]
    gitlab_target_project_id: 1
    gitlab_milestone_id: 1
    gitlab_remove_source_branch: true
    gitlab_allow_collaboration: true
    gitlab_allow_maintainer_to_push: true
    gitlab_squash: true
```

If you don't need a `deps.yml` then you can also configure these settings via environment variables.
This is an easy way to put settings directly in your CI config.

Note that they'll need to be in the format of a JSON-encoded string,
with an uppercase name prefixed by `DEPS_SETTING_`.

```console
$ DEPS_SETTING_GITLAB_LABELS='["dependencies"]' deps ci
```
