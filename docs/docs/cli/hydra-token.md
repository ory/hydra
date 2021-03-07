---
id: hydra-token
title: hydra token
description: hydra token Issue and Manage OAuth2 tokens
---

<!--
This file is auto-generated.

To improve this file please make your change against the appropriate "./cmd/*.go" file.
-->

## hydra token

Issue and Manage OAuth2 tokens

### Synopsis

Issue and Manage OAuth2 tokens

### Options

```
      --fail-after duration    Stop retrying after the specified duration (default 1m0s)
      --fake-tls-termination   fake tls termination by adding "X-Forwarded-Proto: https" to http headers
  -h, --help                   help for token
      --skip-tls-verify        Foolishly accept TLS certificates signed by unknown certificate authorities
```

### SEE ALSO

- [hydra](hydra) - Run and manage ORY Hydra
- [hydra token client](hydra-token-client) - An exemplary OAuth 2.0 Client
  performing the OAuth 2.0 Client Credentials Flow
- [hydra token delete](hydra-token-delete) - Deletes access tokens of a client
- [hydra token flush](hydra-token-flush) - Removes inactive access tokens from
  the database
- [hydra token introspect](hydra-token-introspect) - Introspect an access or
  refresh token
- [hydra token revoke](hydra-token-revoke) - Revoke an access or refresh token
- [hydra token user](hydra-token-user) - An exemplary OAuth 2.0 Client
  performing the OAuth 2.0 Authorize Code Flow
