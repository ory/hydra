package policy

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/internal"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/ladon"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
)

var managers = map[string]ladon.Manager{}

func init() {
	localWarden, httpClient := internal.NewFirewall("hydra", "alice", fosite.Arguments{scope},
		&ladon.DefaultPolicy{
			ID:        "1",
			Subjects:  []string{"alice"},
			Resources: []string{"rn:hydra:policies<.*>"},
			Actions:   []string{"create", "get", "delete", "search"},
			Effect:    ladon.AllowAccess,
		},
	)

	h := &Handler{
		Manager: &ladon.MemoryManager{
			Policies: map[string]ladon.Policy{},
		},
		W: localWarden,
		H: new(herodot.JSON),
	}

	r := httprouter.New()
	h.SetRoutes(r)
	ts := httptest.NewServer(r)

	u, _ := url.Parse(ts.URL + endpoint)
	managers["http"] = &HTTPManager{
		Endpoint: u,
		Client:   httpClient,
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
