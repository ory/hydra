package memory

import (
	"context"
	"io"

	"github.com/gobuffalo/pop/v5"

	"github.com/ory/hydra/persistence"
)

var _ persistence.Persister = new(Persister)

type Persister struct{}

func (*Persister) MigrationStatus(_ context.Context, _ io.Writer) error {
	return nil
}

func (*Persister) MigrateDown(_ context.Context, steps int) error {
	return nil
}

func (*Persister) MigrateUp(_ context.Context) error {
	return nil
}

func (*Persister) PrepareMigration(context.Context) error {
	return nil
}

func (*Persister) Connection(_ context.Context) *pop.Connection {
	return nil
}
