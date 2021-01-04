---
id: hydra-clients
title: hydra clients
description: hydra clients Manage OAuth 2.0 Clients
---

<!--
This file is auto-generated.

To improve this file please make your change against the appropriate "./cmd/*.go" file.
-->

## hydra clients

Manage OAuth 2.0 Clients

### Synopsis

Manage OAuth 2.0 Clients

### Options

```
      --access-token string    Set an access token to be used in the Authorization header, defaults to environment variable OAUTH2_ACCESS_TOKEN
      --endpoint string        Set the URL where ORY Hydra is hosted, defaults to environment variable HYDRA_ADMIN_URL. A unix socket can be set in the form unix:///path/to/socket
      --fail-after duration    Stop retrying after the specified duration (default 1m0s)
      --fake-tls-termination   Fake tls termination by adding "X-Forwarded-Proto: https" to http headers
  -h, --help                   help for clients
      --skip-tls-verify        Foolishly accept TLS certificates signed by unknown certificate authorities
```

### SEE ALSO

- [hydra](hydra) - Run and manage ORY Hydra
- [hydra clients create](hydra-clients-create) - Create a new OAuth 2.0 Client
- [hydra clients delete](hydra-clients-delete) - Delete an OAuth 2.0 Client
- [hydra clients get](hydra-clients-get) - Get an OAuth 2.0 Client
- [hydra clients import](hydra-clients-import) - Import OAuth 2.0 Clients from
  one or more JSON files
- [hydra clients list](hydra-clients-list) - List OAuth 2.0 Clients
- [hydra clients update](hydra-clients-update) - Update an entire OAuth 2.0
  Client
