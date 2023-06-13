// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"github.com/ory/hydra/v2/aead"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/x"
)

type InternalRegistry interface {
	x.RegistryWriter
	x.RegistryLogger
	Registry
}

type Registry interface {
	config.Provider
	KeyManager() Manager
	SoftwareKeyManager() Manager
	KeyCipher() *aead.AESGCM
}
