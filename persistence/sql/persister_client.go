package sql

import (
	"context"

	"github.com/ory/x/errorsx"

	"github.com/gobuffalo/pop/v5"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/x/sqlcon"
)

func (p *Persister) GetConcreteClient(ctx context.Context, id string) (*client.Client, error) {
	var cl client.Client
	return &cl, sqlcon.HandleError(p.Connection(ctx).Where("id = ?", id).First(&cl))
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

		return sqlcon.HandleError(c.Update(cl))
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
	return sqlcon.HandleError(p.Connection(ctx).Create(c, "pk"))
}

func (p *Persister) DeleteClient(ctx context.Context, id string) error {
	cl, err := p.GetConcreteClient(ctx, id)
	if err != nil {
		return err
	}

	return sqlcon.HandleError(p.Connection(ctx).Destroy(&client.Client{ID: cl.ID}))
}

func (p *Persister) GetClients(ctx context.Context, limit, offset int) ([]client.Client, error) {
	cs := make([]client.Client, 0)
	return cs, sqlcon.HandleError(
		p.Connection(ctx).
			Paginate(offset/limit+1, limit).
			Order("id").
			All(&cs))
}

func (p *Persister) CountClients(ctx context.Context) (int, error) {
	n, err := p.Connection(ctx).Count(&client.Client{})
	return n, sqlcon.HandleError(err)
}
