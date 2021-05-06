---
title: Using deps locally
description: Interactively update dependencies with the deps command line tool.
---

# Using deps locally

Locally `deps` can:

- Install dependencies with one command
- Interactively upgrade dependencies (including out-of-range updates)
- Tell you when to re-install (new yarn.lock after `git pull`, etc.)

[Setting up `deps` in CI](/ci/) will automatically create pull requests for new updates.

## Installing

To install `deps`,
[manually download a release](https://github.com/dropseed/deps/releases) or use the auto-install script:
```console
$ curl https://deps.app/install.sh | bash -s -- -b $HOME/bin
```

## Checking for dependency updates

With `deps` in your `$PATH`, you can check a repo for dependency updates.
You'll get an interactive prompt where you can choose and install updates one-by-one:

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

Using `deps upgrade` will perform the same steps as `deps ci`,
but it will *not* commit changes or create pull requests.

If your dependencies were not found automatically,
or you need a more advanced configuration,
[take a look at `deps.yml`](/config/).

## Shell hook

<div class="mb-6 aspect-w-16 aspect-h-9">
  <iframe src="https://www.youtube.com/embed/bNChNdpBroQ" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
</div>

The optional deps shell hook will help ensure your local installations are actually accurate.
After you change git branches or run git pull, it's easy to miss dependency changes that you actually need to install.
The deps shell hook will run (quickly!) before every bash/zsh prompt and let you know if you forget to install dependency updates.

```bash
# For ZSH, add this to the end of .zshrc
eval "$(deps shellhook zsh)"

# For BASH, add this to the end of .bashrc (or .bash_profile)
eval "$(deps shellhook bash)"
```

Now when you switch branches or pull dependency commits from team members and bots,
deps will automatically remind you to re-install dependencies so that your installation matches the lockfile.

```console
$ cd project
$ git pull
Updating 436b311..49c2cc3
Fast-forward
 package-lock.json  |  47 ++---

[Run `deps install` to sync with poetry.lock]

$ deps install
```
