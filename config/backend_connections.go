package config

import (
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory/hydra/pkg"
	"github.com/pkg/errors"
)

type MemoryConnection struct{}

type SQLConnection struct {
	db  *sqlx.DB
	URL *url.URL
	L   logrus.FieldLogger
}

func cleanURLQuery(c *url.URL) *url.URL {
	cleanurl := new(url.URL)
	*cleanurl = *c

	q := cleanurl.Query()
	q.Del("max_conns")
	q.Del("max_idle_conns")
	q.Del("max_conn_lifetime")

	cleanurl.RawQuery = q.Encode()
	return cleanurl
}

func (c *SQLConnection) GetDatabase() *sqlx.DB {
	if c.db != nil {
		return c.db
	}

	var err error
	clean := cleanURLQuery(c.URL)

	if err = pkg.Retry(c.L, time.Second*15, time.Minute*2, func() error {
		c.L.Infof("Connecting with %s", c.URL.Scheme+"://*:*@"+c.URL.Host+c.URL.Path+"?"+clean.RawQuery)
		u := clean.String()
		if clean.Scheme == "mysql" {
			u = strings.Replace(u, "mysql://", "", -1)
		}

		if c.db, err = sqlx.Open(clean.Scheme, u); err != nil {
			return errors.Errorf("Could not connect to SQL: %s", err)
		} else if err := c.db.Ping(); err != nil {
			return errors.Errorf("Could not connect to SQL: %s", err)
		}

		c.L.Infof("Connected to SQL!")
		return nil
	}); err != nil {
		c.L.Fatalf("Could not connect to SQL: %s", err)
	}

	maxConns := maxParallelism() * 2
	if v := c.URL.Query().Get("max_conns"); v != "" {
		s, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			c.L.Warnf("max_conns value %s could not be parsed to int: %s", v, err)
		} else {
			maxConns = int(s)
		}
	}

	maxIdleConns := maxParallelism()
	if v := c.URL.Query().Get("max_idle_conns"); v != "" {
		s, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			c.L.Warnf("max_idle_conns value %s could not be parsed to int: %s", v, err)
		} else {
			maxIdleConns = int(s)
		}
	}

	maxConnLifetime := time.Duration(0)
	if v := c.URL.Query().Get("max_conn_lifetime"); v != "" {
		s, err := time.ParseDuration(v)
		if err != nil {
			c.L.Warnf("max_conn_lifetime value %s could not be parsed to int: %s", v, err)
		} else {
			maxConnLifetime = s
		}
	}

	c.db.SetMaxOpenConns(maxConns)
	c.db.SetMaxIdleConns(maxIdleConns)
	c.db.SetConnMaxLifetime(maxConnLifetime)

	return c.db
}

func maxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}
