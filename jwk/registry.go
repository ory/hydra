// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"github.com/ory/x/httpx"
	"github.com/ory/x/logrusx"

	"github.com/ory/hydra/v2/driver/config"
)

type InternalRegistry interface {
	httpx.WriterProvider
	logrusx.Provider
	Registry
}

type Registry interface {
	config.Provider
	ManagerProvider
}
