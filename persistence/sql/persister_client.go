package sql

import (
	"context"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/ory/x/errorsx"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/x/sqlcon"
)

func (p *Persister) GetConcreteClient(ctx context.Context, id string) (*client.Client, error) {
	var cl client.Client
	return &cl, sqlcon.HandleError(p.QueryWithNetwork(ctx).Where("id = ?", id).First(&cl))
}

func (p *Persister) GetClient(ctx context.Context, id string) (fosite.Client, error) {
	return p.GetConcreteClient(ctx, id)
}

func (p *Persister) UpdateClient(ctx context.Context, cl *client.Client) error {
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

		return sqlcon.HandleError(p.UpdateWithNetwork(ctx, cl))
	})
}

func (p *Persister) Authenticate(ctx context.Context, id string, secret []byte) (*client.Client, error) {
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
	h, err := p.r.ClientHasher().Hash(ctx, []byte(c.Secret))
	if err != nil {
		return err
	}

	c.Secret = string(h)
	c.ID = uuid.Must(uuid.NewV4())
	return sqlcon.HandleError(p.CreateWithNetwork(ctx, c))
}

func (p *Persister) DeleteClient(ctx context.Context, id string) error {
	cl, err := p.GetConcreteClient(ctx, id)
	if err != nil {
		return err
	}

	return sqlcon.HandleError(p.QueryWithNetwork(ctx).Where("pk = ?", cl.ID).Delete(&client.Client{}))
}

func (p *Persister) GetClients(ctx context.Context, filters client.Filter) ([]client.Client, error) {
	cs := make([]client.Client, 0)

	query := p.QueryWithNetwork(ctx).
		Paginate(filters.Offset/filters.Limit+1, filters.Limit).
		Order("id")

	if filters.Name != "" {
		query.Where("client_name = ?", filters.Name)
	}
	if filters.Owner != "" {
		query.Where("owner = ?", filters.Owner)
	}

	return cs, sqlcon.HandleError(query.All(&cs))
}

func (p *Persister) CountClients(ctx context.Context) (int, error) {
	n, err := p.QueryWithNetwork(ctx).Count(&client.Client{})
	return n, sqlcon.HandleError(err)
}
