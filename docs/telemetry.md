# Telemetry

Our goal is to have the most reliable and fastest OAuth2 and OpenID Connect server. To achieve this goal,
we meter endpoint performance and send a **summarized, anonymized** telemetry report ("anonymous usage statistics") to our servers. This helps us
identify how changes impact performance and stability of ORY Hydra.

We are fully transparent on this issue. The source code of the telemetry package is completely open source
and located [here](https://github.com/ory/hydra/tree/master/metrics). It is possible to [turn this feature off](#disabling-telemetry), but
we kindly ask you to keep it enabled.

**We are unable to link telemetry data to ip addresses.** The host's ip address is fully anonymized and never available to
us. We collect but a few metrics, including latency, uptime, overall responses and information on the host such as number of
CPUs. We do not collect user data, we do not collect data from the database. We filter all data out of request URLs and group
them together by feature. For example:

* `GET /clients/1235` becomes `GET /clients`
* `GET /oauth2/token?param=1` becomes `GET /oauth2/token`
* `POST /clients` strips out all post data and becomes `POST /clients`

## Identification

The ORY Hydra instance is identified using an unique identifier which is set every time ORY Hydra starts. The identifier
is Universally Unique Identifier (V4) and is crypto-random. Identification is triggered when the instance has been
running for more than 15 minutes.

**We never identify, or share the ip address of the host.** The ip address is anonymized to `0.0.0.0` and is not available
to us.

We collect the following system metrics:

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
* Total number of requests, responses, response latencies and response sizes per HTTP method.
* Total number of requests, responses, response latencies and response sizes per anonymized API endpoint and HTTP method.
* Total number of requests, responses, response latencies and response sizes per anonymized API endpoint.

A raw data example can be found [here](https://github.com/ory/hydra/tree/master/docs/metrics/telemetry-example.json).

## Disabling telemetry

You can disable telemetry with `hydra host --disable-telemetry`, using the [oryd/hydra:{tag}-without-telemetry](https://hub.docker.com/r/oryd/hydra/tags/) docker image, by
setting ` export DISABLE_TELEMETRY=1`, or `DISABLE_TELEMETRY=1 hydra host`.

Please be aware that disabling telemetry also disables metrics on the `/health` endpoint.
