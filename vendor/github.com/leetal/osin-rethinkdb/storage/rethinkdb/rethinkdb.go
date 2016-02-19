// Copyright (C) 2016 Alexander Widerberg.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package rethinkdb is a osin storage implementation for rethinkdb.
package rethinkdb

import (
	"fmt"
	"time"

	"github.com/RangelReale/osin"
	rdb "github.com/dancannon/gorethink"
	"github.com/go-errors/errors"
	"github.com/ory-am/common/pkg"
)

const (
	clientsTable      = "oauth_clients"
	authorizeTable    = "oauth_authorize_data"
	accessTable       = "oauth_access_data"
	accessTokenField  = "accessToken"
	refreshTokenField = "refreshToken"
)

/***************************** OSIN Overloads *********************************/

// rethinkClient - Overload of osin.Client so we can add tags for gorethink
type rethinkClient struct {
	ID          string `gorethink:"id"`
	Secret      string `gorethink:"secret"`
	RedirectURI string `gorethink:"redirectUri"`
	UserData    string `gorethink:"userData"`
}

// rethinkAuthorizeData - Overload of osin.AuthorizeData so we can add tags for gorethink
type rethinkAuthorizeData struct {
	// Client information
	Client rethinkClient `gorethink:"client"`

	// Authorization code
	Code string `gorethink:"code"`

	// Token expiration in seconds
	ExpiresIn int32 `gorethink:"expiresIn"`

	// Requested scope
	Scope string `gorethink:"scope"`

	// Redirect Uri from request
	RedirectURI string `gorethink:"redirectURI"`

	// State data from request
	State string `gorethink:"state"`

	// Date created
	CreatedAt time.Time `gorethink:"createdAt"`

	// Data to be passed to storage. Not used by the library.
	UserData string `gorethink:"userData"`
}

// rethinkAccessData - Overload of osin.AccessData so we can add tags for gorethink
type rethinkAccessData struct {
	// Client information
	Client rethinkClient `gorethink:"client"`

	// Authorize data, for authorization code
	AuthorizeData rethinkAuthorizeData `gorethink:"authorizeData"`

	// Previous access data, for refresh token
	PreviousAccessDataToken string `gorethink:"previousAccessDataToken"`

	// Access token
	AccessToken string `gorethink:"accessToken"`

	// Refresh Token. Can be blank
	RefreshToken string `gorethink:"refreshToken"`

	// Token expiration in seconds
	ExpiresIn int32 `gorethink:"expiresIn"`

	// Requested scope
	Scope string `gorethink:"scope"`

	// Redirect Uri from request
	RedirectURI string `gorethink:"redirectURI"`

	// Date created
	CreatedAt time.Time `gorethink:"createdAt"`

	// Data to be passed to storage. Not used by the library.
	UserData string `gorethink:"userData"`
}

/******************************************************************************/

type Storage struct {
	session *rdb.Session
}

// New returns a new rethinkdb storage instance.
func New(session *rdb.Session) *Storage {
	return &Storage{session}
}

// CreateTables creates the tables. Returns an error if something went wrong.
func (s *Storage) CreateTables() error {

	exists, err := s.tableExists(clientsTable)
	if err == nil && !exists {
		_, err := rdb.TableCreate(clientsTable).RunWrite(s.session)
		if err != nil {
			return errors.New(err)
		}
	}

	exists, err = s.tableExists(authorizeTable)
	if err == nil && !exists {
		_, err = rdb.TableCreate(authorizeTable).RunWrite(s.session)
		if err != nil {
			return errors.New(err)
		}

		// Setup secondary indexes (accessTokenField)
		res, err := rdb.Table(authorizeTable).IndexCreateFunc("code", func(row rdb.Term) interface{} {
			return row.Field("code")
		}).RunWrite(s.session)
		if err != nil {
			return errors.New(err)
		} else if res.Errors > 0 {
			return errors.New(res.FirstError)
		}
	}

	exists, err = s.tableExists(accessTable)
	if err == nil && !exists {
		_, err = rdb.TableCreate(accessTable).RunWrite(s.session)
		if err != nil {
			return errors.New(err)
		}

		// Setup secondary indexes (accessTokenField)
		res, err := rdb.Table(accessTable).IndexCreateFunc(accessTokenField, func(row rdb.Term) interface{} {
			return row.Field(accessTokenField)
		}).RunWrite(s.session)
		if err != nil {
			return errors.New(err)
		} else if res.Errors > 0 {
			return errors.New(res.FirstError)
		}

		// Setup secondary indexes (refreshTokenField)
		res, err = rdb.Table(accessTable).IndexCreateFunc(refreshTokenField, func(row rdb.Term) interface{} {
			return row.Field(refreshTokenField)
		}).RunWrite(s.session)
		if err != nil {
			return errors.New(err)
		} else if res.Errors > 0 {
			return errors.New(res.FirstError)
		}
	}

	return nil
}

