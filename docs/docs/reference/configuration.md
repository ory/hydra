---
id: configuration
title: Configuration
---

<!-- THIS FILE IS BEING AUTO-GENERATED. DO NOT MODIFY IT AS ALL CHANGES WILL BE OVERWRITTEN.
OPEN AN ISSUE IF YOU WOULD LIKE TO MAKE ADJUSTMENTS HERE AND MAINTAINERS WILL HELP YOU LOCATE THE RIGHT
FILE -->

If file `$HOME/.hydra.yaml` exists, it will be used as a configuration file
which supports all configuration settings listed below.

You can load the config file from another source using the
`-c path/to/config.yaml` or `--config path/to/config.yaml` flag:
`hydra --config path/to/config.yaml`.

Config files can be formatted as JSON, YAML and TOML. Some configuration values
support reloading without server restart. All configuration values can be set
using environment variables, as documented below.

This reference configuration documents all keys, also deprecated ones! It is a
reference for all possible configuration values.

If you are looking for an example configuration, it is better to try out the
quickstart.

To find out more about edge cases like setting string array values through
environmental variables head to the
[Configuring ORY services](https://www.ory.sh/docs/ecosystem/configuring)
section.

```yaml
## ORY Hydra Configuration
#

## log ##
#
# Configures the logger
#
log:
  ## level ##
  #
  # Sets the log level.
  #
  # Default value: info
  #
  # One of:
  # - panic
  # - fatal
  # - error
  # - warn
  # - info
  # - debug
  # - trace
  #
  # Set this value using environment variables on
  # - Linux/macOS:
  #    $ export LOG_LEVEL=<value>
  # - Windows Command Line (CMD):
  #    > set LOG_LEVEL=<value>
  #
  level: debug

  ## leak_sensitive_values ##
  #
  # Logs sensitive values such as cookie and URL parameter.
  #
  # Default value: false
  #
  # Set this value using environment variables on
  # - Linux/macOS:
  #    $ export LOG_LEAK_SENSITIVE_VALUES=<value>
  # - Windows Command Line (CMD):
  #    > set LOG_LEAK_SENSITIVE_VALUES=<value>
  #
  leak_sensitive_values: true

  ## format ##
  #
  # Sets the log format.
  #
  # Default value: text
  #
  # One of:
  # - json
  # - json_pretty
  # - text
  #
  # Set this value using environment variables on
  # - Linux/macOS:
  #    $ export LOG_FORMAT=<value>
  # - Windows Command Line (CMD):
  #    > set LOG_FORMAT=<value>
  #
  format: text

## serve ##
#
# Controls the configuration for the http(s) daemon(s).
#
serve:
  ## public ##
  #
  # Controls the public daemon serving public API endpoints like /oauth2/auth, /oauth2/token, /.well-known/jwks.json
  #
  public:
    ## port ##
    #
    # Default value: 4444
    #
    # Set this value using environment variables on
    # - Linux/macOS:
    #    $ export SERVE_PUBLIC_PORT=<value>
    # - Windows Command Line (CMD):
    #    > set SERVE_PUBLIC_PORT=<value>
    #
    port: 50196

    ## host ##
    #
    # The interface or unix socket ORY Hydra should listen and handle public API requests on. Use the prefix "unix:" to specify a path to a unix socket. Leave empty to listen on all interfaces.
    #
    # Examples:
    # - localhost
    #
    # Set this value using environment variables on
    # - Linux/macOS:
    #    $ export SERVE_PUBLIC_HOST=<value>
    # - Windows Command Line (CMD):
    #    > set SERVE_PUBLIC_HOST=<value>
    #
    host: localhost

    ## cors ##
    #
    # Configures Cross Origin Resource Sharing for public endpoints.
    #
    cors:
      ## enabled ##
      #
      # Sets whether CORS is enabled.
      #
      # Default value: false
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export SERVE_PUBLIC_CORS_ENABLED=<value>
      # - Windows Command Line (CMD):
      #    > set SERVE_PUBLIC_CORS_ENABLED=<value>
      #
      enabled: false

      ## allowed_origins ##
      #
      # A list of origins a cross-domain request can be executed from. If the special * value is present in the list, all origins will be allowed. An origin may contain a wildcard (*) to replace 0 or more characters (i.e.: http://*.domain.com). Only one wildcard can be used per origin.
      #
      # Default value: *
      #
      # Examples:
      # - - https://example.com
      #   - https://*.example.com
      #   - https://*.foo.example.com
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export SERVE_PUBLIC_CORS_ALLOWED_ORIGINS=<value>
      # - Windows Command Line (CMD):
      #    > set SERVE_PUBLIC_CORS_ALLOWED_ORIGINS=<value>
      #
      allowed_origins:
        - '*'

      ## allowed_methods ##
      #
      # A list of HTTP methods the user agent is allowed to use with cross-domain requests.
      #
      # Default value: POST,GET,PUT,PATCH,DELETE
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export SERVE_PUBLIC_CORS_ALLOWED_METHODS=<value>
      # - Windows Command Line (CMD):
      #    > set SERVE_PUBLIC_CORS_ALLOWED_METHODS=<value>
      #
      allowed_methods:
        - TRACE
        - POST
        - HEAD

      ## allowed_headers ##
      #
      # A list of non simple headers the client is allowed to use with cross-domain requests.
      #
      # Default value: Authorization,Content-Type
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export SERVE_PUBLIC_CORS_ALLOWED_HEADERS=<value>
      # - Windows Command Line (CMD):
      #    > set SERVE_PUBLIC_CORS_ALLOWED_HEADERS=<value>
      #
      allowed_headers:
        - sit dolore labore aliqua
        - aliqua
        - aute sunt pariatur proident commodo

      ## exposed_headers ##
      #
      # Sets which headers are safe to expose to the API of a CORS API specification.
      #
      # Default value: Content-Type
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export SERVE_PUBLIC_CORS_EXPOSED_HEADERS=<value>
      # - Windows Command Line (CMD):
      #    > set SERVE_PUBLIC_CORS_EXPOSED_HEADERS=<value>
      #
      exposed_headers:
        - ipsum veniam sit incididunt

      ## allow_credentials ##
      #
      # Sets whether the request can include user credentials like cookies, HTTP authentication or client side SSL certificates.
      #
      # Default value: true
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export SERVE_PUBLIC_CORS_ALLOW_CREDENTIALS=<value>
      # - Windows Command Line (CMD):
      #    > set SERVE_PUBLIC_CORS_ALLOW_CREDENTIALS=<value>
      #
      allow_credentials: true

      ## options_passthrough ##
      #
      # TODO
      #
      # Default value: false
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export SERVE_PUBLIC_CORS_OPTIONS_PASSTHROUGH=<value>
      # - Windows Command Line (CMD):
      #    > set SERVE_PUBLIC_CORS_OPTIONS_PASSTHROUGH=<value>
      #
      options_passthrough: true

      ## max_age ##
      #
      # Sets how long (in seconds) the results of a preflight request can be cached. If set to 0, every request is preceded by a preflight request.
      #
      # Minimum value: 0
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export SERVE_PUBLIC_CORS_MAX_AGE=<value>
      # - Windows Command Line (CMD):
      #    > set SERVE_PUBLIC_CORS_MAX_AGE=<value>
      #
      max_age: 65427458

      ## debug ##
      #
      # Adds additional log output to debug server side CORS issues.
      #
      # Default value: false
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export SERVE_PUBLIC_CORS_DEBUG=<value>
      # - Windows Command Line (CMD):
      #    > set SERVE_PUBLIC_CORS_DEBUG=<value>
      #
      debug: true

    ## socket ##
    #
    # Sets the permissions of the unix socket
    #
    socket:
      ## owner ##
      #
      # Owner of unix socket. If empty, the owner will be the user running hydra.
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export SERVE_PUBLIC_SOCKET_OWNER=<value>
      # - Windows Command Line (CMD):
      #    > set SERVE_PUBLIC_SOCKET_OWNER=<value>
      #
      owner: enim

      ## group ##
      #
      # Group of unix socket. If empty, the group will be the primary group of the user running hydra.
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export SERVE_PUBLIC_SOCKET_GROUP=<value>
      # - Windows Command Line (CMD):
      #    > set SERVE_PUBLIC_SOCKET_GROUP=<value>
      #
      group: incididunt culpa ut dolor

      ## mode ##
      #
      # Mode of unix socket in numeric form
      #
      # Default value: 493
      #
      # Minimum value: 0
      #
      # Maximum value: 511
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export SERVE_PUBLIC_SOCKET_MODE=<value>
      # - Windows Command Line (CMD):
      #    > set SERVE_PUBLIC_SOCKET_MODE=<value>
      #
      mode: 188

    ## access_log ##
    #
    # Access Log configuration for public server.
    #
    access_log:
      ## disable_for_health ##
      #
      # Disable access log for health endpoints.
      #
      # Default value: false
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export SERVE_PUBLIC_ACCESS_LOG_DISABLE_FOR_HEALTH=<value>
      # - Windows Command Line (CMD):
      #    > set SERVE_PUBLIC_ACCESS_LOG_DISABLE_FOR_HEALTH=<value>
      #
      disable_for_health: false

  ## admin ##
  #
  admin:
    ## port ##
    #
    # Default value: 4445
    #
    # Set this value using environment variables on
    # - Linux/macOS:
    #    $ export SERVE_ADMIN_PORT=<value>
    # - Windows Command Line (CMD):
    #    > set SERVE_ADMIN_PORT=<value>
    #
    port: 12206

    ## host ##
    #
    # The interface or unix socket ORY Hydra should listen and handle administrative API requests on. Use the prefix "unix:" to specify a path to a unix socket. Leave empty to listen on all interfaces.
    #
    # Examples:
    # - localhost
    #
    # Set this value using environment variables on
    # - Linux/macOS:
    #    $ export SERVE_ADMIN_HOST=<value>
    # - Windows Command Line (CMD):
    #    > set SERVE_ADMIN_HOST=<value>
    #
    host: ''

    ## cors ##
    #
    # Configures Cross Origin Resource Sharing for public endpoints.
    #
    cors:
      ## enabled ##
      #
      # Sets whether CORS is enabled.
      #
      # Default value: false
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export SERVE_ADMIN_CORS_ENABLED=<value>
      # - Windows Command Line (CMD):
      #    > set SERVE_ADMIN_CORS_ENABLED=<value>
      #
      enabled: true

      ## allowed_origins ##
      #
      # A list of origins a cross-domain request can be executed from. If the special * value is present in the list, all origins will be allowed. An origin may contain a wildcard (*) to replace 0 or more characters (i.e.: http://*.domain.com). Only one wildcard can be used per origin.
      #
      # Default value: *
      #
      # Examples:
      # - - https://example.com
      #   - https://*.example.com
      #   - https://*.foo.example.com
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export SERVE_ADMIN_CORS_ALLOWED_ORIGINS=<value>
      # - Windows Command Line (CMD):
      #    > set SERVE_ADMIN_CORS_ALLOWED_ORIGINS=<value>
      #
      allowed_origins:
        - '*'

      ## allowed_methods ##
      #
      # A list of HTTP methods the user agent is allowed to use with cross-domain requests.
      #
      # Default value: POST,GET,PUT,PATCH,DELETE
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export SERVE_ADMIN_CORS_ALLOWED_METHODS=<value>
      # - Windows Command Line (CMD):
      #    > set SERVE_ADMIN_CORS_ALLOWED_METHODS=<value>
      #
      allowed_methods:
        - PUT
        - OPTIONS
        - CONNECT

      ## allowed_headers ##
      #
      # A list of non simple headers the client is allowed to use with cross-domain requests.
      #
      # Default value: Authorization,Content-Type
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export SERVE_ADMIN_CORS_ALLOWED_HEADERS=<value>
      # - Windows Command Line (CMD):
      #    > set SERVE_ADMIN_CORS_ALLOWED_HEADERS=<value>
      #
      allowed_headers:
        - nostrud eiusmod nisi
        - id
        - ea non fugiat ut irure
        - sit dolor commodo

      ## exposed_headers ##
      #
      # Sets which headers are safe to expose to the API of a CORS API specification.
      #
      # Default value: Content-Type
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export SERVE_ADMIN_CORS_EXPOSED_HEADERS=<value>
      # - Windows Command Line (CMD):
      #    > set SERVE_ADMIN_CORS_EXPOSED_HEADERS=<value>
      #
      exposed_headers:
        - adipisicing et ex nulla consequat
        - dolor deserunt
        - nisi incididunt
        - dolor dolore

      ## allow_credentials ##
      #
      # Sets whether the request can include user credentials like cookies, HTTP authentication or client side SSL certificates.
      #
      # Default value: true
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export SERVE_ADMIN_CORS_ALLOW_CREDENTIALS=<value>
      # - Windows Command Line (CMD):
      #    > set SERVE_ADMIN_CORS_ALLOW_CREDENTIALS=<value>
      #
      allow_credentials: true

      ## options_passthrough ##
      #
      # TODO
      #
      # Default value: false
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export SERVE_ADMIN_CORS_OPTIONS_PASSTHROUGH=<value>
      # - Windows Command Line (CMD):
      #    > set SERVE_ADMIN_CORS_OPTIONS_PASSTHROUGH=<value>
      #
      options_passthrough: false

      ## max_age ##
      #
      # Sets how long (in seconds) the results of a preflight request can be cached. If set to 0, every request is preceded by a preflight request.
      #
      # Minimum value: 0
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export SERVE_ADMIN_CORS_MAX_AGE=<value>
      # - Windows Command Line (CMD):
      #    > set SERVE_ADMIN_CORS_MAX_AGE=<value>
      #
      max_age: 51719476

      ## debug ##
      #
      # Adds additional log output to debug server side CORS issues.
      #
      # Default value: false
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export SERVE_ADMIN_CORS_DEBUG=<value>
      # - Windows Command Line (CMD):
      #    > set SERVE_ADMIN_CORS_DEBUG=<value>
      #
      debug: false

    ## socket ##
    #
    # Sets the permissions of the unix socket
    #
    socket:
      ## owner ##
      #
      # Owner of unix socket. If empty, the owner will be the user running hydra.
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export SERVE_ADMIN_SOCKET_OWNER=<value>
      # - Windows Command Line (CMD):
      #    > set SERVE_ADMIN_SOCKET_OWNER=<value>
      #
      owner: elit

      ## group ##
      #
      # Group of unix socket. If empty, the group will be the primary group of the user running hydra.
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export SERVE_ADMIN_SOCKET_GROUP=<value>
      # - Windows Command Line (CMD):
      #    > set SERVE_ADMIN_SOCKET_GROUP=<value>
      #
      group: aute sint ea adipisicing

      ## mode ##
      #
      # Mode of unix socket in numeric form
      #
      # Default value: 493
      #
      # Minimum value: 0
      #
      # Maximum value: 511
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export SERVE_ADMIN_SOCKET_MODE=<value>
      # - Windows Command Line (CMD):
      #    > set SERVE_ADMIN_SOCKET_MODE=<value>
      #
      mode: 322

    ## access_log ##
    #
    # Access Log configuration for admin server.
    #
    access_log:
      ## disable_for_health ##
      #
      # Disable access log for health endpoints.
      #
      # Default value: false
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export SERVE_ADMIN_ACCESS_LOG_DISABLE_FOR_HEALTH=<value>
      # - Windows Command Line (CMD):
      #    > set SERVE_ADMIN_ACCESS_LOG_DISABLE_FOR_HEALTH=<value>
      #
      disable_for_health: false

  ## tls ##
  #
  # Configures HTTPS (HTTP over TLS). If configured, the server automatically supports HTTP/2.
  #
  tls:
    ## key ##
    #
    # Configures the private key (pem encoded).
    #
    key:
      ## base64 ##
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export SERVE_TLS_KEY_BASE64=<value>
      # - Windows Command Line (CMD):
      #    > set SERVE_TLS_KEY_BASE64=<value>
      #
      base64: b3J5IGh5ZHJhIGlzIGF3ZXNvbWUK

    ## cert ##
    #
    # Configures the private key (pem encoded).
    #
    cert:
      ## path ##
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export SERVE_TLS_CERT_PATH=<value>
      # - Windows Command Line (CMD):
      #    > set SERVE_TLS_CERT_PATH=<value>
      #
      path: /path/to/file.pem

    ## allow_termination_from ##
    #
    # Whitelist one or multiple CIDR address ranges and allow them to terminate TLS connections. Be aware that the X-Forwarded-Proto header must be set and must never be modifiable by anyone but your proxy / gateway / load balancer. Supports ipv4 and ipv6. Hydra serves http instead of https when this option is set.
    #
    # Set this value using environment variables on
    # - Linux/macOS:
    #    $ export SERVE_TLS_ALLOW_TERMINATION_FROM=<value>
    # - Windows Command Line (CMD):
    #    > set SERVE_TLS_ALLOW_TERMINATION_FROM=<value>
    #
    allow_termination_from:
      - 127.0.0.1/32

  ## cookies ##
  #
  cookies:
    ## same_site_mode ##
    #
    # Specify the SameSite mode that cookies should be sent with.
    #
    # Default value: None
    #
    # One of:
    # - Strict
    # - Lax
    # - None
    #
    # Set this value using environment variables on
    # - Linux/macOS:
    #    $ export SERVE_COOKIES_SAME_SITE_MODE=<value>
    # - Windows Command Line (CMD):
    #    > set SERVE_COOKIES_SAME_SITE_MODE=<value>
    #
    same_site_mode: Lax

    ## same_site_legacy_workaround ##
    #
    # Some older browser versions don’t work with SameSite=None. This option enables the workaround defined in https://web.dev/samesite-cookie-recipes/ which essentially stores a second cookie without SameSite as a fallback.
    #
    # Default value: false
    #
    # Examples:
    # - true
    #
    # Set this value using environment variables on
    # - Linux/macOS:
    #    $ export SERVE_COOKIES_SAME_SITE_LEGACY_WORKAROUND=<value>
    # - Windows Command Line (CMD):
    #    > set SERVE_COOKIES_SAME_SITE_LEGACY_WORKAROUND=<value>
    #
    same_site_legacy_workaround: true

## dsn ##
#
# Sets the data source name. This configures the backend where ORY Hydra persists data. If dsn is "memory", data will be written to memory and is lost when you restart this instance. ORY Hydra supports popular SQL databases. For more detailed configuration information go to: https://www.ory.sh/docs/hydra/dependencies-environment#sql
#
# Set this value using environment variables on
# - Linux/macOS:
#    $ export DSN=<value>
# - Windows Command Line (CMD):
#    > set DSN=<value>
#
dsn: Excepteur labore laboris ut est

## webfinger ##
#
# Configures ./well-known/ settings.
#
webfinger:
  ## jwks ##
  #
  # Configures the /.well-known/jwks.json endpoint.
  #
  jwks:
    ## broadcast_keys ##
    #
    # A list of JSON Web Keys that should be exposed at that endpoint. This is usually the public key for verifying OpenID Connect ID Tokens. However, you might want to add additional keys here as well.
    #
    # Default value: hydra.openid.id-token
    #
    # Examples:
    # - hydra.jwt.access-token
    #
    # Set this value using environment variables on
    # - Linux/macOS:
    #    $ export WEBFINGER_JWKS_BROADCAST_KEYS=<value>
    # - Windows Command Line (CMD):
    #    > set WEBFINGER_JWKS_BROADCAST_KEYS=<value>
    #
    broadcast_keys: hydra.jwt.access-token

  ## oidc_discovery ##
  #
  # Configures OpenID Connect Discovery (/.well-known/openid-configuration).
  #
  oidc_discovery:
    ## client_registration_url ##
    #
    # Examples:
    # - https://my-service.com/clients
    #
    # Set this value using environment variables on
    # - Linux/macOS:
    #    $ export WEBFINGER_OIDC_DISCOVERY_CLIENT_REGISTRATION_URL=<value>
    # - Windows Command Line (CMD):
    #    > set WEBFINGER_OIDC_DISCOVERY_CLIENT_REGISTRATION_URL=<value>
    #
    client_registration_url: https://my-service.com/clients

    ## supported_claims ##
    #
    # A list of supported claims to be broadcasted. Claim "sub" is always included.
    #
    # Examples:
    # - - email
    #   - username
    #
    # Set this value using environment variables on
    # - Linux/macOS:
    #    $ export WEBFINGER_OIDC_DISCOVERY_SUPPORTED_CLAIMS=<value>
    # - Windows Command Line (CMD):
    #    > set WEBFINGER_OIDC_DISCOVERY_SUPPORTED_CLAIMS=<value>
    #
    supported_claims:
      - email
      - username

    ## supported_scope ##
    #
    # The scope OAuth 2.0 Clients may request. Scope `offline`, `offline_access`, and `openid` are always included.
    #
    # Examples:
    # - - email
    #   - whatever
    #   - read.photos
    #
    # Set this value using environment variables on
    # - Linux/macOS:
    #    $ export WEBFINGER_OIDC_DISCOVERY_SUPPORTED_SCOPE=<value>
    # - Windows Command Line (CMD):
    #    > set WEBFINGER_OIDC_DISCOVERY_SUPPORTED_SCOPE=<value>
    #
    supported_scope:
      - email
      - whatever
      - read.photos

    ## userinfo_url ##
    #
    # A URL of the userinfo endpoint to be advertised at the OpenID Connect Discovery endpoint /.well-known/openid-configuration. Defaults to ORY Hydra's userinfo endpoint at /userinfo. Set this value if you want to handle this endpoint yourself.
    #
    # Examples:
    # - https://example.org/my-custom-userinfo-endpoint
    #
    # Set this value using environment variables on
    # - Linux/macOS:
    #    $ export WEBFINGER_OIDC_DISCOVERY_USERINFO_URL=<value>
    # - Windows Command Line (CMD):
    #    > set WEBFINGER_OIDC_DISCOVERY_USERINFO_URL=<value>
    #
    userinfo_url: https://example.org/my-custom-userinfo-endpoint

## oidc ##
#
# Configures OpenID Connect features.
#
oidc:
  ## subject_identifiers ##
  #
  # Configures the Subject Identifier algorithm. For more information please head over to the documentation: https://www.ory.sh/docs/hydra/advanced#subject-identifier-algorithms
  #
  # Examples:
  # - enabled:
  #     - public
  #     - pairwise
  #   pairwise:
  #     salt: some-random-salt
  #
  subject_identifiers:
    ## enabled ##
    #
    # A list of algorithms to enable.
    #
    # Set this value using environment variables on
    # - Linux/macOS:
    #    $ export OIDC_SUBJECT_IDENTIFIERS_ENABLED=<value>
    # - Windows Command Line (CMD):
    #    > set OIDC_SUBJECT_IDENTIFIERS_ENABLED=<value>
    #
    enabled:
      - public
      - pairwise

    ## pairwise ##
    #
    # Configures the pairwise algorithm.
    #
    pairwise:
      ## salt ##
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export OIDC_SUBJECT_IDENTIFIERS_PAIRWISE_SALT=<value>
      # - Windows Command Line (CMD):
      #    > set OIDC_SUBJECT_IDENTIFIERS_PAIRWISE_SALT=<value>
      #
      salt: some-random-salt

  ## dynamic_client_registration ##
  #
  # Configures OpenID Connect Dynamic Client Registration (exposed as admin endpoints /clients/...).
  #
  dynamic_client_registration:
    ## default_scope ##
    #
    # The OpenID Connect Dynamic Client Registration specification has no concept of whitelisting OAuth 2.0 Scope. If you want to expose Dynamic Client Registration, you should set the default scope enabled for newly registered clients. Keep in mind that users can overwrite this default by setting the "scope" key in the registration payload, effectively disabling the concept of whitelisted scopes.
    #
    # Examples:
    # - - openid
    #   - offline
    #   - offline_access
    #
    # Set this value using environment variables on
    # - Linux/macOS:
    #    $ export OIDC_DYNAMIC_CLIENT_REGISTRATION_DEFAULT_SCOPE=<value>
    # - Windows Command Line (CMD):
    #    > set OIDC_DYNAMIC_CLIENT_REGISTRATION_DEFAULT_SCOPE=<value>
    #
    default_scope:
      - openid
      - offline
      - offline_access

## urls ##
#
urls:
  ## self ##
  #
  self:
    ## issuer ##
    #
    # This value will be used as the "issuer" in access and ID tokens. It must be specified and using HTTPS protocol, unless --dangerous-force-http is set. This should typically be equal to the public value.
    #
    # Examples:
    # - https://localhost:4444/
    #
    # Set this value using environment variables on
    # - Linux/macOS:
    #    $ export URLS_SELF_ISSUER=<value>
    # - Windows Command Line (CMD):
    #    > set URLS_SELF_ISSUER=<value>
    #
    issuer: https://localhost:4444/

    ## public ##
    #
    # This is the base location of the public endpoints of your ORY Hydra installation. This should typically be equal to the issuer value. If left unspecified, it falls back to the issuer value.
    #
    # Examples:
    # - https://localhost:4444/
    #
    # Set this value using environment variables on
    # - Linux/macOS:
    #    $ export URLS_SELF_PUBLIC=<value>
    # - Windows Command Line (CMD):
    #    > set URLS_SELF_PUBLIC=<value>
    #
    public: https://localhost:4444/

  ## login ##
  #
  # Sets the login endpoint of the User Login & Consent flow. Defaults to an internal fallback URL showing an error.
  #
  # Examples:
  # - https://my-login.app/login
  #
  # Set this value using environment variables on
  # - Linux/macOS:
  #    $ export URLS_LOGIN=<value>
  # - Windows Command Line (CMD):
  #    > set URLS_LOGIN=<value>
  #
  login: https://my-login.app/login

  ## consent ##
  #
  # Sets the consent endpoint of the User Login & Consent flow. Defaults to an internal fallback URL showing an error.
  #
  # Examples:
  # - https://my-consent.app/consent
  #
  # Set this value using environment variables on
  # - Linux/macOS:
  #    $ export URLS_CONSENT=<value>
  # - Windows Command Line (CMD):
  #    > set URLS_CONSENT=<value>
  #
  consent: https://my-consent.app/consent

  ## logout ##
  #
  # Sets the logout endpoint. Defaults to an internal fallback URL showing an error.
  #
  # Examples:
  # - https://my-logout.app/logout
  #
  # Set this value using environment variables on
  # - Linux/macOS:
  #    $ export URLS_LOGOUT=<value>
  # - Windows Command Line (CMD):
  #    > set URLS_LOGOUT=<value>
  #
  logout: https://my-logout.app/logout

  ## error ##
  #
  # Sets the error endpoint. The error ui will be shown when an OAuth2 error occurs that which can not be sent back to the client. Defaults to an internal fallback URL showing an error.
  #
  # Examples:
  # - https://my-error.app/error
  #
  # Set this value using environment variables on
  # - Linux/macOS:
  #    $ export URLS_ERROR=<value>
  # - Windows Command Line (CMD):
  #    > set URLS_ERROR=<value>
  #
  error: https://my-error.app/error

  ## post_logout_redirect ##
  #
  # When a user agent requests to logout, it will be redirected to this url afterwards per default.
  #
  # Examples:
  # - https://my-example.app/logout-successful
  #
  # Set this value using environment variables on
  # - Linux/macOS:
  #    $ export URLS_POST_LOGOUT_REDIRECT=<value>
  # - Windows Command Line (CMD):
  #    > set URLS_POST_LOGOUT_REDIRECT=<value>
  #
  post_logout_redirect: https://my-example.app/logout-successful

## strategies ##
#
strategies:
  ## scope ##
  #
  # Defines how scopes are matched. For more details have a look at https://github.com/ory/fosite#scopes
  #
  # Default value: wildcard
  #
  # One of:
  # - exact
  # - wildcard
  # - DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY
  #
  # Set this value using environment variables on
  # - Linux/macOS:
  #    $ export STRATEGIES_SCOPE=<value>
  # - Windows Command Line (CMD):
  #    > set STRATEGIES_SCOPE=<value>
  #
  scope: DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY

  ## access_token ##
  #
  # Defines access token type. jwt is a bad idea, see https://www.ory.sh/docs/hydra/advanced#json-web-tokens
  #
  # One of:
  # - opaque
  # - jwt
  #
  # Set this value using environment variables on
  # - Linux/macOS:
  #    $ export STRATEGIES_ACCESS_TOKEN=<value>
  # - Windows Command Line (CMD):
  #    > set STRATEGIES_ACCESS_TOKEN=<value>
  #
  access_token: opaque

## ttl ##
#
# Configures time to live.
#
ttl:
  ## login_consent_request ##
  #
  # Configures how long a user login and consent flow may take.
  #
  # Default value: 1h
  #
  # Set this value using environment variables on
  # - Linux/macOS:
  #    $ export TTL_LOGIN_CONSENT_REQUEST=<value>
  # - Windows Command Line (CMD):
  #    > set TTL_LOGIN_CONSENT_REQUEST=<value>
  #
  login_consent_request: 1h

  ## access_token ##
  #
  # Configures how long access tokens are valid.
  #
  # Default value: 1h
  #
  # Set this value using environment variables on
  # - Linux/macOS:
  #    $ export TTL_ACCESS_TOKEN=<value>
  # - Windows Command Line (CMD):
  #    > set TTL_ACCESS_TOKEN=<value>
  #
  access_token: 1h

  ## refresh_token ##
  #
  # Configures how long refresh tokens are valid. Set to -1 for refresh tokens to never expire.
  #
  # Default value: 720h
  #
  # Set this value using environment variables on
  # - Linux/macOS:
  #    $ export TTL_REFRESH_TOKEN=<value>
  # - Windows Command Line (CMD):
  #    > set TTL_REFRESH_TOKEN=<value>
  #
  refresh_token: -1

  ## id_token ##
  #
  # Configures how long id tokens are valid.
  #
  # Default value: 1h
  #
  # Set this value using environment variables on
  # - Linux/macOS:
  #    $ export TTL_ID_TOKEN=<value>
  # - Windows Command Line (CMD):
  #    > set TTL_ID_TOKEN=<value>
  #
  id_token: 1h

  ## auth_code ##
  #
  # Configures how long auth codes are valid.
  #
  # Default value: 10m
  #
  # Set this value using environment variables on
  # - Linux/macOS:
  #    $ export TTL_AUTH_CODE=<value>
  # - Windows Command Line (CMD):
  #    > set TTL_AUTH_CODE=<value>
  #
  auth_code: 1h

## oauth2 ##
#
oauth2:
  ## expose_internal_errors ##
  #
  # Set this to true if you want to share error debugging information with your OAuth 2.0 clients. Keep in mind that debug information is very valuable when dealing with errors, but might also expose database error codes and similar errors.
  #
  # Default value: false
  #
  # Examples:
  # - true
  #
  # Set this value using environment variables on
  # - Linux/macOS:
  #    $ export OAUTH2_EXPOSE_INTERNAL_ERRORS=<value>
  # - Windows Command Line (CMD):
  #    > set OAUTH2_EXPOSE_INTERNAL_ERRORS=<value>
  #
  expose_internal_errors: false

  ## hashers ##
  #
  # Configures hashing algorithms. Supports only BCrypt at the moment.
  #
  hashers:
    ## bcrypt ##
    #
    # Configures the BCrypt hashing algorithm used for hashing Client Secrets.
    #
    bcrypt:
      ## cost ##
      #
      # Sets the BCrypt cost. The higher the value, the more CPU time is being used to generate hashes.
      #
      # Default value: 10
      #
      # Minimum value: 8
      #
      # Maximum value: 31
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export OAUTH2_HASHERS_BCRYPT_COST=<value>
      # - Windows Command Line (CMD):
      #    > set OAUTH2_HASHERS_BCRYPT_COST=<value>
      #
      cost: 12

  ## pkce ##
  #
  pkce:
    ## enforced ##
    #
    # Sets whether PKCE should be enforced for all clients.
    #
    # Examples:
    # - true
    #
    # Set this value using environment variables on
    # - Linux/macOS:
    #    $ export OAUTH2_PKCE_ENFORCED=<value>
    # - Windows Command Line (CMD):
    #    > set OAUTH2_PKCE_ENFORCED=<value>
    #
    enforced: true

    ## enforced_for_public_clients ##
    #
    # Sets whether PKCE should be enforced for public clients.
    #
    # Examples:
    # - true
    #
    # Set this value using environment variables on
    # - Linux/macOS:
    #    $ export OAUTH2_PKCE_ENFORCED_FOR_PUBLIC_CLIENTS=<value>
    # - Windows Command Line (CMD):
    #    > set OAUTH2_PKCE_ENFORCED_FOR_PUBLIC_CLIENTS=<value>
    #
    enforced_for_public_clients: true

  ## client_credentials ##
  #
  client_credentials:
    ## default_grant_allowed_scope ##
    #
    # Defines how scopes are added if the request doesn't contains any scope
    #
    # Examples:
    # - false
    #
    # Set this value using environment variables on
    # - Linux/macOS:
    #    $ export OAUTH2_CLIENT_CREDENTIALS_DEFAULT_GRANT_ALLOWED_SCOPE=<value>
    # - Windows Command Line (CMD):
    #    > set OAUTH2_CLIENT_CREDENTIALS_DEFAULT_GRANT_ALLOWED_SCOPE=<value>
    #
    default_grant_allowed_scope: false

## secrets ##
#
# The secrets section configures secrets used for encryption and signing of several systems. All secrets can be rotated, for more information on this topic go to: https://www.ory.sh/docs/hydra/advanced#rotation-of-hmac-token-signing-and-database-and-cookie-encryption-keys
#
secrets:
  ## system ##
  #
  # The system secret must be at least 16 characters long. If none is provided, one will be generated. They key is used to encrypt sensitive data using AES-GCM (256 bit) and validate HMAC signatures. The first item in the list is used for signing and encryption. The whole list is used for verifying signatures and decryption.
  #
  # Examples:
  # - - this-is-the-primary-secret
  #   - this-is-an-old-secret
  #   - this-is-another-old-secret
  #
  # Set this value using environment variables on
  # - Linux/macOS:
  #    $ export SECRETS_SYSTEM=<value>
  # - Windows Command Line (CMD):
  #    > set SECRETS_SYSTEM=<value>
  #
  system:
    - this-is-the-primary-secret
    - this-is-an-old-secret
    - this-is-another-old-secret

  ## cookie ##
  #
  # A secret that is used to encrypt cookie sessions. Defaults to secrets.system. It is recommended to use a separate secret in production. The first item in the list is used for signing and encryption. The whole list is used for verifying signatures and decryption.
  #
  # Examples:
  # - - this-is-the-primary-secret
  #   - this-is-an-old-secret
  #   - this-is-another-old-secret
  #
  # Set this value using environment variables on
  # - Linux/macOS:
  #    $ export SECRETS_COOKIE=<value>
  # - Windows Command Line (CMD):
  #    > set SECRETS_COOKIE=<value>
  #
  cookie:
    - this-is-the-primary-secret
    - this-is-an-old-secret
    - this-is-another-old-secret

## profiling ##
#
# Enables profiling if set. For more details on profiling, head over to: https://blog.golang.org/profiling-go-programs
#
# One of:
# - cpu
# - mem
#
# Examples:
# - cpu
#
# Set this value using environment variables on
# - Linux/macOS:
#    $ export PROFILING=<value>
# - Windows Command Line (CMD):
#    > set PROFILING=<value>
#
profiling: cpu

## tracing ##
#
# ORY Hydra supports distributed tracing.
#
tracing:
  ## provider ##
  #
  # Set this to the tracing backend you wish to use. Supports Jaeger, Zipkin and DataDog. If omitted or empty, tracing will be disabled. Use environment variables to configure DataDog (see https://docs.datadoghq.com/tracing/setup/go/#configuration).
  #
  # One of:
  # - jaeger
  # - zipkin
  # - datadog
  # - elastic-apm
  #
  # Examples:
  # - jaeger
  #
  # Set this value using environment variables on
  # - Linux/macOS:
  #    $ export TRACING_PROVIDER=<value>
  # - Windows Command Line (CMD):
  #    > set TRACING_PROVIDER=<value>
  #
  provider: jaeger

  ## service_name ##
  #
  # Specifies the service name to use on the tracer.
  #
  # Examples:
  # - ORY Hydra
  #
  # Set this value using environment variables on
  # - Linux/macOS:
  #    $ export TRACING_SERVICE_NAME=<value>
  # - Windows Command Line (CMD):
  #    > set TRACING_SERVICE_NAME=<value>
  #
  service_name: ORY Hydra

  ## providers ##
  #
  providers:
    ## jaeger ##
    #
    # Configures the jaeger tracing backend.
    #
    jaeger:
      ## local_agent_address ##
      #
      # The address of the jaeger-agent where spans should be sent to.
      #
      # Examples:
      # - 127.0.0.1:6831
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export TRACING_PROVIDERS_JAEGER_LOCAL_AGENT_ADDRESS=<value>
      # - Windows Command Line (CMD):
      #    > set TRACING_PROVIDERS_JAEGER_LOCAL_AGENT_ADDRESS=<value>
      #
      local_agent_address: 127.0.0.1:6831

      ## propagation ##
      #
      # The tracing header format
      #
      # Examples:
      # - jaeger
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export TRACING_PROVIDERS_JAEGER_PROPAGATION=<value>
      # - Windows Command Line (CMD):
      #    > set TRACING_PROVIDERS_JAEGER_PROPAGATION=<value>
      #
      propagation: jaeger

      ## sampling ##
      #
      # Examples:
      # - type: const
      #   value: 1
      #   server_url: http://localhost:5778/sampling
      #
      sampling:
        ## type ##
        #
        # Set this value using environment variables on
        # - Linux/macOS:
        #    $ export TRACING_PROVIDERS_JAEGER_SAMPLING_TYPE=<value>
        # - Windows Command Line (CMD):
        #    > set TRACING_PROVIDERS_JAEGER_SAMPLING_TYPE=<value>
        #
        type: const

        ## value ##
        #
        # Set this value using environment variables on
        # - Linux/macOS:
        #    $ export TRACING_PROVIDERS_JAEGER_SAMPLING_VALUE=<value>
        # - Windows Command Line (CMD):
        #    > set TRACING_PROVIDERS_JAEGER_SAMPLING_VALUE=<value>
        #
        value: 1

        ## server_url ##
        #
        # Set this value using environment variables on
        # - Linux/macOS:
        #    $ export TRACING_PROVIDERS_JAEGER_SAMPLING_SERVER_URL=<value>
        # - Windows Command Line (CMD):
        #    > set TRACING_PROVIDERS_JAEGER_SAMPLING_SERVER_URL=<value>
        #
        server_url: http://localhost:5778/sampling

    ## zipkin ##
    #
    # Configures the zipkin tracing backend.
    #
    # Examples:
    # - server_url: http://localhost:9411/api/v2/spans
    #
    zipkin:
      ## server_url ##
      #
      # The address of Zipkin server where spans should be sent to.
      #
      # Set this value using environment variables on
      # - Linux/macOS:
      #    $ export TRACING_PROVIDERS_ZIPKIN_SERVER_URL=<value>
      # - Windows Command Line (CMD):
      #    > set TRACING_PROVIDERS_ZIPKIN_SERVER_URL=<value>
      #
      server_url: http://localhost:9411/api/v2/spans

## The Hydra version this config is written for. ##
#
# SemVer according to https://semver.org/ prefixed with `v` as in our releases.
#
# Set this value using environment variables on
# - Linux/macOS:
#    $ export VERSION=<value>
# - Windows Command Line (CMD):
#    > set VERSION=<value>
#
version: v0.0.933279065-0.43qhs.5497.0.027721lS.607306avZ3YQBW.22678684644.511350L.0+BD57mbvGUG8.i.ZpbWuGHH6.k2OT.aU3lZ.AMc

## cgroups ##
#
# ORY Hydra can respect Linux container CPU quota
#
cgroups:
  ## v1 ##
  #
  # Configures parameters using cgroups v1 hierarchy
  #
  v1:
    ## auto_max_procs_enabled ##
    #
    # Set GOMAXPROCS automatically according to cgroups limits
    #
    # Default value: false
    #
    # Examples:
    # - true
    #
    # Set this value using environment variables on
    # - Linux/macOS:
    #    $ export CGROUPS_V1_AUTO_MAX_PROCS_ENABLED=<value>
    # - Windows Command Line (CMD):
    #    > set CGROUPS_V1_AUTO_MAX_PROCS_ENABLED=<value>
    #
    auto_max_procs_enabled: true
```
