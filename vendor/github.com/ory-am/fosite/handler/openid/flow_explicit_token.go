package openid

import (
	"net/http"

	"github.com/ory-am/fosite"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func (c *OpenIDConnectExplicitHandler) HandleTokenEndpointRequest(ctx context.Context, r *http.Request, request fosite.AccessRequester) error {
	return fosite.ErrUnknownRequest
}

func (c *OpenIDConnectExplicitHandler) PopulateTokenEndpointResponse(ctx context.Context, req *http.Request, requester fosite.AccessRequester, responder fosite.AccessResponder) error {
	if !requester.GetGrantTypes().Exact("authorization_code") {
		return errors.WithStack(fosite.ErrUnknownRequest)
	}

	authorize, err := c.OpenIDConnectRequestStorage.GetOpenIDConnectSession(ctx, requester.GetRequestForm().Get("code"), requester)
	if errors.Cause(err) == ErrNoSessionFound {
		return errors.Wrap(fosite.ErrUnknownRequest, err.Error())
	} else if err != nil {
		return errors.Wrap(fosite.ErrServerError, err.Error())
	}

	if !authorize.GetGrantedScopes().Has("openid") {
		return errors.Wrap(fosite.ErrMisconfiguration, "The an openid connect session was found but the openid scope is missing in it")
	}

	if !requester.GetClient().GetGrantTypes().Has("authorization_code") {
		return errors.Wrap(fosite.ErrInvalidGrant, "The client is not allowed to use the authorization_code grant type")
	}

	if !requester.GetClient().GetResponseTypes().Has("id_token") {
		return errors.Wrap(fosite.ErrInvalidGrant, "The client is not allowed to use response type id_token")
	}

	return c.IssueExplicitIDToken(ctx, req, authorize, responder)
}
