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
 *
 */

// Package hmac is the default implementation for generating and validating challenges. It uses HMAC-SHA256 to
// generate and validate challenges.

package hmac

import (
	"encoding/base64"
	"fmt"
	"strings"
	"sync"

	"github.com/gtank/cryptopasta"
	"github.com/pkg/errors"

	"github.com/ory/fosite"
)

// HMACStrategy is responsible for generating and validating challenges.
type HMACStrategy struct {
	AuthCodeEntropy      int
	GlobalSecret         []byte
	RotatedGlobalSecrets [][]byte
	sync.Mutex
}

const (
	// key should be at least 256 bit long, making it
	minimumEntropy = 32

	// the secrets (client and global) should each have at least 16 characters making it harder to guess them
	minimumSecretLength = 32
)

var b64 = base64.URLEncoding.WithPadding(base64.NoPadding)

// Generate generates a token and a matching signature or returns an error.
// This method implements rfc6819 Section 5.1.4.2.2: Use High Entropy for Secrets.
func (c *HMACStrategy) Generate() (string, string, error) {
	c.Lock()
	defer c.Unlock()

	if len(c.GlobalSecret) < minimumSecretLength {
		return "", "", errors.Errorf("secret for signing HMAC-SHA256 is expected to be 32 byte long, got %d byte", len(c.GlobalSecret))
	}

	var signingKey [32]byte
	copy(signingKey[:], c.GlobalSecret)

	if c.AuthCodeEntropy < minimumEntropy {
		c.AuthCodeEntropy = minimumEntropy
	}

	// When creating secrets not intended for usage by human users (e.g.,
	// client secrets or token handles), the authorization server should
	// include a reasonable level of entropy in order to mitigate the risk
	// of guessing attacks.  The token value should be >=128 bits long and
	// constructed from a cryptographically strong random or pseudo-random
	// number sequence (see [RFC4086] for best current practice) generated
	// by the authorization server.
	tokenKey, err := RandomBytes(c.AuthCodeEntropy)
	if err != nil {
		return "", "", errors.WithStack(err)
	}

	signature := cryptopasta.GenerateHMAC(tokenKey, &signingKey)

	encodedSignature := b64.EncodeToString(signature)
	encodedToken := fmt.Sprintf("%s.%s", b64.EncodeToString(tokenKey), encodedSignature)
	return encodedToken, encodedSignature, nil
}

// Validate validates a token and returns its signature or an error if the token is not valid.
func (c *HMACStrategy) Validate(token string) error {
	keys := append([][]byte{c.GlobalSecret}, c.RotatedGlobalSecrets...)
	for _, key := range keys {
		if err := c.validate(key, token); err == nil {
			return nil
		} else if errors.Cause(err) == fosite.ErrTokenSignatureMismatch {
		} else {
			return err
		}
	}

	return errors.New("a secret for signing HMAC-SHA256 is expected to be defined, but none were")
}

func (c *HMACStrategy) validate(secret []byte, token string) error {
	if len(secret) < minimumSecretLength {
		return errors.Errorf("secret for signing HMAC-SHA256 is expected to be 32 byte long, got %d byte", len(secret))
	}

	var signingKey [32]byte
	copy(signingKey[:], secret)

	split := strings.Split(token, ".")
	if len(split) != 2 {
		return errors.WithStack(fosite.ErrInvalidTokenFormat)
	}

	tokenKey := split[0]
	tokenSignature := split[1]
	if tokenKey == "" || tokenSignature == "" {
		return errors.WithStack(fosite.ErrInvalidTokenFormat)
	}

	decodedTokenSignature, err := b64.DecodeString(tokenSignature)
	if err != nil {
		return errors.WithStack(err)
	}

	decodedTokenKey, err := b64.DecodeString(tokenKey)
	if err != nil {
		return errors.WithStack(err)
	}

	if !cryptopasta.CheckHMAC(decodedTokenKey, decodedTokenSignature, &signingKey) {
		// Hash is invalid
		return errors.WithStack(fosite.ErrTokenSignatureMismatch)
	}

	return nil
}

func (c *HMACStrategy) Signature(token string) string {
	split := strings.Split(token, ".")

	if len(split) != 2 {
		return ""
	}

	return split[1]
}
