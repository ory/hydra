//go:build go1.18
// +build go1.18

// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package hydra_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"hash"
	"testing"

	tokenhmac "github.com/ory/hydra/v2/fosite/token/hmac"
	fositejwt "github.com/ory/hydra/v2/fosite/token/jwt"
)

// testConfig satisfies all HMAC strategy config interfaces.
type testConfig struct {
	secret []byte
}

func (c *testConfig) GetTokenEntropy(_ context.Context) int             { return 32 }
func (c *testConfig) GetGlobalSecret(_ context.Context) ([]byte, error) { return c.secret, nil }
func (c *testConfig) GetRotatedGlobalSecrets(_ context.Context) ([][]byte, error) {
	return nil, nil
}
func (c *testConfig) GetHMACHasher(_ context.Context) func() hash.Hash {
	return sha512.New512_256
}

func mustMakeSecret(tb testing.TB) []byte {
	tb.Helper()
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		tb.Fatal(err)
	}
	return secret
}

// FuzzHMACTokenValidate tests HMAC token validation with arbitrary
// attacker-controlled token strings. This is the pre-auth boundary
// for every OAuth2 access token, refresh token, and authorization
// code in Hydra. The token comes from the Authorization: Bearer header.
func FuzzHMACTokenValidate(f *testing.F) {
	secret := mustMakeSecret(f)
	cfg := &testConfig{secret: secret}
	strategy := &tokenhmac.HMACStrategy{Config: cfg}

	// Generate a valid token for the seed corpus
	validToken, _, err := strategy.Generate(context.Background())
	if err == nil {
		f.Add(validToken)
	}

	f.Add("")                          // empty
	f.Add(".")                         // just separator
	f.Add("abc.def")                   // invalid base64
	f.Add("AAAA.AAAA")                 // valid base64, wrong signature
	f.Add("not-a-token")               // no separator
	f.Add(string(make([]byte, 10000))) // large input

	f.Fuzz(func(t *testing.T, token string) {
		if len(token) > 1<<16 {
			return
		}
		// Validate should never panic, only return errors
		_ = strategy.Validate(context.Background(), token)
	})
}

// FuzzHMACTokenSignature tests the Signature extraction method
// with arbitrary token strings. Signature is used to look up
// tokens in storage — a crash here is a DoS vector.
func FuzzHMACTokenSignature(f *testing.F) {
	secret := mustMakeSecret(f)
	cfg := &testConfig{secret: secret}
	strategy := &tokenhmac.HMACStrategy{Config: cfg}

	// Create a real token for valid case
	validToken, _, err := strategy.Generate(context.Background())
	if err == nil {
		f.Add(validToken, validToken)
	}

	f.Fuzz(func(t *testing.T, token, compare string) {
		_ = strategy.Signature(token)
	})
}

// FuzzHMACGenerateForString tests HMAC generation for arbitrary
// string inputs. Used to hash client secrets and other sensitive
// strings — a crash here can break secret hashing.
func FuzzHMACGenerateForString(f *testing.F) {
	secret := mustMakeSecret(f)
	cfg := &testConfig{secret: secret}
	strategy := &tokenhmac.HMACStrategy{Config: cfg}

	f.Add("test-string")
	f.Add("")
	f.Add(string(make([]byte, 10000)))

	f.Fuzz(func(t *testing.T, text string) {
		_, _ = strategy.GenerateHMACForString(context.Background(), text)
	})
}

// FuzzJWTValidate tests JWT token validation with arbitrary
// attacker-controlled JWT strings. This is the pre-auth boundary
// for JWT access tokens and OpenID Connect ID tokens.
func FuzzJWTValidate(f *testing.F) {
	// Generate a test RSA key
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		f.Fatal(err)
	}

	signer := &fositejwt.DefaultSigner{
		GetPrivateKey: func(_ context.Context) (interface{}, error) {
			return key, nil
		},
	}

	// Generate valid token as seed
	claims := fositejwt.MapClaims{
		"sub": "test-user",
		"iss": "https://hydra.example.com",
	}
	validToken, _, err := signer.Generate(context.Background(), claims, &fositejwt.Headers{})
	if err == nil {
		f.Add(validToken)
	}

	f.Add("")               // empty
	f.Add("eyJ...")         // garbage
	f.Add("a.b.c")          // 3-part but invalid
	f.Add("header.payload") // 2-part

	f.Fuzz(func(t *testing.T, token string) {
		if len(token) > 1<<16 {
			return
		}
		// Validate should never panic
		_, _ = signer.Validate(context.Background(), token)
	})
}

// FuzzJWTDecode tests JWT token decoding with arbitrary token strings.
// Decode is used even before Validate in many code paths.
func FuzzJWTDecode(f *testing.F) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		f.Fatal(err)
	}

	signer := &fositejwt.DefaultSigner{
		GetPrivateKey: func(_ context.Context) (interface{}, error) {
			return key, nil
		},
	}

	claims := fositejwt.MapClaims{"sub": "test"}
	validToken, _, err := signer.Generate(context.Background(), claims, &fositejwt.Headers{})
	if err == nil {
		f.Add(validToken)
	}

	f.Add("")
	f.Add("...")
	f.Add(string(make([]byte, 100000)))

	f.Fuzz(func(t *testing.T, token string) {
		if len(token) > 1<<16 {
			return
		}
		// Decode should never panic
		_, _ = signer.Decode(context.Background(), token)
	})
}

// FuzzJWTGetSignature tests JWT signature extraction with arbitrary
// token strings. The signature is used to identify tokens for revocation.
func FuzzJWTGetSignature(f *testing.F) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		f.Fatal(err)
	}

	signer := &fositejwt.DefaultSigner{
		GetPrivateKey: func(_ context.Context) (interface{}, error) {
			return key, nil
		},
	}

	validToken, _, err := signer.Generate(
		context.Background(),
		fositejwt.MapClaims{"sub": "test"},
		&fositejwt.Headers{},
	)
	if err == nil {
		f.Add(validToken)
	}

	f.Fuzz(func(t *testing.T, token string) {
		if len(token) > 1<<16 {
			return
		}
		_, _ = signer.GetSignature(context.Background(), token)
	})
}

// Ensure deterministic per-fuzz, not per-package
var _ = tokenhmac.RandomBytes
