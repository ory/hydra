package postgres

import (
	"database/sql"

	log "github.com/Sirupsen/logrus"
	"github.com/go-errors/errors"
	"github.com/ory-am/hydra/connection"
	"github.com/ory-am/hydra/pkg"
)

var schemata = []string{
	`CREATE TABLE IF NOT EXISTS hydra_oauth_link (
id       		text NOT NULL PRIMARY KEY,
provider		text NOT NULL,
subject_local  	text NOT NULL,
subject_remote  text NOT NULL,

CONSTRAINT u_hydra_oauth_link_token UNIQUE (provider, subject_remote)
)`}

type Manager struct {
	DB *sql.DB
}

func (s *Manager) CreateSchemas() error {
	for _, schema := range schemata {
		if _, err := s.DB.Exec(schema); err != nil {
			log.Warnf("Error creating schema %s: %s", schema, err)
			return err
		}
	}
	return nil
}

func (s *Manager) Create(c connection.Connection) error {
	_, err := s.DB.Exec(
		"INSERT INTO hydra_oauth_link (id, provider, subject_local, subject_remote) VALUES ($1, $2, $3, $4)",
		c.GetID(),
		c.GetProvider(),
		c.GetLocalSubject(),
		c.GetRemoteSubject(),
	)
	return err
}

func (s *Manager) Delete(id string) error {
	_, err := s.DB.Exec("DELETE FROM hydra_oauth_link WHERE id=$1", id)
	return err
}

func (s *Manager) Get(id string) (connection.Connection, error) {
	var c connection.DefaultConnection

	row := s.DB.QueryRow("SELECT id, provider, subject_local, subject_remote FROM hydra_oauth_link WHERE id=$1 LIMIT 1", id)
	if err := row.Scan(&c.ID, &c.Provider, &c.LocalSubject, &c.RemoteSubject); err == sql.ErrNoRows {
		return nil, errors.New(pkg.ErrNotFound)
	} else if err != nil {
		return nil, errors.New(err)
	}
	return &c, nil
}

func (s *Manager) FindByRemoteSubject(provider, subject string) (connection.Connection, error) {
	var c connection.DefaultConnection

	row := s.DB.QueryRow("SELECT id, provider, subject_local, subject_remote FROM hydra_oauth_link WHERE provider=$1 AND subject_remote=$2 LIMIT 1", provider, subject)
	if err := row.Scan(&c.ID, &c.Provider, &c.LocalSubject, &c.RemoteSubject); err == sql.ErrNoRows {
		return nil, errors.New(pkg.ErrNotFound)
	} else if err != nil {
		return nil, errors.New(err)
	}
	return &c, nil
}

func (s *Manager) FindAllByLocalSubject(subject string) (cs []connection.Connection, err error) {
	rows, err := s.DB.Query("SELECT id, provider, subject_local, subject_remote FROM hydra_oauth_link WHERE subject_local=$1", subject)
	if err != nil {
		return []connection.Connection{}, errors.New(err)
	}
	defer rows.Close()

	for rows.Next() {
		var c connection.DefaultConnection
		if err := rows.Scan(&c.ID, &c.Provider, &c.LocalSubject, &c.RemoteSubject); err == sql.ErrNoRows {
			return []connection.Connection{}, errors.New(pkg.ErrNotFound)
		} else if err != nil {
			return []connection.Connection{}, errors.New(err)
		}
		cs = append(cs, &c)
	}

	return cs, nil
}
