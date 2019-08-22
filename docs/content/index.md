---
description: "Documentation and examples for automatically updating your dependencies with deps"
---

# Overview

Deps is a command line tool written in Go.
You can install and run it on your own machine to make one-off updates and see how things work.
Running it locally is also the easiest way to test more advanced configurations *before* committing and pushing it to CI.

**To automate dependency updates,
you'll set up deps to run directly in your CI.**
Each CI provider is slightly different,
but generally all you need to do beyond your existing setup is install the tool and set some environment variables!
When you run `deps ci`,
it will automatically commit dependency updates to new branches and send pull requests for you to review.

## Using deps on your machine

Running deps locally is completely free and doesn't require and API token.

To install it,
[manually download a release](https://github.com/dropseed/deps/releases) or use the following auto-install script:
```sh
$ curl https://www.dependencies.io/install.sh | bash -s -- -b $HOME/bin
```

With `deps` in your `$PATH`, you can check a repo for dependency updates:
```sh
$ cd project
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

If your dependencies were not automatically found,
or you need a more advanced configuration,
take a look at `deps.yml`.

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
