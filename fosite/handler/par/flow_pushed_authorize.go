// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package par

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"time"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/token/hmac"
	"github.com/ory/x/errorsx"
)

const (
	defaultPARKeyLength = 32
)

var b64 = base64.URLEncoding.WithPadding(base64.NoPadding)

// PushedAuthorizeHandler handles the PAR request
type PushedAuthorizeHandler struct {
	Storage fosite.PARStorageProvider
	Config  fosite.Configurator
}

// HandlePushedAuthorizeEndpointRequest handles a pushed authorize endpoint request. To extend the handler's capabilities, the http request
// is passed along, if further information retrieval is required. If the handler feels that he is not responsible for
// the pushed authorize request, he must return nil and NOT modify session nor responder neither requester.
func (c *PushedAuthorizeHandler) HandlePushedAuthorizeEndpointRequest(ctx context.Context, ar fosite.AuthorizeRequester, resp fosite.PushedAuthorizeResponder) error {
	configProvider, ok := c.Config.(fosite.PushedAuthorizeRequestConfigProvider)
	if !ok {
		return errorsx.WithStack(fosite.ErrServerError.WithHint(fosite.ErrorPARNotSupported).WithDebug(fosite.DebugPARConfigMissing))
	}

	if !ar.GetResponseTypes().HasOneOf("token", "code", "id_token") {
		return nil
	}

	if !c.secureChecker(ctx, ar.GetRedirectURI()) {
		return errorsx.WithStack(fosite.ErrInvalidRequest.WithHint("Redirect URL is using an insecure protocol, http is only allowed for hosts with suffix 'localhost', for example: http://myapp.localhost/."))
	}

	client := ar.GetClient()
	for _, scope := range ar.GetRequestedScopes() {
		if !c.Config.GetScopeStrategy(ctx)(client.GetScopes(), scope) {
			return errorsx.WithStack(fosite.ErrInvalidScope.WithHintf("The OAuth 2.0 Client is not allowed to request scope '%s'.", scope))
		}
	}

	if err := c.Config.GetAudienceStrategy(ctx)(client.GetAudience(), ar.GetRequestedAudience()); err != nil {
		return err
	}

	expiresIn := configProvider.GetPushedAuthorizeContextLifespan(ctx)
	if ar.GetSession() != nil {
		ar.GetSession().SetExpiresAt(fosite.PushedAuthorizeRequestContext, time.Now().UTC().Add(expiresIn))
	}

	// generate an ID
	stateKey, err := hmac.RandomBytes(defaultPARKeyLength)
	if err != nil {
		return errorsx.WithStack(fosite.ErrInsufficientEntropy.WithHint("Unable to generate the random part of the request_uri.").WithWrap(err).WithDebug(err.Error()))
	}

	requestURI := fmt.Sprintf("%s%s", configProvider.GetPushedAuthorizeRequestURIPrefix(ctx), b64.EncodeToString(stateKey))

	// store
	if err = c.Storage.PARStorage().CreatePARSession(ctx, requestURI, ar); err != nil {
		return errorsx.WithStack(fosite.ErrServerError.WithHint("Unable to store the PAR session").WithWrap(err).WithDebug(err.Error()))
	}

	resp.SetRequestURI(requestURI)
	resp.SetExpiresIn(int(expiresIn.Seconds()))
	return nil
}

func (c *PushedAuthorizeHandler) secureChecker(ctx context.Context, u *url.URL) bool {
	isRedirectURISecure := c.Config.GetRedirectSecureChecker(ctx)
	if isRedirectURISecure == nil {
		isRedirectURISecure = fosite.IsRedirectURISecure
	}
	return isRedirectURISecure(ctx, u)
}
