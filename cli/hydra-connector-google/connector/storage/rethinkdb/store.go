package rethinkdb

import (
	rdb "github.com/dancannon/gorethink"
	"github.com/ory-am/hydra/handler/connector/storage"

	"github.com/go-errors/errors"
	pkg "github.com/ory-am/common/pkg"
)

const storageStateTable = "hydra_state_data"

type Store struct {
	session *rdb.Session
}

func New(session *rdb.Session) *Store {
	return &Store{session}
}

func (s *Store) CreateTables() error {
	exists, err := s.tableExists(storageStateTable)
	if err == nil && !exists {
		_, err := rdb.TableCreate(storageStateTable).RunWrite(s.session)
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

func (s *Store) SaveStateData(sd *storage.StateData) error {
	res, err := rdb.Table(storageStateTable).Insert(sd).RunWrite(s.session)

	if err != nil {
		return errors.New(err.Error())
	} else if res.Errors > 0 {
		return errors.New(res.FirstError)
	}

	return nil
}

func (s *Store) GetStateData(id string) (*storage.StateData, error) {
	// Query state data
	result, err := rdb.Table(storageStateTable).Get(id).Run(s.session)

	if err != nil {
		return nil, errors.New(err)
	} else if result.IsNil() {
		return nil, pkg.ErrNotFound
	}

	defer result.Close()

	var sd storage.StateData
	err = result.One(&sd)
	if err != nil {
		return nil, errors.New(err)
	}

	return &sd, nil

}
