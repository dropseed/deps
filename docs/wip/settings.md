# Dependencies.io component settings

When someone uses your component, they can pass in settings by adding them to
their `dependencies.yml` config. Those settings are given to the container as
JSON-encoded strings, with the name of `SETTING_{all-caps name from config}`.

For example:
```yaml
version: 2
dependencies:
- type: python
  path: requirements.txt
  settings:
    python_version: "3.5.6"
```

The container will be given `SETTING_PYTHON_VERSION="3.5.6"` which can be used
at runtime to further customize the behavior of the component.

A number of other settings are added to your component automatically by deps, like
`commit_message_prefix: "(chore) "` and other PR-specific options.
