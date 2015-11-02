package handler

import (
	"database/sql"
	accounts "github.com/ory-am/hydra/account/postgres"
	connections "github.com/ory-am/hydra/oauth/connection/postgres"
	policies "github.com/ory-am/ladon/policy/postgres"
	osins "github.com/ory-am/osin-storage/storage/postgres"

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

	ctx := &Context{
		DB:          db,
		Accounts:    accounts.New(&hash.BCrypt{wf}, db),
		Connections: connections.New(db),
		Policies:    policies.New(db),
		Osins:       osins.New(db),
	}

	if err := ctx.Accounts.CreateSchemas(); err != nil {
		log.Fatal(err)
	} else if err := ctx.Connections.CreateSchemas(); err != nil {
		log.Fatal(err)
	} else if err := ctx.Policies.CreateSchemas(); err != nil {
		log.Fatal(err)
	} else if err := ctx.Osins.CreateSchemas(); err != nil {
		log.Fatal(err)
	}
}
