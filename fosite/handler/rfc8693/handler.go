// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package rfc8693

import (
	"context"
	"time"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/oauth2"
	"github.com/ory/x/errorsx"
)

// RFC 8693 token type identifiers (section 3).
const (
	TokenTypeAccessToken = "urn:ietf:params:oauth:token-type:access_token"
	TokenTypeJWT         = "urn:ietf:params:oauth:token-type:jwt"
)

var _ fosite.TokenEndpointHandler = (*Handler)(nil)

type Handler struct {
	Storage  oauth2.AccessTokenStorageProvider
	Strategy oauth2.AccessTokenStrategyProvider
	Config   interface {
		fosite.AccessTokenLifespanProvider
		fosite.ScopeStrategyProvider
		fosite.AudienceStrategyProvider
		fosite.GrantTypeTokenExchangeCanSkipClientAuthProvider
	}
}

// CanHandleTokenEndpointRequest returns true when grant_type is token-exchange (RFC 8693).
func (c *Handler) CanHandleTokenEndpointRequest(ctx context.Context, requester fosite.AccessRequester) bool {
	return requester.GetGrantTypes().ExactOne(string(fosite.GrantTypeTokenExchange))
}

// CanSkipClientAuth returns whether client authentication can be skipped for token exchange.
func (c *Handler) CanSkipClientAuth(ctx context.Context, requester fosite.AccessRequester) bool {
	return c.Config.GetGrantTypeTokenExchangeCanSkipClientAuth(ctx)
}

// CheckRequest validates that the request is a token exchange and the client is allowed.
func (c *Handler) CheckRequest(ctx context.Context, request fosite.AccessRequester) error {
	if !c.CanHandleTokenEndpointRequest(ctx, request) {
		return errorsx.WithStack(fosite.ErrUnknownRequest)
	}
	if !c.CanSkipClientAuth(ctx, request) && request.GetClient() != nil && !request.GetClient().GetGrantTypes().Has(string(fosite.GrantTypeTokenExchange)) {
		return errorsx.WithStack(fosite.ErrUnauthorizedClient.WithHintf("The OAuth 2.0 Client is not allowed to use authorization grant \"%s\".", fosite.GrantTypeTokenExchange))
	}
	return nil
}

// HandleTokenEndpointRequest implements RFC 8693 token exchange request handling.
func (c *Handler) HandleTokenEndpointRequest(ctx context.Context, request fosite.AccessRequester) error {
	if err := c.CheckRequest(ctx, request); err != nil {
		return err
	}

	form := request.GetRequestForm()
	subjectToken := form.Get("subject_token")
	subjectTokenType := form.Get("subject_token_type")

	if subjectToken == "" {
		return errorsx.WithStack(fosite.ErrInvalidRequest.WithHint("The subject_token request parameter must be set when using grant_type token-exchange."))
	}
	if subjectTokenType == "" {
		return errorsx.WithStack(fosite.ErrInvalidRequest.WithHint("The subject_token_type request parameter must be set when using grant_type token-exchange."))
	}

	// Support access_token and jwt types; for our server both resolve via access token introspection.
	if subjectTokenType != TokenTypeAccessToken && subjectTokenType != TokenTypeJWT {
		return errorsx.WithStack(fosite.ErrInvalidRequest.WithHintf("Unsupported subject_token_type: %s.", subjectTokenType))
	}

	// Resolve subject token via access token storage (introspection path).
	sig := c.Strategy.AccessTokenStrategy().AccessTokenSignature(ctx, subjectToken)
	subjectRequester, err := c.Storage.AccessTokenStorage().GetAccessTokenSession(ctx, sig, request.GetSession())
	if err != nil {
		return errorsx.WithStack(fosite.ErrInvalidGrant.WithHint("The subject_token is invalid, expired, or revoked.").WithWrap(err).WithDebug(err.Error()))
	}
	if err := c.Strategy.AccessTokenStrategy().ValidateAccessToken(ctx, subjectRequester, subjectToken); err != nil {
		return err
	}

	// Copy subject and granted rights from subject token; restrict to requested scope/audience.
	request.GetSession().SetExpiresAt(fosite.AccessToken, subjectRequester.GetSession().GetExpiresAt(fosite.AccessToken))
	if sub := subjectRequester.GetSession().GetSubject(); sub != "" {
		if s, ok := request.GetSession().(interface{ SetSubject(string) }); ok {
			s.SetSubject(sub)
		}
	}

	requestedScopes := request.GetRequestedScopes()
	requestedAudience := request.GetRequestedAudience()
	subjectGrantedScopes := subjectRequester.GetGrantedScopes()
	subjectGrantedAudience := subjectRequester.GetGrantedAudience()

	if len(requestedScopes) == 0 {
		for _, s := range subjectGrantedScopes {
			request.GrantScope(s)
		}
	} else {
		for _, scope := range requestedScopes {
			if c.Config.GetScopeStrategy(ctx)(subjectGrantedScopes, scope) {
				request.GrantScope(scope)
			}
		}
	}

	if len(requestedAudience) == 0 {
		for _, a := range subjectGrantedAudience {
			request.GrantAudience(a)
		}
	} else {
		for _, aud := range requestedAudience {
			if fosite.DefaultAudienceMatchingStrategy(subjectGrantedAudience, []string{aud}) == nil {
				request.GrantAudience(aud)
			}
		}
	}

	atLifespan := fosite.GetEffectiveLifespan(request.GetClient(), fosite.GrantTypeTokenExchange, fosite.AccessToken, c.Config.GetAccessTokenLifespan(ctx))
	request.GetSession().SetExpiresAt(fosite.AccessToken, time.Now().UTC().Add(atLifespan).Round(time.Second))

	return nil
}

// PopulateTokenEndpointResponse issues the access token and sets issued_token_type (RFC 8693).
func (c *Handler) PopulateTokenEndpointResponse(ctx context.Context, request fosite.AccessRequester, response fosite.AccessResponder) error {
	if err := c.CheckRequest(ctx, request); err != nil {
		return err
	}

	atLifespan := fosite.GetEffectiveLifespan(request.GetClient(), fosite.GrantTypeTokenExchange, fosite.AccessToken, c.Config.GetAccessTokenLifespan(ctx))
	_, err := c.issueAccessToken(ctx, atLifespan, request, response)
	return err
}

func (c *Handler) issueAccessToken(ctx context.Context, atLifespan time.Duration, requester fosite.AccessRequester, responder fosite.AccessResponder) (signature string, err error) {
	token, signature, err := c.Strategy.AccessTokenStrategy().GenerateAccessToken(ctx, requester)
	if err != nil {
		return "", err
	}
	if err := c.Storage.AccessTokenStorage().CreateAccessTokenSession(ctx, signature, requester.Sanitize([]string{})); err != nil {
		return "", err
	}

	if !requester.GetSession().GetExpiresAt(fosite.AccessToken).IsZero() {
		atLifespan = time.Duration(requester.GetSession().GetExpiresAt(fosite.AccessToken).UnixNano() - time.Now().UTC().UnixNano())
	}

	responder.SetAccessToken(token)
	responder.SetTokenType("bearer")
	responder.SetExpiresIn(atLifespan)
	responder.SetScopes(requester.GetGrantedScopes())
	// RFC 8693 section 2.2.1: issued_token_type is REQUIRED in the response.
	responder.SetExtra("issued_token_type", TokenTypeAccessToken)

	return signature, nil
}
