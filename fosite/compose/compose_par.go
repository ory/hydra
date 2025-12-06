// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package compose

import (
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/par"
)

// PushedAuthorizeHandlerFactory creates the basic PAR handler
func PushedAuthorizeHandlerFactory(config fosite.Configurator, storage fosite.Storage, _ interface{}) interface{} {
	return &par.PushedAuthorizeHandler{
		Storage: storage.(fosite.PARStorageProvider),
		Config:  config,
	}
}
