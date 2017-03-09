package openid

import (
	"net/http"

	"github.com/ory-am/fosite"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type OpenIDConnectExplicitHandler struct {
	// OpenIDConnectRequestStorage is the storage for open id connect sessions.
	OpenIDConnectRequestStorage OpenIDConnectRequestStorage

	*IDTokenHandleHelper
}

func (c *OpenIDConnectExplicitHandler) HandleAuthorizeEndpointRequest(ctx context.Context, req *http.Request, ar fosite.AuthorizeRequester, resp fosite.AuthorizeResponder) error {
	if !(ar.GetGrantedScopes().Has("openid") && ar.GetResponseTypes().Exact("code")) {
		return nil
	}

	if !ar.GetClient().GetResponseTypes().Has("id_token", "code") {
		return errors.Wrap(fosite.ErrInvalidRequest, "The client is not allowed to use response type id_token and code")
	}

	if len(resp.GetCode()) == 0 {
		return errors.Wrap(fosite.ErrMisconfiguration, "Authorization code has not been issued yet")
	}

	if err := c.OpenIDConnectRequestStorage.CreateOpenIDConnectSession(ctx, resp.GetCode(), ar); err != nil {
		return errors.Wrap(fosite.ErrServerError, err.Error())
	}

	// there is no need to check for https, because it has already been checked by core.explicit

	return nil
}
