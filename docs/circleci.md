# CircleCI

1. Set your environment variables in a [context](https://circleci.com/docs/2.0/contexts/)
1. Add a deps workflow and accompanying jobs to your `.circleci/config.yml`

This example shows two different languages in use,
each running in their own container.

The `--type` option is used to run the specific language updates in their respective containers.

```yml
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
