// Package postgres is a osin storage implementation for postgres.
package postgres

import (
	"database/sql"
	"github.com/RangelReale/osin"
	"github.com/go-errors/errors"
	"github.com/ory-am/common/pkg"
	"log"
	"fmt"
)

var schemas = []string{`CREATE TABLE IF NOT EXISTS client (
	id           text NOT NULL PRIMARY KEY,
	secret 		 text NOT NULL,
	extra 		 text NOT NULL,
	redirect_uri text NOT NULL
)`, `CREATE TABLE IF NOT EXISTS authorize (
	client       text NOT NULL,
	code         text NOT NULL PRIMARY KEY,
	expires_in   int NOT NULL,
	scope        text NOT NULL,
	redirect_uri text NOT NULL,
	state        text NOT NULL,
	extra 		 text NOT NULL,
	created_at   timestamp with time zone NOT NULL
)`, `CREATE TABLE IF NOT EXISTS access (
	client        text NOT NULL,
	authorize     text NOT NULL,
	previous      text NOT NULL,
	access_token  text NOT NULL PRIMARY KEY,
	refresh_token text NOT NULL,
	expires_in    int NOT NULL,
	scope         text NOT NULL,
	redirect_uri  text NOT NULL,
	extra 		  text NOT NULL,
	created_at    timestamp with time zone NOT NULL
)`, `CREATE TABLE IF NOT EXISTS refresh (
	token         text NOT NULL PRIMARY KEY,
	access        text NOT NULL
)`}

// Storage implements interface "github.com/RangelReale/osin".Storage and interface "github.com/ory-am/osin-storage".Storage
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

// Clone the storage if needed. For example, using mgo, you can clone the session with session.Clone
// to avoid concurrent access problems.
// This is to avoid cloning the connection at each method access.
// Can return itself if not a problem.
func (s *Storage) Clone() osin.Storage {
	return s
}

// Close the resources the Storage potentially holds (using Clone for example)
func (s *Storage) Close() {
}

// GetClient loads the client by id
func (s *Storage) GetClient(id string) (osin.Client, error) {
	row := s.db.QueryRow("SELECT id, secret, redirect_uri, extra FROM client WHERE id=$1", id)
	var c osin.DefaultClient
	var extra string

	if err := row.Scan(&c.Id, &c.Secret, &c.RedirectUri, &extra); err == sql.ErrNoRows {
		return nil, pkg.ErrNotFound
	} else if err != nil {
		return nil, errors.New(err)
	}
	c.UserData = extra
	return &c, nil
}

// UpdateClient updates the client (identified by it's id) and replaces the values with the values of client.
func (s *Storage) UpdateClient(c osin.Client) error {
	data, err := assertToString(c.GetUserData())
	if err != nil {
		return err
	}

	if _, err := s.db.Exec("UPDATE client SET (secret, redirect_uri, extra) = ($2, $3, $4) WHERE id=$1", c.GetId(), c.GetSecret(), c.GetRedirectUri(), data); err != nil {
		return errors.New(err)
	}
	return nil
}

// CreateClient stores the client in the database and returns an error, if something went wrong.
func (s *Storage) CreateClient(c osin.Client) error {
	data, err := assertToString(c.GetUserData())
	if err != nil {
		return err
	}

	if _, err := s.db.Exec("INSERT INTO client (id, secret, redirect_uri, extra) VALUES ($1, $2, $3, $4)", c.GetId(), c.GetSecret(), c.GetRedirectUri(), data); err != nil {
		return errors.New(err)
	}
	return nil
}

// RemoveClient removes a client (identified by id) from the database. Returns an error if something went wrong.
func (s *Storage) RemoveClient(id string) (err error) {
	if _, err = s.db.Exec("DELETE FROM client WHERE id=$1", id); err != nil {
		return errors.New(err)
	}
	return nil
}

// SaveAuthorize saves authorize data.
func (s *Storage) SaveAuthorize(data *osin.AuthorizeData) (err error) {
	extra, err := assertToString(data.UserData)
	if err != nil {
		return err
	}

	if _, err = s.db.Exec(
		"INSERT INTO authorize (client, code, expires_in, scope, redirect_uri, state, created_at, extra) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		data.Client.GetId(),
		data.Code,
		data.ExpiresIn,
		data.Scope,
		data.RedirectUri,
		data.State,
		data.CreatedAt,
		extra,
	); err != nil {
		return errors.New(err)
	}
	return nil
}

