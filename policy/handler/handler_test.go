package handler

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/ory-am/dockertest"
	hcon "github.com/ory-am/hydra/context"
	hjwt "github.com/ory-am/hydra/jwt"
	"github.com/ory-am/hydra/middleware"
	"github.com/ory-am/ladon/guard"
	"github.com/ory-am/ladon/guard/operator"
	"github.com/ory-am/ladon/policy"
	"github.com/ory-am/ladon/policy/postgres"
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
	store *postgres.Store
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
	createData policy.DefaultPolicy
	//allowedPayload payload

	statusCreate int
	statusGet    int
	//expectsAllowed       bool
	statusDelete         int
	statusGetAfterDelete int
}

func mockAuthorization(c test) func(h hcon.ContextHandler) hcon.ContextHandler {
	return func(h hcon.ContextHandler) hcon.ContextHandler {
		return hcon.ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			claims := hjwt.NewClaimsCarrier(uuid.New(), "hydra", c.subject, "tests", time.Now(), time.Now())
			ctx = hcon.NewContextFromAuthValues(ctx, claims, &c.token, c.policies)
			h.ServeHTTPContext(ctx, rw, req)
		})
	}
}

var policies = map[string]policy.Policy{
	"pass-all":    &policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"rn:hydra:policies<.*>"}, []string{"<.*>"}, nil},
	"pass-create": &policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"rn:hydra:policies"}, []string{"create"}, nil},
	"pass-get":    &policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"rn:hydra:policies:<.*>"}, []string{"get"}, nil},
	"pass-delete": &policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"rn:hydra:policies:<.*>"}, []string{"delete"}, nil},
	"fail":        &policy.DefaultPolicy{},
}

var payloads = []policy.DefaultPolicy{
	policy.DefaultPolicy{
		"",
		"description",
		[]string{"max"},
		policy.AllowAccess,
		[]string{"resource"},
		[]string{"<.*>"},
		nil,
	},
	policy.DefaultPolicy{
		"",
		"Should allow peter all permissions on resource article",
		[]string{"peter"},
		policy.AllowAccess,
		[]string{"article"},
		[]string{"<.*>"},
		nil,
	},
}

func TestGrantedEndpoint(t *testing.T) {
	c := test{
		subject:  "peter",
		token:    jwt.Token{Valid: true},
		policies: []policy.Policy{policies["pass-all"]},
	}

	handler := &Handler{s: store, m: mw, g: &guard.Guard{}}
	router := mux.NewRouter()
	handler.SetRoutes(router, mockAuthorization(c))
	ts := httptest.NewServer(router)
	defer ts.Close()

	request := gorequest.New()
	resp, _, _ := request.Post(ts.URL + "/policies").Send(payloads[1]).End()
	require.Equal(t, 200, resp.StatusCode)

	do := func(p payload, shouldAllow bool) {
		resp, body, _ := request.Post(ts.URL + "/granted").Send(p).End()
		require.Equal(t, 200, resp.StatusCode)

		var isAllowed struct {
			Allowed bool `json:"allowed"`
		}
		json.Unmarshal([]byte(body), &isAllowed)
		require.Equal(t, 200, resp.StatusCode)
		require.Equal(t, shouldAllow, isAllowed.Allowed)
	}

	do(payload{
		Resource:   "article",
		Subject:    "peter",
		Permission: "random",
		Context: &operator.Context{
			Owner: "peter",
		},
	}, true)

	do(payload{
		Resource:   "article",
		Subject:    "peter",
		Permission: "foobar",
		Context: &operator.Context{
			Owner: "peter",
		},
	}, true)

	do(payload{
		Resource:   "foobar",
		Subject:    "peter",
		Permission: "random",
		Context: &operator.Context{
			Owner: "peter",
		},
	}, false)

	do(payload{
		Resource:   "article",
		Subject:    "max",
		Permission: "random",
		Context: &operator.Context{
			Owner: "peter",
		},
	}, false)
}

func TestCreateGetDeleteGet(t *testing.T) {
	for k, c := range []test{
		{
			subject: "peter", token: jwt.Token{Valid: false},
			policies:     []policy.Policy{policies["fail"]},
			createData:   payloads[0],
			statusCreate: http.StatusUnauthorized,
		},
		{
			subject: "peter", token: jwt.Token{Valid: true},
			policies:     []policy.Policy{policies["fail"]},
			createData:   payloads[0],
			statusCreate: http.StatusForbidden,
		},
		{
			subject: "peter", token: jwt.Token{Valid: true},
			policies:     []policy.Policy{policies["pass-create"]},
			createData:   payloads[0],
			statusCreate: http.StatusOK, statusGet: http.StatusForbidden,
		},
		{
			subject: "peter", token: jwt.Token{Valid: true},
			policies:     []policy.Policy{policies["pass-create"], policies["pass-get"]},
			createData:   payloads[0],
			statusCreate: http.StatusOK, statusGet: http.StatusOK, statusDelete: http.StatusForbidden,
		},
		{
			subject: "peter", token: jwt.Token{Valid: true},
			policies:     []policy.Policy{policies["pass-all"]},
			createData:   payloads[0],
			statusCreate: http.StatusOK, statusGet: http.StatusOK, statusDelete: http.StatusAccepted, statusGetAfterDelete: http.StatusNotFound,
		},
	} {
		handler := &Handler{s: store, m: mw}
		router := mux.NewRouter()
		handler.SetRoutes(router, mockAuthorization(c))
		ts := httptest.NewServer(router)
		defer ts.Close()

		request := gorequest.New()
		resp, body, _ := request.Post(ts.URL + "/policies").Send(c.createData).End()
		require.Equal(t, c.statusCreate, resp.StatusCode, "case %d: %s", k, body)
		if resp.StatusCode != http.StatusOK {
			continue
		}

		var pol policy.DefaultPolicy
		json.Unmarshal([]byte(body), &pol)

		resp, body, _ = request.Get(ts.URL + "/policies/" + pol.ID).End()
		require.Equal(t, c.statusGet, resp.StatusCode, "case %d: %s", k, body)
		if resp.StatusCode != http.StatusOK {
			continue
		}

		resp, body, _ = request.Delete(ts.URL + "/policies/" + pol.ID).End()
		require.Equal(t, c.statusDelete, resp.StatusCode, "case %d: %s", k, body)
		if resp.StatusCode != http.StatusAccepted {
			continue
		}

		resp, body, _ = request.Get(ts.URL + "/policies/" + pol.ID).End()
		require.Equal(t, c.statusGetAfterDelete, resp.StatusCode, "case %d: %s", k, body)
	}
}

var allowedPayloads = map[string]payload{
	"create-grant": payload{
		Resource:   "rn:hydra:policies",
		Subject:    "peter",
		Permission: "create",
	},
}
