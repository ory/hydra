package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/RangelReale/osin"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	chd "github.com/ory-am/common/handler"
	authcon "github.com/ory-am/hydra/context"
	hjwt "github.com/ory-am/hydra/jwt"
	middleware "github.com/ory-am/hydra/middleware/host"
	"github.com/ory-am/ladon/guard"
	"github.com/ory-am/ladon/guard/operator"
	"github.com/ory-am/ladon/policy"
	"github.com/ory-am/ladon/policy/postgres"
	opg "github.com/ory-am/osin-storage/storage/postgres"
	"github.com/parnurzeal/gorequest"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"gopkg.in/ory-am/dockertest.v2"
)

var (
	mw    *middleware.Middleware
	store *postgres.Store
	o     *opg.Storage
	db    *sql.DB
)

var jwtService = hjwt.New([]byte(hjwt.TestCertificates[0][1]), []byte(hjwt.TestCertificates[1][1]))

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
	o = opg.New(db)
	if err := store.CreateSchemas(); err != nil {
		log.Fatalf("Could not create schemas: %s", err)
	}
	if err := o.CreateSchemas(); err != nil {
		log.Fatalf("Could not create schemas: %s", err)
	}

	if err := o.CreateClient(&osin.DefaultClient{
		Id:     "app",
		Secret: "secret",
	}); err != nil {
		log.Fatalf("Could not create app: %s", err)
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
	createData policy.DefaultPolicy
	//allowedPayload payload

	statusCreate int
	statusGet    int
	//expectsAllowed       bool
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
	"pass-all":    &policy.DefaultPolicy{"", "", []string{"api-app"}, policy.AllowAccess, []string{"rn:hydra:policies<.*>"}, []string{"<.*>"}, nil},
	"pass-create": &policy.DefaultPolicy{"", "", []string{"api-app"}, policy.AllowAccess, []string{"rn:hydra:policies"}, []string{"create"}, nil},
	"pass-get":    &policy.DefaultPolicy{"", "", []string{"api-app"}, policy.AllowAccess, []string{"rn:hydra:policies:<.*>"}, []string{"get"}, nil},
	"pass-delete": &policy.DefaultPolicy{"", "", []string{"api-app"}, policy.AllowAccess, []string{"rn:hydra:policies:<.*>"}, []string{"delete"}, nil},
	"fail":        &policy.DefaultPolicy{},
}

var payloads = []policy.DefaultPolicy{
	{
		Description: "description",
		Subjects:    []string{"max"},
		Effect:      policy.AllowAccess,
		Resources:   []string{"resource"},
		Permissions: []string{"<.*>"},
		Conditions:  nil,
	},
	{
		Description: "Should allow peter all permissions on resource article",
		Subjects:    []string{"peter"},
		Effect:      policy.AllowAccess,
		Resources:   []string{"article"},
		Permissions: []string{"<.*>"},
		Conditions: []policy.DefaultCondition{
			{
				Operator: "SubjectIsOwner",
				Extra:    map[string]interface{}{"bar": "foo"},
			},
		},
	},
}

