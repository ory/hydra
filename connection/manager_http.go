package connection

import (
	"net/http"
	"net/url"

	"github.com/ory-am/hydra/pkg"
)

type HTTPManager struct {
	Endpoint *url.URL
	Client   *http.Client
}

func (m *HTTPManager) Create(connection *Connection) error {
	var r = pkg.NewSuperAgent(m.Endpoint.String())
	r.Client = m.Client
	return r.Create(connection)
}

func (m *HTTPManager) Get(id string) (*Connection, error) {
	var connection Connection
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, id).String())
	r.Client = m.Client
	if err := r.Get(&connection); err != nil {
		return nil, err
	}

	return &connection, nil
}

func (m *HTTPManager) Delete(id string) error {
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, id).String())
	r.Client = m.Client
	return r.Delete()
}

func (m *HTTPManager) FindAllByLocalSubject(subject string) ([]*Connection, error) {
	var connection []*Connection
	var u = pkg.CopyURL(m.Endpoint)
	var q = u.Query()

	q.Add("local_subject", subject)
	u.RawQuery = q.Encode()

	var r = pkg.NewSuperAgent(u.String())
	r.Client = m.Client
	if err := r.Get(&connection); err != nil {
		return nil, err
	}

	return connection, nil
}

func (m *HTTPManager) FindByRemoteSubject(provider, subject string) (*Connection, error) {
	var connection Connection
	var u = pkg.CopyURL(m.Endpoint)
	var q = u.Query()
	q.Add("remote_subject", subject)
	q.Add("provider", provider)
	u.RawQuery = q.Encode()

	var r = pkg.NewSuperAgent(u.String())
	r.Client = m.Client
	if err := r.Get(&connection); err != nil {
		return nil, err
	}

	return &connection, nil
}
