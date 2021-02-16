# Using deps locally

Running `deps` on your own machine is the best way to...

- See how deps works *before* setting up CI
- Make one-off updates
- Test advanced configuration without committing to your repo

## Installing

To install `deps`,
[manually download a release](https://github.com/dropseed/deps/releases) or use the auto-install script:
```console
$ curl https://deps.app/install.sh | bash -s -- -b $HOME/bin
```

## Checking for dependency updates

With `deps` in your `$PATH`, you can check a repo for dependency updates.
You'll get an interactive prompt where you can choose and install updates one-by-one:

