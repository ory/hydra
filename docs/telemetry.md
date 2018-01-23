# Telemetry

Our goal is to have the fastest and most reliable OAuth2 and OpenID Connect server. To achieve this goal,
we collect metrics on endpoint performance and send a **fully anonymized** telemetry report
("anonymous usage statistics") to our servers. This data helps us understand how changes impact performance
and stability of ORY Hydra and identify potential issues.

We are committed to full transparency on what data we transmit why and how. The source code of the telemetry package is completely open source
and located [here](https://github.com/ory/hydra/tree/master/metrics). If you do not wish to help us improving ORY Hydra
by sharing telemetry data, it is possible to [turn this feature off](#disabling-telemetry).

To protect your privacy, we filter out any data that could identify you or your users. We are taking the following
measures to protect your privacy:

1. We only transmit information on how often endpoints are requested, how fast they respond and what HTTP status code was sent.
2. We filter out any query parameters, headers, response and request bodies and path parameters. A full list of transmitted
URL paths is listed in section [Request telemetry](#request-telemetry).
4. **We are unable to see or store the IP address of your host**, as the
[IP is set to `0.0.0.0`](https://github.com/ory/hydra/tree/master/metrics/middleware.go) when transmitting data to our metrics aggregator.
5. We do not transmit any environment information from the host, except:
  * Operating system id (windows, linux, osx)
  * The target architecture (amd64, darwin, ...)
  * Number of CPUs available
  * Build time, hash and version of ORY Hydra
  * Memory consumption of ORY Hydra's process

## Identification

To identify an installation and group together clusters, we create a SHA-256 hash of the Issuer URL for identification.
Additionally, each running instance is identified using an unique identifier which is set every time ORY Hydra starts. The identifier
is a Universally Unique Identifier (V4) and is thus a cryptographically safe random string. Identification is triggered
when we are confident that the instance is not a test instance (e.g. one of our tutorials).

We collect the following system metrics:

* `goarch`: The target architecture of the ORY Hydra binary.
* `goos`: The target system of the ORY Hydra binary.
* `numCpu`: The number of CPUs available.
* `runtimeVersion`: The go version used to create the binary.
* `version`: The version of this binary.
* `hash`: The git hash of this binary.
* `buildTime`: The build time of this binary.

## Request telemetry

We collect telemetry data from the following endpoints:

```
"/.well-known/jwks.json",
"/.well-known/openid-configuration",
"/clients",
"/health",
"/keys",
"/userinfo",
"/oauth2/consent-fallback",
"/oauth2/auth",
"/oauth2/introspect",
"/oauth2/revoke",
"/oauth2/token",
"/oauth2/consent/requests",
"/policies",
"/warden/allowed",
"/warden/groups",
"/warden/token/allowed",
"/",
```

The ip addresses of both host and client are anonymized to `0.0.0.0`. Any identifiable information in the URL path and query is hashed with
sha256 using a randomly assigned uuid v4 salt:

* `/clients/foo` with salt `ABCDEFGH` becomes `/clients/sha256("foo|ABCDEFGH")`: `/clients/0301424a80469ad03a208de925563a97ec6ab2f9dc7a2ad71b2ded85a7f7a7af`
* `/policies?owner=foo` with salt `ABCDEFGH` becomes `/policies?owner=sha256("foo|ABCDEFGH")`: `/policies?owner=0301424a80469ad03a208de925563a97ec6ab2f9dc7a2ad71b2ded85a7f7a7af`). We do not

We do not collect the hostname, headers, payloads. Here is what we do collect:

```go
path := anonymizePath(r.URL.Path)
query := anonymizeQuery(r.URL.Query())
analytics.
			NewProperties().
			SetURL(scheme + "//" + sw.ID + path + "?" + query). // sw.ID is the SHA-256 hash of the Issuer URL
			SetPath(path).
			SetName(path).
			Set("status", status).
			Set("size", size).
			Set("latency", latency).
			Set("instance", sw.InstanceID).
			Set("method", r.Method),
```

A raw data example can be found [here](https://github.com/ory/hydra/tree/master/docs/metrics/telemetry-example.json).

## Data processing

Once the data was transmitted to [Segment.com](http://segment.com/) it is then fed to an encrypted AWS S3 bucket and stored
for analysis using Apache Flink.

We analyze the data from ORY Hydra deployments with the following goals:

1. Be able to say how many production deployments exist.
2. Understand how much throughput deployments handle.
3. Evaluate how frequently specific features (e.g. policies) are used.
4. Detect issues introduced by new features (e.g. buggy releases). For example:
  * After release 0.X.Y, all instances show 25% increase in response times for Warden API calls.
5. Identify real-world problems caused by things such as slow queries. For example:
  * Searching for policies by owners takes causes high response times.
  * Running the deployment for several months and high traffic causes slow response times.

## Disabling telemetry

You can disable telemetry with `hydra host --disable-telemetry`, using the [oryd/hydra:{tag}-without-telemetry](https://hub.docker.com/r/oryd/hydra/tags/) docker image, by
setting `export DISABLE_TELEMETRY=1`, or `DISABLE_TELEMETRY=1 hydra host`.

Disabling telemetry does not have any downsides, except for us not being able to detect issues automatically ;)
