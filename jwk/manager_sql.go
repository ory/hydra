package jwk

import (
	"database/sql"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/ory-am/hydra/pkg"
	"github.com/pkg/errors"
	"github.com/square/go-jose"
)

type SQLManager struct {
	DB     *sqlx.DB
	Cipher *AEAD
}

var sqlSchema = []string{
	`CREATE TABLE IF NOT EXISTS hydra_jwk (
	sid     varchar(255) NOT NULL,
	kid 	varchar(255) NOT NULL,
	version int NOT NULL DEFAULT 0,
	keydata text NOT NULL,
	PRIMARY KEY (sid, kid)
)`,
}

type sqlData struct {
	Set     string `db:"sid"`
	KID     string `db:"kid"`
	Version int    `db:"version"`
	Key     string `db:"keydata"`
}

func (s *SQLManager) CreateSchemas() error {
	for _, query := range sqlSchema {
		if _, err := s.DB.Exec(query); err != nil {
			return errors.Wrapf(err, "Could not create schema:\n%s", query)
		}
	}
	return nil
}

func (m *SQLManager) AddKey(set string, key *jose.JsonWebKey) error {
	out, err := json.Marshal(key)
	if err != nil {
		return errors.Wrap(err, "")
	}

	encrypted, err := m.Cipher.Encrypt(out)
	if err != nil {
		return errors.Wrap(err, "")
	}

	if _, err = m.DB.NamedExec(`INSERT INTO hydra_jwk (sid, kid, version, keydata) VALUES (:sid, :kid, :version, :keydata)`, &sqlData{
		Set:     set,
		KID:     key.KeyID,
		Version: 0,
		Key:     encrypted,
	}); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func (m *SQLManager) AddKeySet(set string, keys *jose.JsonWebKeySet) error {
	tx, err := m.DB.Beginx()
	if err != nil {
		return errors.Wrap(err, "")
	}

	for _, key := range keys.Keys {
		out, err := json.Marshal(key)
		if err != nil {
			return errors.Wrap(err, "")
		}

		encrypted, err := m.Cipher.Encrypt(out)
		if err != nil {
			return errors.Wrap(err, "")
		}

		if _, err = tx.NamedExec(`INSERT INTO hydra_jwk (sid, kid, version, keydata) VALUES (:sid, :kid, :version, :keydata)`, &sqlData{
			Set:     set,
			KID:     key.KeyID,
			Version: 0,
			Key:     encrypted,
		}); err != nil {
			return errors.Wrap(err, "")
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func (m *SQLManager) GetKey(set, kid string) (*jose.JsonWebKeySet, error) {
	var d sqlData
	if err := m.DB.Get(&d, m.DB.Rebind("SELECT * FROM hydra_jwk WHERE sid=? AND kid=?"), set, kid); err == sql.ErrNoRows {
		return nil, errors.Wrap(pkg.ErrNotFound, "")
	} else if err != nil {
		return nil, errors.Wrap(err, "")
	}

	key, err := m.Cipher.Decrypt(d.Key)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	var c jose.JsonWebKey
	if err := json.Unmarshal(key, &c); err != nil {
		return nil, errors.Wrap(err, "")
	}

	return &jose.JsonWebKeySet{
		Keys: []jose.JsonWebKey{c},
	}, nil
}

func (m *SQLManager) GetKeySet(set string) (*jose.JsonWebKeySet, error) {
	var ds []sqlData
	if err := m.DB.Select(&ds, m.DB.Rebind("SELECT * FROM hydra_jwk WHERE sid=?"), set); err == sql.ErrNoRows {
		return nil, errors.Wrap(pkg.ErrNotFound, "")
	} else if err != nil {
		return nil, errors.Wrap(err, "")
	}

	if len(ds) == 0 {
		return nil, errors.Wrap(pkg.ErrNotFound, "")
	}

	keys := &jose.JsonWebKeySet{Keys: []jose.JsonWebKey{}}
	for _, d := range ds {
		key, err := m.Cipher.Decrypt(d.Key)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}

		var c jose.JsonWebKey
		if err := json.Unmarshal(key, &c); err != nil {
			return nil, errors.Wrap(err, "")
		}
		keys.Keys = append(keys.Keys, c)
	}

	return keys, nil
}

func (m *SQLManager) DeleteKey(set, kid string) error {
	if _, err := m.DB.Exec(m.DB.Rebind(`DELETE FROM hydra_jwk WHERE sid=? AND kid=?`), set, kid); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func (m *SQLManager) DeleteKeySet(set string) error {
	if _, err := m.DB.Exec(m.DB.Rebind(`DELETE FROM hydra_jwk WHERE sid=?`), set); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}
