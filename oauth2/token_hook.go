// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/ory/hydra/v2/flow"
	"github.com/ory/hydra/v2/x"

	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/fosite"
)

// AccessRequestHook is called when an access token request is performed.
type AccessRequestHook func(ctx context.Context, requester fosite.AccessRequester) error

// Request is a token endpoint's request context.
//
// swagger:ignore
type Request struct {
	// ClientID is the identifier of the OAuth 2.0 client.
	ClientID string `json:"client_id"`
	// RequestedScopes is the list of scopes requested to the OAuth 2.0 client.
	RequestedScopes []string `json:"requested_scopes"`
	// GrantedScopes is the list of scopes granted to the OAuth 2.0 client.
	GrantedScopes []string `json:"granted_scopes"`
	// GrantedAudience is the list of audiences granted to the OAuth 2.0 client.
	GrantedAudience []string `json:"granted_audience"`
	// GrantTypes is the requests grant types.
	GrantTypes []string `json:"grant_types"`
	// Payload is the requests payload.
	Payload map[string][]string `json:"payload"`
}

// TokenHookRequest is the request body sent to the token hook.
//
// swagger:ignore
type TokenHookRequest struct {
	// Session is the request's session..
	Session *Session `json:"session"`
	// Requester is a token endpoint's request context.
	Request Request `json:"request"`
}

// TokenHookResponse is the response body received from the token hook.
//
// swagger:ignore
type TokenHookResponse struct {
	// Session is the session data returned by the hook.
	Session flow.AcceptOAuth2ConsentRequestSession `json:"session"`
}

type APIKeyAuthConfig struct {
	In    string `json:"in"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

func applyAuth(req *retryablehttp.Request, auth *config.Auth) error {
	if auth == nil {
		return nil
	}

	switch auth.Type {
	case "api_key":
		switch auth.Config.In {
		case "header":
			req.Header.Set(auth.Config.Name, auth.Config.Value)
		case "cookie":
			req.AddCookie(&http.Cookie{Name: auth.Config.Name, Value: auth.Config.Value})
		}
	default:
		return errors.Errorf("unsupported auth type %q", auth.Type)
	}
	return nil
}

func executeHookAndUpdateSession(ctx context.Context, reg x.HTTPClientProvider, hookConfig *config.HookConfig, reqBodyBytes []byte, session *Session) error {
	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodPost, hookConfig.URL, bytes.NewReader(reqBodyBytes))
	if err != nil {
		return errors.WithStack(
			fosite.ErrServerError.
				WithWrap(err).
				WithDescription("An error occurred while preparing the token hook.").
				WithDebugf("Unable to prepare the HTTP Request: %s", err),
		)
	}
	if err := applyAuth(req, hookConfig.Auth); err != nil {
		return errors.WithStack(
			fosite.ErrServerError.
				WithWrap(err).
				WithDescription("An error occurred while applying the token hook authentication.").
				WithDebugf("Unable to apply the token hook authentication: %s", err))
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := reg.HTTPClient(ctx).Do(req)
	if err != nil {
		return errors.WithStack(
			fosite.ErrServerError.
				WithWrap(err).
				WithDescription("An error occurred while executing the token hook.").
				WithDebugf("Unable to execute HTTP Request: %s", err),
		)
	}
	defer resp.Body.Close() //nolint:errcheck
	resp.Body = io.NopCloser(io.LimitReader(resp.Body, 5<<20 /* 5 MiB */))

	switch resp.StatusCode {
	case http.StatusOK:
		// Token permitted with new session data
	case http.StatusNoContent:
		// Token is permitted without overriding session data
		return nil
	case http.StatusForbidden:
		return errors.WithStack(
			fosite.ErrAccessDenied.
				WithDescription("The token hook target responded with an error.").
				WithDebugf("Token hook responded with HTTP status code: %s", resp.Status),
		)
	default:
		return errors.WithStack(
			fosite.ErrServerError.
				WithDescription("The token hook target responded with an error.").
				WithDebugf("Token hook responded with HTTP status code: %s", resp.Status),
		)
	}

	var respBody TokenHookResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return errors.WithStack(
			fosite.ErrServerError.
				WithWrap(err).
				WithDescription("The token hook target responded with an error.").
				WithDebugf("Response from token hook could not be decoded: %s", err),
		)
	}

	// Overwrite existing session data (extra claims).
	session.Extra = respBody.Session.AccessToken
	idTokenClaims := session.IDTokenClaims()
	idTokenClaims.Extra = respBody.Session.IDToken
	return nil
}

// TokenHook is an AccessRequestHook called for all grant types.
func TokenHook(reg interface {
	config.Provider
	x.HTTPClientProvider
}) AccessRequestHook {
	return func(ctx context.Context, requester fosite.AccessRequester) error {
		hookConfig := reg.Config().TokenHookConfig(ctx)
		if hookConfig == nil {
			return nil
		}

		session, ok := requester.GetSession().(*Session)
		if !ok {
			return nil
		}

		request := Request{
			ClientID:        requester.GetClient().GetID(),
			RequestedScopes: requester.GetRequestedScopes(),
			GrantedScopes:   requester.GetGrantedScopes(),
			GrantedAudience: requester.GetGrantedAudience(),
			GrantTypes:      requester.GetGrantTypes(),
			Payload:         requester.Sanitize([]string{"assertion"}).GetRequestForm(),
		}

		reqBody := TokenHookRequest{
			Session: session,
			Request: request,
		}

		reqBodyBytes, err := json.Marshal(&reqBody)
		if err != nil {
			return errors.WithStack(
				fosite.ErrServerError.
					WithWrap(err).
					WithDescription("An error occurred while encoding the token hook.").
					WithDebugf("Unable to encode the token hook body: %s", err),
			)
		}

		err = executeHookAndUpdateSession(ctx, reg, hookConfig, reqBodyBytes, session)
		if err != nil {
			return err
		}

		return nil
	}
}
