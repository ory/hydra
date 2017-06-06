# How should I run migrations?

Since ORY Hydra 0.8.0, migrations are no longer run automatically on boot. This is required in production environments,
because:

1. Although SQL migrations are tested, migrating schemas can cause data loss and should only be done consciously with
prior back ups.
2. Running a production system with a user that has right such as ALTER TABLE is a security anti-pattern.

Thus, to initialize the database schemas, it is required to run `hydra migrate sql driver://user:password@host:port/db` before running
`hydra host`.

## What does the installation process look like?

1. Run `hydra migrate sql ...` on a host close to the database (e.g. a virtual machine with access to the SQL instance).

## What does a migration process look like?

1. Make sure a database update is required by checking the release notes.
2. Make a back up of the database.
3. Run the migration script on a host close to the database (e.g. a virtual machine with access to the SQL instance).
Schemas are usually backwards compatible, so instances running previous versions of ORY Hydra should keep working fine.
If backwards compatibility is not given, this will be addressed in the patch notes.
4. Upgrade all ORY Hydra instances.

## How can I do this in docker?

Many deployments of ORY Hydra use Docker. Although several options are available, we advise to extend the ORY Hydra Docker
image

**Dockerfile**
```
FROM oryd/hydra:tag

ENTRYPOINT /go/bin/hydra migrate sql $DATABASE_URL
```

and run it in your infrastructure once.

Additionally, *but not recommended*, it is possible to override the entry point of the ORY Hydra Docker image using CLI flag
`--entrypoint "hydra migrate sql $DATABASE_URL; hydra host"` or with `entrypoint: hydra migrate sql $DATABASE_URL; hydra host`
set in your docker compose config.
