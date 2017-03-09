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

type OpenIDConnectHybridHandler struct {
	AuthorizeImplicitGrantTypeHandler *oauth2.AuthorizeImplicitGrantTypeHandler
	AuthorizeExplicitGrantHandler     *oauth2.AuthorizeExplicitGrantHandler
	IDTokenHandleHelper               *IDTokenHandleHelper
	ScopeStrategy                     fosite.ScopeStrategy

	Enigma *jwt.RS256JWTStrategy
}

func (c *OpenIDConnectHybridHandler) HandleAuthorizeEndpointRequest(ctx context.Context, req *http.Request, ar fosite.AuthorizeRequester, resp fosite.AuthorizeResponder) error {
	if len(ar.GetResponseTypes()) < 2 {
		return nil
	}

	if !(ar.GetResponseTypes().Matches("token", "id_token", "code") || ar.GetResponseTypes().Matches("token", "code") || ar.GetResponseTypes().Matches("id_token", "code")) {
		return nil
	}

	if ar.GetResponseTypes().Matches("token") && !ar.GetClient().GetResponseTypes().Has("token") {
		return errors.Wrap(fosite.ErrInvalidGrant, "The client is not allowed to use the token response type")
	} else if ar.GetResponseTypes().Matches("code") && !ar.GetClient().GetResponseTypes().Has("code") {
		return errors.Wrap(fosite.ErrInvalidGrant, "The client is not allowed to use the code response type")
	} else if ar.GetResponseTypes().Matches("id_token") && !ar.GetClient().GetResponseTypes().Has("id_token") {
		return errors.Wrap(fosite.ErrInvalidGrant, "The client is not allowed to use the id_token response type")
	}

	sess, ok := ar.GetSession().(Session)
	if !ok {
		return errors.WithStack(ErrInvalidSession)
	}

	client := ar.GetClient()
	for _, scope := range ar.GetRequestedScopes() {
		if !c.ScopeStrategy(client.GetScopes(), scope) {
			return errors.Wrap(fosite.ErrInvalidScope, fmt.Sprintf("The client is not allowed to request scope %s", scope))
		}
	}

	claims := sess.IDTokenClaims()
	if ar.GetResponseTypes().Has("code") {
		if !ar.GetClient().GetGrantTypes().Has("authorization_code") {
			return errors.Wrap(fosite.ErrInvalidGrant, "The client is not allowed to use the authorization_code grant type")
		}

		code, signature, err := c.AuthorizeExplicitGrantHandler.AuthorizeCodeStrategy.GenerateAuthorizeCode(ctx, ar)
		if err != nil {
			return errors.Wrap(fosite.ErrServerError, err.Error())
		} else if err := c.AuthorizeExplicitGrantHandler.AuthorizeCodeGrantStorage.CreateAuthorizeCodeSession(ctx, signature, ar); err != nil {
			return errors.Wrap(fosite.ErrServerError, err.Error())
		}

		resp.AddFragment("code", code)
		ar.SetResponseTypeHandled("code")

		hash, err := c.Enigma.Hash([]byte(resp.GetFragment().Get("code")))
		if err != nil {
			return err
		}
		claims.CodeHash = base64.RawURLEncoding.EncodeToString([]byte(hash[:c.Enigma.GetSigningMethodLength()/2]))
	}

	if ar.GetResponseTypes().Has("token") {
		if !ar.GetClient().GetGrantTypes().Has("implicit") {
			return errors.Wrap(fosite.ErrInvalidGrant, "The client is not allowed to use the implicit grant type")
		} else if err := c.AuthorizeImplicitGrantTypeHandler.IssueImplicitAccessToken(ctx, req, ar, resp); err != nil {
			return errors.Wrap(err, err.Error())
		}
		ar.SetResponseTypeHandled("token")

		hash, err := c.Enigma.Hash([]byte(resp.GetFragment().Get("access_token")))
		if err != nil {
			return err
		}
		claims.AccessTokenHash = base64.RawURLEncoding.EncodeToString([]byte(hash[:c.Enigma.GetSigningMethodLength()/2]))
	}

	if resp.GetFragment().Get("state") == "" {
		resp.AddFragment("state", ar.GetState())
	}

	if !ar.GetGrantedScopes().Has("openid") || !ar.GetResponseTypes().Has("id_token") {
		ar.SetResponseTypeHandled("id_token")
		return nil
	}

	if err := c.IDTokenHandleHelper.IssueImplicitIDToken(ctx, req, ar, resp); err != nil {
		return errors.Wrap(err, err.Error())
	}

	ar.SetResponseTypeHandled("id_token")
	return nil
	// there is no need to check for https, because implicit flow does not require https
	// https://tools.ietf.org/html/rfc6819#section-4.4.2
}
