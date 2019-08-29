# Using deps locally

You can use deps locally for free and without setting up an API token!

Running it on your machine is the best way to...

- Test advanced configuration without committing to your repo
- Make one-off updates
- See how deps works *before* setting up CI

## Installing

To install `deps`,
[manually download a release](https://github.com/dropseed/deps/releases) or use the following auto-install script:
```sh
$ curl https://www.dependencies.io/install.sh | bash -s -- -b $HOME/bin
```

## Checking for dependency updates

With `deps` in your `$PATH`, you can check a repo for dependency updates.
You'll get an interactive prompt where you can choose and install updates one-by-one:

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

Using `deps run` will perform the same steps as `deps ci`,
but it will *not* commit changes or create pull requests.

If your dependencies were not found automatically,
or you need a more advanced configuration,
[take a look at `deps.yml`](/config/).
