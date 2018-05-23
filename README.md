<h1 align="center"><img src="./docs/images/banner_hydra.png" alt="ORY Hydra - Open Source OAuth 2 and OpenID Connect server"></h1>

<h4 align="center">
    <a href="https://discord.gg/PAMQWkr">Chat</a> |
    <a href="https://community.ory.am/">Forums</a> |
    <a href="http://eepurl.com/bKT3N9">Newsletter</a><br/><br/>
    <a href="https://www.ory.sh/docs/1-hydra/">Guide</a> |
    <a href="https://www.ory.sh/docs/api/hydra">API Docs</a> |
    <a href="https://godoc.org/github.com/ory/hydra">Code Docs</a><br/><br/>
    <a href="https://opencollective.com/ory-hydra">Support this project!</a>
</h4>
 
---

ORY Hydra is a hardened OAuth2 and OpenID Connect server optimized for low-latency, high throughput,
and low resource consumption. ORY Hydra *is not* an identity provider (user sign up, user log in, password reset flow),
but connects to your existing identity provider through a [consent app](https://www.ory.sh/docs/1-hydra/2-overview/1-oauth2#consent-flow).
Implementing the consent app in a different language is easy, and exemplary consent apps
([Go](https://github.com/ory/hydra-consent-app-go), [Node](https://github.com/ory/hydra-consent-app-express)) and
[SDKs](https://www.ory.sh/docs/1-hydra/7-sdk/0-readme) are provided.

Besides mitigating various attack vectors, such as database compromisation and OAuth 2.0 weaknesses, ORY Hydra is also
able to securely manage JSON Web Keys.
[Click here](https://www.ory.sh/docs/1-hydra/4-security/0-readme) to read more about security.

<p align="left">
    <a href="https://circleci.com/gh/ory/hydra/tree/master"><img src="https://circleci.com/gh/ory/hydra/tree/master.svg?style=shield" alt="Build Status"></a>
    <a href="https://coveralls.io/github/ory/hydra?branch=master"> <img src="https://coveralls.io/repos/ory/hydra/badge.svg?branch=master&service=github" alt="Coverage Status"></a>
    <a href="https://goreportcard.com/report/github.com/ory/hydra"><img src="https://goreportcard.com/badge/github.com/ory/hydra" alt="Go Report Card"></a>
    <a href="https://bestpractices.coreinfrastructure.org/projects/364"><img src="https://bestpractices.coreinfrastructure.org/projects/364/badge" alt="CII Best Practices"></a>
    <a href="#backers" alt="sponsors on Open Collective"><img src="https://opencollective.com/ory-hydra/backers/badge.svg" /></a> <a href="#sponsors" alt="Sponsors on Open Collective"><img src="https://opencollective.com/ory-hydra/sponsors/badge.svg" /></a>
</p>

---

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [What is ORY Hydra?](#what-is-ory-hydra)
  - [OAuth2 and OpenID Connect: Open Standards!](#oauth2-and-openid-connect-open-standards)
- [Quickstart](#quickstart)
  - [5 minutes tutorial: Run your very own OAuth2 environment](#5-minutes-tutorial-run-your-very-own-oauth2-environment)
  - [Installation](#installation)
    - [Download binaries](#download-binaries)
    - [Using Docker](#using-docker)
    - [Building from source](#building-from-source)
- [Ecosystem](#ecosystem)
  - [ORY Security Console: Administrative User Interface](#ory-security-console-administrative-user-interface)
  - [ORY Oathkeeper: Identity & Access Proxy](#ory-oathkeeper-identity-&-access-proxy)
  - [ORY Keto: Access Control Policies as a Server](#ory-keto-access-control-policies-as-a-server)
- [Security](#security)
  - [Disclosing vulnerabilities](#disclosing-vulnerabilities)
- [Telemetry](#telemetry)
- [Documentation](#documentation)
  - [Guide](#guide)
  - [HTTP API documentation](#http-api-documentation)
  - [Upgrading and Changelog](#upgrading-and-changelog)
  - [Command line documentation](#command-line-documentation)
  - [Develop](#develop)
- [Libraries and third-party projects](#libraries-and-third-party-projects)
- [Blog posts & articles](#blog-posts-&-articles)
- [Contributors](#contributors)
- [Backers](#backers)
- [Sponsors](#sponsors)

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

### OAuth2 and OpenID Connect: Open Standards!

ORY Hydra implements Open Standards set by the IETF:

* [The OAuth 2.0 Authorization Framework](https://tools.ietf.org/html/rfc6749)
* [OAuth 2.0 Threat Model and Security Considerations](https://tools.ietf.org/html/rfc6819)
* [OAuth 2.0 Token Revocation](https://tools.ietf.org/html/rfc7009)
* [OAuth 2.0 Token Introspection](https://tools.ietf.org/html/rfc7662)
* [OAuth 2.0 Dynamic Client Registration Protocol](https://tools.ietf.org/html/rfc7591)
* [OAuth 2.0 Dynamic Client Registration Management Protocol](https://tools.ietf.org/html/rfc7592)
* [OAuth 2.0 for Native Apps](https://tools.ietf.org/html/draft-ietf-oauth-native-apps-10)
* [Proof Key for Code Exchange by OAuth Public Clients](https://tools.ietf.org/html/rfc7636)

and the OpenID Foundation:

* [OpenID Connect Core 1.0](http://openid.net/specs/openid-connect-core-1_0.html)
* [OpenID Connect Discovery 1.0](https://openid.net/specs/openid-connect-discovery-1_0.html)

## Quickstart

This section is a quickstart guide to working with ORY Hydra. In-depth docs are available as well:

* The documentation is available [here](https://www.ory.sh/docs/1-hydra/).
* The REST API documentation is available [here](https://www.ory.sh/docs/api/hydra).

### 5 minutes tutorial: Run your very own OAuth2 environment

The **[tutorial](https://www.ory.sh/docs/1-hydra/0-tutorial/0-readme)** teaches you to set up ORY Hydra,
a Postgres instance and an exemplary identity provider written in React using docker compose.
It will take you about 5 minutes to complete the **[tutorial](https://www.ory.sh/docs/1-hydra/0-tutorial/0-readme)**.

<img src="docs/images/oauth2-flow.gif" alt="OAuth2 Flow">

<br clear="all">

### Installation

There are various ways of installing ORY Hydra on your system.

#### Download binaries

The client and server **binaries are downloadable at [releases](https://github.com/ory/hydra/releases)**.
There is currently no installer available. You have to add the ORY Hydra binary to the PATH environment variable yourself or put
the binary in a location that is already in your path (`/usr/bin`, ...).
If you do not understand what that all of this means, ask in our [chat channel](https://www.ory.sh/chat). We are happy to help.

#### Using Docker

**Starting the host** is easiest with docker. The host process handles HTTP requests and is backed by a database.
Read how to install docker on [Linux](https://docs.docker.com/linux/), [OSX](https://docs.docker.com/mac/) or
[Windows](https://docs.docker.com/windows/). ORY Hydra is available on [Docker Hub](https://hub.docker.com/r/oryd/hydra/).

You can use ORY Hydra without a database, but be aware that restarting, scaling
or stopping the container will **lose all data**:

```
$ docker run -e "DATABASE_URL=memory" -e "ISSUER=https://localhost:4444/" -d --name my-hydra -p 4444:4444 oryd/hydra
ec91228cb105db315553499c81918258f52cee9636ea2a4821bdb8226872f54b
```

*Note: We had to create a new docker hub repository. Tags prior to 0.7.5 are available [here](https://hub.docker.com/r/ory-am/hydra/).*

**Using the client command line interface:** You can enter into the ORY Hydra container
and execute the ORY Hydra command from there:

```
$ docker exec -i -t <hydra-container-id> /bin/sh
# e.g. docker exec -i -t ec91228 /bin/sh

root@ec91228cb105:/go/src/github.com/ory/hydra# hydra
Hydra is a twelve factor OAuth2 and OpenID Connect provider

[...]
```

#### Building from source

If you wish to compile ORY Hydra yourself, you need to install and set up [Go 1.9+](https://golang.org/) and add `$GOPATH/bin`
to your `$PATH` as well as [golang/dep](http://github.com/golang/dep). To do so, run the following commands in a shell (bash, sh, cmd.exe, ...):

```
go get -d -u github.com/ory/hydra
cd $GOPATH/src/github.com/ory/hydra
dep ensure
go install github.com/ory/hydra
hydra
```

**Notes**

* We changed organization name from `ory-am` to `ory`. In order to keep backwards compatibility, we did not rename Go packages.
* You can ignore warnings similar to `package github.com/ory/hydra/cmd/server: case-insensitive import collision: "github.com/Sirupsen/logrus" and "github.com/sirupsen/logrus"`.

## Ecosystem

<a href="https://console.ory.am/auth/login">
    <img align="right" width="30%" src="docs/images/sec-console.png" alt="ORY Security Console">
</a>

### ORY Security Console: Administrative User Interface

The [ORY Security Console](https://console.ory.am/auth/login) is a visual admin interface for managing ORY Hydra,
ORY Oathkeeper, and ORY Keto.

### ORY Oathkeeper: Identity & Access Proxy

[ORY Oathkeeper](https://github.com/ory/oathkeeper) is a BeyondCorp/Zero Trust Identity & Access Proxy (IAP) built
on top of OAuth2 and ORY Hydra.

### ORY Keto: Access Control Policies as a Server

[ORY Keto](https://github.com/ory/keto) is a policy decision point. It uses a set of access control policies, similar
to AWS IAM Policies, in order to determine whether a subject (user, application, service, car, ...) is authorized to
perform a certain action on a resource.

## Security

*Why should I use ORY Hydra? It's not that hard to implement two OAuth2 endpoints and there are numerous SDKs out there!*

OAuth2 and OAuth2 related specifications are over 400 written pages. Implementing OAuth2 is easy, getting it right is hard.
ORY Hydra is trusted by companies all around the world, has a vibrant community and faces millions of requests in production
each day. Of course, we also compiled a security guide with more details on cryptography and security concepts.
Read [the security guide now](https://www.ory.sh/docs/1-hydra/0-tutorial/0-readme).

### Disclosing vulnerabilities

If you think you found a security vulnerability, please refrain from posting it publicly on the forums, the chat, or GitHub
and send us an email to [hi@ory.am](mailto:hi@ory.am) instead.

## Telemetry

Our services collect summarized, anonymized data which can optionally be turned off. Click
[here](https://www.ory.sh/docs/guides/latest/9-telemetry) to learn more.

## Documentation

### Guide

The Guide is available [here](https://www.ory.sh/docs/1-hydra/).

### HTTP API documentation

The HTTP API is documented [here](https://www.ory.sh/docs/api/hydra).

### Upgrading and Changelog

New releases might introduce breaking changes. To help you identify and incorporate those changes, we document these
changes in [UPGRADE.md](./UPGRADE.md) and [CHANGELOG.md](./CHANGELOG.md).

### Command line documentation

Run `hydra -h` or `hydra help`.

### Develop

Developing with ORY Hydra is as easy as:

```
go get -d -u github.com/ory/hydra
cd $GOPATH/src/github.com/ory/hydra
dep ensure
go test ./...
```

Then run it with in-memory database:

```
DATABASE_URL=memory go run main.go host
```

**Notes**

* We changed organization name from `ory-am` to `ory`. In order to keep backwards compatibility, we did not rename Go packages.
* You can ignore warnings similar to `package github.com/ory/hydra/cmd/server: case-insensitive import collision: "github.com/sirupsen/logrus" and "github.com/sirupsen/logrus"`.

## Libraries and third-party projects

Official:
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

## Contributors

This project exists thanks to all the people who contribute. [[Contribute](CONTRIBUTING.md)].
<a href="graphs/contributors"><img src="https://opencollective.com/ory-hydra/contributors.svg?width=890&button=false" /></a>


## Backers

Thank you to all our backers! 🙏 [[Become a backer](https://opencollective.com/ory-hydra#backer)]

<a href="https://opencollective.com/ory-hydra#backers" target="_blank"><img src="https://opencollective.com/ory-hydra/backers.svg?width=890"></a>

We would also like to thank (past & current) supporters (in alphabetical order) on [Patreon](https://www.patreon.com/_ory):

* Alexander Alimovs
* Chancy Kennedy
* Drozzy
* Oz Haven
* TheCrealm

## Sponsors

Support this project by becoming a sponsor. Your logo will show up here with a link to your website. [[Become a sponsor](https://opencollective.com/ory-hydra#sponsor)]

<a href="https://opencollective.com/ory-hydra/sponsor/0/website" target="_blank"><img src="https://opencollective.com/ory-hydra/sponsor/0/avatar.svg"></a>
<a href="https://opencollective.com/ory-hydra/sponsor/1/website" target="_blank"><img src="https://opencollective.com/ory-hydra/sponsor/1/avatar.svg"></a>
<a href="https://opencollective.com/ory-hydra/sponsor/2/website" target="_blank"><img src="https://opencollective.com/ory-hydra/sponsor/2/avatar.svg"></a>
<a href="https://opencollective.com/ory-hydra/sponsor/3/website" target="_blank"><img src="https://opencollective.com/ory-hydra/sponsor/3/avatar.svg"></a>
<a href="https://opencollective.com/ory-hydra/sponsor/4/website" target="_blank"><img src="https://opencollective.com/ory-hydra/sponsor/4/avatar.svg"></a>
<a href="https://opencollective.com/ory-hydra/sponsor/5/website" target="_blank"><img src="https://opencollective.com/ory-hydra/sponsor/5/avatar.svg"></a>
<a href="https://opencollective.com/ory-hydra/sponsor/6/website" target="_blank"><img src="https://opencollective.com/ory-hydra/sponsor/6/avatar.svg"></a>
<a href="https://opencollective.com/ory-hydra/sponsor/7/website" target="_blank"><img src="https://opencollective.com/ory-hydra/sponsor/7/avatar.svg"></a>
<a href="https://opencollective.com/ory-hydra/sponsor/8/website" target="_blank"><img src="https://opencollective.com/ory-hydra/sponsor/8/avatar.svg"></a>
<a href="https://opencollective.com/ory-hydra/sponsor/9/website" target="_blank"><img src="https://opencollective.com/ory-hydra/sponsor/9/avatar.svg"></a>

A special thanks goes out to **Wayne Robinson** for supporting this project with $200 every month since Oktober 2016 [on Patreon](https://www.patreon.com/_ory).
