package handler

import (
	"database/sql"

	"github.com/ory-am/hydra/account"
	accounts "github.com/ory-am/hydra/account/postgres"
	"github.com/ory-am/hydra/endpoint/connection"
	connections "github.com/ory-am/hydra/endpoint/connection/postgres"
	statesStorage "github.com/ory-am/hydra/endpoint/connector/storage"
	states "github.com/ory-am/hydra/endpoint/connector/storage/postgres"
	"github.com/ory-am/ladon/policy"
	policies "github.com/ory-am/ladon/policy/postgres"
	"github.com/ory-am/osin-storage/storage"
	osins "github.com/ory-am/osin-storage/storage/postgres"

	"strconv"

	_ "github.com/lib/pq"
	"github.com/ory-am/hydra/hash"
)

type PostgresContext struct {
	DB          *sql.DB
	Accounts    *accounts.Store
	Connections *connections.Store
	Policies    *policies.Store
	Osins       *osins.Storage
	States      *states.Store
}

// Init - Initializes the backend context
func (c *PostgresContext) Init() (Context, error) {
	return c, nil
}

// Start - Starts the backend context
func (c *PostgresContext) Start() error {
	// Load the environment variables
	getEnv()

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return err
	} else if err := db.Ping(); err != nil {
		return err
	}

	wf, err := strconv.Atoi(bcryptWorkFactor)
	if err != nil {
		return err
	}

	c.DB = db
	c.Accounts = accounts.New(&hash.BCrypt{wf}, db)
	c.Connections = connections.New(db)
	c.Policies = policies.New(db)
	c.Osins = osins.New(db)
	c.States = states.New(db)

	if err := c.Accounts.CreateSchemas(); err != nil {
		return err
	} else if err := c.Connections.CreateSchemas(); err != nil {
		return err
	} else if err := c.Policies.CreateSchemas(); err != nil {
		return err
	} else if err := c.Osins.CreateSchemas(); err != nil {
		return err
	} else if err := c.States.CreateSchemas(); err != nil {
		return err
	}

	return nil
}

func (c *PostgresContext) Close() {

}

func (c *PostgresContext) GetAccounts() account.Storage {
	return c.Accounts
}

func (c *PostgresContext) GetConnections() connection.Storage {
	return c.Connections
}

func (c *PostgresContext) GetPolicies() policy.Storage {
	return c.Policies
}

func (c *PostgresContext) GetOsins() storage.Storage {
	return c.Osins
}

func (c *PostgresContext) GetStates() statesStorage.Storage {
	return c.States
}
