// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sqlcon

// HelpMessage returns a string explaining how to set up SQL using environment variables.
func HelpMessage() string {
	return `- DATABASE_URL: A DSN to a persistent backend. Various backends are supported:

  - Changes are lost on process death (ephemeral storage):

	- Memory: If DATABASE_URL is "memory", data will be written to memory and is lost when you restart this instance.
	  Example: DATABASE_URL=memory

  - Changes are kept after process death (persistent storage):

    - SQL Databases: Officially, PostgreSQL, MySQL and CockroachDB are supported. This project works best with PostgreSQL.

	  - PostgreSQL: If DATABASE_URL is a DSN starting with postgres://, PostgreSQL will be used as storage backend.
		Example: DATABASE_URL=postgres://user:password@host:123/database

		Additionally, the following query/DSN parameters are supported:

      	* max_conns (number): Sets the maximum number of open connections to the database. Defaults to the number of CPU cores times 2.
		* max_idle_conns (number): Sets the maximum number of connections in the idle. Defaults to the number of CPU cores.
        * max_conn_lifetime (duratino): Sets the maximum amount of time ("ms", "s", "m", "h") a connection may be reused.
		  Defaults to 0s (disabled).
		* sslmode (string): Whether or not to use SSL (default is require)
		  * disable - No SSL
		  * require - Always SSL (skip verification)
		  * verify-ca - Always SSL (verify that the certificate presented by the
		    server was signed by a trusted CA)
		  * verify-full - Always SSL (verify that the certification presented by
		    the server was signed by a trusted CA and the server host name
		    matches the one in the certificate)
		* fallback_application_name (string): An application_name to fall back to if one isn't provided.
		* connect_timeout (number): Maximum wait for connection, in seconds. Zero or
		  not specified means wait indefinitely.
		* sslcert (string): Cert file location. The file must contain PEM encoded data.
		* sslkey (string): Key file location. The file must contain PEM encoded data.
		* sslrootcert (string): The location of the root certificate file. The file
		  must contain PEM encoded data.
		Example: DATABASE_URL=postgres://user:password@host:123/database?sslmode=verify-full

	  - MySQL: If DATABASE_URL is a DSN starting with mysql:// MySQL will be used as storage backend.
		Be aware that the ?parseTime=true parameter is mandatory, or timestamps will not work.
		Example: DATABASE_URL=mysql://user:password@tcp(host:123)/database?parseTime=true

		Additionally, the following query/DSN parameters are supported:
		* collation (string): Sets the collation used for client-server interaction on connection. In contrast to charset, 
		  collation does not issue additional queries. If the specified collation is unavailable on the target server,
		  the connection will fail.
		* loc (string): Sets the location for time.Time values. Note that this sets the location for time.Time values
		  but does not change MySQL's time_zone setting. For that set the time_zone DSN parameter. Please keep in mind,
		  that param values must be url.QueryEscape'ed. Alternatively you can manually replace the / with %2F.
		  For example US/Pacific would be loc=US%2FPacific.
		* maxAllowedPacket (number): Max packet size allowed in bytes. The default value is 4 MiB and should be
		  adjusted to match the server settings. maxAllowedPacket=0 can be used to automatically fetch the max_allowed_packet variable from server on every connection.
		* readTimeout (duration): I/O read timeout. The value must be a decimal number with a unit suffix
		  ("ms", "s", "m", "h"), such as "30s", "0.5m" or "1m30s".
		* timeout (duration): Timeout for establishing connections, aka dial timeout. The value must be a decimal number with a unit suffix
		  ("ms", "s", "m", "h"), such as "30s", "0.5m" or "1m30s".
		* tls (bool / string): tls=true enables TLS / SSL encrypted connection to the server. Use skip-verify if
		  you want to use a self-signed or invalid certificate (server side).
		* writeTimeout (duration): I/O write timeout. The value must be a decimal number with a unit suffix
		  ("ms", "s", "m", "h"), such as "30s", "0.5m" or "1m30s".
		Example: DATABASE_URL=mysql://user:password@tcp(host:123)/database?parseTime=true&writeTimeout=123s

	  - CockroachDB: If DATABASE_URL is a DSN starting with cockroach://, CockroachDB will be used as storage backend.
		Example: DATABASE_URL=cockroach://user:password@host:123/database

		Additionally, the following query/DSN parameters are supported:
		* sslmode (string): Whether or not to use SSL (default is require)
		  * disable - No SSL
		  * require - Always SSL (skip verification)
		  * verify-ca - Always SSL (verify that the certificate presented by the
		    server was signed by a trusted CA)
		  * verify-full - Always SSL (verify that the certification presented by
		    the server was signed by a trusted CA and the server host name
		    matches the one in the certificate)
		* application_name (string): An initial value for the application_name session variable.
		* sslcert (string): Cert file location. The file must contain PEM encoded data.
		* sslkey (string): Key file location. The file must contain PEM encoded data.
		* sslrootcert (string): The location of the root certificate file. The file
		  must contain PEM encoded data.
		Example: DATABASE_URL=cockroach://user:password@host:123/database?sslmode=verify-full`
}
