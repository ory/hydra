# Install, Configure and Run ORY Hydra

The goal of this section is to familiarize you with the specifics of setting up ORY Hydra in your environment.
Before starting with this section, please check out the [tutorial](./tutorial.md). It will teach you the most important flows
and settings for Hydra.

This guide will:

1. Download and run a PostgreSQL container in docker.
2. Download and run ORY Hydra in docker.
3. Create OAuth 2.0 Client for the consent app, and apply necessary policies.
4. Download and run an exemplary consent app ([node](https://github.com/ory/hydra-consent-app-express), [golang](hydra-consent-app-go)).
5. Create an OAuth 2.0 consumer app.
6. Download and run an OAuth 2.0 consumer using docker.

Before starting with this guide, please install the most recent version of [Docker](https://www.docker.com/community-edition#/download).

## Start a ostgreSQL container

For the purpose of this tutorial, we will use a PostgreSQL container. Never run databases in Docker in production!

```
$ docker run \
  --name ory-hydra-example--postgres \
  -e POSTGRES_USER=hydra \
  -e POSTGRES_PASSWORD=secret \
  -e POSTGRES_DB=hydra \
  -d postgres:9.6
```

This command wil start a postgres instance with name `ory-hydra-example--postgres`, set up a database called `hydra`
and create a user `hydra` with password `secret`.

## Install and run ORY Hydra

We highly recommend using Docker to run Hydra, as installing, configuring and running Hydra is easiest with docker. The host process
handles HTTP requests and is backed by a database. ORY Hydra is available on [Docker Hub](https://hub.docker.com/r/oryd/hydra/).

```
# The system secret can only be set against a fresh database. Key rotation is currently not supported. This
# secret is used to encrypt the database and needs to be set to the same value every time.
$ export SYSTEM_SECRET=this_needs_to_be_the_same_always

# The database url points us at the postgres instance. This could also be an ephermal in-memory database (`export DATABASE_URL=memory`)
# or a MySQL url.
$ export DATABASE_URL=postgres://hydra:secret@postgres:5432/hydra?sslmode=disable

# This command will show you all the environment variables that you can set. Read this carefully.
# It is the equivalent to `hydra help host`.
$ docker run -it --entrypoint hydra oryd/hydra:latest help host

Starts all HTTP/2 APIs and connects to a database backend.
[...]

# ORY Hydra does not do magic, it requires concious decisions, for example, when running SQL migrations, which is required
# when installing a new version of ORY Hydra, or upgrading an existing installation.
# It is the equivalent to `hydra migrate sql postgres://hydra:secret@postgres:5432/hydra?sslmode=disable`
$ docker run --link ory-hydra-example--postgres:postgres -it --entrypoint hydra oryd/hydra:latest migrate sql $DATABASE_URL

Applying `ladon` SQL migrations...
Applied 3 `ladon` SQL migrations.
Applying `client` SQL migrations...
[...]
Migration successful!

# Let's run our docker server (settings explained below):
$ docker run -d \
  --name ory-hydra-example--hydra \
  --link ory-hydra-example--postgres:postgres \
  -p 9000:4444 \
  -e SYSTEM_SECRET=$SYSTEM_SECRET \
  -e DATABASE_URL=$DATABASE_URL \
  -e ISSUER=https://localhost:9000/ \
  -e CONSENT_URL=http://localhost:9020/consent \
  -e FORCE_ROOT_CLIENT_CREDENTIALS=admin:demo-password \
  oryd/hydra:latest

# And check if it's running:
$ docker logs ory-hydra-example--hydra

time="2017-06-29T21:26:26Z" level=info msg="Connecting with postgres://*:*@postgres:5432/hydra?sslmode=disable"
time="2017-06-29T21:26:26Z" level=info msg="Connected to SQL!"
[...]
time="2017-06-29T21:26:34Z" level=info msg="Setting up http server on :4444"
```

Let's dive into the various settings:

* `--link ory-hydra-example--postgres:postgres` connects this instance to postgres. Attention, this feature will
be deprecated in docker in the future.
* `-p 9000:4444` exposes ORY Hydra on `https://localhost:9000/`.
* `-e SYSTEM_SECRET=$SYSTEM_SECRET` sets the system secret environment variable.
* `-e DATABASE_URL=$DATABASE_URL` sets the database url environment variable.
* `-e ISSUER=https://localhost:9000/` set issuer to the publicly accessible url.
* `-e CONSENT_URL=http://localhost:9020/consent` set the url of the consent app to this one. We will set up the consent
app in the following sections.
* `-e FORCE_ROOT_CLIENT_CREDENTIALS=admin:demo-password` sets the credentials of the root account. Use the root
account to manage your ORY Hydra instance. If this is not set, ORY Hydra will auto-generate a client and display
the credentials in the logs.

To confirm that the instance is running properly, [open the health check](https://localhost:4444/health). If asked,
accept the self signed certificate in your browser. You should simply see `ok`.

### Install ORY Hydra without Docker

You can also install ORY Hydra without docker. For the purpose of this tutorial, please skip this section for now, and read
it later.

#### Download binaries

The client and server **binaries are downloadable at the [releases tab](https://github.com/ory/hydra/releases)**.
There is currently no installer available. You have to add the Hydra binary to the PATH environment variable yourself or put
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

#### Build from source

If you wish to compile ORY Hydra yourself, you need to install and set up [Go 1.8+](https://golang.org/) and add `$GOPATH/bin`
to your `$PATH`. To do so, run the following commands in a shell (bash, sh, cmd.exe, ...):

```
$ go get -d -u github.com/ory/hydra
$ go get github.com/Masterminds/glide
$ cd $GOPATH/src/github.com/ory/hydra
$ glide install
$ go install github.com/ory/hydra
$ hydra

Hydra is a cloud native high throughput OAuth2 and OpenID Connect provider

Usage:
  hydra [command]

Available Commands:
  clients     Manage OAuth2 clients
...
```

## Configure ORY Hydra

Next we need to connect to ORY Hydra and set up OAuth 2.0 Clients and Access Control Policies. In particular:

1. We create an OAuth 2.0 Client to use for the Consent App `consent-app`.
2. We set up access control policies for `consent-app`.
3. We add a policy that allows anybody to access the public keys for validating OpenID Connect ID Tokens.
3. We create another client capable of performing the client credentials grant and the authorize code grant.

```
# We run a shell with ORY Hydra installed. We also expose port 4445 which we will use later to perform
# the authorize code flow. Also, we connect this container to our ORY Hydra instance.
$ docker run -p 9010:4445 --link ory-hydra-example--hydra:hydra -it --entrypoint "/bin/sh" oryd/hydra:latest

# We connect to the ORY Hydra cluster
$ hydra connect

Cluster URL []: https://hydra:4444
Client ID []: admin
Client Secret [empty]: demo-password
Persisting config in file /root/.hydra.yml

# And issue an access token to validate that everything is working (ps: we need to disable TLS verification
# because the TLS certificate is self-signed).
$ hydra token client --skip-tls-verify
tY9tGakiYAUn8VIGn_yCDlTahckSfGbDQIlXahjXtX0.BQlCxRDL3ngag6hdsSl9N2qrz7R399cQMfld8aI2Mlg

# We can also try to validate an token:
$ hydra token validate --skip-tls-verify $(hydra token client --skip-tls-verify)

{
        "active": true,
        "scope": "hydra",
        "client_id": "admin",
        "sub": "admin",
        "exp": 1498775665,
        "iat": 1498772065,
        "aud": "admin",
        "iss": "https://localhost:9000"
}
```

### Consent app set up

The consent app is a bridge between ORY Hydra and your authentication services. The consent app requires
a registered OAuth 2.0 Client at ORY Hydra with the following properties. Do not re-use this client for any
other purposes:

```
# First we need to create the client
$ hydra clients create --skip-tls-verify \
  --id consent-app \
  --secret consent-secret \
  --name "Consent App Client" \
  --grant-types client_credentials \
  --response-types token \
  --allowed-scopes hydra.keys.get
```

Let's dive into the arguments:

* `--id consent-app` is the id of the client.
* `--secret consent-secret` sets the secret of the client. If no secret is provided, one will be generated.
The secret is visible *only once and can not be retrieved again*.
* `--name "Consent App Client"` a human-readable name.
* `--grant-types client_credentials` this client performs the OAuth 2.0 Client Credentials grant only.
* `--response-types token` this client only requires access tokens, no authorize codes or refresh tokens.
* `--scope hydra.keys.get` in order to access cryptographic keys, the client needs the keys scope.

Cool, next we need to create a policy for this client as well. The policy allows `consent-app` to access the
required cryptographic keys for validating and signing the consent challenge and response:

```
# For more information on access control policies, please read
# https://ory.gitbooks.io/hydra/content/security.html#access-control-policies
$ hydra policies create --skip-tls-verify \
  --actions get \
  --description "Allow consent-app to access the cryptographic keys for signing and validating the consent challenge and response" \
  --allow \
  --id consent-app-policy \
  --resources rn:hydra:keys:hydra.consent.challenge:public,rn:hydra:keys:hydra.consent.response:private \
  --subjects consent-app

Created policy consent-app-policy.
```

Let's take a look at the arguments:

* `--actions get` we need to access the keys
* `--allow` sets the policy effect to `allow`. Omit to set this for `deny`.
* `--id consent-app-policy` a unique identifier.
* `--resources rn:hydra:keys:hydra.consent.challenge:public,rn:hydra:keys:hydra.consent.response:private` an array
of comma-separated resource names. These two are fixed in ORY Hydra.
* `--subjects consent-app` the subject ("user") of this policy is our consent app.

Awesome! Next we will run the [ORY Hydra Consent App Example (NodeJS)](https://github.com/ory/hydra-consent-app-express).
This app is also available in [Golang](https://github.com/ory/hydra-consent-app-go), but for the purpose of this
tutorial we will use the NodeJS one. In a new shell, run:

```
$ docker run -d \
  --name ory-hydra-example--consent \
  --link ory-hydra-example--hydra:hydra \
  -p 9020:3000 \
  -e HYDRA_CLIENT_ID=consent-app \
  -e HYDRA_CLIENT_SECRET=consent-secret \
  -e HYDRA_URL=https://hydra:4444 \
  -e NODE_TLS_REJECT_UNAUTHORIZED=0 \
  oryd/hydra-consent-app-express:latest

# Let's check if it's running ok:
$ docker logs ory-hydra-example--consent
```

Let's take a look at the arguments:
* `-p 9020:3000` exposes this service at port 9010. If you remember, that's the port of the `CONSENT_URL` value
from the ORY Hydra docker container (`CONSENT_URL=http://localhost:9020/consent`).
* `-e HYDRA_CLIENT_ID=consent-app` this is the client id we created in the steps above.
* `-e HYDRA_CLIENT_SECRET=consent-secret` this is the client secret we set in the steps above.
* `HYDRA_URL=http://hydra:4444` point to the ORY Hydra container.
* `NODE_TLS_REJECT_UNAUTHORIZED=0` disables TLS verification, because we are using self-signed certificates.

Coming back to the ORY Hydra shell (`docker run -p 9010:4445 --link ory-hydra-example--hydra:hydra -it --entrypoint "/bin/sh" oryd/hydra:latest`),
we can no

## Perform OAuth 2.0 Flow

Our infrastructure is now set up. Now it's time to set up an OAuth 2.0 Consumer and perform the OAuth 2.0 Authorize Code flow.
To do so, we will create a new client and also create an access control policy that allows everybody to read the
public key of the OpenID Connect ID Token. Please use the SSH container from above to perform these commands:

```
$ hydra clients create --skip-tls-verify \
  --id some-consumer \
  --secret consumer-secret \
  --grant-types authorize_code,refresh_token,client_credentials,implicit \
  --response-types token,code,id_token \
  --allowed-scopes openid,offline,hydra.clients \
  --callbacks http://localhost:9010/callback

Client ID: some-consumer
Client Secret: consumer-secret
```

Let's dive into some of the arguments:
* `--grant-types authorize_code,refresh_token,client_credentials,implicit` we want to be able to perform all of these
OAuth 2.0 flows.
* `--response-types token,code,id_token` allows us to receive authorize codes, access and refresh tokens, and
OpenID Connect ID Tokens.
* `--allowed-scopes hydra.clients` allows this client to request scope `hydra.clients` and all scopes prefixed with `hydra.clients.`, for example `hydra.clients.get`.

Also, we want to allow everyone (not only our consumer) access to the public key of the OpenID Connect ID Token, which can be achieved with:

```
$ hydra policies create --skip-tls-verify \
  --actions get \
  --description "Allow everyone to read the OpenID Connect ID Token public key" \
  --allow \
  --id openid-id_token-policy \
  --resources rn:hydra:keys:hydra.openid.id-token:public \
  --subjects "<.*>"

Created policy openid-id_token-policy.
```

Perfect, let's perform an exemplary OAuth 2.0 Authorize Code Flow! To make this easy, the ORY Hydra CLI provides
a helper command called `hydra token user`. Just imagine this being, for example, passport.js that is generating
an auth code url, redirecting the browser to it, and then exchanging the authorize code for an access token. The
same thing happens with this command:

```
$ hydra token user --skip-tls-verify \
  --auth-url https://localhost:9000/oauth2/auth \
  --token-url https://hydra:4444/oauth2/token \
  --id some-consumer \
  --secret consumer-secret \
  --scopes openid,offline,hydra.clients \
  --redirect http://localhost:9010/callback

Setting up callback listener on http://localhost:4445/callback
Press ctrl + c on Linux / Windows or cmd + c on OSX to end the process.
If your browser does not open automatically, navigate to:

        https://localhost:9000/oauth2/auth?client_id=some-consumer&redirect_uri=http%3A%2F%2Flocalhost%3A9020%2Fcallback&response_type=code&scope=openid+offline+hydra.clients&state=hfcyxoqoctwbnvrxrsuwgzfu&nonce=lbeouolavuvcdhjefcnzlqur
```

open the link, as prompted, in your browser, and follow the steps shown there.

## Old guide

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

1. If no system secret was given, a random one is generated. Note this down, otherwise you won't be able to restart Hydra.
2. Cryptographic keys for JWT signing are being generated.
3. If the OAuth 2.0 Client database table is empty, a new root client with random credentials is created. Root clients
have access to all APIs, OAuth 2.0 flows and are allowed to do everything. If the `FORCE_ROOT_CLIENT_CREDENTIALS` environment.
is set, those credentials will be used instead.
4. A self signed certificate for serving HTTP over TLS is created.

Hydra can be managed using the Hydra Command Line Interface (CLI) client. This client has to log on before it is
allowed to do anything. When Hydra host process detects a new installation, a new temporary root client is
created and its credentials are printed to the container logs.

```
mhydra   | mtime="2016-05-17T18:09:29Z" level=warning msg="client_id: d9227bd5-5d47-4557-957d-2fd3bee11035"
mhydra   | mtime="2016-05-17T18:09:29Z" level=warning msg="client_secret: ,IvxGt02uNjv1ur9"
```

The system secret is a global secret assigned to every Hydra instance. It is used to encrypt data at rest. You can
set the system secret through the `SYSTEM_SECRET` environment variable. When no secret is set, Hydra generates one:

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

Great! You installed Hydra, connected the CLI, created a client and completed two authentication flows!
