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

package jwk

import (
	"crypto/rand"
	"io"
	"testing"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// RandomBytes returns n random bytes by reading from crypto/rand.Reader
func randomBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return []byte{}, errors.WithStack(err)
	}
	return bytes, nil
}

func TestAEAD(t *testing.T) {
	key, err := randomBytes(32)
	require.NoError(t, err)

	a := &AEAD{
		Key: key,
	}

	for i := 0; i < 100; i++ {
		plain := []byte(uuid.New())
		ct, err := a.Encrypt(plain)
		require.NoError(t, err)

		res, err := a.Decrypt(ct)
		require.NoError(t, err)
		assert.Equal(t, plain, res)
	}
}
