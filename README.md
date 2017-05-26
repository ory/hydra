# ![ORY Hydra](docs/images/logo.png)

[![Join the chat at https://gitter.im/ory-am/hydra](https://img.shields.io/badge/join-chat-00cc99.svg)](https://gitter.im/ory-am/hydra?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Join newsletter](https://img.shields.io/badge/join-newsletter-00cc99.svg)](http://eepurl.com/bKT3N9)
[![Follow twitter](https://img.shields.io/badge/follow-twitter-00cc99.svg)](https://twitter.com/_aeneasr)
[![Follow GitHub](https://img.shields.io/badge/follow-github-00cc99.svg)](https://github.com/arekkas)
[![Become a patron!](https://img.shields.io/badge/support%20us-on%20patreon-green.svg)](https://patreon.com/user?u=4298803)

[![Build Status](https://travis-ci.org/ory/hydra.svg?branch=master)](https://travis-ci.org/ory/hydra)
[![Coverage Status](https://coveralls.io/repos/ory/hydra/badge.svg?branch=master&service=github)](https://coveralls.io/github/ory/hydra?branch=master)
[![Code Climate](https://codeclimate.com/github/ory/hydra/badges/gpa.svg)](https://codeclimate.com/github/ory/hydra)
[![Go Report Card](https://goreportcard.com/badge/github.com/ory/hydra)](https://goreportcard.com/report/github.com/ory/hydra)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/364/badge)](https://bestpractices.coreinfrastructure.org/projects/364)

[![Docs Guide](https://img.shields.io/badge/docs-guide-blue.svg)](https://ory.gitbooks.io/hydra/content/)
[![HTTP API Documentation](https://img.shields.io/badge/docs-http%20api-blue.svg)](http://docs.hydra13.apiary.io/)
[![Code Documentation](https://img.shields.io/badge/docs-godoc-blue.svg)](https://godoc.org/github.com/ory/hydra)

ORY Hydra offers OAuth 2.0 and OpenID Connect Core 1.0 capabilities as a service and is built on top of the security-first
OAuth2 and OpenID Connect SDK [ORY Fosite](https://github.com/ory/fosite) and the access control
SDK [ORY Ladon](https://github.com/ory/ladon). ORY Hydra is different, because it works with
any existing authentication infrastructure, not just LDAP or SAML. By implementing a consent app (works with any programming language)
you build a bridge between ORY Hydra and your authentication infrastructure.

ORY Hydra is able to securely manage JSON Web Keys, and has a sophisticated policy-based access control you can use if you want to.

ORY Hydra is suitable for green- (new) and brownfield (existing) projects. If you are not familiar with OAuth 2.0 and are working
on a greenfield project, we recommend evaluating if OAuth 2.0 really serves your purpose.
**Knowledge of OAuth 2.0 is imperative in understanding what ORY Hydra does and how it works.**

Join the [ORY Hydra Newsletter](http://eepurl.com/bKT3N9) to stay on top of new developments. ORY Hydra has a lovely, active
community on [Gitter](https://gitter.im/ory-am/hydra). For advanced use cases, check out the
[Enterprise Edition](#enterprise-edition) section.

---

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [What is ORY Hydra?](#what-is-ory-hydra)
  - [ORY Hydra implements open standards](#ory-hydra-implements-open-standards)
- [Sponsors & Adopters](#sponsors-&-adopters)
  - [Sponsors](#sponsors)
  - [Adopters](#adopters)
- [ORY Hydra for Enterprise](#ory-hydra-for-enterprise)
- [Quickstart](#quickstart)
  - [5 minutes tutorial: Run your very own OAuth2 environment](#5-minutes-tutorial-run-your-very-own-oauth2-environment)
  - [Installation](#installation)
    - [Download binaries](#download-binaries)
    - [Using Docker](#using-docker)
    - [Building from source](#building-from-source)
- [Security](#security)
- [Documentation](#documentation)
  - [Guide](#guide)
  - [HTTP API documentation](#http-api-documentation)
  - [Command line documentation](#command-line-documentation)
  - [Develop](#develop)
- [Reception](#reception)
- [Libraries and third-party projects](#libraries-and-third-party-projects)
- [Blog posts & articles](#blog-posts-&-articles)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## What is ORY Hydra?

ORY Hydra is a server implementation of the OAuth 2.0 authorization framework and the OpenID Connect Core 1.0. Existing OAuth2
implementations usually ship as libraries or SDKs such as [node-oauth2-server](https://github.com/oauthjs/node-oauth2-server)
or [fosite](https://github.com/ory/fosite/issues), or as fully featured identity solutions with user
management and user interfaces, such as [Dex](https://github.com/coreos/dex).

Implementing and using OAuth2 without understanding the whole specification is challenging and prone to errors, even when
SDKs are being used. The primary goal of ORY Hydra is to make OAuth 2.0 and OpenID Connect 1.0 better accessible.

ORY Hydra implements the flows described in OAuth2 and OpenID Connect 1.0 without forcing you to use a "Hydra User Management"
or some template engine or a predefined front-end. Instead it relies on HTTP redirection and cryptographic methods
to verify user consent allowing you to use ORY Hydra with any authentication endpoint, be it [authboss](https://github.com/go-authboss/authboss),
[auth0.com](https://auth0.com/) or your proprietary PHP authentication.

### ORY Hydra implements open standards

ORY Hydra implements Open Standards set by the IETF:

* [The OAuth 2.0 Authorization Framework](https://tools.ietf.org/html/rfc6749)
* [OAuth 2.0 Threat Model and Security Considerations](https://tools.ietf.org/html/rfc6819)
* [OAuth 2.0 Token Revocation](https://tools.ietf.org/html/rfc7009)
* [OAuth 2.0 Token Introspection](https://tools.ietf.org/html/rfc7662)
* [OAuth 2.0 Dynamic Client Registration Protocol](https://tools.ietf.org/html/rfc7591)
* [OAuth 2.0 Dynamic Client Registration Management Protocol](https://tools.ietf.org/html/rfc7592)
* [OAuth 2.0 for Native Apps](https://tools.ietf.org/html/draft-ietf-oauth-native-apps-10)

and the OpenID Foundation:

* [OpenID Connect Core 1.0](http://openid.net/specs/openid-connect-core-1_0.html)
* [OpenID Connect Discovery 1.0](https://openid.net/specs/openid-connect-discovery-1_0.html)

## Sponsors & Adopters

This is a cureated list of Hydra sponsors and adopters. If you want to be on this list, [contact us](mailto:hi@ory.am).

### Sponsors

<img src="docs/images/sponsors/auth0.png" align="left" width="30%" alt="Auth0.com" />

We are proud to have [Auth0](https://auth0.com) as a **gold sponsor** for ORY Hydra. [Auth0](https://auth0.com) solves
the most complex identity use cases with an extensible and easy to integrate platform that secures billions of logins
every year. At ORY, we use [Auth0](https://auth0.com) in conjunction with ORY Hydra for various internal projects.

<br clear="all"/>

### Adopters

ORY Hydra is battle-tested in production systems. This is a curated list of ORY Hydra adopters.

<img src="https://www.arduino.cc/en/uploads/Trademark/ArduinoCommunityLogo.png" align="left" width="30%" alt="arduino.cc"/> 

<p>Arduino is an open-source electronics platform based on easy-to-use hardware and software. It's intended
for anyone making interactive projects. ORY Hydra secures Arduino's developer platform.</p>

<br clear="all"/>

## ORY Hydra for Enterprise

ORY Hydra is available as an Apache 2.0-licensed Open Source technology. In enterprise environments however,
there are numerous demands, such as

* OAuth 2.0 and OpenID Connect consulting.
* security auditing and certification.
* auditable log trails.
* guaranteed performance metrics, such as throughput per second.
* management user interfaces.
* ... and a wide range of narrow use cases specific to each business demands.

Gain access to more features and our security experts with ORY Hydra for Enterprise! [Request details now!](https://docs.google.com/forms/d/e/1FAIpQLSf53GJwQxzIatTSEM7sXhpkWRh6kddKxzNfNAQ9GsLNEfuFRA/viewform)

## Quickstart

This section is a quickstart guide to working with ORY Hydra. In-depth docs are available as well:

* The documentation is available on [GitBook](https://ory.gitbooks.io/hydra/content/).
* The REST API documentation is available at [Apiary](http://docs.hydra13.apiary.io/).

### 5 minutes tutorial: Run your very own OAuth2 environment

The **[tutorial](https://ory.gitbooks.io/hydra/content/tutorial.html)** teaches you to set up ORY Hydra,
a Postgres instance and an exemplary identity provider written in React using docker compose.
It will take you about 5 minutes to complete the **[tutorial](https://ory.gitbooks.io/hydra/content/tutorial.html)**.

<img src="docs/images/oauth2-flow.gif" alt="OAuth2 Flow">

<br clear="all">

### Installation

There are various ways of installing ORY Hydra on your system.

#### Download binaries

The client and server **binaries are downloadable at [releases](https://github.com/ory/hydra/releases)**.
There is currently no installer available. You have to add the ORY Hydra binary to the PATH environment variable yourself or put
the binary in a location that is already in your path (`/usr/bin`, ...). 
If you do not understand what that all of this means, ask in our [chat channel](https://gitter.im/ory-am/hydra). We are happy to help.

#### Using Docker

**Starting the host** is easiest with docker. The host process handles HTTP requests and is backed by a database.
Read how to install docker on [Linux](https://docs.docker.com/linux/), [OSX](https://docs.docker.com/mac/) or
[Windows](https://docs.docker.com/windows/). ORY Hydra is available on [Docker Hub](https://hub.docker.com/r/oryd/hydra/).

You can use ORY Hydra without a database, but be aware that restarting, scaling
or stopping the container will **lose all data**:

```
$ docker run -d --name my-hydra -p 4444:4444 oryd/hydra
ec91228cb105db315553499c81918258f52cee9636ea2a4821bdb8226872f54b
```

*Note: We had to create a new docker hub repository. Tags prior to 0.7.5 are available [here](https://hub.docker.com/r/ory-am/hydra/).*

**Using the client command line interface:** You can ssh into the ORY Hydra container
and execute the ORY Hydra command from there:

```
$ docker exec -i -t <hydra-container-id> /bin/bash
# e.g. docker exec -i -t ec91228 /bin/bash

root@ec91228cb105:/go/src/github.com/ory/hydra# hydra
Hydra is a twelve factor OAuth2 and OpenID Connect provider

[...]
```

#### Building from source

If you wish to compile ORY Hydra yourself, you need to install and set up [Go 1.8+](https://golang.org/) and add `$GOPATH/bin`
to your `$PATH`. To do so, run the following commands in a shell (bash, sh, cmd.exe, ...):

```
go get -d github.com/ory/hydra
go get github.com/Masterminds/glide
cd $GOPATH/src/github.com/ory/hydra
glide install
go install github.com/ory/hydra
hydra
```

**Note:** We changed organization name from `ory-am` to `ory`. In order to keep backwards compatibility, we did not
rename Go packages.

## Security

*Why should I use ORY Hydra? It's not that hard to implement two OAuth2 endpoints and there are numerous SDKs out there!*

OAuth2 and OAuth2 related specifications are over 400 written pages. Implementing OAuth2 is easy, getting it right is hard.
ORY Hydra is trusted by companies all around the world, has a vibrant community and faces millions of requests in production
each day. Of course, we also compiled a security guide with more details on cryptography and security concepts.
Read [the security guide now](https://ory.gitbooks.io/hydra/content/security.html).

## Documentation

### Guide

The Guide is available on [GitBook](https://ory.gitbooks.io/hydra/content/).

### HTTP API documentation

The HTTP API is documented at [Apiary](http://docs.hydra13.apiary.io/).

### Command line documentation

Run `hydra -h` or `hydra help`.

### Develop

Developing with ORY Hydra is as easy as:

```
go get -d github.com/ory/hydra
go get github.com/Masterminds/glide
cd $GOPATH/src/github.com/ory/hydra
glide install
go test $(glide novendor)
```

Then run it with in-memory database:

```
DATABASE_URL=memory go run main.go host
```

## Reception

Hydra has received a lot of positive feedback. Let's see what the community is saying:

> Nice! Lowering barriers to the use of technologies like these is important.

[Pyxl101](https://news.ycombinator.com/item?id=11798641)

> OAuth is a framework not a protocol. The security it provides can vary greatly between implementations.
Fosite (which is what this is based on) is a very good implementation from a security perspective: https://github.com/ory/fosite#a-word-on-security

[abritishguy](https://news.ycombinator.com/item?id=11800515)

> [...] Thanks for releasing this by the way, looks really well engineered. [...]

## Libraries and third-party projects

Official:
* [Consent App SDK For NodeJS](https://github.com/ory/hydra-js)
* [Consent App Example written in Go](https://github.com/ory/hydra-consent-app-go)
* [Exemplary Consent App with Express and NodeJS](https://github.com/ory/hydra-consent-app-express)

Community:
* [Consent App SDK for Go](https://github.com/janekolszak/idp)
* [ORY Hydra middleware for Gin](https://github.com/janekolszak/gin-hydra)
* [Kubernetes helm chart](https://github.com/kubernetes/charts/pull/1022)

## Blog posts & articles

* [Creating an oauth2 custom lamda authorizer for use with Amazons (AWS) API Gateway using Hydra](https://blogs.edwardwilde.com/2017/01/12/creating-an-oauth2-custom-lamda-authorizer-for-use-with-amazons-aws-api-gateway-using-hydra/)
* Warning, ORY Hydra has changed almost everything since writing this
article: [Hydra: Run your own Identity and Access Management service in <5 Minutes](https://blog.gopheracademy.com/advent-2015/hydra-auth/)
