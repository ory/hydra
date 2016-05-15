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
- [Motivation](#motivation)
- [Quickstart](#quickstart)
  - [Installation](#installation)
  - [Run Hydra](#run-hydra)
- [Documentation](#documentation)
  - [Guide](#guide)
  - [REST API Documentation](#rest-api-documentation)
  - [CLI Documentation](#cli-documentation)
  - [Develop](#develop)
- [Frequently Asked Questions](#frequently-asked-questions)
  - [Deploy using buildpacks (Heroku, Cloud Foundry, ...)](#deploy-using-buildpacks-heroku-cloud-foundry-)
- [Attributions](#attributions)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## What is Hydra?

1. Hydra is an OAuth2 and OpenID Connect provider built for availability. The distributed in-memory architecture allows for heavy duty workloads.
2. Hydra works with every Identity Provider. The deprecated php-3.0 authentication service your intern wrote? It works with that too, don't worry.
3. Hydra does not use any templates, it is up to you what your front end should look like.
4. Hydra comes with two factor authentication, key management, social log on, policy management and access control.

## Motivation

At first, there was the monolith. The monolith worked well with the customized joomla authentication module. Then, the web evolved into an elastic cloud that serves thousands of different user agents in every part of the world. Hydra is driven by the need for an easy scalable, in memory OAuth2 and OpenID Connect provider, that integrates with every Identity Provider you can imagine. 

Hydra uses pub/sub to always have the latest data available in memory. Hydra scales effortlessly on every platform you can imagine, including Heroku, Cloud Foundry, Docker, Google Container Engine and many more.

## Quickstart

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
$ docker run -d -p 4444:4444 oryam/hydra --name my-hydra
ec91228cb105db315553499c81918258f52cee9636ea2a4821bdb8226872f54b
```

**CLI Client (Docker)**

If you are running docker locally, you can use the CLI by connecting to it:

```
$ docker exec -i -t <container> /bin/bash
# e.g. docker exec -i -t ec /bin/bash

root@ec91228cb105:/go/src/github.com/ory-am/hydra# hydra
Hydra is a twelve factor OAuth2 and OpenID Connect provider

Usage:
  hydra [command]

[...]
```

**CLI Client (Binary)**
The CLI client is available at [gobuild.io](https://gobuild.io/ory-am/hydra) or through
the [releases tab](https://github.com/ory-am/hydra/releases).

There is currently no installer which adds the client to your path automatically. You have to set up the path yourself.
If you do not understand what that means, ask on our [Gitter channel](https://gitter.im/ory-am/hydra).

If you wish to compile the CLI yourself, you need to install and set up [Go](https://golang.org/) and add `$GOPATH/bin`
to your `$PATH`. Here is a [comprehensive Go installation guide](https://github.com/ory-am/workshop-dbg#googles-go-language) with screenshots.

```
go install github.com/ory-am/hydra
hydra
```

### Run Hydra

Once you have [set up docker and installed the CLI](#installation) run:

```
$ docker run -d -p 4444:4444 oryam/hydra --name my-hydra
# You will receive a different container id.
# You can concatenate most of the id and use the two to three letters.
# In this case, that could be `ec9`.
ec91228cb105db315553499c81918258f52cee9636ea2a4821bdb8226872f54b
```

You have now a running hydra docker container! It is not backed by any database and runs completely in memory. Rebooting
or any other sort of disruption will purge all data.

There are two interesting flags used above:
* **-d** runs the docker in daemon mode.
* **-p** publishes port 4444.

**TBD:** Provision with RethinkDB.

Hydra can be managed with the hydra cli client. The client hast to log on before it is allowed to do anything.
When hydra detects a new installation, a new temporary root client is created. The client credentials will be available
from `docker logs`.

```
$ docker logs ec9
Pointing cluster at https://localhost:4444
time="2016-05-15T14:56:34Z" level=warning msg="No system secret specified."
time="2016-05-15T14:56:34Z" level=warning msg="Generated system secret: (.UL_&77zy8/v9<sUsWLKxLwuld?.82B"
time="2016-05-15T14:56:34Z" level=warning msg="Do not auto-generate system secrets in production."
time="2016-05-15T14:56:34Z" level=warning msg="Could not find OpenID Connect singing keys. Generating a new keypair..."
time="2016-05-15T14:56:34Z" level=warning msg="Keypair generated."
time="2016-05-15T14:56:34Z" level=warning msg="WARNING: Automated key creation causes low entropy. Replace the keys as soon as possible."
time="2016-05-15T14:56:34Z" level=warning msg="No clients were found. Creating a temporary root client..."
time="2016-05-15T14:56:34Z" level=warning msg="Temporary root client created."
time="2016-05-15T14:56:34Z" level=warning msg="client_id: ad586b43-eb85-433c-8e46-8264bf0407b3"
time="2016-05-15T14:56:34Z" level=warning msg="client_secret: -,ak$P_qLjijKa,5"
time="2016-05-15T14:56:34Z" level=warning msg="The root client must be removed in production. The root's credentials could be accidentally logged."
time="2016-05-15T14:56:34Z" level=warning msg="Key for TLS not found. Creating new one."
time="2016-05-15T14:56:34Z" level=warning msg="Temporary key created."
time="2016-05-15T14:56:34Z" level=info msg="Starting server on :4444"
```

As you can see, various keys are being generated, when hydra is started against an empty database.

The system secret is a global secret assigned to every hydra instance. It is used to encrypt data at rest. You can
set the system secret through the `$SYSTEM_SECRET` environment variable. When no secret is set, hydra generates one:

```
time="2016-05-15T14:56:34Z" level=warning msg="Generated system secret: (.UL_&77zy8/v9<sUsWLKxLwuld?.82B"
```

Our temporary root client was generated as well:

```
time="2016-05-15T14:56:34Z" level=warning msg="client_id: ad586b43-eb85-433c-8e46-8264bf0407b3"
time="2016-05-15T14:56:34Z" level=warning msg="client_secret: -,ak$P_qLjijKa,5"
```

**Important note:** Please be aware that logging passwords should never be done on a production server. Either prune
the logs, set the required parameters, or replace the credentials with other ones.

Now you know which credentials you need to use. Next, we log in.

**Note:** If you are using docker toolbox, please use the ip address provided by `docker-machine ip default` as cluster url host.

```
$ hydra connect
Cluster URL: https://localhost:4444
Client ID: ad586b43-eb85-433c-8e46-8264bf0407b3
Client Secret: -,ak$P_qLjijKa,5
Done.
```

Great! You are now connected to Hydra and can start by creating a new client:

```
$ hydra clients create --skip-tls-verify
Warning: Skipping TLS Certificate Verification.
Client ID: c003830f-a090-4721-9463-92424270ce91
Client Secret: Z2pJ0>Tp7.ggn>EE&rhnOzdt1
```

**Important note:** Hydra is using self signed TLS certificates for HTTPS, if no certificate was provided. This should
never be done in production. To skip the TLS verification step on the client, provide the `--skip-tls-verify` flag.

Great! You installed hydra, connected the CLI and created a client. Your next stop should be the [Documentation](#documentation).

## Documentation

### Guide

The Guide is available on [GitBook](https://ory-am.gitbooks.io/hydra/content/).

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
