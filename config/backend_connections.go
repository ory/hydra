package config

import (
	"net/url"

	"time"

	"github.com/Sirupsen/logrus"
	r "github.com/dancannon/gorethink"
	"github.com/go-errors/errors"
	"github.com/ory-am/hydra/pkg"
)

type MemoryConnection struct{}

type RethinkDBConnection struct {
	session *r.Session
	URL     *url.URL
}

func (c *RethinkDBConnection) GetSession() *r.Session {
	if c.session != nil {
		return c.session
	}

	var err error
	var username, password string
	database := c.URL.Path[1:]
	if c.URL.User != nil {
		password, _ = c.URL.User.Password()
		username = c.URL.User.Username()
	}

	if err := pkg.Retry(time.Second*15, time.Minute*2, func() error {
		logrus.Infof("Connecting with RethinkDB: %s (%s) (%s)", c.URL.String(), c.URL.Host, database)
		if c.session, err = r.Connect(r.ConnectOpts{
			Address:  c.URL.Host,
			Username: username,
			Password: password,
		}); err != nil {
			return errors.Errorf("Could not connect to RethinkDB: %s", err)
		}

		if _, err := r.DBList().Contains(database).Do(func(e r.Term) r.Term {
			return r.Branch(
				e,
				map[string]interface{}{"dbs_created": 0},
				r.DBCreate(database),
			)
		}).RunWrite(c.session); err != nil {
			return errors.Errorf("Could not create database: %s", err)
		}

		c.session.Use(database)
		logrus.Infof("Connected to RethinkDB!")
		return nil
	}); err != nil {
		logrus.Fatalf("Could not connect to RethinkDB: %s", err)
	}

	return c.session
}

func (c *RethinkDBConnection) CreateTableIfNotExists(table string) {
	if _, err := r.TableList().Contains(table).Do(func(e r.Term) r.Term {
		return r.Branch(
			e,
			map[string]interface{}{"tables_created": 0},
			r.TableCreate(table),
		)
	}).RunWrite(c.GetSession()); err != nil {
		logrus.Fatalf("Could not create table: %s", err)
	}
}
