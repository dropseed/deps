---
title: deps.yml
description: Advanced configuration of automated dependency updates using deps.yml.
---

# deps.yml

A `deps.yml` is only required if you need to make changes beyond what is configured automatically.

```yaml
# deps.yml
version: 3  # required!
dependencies:
- type: python
  path: app/server/requirements.txt
- type: js
```

## Lockfile updates

Most modern dependency managers have the concept of a "lockfile" (yarn.lock).
This is how you save the *exact* version of your direct and transitive dependencies that your app should be using.

**When your lockfile is outdated,
deps will send you a single pull request that updates the entire lockfile.**
This single pull request will include in-range updates to all of your direct *and* transitive dependencies.

[![Lockfile update pull request made by deps](/assets/img/screenshots/deps-lockfile-pr.png)](/assets/img/screenshots/deps-lockfile-pr.png)

To disable lockfile updates, you can set `enabled: false` in your `deps.yml`.

```yaml
# deps.yml
version: 3
dependencies:
- type: js
  lockfile_updates:
    enabled: false
```

### Examples of supported lockfiles

- `yarn.lock` in [Yarn](https://yarnpkg.com/)
- `package-lock.json` in [npm](https://www.npmjs.com/)
- `Pipfile.lock` in [Pipenv](https://docs.pipenv.org/)
- `poetry.lock` in [poetry](https://python-poetry.org/)
- `composer.lock` in [Composer](https://getcomposer.org/)

## Manifest updates

A manifest is where you define your *direct* dependencies (like in `package.json`).

**When an out-of-range update is available for a direct dependency,
you'll get a pull request suggesting a new constraint to use.**
In-range updates will be delivered as [lockfile updates](#lockfile-updates).

[![Manifest update pull request made by deps](/assets/img/screenshots/deps-manifest-pr.png)](/assets/img/screenshots/deps-manifest-pr.png)

You can disable manifest updates entirely:

```yaml
# deps.yml
version: 3
dependencies:
- type: python
  manifest_updates:
    enabled: false
```

### Examples of supported manifests

- `package.json` in [npm](https://www.npmjs.com/)
- `Pipfile` in [Pipenv](https://docs.pipenv.org/)
- `requirements.txt` in [Python/pip](https://pip.pypa.io/en/stable/user_guide/)
- `pyproject.toml` in [poetry](https://python-poetry.org/)
- `composer.json` in [Composer](https://getcomposer.org/)

### Disabling updates for a direct dependency

Use `manifest_updates.filters` to enable or disable updates on a per-dependency basis.

```yaml
# deps.yml
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
# deps.yml
version: 3
dependencies:
- type: python
  manifest_updates:
    filters:
    - name: react.*
      group: true
    - name: .*
```

## Injecting commands (hooks)

```yaml
# deps.yml
version: 3
dependencies:
- type: js
  settings:
    before_commit: npm run compile  # Only runs in CI
```

## Customizing commit messages

Add commit message prefixes, suffixes, and trailers by providing your own template for the commit message.
The template is rendered using [Go's text/template package](https://golang.org/pkg/text/template/).

```yaml
# deps.yml
version: 3
dependencies:
- type: js
  settings:
    ## Variables
    # Single line subject (ex. "Update x from 1.0 to 2.0")
    # {{.Subject}}
    # Expanded body description (if available)
    # {{.Body}}
    # Combined subject + \n\n + optional body
    # {{.SubjectAndBody}}

    # Default
    commit_message_template: "{{.SubjectAndBody}}"

    # Subject prefix example
    commit_message_template: "deps: {{.SubjectAndBody}}"

    # Simplified subject w/ suffix example
    commit_message_template: "{{.Subject}} (skip ci)"

    # Trailer example
    commit_message_template: |-
      {{.SubjectAndBody}}

      Changelog: updated
```

## Environment variables

For each dependency type,
you can set `env` variables that will be set when that component runs.

*These must be strings!*

```yaml
# deps.yml
version: 3
dependencies:
- type: js
  env:
    NODE_ENV: production
```

## Settings

Most components have `settings` to further specify how they work.

```yaml
# deps.yml
version: 3
dependencies:
- type: js
  settings:
    github_labels:
    - dependencies
```

Settings can be more complex types and will be passed to the component as `DEPS_SETTING_{NAME}={JSON encoded value}`.

If you do *not* have a `deps.yml`,
you can also pass settings manually
(and for every component)
by using an env variable in your CI.
This is an easy way to apply the same GitHub PR labels to all updates, for example:
```console
$ DEPS_SETTING_GITHUB_LABELS='["dependencies"]' deps ci
```

### Filter settings

Settings can also be configured for specific dependencies via filters.

```yaml
# deps.yml
version: 3
dependencies:
- type: python
  manifest_updates:
    filters:
    - name: requests
      enabled: false
      settings:
        github_labels:
        - requests
    - name: .*
      enabled: true
```
