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

## Configuring Hydra

Running the default Hydra environment is as easy as:
 
```
$ hydra host
time="2016-10-13T10:04:01+02:00" level=info msg="DATABASE_URL not set, connecting to ephermal in-memory database."
time="2016-10-13T10:04:01+02:00" level=warning msg="Expected system secret to be at least 32 characters long, got 0 characters."
time="2016-10-13T10:04:01+02:00" level=info msg="Generating a random system secret..."
[...]
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
[...]
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

## Running Hydra

On first run, Hydra initializes various settings:

```
$ hydra host
[...]
mtime="2016-05-17T18:09:28Z" level=warning msg="Generated system secret: MnjFP5eLIr60h?hLI1h-!<4(TlWjAHX7"
[...]
time="2016-10-25T09:58:54+02:00" level=info msg="Key pair for signing hydra.openid.id-token is missing. Creating new one."
time="2016-10-25T09:58:56+02:00" level=info msg="Key pair for signing hydra.consent.response is missing. Creating new one."
time="2016-10-25T09:59:02+02:00" level=info msg="Key pair for signing hydra.consent.challenge is missing. Creating new one."
[...]
time="2016-10-25T09:59:04+02:00" level=warning msg="No clients were found. Creating a temporary root client..."
mtime="2016-05-17T18:09:29Z" level=warning msg="client_id: d9227bd5-5d47-4557-957d-2fd3bee11035"
mtime="2016-05-17T18:09:29Z" level=warning msg="client_secret: ,IvxGt02uNjv1ur9"
[...]
time="2016-10-25T09:59:04+02:00" level=warning msg="No TLS Key / Certificate for HTTPS found. Generating self-signed certificate."
```

1. If no system secret was given, a random one is generated
2. Cryptographic keys for JWT signing are being generated
3. If the OAuth 2.0 Client database table is empty, a new root client with random credentials is created. Root clients
have access to all APIs, OAuth 2.0 flows and are allowed to do everything. If the `FORCE_ROOT_CLIENT_CREDENTIALS` environment
is set, those credentials will be used instead.
4. A self signed certificate for serving HTTP over TLS is created.

Hydra can be managed using the Hydra Command Line Interface (CLI) client. This client has to log on before it is
allowed to do anything. When Hydra host process detects a new installation, a new temporary root client is
created and its credentials are printed to the container logs.

```
mhydra   | mtime="2016-05-17T18:09:29Z" level=warning msg="client_id: d9227bd5-5d47-4557-957d-2fd3bee11035"
mhydra   | mtime="2016-05-17T18:09:29Z" level=warning msg="client_secret: ,IvxGt02uNjv1ur9"
```

The system secret is a global secret assigned to every hydra instance. It is used to encrypt data at rest. You can
set the system secret through the `SYSTEM_SECRET` environment variable. When no secret is set, hydra generates one:

```
time="2016-05-15T14:56:34Z" level=warning msg="Generated system secret: (.UL_&77zy8/v9<sUsWLKxLwuld?.82B"
```

If you are using the Hydra CLI locally or on a different host, you need to use the credentials from above to log in.

```
$ hydra connect
Cluster URL: https://localhost:4444
Client ID: d9227bd5-5d47-4557-957d-2fd3bee11035
Client Secret: ,IvxGt02uNjv1ur9
Done.
```

Great! You are now connected to Hydra and can start by creating a new client:

```
$ hydra clients create
Client ID: c003830f-a090-4721-9463-92424270ce91
Client Secret: Z2pJ0>Tp7.ggn>EE&rhnOzdt1
```

Now, let us issue an access token for your OAuth2 client!

```
$ hydra token client
JLbnRS9GQmzUBT4x7ESNw0kj2wc0ffbMwOv3QQZW4eI.qkP-IQXn6guoFew8TvaMFUD-SnAyT8GmWuqGi3wuWXg
```

Great! You installed hydra, connected the CLI, created a client and completed two authentication flows!
