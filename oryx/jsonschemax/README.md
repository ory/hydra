# JSON Schema Helpers

This package contains utilities for working with JSON Schemas.

## Listing all Possible JSON Schema Paths

Using `jsonschemax.ListPaths()` you can get a list of all possible JSON paths in
a JSON Schema.

```go
package main

import (
	"bytes"
	"fmt"
	"github.com/ory/jsonschema/v3"
	"github.com/ory/x/jsonschemax"
)

var schema = "..."

func main() {
	c := jsonschema.NewCompiler()
	_ = c.AddResource("test.json", bytes.NewBufferString(schema))
	paths, _ := jsonschemax.ListPaths("test.json", c)
	fmt.Printf("%+v", paths)
}
```

All keys are delimited using `.`. Please note that arrays are denoted with `#`
when `ListPathsWithArraysIncluded` is used. For example, the JSON Schema

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "properties": {
    "providers": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "id": {
            "type": "string"
          }
        }
      }
    }
  }
}
```

Results in paths:

```json
[
  {
    "Title": "",
    "Description": "",
    "Examples": null,
    "Name": "providers",
    "Default": null,
    "Type": [],
    "TypeHint": 5,
    "Format": "",
    "Pattern": null,
    "Enum": null,
    "Constant": null,
    "ReadOnly": false,
    "MinLength": -1,
    "MaxLength": -1,
    "Required": false,
    "Minimum": null,
    "Maximum": null,
    "MultipleOf": null,
    "CustomProperties": null
  },
  {
    "Title": "",
    "Description": "",
    "Examples": null,
    "Name": "providers.#",
    "Default": null,
    "Type": {},
    "TypeHint": 5,
    "Format": "",
    "Pattern": null,
    "Enum": null,
    "Constant": null,
    "ReadOnly": false,
    "MinLength": -1,
    "MaxLength": -1,
    "Required": false,
    "Minimum": null,
    "Maximum": null,
    "MultipleOf": null,
    "CustomProperties": null
  },
  {
    "Title": "",
    "Description": "",
    "Examples": null,
    "Name": "providers.#.id",
    "Default": null,
    "Type": "",
    "TypeHint": 1,
    "Format": "",
    "Pattern": null,
    "Enum": null,
    "Constant": null,
    "ReadOnly": false,
    "MinLength": -1,
    "MaxLength": -1,
    "Required": false,
    "Minimum": null,
    "Maximum": null,
    "MultipleOf": null,
    "CustomProperties": null
  }
]
```
