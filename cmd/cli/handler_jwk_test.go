// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"testing"

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
