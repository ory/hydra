package integration_test

import (
	"net/http"
	"testing"

	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/handler/oauth2"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func tokenRevocationHandler(t *testing.T, oauth2 fosite.OAuth2Provider, session fosite.Session) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx := fosite.NewContext()
		err := oauth2.NewRevocationRequest(ctx, req)
		if err != nil {
			t.Logf("Revoke request failed because %s.", err.Error())
			t.Logf("Stack: %v", err.(stackTracer).StackTrace())
		}
		oauth2.WriteRevocationResponse(rw, err)
	}
}

func tokenIntrospectionHandler(t *testing.T, oauth2 fosite.OAuth2Provider, session fosite.Session) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx := fosite.NewContext()
		ar, err := oauth2.NewIntrospectionRequest(ctx, req, session)
		if err != nil {
			t.Logf("Introspection request failed because %s.", err.Error())
			t.Logf("Stack: %s", err.(stackTracer).StackTrace())
			oauth2.WriteIntrospectionError(rw, err)
			return
		}

		oauth2.WriteIntrospectionResponse(rw, ar)
	}
}

func tokenInfoHandler(t *testing.T, oauth2 fosite.OAuth2Provider, session fosite.Session) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx := fosite.NewContext()
		if _, err := oauth2.IntrospectToken(ctx, fosite.AccessTokenFromRequest(req), fosite.AccessToken, session); err != nil {
			rfcerr := fosite.ErrorToRFC6749Error(err)
			t.Logf("Info request failed because `%s`.", err.Error())
			t.Logf("Stack: %s", err.(stackTracer).StackTrace())
			http.Error(rw, rfcerr.Description, rfcerr.StatusCode)
			return
		}

		rw.WriteHeader(http.StatusNoContent)
	}
}

func authEndpointHandler(t *testing.T, oauth2 fosite.OAuth2Provider, session fosite.Session) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx := fosite.NewContext()

		ar, err := oauth2.NewAuthorizeRequest(ctx, req)
		if err != nil {
			t.Logf("Access request failed because %s.", err.Error())
			t.Logf("Request: %s.", ar)
			t.Logf("Stack: %s.", err.(stackTracer).StackTrace())
			oauth2.WriteAuthorizeError(rw, ar, err)
			return
		}

		if ar.GetRequestedScopes().Has("fosite") {
			ar.GrantScope("fosite")
		}
		if ar.GetRequestedScopes().Has("offline") {
			ar.GrantScope("offline")
		}
		if ar.GetRequestedScopes().Has("openid") {
			ar.GrantScope("openid")
		}

		// Normally, this would be the place where you would check if the user is logged in and gives his consent.
		// For this test, let's assume that the user exists, is logged in, and gives his consent...

		response, err := oauth2.NewAuthorizeResponse(ctx, req, ar, session)
		if err != nil {
			t.Logf("Access request failed because %s.", err.Error())
			t.Logf("Request: %s.", ar)
			t.Logf("Stack: %s.", err.(stackTracer).StackTrace())
			oauth2.WriteAuthorizeError(rw, ar, err)
			return
		}

		oauth2.WriteAuthorizeResponse(rw, ar, response)
	}
}

func authCallbackHandler(t *testing.T) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		if q.Get("code") == "" && q.Get("error") == "" {
			assert.NotEmpty(t, q.Get("code"))
			assert.NotEmpty(t, q.Get("error"))
		}

		if q.Get("code") != "" {
			rw.Write([]byte("code: ok"))
		}
		if q.Get("error") != "" {
			rw.WriteHeader(http.StatusNotAcceptable)
			rw.Write([]byte("error: " + q.Get("error")))
		}

	}
}

func tokenEndpointHandler(t *testing.T, provider fosite.OAuth2Provider) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		ctx := fosite.NewContext()

		accessRequest, err := provider.NewAccessRequest(ctx, req, &oauth2.JWTSession{})
		if err != nil {
			t.Logf("Access request failed because %s.", err.Error())
			t.Logf("Request: %s.", accessRequest)
			t.Logf("Stack: %v.", err.(stackTracer).StackTrace())
			provider.WriteAccessError(rw, accessRequest, err)
			return
		}

		if accessRequest.GetRequestedScopes().Has("fosite") {
			accessRequest.GrantScope("fosite")
		}

		response, err := provider.NewAccessResponse(ctx, req, accessRequest)
		if err != nil {
			t.Logf("Access request failed because %s.", err.Error())
			t.Logf("Request: %s.", accessRequest)
			t.Logf("Stack: %v.", err.(stackTracer).StackTrace())
			provider.WriteAccessError(rw, accessRequest, err)
			return
		}

		provider.WriteAccessResponse(rw, accessRequest, response)
	}
}
