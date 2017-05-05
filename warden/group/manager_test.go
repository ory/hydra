package group_test

import (
	"log"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"fmt"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/ory/fosite"
	"github.com/ory/herodot"
	"github.com/ory/hydra/compose"
	"github.com/ory/hydra/integration"
	. "github.com/ory/hydra/warden/group"
	"github.com/ory/ladon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var clientManagers = map[string]Manager{}
var ts *httptest.Server

func init() {
	clientManagers["memory"] = &MemoryManager{
		Groups: map[string]Group{},
	}

	localWarden, httpClient := compose.NewMockFirewall("foo", "alice", fosite.Arguments{Scope}, &ladon.DefaultPolicy{
		ID:        "1",
		Subjects:  []string{"alice"},
		Resources: []string{"rn:hydra:warden<.*>"},
		Actions:   []string{"create", "get", "delete", "update", "members.add", "members.remove"},
		Effect:    ladon.AllowAccess,
	})

	s := &Handler{
		Manager: &MemoryManager{
			Groups: map[string]Group{},
		},
		H: herodot.NewJSONWriter(nil),
		W: localWarden,
	}

	routing := httprouter.New()
	s.SetRoutes(routing)
	ts = httptest.NewServer(routing)

	u, _ := url.Parse(ts.URL + GroupsHandlerPath)
	clientManagers["http"] = &HTTPManager{
		Client:   httpClient,
		Endpoint: u,
	}
}

func TestMain(m *testing.M) {
	connectToPG()
	connectToMySQL()

	s := m.Run()
	integration.KillAll()
	os.Exit(s)
}

func connectToMySQL() {
	var db = integration.ConnectToMySQL()
	s := &SQLManager{DB: db}
	if _, err := s.CreateSchemas(); err != nil {
		log.Fatalf("Could not create postgres schema: %v", err)
	}

	clientManagers["mysql"] = s
}

func connectToPG() {
	var db = integration.ConnectToPostgres()
	s := &SQLManager{DB: db}

	if _, err := s.CreateSchemas(); err != nil {
		log.Fatalf("Could not create postgres schema: %v", err)
	}

	clientManagers["postgres"] = s
}

func TestManagers(t *testing.T) {
	for k, m := range clientManagers {
		t.Run(fmt.Sprintf("case=%s", k), func(t *testing.T) {
			_, err := m.GetGroup("4321")
			assert.NotNil(t, err)

			c := &Group{
				ID:      "1",
				Members: []string{"bar", "foo"},
			}
			assert.Nil(t, m.CreateGroup(c))
			assert.Nil(t, m.CreateGroup(&Group{
				ID:      "2",
				Members: []string{"foo"},
			}))

			d, err := m.GetGroup("1")
			require.Nil(t, err)
			assert.EqualValues(t, c.Members, d.Members)
			assert.EqualValues(t, c.ID, d.ID)

			ds, err := m.FindGroupNames("foo")
			require.Nil(t, err)
			assert.Len(t, ds, 2)

			assert.Nil(t, m.AddGroupMembers("1", []string{"baz"}))

			ds, err = m.FindGroupNames("baz")
			require.Nil(t, err)
			assert.Len(t, ds, 1)

			assert.Nil(t, m.RemoveGroupMembers("1", []string{"baz"}))
			ds, err = m.FindGroupNames("baz")
			require.Nil(t, err)
			assert.Len(t, ds, 0)

			assert.Nil(t, m.DeleteGroup("1"))
			_, err = m.GetGroup("1")
			require.NotNil(t, err)
		})
	}
}
