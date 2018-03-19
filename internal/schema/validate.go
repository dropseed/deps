package schema

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/xeipuuv/gojsonschema"
)

const schemaString string = `
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "definitions": {
    "dependencies_schema.v2.schema.Lockfile": {
      "type": "object",
      "properties": {
        "current": {
          "$ref": "#/definitions/dependencies_schema.v2.schema.LockfileVersion"
        },
        "updated": {
          "$ref": "#/definitions/dependencies_schema.v2.schema.LockfileVersion"
        }
      },
      "required": [
        "current"
      ],
      "additionalProperties": false
    },
    "dependencies_schema.v2.schema.LockfileDependency": {
      "type": "object",
      "properties": {
        "source": {
          "type": "string"
        },
        "installed": {
          "$ref": "#/definitions/dependencies_schema.v2.schema.Version"
        },
        "constraint": {
          "type": "string"
        },
        "is_transitive": {
          "type": "boolean",
          "default": false
        }
      },
      "required": [
        "source",
        "installed"
      ],
      "additionalProperties": false
    },
    "dependencies_schema.v2.schema.LockfileVersion": {
      "type": "object",
      "properties": {
        "fingerprint": {
          "type": "string"
        },
        "dependencies": {
          "type": "object",
          "patternProperties": {
            ".*": {
              "$ref": "#/definitions/dependencies_schema.v2.schema.LockfileDependency"
            }
          }
        }
      },
      "required": [
        "fingerprint",
        "dependencies"
      ],
      "additionalProperties": false
    },
    "dependencies_schema.v2.schema.Manifest": {
      "type": "object",
      "properties": {
        "current": {
          "$ref": "#/definitions/dependencies_schema.v2.schema.ManifestVersion"
        },
        "updated": {
          "$ref": "#/definitions/dependencies_schema.v2.schema.ManifestVersion"
        },
        "lockfile_path": {
          "type": "string"
        }
      },
      "required": [
        "current"
      ],
      "additionalProperties": false
    },
    "dependencies_schema.v2.schema.ManifestDependency": {
      "type": "object",
      "properties": {
        "source": {
          "type": "string"
        },
        "constraint": {
          "type": "string"
        },
        "available": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/dependencies_schema.v2.schema.Version"
          }
        }
      },
      "required": [
        "source",
        "constraint"
      ],
      "additionalProperties": false
    },
    "dependencies_schema.v2.schema.ManifestVersion": {
      "type": "object",
      "properties": {
        "dependencies": {
          "type": "object",
          "patternProperties": {
            ".*": {
              "$ref": "#/definitions/dependencies_schema.v2.schema.ManifestDependency"
            }
          }
        }
      },
      "additionalProperties": false
    },
    "dependencies_schema.v2.schema.Version": {
      "type": "object",
      "properties": {
        "name": {
          "type": "string"
        },
        "content": {
          "type": "string"
        }
      },
      "required": [
        "name"
      ],
      "additionalProperties": false
    }
  },
  "type": "object",
  "properties": {
    "lockfiles": {
      "type": "object",
      "patternProperties": {
        ".*": {
          "$ref": "#/definitions/dependencies_schema.v2.schema.Lockfile"
        }
      }
    },
    "manifests": {
      "type": "object",
      "patternProperties": {
        ".*": {
          "$ref": "#/definitions/dependencies_schema.v2.schema.Manifest"
        }
      }
    }
  },
  "additionalProperties": false
}
`

func ValidateDependenciesJSONPath(path string) error {
	documentPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	documentLoader := gojsonschema.NewReferenceLoader("file://" + documentPath)

	schemaLoader := gojsonschema.NewStringLoader(schemaString)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err
	}
	if !result.Valid() {
		s := "The document is not valid. see errors:\n"
		for _, err := range result.Errors() {
			// Err implements the ResultError interface
			s = s + fmt.Sprintf("- %s\n", err)
		}
		return errors.New(s)
	}

	return nil
}
