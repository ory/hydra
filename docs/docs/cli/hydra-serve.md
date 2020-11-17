---
id: hydra-serve
title: hydra serve
description:
  hydra serve Parent command for starting public and administrative HTTP/2 APIs
---

<!--
This file is auto-generated.

To improve this file please make your change against the appropriate "./cmd/*.go" file.
-->

## hydra serve

Parent command for starting public and administrative HTTP/2 APIs

### Synopsis

ORY Hydra exposes two ports, a public and an administrative port. The public
port is responsible for handling requests from the public internet, such as the
OAuth 2.0 Authorize and Token URLs. The administrative port handles
administrative requests like creating OAuth 2.0 Clients, managing JSON Web Keys,
and managing User Login and Consent sessions.

It is recommended to run "hydra serve all". If you need granular control over
CORS settings or similar, you may want to run "hydra serve admin" and "hydra
serve public" separately.

To learn more about each individual command, run:

- hydra help serve all
- hydra help serve admin
- hydra help serve public

All sub-commands share command line flags and configuration options.

## Configuration

ORY Hydra can be configured using environment variables as well as a
configuration file. For more information on configuration options, open the
configuration documentation:

> > https://github.com/ory/hydra/blob/undefined/docs/docs/reference/configuration.md
> > <<

### Options

```
  -c, --config string                                    Config file (default is $HOME/hydra.yaml)
      --dangerous-allow-insecure-redirect-urls strings   DO NOT USE THIS IN PRODUCTION - Disable HTTPS enforcement for the provided redirect URLs
      --dangerous-force-http                             DO NOT USE THIS IN PRODUCTION - Disables HTTP/2 over TLS (HTTPS) and serves HTTP instead
      --disable-telemetry                                Disable anonymized telemetry reports - for more information please visit https://www.ory.sh/docs/ecosystem/sqa
  -h, --help                                             help for serve
      --sqa-opt-out                                      Disable anonymized telemetry reports - for more information please visit https://www.ory.sh/docs/ecosystem/sqa
```

### SEE ALSO

- [hydra](hydra) - Run and manage ORY Hydra
- [hydra serve admin](hydra-serve-admin) - Serves Administrative HTTP/2 APIs
- [hydra serve all](hydra-serve-all) - Serves both public and administrative
  HTTP/2 APIs
- [hydra serve public](hydra-serve-public) - Serves Public HTTP/2 APIs
