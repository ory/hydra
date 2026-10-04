// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package consent

import (
	"github.com/ory/hydra/v2/aead"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/fosite/handler/openid"
	"github.com/ory/hydra/v2/internal/kratos"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/httpx"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/otelx"
)

type InternalRegistry interface {
	httpx.WriterProvider
	x.RegistryCookieStore
	logrusx.Provider
	httpx.ClientProvider
	otelx.Provider
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
