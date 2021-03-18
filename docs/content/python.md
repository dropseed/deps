---
title: Python
description: Automated updates and pull requests for Python dependencies.
---

# Python

Currently supports:

- `Pipfile`
- `poetry.lock`
- `requirements.txt` (filename doesn't matter)

## Example `deps.yml`

```yaml
version: 3
dependencies:
- type: python
  path: app/requirements.txt
  settings:
    # Enable updates for specific sections in Pipfile
    #
    # Default: ["packages", "dev-packages"]
    pipfile_sections:
    - packages

    # Enable updates for specific sections in Pipfile.lock
    #
    # Default: ["default", "develop"]
    pipfilelock_sections:
    - default

    # If your virtualenv is in a specific location that
    # isn't found automatically
    #
    # Default: ""
    pip_path: myenv/bin/pip
```

## Support

Any questions or issues with this specific component should be discussed in [GitHub issues](https://github.com/dropseed/deps-python/issues).

If there is private information which needs to be shared then please use the private support channels in [dependencies.io](https://www.dependencies.io/contact/).
