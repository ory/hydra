package connection

import (
	"testing"

	"net/http/httptest"
	"net/url"

	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/internal"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/ladon"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	r "github.com/dancannon/gorethink"
	"time"
	"os"
	"gopkg.in/ory-am/dockertest.v2"
	"log"
	"golang.org/x/net/context"
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

var rethinkManager *RethinkManager

func TestMain(m *testing.M) {
	var session *r.Session
	var err error

	c, err := dockertest.ConnectToRethinkDB(20, time.Second, func(url string) bool {
		if session, err = r.Connect(r.ConnectOpts{Address:  url, Database: "hydra"}); err != nil {
			return false
		} else if _, err = r.DBCreate("hydra").RunWrite(session); err != nil {
			log.Printf("Database exists: %s", err)
			return false
		} else if _, err = r.TableCreate("hydra_clients").RunWrite(session); err != nil {
			log.Printf("Could not create table: %s", err)
			return false
		}

		rethinkManager = &RethinkManager{
			Session: session,
			Table: r.Table("hydra_clients"),
			Connections:make(map[string]*Connection),
		}
		err := rethinkManager.Watch(context.Background())
		if err != nil {
			log.Printf("Could not watch: %s", err)
			return false
		}
		return true
	})
	if session != nil {
		defer session.Close()
	}
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	managers["rethink"] = rethinkManager

	retCode := m.Run()
	c.KillRemove()
	os.Exit(retCode)
}

func BenchmarkRethinkGet(b *testing.B) {
	b.StopTimer()
	m := rethinkManager
	var err error
	err = m.Create(&Connection{
		ID:            "someid",
		LocalSubject:  "peter",
		RemoteSubject: "peterson",
		Provider:      "google",
	}, )
	if err != nil {
		b.Fatalf("%s", err)
	}
	time.Sleep(500 * time.Millisecond)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = m.Get("someid")
	}
}

func TestCreateGetFindDelete(t *testing.T) {
	for _, store := range managers {
		for _, c := range connections {
			_, err := store.Get("asdf")
			pkg.RequireError(t, true, err)

			err = store.Create(c)
			pkg.RequireError(t, false, err)

			time.Sleep(500 * time.Millisecond)

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

			time.Sleep(500 * time.Millisecond)

			_, err = store.Get(c.GetID())
			pkg.RequireError(t, true, err)
		}
	}
}
