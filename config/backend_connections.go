package config

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/url"

	"time"

	"github.com/Sirupsen/logrus"
	"github.com/go-errors/errors"
	"github.com/ory-am/hydra/pkg"
	"github.com/spf13/viper"
	r "gopkg.in/dancannon/gorethink.v2"
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

		options := r.ConnectOpts{
			Address:  c.URL.Host,
			Username: username,
			Password: password,
		}

		if cert := loadCertificate(); cert != nil {
			roots := x509.NewCertPool()
			roots.AppendCertsFromPEM(cert)
			options.TLSConfig = &tls.Config{
				RootCAs: roots,
			}
		}

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

// loadCertificate reads a certificate either from a file or from a string
func loadCertificate() []byte {
	var cert []byte
	certPath := viper.GetString("RETHINK_TLS_CERT_PATH")
	if certPath != "" {
		var err error
		cert, err = ioutil.ReadFile(certPath)
		if err != nil {
			logrus.Debugf("Could not read rethinkdb certificate: %s", err)
			return nil
		}
	}

	certString := viper.GetString("RETHINK_TLS_CERT")
	if certString != "" {
		cert = []byte(certString)
	}

	return cert

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
