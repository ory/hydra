package postgres

import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
	. "github.com/ory-am/common/pkg"
	. "github.com/ory-am/hydra/endpoint/connection"
)

var schemata = []string{
	`CREATE TABLE IF NOT EXISTS hydra_oauth_link (
	id       		text NOT NULL PRIMARY KEY,
	provider		text NOT NULL,
    subject_local  	text NOT NULL,
    subject_remote  text NOT NULL,

	CONSTRAINT u_hydra_oauth_link_token UNIQUE (provider, subject_remote)
)`}

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateSchemas() error {
	for _, schema := range schemata {
		if _, err := s.db.Exec(schema); err != nil {
			log.Warnf("Error creating schema %s: %s", schema, err)
			return err
		}
	}
	return nil
}

func (s *Store) Create(c Connection) error {
	_, err := s.db.Exec(
		"INSERT INTO hydra_oauth_link (id, provider, subject_local, subject_remote) VALUES ($1, $2, $3, $4)",
		c.GetID(),
		c.GetProvider(),
		c.GetLocalSubject(),
		c.GetRemoteSubject(),
	)
	return err
}

func (s *Store) Delete(id string) error {
	_, err := s.db.Exec("DELETE FROM hydra_oauth_link WHERE id=$1", id)
	return err

}

func (s *Store) Get(id string) (Connection, error) {
	var c DefaultConnection
	row := s.db.QueryRow("SELECT id, provider, subject_local, subject_remote FROM hydra_oauth_link WHERE id=$1 LIMIT 1", id)
	if err := row.Scan(&c.ID, &c.Provider, &c.LocalSubject, &c.RemoteSubject); err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *Store) FindByRemoteSubject(provider, subject string) (Connection, error) {
	var c DefaultConnection
	row := s.db.QueryRow("SELECT id, provider, subject_local, subject_remote FROM hydra_oauth_link WHERE provider=$1 AND subject_remote=$2 LIMIT 1", provider, subject)
	if err := row.Scan(&c.ID, &c.Provider, &c.LocalSubject, &c.RemoteSubject); err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *Store) FindAllByLocalSubject(subject string) (cs []Connection, err error) {
	rows, err := s.db.Query("SELECT id, provider, subject_local, subject_remote FROM hydra_oauth_link WHERE subject_local=$1", subject)
	if err != nil {
		return []Connection{}, err
	}

	defer rows.Close()
	for rows.Next() {
		var c DefaultConnection
		if err := rows.Scan(&c.ID, &c.Provider, &c.LocalSubject, &c.RemoteSubject); err == sql.ErrNoRows {
			return []Connection{}, ErrNotFound
		} else if err != nil {
			return []Connection{}, err
		}
		cs = append(cs, &c)
	}

	return cs, nil
}
