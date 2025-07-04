// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package configx

import (
	"crypto/sha256"

	"github.com/ory/x/jsonschemax"

	"github.com/dgraph-io/ristretto/v2"

	"github.com/ory/jsonschema/v3"
)

var schemaPathCacheConfig = &ristretto.Config[[]byte, []jsonschemax.Path]{
	// Hold up to 25 schemas in cache. Usually we only need one.
	MaxCost:            250,
	NumCounters:        2500,
	BufferItems:        64,
	Metrics:            false,
	IgnoreInternalCost: true,
}

var schemaPathCache, _ = ristretto.NewCache[[]byte, []jsonschemax.Path](schemaPathCacheConfig)

func getSchemaPaths(rawSchema []byte, schema *jsonschema.Schema) ([]jsonschemax.Path, error) {
	key := sha256.Sum256(rawSchema)
	if val, found := schemaPathCache.Get(key[:]); found {
		return val, nil
	}

	keys, err := jsonschemax.ListPathsWithInitializedSchemaAndArraysIncluded(schema)
	if err != nil {
		return nil, err
	}

	schemaPathCache.Set(key[:], keys, 1)
	schemaPathCache.Wait()
	return keys, nil
}
