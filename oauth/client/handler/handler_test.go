package handler

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	chd "github.com/ory-am/common/handler"
	"github.com/ory-am/dockertest"
	authcon "github.com/ory-am/hydra/context"
	hjwt "github.com/ory-am/hydra/jwt"
	middleware "github.com/ory-am/hydra/middleware/host"
	"github.com/ory-am/ladon/policy"
	"github.com/ory-am/osin-storage/storage/postgres"
	"github.com/parnurzeal/gorequest"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var (
	mw    *middleware.Middleware
	store *postgres.Storage
)

func TestMain(m *testing.M) {
	c, db, err := dockertest.OpenPostgreSQLContainerConnection(15, time.Second)
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	defer c.KillRemove()

	store = postgres.New(db)
	mw = &middleware.Middleware{}
	if err := store.CreateSchemas(); err != nil {
		log.Fatalf("COuld not create schemas: %s", err)
	}

	os.Exit(m.Run())
}

type test struct {
	subject    string
	token      jwt.Token
	policies   []policy.Policy
	createData payload

	statusGet            int
	statusCreate         int
	statusDelete         int
	statusGetAfterDelete int
}

func mockAuthorization(c test) func(h chd.ContextHandler) chd.ContextHandler {
	return func(h chd.ContextHandler) chd.ContextHandler {
		return chd.ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			claims := hjwt.NewClaimsCarrier(uuid.New(), "hydra", c.subject, "tests", time.Now().Add(time.Hour), time.Now(), time.Now())
			ctx = authcon.NewContextFromAuthValues(ctx, claims, &c.token, c.policies)
			h.ServeHTTPContext(ctx, rw, req)
		})
	}
}

var policies = map[string]policy.Policy{
	"pass-all":    &policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"rn:hydra:clients<.*>"}, []string{"<.*>"}, nil},
	"pass-create": &policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"rn:hydra:clients"}, []string{"create"}, nil},
	"pass-get":    &policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"rn:hydra:clients:<.*>"}, []string{"get"}, nil},
	"pass-delete": &policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"rn:hydra:clients:<.*>"}, []string{"delete"}, nil},
	"fail":        &policy.DefaultPolicy{},
}

func TestCreateGetDeleteGet(t *testing.T) {
	for k, c := range []test{
		{subject: "peter", token: jwt.Token{Valid: false}, policies: []policy.Policy{policies["fail"]}, createData: payload{RedirectURIs: "redir"}, statusCreate: http.StatusUnauthorized},
		{subject: "peter", token: jwt.Token{Valid: true}, policies: []policy.Policy{policies["fail"]}, createData: payload{RedirectURIs: "redir"}, statusCreate: http.StatusForbidden},
		{subject: "peter", token: jwt.Token{Valid: true}, policies: []policy.Policy{policies["pass-create"]}, createData: payload{RedirectURIs: "redir"}, statusCreate: http.StatusOK, statusGet: http.StatusForbidden},
		{subject: "peter", token: jwt.Token{Valid: true}, policies: []policy.Policy{policies["pass-create"], policies["pass-get"]}, createData: payload{RedirectURIs: "redir"}, statusCreate: http.StatusOK, statusGet: http.StatusOK, statusDelete: http.StatusForbidden},
		{subject: "peter", token: jwt.Token{Valid: true}, policies: []policy.Policy{policies["pass-all"]}, createData: payload{RedirectURIs: "redir"}, statusCreate: http.StatusOK, statusGet: http.StatusOK, statusDelete: http.StatusAccepted, statusGetAfterDelete: http.StatusNotFound},
	} {
		handler := &Handler{s: store, m: mw}
		router := mux.NewRouter()
		handler.SetRoutes(router, mockAuthorization(c))
		ts := httptest.NewServer(router)
		defer ts.Close()

		request := gorequest.New()
		resp, body, _ := request.Post(ts.URL + "/clients").Send(c.createData).End()
		require.Equal(t, c.statusCreate, resp.StatusCode, "case %d: %s", k, body)
		if resp.StatusCode != http.StatusOK {
			continue
		}

		var client payload
		json.Unmarshal([]byte(body), &client)

		resp, body, _ = request.Get(ts.URL + "/clients/" + client.ID).End()
		require.Equal(t, c.statusGet, resp.StatusCode, "case %d: %s", k, body)
		if resp.StatusCode != http.StatusOK {
			continue
		}

		resp, body, _ = request.Delete(ts.URL + "/clients/" + client.ID).End()
		require.Equal(t, c.statusDelete, resp.StatusCode, "case %d: %s", k, body)
		if resp.StatusCode != http.StatusAccepted {
			continue
		}

		resp, body, _ = request.Get(ts.URL + "/clients/" + client.ID).End()
		require.Equal(t, c.statusGetAfterDelete, resp.StatusCode, "case %d: %s", k, body)
	}
}
