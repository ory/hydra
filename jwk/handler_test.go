package jwk_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/fosite"
	"github.com/ory/herodot"
	"github.com/ory/hydra/compose"
	. "github.com/ory/hydra/jwk"
	"github.com/ory/ladon"
	"github.com/square/go-jose"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testServer *httptest.Server
var IDKS *jose.JsonWebKeySet

func init() {
	localWarden, _ := compose.NewMockFirewall(
		"tests",
		"alice",
		fosite.Arguments{
			"hydra.keys.create",
			"hydra.keys.get",
			"hydra.keys.delete",
			"hydra.keys.update",
		}, &ladon.DefaultPolicy{
			ID:        "1",
			Subjects:  []string{"<.*>"},
			Resources: []string{"rn:hydra:keys:<[^:]+>:public"},
			Actions:   []string{"get"},
			Effect:    ladon.AllowAccess,
		},
	)
	router := httprouter.New()
	IDKS, _ = testGenerator.Generate("")

	h := Handler{
		Manager: &MemoryManager{},
		W:       localWarden,
		H:       herodot.NewJSONWriter(nil),
	}
	h.Manager.AddKeySet(IDTokenKeyName, IDKS)
	h.SetRoutes(router)
	testServer = httptest.NewServer(router)
}

func TestHandlerWellKnown(t *testing.T) {

	JWKPath := "/.well-known/jwks.json"
	res, err := http.Get(testServer.URL + JWKPath)
	require.NoError(t, err, "problem in http request")
	defer res.Body.Close()

	var known jose.JsonWebKeySet
	err = json.NewDecoder(res.Body).Decode(&known)
	require.NoError(t, err, "problem in decoding response")

	resp := known.Key("public")
	require.NotNil(t, resp, "Could not find key public")
	assert.Equal(t, resp, IDKS.Key("public"))
}