// LoadAuthorize looks up AuthorizeData by a code.
// Client information MUST be loaded together.
// Optionally can return error if expired.
func (s *Storage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	var data osin.AuthorizeData
	var extra string
	var cid string
	if err := s.db.QueryRow("SELECT client, code, expires_in, scope, redirect_uri, state, created_at, extra FROM authorize WHERE code=$1 LIMIT 1", code).Scan(&cid, &data.Code, &data.ExpiresIn, &data.Scope, &data.RedirectUri, &data.State, &data.CreatedAt, &extra); err == sql.ErrNoRows {
		return nil, pkg.ErrNotFound
	} else if err != nil {
		return nil, errors.New(err)
	}
	data.UserData = extra

	c, err := s.GetClient(cid)
	if err != nil {
		return nil, err
	}

	data.Client = c
	return &data, nil
}

// RemoveAuthorize revokes or deletes the authorization code.
func (s *Storage) RemoveAuthorize(code string) (err error) {
	if _, err = s.db.Exec("DELETE FROM authorize WHERE code=$1", code); err != nil {
		return errors.New(err)
	}
	return nil
}

// SaveAccess writes AccessData.
// If RefreshToken is not blank, it must save in a way that can be loaded using LoadRefresh.
func (s *Storage) SaveAccess(data *osin.AccessData) (err error) {
	prev := ""
	authorizeData := &osin.AuthorizeData{}

	if data.AccessData != nil {
		prev = data.AccessData.AccessToken
	}

	if data.AuthorizeData != nil {
		authorizeData = data.AuthorizeData
	}

	extra, err := assertToString(data.UserData)
	if err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return errors.New(err)
	}

	if data.RefreshToken != "" {
		if err := s.saveRefresh(tx, data.RefreshToken, data.AccessToken); err != nil {
			return err
		}
	}

	if data.Client == nil {
		return errors.New("data.Client must not be nil")
	}

	_, err = tx.Exec("INSERT INTO access (client, authorize, previous, access_token, refresh_token, expires_in, scope, redirect_uri, created_at, extra) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", data.Client.GetId(), authorizeData.Code, prev, data.AccessToken, data.RefreshToken, data.ExpiresIn, data.Scope, data.RedirectUri, data.CreatedAt, extra)
	if err != nil {
		if rbe := tx.Rollback(); rbe != nil {
			return errors.New(rbe)
		}
		return errors.New(err)
	}

	if err = tx.Commit(); err != nil {
		return errors.New(err)
	}

	return nil
}

// LoadAccess retrieves access data by token. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (s *Storage) LoadAccess(code string) (*osin.AccessData, error) {
	var extra, cid, prevAccessToken, authorizeCode string
	var result osin.AccessData

	if err := s.db.QueryRow(
		"SELECT client, authorize, previous, access_token, refresh_token, expires_in, scope, redirect_uri, created_at, extra FROM access WHERE access_token=$1 LIMIT 1",
		code,
	).Scan(
		&cid,
		&authorizeCode,
		&prevAccessToken,
		&result.AccessToken,
		&result.RefreshToken,
		&result.ExpiresIn,
		&result.Scope,
		&result.RedirectUri,
		&result.CreatedAt,
		&extra,
	); err == sql.ErrNoRows {
		return nil, pkg.ErrNotFound
	} else if err != nil {
		return nil, errors.New(err)
	}
	result.UserData = extra

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

	return &result, nil
}

// RemoveAccess revokes or deletes an AccessData.
func (s *Storage) RemoveAccess(code string) (err error) {
	_, err = s.db.Exec("DELETE FROM access WHERE access_token=$1", code)
	if err != nil {
		return errors.New(err)
	}
	return nil
}

// LoadRefresh retrieves refresh AccessData. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (s *Storage) LoadRefresh(code string) (*osin.AccessData, error) {
	row := s.db.QueryRow("SELECT access FROM refresh WHERE token=$1 LIMIT 1", code)
	var access string
	if err := row.Scan(&access); err == sql.ErrNoRows {
		return nil, pkg.ErrNotFound
	} else if err != nil {
		return nil, errors.New(err)
	}
	return s.LoadAccess(access)
}

// RemoveRefresh revokes or deletes refresh AccessData.
func (s *Storage) RemoveRefresh(code string) error {
	_, err := s.db.Exec("DELETE FROM refresh WHERE token=$1", code)
	if err != nil {
		return errors.New(err)
	}
	return nil
}

func (s *Storage) saveRefresh(tx *sql.Tx, refresh, access string) (err error) {
	_, err = tx.Exec("INSERT INTO refresh (token, access) VALUES ($1, $2)", refresh, access)
	if err != nil {
		if rbe := tx.Rollback(); rbe != nil {
			return errors.New(rbe)
		}
		return errors.New(err)
	}
	return nil
}

func assertToString(in interface{}) (string, error) {
	var ok bool
	var data string
	if in == nil {
		return "", nil
	} else if data, ok = in.(string); ok {
		return data, nil
	} else if str, ok := in.(fmt.Stringer); ok {
		return str.String(), nil
	}
	return "", errors.Errorf(`Could not assert "%v" to string`, in)
}
