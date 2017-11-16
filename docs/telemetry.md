# Telemetry

Our goal is to have the most reliable and fastest OAuth2 and OpenID Connect server. To achieve this goal,
ORY Hydra collects metrics on endpoint performance and sends a **fully anonymized** telemetry report
("anonymous usage statistics") to our servers. This data helps us understand how changes impact performance
and stability of ORY Hydra and identify potential issues.

ORY is a company registered in Germany, and we are very aware and sensible of and to privacy laws. Thus, we are fully
transparent on what data we transmit why and how. The source code of the telemetry package is completely open source
and located [here](https://github.com/ory/hydra/tree/master/metrics). If you do not wish to help us improving ORY Hydra
by sharing telemetry data, it is possible to [turn this feature off](#disabling-telemetry).

To protect your privacy, we filter out any data that could identify you, or your users. We are taking the following
measures to protect your privacy:

1. We only transmit information on how often endpoints are requested, how fast they respond and what http status code
was sent.
2. We filter out any query parameters, headers, response and request bodies and path parameters. A full list of transmitted
URL paths is listed in section [Request telemetry](#request-telemetry). For example:
  * `GET /clients/1235` becomes `GET /clients`
  * `GET /oauth2/token?param=1` becomes `GET /oauth2/token`
  * `POST /clients` strips out all post data and becomes `POST /clients`
4. **We are unable to see or store the IP address of your host**, as the
[IP is set to `0.0.0.0`](https://github.com/ory/hydra/tree/master/metrics/middleware.go) when transmitting data to Segment.
5. We do not transmit any environment information from the host, except the operating system id (windows, linux, osx),
the target architecture (amd64, darwin, ...), and the number of CPUs available on the host.

## Identification

To identify an installation and group together clusters, we create a SHA-512 hash of the Issuer URL for identification.
Additionally, each running instance is identified using an unique identifier which is set every time ORY Hydra starts. The identifier
is a Universally Unique Identifier (V4) and is thus a cryptographically safe random string. Identification is triggered
when the instance has been running for more than 15 minutes.

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

Additionally, the following data is submitted to us:

* Total number of requests, responses, response latencies and response sizes.
* Total number of requests, responses, response latencies and response sizes per HTTP method.
* Total number of requests, responses, response latencies and response sizes per anonymized API endpoint and HTTP method.
* Total number of requests, responses, response latencies and response sizes per anonymized API endpoint.
* Memory statistics such as heap allocations, gc cycles, and other.

A raw data example can be found [here](https://github.com/ory/hydra/tree/master/docs/metrics/telemetry-example.json).

## Keep alive

A keep-alive containing no information except the instance id is sent every 5 minutes.

## Data processing

Once the data was transmitted to [Segment.com](http://segment.com/) it is then fed to an AWS S3 bucket and stored
there for later analysis. At the moment, we are working on a python / numpy toolchain to help us analyze the data
we get.

## Disabling telemetry

You can disable telemetry with `hydra host --disable-telemetry`, using the [oryd/hydra:{tag}-without-telemetry](https://hub.docker.com/r/oryd/hydra/tags/) docker image, by
setting ` export DISABLE_TELEMETRY=1`, or `DISABLE_TELEMETRY=1 hydra host`.

Please be aware that disabling telemetry also disables metrics on the `/health` endpoint.
