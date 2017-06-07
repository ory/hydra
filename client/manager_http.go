package client

import (
	"net/http"
	"net/url"

	"context"

	"github.com/ory/fosite"
	"github.com/ory/hydra/pkg"
	"github.com/pkg/errors"
)

type HTTPManager struct {
	Client             *http.Client
	Endpoint           *url.URL
	Dry                bool
	FakeTLSTermination bool
}

func (m *HTTPManager) GetConcreteClient(id string) (*Client, error) {
	var c Client
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, id).String())
	r.Client = m.Client
	r.Dry = m.Dry
	r.FakeTLSTermination = m.FakeTLSTermination

	if err := r.Get(&c); err != nil {
		return nil, errors.WithStack(err)
	}

	return &c, nil
}

func (m *HTTPManager) GetClient(_ context.Context, id string) (fosite.Client, error) {
	return m.GetConcreteClient(id)
}

func (m *HTTPManager) UpdateClient(c *Client) error {
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, c.ID).String())
	r.Client = m.Client
	r.Dry = m.Dry
	r.FakeTLSTermination = m.FakeTLSTermination
	return r.Update(c)
}

func (m *HTTPManager) CreateClient(c *Client) error {
	var r = pkg.NewSuperAgent(m.Endpoint.String())
	r.Client = m.Client
	r.Dry = m.Dry
	r.FakeTLSTermination = m.FakeTLSTermination
	return r.Create(c)
}

func (m *HTTPManager) DeleteClient(id string) error {
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, id).String())
	r.Client = m.Client
	r.Dry = m.Dry
	r.FakeTLSTermination = m.FakeTLSTermination
	return r.Delete()
}

func (m *HTTPManager) GetClients() (map[string]Client, error) {
	cs := make(map[string]Client)
	var r = pkg.NewSuperAgent(m.Endpoint.String())
	r.Client = m.Client
	r.Dry = m.Dry
	r.FakeTLSTermination = m.FakeTLSTermination
	if err := r.Get(&cs); err != nil {
		return nil, errors.WithStack(err)
	}

	return cs, nil
}
