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
	"encoding/base64"

	"github.com/gtank/cryptopasta"
	"github.com/pkg/errors"
)

type AEAD struct {
	Key []byte
}

func (c *AEAD) Encrypt(plaintext []byte) (string, error) {
	if len(c.Key) < 32 {
		return "", errors.Errorf("Key must be 32 bytes, got %d bytes", len(c.Key))
	}

	var key [32]byte
	copy(key[:], c.Key[:32])

	ciphertext, err := cryptopasta.Encrypt(plaintext, &key)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (c *AEAD) Decrypt(ciphertext string) ([]byte, error) {
	if len(c.Key) < 32 {
		return []byte{}, errors.Errorf("Key must be longer 32 bytes, got %d bytes", len(c.Key))
	}

	var key [32]byte
	copy(key[:], c.Key[:32])

	raw, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return []byte{}, errors.WithStack(err)
	}

	plaintext, err := cryptopasta.Decrypt(raw, &key)
	return plaintext, nil
}
