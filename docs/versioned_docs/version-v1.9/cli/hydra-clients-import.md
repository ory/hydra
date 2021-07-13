---
id: hydra-clients-import
title: hydra clients import
description:
  hydra clients import Import OAuth 2.0 Clients from one or more JSON files
---

<!--
This file is auto-generated.

To improve this file please make your change against the appropriate "./cmd/*.go" file.
-->

## hydra clients import

Import OAuth 2.0 Clients from one or more JSON files

### Synopsis

This command reads in each listed JSON file and imports their contents as OAuth
2.0 Clients.

The format for the JSON file is:

{ "client_id": "...", "client_secret": "...", // ... all other fields of the
OAuth 2.0 Client model are allowed here }

Please be aware that this command does not update existing clients. If the
client exists already, this command will fail.

Example: hydra clients import client-1.json

To encrypt auto generated client secret, use "--pgp-key", "--pgp-key-url" or
"--keybase" flag, for example: hydra clients import client-1.json --keybase
keybase_username

```
hydra clients import <path/to/file.json> [<path/to/other/file.json>...] [flags]
```

### Options

```
  -h, --help                 help for import
      --keybase string       Keybase username for encrypting client secret
      --pgp-key string       Base64 encoded PGP encryption key for encrypting client secret
      --pgp-key-url string   PGP encryption key URL for encrypting client secret
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
