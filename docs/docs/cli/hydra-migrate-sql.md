---
id: hydra-migrate-sql
title: hydra migrate sql
description: hydra migrate sql Create SQL schemas and apply migration plans
---

<!--
This file is auto-generated.

To improve this file please make your change against the appropriate "./cmd/*.go" file.
-->

## hydra migrate sql

Create SQL schemas and apply migration plans

### Synopsis

Run this command on a fresh SQL installation and when you upgrade Hydra to a new
minor version. For example, upgrading Hydra 0.7.0 to 0.8.0 requires running this
command.

It is recommended to run this command close to the SQL instance (e.g. same
subnet) instead of over the public internet. This decreases risk of failure and
decreases time required.

You can read in the database URL using the -e flag, for example: export DSN=...
hydra migrate sql -e

### WARNING

Before running this command on an existing database, create a back up!

```
hydra migrate sql <database-url> [flags]
```

### Options

```
  -h, --help            help for sql
  -e, --read-from-env   If set, reads the database connection string from the environment variable DSN or config file key dsn.
  -y, --yes             If set all confirmation requests are accepted without user interaction.
```

### Options inherited from parent commands

```
  -c, --config string   Config file (default is $HOME/hydra.yaml)
```

### SEE ALSO

- [hydra migrate](hydra-migrate) - Various migration helpers
