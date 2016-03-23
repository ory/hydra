package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	chd "github.com/ory-am/common/handler"
	authcon "github.com/ory-am/hydra/context"
	hjwt "github.com/ory-am/hydra/jwt"
	middleware "github.com/ory-am/hydra/middleware/host"
	. "github.com/ory-am/hydra/endpoint/connection"
	"github.com/ory-am/hydra/endpoint/connection/postgres"
	"github.com/ory-am/ladon/policy"
	"github.com/parnurzeal/gorequest"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"gopkg.in/ory-am/dockertest.v2"
)

var (
	mw    *middleware.Middleware
	store *postgres.Store
	db    *sql.DB
)

func TestMain(m *testing.M) {
	c, err := dockertest.ConnectToPostgreSQL(15, time.Second, func(url string) bool {
		var err error
		db, err = sql.Open("postgres", url)
		if err != nil {
			return false
		}
		return db.Ping() == nil
	})

	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	store = postgres.New(db)
	mw = &middleware.Middleware{}
	if err := store.CreateSchemas(); err != nil {
		log.Fatalf("Could not create schemas: %s", err)
	}

	retCode := m.Run()

	// force teardown
	tearDown(c)

	os.Exit(retCode)
}

func tearDown(c dockertest.ContainerID) {
	db.Close()
	c.KillRemove()
}

type test struct {
	subject    string
	token      jwt.Token
	policies   []policy.Policy
	createData *DefaultConnection

	statusCreate         int
	statusGet            int
	statusFind           int
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
	"pass-all":    &policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"rn:hydra:oauth2:connections<.*>"}, []string{"<.*>"}, nil},
	"pass-create": &policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"rn:hydra:oauth2:connections<.*>"}, []string{"create"}, nil},
	"pass-get":    &policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"rn:hydra:oauth2:connections<.*>"}, []string{"get"}, nil},
	"pass-delete": &policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"rn:hydra:oauth2:connections<.*>"}, []string{"delete"}, nil},
	"fail":        &policy.DefaultPolicy{},
}

var data = map[string]*DefaultConnection{
	"ok": {
		Provider:      "google",
		RemoteSubject: "peter@gmail.com",
		LocalSubject:  "peter",
	},
	"ok-max": {
		Provider:      "google",
		RemoteSubject: "max@gmail.com",
		LocalSubject:  "max",
	},
	"ok-zac": {
		Provider:      "google",
		RemoteSubject: "zac@gmail.com",
		LocalSubject:  "zac",
	},
	"ok-steve": {
		Provider:      "google",
		RemoteSubject: "steve@gmail.com",
		LocalSubject:  "steve",
	},
	"fail": {},
	"fail-validation": {
		Provider:      "",
		RemoteSubject: "",
		LocalSubject:  "",
	},
}

func TestCreateGetDeleteGet(t *testing.T) {
	for k, c := range []test{
		{subject: "peter", token: jwt.Token{Valid: false}, policies: []policy.Policy{policies["fail"]}, createData: data["fail"], statusCreate: http.StatusUnauthorized},
		{subject: "peter", token: jwt.Token{Valid: true}, policies: []policy.Policy{policies["fail"]}, createData: data["fail"], statusCreate: http.StatusForbidden},
		{subject: "peter", token: jwt.Token{Valid: true}, policies: []policy.Policy{policies["pass-create"]}, createData: data["ok-max"], statusCreate: http.StatusCreated, statusGet: http.StatusForbidden},
		{subject: "peter", token: jwt.Token{Valid: true}, policies: []policy.Policy{policies["pass-create"]}, createData: data["fail"], statusCreate: http.StatusBadRequest},
		{subject: "peter", token: jwt.Token{Valid: true}, policies: []policy.Policy{policies["pass-create"]}, createData: data["fail-validation"], statusCreate: http.StatusBadRequest},
		{subject: "peter", token: jwt.Token{Valid: true}, policies: []policy.Policy{policies["pass-create"], policies["pass-get"]}, createData: data["ok-zac"], statusCreate: http.StatusCreated, statusGet: http.StatusOK, statusDelete: http.StatusForbidden},
		{subject: "peter", token: jwt.Token{Valid: true}, policies: []policy.Policy{policies["pass-all"]}, createData: data["ok-steve"], statusCreate: http.StatusCreated, statusGet: http.StatusOK, statusDelete: http.StatusAccepted, statusGetAfterDelete: http.StatusNotFound},
	} {
		handler := &Handler{s: store, m: mw}
		router := mux.NewRouter()
		handler.SetRoutes(router, mockAuthorization(c))
		ts := httptest.NewServer(router)
		defer ts.Close()

		request := gorequest.New()
		connectionsURL := fmt.Sprintf("%s/oauth2/connections", ts.URL)

		resp, body, _ := request.Post(connectionsURL).Send(*c.createData).End()
		require.Equal(t, c.statusCreate, resp.StatusCode, "case %d: %s", k, body)
		if resp.StatusCode != http.StatusCreated {
			continue
		}

		var conn DefaultConnection
		assert.Nil(t, json.Unmarshal([]byte(body), &conn))

		resp, body, _ = request.Get(connectionsURL + "?subject=" + c.subject).End()
		require.Equal(t, c.statusGet, resp.StatusCode, "case %d: %s", k, body)
		if resp.StatusCode != http.StatusOK {
			continue
		}

		resp, body, _ = request.Get(fmt.Sprintf("%s/oauth2/connections/%s", ts.URL, conn.ID)).End()
		require.Equal(t, c.statusGet, resp.StatusCode, "case %d: %s", k, body)
		if resp.StatusCode != http.StatusOK {
			continue
		}

		resp, body, _ = request.Post(connectionsURL).Send(*c.createData).End()
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode, "case %d: %s", k, body)

		resp, body, _ = request.Delete(fmt.Sprintf("%s/oauth2/connections/%s", ts.URL, conn.ID)).End()
		require.Equal(t, c.statusDelete, resp.StatusCode, "case %d: %s", k, body)
		if resp.StatusCode != http.StatusAccepted {
			continue
		}

		resp, body, _ = request.Get(fmt.Sprintf("%s/oauth2/connections/%s", ts.URL, conn.ID)).End()
		require.Equal(t, c.statusGetAfterDelete, resp.StatusCode, "case %d: %s", k, body)

		resp, body, _ = request.Delete(fmt.Sprintf("%s/oauth2/connections/%s", ts.URL, conn.ID)).End()
		require.Equal(t, http.StatusNotFound, resp.StatusCode, "case %d: %s", k, body)
	}
}
