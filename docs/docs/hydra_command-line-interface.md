---
id: hydra_cli
title: Ory Hydra Command Line Interface
---

## Command Line Interface
Ory has a Command Line Interface (CLI) that ...

## Using the CLI with commands and arguments


## Commands
clients     Manage OAuth 2.0 Clients
help        Help about any command
keys        Manage JSON Web Keys
migrate     Various migration helpers
serve       Parent command for starting public and administrative HTTP/2 APIs
token       Issue and Manage OAuth2 tokens
version     Display this binary's version, build time and git hash of this build

## Reference

### `clients`

Alias: `clients`

| Options | Default | Description |
| --- | --- | --- |
| `create` | `-` | Create a new OAuth 2.0 Client |
| `delete` | `-` | Delete an OAuth 2.0 Client |
| `get` | `-` | Get an OAuth 2.0 Client |
| `import` | `-` | Import OAuth 2.0 Clients from one or more JSON files |
| `list` | `-` |  List OAuth 2.0 Clients |

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--access-token string` | `local` | `-` |  Set an access token to be used in the Authorization header, defaults to environment variable OAUTH2_ACCESS_TOKEN |
| `--endpoint string` | `local` | `-` | Set the URL where ORY Hydra is hosted, defaults to environment variable HYDRA_ADMIN_URL|
| `--fail-after duration` | `local` | `-` | Stop retrying after the specified duration (default 1m0s) |
| `--fake-tls-termination` | `local` | `-` | Fake tls termination by adding "X-Forwarded-Proto: https" to http headers |
| `-h, --help` | `local` | `-` | help for clients |
| `--config string` | `global` | `$HOME/.hydra.yaml` | Config file |
| `--skip-tls-verify` | `global` | `$HOME/.hydra.yaml` |  Accept TLS certificates signed by unkown certificate authorities !Foolishly|

The `clients` command manages OAUTH2 clients

**Example**
```bash
hydra clients create --fail-after duration
```

Use "hydra [command] --help" for more information about a command.
