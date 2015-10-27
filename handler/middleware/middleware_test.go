package middleware_test

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	hcon "github.com/ory-am/hydra/context"
	. "github.com/ory-am/hydra/handler/middleware"
	hjwt "github.com/ory-am/hydra/jwt"
	"github.com/ory-am/ladon/policy"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type test struct {
	subject            string
	token              *jwt.Token
	policies           []policy.Policy
	resource           string
	permission         string
	passAuthentication bool
	passAuthorization  bool
}

var cases = []test{
	test{
		"max",
		&jwt.Token{Valid: false},
		[]policy.Policy{},
		"", "",
		false, false,
	},
	test{
		"max",
		&jwt.Token{Valid: true},
		[]policy.Policy{},
		"", "",
		true, false,
	},
	test{
		"peter",
		&jwt.Token{Valid: true},
		[]policy.Policy{
			&policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"/articles/74251"}, []string{"create"}},
		},
		"/articles/74251", "create",
		true, true,
	},
	test{
		"peter",
		&jwt.Token{Valid: true},
		[]policy.Policy{
			&policy.DefaultPolicy{"", "", []string{"peter"}, policy.DenyAccess, []string{"/articles/74251"}, []string{"create"}},
		},
		"/articles/74251", "create",
		true, false,
	},
	test{
		"max",
		&jwt.Token{Valid: true},
		[]policy.Policy{
			&policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"/articles/74251"}, []string{"create"}},
		},
		"/articles/74251", "create",
		true, false,
	},
	test{
		"max",
		nil,
		[]policy.Policy{
			&policy.DefaultPolicy{"", "", []string{"peter"}, policy.AllowAccess, []string{"/articles/74251"}, []string{"create"}},
		},
		"/articles/74251", "create",
		false, false,
	},
	test{
		"",
		&jwt.Token{Valid: true},
		[]policy.Policy{},
		"", "",
		false, false,
	},
	test{
		"max",
		&jwt.Token{Valid: true},
		nil,
		"", "",
		true, false,
	},
}

func mockContext(h hcon.ContextHandler, c test) hcon.ContextHandler {
	return hcon.ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		claims := hjwt.NewClaimsCarrier(uuid.New(), "hydra", c.subject, "tests", time.Now(), time.Now())
		ctx = hcon.NewContextFromAuthValues(ctx, claims, c.token, c.policies)
		h.ServeHTTPContext(ctx, rw, req)
	})
}

func handler(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(rw, "ok")
}

func TestIsAuthenticated(t *testing.T) {
	m := &Middleware{}

	for k, c := range cases {
		h := &hcon.ContextAdapter{
			Ctx: context.Background(),
			Handler: mockContext(
				m.IsAuthenticated(
					m.IsAuthorized(hcon.ContextHandlerFunc(handler),
						c.resource,
						c.permission,
					),
				), c,
			),
		}

		ts := httptest.NewServer(h)
		defer ts.Close()

		res, err := http.Get(ts.URL)
		require.Nil(t, err)
		res.Body.Close()

		if !c.passAuthentication {
			assert.Equal(t, http.StatusUnauthorized, res.StatusCode, "Case %d", k)
		} else if !c.passAuthorization {
			assert.Equal(t, http.StatusForbidden, res.StatusCode, "Case %d", k)
		} else {
			assert.Equal(t, http.StatusOK, res.StatusCode, "Case %d", k)
		}
	}
}
