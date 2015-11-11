# Hydra
[![Build Status](https://travis-ci.org/ory-am/hydra.svg?branch=master)](https://travis-ci.org/ory-am/hydra)
[![Coverage Status](https://coveralls.io/repos/ory-am/hydra/badge.svg?branch=master&service=github)](https://coveralls.io/github/ory-am/hydra?branch=master)

![Hydra](hydra.png)

Hydra is a twelve factor authentication, authorization and account management service, ready for you to use in your micro service architecture.
Hydra is written in go and backed by PostgreSQL or any implementation of [account/storage.go](account/storage.go).

*Please be aware that Hydra is not ready for production just yet and has not been tested on a production system.
If time schedule holds, we will use it in production in Q1 2016 for an awesome business app that has yet to be revealed.*

## What is Hydra?

Authentication, authorization and user account management are always lengthy to plan and implement. If you're building a micro service app
in need of these three, you are in the right place.

## Motivation

Many authentication, authorization and user management solutions exist. Some are outdated, some come with a crazy stack, some enforce patterns you might dislike and others like [auth0.com](http://auth0.com) or [oauth.io](http://oauth.io) cost good money if you're out to scale.

Hydra was written because we needed a scalable 12factor OAuth2 consumer and provider with enterprise grade authorization and interoperability without a ton of dependencies or crazy features. That is why hydra only depends on [Go](http://golang.org) and PostgreSQL. If you don't like PostgreSQL you can easily implement other databases and use them instead. 
Hydra is completely RESTful and does not serve any template (check [caveats](#caveats) why this might affect you).

Hydra is the open source alternative to proprietary authorization solutions in the age of microservices.

*Use it, enjoy it and contribute!*

## Features

Hydra is a RESTful service providing you with things like:

* **Account Management**: Sign up, settings, password recovery
* **Access Control / Policy Management** backed by [ladon](https://github.com/ory-am/ladon)
* Hydra comes with a rich set of **OAuth2** features:
  * Hydra implements OAuth2 as specified at [rfc6749](http://tools.ietf.org/html/rfc6749) and [draft-ietf-oauth-v2-10](http://tools.ietf.org/html/draft-ietf-oauth-v2-10).
  * Hydra uses self-contained Acccess Tokens as suggessted in [rfc6794#section-1.4](http://tools.ietf.org/html/rfc6749#section-1.4) by issuing JSON Web Tokens as specified at
   [https://tools.ietf.org/html/rfc7519](https://tools.ietf.org/html/rfc7519) with [RSASSA-PKCS1-v1_5 SHA-256](https://tools.ietf.org/html/rfc7519#section-8) hashing algorithm, Hydra reduces database roundtrips.
  * Hydra implements **OAuth2 Introspection** as specified in [rfc7662](https://tools.ietf.org/html/rfc7662)

## Caveats

To make hydra suitable for every usecase we decided to exclude any sort of HTML templates. Hydra speaks only JSON. This obviously prevents Hydra from delivering a dedicated login and authorization ("Do you want to grant App Foobar access to all of your data?") page.

At this moment, the */oauth2/auth* endpoint only works, if a provider is given, for example:  
```
/oauth2/auth?provider=google&client_id=123&response_type=code&redirect_uri=/callback&state=randomstate
```

A provider should be an OAuth2 */authorization* endpoint.

To log in a user you have to use the [password grant type](https://aaronparecki.com/articles/2012/07/29/1/oauth2-simplified#others). At this moment, the password grant is allowed to *all clients*. This will be changed in the future.

We will provide an exemplary provider implementation in NodeJS which uses the password grant type to log users in and is easy to customize.

The provider workflow is not standardized by any authority, has not yet been subject to a security audit and is therefore subject to change. Unfortunately most providers do not support SSO provider endpoints so we might have to rely on the OAuth2 provider workflow for a while.

## hydra-host

Hydra host is the server side of things.

### Set up PostgreSQL locally

**On Windows and Max OS X**, download and install [docker-toolbox(https://www.docker.com/docker-toolbox). After starting the *Docker Quickstart Terminal*,
do the following:

```
> docker-machine ssh default # if you're not already ssh'ed into it
> docker run --name hydra-postgres -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres
> exit
> docker-machine ip default
# This should give you something like: 192.168.99.100

> # On Windows
> set DATABASE_URL=postgres://postgres:secret@{ip from above}:5432/postgres?sslmode=disable

> # On Mac OSX
> export DATABASE_URL=postgres://postgres:secret@{ip from above}:5432/postgres?sslmode=disable
```

**On Linux** download and install [Docker](https://www.docker.com/):

```
> docker run --name hydra-postgres -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres
> export DATABASE_URL=postgres://postgres:secret@localhost:5432/postgres?sslmode=disable
```

*Warning: This uses the postgres database, which is reserved.
For brevity the guide to creating a new database in Postgres has been skipped.*

### Run as executable

```
> go install github.com/ory-am/hydra/cli/hydra-host
> hydra-host start
```

*Note: For this to work, $GOPATH/bin [must be in your path](https://golang.org/doc/code.html#GOPATH)*

### Run from sourcecode

```
> go get -u github.com/ory-am/hydra
> # cd to project root, usually in $GOPATH/src/github.com/ory-am/hydra
> cd cli
> cd hydra-host
> go run main.go start
```

### Environment

The CLI currently requires two environment variables:

| Variable          | Description               | Format                                        | Default   |
| ----------------- | ------------------------- | --------------------------------------------- | --------- |
| DATABASE_URL      | PostgreSQL Database URL   | `postgres://user:password@host:port/database` | empty     |
| BCRYPT_WORKFACTOR | BCrypt Strength           | number                                        | `10`      |
| SIGNUP_URL        | [Sign up URL](#sign-up)   | url                                           | empty     |
| SIGNIN_URL        | [Sign in URL](#sign-in)   | url                                           | empty     |
| DROPBOX_CLIENT    | Dropbox Client ID         | string                                        | empty     |
| DROPBOX_SECRET    | Dropbox Client Secret     | string                                        | empty     |
| DROPBOX_CALLBACK  | Dropbox Redirect URL      | url                                           | empty     |


### CLI Usage

```
NAME:
   hydra-host - Dragons guard your resources.

USAGE:
   hydra-host [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
   client       client actions
   user         user actions
   start        start hydra-host
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h           show help
   --version, -v        print the version
```

#### Start server

```
NAME:
   hydra-host start - start hydra-host

USAGE:
   hydra-host start [arguments...]
```

#### Create client

```
NAME:
   hydra-host client create - Create a new client.

USAGE:
   hydra-host client create [command options] [arguments...]

OPTIONS:
   -i, --id             Set client's id
   -s, --secret         The client's secret
   -r, --redirect-url   A list of allowed redirect URLs: https://foobar.com/callback|https://bazbar.com/cb|http://localhost:3000/authcb
   --as-superuser       Grant superuser privileges to the client
```

#### Create user

```
NAME:
   hydra-host user create - create a new user

USAGE:
   hydra-host user create [command options] <email>

OPTIONS:
   --password           the user's password
   --as-superuser       grant superuser privileges to the user
```

### API

The API is loosely described at [apiary](http://docs.hydra6.apiary.io/#).

## Good to know

This section covers information necessary for understanding how hydra works.

### No templates

Hydra never responds with HTML. There is no way to set up HTML templates for signing in, up or granting access.

### Sign up

Hydra offers capabilities to sign users up. First, a registered client has to acquire an access token through the OAuth2 Workflow.
Second, the client sets up a user account through the `/accounts` endpoint.

You can set up a environment variable called `SIGNUP_URL` for Hydra to redirect users to,
when the user successfully authenticated via the OAuth2 Provider Workflow but has not an account in hydra yet.
If you leave this variable empty, a 401 Unauthorized Error will be shown instead.

### Sign in

Hydra offers capabilities to sign users in. To do so, use the PASSWORD grant type.

You can set up an environment variable called `SIGNIN_URL` for Hydra to redirect users to,
when a client requests authorization through the `/oauth2/auth` endpoint but is not yet authenticated.

**This feature is not implemented yet.**

### Visually confirm authorization

When a client is not allowed to bypass the authorization screen *("Do you want to grant app XYZ access to your private information?")*,
he will be redirected to the value of the environment variable `AUTHORIZE_URL`.

**This feature is not implemented yet.**

### Principles

* Authorization and authentication require verbose logging.
* Logging should *never* include credentials, neither passwords, secrets nor tokens.

## Attributions

* [Logo source](https://www.flickr.com/photos/pathfinderlinden/7161293044/)
