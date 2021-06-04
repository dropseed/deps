<p align="center">
  <a href="https://dependencies.io/?utm_source=github&utm_medium=logo" target="_blank">
    <img src="https://user-images.githubusercontent.com/649496/111808669-375f3d00-88a2-11eb-9a25-2ee66a469b66.png" alt="Deps" height="140">
  </a>
</p>

# deps [![GitHub release](https://img.shields.io/github/release/dropseed/deps.svg)](https://github.com/dropseed/deps/releases)

**Deps is a command line tool for staying on top of dependencies. It runs updates, automates pull requests, and keeps your local installations in check.**

This repo contains the code for the `deps` command line tool itself (written in Go),
but each language/ecosystem has it's own repo which often uses the native language and acts as a light wrapper around the native dependency management tools.
This way deps can automate updates using the same tools that you would use in your terminal.

[Read the docs →](https://www.dependencies.io)

<!-- Edit in Excalidraw: https://excalidraw.com/#json=6195990008692736,unm4UYeUzdmbXTk_xcvIdQ -->
![deps overview flowchart](https://user-images.githubusercontent.com/649496/111809843-675b1000-88a3-11eb-8b5c-d85d71cb5a25.png)

The key features of deps are:

- **Native languages and tools**: The goal is to wrap the native package managers when possible (npm, yarn, pipenv, composer, etc.), so the updates delivered by deps are the same as updates you would make yourself on the command line.
- **Manifests vs Lockfiles**: If you use `"react": "^17.0.0"` in your package.json, we'll send you a pull request when 18.0.0 comes out. This is an out-of-range update to a direct dependency. But when react 17.1.2 is released, all you need to do is update your lockfile (package-lock.json or yarn.lock). In JavaScript, for example, your lockfile can be outdated daily between all of your direct and indirect (transitive) dependencies. These are in-range updates to direct and indirect (transitive) dependencies, and deps will send you a single rolling pull request to keep your lockfile up-to-date.
  ![Lockfile and manifest pull requests](https://user-images.githubusercontent.com/649496/119998663-87d7d280-bf96-11eb-8e73-4c686cc08c34.png)
- **Runs in an environment you control**: Deps runs in the same CI environment that you use for testing. You have full control over the container/host and system requirements.
- **Pluggable ecosystem**: We maintain a set of "official" components, but new or bespoke dependency types can be supported by pointing to a different component repo.

## Quick install

```console
$ curl https://deps.app/install.sh | bash -s -- -b $HOME/bin
```

## Official components

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
