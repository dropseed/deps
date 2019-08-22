---
title: "TravisCI"
description: "Examples and instructions for setting up deps in TravisCI"
---

# TravisCI

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
