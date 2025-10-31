// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package hmac

import (
	"context"
	"crypto/sha512"
	"fmt"
	"testing"

	"github.com/ory/hydra/v2/fosite"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateFailsWithShortCredentials(t *testing.T) {
	cg := HMACStrategy{Config: &fosite.Config{GlobalSecret: []byte("foo")}}
	challenge, signature, err := cg.Generate(context.Background())
	require.Error(t, err)
	require.Empty(t, challenge)
	require.Empty(t, signature)
}

func TestGenerate(t *testing.T) {
	ctx := context.Background()
	config := &fosite.Config{
		GlobalSecret: []byte("1234567890123456789012345678901234567890"),
	}
	cg := HMACStrategy{Config: config}

	for _, entropy := range []int{32, 64} {
		t.Run(fmt.Sprintf("entropy=%d", entropy), func(t *testing.T) {
			config.TokenEntropy = entropy

			token, signature, err := cg.Generate(ctx)
			require.NoError(t, err)
			require.NotEmpty(t, token)
			require.NotEmpty(t, signature)

			err = cg.Validate(ctx, token)
			require.NoError(t, err)

			actualSignature := cg.Signature(token)
			assert.Equal(t, signature, actualSignature)

			config.GlobalSecret = append([]byte("not"), config.GlobalSecret...)
			err = cg.Validate(ctx, token)
			assert.ErrorIs(t, err, fosite.ErrTokenSignatureMismatch)
		})
	}
}

func TestSignature(t *testing.T) {
	cg := HMACStrategy{}

	for token, expected := range map[string]string{
		"":            "",
		"foo":         "",
		"foo.bar":     "bar",
		"foo.bar.baz": "",
		".":           "",
	} {
		assert.Equal(t, expected, cg.Signature(token))
	}
}

func TestValidateSignatureRejects(t *testing.T) {
	cg := HMACStrategy{
		Config: &fosite.Config{GlobalSecret: []byte("1234567890123456789012345678901234567890")},
	}
	for k, c := range []string{
		"",
		" ",
		".",
		"foo.",
		".foo",
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			err := cg.Validate(context.Background(), c)
			assert.ErrorIs(t, err, fosite.ErrInvalidTokenFormat)
		})
	}

	err := cg.Validate(context.Background(), "foo.bar")
	assert.ErrorIs(t, err, fosite.ErrTokenSignatureMismatch)
}

func TestValidateWithRotatedKey(t *testing.T) {
	ctx := context.Background()
	oldGlobalSecret := []byte("1234567890123456789012345678901234567890")
	old := HMACStrategy{Config: &fosite.Config{GlobalSecret: oldGlobalSecret}}
	now := HMACStrategy{Config: &fosite.Config{
		GlobalSecret: []byte("0000000090123456789012345678901234567890"),
		RotatedGlobalSecrets: [][]byte{
			[]byte("abcdefgh90123456789012345678901234567890"),
			oldGlobalSecret,
		},
	}}

	token, _, err := old.Generate(ctx)
	require.NoError(t, err)

	assert.ErrorIs(t, now.Validate(ctx, "thisisatoken.withaninvalidsignature"), fosite.ErrTokenSignatureMismatch)
	assert.NoError(t, now.Validate(ctx, token))
}

func TestValidateWithRotatedKeyInvalid(t *testing.T) {
	ctx := context.Background()
	oldGlobalSecret := []byte("1234567890123456789012345678901234567890")
	old := HMACStrategy{Config: &fosite.Config{GlobalSecret: oldGlobalSecret}}
	now := HMACStrategy{Config: &fosite.Config{
		GlobalSecret: []byte("0000000090123456789012345678901234567890"),
		RotatedGlobalSecrets: [][]byte{
			[]byte("abcdefgh90123456789012345678901"),
			oldGlobalSecret,
		}},
	}

	token, _, err := old.Generate(ctx)
	require.NoError(t, err)

	require.EqualError(t, now.Validate(ctx, token), "secret for signing HMAC-SHA512/256 is expected to be 32 byte long, got 31 byte")

	require.EqualError(t, (&HMACStrategy{Config: &fosite.Config{}}).Validate(ctx, token), "a secret for signing HMAC-SHA512/256 is expected to be defined, but none were")
}

func TestCustomHMAC(t *testing.T) {
	ctx := context.Background()
	globalSecret := []byte("1234567890123456789012345678901234567890")
	defaultHasher := HMACStrategy{Config: &fosite.Config{
		GlobalSecret: globalSecret,
	}}
	sha512Hasher := HMACStrategy{Config: &fosite.Config{
		GlobalSecret: globalSecret,
		HMACHasher:   sha512.New,
	}}

	token, _, err := defaultHasher.Generate(ctx)
	require.NoError(t, err)
	require.ErrorIs(t, sha512Hasher.Validate(ctx, token), fosite.ErrTokenSignatureMismatch)

	token512, _, err := sha512Hasher.Generate(ctx)
	require.NoError(t, err)
	require.NoError(t, sha512Hasher.Validate(ctx, token512))
	require.ErrorIs(t, defaultHasher.Validate(ctx, token512), fosite.ErrTokenSignatureMismatch)
}

func TestGenerateFromString(t *testing.T) {
	cg := HMACStrategy{Config: &fosite.Config{
		GlobalSecret: []byte("1234567890123456789012345678901234567890")},
	}
	for _, c := range []struct {
		text string
		hash string
	}{
		{
			text: "",
			hash: "-n7EqD-bXkY3yYMH-ctEAGV8XLkU7Y6Bo6pbyT1agGA",
		},
		{
			text: " ",
			hash: "zXJvonHTNSOOGj_QKl4RpIX_zXgD2YfXUfwuDKaTTIg",
		},
		{
			text: "Test",
			hash: "TMeEaHS-cDC2nijiesCNtsOyBqHHtzWqAcWvceQT50g",
		},
		{
			text: "AnotherTest1234",
			hash: "zHYDOZGjzhVjx5r8RlBhpnJemX5JxEEBUjVT01n3IFM",
		},
	} {
		hash, _ := cg.GenerateHMACForString(context.Background(), c.text)
		assert.Equal(t, c.hash, hash)
	}
}
