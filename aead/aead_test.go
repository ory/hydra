// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package aead_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"testing"

	"github.com/ory/hydra/v2/internal/testhelpers"

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/aead"
	"github.com/ory/hydra/v2/driver/config"
)

func secret(t *testing.T) string {
	bytes := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, bytes)
	require.NoError(t, err)
	return fmt.Sprintf("%X", bytes)
}

func TestAEAD(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		new  func(aead.Dependencies) aead.Cipher
	}{
		{"AES-GCM", func(d aead.Dependencies) aead.Cipher { return aead.NewAESGCM(d) }},
		{"XChaChaPoly", func(d aead.Dependencies) aead.Cipher { return aead.NewXChaCha20Poly1305(d) }},
	} {
		tc := tc

		t.Run("cipher="+tc.name, func(t *testing.T) {
			NewCipher := tc.new

			t.Run("case=without-rotation", func(t *testing.T) {
				t.Parallel()
				ctx := context.Background()
				c := testhelpers.NewConfigurationWithDefaults()
				c.MustSet(ctx, config.KeyGetSystemSecret, []string{secret(t)})
				a := NewCipher(c)

				plain := []byte(uuid.New())
				ct, err := a.Encrypt(ctx, plain, nil)
				assert.NoError(t, err)

				ct2, err := a.Encrypt(ctx, plain, nil)
				assert.NoError(t, err)
				assert.NotEqual(t, ct, ct2, "ciphertexts for the same plaintext must be different each time")

				res, err := a.Decrypt(ctx, ct, nil)
				assert.NoError(t, err)
				assert.Equal(t, plain, res)
			})

			t.Run("case=wrong-secret", func(t *testing.T) {
				t.Parallel()
				ctx := context.Background()
				c := testhelpers.NewConfigurationWithDefaults()
				c.MustSet(ctx, config.KeyGetSystemSecret, []string{secret(t)})
				a := NewCipher(c)

				ct, err := a.Encrypt(ctx, []byte(uuid.New()), nil)
				require.NoError(t, err)

				c.MustSet(ctx, config.KeyGetSystemSecret, []string{secret(t)})
				_, err = a.Decrypt(ctx, ct, nil)
				require.Error(t, err)
			})

			t.Run("case=with-rotation", func(t *testing.T) {
				t.Parallel()
				ctx := context.Background()
				c := testhelpers.NewConfigurationWithDefaults()
				old := secret(t)
				c.MustSet(ctx, config.KeyGetSystemSecret, []string{old})
				a := NewCipher(c)

				plain := []byte(uuid.New())
				ct, err := a.Encrypt(ctx, plain, nil)
				require.NoError(t, err)

				// Sets the old secret as a rotated secret and creates a new one.
				c.MustSet(ctx, config.KeyGetSystemSecret, []string{secret(t), old})
				res, err := a.Decrypt(ctx, ct, nil)
				require.NoError(t, err)
				assert.Equal(t, plain, res)

				// THis should also work when we re-encrypt the same plain text.
				ct2, err := a.Encrypt(ctx, plain, nil)
				require.NoError(t, err)
				assert.NotEqual(t, ct2, ct)

				res, err = a.Decrypt(ctx, ct, nil)
				require.NoError(t, err)
				assert.Equal(t, plain, res)
			})

			t.Run("case=with-rotation-wrong-secret", func(t *testing.T) {
				t.Parallel()
				ctx := context.Background()
				c := testhelpers.NewConfigurationWithDefaults()
				c.MustSet(ctx, config.KeyGetSystemSecret, []string{secret(t)})
				a := NewCipher(c)

				plain := []byte(uuid.New())
				ct, err := a.Encrypt(ctx, plain, nil)
				require.NoError(t, err)

				// When the secrets do not match, an error should be thrown during decryption.
				c.MustSet(ctx, config.KeyGetSystemSecret, []string{secret(t), secret(t)})
				_, err = a.Decrypt(ctx, ct, nil)
				require.Error(t, err)
			})

			t.Run("suite=with additional data", func(t *testing.T) {
				t.Parallel()
				ctx := context.Background()
				c := testhelpers.NewConfigurationWithDefaults()
				c.MustSet(ctx, config.KeyGetSystemSecret, []string{secret(t)})
				a := NewCipher(c)

				plain := []byte(uuid.New())
				ct, err := a.Encrypt(ctx, plain, []byte("additional data"))
				assert.NoError(t, err)

				t.Run("case=additional data matches", func(t *testing.T) {
					res, err := a.Decrypt(ctx, ct, []byte("additional data"))
					assert.NoError(t, err)
					assert.Equal(t, plain, res)
				})

				t.Run("case=additional data does not match", func(t *testing.T) {
					res, err := a.Decrypt(ctx, ct, []byte("wrong data"))
					assert.Error(t, err)
					assert.Nil(t, res)
				})

				t.Run("case=missing additional data", func(t *testing.T) {
					res, err := a.Decrypt(ctx, ct, nil)
					assert.Error(t, err)
					assert.Nil(t, res)
				})
			})
		})
	}
}
