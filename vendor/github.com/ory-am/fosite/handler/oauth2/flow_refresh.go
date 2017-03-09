package oauth2

import (
	"net/http"
	"time"

	"github.com/ory-am/fosite"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type RefreshTokenGrantHandler struct {
	AccessTokenStrategy AccessTokenStrategy

	RefreshTokenStrategy RefreshTokenStrategy

	// RefreshTokenGrantStorage is used to persist session data across requests.
	RefreshTokenGrantStorage RefreshTokenGrantStorage

	// AccessTokenLifespan defines the lifetime of an access token.
	AccessTokenLifespan time.Duration
}

// HandleTokenEndpointRequest implements https://tools.ietf.org/html/rfc6749#section-6
func (c *RefreshTokenGrantHandler) HandleTokenEndpointRequest(ctx context.Context, req *http.Request, request fosite.AccessRequester) error {
	// grant_type REQUIRED.
	// Value MUST be set to "refresh_token".
	if !request.GetGrantTypes().Exact("refresh_token") {
		return errors.WithStack(fosite.ErrUnknownRequest)
	}

	if !request.GetClient().GetGrantTypes().Has("refresh_token") {
		return errors.Wrap(fosite.ErrInvalidGrant, "The client is not allowed to use grant type refresh_token")
	}

	refresh := req.PostForm.Get("refresh_token")
	signature := c.RefreshTokenStrategy.RefreshTokenSignature(refresh)
	originalRequest, err := c.RefreshTokenGrantStorage.GetRefreshTokenSession(ctx, signature, request.GetSession())
	if errors.Cause(err) == fosite.ErrNotFound {
		return errors.Wrap(fosite.ErrInvalidRequest, err.Error())
	} else if err != nil {
		return errors.Wrap(fosite.ErrServerError, err.Error())
	}

	if !originalRequest.GetGrantedScopes().Has("offline") {
		return errors.Wrap(fosite.ErrScopeNotGranted, "The client is not allowed to use grant type refresh_token")

	}

	// The authorization server MUST ... validate the refresh token.
	if err := c.RefreshTokenStrategy.ValidateRefreshToken(ctx, request, refresh); err != nil {
		return errors.Wrap(fosite.ErrInvalidRequest, err.Error())
	}

	// The authorization server MUST ... and ensure that the refresh token was issued to the authenticated client
	if originalRequest.GetClient().GetID() != request.GetClient().GetID() {
		return errors.Wrap(fosite.ErrInvalidRequest, "Client ID mismatch")
	}

	request.SetSession(originalRequest.GetSession().Clone())
	request.SetRequestedScopes(originalRequest.GetRequestedScopes())
	for _, scope := range originalRequest.GetGrantedScopes() {
		request.GrantScope(scope)
	}

	request.GetSession().SetExpiresAt(fosite.AccessToken, time.Now().Add(c.AccessTokenLifespan))
	return nil
}

// PopulateTokenEndpointResponse implements https://tools.ietf.org/html/rfc6749#section-6
func (c *RefreshTokenGrantHandler) PopulateTokenEndpointResponse(ctx context.Context, req *http.Request, requester fosite.AccessRequester, responder fosite.AccessResponder) error {
	if !requester.GetGrantTypes().Exact("refresh_token") {
		return errors.WithStack(fosite.ErrUnknownRequest)
	}

	accessToken, accessSignature, err := c.AccessTokenStrategy.GenerateAccessToken(ctx, requester)
	if err != nil {
		return errors.Wrap(fosite.ErrServerError, err.Error())
	}

	refreshToken, refreshSignature, err := c.RefreshTokenStrategy.GenerateRefreshToken(ctx, requester)
	if err != nil {
		return errors.Wrap(fosite.ErrServerError, err.Error())
	}

	signature := c.RefreshTokenStrategy.RefreshTokenSignature(req.PostForm.Get("refresh_token"))
	if err := c.RefreshTokenGrantStorage.PersistRefreshTokenGrantSession(ctx, signature, accessSignature, refreshSignature, requester); err != nil {
		return errors.Wrap(fosite.ErrServerError, err.Error())
	}

	responder.SetAccessToken(accessToken)
	responder.SetTokenType("bearer")
	responder.SetExpiresIn(getExpiresIn(requester, fosite.AccessToken, c.AccessTokenLifespan, time.Now()))
	responder.SetScopes(requester.GetGrantedScopes())
	responder.SetExtra("refresh_token", refreshToken)
	return nil
}
