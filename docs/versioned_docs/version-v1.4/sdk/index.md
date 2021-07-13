---
id: index
title: Overview
---

All SDKs use automated code generation provided by
[`openapi-generator`](https://github.com/OpenAPITools/openapi-generator).
Unfortunately, `openapi-generator` has serious breaking changes in the generated
code when upgrading versions. Therefore, we do not make backwards compatibility
promises with regards to the generated SDKs. We hope to improve this process in
the future.

Before you check out the SDKs, head over to the [REST API](../reference/api)
documentation which includes code samples for common programming languages for
each REST endpoint.

> The SDKs do not provide a good API for dealing with OAuth 2.0 Flows (e.g.
> Authorize Code Flow, Refresh Flow, ...). But there are tons of
> [libraries available for consuming OAuth 2.0](https://oauth.net/code/). Do not
> write your own OAuth 2.0 Library!

We publish our SDKs for popular languages in their respective package
repositories:

- [Python](https://pypi.org/project/ory-hydra-client/)
- [PHP](https://packagist.org/packages/ory/hydra-client)
- [Go](https://github.com/ory/hydra-client-go)
- [NodeJS](https://www.npmjs.com/package/@ory/hydra-client) (with TypeScript)
- [Java](https://search.maven.org/artifact/sh.ory.hydra/hydra-client)
- [Ruby](https://rubygems.org/gems/ory-hydra-client)

We also provide more info for these SDKs:

- [Golang](go)
- [JavaScript](js)

Missing your programming language?
[Create an issue](https://github.com/ory/hydra/issues) and help us build, test
and publish the SDK for your programming language!
