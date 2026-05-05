// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package configx

import (
	"strings"

	"github.com/knadh/koanf/maps"
	"github.com/pkg/errors"

	"github.com/ory/jsonschema/v3"
	"github.com/ory/x/jsonschemax"
)

type KoanfSchemaDefaults struct {
	keys []jsonschemax.Path
}

func NewKoanfSchemaDefaults(rawSchema []byte, schema *jsonschema.Schema) (*KoanfSchemaDefaults, error) {
	keys, err := getSchemaPaths(rawSchema, schema)
	if err != nil {
		return nil, err
	}

	return &KoanfSchemaDefaults{keys: keys}, nil
}

func (k *KoanfSchemaDefaults) ReadBytes() ([]byte, error) {
	return nil, errors.New("schema defaults provider does not support this method")
}

func (k *KoanfSchemaDefaults) Read() (map[string]interface{}, error) {
	values := map[string]interface{}{}
	for _, key := range k.keys {
		// It's an array!
		if strings.Contains(key.Name, "#") {
			continue
		}

		if key.Default != nil {
			values[key.Name] = key.Default
		}
	}

	// Deep-copy: getSchemaPaths caches Path entries by schema hash, so
	// key.Default may be a map/slice shared across all callers. koanf.Merge
	// retains references when merging, and recursively writes into shared
	// sub-maps when a later provider overlaps the default's path —
	// without this copy, concurrent configx.New calls race on the shared
	// default value.
	return maps.Unflatten(maps.Copy(values), "."), nil
}
