package group

import (
	"github.com/jmoiron/sqlx"
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

func (s *SQLManager) CreateSchemas() error {
	migrate.SetTable("hydra_groups_migration")
	n, err := migrate.Exec(s.DB.DB, s.DB.DriverName(), migrations, migrate.Up)
	if err != nil {
		return errors.Wrapf(err, "Could not migrate sql schema, applied %d migrations", n)
	}

	return nil
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
	if err := m.DB.Select(&q, m.DB.Rebind("SELECT member from hydra_warden_group_member WHERE group_id = ?"), found); err != nil {
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
			return errors.WithStack(err)
		}
	}

	if err := tx.Commit(); err != nil {
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
			return errors.WithStack(err)
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "Could not commit transaction")
	}
	return nil
}

func (m *SQLManager) FindGroupNames(subject string) ([]string, error) {
	var q []string
	if err := m.DB.Select(&q, m.DB.Rebind("SELECT group_id from hydra_warden_group_member WHERE member = ? GROUP BY group_id"), subject); err != nil {
		return nil, errors.WithStack(err)
	}

	return q, nil
}