// TableExists check if table(s) exists in database
func (s *Storage) tableExists(table string) (bool, error) {

	res, err := rdb.TableList().Run(s.session)
	if err != nil {
		return false, errors.New(err)
	} else if res.IsNil() {
		return false, nil
	}

	defer res.Close()

	var tableDB string
	for res.Next(&tableDB) {
		if table == tableDB {
			return true, nil
		}
	}

	return false, nil
}

// Clone the storage if needed.
func (s *Storage) Clone() osin.Storage {
	return s
}

// Close the resources the Storage potentially holds
func (s *Storage) Close() {}

// CreateClient inserts a new client
func (s *Storage) CreateClient(c osin.Client) error {
	// Make sure that the extra data is in correct format
	userdata, err := assertToString(c.GetUserData())
	if err != nil {
		return errors.New(err)
	}

	r := rethinkClient{
		ID:          c.GetId(),
		Secret:      c.GetSecret(),
		RedirectURI: c.GetRedirectUri(),
		UserData:    userdata,
	}

	res, err := rdb.Table(clientsTable).Insert(r).RunWrite(s.session)
	if err != nil {
		return errors.New(err)
	} else if res.Errors > 0 {
		return errors.New(res.FirstError)
	}
	return nil
}

func (s *Storage) Contains(table, field, value string) (bool, error) {
	res, err := rdb.Table(table).Field(field).Contains(value).Run(s.session)
	if err != nil {
		return false, errors.New(err)
	} else if res.IsNil() {
		return false, pkg.ErrNotFound
	}

	defer res.Close()

	var found bool
	err = res.One(&found)

	if err != nil {
		return false, errors.New(err)
	}

	return found, nil
}

func (s *Storage) internalGetClient(clientID string) (rethinkClient, error) {
	result, err := rdb.Table(clientsTable).Get(clientID).Run(s.session)
	if err != nil {
		return rethinkClient{}, errors.New(err)
	} else if result.IsNil() {
		return rethinkClient{}, pkg.ErrNotFound
	}

	defer result.Close()

	var clientStruct rethinkClient
	err = result.One(&clientStruct)
	if err != nil {
		return rethinkClient{}, errors.New(err)
	}

	return clientStruct, nil
}

func (s *Storage) convertInternalClientToExternal(internal rethinkClient) osin.Client {
	return &osin.DefaultClient{
		Id:          internal.ID,
		Secret:      internal.Secret,
		RedirectUri: internal.RedirectURI,
		UserData:    internal.UserData,
	}
}

// GetClient returns client with given ID
func (s *Storage) GetClient(clientID string) (osin.Client, error) {

	clientStruct, err := s.internalGetClient(clientID)
	if err != nil {
		return nil, err
	}

	return s.convertInternalClientToExternal(clientStruct), nil
}

// UpdateClient updates given client
func (s *Storage) UpdateClient(c osin.Client) error {

	// Make sure that the extra data is in correct format
	userdata, err := assertToString(c.GetUserData())
	if err != nil {
		return errors.New(err)
	}

	r := rethinkClient{
		ID:          c.GetId(),
		Secret:      c.GetSecret(),
		RedirectURI: c.GetRedirectUri(),
		UserData:    userdata,
	}

	res, err := rdb.Table(clientsTable).Get(c.GetId()).Update(r).RunWrite(s.session)
	if err != nil {
		return errors.New(err)
	} else if res.Errors > 0 {
		return errors.New(res.FirstError)
	}

	return nil
}

// RemoveClient deletes given client
func (s *Storage) RemoveClient(id string) error {
	if _, err := rdb.Table(clientsTable).Get(id).Delete().RunWrite(s.session); err != nil {
		return errors.New(err)
	}
	return nil
}

