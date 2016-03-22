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
	"github.com/ory-am/hydra/account"
	hydra "github.com/ory-am/hydra/account/postgres"
	authcon "github.com/ory-am/hydra/context"
	"github.com/ory-am/hydra/hash"
	hjwt "github.com/ory-am/hydra/jwt"
	middleware "github.com/ory-am/hydra/middleware/host"
	"github.com/ory-am/ladon/policy"
	"github.com/parnurzeal/gorequest"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"gopkg.in/ory-am/dockertest.v2"
)

var (
	hd *Handler
)

var db *sql.DB

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

	s := hydra.New(&hash.BCrypt{10}, db)
	if err := s.CreateSchemas(); err != nil {
		log.Fatalf("Could not set up schemas: %v", err)
	}

	mw := &middleware.Middleware{}
	hd = &Handler{s, mw}

	retCode := m.Run()

	// force teardown
	tearDown(c)

	os.Exit(retCode)
}

func tearDown(c dockertest.ContainerID) {
	db.Close()
	c.KillRemove()
}

type result struct {
	create int
	get    int
	delete int
}

type payload struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
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
			claims := hjwt.NewClaimsCarrier(uuid.New(), "hydra", c.subject, "tests", time.Now().Add(time.Hour), time.Now(), time.Now())
			ctx = authcon.NewContextFromAuthValues(ctx, claims, c.token, c.policies)
			h.ServeHTTPContext(ctx, rw, req)
		})
	}
}

func assertAccount(t *testing.T, c test, data string) account.Account {
	var acc account.DefaultAccount
	require.Nil(t, json.Unmarshal([]byte(data), &acc))
	assert.Equal(t, c.payload.Username, acc.Username)
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
		&policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"rn:hydra:accounts:<.*>"}, []string{"<.*>"}, nil},
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
			payload:  payload{Username: uuid.New() + "@foobar.com", Data: "{}"},
			policies: policies["allow-create"],
			expected: result{
				create: http.StatusBadRequest, get: 0, delete: 0,
			},
		},
		{
			subject: "peter", token: &jwt.Token{Valid: true},
			payload:  payload{Username: uuid.New() + "@bar.com", Password: "", Data: "{}"},
			policies: policies["allow-create"],
			expected: result{
				create: http.StatusBadRequest, get: 0, delete: 0,
			},
		},
		{
			subject: "peter", token: &jwt.Token{Valid: true},
			payload:  payload{Username: uuid.New() + "@bar.com", Password: "secret", Data: "not json"},
			policies: policies["allow-create"],
			expected: result{
				create: http.StatusBadRequest, get: 0, delete: 0,
			},
		},
		{
			subject: "peter", token: &jwt.Token{Valid: true},
			payload:  payload{Username: uuid.New() + "@bar.com", Password: "secret", Data: "{}"},
			policies: policies["allow-create"],
			expected: result{
				create: http.StatusCreated, get: http.StatusForbidden, delete: http.StatusForbidden,
			},
		},
		{
			subject: "peter", token: &jwt.Token{Valid: true},
			payload:  payload{Username: uuid.New() + "@bar.com", Password: "secret", Data: "{}"},
			policies: policies["allow-create-get"],
			expected: result{
				create: http.StatusCreated, get: http.StatusOK, delete: http.StatusForbidden,
			},
		},
		{
			subject: "peter", token: &jwt.Token{Valid: true},
			payload:  payload{Username: uuid.New() + "@bar.com", Password: "secret", Data: "{}"},
			policies: policies["allow-all"],
			expected: result{
				create: http.StatusCreated, get: http.StatusOK, delete: http.StatusAccepted,
			},
		},
	} {
		router := mux.NewRouter()
		hd.SetRoutes(router, mock(c))
		ts := httptest.NewServer(router)
		defer ts.Close()

		t.Logf(ts.URL + "/accounts")

		request := gorequest.New()
		resp, body, errs := request.Post(ts.URL + "/accounts").Send(c.payload).End()
		require.Len(t, errs, 0, "%s", errs)
		require.Equal(t, c.expected.create, resp.StatusCode, "case %d: %s", k, body)
		if resp.StatusCode != http.StatusCreated {
			continue
		}
		user := assertAccount(t, c, body)

		resp, body, errs = request.Get(ts.URL + "/accounts/" + user.GetID()).End()
		require.Len(t, errs, 0, "%s", errs)
		require.Equal(t, c.expected.get, resp.StatusCode, "case %d: %s", k, body)
		if resp.StatusCode != http.StatusOK {
			continue
		}
		user = assertAccount(t, c, body)

		resp, body, errs = request.Delete(ts.URL + "/accounts/" + user.GetID()).End()
		require.Len(t, errs, 0, "%s", errs)
		require.Equal(t, c.expected.delete, resp.StatusCode, "case %d: %s", k, body)
		if resp.StatusCode != http.StatusAccepted {
			continue
		}

		resp, body, errs = request.Get(ts.URL + "/accounts/" + user.GetID()).End()
		require.Len(t, errs, 0, "%s", errs)
		require.Equal(t, http.StatusNotFound, resp.StatusCode, "case %d: %s", k, body)
	}
}

