---
description: "Documentation and examples of how to use Dependencies.io"
---

# Overview

Deps is a command line tool written in Go.
You can install it on your own machine to manually make dependency updates and see how things work or tweak your configuration.

To automate dependency updates,
you'll set up deps to run directly in your CI.

## Using deps on your machine

An easy way to get started with `deps` is to try it on your own machine!
You don't need an API token and it's completely free to run locally.

To install it,
[download a release](https://github.com/dropseed/deps/releases) or use the following install script:
```sh
$ curl https://www.dependencies.io/install.sh | bash -s -- -b $HOME/bin
```

Once `deps` is in your `$PATH`, you can check a local project for dependency updates:
```sh
$ cd MY_GIT_REPO
$ deps run

> No local config found, detecting your dependencies automatically
---
version: 3
dependencies:
- type: python
  path: Pipfile
- type: js
  path: .

---
> Collecting with deps-python
> Collecting with deps-js

> 2 new updates to be made
> [a711128] Update package-lock.json
> [2231fa2] Update tailwindcss in package.json from 1.0.5 to 1.1.2

Use the arrow keys to navigate: ↓ ↑ → ←
? Choose an update to make:
  ▸ Update package-lock.json
    Update tailwindcss in package.json from 1.0.5 to 1.1.2
    Skip
```

Running deps locally will run the same steps as `deps ci`,
but it will *not* commit changes or create pull requests.

## In CI

Most CI systems have a cron or scheduled job function,
which is how we recommend using it.
This way you don't have to run deps on every commit and you can also decide how often you want updates (daily, weekly, monthly, etc.).

You can run deps in almost any combination of git host (GitHub, GitLab, BitBucket, etc.) and CI provider (CircleCI, TravisCI, etc.).

You will need an API token to run `deps ci`.
You can get one at [3.dependencies.io](https://3.dependencies.io/).

### CI providers

- [CircleCI](/circleci/)

### Git hosts

- [GitHub](/github/)

## deps.yml

Deps can automatically detect most dependency types.
For more advanced configurations,
you can put a `deps.yml` file at the root of your repo.
