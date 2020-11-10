package sql

import (
	"context"
	"database/sql"

	"github.com/ory/x/errorsx"

	"github.com/gobuffalo/x/randx"

	"github.com/ory/fosite/storage"

	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"

	"github.com/ory/fosite"
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/persistence"
	"github.com/ory/x/logrusx"
)

var _ persistence.Persister = new(Persister)
var _ storage.Transactional = new(Persister)

const transactionContextKey transactionContextType = "transactionConnection"

var (
	migrations = packr.New("migrations", "migrations")

	ErrTransactionOpen   = errors.New("There is already a transaction in this context.")
	ErrNoTransactionOpen = errors.New("There is no transaction in this context.")
)

type (
	Persister struct {
		conn   *pop.Connection
		mb     pop.MigrationBox
		r      Dependencies
		config configuration.Provider
		l      *logrusx.Logger
	}
	Dependencies interface {
		ClientHasher() fosite.Hasher
		KeyCipher() *jwk.AEAD
	}
	transactionContextType string
)

func (p *Persister) BeginTX(ctx context.Context) (context.Context, error) {
	_, ok := ctx.Value(transactionContextKey).(*pop.Connection)
	if ok {
		return ctx, errorsx.WithStack(ErrTransactionOpen)
	}

	tx, err := p.conn.Store.TransactionContextOptions(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	})
	c := &pop.Connection{
		TX:      tx,
		Store:   tx,
		ID:      randx.String(30),
		Dialect: p.conn.Dialect,
	}
	return context.WithValue(ctx, transactionContextKey, c), err
}

func (p *Persister) Commit(ctx context.Context) error {
	c, ok := ctx.Value(transactionContextKey).(*pop.Connection)
	if !ok || c.TX == nil {
		return errorsx.WithStack(ErrNoTransactionOpen)
	}

	return c.TX.Commit()
}

func (p *Persister) Rollback(ctx context.Context) error {
	c, ok := ctx.Value(transactionContextKey).(*pop.Connection)
	if !ok || c.TX == nil {
		return errorsx.WithStack(ErrNoTransactionOpen)
	}

	return c.TX.Rollback()
}

func NewPersister(c *pop.Connection, r Dependencies, config configuration.Provider, l *logrusx.Logger) (*Persister, error) {
	mb, err := pop.NewMigrationBox(migrations, c)
	if err != nil {
		return nil, errorsx.WithStack(err)
	}

	return &Persister{
		conn:   c,
		mb:     mb,
		r:      r,
		config: config,
		l:      l,
	}, nil
}

func (p *Persister) Connection(ctx context.Context) *pop.Connection {
	if c, ok := ctx.Value(transactionContextKey).(*pop.Connection); ok {
		return c.WithContext(ctx)
	}
	return p.conn.WithContext(ctx)
}

func (p *Persister) transaction(ctx context.Context, f func(ctx context.Context, c *pop.Connection) error) error {
	isNested := true
	c, ok := ctx.Value(transactionContextKey).(*pop.Connection)
	if !ok {
		isNested = false

		var err error
		c, err = p.conn.WithContext(ctx).NewTransaction()

		if err != nil {
			return errorsx.WithStack(err)
		}
	}

	if err := f(context.WithValue(ctx, transactionContextKey, c), c); err != nil {
		if !isNested {
			if err := c.TX.Rollback(); err != nil {
				return errorsx.WithStack(err)
			}
		}
		return err
	}

	// commit if there is no wrapping transaction
	if !isNested {
		return errorsx.WithStack(c.TX.Commit())
	}

	return nil
}
