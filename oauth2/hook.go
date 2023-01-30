// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/ory/hydra/v2/x"

	"github.com/ory/fosite"
	"github.com/ory/hydra/v2/consent"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/x/errorsx"
)

// AccessRequestHook is called when an access token is being refreshed.
type AccessRequestHook func(ctx context.Context, requester fosite.AccessRequester) error

// Requester is a token endpoint's request context.
//
// swagger:ignore
type Requester struct {
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

// TokenHookRequest is the request body sent to token hooks.
//
// swagger:ignore
type TokenHookRequest struct {
	// Subject is the identifier of the authenticated end-user.
	Subject string `json:"subject"`
	// Session is the request's session..
	Session *Session `json:"session"`
	// Requester is a token endpoint's request context.
	Requester Requester `json:"requester"`
	// ClientID is the identifier of the OAuth 2.0 client.
	ClientID string `json:"client_id"`
	// GrantedScopes is the list of scopes granted to the OAuth 2.0 client.
	GrantedScopes []string `json:"granted_scopes"`
	// GrantedAudience is the list of audiences granted to the OAuth 2.0 client.
	GrantedAudience []string `json:"granted_audience"`
}

// TokenHookResponse is the response body received from token hooks.
//
// swagger:ignore
type TokenHookResponse struct {
	// Session is the session data returned by the hook.
	Session consent.AcceptOAuth2ConsentRequestSession `json:"session"`
}

// RefreshTokenHook is an AccessRequestHook called for `refresh_token` grant type.
func RefreshTokenHook(reg interface {
	config.Provider
	x.HTTPClientProvider
},
) AccessRequestHook {
	return func(ctx context.Context, requester fosite.AccessRequester) error {
		hookURL := reg.Config().TokenRefreshHookURL(ctx)
		if hookURL == nil {
			return nil
		}
		return callHook(ctx, reg, requester, "refresh_token", hookURL)
	}
}

// AuthorizationCodeHook is an AccessRequestHook called for `authorization_code` grant type.
func AuthorizationCodeHook(reg interface {
	config.Provider
	x.HTTPClientProvider
},
) AccessRequestHook {
	return func(ctx context.Context, requester fosite.AccessRequester) error {
		hookURL := reg.Config().AuthorizationCodeHookURL(ctx)
		if hookURL == nil {
			return nil
		}
		return callHook(ctx, reg, requester, "authorization_code", hookURL)
	}
}

// ClientCredentialsHook is an AccessRequestHook called for `client_credentials` grant type.
func ClientCredentialsHook(reg interface {
	config.Provider
	x.HTTPClientProvider
},
) AccessRequestHook {
	return func(ctx context.Context, requester fosite.AccessRequester) error {
		hookURL := reg.Config().ClientCredentialsHookURL(ctx)
		if hookURL == nil {
			return nil
		}
		return callHook(ctx, reg, requester, "client_credentials", hookURL)
	}
}

// JWTBearerHook is an AccessRequestHook called for `urn:ietf:params:oauth:grant-type:jwt-bearer` grant type.
func JWTBearerHook(reg interface {
	config.Provider
	x.HTTPClientProvider
},
) AccessRequestHook {
	return func(ctx context.Context, requester fosite.AccessRequester) error {
		hookURL := reg.Config().JWTBearerRefreshHookURL(ctx)
		if hookURL == nil {
			return nil
		}
		return callHook(ctx, reg, requester, "urn:ietf:params:oauth:grant-type:jwt-bearer", hookURL)
	}
}

func callHook(ctx context.Context, reg x.HTTPClientProvider, requester fosite.AccessRequester, hookType string, hookURL *url.URL) error {
	if !requester.GetGrantTypes().ExactOne(hookType) {
		return nil
	}

	session, ok := requester.GetSession().(*Session)
	if !ok {
		return nil
	}

	payload := map[string][]string{}

	if requester.GetGrantTypes().ExactOne("urn:ietf:params:oauth:grant-type:jwt-bearer") || requester.GetGrantTypes().ExactOne("client_credentials") {
		payload = requester.GetRequestForm()
	}

	requesterInfo := Requester{
		ClientID:        requester.GetClient().GetID(),
		GrantedScopes:   requester.GetGrantedScopes(),
		GrantedAudience: requester.GetGrantedAudience(),
		GrantTypes:      requester.GetGrantTypes(),
		Payload:         payload,
	}

	reqBody := TokenHookRequest{
		Session:         session,
		Requester:       requesterInfo,
		Subject:         session.GetSubject(),
		ClientID:        requester.GetClient().GetID(),
		GrantedScopes:   requester.GetGrantedScopes(),
		GrantedAudience: requester.GetGrantedAudience(),
	}
	reqBodyBytes, err := json.Marshal(&reqBody)
	if err != nil {
		return errorsx.WithStack(
			fosite.ErrServerError.
				WithWrap(err).
				WithDescription(fmt.Sprintf("An error occurred while encoding the %s hook.", hookType)).
				WithDebugf("Unable to encode the %v hook body: %s", hookType, err),
		)
	}

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodPost, hookURL.String(), bytes.NewReader(reqBodyBytes))
	if err != nil {
		return errorsx.WithStack(
			fosite.ErrServerError.
				WithWrap(err).
				WithDescription(fmt.Sprintf("An error occurred while preparing the %s hook.", hookType)).
				WithDebugf("Unable to prepare the HTTP Request: %s", err),
		)
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := reg.HTTPClient(ctx).Do(req)
	if err != nil {
		return errorsx.WithStack(
			fosite.ErrServerError.
				WithWrap(err).
				WithDescription(fmt.Sprintf("An error occurred while executing the %s hook.", hookType)).
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
				WithDescription(fmt.Sprintf("The %s hook target responded with an error.", hookType)).
				WithDebugf(fmt.Sprintf("%s hook responded with HTTP status code: %s", hookType, resp.Status)),
		)
	default:
		return errorsx.WithStack(
			fosite.ErrServerError.
				WithDescription(fmt.Sprintf("The %s hook target responded with an error.", hookType)).
				WithDebugf(fmt.Sprintf("%s hook responded with HTTP status code: %s", hookType, resp.Status)),
		)
	}

	var respBody TokenHookResponse
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return errorsx.WithStack(
			fosite.ErrServerError.
				WithWrap(err).
				WithDescription(fmt.Sprintf("The %s hook target responded with an error.", hookType)).
				WithDebugf(fmt.Sprintf("Response from %s hook could not be decoded: %s", hookType, err)),
		)
	}

	// Overwrite existing session data (extra claims).
	session.Extra = respBody.Session.AccessToken
	idTokenClaims := session.IDTokenClaims()
	idTokenClaims.Extra = respBody.Session.IDToken
	return nil
}
