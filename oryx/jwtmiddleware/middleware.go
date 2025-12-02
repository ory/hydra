// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwtmiddleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"

	"github.com/ory/herodot"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/urfave/negroni"

	"github.com/ory/x/jwksx"
)

// Deprecated: use jwtmiddleware.ContextKey{} instead.
var SessionContextKey = jwtmiddleware.ContextKey{}

type Middleware struct {
	o   *middlewareOptions
	wku string
	jm  *jwtmiddleware.JWTMiddleware
}

type middlewareOptions struct {
	Debug         bool
	ExcludePaths  []string
	SigningMethod jwt.SigningMethod
	ErrorWriter   herodot.Writer
}

type MiddlewareOption func(*middlewareOptions)

func SessionFromContext(ctx context.Context) (json.RawMessage, error) {
	raw := ctx.Value(jwtmiddleware.ContextKey{})
	if raw == nil {
		return nil, errors.WithStack(herodot.ErrUnauthorized.WithReasonf("Could not find credentials in the request."))
	}

	token, ok := raw.(*jwt.Token)
	if !ok {
		return nil, errors.WithStack(herodot.ErrInternalServerError.WithDebugf(`Expected context key "%s" to transport value of type *jwt.MapClaims but got type: %T`, SessionContextKey, raw))
	}

	session, err := json.Marshal(token.Claims)
	if err != nil {
		return nil, errors.WithStack(herodot.ErrInternalServerError.WithDebugf("Unable to encode session data: %s", err))
	}

	return session, nil
}

func MiddlewareDebugEnabled() MiddlewareOption {
	return func(o *middlewareOptions) {
		o.Debug = true
	}
}

func MiddlewareExcludePaths(paths ...string) MiddlewareOption {
	return func(o *middlewareOptions) {
		o.ExcludePaths = append(o.ExcludePaths, paths...)
	}
}

func MiddlewareAllowSigningMethod(method jwt.SigningMethod) MiddlewareOption {
	return func(o *middlewareOptions) {
		o.SigningMethod = method
	}
}

func MiddlewareErrorWriter(w herodot.Writer) MiddlewareOption {
	return func(o *middlewareOptions) {
		o.ErrorWriter = w
	}
}

func NewMiddleware(
	wellKnownURL string,
	opts ...MiddlewareOption,
) *Middleware {
	c := &middlewareOptions{
		SigningMethod: jwt.SigningMethodES256,
		ErrorWriter:   herodot.NewJSONWriter(nil),
	}

	for _, o := range opts {
		o(c)
	}
	jc := jwksx.NewFetcher(wellKnownURL)
	return &Middleware{
		o:   c,
		wku: wellKnownURL,
		jm: jwtmiddleware.New(
			func(ctx context.Context, rawToken string) (any, error) {
				return jwt.NewParser(
					jwt.WithValidMethods([]string{c.SigningMethod.Alg()}),
				).Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
					if raw, ok := token.Header["kid"]; !ok {
						return nil, errors.New(`jwt from authorization HTTP header is missing value for "kid" in token header`)
					} else if kid, ok := raw.(string); !ok {
						return nil, fmt.Errorf(`jwt from authorization HTTP header is expecting string value for "kid" in tokenWithoutKid header but got: %T`, raw)
					} else if k, err := jc.GetKey(kid); err != nil {
						return nil, err
					} else {
						return k.Key, nil
					}
				})
			},
			jwtmiddleware.WithCredentialsOptional(false),
			jwtmiddleware.WithTokenExtractor(func(r *http.Request) (string, error) {
				// wrapping the extractor to get a herodot.ErrorContainer
				token, err := jwtmiddleware.AuthHeaderTokenExtractor(r)
				if err != nil {
					return "", herodot.ErrUnauthorized.WithReason(err.Error())
				}
				return token, nil
			}),
			jwtmiddleware.WithErrorHandler(func(w http.ResponseWriter, r *http.Request, err error) {
				switch {
				case errors.Is(err, jwtmiddleware.ErrJWTInvalid):
					reason := "The token is invalid or expired."
					if err := errors.Unwrap(err); err != nil {
						reason = err.Error()
					}
					c.ErrorWriter.WriteError(w, r, errors.WithStack(herodot.ErrUnauthorized.WithReason(reason)))
				case errors.Is(err, jwtmiddleware.ErrJWTMissing):
					c.ErrorWriter.WriteError(w, r, errors.WithStack(herodot.ErrUnauthorized.WithReason("The token is missing.")))
				default:
					c.ErrorWriter.WriteError(w, r, err)
				}
			}),
		),
	}
}

// Deprecated: use Middleware as a negroni.Handler directly instead.
func (h *Middleware) NegroniHandler() negroni.Handler {
	return negroni.HandlerFunc(h.ServeHTTP)
}

func (h *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	for _, excluded := range h.o.ExcludePaths {
		if strings.HasPrefix(r.URL.Path, excluded) {
			next(w, r)
			return
		}
	}

	h.jm.CheckJWT(next).ServeHTTP(w, r)
}
