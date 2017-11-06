// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/ory/fosite"
	"github.com/ory/hydra/pkg"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/rubenv/sql-migrate"
)

var migrations = &migrate.MemoryMigrationSource{
	Migrations: []*migrate.Migration{
		{
			Id: "1",
			Up: []string{`CREATE TABLE IF NOT EXISTS hydra_client (
	id      	varchar(255) NOT NULL PRIMARY KEY,
	client_name  	text NOT NULL,
	client_secret  	text NOT NULL,
	redirect_uris  	text NOT NULL,
	grant_types  	text NOT NULL,
	response_types  text NOT NULL,
	scope  			text NOT NULL,
	owner  			text NOT NULL,
	policy_uri  	text NOT NULL,
	tos_uri  		text NOT NULL,
	client_uri  	text NOT NULL,
	logo_uri  		text NOT NULL,
	contacts  		text NOT NULL,
	public  		boolean NOT NULL
)`},
			Down: []string{
				"DROP TABLE hydra_client",
			},
		},
	},
}

type SQLManager struct {
	Hasher fosite.Hasher
	DB     *sqlx.DB
}

type sqlData struct {
	ID                string `db:"id"`
	Name              string `db:"client_name"`
	Secret            string `db:"client_secret"`
	RedirectURIs      string `db:"redirect_uris"`
	GrantTypes        string `db:"grant_types"`
	ResponseTypes     string `db:"response_types"`
	Scope             string `db:"scope"`
	Owner             string `db:"owner"`
	PolicyURI         string `db:"policy_uri"`
	TermsOfServiceURI string `db:"tos_uri"`
	ClientURI         string `db:"client_uri"`
	LogoURI           string `db:"logo_uri"`
	Contacts          string `db:"contacts"`
	Public            bool   `db:"public"`
}

var sqlParams = []string{
	"id",
	"client_name",
	"client_secret",
	"redirect_uris",
	"grant_types",
	"response_types",
	"scope",
	"owner",
	"policy_uri",
	"tos_uri",
	"client_uri",
	"logo_uri",
	"contacts",
	"public",
}

func sqlDataFromClient(d *Client) *sqlData {
	return &sqlData{
		ID:                d.ID,
		Name:              d.Name,
		Secret:            d.Secret,
		RedirectURIs:      strings.Join(d.RedirectURIs, "|"),
		GrantTypes:        strings.Join(d.GrantTypes, "|"),
		ResponseTypes:     strings.Join(d.ResponseTypes, "|"),
		Scope:             d.Scope,
		Owner:             d.Owner,
		PolicyURI:         d.PolicyURI,
		TermsOfServiceURI: d.TermsOfServiceURI,
		ClientURI:         d.ClientURI,
		LogoURI:           d.LogoURI,
		Contacts:          strings.Join(d.Contacts, "|"),
		Public:            d.Public,
	}
}

func (d *sqlData) ToClient() *Client {
	return &Client{
		ID:                d.ID,
		Name:              d.Name,
		Secret:            d.Secret,
		RedirectURIs:      pkg.SplitNonEmpty(d.RedirectURIs, "|"),
		GrantTypes:        pkg.SplitNonEmpty(d.GrantTypes, "|"),
		ResponseTypes:     pkg.SplitNonEmpty(d.ResponseTypes, "|"),
		Scope:             d.Scope,
		Owner:             d.Owner,
		PolicyURI:         d.PolicyURI,
		TermsOfServiceURI: d.TermsOfServiceURI,
		ClientURI:         d.ClientURI,
		LogoURI:           d.LogoURI,
		Contacts:          pkg.SplitNonEmpty(d.Contacts, "|"),
		Public:            d.Public,
	}
}

func (s *SQLManager) CreateSchemas() (int, error) {
	migrate.SetTable("hydra_client_migration")
	n, err := migrate.Exec(s.DB.DB, s.DB.DriverName(), migrations, migrate.Up)
	if err != nil {
		return 0, errors.Wrapf(err, "Could not migrate sql schema, applied %d migrations", n)
	}
	return n, nil
}

func (m *SQLManager) GetConcreteClient(id string) (*Client, error) {
	var d sqlData
	if err := m.DB.Get(&d, m.DB.Rebind("SELECT * FROM hydra_client WHERE id=?"), id); err == sql.ErrNoRows {
		return nil, errors.Wrap(pkg.ErrNotFound, "")
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	return d.ToClient(), nil
}

func (m *SQLManager) GetClient(_ context.Context, id string) (fosite.Client, error) {
	return m.GetConcreteClient(id)
}

func (m *SQLManager) UpdateClient(c *Client) error {
	o, err := m.GetClient(context.Background(), c.ID)
	if err != nil {
		return errors.WithStack(err)
	}

	if c.Secret == "" {
		c.Secret = string(o.GetHashedSecret())
	} else {
		h, err := m.Hasher.Hash([]byte(c.Secret))
		if err != nil {
			return errors.WithStack(err)
		}
		c.Secret = string(h)
	}

	s := sqlDataFromClient(c)
	var query []string
	for _, param := range sqlParams {
		query = append(query, fmt.Sprintf("%s=:%s", param, param))
	}

	if _, err := m.DB.NamedExec(fmt.Sprintf(`UPDATE hydra_client SET %s WHERE id=:id`, strings.Join(query, ", ")), s); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (m *SQLManager) Authenticate(id string, secret []byte) (*Client, error) {
	c, err := m.GetConcreteClient(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := m.Hasher.Compare(c.GetHashedSecret(), secret); err != nil {
		return nil, errors.WithStack(err)
	}

	return c, nil
}

func (m *SQLManager) CreateClient(c *Client) error {
	if c.ID == "" {
		c.ID = uuid.New()
	}

	h, err := m.Hasher.Hash([]byte(c.Secret))
	if err != nil {
		return errors.WithStack(err)
	}
	c.Secret = string(h)

	data := sqlDataFromClient(c)
	if _, err := m.DB.NamedExec(fmt.Sprintf(
		"INSERT INTO hydra_client (%s) VALUES (%s)",
		strings.Join(sqlParams, ", "),
		":"+strings.Join(sqlParams, ", :"),
	), data); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (m *SQLManager) DeleteClient(id string) error {
	if _, err := m.DB.Exec(m.DB.Rebind(`DELETE FROM hydra_client WHERE id=?`), id); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (m *SQLManager) GetClients() (clients map[string]Client, err error) {
	var d = []sqlData{}
	clients = make(map[string]Client)

	if err := m.DB.Select(&d, "SELECT * FROM hydra_client"); err != nil {
		return nil, errors.WithStack(err)
	}

	for _, k := range d {
		clients[k.ID] = *k.ToClient()
	}
	return clients, nil
}