func setUpAccountAndServer(t *testing.T) (*httptest.Server, account.DefaultAccount) {
	router := mux.NewRouter()
	hd.SetRoutes(router, mock(test{
		subject:  "peter",
		token:    &jwt.Token{Valid: true},
		policies: policies["allow-all"],
	}))
	ts := httptest.NewServer(router)

	var user account.DefaultAccount
	request := gorequest.New()
	resp, body, errs := request.Post(ts.URL + "/accounts").Send(account.CreateAccountRequest{
		Username: uuid.New(),
		Password: "secret",
	}).End()
	require.Len(t, errs, 0, "%s", errs)
	require.Equal(t, http.StatusCreated, resp.StatusCode, "%s", body)
	require.Nil(t, json.Unmarshal([]byte(body), &user))
	return ts, user
}

func TestUpdatePassword(t *testing.T) {
	ts, user := setUpAccountAndServer(t)
	defer ts.Close()
	request := gorequest.New()

	resp, body, errs := request.Put(fmt.Sprintf("%s/accounts/%s/password", ts.URL, user.ID)).Send(account.UpdatePasswordRequest{
		CurrentPassword: "wrong",
		NewPassword:     "secret",
	}).End()
	require.Len(t, errs, 0, "%s", errs)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode, "%s", body)

	resp, body, errs = request.Put(fmt.Sprintf("%s/accounts/%s/password", ts.URL, user.ID)).Send(account.UpdatePasswordRequest{
		CurrentPassword: "secret",
		NewPassword:     "new secret",
	}).End()
	require.Len(t, errs, 0, "%s", errs)
	require.Equal(t, http.StatusOK, resp.StatusCode, "%s", body)
}

func TestUpdateUsername(t *testing.T) {
	ts, user := setUpAccountAndServer(t)
	defer ts.Close()
	request := gorequest.New()

	resp, body, errs := request.Put(fmt.Sprintf("%s/accounts/%s/username", ts.URL, user.ID)).Send(account.UpdateUsernameRequest{
		Password: "wrong",
		Username: "secret",
	}).End()
	require.Len(t, errs, 0, "%s", errs)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode, "%s", body)

	resp, body, errs = request.Put(fmt.Sprintf("%s/accounts/%s/username", ts.URL, user.ID)).Send(account.UpdateUsernameRequest{
		Password: "secret",
		Username: "new-username",
	}).End()
	require.Len(t, errs, 0, "%s", errs)
	require.Equal(t, http.StatusOK, resp.StatusCode, "%s", body)

	assert.Nil(t, json.Unmarshal([]byte(body), &user))
	assert.Equal(t, "new-username", user.Username)
}

func TestUpdateData(t *testing.T) {
	updateData := `{"update": "data"}`
	ts, user := setUpAccountAndServer(t)
	defer ts.Close()
	request := gorequest.New()

	resp, body, errs := request.Put(fmt.Sprintf("%s/accounts/%s/data", ts.URL, user.ID)).Send(account.UpdateDataRequest{Data: "not json"}).End()
	require.Len(t, errs, 0, "%s", errs)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode, "%s", body)

	resp, body, errs = request.Put(fmt.Sprintf("%s/accounts/%s/data", ts.URL, user.ID)).Send(account.UpdateDataRequest{Data: updateData}).End()
	require.Len(t, errs, 0, "%s", errs)
	require.Equal(t, http.StatusOK, resp.StatusCode, "%s", body)

	assert.Nil(t, json.Unmarshal([]byte(body), &user))
	assert.Equal(t, updateData, user.Data)
}
