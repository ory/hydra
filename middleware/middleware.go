package middleware

import (
	log "github.com/Sirupsen/logrus"
	hydcon "github.com/ory-am/hydra/context"
	"github.com/ory-am/hydra/jwt"
	ladonGuard "github.com/ory-am/ladon/guard"
	"github.com/ory-am/ladon/policy"
	"golang.org/x/net/context"
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

func (m *Middleware) ExtractAuthentication(next hydcon.ContextHandler) hydcon.ContextHandler {
	return hydcon.ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		ctx = hydcon.NewContextFromAuthorization(ctx, req, m.jwtService, m.policyStore)
		next.ServeHTTPContext(ctx, rw, req)
	})
}

func (m *Middleware) IsAuthorized(resource, permission string, environment *env) func(hydcon.ContextHandler) hydcon.ContextHandler {
	return func(next hydcon.ContextHandler) hydcon.ContextHandler {
		return hydcon.ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			if environment == nil {
				environment = Env(req)
			}

			policies, err := hydcon.PoliciesFromContext(ctx)
			if err != nil {
				log.WithFields(log.Fields{
					"authorization": "forbidden",
					"error":         err,
				}).Warnf(`Policy extraction failed.`)
				errorHandler(rw, req, http.StatusForbidden)
				return
			}

			subject, err := hydcon.SubjectFromContext(ctx)
			if err != nil {
				log.WithFields(log.Fields{
					"authorization": "forbidden",
					"error":         err,
				}).Warnf(`Subject extraction failed.`)
				errorHandler(rw, req, http.StatusForbidden)
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
				errorHandler(rw, req, http.StatusForbidden)
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

func (m *Middleware) IsAuthenticated(next hydcon.ContextHandler) hydcon.ContextHandler {
	return hydcon.ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		if !hydcon.IsAuthenticatedFromContext(ctx) {
			log.WithFields(log.Fields{"authentication": "fail"}).Warn(`Not able to get authorization from context.`)
			errorHandler(rw, req, http.StatusUnauthorized)
			return
		}

		subject, err := hydcon.SubjectFromContext(ctx)
		if err != nil {
			log.WithFields(log.Fields{"authentication": "fail"}).Warnf("Subject extraction failed: %s", err)
			errorHandler(rw, req, http.StatusUnauthorized)
			return
		} else if subject == "" {
			log.WithFields(log.Fields{"authentication": "fail"}).Warnf("No subject given.")
			errorHandler(rw, req, http.StatusUnauthorized)
			return
		}

		log.WithFields(log.Fields{"authentication": "success"}).Infof(`Authenticated subject "%s".`, subject)
		next.ServeHTTPContext(ctx, rw, req)
	})
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
}
