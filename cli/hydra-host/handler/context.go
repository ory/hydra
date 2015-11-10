package handler

import (
	"database/sql"
	policies "github.com/ory-am/ladon/policy/postgres"
	osins "github.com/ory-am/osin-storage/storage/postgres"
	accounts "github.com/ory-am/hydra/account/postgres"
	connections "github.com/ory-am/hydra/oauth/connection/postgres"

	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
	"github.com/ory-am/hydra/hash"
	"strconv"
)

type Context struct {
	DB          *sql.DB
	Accounts    *accounts.Store
	Connections *connections.Store
	Policies    *policies.Store
	Osins       *osins.Storage
}

func (c *Context) Start() {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal(err)
	} else if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	wf, err := strconv.Atoi(bcryptWorkFactor)
	if err != nil {
		log.Fatal(err)
	}

	c.DB = db
	c.Accounts = accounts.New(&hash.BCrypt{wf}, db)
	c.Connections = connections.New(db)
	c.Policies = policies.New(db)
	c.Osins = osins.New(db)

	if err := c.Accounts.CreateSchemas(); err != nil {
		log.Fatal(err)
	} else if err := c.Connections.CreateSchemas(); err != nil {
		log.Fatal(err)
	} else if err := c.Policies.CreateSchemas(); err != nil {
		log.Fatal(err)
	} else if err := c.Osins.CreateSchemas(); err != nil {
		log.Fatal(err)
	}
}
