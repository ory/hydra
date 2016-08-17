### 5 minutes tutorial: Run your very own OAuth2 environment

In this example, you will set up Hydra, a RethinkDB instance and an exemplary identity provider written in React using docker compose.
It will take you about 5 minutes to get complete this tutorial.

<img src="dist/images/oauth2-flow.gif" alt="OAuth2 Flow">

<img alt="Running the example" align="right" width="35%" src="dist/images/run-the-example.gif">

Install the [CLI and Docker](https://github.com/ory-am/hydra#installation). Make sure you install Docker Compose as well.

We will use a dummy password as the system secret: `SYSTEM_SECRET=passwordtutorialpasswordtutorial`. Use a very secure secret in production.

```
$ go get github.com/ory-am/hydra
$ cd $GOPATH/src/github.com/ory-am/hydra
$ docker-compose build
Building hydra
[...]
$ SYSTEM_SECRET=passwordtutorial DOCKER_IP=localhost docker-compose up
Starting hydra_rethinkdb_1
[...]
mhydra   | mtime="2016-05-17T18:09:28Z" level=warning msg="Generated system secret: MnjFP5eLIr60h?hLI1h-!<4(TlWjAHX7"
[...]
mhydra   | mtime="2016-05-17T18:09:29Z" level=warning msg="client_id: d9227bd5-5d47-4557-957d-2fd3bee11035"
mhydra   | mtime="2016-05-17T18:09:29Z" level=warning msg="client_secret: ,IvxGt02uNjv1ur9"
[...]
```

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
Client ID: c003830f-a090-4721-9463-92424270ce91
Client Secret: Z2pJ0>Tp7.ggn>EE&rhnOzdt1
```

**Important note:** if no certificate is provided, Hydra uses self-signed TLS certificates for HTTPS. This should
never be done in production. To skip the TLS verification step on the client, provide the `--skip-tls-verify` flag.

Why not issue an access token for your client?

```
$ hydra token client --skip-tls-verify
JLbnRS9GQmzUBT4x7ESNw0kj2wc0ffbMwOv3QQZW4eI.qkP-IQXn6guoFew8TvaMFUD-SnAyT8GmWuqGi3wuWXg
```

Let's try this with the authorize code grant!

```
$ hydra token user --skip-tls-verify
If your browser does not open automatically, navigate to: https://192.168.99.100:4444/oauth2/...
Setting up callback listener on http://localhost:4445/callback
Press ctrl + c on Linux / Windows or cmd + c on OSX to end the process.
```

Great! You installed hydra, connected the CLI, created a client and completed two authentication flows!
