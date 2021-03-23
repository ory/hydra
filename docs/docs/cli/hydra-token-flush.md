---
id: hydra-token-flush
title: hydra token flush
description: hydra token flush Removes inactive access tokens from the database
---

<!--
This file is auto-generated.

To improve this file please make your change against the appropriate "./cmd/*.go" file.
-->

## hydra token flush

Removes inactive access tokens from the database

```
hydra token flush [flags]
```

### Options

```
      --access-token string   Set an access token to be used in the Authorization header, defaults to environment variable OAUTH2_ACCESS_TOKEN
      --endpoint string       Set the URL where ORY Hydra is hosted, defaults to environment variable HYDRA_ADMIN_URL
  -h, --help                  help for flush
      --min-age duration      Skip removing tokens which do not satisfy the minimum age (1s, 1m, 1h)
```

### Options inherited from parent commands

```
      --fail-after duration    Stop retrying after the specified duration (default 1m0s)
      --fake-tls-termination   fake tls termination by adding "X-Forwarded-Proto: https" to http headers
      --skip-tls-verify        Foolishly accept TLS certificates signed by unknown certificate authorities
```

### SEE ALSO

- [hydra token](hydra-token) - Issue and Manage OAuth2 tokens
