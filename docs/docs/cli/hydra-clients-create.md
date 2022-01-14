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

Ory Hydra implements the OpenID Connect Dynamic Client registration
specification. Most flags are supported by this command as well.

Example: hydra clients create -n &#34;my app&#34; -c http://localhost/cb -g
authorization_code -r code -a core,foobar

To encrypt auto generated client secret, use &#34;--pgp-key&#34;,
&#34;--pgp-key-url&#34; or &#34;--keybase&#34; flag, for example: hydra clients
create -n &#34;my app&#34; -g client_credentials -r token -a core,foobar
--keybase keybase_username

```
hydra clients create [flags]
```

### Options

```
      --allowed-cors-origins strings           The list of URLs allowed to make CORS requests. Requires CORS_ENABLED.
      --audience strings                       The audience this client is allowed to request
      --backchannel-logout-callback string     Client URL that will cause the client to log itself out when sent a Logout Token by Hydra.
      --backchannel-logout-session-required    Boolean flag specifying whether the client requires that a sid (session ID) Claim be included in the Logout Token to identify the client session with the OP when the backchannel-logout-callback is used. If omitted, the default value is false.
  -c, --callbacks strings                      REQUIRED list of allowed callback URLs
      --client-uri string                      A URL string of a web page providing information about the client
      --frontchannel-logout-callback string    Client URL that will cause the client to log itself out when rendered in an iframe by Hydra.
      --frontchannel-logout-session-required   Boolean flag specifying whether the client requires that a sid (session ID) Claim be included in the Logout Token to identify the client session with the OP when the frontchannel-logout-callback is used. If omitted, the default value is false.
  -g, --grant-types strings                    A list of allowed grant types (default [authorization_code])
  -h, --help                                   help for create
      --id string                              Give the client this id
      --jwks-uri string                        Define the URL where the JSON Web Key Set should be fetched from when performing the &#34;private_key_jwt&#34; client authentication method
      --keybase string                         Keybase username for encrypting client secret
      --logo-uri string                        A URL string that references a logo for the client
  -n, --name string                            The client&#39;s name
      --pgp-key string                         Base64 encoded PGP encryption key for encrypting client secret
      --pgp-key-url string                     PGP encryption key URL for encrypting client secret
      --policy-uri string                      A URL string that points to a human-readable privacy policy document that describes how the deployment organization collects, uses, retains, and discloses personal data
      --post-logout-callbacks strings          List of allowed URLs to be redirected to after a logout
  -r, --response-types strings                 A list of allowed response types (default [code])
  -a, --scope strings                          The scope the client is allowed to request
      --secret string                          Provide the client&#39;s secret
      --subject-type string                    A identifier algorithm. Valid values are &#34;public&#34; and &#34;pairwise&#34; (default &#34;public&#34;)
      --token-endpoint-auth-method string      Define which authentication method the client may use at the Token Endpoint. Valid values are &#34;client_secret_post&#34;, &#34;client_secret_basic&#34;, &#34;private_key_jwt&#34;, and &#34;none&#34; (default &#34;client_secret_basic&#34;)
      --tos-uri string                         A URL string that points to a human-readable terms of service document for the client that describes a contractual relationship between the end-user and the client that the end-user accepts when authorizing the client
```

### Options inherited from parent commands

```
      --access-token string    Set an access token to be used in the Authorization header, defaults to environment variable OAUTH2_ACCESS_TOKEN
      --endpoint string        Set the URL where Ory Hydra is hosted, defaults to environment variable HYDRA_ADMIN_URL. A unix socket can be set in the form unix:///path/to/socket
      --fail-after duration    Stop retrying after the specified duration (default 1m0s)
      --fake-tls-termination   Fake tls termination by adding &#34;X-Forwarded-Proto: https&#34; to http headers
      --skip-tls-verify        Foolishly accept TLS certificates signed by unknown certificate authorities
```

### SEE ALSO

- [hydra clients](hydra-clients) - Manage OAuth 2.0 Clients
