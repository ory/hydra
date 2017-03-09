package oauth2

import (
	"github.com/ory-am/fosite"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type CoreValidator struct {
	CoreStrategy
	CoreStorage
	ScopeStrategy fosite.ScopeStrategy
}

func (c *CoreValidator) IntrospectToken(ctx context.Context, token string, tokenType fosite.TokenType, accessRequest fosite.AccessRequester, scopes []string) (err error) {
	switch tokenType {
	case fosite.RefreshToken:
		if err = c.introspectRefreshToken(ctx, token, accessRequest); err == nil {
			return err
		} else if err = c.introspectAuthorizeCode(ctx, token, accessRequest); err == nil {
			return err
		} else if err = c.introspectAccessToken(ctx, token, accessRequest, scopes); err == nil {
			return err
		}
		return err
	case fosite.AuthorizeCode:
		if err = c.introspectAuthorizeCode(ctx, token, accessRequest); err == nil {
			return err
		} else if err := c.introspectAccessToken(ctx, token, accessRequest, scopes); err == nil {
			return err
		} else if err := c.introspectRefreshToken(ctx, token, accessRequest); err == nil {
			return err
		}
		return err
	}
	if err = c.introspectAccessToken(ctx, token, accessRequest, scopes); err == nil {
		return err
	} else if err := c.introspectRefreshToken(ctx, token, accessRequest); err == nil {
		return err
	} else if err := c.introspectAuthorizeCode(ctx, token, accessRequest); err == nil {
		return err
	}
	return err
}

func (c *CoreValidator) introspectAccessToken(ctx context.Context, token string, accessRequest fosite.AccessRequester, scopes []string) error {
	sig := c.CoreStrategy.AccessTokenSignature(token)
	or, err := c.CoreStorage.GetAccessTokenSession(ctx, sig, accessRequest.GetSession())
	if err != nil {
		return errors.Wrap(fosite.ErrRequestUnauthorized, err.Error())
	} else if err := c.CoreStrategy.ValidateAccessToken(ctx, or, token); err != nil {
		return err
	}

	for _, scope := range scopes {
		if scope == "" {
			continue
		}

		if !c.ScopeStrategy(or.GetGrantedScopes(), scope) {
			return errors.WithStack(fosite.ErrInvalidScope)
		}
	}

	accessRequest.Merge(or)
	return nil
}

func (c *CoreValidator) introspectRefreshToken(ctx context.Context, token string, accessRequest fosite.AccessRequester) error {
	sig := c.CoreStrategy.RefreshTokenSignature(token)
	if or, err := c.CoreStorage.GetRefreshTokenSession(ctx, sig, accessRequest.GetSession()); err != nil {
		return errors.Wrap(fosite.ErrRequestUnauthorized, err.Error())
	} else if err := c.CoreStrategy.ValidateRefreshToken(ctx, or, token); err != nil {
		return err
	} else {
		accessRequest.Merge(or)
	}

	return nil
}

func (c *CoreValidator) introspectAuthorizeCode(ctx context.Context, token string, accessRequest fosite.AccessRequester) error {
	sig := c.CoreStrategy.AuthorizeCodeSignature(token)
	if or, err := c.CoreStorage.GetAuthorizeCodeSession(ctx, sig, accessRequest.GetSession()); err != nil {
		return errors.Wrap(err, fosite.ErrRequestUnauthorized.Error())
	} else if err := c.CoreStrategy.ValidateAuthorizeCode(ctx, or, token); err != nil {
		return err
	} else {
		accessRequest.Merge(or)
	}

	return nil
}
