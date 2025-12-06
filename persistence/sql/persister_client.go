// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"context"

	"github.com/gofrs/uuid"
	"go.opentelemetry.io/otel/trace"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/x/events"
	"github.com/ory/pop/v6"
	"github.com/ory/x/otelx"
	keysetpagination "github.com/ory/x/pagination/keysetpagination_v2"
	"github.com/ory/x/sqlcon"
)

// AuthenticateClient implements client.Manager.
func (p *Persister) AuthenticateClient(ctx context.Context, id string, secret []byte) (_ *client.Client, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.AuthenticateClient",
		trace.WithAttributes(events.ClientID(id)),
	)
	defer otelx.End(span, &err)

	c, err := p.GetConcreteClient(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := p.r.ClientHasher().Compare(ctx, c.GetHashedSecret(), secret); err != nil {
		return nil, err
	}

	return c, nil
}

// CreateClient implements client.Storage.
func (p *Persister) CreateClient(ctx context.Context, c *client.Client) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.CreateClient")
	defer otelx.End(span, &err)

	h, err := p.r.ClientHasher().Hash(ctx, []byte(c.Secret))
	if err != nil {
		return err
	}

	c.Secret = string(h)
	if c.ID == "" {
		c.ID = uuid.Must(uuid.NewV4()).String()
	}
	if err := sqlcon.HandleError(p.CreateWithNetwork(ctx, c)); err != nil {
		return err
	}

	events.Trace(ctx, events.ClientCreated,
		events.WithClientID(c.ID),
		events.WithClientName(c.Name))

	return nil
}

// UpdateClient implements client.Storage.
func (p *Persister) UpdateClient(ctx context.Context, cl *client.Client) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.UpdateClient",
		trace.WithAttributes(events.ClientID(cl.ID)),
	)
	defer otelx.End(span, &err)

	return p.Transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		o, err := p.GetConcreteClient(ctx, cl.GetID())
		if err != nil {
			return err
		}

		if cl.Secret == "" {
			cl.Secret = string(o.GetHashedSecret())
		} else {
			h, err := p.r.ClientHasher().Hash(ctx, []byte(cl.Secret))
			if err != nil {
				return err
			}
			cl.Secret = string(h)
		}

		// Ensure ID is the same
		cl.ID = o.ID

		if err = cl.BeforeSave(c); err != nil {
			return sqlcon.HandleError(err)
		}

		count, err := p.UpdateWithNetwork(ctx, cl)
		if err != nil {
			return sqlcon.HandleError(err)
		} else if count == 0 {
			return sqlcon.HandleError(sqlcon.ErrNoRows)
		}

		events.Trace(ctx, events.ClientUpdated,
			events.WithClientID(cl.ID),
			events.WithClientName(cl.Name))

		return sqlcon.HandleError(err)
	})
}

// DeleteClient implements client.Storage.
func (p *Persister) DeleteClient(ctx context.Context, id string) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.DeleteClient",
		trace.WithAttributes(events.ClientID(id)),
	)
	defer otelx.End(span, &err)

	c, err := p.GetConcreteClient(ctx, id)
	if err != nil {
		return err
	}

	if err := sqlcon.HandleError(p.QueryWithNetwork(ctx).Where("id = ?", id).Delete(&client.Client{})); err != nil {
		return err
	}

	events.Trace(ctx, events.ClientDeleted,
		events.WithClientID(c.ID),
		events.WithClientName(c.Name))

	return nil
}

// GetClients implements client.Storage.
func (p *Persister) GetClients(ctx context.Context, filters client.Filter) (cs []client.Client, _ *keysetpagination.Paginator, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetClients")
	defer otelx.End(span, &err)

	paginator := keysetpagination.NewPaginator(append(filters.PageOpts,
		keysetpagination.WithDefaultToken(keysetpagination.NewPageToken(keysetpagination.Column{Name: "id", Value: ""})),
	)...)

	query := p.QueryWithNetwork(ctx).Scope(
		keysetpagination.Paginate[client.Client](paginator))

	if filters.Name != "" {
		query.Where("client_name = ?", filters.Name)
	}
	if filters.Owner != "" {
		query.Where("owner = ?", filters.Owner)
	}
	if len(filters.IDs) > 0 {
		query.Where("id IN (?)", filters.IDs)
	}

	if err := query.All(&cs); err != nil {
		return nil, nil, sqlcon.HandleError(err)
	}
	cs, nextPage := keysetpagination.Result(cs, paginator)
	return cs, nextPage, nil
}

// GetConcreteClient implements client.Storage.
func (p *Persister) GetConcreteClient(ctx context.Context, id string) (c *client.Client, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetConcreteClient",
		trace.WithAttributes(events.ClientID(id)),
	)
	defer otelx.End(span, &err)

	var cl client.Client
	if err := p.QueryWithNetwork(ctx).Where("id = ?", id).First(&cl); err != nil {
		return nil, sqlcon.HandleError(err)
	}
	return &cl, nil
}

// GetClient implements fosite.ClientManager.
func (p *Persister) GetClient(ctx context.Context, id string) (fosite.Client, error) {
	return p.GetConcreteClient(ctx, id)
}
