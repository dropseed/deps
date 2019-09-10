---
title: GitLab CI
description: Examples and instructions for setting up deps in GitLab CI
---

# GitLab CI

## 1. Create a pipeline schedule

You will need a `DEPS_TOKEN` from [3.dependencies.io](https://3.dependencies.io) and
a `DEPS_GITLAB_TOKEN` if you are running this on a [GitLab repo](/gitlab/).

[![GitLab pipeline schedule for deps](/assets/img/screenshots/gitlab-ci-pipeline-schedule.png)](/assets/img/screenshots/gitlab-ci-pipeline-schedule.png)

## 2. Add deps jobs to .gitlab-ci.yml

This example shows two different languages in use,
each running in their own container.

The `--type` option is used to run the specific language updates in their respective containers.

> Note: GitLab CI is supported by [CI autoconfigure](/ci/#autoconfigure).

A minimal example of using deps in `.gitlab-ci.yml` would look like this:

```yaml
deps-js:
  image: "node:latest"
  only: [schedules]
  script:
    - curl https://deps.app/install.sh | bash -s -- -b $HOME/bin
    - yarn install
    - $HOME/bin/deps ci --type js

deps-python:
  image: "python:3"
  only: [schedules]
  script:
    - curl https://deps.app/install.sh | bash -s -- -b $HOME/bin
    - pipenv sync --dev
    - $HOME/bin/deps ci --type python
```
