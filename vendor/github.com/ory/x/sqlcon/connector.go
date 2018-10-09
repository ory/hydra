/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package sqlcon

import (
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type SQLConnection struct {
	db  *sqlx.DB
	URL *url.URL
	L   logrus.FieldLogger
}

func NewSQLConnection(db string, l logrus.FieldLogger) (*SQLConnection, error) {
	u, err := url.Parse(db)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if l == nil {
		logger := logrus.New()

		// Basically avoids any logging because no one uses panics
		logger.Level = logrus.PanicLevel

		l = logger
	}

	return &SQLConnection{
		URL: u,
		L:   l,
	}, nil
}

func cleanURLQuery(c *url.URL) *url.URL {
	cleanurl := new(url.URL)
	*cleanurl = *c

	q := cleanurl.Query()
	q.Del("max_conns")
	q.Del("max_idle_conns")
	q.Del("max_conn_lifetime")
	q.Del("parseTime")

	cleanurl.RawQuery = q.Encode()
	return cleanurl
}

func (c *SQLConnection) GetDatabaseRetry(maxWait time.Duration, failAfter time.Duration) (*sqlx.DB, error) {
	if err := retry(c.L, maxWait, failAfter, func() error {
		if err := c.GetDatabase().Ping(); err != nil {
			return errors.WithStack(err)
		}
		return nil
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	return c.db, nil
}

func (c *SQLConnection) GetDatabase() *sqlx.DB {
	if c.db != nil {
		return c.db
	}

	var err error

	clean := cleanURLQuery(c.URL)

	if err = retry(c.L, time.Second*15, time.Minute*2, func() error {
		c.L.Infof("Connecting with %s", c.URL.Scheme+"://*:*@"+c.URL.Host+c.URL.Path+"?"+clean.RawQuery)

		u := connectionString(clean)
		if c.db, err = sqlx.Open(clean.Scheme, u); err != nil {
			return errors.Errorf("Could not Connect to SQL: %s", err)
		} else if err := c.db.Ping(); err != nil {
			return errors.Errorf("Could not Connect to SQL: %s", err)
		}

		c.L.Infof("Connected to SQL!")
		return nil
	}); err != nil {
		c.L.Fatalf("Could not Connect to SQL: %s", err)
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

func connectionString(clean *url.URL) string {
	if clean.Scheme == "mysql" {
		q := clean.Query()
		q.Set("parseTime", "true")
		clean.RawQuery = q.Encode()
	}

	username := clean.User.Username()
	userinfo := username
	password, hasPassword := clean.User.Password()
	if hasPassword {
		userinfo = userinfo + ":" + password
	}
	clean.User = nil
	u := clean.String()
	clean.User = url.UserPassword(username, password)

	if strings.HasPrefix(u, clean.Scheme+"://") {
		u = strings.Replace(u, clean.Scheme+"://", clean.Scheme+"://"+userinfo+"@", 1)
	}
	if clean.Scheme == "mysql" {
		u = strings.Replace(u, "mysql://", "", -1)
	}
	return u
}
