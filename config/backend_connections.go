package config

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/url"

	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory-am/hydra/pkg"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	r "gopkg.in/dancannon/gorethink.v2"
	"gopkg.in/redis.v5"
	"strings"
	"runtime"
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



	maxConns := maxParallelism() * 2
	if v := c.URL.Query().Get("max_conns"); v != "" {
		s, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			logrus.Warnf("max_conns value %s could not be parsed to int: %s", v, err)
		} else {
			maxConns = int(s)
		}
	}

	maxIdleConns := maxParallelism()
	if v := c.URL.Query().Get("max_idle_conns"); v != "" {
		s, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			logrus.Warnf("max_idle_conns value %s could not be parsed to int: %s", v, err)
		} else {
			maxIdleConns = int(s)
		}
	}

	maxConnLifetime := time.Duration(0)
	if v := c.URL.Query().Get("max_conn_lifetime"); v != "" {
		s, err := time.ParseDuration(v)
		if err != nil {
			logrus.Warnf("max_conn_lifetime value %s could not be parsed to int: %s", v, err)
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
	if len(c.URL.Path) <= 1 {
		logrus.Fatalf("Database hostname specified, but database name is missing.")
	}

	database := c.URL.Path[1:]
	if c.URL.User != nil {
		password, _ = c.URL.User.Password()
		username = c.URL.User.Username()
	}

	if err := pkg.Retry(time.Second*15, time.Minute*2, func() error {
		logrus.Infof("Connecting with RethinkDB: %s@%s/%s", username, c.URL.Host, database)

		options := r.ConnectOpts{
			Address:         c.URL.Host,
			Username:        username,
			Password:        password,
			KeepAlivePeriod: 10 * time.Second,
		}

		importRethinkDBRootCA(&options)

		if c.session, err = r.Connect(options); err != nil {
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

// importRethinkDBRootCA checks for the configuration values RETHINK_TLS_CERT_PATH
// or RETHINK_TLS_CERT and adds the certificate to the connect options
func importRethinkDBRootCA(opts *r.ConnectOpts) {
	var cert []byte
	certPath := viper.GetString("RETHINK_TLS_CERT_PATH")
	if certPath != "" {
		var err error
		cert, err = ioutil.ReadFile(certPath)
		if err != nil {
			logrus.Warningf("Could not read rethinkdb certificate: %s", err)
			return
		}
	}

	certString := viper.GetString("RETHINK_TLS_CERT")
	if certString != "" {
		cert = []byte(certString)
	}

	if cert != nil {
		roots := x509.NewCertPool()
		roots.AppendCertsFromPEM(cert)
		opts.TLSConfig = &tls.Config{
			RootCAs: roots,
		}
		logrus.Warnln("Loaded self-signed certificate for rethinkdb")
	}
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

type RedisConnection struct {
	db  *redis.Client
	URL *url.URL
}

func (c *RedisConnection) RedisSession() *redis.Client {
	if c.db != nil {
		return c.db
	}

	var password string
	var database int
	var err error

	if len(c.URL.Path) <= 1 {
		logrus.Infof("Defaulting database to 0.")
		database = 0
	} else {
		database, err = strconv.Atoi(c.URL.Path[1:])
		if err != nil {
			logrus.Fatalf("Database must be an integer.")
		}
	}

	if c.URL.User != nil {
		if p, exists := c.URL.User.Password(); exists {
			password = p
		} else {
			// No username, so first value is taken as password
			password = c.URL.User.Username()
		}
	}

	options := &redis.Options{
		Addr:     c.URL.Host,
		Password: password,
		DB:       database,
	}

	c.db = redis.NewClient(options)
	return c.db
}
