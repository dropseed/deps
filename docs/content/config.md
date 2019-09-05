# deps.yml

A `deps.yml` is only required if you need to make changes beyond what deps detects automatically.

```yaml
version: 3  # required!
dependencies:
- type: python
  path: app/server/requirements.txt
- type: js
```

*Note, this can also be named `.deps.yml` instead of `deps.yml`.*

## Manifest updates

A manifest is where you define your *direct* dependencies (like in `package.json`).

**When an out-of-range update is available for a direct dependency,
you'll get a pull request suggesting a new constraint to use.**

You can disable manifest updates entirely:

```yaml
version: 3
dependencies:
- type: python
  manifest_updates:
    enabled: false
```

### Examples of manifests

- `package.json` in npm
- `Pipfile` in Pipenv
- `composer.json` in Composer
- `Gemfile` in Bundler
- `Cargo.toml` in Cargo

### Disabling updates for a direct dependency

Use `manifest_updates.filters` to enable or disable updates on a per-dependency basis.

```yaml
version: 3
dependencies:
- type: python
  manifest_updates:
    # Filters are evaluated *in order*
    # so each dependency will use the first rule that it matches
    filters:
    - name: requests
      enabled: false
    # Typically your last filter will look like this,
    # which says any remaining matches should have updates enabled
    - name: .*
      enabled: true
```

### Grouping related updates

You can also use `manifest_updates.filters` to group related updates,
such as "react" and "react-dom". This way you'll get a *single* pull request that updates all of the react packages.

For example:

```yaml
version: 3
dependencies:
- type: python
  manifest_updates:
    filters:
    - name: react.*
      group: true
    - name: .*
```

## Lockfile updates

Most modern dependency managers have the concept of a "lockfile" (yarn.lock).
This is how you save the *exact* version of your direct and transitive dependencies that your app should be using.

**Whenever your lockfile is out of date,
deps will send you a single pull request that updates the entire lockfile.**

To disable lockfile updates, you can set `enabled: false` in your `deps.yml`.

```yaml
version: 3
dependencies:
- type: js
  lockfile_updates:
    enabled: false
```

### Examples of lockfiles

- `yarn.lock` in Yarn
- `package-lock.json` in npm
- `Pipfile.lock` in Pipenv
- `composer.lock` in Composer
- `Gemfile.lock` in Bundler
- `Cargo.lock` in Cargo

## Injecting commands (hooks)

WIP

## Environment variables

For each dependency type,
you can set `env` variables that will be set when that component runs.

*These must be strings!*

```yaml
version: 3
dependencies:
- type: js
  env:
    NODE_ENV: production
```

## Settings

TODO - as env vars or config settings
