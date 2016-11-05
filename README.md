# ![Ory/Hydra](docs/images/logo.png)

[![Join the chat at https://gitter.im/ory-am/hydra](https://img.shields.io/badge/join-chat-00cc99.svg)](https://gitter.im/ory-am/hydra?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Join mailinglist](https://img.shields.io/badge/join-mailinglist-00cc99.svg)](https://groups.google.com/forum/#!forum/ory-hydra/new)
[![Join newsletter](https://img.shields.io/badge/join-newsletter-00cc99.svg)](http://eepurl.com/bKT3N9)
[![Follow twitter](https://img.shields.io/badge/follow-twitter-00cc99.svg)](https://twitter.com/_aeneasr)
[![Follow GitHub](https://img.shields.io/badge/follow-github-00cc99.svg)](https://github.com/arekkas)
[![Become a patron!](https://img.shields.io/badge/support%20us-on%20patreon-green.svg)](https://patreon.com/user?u=4298803)

[![Build Status](https://travis-ci.org/ory-am/hydra.svg?branch=master)](https://travis-ci.org/ory-am/hydra)
[![Coverage Status](https://coveralls.io/repos/ory-am/hydra/badge.svg?branch=master&service=github)](https://coveralls.io/github/ory-am/hydra?branch=master)
[![Code Climate](https://codeclimate.com/github/ory-am/hydra/badges/gpa.svg)](https://codeclimate.com/github/ory-am/hydra)
[![Go Report Card](https://goreportcard.com/badge/github.com/ory-am/hydra)](https://goreportcard.com/report/github.com/ory-am/hydra)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/364/badge)](https://bestpractices.coreinfrastructure.org/projects/364)

[![Docs Guide](https://img.shields.io/badge/docs-guide-blue.svg)](https://ory-am.gitbooks.io/hydra/content/)
[![HTTP API Documentation](https://img.shields.io/badge/docs-http%20api-blue.svg)](http://docs.hdyra.apiary.io/)
[![Code Documentation](https://img.shields.io/badge/docs-godoc-blue.svg)](https://godoc.org/github.com/ory-am/hydra)

Hydra is a runnable server implementation of the OAuth2 2.0 authorization framework and the OpenID Connect Core 1.0.

Join our [newsletter](http://eepurl.com/bKT3N9) to stay on top of new developments.
We answer basic support requests on [Google Groups](https://groups.google.com/forum/#!forum/ory-hydra/new) and [Gitter](https://gitter.im/ory-am/hydra)
and offer [premium services](http://www.ory.am/products/hydra) around Hydra.

Hydra uses the security first OAuth2 and OpenID Connect SDK [Fosite](https://github.com/ory-am/fosite) and
the access control SDK [Ladon](https://github.com/ory-am/ladon).

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [What is Hydra?](#what-is-hydra)
- [Quickstart](#quickstart)
  - [5 minutes tutorial: Run your very own OAuth2 environment](#5-minutes-tutorial-run-your-very-own-oauth2-environment)
  - [Installation](#installation)
    - [Download binaries](#download-binaries)
    - [Using Docker](#using-docker)
    - [Building from source](#building-from-source)
- [Security](#security)
- [Reception](#reception)
- [Documentation](#documentation)
  - [Guide](#guide)
  - [HTTP API Documentation](#http-api-documentation)
  - [Command Line Documentation](#command-line-documentation)
  - [Develop](#develop)
- [Third-party libraries and projects](#third-party-libraries-and-projects)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## What is Hydra?

Hydra is a server implementation of the OAuth 2.0 authorization framework and the OpenID Connect Core 1.0. Existing OAuth2
implementations usually ship as libraries or SDKs such as [node-oauth2-server](https://github.com/oauthjs/node-oauth2-server)
or [fosite](https://github.com/ory-am/fosite/issues), or as fully featured identity solutions with user
management and user interfaces, such as [Dex](https://github.com/coreos/dex).

Implementing and using OAuth2 without understanding the whole specification is challenging and prone to errors, even when
SDKs are being used. The primary goal of Hydra is to make OAuth 2.0 and OpenID Connect 1.0 better accessible.

Hydra implements the flows described in OAuth2 and OpenID Connect 1.0 without forcing you to use a "Hydra User Management"
or some template engine or a predefined front-end. Instead it relies on HTTP redirection and cryptographic methods
to verify user consent allowing you to use Hydra with any authentication endpoint, be it [authboss](https://github.com/go-authboss/authboss),
[auth0.com](https://auth0.com/) or your proprietary PHP authentication.

## Quickstart

This section is a quickstart guide to working with Hydra. In-depth docs are available as well:

* The documentation is available on [GitBook](https://ory-am.gitbooks.io/hydra/content/).
* The REST API documentation is available at [Apiary](http://docs.hdyra.apiary.io).

### 5 minutes tutorial: Run your very own OAuth2 environment

The **[tutorial](https://ory-am.gitbooks.io/hydra/content/tutorial.html)** teaches you to set up Hydra,
a Postgres instance and an exemplary identity provider written in React using docker compose.
It will take you about 5 minutes to complete the **[tutorial](https://ory-am.gitbooks.io/hydra/content/tutorial.html)**.

<img src="docs/images/oauth2-flow.gif" alt="OAuth2 Flow">

<br clear="all">

### Installation

There are various ways of installing hydra on your system.

#### Download binaries

The client and server **binaries are downloadable at [releases](https://github.com/ory-am/hydra/releases)**.
There is currently no installer available. You have to add the hydra binary to the PATH environment variable yourself or put
the binary in a location that is already in your path (`/usr/bin`, ...). 
If you do not understand what that all of this means, ask in our [chat channel](https://gitter.im/ory-am/hydra). We are happy to help.

#### Using Docker

**Starting the host** is easiest with docker. The host process handles HTTP requests and is backed by a database.
Read how to install docker on [Linux](https://docs.docker.com/linux/), [OSX](https://docs.docker.com/mac/) or
[Windows](https://docs.docker.com/windows/). Hydra is available on [Docker Hub](https://hub.docker.com/r/oryam/hydra/).

You can use Hydra without a database, but be aware that restarting, scaling
or stopping the container will **lose all data**:

```
$ docker run -d -p 4444:4444 oryam/hydra --name my-hydra
ec91228cb105db315553499c81918258f52cee9636ea2a4821bdb8226872f54b
```

**Using the client command line interface:** You can ssh into the hydra container
and execute the hydra command from there:

```
$ docker exec -i -t <hydra-container-id> /bin/bash
# e.g. docker exec -i -t ec91228 /bin/bash

root@ec91228cb105:/go/src/github.com/ory-am/hydra# hydra
Hydra is a twelve factor OAuth2 and OpenID Connect provider

[...]
```

#### Building from source

If you wish to compile hydra yourself, you need to install and set up [Go 1.5+](https://golang.org/) and add `$GOPATH/bin`
to your `$PATH`. To do so, run the following commands in a shell (bash, sh, cmd.exe, ...):

```
go get github.com/ory-am/hydra
go get github.com/Masterminds/glide
cd $GOPATH/src/github.com/ory-am/hydra
glide install
go install github.com/ory-am/hydra
hydra
```

## Security

*Why should I use Hydra? It's not that hard to implement two OAuth2 endpoints and there are numerous SDKs out there!*

OAuth2 and OAuth2 related specifications are over 200 written pages. Implementing OAuth2 is easy, getting it right is hard.
Even if you use a secure SDK (there are numerous SDKs not secure by design in the wild), messing up the implementation
is a real threat - no matter how good you or your team is. To err is human.

An in-depth list of security features is listed [in the security guide](https://ory-am.gitbooks.io/hydra/content/faq/security.html).

## Reception

Hydra has received a lot of positive feedback. Let's see what the community is saying:

> Nice! Lowering barriers to the use of technologies like these is important.

[Pyxl101](https://news.ycombinator.com/item?id=11798641)

> OAuth is a framework not a protocol. The security it provides can vary greatly between implementations.
Fosite (which is what this is based on) is a very good implementation from a security perspective: https://github.com/ory-am/fosite#a-word-on-security

[abritishguy](https://news.ycombinator.com/item?id=11800515)

> [...] Thanks for releasing this by the way, looks really well engineered. [...]

[olalonde](https://news.ycombinator.com/item?id=11798831)

## Documentation

### Guide

The Guide is available on [GitBook](https://ory-am.gitbooks.io/hydra/content/).

### HTTP API Documentation

The HTTP API is documented at [Apiary](http://docs.hdyra.apiary.io).

### Command Line Documentation

Run `hydra -h` or `hydra help`.

### Develop

Developing with Hydra is as easy as:

```
go get github.com/ory-am/hydra
go get github.com/Masterminds/glide
cd $GOPATH/src/github.com/ory-am/hydra
glide install
go test $(glide novendor)
```

If you want to run a Hydra instance, there are two possibilities:

Run without Database:
```
go run main.go host
```
 
Run against RethinkDB using Docker:
```
docker run --name some-rethink -d -p 8080:8080 -p 28015:28015 rethinkdb
DATABASE_URL=rethinkdb://localhost:28015/hydra go run main.go host
```

## Third-party libraries and projects

* [Hydra middleware for Gin](https://github.com/janekolszak/gin-hydra)
