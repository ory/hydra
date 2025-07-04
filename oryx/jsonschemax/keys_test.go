// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jsonschemax

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/ory/x/snapshotx"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/require"

	"github.com/ory/jsonschema/v3"
)

const recursiveSchema = `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "test.json",
  "definitions": {
    "foo": {
      "type": "object",
      "properties": {
		"bars": {
			"type": "string",
			"format": "email",
			"pattern": ".*"
		},
        "bar": {
          "$ref": "#/definitions/bar"
        }
      },
      "required":["bars"]
    },
    "bar": {
      "type": "object",
      "properties": {
		"foos": {
		  "type": "string",
		  "minLength": 1,
		  "maxLength": 10
		},
        "foo": {
          "$ref": "#/definitions/foo"
        }
      }
    }
  },
  "type": "object",
  "properties": {
    "bar": {
      "$ref": "#/definitions/bar"
    }
  }
}`

func readFile(t *testing.T, path string) string {
	schema, err := os.ReadFile(path)
	require.NoError(t, err)
	return string(schema)
}

const fooExtensionName = "fooExtension"

type (
	extensionConfig struct {
		NotAJSONSchemaKey string `json:"not-a-json-schema-key"`
	}
)

func fooExtensionCompile(_ jsonschema.CompilerContext, m map[string]interface{}) (interface{}, error) {
	if raw, ok := m[fooExtensionName]; ok {
		var b bytes.Buffer
		if err := json.NewEncoder(&b).Encode(raw); err != nil {
			return nil, errors.WithStack(err)
		}

		var e extensionConfig
		if err := json.NewDecoder(&b).Decode(&e); err != nil {
			return nil, errors.WithStack(err)
		}

		return &e, nil
	}
	return nil, nil
}

func fooExtensionValidate(_ jsonschema.ValidationContext, _, _ interface{}) error {
	return nil
}

func (ec *extensionConfig) EnhancePath(p Path) map[string]interface{} {
	if ec.NotAJSONSchemaKey != "" {
		fmt.Printf("enhancing path: %s with custom property %s\n", p.Name, ec.NotAJSONSchemaKey)
		return map[string]interface{}{
			ec.NotAJSONSchemaKey: p.Name,
		}
	}
	return nil
}

func TestListPathsWithRecursion(t *testing.T) {
	for k, tc := range []struct {
		recursion uint8
		expected  interface{}
	}{
		{
			recursion: 5,
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			c := jsonschema.NewCompiler()
			require.NoError(t, c.AddResource("test.json", bytes.NewBufferString(recursiveSchema)))
			actual, err := ListPathsWithRecursion(context.Background(), "test.json", c, tc.recursion)
			require.NoError(t, err)

			snapshotx.SnapshotT(t, actual)
		})
	}
}

func TestListPaths(t *testing.T) {
	for k, tc := range []struct {
		schema    string
		expectErr bool
		extension *jsonschema.Extension
	}{
		{
			schema: readFile(t, "./stub/.oathkeeper.schema.json"),
		},
		{
			schema: readFile(t, "./stub/nested-simple-array.schema.json"),
		},
		{
			schema: readFile(t, "./stub/config.schema.json"),
		},
		{
			schema: readFile(t, "./stub/nested-array.schema.json"),
		},
		{
			// this should fail because of recursion
			schema:    recursiveSchema,
			expectErr: true,
		},
		{
			schema: `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "test.json",
  "oneOf": [
    {
      "type": "object",
      "properties": {
        "list": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "foo": {
          "default": false,
          "type": "boolean"
        },
        "bar": {
          "type": "boolean",
          "default": "asdf",
          "readOnly": true
        }
      }
    },
    {
      "type": "object",
      "properties": {
        "foo": {
          "type": "boolean"
        }
      }
    }
  ]
}`,
		},
		{
			schema: `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "test.json",
  "type": "object",
  "required": ["foo"],
  "properties": {
    "foo": {
      "type": "boolean"
    },
    "bar": {
      "type": "string",
      "fooExtension": {
        "not-a-json-schema-key": "foobar"
      }
    }
  }
}`,
			extension: &jsonschema.Extension{
				Meta:     nil,
				Compile:  fooExtensionCompile,
				Validate: fooExtensionValidate,
			},
		},
		{
			schema: `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "test.json",
  "type": "object",
  "definitions": {
    "foo": {
      "type": "string"
    }
  },
  "properties": {
    "bar": {
      "type": "array",
      "items": {
        "$ref": "#/definitions/foo"
      }
    }
  }
}`,
		},
		{
			schema: `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "test.json",
  "type": "object",
  "definitions": {
    "foo": {
      "type": "string"
    },
    "bar": {
      "type": "array",
      "items": {
        "$ref": "#/definitions/foo"
      },
      "required": ["foo"]
    }
  },
  "properties": {
    "baz": {
      "type": "array",
      "items": {
        "$ref": "#/definitions/bar"
      }
    }
  }
}`,
		},
		{
			schema: `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "test.json",
  "type": "object",
  "definitions": {
    "foo": {
      "type": "string"
    },
    "bar": {
      "type": "object",
      "properties": {
        "foo": {
          "$ref": "#/definitions/foo"
        }
      },
      "required": ["foo"]
    }
  },
  "properties": {
    "baz": {
      "type": "array",
      "items": {
        "$ref": "#/definitions/bar"
      }
    }
  }
}`,
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			c := jsonschema.NewCompiler()
			if tc.extension != nil {
				c.Extensions[fooExtensionName] = *tc.extension
			}

			require.NoError(t, c.AddResource("test.json", bytes.NewBufferString(tc.schema)))
			actual, err := ListPathsWithArraysIncluded(context.Background(), "test.json", c)
			if tc.expectErr {
				require.Error(t, err, "%+v", actual)
				return
			}
			require.NoError(t, err)

			snapshotx.SnapshotT(t, actual)
		})
	}
}
