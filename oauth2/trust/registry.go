// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package trust

import (
	"github.com/ory/hydra/v2/x"
)

type InternalRegistry interface {
	x.RegistryWriter
	x.RegistryLogger
	Registry
}

type Registry interface {
	GrantManager() GrantManager
	GrantValidator() *GrantValidator
}
