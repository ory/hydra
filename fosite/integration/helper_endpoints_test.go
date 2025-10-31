// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package integration_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/oauth2"
)

func tokenRevocationHandler(t *testing.T, oauth2 fosite.OAuth2Provider, session fosite.Session) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx := fosite.NewContext()
		err := oauth2.NewRevocationRequest(ctx, req)
		if err != nil {
			t.Logf("Revoke request failed because %+v", err)
		}
		oauth2.WriteRevocationResponse(req.Context(), rw, err)
	}
}

func tokenIntrospectionHandler(t *testing.T, oauth2 fosite.OAuth2Provider, session fosite.Session) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx := fosite.NewContext()
		ar, err := oauth2.NewIntrospectionRequest(ctx, req, session)
		if err != nil {
			t.Logf("Introspection request failed because: %+v", err)
			oauth2.WriteIntrospectionError(req.Context(), rw, err)
			return
		}

		oauth2.WriteIntrospectionResponse(req.Context(), rw, ar)
	}
}

func tokenInfoHandler(t *testing.T, oauth2 fosite.OAuth2Provider, session fosite.Session) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx := fosite.NewContext()
		_, resp, err := oauth2.IntrospectToken(ctx, fosite.AccessTokenFromRequest(req), fosite.AccessToken, session)
		if err != nil {
			t.Logf("Info request failed because: %+v", err)
			var e *fosite.RFC6749Error
			require.True(t, errors.As(err, &e))
			http.Error(rw, e.DescriptionField, e.CodeField)
			return
		}

		t.Logf("Introspecting caused: %+v", resp)

		if err := json.NewEncoder(rw).Encode(resp); err != nil {
			panic(err)
		}
	}
}

func authEndpointHandler(t *testing.T, oauth2 fosite.OAuth2Provider, session fosite.Session) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx := fosite.NewContext()

		ar, err := oauth2.NewAuthorizeRequest(ctx, req)
		if err != nil {
			t.Logf("Access request failed because: %+v", err)
			t.Logf("Request: %+v", ar)
			oauth2.WriteAuthorizeError(req.Context(), rw, ar, err)
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

		for _, a := range ar.GetRequestedAudience() {
			ar.GrantAudience(a)
		}

		// Normally, this would be the place where you would check if the user is logged in and gives his consent.
		// For this test, let's assume that the user exists, is logged in, and gives his consent...

		response, err := oauth2.NewAuthorizeResponse(ctx, ar, session)
		if err != nil {
			t.Logf("Access request failed because: %+v", err)
			t.Logf("Request: %+v", ar)
			oauth2.WriteAuthorizeError(req.Context(), rw, ar, err)
			return
		}

		oauth2.WriteAuthorizeResponse(req.Context(), rw, ar, response)
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
		req.ParseMultipartForm(1 << 20)
		ctx := fosite.NewContext()

		accessRequest, err := provider.NewAccessRequest(ctx, req, &oauth2.JWTSession{})
		if err != nil {
			t.Logf("Access request failed because: %+v", err)
			t.Logf("Request: %+v", accessRequest)
			provider.WriteAccessError(req.Context(), rw, accessRequest, err)
			return
		}

		if accessRequest.GetRequestedScopes().Has("fosite") {
			accessRequest.GrantScope("fosite")
		}

		response, err := provider.NewAccessResponse(ctx, accessRequest)
		if err != nil {
			t.Logf("Access request failed because: %+v", err)
			t.Logf("Request: %+v", accessRequest)
			provider.WriteAccessError(req.Context(), rw, accessRequest, err)
			return
		}

		provider.WriteAccessResponse(req.Context(), rw, accessRequest, response)
	}
}

func pushedAuthorizeRequestHandler(t *testing.T, oauth2 fosite.OAuth2Provider, session fosite.Session) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx := fosite.NewContext()

		ar, err := oauth2.NewPushedAuthorizeRequest(ctx, req)
		if err != nil {
			t.Logf("PAR request failed because: %+v", err)
			t.Logf("Request: %+v", ar)
			oauth2.WritePushedAuthorizeError(ctx, rw, ar, err)
			return
		}

		response, err := oauth2.NewPushedAuthorizeResponse(ctx, ar, session)
		if err != nil {
			t.Logf("PAR response failed because: %+v", err)
			t.Logf("Request: %+v", ar)
			oauth2.WritePushedAuthorizeError(ctx, rw, ar, err)
			return
		}

		oauth2.WritePushedAuthorizeResponse(ctx, rw, ar, response)
	}
}

func deviceAuthorizationEndpointHandler(t *testing.T, oauth2 fosite.OAuth2Provider, session fosite.Session) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx := fosite.NewContext()

		r, err := oauth2.NewDeviceRequest(ctx, req)
		if err != nil {
			t.Logf("Device auth request failed because: %+v", err)
			t.Logf("Request: %+v", r)
			oauth2.WriteAccessError(ctx, rw, r, err)
			return
		}

		if r.GetRequestedScopes().Has("fosite") {
			r.GrantScope("fosite")
		}

		if r.GetRequestedScopes().Has("offline") {
			r.GrantScope("offline")
		}

		if r.GetRequestedScopes().Has("openid") {
			r.GrantScope("openid")
		}

		for _, a := range r.GetRequestedAudience() {
			r.GrantAudience(a)
		}

		response, err := oauth2.NewDeviceResponse(ctx, r, session)
		if err != nil {
			t.Logf("Device auth response failed because: %+v", err)
			t.Logf("Request: %+v", r)
			oauth2.WriteAccessError(ctx, rw, r, err)
			return
		}

		oauth2.WriteDeviceResponse(ctx, rw, r, response)
	}
}
