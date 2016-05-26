package config

import (
	r "github.com/dancannon/gorethink"
	"net/url"
	"github.com/Sirupsen/logrus"
)

type MemoryConnection struct{}

type RethinkDBConnection struct {
	Session *r.Session
	URL     *url.URL
}

func (c *RethinkDBConnection) GetSession() *r.Session {
	if c.Session != nil {
		return c.Session
	}

	var err error
	password, _ := c.URL.User.Password()
	if c.Session, err = r.Connect(r.ConnectOpts{
		Address: c.URL.Host,
		Database: c.URL.Path,
		Username: c.URL.User.Username(),
		Password: password,
	}); err != nil {
		logrus.Fatalf("Could not connect to RethinkDB: %s", err)
		return nil
	} else if _, err = r.DBCreate("hydra").RunWrite(c.Session); err != nil {
		logrus.Fatalf("Database could not be created: %s", err)
		return nil
	}
	return c.Session
}

func (c *RethinkDBConnection) CreateTableIfNotExists(table string) {
	if _, err := r.TableCreate(table).RunWrite(c.GetSession()); err != nil {
		logrus.Fatalf("Could not create table: %s", err)
		return
	}
}
