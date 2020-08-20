---
id: milestones
title: Milestones and Roadmap
---

## [v1.7.1](https://github.com/ory/hydra/milestone/40)

*This milestone does not have a description.*

### [Bug](https://github.com/ory/hydra/labels/bug)

Something is not working.

#### Issues

* [ ] Slow consent revocation request ([hydra#1997](https://github.com/ory/hydra/issues/1997))

### [Feat](https://github.com/ory/hydra/labels/feat)

New feature or request.

#### Pull Requests

* [ ] perf: add (client_id, subject) index to access and refresh tables to improve revocation performance ([hydra#2001](https://github.com/ory/hydra/pull/2001)) - [@hackerman](https://github.com/aeneasr)

### [Docs](https://github.com/ory/hydra/labels/docs)

Affects documentation.

#### Pull Requests

* [x] docs: remove introspect security spec ([hydra#2002](https://github.com/ory/hydra/pull/2002))

### [Ci](https://github.com/ory/hydra/labels/ci)

Affects Continuous Integration (CI).

#### Pull Requests

* [x] ci: fix etcd CVEs ([hydra#2003](https://github.com/ory/hydra/pull/2003)) - [@hackerman](https://github.com/aeneasr)

## [v1.8.0](https://github.com/ory/hydra/milestone/39)

*This milestone does not have a description.*

### [Bug](https://github.com/ory/hydra/labels/bug)

Something is not working.

#### Issues

* [ ] client_id case sensitivity is not properly enforced when using MySQL ([hydra#1644](https://github.com/ory/hydra/issues/1644)) - [@Patrik](https://github.com/zepatrik)
* [ ] Introspection Response: `access_token` and `refresh_token` are not valid `token_type` ([hydra#1762](https://github.com/ory/hydra/issues/1762))
* [ ] Make cookies with SameSite=None secure by default or using the configuration flag ([hydra#1844](https://github.com/ory/hydra/issues/1844))
* [ ] RSA key generation is slow on ARM ([hydra#1989](https://github.com/ory/hydra/issues/1989))

### [Feat](https://github.com/ory/hydra/labels/feat)

New feature or request.

#### Issues

* [ ] consent: Improve remember for consent ([hydra#1006](https://github.com/ory/hydra/issues/1006))
* [ ] [Feature] Enhance Security Middleware ([hydra#1029](https://github.com/ory/hydra/issues/1029))
* [ ] Add API versioning for administrative APIs ([hydra#1050](https://github.com/ory/hydra/issues/1050))
* [ ] consent: Allow removing tokens without revoking consent ([hydra#1142](https://github.com/ory/hydra/issues/1142)) - [@hackerman](https://github.com/aeneasr)
* [ ] OAuth Client authentication creation CLI jwks client field not present ([hydra#1404](https://github.com/ory/hydra/issues/1404))
* [ ] Add oAuth2Client to logoutRequest similar to loginRequest. ([hydra#1483](https://github.com/ory/hydra/issues/1483))
* [ ] Add a way to filter/sort the list of clients ([hydra#1485](https://github.com/ory/hydra/issues/1485)) - [@hackerman](https://github.com/aeneasr)
* [ ] Remove "not before" claim "nbf" from JWT access token ([hydra#1542](https://github.com/ory/hydra/issues/1542))
* [ ] No way to handle 409 GetLoginRequestConflict. ([hydra#1569](https://github.com/ory/hydra/issues/1569)) - [@Patrik](https://github.com/zepatrik)
* [ ] Auth session cannot be prolonged even if the user is active ([hydra#1690](https://github.com/ory/hydra/issues/1690))
* [ ] Add endpoint to Admin API to revoke access tokens ([hydra#1728](https://github.com/ory/hydra/issues/1728))
* [ ] Migrate to gobuffalo/pop ([hydra#1730](https://github.com/ory/hydra/issues/1730)) - [@Patrik](https://github.com/zepatrik)
* [ ] Rename DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY ([hydra#1760](https://github.com/ory/hydra/issues/1760)) - [@hackerman](https://github.com/aeneasr)
* [ ] CLI Migration Down ([hydra#1763](https://github.com/ory/hydra/issues/1763))
* [ ] Move to go-jose key generation ([hydra#1825](https://github.com/ory/hydra/issues/1825))
* [ ] Make cookies with SameSite=None secure by default or using the configuration flag ([hydra#1844](https://github.com/ory/hydra/issues/1844))
* [ ] Split HTTPS handling for public/admin ([hydra#1962](https://github.com/ory/hydra/issues/1962))
* [ ] Token claims customization with Jsonnet ([hydra#1748](https://github.com/ory/hydra/issues/1748)) - [@hackerman](https://github.com/aeneasr)
* [x] cmd: Add upsert command for client CLI ([hydra#1086](https://github.com/ory/hydra/issues/1086)) - [@hackerman](https://github.com/aeneasr)
* [x] oauth2: Make cleaning up refresh and authz codes possible ([hydra#1130](https://github.com/ory/hydra/issues/1130)) - [@hackerman](https://github.com/aeneasr)

### [Help wanted](https://github.com/ory/hydra/labels/help%20wanted)

We are looking for help on this one.

#### Issues

* [ ] Device Authorization Grant ([hydra#1553](https://github.com/ory/hydra/issues/1553))
* [ ] client_id case sensitivity is not properly enforced when using MySQL ([hydra#1644](https://github.com/ory/hydra/issues/1644)) - [@Patrik](https://github.com/zepatrik)
* [ ] Add endpoint to Admin API to revoke access tokens ([hydra#1728](https://github.com/ory/hydra/issues/1728))
* [ ] Migrate to gobuffalo/pop ([hydra#1730](https://github.com/ory/hydra/issues/1730)) - [@Patrik](https://github.com/zepatrik)
* [ ] CLI Migration Down ([hydra#1763](https://github.com/ory/hydra/issues/1763))
* [ ] Move to go-jose key generation ([hydra#1825](https://github.com/ory/hydra/issues/1825))
* [ ] Introspection Response: `access_token` and `refresh_token` are not valid `token_type` ([hydra#1762](https://github.com/ory/hydra/issues/1762))
* [ ] Split HTTPS handling for public/admin ([hydra#1962](https://github.com/ory/hydra/issues/1962))

### [Rfc](https://github.com/ory/hydra/labels/rfc)

A request for comments to discuss and share ideas.

#### Issues

* [ ] Split HTTPS handling for public/admin ([hydra#1962](https://github.com/ory/hydra/issues/1962))