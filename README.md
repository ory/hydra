# Ory/Hydra

[![Join the chat at https://gitter.im/ory-am/hydra](https://badges.gitter.im/ory-am/hydra.svg)](https://gitter.im/ory-am/hydra?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Build Status](https://travis-ci.org/ory-am/hydra.svg?branch=master)](https://travis-ci.org/ory-am/hydra)
[![Coverage Status](https://coveralls.io/repos/ory-am/hydra/badge.svg?branch=master&service=github)](https://coveralls.io/github/ory-am/hydra?branch=master)

![Hydra](hydra.png)

Hydra is being developed at [Ory](https://ory.am). Join our [mailinglist](http://eepurl.com/bKT3N9) to stay on top of new developments.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [What is Hydra?](#what-is-hydra)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## What is Hydra?

1. Hydra is an OAuth2 and OpenID Connect provider built for availability. The distributed in-memory architecture allows for heavy duty workloads.
2. Hydra works with every Identity Provider. The deprecated php-3.0 authentication service your intern wrote? It works with that too, don't worry.
3. Hydra does not use any templates, it is up to you what your front end should look like.
4. Hydra comes with two factor authentication, key management, social log on, policy management and access control.

## Motivation

At first, there was the monolith. The monolith worked well with the customized joomla authentication module. Then, the web evolved into an elastic cloud that serves thousands of different user agents in every part of the world. Hydra is driven by the need for an easy scalable, in memory OAuth2 and OpenID Connect provider, that integrates with every Identity Provider you can imagine. 

Hydra uses pub/sub to always have the latest data available in memory. Hydra scales effortlessly on every platform you can imagine, including Heroku, Cloud Foundry, Docker, Google Container Engine and many more.

## First steps

This section is a quickstart guide to working with Hydra. In-depth docs are available as well:

* The documentation is available on [GitBook](https://ory-am.gitbooks.io/hydra/content/).
* The REST API documentation is available at [Apiary](http://docs.hdyra.apiary.io).

### Installation

**Starting the host** is easiest with docker. The host process handles HTTP requests and is backed by a database.
Read how to install docker on [Linux](https://docs.docker.com/linux/), [OSX](https://docs.docker.com/mac/) or
[Windows](https://docs.docker.com/windows/).

The easiest way to start docker is without a database. Hydra will keep all changes in memory. But be aware! Restarting, scaling
or stopping the container will make you **lose all data**.

```
docker run oryam/hydra -p 4444:4444
open http://$(docker-machine ip default):444
```

**CLI Client**

The CLI client is available at [gobuild.io](https://gobuild.io/ory-am/hydra) or through
the [releases tab](https://github.com/ory-am/hydra/releases).

If you wish to compile the CLI yourself, you need to install and set up [Go](https://golang.org/) and add `$GOPATH/bin`
to your `$PATH`.

```
go install github.com/ory-am/hydra
hydra
```

### First steps



## Documentation

The documentation is available on [GitBook](https://ory-am.gitbooks.io/hydra/content/).

### REST API Documentation

The REST API is documented at [Apiary](http://docs.hdyra.apiary.io).

### CLI Documentation

The CLI help is verbose. To see it, run `hydra -h` or `hydra [command] -h`.

### Develop

Unless you want to test Hydra against a database, developing with Hydra is as easy as

```
go get github.com/ory-am/hydra
cd $GOPATH/src/github.com/ory-am/hydra
git checkout -b develop
go test ./... -race
go run main.go
```

## Frequently Asked Questions

### Deploy using buildpacks (Heroku, Cloud Foundry, ...)

Hydra runs pretty much out of the box when using a Platform as a Service (PaaS).
Here are however a few notes which might assist you in your task:
* Heroku (and probably Cloud Foundry as well) *force* TLS termination, meaning that Hydra must be configured with `DANGEROUSLY_FORCE_HTTP=force`.
* Using bash, you can easily add multi-line environment variables to Heroku using `heroku config:set JWT_PUBLIC_KEY="$(my-public-key.pem)"`.
  This does not work on Windows!

## Attributions

* [Logo source](https://www.flickr.com/photos/pathfinderlinden/7161293044/)
