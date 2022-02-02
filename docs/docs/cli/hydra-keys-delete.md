---
id: hydra-keys-delete
title: hydra keys delete
description: hydra keys delete Delete a new JSON Web Key Set
---

<!--
This file is auto-generated.

To improve this file please make your change against the appropriate "./cmd/*.go" file.
-->

## hydra keys delete

Delete a new JSON Web Key Set

```
hydra keys delete &lt;set&gt; [flags]
```

### Options

```
  -h, --help   help for delete
```

### Options inherited from parent commands

```
      --access-token string    Set an access token to be used in the Authorization header, defaults to environment variable OAUTH2_ACCESS_TOKEN
      --endpoint string        Set the URL where Ory Hydra is hosted, defaults to environment variable HYDRA_ADMIN_URL
      --fail-after duration    Stop retrying after the specified duration (default 1m0s)
      --fake-tls-termination   fake tls termination by adding &#34;X-Forwarded-Proto: https&#34; to http headers
      --skip-tls-verify        Foolishly accept TLS certificates signed by unknown certificate authorities
```

### SEE ALSO

- [hydra keys](hydra-keys) - Manage JSON Web Keys
