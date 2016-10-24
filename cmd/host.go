package cmd

import (
	"github.com/ory-am/hydra/cmd/server"
	"github.com/spf13/cobra"
)

// hostCmd represents the host command
var hostCmd = &cobra.Command{
	Use:   "host",
	Short: "Start the HTTP/2 host service",
	Long: `Starts all HTTP/2 APIs and connects to a database backend.

This command exposes a variety of controls via environment variables. You can
set environments using "export KEY=VALUE" (Linux/macOS) or "set KEY=VALUE" (Windows). On Linux,
you can also set environments by prepending key value pairs: "KEY=VALUE KEY2=VALUE2 hydra"

All possible controls are listed below. The host process additionally exposes a few flags, which are listed below
the controls section.

CORE CONTROLS
=============

- DATABASE_URL: A URL to a persistent backend. Hydra supports various backends:
  - None: If DATABASE_URL is empty, all data will be lost when the command is killed.
  - Postgres: If DATABASE_URL is a DSN starting with postgres:// PostgreSQL will be used as storage backend.
	Example: DATABASE_URL=rethinkdb://user:password@host:123/database

	If PostgreSQL is not serving TLS, append ?sslmode=disable to the url:
	DATABASE_URL=rethinkdb://user:password@host:123/database?sslmode=disable

  - MySQL: If DATABASE_URL is a DSN starting with mysql:// MySQL will be used as storage backend.
	Example: DATABASE_URL=mysql://user:password@tcp(host:123)/database?parseTime=true

	Be aware that the ?parseTime=true parameter is mandatory, or timestamps will not work.

  - RethinkDB: If DATABASE_URL is a DSN starting with rethinkdb:// RethinkDB will be used as storage backend.
	Example: DATABASE_URL=rethinkdb://user:password@host:123/database

	Additionally, these controls are available when using RethinkDB:
	- RETHINK_TLS_CERT_PATH: The path to the TLS certificate (pem encoded) used to connect to rethinkdb.
		Example: RETHINK_TLS_CERT_PATH=~/rethink.pem

	- RETHINK_TLS_CERT: A pem encoded TLS certificate passed as string. Can be used instead of RETHINK_TLS_CERT_PATH.
		Example: RETHINK_TLS_CERT_PATH="-----BEGIN CERTIFICATE-----\nMIIDZTCCAk2gAwIBAgIEV5xOtDANBgkqhkiG9w0BAQ0FADA0MTIwMAYDVQQDDClP..."

- SYSTEM_SECRET: A secret that is at least 16 characters long. If none is provided, one will be generated. They key
	is used to encrypt sensitive data using AES-GCM (256 bit) and validate HMAC signatures.
	Example: SYSTEM_SECRET=jf89-jgklAS9gk3rkAF90dfsk

- FORCE_ROOT_CLIENT_CREDENTIALS: On first start up, Hydra generates a root client with random id and secret. Use
	this environment variable in the form of "FORCE_ROOT_CLIENT_CREDENTIALS=id:secret" to set
	the client id and secret yourself.
	Example: FORCE_ROOT_CLIENT_CREDENTIALS=admin:kf0AKfm12fas3F-.f

- PORT: The port hydra should listen on.
	Defaults to PORT=4444

- HOST: The host interface hydra should listen on. Leave empty to listen on all interfaces.
	Example: HOST=localhost

- BCRYPT_COST: Set the bcrypt hashing cost. This is a trade off between
	security and performance. Range is 4 =< x =< 31.
	Defaults to BCRYPT_COST=10


OAUTH2 CONTROLS
===============

- CONSENT_URL: The uri of the consent endpoint.
	Example: CONSENT_URL=https://id.myapp.com/consent

- ISSUER: The issuer is used for identification in all OAuth2 tokens.
	Defaults to ISSUER=hydra.localhost

- AUTH_CODE_LIFESPAN: Lifespan of OAuth2 authorize codes. Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
	Defaults to AUTH_CODE_LIFESPAN=10m

- ID_TOKEN_LIFESPAN: Lifespan of OpenID Connect ID Tokens. Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
	Defaults to ID_TOKEN_LIFESPAN=1h

- ACCESS_TOKEN_LIFESPAN: Lifespan of OAuth2 access tokens. Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
	Defaults to ACCESS_TOKEN_LIFESPAN=1h

- CHALLENGE_TOKEN_LIFESPAN: Lifespan of OAuth2 consent tokens. Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
	Defaults to CHALLENGE_TOKEN_LIFESPAN=10m


HTTPS CONTROLS
==============

- HTTPS_ALLOW_TERMINATION_FROM: Whitelist one or multiple CIDR address ranges and allow them to terminate TLS connections.
	Be aware that the X-Forwarded-Proto header must be set and must never be modifiable by anyone but
	your proxy / gateway / load balancer. Supports ipv4 and ipv6.
	Hydra serves http instead of https when this option is set.
	Example: HTTPS_ALLOW_TERMINATION_FROM=127.0.0.1/32,192.168.178.0/24,2620:0:2d0:200::7/32

- HTTPS_TLS_CERT_PATH: The path to the TLS certificate (pem encoded).
	Example: HTTPS_TLS_CERT_PATH=~/cert.pem

- HTTPS_TLS_KEY_PATH: The path to the TLS private key (pem encoded).
	Example: HTTPS_TLS_KEY_PATH=~/key.pem

- HTTPS_TLS_CERT: A pem encoded TLS certificate passed as string. Can be used instead of HTTPS_TLS_CERT_PATH.
	Example: HTTPS_TLS_CERT="-----BEGIN CERTIFICATE-----\nMIIDZTCCAk2gAwIBAgIEV5xOtDANBgkqhkiG9w0BAQ0FADA0MTIwMAYDVQQDDClP..."

- HTTPS_TLS_KEY: A pem encoded TLS key passed as string. Can be used instead of HTTPS_TLS_KEY_PATH.
	Example: HTTPS_TLS_KEY="-----BEGIN ENCRYPTED PRIVATE KEY-----\nMIIFDjBABgkqhkiG9w0BBQ0wMzAbBgkqhkiG9w0BBQwwDg..."


DEBUG CONTROLS
==============

- PROFILING: Set "PROFILING=cpu" to enable cpu profiling and "PROFILING=memory" to enable memory profiling.
	It is not possible to do both at the same time.
	Example: PROFILING=cpu
`,
	Run: server.RunHost(c),
}

func init() {
	RootCmd.AddCommand(hostCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hostCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	hostCmd.Flags().BoolVar(&c.ForceHTTP, "dangerous-force-http", false, "Disable HTTP/2 over TLS (HTTPS) and serve HTTP instead. Never use this in production.")
	hostCmd.Flags().Bool("dangerous-auto-logon", false, "Stores the root credentials in ~/.hydra.yml. Do not use in production.")
	hostCmd.Flags().String("https-tls-key-path", "", "Path to the key file for HTTP/2 over TLS (https). You can set HTTPS_TLS_KEY_PATH or HTTPS_TLS_KEY instead.")
	hostCmd.Flags().String("https-tls-cert-path", "", "Path to the certificate file for HTTP/2 over TLS (https). You can set HTTPS_TLS_CERT_PATH or HTTPS_TLS_CERT instead.")
	hostCmd.Flags().String("rethink-tls-cert-path", "", "Path to the certificate file to connect to rethinkdb over TLS (https). You can set RETHINK_TLS_CERT_PATH or RETHINK_TLS_CERT instead.")
}
