package oauth2

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/go-cleanhttp"

	"github.com/ory/fosite"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/x/errorsx"
)

// AccessRequestHook is called when an access token is being refreshed.
type AccessRequestHook func(ctx context.Context, requester fosite.AccessRequester) error

// RefreshTokenHookRequest is the request body sent to the refresh token hook.
//
// swagger:model refreshTokenHookRequest
type RefreshTokenHookRequest struct {
	// Subject is the identifier of the authenticated end-user.
	Subject string `json:"subject"`
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
func RefreshTokenHook(config *config.Provider) AccessRequestHook {
	client := cleanhttp.DefaultPooledClient()

	return func(ctx context.Context, requester fosite.AccessRequester) error {
		hookURL := config.TokenRefreshHookURL()
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

		reqBody := RefreshTokenHookRequest{
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

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, hookURL.String(), bytes.NewReader(reqBodyBytes))
		if err != nil {
			return errorsx.WithStack(
				fosite.ErrServerError.
					WithWrap(err).
					WithDebug("refresh token hook: new http request"),
			)
		}
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")

		resp, err := client.Do(req)
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
			// We only accept '200 OK' here. Any other status code is considered an error.
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
