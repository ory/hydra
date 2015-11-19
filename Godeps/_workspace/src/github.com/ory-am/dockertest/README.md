# Dockertest

[![Build Status](https://travis-ci.org/ory-am/dockertest.svg)](https://travis-ci.org/ory-am/dockertest)

Use Docker to run your Go language integration tests against persistent data storage services like **MySQL, Postgres or MongoDB** on **Microsoft Windows, Mac OSX and Linux**! Dockertest uses [docker-machine](https://docs.docker.com/machine/) (aka [Docker Toolbox](https://www.docker.com/toolbox)) to spin up images on Windows and Mac OSX as well.

A suite for testing with Docker. Based on  [docker.go](https://github.com/camlistore/camlistore/blob/master/pkg/test/dockertest/docker.go) from [camlistore](https://github.com/camlistore/camlistore).
This fork detects automatically, if [Docker Toolbox](https://www.docker.com/toolbox) is installed. If it is, Docker integration on Windows and Mac OSX can be used without any additional work. To avoid port collisions when using docker-machine, Dockertest chooses a random port to bind the requested image.

## Why should I use Dockertest?

When developing applications, it is often necessary to use services that talk to a database system. Unit Testing these services can be cumbersome because mocking database/DBAL is strenuous. Making slight changes to the schema implies rewriting at least some, if not all of the mocks. The same goes for API changes in the DBAL.  
To avoid this, it is smarter to test these specific services against a real database that is destroyed after testing. Docker is the perfect system for running unit tests as you can spin up containers in a few seconds and kill them when the test completes. The Dockertest library provides easy to use commands for spinning up Docker containers and using them for your tests.

## Using Dockertest

Using Dockertest is straightforward and  simple. At present, Dockertest supports MongoDB, Postgres and MySQL containers out of the box. Feel free to extend this list by contributing to this project.

**Note:** When using the Docker Toolbox (Windows / OSX), make sure that the VM is started by running `docker-machine start default`.

### MongoDB Container

```go
import "github.com/ory-am/dockertest"
import "gopkg.in/mgo.v2"
import "time"

func Foobar() {
  // Start MongoDB Docker container. Wait 1 second for the image to load.
  containerID, ip, port, err := dockertest.SetupMongoContainer()

  if err != nil {
    return err
  }

  // kill the container on deference
  defer containerID.KillRemove()

  url := fmt.Sprintf("%s:%d", ip, port)
  sess, err := mgo.Dial(url)
  if err != nil {
    return err
  }

  defer sess.Close()
  // ...
}
```

### MySQL Container

```go
import "github.com/ory-am/dockertest"
import "github.com/go-sql-driver/mysql"
import "database/sql"
import "time"

func Foobar() {
    // Wait 10 seconds for the image to load.
    c, ip, port, err := dockertest.SetupMySQLContainer()
    if err != nil {
        return
    }
    defer c.KillRemove()

    url := fmt.Sprintf("mysql://%s:%s@%s:%d/", dockertest.MySQLUsername, dockertest.MySQLPassword, ip, port)
    db, err := sql.Open("mysql", url)
    if err != nil {
        return
    }

    defer db.Close()
    // ...
}
```
### Postgres Container

```go
import "github.com/ory-am/dockertest"
import "github.com/lib/pq"
import "database/sql"
import "time"

func Foobar() {
    // Wait 10 seconds for the image to load.
    c, ip, port, err := dockertest.SetupPostgreSQLContainer()
    if err != nil {
        return
    }
    defer c.KillRemove()

    url := fmt.Sprintf("postgres://%s:%s@%s:%d/", dockertest.PostgresUsername, dockertest.PostgresPassword, ip, port)
    db, err := sql.Open("postgres", url)
    if err != nil {
        return
    }

    defer db.Close()
    // ...
}
```

## Usage in tests

It is a good idea to start up the container only once when running tests. To achieve this for example:

```go

import (
	"fmt"
	"testing"
    "log"
	"os"

	"database/sql"
	_ "github.com/lib/pq"
	"github.com/ory-am/dockertest"
)

var db *sql.DB

func TestMain(m *testing.M) {
	c, ip, port, err := dockertest.SetupPostgreSQLContainer()
	if err != nil {
		log.Fatalf("Could not set up PostgreSQL container: %v", err)
	}
	defer c.KillRemove()

	url := fmt.Sprintf("postgres://%s:%s@%s:%d/postgres?sslmode=disable", dockertest.PostgresUsername, dockertest.PostgresPassword, ip, port)
	db, err = sql.Open("postgres", url)
	if err != nil {
		log.Fatalf("Could not set up PostgreSQL container: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Could not ping database: %v", err)
	}

	os.Exit(m.Run())
}

func TestFunction(t *testing.T) {
    // ...
}
```

### Setting up Travis-CI

You can run the Docker integration on Travis easily:

```yml
# Sudo is required for docker
sudo: required

# Enable docker
services:
  - docker

# In Travis, we need to bind to 127.0.0.1 in order to get a working connection. This environment variable
# tells dockertest to do that.
env:
  - DOCKER_BIND_LOCALHOST=true

```

Thanks to our sponsors: Ory GmbH & Imarum GmbH
