package jwk

import (
	"net/http"
	"net/url"

	"github.com/ory-am/hydra/pkg"
	"github.com/square/go-jose"
)

type HTTPManager struct {
	Client   *http.Client
	Endpoint *url.URL
	Dry      bool
}

func (m *HTTPManager) CreateKeys(set, algorithm string) (*jose.JsonWebKeySet, error) {
	var c = struct {
		Algorithm string            `json:"alg"`
		Keys      []jose.JsonWebKey `json:"keys"`
	}{
		Algorithm: algorithm,
	}

	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, set).String())
	r.Client = m.Client
	r.Dry = m.Dry
	if err := r.Create(&c); err != nil {
		return nil, err
	}

	return &jose.JsonWebKeySet{
		Keys: c.Keys,
	}, nil
}

func (m *HTTPManager) AddKey(set string, key *jose.JsonWebKey) error {
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, set, key.KeyID).String())
	r.Client = m.Client
	r.Dry = m.Dry
	return r.Update(key)
}

func (m *HTTPManager) AddKeySet(set string, keys *jose.JsonWebKeySet) error {
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, set).String())
	r.Client = m.Client
	r.Dry = m.Dry
	return r.Update(keys)
}

func (m *HTTPManager) GetKey(set, kid string) (*jose.JsonWebKeySet, error) {
	var c jose.JsonWebKeySet
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, set, kid).String())
	r.Client = m.Client
	r.Dry = m.Dry
	if err := r.Get(&c); err != nil {
		return nil, err
	}

	return &c, nil
}

func (m *HTTPManager) GetKeySet(set string) (*jose.JsonWebKeySet, error) {
	var c jose.JsonWebKeySet
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, set).String())
	r.Client = m.Client
	r.Dry = m.Dry
	if err := r.Get(&c); err != nil {
		return nil, err
	}

	return &c, nil
}

func (m *HTTPManager) DeleteKey(set, kid string) error {
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, set, kid).String())
	r.Client = m.Client
	r.Dry = m.Dry
	return r.Delete()
}

func (m *HTTPManager) DeleteKeySet(set string) error {
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, set).String())
	r.Client = m.Client
	r.Dry = m.Dry
	return r.Delete()
}
