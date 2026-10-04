// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"context"
	"reflect"

	"github.com/gofrs/uuid"

	"github.com/ory/hydra/v2/aead"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/internal/kratos"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/persistence"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/pop/v6"
	"github.com/ory/x/contextx"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/networkx"
	"github.com/ory/x/otelx"
	"github.com/ory/x/popx"
)

var (
	_ persistence.Persister     = (*Persister)(nil)
	_ fosite.ClientManager      = (*Persister)(nil)
	_ oauth2.AssertionJWTReader = (*Persister)(nil)
	_ x.FositeStorer            = (*Persister)(nil)
)

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
		logrusx.Provider
		otelx.Provider
		config.Provider
	}
	BasePersister struct {
		c           *pop.Connection
		fallbackNID uuid.UUID
		d           baseDependencies
	}
	baseDependencies interface {
		logrusx.Provider
		otelx.Provider
		contextx.Provider
		config.Provider
		jwk.ManagerProvider
	}
	BasePersisterProvider interface {
		BasePersister() *BasePersister
	}
)

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

	if rv.Kind() != reflect.Pointer || (rv.Kind() == reflect.Pointer && rv.Elem().Kind() != reflect.Struct) {
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
