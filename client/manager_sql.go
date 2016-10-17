package client

import (
	"encoding/json"
	"github.com/imdario/mergo"
	"github.com/jmoiron/sqlx"
	"github.com/ory-am/fosite"
	"github.com/pkg/errors"
)

var sqlSchema = []string{
	`CREATE TABLE IF NOT EXISTS hydra_client (
	id      varchar(255) NOT NULL PRIMARY KEY,
	version int NOT NULL DEFAULT 0,
	client  json NOT NULL)`,
}

type SQLManager struct {
	Hasher fosite.Hasher
	DB     *sqlx.DB
}

type sqlData struct {
	ID      string `db:"id"`
	Version int    `db:"version"`
	Client  []byte `db:"client"`
}

// CreateSchemas creates ladon_policy tables
func (s *SQLManager) CreateSchemas() error {
	for _, query := range sqlSchema {
		if _, err := s.DB.Exec(query); err != nil {
			return errors.Wrapf(err, "Could not create schema:\n%s", query)
		}
	}
	return nil
}

func (m *SQLManager) GetConcreteClient(id string) (*Client, error) {
	var d sqlData
	var c Client
	if err := m.DB.Get(&d, m.DB.Rebind("SELECT * FROM hydra_client WHERE id=?"), id); err != nil {
		return nil, errors.Wrap(err, "")
	} else if err := json.Unmarshal(d.Client, &c); err != nil {
		return nil, errors.Wrap(err, "")
	}

	return &c, nil
}

func (m *SQLManager) GetClient(id string) (fosite.Client, error) {
	return m.GetConcreteClient(id)
}

func (m *SQLManager) UpdateClient(c *Client) error {
	o, err := m.GetClient(c.ID)
	if err != nil {
		return err
	}

	if c.Secret == "" {
		c.Secret = string(o.GetHashedSecret())
	} else {
		h, err := m.Hasher.Hash([]byte(c.Secret))
		if err != nil {
			return errors.Wrap(err, "")
		}
		c.Secret = string(h)
	}
	if err := mergo.Merge(c, o); err != nil {
		return errors.Wrap(err, "")
	}

	b, err := json.Marshal(c)
	if err != nil {
		return errors.Wrap(err, "")
	}

	if _, err := m.DB.NamedExec(`UPDATE hydra_client SET id=:id, client=:client WHERE id=:id`, &sqlData{
		ID:     c.ID,
		Client: b,
	}); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func (m *SQLManager) Authenticate(id string, secret []byte) (*Client, error) {
	c, err := m.GetConcreteClient(id)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	if err := m.Hasher.Compare(c.GetHashedSecret(), secret); err != nil {
		return nil, errors.Wrap(err, "")
	}

	return c, nil
}

func (m *SQLManager) CreateClient(c *Client) error {
	b, err := json.Marshal(c)
	if err != nil {
		return errors.Wrap(err, "")
	}

	if _, err := m.DB.NamedExec(`INSERT INTO hydra_client (id, client, version) VALUES (:id, :client, :version)`, &sqlData{
		ID:      c.ID,
		Client:  b,
		Version: 0,
	}); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func (m *SQLManager) DeleteClient(id string) error {
	if _, err := m.DB.Exec(m.DB.Rebind(`DELETE FROM hydra_client WHERE id=?`), id); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func (m *SQLManager) GetClients() (clients map[string]Client, err error) {
	var d = []sqlData{}
	clients = make(map[string]Client)

	if err := m.DB.Select(&d, "SELECT * FROM hydra_client"); err != nil {
		return nil, errors.Wrap(err, "")
	}

	for _, k := range d {
		var c Client
		if err := json.Unmarshal(k.Client, &c); err != nil {
			return nil, errors.Wrap(err, "")
		}
		clients[k.ID] = c
	}
	return clients, nil
}
