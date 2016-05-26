package policy

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/julienschmidt/httprouter"
	r "github.com/dancannon/gorethink"
	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/internal"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/ladon"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"time"
	"os"
	"gopkg.in/ory-am/dockertest.v2"
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	"github.com/stretchr/testify/require"
)

var managers = map[string]ladon.Manager{}

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
		} else if _, err = r.TableCreate("hydra_policies").RunWrite(session); err != nil {
			log.Printf("Could not create table: %s", err)
			return false
		}

		rethinkManager = &RethinkManager{
			Session: session,
			Table: r.Table("hydra_policies"),
			Policies: make(map[string]ladon.Policy),
		}

		if err := rethinkManager.Watch(context.Background()); err != nil {
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

func init() {
	localWarden, httpClient := internal.NewFirewall("hydra", "alice", fosite.Arguments{scope},
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
