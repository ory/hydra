package rethinkdb

import (
	rdb "github.com/dancannon/gorethink"
	"github.com/mitchellh/mapstructure"

	"github.com/asaskevich/govalidator"
	"github.com/go-errors/errors"
	"github.com/ory-am/common/pkg"
	"github.com/ory-am/hydra/account"
	"github.com/ory-am/hydra/hash"
	"github.com/pborman/uuid"
)

const accountsTable = "hydra_account"

type Store struct {
	hasher  hash.Hasher
	session *rdb.Session
}

func New(h hash.Hasher, session *rdb.Session) *Store {
	return &Store{h, session}
}

func (s *Store) CreateTables() error {
	exists, err := s.tableExists(accountsTable)
	if err == nil && !exists {
		_, err := rdb.TableCreate(accountsTable).RunWrite(s.session)
		if err != nil {
			return errors.New(err)
		}
	}
	return nil
}

// TableExists check if table(s) exists in database
func (s *Store) tableExists(table string) (bool, error) {

	res, err := rdb.TableList().Run(s.session)
	if err != nil {
		return false, errors.New(err)
	}

	if res.IsNil() {
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
		return nil, errors.New(err)
	}

	// Make sure that username is unique
	found, err := s.Contains("username", r.Username)

	if err != nil || found {
		return nil, errors.New("Username not unique")
	}

	acc := account.DefaultAccount{
		ID:       r.ID,
		Username: r.Username,
		Password: r.Password,
		Data:     r.Data,
	}

	res, err := rdb.Table(accountsTable).Insert(acc).RunWrite(s.session)

	if err != nil {
		return nil, errors.New(err)
	} else if res.Errors > 0 {
		return nil, errors.New(res.FirstError)
	}

	return &acc, nil
}

func (s *Store) Contains(field string, value string) (bool, error) {
	res, err := rdb.Table(accountsTable).Field(field).Contains(value).Run(s.session)
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

func (s *Store) Get(id string) (account.Account, error) {
	// Query account
	result, err := rdb.Table(accountsTable).Get(id).Run(s.session)
	defer result.Close()

	if err != nil {
		return nil, errors.New(err)
	} else if result.IsNil() {
		return nil, pkg.ErrNotFound
	}

	var a account.DefaultAccount
	err = result.One(&a)
	if err != nil {
		return nil, errors.New(err)
	}

	return &a, nil
}

func (s *Store) UpdatePassword(id string, r account.UpdatePasswordRequest) (acc account.Account, err error) {
	if err := validate(r); err != nil {
		return nil, err
	}

	if acc, err = s.authenticateWithIDAndPassword(id, r.CurrentPassword); err != nil {
		return nil, errors.New(err)
	}

	// Hash the new password
	r.NewPassword, err = s.hasher.Hash(r.NewPassword)
	if err != nil {
		return nil, errors.New(err)
	}

	account := account.DefaultAccount{
		ID:       acc.GetID(),
		Username: acc.GetUsername(),
		Password: r.NewPassword,
		Data:     acc.GetData(),
	}

	if _, err = rdb.Table(accountsTable).Get(id).Update(account).RunWrite(s.session); err != nil {
		return nil, errors.New(err)
	}

	return &account, nil
}

func (s *Store) UpdateUsername(id string, r account.UpdateUsernameRequest) (acc account.Account, err error) {
	if err := validate(r); err != nil {
		return nil, err
	}

	if acc, err = s.authenticateWithIDAndPassword(id, r.Password); err != nil {
		return nil, errors.New(err)
	}

	account := account.DefaultAccount{
		ID:       acc.GetID(),
		Username: r.Username,
		Password: acc.GetPassword(),
		Data:     acc.GetData(),
	}

	// Execute SQL statement
	if _, err = rdb.Table(accountsTable).Get(id).Update(account).RunWrite(s.session); err != nil {
		return nil, errors.New(err)
	}

	return &account, nil
}

func (s *Store) Delete(id string) (err error) {
	if _, err = rdb.Table(accountsTable).Get(id).Delete().RunWrite(s.session); err != nil {
		return errors.New(err)
	}
	return nil
}

func (s *Store) Authenticate(username, password string) (account.Account, error) {
	var a account.DefaultAccount
	// Query account
	result, err := rdb.Table(accountsTable).Filter(rdb.Row.Field("username").Eq(username)).Run(s.session)
	if err != nil {
		return nil, errors.New(err)
	}
	defer result.Close()

	var accountMap map[string]interface{}
	err = result.One(&accountMap)
	if err != nil {
		return nil, errors.New(err)
	}

	err = mapstructure.Decode(accountMap, &a)
	if err != nil {
		return nil, errors.New(err)
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

	if acc, err = s.Get(id); err != nil {
		return nil, err
	}

	account := account.DefaultAccount{
		ID:       acc.GetID(),
		Username: acc.GetUsername(),
		Password: acc.GetPassword(),
		Data:     r.Data,
	}

	// Execute SQL statement
	if _, err = rdb.Table(accountsTable).Get(id).Update(account).RunWrite(s.session); err != nil {
		return nil, errors.New(err)
	}

	return &account, nil
}

func (s *Store) authenticateWithIDAndPassword(id, password string) (account.Account, error) {
	// Look up account
	acc, err := s.Get(id)

	if err != nil {
		return nil, err
	}

	// Compare the given password with the hashed password stored in the database
	if err := s.hasher.Compare(acc.GetPassword(), password); err != nil {
		return nil, pkg.ErrInvalidPayload
	}

	return acc, nil
}
