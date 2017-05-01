package jwk_test

import (
	"net/http/httptest"
	"testing"
	. "github.com/ory-am/hydra/jwk"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/compose"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/ladon"
	"net/http"
	"encoding/json"
	"github.com/square/go-jose"
	"github.com/docker/docker/pkg/testutil/assert"
)


var testServer *httptest.Server
var IDKS, CCKS, CRKS *jose.JsonWebKeySet

func init() {
	localWarden, _ := compose.NewFirewall(
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
			Resources: []string{"rn:hydra:keys:<[^:]+>:public:<[^:]+>"},
			Actions:   []string{"get"},
			Effect:    ladon.AllowAccess,
		},
	)
	router := httprouter.New()
	IDKS, _ = testGenerator.Generate(IDTokenKeyName)
	CCKS, _ = testGenerator.Generate(ConsentChallengeKeyName)
	CRKS, _ = testGenerator.Generate(ConsentResponseKeyName)

	h := Handler{
		Manager: &MemoryManager{},
		W:       localWarden,
		H:       &herodot.JSON{},
	}
	h.Manager.AddKeySet(IDTokenKeyName, IDKS)
	h.Manager.AddKeySet(ConsentChallengeKeyName, CCKS)
	h.Manager.AddKeySet(ConsentResponseKeyName, CRKS)
	h.SetRoutes(router)
	testServer = httptest.NewServer(router)
}



func TestHandlerWellKnown(t *testing.T) {

	JWKPath := "/.well-known/jwks.json"
	res, err := http.Get(testServer.URL + JWKPath)
	if err != nil {
		t.Errorf("problem in http request: %v", err)
	}
	defer res.Body.Close()

	var known jose.JsonWebKeySet
	err = json.NewDecoder(res.Body).Decode(&known)
	if err != nil {
		t.Errorf("problem decoding well known response: %v", err)
	}
	sets := map[string]*jose.JsonWebKeySet{
		ConsentChallengeKeyName: CCKS,
		IDTokenKeyName: IDKS,
		ConsentResponseKeyName: CRKS,
	}
	for k, v := range sets {
		resp := known.Key("public:" + k)
		if resp == nil {
			t.Errorf("could not find key public: %v", k)
		}
		assert.DeepEqual(t, resp, v.Key("public:" + k))
	}
}
