package config

import (
	"net/url"

	"github.com/Sirupsen/logrus"
	r "github.com/dancannon/gorethink"
	"time"
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
	var username, password string
	database := c.URL.Path[1:]
	if c.URL.User != nil {
		password, _ = c.URL.User.Password()
		username = c.URL.User.Username()
	}
	for i := 0; i < 10; i++ {
		logrus.Infof("Connecting with RethinkDB: %s (%s) (%s)", c.URL.String(), c.URL.Host, database)
		if c.Session, err = r.Connect(r.ConnectOpts{
			Address:  c.URL.Host,
			Username: username,
			Password: password,
		}); err != nil {
			logrus.Warnf("Could not connect to RethinkDB: %s", err)
			logrus.Warnf("Retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}

		if _, err := r.DBList().Contains(database).Do(func(e r.Term) r.Term {
			return r.Branch(
				e,
				map[string]interface{}{"dbs_created": 0},
				r.DBCreate(database),
			)
		}).RunWrite(c.Session); err != nil {
			logrus.Fatalf("Could not create database: %s", err)
			logrus.Warnf("Retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		} else {
			c.Session.Use(database)
			logrus.Infof("Connected to RethinkDB!")
			return c.Session
		}

	}

	logrus.Fatalf("Could not connect to RethinkDB: %s", err)
	return nil
}

func (c *RethinkDBConnection) CreateTableIfNotExists(table string) {
	if _, err := r.TableList().Contains(table).Do(func(e r.Term) r.Term {
		return r.Branch(
			e,
			map[string]interface{}{"tables_created": 0},
			r.TableCreate(table),
		)
	}).RunWrite(c.Session); err != nil {
		logrus.Fatalf("Could not create table: %s", err)
	}
}
