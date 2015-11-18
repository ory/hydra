package host_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"

	chd "github.com/ory-am/common/handler"
	authcon "github.com/ory-am/hydra/context"
	hjwt "github.com/ory-am/hydra/jwt"
	mwroot "github.com/ory-am/hydra/middleware"
	. "github.com/ory-am/hydra/middleware/host"
	. "github.com/ory-am/ladon/policy"
)

type test struct {
	subject     string
	token       *jwt.Token
	policies    []Policy
	resource    string
	permission  string
	owner       string
	expectAuthN bool
	expectAuthZ bool
}

var cases = []test{
	{
		subject:  "max",
		token:    &jwt.Token{Valid: false},
		policies: []Policy{},
		resource: "", permission: "",
		expectAuthN: false, expectAuthZ: false,
	},
	{
		subject:  "max",
		token:    &jwt.Token{Valid: true},
		policies: []Policy{},
		resource: "", permission: "",
		expectAuthN: true, expectAuthZ: false,
	},
	{
		subject: "peter",
		token:   &jwt.Token{Valid: true},
		policies: []Policy{
			&DefaultPolicy{"", "", []string{"peter"}, AllowAccess, []string{"/articles/74251"}, []string{"create"}, nil},
		},
		resource: "/articles/74251", permission: "create",
		expectAuthN: true, expectAuthZ: true,
	},
	{
		subject: "peter",
		token:   &jwt.Token{Valid: true},
		policies: []Policy{
			&DefaultPolicy{"", "", []string{"peter"}, DenyAccess, []string{"/articles/74251"}, []string{"create"}, nil},
		},
		resource: "/articles/74251", permission: "create",
		expectAuthN: true, expectAuthZ: false,
	},
	{
		subject: "max",
		token:   &jwt.Token{Valid: true},
		policies: []Policy{
			&DefaultPolicy{"", "", []string{"peter"}, AllowAccess, []string{"/articles/74251"}, []string{"create"}, nil},
		},
		resource: "/articles/74251", permission: "create",
		expectAuthN: true, expectAuthZ: false,
	},
	{
		subject: "max",
		token:   nil,
		policies: []Policy{
			&DefaultPolicy{"", "", []string{"peter"}, AllowAccess, []string{"/articles/74251"}, []string{"create"}, nil},
		},
		resource: "/articles/74251", permission: "create",
		expectAuthN: false, expectAuthZ: false,
	},
	{
		subject: "max",
		token:   &jwt.Token{Valid: true},
		owner:   "max",
		policies: []Policy{
			&DefaultPolicy{"", "", []string{"<.*>"}, AllowAccess, []string{"/articles/74251"}, []string{"get"}, []Condition{
				&DefaultCondition{Operator: "SubjectIsOwner"},
			}},
		},
		resource: "/articles/74251", permission: "get",
		expectAuthN: true, expectAuthZ: true,
	},
	{
		subject: "max",
		token:   &jwt.Token{Valid: true},
		owner:   "max",
		policies: []Policy{
			&DefaultPolicy{"", "", []string{"<.*>"}, AllowAccess, []string{"/articles/74251"}, []string{"get"}, []Condition{
				&DefaultCondition{Operator: "SubjectIsNotOwner"},
			}},
		},
		resource: "/articles/74251", permission: "get",
		expectAuthN: true, expectAuthZ: false,
	},
	{
		subject: "max",
		token:   &jwt.Token{Valid: true},
		owner:   "max",
		policies: []Policy{
			&DefaultPolicy{"", "", []string{"<.*>"}, AllowAccess, []string{"/articles/74251"}, []string{"get"}, []Condition{
				&DefaultCondition{Operator: "ThisOperatorDoesNotExist"},
			}},
		},
		resource: "/articles/74251", permission: "get",
		expectAuthN: true, expectAuthZ: false,
	},
	{
		subject:  "",
		token:    &jwt.Token{Valid: true},
		policies: []Policy{},
		resource: "", permission: "",
		expectAuthN: false, expectAuthZ: false,
	},
	{
		subject:  "max",
		token:    &jwt.Token{Valid: true},
		policies: nil,
		resource: "", permission: "",
		expectAuthN: true, expectAuthZ: false,
	},
}

func mockContext(c test) func(chd.ContextHandler) chd.ContextHandler {
	return func(next chd.ContextHandler) chd.ContextHandler {
		return chd.ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			claims := hjwt.NewClaimsCarrier(uuid.New(), "hydra", c.subject, "tests", time.Now().Add(time.Hour), time.Now(), time.Now())
			ctx = authcon.NewContextFromAuthValues(ctx, claims, c.token, c.policies)
			next.ServeHTTPContext(ctx, rw, req)
		})
	}
}

func handler(m *Middleware, c test) func(context.Context, http.ResponseWriter, *http.Request) {
	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		m.IsAuthorized(c.resource, c.permission, mwroot.NewEnv(req).Owner(c.owner))(chd.ContextHandlerFunc(
			func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
				fmt.Fprintln(rw, "ok")
			},
		)).ServeHTTPContext(ctx, rw, req)
	}
}

func TestMiddleware(t *testing.T) {
	m := &Middleware{}

	for k, c := range cases {
		h := chd.NewContextAdapter(
			context.Background(),
			mockContext(c),
			m.IsAuthenticated,
		).ThenFunc(chd.ContextHandlerFunc(handler(m, c)))

		ts := httptest.NewServer(h)
		defer ts.Close()

		res, err := http.Get(ts.URL)
		require.Nil(t, err)
		res.Body.Close()

		if !c.expectAuthN {
			assert.Equal(t, http.StatusUnauthorized, res.StatusCode, "Authentication failed case %d", k)
		} else if !c.expectAuthZ {
			assert.Equal(t, http.StatusForbidden, res.StatusCode, "Authorization failed case %d", k)
		} else {
			assert.Equal(t, http.StatusOK, res.StatusCode, "Case %d should be authorized but wasn't.", k)
		}
	}
}
