package handler

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/ory-am/hydra/account"
	"github.com/ory-am/hydra/handler/connection"
	statesStorage "github.com/ory-am/hydra/handler/connector/storage"
	"github.com/ory-am/ladon/policy"
	"github.com/ory-am/osin-storage/storage"
)

// DefaultSystemContext - This is the container for the context that we use
// throughout Hydra
type DefaultSystemContext struct {
	Ctx Context
	sync.Mutex
}

// GetSystemContext - Fetches the system context by lazy-loading it
func (d *DefaultSystemContext) GetSystemContext() Context {

	// Lazy-load the system context (only done once!)
	// Lock just in case we rewrite the instantiation in the future..
	d.Lock()
	if d.Ctx == nil {
		// Load the environment variables
		getEnv()

		// This is done only once!
		var ctx Context
		var err error

		if strings.Contains(databaseURL, "rethinkdb://") {
			// Strip unwanted "rethinkdb://"
			databaseURL = strings.Replace(databaseURL, "rethinkdb://", "", 1)
			os.Setenv("DATABASE_URL", databaseURL)
			// Init the context
			ctx, err = new(RethinkContext).Init()
		} else {
			// Fall back to postgres
			// Init the context
			ctx, err = new(PostgresContext).Init()
		}

		if err != nil {
			log.Fatal(err.Error())
		}
		d.Ctx = ctx
	}
	d.Unlock()

	return d.Ctx
}

// Context - Context is an interface we use in Hydra to utilize different storage
// backends.
type Context interface {

	// Initializes the context and sets default values (if needed)
	Init() (Context, error)

	// Start - Starts the context (usually creats a connection pool to the DB and so forth)
	Start() error

	// Close - Closes connection to DB and other things that might have to be done on shutdown
	Close()

	GetAccounts() account.Storage

	GetConnections() connection.Storage

	GetPolicies() policy.Storage

	GetOsins() storage.Storage

	GetStates() statesStorage.Storage
}
