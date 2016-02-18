package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"

	"database/sql"
	"log"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
	"github.com/ory-am/common/env"
	"gopkg.in/ory-am/dockertest.v2"
)

var tmpDir = os.TempDir()

var db *sql.DB

func TestMain(m *testing.M) {
	c, err := dockertest.ConnectToPostgreSQL(15, time.Second, func(url string) bool {
		var err error
		db, err = sql.Open("postgres", url)
		if err != nil {
			log.Printf("Could not connect to database because %s", err)
			return false
		} else if err := db.Ping(); err != nil {
			log.Printf("Could not ping database because %s", err)
			return false
		}

		// Database is now available, let's continue!
		os.Setenv("DATABASE_URL", url)
		if env.Getenv("DATABASE_URL", "") != url {
			log.Fatalf("Could not set DATABASE_URL environment variable: %s != %s", env.Getenv("DATABASE_URL", ""), url)
			return false
		}
		return true
	})

	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	retCode := m.Run()

	// force teardown
	tearDown(c)

	os.Exit(retCode)
}

func tearDown(c dockertest.ContainerID) {
	db.Close()
	c.KillRemove()
}

func TestRunCLITests(t *testing.T) {
	oldStderr := os.Stderr // keep backup of the real stdout

	for k, c := range []struct {
		args []string
		pass bool
	}{
		{args: []string{"hydra-host", "account", "create", "foo@bar", "--password", "secret"}, pass: true},
		{args: []string{"hydra-host", "account", "create", "foo@bar", "--password", "secret"}, pass: false},
		{args: []string{"hydra-host", "account", "create", "bar@baz", "--password", "secret", "--as-superuser"}, pass: true},
		{args: []string{"hydra-host", "client", "create", "-i", "foo-app", "-s", "secret", "-r", "http://localhost/"}, pass: true},
		{args: []string{"hydra-host", "client", "create", "-i", "foo-app", "-s", "secret", "-r", "http://localhost/"}, pass: false},
		{args: []string{"hydra-host", "client", "create", "-i", "bar-app", "-s", "secret", "-r", "http://localhost/", "--as-superuser"}, pass: true},
		{args: []string{"hydra-host", "policy", "import", "../../example/policies.json"}, pass: true},
	} {
		r, w, _ := os.Pipe()
		os.Stderr = w
		outC := make(chan string)
		go func() {
			var buf bytes.Buffer
			io.Copy(&buf, r)
			outC <- buf.String()
		}()

		os.Args = c.args
		err := NewApp().Run(os.Args)
		assert.Nil(t, err, "Case %d: %s", k, err)

		w.Close()
		if c.pass {
			assert.Empty(t, <-outC)
		} else {
			assert.NotEmpty(t, <-outC)
		}
	}
	os.Stdout = oldStderr
}

func TestJWTGen(t *testing.T) {
	priv := filepath.Join(tmpDir, uuid.New())
	pub := filepath.Join(tmpDir, uuid.New())
	os.Args = []string{"hydra-host", "jwt", "generate-keypair", "-s", priv, "-p", pub}
	assert.Nil(t, NewApp().Run(os.Args))
	assertAndRemoveFile(t, priv)
	assertAndRemoveFile(t, pub)
}

func TestTLSGen(t *testing.T) {
	priv := filepath.Join(tmpDir, uuid.New())
	pub := filepath.Join(tmpDir, uuid.New())
	os.Args = []string{"hydra-host", "tls", "generate-dummy-certificate", "-c", priv, "-k", pub, "-u", "localhost", "--sd", "Jan 1 15:04:05 2011", "-d", "8760h0m0s", "--ca", "--rb", "4069", "--ec", "P521"}
	assert.Nil(t, NewApp().Run(os.Args))
	assertAndRemoveFile(t, priv)
	assertAndRemoveFile(t, pub)
}

func assertAndRemoveFile(t *testing.T, file string) {
	_, err := os.Stat(file)
	assert.Nil(t, err)
	assert.Nil(t, os.Remove(file))
}
