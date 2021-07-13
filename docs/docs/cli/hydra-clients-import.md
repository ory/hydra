---
id: hydra-clients-import
title: hydra clients import
description:
  hydra clients import Imports cryptographic keys of any format to the JSON Web
  Key Store
---

<!--
This file is auto-generated.

To improve this file please make your change against the appropriate "./cmd/*.go" file.
-->

## hydra clients import

Imports cryptographic keys of any format to the JSON Web Key Store

### Synopsis

This command allows you to import cryptographic keys to the JSON Web Key Store.

Currently supported formats are raw JSON Web Keys or PEM/DER encoded data. If
the JSON Web Key Set exists already, the imported keys will be added to that
set. Otherwise, a new set will be created.

Please be aware that importing a private key does not automatically import its
public key as well.

Examples: hydra keys import my-set ./path/to/jwk.json ./path/to/jwk-2.json hydra
keys import my-set ./path/to/rsa.key ./path/to/rsa.pub --default-key-id
cae6b214-fb1e-4ebc-9019-95286a62eabc

```
hydra clients import &lt;set&gt; &lt;file-1&gt; [&lt;file-2&gt; [&lt;file-3 [&lt;...&gt;]]] [flags]
```

### Options

```
      --default-key-id string   A fallback value for keys without &#34;kid&#34; attribute to be stored with a common &#34;kid&#34;, e.g. private/public keypairs
  -h, --help                    help for import
      --use string              Sets the &#34;use&#34; value of the JSON Web Key if not &#34;use&#34; value was defined by the key itself (default &#34;sig&#34;)
```

### Options inherited from parent commands

```
      --access-token string    Set an access token to be used in the Authorization header, defaults to environment variable OAUTH2_ACCESS_TOKEN
      --endpoint string        Set the URL where ORY Hydra is hosted, defaults to environment variable HYDRA_ADMIN_URL. A unix socket can be set in the form unix:///path/to/socket
      --fail-after duration    Stop retrying after the specified duration (default 1m0s)
      --fake-tls-termination   Fake tls termination by adding &#34;X-Forwarded-Proto: https&#34; to http headers
      --skip-tls-verify        Foolishly accept TLS certificates signed by unknown certificate authorities
```

### SEE ALSO

- [hydra clients](hydra-clients) - Manage OAuth 2.0 Clients
