package connection

import (
	"os"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/handler/core"
	"github.com/ory-am/hydra/herodot"
	ioa2 "github.com/ory-am/hydra/oauth2"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/hydra/warden"
	"github.com/ory-am/ladon"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"gopkg.in/ory-am/dockertest.v2"
	"net/http/httptest"
	"net/url"
	"time"
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

var containers = []dockertest.ContainerID{}

var fositeStore = pkg.FositeStore()

var ladonWarden = pkg.LadonWarden(map[string]ladon.Policy{
	"1": &ladon.DefaultPolicy{
		ID:        "1",
		Subjects:  []string{"alice"},
		Resources: []string{"rn:hydra:clients<.*>"},
		Actions:   []string{"create", "get", "delete"},
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

var ts *httptest.Server

var tokens = pkg.Tokens(1)

func init() {
	ar := fosite.NewAccessRequest(&ioa2.Session{Subject: "alice"})
	ar.GrantedScopes = fosite.Arguments{scope}
	fositeStore.CreateAccessTokenSession(nil, tokens[0][0], ar)

	managers["memory"] = &MemoryManager{Connections: map[string]*Connection{}}

	s := &Handler{
		Manager: &MemoryManager{Clients: map[string]*Connection{}},
		H:       &herodot.JSON{},
		W:       localWarden,
	}
	r := httprouter.New()
	s.SetRoutes(r)
	ts = httptest.NewServer(r)
	conf := &oauth2.Config{Scopes: []string{}, Endpoint: oauth2.Endpoint{}}

	u, _ := url.Parse(ts.URL + "/oauth2/connections/")
	managers["http"] = &HTTPManager{
		Client: conf.Client(oauth2.NoContext, &oauth2.Token{
			AccessToken: tokens[0][1],
			Expiry:      time.Now().Add(time.Hour),
			TokenType:   "bearer",
		}),
		Endpoint: u,
	}
}

func TestMain(m *testing.M) {
	retCode := m.Run()
	for _, c := range containers {
		c.KillRemove()
	}

	os.Exit(retCode)
}

func TestNotFound(t *testing.T) {
	for _, store := range managers {
		_, err := store.Get("asdf")
		assert.EqualError(t, err, pkg.ErrNotFound.Error())
	}
}

func TestCreateGetFindDelete(t *testing.T) {
	for _, store := range managers {
		for _, c := range connections {
			require.Nil(t, store.Create(c))

			res, err := store.Get(c.GetID())
			require.Nil(t, err)
			require.Equal(t, c, res)

			res, err = store.FindByRemoteSubject("google", "peterson")
			require.Nil(t, err)
			require.Equal(t, c, res)

			cs, err := store.FindAllByLocalSubject("peter")
			require.Nil(t, err)
			assert.Len(t, cs, 1)
			require.Equal(t, c, cs[0])

			require.Nil(t, store.Delete(c.GetID()))
			_, err = store.Get(c.GetID())
			require.NotNil(t, err)
		}
	}
}
