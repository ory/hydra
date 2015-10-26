package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/ory-am/dockertest"
	"github.com/ory-am/hydra/account"
	hydra "github.com/ory-am/hydra/account/postgres"
	hcon "github.com/ory-am/hydra/context"
	"github.com/ory-am/hydra/handler/middleware"
	"github.com/ory-am/hydra/hash"
	"github.com/ory-am/ladon/policy"
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
	mw *middleware.Middleware
	hd *Handler
	s  *hydra.Store
	db *sql.DB
)

func TestMain(m *testing.M) {
	var err error
	var c dockertest.ContainerID
	c, db, err = dockertest.OpenPostgreSQLContainerConnection(15, time.Second)
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	defer c.KillRemove()

	s = hydra.New(&hash.BCrypt{10}, db)
	if err := s.CreateSchemas(); err != nil {
		log.Fatalf("Could not set up schemas: %v", err)
	}

	mw = &middleware.Middleware{}
	hd = &Handler{s, mw}
	os.Exit(m.Run())
}

type payload struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Data     string `json:"data"`
}

type test struct {
	subject      account.Account
	token        *jwt.Token
	policies     []policy.Policy
	createData   *payload
	statusCreate int
	statusGet    int
	StatusDelete int
}

var cases = []*test{
	&test{
		&account.DefaultAccount{},
		&jwt.Token{Valid: false},
		[]policy.Policy{},
		&payload{},
		http.StatusUnauthorized, 0, 0,
	},
	&test{
		&account.DefaultAccount{},
		&jwt.Token{Valid: true},
		[]policy.Policy{},
		&payload{},
		http.StatusForbidden, 0, 0,
	},
	&test{
		&account.DefaultAccount{},
		&jwt.Token{Valid: true},
		[]policy.Policy{},
		&payload{},
		http.StatusForbidden, 0, 0,
	},
	&test{
		&account.DefaultAccount{ID: "max"},
		&jwt.Token{Valid: true},
		[]policy.Policy{},
		&payload{},
		http.StatusForbidden, 0, 0,
	},
	&test{
		&account.DefaultAccount{ID: "max"},
		&jwt.Token{Valid: true},
		[]policy.Policy{
			&policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"/users"}, []string{"create"}},
		},
		&payload{},
		http.StatusForbidden, 0, 0,
	},
	&test{
		&account.DefaultAccount{ID: "peter"},
		&jwt.Token{Valid: true},
		[]policy.Policy{
			&policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"/users"}, []string{"create"}},
		},
		&payload{Email: "1@bar.com", Password: "secret", Data: "{}"},
		http.StatusOK, http.StatusForbidden, http.StatusForbidden,
	},
	&test{
		&account.DefaultAccount{ID: "peter"},
		&jwt.Token{Valid: true},
		[]policy.Policy{
			&policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"/users"}, []string{"create"}},
			&policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{".*"}, []string{"get"}},
		},
		&payload{Email: "2@bar.com", Password: "secret", Data: "{}"},
		http.StatusOK, http.StatusOK, http.StatusForbidden,
	},
	&test{
		&account.DefaultAccount{ID: "peter"},
		&jwt.Token{Valid: true},
		[]policy.Policy{
			&policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"/users"}, []string{"create"}},
			&policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{".*"}, []string{"get"}},
			&policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"/users/.*"}, []string{"delete"}},
		},
		&payload{Email: "3@bar.com", Password: "secret", Data: "{}"},
		http.StatusOK, http.StatusOK, http.StatusAccepted,
	},
}

func mock(c *test) func(h hcon.ContextHandler) hcon.ContextHandler {
	return func(h hcon.ContextHandler) hcon.ContextHandler {
		return hcon.ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			ctx = hcon.NewContextFromAuthValues(ctx, c.subject, c.token, c.policies)
			h.ServeHTTPContext(ctx, rw, req)
		})
	}
}

func TestCreateGetDelete(t *testing.T) {
	run := func(t *testing.T, name string, k int, testCase *test, code int, request func(c *test) *http.Request, finish func(c *test, res *httptest.ResponseRecorder)) {
		router := mux.NewRouter()
		hd.SetRoutes(router, mock(testCase))
		req := request(testCase)
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)
		assert.Equal(t, code, res.Code, "Case %d, %s", k, name)
		if http.StatusOK == res.Code || http.StatusAccepted == res.Code {
			finish(testCase, res)
		}
	}

	for k, c := range cases {
		var p payload
		var code int
		run(t, "create", k, c, c.statusCreate, func(c *test) *http.Request {
			data, _ := json.Marshal(c.createData)
			req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(data))
			req.Header.Set("Content-Type", "application/json")
			return req
		}, func(c *test, res *httptest.ResponseRecorder) {
			code = res.Code
			result := res.Body.Bytes()
			log.Printf("POST ok /users: %s", result)
			require.Nil(t, json.Unmarshal(result, &p))
			assert.Equal(t, c.createData.Email, p.Email)
			assert.Equal(t, c.createData.Data, p.Data)
			assert.Empty(t, p.Password)
		})

		if code != http.StatusOK {
			continue
		}

		run(t, "get", k, c, c.statusGet, func(c *test) *http.Request {
			req, _ := http.NewRequest("GET", "/users/"+p.ID, nil)
			return req
		}, func(c *test, res *httptest.ResponseRecorder) {
			code = res.Code
			result := res.Body.Bytes()
			log.Printf("GET ok /users/%s: %s", p.ID, result)
			require.Nil(t, json.Unmarshal(result, &p))
			assert.Equal(t, c.createData.Email, p.Email)
			assert.Equal(t, c.createData.Data, p.Data)
			assert.Empty(t, p.Password)
		})

		if code != http.StatusOK {
			continue
		}

		run(t, "delete", k, c, c.StatusDelete, func(c *test) *http.Request {
			req, _ := http.NewRequest("DELETE", "/users/"+p.ID, nil)
			return req
		}, func(c *test, res *httptest.ResponseRecorder) {
			code = res.Code
			log.Printf("DELETE ok /users/%s", p.ID)
		})

		if code != http.StatusAccepted {
			continue
		}

		run(t, "get after delete", k, c, http.StatusNotFound, func(c *test) *http.Request {
			req, _ := http.NewRequest("GET", "/users/"+p.ID, nil)
			return req
		}, func(c *test, res *httptest.ResponseRecorder) {
			assert.True(t, false)
		})
	}
}
