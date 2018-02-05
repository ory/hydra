# Dependencies & Environment

ORY Hydra is built cloud native and implements [12factor](http://12factor.net) principles. The Docker Image is 5 MB light
and thoroughly versioned with [verbose upgrade instructions](https://github.com/ory/hydra/blob/master/UPGRADE.md)
and [detailed changelogs](https://github.com/ory/hydra/blob/master/CHANGELOG.md). Auto-scaling, migrations, health checks,
it all works with zero additional work required. It is possible to run ORY Hydra on any platform, including but not limited
to OSX, Linux, Windows, ARM, FreeBSD and more.

ORY Hydra has two operational modes:

* In-memory: This mode does not work with more than one instance ("cluster") and any state is lost after restarting the instance.
* SQL: This mode works with more than one instance and state is not lost after restarts.

The SQL adapter supports two DBMS: PostgreSQL 9.6+ and MySQL 5.7+. Please note that
older MySQL versions have issues with ORY Hydra's database schema. For more information [go here](https://github.com/ory/hydra/issues/377).

No further dependencies are required for a production-ready instance.
