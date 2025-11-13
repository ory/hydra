// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package consent

import (
	"github.com/ory/hydra/v2/aead"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/fosite/handler/openid"
	"github.com/ory/hydra/v2/internal/kratos"
	"github.com/ory/hydra/v2/x"
)

type InternalRegistry interface {
	x.RegistryWriter
	x.RegistryCookieStore
	x.RegistryLogger
	x.HTTPClientProvider
	x.TracingProvider
	x.NetworkProvider
	kratos.Provider
	Registry
	client.Registry

	FlowCipher() *aead.XChaCha20Poly1305
	OAuth2Storage() x.FositeStorer
	OpenIDConnectRequestValidator() *openid.OpenIDConnectRequestValidator
}

type Registry interface {
	ManagerProvider
	ObfuscatedSubjectManagerProvider
	LoginManagerProvider
	LogoutManagerProvider

	ConsentStrategy() Strategy
}
