---
title: Bitbucket Pipeilnes
description: Examples and instructions for setting up deps in Bitbucket Pipelines
---

# Bitbucket Pipelines

## 1. Add deps to bitbucket-pipelines.yml

Add a custom pipeline so that it will only run from a schedule.

This example shows two different languages in use,
each running in their own container. The `--type` option is used to run the specific language updates in their respective containers.

> Note: Bitbucket Pipelines are supported by [CI autoconfigure](/ci/#autoconfigure).

A minimal example of using deps in `bitbucket-pipelines.yml` would look like this:

```yaml
clone:
  depth: full

pipelines:
  custom:
    deps:
      - parallel:
        - step:
            image: "python:3.7"
            script:
            - curl https://www.dependencies.io/install.sh | bash -s -- -b $HOME/bin
            - $HOME/bin/deps ci --type python
        - step:
            image: "node:latest"
            script:
            - curl https://www.dependencies.io/install.sh | bash -s -- -b $HOME/bin
            - $HOME/bin/deps ci --type js
```

## 2. Set the pipeline repository variables

For a standard Bitbucket repo, you will need a `DEPS_TOKEN`, `DEPS_BITBUCKET_USERNAME`, and `DEPS_BITBUCKET_PASSWORD`.

[![Bitbucket pipeline variables for deps](/assets/img/screenshots/bitbucket-pipeline-variables.png)](/assets/img/screenshots/bitbucket-pipeline-variables.png)

## 3. Create a pipeline schedule

Create a daily or weekly schedule to run your new deps pipeline.

[![Bitbucket pipeline schedule for deps](/assets/img/screenshots/bitbucket-pipeline-schedule.png)](/assets/img/screenshots/bitbucket-pipeline-schedule.png)

## 4. Test the pipeline manually

If you want to test your new pipeline without waiting for the schedule,
just navigate to the commit view and click "run pipeline".

[![Bitbucket pipeline manual run for deps](/assets/img/screenshots/bitbucket-pipeline-manual.png)](/assets/img/screenshots/bitbucket-pipeline-manual.png)
