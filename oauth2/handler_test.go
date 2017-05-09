package oauth2

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/json"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/herodot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerWellKnown(t *testing.T) {
	h := &Handler{
		H:      herodot.NewJSONWriter(nil),
		Issuer: "http://hydra.localhost",
	}

	AuthPathT := "/oauth2/auth"
	TokenPathT := "/oauth2/token"
	JWKPathT := "/.well-known/jwks.json"

	r := httprouter.New()
	h.SetRoutes(r)
	ts := httptest.NewServer(r)

	res, err := http.Get(ts.URL + "/.well-known/openid-configuration")

	defer res.Body.Close()

	trueConfig := WellKnown{
		Issuer:        h.Issuer,
		AuthURL:       h.Issuer + AuthPathT,
		TokenURL:      h.Issuer + TokenPathT,
		JWKsURI:       h.Issuer + JWKPathT,
		SubjectTypes:  []string{"pairwise", "public"},
		SigningAlgs:   []string{"RS256"},
		ResponseTypes: []string{"code", "code id_token", "id_token", "token id_token", "token"},
	}
	var wellKnownResp WellKnown
	err = json.NewDecoder(res.Body).Decode(&wellKnownResp)
	require.NoError(t, err, "problem decoding wellknown json response: %+v", err)
	assert.Equal(t, trueConfig, wellKnownResp)
}
