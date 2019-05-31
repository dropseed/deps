# Dependencies.io JSON schema

The data that a component collects needs to follow a specific, universal format.
That format is defined by a [JSON schema](https://json-schema.org/) and the data
is automatically validated using the schema when you run `deps component
collect`.

```json
{
  "lockfiles": {
    "example_lockfile.json": {
      "current": {
        "dependencies": {
          "package1": {
            "installed": {
              "name": "1.1.0"
            },
            "source": "example-package-manager"
          }
        },
        "fingerprint": "d8db5538e62deadd2174b03d7b4ef7e2"
      },
      "updated": {
        "dependencies": {
          "package1": {
            "installed": {
              "name": "1.2.0"
            },
            "source": "example-package-manager"
          }
        },
        "fingerprint": "42c294a77caca9723baf339634a6b9ec"
      }
    }
  },
  "manifests": {
    "example_manifest.json": {
      "current": {
        "dependencies": {
          "package1": {
            "available": [
              {
                "name": "2.0.0"
              }
            ],
            "constraint": "> 1.0.0",
            "source": "example-package-manager"
          }
        }
      }
    }
  }
}
```
