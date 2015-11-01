package middleware

import (
	hydcon "github.com/ory-am/hydra/context"
	"github.com/ory-am/hydra/jwt"
	ladonGuard "github.com/ory-am/ladon/guard"
	"github.com/ory-am/ladon/policy"
	"golang.org/x/net/context"
	"log"
	"net/http"
)

type Middleware struct {
	policyStore policy.Storer
	jwtService  *jwt.JWT
}

func New(policyStore policy.Storer, jwtService *jwt.JWT) *Middleware {
	return &Middleware{policyStore, jwtService}
}

var guard = &ladonGuard.Guard{}

func (m *Middleware) ExtractAuthentication(next hydcon.ContextHandler) hydcon.ContextHandler {
	return hydcon.ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		ctx = hydcon.NewContextFromAuthorization(ctx, req, m.jwtService, m.policyStore)
		next.ServeHTTPContext(ctx, rw, req)
	})
}

func (m *Middleware) IsAuthorized(resource, permission string) func(hydcon.ContextHandler) hydcon.ContextHandler {
	return func(next hydcon.ContextHandler) hydcon.ContextHandler {
		return hydcon.ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			policies, err := hydcon.PoliciesFromContext(ctx)
			if err != nil {
				log.Printf("Unauthorized: Policy extraction failed: %s", err)
				errorHandler(rw, req, http.StatusForbidden)
				return
			}

			subject, err := hydcon.SubjectFromContext(ctx)
			if err != nil {
				log.Printf("Unauthorized: Subject extraction failed: %s", err)
				errorHandler(rw, req, http.StatusForbidden)
				return
			}

			ok, err := guard.IsGranted(resource, permission, subject, policies)
			if err != nil || !ok {
				log.Printf(`Unauthorized: Subject "%s" is not being granted access "%s" to resource "%s"`, subject, permission, resource)
				errorHandler(rw, req, http.StatusForbidden)
				return
			}

			log.Printf(`Authorized subject "%s"`, subject)
			next.ServeHTTPContext(ctx, rw, req)
		})
	}
}

func (m *Middleware) IsAuthenticated(next hydcon.ContextHandler) hydcon.ContextHandler {
	return hydcon.ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		if !hydcon.IsAuthenticatedFromContext(ctx) {
			log.Printf("Not authenticated: %s", ctx)
			errorHandler(rw, req, http.StatusUnauthorized)
			return
		}

		subject, err := hydcon.SubjectFromContext(ctx)
		if err != nil {
			log.Printf("Not authenticated: Subject extraction failed: %s", err)
			errorHandler(rw, req, http.StatusUnauthorized)
			return
		} else if subject == "" {
			log.Println("Unauthorized: No subject given.")
			errorHandler(rw, req, http.StatusUnauthorized)
			return
		}

		log.Printf(`Authenticated subject "%s"`, subject)
		next.ServeHTTPContext(ctx, rw, req)
	})
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
}
