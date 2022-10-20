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
}

// RefreshTokenHookRequest is the request body sent to the refresh token hook.
//
// swagger:ignore
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
// swagger:ignore
type RefreshTokenHookResponse struct {
	// Session is the session data returned by the hook.
	Session consent.AcceptOAuth2ConsentRequestSession `json:"session"`
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
					WithDescription("An error occurred while encoding the refresh token hook.").
					WithDebugf("Unable to encode the refresh token hook body: %s", err),
			)
		}

		req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodPost, hookURL.String(), bytes.NewReader(reqBodyBytes))
		if err != nil {
			return errorsx.WithStack(
				fosite.ErrServerError.
					WithWrap(err).
					WithDescription("An error occurred while preparing the refresh token hook.").
					WithDebugf("Unable to prepare the HTTP Request: %s", err),
			)
		}
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")

		resp, err := reg.HTTPClient(ctx).Do(req)
		if err != nil {
			return errorsx.WithStack(
				fosite.ErrServerError.
					WithWrap(err).
					WithDescription("An error occurred while executing the refresh token hook.").
					WithDebugf("Unable to execute HTTP Request: %s", err),
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
					WithDescription("The refresh token hook target responded with an error.").
					WithDebugf("Refresh token hook responded with HTTP status code: %s", resp.Status),
			)
		default:
			return errorsx.WithStack(
				fosite.ErrServerError.
					WithDescription("The refresh token hook target responded with an error.").
					WithDebugf("Refresh token hook responded with HTTP status code: %s", resp.Status),
			)
		}

		var respBody RefreshTokenHookResponse
		if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
			return errorsx.WithStack(
				fosite.ErrServerError.
					WithWrap(err).
					WithDescription("The refresh token hook target responded with an error.").
					WithDebugf("Response from refresh token hook could not be decoded: %s", err),
			)
		}

		// Overwrite existing session data (extra claims).
		session.Extra = respBody.Session.AccessToken
		idTokenClaims := session.IDTokenClaims()
		idTokenClaims.Extra = respBody.Session.IDToken
		return nil
	}
}
