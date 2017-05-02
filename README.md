# ![Ory/Hydra](docs/images/logo.png)

[![Join the chat at https://gitter.im/ory-am/hydra](https://img.shields.io/badge/join-chat-00cc99.svg)](https://gitter.im/ory-am/hydra?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Join mailinglist](https://img.shields.io/badge/join-mailinglist-00cc99.svg)](https://groups.google.com/forum/#!forum/ory-hydra/new)
[![Join newsletter](https://img.shields.io/badge/join-newsletter-00cc99.svg)](http://eepurl.com/bKT3N9)
[![Follow twitter](https://img.shields.io/badge/follow-twitter-00cc99.svg)](https://twitter.com/_aeneasr)
[![Follow GitHub](https://img.shields.io/badge/follow-github-00cc99.svg)](https://github.com/arekkas)
[![Become a patron!](https://img.shields.io/badge/support%20us-on%20patreon-green.svg)](https://patreon.com/user?u=4298803)

[![Build Status](https://travis-ci.org/ory/hydra.svg?branch=master)](https://travis-ci.org/ory/hydra)
[![Coverage Status](https://coveralls.io/repos/ory/hydra/badge.svg?branch=master&service=github)](https://coveralls.io/github/ory/hydra?branch=master)
[![Code Climate](https://codeclimate.com/github/ory/hydra/badges/gpa.svg)](https://codeclimate.com/github/ory/hydra)
[![Go Report Card](https://goreportcard.com/badge/github.com/ory/hydra)](https://goreportcard.com/report/github.com/ory/hydra)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/364/badge)](https://bestpractices.coreinfrastructure.org/projects/364)

---

[![Docs Guide](https://img.shields.io/badge/docs-guide-blue.svg)](https://ory.gitbooks.io/hydra/content/)
[![HTTP API Documentation](https://img.shields.io/badge/docs-http%20api-blue.svg)](http://docs.hydra13.apiary.io/)
[![Code Documentation](https://img.shields.io/badge/docs-godoc-blue.svg)](https://godoc.org/github.com/ory/hydra)

Hydra offers OAuth 2.0 and OpenID Connect Core 1.0 capabilities as a service. Hydra is different, because it works with
any existing authentication infrastructure, not just LDAP or SAML. By implementing a consent app (works with any programming language)
you build a bridge between Hydra and your authentication infrastructure.

Hydra is able to securely manage JSON Web Keys, and has a sophisticated policy-based access control you can use if you want to.

Hydra is suitable for green- (new) and brownfield (existing) projects. If you are not familiar with OAuth 2.0 and are working
on a greenfield project, we recommend evaluating if OAuth 2.0 really serves your purpose. **Knowledge of OAuth 2.0 is imperative in understanding what Hydra does and how it works.**

---

Join our [newsletter](http://eepurl.com/bKT3N9) to stay on top of new developments.
We answer basic support requests on [Google Groups](https://groups.google.com/forum/#!forum/ory-hydra/new) and [Gitter](https://gitter.im/ory-am/hydra)
and offer [premium services](http://www.ory.am/products/hydra) around Hydra.

Hydra uses the security first OAuth2 and OpenID Connect SDK [Fosite](https://github.com/ory/fosite) and
the access control SDK [Ladon](https://github.com/ory/ladon).

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [What is Hydra?](#what-is-hydra)
- [Enterprise Edition](#enterprise-edition)
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
  - [HTTP API documentation](#http-api-documentation)
  - [Command line documentation](#command-line-documentation)
  - [Develop](#develop)
- [Sponsors](#sponsors)
- [Libraries and third-party projects](#libraries-and-third-party-projects)
- [Blog posts & articles](#blog-posts-&-articles)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## What is Hydra?

Hydra is a server implementation of the OAuth 2.0 authorization framework and the OpenID Connect Core 1.0. Existing OAuth2
implementations usually ship as libraries or SDKs such as [node-oauth2-server](https://github.com/oauthjs/node-oauth2-server)
or [fosite](https://github.com/ory/fosite/issues), or as fully featured identity solutions with user
management and user interfaces, such as [Dex](https://github.com/coreos/dex).

Implementing and using OAuth2 without understanding the whole specification is challenging and prone to errors, even when
SDKs are being used. The primary goal of Hydra is to make OAuth 2.0 and OpenID Connect 1.0 better accessible.

Hydra implements the flows described in OAuth2 and OpenID Connect 1.0 without forcing you to use a "Hydra User Management"
or some template engine or a predefined front-end. Instead it relies on HTTP redirection and cryptographic methods
to verify user consent allowing you to use Hydra with any authentication endpoint, be it [authboss](https://github.com/go-authboss/authboss),
[auth0.com](https://auth0.com/) or your proprietary PHP authentication.

## Enterprise Edition

Hydra is available as an Apache 2.0-licensed Open Source technology. In enterprise environments however,
there are numerous demands, such as

* OAuth 2.0 and OpenID Connect consulting.
* security auditing and certification.
* auditable log trails.
* guaranteed performance metrics, such as throughput per second.
* management user interfaces.
* ... and a wide range of narrow use cases specific to each business demands.

Gain access to more features and our security experts with our new enterprise edition of Hydra. **[Contact us now](mailto:hi@ory.am) for more details.**

## Quickstart

This section is a quickstart guide to working with Hydra. In-depth docs are available as well:

* The documentation is available on [GitBook](https://ory.gitbooks.io/hydra/content/).
* The REST API documentation is available at [Apiary](http://docs.hydra13.apiary.io/).

### 5 minutes tutorial: Run your very own OAuth2 environment

The **[tutorial](https://ory.gitbooks.io/hydra/content/tutorial.html)** teaches you to set up Hydra,
a Postgres instance and an exemplary identity provider written in React using docker compose.
It will take you about 5 minutes to complete the **[tutorial](https://ory.gitbooks.io/hydra/content/tutorial.html)**.

<img src="docs/images/oauth2-flow.gif" alt="OAuth2 Flow">

<br clear="all">

### Installation

There are various ways of installing hydra on your system.

#### Download binaries

The client and server **binaries are downloadable at [releases](https://github.com/ory/hydra/releases)**.
There is currently no installer available. You have to add the hydra binary to the PATH environment variable yourself or put
the binary in a location that is already in your path (`/usr/bin`, ...). 
If you do not understand what that all of this means, ask in our [chat channel](https://gitter.im/ory-am/hydra). We are happy to help.

#### Using Docker

**Starting the host** is easiest with docker. The host process handles HTTP requests and is backed by a database.
Read how to install docker on [Linux](https://docs.docker.com/linux/), [OSX](https://docs.docker.com/mac/) or
[Windows](https://docs.docker.com/windows/). Hydra is available on [Docker Hub](https://hub.docker.com/r/oryd/hydra/).

You can use Hydra without a database, but be aware that restarting, scaling
or stopping the container will **lose all data**:

```
$ docker run -d --name my-hydra -p 4444:4444 oryd/hydra
ec91228cb105db315553499c81918258f52cee9636ea2a4821bdb8226872f54b
```

*Note: We had to create a new docker hub repository. Tags prior to 0.7.5 are available [here](https://hub.docker.com/r/ory-am/hydra/).*

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

**Note:** We changed organization name from `ory-am` to `ory`. In order to keep backwards compatibility, we did not
rename Go packages.

## Security

*Why should I use Hydra? It's not that hard to implement two OAuth2 endpoints and there are numerous SDKs out there!*

OAuth2 and OAuth2 related specifications are over 200 written pages. Implementing OAuth2 is easy, getting it right is hard.
Even if you use a secure SDK (there are numerous SDKs not secure by design in the wild), messing up the implementation
is a real threat - no matter how good you or your team is. To err is human.

An in-depth list of security features is listed [in the security guide](https://ory.gitbooks.io/hydra/content/faq/security.html).

## Reception

Hydra has received a lot of positive feedback. Let's see what the community is saying:

> Nice! Lowering barriers to the use of technologies like these is important.

[Pyxl101](https://news.ycombinator.com/item?id=11798641)

> OAuth is a framework not a protocol. The security it provides can vary greatly between implementations.
Fosite (which is what this is based on) is a very good implementation from a security perspective: https://github.com/ory/fosite#a-word-on-security

[abritishguy](https://news.ycombinator.com/item?id=11800515)

> [...] Thanks for releasing this by the way, looks really well engineered. [...]

[olalonde](https://news.ycombinator.com/item?id=11798831)

## Documentation

### Guide

The Guide is available on [GitBook](https://ory.gitbooks.io/hydra/content/).

### HTTP API documentation

The HTTP API is documented at [Apiary](http://docs.hydra13.apiary.io/).

### Command line documentation

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

## Sponsors

<img src="docs/images/sponsors/auth0.png" align="left" width="30%" alt="Auth0.com"> We are proud to have [Auth0](https://auth0.com) as a **gold sponsor** for Hydra. [Auth0](https://auth0.com) solves the most complex identity use cases with an extensible and easy to integrate platform that secures billions of logins every year. At ORY, we use [Auth0](https://auth0.com) in conjunction with Hydra for various internal projects.

## Libraries and third-party projects

Official:
* [Consent App SDK For NodeJS](https://github.com/ory/hydra-js)
* [Consent App Example written in Go](https://github.com/ory/hydra-consent-app-go)
* [Exemplary Consent App with Express and NodeJS](https://github.com/ory/hydra-consent-app-express)

Community:
* [Consent App SDK for Go](https://github.com/janekolszak/idp)
* [Hydra middleware for Gin](https://github.com/janekolszak/gin-hydra)

## Blog posts & articles

* [Creating an oauth2 custom lamda authorizer for use with Amazons (AWS) API Gateway using Hydra](https://blogs.edwardwilde.com/2017/01/12/creating-an-oauth2-custom-lamda-authorizer-for-use-with-amazons-aws-api-gateway-using-hydra/)
* Warning, Hydra has changed almost everything since writing this article: [Hydra: Run your own Identity and Access Management service in <5 Minutes](https://blog.gopheracademy.com/advent-2015/hydra-auth/)
