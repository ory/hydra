package rethinkdb

import (
	"github.com/go-errors/errors"

	rdb "github.com/dancannon/gorethink"
	"github.com/mitchellh/mapstructure"
	pkg "github.com/ory-am/common/pkg"
	cn "github.com/ory-am/hydra/endpoint/connection"
)

const connectionTable = "hydra_oauth_link"

type rdbFunction func(rdb.Term) rdb.Term

type Store struct {
	session *rdb.Session
}

func New(session *rdb.Session) *Store {
	return &Store{session: session}
}

func (s *Store) CreateTables() error {
	exists, err := s.tableExists(connectionTable)
	if err == nil && !exists {
		_, err := rdb.TableCreate(connectionTable).RunWrite(s.session)
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

func (s *Store) Contains(matchFunction rdbFunction) (bool, error) {

	res, err := rdb.Table(connectionTable).Contains(matchFunction).Run(s.session)
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

func (s *Store) Create(c cn.Connection) error {

	// Create a matching function that looks for "provider, subject_remote" UNIQUE pairs
	fn := func(connection rdb.Term) rdb.Term {
		return rdb.And(connection.Field("provider").Eq(c.GetProvider()),
			connection.Field("remotesubject").Eq(c.GetRemoteSubject()))
	}

	notUnique, err := s.Contains(fn)

	if err != nil || notUnique {
		return errors.New("Duplicate values in provider and subject_remote")
	}

	res, err := rdb.Table(connectionTable).Insert(c).RunWrite(s.session)

	if err != nil {
		return errors.New(err.Error())
	} else if res.Errors > 0 {
		return errors.New(res.FirstError)
	}

	return err
}

func (s *Store) Delete(id string) error {
	if _, err := rdb.Table(connectionTable).Get(id).Delete().RunWrite(s.session); err != nil {
		return errors.New(err)
	}
	return nil
}

func (s *Store) Get(id string) (cn.Connection, error) {
	// Query connection
	result, err := rdb.Table(connectionTable).Get(id).Run(s.session)
	defer result.Close()

	if err != nil {
		return nil, err
	} else if result.IsNil() {
		return nil, pkg.ErrNotFound
	}

	var c cn.DefaultConnection
	err = result.One(&c)
	if err != nil {
		return nil, errors.New(err)
	}

	return &c, nil
}

func (s *Store) FindByRemoteSubject(provider, subject string) (cn.Connection, error) {
	// Query connection
	result, err := rdb.Table(connectionTable).Filter(map[string]interface{}{
		"provider":      provider,
		"remotesubject": subject,
	}).Run(s.session)
	defer result.Close()

	if err != nil {
		return nil, errors.New(err)
	} else if result.IsNil() {
		return nil, pkg.ErrNotFound
	}

	var connectionMap map[string]interface{}
	err = result.One(&connectionMap)
	if err != nil {
		return nil, errors.New(err)
	}

	var c cn.DefaultConnection

	err = mapstructure.Decode(connectionMap, &c)
	if err != nil {
		return nil, errors.New(err)
	}

	return &c, nil
}

func (s *Store) FindAllByLocalSubject(subject string) (cs []cn.Connection, err error) {
	// Query connection
	result, err := rdb.Table(connectionTable).Filter(map[string]interface{}{
		"localsubject": subject,
	}).Run(s.session)
	defer result.Close()

	if err != nil {
		return []cn.Connection{}, errors.New(err)
	} else if result.IsNil() {
		return []cn.Connection{}, pkg.ErrNotFound
	}

	var connectionMap []map[string]interface{}
	err = result.All(&connectionMap)
	if err != nil {
		return []cn.Connection{}, errors.New(err)
	}

	for _, data := range connectionMap {
		var c cn.DefaultConnection
		err = mapstructure.Decode(data, &c)
		if err != nil {
			return []cn.Connection{}, errors.New(err)
		}
		cs = append(cs, &c)
	}

	return cs, nil
}
