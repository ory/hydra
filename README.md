# Ory/Hydra

[![Join the chat at https://gitter.im/ory-am/hydra](https://badges.gitter.im/ory-am/hydra.svg)](https://gitter.im/ory-am/hydra?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Build Status](https://travis-ci.org/ory-am/hydra.svg?branch=master)](https://travis-ci.org/ory-am/hydra)
[![Coverage Status](https://coveralls.io/repos/ory-am/hydra/badge.svg?branch=master&service=github)](https://coveralls.io/github/ory-am/hydra?branch=master)

![Hydra](hydra.png)

Hydra is a twelve factor authentication, authorization and account management service, ready for you to use in your micro service architecture. Hydra is written in go and backed by PostgreSQL or any implementation of [account/storage.go](account/storage.go).

**Note:** Don't worry, Hydra development is not halted. We are simply working on a more secure OAuth2 framework to back Hydra. Check out the [fosite project](https://github.com/ory-am/fosite). We encourage contributions!

Hydra implements TLS, different OAuth2 IETF standards and supports HTTP/2. To make things as easy as possible, hydra
comes with tools to generate TLS and RS256 PEM files, leaving you with almost zero trouble to set up.

![Hydra implements HTTP/2 and TLS.](h2tls.png)

Please be aware that Hydra is not ready for production just yet and has not been tested on a production system.
If time schedule holds, we will use it in production in Q1 2016 for an awesome business app that has yet to be revealed.
This should however not discourage you from trying out or using Hydra. Most of the HTTP endpoints have reached a stable status and should not change in a major way until the first release (0.1).

Current status:

- Development: 95% of principal components
- HTTP API: 80% (in review)
- Real world use: 20%

Hydra is being developed at [Ory](https://ory.am) because we need a lightweight and clean IAM solution for our customers.  
Join our [mailinglist](http://eepurl.com/bKT3N9) to stay on top of new developments.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [What is Hydra?](#what-is-hydra)
- [Motivation](#motivation)
- [Features](#features)
- [What do you mean by *Hydra is backend*?](#what-do-you-mean-by-hydra-is-backend)
- [HTTP/2 RESTful API](#http2-restful-api)
- [Run hydra-host](#run-hydra-host)
  - [With vagrant](#with-vagrant)
  - [Set up PostgreSQL locally](#set-up-postgresql-locally)
  - [Run as executable](#run-as-executable)
  - [Run from sourcecode](#run-from-sourcecode)
  - [Available Environment Variables](#available-environment-variables)
  - [CLI Usage](#cli-usage)
    - [Start server](#start-server)
    - [Create client](#create-client)
    - [Create user](#create-user)
    - [Create JWT RSA Key Pair](#create-jwt-rsa-key-pair)
    - [Create a TLS certificate](#create-a-tls-certificate)
    - [Import policies](#import-policies)
- [Security considerations](#security-considerations)
- [Good to know](#good-to-know)
  - [Policies](#policies)
  - [Everything is RESTful. No HTML. No Templates.](#everything-is-restful-no-html-no-templates)
  - [Sign up workflow](#sign-up-workflow)
  - [Sign in workflow](#sign-in-workflow)
    - [Authenticate with Google, Dropbox, ...](#authenticate-with-google-dropbox-)
    - [Authenticate with a hydra account](#authenticate-with-a-hydra-account)
  - [Visually confirm authorization](#visually-confirm-authorization)
  - [Principles](#principles)
- [Attributions](#attributions)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## What is Hydra?

Authentication, authorization and user account management are always lengthy to plan and implement. If you're building a micro service app
in need of these three, you are in the right place.

## Motivation

We develop Hydra because Hydra we need a lightweight and clean IAM solution for our customers. We believe that security and simplicity come together. This is why Hydra only relies on Google's Go Language, PostgreSQL and a slim dependency tree. Hydra is the simple, open source alternative to proprietary authorization solutions suited best for your micro service eco system.

*Use it, enjoy it and contribute!*

## Features

Hydra's core features in a nutshell:

* **Account Management**: Sign up, settings, password recovery
* **Access Control / Policy Decision Point / Policy Storage Point** backed by [Ladon](https://github.com/ory-am/ladon).
* Rich set of **OAuth2** features:
  * Hydra implements OAuth2 as specified at [rfc6749](http://tools.ietf.org/html/rfc6749) and [draft-ietf-oauth-v2-10](http://tools.ietf.org/html/draft-ietf-oauth-v2-10) using [osin](https://github.com/RangelReale/osin) and [osin-storage](https://github.com/ory-am/osin-storage)
  * Hydra uses self-contained Acccess Tokens as suggessted in [rfc6794#section-1.4](http://tools.ietf.org/html/rfc6749#section-1.4) by issuing JSON Web Tokens as specified at
   [https://tools.ietf.org/html/rfc7519](https://tools.ietf.org/html/rfc7519) with [RSASSA-PKCS1-v1_5 SHA-256](https://tools.ietf.org/html/rfc7519#section-8) hashing algorithm.
  * Hydra implements **OAuth2 Introspection** ([rfc7662](https://tools.ietf.org/html/rfc7662)) and **OAuth2 Revokation** ([rfc7009](https://tools.ietf.org/html/rfc7009)).
  * Hydra is able to sign users up and in through OAuth2 providers like Dropbox, LinkedIn, Google, you name it.
* Hydra does not speak HTML. We believe that the design decision to keep templates out of Hydra is a core feature. *Hydra is backend, not frontend.*
* Easy command line tools like `hydra-host jwt` for generating jwt signing key pairs or `hydra-host client create`.
* Hydra works both over HTTP (use only in development) and HTTP/2 with TLS (use in production).
* Hydra is unit and integration tested. We use [dockertest](https://github.com/ory-am/dockertest)

## What do you mean by *Hydra is backend*?

Hydra does not offer a sign in, sign up or authorize HTML page. Instead, if such action is required, Hydra redirects the user
to a predefined URL, for example `http://sign-up-app.yourservice.com/sign-up` or `http://sign-in-app.yourservice.com/sign-in`.
Additionally, a user can authenticate through another OAuth2 Provider, for example Dropbox or Google.

Take a look at the example sign up/in endpoint implementations [hydra-signin](https://github.com/ory-am/hydra/blob/master/cli/hydra-signup/main.go)
and [hydra-signup](https://github.com/ory-am/hydra/blob/master/cli/hydra-signup/main.go).

## HTTP/2 RESTful API

The API is described at [apiary](http://docs.hydra6.apiary.io/#). The API Documentation is still work in progress.

## Run hydra-host

### With vagrant

You'll need [Vagrant](https://www.vagrantup.com/), [VirtualBox](https://www.virtualbox.org/) and [Git](https://git-scm.com/)
installed on your system.

```
git clone https://github.com/ory-am/hydra.git
cd hydra
vagrant up
# Get a cup of coffee
```

You should now have a running Hydra instance! Vagrant exposes ports 9000 (HTTPS - Hydra) and 9001 (Postgres) on your localhost.
Open [https://localhost:9000/](https://localhost:9000/) to confirm that Hydra is running. You will probably have to add an exception for the
HTTP certificate because it is self-signed, but after that you should see a 404 error indicating that Hydra is running!

*hydra-host* offers different capabilities for managing your Hydra instance.
Check the [this section](#cli-usage) if you want to find out more.

You can also always access hydra-host through vagrant:

```
# Assuming, that your current working directory is /where/you/cloned/hydra
vagrant ssh
hydra-host help
```

### Set up PostgreSQL locally

**On Windows and Max OS X**, download and install [Docker Toolbox](https://www.docker.com/docker-toolbox). After starting the *Docker Quickstart Terminal*,
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
> go get -d -v github.com/ory-am/hydra/...
> go install github.com/ory-am/hydra/cli/hydra-host
> hydra-host start
```

*Note: For this to work, $GOPATH/bin [must be in your path](https://golang.org/doc/code.html#GOPATH)*

### Run from sourcecode

```
> go get -d -v github.com/ory-am/hydra/...
> # cd to project root, usually in $GOPATH/src/github.com/ory-am/hydra
> cd cli
> cd hydra-host
> go run main.go start
```

### Available Environment Variables

The CLI currently requires two environment variables:

| Variable             | Description               | Format                                        | Default   |
| -------------------- | ------------------------- | --------------------------------------------- | --------- |
| PORT                 | Which port to listen on   | number                                        | 443       |
| HOST                 | Which host to listen on   | ip or hostname                                | empty (all) |
| HOST_URL             | Hydra's host URL          | url                                           | "https://localhost:4443" |
| DATABASE_URL         | PostgreSQL Database URL   | `postgres://user:password@host:port/database` | empty     |
| BCRYPT_WORKFACTOR    | BCrypt Strength           | number                                        | `10`      |
| SIGNUP_URL           | [Sign up URL](#sign-up)   | url                                           | empty     |
| SIGNIN_URL           | [Sign in URL](#sign-in)   | url                                           | empty     |
| DROPBOX_CLIENT       | Dropbox Client ID         | string                                        | empty     |
| DROPBOX_SECRET       | Dropbox Client Secret     | string                                        | empty     |
| JWT_PUBLIC_KEY_PATH  | JWT Signing Public Key    | `./cert/rs256-public.pem` (local path)        | "../../example/cert/rs256-public.pem"  |
| JWT_PRIVATE_KEY_PATH | JWT Signing Private Key   | `./cert/rs256-private.pem` (local path)       | "../../example/cert/rs256-private.pem" |
| TLS_CERT_PATH        | TLS Certificate Path      | `./cert/cert.pem`                             | "../../example/cert/tls-cert.pem"      |
| TLS_KEY_PATH         | TLS Key Path              | `./cert/key.pem`                              | "../../example/cert/tls-key.pem"       |
| DANGEROUSLY_FORCE_HTTP | Disable HTTPS           | `force`                                       | disabled  |

### CLI Usage

```
NAME:
   hydra-host - Dragons guard your resources

USAGE:
   hydra-host [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
   client       Client actions
   account      Account actions
   start        Start the host service
   jwt          JWT actions
   tls          JWT actions
   policy       Policy actions
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h                   show help
   --generate-bash-completion
   --version, -v                print the version

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

#### Create Account

```
NAME:
   hydra-host account create - create a new account

USAGE:
   hydra-host account create [command options] <username>

OPTIONS:
   --password           the user's password
   --as-superuser       grant superuser privileges to the user
```

#### Create JWT RSA Key Pair

To generate files *rs256-private.pem* and *rs256-public.pem* in the current directory, run:

```
NAME:
   hydra-host jwt generate-keypair - Create a JWT PEM keypair.

   You can use these files by providing the environment variables JWT_PRIVATE_KEY_PATH and JWT_PUBLIC_KEY_PATH

USAGE:
   hydra-host jwt generate-keypair [command options] [arguments...]

OPTIONS:
   -s, --private-file-path "rs256-private.pem"  Where to save the private key PEM file
   -p, --public-file-path "rs256-public.pem"    Where to save the private key PEM file

```

#### Create a TLS certificate

```
NAME:
   hydra-host tls generate-dummy-certificate - Create a dummy TLS certificate and private key.

   You can use these files (in development!) by providing the environment variables TLS_CERT_PATH and TLS_KEY_PATH

USAGE:
   hydra-host tls generate-dummy-certificate [command options] [arguments...]

OPTIONS:
   -c, --certificate-file-path "tls-cert.pem"   Where to save the private key PEM file
   -k, --key-file-path "tls-key.pem"            Where to save the private key PEM file
   -u, --host                                   Comma-separated hostnames and IPs to generate a certificate for
   --sd, --start-date                           Creation date formatted as Jan 1 15:04:05 2011
   -d, --duration "8760h0m0s"                   Duration that certificate is valid for
   --ca                                         whether this cert should be its own Certificate Authority
   --rb, --rsa-bits "2048"                      Size of RSA key to generate. Ignored if --ecdsa-curve is set
   --ec, --ecdsa-curve                          ECDSA curve to use to generate a key. Valid values are P224, P256, P384, P521

```

#### Import policies

You can import policies from json files.

```
NAME:
   hydra-host policy import - Import a json file which defines an array of policies

USAGE:
   hydra-host policy import <policies1.json> <policies2.json> <policies3.json>
```

Here's an exemplary *policies.json:*

```json
[
  {
    "description": "Allow editing and deleting of personal articles and all sub resources.",
    "subject": ["{edit|delete}"],
    "effect": "allow",
    "resources": [
      "urn:flitt.net:articles:{.*}"
    ],
    "permissions": [
      "edit"
    ],
    "conditions": [
      {
        "op": "SubjectIsOwner"
      }
    ]
  },
  {
    "description": "Allow creation of personal articles and all sub resources.",
    "subject": ["create"],
    "effect": "allow",
    "resources": [
      "urn:flitt.net:articles"
    ],
    "permissions": [
      "edit"
    ],
    "conditions": [
      {
        "op": "SubjectIsOwner"
      }
    ]
  }
]
```

## Security considerations

[rfc6819](https://tools.ietf.org/html/rfc6819) provides good guidelines to keep your apps and environment secure. It is recommended to read:
* [Section 5.3](https://tools.ietf.org/html/rfc6819#section-5.3) on client app security.

## Good to know

This section covers information necessary for understanding how hydra works.

### Policies

Policies are something very powerful. I have to admit that I am a huge fan of how AWS handles policies and adopted their architecture for Hydra. Please find a more in depth documentation
at the [Ladon GitHub Repository](https://github.com/ory-am/ladon).

```
{
    // This should be a unique ID. This ID is required for database retrieval.
    id: "68819e5a-738b-41ec-b03c-b58a1b19d043",

    // A human readable description. Not required
    description: "something humanly readable",

    // Which identity does this policy affect?
    // As you can see here, you can use regular expressions inside < >.
    subjects: ["max", "peter", "<zac|ken>"],

    // Should the policy allow or deny access?
    effect: "allow",

    // Which resources this policy affects.
    // Again, you can put regular expressions in inside < >.
    resources: ["urn:something:resource_a", "urn:something:resource_b", "urn:something:foo:<.+>"],

    // Which permissions this policy affects. Supports RegExp
    // Again, you can put regular expressions in inside < >.
    permissions: ["<create|delete>", "get"],

    // Under which conditions this policy is active.
    conditions: [
        // Currently, only an exemplary SubjectIsOwner condition is available.
        {
            "op": "SubjectIsOwner"
        }
    ]
}
```

This is what a policy looks like. As you can see, we have various attributes:

* A **Subject** could be an account or an client app
* A **Resource** could be an online article or a file in a cloud drive
* A **Permission** can also be referred to as "Action" ("create" something, "delete" something, ...)
* A **Condition** can be an intelligent assertion *(e.g. is the Subject requesting access also the Resource Owner?)*. Right now, only the SubjectIsOwner Condition is defined. In the future, many more (e.g. IPAddressMatches or UserAgentMatches) will be added.
* The **Effect**, which can only be **allow** or **deny** (deny *always* overrides).

Hydra needs the following information to decide if a access request is allowed:
* Resource: Which resource is affected
* Permission: Which permission is requested
* Token: What access token is trying to perform this action
* Context: The context, for example the user ID.
* Header `Authorization: Bearer <token>` with a valid access token, so this endpoint can't be scanned by malicious anonymous users.

### Everything is RESTful. No HTML. No Templates.

Hydra never responds with HTML. There is no way to set up HTML templates for signing in, up or granting access.

### Sign up workflow

Hydra offers capabilities to sign users up. First, a registered client has to acquire an access token through the OAuth2 Workflow.
Second, the client sets up a user account through the `/accounts` endpoint.

You can set up a environment variable called `SIGNUP_URL` for Hydra to redirect users to,
when the user successfully authenticated via the OAuth2 Provider Workflow but has not an account in hydra yet.
If you leave this variable empty, a 401 Unauthorized Error will be shown instead.

### Sign in workflow

Hydra offers different methods to sign users in.

#### Authenticate with Google, Dropbox, ...

You can authenticate a user through any other OAuth2 provider, such as Google, Dropbox or Facebook. To do so, simply add
a provider query parameter to the authentication url endpoint:
```
/oauth2/auth?provider=google&client_id=123&response_type=code&redirect_uri=/callback&state=randomstate
```

The provider workflow is not standardized by any authority, has not yet been subject to a security audit and is therefore subject to change.
Unfortunately most providers do not support SSO provider endpoints so we might have to rely on the OAuth2 provider workflow for a while.

We will soon document how you can add more providers (currently **only Dropbox is supported**).

#### Authenticate with a hydra account

There are multiple ways to authenticate a hydra account:
* **Password grant type:** To do so, use the [OAuth2 PASSWORD grant type](http://oauthlib.readthedocs.org/en/latest/oauth2/grants/password.html).
At this moment, the password grant is allowed to *all clients*. This will be changed in the future.
* **Callback:** You can set up an environment variable called `SIGNIN_URL` for Hydra to redirect users to,
when a client requests authorization through the `/oauth2/auth` endpoint but is not yet authenticated.

### Visually confirm authorization

When a client is not allowed to bypass the authorization screen *("Do you want to grant app XYZ access to your private information?")*,
he will be redirected to the value of the environment variable `AUTHORIZE_URL`.

**This feature is not implemented yet.**

### Principles

* Authorization and authentication require verbose logging.
* Logging should *never* include credentials, neither passwords, secrets nor tokens.

## Attributions

* [Logo source](https://www.flickr.com/photos/pathfinderlinden/7161293044/)
