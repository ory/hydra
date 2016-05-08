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

1. Hydra is an OAuth2 and OpenID Connect provider with a few extras. The distributed in-memory design allows for heavy duty throughput.
2. Hydra works with every Identity Provider, even with that deprecated php-3.0 authentication service your intern wrote.
3. Hydra does not use any templates, it is up to you what your frontend should look like.

## Motivation

At first, there was the monolith. The monolith worked well with the customized joomla authentication module. Then, the web evolved into an elastic cloud that serves thousands of different user agents in every part of the world. Hydra is driven by the need for an easy scalable, in memory OAuth2 and OpenID Connect provider, that integrates with every Identity Provider you can imagine. 

Hydra uses pub/sub to always have the latest data available in memory. Hydra scales effortlessly on every platform you can imagine, including Heroku, Cloud Foundry, Docker, Google Container Engine and many more.

## Installation

**Host**

Docker
```
docker run -d oryam/hydra
```

With backend


**CLI Client**

https://gobuild.io/ory-am/hydra

Install go

Set up go path

Set up go binarypath
```
go install github.com/ory-am/hydra
hydra connect
```

## Documentation

Git Book link

### REST API

The REST API is documentet at [docs.hdyra.apiary.io](http://docs.hdyra.apiary.io).

### CLI

The CLI help is well documented. To see it, run `hydra -h` or `hydra [command] -h`.

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