func TestGrantedEndpoint(t *testing.T) {
	c := test{
		subject:  "api-app",
		token:    jwt.Token{Valid: true},
		policies: []policy.Policy{policies["pass-all"]},
	}

	handler := &Handler{s: store, m: mw, g: &guard.Guard{}, j: jwtService, o: o}
	router := mux.NewRouter()
	handler.SetRoutes(router, mockAuthorization(c))
	ts := httptest.NewServer(router)
	defer ts.Close()

	request := gorequest.New()
	resp, _, _ := request.Post(ts.URL + "/policies").Send(payloads[1]).End()
	require.Equal(t, 201, resp.StatusCode)

	num := 0
	var isAllowed struct {
		Allowed bool `json:"allowed"`
	}
	do := func(p GrantedPayload, username, password string, shouldAllow bool) {
		resp, body, _ := request.Post(ts.URL+"/guard/allowed").SetBasicAuth(username, password).Send(p).End()
		if username != "app" || password != "secret" {
			require.Equal(t, 401, resp.StatusCode, "Case %d", num)
			return
		}

		require.Equal(t, 200, resp.StatusCode, "Case %d", num)
		json.Unmarshal([]byte(body), &isAllowed)
		require.Equal(t, 200, resp.StatusCode, "Case %d", num)
		require.Equal(t, shouldAllow, isAllowed.Allowed, "Case %d", num)
		num++
	}

	do(GrantedPayload{
		Resource: "article",
		// sub: api-app
		Token:      "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhcGktYXBwIiwiZXhwIjoxOTI0OTc1NjE5fQ.jipsxS1s5xnyZ2K9EqL33y9B6dWuDB6gzgA3M0rLUS1bcOcSj9hVQMAxcl6Udezid057denHH6a5LrbcuGqwTi7bMlSCs_eWIoTQ5WKTvd0PxEMJGyjw9MUWStHJWna2Drp_vXhZGVvkUbXCRAkVO8KCkKWUB5-wNfoNh6ba-_c7zppcyIV7aRwSFJ5Eu2Gq_dwlNWmu-GB8hTbhHEcXTkBDjRsy6oITfpwGRkxvzmJmYXJKRUFsNlt8DJaWHguOszWGEjfJeOhooybnrUHiwgEwVuciHptI50UaQYDjvBQolLUrcnkf98bQXJsALoBYkaHFC87mVzv0ZR_ZPTzb2A",
		Permission: "random",
		Context: &operator.Context{
			Owner: "peter",
		},
	}, "app", "secret", false)

	do(GrantedPayload{
		Resource: "article",
		// sub: peter
		Token:      "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJwZXRlciIsImV4cCI6MTkyNDk3NTYxOX0.GVn0YAQTFFoIa-fcsqQgq3pgWBAYNsbd9SqoXUPt7EK63zqiZ0yVqWgQCBEXU5NyT96Alg1Se6Pq6wzAC4ydof-MN3nQhcoNhx6QEHBGFDwwsHwMVyi-51S0NXzYXSV-gGrPoOloCkOSoyab-RWdMZ6LrgV5WQOW4WAfYL0nJ0I-WxlXcoKi-8MJ1GqScqC_E0v9cn4iNAT5e1tPMT49KdjOo_HYPQlJQjcJ724USdDWywPxZy5AmYxG5A2XeaY41Ly0O0HJ8Q56I2ukPMfXiTpnm5mnb9mRbK99HnvlAvtEKJ-Lf0w_BTurL_3ZmONKSYR0HHIMZC0hO9NJNNTS1Q",
		Permission: "random",
		Context: &operator.Context{
			Owner: "peter",
		},
	}, "app", "secret", true)

	do(GrantedPayload{
		Resource: "article",
		// sub: peter
		Token:      "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJwZXRlciIsImV4cCI6MTkyNDk3NTYxOX0.GVn0YAQTFFoIa-fcsqQgq3pgWBAYNsbd9SqoXUPt7EK63zqiZ0yVqWgQCBEXU5NyT96Alg1Se6Pq6wzAC4ydof-MN3nQhcoNhx6QEHBGFDwwsHwMVyi-51S0NXzYXSV-gGrPoOloCkOSoyab-RWdMZ6LrgV5WQOW4WAfYL0nJ0I-WxlXcoKi-8MJ1GqScqC_E0v9cn4iNAT5e1tPMT49KdjOo_HYPQlJQjcJ724USdDWywPxZy5AmYxG5A2XeaY41Ly0O0HJ8Q56I2ukPMfXiTpnm5mnb9mRbK99HnvlAvtEKJ-Lf0w_BTurL_3ZmONKSYR0HHIMZC0hO9NJNNTS1Q",
		Permission: "foobar",
		Context: &operator.Context{
			Owner: "peter",
		},
	}, "app", "secret", true)

	do(GrantedPayload{
		Resource: "article",
		// sub: peter
		Token:      "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJwZXRlciIsImV4cCI6MTkyNDk3NTYxOX0.GVn0YAQTFFoIa-fcsqQgq3pgWBAYNsbd9SqoXUPt7EK63zqiZ0yVqWgQCBEXU5NyT96Alg1Se6Pq6wzAC4ydof-MN3nQhcoNhx6QEHBGFDwwsHwMVyi-51S0NXzYXSV-gGrPoOloCkOSoyab-RWdMZ6LrgV5WQOW4WAfYL0nJ0I-WxlXcoKi-8MJ1GqScqC_E0v9cn4iNAT5e1tPMT49KdjOo_HYPQlJQjcJ724USdDWywPxZy5AmYxG5A2XeaY41Ly0O0HJ8Q56I2ukPMfXiTpnm5mnb9mRbK99HnvlAvtEKJ-Lf0w_BTurL_3ZmONKSYR0HHIMZC0hO9NJNNTS1Q",
		Permission: "foobar",
		Context: &operator.Context{
			Owner: "peter",
		},
	}, "foo", "secret", false)

	do(GrantedPayload{
		Resource: "article",
		// sub: peter
		Token:      "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJwZXRlciIsImV4cCI6MTkyNDk3NTYxOX0.GVn0YAQTFFoIa-fcsqQgq3pgWBAYNsbd9SqoXUPt7EK63zqiZ0yVqWgQCBEXU5NyT96Alg1Se6Pq6wzAC4ydof-MN3nQhcoNhx6QEHBGFDwwsHwMVyi-51S0NXzYXSV-gGrPoOloCkOSoyab-RWdMZ6LrgV5WQOW4WAfYL0nJ0I-WxlXcoKi-8MJ1GqScqC_E0v9cn4iNAT5e1tPMT49KdjOo_HYPQlJQjcJ724USdDWywPxZy5AmYxG5A2XeaY41Ly0O0HJ8Q56I2ukPMfXiTpnm5mnb9mRbK99HnvlAvtEKJ-Lf0w_BTurL_3ZmONKSYR0HHIMZC0hO9NJNNTS1Q",
		Permission: "foobar",
		Context: &operator.Context{
			Owner: "peter",
		},
	}, "app", "wrong-secret", false)

	do(GrantedPayload{
		Resource: "foobar",
		// sub: peter
		Token:      "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJwZXRlciIsImV4cCI6MTkyNDk3NTYxOX0.GVn0YAQTFFoIa-fcsqQgq3pgWBAYNsbd9SqoXUPt7EK63zqiZ0yVqWgQCBEXU5NyT96Alg1Se6Pq6wzAC4ydof-MN3nQhcoNhx6QEHBGFDwwsHwMVyi-51S0NXzYXSV-gGrPoOloCkOSoyab-RWdMZ6LrgV5WQOW4WAfYL0nJ0I-WxlXcoKi-8MJ1GqScqC_E0v9cn4iNAT5e1tPMT49KdjOo_HYPQlJQjcJ724USdDWywPxZy5AmYxG5A2XeaY41Ly0O0HJ8Q56I2ukPMfXiTpnm5mnb9mRbK99HnvlAvtEKJ-Lf0w_BTurL_3ZmONKSYR0HHIMZC0hO9NJNNTS1Q",
		Permission: "random",
		Context: &operator.Context{
			Owner: "peter",
		},
	}, "app", "secret", false)

	do(GrantedPayload{
		Resource: "article",
		// sub: foobar
		Token:      "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmb29iYXIiLCJleHAiOjE5MjQ5NzU2MTl9.d4Z9sEB52LWysYXto_mlT41uaLgAETTQJS4iSXjBc7U1lzmT7vsaMpMVNVKhYCe_2ptx7uZcW4pDy8njjQMtFoesAmbUK-finVslYpqjQmyre9eqWURhIXgDu95w2hP9EoSfjXpyE8EUct3a5pkm6rje4C5y-16MrAQpuq3IZVYTPwdS6Gl33BG3Obw3sXheBGMcnmtcGtSQe6ekTqgF-NkVTe5bQPGL6DxGdRLbHOg_nky91JWs4lLO526KVTbDrwM7SVGex5w1rPcn2Qg8RUefbWF2x-KuoAGlTnStfN3tOgw6DW3Q-35fcGesyvy7DAP-Zy68vZ6W7h2rIy6wiQ",
		Permission: "random",
		Context: &operator.Context{
			Owner: "peter",
		},
	}, "app", "secret", false)
}

