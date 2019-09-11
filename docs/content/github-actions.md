---
title: GitHub Actions
description: Examples and instructions for setting up deps in GitHub Actions
---

# GitHub Actions

## 1. Add `DEPS_TOKEN` to your secrets

You can get your `DEPS_TOKEN` from [3.dependencies.io](https://3.dependencies.io).

[![GitHub actions secrets for deps](/assets/img/screenshots/github-actions-secrets.png)](/assets/img/screenshots/github-actions-secrets.png)

## 2. Add a deps workflow

A example of using deps with Python dependencies.
You can put your deps workflow in its own file,
such as `.github/workflows/deps.yml`.
You will need the `DEPS_TOKEN` and `GITHUB_TOKEN` at a minimum.

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
