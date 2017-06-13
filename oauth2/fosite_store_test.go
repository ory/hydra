package oauth2

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/integration"
	"github.com/ory/hydra/pkg"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var clientManagers = map[string]pkg.FositeStorer{}
var clientManager = &client.MemoryManager{
	Clients: map[string]client.Client{"foobar": {ID: "foobar"}},
	Hasher:  &fosite.BCrypt{},
}

func init() {
	clientManagers["memory"] = &FositeMemoryStore{
		AuthorizeCodes: make(map[string]fosite.Requester),
		IDSessions:     make(map[string]fosite.Requester),
		AccessTokens:   make(map[string]fosite.Requester),
		RefreshTokens:  make(map[string]fosite.Requester),
	}
}

func TestMain(m *testing.M) {
	s := m.Run()
	integration.KillAll()
	os.Exit(s)
}

func connectToPG() {
	var db = integration.ConnectToPostgres()
	s := &FositeSQLStore{DB: db, Manager: clientManager, L: logrus.New()}
	if _, err := s.CreateSchemas(); err != nil {
		logrus.Fatalf("Could not create postgres schema: %v", err)
	}

	clientManagers["postgres"] = s
}

func connectToMySQL() {
	var db = integration.ConnectToMySQL()
	s := &FositeSQLStore{DB: db, Manager: clientManager, L: logrus.New()}
	if _, err := s.CreateSchemas(); err != nil {
		logrus.Fatalf("Could not create postgres schema: %v", err)
	}

	clientManagers["mysql"] = s
}

// This needs to be the first test!!
func TestConnectToStores(t *testing.T) {
	connectToPG()
	connectToMySQL()
}

func TestCreateGetDeleteAuthorizeCodes(t *testing.T) {
	for k, m := range clientManagers {
		t.Run(fmt.Sprintf("case=%s", k), TestHelperCreateGetDeleteAuthorizeCodes(m))
	}
}

func TestCreateGetDeleteAccessTokenSession(t *testing.T) {
	for k, m := range clientManagers {
		t.Run(fmt.Sprintf("case=%s", k), TestHelperCreateGetDeleteAccessTokenSession(m))
	}
}

func TestCreateGetDeleteOpenIDConnectSession(t *testing.T) {
	ctx := context.Background()
	for k, m := range clientManagers {
		t.Run(fmt.Sprintf("case=%s", k), func(t *testing.T) {
			_, err := m.GetOpenIDConnectSession(ctx, "4321", &fosite.Request{})
			assert.NotNil(t, err)

			err = m.CreateOpenIDConnectSession(ctx, "4321", &defaultRequest)
			require.Nil(t, err)

			res, err := m.GetOpenIDConnectSession(ctx, "4321", &fosite.Request{Session: &fosite.DefaultSession{}})
			require.Nil(t, err)
			AssertObjectKeysEqual(t, &defaultRequest, res, "Scopes", "GrantedScopes", "Form", "Session")

			err = m.DeleteOpenIDConnectSession(ctx, "4321")
			require.Nil(t, err)

			time.Sleep(100 * time.Millisecond)

			_, err = m.GetOpenIDConnectSession(ctx, "4321", &fosite.Request{})
			assert.NotNil(t, err)
		})
	}
}

func TestCreateGetDeleteRefreshTokenSession(t *testing.T) {
	ctx := context.Background()
	for k, m := range clientManagers {
		t.Run(fmt.Sprintf("case=%s", k), func(t *testing.T) {
			_, err := m.GetRefreshTokenSession(ctx, "4321", &fosite.DefaultSession{})
			assert.NotNil(t, err)

			err = m.CreateRefreshTokenSession(ctx, "4321", &defaultRequest)
			require.Nil(t, err)

			res, err := m.GetRefreshTokenSession(ctx, "4321", &fosite.DefaultSession{})
			require.Nil(t, err)
			AssertObjectKeysEqual(t, &defaultRequest, res, "Scopes", "GrantedScopes", "Form", "Session")

			err = m.DeleteRefreshTokenSession(ctx, "4321")
			require.Nil(t, err)

			time.Sleep(100 * time.Millisecond)

			_, err = m.GetRefreshTokenSession(ctx, "4321", &fosite.DefaultSession{})
			assert.NotNil(t, err)
		})
	}
}

func TestRevokeRefreshToken(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	for k, m := range clientManagers {
		t.Run(fmt.Sprintf("case=%s", k), func(t *testing.T) {
			_, err := m.GetRefreshTokenSession(ctx, "1111", &fosite.DefaultSession{})
			assert.NotNil(t, err)

			err = m.CreateRefreshTokenSession(ctx, "1111", &fosite.Request{ID: id, Client: &client.Client{ID: "foobar"}, RequestedAt: time.Now().Round(time.Second)})
			require.Nil(t, err)

			err = m.CreateRefreshTokenSession(ctx, "1122", &fosite.Request{ID: id, Client: &client.Client{ID: "foobar"}, RequestedAt: time.Now().Round(time.Second)})
			require.Nil(t, err)

			_, err = m.GetRefreshTokenSession(ctx, "1111", &fosite.DefaultSession{})
			require.Nil(t, err)

			err = m.RevokeRefreshToken(ctx, id)
			require.Nil(t, err)

			time.Sleep(100 * time.Millisecond)

			_, err = m.GetRefreshTokenSession(ctx, "1111", &fosite.DefaultSession{})
			assert.NotNil(t, err)

			_, err = m.GetRefreshTokenSession(ctx, "1122", &fosite.DefaultSession{})
			assert.NotNil(t, err)
		})
	}
}
