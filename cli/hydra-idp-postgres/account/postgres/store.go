package postgres

import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
	"github.com/asaskevich/govalidator"
	"github.com/go-errors/errors"
	"github.com/ory-am/common/pkg"
	"github.com/ory-am/hydra/account"
	"github.com/ory-am/hydra/hash"
	"github.com/pborman/uuid"
)

const accountSchema = `CREATE TABLE IF NOT EXISTS hydra_account (
	id           text NOT NULL PRIMARY KEY,
	username	 text NOT NULL UNIQUE,
	password     text NOT NULL,
	data		 json
)`

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
		return errors.New(err)
	}
	return nil
}

func validate(r interface{}) error {
	if v, err := govalidator.ValidateStruct(r); !v {
		return pkg.ErrInvalidPayload
	} else if err != nil {
		return pkg.ErrInvalidPayload
	}
	return nil
}

func (s *Store) Create(r account.CreateAccountRequest) (account.Account, error) {
	var err error

	if r.ID == "" {
		r.ID = uuid.New()
	}

	if r.Data == "" {
		r.Data = "{}"
	}

	if err := validate(r); err != nil {
		return nil, err
	}

	// Hash the password
	if r.Password, err = s.hasher.Hash(r.Password); err != nil {
		return nil, err
	}

	if _, err = s.db.Exec("INSERT INTO hydra_account (id, username, password, data) VALUES ($1, $2, $3, $4)", r.ID, r.Username, r.Password, r.Data); err != nil {
		return nil, errors.New(err)
	}

	return &account.DefaultAccount{
		ID:       r.ID,
		Username: r.Username,
		Password: r.Password,
		Data:     r.Data,
	}, nil
}

func (s *Store) Get(id string) (account.Account, error) {
	var a account.DefaultAccount
	// Query account
	row := s.db.QueryRow("SELECT id, username, password, data FROM hydra_account WHERE id=$1 LIMIT 1", id)

	// Hydrate struct with data
	if err := row.Scan(&a.ID, &a.Username, &a.Password, &a.Data); err == sql.ErrNoRows {
		return nil, pkg.ErrNotFound
	} else if err != nil {
		return nil, errors.New(err)
	}
	return &a, nil
}

func (s *Store) UpdatePassword(id string, r account.UpdatePasswordRequest) (acc account.Account, err error) {
	if err := validate(r); err != nil {
		return nil, err
	}

	if acc, err = s.authenticateWithIDAndPassword(id, r.CurrentPassword); err != nil {
		return nil, err
	}

	// Hash the new password
	r.NewPassword, err = s.hasher.Hash(r.NewPassword)
	if err != nil {
		return nil, err
	}

	if _, err = s.db.Exec("UPDATE hydra_account SET (password) = ($2) WHERE id=$1", id, r.NewPassword); err != nil {
		return nil, errors.New(err)
	}

	return &account.DefaultAccount{
		ID:       acc.GetID(),
		Username: acc.GetUsername(),
		Password: r.NewPassword,
		Data:     acc.GetData(),
	}, nil
}

func (s *Store) UpdateUsername(id string, r account.UpdateUsernameRequest) (acc account.Account, err error) {
	if err := validate(r); err != nil {
		return nil, err
	}

	if acc, err = s.authenticateWithIDAndPassword(id, r.Password); err != nil {
		return nil, err
	}

	// Execute SQL statement
	if _, err = s.db.Exec("UPDATE hydra_account SET (username) = ($2) WHERE id=$1", id, r.Username); err != nil {
		return nil, errors.New(err)
	}

	return &account.DefaultAccount{
		ID:       acc.GetID(),
		Username: r.Username,
		Password: acc.GetPassword(),
		Data:     acc.GetData(),
	}, nil
}

func (s *Store) Delete(id string) (err error) {
	if _, err = s.db.Exec("DELETE FROM hydra_account WHERE id=$1", id); err != nil {
		return errors.New(err)
	}
	return nil
}

func (s *Store) Authenticate(username, password string) (account.Account, error) {
	var a account.DefaultAccount
	// Query account
	row := s.db.QueryRow("SELECT id, username, password, data FROM hydra_account WHERE username=$1", username)

	// Hydrate struct with data
	if err := row.Scan(&a.ID, &a.Username, &a.Password, &a.Data); err == sql.ErrNoRows {
		return nil, pkg.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	// Compare the given password with the hashed password stored in the database
	if err := s.hasher.Compare(a.Password, password); err != nil {
		return nil, pkg.ErrInvalidPayload
	}

	return &a, nil
}

func (s *Store) UpdateData(id string, r account.UpdateDataRequest) (acc account.Account, err error) {
	if err := validate(r); err != nil {
		return nil, err
	}

	// Execute SQL statement
	if _, err = s.db.Exec("UPDATE hydra_account SET (data) = ($2) WHERE id=$1", id, r.Data); err != nil {
		return nil, errors.New(err)
	}

	return s.Get(id)
}

func (s *Store) authenticateWithIDAndPassword(id, password string) (account.Account, error) {
	// Look up account
	acc, err := s.Get(id)
	if err != nil {
		return nil, errors.New(err)
	}

	// Compare the given password with the hashed password stored in the database
	if err := s.hasher.Compare(acc.GetPassword(), password); err != nil {
		return nil, pkg.ErrInvalidPayload
	}

	return acc, nil
}
