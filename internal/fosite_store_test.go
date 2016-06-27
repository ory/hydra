package internal

import (
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	c "github.com/ory-am/common/pkg"
	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/client"
	"github.com/ory-am/hydra/pkg"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
	r "gopkg.in/dancannon/gorethink.v2"
	"gopkg.in/ory-am/dockertest.v2"
)

var rethinkManager *FositeRehinkDBStore

var clientManagers = map[string]pkg.FositeStorer{}

func init() {
	clientManagers["memory"] = &FositeMemoryStore{
		AuthorizeCodes: make(map[string]fosite.Requester),
		IDSessions:     make(map[string]fosite.Requester),
		AccessTokens:   make(map[string]fosite.Requester),
		Implicit:       make(map[string]fosite.Requester),
		RefreshTokens:  make(map[string]fosite.Requester),
	}

}

func TestMain(m *testing.M) {
	var session *r.Session
	var err error

	c, err := dockertest.ConnectToRethinkDB(20, time.Millisecond*500, func(url string) bool {
		if session, err = r.Connect(r.ConnectOpts{Address: url, Database: "hydra"}); err != nil {
			return false
		} else if _, err = r.DBCreate("hydra").RunWrite(session); err != nil {
			logrus.Printf("Database exists: %s", err)
			return false
		} else if _, err = r.TableCreate("hydra_authorize_code").RunWrite(session); err != nil {
			logrus.Printf("Could not create table: %s", err)
			return false
		} else if _, err = r.TableCreate("hydra_id_sessions").RunWrite(session); err != nil {
			logrus.Printf("Could not create table: %s", err)
			return false
		} else if _, err = r.TableCreate("hydra_access_token").RunWrite(session); err != nil {
			logrus.Printf("Could not create table: %s", err)
			return false
		} else if _, err = r.TableCreate("hydra_implicit").RunWrite(session); err != nil {
			logrus.Printf("Could not create table: %s", err)
			return false
		} else if _, err = r.TableCreate("hydra_refresh_token").RunWrite(session); err != nil {
			logrus.Printf("Could not create table: %s", err)
			return false
		}

		rethinkManager = &FositeRehinkDBStore{
			Session:             session,
			AuthorizeCodesTable: r.Table("hydra_authorize_code"),
			IDSessionsTable:     r.Table("hydra_id_sessions"),
			AccessTokensTable:   r.Table("hydra_access_token"),
			ImplicitTable:       r.Table("hydra_implicit"),
			RefreshTokensTable:  r.Table("hydra_refresh_token"),
			AuthorizeCodes:      make(RDBItems),
			IDSessions:          make(RDBItems),
			AccessTokens:        make(RDBItems),
			Implicit:            make(RDBItems),
			RefreshTokens:       make(RDBItems),
		}
		rethinkManager.Watch(context.Background())
		time.Sleep(500 * time.Millisecond)
		return true
	})
	if session != nil {
		defer session.Close()
	}
	if err != nil {
		logrus.Fatalf("Could not connect to database: %s", err)
	}
	clientManagers["rethink"] = rethinkManager

	retCode := m.Run()
	c.KillRemove()
	os.Exit(retCode)
}

type testSession struct {
	Foo string `json:"foo" gorethink:"foo"`
}

var defaultRequest = fosite.Request{
	RequestedAt:   time.Now().Round(time.Second),
	Client:        &client.Client{ID: "foobar"},
	Scopes:        fosite.Arguments{"fa", "ba"},
	GrantedScopes: fosite.Arguments{"fa", "ba"},
	Form:          url.Values{"foo": []string{"bar", "baz"}},
	Session:       &testSession{Foo: "bar"},
}

func TestColdStartRethinkManager(t *testing.T) {
	ctx := context.Background()
	m := rethinkManager
	id := uuid.New()

	err := m.CreateAuthorizeCodeSession(ctx, id, &defaultRequest)
	pkg.AssertError(t, false, err)
	err = m.CreateAccessTokenSession(ctx, id, &defaultRequest)
	pkg.AssertError(t, false, err)

	time.Sleep(100 * time.Millisecond)

	_, err = m.GetAuthorizeCodeSession(ctx, id, &testSession{})
	pkg.AssertError(t, false, err)
	_, err = m.GetAccessTokenSession(ctx, id, &testSession{})
	pkg.AssertError(t, false, err)

	delete(rethinkManager.AuthorizeCodes, id)
	delete(rethinkManager.AccessTokens, id)

	_, err = m.GetAuthorizeCodeSession(ctx, id, &testSession{})
	pkg.AssertError(t, true, err)
	_, err = m.GetAccessTokenSession(ctx, id, &testSession{})
	pkg.AssertError(t, true, err)

	err = rethinkManager.ColdStart()
	pkg.AssertError(t, false, err)

	_, err = m.GetAuthorizeCodeSession(ctx, id, &testSession{})
	pkg.AssertError(t, false, err)
	_, err = m.GetAccessTokenSession(ctx, id, &testSession{})
	pkg.AssertError(t, false, err)
}

