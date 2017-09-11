package policy

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/fosite"
	"github.com/ory/herodot"
	"github.com/ory/hydra/compose"
	"github.com/ory/ladon"
	"github.com/ory/ladon/manager/memory"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var managers = map[string]Manager{}

func init() {
	localWarden, httpClient := compose.NewMockFirewall("hydra", "alice", fosite.Arguments{scope},
		&ladon.DefaultPolicy{
			ID:        "1",
			Subjects:  []string{"alice"},
			Resources: []string{"rn:hydra:policies<.*>"},
			Actions:   []string{"create", "get", "delete", "list", "update"},
			Effect:    ladon.AllowAccess,
		},
	)

	h := &Handler{
		Manager: &memory.MemoryManager{
			Policies: map[string]ladon.Policy{},
		},
		W: localWarden,
		H: herodot.NewJSONWriter(nil),
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
		t.Run("manager="+k, func(t *testing.T) {
			_, err := m.Get(p.ID)
			require.Error(t, err)
			require.NoError(t, m.Create(p))

			res, err := m.Get(p.ID)
			require.NoError(t, err)
			assert.Equal(t, p, res)

			p.Subjects = []string{"stan"}
			require.NoError(t, m.Update(p))

			pols, err := m.List(10, 0)
			require.NoError(t, err)
			assert.Len(t, pols, 1)

			res, err = m.Get(p.ID)
			require.NoError(t, err)
			assert.Equal(t, p, res)

			require.NoError(t, m.Delete(p.ID))

			_, err = m.Get(p.ID)
			assert.Error(t, err)
		})
	}
}
