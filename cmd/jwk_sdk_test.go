// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"encoding/json"
	"testing"

	jose "github.com/go-jose/go-jose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/x/josex"
)

func Test_toSDKFriendlyJSONWebKey(t *testing.T) {
	publicJWK := []byte(`{
		"kty": "RSA",
		"e": "AQAB",
		"use": "sig",
		"kid": "7a5ff76a-6766-11ea-bc55-0242ac130003",
		"alg": "RS256",
		"n": "l80jJJqcc1PpefIGVIjuPvA1D7NscnuF9aQqLa7I9rDUK4IaSOO3kL_EF13k-jTzcA5q4OZn5dR0kmqIMZT2gQ"
	}`)

	publicPEM := []byte(`
		-----BEGIN PUBLIC KEY-----
		MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBAPf64dykufSkwnvUiBAwd5Si0K6t4m5i
		qJD8TmLJCmFjKaOUa6nszcFt/FkAuORfdlrD9mEZLPrPx74RSluyTBMCAwEAAQ==
		-----END PUBLIC KEY-----
	`)

	type args struct {
		key []byte
		kid string
		use string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "JWK with algorithm",
			args: args{
				key: publicJWK,
				kid: "public:7a5ff76a-6766-11ea-bc55-0242ac130003",
				use: "sig",
			},
			want: "RS256",
		},
		{
			name: "PEM key without algorithm",
			args: args{
				key: publicPEM,
				kid: "public:7a5ff76a-6766-11ea-bc55-0242ac130003",
				use: "sig",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, _ := josex.LoadPublicKey(tt.args.key)
			if got := ToSDKFriendlyJSONWebKey(key, tt.args.kid, tt.args.use); got.Algorithm != tt.want {
				t.Errorf("toSDKFriendlyJSONWebKey() = %v, want %v", got.Algorithm, tt.want)
			}
		})
	}
}

func TestOnlyPublicSDKKeys(t *testing.T) {
	set, err := jwk.GenerateJWK(jose.RS256, "test-id-1", "sig")
	require.NoError(t, err)

	out, err := json.Marshal(set.Keys)
	require.NoError(t, err)

	var sdkSet []hydra.JsonWebKey
	require.NoError(t, json.Unmarshal(out, &sdkSet))

	assert.NotEmpty(t, sdkSet[0].P)
	result, err := OnlyPublicSDKKeys(sdkSet)
	require.NoError(t, err)

	assert.Empty(t, result[0].P)
}
