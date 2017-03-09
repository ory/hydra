package compose

import (
	"crypto/rsa"

	"github.com/ory-am/fosite"
)

type Factory func(config *Config, storage interface{}, strategy interface{}) interface{}

// Compose takes a config, a storage, a strategy and handlers to instantiate an OAuth2Provider:
//
//  import "github.com/ory-am/fosite/compose"
//
//  // var storage = new(MyFositeStorage)
//  var config = Config {
//  	AccessTokenLifespan: time.Minute * 30,
// 	// check Config for further configuration options
//  }
//
//  var strategy = NewOAuth2HMACStrategy(config)
//
//  var oauth2Provider = Compose(
//  	config,
// 	storage,
// 	strategy,
//	NewOAuth2AuthorizeExplicitHandler,
//	OAuth2ClientCredentialsGrantFactory,
// 	// for a complete list refer to the docs of this package
//  )
//
// Compose makes use of interface{} types in order to be able to handle a all types of stores, strategies and handlers.
func Compose(config *Config, storage interface{}, strategy interface{}, factories ...Factory) fosite.OAuth2Provider {
	f := &fosite.Fosite{
		Store: storage.(fosite.Storage),
		AuthorizeEndpointHandlers:  fosite.AuthorizeEndpointHandlers{},
		TokenEndpointHandlers:      fosite.TokenEndpointHandlers{},
		TokenIntrospectionHandlers: fosite.TokenIntrospectionHandlers{},
		RevocationHandlers:         fosite.RevocationHandlers{},
		Hasher:                     &fosite.BCrypt{WorkFactor: config.GetHashCost()},
		ScopeStrategy:              fosite.HierarchicScopeStrategy,
	}

	for _, factory := range factories {
		res := factory(config, storage, strategy)
		if ah, ok := res.(fosite.AuthorizeEndpointHandler); ok {
			f.AuthorizeEndpointHandlers.Append(ah)
		}
		if th, ok := res.(fosite.TokenEndpointHandler); ok {
			f.TokenEndpointHandlers.Append(th)
		}
		if tv, ok := res.(fosite.TokenIntrospector); ok {
			f.TokenIntrospectionHandlers.Append(tv)
		}
		if rh, ok := res.(fosite.RevocationHandler); ok {
			f.RevocationHandlers.Append(rh)
		}
	}

	return f
}

// ComposeAllEnabled returns a fosite instance with all OAuth2 and OpenID Connect handlers enabled.
func ComposeAllEnabled(config *Config, storage interface{}, secret []byte, key *rsa.PrivateKey) fosite.OAuth2Provider {
	return Compose(
		config,
		storage,
		&CommonStrategy{
			CoreStrategy:               NewOAuth2HMACStrategy(config, secret),
			OpenIDConnectTokenStrategy: NewOpenIDConnectStrategy(key),
		},
		OAuth2AuthorizeExplicitFactory,
		OAuth2AuthorizeImplicitFactory,
		OAuth2ClientCredentialsGrantFactory,
		OAuth2RefreshTokenGrantFactory,
		OAuth2ResourceOwnerPasswordCredentialsFactory,

		OpenIDConnectExplicitFactory,
		OpenIDConnectImplicitFactory,
		OpenIDConnectHybridFactory,

		OAuth2TokenIntrospectionFactory,
	)
}
