// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import (
	"context"
	"reflect"
)

const MinParameterEntropy = 8

var defaultResponseModeHandler = &DefaultResponseModeHandler{}

// AuthorizeEndpointHandlers is a list of AuthorizeEndpointHandler
type AuthorizeEndpointHandlers []AuthorizeEndpointHandler

// Append adds an AuthorizeEndpointHandler to this list. Ignores duplicates based on reflect.TypeOf.
func (a *AuthorizeEndpointHandlers) Append(h AuthorizeEndpointHandler) {
	for _, this := range *a {
		if reflect.TypeOf(this) == reflect.TypeOf(h) {
			return
		}
	}

	*a = append(*a, h)
}

// TokenEndpointHandlers is a list of TokenEndpointHandler
type TokenEndpointHandlers []TokenEndpointHandler

// Append adds an TokenEndpointHandler to this list. Ignores duplicates based on reflect.TypeOf.
func (t *TokenEndpointHandlers) Append(h TokenEndpointHandler) {
	for _, this := range *t {
		if reflect.TypeOf(this) == reflect.TypeOf(h) {
			return
		}
	}

	*t = append(*t, h)
}

// TokenIntrospectionHandlers is a list of TokenValidator
type TokenIntrospectionHandlers []TokenIntrospector

// Append adds an AccessTokenValidator to this list. Ignores duplicates based on reflect.TypeOf.
func (t *TokenIntrospectionHandlers) Append(h TokenIntrospector) {
	for _, this := range *t {
		if reflect.TypeOf(this) == reflect.TypeOf(h) {
			return
		}
	}

	*t = append(*t, h)
}

// RevocationHandlers is a list of RevocationHandler
type RevocationHandlers []RevocationHandler

// Append adds an RevocationHandler to this list. Ignores duplicates based on reflect.TypeOf.
func (t *RevocationHandlers) Append(h RevocationHandler) {
	for _, this := range *t {
		if reflect.TypeOf(this) == reflect.TypeOf(h) {
			return
		}
	}

	*t = append(*t, h)
}

// PushedAuthorizeEndpointHandlers is a list of PushedAuthorizeEndpointHandler
type PushedAuthorizeEndpointHandlers []PushedAuthorizeEndpointHandler

// Append adds an AuthorizeEndpointHandler to this list. Ignores duplicates based on reflect.TypeOf.
func (a *PushedAuthorizeEndpointHandlers) Append(h PushedAuthorizeEndpointHandler) {
	for _, this := range *a {
		if reflect.TypeOf(this) == reflect.TypeOf(h) {
			return
		}
	}

	*a = append(*a, h)
}

// DeviceEndpointHandlers is a list of DeviceEndpointHandler
type DeviceEndpointHandlers []DeviceEndpointHandler

// Append adds an DeviceEndpointHandlers to this list. Ignores duplicates based on reflect.TypeOf.
func (a *DeviceEndpointHandlers) Append(h DeviceEndpointHandler) {
	for _, this := range *a {
		if reflect.TypeOf(this) == reflect.TypeOf(h) {
			return
		}
	}

	*a = append(*a, h)
}

var _ OAuth2Provider = (*Fosite)(nil)

type Configurator interface {
	IDTokenIssuerProvider
	IDTokenLifespanProvider
	AllowedPromptsProvider
	EnforcePKCEProvider
	EnforcePKCEForPublicClientsProvider
	EnablePKCEPlainChallengeMethodProvider
	GrantTypeJWTBearerCanSkipClientAuthProvider
	GrantTypeJWTBearerIDOptionalProvider
	GrantTypeJWTBearerIssuedDateOptionalProvider
	GetJWTMaxDurationProvider
	AudienceStrategyProvider
	ScopeStrategyProvider
	RedirectSecureCheckerProvider
	OmitRedirectScopeParamProvider
	SanitationAllowedProvider
	JWTScopeFieldProvider
	AccessTokenIssuerProvider
	DisableRefreshTokenValidationProvider
	RefreshTokenScopesProvider
	AccessTokenLifespanProvider
	RefreshTokenLifespanProvider
	VerifiableCredentialsNonceLifespanProvider
	AuthorizeCodeLifespanProvider
	DeviceAndUserCodeLifespanProvider
	TokenEntropyProvider
	RotatedGlobalSecretsProvider
	GlobalSecretProvider
	JWKSFetcherStrategyProvider
	HTTPClientProvider
	ScopeStrategyProvider
	AudienceStrategyProvider
	MinParameterEntropyProvider
	HMACHashingProvider
	ClientAuthenticationStrategyProvider
	ResponseModeHandlerExtensionProvider
	SendDebugMessagesToClientsProvider
	JWKSFetcherStrategyProvider
	ClientAuthenticationStrategyProvider
	ResponseModeHandlerExtensionProvider
	MessageCatalogProvider
	FormPostHTMLTemplateProvider
	TokenURLProvider
	GetSecretsHashingProvider
	AuthorizeEndpointHandlersProvider
	TokenEndpointHandlersProvider
	TokenIntrospectionHandlersProvider
	RevocationHandlersProvider
	UseLegacyErrorFormatProvider
	DeviceEndpointHandlersProvider
	UserCodeProvider
	DeviceProvider
}

func NewOAuth2Provider(s Storage, c Configurator) *Fosite {
	return &Fosite{Store: s, Config: c}
}

// Fosite implements OAuth2Provider.
type Fosite struct {
	Store Storage

	Config Configurator
}

// GetMinParameterEntropy returns MinParameterEntropy if set. Defaults to fosite.MinParameterEntropy.
func (f *Fosite) GetMinParameterEntropy(ctx context.Context) int {
	if mp := f.Config.GetMinParameterEntropy(ctx); mp > 0 {
		return mp
	}

	return MinParameterEntropy
}

func (f *Fosite) ResponseModeHandler(ctx context.Context) ResponseModeHandler {
	if ext := f.Config.GetResponseModeHandlerExtension(ctx); ext != nil {
		return ext
	}
	return defaultResponseModeHandler
}