func TestCreateGetDeleteGet(t *testing.T) {
	for k, c := range []test{
		{
			subject: "api-app", token: jwt.Token{Valid: false},
			policies:     []policy.Policy{policies["fail"]},
			createData:   payloads[0],
			statusCreate: http.StatusUnauthorized,
		},
		{
			subject: "api-app", token: jwt.Token{Valid: true},
			policies:     []policy.Policy{policies["fail"]},
			createData:   payloads[0],
			statusCreate: http.StatusForbidden,
		},
		{
			subject: "api-app", token: jwt.Token{Valid: true},
			policies:     []policy.Policy{policies["pass-create"]},
			createData:   payloads[0],
			statusCreate: http.StatusCreated, statusGet: http.StatusForbidden,
		},
		{
			subject: "api-app", token: jwt.Token{Valid: true},
			policies:     []policy.Policy{policies["pass-create"], policies["pass-get"]},
			createData:   payloads[0],
			statusCreate: http.StatusCreated, statusGet: http.StatusOK, statusDelete: http.StatusForbidden,
		},
		{
			subject: "api-app", token: jwt.Token{Valid: true},
			policies:     []policy.Policy{policies["pass-all"]},
			createData:   payloads[0],
			statusCreate: http.StatusCreated, statusGet: http.StatusOK, statusDelete: http.StatusAccepted, statusGetAfterDelete: http.StatusNotFound,
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
		if resp.StatusCode != http.StatusCreated {
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

var allowedPayloads = map[string]GrantedPayload{
	"create-grant": {
		Resource:   "rn:hydra:policies",
		Token:      "some.token",
		Permission: "create",
	},
}
