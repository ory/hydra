// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ory/x/hasherx"
	"github.com/ory/x/otelx"
)

type hasherConfig struct {
	cost uint32
}

func (c hasherConfig) HasherPBKDF2Config(_ context.Context) *hasherx.PBKDF2Config {
	return &hasherx.PBKDF2Config{}
}

func (c hasherConfig) HasherBcryptConfig(_ context.Context) *hasherx.BCryptConfig {
	return &hasherx.BCryptConfig{Cost: c.cost}
}

func (c hasherConfig) GetHasherAlgorithm(_ context.Context) string { return hashAlgorithmPBKDF2 }
func (c hasherConfig) Tracer(_ context.Context) *otelx.Tracer      { return otelx.NewNoop(nil, nil) }

func TestHasher(t *testing.T) {
	for _, cost := range []uint32{1, 8, 10} {
		c := &hasherConfig{cost: cost}
		result, err := NewHasher(c, c).Hash(t.Context(), []byte("foobar"))
		require.NoError(t, err)
		require.NotEmpty(t, result)
	}
}

// TestBackwardsCompatibility confirms that hashes generated with v1.x work with v2.x.
func TestBackwardsCompatibility(t *testing.T) {
	c := new(hasherConfig)
	h := NewHasher(c, c)
	require.NoError(t, h.Compare(context.Background(), []byte("$2a$10$lsrJjLPOUF7I75s3339R2uwqpjSlYGfhFyg7YsPtrSoITVy5UF3B2"), []byte("secret")))
	require.NoError(t, h.Compare(context.Background(), []byte("$2a$10$O1jZhd3U0azpLXwTu0cHHuTDWsBFnTJVbeHTADNQJWPR4Zqs8ATKS"), []byte("secret")))
	require.Error(t, h.Compare(context.Background(), []byte("$2a$10$lsrJjLPOUF7I75s3339R2uwqpjSlYGfhFyg7YsPtrSoITVy5UF3B3"), []byte("secret")))
}

func BenchmarkHasher(b *testing.B) {
	for cost := uint32(1); cost <= 16; cost++ {
		b.Run(fmt.Sprintf("cost=%d", cost), func(b *testing.B) {
			for range b.N {
				c := &hasherConfig{cost: cost}
				result, err := NewHasher(c, c).Hash(b.Context(), []byte("foobar"))
				require.NoError(b, err)
				require.NotEmpty(b, result)
			}
		})
	}
}
