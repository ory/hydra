package handler

import (
	rethinkOsins "github.com/leetal/osin-rethinkdb/storage/rethinkdb"
	"github.com/ory-am/hydra/account"
	accounts "github.com/ory-am/hydra/account/rethinkdb"
	"github.com/ory-am/hydra/endpoint/connection"
	connections "github.com/ory-am/hydra/endpoint/connection/rethinkdb"
	statesStorage "github.com/ory-am/hydra/endpoint/connector/storage"
	states "github.com/ory-am/hydra/endpoint/connector/storage/rethinkdb"
	"github.com/ory-am/ladon/policy"
	policies "github.com/ory-am/ladon/policy/rethinkdb"
	"github.com/ory-am/osin-storage/storage"

	rdb "github.com/dancannon/gorethink"

	"strconv"

	"github.com/ory-am/hydra/hash"
)

// RethinkContext - Service backend context of RethinkDB
type RethinkContext struct {
	Context
	DB          *rdb.Session
	Accounts    *accounts.Store
	Connections *connections.Store
	Policies    *policies.Store
	Osins       *rethinkOsins.Storage
	States      *states.Store
}

const databaseName = "hydra"

// Init - Initializes the backend context
func (c *RethinkContext) Init() (Context, error) {
	return c, nil
}

// Start - Starts the backend context
func (c *RethinkContext) Start() error {
	// Load the environment variables
	getEnv()

	session, err := rdb.Connect(rdb.ConnectOpts{
		Address:  databaseURL,
		Database: databaseName,
	})

	c.DB = session

	if err != nil {
		return err
	}

	// Make sure that the database actually exists! Otherwise, create it!
	exists, err := c.databaseExists(databaseName)
	if err == nil && !exists {
		rdb.DBCreate(databaseName).RunWrite(session)
	} else if err != nil {
		return err
	}

	wf, err := strconv.Atoi(bcryptWorkFactor)
	if err != nil {
		return err
	}

	c.Accounts = accounts.New(&hash.BCrypt{WorkFactor: wf}, session)
	c.Connections = connections.New(session)
	c.Policies = policies.New(session)
	c.Osins = rethinkOsins.New(session)
	c.States = states.New(session)

	if err := c.Accounts.CreateTables(); err != nil {
		return err
	} else if err := c.Connections.CreateTables(); err != nil {
		return err
	} else if err := c.Policies.CreateTables(); err != nil {
		return err
	} else if err := c.Osins.CreateTables(); err != nil {
		return err
	} else if err := c.States.CreateTables(); err != nil {
		return err
	}

	return nil
}

func (c *RethinkContext) Close() {

}

func (c *RethinkContext) GetAccounts() account.Storage {
	return c.Accounts
}

func (c *RethinkContext) GetConnections() connection.Storage {
	return c.Connections
}

func (c *RethinkContext) GetPolicies() policy.Storage {
	return c.Policies
}

func (c *RethinkContext) GetOsins() storage.Storage {
	return c.Osins
}

func (c *RethinkContext) GetStates() statesStorage.Storage {
	return c.States
}

// DatabaseExists : Check if database exists
func (c *RethinkContext) databaseExists(name string) (bool, error) {
	res, err := rdb.DBList().Run(c.DB)
	if err != nil {
		return false, err
	}
	defer res.Close()

	if res.IsNil() {
		return false, nil
	}

	var database string
	for res.Next(&database) {
		if database == name {
			return true, nil
		}
	}

	return false, nil
}
