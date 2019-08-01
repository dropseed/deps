# How to use deps

## Locally

An easy way to get started with `deps` is to try it on your own machine!

To install it,
[download a release](https://github.com/dropseed/deps/releases) or use the following install script:
```sh
$ curl https://www.dependencies.io/install.sh | bash -s -- -b $HOME/bin
```

Once `deps` is in your `$PATH`, you can check a local project for dependency updates:
```sh
$ cd MY_GIT_REPO
$ deps run
```

## In CI

When you are comfortable with how deps works,
you'll want to set it up in your CI provider to automate your dependency updates.
Most CI systems have a cron or scheduled job function,
which is how we recommend using it.
This way your regular CI pipeline can run uninterruped by deps during your normal workflows.

Below are some simple examples of how you might configure this,
but you should adjust the setup to work how you want within your existing CI configuration.

You will also need to give deps access to your git repo. The easiest way to do this is usually with personal access tokens.

### TravisCI

1. Set your environment variables
1. Add a daily/weekly/monthly [cron job in the Travis settings UI](https://docs.travis-ci.com/user/cron-jobs/)
1. Define `jobs` in your `.travis.yml` and use `if` statements to define a deps job that only runs via cron

```yaml
language: python

jobs:
  include:
  - name: test
    script: echo "Tests passing here..."
    install: ./scripts/install_requirements
  - name: deps
    if: branch = master AND type = cron
    git:
      depth: false  # required for existing deps branches to be available
    install:
    - ./scripts/install_requirements
    - curl https://www.dependencies.io/install.sh | bash -s -- -b $HOME/bin
    script: deps ci
```

### CircleCI

1. Set your environment variables in a [context](https://circleci.com/docs/2.0/contexts/)
1. Add a deps workflow and accompanying jobs to your `.circleci/config.yml`

This example shows two different languages in use,
each running in their own container.
The `--type` option is used to run the specific language updates in their respective containers.

```yaml
version: 2
jobs:

  # in addition to your existing jobs

  deps-python:
    docker:
      - image: circleci/python:3.7
    steps:
      - checkout
      - run: curl https://www.dependencies.io/install.sh | bash -s -- -b $HOME/bin
      - run: $HOME/bin/deps ci --type python
  deps-node:
    docker:
      - image: circleci/node:12
    steps:
      - checkout
      - run: curl https://www.dependencies.io/install.sh | bash -s -- -b $HOME/bin
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

### BitBucket Pipelines

TBD

### Jenkins

TBD

### GitHub Actions

TBD
