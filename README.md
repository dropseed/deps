# deps [![GitHub release](https://img.shields.io/github/release/dropseed/deps.svg)](https://github.com/dropseed/deps/releases)

Deps is a command line tool for managing dependencies that can run locally and in CI.

This repo contains the code for the `deps` command line tool itself (written in Go),
but each language/ecosystem has it's own repo which often uses the native language and acts as a light wrapper around the native dependency management tools.
This way deps can automate updates using the same tools that you would use in your terminal.

[Read the docs →](https://www.dependencies.io)

## Quick install

```console
$ curl https://deps.app/install.sh | bash -s -- -b $HOME/bin
```

## Official Components

- [dropseed/deps-python](https://github.com/dropseed/deps-python)
- [dropseed/deps-js](https://github.com/dropseed/deps-js)
- [dropseed/deps-php](https://github.com/dropseed/deps-php)
- [dropseed/deps-git](https://github.com/dropseed/deps-git)
- [dropseed/deps-manual](https://github.com/dropseed/deps-manual)
- [dropseed/deps-wordpress-core](https://github.com/dropseed/deps-wordpress-core)
- [dropseed/deps-wordpress-plugins](https://github.com/dropseed/deps-wordpress-plugins)
- [dropseed/deps-wordpress-themes](https://github.com/dropseed/deps-wordpress-themes)
