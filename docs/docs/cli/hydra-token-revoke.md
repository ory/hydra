---
id: hydra-token-revoke
title: hydra token revoke
description: hydra token revoke Revoke an access or refresh token
---

<!--
This file is auto-generated.

To improve this file please make your change against the appropriate "./cmd/*.go" file.
-->

## hydra token revoke

Revoke an access or refresh token

```
hydra token revoke &lt;token&gt; [flags]
```

### Options

```
      --client-id string       Use the provided OAuth 2.0 Client ID, defaults to environment variable OAUTH2_CLIENT_ID
      --client-secret string   Use the provided OAuth 2.0 Client Secret, defaults to environment variable OAUTH2_CLIENT_SECRET
      --endpoint string        Set the URL where Ory Hydra is hosted, defaults to environment variable HYDRA_URL
  -h, --help                   help for revoke
```

### Options inherited from parent commands

```
      --fail-after duration    Stop retrying after the specified duration (default 1m0s)
      --fake-tls-termination   fake tls termination by adding &#34;X-Forwarded-Proto: https&#34; to http headers
      --skip-tls-verify        Foolishly accept TLS certificates signed by unknown certificate authorities
```

### SEE ALSO

- [hydra token](hydra-token) - Issue and Manage OAuth2 tokens
