# PHP

Currently supports:

- `composer.json`
- `composer.lock`

## Example `deps.yml`

```yaml
version: 3
dependencies:
- type: php
  path: app
  settings:
    # Set the options for composer install/update/require
    #
    # Default: "--no-progress --no-suggest"
    composer_options: "--no-scripts"
```

## Support

Any questions or issues with this specific component should be discussed in [GitHub issues](https://github.com/dropseed/deps-php/issues).

If there is private information which needs to be shared then please use the private support channels in [dependencies.io](https://www.dependencies.io/contact/).
