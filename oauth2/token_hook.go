// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/ory/hydra/v2/flow"
	"github.com/ory/hydra/v2/x"

	"github.com/ory/fosite"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/x/errorsx"
)

// AccessRequestHook is called when an access token request is performed.
type AccessRequestHook func(ctx context.Context, requester fosite.AccessRequester) error

// Request is a token endpoint's request context.
//
// swagger:ignore
type Request struct {
	// ClientID is the identifier of the OAuth 2.0 client.
	ClientID string `json:"client_id"`
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

func executeHookAndUpdateSession(ctx context.Context, reg x.HTTPClientProvider, hookURL *url.URL, reqBodyBytes []byte, session *Session) error {
	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodPost, hookURL.String(), bytes.NewReader(reqBodyBytes))
	if err != nil {
		return errorsx.WithStack(
			fosite.ErrServerError.
				WithWrap(err).
				WithDescription("An error occurred while preparing the token hook.").
				WithDebugf("Unable to prepare the HTTP Request: %s", err),
		)
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := reg.HTTPClient(ctx).Do(req)
	if err != nil {
		return errorsx.WithStack(
			fosite.ErrServerError.
				WithWrap(err).
				WithDescription("An error occurred while executing the token hook.").
				WithDebugf("Unable to execute HTTP Request: %s", err),
		)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// Token permitted with new session data
	case http.StatusNoContent:
		// Token is permitted without overriding session data
		return nil
	case http.StatusForbidden:
		return errorsx.WithStack(
			fosite.ErrAccessDenied.
				WithDescription("The token hook target responded with an error.").
				WithDebugf("Token hook responded with HTTP status code: %s", resp.Status),
		)
	default:
		return errorsx.WithStack(
			fosite.ErrServerError.
				WithDescription("The token hook target responded with an error.").
				WithDebugf("Token hook responded with HTTP status code: %s", resp.Status),
		)
	}

	var respBody TokenHookResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return errorsx.WithStack(
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
		hookURL := reg.Config().TokenHookURL(ctx)
		if hookURL == nil {
			return nil
		}

		session, ok := requester.GetSession().(*Session)
		if !ok {
			return nil
		}

		request := Request{
			ClientID:        requester.GetClient().GetID(),
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
			return errorsx.WithStack(
				fosite.ErrServerError.
					WithWrap(err).
					WithDescription("An error occurred while encoding the token hook.").
					WithDebugf("Unable to encode the token hook body: %s", err),
			)
		}

		err = executeHookAndUpdateSession(ctx, reg, hookURL, reqBodyBytes, session)
		if err != nil {
			return err
		}

		return nil
	}
}
