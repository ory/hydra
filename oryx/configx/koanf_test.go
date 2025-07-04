// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package configx

import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/spf13/pflag"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/stretchr/testify/require"
)

func newKoanf(ctx context.Context, schemaPath string, configPaths []string, modifiers ...OptionModifier) (*Provider, error) {
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, err
	}

	f := pflag.NewFlagSet("config", pflag.ContinueOnError)
	f.StringSliceP("config", "c", configPaths, "")

	modifiers = append(modifiers, WithFlags(f))
	k, err := New(ctx, schema, modifiers...)
	if err != nil {
		return nil, err
	}

	return k, nil
}

func setEnvs(t testing.TB, envs [][2]string) {
	for _, v := range envs {
		require.NoError(t, os.Setenv(v[0], v[1]))
	}
	t.Cleanup(func() {
		for _, v := range envs {
			_ = os.Unsetenv(v[0])
		}
	})
}

func BenchmarkNewKoanf(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	setEnvs(b, [][2]string{{"MUTATORS_HEADER_ENABLED", "true"}})
	schemaPath := path.Join("stub/benchmark/schema.config.json")
	for i := 0; i < b.N; i++ {
		_, err := newKoanf(ctx, schemaPath, []string{}, WithValues(map[string]interface{}{
			"dsn": "memory",
		}))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkKoanf(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	setEnvs(b, [][2]string{{"MUTATORS_HEADER_ENABLED", "true"}})
	schemaPath := path.Join("stub/benchmark/schema.config.json")
	k, err := newKoanf(ctx, schemaPath, []string{"stub/benchmark/benchmark.yaml"})
	require.NoError(b, err)

	keys := k.Koanf.Keys()
	numKeys := len(keys)

	b.Run("cache=false", func(b *testing.B) {
		var key string

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key = keys[i%numKeys]

			if k.Koanf.Get(key) == nil {
				b.Fatalf("cachedFind returned a nil value for key: %s", key)
			}
		}
	})

	b.Run("cache=true", func(b *testing.B) {
		for i, c := range []*ristretto.Config[string, any]{
			{
				NumCounters: int64(numKeys),
				MaxCost:     500000,
				BufferItems: 64,
			},
			{
				NumCounters: int64(numKeys * 10),
				MaxCost:     1000000,
				BufferItems: 64,
			},
			{
				NumCounters: int64(numKeys * 10),
				MaxCost:     5000000,
				BufferItems: 64,
			},
		} {
			cache, err := ristretto.NewCache[string, any](c)
			require.NoError(b, err)

			b.Run(fmt.Sprintf("config=%d", i), func(b *testing.B) {
				b.ResetTimer()
				for i := range b.N {
					key := keys[i%numKeys]

					val, found := cache.Get(key)
					if !found {
						val = k.Koanf.Get(key)
						_ = cache.Set(key, val, 0)
					}

					if val == nil {
						b.Fatalf("cachedFind returned a nil value for key: %s", key)
					}
				}
			})
		}
	})
}
