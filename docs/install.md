# Installing, Configuring and Running Hydra

Before starting with this section, please check out the [tutorial](./demo.md). It will teach you the most important flows
and settings for Hydra.

## Installing Hydra

You can install Hydra using multiple methods.

### Using Docker

Installing, configuring and running Hydra is easiest with docker. The host process
handles HTTP requests and is backed by a database.
Read how to install docker on [Linux](https://docs.docker.com/linux/), [OSX](https://docs.docker.com/mac/) or
[Windows](https://docs.docker.com/windows/). Hydra is available on [Docker Hub](https://hub.docker.com/r/oryam/hydra/).

In this minimalistic example, we will use Hydra without a database. Bee aware that restarting, scaling
or stopping the container will **lose all the data**.

```
$ docker run -d -p 4444:4444 oryam/hydra --name my-hydra
ec91228cb105db315553499c81918258f52cee9636ea2a4821bdb8226872f54b
```

Now, you should be able to open [https://localhost:4444](https://localhost:4444). If asked, accept the self signed
certificate in your browser.

**Using the client command line interface** can be achieved by ssh'ing into the hydra container
and execute the hydra command from there:

```
$ docker exec -i -t <container-id> /bin/bash
# e.g. docker exec -i -t ec91228 /bin/bash

root@ec91228cb105:/go/src/github.com/ory-am/hydra# hydra
Hydra is a twelve factor OAuth2 and OpenID Connect provider

[...]
```

### Download Binaries

The client and server **binaries are downloadable at the [releases tab](https://github.com/ory-am/hydra/releases)**.
There is currently no installer available. You have to add the hydra binary to the PATH environment variable yourself or put
the binary in a location that is already in your path (`/usr/bin`, ...). 
If you do not understand what that all of this means, ask in our [chat channel](https://gitter.im/ory-am/hydra). We are happy to help.

Once installed, you should be able to run:

```
$ hydra help
Hydra is a cloud native high throughput OAuth2 and OpenID Connect provider

Usage:
  hydra [command]

Available Commands:
  clients     Manage OAuth2 clients
...
```

### Build from source

If you wish to compile hydra yourself, you need to install and set up [Go 1.5+](https://golang.org/) and add `$GOPATH/bin`
to your `$PATH`. To do so, run the following commands in a shell (bash, sh, cmd.exe, ...):

```
$ go get github.com/ory-am/hydra
$ go get github.com/Masterminds/glide
$ cd $GOPATH/src/github.com/ory-am/hydra
$ glide install
$ go install github.com/ory-am/hydra
$ hydra
Hydra is a cloud native high throughput OAuth2 and OpenID Connect provider

Usage:
  hydra [command]

Available Commands:
  clients     Manage OAuth2 clients
...
```

## Configuring and Running Hydra

Running the default Hydra environment is as easy as:
 
```
$ hydra host
time="2016-10-13T10:04:01+02:00" level=info msg="DATABASE_URL not set, connecting to ephermal in-memory database."
time="2016-10-13T10:04:01+02:00" level=warning msg="Expected system secret to be at least 32 characters long, got 0 characters."
time="2016-10-13T10:04:01+02:00" level=info msg="Generating a random system secret..."
time="2016-10-13T10:04:01+02:00" level=info msg="Generated system secret: gbaMR?mYdTB-g/$YhRXjmX!,h08.t07/"
time="2016-10-13T10:04:01+02:00" level=warning msg="WARNING: DO NOT generate system secrets in production. The secret will be leaked to the logs."
...
```

Hydra relies on a third party for storing data, such as Postgres or MySQL (officially supported) and RethinkDB
(community supported). If no storage is set, data will be written to memory and is lost when the process is killed.
The `hydra help host` command will give you an insight into the different configuration settings. The following section
might be outdated and is only for demonstration purposes, please run the `hydra help host` command on your local
machine to get the latest documentation:

```
$ hydra help host
Starts all HTTP/2 APIs and connects to a database backend.

This command exposes a variety of controls via environment variables. You can
set environments using "export KEY=VALUE" (Linux/macOS) or "set KEY=VALUE" (Windows). On Linux,
you can also set environments by prepending key value pairs: "KEY=VALUE KEY2=VALUE2 hydra"

All possible controls are listed below. The host process additionally exposes a few flags, which are listed below
the controls section.

CORE CONTROLS
=============

- DATABASE_URL: A URL to a persistent backend. Hydra supports various backends:
  - None: If DATABASE_URL is empty, all data will be lost when the command is killed.
  - RethinkDB: If DATABASE_URL is a DSN starting with rethinkdb://, RethinkDB will be used as storage backend.
        Example: DATABASE_URL=rethinkdb://user:password@host:123/database

        Additionally, these controls are available when using RethinkDB:
        - RETHINK_TLS_CERT_PATH: The path to the TLS certificate (pem encoded) used to connect to rethinkdb.
                Example: RETHINK_TLS_CERT_PATH=~/rethink.pem

        - RETHINK_TLS_CERT: A pem encoded TLS certificate passed as string. Can be used instead of RETHINK_TLS_CERT_PATH.
                Example: RETHINK_TLS_CERT_PATH="-----BEGIN CERTIFICATE-----\nMIIDZTCCAk2gAwIBAgIEV5xOtDANBgkqhkiG9w0BAQ0FADA0MTIwMAYDVQQDDClP..."

- SYSTEM_SECRET: A secret that is at least 16 characters long. If none is provided, one will be generated. They key
        is used to encrypt sensitive data using AES-GCM (256 bit) and validate HMAC signatures.
        Example: SYSTEM_SECRET=jf89-jgklAS9gk3rkAF90dfsk

- FORCE_ROOT_CLIENT_CREDENTIALS: On first start up, Hydra generates a root client with random id and secret. Use
        this environment variable in the form of "FORCE_ROOT_CLIENT_CREDENTIALS=id:secret" to set
        the client id and secret yourself.
        Example: FORCE_ROOT_CLIENT_CREDENTIALS=admin:kf0AKfm12fas3F-.f

- PORT: The port hydra should listen on.
        Defaults to PORT=4444

- HOST: The port hydra should listen on.
        Example: PORT=localhost

- BCRYPT_COST: Set the bcrypt hashing cost. This is a trade off between
        security and performance. Range is 4 =< x =< 31.
        Defaults to BCRYPT_COST=10


OAUTH2 CONTROLS
===============

- CONSENT_URL: The uri of the consent endpoint.
        Example: CONSENT_URL=https://id.myapp.com/consent

- ISSUER: The issuer is used for identification in all OAuth2 tokens.
        Defaults to ISSUER=hydra.localhost

- AUTH_CODE_LIFESPAN: Lifespan of OAuth2 authorize codes. Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
        Defaults to AUTH_CODE_LIFESPAN=10m

- ID_TOKEN_LIFESPAN: Lifespan of OpenID Connect ID Tokens. Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
        Defaults to ID_TOKEN_LIFESPAN=1h

- ACCESS_TOKEN_LIFESPAN: Lifespan of OAuth2 access tokens. Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
        Defaults to ACCESS_TOKEN_LIFESPAN=1h

- CHALLENGE_TOKEN_LIFESPAN: Lifespan of OAuth2 consent tokens. Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
        Defaults to CHALLENGE_TOKEN_LIFESPAN=10m


HTTPS CONTROLS
==============

- HTTPS_ALLOW_TERMINATION_FROM: Whitelist one or multiple CIDR address ranges and allow them to terminate TLS connections.
        Be aware that the X-Forwarded-Proto header must be set and must never be modifiable by anyone but
        your proxy / gateway / load balancer. Supports ipv4 and ipv6.
        Hydra serves http instead of https when this option is set.
        Example: HTTPS_ALLOW_TERMINATION_FROM=127.0.0.1/32,192.168.178.0/24,2620:0:2d0:200::7/32

- HTTPS_TLS_CERT_PATH: The path to the TLS certificate (pem encoded).
        Example: HTTPS_TLS_CERT_PATH=~/cert.pem

- HTTPS_TLS_KEY_PATH: The path to the TLS private key (pem encoded).
        Example: HTTPS_TLS_KEY_PATH=~/key.pem

- HTTPS_TLS_CERT: A pem encoded TLS certificate passed as string. Can be used instead of HTTPS_TLS_CERT_PATH.
        Example: HTTPS_TLS_CERT="-----BEGIN CERTIFICATE-----\nMIIDZTCCAk2gAwIBAgIEV5xOtDANBgkqhkiG9w0BAQ0FADA0MTIwMAYDVQQDDClP..."

- HTTPS_TLS_KEY: A pem encoded TLS key passed as string. Can be used instead of HTTPS_TLS_KEY_PATH.
        Example: HTTPS_TLS_KEY="-----BEGIN ENCRYPTED PRIVATE KEY-----\nMIIFDjBABgkqhkiG9w0BBQ0wMzAbBgkqhkiG9w0BBQwwDg..."


DEBUG CONTROLS
==============

- PROFILING: Set "PROFILING=cpu" to enable cpu profiling and "PROFILING=memory" to enable memory profiling.
        It is not possible to do both at the same time.
        Example: PROFILING=cpu

Usage:
  hydra host [flags]

Flags:
      --dangerous-auto-logon           Stores the root credentials in ~/.hydra.yml. Do not use in production.
      --dangerous-force-http           Disable HTTP/2 over TLS (HTTPS) and serve HTTP instead. Never use this in production.
      --https-tls-cert-path string     Path to the certificate file for HTTP/2 over TLS (https). You can set HTTPS_TLS_CERT_PATH or HTTPS_TLS_CERT instead.
      --https-tls-key-path string      Path to the key file for HTTP/2 over TLS (https). You can set HTTPS_TLS_KEY_PATH or HTTPS_TLS_KEY instead.
      --rethink-tls-cert-path string   Path to the certificate file to connect to rethinkdb over TLS (https). You can set RETHINK_TLS_CERT_PATH or RETHINK_TLS_CERT instead.

Global Flags:
      --config string     config file (default is $HOME/.hydra.yaml)
      --skip-tls-verify   foolishly accept TLS certificates signed by unkown certificate authorities
```

It is quite common to run hydra with the following options:

```
$ export DATABASE_URL=postgres://foo:bar@localhost/hydra
$ export SYSTEM_SECRET=some-very-random-secret-$§123
$ hydra host
```

If you want to check out hydra locally, we recommend setting these options to ease things up. ` --dangerous-auto-logon`
will write the administrator's credentials directly to `~/.hydra.yml`, no `hydra connect` required. The section option
`--dangerous-force-http` disables https and serves Hydra over http instead:

```
$ hydra host --dangerous-auto-logon --dangerous-force-http
```
