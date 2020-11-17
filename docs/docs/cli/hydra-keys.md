---
id: hydra-keys
title: hydra keys
description: hydra keys Manage JSON Web Keys
---

<!--
This file is auto-generated.

To improve this file please make your change against the appropriate "./cmd/*.go" file.
-->

## hydra keys

Manage JSON Web Keys

### Synopsis

Manage JSON Web Keys

### Options

```
      --access-token string    Set an access token to be used in the Authorization header, defaults to environment variable OAUTH2_ACCESS_TOKEN
      --endpoint string        Set the URL where ORY Hydra is hosted, defaults to environment variable HYDRA_ADMIN_URL
      --fail-after duration    Stop retrying after the specified duration (default 1m0s)
      --fake-tls-termination   fake tls termination by adding "X-Forwarded-Proto: https" to http headers
  -h, --help                   help for keys
      --skip-tls-verify        Foolishly accept TLS certificates signed by unknown certificate authorities
```

### SEE ALSO

- [hydra](hydra) - Run and manage ORY Hydra
- [hydra keys create](hydra-keys-create) - Create a new JSON Web Key Set
- [hydra keys delete](hydra-keys-delete) - Delete a new JSON Web Key Set
- [hydra keys get](hydra-keys-get) - Get a new JSON Web Key Set
- [hydra keys import](hydra-keys-import) - Imports cryptographic keys of any
  format to the JSON Web Key Store
