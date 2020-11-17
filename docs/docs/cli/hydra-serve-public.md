---
id: hydra-serve-public
title: hydra serve public
description: hydra serve public Serves Public HTTP/2 APIs
---

<!--
This file is auto-generated.

To improve this file please make your change against the appropriate "./cmd/*.go" file.
-->

## hydra serve public

Serves Public HTTP/2 APIs

### Synopsis

This command opens one port and listens to HTTP/2 API requests. The exposed API
handles requests coming from the public internet, like OAuth 2.0 Authorization
and Token requests, OpenID Connect UserInfo, OAuth 2.0 Token Revokation, and
OpenID Connect Discovery.

This command is configurable using the same options available to "serve admin"
and "serve all".

It is generally recommended to use this command only if you require granular
control over the privileged and public APIs. For example, you might want to run
different TLS certificates or CORS settings on the public and privileged API.

This command does not work with the "memory" database. Both services
(privileged, public) MUST use the same database connection to be able to
synchronize.

## Configuration

ORY Hydra can be configured using environment variables as well as a
configuration file. For more information on configuration options, open the
configuration documentation:

> > https://github.com/ory/hydra/blob/undefined/docs/docs/reference/configuration.md
> > <<

```
hydra serve public [flags]
```

### Options

```
  -h, --help   help for public
```

### Options inherited from parent commands

```
  -c, --config string                                    Config file (default is $HOME/hydra.yaml)
      --dangerous-allow-insecure-redirect-urls strings   DO NOT USE THIS IN PRODUCTION - Disable HTTPS enforcement for the provided redirect URLs
      --dangerous-force-http                             DO NOT USE THIS IN PRODUCTION - Disables HTTP/2 over TLS (HTTPS) and serves HTTP instead
      --disable-telemetry                                Disable anonymized telemetry reports - for more information please visit https://www.ory.sh/docs/ecosystem/sqa
      --sqa-opt-out                                      Disable anonymized telemetry reports - for more information please visit https://www.ory.sh/docs/ecosystem/sqa
```

### SEE ALSO

- [hydra serve](hydra-serve) - Parent command for starting public and
  administrative HTTP/2 APIs
