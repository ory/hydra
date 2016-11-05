### 5 Minute Tutorial

In this example, you will set up Hydra, a Postgres instance and an exemplary identity provider written in React using docker compose. It will take you about 5 minutes to get complete this tutorial.

<img src="images/oauth2-flow.gif" alt="OAuth2 Flow">

<img alt="Running the example" align="right" width="35%" src="images/run-the-example.gif">

Install the [CLI and Docker](https://github.com/ory-am/hydra#installation). Make sure you install Docker Compose as well.

We will use a dummy password as the system secret: `SYSTEM_SECRET=passwordtutorialpasswordtutorial`.
Use a very secure secret in production.

```
$ go get github.com/ory-am/hydra
$ cd $GOPATH/src/github.com/ory-am/hydra
$ SYSTEM_SECRET=passwordtutorial DOCKER_IP=localhost docker-compose up --build
Starting hydra_mysqld_1
Starting hydra_postgresd_1
Starting hydra_hydra_1

[...]
```

You now have a running hydra docker container! Additionally, a Postgres image was deployed as well as a consent app.
Next, let us manage the host process. You can use the Hydra CLI by ssh'ing to the docker container:

```
$ docker exec -i -t hydra_hydra_1 /bin/bash
root@b4403bb4147f:/go/src/github.com/ory-am/hydra#
```

Let's start by creating a new client:

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

Let's try this with the authorize code grant!

```
$ hydra token user
Setting up callback listener on http://localhost:4445/callback
Press ctrl + c on Linux / Windows or cmd + c on OSX to end the process.
If your browser does not open automatically, navigate to:

    https://192.168.99.100:4444/oauth2/...
```

Great! You installed hydra, connected the CLI, created a client and completed two authentication flows!
