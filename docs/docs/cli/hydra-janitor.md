---
id: hydra-janitor
title: hydra janitor
description:
  hydra janitor Clean the database of old tokens and login/consent requests
---

<!--
This file is auto-generated.

To improve this file please make your change against the appropriate "./cmd/*.go" file.
-->

## hydra janitor

Clean the database of old tokens and login/consent requests

### Synopsis

This command will cleanup any expired oauth2 tokens as well as login/consent
requests. This will select records to delete with a limit and delete records in
batch to ensure that no table locking issues arise in big production databases.

### Warning

This command is in beta. Proceed with caution!

This is a destructive command and will purge data directly from the database.
Please use this command with caution if you need to keep historic data for any
reason.

###############

Janitor can be used in several ways.

1.  By passing the database connection string (DSN) as an argument Pass the
    database url (dsn) as an argument to janitor. E.g. janitor
    &lt;database-url&gt;
2.  By passing the DSN as an environment variable

        export DSN=...
        janitor -e

3.  By passing a configuration file containing the DSN janitor -c
    /path/to/conf.yml
4.  Extra _optional_ parameters can also be added such as

        janitor --keep-if-younger 23h --access-lifespan 1h --refresh-lifespan 40h --consent-request-lifespan 10m &lt;database-url&gt;

5.  Running only a certain cleanup

        janitor --tokens &lt;database-url&gt;

    or

        janitor --requests &lt;database-url&gt;

    or both

        janitor --tokens --requests &lt;database-url&gt;

```
hydra janitor [&lt;database-url&gt;] [flags]
```

### Options

```
      --access-lifespan duration            Set the access token lifespan e.g. 1s, 1m, 1h.
      --batch-size int                      Define how many records are deleted with each iteration. (default 100)
  -c, --config strings                      Path to one or more .json, .yaml, .yml, .toml config files. Values are loaded in the order provided, meaning that the last config file overwrites values from the previous config file.
      --consent-request-lifespan duration   Set the login/consent request lifespan e.g. 1s, 1m, 1h
  -h, --help                                help for janitor
      --keep-if-younger duration            Keep database records that are younger than a specified duration e.g. 1s, 1m, 1h.
      --limit int                           Limit the number of records retrieved from database for deletion. (default 10000)
  -e, --read-from-env                       If set, reads the database connection string from the environment variable DSN or config file key dsn.
      --refresh-lifespan duration           Set the refresh token lifespan e.g. 1s, 1m, 1h.
      --requests                            This will only run the cleanup on requests and will skip token cleanup.
      --tokens                              This will only run the cleanup on tokens and will skip requests cleanup.
```

### SEE ALSO

- [hydra](hydra) - Run and manage ORY Hydra
