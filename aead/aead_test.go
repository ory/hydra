// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package aead_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"testing"

	"github.com/ory/hydra/v2/aead"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/internal"

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func secret(t *testing.T) string {
	bytes := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, bytes)
	require.NoError(t, err)
	return fmt.Sprintf("%X", bytes)
}

func TestAEAD(t *testing.T) {
	t.Parallel()
	for _, NewCipher := range []func(aead.Dependencies) aead.Cipher{
		func(d aead.Dependencies) aead.Cipher { return aead.NewAESGCM(d) },
		func(d aead.Dependencies) aead.Cipher { return aead.NewXChaCha20Poly1305(d) },
	} {
		NewCipher := NewCipher

		t.Run("case=without-rotation", func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			c := internal.NewConfigurationWithDefaults()
			c.MustSet(ctx, config.KeyGetSystemSecret, []string{secret(t)})
			a := NewCipher(c)

			plain := []byte(uuid.New())
			ct, err := a.Encrypt(ctx, plain, nil)
			assert.NoError(t, err)
			t.Log(ct)

			res, _, err := a.Decrypt(ctx, ct)
			assert.NoError(t, err)
			assert.Equal(t, plain, res)
		})

		t.Run("case=wrong-secret", func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			c := internal.NewConfigurationWithDefaults()
			c.MustSet(ctx, config.KeyGetSystemSecret, []string{secret(t)})
			a := NewCipher(c)

			ct, err := a.Encrypt(ctx, []byte(uuid.New()), nil)
			require.NoError(t, err)

			c.MustSet(ctx, config.KeyGetSystemSecret, []string{secret(t)})
			_, _, err = a.Decrypt(ctx, ct)
			require.Error(t, err)
		})

		t.Run("case=with-rotation", func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			c := internal.NewConfigurationWithDefaults()
			old := secret(t)
			c.MustSet(ctx, config.KeyGetSystemSecret, []string{old})
			a := NewCipher(c)

			plain := []byte(uuid.New())
			ct, err := a.Encrypt(ctx, plain, nil)
			require.NoError(t, err)

			// Sets the old secret as a rotated secret and creates a new one.
			c.MustSet(ctx, config.KeyGetSystemSecret, []string{secret(t), old})
			res, _, err := a.Decrypt(ctx, ct)
			require.NoError(t, err)
			assert.Equal(t, plain, res)

			// THis should also work when we re-encrypt the same plain text.
			ct2, err := a.Encrypt(ctx, plain, nil)
			require.NoError(t, err)
			assert.NotEqual(t, ct2, ct)

			res, _, err = a.Decrypt(ctx, ct)
			require.NoError(t, err)
			assert.Equal(t, plain, res)
		})

		t.Run("case=with-rotation-wrong-secret", func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			c := internal.NewConfigurationWithDefaults()
			c.MustSet(ctx, config.KeyGetSystemSecret, []string{secret(t)})
			a := NewCipher(c)

			plain := []byte(uuid.New())
			ct, err := a.Encrypt(ctx, plain, nil)
			require.NoError(t, err)

			// When the secrets do not match, an error should be thrown during decryption.
			c.MustSet(ctx, config.KeyGetSystemSecret, []string{secret(t), secret(t)})
			_, _, err = a.Decrypt(ctx, ct)
			require.Error(t, err)
		})

		t.Run("suite=with additional data", func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			c := internal.NewConfigurationWithDefaults()
			c.MustSet(ctx, config.KeyGetSystemSecret, []string{secret(t)})
			a := NewCipher(c)

			plain := []byte(uuid.New())
			ct, err := a.Encrypt(ctx, plain, []byte("additional data"))
			assert.NoError(t, err)

			t.Run("case=additional data matches", func(t *testing.T) {
				res, aad, err := a.Decrypt(ctx, ct)
				assert.NoError(t, err)
				assert.Equal(t, plain, res)
				assert.Equal(t, []byte("additional data"), aad)
			})
		})
	}
}
