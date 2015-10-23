package middleware

import (
	"github.com/ory-am/hydra/account"
	hydcon "github.com/ory-am/hydra/context"
	"github.com/ory-am/hydra/jwt"
	ladonGuard "github.com/ory-am/ladon/guard"
	"github.com/ory-am/ladon/policy"
	"golang.org/x/net/context"
	"log"
	"net/http"
)

type Middleware struct {
	accountStore account.Storage
	policyStore  policy.Storer
	jwtService   *jwt.JWT
}

func New(accountStore account.Storage, policyStore policy.Storer, jwtService *jwt.JWT) *Middleware {
	return &Middleware{accountStore, policyStore, jwtService}
}

var guard = &ladonGuard.Guard{}

func (m *Middleware) ExtractAuthentication(h hydcon.ContextHandler) hydcon.ContextHandler {
	return hydcon.ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		ctx = hydcon.NewContextFromAuthorization(ctx, req, m.jwtService, m.accountStore, m.policyStore)
		h.ServeHTTPContext(ctx, rw, req)
	})
}

func (m *Middleware) IsAuthorized(h hydcon.ContextHandler, resource, permission string) hydcon.ContextHandler {
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

		ok, err := guard.IsGranted(resource, permission, subject.GetID(), policies)
		if err != nil || !ok {
			log.Printf(`Unauthorized: Subject "%s" is not being granted access "%s" to resource "%s"`, subject.GetID(), permission, resource)
			errorHandler(rw, req, http.StatusForbidden)
			return
		}

		log.Printf(`Authorized subject "%s"`, subject.GetID())
		h.ServeHTTPContext(ctx, rw, req)
	})
}

func (m *Middleware) IsAuthenticated(h hydcon.ContextHandler) hydcon.ContextHandler {
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
		}

		log.Printf(`Authenticated subject "%s"`, subject.GetID())
		h.ServeHTTPContext(ctx, rw, req)
	})
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
}
