# Install, Configure and Run ORY Hydra

The goal of this section is to familiarize you with the specifics of setting up ORY Hydra in your environment.
Before starting with this section, please check out the [tutorial](../0-Tutorial/0-README.md). It will teach you the most important flows
and settings for Hydra.

This guide will:

1. Download and run a PostgreSQL container in Docker.
2. Download and run ORY Hydra in Docker.
3. Create OAuth 2.0 Client for the consent app, and apply necessary policies.
4. Download and run an exemplary consent app ([node](https://github.com/ory/hydra-consent-app-express), [golang](hydra-consent-app-go)).
5. Create an OAuth 2.0 consumer app.
6. Download and run an OAuth 2.0 consumer using docker.

Before starting with this guide, please install the most recent version of [Docker](https://www.docker.com/community-edition#/download).

## Create a network

```
$ docker network create hydraguide
```

## Start a PostgreSQL container

For the purpose of this tutorial, we will use a PostgreSQL container. Never run databases in Docker in production!

```
$ docker run \
  --network hydraguide \
  --name ory-hydra-example--postgres \
  -e POSTGRES_USER=hydra \
  -e POSTGRES_PASSWORD=secret \
  -e POSTGRES_DB=hydra \
  -d postgres:9.6
```

This command wil start a postgres instance with name `ory-hydra-example--postgres`, set up a database called `hydra`
and create a user `hydra` with password `secret`.

## Install and run ORY Hydra

We highly recommend using Docker to run Hydra, as installing, configuring and running Hydra is easiest with Docker.
ORY Hydra is available on [Docker Hub](https://hub.docker.com/r/oryd/hydra/).

```
# The system secret can only be set against a fresh database. Key rotation is currently not supported. This
# secret is used to encrypt the database and needs to be set to the same value every time the process (re-)starts.
$ export SYSTEM_SECRET=this_needs_to_be_the_same_always_and_also_very_$3cuR3-._

# The database url points us at the postgres instance. This could also be an ephermal in-memory database (`export DATABASE_URL=memory`)
# or a MySQL url.
$ export DATABASE_URL=postgres://hydra:secret@ory-hydra-example--postgres:5432/hydra?sslmode=disable

# Before starting, let's pull the latest ORY Hydra tag from docker.
$ docker pull oryd/hydra:v0.11.6

# This command will show you all the environment variables that you can set. Read this carefully.
# It is the equivalent to `hydra help host`.
$ docker run -it --rm --entrypoint hydra oryd/hydra:v0.11.6 help host

Starts all HTTP/2 APIs and connects to a database backend.
[...]

# ORY Hydra does not do magic, it requires concious decisions, for example, when running SQL migrations, which is required
# when installing a new version of ORY Hydra, or upgrading an existing installation.
# It is the equivalent to `hydra migrate sql postgres://hydra:secret@ory-hydra-example--postgres:5432/hydra?sslmode=disable`
$ docker run -it --rm \
  --network hydraguide \
  oryd/hydra:v0.11.6 \
  migrate sql $DATABASE_URL

Applying `ladon` SQL migrations...
Applied 3 `ladon` SQL migrations.
Applying `client` SQL migrations...
[...]
Migration successful!

# Let's run the server (settings explained below):
$ docker run -d \
  --name ory-hydra-example--hydra \
  --network hydraguide \
  -p 9000:4444 \
  -e SYSTEM_SECRET=$SYSTEM_SECRET \
  -e DATABASE_URL=$DATABASE_URL \
  -e ISSUER=https://localhost:9000/ \
  -e CONSENT_URL=http://localhost:9020/consent \
  -e FORCE_ROOT_CLIENT_CREDENTIALS=admin:demo-password \
  oryd/hydra:v0.11.6

# And check if it's running:
$ docker logs ory-hydra-example--hydra

time="2017-06-29T21:26:26Z" level=info msg="Connecting with postgres://*:*@postgres:5432/hydra?sslmode=disable"
time="2017-06-29T21:26:26Z" level=info msg="Connected to SQL!"
[...]
time="2017-06-29T21:26:34Z" level=info msg="Setting up http server on :4444"
```

Let's dive into the various settings:

* `--network hydraguide` connects this instance to the network and makes it possible to connect to the PostgreSQL database.
* `-p 9000:4444` exposes ORY Hydra on `https://localhost:9000/`.
* `-e SYSTEM_SECRET=$SYSTEM_SECRET` sets the system secret environment variable **(required)**.
* `-e DATABASE_URL=$DATABASE_URL` sets the database url environment variable **(required)**.
* `-e ISSUER=https://localhost:9000/` set issuer to the publicly accessible url **(required)**.
* `-e CONSENT_URL=http://localhost:9020/consent` set the url of the consent app to this one. We will set up the consent
app in the following sections **(required)**.
* `-e FORCE_ROOT_CLIENT_CREDENTIALS=admin:demo-password` sets the credentials of the root account. Use the root
account to manage your ORY Hydra instance. If this is not set, ORY Hydra will auto-generate a client and display
the credentials in the logs **(optional)**.

To confirm that the instance is running properly, [open the health check](https://localhost:9000/health/status). If asked,
accept the self signed certificate in your browser. You should simply see `ok`.

On start up, ORY Hydra is initializing some values. Let's take a look at the logs:

```
$ docker logs ory-hydra-example--hydra
time="2017-06-30T09:06:34Z" level=info msg="Connecting with postgres://*:*@postgres:5432/hydra?sslmode=disable"
time="2017-06-30T09:06:34Z" level=info msg="Connected to SQL!"
time="2017-06-30T09:06:34Z" level=info msg="Key pair for signing hydra.openid.id-token is missing. Creating new one."
time="2017-06-30T09:06:35Z" level=info msg="Setting up telemetry - for more information please visit https://ory.gitbooks.io/hydra/content/telemetry.html"
time="2017-06-30T09:06:35Z" level=info msg="Key pair for signing hydra.consent.response is missing. Creating new one."
time="2017-06-30T09:06:39Z" level=info msg="Key pair for signing hydra.consent.challenge is missing. Creating new one."
time="2017-06-30T09:06:41Z" level=warning msg="No clients were found. Creating a temporary root client..."
time="2017-06-30T09:06:41Z" level=info msg="Temporary root client created."
time="2017-06-30T09:06:41Z" level=warning msg="No TLS Key / Certificate for HTTPS found. Generating self-signed certificate."
time="2017-06-30T09:06:41Z" level=info msg="Setting up http server on :4444"
```

As you can see, the following steps are performed when running ORY Hydra against a fresh database:

1. If no system secret was given (in our case we provided one), a random one is generated and emitted to the logs.
Note this down, otherwise you won't be able to restart Hydra.
2. Cryptographic keys are generated for the OpenID Connect ID Token, the consent challenge and response, and TLS encryption
using a self-signed certificate, which is why we need to run all commands using --skip-tls-verify.
3. If the OAuth 2.0 Client database table is empty, a new root client with random credentials is created. Root clients
have access to all APIs, OAuth 2.0 flows and are allowed to do everything. If the `FORCE_ROOT_CLIENT_CREDENTIALS` environment.
is set, those credentials will be used instead.

ORY Hydra can be managed using the Hydra Command Line Interface (CLI), which is using ORY Hydra's REST APIs. To
see the available commands, run:

```
$ docker run --rm -it --entrypoint hydra oryd/hydra:v0.11.6 help
Hydra is a cloud native high throughput OAuth2 and OpenID Connect provider

Usage:
  hydra [command]

[...]
```

### Install ORY Hydra without Docker

You can also install ORY Hydra without docker. For the purpose of this tutorial, [please skip this section for now](#configure-ory-hydra), and read
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
$ dep ensure --vendor-only
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

Next we will take a look at configuring ORY Hydra, specifically:

1. we create an OAuth 2.0 Client to use in the consent app (id: `consent-app`).
2. we set up access control policies for `consent-app`.
3. we add a policy that allows anybody to access the public keys for validating OpenID Connect ID Tokens.
3. we create another client capable of performing the client credentials grant and the authorize code grant.

To do so, we will repeatedly use the same docker command layout

```
$ docker run --rm -it \
  -e CLUSTER_URL=https://ory-hydra-example--hydra:4444 \
  -e CLIENT_ID=admin \
  -e CLIENT_SECRET=demo-password \
  --network hydraguide \
  oryd/hydra:v0.11.6 \
  <command>
```

where

* `--rm` removes the container after it is done.
* `-it` allows interactive processes.
* `-e CLUSTER_URL=https://ory-hydra-example--hydra:4444`, `-e CLIENT_ID=admin`, `-e CLIENT_SECRET=demo-password` tell
the command what credentials to use and where hydra is hosted. If you use the ORY Hydra CLI locally, you can skip this
step by running `hydra connect`.
* `--network hydraguide` connects the container to the network, so it is actually able to speak to our ORY Hydra host process.

Ready? Let's go!

```
# Issue an access token to validate that everything is working (ps: we need to disable TLS verification
# because the TLS certificate is self-signed).
$ docker run --rm -it \
  -e CLUSTER_URL=https://ory-hydra-example--hydra:4444 \
  -e CLIENT_ID=admin \
  -e CLIENT_SECRET=demo-password \
  --network hydraguide \
  oryd/hydra:v0.11.6 \
  token client --skip-tls-verify

tY9tGakiYAUn8VIGn_yCDlTahckSfGbDQIlXahjXtX0.BQlCxRDL3ngag6hdsSl9N2qrz7R399cQMfld8aI2Mlg

# We can also validate the token. Please make sure to copy the previous value here.
$ docker run --rm -it \
  -e CLUSTER_URL=https://ory-hydra-example--hydra:4444 \
  -e CLIENT_ID=admin \
  -e CLIENT_SECRET=demo-password \
  -e TOKEN=$token \
  --network hydraguide \
  oryd/hydra:v0.11.6 \
  token validate --skip-tls-verify \  
  tY9tGakiYAUn8VIGn_yCDlTahckSfGbDQIlXahjXtX0.BQlCxRDL3ngag6hdsSl9N2qrz7R399cQMfld8aI2Mlg

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

### Setting up the consent app

The consent app is a bridge between ORY Hydra and your identity provider (log in, sign up, reset password, ...). The consent app requires
a registered OAuth 2.0 Client at ORY Hydra. Do not re-use this client for any other purposes. To create the client,
use the following command:

```
# Create the client for consent-app
$ docker run --rm -it \
  -e CLUSTER_URL=https://ory-hydra-example--hydra:4444 \
  -e CLIENT_ID=admin \
  -e CLIENT_SECRET=demo-password \
  --network hydraguide \
  -p 9010:4445 \
  oryd/hydra:v0.11.6 \
  clients create --skip-tls-verify \
    --id consent-app \
    --secret consent-secret \
    --name "Consent App Client" \
    --grant-types client_credentials \
    --response-types token \
    --allowed-scopes hydra.consent
```

Let's dive into the arguments:

* `--id consent-app` is the id of the client.
* `--secret consent-secret` sets the secret of the client. If no secret is provided, one will be generated.
The secret is visible *only once and can not be retrieved again*.
* `--name "Consent App Client"` is a human-readable name.
* `--grant-types client_credentials` allows this client to perform the OAuth 2.0 Client Credentials grant.
* `--response-types token` allows this client to request access tokens, but not authorize codes or refresh tokens.
* `--scope hydra.consent` allows this client to request access tokens capable of managing consent requests.

Cool, next we need to create a policy for this client as well. ORY Hydra uses policies to decide whether a user
is allowed to do something in the system or not. It is different from OAuth 2.0 Scopes, as those apply only to the
access token itself, not the user. The following policy allows `consent-app` to access the
required cryptographic keys for validating and signing the consent challenge and response:

```
# For more information on access control policies, please read
# https://ory.gitbooks.io/hydra/content/security.html#access-control-policies
$ docker run --rm -it \
  -e CLUSTER_URL=https://ory-hydra-example--hydra:4444 \
  -e CLIENT_ID=admin \
  -e CLIENT_SECRET=demo-password \
  --network hydraguide \
  -p 9010:4445 \
  oryd/hydra:v0.11.6 \
  policies create --skip-tls-verify \
    --actions get,accept,reject \
    --description "Allow consent-app to manage OAuth2 consent requests." \
    --allow \
    --id consent-app-policy \
    --resources "rn:hydra:oauth2:consent:requests:<.*>" \
    --subjects consent-app

Created policy consent-app-policy.
```

Let's take a look at the arguments:

* `--actions get,accept,reject` we need to be able to get, accept and reject consent requests.
* `--allow` sets the policy effect to `allow`. Omit to set this for `deny`.
* `--id consent-app-policy` a unique identifier.
* `--resources "rn:hydra:oauth2:consent:requests:<.*>" ` we need to be able to access the consent request resource.
* `--subjects consent-app` the subject ("user") of this policy is our consent app.

Awesome! Next we will run the [ORY Hydra Consent App Example (NodeJS)](https://github.com/ory/hydra-consent-app-express).
This app is also available in [Golang](https://github.com/ory/hydra-consent-app-go), but for the purpose of this
tutorial we will use the NodeJS one:

```
$ docker run -d \
  --name ory-hydra-example--consent \
  -p 9020:3000 \
  --network hydraguide \
  -e HYDRA_CLIENT_ID=consent-app \
  -e HYDRA_CLIENT_SECRET=consent-secret \
  -e HYDRA_URL=https://ory-hydra-example--hydra:4444 \
  -e NODE_TLS_REJECT_UNAUTHORIZED=0 \
  oryd/hydra-consent-app-express:v0.10.10

# Let's check if it's running ok:
$ docker logs ory-hydra-example--consent
```

Let's take a look at the arguments:
* `-p 9020:3000` exposes this service at port 9020. If you remember, that's the port of the `CONSENT_URL` value
from the ORY Hydra docker container (`CONSENT_URL=http://localhost:9020/consent`).
* `-e HYDRA_CLIENT_ID=consent-app` this is the client id we created in the steps above.
* `-e HYDRA_CLIENT_SECRET=consent-secret` this is the client secret we set in the steps above.
* `HYDRA_URL=http://hydra:4444` point to the ORY Hydra container.
* `NODE_TLS_REJECT_UNAUTHORIZED=0` disables TLS verification, because we are using self-signed certificates.

## Perform OAuth 2.0 Flow

Awesome, our infrastructure is set up! Now it's time to create an OAuth 2.0 Consumer and perform the OAuth 2.0 Authorize Code flow.
To do so, we will create a new client:

```
$ docker run --rm -it \
  -e CLUSTER_URL=https://ory-hydra-example--hydra:4444 \
  -e CLIENT_ID=admin \
  -e CLIENT_SECRET=demo-password \
  --network hydraguide \
  oryd/hydra:v0.11.6 \
  clients create --skip-tls-verify \
    --id some-consumer \
    --secret consumer-secret \
    --grant-types authorization_code,refresh_token,client_credentials,implicit \
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
* `--allowed-scopes hydra.clients.*` allows this client to request scope `hydra.clients` and all scopes prefixed with `hydra.clients.`, for example `hydra.clients.get`.
* `--callbacks http://localhost:9010/callback` allows the client to request this redirect uri.

Also, we want to allow everyone (not only our consumer) access to the public key of the OpenID Connect ID Token, which can be achieved with:

```
$ docker run --rm -it \
  -e CLUSTER_URL=https://ory-hydra-example--hydra:4444 \
  -e CLIENT_ID=admin \
  -e CLIENT_SECRET=demo-password \
  --network hydraguide \
  oryd/hydra:v0.11.6 \
  policies create --skip-tls-verify \
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
$ docker run --rm -it \
  --network hydraguide \
  -p 9010:4445 \
  oryd/hydra:v0.11.6 \
  token user --skip-tls-verify \
    --auth-url https://localhost:9000/oauth2/auth \
    --token-url https://ory-hydra-example--hydra:4444/oauth2/token \
    --id some-consumer \
    --secret consumer-secret \
    --scopes openid,offline,hydra.clients \
    --redirect http://localhost:9010/callback

Setting up callback listener on http://localhost:4445/callback
Press ctrl + c on Linux / Windows or cmd + c on OSX to end the process.
If your browser does not open automatically, navigate to:

        https://localhost:9000/oauth2/auth?client_id=some-consumer&redirect_uri=http%3A%2F%2Flocalhost%3A9020%2Fcallback&response_type=code&scope=openid+offline+hydra.clients&state=hfcyxoqoctwbnvrxrsuwgzfu&nonce=lbeouolavuvcdhjefcnzlqur
```

open the link, as prompted, in your browser, and follow the steps shown there. When completed, you should land
at a screen that looks like this one:

![OAuth 2.0 result](../images/install-result.png)
