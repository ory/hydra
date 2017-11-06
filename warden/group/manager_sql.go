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

package group

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/ory/hydra/pkg"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/rubenv/sql-migrate"
)

var migrations = &migrate.MemoryMigrationSource{
	Migrations: []*migrate.Migration{
		{
			Id: "1",
			Up: []string{`CREATE TABLE IF NOT EXISTS hydra_warden_group (
	id      	varchar(255) NOT NULL PRIMARY KEY
)`, `CREATE TABLE IF NOT EXISTS hydra_warden_group_member (
	member		varchar(255) NOT NULL,
	group_id	varchar(255) NOT NULL,
	FOREIGN KEY (group_id) REFERENCES hydra_warden_group(id) ON DELETE CASCADE,
	PRIMARY KEY (member, group_id)
)`},
			Down: []string{
				"DROP TABLE hydra_warden_group",
				"DROP TABLE hydra_warden_group_member",
			},
		},
	},
}

type SQLManager struct {
	DB *sqlx.DB
}

func (m *SQLManager) CreateSchemas() (int, error) {
	migrate.SetTable("hydra_groups_migration")
	n, err := migrate.Exec(m.DB.DB, m.DB.DriverName(), migrations, migrate.Up)
	if err != nil {
		return 0, errors.Wrapf(err, "Could not migrate sql schema, applied %d migrations", n)
	}
	return n, nil
}

func (m *SQLManager) CreateGroup(g *Group) error {
	if g.ID == "" {
		g.ID = uuid.New()
	}

	if _, err := m.DB.Exec(m.DB.Rebind("INSERT INTO hydra_warden_group (id) VALUES (?)"), g.ID); err != nil {
		return errors.WithStack(err)
	}

	return m.AddGroupMembers(g.ID, g.Members)
}

func (m *SQLManager) GetGroup(id string) (*Group, error) {
	var found string
	if err := m.DB.Get(&found, m.DB.Rebind("SELECT id from hydra_warden_group WHERE id = ?"), id); err != nil {
		return nil, errors.WithStack(err)
	}

	var q []string
	if err := m.DB.Select(&q, m.DB.Rebind("SELECT member from hydra_warden_group_member WHERE group_id = ?"), found); err == sql.ErrNoRows {
		return nil, errors.WithStack(pkg.ErrNotFound)
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Group{
		ID:      found,
		Members: q,
	}, nil
}

func (m *SQLManager) DeleteGroup(id string) error {
	if _, err := m.DB.Exec(m.DB.Rebind("DELETE FROM hydra_warden_group WHERE id=?"), id); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (m *SQLManager) AddGroupMembers(group string, subjects []string) error {
	tx, err := m.DB.Beginx()
	if err != nil {
		return errors.Wrap(err, "Could not begin transaction")
	}

	for _, subject := range subjects {
		if _, err := tx.Exec(m.DB.Rebind("INSERT INTO hydra_warden_group_member (group_id, member) VALUES (?, ?)"), group, subject); err != nil {
			if err := tx.Rollback(); err != nil {
				return errors.WithStack(err)
			}
			return errors.WithStack(err)
		}
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return errors.WithStack(err)
		}
		return errors.Wrap(err, "Could not commit transaction")
	}
	return nil
}

func (m *SQLManager) RemoveGroupMembers(group string, subjects []string) error {
	tx, err := m.DB.Beginx()
	if err != nil {
		return errors.Wrap(err, "Could not begin transaction")
	}
	for _, subject := range subjects {
		if _, err := m.DB.Exec(m.DB.Rebind("DELETE FROM hydra_warden_group_member WHERE member=? AND group_id=?"), subject, group); err != nil {
			if err := tx.Rollback(); err != nil {
				return errors.WithStack(err)
			}
			return errors.WithStack(err)
		}
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return errors.WithStack(err)
		}
		return errors.Wrap(err, "Could not commit transaction")
	}
	return nil
}

func (m *SQLManager) FindGroupsByMember(subject string) ([]Group, error) {
	var ids []string
	if err := m.DB.Select(&ids, m.DB.Rebind("SELECT group_id from hydra_warden_group_member WHERE member = ? GROUP BY group_id"), subject); err == sql.ErrNoRows {
		return nil, errors.WithStack(pkg.ErrNotFound)
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	var groups = make([]Group, len(ids))
	for k, id := range ids {
		group, err := m.GetGroup(id)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		groups[k] = *group
	}

	return groups, nil
}
