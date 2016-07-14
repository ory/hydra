# ![Ory/Hydra](dist/logo.png)

[![Build Status](https://travis-ci.org/ory-am/hydra.svg?branch=master)](https://travis-ci.org/ory-am/hydra)
[![Coverage Status](https://coveralls.io/repos/ory-am/hydra/badge.svg?branch=master&service=github)](https://coveralls.io/github/ory-am/hydra?branch=master)
[![Code Climate](https://codeclimate.com/github/ory-am/hydra/badges/gpa.svg)](https://codeclimate.com/github/ory-am/hydra)
[![Go Report Card](https://goreportcard.com/badge/github.com/ory-am/hydra)](https://goreportcard.com/report/github.com/ory-am/hydra)

[![Join the chat at https://gitter.im/ory-am/hydra](https://img.shields.io/badge/join-chat-00cc99.svg)](https://gitter.im/ory-am/hydra?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Join mailinglist](https://img.shields.io/badge/join-mailinglist-00cc99.svg)](https://groups.google.com/forum/#!forum/ory-hydra/new)
[![Join newsletter](https://img.shields.io/badge/join-newsletter-00cc99.svg)](http://eepurl.com/bKT3N9)
[![Follow newsletter](https://img.shields.io/badge/follow-twitter-00cc99.svg)](https://twitter.com/_aeneasr)
[![Follow GitHub](https://img.shields.io/badge/follow-github-00cc99.svg)](https://github.com/arekkas)

Hydra is being developed by german-based company [Ory](https://ory.am). Join our [newsletter](http://eepurl.com/bKT3N9) to stay on top of new developments.
We respond on [Google Groups](https://groups.google.com/forum/#!forum/ory-hydra/new) and [Gitter](https://gitter.im/ory-am/hydra). Our timezone is CET.

Hydra uses the security first OAuth2 and OpenID Connect SDK [Fosite](https://github.com/ory-am/fosite) and [Ladon](https://github.com/ory-am/ladon) for policy-based access control.

:fire: Don't want to worry about security updates, backups, migration and scaling? Do you need enterprise support? [Contact us now](mailto:sales@ory.am) and become an early adopter. :fire:

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [What is Hydra?](#what-is-hydra)
  - [Feature Overview](#feature-overview)
  - [Can I use Hydra in my new or existing app?](#can-i-use-hydra-in-my-new-or-existing-app)
- [Security](#security)
- [Quickstart](#quickstart)
  - [Installation](#installation)
    - [Server](#server)
      - [Downloading an image form the Hub](#downloading-an-image-form-the-hub)
      - [Build a docker image from the sources](#build-a-docker-image-from-the-sources)
    - [Client](#client)
      - [Using gopack](#using-gopack)
      - [Building from the sources](#building-from-the-sources)
      - [From the Docker container (not recommended)](#from-the-docker-container-not-recommended)
  - [Run the example](#run-the-example)
- [Documentation](#documentation)
  - [Guide](#guide)
  - [REST API Documentation](#rest-api-documentation)
  - [CLI Documentation](#cli-documentation)
  - [Develop](#develop)
- [FAQ](#faq)
  - [I have troubles with the redirect URI.](#i-have-troubles-with-the-redirect-uri)
  - [How can I validate tokens?](#how-can-i-validate-tokens)
  - [How can I import TLS certificates?](#how-can-i-import-tls-certificates)
  - [I want to disable HTTPS for testing](#i-want-to-disable-https-for-testing)
  - [Is there a client library / SDK?](#is-there-a-client-library--sdk)
- [Hall of Fame](#hall-of-fame)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## What is Hydra?

At first, there was the monolith. The monolith worked well with the bespoke authentication module.
Then, the web evolved into an elastic cloud that serves thousands of different user agents
in every part of the world.

Hydra is driven by the need for a **scalable, low-latency, in memory
Access Control, OAuth2, and OpenID Connect layer** that integrates with every identity provider you can imagine.

Hydra is available through [Docker](https://hub.docker.com/r/oryam/hydra/) and relies on RethinkDB for persistence.
Database drivers are extensible in case you want to use RabbitMQ, MySQL, MongoDB, or some other database instead.

Hydra is built for high throughput environments. Check out the below siege benchmark on a Macbook Pro Late 2013, connected to RethinkDB validating access tokens.

*Note*: There is no cache involved. You could use a different access token for every request and still get the same benchmark results.

```
# 500 concurrent connections, 5 minutes stress test
$ siege -t5m -c=500 -b --content-type="application/json" --header="Authorization: bearer 14Q28Y5PILTSKhByweAmQ1JddECMtp4j65yvlKs5Rhk.mqkREUzhqKexxGct5tA4UhJSnL8RYMpwMeFlEF-10OM" 'https://localhost:4444/warden/authorized POST {"assertion": "14Q28Y5PILTSKhByweAmQ1JddECMtp4j65yvlKs5Rhk.mqkREUzhqKexxGct5tA4UhJSnL8RYMpwMeFlEF-10OM"}' | grep '^[^HTTP]+'

Lifting the server siege...
Transactions:           155455 hits
Availability:           100.00 %
Elapsed time:           299.67 secs
Data transferred:        32.47 MB
Response time:	          0.00 secs
Transaction rate:       518.75 trans/sec
Throughput:               0.11 MB/sec
Concurrency:  	           0.96
Successful transactions: 155455
Failed transactions:          0
Longest transaction:       6.80
Shortest transaction:      0.00
```

### Feature Overview

1. **Availability:** Hydra uses pub/sub to have the latest data available in memory. The in-memory architecture allows for heavy duty workloads.
2. **Scalability:** Hydra scales effortlessly on every platform you can imagine, including Heroku, Cloud Foundry, Docker,
Google Container Engine and many more.
3. **Integration:** Hydra wraps your existing stack like a blanket and keeps it safe. Hydra uses cryptographic tokens to authenticate users and request their consent, no APIs required.
The deprecated php-3.0 authentication service your intern wrote? It works with that too, don't worry.
We wrote an example with React to show you what this could look like: [React.js Identity Provider Example App](https://github.com/ory-am/hydra-idp-react).
4. **Security:** Hydra leverages the security first OAuth2 framework **[Fosite](https://github.com/ory-am/fosite)**,
encrypts important data at rest, and supports HTTP over TLS (https) out of the box.
5. **Ease of use:** Developers and operators are human. Therefore, Hydra is easy to install and manage. Hydra does not care if you use React, Angular, or Cocoa for your user interface.
To support you even further, there are APIs available for *cryptographic key management, social log on, policy based access control, policy management, and two factor authentication (tbd).*
Hydra is packaged using [Docker](https://hub.docker.com/r/oryam/hydra/).
6. **Open Source:** Hydra is licensed under Apache Version 2.0
7. **Professional:** Hydra implements peer reviewed open standards published by [The Internet Engineering Task Force (IETF®)](https://www.ietf.org/) and the [OpenID Foundation](https://openid.net/)
and under supervision of the [LMU Teaching and Research Unit Programming and Modelling Languages](http://www.en.pms.ifi.lmu.de). No funny business.
8. **Real Time:** Operation is a lot easier with real time monitoring. Because Hydra leverages RethinkDB, you get real time monitoring for free:
![monitoring.gif](dist/monitoring.gif)

### Can I use Hydra in my new or existing app?

OAuth2 and OpenID Connect are tricky to understand. It is important to understand that OAuth2 is
a delegation protocol. It makes sense to use Hydra in new and existing projects. A use case covering an existing project
explains how one would use Hydra in a new one as well. So let's look at a use case!

Let's assume we are running a ToDo List App (todo24.com). ToDo24 has a login endpoint (todo24.com/login).
The login endpoint is written in node and uses MongoDB to store user information (email + password + settings). Of course,
todo24 has other services as well: list management (todo24.com/lists/manage: close, create, move), item management (todo24.com/lists/items/manage: mark solved, add), and so on.
You are using cookies to see which user is performing the request.

Now you decide to use OAuth2 on top of your current infrastructure. There are many reasons to do this:
* You want to open up your APIs to third-party developers. Their apps will be using OAuth2 Access Tokens to access a user's to do list.
* You want a mobile client. Because you can not store secrets on devices (they can be reverse engineered and stolen), you use OAuth2 Access Tokens instead.
* You have Cross Origin Requests. Making cookies work with Cross Origin Requests weakens or even disables important anti-CSRF measures.
* You want to write an in-browser client. This is the same case as in a mobile client (you can't store secrets in a browser).

These are only a couple of reasons to use OAuth2. You might decide to use OAuth2 as your single source of authorization, thus maintaining
only one authorization protocol and being able to open up to third party devs in no time. With OpenID Connect, you are able to delegate authentication as well as authorization!

Your decision is final. You want to use OAuth2 and you want Hydra to do the job. You install Hydra in your cluster using docker.
Next, you set up some exemplary OAuth2 clients. Clients can act on their own, but most of the time they need to access a user's todo lists.
To do so, the client initiates an OAuth2 request. This is where [Hydra's authentication flow](https://ory-am.gitbooks.io/hydra/content/oauth2.html#authentication-flow) comes in to play.
Before Hydra can issue an access token, we need to know WHICH user is giving consent. To do so, Hydra redirects the user agent (e.g. browser, mobile device)
to the login endpoint alongside with a challenge that contains an expiry time and other information. The login endpoint (todo24.com/login) authenticates the
user as usual, e.g. by username & password, session cookie or other means. Upon successful authentication, the login endpoint asks for the user's consent:
*"Do you want to grant MyCoolAnalyticsApp read & write access to all your todo lists? [Yes] [No]"*. Once the user clicks *Yes* and gives consent,
the login endpoint redirects back to hydra and appends something called a *consent token*. The consent token is a cryptographically signed
string that contains information about the user, specifically the user's unique id. Hydra validates the signature's trustworthiness
and issues an OAuth2 access token and optionally a refresh or OpenID token.

Every time a request containing an access token hits a resource server (todo24.com/lists/manage), you make a request to Hydra asking who the token's
subject (the user who authorized the client to create a token on its behalf) is and whether the token is valid or not. You may optionally
ask if the token has permission to perform a certain action.

## Security

*Why should I use Hydra? It's not that hard to implement two OAuth2 endpoints and there are numerous SDKs out there!*

OAuth2 and OAuth2 related specifications are over 200 written pages. Implementing OAuth2 is easy, getting it right is hard.
Even if you use a secure SDK (there are numerous SDKs not secure by design in the wild), messing up the implementation
is a real threat - no matter how good you or your team is. To err is human.

Let's take a look at security in Hydra:
* Hydra uses [Fosite](https://github.com/ory-am/fosite#a-word-on-security), a secure-by-design OAuth2 SDK. Fosite implements
best practices proposed by the IETF:
    * [No Cleartext Storage of Credentials](https://tools.ietf.org/html/rfc6819#section-5.1.4.1.3)
    * [Encryption of Credentials](https://tools.ietf.org/html/rfc6819#section-5.1.4.1.4)
    * [Use Short Expiration Time](https://tools.ietf.org/html/rfc6819#section-5.1.5.3)
    * [Limit Number of Usages or One-Time Usage](https://tools.ietf.org/html/rfc6819#section-5.1.5.4)
    * [Bind Token to Client id](https://tools.ietf.org/html/rfc6819#section-5.1.5.8)
    * [Automatic Revocation of Derived Tokens If Abuse Is Detected](https://tools.ietf.org/html/rfc6819#section-5.2.1.1)
    * [Binding of Refresh Token to "client_id"](https://tools.ietf.org/html/rfc6819#section-5.2.2.2)
    * [Refresh Token Rotation](https://tools.ietf.org/html/rfc6819#section-5.2.2.3)
    * [Revocation of Refresh Tokens](https://tools.ietf.org/html/rfc6819#section-5.2.2.4)
    * [Validate Pre-Registered "redirect_uri"](https://tools.ietf.org/html/rfc6819#section-5.2.3.5)
    * [Binding of Authorization "code" to "client_id"](https://tools.ietf.org/html/rfc6819#section-5.2.4.4)
    * [Binding of Authorization "code" to "redirect_uri"](https://tools.ietf.org/html/rfc6819#section-5.2.4.6)
    * [Opaque access tokens](https://tools.ietf.org/html/rfc6749#section-1.4)
    * [Opaque refresh tokens](https://tools.ietf.org/html/rfc6749#section-1.5)
    * [Ensure Confidentiality of Requests](https://tools.ietf.org/html/rfc6819#section-5.1.1)
    * [Use of Asymmetric Cryptography](https://tools.ietf.org/html/rfc6819#section-5.1.4.1.5)
    * **Enforcing random states:** Without a random-looking state or OpenID Connect nonce the request will fail.
    * **Advanced Token Validation:** Tokens are laid out as `<key>.<signature>` where `<signature>` is created using HMAC-SHA256
     and a global secret. This is what a token can look like: `/tgBeUhWlAT8tM8Bhmnx+Amf8rOYOUhrDi3pGzmjP7c=.BiV/Yhma+5moTP46anxMT6cWW8gz5R5vpC9RbpwSDdM=`
    * **Enforcing scopes:** By default, you always need to include the `core` scope or Hydra will not execute the request.
* Hydra uses [Ladon](https://github.com/ory-am/ladon) for policy management and access control. Ladon's API is minimalistic
and well tested.
* Hydra encrypts symmetric and asymmetric keys at rest using AES-GCM 256bit.
* Hydra does not store tokens, only their signatures. An attacker gaining database access is neither able to steal tokens nor
to issue new ones.
* Hydra has automated unit and integration tests.
* Hydra does not use hacks. We would rather rewrite the whole thing instead of introducing a hack.
* APIs are uniform, well documented and secured using the warden's access control infrastructure.
* Hydra is open source and can be reviewed by anyone.
* Hydra is designed by a [security enthusiast](https://github.com/arekkas), who has written and participated in numerous auth* projects.

Additionally to the claims above, Hydra has received a lot of positive feedback. Let's see what the community is saying:

> Nice! Lowering barriers to the use of technologies like these is important.

[Pyxl101](https://news.ycombinator.com/item?id=11798641)

> OAuth is a framework not a protocol. The security it provides can vary greatly between implementations.
Fosite (which is what this is based on) is a very good implementation from a security perspective: https://github.com/ory-am/fosite#a-word-on-security

[abritishguy](https://news.ycombinator.com/item?id=11800515)

> [...] Thanks for releasing this by the way, looks really well engineered. [...]

[olalonde](https://news.ycombinator.com/item?id=11798831)

## Quickstart

This section is a quickstart guide to working with Hydra. In-depth docs are available as well:

* The documentation is available on [GitBook](https://ory-am.gitbooks.io/hydra/content/).
* The REST API documentation is available at [Apiary](http://docs.hdyra.apiary.io).

### Installation

#### Server

**Starting the host** is easiest with docker. The host process handles HTTP requests and is backed by a database.
Read how to install docker on [Linux](https://docs.docker.com/linux/), [OSX](https://docs.docker.com/mac/) or
[Windows](https://docs.docker.com/windows/). Hydra is available on [Docker Hub](https://hub.docker.com/r/oryam/hydra/).

##### Downloading an image from the Hub

Hydra will keep all changes in memory. But be aware: Restarting, scaling
or stopping the container will make you **lose all data**.

To download and start the image:
```
$ docker run -d -p 4444:4444 oryam/hydra --name my-hydra
ec91228cb105db315553499c81918258f52cee9636ea2a4821bdb8226872f54b
```

##### Build a docker image from the source

Install the [CLI and Docker Toolbox](#installation). Make sure you install Docker Compose. On OSX and Windows,
open the Docker Quickstart Terminal. On Linux, open any terminal.

We will use a dummy password as the system secret: `SYSTEM_SECRET=passwordtutorialpasswordtutorial`. Use a very secure secret in production.

**On OSX and Windows** using the Docker Quickstart Terminal:
```
$ go get github.com/ory-am/hydra
$ cd $GOPATH/src/github.com/ory-am/hydra
$ DOCKER_IP=$(docker-machine ip default) docker-compose build
WARNING: The SYSTEM_SECRET variable is not set. Defaulting to a blank string.
rethinkdb uses an image, skipping
Building hydra
[...]
$ SYSTEM_SECRET=passwordtutorial DOCKER_IP=$(docker-machine ip default) docker-compose up
Starting hydra_rethinkdb_1
Recreating hydra_hydra_1
Recreating hydra_consent_1
Attaching to hydra_rethinkdb_1, hydra_hydra_1, hydra_consent_1
[...]
```

**On Linux:**
```
$ go get github.com/ory-am/hydra
$ cd $GOPATH/src/github.com/ory-am/hydra
$ DOCKER_IP=localhost docker-compose build
WARNING: The SYSTEM_SECRET variable is not set. Defaulting to a blank string.
rethinkdb uses an image, skipping
Building hydra
[...]
$ SYSTEM_SECRET=passwordtutorial DOCKER_IP=localhost docker-compose up
Starting hydra_rethinkdb_1
Recreating hydra_hydra_1
Recreating hydra_consent_1
Attaching to hydra_rethinkdb_1, hydra_hydra_1, hydra_consent_1
[...]
mhydra   | mtime="2016-05-17T18:09:28Z" level=warning msg="Generated system secret: MnjFP5eLIr60h?hLI1h-!<4(TlWjAHX7"
[...]
mhydra   | mtime="2016-05-17T18:09:29Z" level=warning msg="Temporary root client created."
mhydra   | mtime="2016-05-17T18:09:29Z" level=warning msg="client_id: d9227bd5-5d47-4557-957d-2fd3bee11035"
mhydra   | mtime="2016-05-17T18:09:29Z" level=warning msg="client_secret: ,IvxGt02uNjv1ur9"
[...]
```
#### Client

##### Using gopack

**The CLI client is available at [gobuild.io](https://gobuild.io/ory-am/hydra)**.

There is currently no installer which adds the CLI to your path automatically. You have to set up the path yourself.
If you do not understand what that means, ask on our [Gitter channel](https://gitter.im/ory-am/hydra).

##### Building from the sources

If you wish to compile the CLI yourself, you need to install and set up [Go](https://golang.org/) and add `$GOPATH/bin`
to your `$PATH`. Here is a [comprehensive Go installation guide](https://github.com/ory-am/workshop-dbg#googles-go-language) with screenshots.

To install the CLI from source, execute:

```
go get github.com/ory-am/hydra
go get github.com/Masterminds/glide
cd $GOPATH/src/github.com/ory-am/hydra
glide install
go install github.com/ory-am/hydra
hydra
```
##### From the Docker container (not recommended):

```
$ docker exec -i -t <hydra-container-id> /bin/bash
# e.g. docker exec -i -t ec12 /bin/bash

root@ec91228cb105:/go/src/github.com/ory-am/hydra# hydra
Hydra is a twelve factor OAuth2 and OpenID Connect provider

Usage:
  hydra [command]

[...]
```

### Run the example

![Run the example](dist/run-the-example.gif)

You now have a running hydra docker container! Additionally, a RethinkDB image was deployed as well as a consent app.

Hydra can be managed with the hydra CLI client. The client has to log on before it is allowed to do anything.
When hydra detects a new installation, a new temporary root client is created. The client credentials are printed in 
the container logs.

```
mhydra   | mtime="2016-05-17T18:09:29Z" level=warning msg="client_id: d9227bd5-5d47-4557-957d-2fd3bee11035"
mhydra   | mtime="2016-05-17T18:09:29Z" level=warning msg="client_secret: ,IvxGt02uNjv1ur9"
```

The system secret is a global secret assigned to every hydra instance. It is used to encrypt data at rest. You can
set the system secret through the `$SYSTEM_SECRET` environment variable. When no secret is set, hydra generates one:

```
time="2016-05-15T14:56:34Z" level=warning msg="Generated system secret: (.UL_&77zy8/v9<sUsWLKxLwuld?.82B"
```

**Important note:** Please be aware that logging passwords should never be done on a production server. Either prune
the logs, set the required parameters, or replace the credentials with other ones.

Now you know which credentials you need to use. Next, we log in.

**Note:** If you are using docker toolbox, please use the IP address provided by `docker-machine ip default` as the cluster URL host.

```
$ hydra connect
Cluster URL: https://localhost:4444
Client ID: d9227bd5-5d47-4557-957d-2fd3bee11035
Client Secret: ,IvxGt02uNjv1ur9
Done.
```

Great! You are now connected to Hydra and can start by creating a new client:

```
$ hydra clients create --skip-tls-verify
Warning: Skipping TLS Certificate Verification.
Client ID: c003830f-a090-4721-9463-92424270ce91
Client Secret: Z2pJ0>Tp7.ggn>EE&rhnOzdt1
```

**Important note:** if no certificate is provided, Hydra uses self-signed TLS certificates for HTTPS. This should
never be done in production. To skip the TLS verification step on the client, provide the `--skip-tls-verify` flag.

Why not issue an access token for your client?

```
$ hydra token client --skip-tls-verify
Warning: Skipping TLS Certificate Verification.
JLbnRS9GQmzUBT4x7ESNw0kj2wc0ffbMwOv3QQZW4eI.qkP-IQXn6guoFew8TvaMFUD-SnAyT8GmWuqGi3wuWXg
```

Let's try this with the authorize code grant!

```
$ hydra token user --skip-tls-verify
Warning: Skipping TLS Certificate Verification.
If your browser does not open automatically, navigate to: https://192.168.99.100:4444/oauth2/auth?client_id=d9227bd5-5d47-4557-957d-2fd3bee11035&response_type=code&scope=core+hydra&state=sbnwdelqzxyedwtqinxnolbr&nonce=sffievieeesltbjkwxyhycyq
Setting up callback listener on http://localhost:4445/callback
Press ctrl + c on Linux / Windows or cmd + c on OSX to end the process.
```

![OAuth2 Flow](dist/oauth2-flow.gif)

Great! You installed hydra, connected the CLI, created a client and completed two authentication flows!
Your next stop should be the [Guide](#guide).

## Documentation

### Guide

The Guide is available on [GitBook](https://ory-am.gitbooks.io/hydra/content/).

### REST API Documentation

The REST API is documented at [Apiary](http://docs.hdyra.apiary.io).

### CLI Documentation

The CLI help is verbose. To see it, run `hydra -h` or `hydra [command] -h`.

### Develop

Unless you want to test Hydra against a database, developing with Hydra is as easy as:

```
go get github.com/ory-am/hydra
go get github.com/Masterminds/glide
cd $GOPATH/src/github.com/ory-am/hydra
glide install
go test ./...
go run main.go
```

If you want to run Hydra against RethinkDB, you can do so by using docker:

```
docker run --name some-rethink -d -p 8080:8080 -p 28015:28015 rethinkdb

# Linux
DATABASE_URL=rethinkdb://localhost:28015/hydra go run main.go

# Docker Terminal
DATABASE_URL=rethinkdb://$(docker-machine ip default):28015/hydra go run main.go
```

## FAQ

### I'm having trouble with the redirect URI.

Hydra enforces HTTPS for all hosts except localhost. Also make sure that the path is an exact match. `http://localhost:123/`
is not the same as `http://localhost:123`.

### How can I validate tokens?

Please use the Warden API. There is a go client library available [here](https://github.com/ory-am/hydra/blob/master/warden/warden_http.go).
The Warden API is documented [here](http://docs.hdyra.apiary.io/#reference/warden) and [here](https://ory-am.gitbooks.io/hydra/content/policy.html).

### How can I import TLS certificates?

You can import TLS certificates when running `hydra host`. This can be done by setting the following environment variables:

**Read from file**
- `HTTPS_TLS_CERT_PATH`: The path to the TLS certificate (pem encoded).
- `HTTPS_TLS_KEY_PATH`: The path to the TLS private key (pem encoded).

**Embedded**
- `HTTPS_TLS_CERT`: A pem encoded TLS certificate passed as string. Can be used instead of TLS_CERT_PATH.
- `HTTPS_TLS_KEY`: A pem encoded TLS key passed as string. Can be used instead of TLS_KEY_PATH.

Or by specifying the following flags:

```
--https-tls-cert-path string   Path to the certificate file for HTTP/2 over TLS (https). You can set HTTPS_TLS_KEY_PATH or HTTPS_TLS_KEY instead.
--https-tls-key-path string    Path to the key file for HTTP/2 over TLS (https). You can set HTTPS_TLS_KEY_PATH or HTTPS_TLS_KEY instead.
```

### I want to disable HTTPS for testing

You can do so by running `hydra host --force-dangerous-http`.

### Is there a client library / SDK?

Yes, for Go! It is available at `github.com/ory-am/hydra/sdk`.

Connect the SDK to Hydra:
```go
import "github.com/ory-am/hydra/sdk"

hydra, err := sdk.Connect(
    sdk.ClientID("client-id"),
    sdk.ClientSecret("client-secret"),
    sdk.ClustURL("https://localhost:4444"),
)
```

Manage OAuth Clients using [`ory-am/hydra/client.HTTPManager`](/client/manager_http.go):

```go
import "github.com/ory-am/hydra/client"

// Create a new OAuth2 client
newClient, err := hydra.Client.CreateClient(&client.Client{
    ID:                "deadbeef",
	Secret:            "sup3rs3cret",
	RedirectURIs:      []string{"http://yourapp/callback"},
    // ...
})

// Retrieve newly created client
newClient, err = hydra.Client.GetClient(newClient.ID)

// Remove the newly created client
err = hydra.Client.DeleteClient(newClient.ID)

// Retrieve list of all clients
clients, err := hydra.Client.GetClients()
```

Manage SSO Connections using [`ory-am/hydra/connection.HTTPManager`](connection/manager_http.go):
```go
import "github.com/ory-am/hydra/connection"

// Create a new connection
newSSOConn, err := hydra.SSO.Create(&connection.Connection{
    Provider: "login.google.com",
    LocalSubject: "bob",
    RemoteSubject: "googleSubjectID",
})

// Retrieve newly created connection
ssoConn, err := hydra.SSO.Get(newSSOConn.ID)

// Delete connection
ssoConn, err := hydra.SSO.Delete(newSSOConn.ID)

// Find a connection by subject
ssoConns, err := hydra.SSO.FindAllByLocalSubject("bob")
ssoConns, err := hydra.SSO.FindByRemoteSubject("login.google.com", "googleSubjectID")
```

Manage policies using [`ory-am/hydra/policy.HTTPManager`](policy/manager_http.go):
```go
import "github.com/ory-am/ladon"

// Create a new policy
// allow user to view his/her own photos
newPolicy, err := hydra.Policy.Create(ladon.DefaultPolicy{
    ID: "1234", // ID is not required
    Subjects: []string{"bob"},
    Resources: []string{"urn:media:images"},
    Actions: []string{"get", "find"},
    Effect: ladon.AllowAccess,
    Conditions: ladon.Conditions{
        "owner": &ladon.EqualSubjectCondition{},
    }
})

// Retrieve a stored policy
policy, err := hydra.Policy.Get("1234")

// Delete a policy
err := hydra.Policy.Delete("1234")

// Retrieve all policies for a subject
policies, err := hydra.Policy.FindPoliciesForSubject("bob")
```

Manage JSON Web Keys using [`ory-am/hydra/jwk.HTTPManager`](jwk/manager_http.go):

```go
// Generate new key set
keySet, err := hydra.JWK.CreateKeys("app-tls-keys", "HS256")

// Retrieve key set
keySet, err := hydra.JWK.GetKeySet("app-tls-keys")

// Delete key set
err := hydra.JWK.DeleteKeySet("app-tls-keys")
```

Validate requests with the Warden, uses [`ory-am/hydra/warden.HTTPWarden`](warden/warden_http.go):

```go
import "github.com/ory-am/ladon"

// Check if action is allowed
hydra.Warden.HTTPActionAllowed(ctx, req, &ladon.Request{
    Resource: "urn:media:images",
    Action: "get",
    Subject: "bob",
}, "media.images")

// Check if request is authorized
hydra.Warden.HTTPAuthorized(ctx, req, "media.images")
```

## Hall of Fame

A list of extraordinary contributors and [bug hunters](https://github.com/ory-am/hydra/issues/84).

* [Alexander Widerberg (leetal)](https://github.com/leetal) for implementing the prototype RethinkDB adapters.
