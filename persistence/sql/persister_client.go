package sql

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/gobuffalo/pop/v6"

	"github.com/ory/x/errorsx"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/x/sqlcon"
)

func (p *Persister) GetConcreteClient(ctx context.Context, id string) (*client.Client, error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetConcreteClient")
	defer span.End()

	var cl client.Client
	return &cl, sqlcon.HandleError(p.QueryWithNetwork(ctx).Where("id = ?", id).First(&cl))
}

func (p *Persister) GetClient(ctx context.Context, id string) (fosite.Client, error) {
	return p.GetConcreteClient(ctx, id)
}

func (p *Persister) UpdateClient(ctx context.Context, cl *client.Client) error {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.UpdateClient")
	defer span.End()

	return p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		o, err := p.GetConcreteClient(ctx, cl.GetID())
		if err != nil {
			return err
		}

		if cl.Secret == "" {
			cl.Secret = string(o.GetHashedSecret())
		} else {
			h, err := p.r.ClientHasher().Hash(ctx, []byte(cl.Secret))
			if err != nil {
				return errorsx.WithStack(err)
			}
			cl.Secret = string(h)
		}
		// set the internal primary key
		cl.ID = o.ID

		// Set the legacy client ID
		cl.LegacyClientID = o.LegacyClientID

		if err = cl.BeforeSave(c); err != nil {
			return sqlcon.HandleError(err)
		}

		count, err := p.UpdateWithNetwork(ctx, cl)
		if err != nil {
			return sqlcon.HandleError(err)
		} else if count == 0 {
			return sqlcon.HandleError(sqlcon.ErrNoRows)
		}
		return sqlcon.HandleError(err)
	})
}

func (p *Persister) Authenticate(ctx context.Context, id string, secret []byte) (*client.Client, error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.Authenticate")
	defer span.End()

	c, err := p.GetConcreteClient(ctx, id)
	if err != nil {
		return nil, errorsx.WithStack(err)
	}

	if err := p.r.ClientHasher().Compare(ctx, c.GetHashedSecret(), secret); err != nil {
		return nil, errorsx.WithStack(err)
	}

	return c, nil
}

func (p *Persister) CreateClient(ctx context.Context, c *client.Client) error {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.CreateClient")
	defer span.End()

	h, err := p.r.ClientHasher().Hash(ctx, []byte(c.Secret))
	if err != nil {
		return err
	}

	c.Secret = string(h)
	if c.ID == uuid.Nil {
		c.ID = uuid.Must(uuid.NewV4())
	}
	if c.LegacyClientID == "" {
		c.LegacyClientID = c.ID.String()
	}
	return sqlcon.HandleError(p.CreateWithNetwork(ctx, c))
}

func (p *Persister) DeleteClient(ctx context.Context, id string) error {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.DeleteClient")
	defer span.End()

	_, err := p.GetConcreteClient(ctx, id)
	if err != nil {
		return err
	}

	return sqlcon.HandleError(p.QueryWithNetwork(ctx).Where("id = ?", id).Delete(&client.Client{}))
}

func (p *Persister) GetClients(ctx context.Context, filters client.Filter) ([]client.Client, error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetClients")
	defer span.End()

	cs := make([]client.Client, 0)

	query := p.QueryWithNetwork(ctx).
		Paginate(filters.Offset/filters.Limit+1, filters.Limit).
		Order("pk")

	if filters.Name != "" {
		query.Where("client_name = ?", filters.Name)
	}
	if filters.Owner != "" {
		query.Where("owner = ?", filters.Owner)
	}

	return cs, sqlcon.HandleError(query.All(&cs))
}

func (p *Persister) CountClients(ctx context.Context) (int, error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.CountClients")
	defer span.End()

	n, err := p.QueryWithNetwork(ctx).Count(&client.Client{})
	return n, sqlcon.HandleError(err)
}
