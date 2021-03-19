<p align="center">
  <a href="https://dependencies.io/?utm_source=github&utm_medium=logo" target="_blank">
    <img src="https://user-images.githubusercontent.com/649496/111808669-375f3d00-88a2-11eb-9a25-2ee66a469b66.png" alt="Deps" height="140">
  </a>
</p>

# deps [![GitHub release](https://img.shields.io/github/release/dropseed/deps.svg)](https://github.com/dropseed/deps/releases)

Deps is a command line tool for managing dependencies that can run locally and in CI.

This repo contains the code for the `deps` command line tool itself (written in Go),
but each language/ecosystem has it's own repo which often uses the native language and acts as a light wrapper around the native dependency management tools.
This way deps can automate updates using the same tools that you would use in your terminal.

[Read the docs →](https://www.dependencies.io)

<!-- Edit in Excalidraw: https://excalidraw.com/#json=6195990008692736,unm4UYeUzdmbXTk_xcvIdQ -->
![deps overview flowchart](https://user-images.githubusercontent.com/649496/111809843-675b1000-88a3-11eb-8b5c-d85d71cb5a25.png)

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

## License

The license for the command line tool itself (this repo) is TBD. The individual components are all open-source, usually MIT licensed.
