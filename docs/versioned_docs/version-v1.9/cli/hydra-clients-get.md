---
id: hydra-clients-get
title: hydra clients get
description: hydra clients get Get an OAuth 2.0 Client
---

<!--
This file is auto-generated.

To improve this file please make your change against the appropriate "./cmd/*.go" file.
-->

## hydra clients get

Get an OAuth 2.0 Client

### Synopsis

This command retrieves an OAuth 2.0 Clients by its ID.

Example: hydra clients get client-1

```
hydra clients get <id> [flags]
```

### Options

```
  -h, --help   help for get
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
