package middleware

import (
	"errors"
	"github.com/RangelReale/osin"
	log "github.com/Sirupsen/logrus"
	chd "github.com/ory-am/common/handler"
	. "github.com/ory-am/hydra/client"
	"github.com/ory-am/hydra/middleware"
	"golang.org/x/net/context"
	"net/http"
)

type Middleware struct {
	Client Client
}

func (m *Middleware) IsAuthorized(resource, permission string, environment *middleware.Env) func(chd.ContextHandler) chd.ContextHandler {
	return func(next chd.ContextHandler) chd.ContextHandler {
		return chd.ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			if environment == nil {
				environment = middleware.NewEnv(req)
			}

			bearer := osin.CheckBearerAuth(req)
			if allowed, err := m.Client.IsAllowed(&AuthorizeRequest{
				Resource:   resource,
				Permission: permission,
				Context:    environment.Ctx(),
				Token:      bearer.Code,
			}); err != nil {
				log.WithFields(log.Fields{
					"authorization": "forbidden",
					"error":         err,
					"valid":         allowed,
					"permission":    permission,
					"resource":      resource,
				}).Warnf(`Subject is not allowed perform this action on this resource.`)
				rw.WriteHeader(http.StatusForbidden)
				return
			} else if !allowed {
				log.WithFields(log.Fields{
					"authorization": "forbidden",
					"error":         nil,
					"valid":         allowed,
					"permission":    permission,
					"resource":      resource,
				}).Warnf(`Subject is not allowed perform this action on this resource.`)
				rw.WriteHeader(http.StatusForbidden)
				return
			}

			log.WithFields(log.Fields{"authorization": "success"}).Info(`Allowed!`)
			next.ServeHTTPContext(ctx, rw, req)
		})
	}
}

func (m *Middleware) IsAuthenticated(next chd.ContextHandler) chd.ContextHandler {
	return chd.ContextHandlerFunc(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
		bearer := osin.CheckBearerAuth(req)
		if bearer == nil {
			log.WithFields(log.Fields{
				"authentication": "invalid",
				"error":          errors.New("No bearer token given"),
				"valid":          false,
			}).Warn(`Authentication invalid.`)
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
		if authenticated, err := m.Client.IsAuthenticated(bearer.Code); err != nil {
			log.WithFields(log.Fields{
				"authentication": "invalid",
				"error":          err,
				"valid":          authenticated,
			}).Warn(`Authentication invalid.`)
			rw.WriteHeader(http.StatusUnauthorized)
			return
		} else if !authenticated {
			log.WithFields(log.Fields{
				"authentication": "invalid",
				"error":          nil,
				"valid":          authenticated,
			}).Warn(`Authentication invalid.`)
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		log.WithFields(log.Fields{"authentication": "success"}).Info(`Authenticated.`)
		next.ServeHTTPContext(ctx, rw, req)
	})
}
