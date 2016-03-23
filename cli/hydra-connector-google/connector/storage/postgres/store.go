package postgres

import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
	"github.com/go-errors/errors"
	. "github.com/ory-am/common/pkg"
	. "github.com/ory-am/hydra/endpoint/connector/storage"
)

const accountSchema = `CREATE TABLE IF NOT EXISTS hydra_state_data (
	id           text NOT NULL PRIMARY KEY,
	client_id	 text NOT NULL,
	redirect_uri text NOT NULL,
	scope 		 text NOT NULL,
	state 		 text NOT NULL,
	type 		 text NOT NULL,
	provider 	 text NOT NULL,
	expires_at	 timestamp with time zone NOT NULL
)`

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{db}
}

func (s *Store) CreateSchemas() error {
	if _, err := s.db.Exec(accountSchema); err != nil {
		log.Warnf("Error creating schema %s: %s", accountSchema, err)
		return errors.New(err)
	}
	return nil
}

func (s *Store) SaveStateData(sd *StateData) error {
	// Execute SQL statement
	_, err := s.db.Exec("INSERT INTO hydra_state_data (id, client_id, redirect_uri, scope, state, type, provider, expires_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", sd.ID, sd.ClientID, sd.RedirectURL, sd.Scope, sd.State, sd.Type, sd.Provider, sd.ExpiresAt)
	if err != nil {
		return errors.New(err)
	}
	return nil
}

func (s *Store) GetStateData(id string) (*StateData, error) {
	var sd StateData
	row := s.db.QueryRow("SELECT id, client_id, redirect_uri, scope, state, type, provider, expires_at FROM hydra_state_data WHERE id=$1 LIMIT 1", id)

	if err := row.Scan(&sd.ID, &sd.ClientID, &sd.RedirectURL, &sd.Scope, &sd.State, &sd.Type, &sd.Provider, &sd.ExpiresAt); err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, errors.New(err)
	}
	return &sd, nil
}
