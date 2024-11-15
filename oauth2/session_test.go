// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/jwt"
	"github.com/ory/x/assertx"
	"github.com/ory/x/snapshotx"

	_ "embed"
)

//go:embed fixtures/v1.11.8-session.json
var v1118Session []byte

//go:embed fixtures/v1.11.9-session.json
var v1119Session []byte

func parseTime(t *testing.T, ts string) time.Time {
	out, err := time.Parse(time.RFC3339Nano, ts)
	require.NoError(t, err)
	return out
}

func TestUnmarshalSession(t *testing.T) {
	expect := &Session{
		DefaultSession: &openid.DefaultSession{
			Claims: &jwt.IDTokenClaims{
				JTI:                                 "",
				Issuer:                              "http://127.0.0.1:4444/",
				Subject:                             "foo@bar.com",
				Audience:                            []string{"auth-code-client"},
				Nonce:                               "mbxojlzlkefzmlecvrzfkmpm",
				ExpiresAt:                           parseTime(t, "0001-01-01T00:00:00Z"),
				IssuedAt:                            parseTime(t, "2022-08-25T09:21:04Z"),
				RequestedAt:                         parseTime(t, "2022-08-25T09:20:54Z"),
				AuthTime:                            parseTime(t, "2022-08-25T09:21:01Z"),
				AccessTokenHash:                     "",
				AuthenticationContextClassReference: "0",
				AuthenticationMethodsReferences:     []string{},
				CodeHash:                            "",
				Extra: map[string]interface{}{
					"sid":       "177e1f44-a1e9-415c-bfa3-8b62280b182d",
					"timestamp": 1723546027,
				},
			},
			Headers: &jwt.Headers{Extra: map[string]interface{}{
				"kid": "public:hydra.openid.id-token",
			}},
			ExpiresAt: map[fosite.TokenType]time.Time{
				fosite.AccessToken:   parseTime(t, "2022-08-25T09:26:05Z"),
				fosite.AuthorizeCode: parseTime(t, "2022-08-25T09:23:04.432089764Z"),
				fosite.RefreshToken:  parseTime(t, "2022-08-26T09:21:05Z"),
			},
			Username: "",
			Subject:  "foo@bar.com",
		},
		Extra:                 map[string]interface{}{},
		KID:                   "public:hydra.jwt.access-token",
		ClientID:              "auth-code-client",
		ConsentChallenge:      "2261efbd447044a1b2f76b05c6aca164",
		ExcludeNotBeforeClaim: false,
		AllowedTopLevelClaims: []string{
			"persona_id",
			"persona_krn",
			"grantType",
			"market",
			"zone",
			"login_session_id",
		},
	}

	t.Run("v1.11.8", func(t *testing.T) {
		var actual Session
		require.NoError(t, json.Unmarshal(v1118Session, &actual))
		assertx.EqualAsJSON(t, expect, &actual)
		snapshotx.SnapshotTExcept(t, &actual, nil)
	})

	t.Run("v1.11.9" /* and later versions */, func(t *testing.T) {
		var actual Session
		require.NoError(t, json.Unmarshal(v1119Session, &actual))
		assertx.EqualAsJSON(t, expect, &actual)
		snapshotx.SnapshotTExcept(t, &actual, nil)
	})
}
