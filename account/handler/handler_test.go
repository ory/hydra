package handler

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	chd "github.com/ory-am/common/handler"
	"github.com/ory-am/dockertest"
	"github.com/ory-am/hydra/account"
	hydra "github.com/ory-am/hydra/account/postgres"
	middleware "github.com/ory-am/hydra/middleware/host"
	authcon "github.com/ory-am/hydra/context"
	"github.com/ory-am/hydra/hash"
	hjwt "github.com/ory-am/hydra/jwt"
	"github.com/ory-am/ladon/policy"
	"github.com/parnurzeal/gorequest"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
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
	hd *Handler
)

func TestMain(m *testing.M) {
	c, db, err := dockertest.OpenPostgreSQLContainerConnection(15, 5*time.Second)
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	defer c.KillRemove()

	s := hydra.New(&hash.BCrypt{10}, db)
	if err := s.CreateSchemas(); err != nil {
		log.Fatalf("Could not set up schemas: %v", err)
	}

	mw := &middleware.Middleware{}
	hd = &Handler{s, mw}
	os.Exit(m.Run())
}

type result struct {
	create int
	get    int
	delete int
}

type payload struct {
	ID       string `json:"id,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Data     string `json:"data,omitempty"`
}

type test struct {
	subject  string
	token    *jwt.Token
	policies []policy.Policy
	payload  payload
	expected result
}

func mock(c test) func(h chd.ContextHandler) chd.ContextHandler {
	return func(h chd.ContextHandler) chd.ContextHandler {
		return chd.ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			claims := hjwt.NewClaimsCarrier(uuid.New(), "hydra", c.subject, "tests", time.Now(), time.Now())
			ctx = authcon.NewContextFromAuthValues(ctx, claims, c.token, c.policies)
			h.ServeHTTPContext(ctx, rw, req)
		})
	}
}

func assertAccount(t *testing.T, c test, data string) account.Account {
	var acc account.DefaultAccount
	require.Nil(t, json.Unmarshal([]byte(data), &acc))
	assert.Equal(t, c.payload.Email, acc.Email)
	assert.Equal(t, c.payload.Data, acc.Data)
	assert.Empty(t, acc.Password)
	return &acc
}

var policies = map[string][]policy.Policy{
	"allow-create": {&policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"rn:hydra:accounts"}, []string{"create"}, nil}},
	"allow-create-get": {
		&policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"rn:hydra:accounts"}, []string{"create"}, nil},
		&policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"<.*>"}, []string{"get"}, nil},
	},
	"allow-all": {
		&policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"rn:hydra:accounts"}, []string{"create"}, nil},
		&policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"rn:hydra:accounts:<.*>"}, []string{"get"}, nil},
		&policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"rn:hydra:accounts:<.*>"}, []string{"delete"}, nil},
	},
	"empty": {},
}

func TestCreateGetDelete(t *testing.T) {
	for k, c := range []test{
		{
			subject: "peter", token: &jwt.Token{Valid: false},
			expected: result{create: http.StatusUnauthorized, get: 0, delete: 0},
		},
		{
			subject: "peter", token: &jwt.Token{Valid: true}, policies: policies["empty"],
			expected: result{create: http.StatusForbidden, get: 0, delete: 0},
		},
		{
			subject: "peter", token: &jwt.Token{Valid: true}, policies: policies["empty"],
			expected: result{create: http.StatusForbidden, get: 0, delete: 0},
		},
		{
			subject: "max", token: &jwt.Token{Valid: true}, policies: policies["empty"],
			expected: result{create: http.StatusForbidden, get: 0, delete: 0},
		},
		{
			subject: "max", token: &jwt.Token{Valid: true}, payload: payload{},
			policies: policies["allow-create"],
			expected: result{
				create: http.StatusForbidden, get: 0, delete: 0,
			},
		},
		{
			subject: "peter", token: &jwt.Token{Valid: true},
			payload:  payload{Email: uuid.New() + "@foobar.com", Data: "{}"},
			policies: policies["allow-create"],
			expected: result{
				create: http.StatusBadRequest, get: 0, delete: 0,
			},
		},
		{
			subject: "peter", token: &jwt.Token{Valid: true},
			payload:  payload{Email: uuid.New() + "@foobar.com", Password: "123", Data: "{}"},
			policies: policies["allow-create"],
			expected: result{
				create: http.StatusBadRequest, get: 0, delete: 0,
			},
		},
		{
			subject: "peter", token: &jwt.Token{Valid: true},
			payload:  payload{Email: "notemail", Password: "secret", Data: "{}"},
			policies: policies["allow-create"],
			expected: result{
				create: http.StatusBadRequest, get: 0, delete: 0,
			},
		},
		{
			subject: "peter", token: &jwt.Token{Valid: true},
			payload:  payload{Email: uuid.New() + "@bar.com", Password: "", Data: "{}"},
			policies: policies["allow-create"],
			expected: result{
				create: http.StatusBadRequest, get: 0, delete: 0,
			},
		},
		{
			subject: "peter", token: &jwt.Token{Valid: true},
			payload:  payload{Email: uuid.New() + "@bar.com", Password: "secret", Data: "not json"},
			policies: policies["allow-create"],
			expected: result{
				create: http.StatusBadRequest, get: 0, delete: 0,
			},
		},
		{
			subject: "peter", token: &jwt.Token{Valid: true},
			payload:  payload{Email: uuid.New() + "@bar.com", Password: "secret", Data: "{}"},
			policies: policies["allow-create"],
			expected: result{
				create: http.StatusOK, get: http.StatusForbidden, delete: http.StatusForbidden,
			},
		},
		{
			subject: "peter", token: &jwt.Token{Valid: true},
			payload:  payload{Email: uuid.New() + "@bar.com", Password: "secret", Data: "{}"},
			policies: policies["allow-create-get"],
			expected: result{
				create: http.StatusOK, get: http.StatusOK, delete: http.StatusForbidden,
			},
		},
		{
			subject: "peter", token: &jwt.Token{Valid: true},
			payload:  payload{Email: uuid.New() + "@bar.com", Password: "secret", Data: "{}"},
			policies: policies["allow-all"],
			expected: result{
				create: http.StatusOK, get: http.StatusOK, delete: http.StatusAccepted,
			},
		},
	} {
		router := mux.NewRouter()
		hd.SetRoutes(router, mock(c))
		ts := httptest.NewServer(router)
		defer ts.Close()

		t.Logf(ts.URL + "/accounts")

		request := gorequest.New()
		resp, body, _ := request.Post(ts.URL + "/accounts").Send(c.payload).End()
		require.Equal(t, c.expected.create, resp.StatusCode, "case %d: %s", k, body)
		if resp.StatusCode != http.StatusOK {
			continue
		}
		user := assertAccount(t, c, body)

		resp, body, _ = request.Get(ts.URL + "/accounts/" + user.GetID()).End()
		require.Equal(t, c.expected.get, resp.StatusCode, "case %d: %s", k, body)
		if resp.StatusCode != http.StatusOK {
			continue
		}
		user = assertAccount(t, c, body)

		resp, body, _ = request.Delete(ts.URL + "/accounts/" + user.GetID()).End()
		require.Equal(t, c.expected.delete, resp.StatusCode, "case %d: %s", k, body)
		if resp.StatusCode != http.StatusAccepted {
			continue
		}

		resp, body, _ = request.Get(ts.URL + "/accounts/" + user.GetID()).End()
		require.Equal(t, http.StatusNotFound, resp.StatusCode, "case %d: %s", k, body)
	}
}
