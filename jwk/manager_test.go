package jwk

import (
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/internal"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/ladon"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"net/url"
	"testing"
)

var managers = map[string]Manager{}

var testGenerator = &RS256Generator{}

var ts *httptest.Server

func init() {
	localWarden, httpClient := internal.NewFirewall(
		"tests",
		"alice",
		fosite.Arguments{
			"hydra.keys.create",
			"hydra.keys.get",
			"hydra.keys.delete",
			"hydra.keys.update",
		}, &ladon.DefaultPolicy{
			ID:        "1",
			Subjects:  []string{"alice"},
			Resources: []string{"rn:hydra:keys<.*>"},
			Actions:   []string{"create", "get", "delete", "update"},
			Effect:    ladon.AllowAccess,
		},
	)

	r := httprouter.New()
	h := Handler{
		Manager: &MemoryManager{},
		W:       localWarden,
		H:       &herodot.JSON{},
	}
	h.SetRoutes(r)
	ts := httptest.NewServer(r)
	u, _ := url.Parse(ts.URL + "/keys")
	managers["memory"] = &MemoryManager{}
	managers["http"] = &HTTPManager{Client: httpClient, Endpoint: u}
}

func TestManagerKey(t *testing.T) {
	ks, _ := testGenerator.Generate("")
	key := &ks.Key("public")[0]

	for name, m := range managers {
		_, err := m.GetKey("faz", "baz")
		pkg.AssertError(t, true, err, name)

		err = m.AddKey("faz", key)
		pkg.AssertError(t, false, err, name)

		got, err := m.GetKey("faz", "public")
		pkg.AssertError(t, false, err, name)
		assert.EqualValues(t, key, got, name)

		err = m.DeleteKey("faz", "public")
		pkg.AssertError(t, false, err, name)

		_, err = m.GetKey("faz", "public")
		pkg.AssertError(t, true, err, name)
	}
}


func TestManagerKeySet(t *testing.T) {
	ks, _ := testGenerator.Generate("")

	for name, m := range managers {
		_, err := m.GetKeySet("foo")
		pkg.AssertError(t, true, err, name)

		err = m.AddKeySet("bar", ks)
		pkg.AssertError(t, false, err, name)

		got, err := m.GetKeySet("bar")
		pkg.AssertError(t, false, err, name)
		assert.EqualValues(t, ks, got, name)

		err = m.DeleteKeySet("bar")
		pkg.AssertError(t, false, err, name)

		_, err = m.GetKeySet("bar")
		pkg.AssertError(t, true, err, name)
	}
}
