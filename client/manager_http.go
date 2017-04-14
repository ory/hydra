package client

import (
	"net/http"
	"net/url"

	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/pkg"
)

type HTTPManager struct {
	Client   *http.Client
	Endpoint *url.URL
	Dry      bool
}

func (m *HTTPManager) GetConcreteClient(id string) (*Client, error) {
	var c Client
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, id).String())
	r.Client = m.Client
	r.Dry = m.Dry
	if err := r.Get(&c); err != nil {
		return nil, err
	}

	return &c, nil
}

func (m *HTTPManager) GetClient(id string) (fosite.Client, error) {
	return m.GetConcreteClient(id)
}

func (m *HTTPManager) UpdateClient(c *Client) error {
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, c.ID).String())
	r.Client = m.Client
	r.Dry = m.Dry
	return r.Update(c)
}

func (m *HTTPManager) CreateClient(c *Client) error {
	var r = pkg.NewSuperAgent(m.Endpoint.String())
	r.Client = m.Client
	r.Dry = m.Dry
	return r.Create(c)
}

func (m *HTTPManager) DeleteClient(id string) error {
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, id).String())
	r.Client = m.Client
	r.Dry = m.Dry
	return r.Delete()
}

func (m *HTTPManager) GetClients() (map[string]Client, error) {
	cs := make(map[string]Client)
	var r = pkg.NewSuperAgent(m.Endpoint.String())
	r.Client = m.Client
	r.Dry = m.Dry
	if err := r.Get(&cs); err != nil {
		return nil, err
	}

	return cs, nil
}
