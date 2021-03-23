---
id: hydra-keys-create
title: hydra keys create
description: hydra keys create Create a new JSON Web Key Set
---

<!--
This file is auto-generated.

To improve this file please make your change against the appropriate "./cmd/*.go" file.
-->

## hydra keys create

Create a new JSON Web Key Set

```
hydra keys create <set> <key> [flags]
```

### Options

```
  -a, --alg string   The algorithm to be used to generated they key. Supports: RS256, ES512, HS256 (default "RS256")
  -h, --help         help for create
  -u, --use string   The intended use of this key (default "sig")
```

### Options inherited from parent commands

```
      --access-token string    Set an access token to be used in the Authorization header, defaults to environment variable OAUTH2_ACCESS_TOKEN
      --endpoint string        Set the URL where ORY Hydra is hosted, defaults to environment variable HYDRA_ADMIN_URL
      --fail-after duration    Stop retrying after the specified duration (default 1m0s)
      --fake-tls-termination   fake tls termination by adding "X-Forwarded-Proto: https" to http headers
      --skip-tls-verify        Foolishly accept TLS certificates signed by unknown certificate authorities
```

### SEE ALSO

- [hydra keys](hydra-keys) - Manage JSON Web Keys
