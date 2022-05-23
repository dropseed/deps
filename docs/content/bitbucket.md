---
title: "Bitbucket"
description: "Examples and instructions for setting up deps in Bitbucket"
---

# Bitbucket

Run deps in Bitbucket via a custom pipeline that only runs from a schedule.

If languages are split across containers,
use the `--type` option to tell deps which updates to run.

```yaml
# bitbucket-pipelines.yml
clone:
  depth: full

pipelines:
  custom:
    deps:
      - parallel:
        - step:
            image: "python:3.7"
            script:
            - curl -sSL https://deps.app/install.sh | bash -s -- -b $HOME/bin
            - python3 -m venv .venv
            - .venv/bin/pip install -r requirements.txt
            - $HOME/bin/deps ci --type python

        - step:
            image: "node:latest"
            script:
            - curl -sSL https://deps.app/install.sh | bash -s -- -b $HOME/bin
            - yarn install
            - $HOME/bin/deps ci --type js
```

## Pipeline repository variables

Deps will need API access through an app password.
You can use your personal account to do this, or a "bot" account that your team has.

1. Log in with the account you want deps to use (this will be the author of deps pull requests)
1. Give it access to the repo you're setting up
1. Generate a new app password with the repositories and pull requests `write` scopes
    [![Bitbucket app password settings for deps](/assets/img/screenshots/bitbucket-app-password.png)](/assets/img/screenshots/bitbucket-app-password.png)
1. Set the required environment variables in your CI
    - `DEPS_BITBUCKET_USERNAME` to the user who owns the app password
    - `DEPS_BITBUCKET_PASSWORD` to the app password from above
    - `DEPS_TOKEN` to the token from [3.dependencies.io](https://3.dependencies.io/)

[![Bitbucket pipeline variables for deps](/assets/img/screenshots/bitbucket-pipeline-variables.png)](/assets/img/screenshots/bitbucket-pipeline-variables.png)

## Pipeline schedule

Create a daily or weekly schedule to run your new deps pipeline.

[![Bitbucket pipeline schedule for deps](/assets/img/screenshots/bitbucket-pipeline-schedule.png)](/assets/img/screenshots/bitbucket-pipeline-schedule.png)

## Test or run pipeline manually

If you want to test your new pipeline without waiting for the schedule,
just navigate to the branches view and click "run pipeline".

[![Bitbucket pipeline manual run for deps](/assets/img/screenshots/bitbucket-pipeline-manual.png)](/assets/img/screenshots/bitbucket-pipeline-manual.png)

## Pull request settings

When working with a Bitbucket repo,
there are a few settings you can use to determine what your pull requests look like.

```yaml
# deps.yml
version: 3
dependencies:
- type: python
  settings:
    bitbucket_destination: "dev"  # branch name
    bitbucket_close_source_branch: true
    bitbucket_reviewers:
    - uuid: "{638373c3b62-8120-4f0c-a7bc-87800b9d6f70}"
```

If you don't need a `deps.yml` then you can also configure these settings via environment variables.
This is an easy way to put settings directly in your CI config.

Note that they'll need to be in the format of a JSON-encoded string,
with an uppercase name prefixed by `DEPS_SETTING_`.

```console
$ DEPS_SETTING_BITBUCKET_CLOSE_SOURCE_BRANCH='true' deps ci
```
