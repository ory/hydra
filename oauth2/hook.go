package oauth2

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/ory/hydra/x"

	"github.com/ory/fosite"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/x/errorsx"
)

// AccessRequestHook is called when an access token is being refreshed.
type AccessRequestHook func(ctx context.Context, requester fosite.AccessRequester) error

// Requester is a token endpoint's request context.
//
// swagger:model oAuth2AccessRequest
type Requester struct {
	// ClientID is the identifier of the OAuth 2.0 client.
	ClientID string `json:"client_id"`
	// GrantedScopes is the list of scopes granted to the OAuth 2.0 client.
	GrantedScopes []string `json:"granted_scopes"`
	// GrantedAudience is the list of audiences granted to the OAuth 2.0 client.
	GrantedAudience []string `json:"granted_audience"`
	// GrantTypes is the requests grant types.
	GrantTypes []string `json:"grant_types"`
}

// RefreshTokenHookRequest is the request body sent to the refresh token hook.
//
// swagger:model refreshTokenHookRequest
type RefreshTokenHookRequest struct {
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

// RefreshTokenHookResponse is the response body received from the refresh token hook.
//
// swagger:model refreshTokenHookResponse
type RefreshTokenHookResponse struct {
	// Session is the session data returned by the hook.
	Session consent.ConsentRequestSessionData `json:"session"`
}

// RefreshTokenHook is an AccessRequestHook called for `refresh_token` grant type.
func RefreshTokenHook(reg interface {
	config.Provider
	x.HTTPClientProvider
}) AccessRequestHook {
	return func(ctx context.Context, requester fosite.AccessRequester) error {
		hookURL := reg.Config().TokenRefreshHookURL(ctx)
		if hookURL == nil {
			return nil
		}

		if !requester.GetGrantTypes().ExactOne("refresh_token") {
			return nil
		}

		session, ok := requester.GetSession().(*Session)
		if !ok {
			return nil
		}

		requesterInfo := Requester{
			ClientID:        requester.GetClient().GetID(),
			GrantedScopes:   requester.GetGrantedScopes(),
			GrantedAudience: requester.GetGrantedAudience(),
			GrantTypes:      requester.GetGrantTypes(),
		}

		reqBody := RefreshTokenHookRequest{
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
					WithDebug("refresh token hook: marshal request body"),
			)
		}

		req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodPost, hookURL.String(), bytes.NewReader(reqBodyBytes))
		if err != nil {
			return errorsx.WithStack(
				fosite.ErrServerError.
					WithWrap(err).
					WithDebug("refresh token hook: new http request"),
			)
		}
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")

		resp, err := reg.HTTPClient(ctx).Do(req)
		if err != nil {
			return errorsx.WithStack(
				fosite.ErrServerError.
					WithWrap(err).
					WithDebug("refresh token hook: do http request"),
			)
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusOK:
			// Token refresh permitted with new session data
		case http.StatusNoContent:
			// Token refresh is permitted without overriding session data
			return nil
		case http.StatusForbidden:
			return errorsx.WithStack(
				fosite.ErrAccessDenied.
					WithDebugf("refresh token hook: %s", resp.Status),
			)
		default:
			return errorsx.WithStack(
				fosite.ErrServerError.
					WithDebugf("refresh token hook: %s", resp.Status),
			)
		}

		var respBody RefreshTokenHookResponse
		if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
			return errorsx.WithStack(
				fosite.ErrServerError.
					WithWrap(err).
					WithDebugf("refresh token hook: unmarshal response body"),
			)
		}

		// Overwrite existing session data (extra claims).
		session.Extra = respBody.Session.AccessToken
		idTokenClaims := session.IDTokenClaims()
		idTokenClaims.Extra = respBody.Session.IDToken

		return nil
	}
}
