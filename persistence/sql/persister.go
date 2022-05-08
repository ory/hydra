package sql

import (
	"context"
	"database/sql"
	"reflect"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/gobuffalo/x/randx"
	"github.com/pkg/errors"

	"github.com/ory/fosite"
	"github.com/ory/fosite/storage"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/persistence"
	"github.com/ory/hydra/x"
	"github.com/ory/hydra/x/contextx"
	"github.com/ory/x/errorsx"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/networkx"
	"github.com/ory/x/popx"
)

var _ persistence.Persister = new(Persister)
var _ storage.Transactional = new(Persister)

const transactionContextKey transactionContextType = "transactionConnection"

var (
	ErrTransactionOpen   = errors.New("There is already a transaction in this context.")
	ErrNoTransactionOpen = errors.New("There is no transaction in this context.")
)

type (
	Persister struct {
		conn        *pop.Connection
		mb          *popx.MigrationBox
		r           Dependencies
		config      *config.Provider
		l           *logrusx.Logger
		fallbackNID uuid.UUID
		p           *networkx.Manager
	}
	Dependencies interface {
		ClientHasher() fosite.Hasher
		KeyCipher() *jwk.AEAD
		KeyGenerators() map[string]jwk.KeyGenerator
		contextx.ContextualizerProvider
		x.RegistryLogger
		x.TracingProvider
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

	return errorsx.WithStack(c.TX.Commit())
}

func (p *Persister) Rollback(ctx context.Context) error {
	c, ok := ctx.Value(transactionContextKey).(*pop.Connection)
	if !ok || c.TX == nil {
		return errorsx.WithStack(ErrNoTransactionOpen)
	}

	return errorsx.WithStack(c.TX.Rollback())
}

func NewPersister(ctx context.Context, c *pop.Connection, r Dependencies, config *config.Provider, l *logrusx.Logger) (*Persister, error) {
	mb, err := popx.NewMigrationBox(migrations, popx.NewMigrator(c, r.Logger(), r.Tracer(ctx), 0))
	if err != nil {
		return nil, errorsx.WithStack(err)
	}

	return &Persister{
		conn:   c,
		mb:     mb,
		r:      r,
		config: config,
		l:      l,
		p:      networkx.NewManager(c, r.Logger(), r.Tracer(ctx)),
	}, nil
}

func (p *Persister) DetermineNetwork(ctx context.Context) (*networkx.Network, error) {
	return p.p.Determine(ctx)
}

func (p Persister) WithFallbackNetworkID(nid uuid.UUID) persistence.Persister {
	p.fallbackNID = nid
	return &p
}

func (p *Persister) CreateWithNetwork(ctx context.Context, v interface{}) error {
	n := p.NetworkID(ctx)
	return p.Connection(ctx).Create(p.mustSetNetwork(n, v))
}

func (p *Persister) UpdateWithNetwork(ctx context.Context, v interface{}) (int64, error) {
	n := p.NetworkID(ctx)
	v = p.mustSetNetwork(n, v)

	m := pop.NewModel(v, ctx)
	var cs []string
	for _, t := range m.Columns().Cols {
		cs = append(cs, t.Name)
	}

	return p.Connection(ctx).Where(m.IDField()+" = ? AND nid = ?", m.ID(), n).UpdateQuery(v, cs...)
}

func (p *Persister) NetworkID(ctx context.Context) uuid.UUID {
	return p.r.Contextualizer().Network(ctx, p.fallbackNID)
}

func (p *Persister) QueryWithNetwork(ctx context.Context) *pop.Query {
	return p.Connection(ctx).Where("nid = ?", p.NetworkID(ctx))
}

func (p *Persister) Connection(ctx context.Context) *pop.Connection {
	if c, ok := ctx.Value(transactionContextKey).(*pop.Connection); ok {
		return c.WithContext(ctx)
	}
	return p.conn.WithContext(ctx)
}

func (p *Persister) mustSetNetwork(nid uuid.UUID, v interface{}) interface{} {
	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Ptr || (rv.Kind() == reflect.Ptr && rv.Elem().Kind() != reflect.Struct) {
		panic("v must be a pointer to a struct")
	}
	nf := rv.Elem().FieldByName("NID")
	if !nf.IsValid() || !nf.CanSet() {
		panic("v must have settable a field 'NID uuid.UUID'")
	}
	nf.Set(reflect.ValueOf(nid))
	return v
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
