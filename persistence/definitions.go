package persistence

import (
	"context"
	"github.com/gobuffalo/pop/v5"
	"io"
)

type (
	Persister interface {
		MigrationStatus(context.Context, io.Writer) error
		MigrateDown(context.Context, int) error
		MigrateUp(context.Context) error
		Connection(context.Context) *pop.Connection
	}
	Provider interface {
		Persister() Persister
	}
)
