package connection

import (
	"github.com/ory-am/hydra/pkg"
	"net/http"
	"net/url"
)

type HTTPManager struct {
	Endpoint *url.URL
	Client   *http.Client
}

func (m *HTTPManager) Create(connection *Connection) error {
	var r = pkg.NewSuperAgent(m.Endpoint.String())
	r.Client = m.Client
	if err := r.POST(connection); err != nil {
		return nil
	}

	return nil
}

func (m *HTTPManager) Get(id string) (*Connection, error) {
	var connection Connection
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, id).String())
	r.Client = m.Client
	if err := r.GET(&connection); err != nil {
		return nil, err
	}

	return &connection, nil

}

func (m *HTTPManager) Delete(id string) error {
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, id).String())
	r.Client = m.Client
	if err := r.DELETE(); err != nil {
		return err
	}
	return nil
}

func (m *HTTPManager) FindAllByLocalSubject(subject string) ([]*Connection, error) {
	var connection []*Connection
	var u = pkg.CopyURL(m.Endpoint)
	var q = u.Query()

	q.Add("local", subject)
	u.RawQuery = q.Encode()

	var r = pkg.NewSuperAgent(u.String())
	r.Client = m.Client
	if err := r.GET(&connection); err != nil {
		return nil, err
	}

	return connection, nil
}

func (m *HTTPManager) FindByRemoteSubject(provider, subject string) (*Connection, error) {
	var connection Connection
	var u = pkg.CopyURL(m.Endpoint)
	var q = u.Query()

	q.Add("remote", subject)
	q.Add("provider", provider)
	u.RawQuery = q.Encode()

	var r = pkg.NewSuperAgent(u.String())
	r.Client = m.Client
	if err := r.GET(&connection); err != nil {
		return nil, err
	}

	return &connection, nil
}
