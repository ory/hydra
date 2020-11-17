---
id: hydra-clients-list
title: hydra clients list
description: hydra clients list List OAuth 2.0 Clients
---

<!--
This file is auto-generated.

To improve this file please make your change against the appropriate "./cmd/*.go" file.
-->

## hydra clients list

List OAuth 2.0 Clients

### Synopsis

This command list an OAuth 2.0 Clients.

Example: hydra clients list

```
hydra clients list [flags]
```

### Options

```
  -h, --help        help for list
      --limit int   The maximum amount of policies returned. (default 20)
      --page int    The number of page. (default 1)
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
