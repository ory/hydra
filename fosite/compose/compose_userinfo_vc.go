// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package compose

import (
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/verifiable"
)

// OIDCUserinfoVerifiableCredentialFactory creates a verifiable credentials
// handler.
func OIDCUserinfoVerifiableCredentialFactory(config fosite.Configurator, storage fosite.Storage, strategy any) any {
	return &verifiable.Handler{
		NonceManagerProvider: storage.(verifiable.NonceManagerProvider),
		Config:               config,
	}
}
