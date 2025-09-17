// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"context"
	"database/sql"
	"reflect"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/ory/fosite"
	"github.com/ory/fosite/storage"
	"github.com/ory/hydra/v2/aead"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/internal/kratos"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/persistence"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/pop/v6"
	"github.com/ory/x/contextx"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/networkx"
	"github.com/ory/x/otelx"
	"github.com/ory/x/popx"
)

var _ persistence.Persister = (*Persister)(nil)
var _ storage.Transactional = (*Persister)(nil)

var (
	ErrNoTransactionOpen = errors.New("There is no Transaction in this context.")
)

type skipCommitContextKey int

const skipCommitKey skipCommitContextKey = 0

type (
	Persister struct {
		*BasePersister
		r Dependencies
		l *logrusx.Logger
	}
	Dependencies interface {
		ClientHasher() fosite.Hasher
		KeyCipher() *aead.AESGCM
		FlowCipher() *aead.XChaCha20Poly1305
		Kratos() kratos.Client
		contextx.Provider
		x.RegistryLogger
		x.TracingProvider
		config.Provider
	}
	BasePersister struct {
		c           *pop.Connection
		fallbackNID uuid.UUID
		d           baseDependencies
	}
	baseDependencies interface {
		x.RegistryLogger
		x.TracingProvider
		contextx.Provider
		config.Provider
		jwk.ManagerProvider
	}
)

func (p *BasePersister) BeginTX(ctx context.Context) (_ context.Context, err error) {
	ctx, span := p.d.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.BeginTX")
	defer otelx.End(span, &err)

	fallback := &pop.Connection{TX: &pop.Tx{}}
	if popx.GetConnection(ctx, fallback).TX != fallback.TX {
		return context.WithValue(ctx, skipCommitKey, true), nil // no-op
	}

	tx, err := p.c.Store.TransactionContextOptions(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
	})
	c := &pop.Connection{
		TX:      tx,
		Store:   tx,
		ID:      uuid.Must(uuid.NewV4()).String(),
		Dialect: p.c.Dialect,
	}
	return popx.WithTransaction(ctx, c), err
}

func (p *BasePersister) Commit(ctx context.Context) (err error) {
	ctx, span := p.d.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.Commit")
	defer otelx.End(span, &err)

	if skip, ok := ctx.Value(skipCommitKey).(bool); ok && skip {
		return nil // we skipped BeginTX, so we also skip Commit
	}

	fallback := &pop.Connection{TX: &pop.Tx{}}
	tx := popx.GetConnection(ctx, fallback)
	if tx.TX == fallback.TX || tx.TX == nil {
		return errors.WithStack(ErrNoTransactionOpen)
	}

	return errors.WithStack(tx.TX.Commit())
}

func (p *BasePersister) Rollback(ctx context.Context) (err error) {
	ctx, span := p.d.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.Rollback")
	defer otelx.End(span, &err)

	if skip, ok := ctx.Value(skipCommitKey).(bool); ok && skip {
		return nil // we skipped BeginTX, so we also skip Rollback
	}

	fallback := &pop.Connection{TX: &pop.Tx{}}
	tx := popx.GetConnection(ctx, fallback)
	if tx.TX == fallback.TX || tx.TX == nil {
		return errors.WithStack(ErrNoTransactionOpen)
	}

	return errors.WithStack(tx.TX.Rollback())
}

func NewPersister(base *BasePersister, r Dependencies) *Persister {
	return &Persister{
		BasePersister: base,
		r:             r,
		l:             r.Logger(),
	}
}

func NewBasePersister(c *pop.Connection, d baseDependencies) *BasePersister {
	return &BasePersister{c: c, d: d}
}

func (p *BasePersister) DetermineNetwork(ctx context.Context) (*networkx.Network, error) {
	return networkx.Determine(p.Connection(ctx))
}

func (p BasePersister) WithFallbackNetworkID(nid uuid.UUID) *BasePersister {
	p.fallbackNID = nid
	return &p
}

func (p *BasePersister) CreateWithNetwork(ctx context.Context, v interface{}) error {
	p.mustSetNetwork(ctx, v)
	return p.Connection(ctx).Create(v)
}

func (p *BasePersister) UpdateWithNetwork(ctx context.Context, v interface{}) (int64, error) {
	p.mustSetNetwork(ctx, v)

	m := pop.NewModel(v, ctx)
	cols := m.Columns()
	cs := make([]string, 0, len(cols.Cols))
	for _, t := range m.Columns().Cols {
		cs = append(cs, t.Name)
	}

	return p.Connection(ctx).Where(m.IDField()+" = ? AND nid = ?", m.ID(), p.NetworkID(ctx)).UpdateQuery(v, cs...)
}

func (p *BasePersister) NetworkID(ctx context.Context) uuid.UUID {
	return p.d.Contextualizer().Network(ctx, p.fallbackNID)
}

func (p *BasePersister) QueryWithNetwork(ctx context.Context) *pop.Query {
	return p.Connection(ctx).Where("nid = ?", p.NetworkID(ctx))
}

func (p *BasePersister) Connection(ctx context.Context) *pop.Connection {
	return popx.GetConnection(ctx, p.c)
}

func (p *BasePersister) Ping(ctx context.Context) error { return p.c.Store.SQLDB().PingContext(ctx) }

func (p *BasePersister) mustSetNetwork(ctx context.Context, v interface{}) {
	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Ptr || (rv.Kind() == reflect.Ptr && rv.Elem().Kind() != reflect.Struct) {
		panic("v must be a pointer to a struct")
	}
	nf := rv.Elem().FieldByName("NID")
	if !nf.IsValid() || !nf.CanSet() {
		panic("v must have settable a field 'NID uuid.UUID'")
	}
	nf.Set(reflect.ValueOf(p.NetworkID(ctx)))
}

func (p *BasePersister) Transaction(ctx context.Context, f func(ctx context.Context, c *pop.Connection) error) error {
	return popx.Transaction(ctx, p.c, f)
}
