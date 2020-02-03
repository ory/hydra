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

package x

import (
	"crypto/sha256"

	"github.com/ory/x/randx"
)

var secretCharSet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_-.~")

func GenerateSecret(length int) ([]byte, error) {
	secret, err := randx.RuneSequence(length, secretCharSet)
	if err != nil {
		return []byte{}, err
	}
	return []byte(string(secret)), nil
}

// HashStringSecret hashes the secret for consumption by the AEAD encryption algorithm which expects exactly 32 bytes.
//
// The system secret is being hashed to always match exactly the 32 bytes required by AEAD, even if the secret is long or
// shorter.
func HashStringSecret(secret string) []byte {
	return HashByteSecret([]byte(secret))
}

// HashByteSecret hashes the secret for consumption by the AEAD encryption algorithm which expects exactly 32 bytes.
//
// The system secret is being hashed to always match exactly the 32 bytes required by AEAD, even if the secret is long or
// shorter.
func HashByteSecret(secret []byte) []byte {
	r := sha256.Sum256(secret)
	return r[:]
}
