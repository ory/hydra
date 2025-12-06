// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

// Package hmac is the default implementation for generating and validating challenges. It uses SHA-512/256 to
// generate and validate challenges.

package hmac

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"strings"
	"sync"

	"github.com/ory/x/errorsx"

	"github.com/pkg/errors"

	"github.com/ory/hydra/v2/fosite"
)

type HMACStrategyConfigurator interface {
	fosite.TokenEntropyProvider
	fosite.GlobalSecretProvider
	fosite.RotatedGlobalSecretsProvider
	fosite.HMACHashingProvider
}

// HMACStrategy is responsible for generating and validating challenges.
type HMACStrategy struct {
	sync.Mutex
	Config HMACStrategyConfigurator
}

const (
	minimumEntropy      = 32
	minimumSecretLength = 32
)

var b64 = base64.URLEncoding.WithPadding(base64.NoPadding)

// Generate generates a token and a matching signature or returns an error.
// This method implements rfc6819 Section 5.1.4.2.2: Use High Entropy for Secrets.
func (c *HMACStrategy) Generate(ctx context.Context) (string, string, error) {
	c.Lock()
	defer c.Unlock()

	globalSecret, err := c.Config.GetGlobalSecret(ctx)
	if err != nil {
		return "", "", err
	}

	if len(globalSecret) < minimumSecretLength {
		return "", "", errors.Errorf("secret for signing HMAC-SHA512/256 is expected to be 32 byte long, got %d byte", len(globalSecret))
	}

	var signingKey [32]byte
	copy(signingKey[:], globalSecret)

	entropy := c.Config.GetTokenEntropy(ctx)
	if entropy < minimumEntropy {
		entropy = minimumEntropy
	}

	// When creating tokens not intended for usage by human users (e.g.,
	// client secrets or token handles), the authorization server should
	// include a reasonable level of entropy in order to mitigate the risk
	// of guessing attacks. The token value should be >=128 bits long and
	// constructed from a cryptographically strong random or pseudo-random
	// number sequence (see [RFC4086] for best current practice) generated
	// by the authorization server.
	tokenKey, err := RandomBytes(entropy)
	if err != nil {
		return "", "", errorsx.WithStack(err)
	}

	signature := c.generateHMAC(ctx, tokenKey, &signingKey)

	encodedSignature := b64.EncodeToString(signature)
	encodedToken := fmt.Sprintf("%s.%s", b64.EncodeToString(tokenKey), encodedSignature)
	return encodedToken, encodedSignature, nil
}

// Validate validates a token and returns its signature or an error if the token is not valid.
func (c *HMACStrategy) Validate(ctx context.Context, token string) (err error) {
	var keys [][]byte

	globalSecret, err := c.Config.GetGlobalSecret(ctx)
	if err != nil {
		return err
	}

	if len(globalSecret) > 0 {
		keys = append(keys, globalSecret)
	}

	rotatedSecrets, err := c.Config.GetRotatedGlobalSecrets(ctx)
	if err != nil {
		return err
	}

	keys = append(keys, rotatedSecrets...)

	if len(keys) == 0 {
		return errors.New("a secret for signing HMAC-SHA512/256 is expected to be defined, but none were")
	}

	for _, key := range keys {
		if err = c.validate(ctx, key, token); err == nil {
			return nil
		} else if errors.Is(err, fosite.ErrTokenSignatureMismatch) {
			// Continue to the next key. The error will be returned if it is the last key.
		} else {
			return err
		}
	}

	return err
}

func (c *HMACStrategy) validate(ctx context.Context, secret []byte, token string) error {
	if len(secret) < minimumSecretLength {
		return errors.Errorf("secret for signing HMAC-SHA512/256 is expected to be 32 byte long, got %d byte", len(secret))
	}

	var signingKey [32]byte
	copy(signingKey[:], secret)

	tokenKey, tokenSignature, ok := strings.Cut(token, ".")
	if !ok {
		return errorsx.WithStack(fosite.ErrInvalidTokenFormat)
	}

	if tokenKey == "" || tokenSignature == "" {
		return errorsx.WithStack(fosite.ErrInvalidTokenFormat)
	}

	decodedTokenSignature, err := b64.DecodeString(tokenSignature)
	if err != nil {
		return errorsx.WithStack(err)
	}

	decodedTokenKey, err := b64.DecodeString(tokenKey)
	if err != nil {
		return errorsx.WithStack(err)
	}

	expectedMAC := c.generateHMAC(ctx, decodedTokenKey, &signingKey)
	if !hmac.Equal(expectedMAC, decodedTokenSignature) {
		// Hash is invalid
		return errorsx.WithStack(fosite.ErrTokenSignatureMismatch)
	}

	return nil
}

func (*HMACStrategy) Signature(token string) string {
	split := strings.Split(token, ".")
	if len(split) != 2 {
		return ""
	}
	return split[1]
}

// GenerateHMACForString returns an HMAC for a string
func (c *HMACStrategy) GenerateHMACForString(ctx context.Context, text string) (string, error) {
	var signingKey [32]byte

	secrets, err := c.Config.GetGlobalSecret(ctx)
	if err != nil {
		return "", err
	}

	if len(secrets) < minimumSecretLength {
		return "", errors.Errorf("secret for signing HMAC-SHA512/256 is expected to be 32 byte long, got %d byte", len(secrets))
	}
	copy(signingKey[:], secrets)

	bytes := []byte(text)
	hashBytes := c.generateHMAC(ctx, bytes, &signingKey)

	b64 := base64.RawURLEncoding.EncodeToString(hashBytes)
	return b64, nil
}

func (c *HMACStrategy) generateHMAC(ctx context.Context, data []byte, key *[32]byte) []byte {
	hasher := c.Config.GetHMACHasher(ctx)
	if hasher == nil {
		hasher = sha512.New512_256
	}
	h := hmac.New(hasher, key[:])
	// sha512.digest.Write() always returns nil for err, the panic should never happen
	_, err := h.Write(data)
	if err != nil {
		panic(err)
	}
	return h.Sum(nil)
}
