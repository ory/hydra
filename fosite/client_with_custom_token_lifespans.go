// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import "time"

// GetEffectiveLifespan either maps GrantType x TokenType to the client's configured lifespan, or returns the fallback value.
func GetEffectiveLifespan(c Client, gt GrantType, tt TokenType, fallback time.Duration) time.Duration {
	if clc, ok := c.(ClientWithCustomTokenLifespans); ok {
		return clc.GetEffectiveLifespan(gt, tt, fallback)
	}
	return fallback
}

type ClientWithCustomTokenLifespans interface {
	// GetEffectiveLifespan either maps GrantType x TokenType to the client's configured lifespan, or returns the fallback value.
	GetEffectiveLifespan(gt GrantType, tt TokenType, fallback time.Duration) time.Duration
}

// ClientLifespanConfig holds default lifespan configuration for the different
// token types that may be issued for the client. This configuration takes
// precedence over fosite's instance-wide default lifespan, but it may be
// overridden by a session's expires_at claim.
//
// The OIDC Hybrid grant type inherits token lifespan configuration from the implicit grant.
type ClientLifespanConfig struct {
	AuthorizationCodeGrantAccessTokenLifespan  *time.Duration `json:"authorization_code_grant_access_token_lifespan"`
	AuthorizationCodeGrantIDTokenLifespan      *time.Duration `json:"authorization_code_grant_id_token_lifespan"`
	AuthorizationCodeGrantRefreshTokenLifespan *time.Duration `json:"authorization_code_grant_refresh_token_lifespan"`
	ClientCredentialsGrantAccessTokenLifespan  *time.Duration `json:"client_credentials_grant_access_token_lifespan"`
	ImplicitGrantAccessTokenLifespan           *time.Duration `json:"implicit_grant_access_token_lifespan"`
	ImplicitGrantIDTokenLifespan               *time.Duration `json:"implicit_grant_id_token_lifespan"`
	JwtBearerGrantAccessTokenLifespan          *time.Duration `json:"jwt_bearer_grant_access_token_lifespan"`
	PasswordGrantAccessTokenLifespan           *time.Duration `json:"password_grant_access_token_lifespan"`
	PasswordGrantRefreshTokenLifespan          *time.Duration `json:"password_grant_refresh_token_lifespan"`
	RefreshTokenGrantIDTokenLifespan           *time.Duration `json:"refresh_token_grant_id_token_lifespan"`
	RefreshTokenGrantAccessTokenLifespan       *time.Duration `json:"refresh_token_grant_access_token_lifespan"`
	RefreshTokenGrantRefreshTokenLifespan      *time.Duration `json:"refresh_token_grant_refresh_token_lifespan"`
	//Hybrid grant tokens are not independently configurable, see the comment above.
}

type DefaultClientWithCustomTokenLifespans struct {
	*DefaultClient
	TokenLifespans *ClientLifespanConfig `json:"token_lifespans"`
}

func (c *DefaultClientWithCustomTokenLifespans) GetTokenLifespans() *ClientLifespanConfig {
	return c.TokenLifespans
}

func (c *DefaultClientWithCustomTokenLifespans) SetTokenLifespans(lifespans *ClientLifespanConfig) {
	c.TokenLifespans = lifespans
}

// GetEffectiveLifespan either maps GrantType x TokenType to the client's configured lifespan, or returns the fallback value.
func (c *DefaultClientWithCustomTokenLifespans) GetEffectiveLifespan(gt GrantType, tt TokenType, fallback time.Duration) time.Duration {
	if c.TokenLifespans == nil {
		return fallback
	}
	var cl *time.Duration
	if gt == GrantTypeAuthorizationCode {
		if tt == AccessToken {
			cl = c.TokenLifespans.AuthorizationCodeGrantAccessTokenLifespan
		} else if tt == IDToken {
			cl = c.TokenLifespans.AuthorizationCodeGrantIDTokenLifespan
		} else if tt == RefreshToken {
			cl = c.TokenLifespans.AuthorizationCodeGrantRefreshTokenLifespan
		}
	} else if gt == GrantTypeClientCredentials {
		if tt == AccessToken {
			cl = c.TokenLifespans.ClientCredentialsGrantAccessTokenLifespan
		}
	} else if gt == GrantTypeImplicit {
		if tt == AccessToken {
			cl = c.TokenLifespans.ImplicitGrantAccessTokenLifespan
		} else if tt == IDToken {
			cl = c.TokenLifespans.ImplicitGrantIDTokenLifespan
		}
	} else if gt == GrantTypeJWTBearer {
		if tt == AccessToken {
			cl = c.TokenLifespans.JwtBearerGrantAccessTokenLifespan
		}
	} else if gt == GrantTypePassword {
		if tt == AccessToken {
			cl = c.TokenLifespans.PasswordGrantAccessTokenLifespan
		} else if tt == RefreshToken {
			cl = c.TokenLifespans.PasswordGrantRefreshTokenLifespan
		}
	} else if gt == GrantTypeRefreshToken {
		if tt == AccessToken {
			cl = c.TokenLifespans.RefreshTokenGrantAccessTokenLifespan
		} else if tt == IDToken {
			cl = c.TokenLifespans.RefreshTokenGrantIDTokenLifespan
		} else if tt == RefreshToken {
			cl = c.TokenLifespans.RefreshTokenGrantRefreshTokenLifespan
		}
	}

	if cl == nil {
		return fallback
	}
	return *cl
}
