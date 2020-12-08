---
id: hydra-token-user
title: hydra token user
description:
  hydra token user An exemplary OAuth 2.0 Client performing the OAuth 2.0
  Authorize Code Flow
---

<!--
This file is auto-generated.

To improve this file please make your change against the appropriate "./cmd/*.go" file.
-->

## hydra token user

An exemplary OAuth 2.0 Client performing the OAuth 2.0 Authorize Code Flow

### Synopsis

Starts an exemplary web server that acts as an OAuth 2.0 Client performing the
Authorize Code Flow. This command will help you to see if ORY Hydra has been
configured properly.

This command must not be used for anything else than manual testing or demo
purposes. The server will terminate on error and success.

```
hydra token user [flags]
```

### Options

```
      --audience strings       Request a specific OAuth 2.0 Access Token Audience
      --auth-url endpoint      Usually it is enough to specify the endpoint flag, but if you want to force the authorization url, use this flag
      --client-id string       Use the provided OAuth 2.0 Client ID, defaults to environment variable OAUTH2_CLIENT_ID
      --client-secret string   Use the provided OAuth 2.0 Client Secret, defaults to environment variable OAUTH2_CLIENT_SECRET
      --endpoint string        Set the URL where ORY Hydra is hosted, defaults to environment variable HYDRA_URL
  -h, --help                   help for user
      --https                  Sets up HTTPS for the endpoint using a self-signed certificate which is re-generated every time you start this command
      --max-age int            Set the OpenID Connect max_age parameter
      --no-open                Do not open the browser window automatically
  -p, --port int               The port on which the server should run (default 4446)
      --prompt strings         Set the OpenID Connect prompt parameter
      --redirect string        Force a redirect url
      --scope strings          Request OAuth2 scope (default [offline,openid])
      --token-url endpoint     Usually it is enough to specify the endpoint flag, but if you want to force the token url, use this flag
```

### Options inherited from parent commands

```
      --fail-after duration    Stop retrying after the specified duration (default 1m0s)
      --fake-tls-termination   fake tls termination by adding "X-Forwarded-Proto: https" to http headers
      --skip-tls-verify        Foolishly accept TLS certificates signed by unknown certificate authorities
```

### SEE ALSO

- [hydra token](hydra-token) - Issue and Manage OAuth2 tokens
