---
id: hydra-token-client
title: hydra token client
description:
  hydra token client An exemplary OAuth 2.0 Client performing the OAuth 2.0
  Client Credentials Flow
---

<!--
This file is auto-generated.

To improve this file please make your change against the appropriate "./cmd/*.go" file.
-->

## hydra token client

An exemplary OAuth 2.0 Client performing the OAuth 2.0 Client Credentials Flow

### Synopsis

Performs the OAuth 2.0 Client Credentials Flow. This command will help you to
see if ORY Hydra has been configured properly.

This command should not be used for anything else than manual testing or demo
purposes. The server will terminate on error and success.

```
hydra token client [flags]
```

### Options

```
      --audience strings       Request a specific OAuth 2.0 Access Token Audience
      --client-id string       Use the provided OAuth 2.0 Client ID, defaults to environment variable OAUTH2_CLIENT_ID
      --client-secret string   Use the provided OAuth 2.0 Client Secret, defaults to environment variable OAUTH2_CLIENT_SECRET
      --endpoint string        Set the URL where ORY Hydra is hosted, defaults to environment variable HYDRA_URL
  -h, --help                   help for client
      --scope strings          OAuth2 scope to request
  -v, --verbose                Toggle verbose output mode
```

### Options inherited from parent commands

```
      --fail-after duration    Stop retrying after the specified duration (default 1m0s)
      --fake-tls-termination   fake tls termination by adding "X-Forwarded-Proto: https" to http headers
      --skip-tls-verify        Foolishly accept TLS certificates signed by unknown certificate authorities
```

### SEE ALSO

- [hydra token](hydra-token) - Issue and Manage OAuth2 tokens
