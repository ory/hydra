---
id: hydra-serve-all
title: hydra serve all
description: hydra serve all Serves both public and administrative HTTP/2 APIs
---

<!--
This file is auto-generated.

To improve this file please make your change against the appropriate "./cmd/*.go" file.
-->

## hydra serve all

Serves both public and administrative HTTP/2 APIs

### Synopsis

Starts a process which listens on two ports for public and administrative HTTP/2
API requests.

If you want more granular control (e.g. different TLS settings) over each API
group (administrative, public) you can run &#34;serve admin&#34; and &#34;serve
public&#34; separately.

This command exposes a variety of controls via environment variables. You can
set environments using &#34;export KEY=VALUE&#34; (Linux/macOS) or &#34;set
KEY=VALUE&#34; (Windows). On Linux, you can also set environments by prepending
key value pairs: &#34;KEY=VALUE KEY2=VALUE2 hydra&#34;

All possible controls are listed below. This command exposes exposes command
line flags, which are listed below the controls section.

## Configuration

ORY Hydra can be configured using environment variables as well as a
configuration file. For more information on configuration options, open the
configuration documentation:

&gt;&gt; https://www.ory.sh/hydra/docs/reference/configuration &lt;&lt;

```
hydra serve all [flags]
```

### Options

```
  -h, --help   help for all
```

### Options inherited from parent commands

```
  -c, --config strings                                   Path to one or more .json, .yaml, .yml, .toml config files. Values are loaded in the order provided, meaning that the last config file overwrites values from the previous config file.
      --dangerous-allow-insecure-redirect-urls strings   DO NOT USE THIS IN PRODUCTION - Disable HTTPS enforcement for the provided redirect URLs
      --dangerous-force-http                             DO NOT USE THIS IN PRODUCTION - Disables HTTP/2 over TLS (HTTPS) and serves HTTP instead
      --sqa-opt-out                                      Disable anonymized telemetry reports - for more information please visit https://www.ory.sh/docs/ecosystem/sqa
```

### SEE ALSO

- [hydra serve](hydra-serve) - Parent command for starting public and
  administrative HTTP/2 APIs
