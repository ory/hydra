package jwk

import (
	"net/http"
	"net/url"

	"github.com/ory/hydra/pkg"
	"github.com/square/go-jose"
)

type HTTPManager struct {
	Client             *http.Client
	Endpoint           *url.URL
	Dry                bool
	FakeTLSTermination bool
}

func (m *HTTPManager) CreateKeys(set, algorithm string) (*jose.JSONWebKeySet, error) {
	var c = struct {
		Algorithm string            `json:"alg"`
		Keys      []jose.JSONWebKey `json:"keys"`
	}{
		Algorithm: algorithm,
	}

	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, set).String())
	r.Client = m.Client
	r.Dry = m.Dry
	r.FakeTLSTermination = m.FakeTLSTermination
	if err := r.Create(&c); err != nil {
		return nil, err
	}

	return &jose.JSONWebKeySet{
		Keys: c.Keys,
	}, nil
}

func (m *HTTPManager) AddKey(set string, key *jose.JSONWebKey) error {
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, set, key.KeyID).String())
	r.Client = m.Client
	r.Dry = m.Dry
	r.FakeTLSTermination = m.FakeTLSTermination
	return r.Update(key)
}

func (m *HTTPManager) AddKeySet(set string, keys *jose.JSONWebKeySet) error {
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, set).String())
	r.Client = m.Client
	r.Dry = m.Dry
	r.FakeTLSTermination = m.FakeTLSTermination
	return r.Update(keys)
}

func (m *HTTPManager) GetKey(set, kid string) (*jose.JSONWebKeySet, error) {
	var c jose.JSONWebKeySet
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, set, kid).String())
	r.Client = m.Client
	r.Dry = m.Dry
	r.FakeTLSTermination = m.FakeTLSTermination
	if err := r.Get(&c); err != nil {
		return nil, err
	}

	return &c, nil
}

func (m *HTTPManager) GetKeySet(set string) (*jose.JSONWebKeySet, error) {
	var c jose.JSONWebKeySet
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, set).String())
	r.Client = m.Client
	r.Dry = m.Dry
	r.FakeTLSTermination = m.FakeTLSTermination
	if err := r.Get(&c); err != nil {
		return nil, err
	}

	return &c, nil
}

func (m *HTTPManager) DeleteKey(set, kid string) error {
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, set, kid).String())
	r.Client = m.Client
	r.Dry = m.Dry
	r.FakeTLSTermination = m.FakeTLSTermination
	return r.Delete()
}

func (m *HTTPManager) DeleteKeySet(set string) error {
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, set).String())
	r.Client = m.Client
	r.Dry = m.Dry
	r.FakeTLSTermination = m.FakeTLSTermination
	return r.Delete()
}
