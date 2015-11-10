// Postgres is a osin storage implementation for postgres.
package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/RangelReale/osin"
	"log"
)

var schemas = []string{`CREATE TABLE IF NOT EXISTS client (
	id           text NOT NULL,
	secret 		 text NOT NULL,
	redirect_uri text NOT NULL,
	extra		 text NOT NULL,

    CONSTRAINT client_pk PRIMARY KEY (id)
)`, `CREATE TABLE IF NOT EXISTS authorize (
	client       text NOT NULL,
	code         text NOT NULL,
	expires_in   int NOT NULL,
	scope        text NOT NULL,
	redirect_uri text NOT NULL,
	state        text NOT NULL,
	extra 		 text NOT NULL,
	created_at   timestamp with time zone NOT NULL,

    CONSTRAINT authorize_pk PRIMARY KEY (code)
)`, `CREATE TABLE IF NOT EXISTS access (
	client        text NOT NULL,
	authorize     text NOT NULL,
	previous      text NOT NULL,
	access_token  text NOT NULL,
	refresh_token text NOT NULL,
	expires_in    int NOT NULL,
	scope         text NOT NULL,
	redirect_uri  text NOT NULL,
	extra 		  text NOT NULL,
	created_at    timestamp with time zone NOT NULL,

    CONSTRAINT access_pk PRIMARY KEY (access_token)
)`, `CREATE TABLE IF NOT EXISTS refresh (
	token         text NOT NULL,
	access        text NOT NULL,

    CONSTRAINT refresh_pk PRIMARY KEY (token)
)`}

type Storage struct {
	db *sql.DB
}

// New returns a new postgres storage instance.
func New(db *sql.DB) *Storage {
	return &Storage{db}
}

// CreateSchemas creates the schemata, if they do not exist yet in the database. Returns an error if something went wrong.
func (s *Storage) CreateSchemas() error {
	for k, schema := range schemas {
		if _, err := s.db.Exec(schema); err != nil {
			log.Printf("Error creating schema %d: %s", k, schema)
			return err
		}
	}
	return nil
}

func (s *Storage) Clone() osin.Storage {
	return s
}

func (s *Storage) Close() {}

func (s *Storage) GetClient(id string) (osin.Client, error) {
	row := s.db.QueryRow("SELECT id, secret, redirect_uri, extra FROM client WHERE id=$1 LIMIT 1", id)
	var c osin.DefaultClient
	var extra string
	if err := row.Scan(&c.Id, &c.Secret, &c.RedirectUri, &extra); err != nil {
		return nil, err
	}
	c.UserData = extra
	return &c, nil
}

func (s *Storage) UpdateClient(c osin.Client) error {
	userData, err := dataToString(c.GetUserData())
	if err != nil {
		return err
	}

	_, err = s.db.Exec("UPDATE client SET (secret, redirect_uri, extra) = ($2, $3, $4) WHERE id=$1", c.GetId(), c.GetSecret(), c.GetRedirectUri(), userData)
	return err
}

func (s *Storage) CreateClient(c osin.Client) error {
	userData, err := dataToString(c.GetUserData())
	if err != nil {
		return err
	}
	_, err = s.db.Exec("INSERT INTO client (id, secret, redirect_uri, extra) VALUES ($1, $2, $3, $4)", c.GetId(), c.GetSecret(), c.GetRedirectUri(), userData)
	return err
}

func (s *Storage) RemoveClient(id string) (err error) {
	_, err = s.db.Exec("DELETE FROM client WHERE id=$1", id)
	return err
}

func (s *Storage) SaveAuthorize(data *osin.AuthorizeData) (err error) {
	userData, err := dataToString(data.UserData)
	if err != nil {
		return err
	}

	_, err = s.db.Exec("INSERT INTO authorize (client, code, expires_in, scope, redirect_uri, state, created_at, extra) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", data.Client.GetId(), data.Code, data.ExpiresIn, data.Scope, data.RedirectUri, data.State, data.CreatedAt, userData)
	return err
}

