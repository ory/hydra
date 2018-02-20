# 5 Minute Tutorial

To start off easy, ORY Hydra provides a docker-compose based example for setting up ORY Hydra, a PostgreSQL instance
and an exemplary consent app (identity provider). You need to have the latest Docker version installed.

<img src="../images/oauth2-flow.gif" alt="OAuth2 Flow">

<img alt="Running the example" align="right" width="35%" src="../images/run-the-example.gif">

Install [Docker and Docker Compose](https://github.com/ory-am/hydra#installation) and either clone the Hydra git repository,
download [this zip file](https://github.com/ory-am/hydra/archive/master.zip) or use `go get github.com/ory/hydra` if you have Go (1.8+) installed on you system.

```
$ git clone https://github.com/ory/hydra.git
$ cd hydra
$ git checkout tags/v0.10.10
$ docker-compose -p hydra up --build -d
Starting hydra_mysqld_1
Starting hydra_postgresd_1
Starting hydra_hydra_1

[...]
```

Perfect, everything is running now! Let's SSH into the ORY Hydra container and play around with some of the commands:

```
$ docker exec -i -t hydra_hydra_1 /bin/sh
root@b4403bb4147f:/go/src/github.com/ory-am/hydra$

# Creates a new OAuth 2.0 client
$ hydra clients create
Client ID: c003830f-a090-4721-9463-92424270ce91
Client Secret: Z2pJ0>Tp7.ggn>EE&rhnOzdt1

# Issues a token for the root client (id: admin)
$ hydra token client
JLbnRS9GQmzUBT4x7ESNw0kj2wc0ffbMwOv3QQZW4eI.qkP-IQXn6guoFew8TvaMFUD-SnAyT8GmWuqGi3wuWXg

# Introspects a token:
$ hydra token validate $(hydra token client)
```

Next, we will perform the OAuth 2.0 Authorization Code Grant:

```
$ hydra token user --auth-url http://localhost:4444/oauth2/auth --token-url http://localhost:4444/oauth2/token
Setting up callback listener on http://localhost:4445/callback
Press ctrl + c on Linux / Windows or cmd + c on OSX to end the process.
If your browser does not open automatically, navigate to:

    https://192.168.99.100:4444/oauth2/...
```

Great! You installed hydra, connected the CLI, created a client and completed two authentication flows!
