---
id: hydra-token-delete
title: hydra token delete
description: hydra token delete Deletes access tokens of a client
---

<!--
This file is auto-generated.

To improve this file please make your change against the appropriate "./cmd/*.go" file.
-->

## hydra token delete

Deletes access tokens of a client

```
hydra token delete [flags]
```

### Options

```
      --access-token string   Set an access token to be used in the Authorization header, defaults to environment variable OAUTH2_ACCESS_TOKEN
      --client-id string      Use the provided OAuth 2.0 Client ID, defaults to environment variable OAUTH2_CLIENT_ID
      --endpoint string       Set the URL where Ory Hydra is hosted, defaults to environment variable HYDRA_URL
  -h, --help                  help for delete
```

### Options inherited from parent commands

```
      --fail-after duration    Stop retrying after the specified duration (default 1m0s)
      --fake-tls-termination   fake tls termination by adding &#34;X-Forwarded-Proto: https&#34; to http headers
      --skip-tls-verify        Foolishly accept TLS certificates signed by unknown certificate authorities
```

### SEE ALSO

- [hydra token](hydra-token) - Issue and Manage OAuth2 tokens
