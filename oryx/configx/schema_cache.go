// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package configx

import (
	"context"
	"crypto/sha256"

	"github.com/dgraph-io/ristretto/v2"

	"github.com/ory/jsonschema/v3"
)

var schemaCacheConfig = &ristretto.Config[[]byte, *jsonschema.Schema]{
	// Hold up to 25 schemas in cache. Usually we only need one.
	MaxCost:            25,
	NumCounters:        250,
	BufferItems:        64,
	Metrics:            false,
	IgnoreInternalCost: true,
	Cost: func(*jsonschema.Schema) int64 {
		return 1
	},
}

var schemaCache, _ = ristretto.NewCache(schemaCacheConfig)

func getSchema(ctx context.Context, schema []byte) (*jsonschema.Schema, error) {
	key := sha256.Sum256(schema)
	if val, found := schemaCache.Get(key[:]); found {
		return val, nil
	}

	schemaID, comp, err := newCompiler(schema)
	if err != nil {
		return nil, err
	}

	validator, err := comp.Compile(ctx, schemaID)
	if err != nil {
		return nil, err
	}

	schemaCache.Set(key[:], validator, 1)
	schemaCache.Wait()
	return validator, nil
}
