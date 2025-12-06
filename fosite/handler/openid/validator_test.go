// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package openid_test

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/openid"
	"github.com/ory/hydra/v2/fosite/token/jwt"
)

func TestValidatePrompt(t *testing.T) {
	config := &fosite.Config{
		MinParameterEntropy: fosite.MinParameterEntropy,
	}
	j := &openid.DefaultStrategy{
		Signer: &jwt.DefaultSigner{
			GetPrivateKey: func(_ context.Context) (interface{}, error) {
				return key, nil
			},
		},
		Config: &fosite.Config{
			MinParameterEntropy: fosite.MinParameterEntropy,
		},
	}

	v := openid.NewOpenIDConnectRequestValidator(j, config)

	genIDToken := func(c jwt.IDTokenClaims) string {
		s, _, err := j.Generate(context.TODO(), c.ToMapClaims(), jwt.NewHeaders())
		require.NoError(t, err)
		return s
	}

	for k, tc := range []struct {
		d           string
		prompt      string
		redirectURL string
		isPublic    bool
		expectErr   bool
		idTokenHint string
		s           *openid.DefaultSession
	}{
		{
			d:           "should fail because prompt=none should not work together with public clients and http non-localhost",
			prompt:      "none",
			isPublic:    true,
			expectErr:   true,
			redirectURL: "http://foo-bar/",
			s: &openid.DefaultSession{
				Subject: "foo",
				Claims: &jwt.IDTokenClaims{
					Subject:     "foo",
					RequestedAt: time.Now().UTC(),
					AuthTime:    time.Now().UTC().Add(-time.Minute),
				},
			},
		},
		{
			d:           "should pass because prompt=none works for public clients and http localhost",
			prompt:      "none",
			isPublic:    true,
			expectErr:   false,
			redirectURL: "http://localhost/",
			s: &openid.DefaultSession{
				Subject: "foo",
				Claims: &jwt.IDTokenClaims{
					Subject:     "foo",
					RequestedAt: time.Now().UTC(),
					AuthTime:    time.Now().UTC().Add(-time.Minute),
				},
			},
		},
		{
			d:           "should pass",
			prompt:      "none",
			isPublic:    true,
			expectErr:   false,
			redirectURL: "https://foo-bar/",
			s: &openid.DefaultSession{
				Subject: "foo",
				Claims: &jwt.IDTokenClaims{
					Subject:     "foo",
					RequestedAt: time.Now().UTC(),
					AuthTime:    time.Now().UTC().Add(-time.Minute),
				},
			},
		},
		{
			d:         "should fail because prompt=none requires an auth time being set",
			prompt:    "none",
			isPublic:  false,
			expectErr: true,
			s: &openid.DefaultSession{
				Subject: "foo",
				Claims: &jwt.IDTokenClaims{
					Subject:     "foo",
					RequestedAt: time.Now().UTC(),
				},
			},
		},
		{
			d:         "should fail because prompt=none and auth time is recent (after requested at)",
			prompt:    "none",
			isPublic:  false,
			expectErr: true,
			s: &openid.DefaultSession{
				Subject: "foo",
				Claims: &jwt.IDTokenClaims{
					Subject:     "foo",
					RequestedAt: time.Now().UTC().Add(-time.Minute),
					AuthTime:    time.Now().UTC(),
				},
			},
		},
		{
			d:         "should pass because prompt=none and auth time is in the past (before requested at)",
			prompt:    "none",
			isPublic:  false,
			expectErr: false,
			s: &openid.DefaultSession{
				Subject: "foo",
				Claims: &jwt.IDTokenClaims{
					Subject:     "foo",
					RequestedAt: time.Now().UTC(),
					AuthTime:    time.Now().UTC().Add(-time.Minute),
				},
			},
		},
		{
			d:         "should fail because prompt=none can not be used together with other prompts",
			prompt:    "none login",
			isPublic:  false,
			expectErr: true,
			s: &openid.DefaultSession{
				Subject: "foo",
				Claims: &jwt.IDTokenClaims{
					Subject:     "foo",
					RequestedAt: time.Now().UTC(),
					AuthTime:    time.Now().UTC(),
				},
			},
		},
		{
			d:         "should fail because prompt=foo is an unknown value",
			prompt:    "foo",
			isPublic:  false,
			expectErr: true,
			s: &openid.DefaultSession{
				Subject: "foo",
				Claims: &jwt.IDTokenClaims{
					Subject:     "foo",
					RequestedAt: time.Now().UTC(),
					AuthTime:    time.Now().UTC(),
				},
			},
		},
		{
			d:         "should pass because requesting consent and login works with public clients",
			prompt:    "login consent",
			isPublic:  true,
			expectErr: false,
			s: &openid.DefaultSession{
				Subject: "foo",
				Claims: &jwt.IDTokenClaims{
					Subject:     "foo",
					RequestedAt: time.Now().UTC().Add(-time.Second * 5),
					AuthTime:    time.Now().UTC().Add(-time.Second),
				},
			},
		},
		{
			d:         "should pass because requesting consent and login works with confidential clients",
			prompt:    "login consent",
			isPublic:  false,
			expectErr: false,
			s: &openid.DefaultSession{
				Subject: "foo",
				Claims: &jwt.IDTokenClaims{
					Subject:     "foo",
					RequestedAt: time.Now().UTC().Add(-time.Second * 5),
					AuthTime:    time.Now().UTC().Add(-time.Second),
				},
			},
		},
		{
			d:         "should fail subject from ID token does not match subject from session",
			prompt:    "login",
			isPublic:  false,
			expectErr: true,
			s: &openid.DefaultSession{
				Subject: "foo",
				Claims: &jwt.IDTokenClaims{
					Subject:     "foo",
					RequestedAt: time.Now().UTC(),
					AuthTime:    time.Now().UTC().Add(-time.Second),
				},
			},
			idTokenHint: genIDToken(jwt.IDTokenClaims{
				Subject:     "bar",
				RequestedAt: time.Now(),
				ExpiresAt:   time.Now().Add(time.Hour),
			}),
		},
		{
			d:         "should pass subject from ID token matches subject from session",
			prompt:    "",
			isPublic:  false,
			expectErr: false,
			s: &openid.DefaultSession{
				Subject: "foo",
				Claims: &jwt.IDTokenClaims{
					Subject:     "foo",
					RequestedAt: time.Now().UTC(),
					AuthTime:    time.Now().UTC().Add(-time.Second),
				},
			},
			idTokenHint: genIDToken(jwt.IDTokenClaims{
				Subject:     "foo",
				RequestedAt: time.Now(),
				ExpiresAt:   time.Now().Add(time.Hour),
			}),
		},
		{
			d:         "should pass subject from ID token matches subject from session even though id token is expired",
			prompt:    "",
			isPublic:  false,
			expectErr: false,
			s: &openid.DefaultSession{
				Subject: "foo",
				Claims: &jwt.IDTokenClaims{
					Subject:     "foo",
					RequestedAt: time.Now().UTC(),
					AuthTime:    time.Now().UTC().Add(-time.Second),
					ExpiresAt:   time.Now().UTC().Add(-time.Second),
				},
			},
			idTokenHint: genIDToken(jwt.IDTokenClaims{
				Subject:     "foo",
				RequestedAt: time.Now(),
				ExpiresAt:   time.Now().Add(time.Hour),
			}),
		},
	} {
		t.Run(fmt.Sprintf("case=%d/description=%s", k, tc.d), func(t *testing.T) {
			t.Logf("%s", tc.idTokenHint)
			err := v.ValidatePrompt(context.TODO(), &fosite.AuthorizeRequest{
				Request: fosite.Request{
					Form:    url.Values{"prompt": {tc.prompt}, "id_token_hint": {tc.idTokenHint}},
					Client:  &fosite.DefaultClient{Public: tc.isPublic},
					Session: tc.s,
				},
				RedirectURI: parse(tc.redirectURL),
			})
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func parse(u string) *url.URL {
	o, _ := url.Parse(u)
	return o
}
