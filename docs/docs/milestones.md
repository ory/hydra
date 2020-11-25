---
id: milestones
title: Milestones and Roadmap
---

## [next](https://github.com/ory/hydra/milestone/41)

_This milestone does not have a description._

### [Feat](https://github.com/ory/hydra/labels/feat)

New feature or request.

#### Issues

- [ ] Reuse Detection in Refresh Token Rotation
      ([hydra#2022](https://github.com/ory/hydra/issues/2022))
- [ ] Device Authorization Grant
      ([hydra#1553](https://github.com/ory/hydra/issues/1553))
- [ ] Support AWS X-Ray as distributed trace
      ([hydra#2130](https://github.com/ory/hydra/issues/2130))

### [Rfc](https://github.com/ory/hydra/labels/rfc)

A request for comments to discuss and share ideas.

#### Issues

- [ ] Multi-region deployment support
      ([hydra#2018](https://github.com/ory/hydra/issues/2018))
- [ ] Device Authorization Grant
      ([hydra#1553](https://github.com/ory/hydra/issues/1553))

## [v1.9.1](https://github.com/ory/hydra/milestone/40)

_This milestone does not have a description._

### [Bug](https://github.com/ory/hydra/labels/bug)

Something is not working.

#### Issues

- [ ] Report expired JWT assertion token to client
      ([hydra#2066](https://github.com/ory/hydra/issues/2066))
- [ ] Client allowed_cors_origins not working
      ([hydra#1754](https://github.com/ory/hydra/issues/1754))
- [ ] CORS error with v1.9 on localhost
      ([hydra#2165](https://github.com/ory/hydra/issues/2165)) -
      [@hackerman](https://github.com/aeneasr)
- [x] Slow consent revocation request
      ([hydra#1997](https://github.com/ory/hydra/issues/1997))
- [x] Client update changes it's PK to 0
      ([hydra#2148](https://github.com/ory/hydra/issues/2148)) -
      [@Patrik](https://github.com/zepatrik)

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

### [Docs](https://github.com/ory/hydra/labels/docs)

Affects documentation.

#### Issues

- [ ] Docs: reference configuration provides invalid bcrypt cost and some other
      outdated values ([hydra#2077](https://github.com/ory/hydra/issues/2077)) -
      [@Vincent](https://github.com/vinckr)
- [ ] Configuration is not recognized by certain docker containers
      ([hydra#2164](https://github.com/ory/hydra/issues/2164)) -
      [@hackerman](https://github.com/aeneasr)
- [x] docs: Hydra supports two concepts of authorization
      ([hydra#1992](https://github.com/ory/hydra/issues/1992)) -
      [@hackerman](https://github.com/aeneasr)
- [x] Document GitLab setup with ORY Hydra as the OpenID Connect Provider
      ([hydra#2007](https://github.com/ory/hydra/issues/2007)) -
      [@Kim Neunert](https://github.com/k9ert)

#### Pull Requests

- [x] docs: remove introspect security spec
      ([hydra#2002](https://github.com/ory/hydra/pull/2002))
- [x] docs: fix omissions in the description
      ([hydra#2173](https://github.com/ory/hydra/pull/2173)) -
      [@Vincent](https://github.com/vinckr)

### [Blocking](https://github.com/ory/hydra/labels/blocking)

Blocks milestones or other issues or pulls.

#### Issues

- [ ] CORS error with v1.9 on localhost
      ([hydra#2165](https://github.com/ory/hydra/issues/2165)) -
      [@hackerman](https://github.com/aeneasr)
- [x] Client update changes it's PK to 0
      ([hydra#2148](https://github.com/ory/hydra/issues/2148)) -
      [@Patrik](https://github.com/zepatrik)

#### Pull Requests

- [x] Deprecate client flags in introspect
      ([hydra#2011](https://github.com/ory/hydra/pull/2011)) -
      [@hackerman](https://github.com/aeneasr)
- [x] ci: resolve ci release issues
      ([hydra#2094](https://github.com/ory/hydra/pull/2094)) -
      [@hackerman](https://github.com/aeneasr)

### [Ci](https://github.com/ory/hydra/labels/ci)

Affects Continuous Integration (CI).

#### Pull Requests

- [x] ci: fix etcd CVEs ([hydra#2003](https://github.com/ory/hydra/pull/2003)) -
      [@hackerman](https://github.com/aeneasr)
- [x] test: resolve conformity test suite concurrency issues
      ([hydra#2181](https://github.com/ory/hydra/pull/2181)) -
      [@hackerman](https://github.com/aeneasr)

## [v1.10.0](https://github.com/ory/hydra/milestone/39)

_This milestone does not have a description._

### [Bug](https://github.com/ory/hydra/labels/bug)

Something is not working.

#### Issues

- [ ] client_id case sensitivity is not properly enforced when using MySQL
      ([hydra#1644](https://github.com/ory/hydra/issues/1644)) -
      [@Patrik](https://github.com/zepatrik)
- [ ] Make cookies with SameSite=None secure by default or using the
      configuration flag
      ([hydra#1844](https://github.com/ory/hydra/issues/1844))
- [ ] RSA key generation is slow on ARM
      ([hydra#1989](https://github.com/ory/hydra/issues/1989))
- [ ] `loginRequest.requested_access_token_audience` should not be `null`
      ([hydra#2039](https://github.com/ory/hydra/issues/2039))
- [ ] Consider customizing 'azp' and 'aud' claims in ID Tokens
      ([hydra#2042](https://github.com/ory/hydra/issues/2042))
- [ ] Docs: rendering issue (?) on reference REST API
      ([hydra#2092](https://github.com/ory/hydra/issues/2092)) -
      [@Vincent](https://github.com/vinckr)
- [ ] Do not return `email` in `id_token` but instead in `userinfo` for specific
      response types ([hydra#2163](https://github.com/ory/hydra/issues/2163)) -
      [@hackerman](https://github.com/aeneasr)
- [x] Introspection Response: `access_token` and `refresh_token` are not valid
      `token_type` ([hydra#1762](https://github.com/ory/hydra/issues/1762))
- [x] Redirect URI should be able to contain plus (+) character
      ([hydra#2055](https://github.com/ory/hydra/issues/2055))

### [Feat](https://github.com/ory/hydra/labels/feat)

New feature or request.

#### Issues

- [ ] consent: Improve remember for consent
      ([hydra#1006](https://github.com/ory/hydra/issues/1006))
- [ ] [Feature] Enhance Security Middleware
      ([hydra#1029](https://github.com/ory/hydra/issues/1029))
- [ ] Add API versioning for administrative APIs
      ([hydra#1050](https://github.com/ory/hydra/issues/1050))
- [ ] consent: Allow removing tokens without revoking consent
      ([hydra#1142](https://github.com/ory/hydra/issues/1142)) -
      [@hackerman](https://github.com/aeneasr)
- [ ] OAuth Client authentication creation CLI jwks client field not present
      ([hydra#1404](https://github.com/ory/hydra/issues/1404))
- [ ] Add oAuth2Client to logoutRequest similar to loginRequest.
      ([hydra#1483](https://github.com/ory/hydra/issues/1483))
- [ ] Add a way to filter/sort the list of clients
      ([hydra#1485](https://github.com/ory/hydra/issues/1485)) -
      [@hackerman](https://github.com/aeneasr)
- [ ] Remove "not before" claim "nbf" from JWT access token
      ([hydra#1542](https://github.com/ory/hydra/issues/1542))
- [ ] No way to handle 409 GetLoginRequestConflict.
      ([hydra#1569](https://github.com/ory/hydra/issues/1569)) -
      [@Patrik](https://github.com/zepatrik)
- [ ] Auth session cannot be prolonged even if the user is active
      ([hydra#1690](https://github.com/ory/hydra/issues/1690))
- [ ] Rename DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY
      ([hydra#1760](https://github.com/ory/hydra/issues/1760)) -
      [@hackerman](https://github.com/aeneasr)
- [ ] CLI Migration Down
      ([hydra#1763](https://github.com/ory/hydra/issues/1763))
- [ ] Move to go-jose key generation
      ([hydra#1825](https://github.com/ory/hydra/issues/1825))
- [ ] Make cookies with SameSite=None secure by default or using the
      configuration flag
      ([hydra#1844](https://github.com/ory/hydra/issues/1844))
- [ ] Split HTTPS handling for public/admin
      ([hydra#1962](https://github.com/ory/hydra/issues/1962))
- [ ] Token claims customization with Jsonnet
      ([hydra#1748](https://github.com/ory/hydra/issues/1748)) -
      [@hackerman](https://github.com/aeneasr)
- [ ] Update clients from cli
      ([hydra#2020](https://github.com/ory/hydra/issues/2020))
- [ ] Find out if a login/consent challenge is still valid
      ([hydra#2057](https://github.com/ory/hydra/issues/2057))
- [x] cmd: Add upsert command for client CLI
      ([hydra#1086](https://github.com/ory/hydra/issues/1086)) -
      [@hackerman](https://github.com/aeneasr)
- [x] oauth2: Make cleaning up refresh and authz codes possible
      ([hydra#1130](https://github.com/ory/hydra/issues/1130)) -
      [@hackerman](https://github.com/aeneasr)
- [x] Add endpoint to Admin API to revoke access tokens
      ([hydra#1728](https://github.com/ory/hydra/issues/1728))
- [x] Migrate to gobuffalo/pop
      ([hydra#1730](https://github.com/ory/hydra/issues/1730)) -
      [@Patrik](https://github.com/zepatrik)
- [x] issueLogoutVerifier should allow POST requests as well
      ([hydra#1993](https://github.com/ory/hydra/issues/1993))
- [x] Expired token is considered an error
      ([hydra#2031](https://github.com/ory/hydra/issues/2031))
- [x] Automatically set GOMAXPROCS according to linux container cpu quota
      ([hydra#2033](https://github.com/ory/hydra/issues/2033))
- [x] Prometheus endpoint should not require x-forwarded-proto header
      ([hydra#2072](https://github.com/ory/hydra/issues/2072))

### [Docs](https://github.com/ory/hydra/labels/docs)

Affects documentation.

#### Issues

- [ ] Docs: rendering issue (?) on reference REST API
      ([hydra#2092](https://github.com/ory/hydra/issues/2092)) -
      [@Vincent](https://github.com/vinckr)
- [ ] Add list of possible errors to REST API documentation
      ([hydra#2166](https://github.com/ory/hydra/issues/2166))
- [x] 'make deps' Not work
      ([hydra#2154](https://github.com/ory/hydra/issues/2154)) -
      [@hackerman](https://github.com/aeneasr)

### [Rfc](https://github.com/ory/hydra/labels/rfc)

A request for comments to discuss and share ideas.

#### Issues

- [ ] Split HTTPS handling for public/admin
      ([hydra#1962](https://github.com/ory/hydra/issues/1962))
