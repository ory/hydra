---
id: hydra-keys-import
title: hydra keys import
description:
  hydra keys import Imports cryptographic keys of any format to the JSON Web Key
  Store
---

<!--
This file is auto-generated.

To improve this file please make your change against the appropriate "./cmd/*.go" file.
-->

## hydra keys import

Imports cryptographic keys of any format to the JSON Web Key Store

### Synopsis

This command allows you to import cryptographic keys to the JSON Web Key Store.

Currently supported formats are raw JSON Web Keys or PEM/DER encoded data. If
the JSON Web Key Set exists already, the imported keys will be added to that
set. Otherwise, a new set will be created.

Please be aware that importing a private key does not automatically import its
public key as well.

Examples: hydra keys import my-set ./path/to/jwk.json ./path/to/jwk-2.json hydra
keys import my-set ./path/to/rsa.key ./path/to/rsa.pub

```
hydra keys import <set> <file-1> [<file-2> [<file-3 [<...>]]] [flags]
```

### Options

```
  -h, --help         help for import
      --use string   Sets the "use" value of the JSON Web Key if not "use" value was defined by the key itself (default "sig")
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