// SaveAuthorize creates a new authorization
func (s *Storage) SaveAuthorize(data *osin.AuthorizeData) error {

	if data.Code == "" {
		return errors.New("Code must be set")
	}

	// Make sure that code is unique
	found, err := s.Contains(authorizeTable, "code", data.Code)

	if err != nil || found {
		return errors.New("Code not unique")
	}

	userdata, err := assertToString(data.Client.GetUserData())
	if err != nil {
		return errors.New(err)
	}

	c := rethinkClient{
		ID:          data.Client.GetId(),
		Secret:      data.Client.GetSecret(),
		RedirectURI: data.Client.GetRedirectUri(),
		UserData:    userdata,
	}

	extra, err := assertToString(data.UserData)
	if err != nil {
		return errors.New(err)
	}

	a := rethinkAuthorizeData{
		Client:      c,
		Code:        data.Code,
		ExpiresIn:   data.ExpiresIn,
		Scope:       data.Scope,
		RedirectURI: data.RedirectUri,
		State:       data.State,
		CreatedAt:   data.CreatedAt,
		UserData:    extra,
	}

	res, err := rdb.Table(authorizeTable).Insert(a).RunWrite(s.session)
	if err != nil {
		return errors.New(err)
	} else if res.Errors > 0 {
		return errors.New(res.FirstError)
	}
	return nil
}

func (s *Storage) internalLoadAuthorize(code string) (rethinkAuthorizeData, error) {
	result, err := rdb.Table(authorizeTable).GetAllByIndex("code", code).Run(s.session)
	if err != nil {
		return rethinkAuthorizeData{}, errors.New(err)
	} else if result.IsNil() {
		return rethinkAuthorizeData{}, pkg.ErrNotFound
	}

	defer result.Close()

	var authorizeStruct rethinkAuthorizeData
	err = result.One(&authorizeStruct)
	if err != nil {
		return rethinkAuthorizeData{}, errors.New(err)
	}

	return authorizeStruct, nil
}

func (s *Storage) convertInternalAuthorizeToExternal(internal rethinkAuthorizeData) *osin.AuthorizeData {

	ret := osin.AuthorizeData{
		Client:      s.convertInternalClientToExternal(internal.Client),
		Code:        internal.Code,
		ExpiresIn:   internal.ExpiresIn,
		Scope:       internal.Scope,
		RedirectUri: internal.RedirectURI,
		State:       internal.State,
		CreatedAt:   internal.CreatedAt,
		UserData:    internal.UserData,
	}

	return &ret
}

