# Dockertest

[![Build Status](https://travis-ci.org/ory-am/dockertest.svg)](https://travis-ci.org/ory-am/dockertest) [![Coverage Status](https://coveralls.io/repos/ory-am/dockertest/badge.svg?branch=master&service=github)](https://coveralls.io/github/ory-am/dockertest?branch=master)

Use Docker to run your Go language integration tests against persistent data storage services like **MySQL, Postgres or MongoDB** on **Microsoft Windows, Mac OSX and Linux**! Dockertest uses [docker-machine](https://docs.docker.com/machine/) (aka [Docker Toolbox](https://www.docker.com/toolbox)) to spin up images on Windows and Mac OSX as well.

A suite for testing with Docker. Based on  [docker.go](https://github.com/camlistore/camlistore/blob/master/pkg/test/dockertest/docker.go) from [camlistore](https://github.com/camlistore/camlistore).
This fork detects automatically, if [Docker Toolbox](https://www.docker.com/toolbox) is installed. If it is, Docker integration on Windows and Mac OSX can be used without any additional work. To avoid port collisions when using docker-machine, Dockertest chooses a random port to bind the requested image.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [Why should I use Dockertest?](#why-should-i-use-dockertest)
- [Using Dockertest](#using-dockertest)
  - [Start a container](#start-a-container)
- [Write awesome tests](#write-awesome-tests)
  - [Setting up Travis-CI](#setting-up-travis-ci)
- [Troubleshoot](#troubleshoot)
  - [Out of disk space](#out-of-disk-space)
  - [Removing old containers](#removing-old-containers)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Why should I use Dockertest?

When developing applications, it is often necessary to use services that talk to a database system. Unit Testing these services can be cumbersome because mocking database/DBAL is strenuous. Making slight changes to the schema implies rewriting at least some, if not all of the mocks. The same goes for API changes in the DBAL.  
To avoid this, it is smarter to test these specific services against a real database that is destroyed after testing. Docker is the perfect system for running unit tests as you can spin up containers in a few seconds and kill them when the test completes. The Dockertest library provides easy to use commands for spinning up Docker containers and using them for your tests.

## Using Dockertest

Using Dockertest is straightforward and  simple. At present, Dockertest supports MongoDB, Postgres and MySQL containers out of the box. Feel free to extend this list by contributing to this project.

**Note:** When using the Docker Toolbox (Windows / OSX), make sure that the VM is started by running `docker-machine start default`.

### Start a container

```go
package main

import "github.com/ory-am/dockertest"
import "gopkg.in/mgo.v2"
import "time"

func main() {
	c, err := ConnectToMongoDB(15, time.Millisecond*500, func(url string) bool {
		db, err := mgo.Dial(url)
		if err != nil {
			return false
		}
		defer db.Close()
		return true
	})
	require.Nil(t, err)
	defer c.KillRemove()
}
```

You can start PostgreSQL and MySQL in a similar fashion with


## Write awesome tests

It is a good idea to start up the container only once when running tests.

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
	if c, err := dockertest.ConnectToPostgreSQL(15, time.Second, func(url string) bool {
		var err error
		db, err = sql.Open("postgres", url)
		if err != nil {
			return false
		}
		return db.Ping() == nil
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	defer c.KillRemove()
	os.Exit(m.Run())
}

func TestFunction(t *testing.T) {
    // db.Exec(...
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
  - DOCKERTEST_BIND_LOCALHOST=true

```

## Troubleshoot

### Out of disk space

Try cleaning up the images with [docker-cleanup-volumes](https://github.com/chadoe/docker-cleanup-volumes).

### Removing old containers

Sometimes container clean up fails. Check out
[this stackoverflow question](http://stackoverflow.com/questions/21398087/how-to-delete-dockers-images) on how to fix this.

*Thanks to our sponsors: Ory GmbH & Imarum GmbH*