func TestCreateGetDeleteAuthorizeCodes(t *testing.T) {
	ctx := context.Background()
	for k, m := range clientManagers {
		_, err := m.GetAuthorizeCodeSession(ctx, "4321", &testSession{})
		pkg.AssertError(t, true, err, "%s", k)

		err = m.CreateAuthorizeCodeSession(ctx, "4321", &defaultRequest)
		pkg.AssertError(t, false, err, "%s", k)

		time.Sleep(100 * time.Millisecond)

		res, err := m.GetAuthorizeCodeSession(ctx, "4321", &testSession{})
		pkg.RequireError(t, false, err, "%s", k)
		c.AssertObjectKeysEqual(t, &defaultRequest, res, "Scopes", "GrantedScopes", "Form", "Session")

		err = m.DeleteAuthorizeCodeSession(ctx, "4321")
		pkg.AssertError(t, false, err, "%s", k)

		time.Sleep(100 * time.Millisecond)

		_, err = m.GetAuthorizeCodeSession(ctx, "4321", &testSession{})
		pkg.AssertError(t, true, err, "%s", k)
	}
}

func TestCreateGetDeleteAccessTokenSession(t *testing.T) {
	ctx := context.Background()
	for k, m := range clientManagers {
		_, err := m.GetAccessTokenSession(ctx, "4321", &testSession{})
		pkg.AssertError(t, true, err, "%s", k)

		err = m.CreateAccessTokenSession(ctx, "4321", &defaultRequest)
		pkg.AssertError(t, false, err, "%s", k)

		time.Sleep(100 * time.Millisecond)

		res, err := m.GetAccessTokenSession(ctx, "4321", &testSession{})
		pkg.RequireError(t, false, err, "%s", k)
		c.AssertObjectKeysEqual(t, &defaultRequest, res, "Scopes", "GrantedScopes", "Form", "Session")

		err = m.DeleteAccessTokenSession(ctx, "4321")
		pkg.AssertError(t, false, err, "%s", k)

		time.Sleep(100 * time.Millisecond)

		_, err = m.GetAccessTokenSession(ctx, "4321", &testSession{})
		pkg.AssertError(t, true, err, "%s", k)
	}
}

func TestCreateGetDeleteOpenIDConnectSession(t *testing.T) {
	ctx := context.Background()
	for k, m := range clientManagers {
		_, err := m.GetOpenIDConnectSession(ctx, "4321", &fosite.Request{})
		pkg.AssertError(t, true, err, "%s", k)

		err = m.CreateOpenIDConnectSession(ctx, "4321", &defaultRequest)
		pkg.AssertError(t, false, err, "%s", k)

		time.Sleep(100 * time.Millisecond)

		res, err := m.GetOpenIDConnectSession(ctx, "4321", &fosite.Request{
			Session: &testSession{},
		})
		pkg.RequireError(t, false, err, "%s", k)
		c.AssertObjectKeysEqual(t, &defaultRequest, res, "Scopes", "GrantedScopes", "Form", "Session")

		err = m.DeleteOpenIDConnectSession(ctx, "4321")
		pkg.AssertError(t, false, err, "%s", k)

		time.Sleep(100 * time.Millisecond)

		_, err = m.GetOpenIDConnectSession(ctx, "4321", &fosite.Request{})
		pkg.AssertError(t, true, err, "%s", k)
	}
}

func TestCreateGetDeleteRefreshTokenSession(t *testing.T) {
	ctx := context.Background()
	for k, m := range clientManagers {
		_, err := m.GetRefreshTokenSession(ctx, "4321", &testSession{})
		pkg.AssertError(t, true, err, "%s", k)

		err = m.CreateRefreshTokenSession(ctx, "4321", &defaultRequest)
		pkg.AssertError(t, false, err, "%s", k)

		time.Sleep(100 * time.Millisecond)

		res, err := m.GetRefreshTokenSession(ctx, "4321", &testSession{})
		pkg.RequireError(t, false, err, "%s", k)
		c.AssertObjectKeysEqual(t, &defaultRequest, res, "Scopes", "GrantedScopes", "Form", "Session")

		err = m.DeleteRefreshTokenSession(ctx, "4321")
		pkg.AssertError(t, false, err, "%s", k)

		time.Sleep(100 * time.Millisecond)

		_, err = m.GetRefreshTokenSession(ctx, "4321", &testSession{})
		pkg.AssertError(t, true, err, "%s", k)
	}
}
