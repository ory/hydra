package policy

import (
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/handler/core"
	"github.com/ory-am/hydra/herodot"
	ioa2 "github.com/ory-am/hydra/oauth2"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/hydra/warden"
	"github.com/ory-am/ladon"
	"github.com/ory-am/ladon/memory"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

var managers = map[string]ladon.Manager{}

var fositeStore = pkg.FositeStore()

var ladonWarden = pkg.LadonWarden(map[string]ladon.Policy{
	"1": &ladon.DefaultPolicy{
		ID:        "1",
		Subjects:  []string{"alice"},
		Resources: []string{"rn:hydra:policies<.*>"},
		Actions:   []string{"create", "get", "delete", "search"},
		Effect:    ladon.AllowAccess,
	},
})

var localWarden = &warden.LocalWarden{
	Warden: ladonWarden,
	TokenValidator: &core.CoreValidator{
		AccessTokenStrategy: pkg.HMACStrategy,
		AccessTokenStorage:  fositeStore,
	},
	Issuer: "tests",
}

var tokens = pkg.Tokens(1)

func init() {
	ar := fosite.NewAccessRequest(&ioa2.Session{Subject: "alice"})
	ar.GrantedScopes = fosite.Arguments{Scope}
	fositeStore.CreateAccessTokenSession(nil, tokens[0][0], ar)

	h := &Handler{
		Manager: &memory.Manager{
			Policies: map[string]ladon.Policy{},
		},
		W: localWarden,
		H: new(herodot.JSON),
	}
	r := httprouter.New()
	h.SetRoutes(r)
	ts := httptest.NewServer(r)
	conf := &oauth2.Config{Scopes: []string{}, Endpoint: oauth2.Endpoint{}}
	u, _ := url.Parse(ts.URL + Endpoint)
	managers["http"] = &HTTPManager{
		Endpoint: u,
		Client: conf.Client(oauth2.NoContext, &oauth2.Token{
			AccessToken: tokens[0][1],
			Expiry:      time.Now().Add(time.Hour),
			TokenType:   "bearer",
		}),
	}
}

func TestManagers(t *testing.T) {
	p := &ladon.DefaultPolicy{
		ID:       uuid.New(),
		Subjects: []string{"peter", "max"},
	}

	for k, m := range managers {
		_, err := m.Get(p.ID)
		pkg.RequireError(t, true, err, k)
		pkg.RequireError(t, false, m.Create(p))

		res, err := m.Get(p.ID)
		pkg.RequireError(t, false, err, k)
		assert.Equal(t, p, res, "%s", k)

		ps, err := m.FindPoliciesForSubject("peter")
		pkg.RequireError(t, false, err, k)
		assert.Len(t, ps, 1, "%s", k)
		assert.Equal(t, p, ps[0], "%s", k)

		ps, err = m.FindPoliciesForSubject("stan")
		pkg.RequireError(t, false, err, k)
		assert.Len(t, ps, 0, "%s", k)

		pkg.RequireError(t, false, m.Delete(p.ID), k)

		_, err = m.Get(p.ID)
		pkg.RequireError(t, true, err, k)
	}
}
