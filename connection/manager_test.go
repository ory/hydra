package connection

import (
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/internal"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/ladon"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"net/url"
)

var connections = []*Connection{
	&Connection{
		ID:            uuid.New(),
		LocalSubject:  "peter",
		RemoteSubject: "peterson",
		Provider:      "google",
	},
}

var managers = map[string]Manager{
	"memory": NewMemoryManager(),
}

var ts *httptest.Server

func init() {
	localWarden, httpClient := internal.NewFirewall("hydra", "alice", fosite.Arguments{scope},
		&ladon.DefaultPolicy{
			ID:        "1",
			Subjects:  []string{"alice"},
			Resources: []string{"rn:hydra:connections<.*>"},
			Actions:   []string{"create", "get", "delete", "find"},
			Effect:    ladon.AllowAccess,
		},
	)

	s := &Handler{
		Manager: &MemoryManager{Connections: map[string]*Connection{}},
		H:       &herodot.JSON{},
		W:       localWarden,
	}

	r := httprouter.New()
	s.SetRoutes(r)
	ts = httptest.NewServer(r)

	u, _ := url.Parse(ts.URL + "/connections")
	managers["http"] = &HTTPManager{
		Client:   httpClient,
		Endpoint: u,
	}
}

func TestCreateGetFindDelete(t *testing.T) {
	for _, store := range managers {
		for _, c := range connections {
			_, err := store.Get("asdf")
			pkg.RequireError(t, true, err)

			err = store.Create(c)
			pkg.RequireError(t, false, err)

			res, err := store.Get(c.GetID())
			pkg.RequireError(t, false, err)
			require.Equal(t, c, res)

			cs, err := store.FindAllByLocalSubject("peter")
			pkg.RequireError(t, false, err)
			assert.Len(t, cs, 1)
			require.Equal(t, c, cs[0])

			res, err = store.FindByRemoteSubject("google", "peterson")
			pkg.RequireError(t, false, err)
			require.Equal(t, c, res)

			err = store.Delete(c.GetID())
			pkg.RequireError(t, false, err)

			_, err = store.Get(c.GetID())
			pkg.RequireError(t, true, err)
		}
	}
}
