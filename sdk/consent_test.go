package sdk

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/sessions"
	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/jwk"
	"github.com/ory-am/hydra/oauth2"
	"github.com/square/go-jose"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func genKey() *jose.JsonWebKeySet {
	g := &jwk.RS256Generator{}
	k, _ := g.Generate("")
	return k
}

func TestConsentHelper(t *testing.T) {
	km := &jwk.MemoryManager{Keys: map[string]*jose.JsonWebKeySet{}}
	km.AddKeySet(oauth2.ConsentChallengeKey, genKey())
	km.AddKeySet(oauth2.ConsentEndpointKey, genKey())

	_, err := km.GetKey(oauth2.ConsentChallengeKey, "private")
	require.Nil(t, err)
	c := Consent{KeyManager: km}
	s := oauth2.DefaultConsentStrategy{
		KeyManager:               km,
		DefaultChallengeLifespan: time.Hour,
	}

	ar := fosite.NewAuthorizeRequest()
	ar.Client = &fosite.DefaultClient{ID: "foobarclient"}
	challenge, err := s.IssueChallenge(ar, "http://hydra/oauth2/auth?client_id=foobarclient", &sessions.Session{Values: map[interface{}]interface{}{}})
	require.Nil(t, err)

	claims, err := c.VerifyChallenge(challenge)
	require.Nil(t, err)
	assert.Equal(t, claims.Audience, "foobarclient")
	assert.Equal(t, claims.RedirectURL, "http://hydra/oauth2/auth?client_id=foobarclient")
	assert.NotEmpty(t, claims.ID)

	resp, err := c.GenerateResponse(&ResponseRequest{
		Challenge: challenge,
		Subject:   "buzz",
		Scopes:    []string{"offline", "openid"},
	})
	require.Nil(t, err)

	var dec map[string]interface{}
	result, err := base64.RawURLEncoding.DecodeString(strings.Split(strings.Replace(resp, "http://hydra/oauth2/auth?client_id=foobarclient&consent=", "", -1), ".")[1])
	require.Nil(t, err)

	require.Nil(t, json.Unmarshal(result, &dec))
	assert.Equal(t, dec["jti"], claims.ID)
	t.Logf("%v", dec["jti"])
	assert.Equal(t, dec["scp"].([]interface{}), []interface{}{"offline", "openid"})
	assert.Equal(t, dec["sub"], "buzz")
}
