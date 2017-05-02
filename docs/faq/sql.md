## How can I control SQL connection limits?

You can configure SQL connection limits by appending parameters `max_conns`, `max_idle_conns`, or `max_conn_lifetime`
to the DSN: `postgres://foo:bar@host:port/database?max_conns=12`.
