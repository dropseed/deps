---
title: Quickstart
description: See how the deps command line tools works, and automate your first pull request.
---

# Quickstart

Follow these steps to see how deps works and get your first automated pull request.

## 1. Install deps locally

[Download a release manually](https://github.com/dropseed/deps/releases), or use the auto-install script:

```console
$ curl https://deps.app/install.sh | bash -s -- -b $HOME/bin
```

## 2. Try deps upgrade

Open a repo and run `deps upgrade` to see what is outdated and interactively upgrade:

```console
$ cd project
$ deps upgrade

No local config found, detecting your dependencies automatically
---
version: 3
dependencies:
- type: python
  path: Pipfile
- type: js
  path: .

---
Collecting with deps-python...
Collecting with deps-js...

2 new updates to be made
[a711128] Update package-lock.json
[2231fa2] Update tailwindcss in package.json from 1.0.5 to 1.1.2

Use the arrow keys to navigate: ↓ ↑ → ←
? Choose an update to make:
  ▸ Update package-lock.json
    Update tailwindcss in package.json from 1.0.5 to 1.1.2
    Skip
```

## 3. Install deps in CI

To automate pull requests, use the `deps ci` command in a cron or scheduled job on your CI provider.
How often the job runs will determine how often pull requests are created and updated.

You will need a `DEPS_TOKEN` in order to run `deps ci` ([see pricing](/pricing/)).

<!-- If your repo is on GitHub, `deps init` will automatically help set up a GitHub Actions workflow! -->

Supported CI providers:

- [GitHub Actions](/github-actions/)
- [GitLab CI](/gitlab-ci/)
- [Bitbucket Pipelines](/bitbucket-pipelines/)
- [CircleCI](/circleci/)
- [TravisCI](/travisci/)
- [Other](/other-ci/)
