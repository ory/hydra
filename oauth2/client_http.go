package oauth2

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"path"

	"github.com/go-errors/errors"
	"github.com/ory-am/fosite/client"
)

type HTTPClientManager struct {
	Client *http.Client

	Endpoint *url.URL
}

func (m *HTTPClientManager) GetClient(id string) (client.Client, error) {
	var ep = &url.URL{}
	*ep = *m.Endpoint
	ep.Path = path.Join(ep.Path, id)
	req, err := http.NewRequest("GET", ep.String(), nil)
	if err != nil {
		return nil, errors.New(err)
	}

	var c OAuth2Client
	resp, err := m.Client.Do(req)
	if err != nil {
		return nil, errors.New(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Could not fetch client")
	}

	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return nil, errors.New(err)
	}

	return &c, nil
}

func (m *HTTPClientManager) CreateClient(c *OAuth2Client) error {
	data, err := json.Marshal(c)
	if err != nil {
		return errors.New(err)
	}

	req, err := http.NewRequest("POST", m.Endpoint.String(), bytes.NewReader(data))
	if err != nil {
		return errors.New(err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := m.Client.Do(req)
	if err != nil {
		return errors.New(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return errors.New("Could not create client")
	}

	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return errors.New(err)
	}

	return nil
}

func (m *HTTPClientManager) DeleteClient(id string) error {
	var ep = &url.URL{}
	*ep = *m.Endpoint
	ep.Path = path.Join(ep.Path, id)
	req, err := http.NewRequest("DELETE", ep.String(), nil)
	if err != nil {
		return errors.New(err)
	}

	resp, err := m.Client.Do(req)
	if err != nil {
		return errors.New(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return errors.New("Could not delete client")
	}

	return nil
}
