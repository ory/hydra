package persistence

import (
	"context"
	"io"

	"github.com/gobuffalo/pop/v5"
)

type (
	Persister interface {
		MigrationStatus(context.Context, io.Writer) error
		MigrateDown(context.Context, int) error
		MigrateUp(context.Context) error
		PrepareMigration(context.Context) error
		Connection(context.Context) *pop.Connection
	}
	Provider interface {
		Persister() Persister
	}
)
