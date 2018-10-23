/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 *
 */

package compose

import (
	"crypto/rsa"

	"github.com/ory/fosite"
	"github.com/ory/fosite/token/jwt"
)

type Factory func(config *Config, storage interface{}, strategy interface{}) interface{}

// Compose takes a config, a storage, a strategy and handlers to instantiate an OAuth2Provider:
//
//  import "github.com/ory/fosite/compose"
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
func Compose(config *Config, storage interface{}, strategy interface{}, hasher fosite.Hasher, factories ...Factory) fosite.OAuth2Provider {
	if hasher == nil {
		hasher = &fosite.BCrypt{WorkFactor: config.GetHashCost()}
	}

	f := &fosite.Fosite{
		Store:                      storage.(fosite.Storage),
		AuthorizeEndpointHandlers:  fosite.AuthorizeEndpointHandlers{},
		TokenEndpointHandlers:      fosite.TokenEndpointHandlers{},
		TokenIntrospectionHandlers: fosite.TokenIntrospectionHandlers{},
		RevocationHandlers:         fosite.RevocationHandlers{},
		Hasher:                     hasher,
		ScopeStrategy:              config.GetScopeStrategy(),
		SendDebugMessagesToClients: config.SendDebugMessagesToClients,
		TokenURL:                   config.TokenURL,
		JWKSFetcherStrategy:        config.GetJWKSFetcherStrategy(),
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
			CoreStrategy:               NewOAuth2HMACStrategy(config, secret, nil),
			OpenIDConnectTokenStrategy: NewOpenIDConnectStrategy(config, key),
			JWTStrategy: &jwt.RS256JWTStrategy{
				PrivateKey: key,
			},
		},
		nil,

		OAuth2AuthorizeExplicitFactory,
		OAuth2AuthorizeImplicitFactory,
		OAuth2ClientCredentialsGrantFactory,
		OAuth2RefreshTokenGrantFactory,
		OAuth2ResourceOwnerPasswordCredentialsFactory,

		OAuth2PKCEFactory,

		OpenIDConnectExplicitFactory,
		OpenIDConnectImplicitFactory,
		OpenIDConnectHybridFactory,
		OpenIDConnectRefreshFactory,

		OAuth2TokenIntrospectionFactory,
	)
}
