package compose

import (
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/handler/oauth2"
	"github.com/ory-am/fosite/handler/openid"
)

// OpenIDConnectExplicitFactory creates an OpenID Connect explicit ("authorize code flow") grant handler. You must add this handler
// *after* you have added an OAuth2 authorize code handler!
func OpenIDConnectExplicitFactory(config *Config, storage interface{}, strategy interface{}) interface{} {
	return &openid.OpenIDConnectExplicitHandler{
		OpenIDConnectRequestStorage: storage.(openid.OpenIDConnectRequestStorage),
		IDTokenHandleHelper: &openid.IDTokenHandleHelper{
			IDTokenStrategy: strategy.(openid.OpenIDConnectTokenStrategy),
		},
	}
}

// OpenIDConnectImplicitFactory creates an OpenID Connect implicit ("implicit flow") grant handler. You must add this handler
// *after* you have added an OAuth2 authorize implicit handler!
func OpenIDConnectImplicitFactory(config *Config, storage interface{}, strategy interface{}) interface{} {
	return &openid.OpenIDConnectImplicitHandler{
		AuthorizeImplicitGrantTypeHandler: &oauth2.AuthorizeImplicitGrantTypeHandler{
			AccessTokenStrategy: strategy.(oauth2.AccessTokenStrategy),
			AccessTokenStorage:  storage.(oauth2.AccessTokenStorage),
			AccessTokenLifespan: config.GetAccessTokenLifespan(),
		},
		ScopeStrategy: fosite.HierarchicScopeStrategy,
		IDTokenHandleHelper: &openid.IDTokenHandleHelper{
			IDTokenStrategy: strategy.(openid.OpenIDConnectTokenStrategy),
		},
	}
}

// OpenIDConnectHybridFactory creates an OpenID Connect hybrid grant handler. You must add this handler
// *after* you have added an OAuth2 authorize code and implicit authorize handler!
func OpenIDConnectHybridFactory(config *Config, storage interface{}, strategy interface{}) interface{} {
	return &openid.OpenIDConnectHybridHandler{
		AuthorizeExplicitGrantHandler: &oauth2.AuthorizeExplicitGrantHandler{
			AccessTokenStrategy:       strategy.(oauth2.AccessTokenStrategy),
			RefreshTokenStrategy:      strategy.(oauth2.RefreshTokenStrategy),
			AuthorizeCodeStrategy:     strategy.(oauth2.AuthorizeCodeStrategy),
			AuthorizeCodeGrantStorage: storage.(oauth2.AuthorizeCodeGrantStorage),
			AuthCodeLifespan:          config.GetAuthorizeCodeLifespan(),
			AccessTokenLifespan:       config.GetAccessTokenLifespan(),
		},
		ScopeStrategy: fosite.HierarchicScopeStrategy,
		AuthorizeImplicitGrantTypeHandler: &oauth2.AuthorizeImplicitGrantTypeHandler{
			AccessTokenStrategy: strategy.(oauth2.AccessTokenStrategy),
			AccessTokenStorage:  storage.(oauth2.AccessTokenStorage),
			AccessTokenLifespan: config.GetAccessTokenLifespan(),
		},
		IDTokenHandleHelper: &openid.IDTokenHandleHelper{
			IDTokenStrategy: strategy.(openid.OpenIDConnectTokenStrategy),
		},
	}
}
