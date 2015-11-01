package postgres

import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
	"github.com/go-errors/errors"
	"github.com/ory-am/hydra/account"
	"github.com/ory-am/hydra/hash"
)

const accountSchema = `CREATE TABLE IF NOT EXISTS account (
	id           text NOT NULL PRIMARY KEY,
	email 		 text NOT NULL UNIQUE,
	password     text NOT NULL,
	data		 json
)`

var ErrNotFound = errors.New("Not found")

type Store struct {
	hasher hash.Hasher
	db     *sql.DB
}

func New(h hash.Hasher, db *sql.DB) *Store {
	return &Store{h, db}
}

func (s *Store) CreateSchemas() error {
	if _, err := s.db.Exec(accountSchema); err != nil {
		log.Warnf("Error creating schema %s: %s", accountSchema, err)
		return err
	}
	return nil
}

func (s *Store) Create(id, email, password, data string) (account.Account, error) {
	// Hash the password
	password, err := s.hasher.Hash(password)
	if err != nil {
		return nil, err
	}

	// Execute SQL statement
	_, err = s.db.Exec("INSERT INTO account (id, email, password, data) VALUES ($1, $2, $3, $4)", id, email, password, data)
	if err != nil {
		return nil, err
	}

	return &account.DefaultAccount{id, email, password, data}, nil
}

func (s *Store) Get(id string) (account.Account, error) {
	var a account.DefaultAccount
	// Query account
	row := s.db.QueryRow("SELECT id, email, password, data FROM account WHERE id=$1 LIMIT 1", id)

	// Hydrate struct with data
	if err := row.Scan(&a.ID, &a.Email, &a.Password, &a.Data); err == sql.ErrNoRows {
		return nil, account.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *Store) UpdatePassword(id, oldPassword, newPassword string) (account.Account, error) {
	acc, err := s.authenticateWithIDAndPassword(id, oldPassword)
	if err != nil {
		return nil, err
	}

	// Hash the new password
	newPassword, err = s.hasher.Hash(newPassword)
	if err != nil {
		return nil, err
	}

	// Execute SQL statement
	if _, err = s.db.Exec("UPDATE account SET (password) = ($2) WHERE id=$1", id, newPassword); err != nil {
		return nil, err
	}

	return &account.DefaultAccount{acc.GetID(), acc.GetEmail(), newPassword, acc.GetData()}, nil
}

func (s *Store) UpdateEmail(id, email, password string) (account.Account, error) {
	acc, err := s.authenticateWithIDAndPassword(id, password)
	if err != nil {
		return nil, err
	}

	// Execute SQL statement
	if _, err = s.db.Exec("UPDATE account SET (email) = ($2) WHERE id=$1", id, email); err != nil {
		return nil, err
	}

	return &account.DefaultAccount{acc.GetID(), email, acc.GetEmail(), acc.GetData()}, nil
}

func (s *Store) Delete(id string) (err error) {
	_, err = s.db.Exec("DELETE FROM account WHERE id=$1", id)
	return err
}

func (s *Store) Authenticate(email, password string) (account.Account, error) {
	var a account.DefaultAccount
	// Query account
	row := s.db.QueryRow("SELECT id, email, password, data FROM account WHERE email=$1", email)

	// Hydrate struct with data
	if err := row.Scan(&a.ID, &a.Email, &a.Password, &a.Data); err != nil {
		return nil, err
	}

	// Compare the given password with the hashed password stored in the database
	if err := s.hasher.Compare(a.Password, password); err != nil {
		return nil, err
	}

	return &a, nil
}

func (s *Store) UpdateData(id string, data string) (account.Account, error) {
	// Execute SQL statement
	if _, err := s.db.Exec("UPDATE account SET (data) = ($2) WHERE id=$1", id, data); err != nil {
		return nil, err
	}

	return s.Get(id)
}

func (s *Store) authenticateWithIDAndPassword(id, password string) (account.Account, error) {
	// Look up account
	acc, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	// Compare the given password with the hashed password stored in the database
	if err := s.hasher.Compare(acc.GetPassword(), password); err != nil {
		return nil, err
	}

	return acc, nil
}
