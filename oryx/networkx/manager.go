// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package networkx

import (
	"context"
	"embed"

	"github.com/pkg/errors"

	"github.com/ory/pop/v6"

	"github.com/ory/x/logrusx"
	"github.com/ory/x/popx"
	"github.com/ory/x/sqlcon"
)

// Migrations of the network manager. Apply by merging with your local migrations using
// fsx.Merge() and then passing all to the migration box.
//
//go:embed migrations/sql/*.sql
var Migrations embed.FS

type Manager struct {
	c *pop.Connection
	l *logrusx.Logger
}

func NewManager(
	c *pop.Connection,
	l *logrusx.Logger,
) *Manager {
	return &Manager{
		c: c,
		l: l,
	}
}

func (m *Manager) Determine(ctx context.Context) (*Network, error) {
	var p Network
	c := m.c.WithContext(ctx)
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

// MigrateUp applies pending up migrations.
//
// Deprecated: use fsx.Merge() instead to merge your local migrations with the ones exported here
func (m *Manager) MigrateUp(ctx context.Context) error {
	mm, err := popx.NewMigrationBox(Migrations, m.c.WithContext(ctx), m.l)
	if err != nil {
		return errors.WithStack(err)
	}

	return sqlcon.HandleError(mm.Up(ctx))
}
