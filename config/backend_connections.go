package config

import (
	"net/url"
	"time"
	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory-am/hydra/pkg"
	"github.com/pkg/errors"
	"strings"
)

type MemoryConnection struct{}

type SQLConnection struct {
	db  *sqlx.DB
	URL *url.URL
}

func (c *SQLConnection) GetDatabase() *sqlx.DB {
	if c.db != nil {
		return c.db
	}

	var err error
	if err = pkg.Retry(time.Second*15, time.Minute*2, func() error {
		logrus.Infof("Connecting with %s", c.URL.Scheme+"://*:*@"+c.URL.Host+c.URL.Path+"?"+c.URL.RawQuery)
		u := c.URL.String()
		if c.URL.Scheme == "mysql" {
			u = strings.Replace(u, "mysql://", "", -1)
		}

		if c.db, err = sqlx.Open(c.URL.Scheme, u); err != nil {
			return errors.Errorf("Could not connect to SQL: %s", err)
		} else if err := c.db.Ping(); err != nil {
			return errors.Errorf("Could not connect to SQL: %s", err)
		}

		logrus.Infof("Connected to SQL!")
		return nil
	}); err != nil {
		logrus.Fatalf("Could not connect to SQL: %s", err)
	}

	return c.db
}
