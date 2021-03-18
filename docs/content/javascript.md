---
title: JavaScript
description: Automated updates and pull requests for JavaScript dependencies.
---

# JavaScript

Currently supports:

- `package.json`
- `package-lock.json`
- `yarn.lock`

## Example `deps.yml`

```yaml
version: 3
dependencies:
- type: js
  path: app  # a directory
  settings:
    # Enable updates for specific kinds of dependencies
    # in package.json.
    #
    # Default: [dependencies, devDependencies]
    manifest_package_types:
    - dependencies
```

## Support

Any questions or issues with this specific component should be discussed in [GitHub issues](https://github.com/dropseed/deps-js/issues).

If there is private information which needs to be shared then please use the private support channels in [dependencies.io](https://www.dependencies.io/contact/).
