// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package networkx

import (
	"embed"

	"github.com/pkg/errors"

	"github.com/ory/pop/v6"
	"github.com/ory/x/sqlcon"
)

// Migrations of the network manager. Apply by merging with your local migrations using
// fsx.Merge() and then passing all to the migration box.
//
//go:embed migrations/sql/*.sql
var Migrations embed.FS

func Determine(c *pop.Connection) (*Network, error) {
	var p Network
	if err := sqlcon.HandleError(c.Q().Order("created_at ASC").First(&p)); err != nil {
		if errors.Is(err, sqlcon.ErrNoRows) {
			np := NewNetwork()
			if err := c.Create(np); err != nil {
				return nil, err
			}
			return np, nil
		}
		return nil, err
	}
	return &p, nil
}
