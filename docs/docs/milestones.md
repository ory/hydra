---
id: milestones
title: Milestones and Roadmap
---

## [v2.0](https://github.com/ory/hydra/milestone/42)

_This milestone does not have a description._

### [Bug](https://github.com/ory/hydra/labels/bug)

Something is not working.

#### Issues

- [ ] Make cookies with SameSite=None secure by default or using the
      configuration flag
      ([hydra#1844](https://github.com/ory/hydra/issues/1844))
- [ ] client_id case sensitivity is not properly enforced when using MySQL
      ([hydra#1644](https://github.com/ory/hydra/issues/1644)) -
      [@Patrik](https://github.com/zepatrik)
- [ ] Client allowed_cors_origins not working
      ([hydra#1754](https://github.com/ory/hydra/issues/1754))
- [ ] Consider customizing 'azp' and 'aud' claims in ID Tokens
      ([hydra#2042](https://github.com/ory/hydra/issues/2042))
- [ ] Do not return `email` in `id_token` but instead in `userinfo` for specific
      response types ([hydra#2163](https://github.com/ory/hydra/issues/2163)) -
      [@hackerman](https://github.com/aeneasr)
- [ ] SQL persister uses | to store scopes and audiences without any escaping
      ([hydra#2859](https://github.com/ory/hydra/issues/2859)) -
      [@Grant Zvolský](https://github.com/grantzvolsky)

### [Feat](https://github.com/ory/hydra/labels/feat)

New feature or request.

#### Issues

- [ ] Refactor client CLI
      ([hydra#2124](https://github.com/ory/hydra/issues/2124)) -
      [@Patrik](https://github.com/zepatrik)
- [ ] Rename DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY
      ([hydra#1760](https://github.com/ory/hydra/issues/1760)) -
      [@hackerman](https://github.com/aeneasr)
- [ ] issuer in discovery document contains trailing '/'
      ([hydra#1482](https://github.com/ory/hydra/issues/1482))
- [ ] Make cookies with SameSite=None secure by default or using the
      configuration flag
      ([hydra#1844](https://github.com/ory/hydra/issues/1844))
- [ ] Consider recreating Hydra V2 database model instead of migrations
      ([hydra#2902](https://github.com/ory/hydra/issues/2902)) -
      [@Grant Zvolský](https://github.com/grantzvolsky),
      [@hackerman](https://github.com/aeneasr)
- [ ] Rename SDK methods to follow our OpenAPI spec guide
      ([hydra#2908](https://github.com/ory/hydra/issues/2908))
- [ ] No longer allow users to set the client ID
      ([hydra#2911](https://github.com/ory/hydra/issues/2911))
- [ ] Move to go-jose key generation
      ([hydra#1825](https://github.com/ory/hydra/issues/1825))
- [ ] Auth session cannot be prolonged even if the user is active
      ([hydra#1690](https://github.com/ory/hydra/issues/1690))
- [ ] Token claims customization with Jsonnet
      ([hydra#1748](https://github.com/ory/hydra/issues/1748)) -
      [@hackerman](https://github.com/aeneasr)
- [ ] Update clients from cli
      ([hydra#2020](https://github.com/ory/hydra/issues/2020))
- [x] Refactor SQL Migration tests to match new system
      ([hydra#2901](https://github.com/ory/hydra/issues/2901)) -
      [@Grant Zvolský](https://github.com/grantzvolsky),
      [@hackerman](https://github.com/aeneasr)

## [next](https://github.com/ory/hydra/milestone/41)

_This milestone does not have a description._

### [Bug](https://github.com/ory/hydra/labels/bug)

Something is not working.

#### Issues

- [x] Space character in secret.system value
      ([hydra#2609](https://github.com/ory/hydra/issues/2609)) -
      [@Patrik](https://github.com/zepatrik),
      [@Jakub Błaszczyk](https://github.com/Demonsthere)

### [Feat](https://github.com/ory/hydra/labels/feat)

New feature or request.

#### Issues

- [x] Reuse Detection in Refresh Token Rotation
      ([hydra#2022](https://github.com/ory/hydra/issues/2022))

### [Rfc](https://github.com/ory/hydra/labels/rfc)

A request for comments to discuss and share ideas.

#### Issues

- [x] Multi-region deployment support
      ([hydra#2018](https://github.com/ory/hydra/issues/2018))

## [v1.10](https://github.com/ory/hydra/milestone/40)

_This milestone does not have a description._

### [Bug](https://github.com/ory/hydra/labels/bug)

Something is not working.

#### Issues

- [x] Slow consent revocation request
      ([hydra#1997](https://github.com/ory/hydra/issues/1997))
- [x] Report expired JWT assertion token to client
      ([hydra#2066](https://github.com/ory/hydra/issues/2066))
- [x] Client update changes it's PK to 0
      ([hydra#2148](https://github.com/ory/hydra/issues/2148)) -
      [@Patrik](https://github.com/zepatrik)
- [x] CORS error with v1.9 on localhost
      ([hydra#2165](https://github.com/ory/hydra/issues/2165)) -
      [@hackerman](https://github.com/aeneasr)
- [x] Invalid json response with get login request
      ([hydra#2515](https://github.com/ory/hydra/issues/2515))
- [x] Invalid TLS config after upgrading to 1.10.2
      ([hydra#2518](https://github.com/ory/hydra/issues/2518))

#### Pull Requests

- [x] Deprecate client flags in introspect
      ([hydra#2011](https://github.com/ory/hydra/pull/2011)) -
      [@hackerman](https://github.com/aeneasr)
- [x] fix: bump ory/fosite to v0.34.1 to address CVEs
      ([hydra#2090](https://github.com/ory/hydra/pull/2090)) -
      [@hackerman](https://github.com/aeneasr)
- [x] ci: resolve ci release issues
      ([hydra#2094](https://github.com/ory/hydra/pull/2094)) -
      [@hackerman](https://github.com/aeneasr)
- [x] Prepare OpenID Connect Conformity test suite with new profiles and
      regression fixes ([hydra#2170](https://github.com/ory/hydra/pull/2170)) -
      [@hackerman](https://github.com/aeneasr)
- [x] test: resolve conformity test suite concurrency issues
      ([hydra#2181](https://github.com/ory/hydra/pull/2181)) -
      [@hackerman](https://github.com/aeneasr)
- [x] test: completely refactor consent tests and resolve logout issue
      ([hydra#2227](https://github.com/ory/hydra/pull/2227)) -
      [@hackerman](https://github.com/aeneasr)

### [Feat](https://github.com/ory/hydra/labels/feat)

New feature or request.

#### Issues

- [x] Publish a generated csharp SDK
      ([hydra#2017](https://github.com/ory/hydra/issues/2017))

#### Pull Requests

- [x] perf: add (client_id, subject) index to access and refresh tables to
      improve revocation performance
      ([hydra#2001](https://github.com/ory/hydra/pull/2001)) -
      [@hackerman](https://github.com/aeneasr)
- [x] Prepare OpenID Connect Conformity test suite with new profiles and
      regression fixes ([hydra#2170](https://github.com/ory/hydra/pull/2170)) -
      [@hackerman](https://github.com/aeneasr)

### [Blocking](https://github.com/ory/hydra/labels/blocking)

Blocks milestones or other issues or pulls.

#### Issues

- [x] Client update changes it's PK to 0
      ([hydra#2148](https://github.com/ory/hydra/issues/2148)) -
      [@Patrik](https://github.com/zepatrik)
- [x] CORS error with v1.9 on localhost
      ([hydra#2165](https://github.com/ory/hydra/issues/2165)) -
      [@hackerman](https://github.com/aeneasr)

#### Pull Requests

- [x] Deprecate client flags in introspect
      ([hydra#2011](https://github.com/ory/hydra/pull/2011)) -
      [@hackerman](https://github.com/aeneasr)
- [x] ci: resolve ci release issues
      ([hydra#2094](https://github.com/ory/hydra/pull/2094)) -
      [@hackerman](https://github.com/aeneasr)

## [v1.11](https://github.com/ory/hydra/milestone/39)

_This milestone does not have a description._

### [Bug](https://github.com/ory/hydra/labels/bug)

Something is not working.

#### Issues

- [x] Introspection Response: `access_token` and `refresh_token` are not valid
      `token_type` ([hydra#1762](https://github.com/ory/hydra/issues/1762))
- [x] RSA key generation is slow on ARM
      ([hydra#1989](https://github.com/ory/hydra/issues/1989))
- [x] `loginRequest.requested_access_token_audience` should not be `null`
      ([hydra#2039](https://github.com/ory/hydra/issues/2039))
- [x] Redirect URI should be able to contain plus (+) character
      ([hydra#2055](https://github.com/ory/hydra/issues/2055))
- [x] Docs: rendering issue (?) on reference REST API
      ([hydra#2092](https://github.com/ory/hydra/issues/2092)) -
      [@Vincent](https://github.com/vinckr)
- [x] Jaeger being unavailable is a fatal error that stops service from starting
      ([hydra#2642](https://github.com/ory/hydra/issues/2642))

### [Feat](https://github.com/ory/hydra/labels/feat)

New feature or request.

#### Issues

- [x] consent: Improve remember for consent
      ([hydra#1006](https://github.com/ory/hydra/issues/1006))
- [x] [Feature] Enhance Security Middleware
      ([hydra#1029](https://github.com/ory/hydra/issues/1029))
- [x] cmd: Add upsert command for client CLI
      ([hydra#1086](https://github.com/ory/hydra/issues/1086)) -
      [@hackerman](https://github.com/aeneasr)
- [x] oauth2: Make cleaning up refresh and authz codes possible
      ([hydra#1130](https://github.com/ory/hydra/issues/1130)) -
      [@hackerman](https://github.com/aeneasr)
- [x] consent: Allow removing tokens without revoking consent
      ([hydra#1142](https://github.com/ory/hydra/issues/1142)) -
      [@hackerman](https://github.com/aeneasr)
- [x] OAuth Client authentication creation CLI jwks client field not present
      ([hydra#1404](https://github.com/ory/hydra/issues/1404))
- [x] Add oAuth2Client to logoutRequest similar to loginRequest.
      ([hydra#1483](https://github.com/ory/hydra/issues/1483))
- [x] Add a way to filter/sort the list of clients
      ([hydra#1485](https://github.com/ory/hydra/issues/1485)) -
      [@hackerman](https://github.com/aeneasr)
- [x] Remove "not before" claim "nbf" from JWT access token
      ([hydra#1542](https://github.com/ory/hydra/issues/1542))
- [x] No way to handle 409 GetLoginRequestConflict.
      ([hydra#1569](https://github.com/ory/hydra/issues/1569)) -
      [@Alano Terblanche](https://github.com/Benehiko)
- [x] Add endpoint to Admin API to revoke access tokens
      ([hydra#1728](https://github.com/ory/hydra/issues/1728))
- [x] Migrate to gobuffalo/pop
      ([hydra#1730](https://github.com/ory/hydra/issues/1730)) -
      [@Patrik](https://github.com/zepatrik)
- [x] CLI Migration Down
      ([hydra#1763](https://github.com/ory/hydra/issues/1763))
- [x] Split HTTPS handling for public/admin
      ([hydra#1962](https://github.com/ory/hydra/issues/1962))
- [x] issueLogoutVerifier should allow POST requests as well
      ([hydra#1993](https://github.com/ory/hydra/issues/1993))
- [x] Expired token is considered an error
      ([hydra#2031](https://github.com/ory/hydra/issues/2031))
- [x] Automatically set GOMAXPROCS according to linux container cpu quota
      ([hydra#2033](https://github.com/ory/hydra/issues/2033))
- [x] Find out if a login/consent challenge is still valid
      ([hydra#2057](https://github.com/ory/hydra/issues/2057))
- [x] Prometheus endpoint should not require x-forwarded-proto header
      ([hydra#2072](https://github.com/ory/hydra/issues/2072))

#### Pull Requests

- [x] feat: OpenID Connect Dynamic Client Registration and OAuth2 Dynamic Client
      Registration Protocol
      ([hydra#2909](https://github.com/ory/hydra/pull/2909)) -
      [@hackerman](https://github.com/aeneasr)

### [Rfc](https://github.com/ory/hydra/labels/rfc)

A request for comments to discuss and share ideas.

#### Issues

- [x] Split HTTPS handling for public/admin
      ([hydra#1962](https://github.com/ory/hydra/issues/1962))

## [v1.7.0](https://github.com/ory/hydra/milestone/38)

_This milestone does not have a description._

### [Feat](https://github.com/ory/hydra/labels/feat)

New feature or request.

#### Issues

- [x] "debug" log level outputs multiline logs
      ([hydra#1958](https://github.com/ory/hydra/issues/1958))

## [v1.6.0](https://github.com/ory/hydra/milestone/37)

_This milestone does not have a description._

### [Bug](https://github.com/ory/hydra/labels/bug)

Something is not working.

#### Issues

- [x] Loopback interface redirection with arbitrary port
      ([hydra#1732](https://github.com/ory/hydra/issues/1732))

### [Feat](https://github.com/ory/hydra/labels/feat)

New feature or request.

#### Issues

- [x] Use consistent field types for logging
      ([hydra#1683](https://github.com/ory/hydra/issues/1683))

## [v1.5.0](https://github.com/ory/hydra/milestone/36)

_This milestone does not have a description._

### [Bug](https://github.com/ory/hydra/labels/bug)

Something is not working.

#### Issues

- [x] Invalid ttl.refresh_token -1 (no expiration)
      ([hydra#1811](https://github.com/ory/hydra/issues/1811)) -
      [@Patrik](https://github.com/zepatrik)
- [x] /userinfo endpoint misses www-authenticate header for 401 response
      ([hydra#1827](https://github.com/ory/hydra/issues/1827))
- [x] Superfluous response.writeHeader
      ([hydra#1842](https://github.com/ory/hydra/issues/1842)) -
      [@hackerman](https://github.com/aeneasr)
- [x] config: scopes_supported doesn't have offline_access
      ([hydra#1866](https://github.com/ory/hydra/issues/1866))

### [Feat](https://github.com/ory/hydra/labels/feat)

New feature or request.

#### Pull Requests

- [x] refactor: move migrations to gobuffalo/fizz
      ([hydra#1775](https://github.com/ory/hydra/pull/1775)) -
      [@hackerman](https://github.com/aeneasr),
      [@Patrik](https://github.com/zepatrik)
