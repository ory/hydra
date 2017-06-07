# Telemetry

Our goal is to have the most reliable and fastest OAuth2 and OpenID Connect server. To achieve this goal,
we meter endpoint performance and send a **summarized, anonymized** telemetry report to our servers. This helps us
identify how changes impact performance and stability of ORY Hydra.

Because this is a critical issue to many, we are fully transparent on this issue. The source code of the telemetry
package is completely open source and located [here](https://github.com/ory/hydra/tree/master/metrics). We also compiled an overview for you.

## Identification

The ORY Hydra instance is identified using an unique identifier which is set every time ORY Hydra starts. The identifier
is Universally Unique Identifier (V4) and is crypto-random.

**We never identify, or share the ip address of the host.** The ip address is anonymized to `0.0.0.0` and is not available
to us.

On start up, we collect the following system metrics:

* `goarch`: The target architecture of the ORY Hydra binary.
* `goos`: The target system of the ORY Hydra binary.
* `numCpu`: The number of CPUs available.
* `runtimeVersion`: The go version used to create the binary.
* `version`: The version of this binary.
* `hash`: The git hash of this binary.
* `buildTime`: The build time of this binary.

## Request telemetry

A summarized telemetry report is sent every fifteen (15) minutes. **We never identify, or share the ip address of the
host.** The ip address is anonymized to `0.0.0.0` and is not available to us.

We share metrics for each HTTP API endpoint. The endpoint is filtered and does not include any query parameters or
other identifying information such as OAuth2 Client IDs. The following is the complete list that is used to anonymize
the endpoints:

```
"/.well-known/jwks.json",
"/.well-known/openid-configuration",
"/clients",
"/health",
"/keys",
"/oauth2/auth",
"/oauth2/introspect",
"/oauth2/revoke",
"/oauth2/token",
"/policies",
"/warden/allowed",
"/warden/groups",
"/warden/token/allowed",
"/",
```

Additionally, the following data is shared:

* Total number of requests, responses, response latencies and response sizes.
* Total number of requests, responses, response latencies and response sizes per HTTP method (GET, POST, DELETE, ...).
* Total number of requests, responses, response latencies and response sizes per API endpoint (fully anonymized, e.g. `/oauth2/token`) HTTP method (GET, POST, DELETE, ...).
* Total number of requests, responses, response latencies and response sizes per API endpoint (fully anonymized, e.g. `/oauth2/token`).

A raw data example can be found [here](https://github.com/ory/hydra/tree/master/docs/metrics/telemetry-example.json).

## Disabling telemetry

You can disable telemetry with `hydra host --disable-telemetry`, using the [oryd/hydra:{tag}-without-telemetry](https://hub.docker.com/r/oryd/hydra/tags/) docker image, by
setting ` export DISABLE_TELEMETRY=1`, or `DISABLE_TELEMETRY=1 hydra host`.

Please be aware that disabling telemetry also disables metrics on the `/health` endpoint.
