package sql

import (
	"context"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/pop/v5"
	"github.com/ory/hydra/persistence"
	"github.com/pkg/errors"
	"io"
)

var _ persistence.Persister = new(Persister)
var migrations = packr.New("migrations", "migrations")

type Persister struct {
	c  *pop.Connection
	mb pop.MigrationBox
}

func NewPersister(c *pop.Connection) (*Persister, error) {
	mb, err := pop.NewMigrationBox(migrations, c)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Persister{
		c,
		mb,
	}, nil
}

func (p *Persister) MigrationStatus(_ context.Context, w io.Writer) error {
	return errors.WithStack(p.mb.Status(w))
}

func (p *Persister) MigrateDown(_ context.Context, steps int) error {
	return errors.WithStack(p.mb.Down(steps))
}

func (p *Persister) MigrateUp(_ context.Context) error {
	return errors.WithStack(p.mb.Up())
}

func (p *Persister) Connection(_ context.Context) *pop.Connection {
	return p.c
}
