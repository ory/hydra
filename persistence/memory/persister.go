package memory

import (
	"context"
	"io"

	"github.com/gobuffalo/pop/v5"

	"github.com/ory/hydra/persistence"
)

var _ persistence.Persister = new(Persister)

type Persister struct{}

func (p *Persister) MigrationStatus(_ context.Context, w io.Writer) error {
	return nil
}

func (p *Persister) MigrateDown(_ context.Context, steps int) error {
	return nil
}

func (p *Persister) MigrateUp(_ context.Context) error {
	return nil
}

func (p *Persister) Connection(_ context.Context) *pop.Connection {
	return nil
}
