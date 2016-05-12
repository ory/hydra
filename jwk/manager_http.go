package jwk

import (
	"github.com/ory-am/hydra/pkg"
	"github.com/square/go-jose"
	"net/http"
	"net/url"
)

type HTTPManager struct {
	Client   *http.Client
	Endpoint *url.URL
}

func (m *HTTPManager) AddKey(set string, key *jose.JsonWebKey) error {
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, set, key.KeyID).String())
	r.Client = m.Client
	return r.PUT(key)
}

func (m *HTTPManager) AddKeySet(set string, keys *jose.JsonWebKeySet) error {
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, set).String())
	r.Client = m.Client
	return r.PUT(keys)
}

func (m *HTTPManager) GetKey(set, kid string) (*jose.JsonWebKey, error) {
	var c jose.JsonWebKey
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, set, kid).String())
	r.Client = m.Client
	if err := r.GET(&c); err != nil {
		return nil, err
	}

	return &c, nil
}

func (m *HTTPManager) GetKeySet(set string) (*jose.JsonWebKeySet, error) {
	var c jose.JsonWebKeySet
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, set).String())
	r.Client = m.Client
	if err := r.GET(&c); err != nil {
		return nil, err
	}

	return &c, nil
}

func (m *HTTPManager) DeleteKey(set, kid string) error {
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, set, kid).String())
	r.Client = m.Client
	return r.DELETE()
}

func (m *HTTPManager) DeleteKeySet(set string) error {
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, set).String())
	r.Client = m.Client
	return r.DELETE()
}
