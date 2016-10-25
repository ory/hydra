package policy

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/compose"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/ladon"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var managers = map[string]ladon.Manager{}

func init() {
	localWarden, httpClient := compose.NewFirewall("hydra", "alice", fosite.Arguments{scope},
		&ladon.DefaultPolicy{
			ID:        "1",
			Subjects:  []string{"alice"},
			Resources: []string{"rn:hydra:policies<.*>"},
			Actions:   []string{"create", "get", "delete", "find"},
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
		ID:          uuid.New(),
		Description: "description",
		Subjects:    []string{"<peter>"},
		Effect:      ladon.AllowAccess,
		Resources:   []string{"<article|user>"},
		Actions:     []string{"view"},
		Conditions: ladon.Conditions{
			"ip": &ladon.CIDRCondition{
				CIDR: "1234",
			},
			"owner": &ladon.EqualsSubjectCondition{},
		},
	}

	for k, m := range managers {
		_, err := m.Get(p.ID)
		pkg.AssertError(t, true, err, k)
		pkg.AssertError(t, false, m.Create(p), k)

		time.Sleep(200 * time.Millisecond)

		res, err := m.Get(p.ID)
		pkg.AssertError(t, false, err, k)
		assert.Equal(t, p, res, "%s", k)

		ps, err := m.FindPoliciesForSubject("peter")
		pkg.RequireError(t, false, err, k)
		require.Len(t, ps, 1, "%s", k)
		assert.Equal(t, p, ps[0], "%s", k)

		ps, err = m.FindPoliciesForSubject("stan")
		pkg.AssertError(t, false, err, k)
		assert.Len(t, ps, 0, "%s", k)

		pkg.AssertError(t, false, m.Delete(p.ID), k)

		_, err = m.Get(p.ID)
		pkg.AssertError(t, true, err, k)
	}
}
