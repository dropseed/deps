# Quickstart

You can use `deps` locally, set it up in CI, or both.

Locally `deps` can:

- Install dependencies with one command
- Interactively upgrade dependencies (including out-of-range updates)
- Tell you when to re-install (new yarn.lock after `git pull`, etc.)

Setting up `deps` in CI will automatically create pull requests for new updates!

## Using deps locally

To install `deps`,
[manually download a release](https://github.com/dropseed/deps/releases) or use the auto-install script:

```console
$ curl https://deps.app/install.sh | bash -s -- -b $HOME/bin
```

Now if you `cd` to one of your repos,
you can see what is out of date and interactively upgrade dependencies with `deps upgrade`:

```console
$ cd project
$ deps upgrade

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

Using `deps upgrade` will perform the same steps as `deps ci`,
but it will *not* commit changes or create pull requests.

If your dependencies were not found automatically,
or you need a more advanced configuration,
[take a look at `deps.yml`](/config/).


optional, set up shell hook

## Automating deps in CI
