// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"github.com/ory/hydra/v2/driver"
)

type Handler struct {
	Migration *MigrateHandler
	Janitor   *JanitorHandler
}

func NewHandler(dOpts []driver.OptionsModifier) *Handler {
	return &Handler{
		Migration: newMigrateHandler(dOpts),
		Janitor:   newJanitorHandler(dOpts),
	}
}
