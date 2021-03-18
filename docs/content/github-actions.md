---
title: GitHub Actions
description: Examples and instructions for setting up deps in GitHub Actions
---

# GitHub Actions

## 1. Add `DEPS_TOKEN` to your secrets

You can get your `DEPS_TOKEN` from [3.dependencies.io](https://3.dependencies.io).

[![GitHub actions secrets for deps](/assets/img/screenshots/github-actions-secrets.png)](/assets/img/screenshots/github-actions-secrets.png)

## 2. Add a deps workflow

Create a [workflow file](https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions) at `.github/workflows/deps.yml`

You will need the `DEPS_TOKEN` and `DEPS_GITHUB_TOKEN` at a minimum.

Here's an example for a Python project:

```yaml
name: deps

on:
  schedule:
  - cron: 0 0 * * *

jobs:
  deps:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v1
    - uses: actions/setup-python@v1
    - run: pip install pipenv
    - run: ./scripts/install
    - run: curl https://deps.app/install.sh | bash -s -- -b $HOME/bin
    - run: $HOME/bin/deps ci
      env:
        DEPS_TOKEN: ${{ secrets.DEPS_TOKEN }}
        DEPS_GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

[You can also authenticate as a GitHub App, and customize your pull requests with labels, assignees, and more →](/github/)
