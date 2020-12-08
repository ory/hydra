---
id: hydra-clients-create
title: hydra clients create
description: hydra clients create Create a new OAuth 2.0 Client
---

<!--
This file is auto-generated.

To improve this file please make your change against the appropriate "./cmd/*.go" file.
-->

## hydra clients create

Create a new OAuth 2.0 Client

### Synopsis

This command creates an OAuth 2.0 Client which can be used to perform various
OAuth 2.0 Flows like the Authorize Code, Implicit, Refresh flow.

ORY Hydra implements the OpenID Connect Dynamic Client registration
specification. Most flags are supported by this command as well.

Example: hydra clients create -n "my app" -c http://localhost/cb -g
authorization_code -r code -a core,foobar

To encrypt auto generated client secret, use "--pgp-key", "--pgp-key-url" or
"--keybase" flag, for example: hydra clients create -n "my app" -g
client_credentials -r token -a core,foobar --keybase keybase_username

```
hydra clients create [flags]
```

### Options

```
      --allowed-cors-origins strings        The list of URLs allowed to make CORS requests. Requires CORS_ENABLED.
      --audience strings                    The audience this client is allowed to request
  -c, --callbacks strings                   REQUIRED list of allowed callback URLs
      --client-uri string                   A URL string of a web page providing information about the client
  -g, --grant-types strings                 A list of allowed grant types (default [authorization_code])
  -h, --help                                help for create
      --id string                           Give the client this id
      --jwks-uri string                     Define the URL where the JSON Web Key Set should be fetched from when performing the "private_key_jwt" client authentication method
      --keybase string                      Keybase username for encrypting client secret
      --logo-uri string                     A URL string that references a logo for the client
  -n, --name string                         The client's name
      --pgp-key string                      Base64 encoded PGP encryption key for encrypting client secret
      --pgp-key-url string                  PGP encryption key URL for encrypting client secret
      --policy-uri string                   A URL string that points to a human-readable privacy policy document that describes how the deployment organization collects, uses, retains, and discloses personal data
      --post-logout-callbacks strings       List of allowed URLs to be redirected to after a logout
  -r, --response-types strings              A list of allowed response types (default [code])
  -a, --scope strings                       The scope the client is allowed to request
      --secret string                       Provide the client's secret
      --subject-type string                 A identifier algorithm. Valid values are "public" and "pairwise" (default "public")
      --token-endpoint-auth-method string   Define which authentication method the client may use at the Token Endpoint. Valid values are "client_secret_post", "client_secret_basic", "private_key_jwt", and "none" (default "client_secret_basic")
      --tos-uri string                      A URL string that points to a human-readable terms of service document for the client that describes a contractual relationship between the end-user and the client that the end-user accepts when authorizing the client
```

### Options inherited from parent commands

```
      --access-token string    Set an access token to be used in the Authorization header, defaults to environment variable OAUTH2_ACCESS_TOKEN
      --endpoint string        Set the URL where ORY Hydra is hosted, defaults to environment variable HYDRA_ADMIN_URL. A unix socket can be set in the form unix:///path/to/socket
      --fail-after duration    Stop retrying after the specified duration (default 1m0s)
      --fake-tls-termination   Fake tls termination by adding "X-Forwarded-Proto: https" to http headers
      --skip-tls-verify        Foolishly accept TLS certificates signed by unknown certificate authorities
```

### SEE ALSO

- [hydra clients](hydra-clients) - Manage OAuth 2.0 Clients