// LoadAuthorize gets authorization data with given code
func (s *Storage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {

	authorizeStruct, err := s.internalLoadAuthorize(code)
	if err != nil {
		return nil, err
	}

	ret := s.convertInternalAuthorizeToExternal(authorizeStruct)

	if ret.ExpireAt().Before(time.Now()) {
		return nil, errors.Errorf("Token expired at %s.", ret.ExpireAt().String())
	}

	return ret, nil
}

// RemoveAuthorize deletes given authorization
func (s *Storage) RemoveAuthorize(code string) error {
	if _, err := rdb.Table(authorizeTable).GetAllByIndex("code", code).Delete().RunWrite(s.session); err != nil {
		return errors.New(err)
	}
	return nil
}

// SaveAccess creates a new access data
func (s *Storage) SaveAccess(data *osin.AccessData) error {

	// Construct the Client data
	if data.Client == nil {
		return errors.New("data.Client must not be nil")
	}

	// Make sure that accessTokenField is unique
	found, err := s.Contains(authorizeTable, accessTokenField, data.AccessToken)

	if err != nil || found {
		return errors.New("AccessToken not unique")
	}

	userdata, err := assertToString(data.Client.GetUserData())
	if err != nil {
		return errors.New(err)
	}

	c := rethinkClient{
		ID:          data.Client.GetId(),
		Secret:      data.Client.GetSecret(),
		RedirectURI: data.Client.GetRedirectUri(),
		UserData:    userdata,
	}

	// Build up the return data
	ret := rethinkAccessData{}
	ret.Client = c

	// Fetch the previous access data (if any)
	if data.AccessData != nil {
		ret.PreviousAccessDataToken = data.AccessData.AccessToken
	}

	// Construct the saved authorize data (if any)
	a := rethinkAuthorizeData{}
	if data.AuthorizeData != nil {
		extraAuthorizeUserData, err := assertToString(data.AuthorizeData.UserData)
		if err != nil {
			return errors.New(err)
		}

		a.Code = data.AuthorizeData.Code
		a.ExpiresIn = data.AuthorizeData.ExpiresIn
		a.Scope = data.AuthorizeData.Scope
		a.RedirectURI = data.AuthorizeData.RedirectUri
		a.State = data.AuthorizeData.State
		a.CreatedAt = data.AuthorizeData.CreatedAt
		a.UserData = extraAuthorizeUserData

		// Construct the inner authorizeData Client
		ac := rethinkClient{}
		if data.AuthorizeData.Client != nil {
			extraAuthorizeClientUserdata, err := assertToString(data.AuthorizeData.Client.GetUserData())
			if err != nil {
				return errors.New(err)
			}

			ac.ID = data.AuthorizeData.Client.GetId()
			ac.RedirectURI = data.AuthorizeData.Client.GetRedirectUri()
			ac.Secret = data.AuthorizeData.Client.GetSecret()
			ac.UserData = extraAuthorizeClientUserdata
		}
		a.Client = ac

	}
	ret.AuthorizeData = a

	userData, err := assertToString(data.UserData)
	if err != nil {
		return errors.New(err)
	}

	// Plain types
	ret.AccessToken = data.AccessToken
	ret.RefreshToken = data.RefreshToken
	ret.ExpiresIn = data.ExpiresIn
	ret.Scope = data.Scope
	ret.RedirectURI = data.RedirectUri
	ret.CreatedAt = data.CreatedAt
	ret.UserData = userData

	res, err := rdb.Table(accessTable).Insert(ret).RunWrite(s.session)
	if err != nil {
		return errors.New(err)
	} else if res.Errors > 0 {
		return errors.New(res.FirstError)
	}
	return nil
}

// LoadAccess gets access data with given access token
func (s *Storage) LoadAccess(accessToken string) (*osin.AccessData, error) {
	return s.getAccessData(accessTokenField, accessToken)
}

// RemoveAccess deletes AccessData with given access token
func (s *Storage) RemoveAccess(accessToken string) error {
	return s.removeAccessData(accessTokenField, accessToken)
}

// LoadRefresh gets access data with given refresh token
func (s *Storage) LoadRefresh(refreshToken string) (*osin.AccessData, error) {
	return s.getAccessData(refreshTokenField, refreshToken)
}

// RemoveRefresh deletes AccessData with given refresh token
func (s *Storage) RemoveRefresh(refreshToken string) error {
	return s.removeAccessData(refreshTokenField, refreshToken)
}

func (s *Storage) convertInternalAccessToExternal(internal *rethinkAccessData) *osin.AccessData {
	return &osin.AccessData{
		Client:        s.convertInternalClientToExternal(internal.Client),
		AuthorizeData: s.convertInternalAuthorizeToExternal(internal.AuthorizeData),
		AccessToken:   internal.AccessToken,
		RefreshToken:  internal.RefreshToken,
		ExpiresIn:     internal.ExpiresIn,
		Scope:         internal.Scope,
		RedirectUri:   internal.RedirectURI,
		CreatedAt:     internal.CreatedAt,
		UserData:      internal.UserData,
	}
}

// getAccessData is a common function to get AccessData by field
func (s *Storage) getAccessData(fieldName, token string) (*osin.AccessData, error) {
	result, err := rdb.Table(accessTable).GetAllByIndex(fieldName, token).Run(s.session)
	if err != nil {
		return nil, errors.New(err)
	} else if result.IsNil() {
		return nil, pkg.ErrNotFound
	}

	defer result.Close()

	var accessDataStruct rethinkAccessData
	err = result.One(&accessDataStruct)
	if err != nil {
		return nil, errors.New(err)
	}

	ret := s.convertInternalAccessToExternal(&accessDataStruct)

	// Recursively fetch the chain of access data
	prevAccess, _ := s.LoadAccess(accessDataStruct.PreviousAccessDataToken)
	ret.AccessData = prevAccess

	return ret, nil
}

// removeAccessData is a common function to remove AccessData by field
func (s *Storage) removeAccessData(fieldName, token string) error {
	if _, err := rdb.Table(accessTable).GetAllByIndex(fieldName, token).Delete().RunWrite(s.session); err != nil {
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
