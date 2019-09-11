---
title: CircleCI
description: Examples and instructions for setting up deps in CircleCI
---

# CircleCI

## 1. Set your environment variables

For deps variables that should be kept secret (such as `DEPS_GITHUB_TOKEN`) you'll want to use a [context](https://circleci.com/docs/2.0/contexts/).
You can get your `DEPS_TOKEN` from [3.dependencies.io](https://3.dependencies.io).

[![CircleCI context for deps](/assets/img/screenshots/circleci-context.png)](/assets/img/screenshots/circleci-context.png)

## 2. Add a deps workflow triggered by cron

This example shows two different languages in use,
each running in their own container.

The `--type` option is used to run the specific language updates in their respective containers.

> Note: CircleCI is supported by [CI autoconfigure](/ci/#autoconfigure).

```yaml
version: 2
jobs:

  # in addition to your existing jobs

  deps-python:
    docker:
      - image: circleci/python:3.7
    steps:
      - checkout
      - run: curl https://deps.app/install.sh | bash -s -- -b $HOME/bin
      - run: pipenv sync --dev
      - run: $HOME/bin/deps ci --type python
  deps-node:
    docker:
      - image: circleci/node:12
    steps:
      - checkout
      - run: curl https://deps.app/install.sh | bash -s -- -b $HOME/bin
      - run: yarn install
      - run: $HOME/bin/deps ci --type js

workflows:
  version: 2
  # in addition to your existing workflows
  deps:
    jobs:
      - deps-python:
          context: deps
      - deps-node:
          context: deps
    triggers:
      - schedule:
          cron: "0 0 * * *"  # nightly
          filters:
            branches:
              only:
                - master
```
