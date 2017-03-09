package openid

import (
	"net/http"

	"fmt"

	"encoding/base64"
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/handler/oauth2"
	"github.com/ory-am/fosite/token/jwt"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type OpenIDConnectImplicitHandler struct {
	AuthorizeImplicitGrantTypeHandler *oauth2.AuthorizeImplicitGrantTypeHandler
	*IDTokenHandleHelper
	ScopeStrategy fosite.ScopeStrategy

	RS256JWTStrategy *jwt.RS256JWTStrategy
}

func (c *OpenIDConnectImplicitHandler) HandleAuthorizeEndpointRequest(ctx context.Context, req *http.Request, ar fosite.AuthorizeRequester, resp fosite.AuthorizeResponder) error {
	if !(ar.GetGrantedScopes().Has("openid") && (ar.GetResponseTypes().Has("token", "id_token") || ar.GetResponseTypes().Exact("id_token"))) {
		return nil
	} else if ar.GetResponseTypes().Has("code") {
		// hybrid flow
		return nil
	}

	if !ar.GetClient().GetGrantTypes().Has("implicit") {
		return errors.Wrap(fosite.ErrInvalidGrant, "The client is not allowed to use implicit grant type")
	}

	if ar.GetResponseTypes().Exact("id_token") && !ar.GetClient().GetResponseTypes().Has("id_token") {
		return errors.Wrap(fosite.ErrInvalidGrant, "The client is not allowed to use response type id_token")
	} else if ar.GetResponseTypes().Matches("token", "id_token") && !ar.GetClient().GetResponseTypes().Has("token", "id_token") {
		return errors.Wrap(fosite.ErrInvalidGrant, "The client is not allowed to use response type token and id_token")
	}

	client := ar.GetClient()
	for _, scope := range ar.GetRequestedScopes() {
		if !c.ScopeStrategy(client.GetScopes(), scope) {
			return errors.Wrap(fosite.ErrInvalidScope, fmt.Sprintf("The client is not allowed to request scope %s", scope))
		}
	}

	sess, ok := ar.GetSession().(Session)
	if !ok {
		return errors.WithStack(ErrInvalidSession)
	}

	claims := sess.IDTokenClaims()
	if ar.GetResponseTypes().Has("token") {
		if err := c.AuthorizeImplicitGrantTypeHandler.IssueImplicitAccessToken(ctx, req, ar, resp); err != nil {
			return errors.Wrap(err, err.Error())
		}

		ar.SetResponseTypeHandled("token")
		hash, err := c.RS256JWTStrategy.Hash([]byte(resp.GetFragment().Get("access_token")))
		if err != nil {
			return err
		}

		claims.AccessTokenHash = base64.RawURLEncoding.EncodeToString([]byte(hash[:c.RS256JWTStrategy.GetSigningMethodLength()/2]))
	} else {
		resp.AddFragment("state", ar.GetState())
	}

	if err := c.IssueImplicitIDToken(ctx, req, ar, resp); err != nil {
		return errors.Wrap(err, err.Error())
	}

	// there is no need to check for https, because implicit flow does not require https
	// https://tools.ietf.org/html/rfc6819#section-4.4.2

	ar.SetResponseTypeHandled("id_token")
	return nil
}
