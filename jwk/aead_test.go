/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package jwk_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"testing"

	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/internal"
	. "github.com/ory/hydra/jwk"

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
	ctx := context.Background()
	c := internal.NewConfigurationWithDefaults()
	t.Run("case=without-rotation", func(t *testing.T) {
		c.MustSet(ctx, config.KeyGetSystemSecret, []string{secret(t)})
		a := NewAEAD(c)

		plain := []byte(uuid.New())
		ct, err := a.Encrypt(ctx, plain)
		assert.NoError(t, err)

		res, err := a.Decrypt(ctx, ct)
		assert.NoError(t, err)
		assert.Equal(t, plain, res)
	})

	t.Run("case=wrong-secret", func(t *testing.T) {
		c.MustSet(ctx, config.KeyGetSystemSecret, []string{secret(t)})
		a := NewAEAD(c)

		ct, err := a.Encrypt(ctx, []byte(uuid.New()))
		require.NoError(t, err)

		c.MustSet(ctx, config.KeyGetSystemSecret, []string{secret(t)})
		_, err = a.Decrypt(ctx, ct)
		require.Error(t, err)
	})

	t.Run("case=with-rotation", func(t *testing.T) {
		old := secret(t)
		c.MustSet(ctx, config.KeyGetSystemSecret, []string{old})
		a := NewAEAD(c)

		plain := []byte(uuid.New())
		ct, err := a.Encrypt(ctx, plain)
		require.NoError(t, err)

		// Sets the old secret as a rotated secret and creates a new one.
		c.MustSet(ctx, config.KeyGetSystemSecret, []string{secret(t), old})
		res, err := a.Decrypt(ctx, ct)
		require.NoError(t, err)
		assert.Equal(t, plain, res)

		// THis should also work when we re-encrypt the same plain text.
		ct2, err := a.Encrypt(ctx, plain)
		require.NoError(t, err)
		assert.NotEqual(t, ct2, ct)

		res, err = a.Decrypt(ctx, ct)
		require.NoError(t, err)
		assert.Equal(t, plain, res)
	})

	t.Run("case=with-rotation-wrong-secret", func(t *testing.T) {
		c.MustSet(ctx, config.KeyGetSystemSecret, []string{secret(t)})
		a := NewAEAD(c)

		plain := []byte(uuid.New())
		ct, err := a.Encrypt(ctx, plain)
		require.NoError(t, err)

		// When the secrets do not match, an error should be thrown during decryption.
		c.MustSet(ctx, config.KeyGetSystemSecret, []string{secret(t), secret(t)})
		_, err = a.Decrypt(ctx, ct)
		require.Error(t, err)
	})
}
