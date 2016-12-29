package group

import (
	"net/http"
	"net/url"

	"bytes"
	"encoding/json"
	"github.com/ory-am/hydra/pkg"
	"github.com/pkg/errors"
	"io/ioutil"
)

type HTTPManager struct {
	Client   *http.Client
	Endpoint *url.URL
	Dry      bool
}

func (m *HTTPManager) CreateGroup(g *Group) error {
	var r = pkg.NewSuperAgent(m.Endpoint.String())
	r.Client = m.Client
	r.Dry = m.Dry
	return r.Create(g)
}

func (m *HTTPManager) GetGroup(id string) (*Group, error) {
	var g Group
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, id).String())
	r.Client = m.Client
	r.Dry = m.Dry
	if err := r.Get(&g); err != nil {
		return nil, err
	}

	return &g, nil
}

func (m *HTTPManager) DeleteGroup(id string) error {
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, id).String())
	r.Client = m.Client
	r.Dry = m.Dry
	return r.Delete()
}

func (m *HTTPManager) AddGroupMembers(group string, members []string) error {
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, group, "members").String())
	r.Client = m.Client
	r.Dry = m.Dry
	return r.Create(&membersRequest{
		Members: members,
	})
}

func (m *HTTPManager) RemoveGroupMembers(group string, members []string) error {
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, group, "members").String())
	r.Client = m.Client
	r.Dry = m.Dry
	send, err := json.Marshal(&membersRequest{Members: members})
	if err != nil {
		return errors.WithStack(err)
	}

	req, err := http.NewRequest("DELETE", r.URL, bytes.NewReader(send))
	if err != nil {
		return errors.WithStack(err)
	}

	if err := r.DoDry(req); err != nil {
		return err
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		return errors.WithStack(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.Errorf("Expected status code %d, got %d.\n%s\n", http.StatusNoContent, resp.StatusCode, body)
	}

	return nil
}

func (m *HTTPManager) FindGroupNames(subject string) ([]string, error) {
	var g []string
	var r = pkg.NewSuperAgent(m.Endpoint.String() + "?member=" + subject)
	r.Client = m.Client
	r.Dry = m.Dry
	if err := r.Get(&g); err != nil {
		return nil, err
	}

	return g, nil
}
