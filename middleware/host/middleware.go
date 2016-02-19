package host

import (
	log "github.com/Sirupsen/logrus"
	"github.com/go-errors/errors"
	chd "github.com/ory-am/common/handler"
	"github.com/ory-am/common/pkg"
	ladonGuard "github.com/ory-am/ladon/guard"
	"github.com/ory-am/ladon/policy"
	"golang.org/x/net/context"
	authcon "github.com/ory-am/hydra/context"
	"github.com/ory-am/hydra/jwt"
	"github.com/ory-am/hydra/middleware"
	"net/http"
)

type Middleware struct {
	policyStore policy.Storage
	jwtService  *jwt.JWT
}

func New(policyStore policy.Storage, jwtService *jwt.JWT) *Middleware {
	return &Middleware{policyStore, jwtService}
}

var guard = &ladonGuard.Guard{}

func (m *Middleware) ExtractAuthentication(next chd.ContextHandler) chd.ContextHandler {
	return chd.ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		ctx = authcon.NewContextFromAuthorization(ctx, req, m.jwtService, m.policyStore)
		next.ServeHTTPContext(ctx, rw, req)
	})
}

func (m *Middleware) IsAuthorized(resource, permission string, environment *middleware.Env) func(chd.ContextHandler) chd.ContextHandler {
	return func(next chd.ContextHandler) chd.ContextHandler {
		return chd.ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			if environment == nil {
				environment = middleware.NewEnv(req)
			}

			policies, err := authcon.PoliciesFromContext(ctx)
			if err != nil {
				log.WithFields(log.Fields{
					"authorization": "forbidden",
					"error":         err,
				}).Warnf(`Policy extraction failed.`)
				pkg.HttpError(rw, errors.New("Forbidden"), http.StatusForbidden)
				return
			}

			subject, err := authcon.SubjectFromContext(ctx)
			if err != nil {
				log.WithFields(log.Fields{
					"authorization": "forbidden",
					"error":         err,
				}).Warnf(`Subject extraction failed.`)
				pkg.HttpError(rw, errors.New("Forbidden"), http.StatusForbidden)
				return
			}

			ok, err := guard.IsGranted(resource, permission, subject, policies, environment.Ctx())
			if err != nil || !ok {
				log.WithFields(log.Fields{
					"authorization": "forbidden",
					"error":         err,
					"valid":         ok,
					"subject":       subject,
					"permission":    permission,
					"resource":      resource,
				}).Warnf(`Subject is not allowed perform this action on this resource.`)
				pkg.HttpError(rw, errors.New("Forbidden"), http.StatusForbidden)
				return
			}

			log.WithFields(log.Fields{
				"authorization": "success",
				"subject":       subject,
				"permission":    permission,
				"resource":      resource,
			}).Infof(`Access granted.`)
			next.ServeHTTPContext(ctx, rw, req)
		})
	}
}

func (m *Middleware) IsAuthenticated(next chd.ContextHandler) chd.ContextHandler {
	return chd.ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		if !authcon.IsAuthenticatedFromContext(ctx) {
			log.WithFields(log.Fields{"authentication": "fail"}).Warn(`Not able to get authorization from context.`)
			pkg.HttpError(rw, errors.New("Unauthorized"), http.StatusUnauthorized)
			return
		}

		subject, err := authcon.SubjectFromContext(ctx)
		if err != nil {
			log.WithFields(log.Fields{"authentication": "fail"}).Warnf("Subject extraction failed: %s", err)
			pkg.HttpError(rw, errors.New("Unauthorized"), http.StatusUnauthorized)
			return
		} else if subject == "" {
			log.WithFields(log.Fields{"authentication": "fail"}).Warnf("No subject given.")
			pkg.HttpError(rw, errors.New("Unauthorized"), http.StatusUnauthorized)
			return
		}

		log.WithFields(log.Fields{"authentication": "success"}).Infof(`Authenticated subject "%s".`, subject)
		next.ServeHTTPContext(ctx, rw, req)
	})
}