func (s *Storage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	var data osin.AuthorizeData
	var userData string
	var cid string
	row := s.db.QueryRow("SELECT client, code, expires_in, scope, redirect_uri, state, created_at, extra FROM authorize WHERE code=$1 LIMIT 1", code)
	if err := row.Scan(&cid, &data.Code, &data.ExpiresIn, &data.Scope, &data.RedirectUri, &data.State, &data.CreatedAt, &userData); err != nil {
		return nil, err
	}
	data.UserData = userData

	c, err := s.GetClient(cid)
	if err != nil {
		return nil, err
	}

	data.Client = c
	return &data, nil
}

func (s *Storage) RemoveAuthorize(code string) (err error) {
	_, err = s.db.Exec("DELETE FROM authorize WHERE code=$1", code)
	return err
}

func dataToString(data interface{}) (string, error) {
	if res, ok := data.(string); ok {
		return res, nil
	} else if stringer, ok := data.(fmt.Stringer); ok {
		return stringer.String(), nil
	} else {
		return "", fmt.Errorf("Could not assert to string: %v", data)
	}
}

func (s *Storage) SaveAccess(data *osin.AccessData) (err error) {
	prev := ""
	authorizeData := &osin.AuthorizeData{}

	if data.AccessData != nil {
		prev = data.AccessData.AccessToken
	}

	if data.AuthorizeData != nil {
		authorizeData = data.AuthorizeData
	}

	userData, err := dataToString(data.UserData)
	if err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	if data.RefreshToken != "" {
		if err := s.saveRefresh(tx, data.RefreshToken, data.AccessToken); err != nil {
			return err
		}
	}

	if data.Client == nil {
		return errors.New("data.Client must not be nil.")
	}

	_, err = tx.Exec("INSERT INTO access (client, authorize, previous, access_token, refresh_token, expires_in, scope, redirect_uri, created_at, extra) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", data.Client.GetId(), authorizeData.Code, prev, data.AccessToken, data.RefreshToken, data.ExpiresIn, data.Scope, data.RedirectUri, data.CreatedAt, userData)
	if err != nil {
		if rbe := tx.Rollback(); rbe != nil {
			return rbe
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Storage) LoadAccess(code string) (*osin.AccessData, error) {
	var userData string
	var cid, prevAccessToken, authorizeCode string
	var result osin.AccessData
	row := s.db.QueryRow("SELECT client, authorize, previous, access_token, refresh_token, expires_in, scope, redirect_uri, created_at, extra FROM access WHERE access_token=$1 LIMIT 1", code)
	err := row.Scan(&cid, &authorizeCode, &prevAccessToken, &result.AccessToken, &result.RefreshToken, &result.ExpiresIn, &result.Scope, &result.RedirectUri, &result.CreatedAt, &userData)
	result.UserData = userData

	client, err := s.GetClient(cid)
	if err != nil {
		return nil, err
	}
	result.Client = client

	authorize, err := s.LoadAuthorize(authorizeCode)
	if err != nil {
		return nil, err
	}
	result.AuthorizeData = authorize

	if prevAccessToken != "" {
		prevAccess, err := s.LoadAccess(prevAccessToken)
		if err != nil {
			return nil, err
		}
		result.AccessData = prevAccess
	}

	return &result, err
}

func (s *Storage) RemoveAccess(code string) (err error) {
	st, err := s.db.Prepare("DELETE FROM access WHERE access_token=$1")
	if err != nil {
		return
	}
	_, err = st.Exec(code)
	return err
}

func (s *Storage) saveRefresh(tx *sql.Tx, refresh, access string) (err error) {
	_, err = tx.Exec("INSERT INTO refresh (token, access) VALUES ($1, $2)", refresh, access)
	if err != nil {
		if rbe := tx.Rollback(); rbe != nil {
			return rbe
		}
	}
	return err
}

func (s *Storage) LoadRefresh(code string) (*osin.AccessData, error) {
	row := s.db.QueryRow("SELECT access FROM refresh WHERE token=$1 LIMIT 1", code)
	var access string
	if err := row.Scan(&access); err != nil {
		return nil, err
	}
	return s.LoadAccess(access)
}

func (s *Storage) RemoveRefresh(code string) error {
	st, err := s.db.Prepare("DELETE FROM refresh WHERE token=$1")
	if err != nil {
		return err
	}
	_, err = st.Exec(code)
	return err
}
