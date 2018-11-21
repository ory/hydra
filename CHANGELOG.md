<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Change Log](#change-log)
  - [Unreleased](#unreleased)
  - [v1.0.0-rc.2+oryOS.9 (2018-11-21)](#v100-rc2oryos9-2018-11-21)
  - [v1.0.0-rc.1+oryOS.9 (2018-11-21)](#v100-rc1oryos9-2018-11-21)
  - [v1.0.0-beta.9 (2018-09-01)](#v100-beta9-2018-09-01)
  - [v1.0.0-beta.8 (2018-08-10)](#v100-beta8-2018-08-10)
  - [v1.0.0-beta.7 (2018-07-16)](#v100-beta7-2018-07-16)
  - [v1.0.0-beta.6 (2018-07-11)](#v100-beta6-2018-07-11)
  - [v1.0.0-beta.5 (2018-07-07)](#v100-beta5-2018-07-07)
  - [v0.11.14 (2018-06-15)](#v01114-2018-06-15)
  - [v1.0.0-beta.4 (2018-06-13)](#v100-beta4-2018-06-13)
  - [v1.0.0-beta.3 (2018-06-13)](#v100-beta3-2018-06-13)
  - [v1.0.0-beta.2 (2018-05-29)](#v100-beta2-2018-05-29)
  - [v1.0.0-beta.1 (2018-05-29)](#v100-beta1-2018-05-29)
  - [v0.11.12 (2018-04-08)](#v01112-2018-04-08)
  - [v0.11.10 (2018-03-19)](#v01110-2018-03-19)
  - [v0.11.9 (2018-03-10)](#v0119-2018-03-10)
  - [v0.11.7 (2018-03-03)](#v0117-2018-03-03)
  - [v0.11.6 (2018-02-07)](#v0116-2018-02-07)
  - [v0.11.4 (2018-01-23)](#v0114-2018-01-23)
  - [v0.11.3 (2018-01-23)](#v0113-2018-01-23)
  - [v0.11.2 (2018-01-22)](#v0112-2018-01-22)
  - [v0.11.1 (2018-01-18)](#v0111-2018-01-18)
  - [v0.11.0 (2018-01-08)](#v0110-2018-01-08)
  - [v0.10.10 (2017-12-16)](#v01010-2017-12-16)
  - [v0.10.9 (2017-12-13)](#v0109-2017-12-13)
  - [v0.10.8 (2017-12-12)](#v0108-2017-12-12)
  - [v0.10.7 (2017-12-09)](#v0107-2017-12-09)
  - [v0.10.6 (2017-12-09)](#v0106-2017-12-09)
  - [v0.10.5 (2017-12-09)](#v0105-2017-12-09)
  - [v0.10.4 (2017-12-09)](#v0104-2017-12-09)
  - [v0.10.3 (2017-12-08)](#v0103-2017-12-08)
  - [v0.10.2 (2017-12-08)](#v0102-2017-12-08)
  - [v0.10.1 (2017-12-08)](#v0101-2017-12-08)
  - [v0.10.0 (2017-12-08)](#v0100-2017-12-08)
  - [v0.10.0-alpha.21 (2017-11-27)](#v0100-alpha21-2017-11-27)
  - [v0.10.0-alpha.20 (2017-11-26)](#v0100-alpha20-2017-11-26)
  - [v0.10.0-alpha.19 (2017-11-26)](#v0100-alpha19-2017-11-26)
  - [v0.10.0-alpha.18 (2017-11-06)](#v0100-alpha18-2017-11-06)
  - [v0.10.0-alpha.17 (2017-11-06)](#v0100-alpha17-2017-11-06)
  - [v0.10.0-alpha.16 (2017-11-06)](#v0100-alpha16-2017-11-06)
  - [v0.10.0-alpha.15 (2017-11-06)](#v0100-alpha15-2017-11-06)
  - [v0.10.0-alpha.14 (2017-11-06)](#v0100-alpha14-2017-11-06)
  - [v0.10.0-alpha.13 (2017-11-06)](#v0100-alpha13-2017-11-06)
  - [v0.10.0-alpha.11 (2017-11-06)](#v0100-alpha11-2017-11-06)
  - [v0.10.0-alpha.12 (2017-11-06)](#v0100-alpha12-2017-11-06)
  - [v0.10.0-alpha.10 (2017-10-26)](#v0100-alpha10-2017-10-26)
  - [v0.10.0-alpha.9 (2017-10-25)](#v0100-alpha9-2017-10-25)
  - [v0.9.16 (2017-10-23)](#v0916-2017-10-23)
  - [v0.10.0-alpha.8 (2017-10-18)](#v0100-alpha8-2017-10-18)
  - [v0.9.15 (2017-10-11)](#v0915-2017-10-11)
  - [v0.9.14 (2017-10-06)](#v0914-2017-10-06)
  - [v0.10.0-alpha.7 (2017-10-06)](#v0100-alpha7-2017-10-06)
  - [v0.10.0-alpha.6 (2017-10-05)](#v0100-alpha6-2017-10-05)
  - [v0.10.0-alpha.5 (2017-10-05)](#v0100-alpha5-2017-10-05)
  - [v0.10.0-alpha.4 (2017-10-05)](#v0100-alpha4-2017-10-05)
  - [v0.10.0-alpha.3 (2017-10-05)](#v0100-alpha3-2017-10-05)
  - [v0.10.0-alpha.2 (2017-10-05)](#v0100-alpha2-2017-10-05)
  - [v0.10.0-alpha.1 (2017-10-05)](#v0100-alpha1-2017-10-05)
  - [v0.9.13 (2017-09-26)](#v0913-2017-09-26)
  - [v0.9.12 (2017-07-06)](#v0912-2017-07-06)
  - [v0.9.11 (2017-06-30)](#v0911-2017-06-30)
  - [v0.9.10 (2017-06-29)](#v0910-2017-06-29)
  - [v0.9.9 (2017-06-17)](#v099-2017-06-17)
  - [v0.9.8 (2017-06-17)](#v098-2017-06-17)
  - [v0.9.7 (2017-06-16)](#v097-2017-06-16)
  - [v0.9.6 (2017-06-15)](#v096-2017-06-15)
  - [v0.9.5 (2017-06-15)](#v095-2017-06-15)
  - [v0.9.4 (2017-06-14)](#v094-2017-06-14)
  - [v0.9.3 (2017-06-14)](#v093-2017-06-14)
  - [v0.9.2 (2017-06-13)](#v092-2017-06-13)
  - [v0.9.1 (2017-06-12)](#v091-2017-06-12)
  - [v0.9.0 (2017-06-07)](#v090-2017-06-07)
  - [v0.8.7 (2017-06-05)](#v087-2017-06-05)
  - [v0.8.6 (2017-06-05)](#v086-2017-06-05)
  - [v0.8.5 (2017-06-01)](#v085-2017-06-01)
  - [v0.8.4 (2017-05-24)](#v084-2017-05-24)
  - [v0.8.3 (2017-05-23)](#v083-2017-05-23)
  - [v0.8.2 (2017-05-10)](#v082-2017-05-10)
  - [v0.8.1 (2017-05-08)](#v081-2017-05-08)
  - [v0.8.0 (2017-05-07)](#v080-2017-05-07)
  - [v0.7.13 (2017-05-03)](#v0713-2017-05-03)
  - [v0.7.12 (2017-04-30)](#v0712-2017-04-30)
  - [v0.7.11 (2017-04-28)](#v0711-2017-04-28)
  - [v0.7.10 (2017-04-14)](#v0710-2017-04-14)
  - [v0.7.9 (2017-04-02)](#v079-2017-04-02)
  - [v0.7.8 (2017-03-24)](#v078-2017-03-24)
  - [v0.7.7 (2017-02-11)](#v077-2017-02-11)
  - [v0.7.4 (2017-02-11)](#v074-2017-02-11)
  - [v0.7.5 (2017-02-11)](#v075-2017-02-11)
  - [v0.7.6 (2017-02-11)](#v076-2017-02-11)
  - [v0.7.3 (2017-01-22)](#v073-2017-01-22)
  - [v0.7.2 (2017-01-02)](#v072-2017-01-02)
  - [v0.7.1 (2016-12-30)](#v071-2016-12-30)
  - [v0.7.0 (2016-12-30)](#v070-2016-12-30)
  - [v0.6.10 (2016-12-26)](#v0610-2016-12-26)
  - [v0.6.9 (2016-12-20)](#v069-2016-12-20)
  - [v0.6.8 (2016-12-06)](#v068-2016-12-06)
  - [v0.6.7 (2016-12-04)](#v067-2016-12-04)
  - [v0.6.6 (2016-12-04)](#v066-2016-12-04)
  - [v0.6.5 (2016-11-28)](#v065-2016-11-28)
  - [v0.6.4 (2016-11-22)](#v064-2016-11-22)
  - [v0.6.3 (2016-11-17)](#v063-2016-11-17)
  - [v0.6.2 (2016-11-05)](#v062-2016-11-05)
  - [v0.6.1 (2016-10-26)](#v061-2016-10-26)
  - [v0.6.0 (2016-10-25)](#v060-2016-10-25)
  - [v0.5.8 (2016-10-06)](#v058-2016-10-06)
  - [v0.5.7 (2016-10-04)](#v057-2016-10-04)
  - [v0.5.6 (2016-10-03)](#v056-2016-10-03)
  - [v0.5.5 (2016-09-29)](#v055-2016-09-29)
  - [v0.5.4 (2016-09-29)](#v054-2016-09-29)
  - [v0.5.3 (2016-09-29)](#v053-2016-09-29)
  - [v0.5.2 (2016-09-23)](#v052-2016-09-23)
  - [v0.5.0 (2016-09-22)](#v050-2016-09-22)
  - [v0.5.1 (2016-09-22)](#v051-2016-09-22)
  - [v0.4.2-alpha.4 (2016-09-03)](#v042-alpha4-2016-09-03)
  - [v0.4.2 (2016-09-03)](#v042-2016-09-03)
  - [v0.4.3 (2016-09-03)](#v043-2016-09-03)
  - [v0.4.2-alpha.3 (2016-09-02)](#v042-alpha3-2016-09-02)
  - [v0.4.2-alpha.2 (2016-09-01)](#v042-alpha2-2016-09-01)
  - [v0.4.2-alpha.1 (2016-09-01)](#v042-alpha1-2016-09-01)
  - [0.4.2-alpha (2016-09-01)](#042-alpha-2016-09-01)
  - [v0.4.1 (2016-08-18)](#v041-2016-08-18)
  - [v0.4.0 (2016-08-17)](#v040-2016-08-17)
  - [v0.3.1 (2016-08-17)](#v031-2016-08-17)
  - [v0.3.0 (2016-08-09)](#v030-2016-08-09)
  - [v0.2.0 (2016-08-09)](#v020-2016-08-09)
  - [0.1-beta.4 (2016-06-26)](#01-beta4-2016-06-26)
  - [0.1-beta.3 (2016-06-20)](#01-beta3-2016-06-20)
  - [0.1-beta.2 (2016-06-14)](#01-beta2-2016-06-14)
  - [0.1-beta1 (2016-05-29)](#01-beta1-2016-05-29)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Change Log

## [Unreleased](https://github.com/ory/hydra/tree/HEAD)

[Full Changelog](https://github.com/ory/hydra/compare/v1.0.0-rc.2+oryOS.9...HEAD)

**Closed issues:**

- sql: Unable to run migrations when coming from beta.9 [\#1185](https://github.com/ory/hydra/issues/1185)

## [v1.0.0-rc.2+oryOS.9](https://github.com/ory/hydra/tree/v1.0.0-rc.2+oryOS.9) (2018-11-21)
[Full Changelog](https://github.com/ory/hydra/compare/v1.0.0-rc.1+oryOS.9...v1.0.0-rc.2+oryOS.9)

**Merged pull requests:**

- sql: Resolve beta.9 -\> rc.1 migration issue [\#1186](https://github.com/ory/hydra/pull/1186) ([aeneasr](https://github.com/aeneasr))

## [v1.0.0-rc.1+oryOS.9](https://github.com/ory/hydra/tree/v1.0.0-rc.1+oryOS.9) (2018-11-21)
[Full Changelog](https://github.com/ory/hydra/compare/v1.0.0-beta.9...v1.0.0-rc.1+oryOS.9)

**Implemented enhancements:**

- cmd: `token user` should be able to set up ssl [\#1147](https://github.com/ory/hydra/issues/1147)
- client: Deleting a client should delete all associated data too [\#1131](https://github.com/ory/hydra/issues/1131)
- Use `-mod=vendor` when building binaries / docker [\#1112](https://github.com/ory/hydra/issues/1112)
- Switch to go mod [\#1074](https://github.com/ory/hydra/issues/1074)
- CORS\_ALLOWED\_ORIGINS doesn't respect wildcards [\#1073](https://github.com/ory/hydra/issues/1073)
- consent: Add authorize code URL to consent and login response payloads [\#1046](https://github.com/ory/hydra/issues/1046)
- \[Feature Request\] Update consent tests to match oauth2/client tests [\#1043](https://github.com/ory/hydra/issues/1043)
- cmd/server: Export useful bootstrap function [\#973](https://github.com/ory/hydra/issues/973)
- sdk: C\# language SDK [\#958](https://github.com/ory/hydra/issues/958)
- Opentracing tracing integration [\#931](https://github.com/ory/hydra/issues/931)
- consent: Add ability to specify Access Token Audience [\#883](https://github.com/ory/hydra/issues/883)
- Prepare v1.0.0-rc.1 release [\#1175](https://github.com/ory/hydra/pull/1175) ([aeneasr](https://github.com/aeneasr))
- vendor: Update fosite to 0.27.3 [\#1164](https://github.com/ory/hydra/pull/1164) ([aeneasr](https://github.com/aeneasr))
- sdk: Document userinfo as GET instead of POST [\#1161](https://github.com/ory/hydra/pull/1161) ([aeneasr](https://github.com/aeneasr))
- oauth2: Add audience and improve refresh flow [\#1156](https://github.com/ory/hydra/pull/1156) ([aeneasr](https://github.com/aeneasr))
- cmd: Improve issuer error message [\#1152](https://github.com/ory/hydra/pull/1152) ([aeneasr](https://github.com/aeneasr))
- oauth2: Add OAuth2 audience claim and improve migrations [\#1145](https://github.com/ory/hydra/pull/1145) ([aeneasr](https://github.com/aeneasr))
- Switch to go modules [\#1077](https://github.com/ory/hydra/pull/1077) ([aeneasr](https://github.com/aeneasr))
- cmd: Fix flaky port finder [\#1076](https://github.com/ory/hydra/pull/1076) ([aeneasr](https://github.com/aeneasr))
- rand: Fix flaky random test [\#1075](https://github.com/ory/hydra/pull/1075) ([aeneasr](https://github.com/aeneasr))

**Fixed bugs:**

- tracing: sql args are added as tags when they should be omitted [\#1181](https://github.com/ory/hydra/issues/1181)
- consent: Require proof of authentication before ending user session [\#1154](https://github.com/ory/hydra/issues/1154)
- oauth2: Audience is potentially not being refreshed [\#1153](https://github.com/ory/hydra/issues/1153)
- Hydra shut down after a race condition [\#1141](https://github.com/ory/hydra/issues/1141)
- oauth2: Tables oidc, code, openid, refresh are missing indices [\#1140](https://github.com/ory/hydra/issues/1140)
- consent: SQL field `subject\_obfuscated` does not have an index [\#1138](https://github.com/ory/hydra/issues/1138)
- Setting up a fresh hydra installation results in panic [\#1137](https://github.com/ory/hydra/issues/1137)
- Copy-paste error in manager\_0\_sql\_migrations\_test.go [\#1135](https://github.com/ory/hydra/issues/1135)
- cmd: Error message regarding IssuerURL should contain environment variable name [\#1133](https://github.com/ory/hydra/issues/1133)
- client: Deleting a client should delete all associated data too [\#1131](https://github.com/ory/hydra/issues/1131)
- CORS\\_ALLOWED\\_ORIGINS doesn't respect wildcards [\#1073](https://github.com/ory/hydra/issues/1073)
- OpenID configuration endpoint returns wrong registration endpoint [\#1072](https://github.com/ory/hydra/issues/1072)
- OAuth2 Token Revoke call results in 404 Not Found [\#1070](https://github.com/ory/hydra/issues/1070)
- Missing database indices [\#1067](https://github.com/ory/hydra/issues/1067)
- Use PKCE with hybrid flow [\#1060](https://github.com/ory/hydra/issues/1060)
- cmd: Consent timeout is currently hardcoded but environment variable exists [\#1057](https://github.com/ory/hydra/issues/1057)
- ACR claim not being set on id token when requested by login accept request [\#1032](https://github.com/ory/hydra/issues/1032)
- List all consent sessions returns 404 [\#1031](https://github.com/ory/hydra/issues/1031)
- Introspect endpoint reports expiration time for refresh tokens [\#1025](https://github.com/ory/hydra/issues/1025)
- sql: Resolve index/fk regression issues [\#1178](https://github.com/ory/hydra/pull/1178) ([aeneasr](https://github.com/aeneasr))
- Prepare v1.0.0-rc.1 release [\#1175](https://github.com/ory/hydra/pull/1175) ([aeneasr](https://github.com/aeneasr))
- consent: Ignore row count in revoke [\#1173](https://github.com/ory/hydra/pull/1173) ([aeneasr](https://github.com/aeneasr))
- vendor: Upgrade to fosite 0.27.4 [\#1171](https://github.com/ory/hydra/pull/1171) ([aeneasr](https://github.com/aeneasr))
- vendor: Update fosite to 0.27.3 [\#1164](https://github.com/ory/hydra/pull/1164) ([aeneasr](https://github.com/aeneasr))
- consent: Properly propagate acr value [\#1160](https://github.com/ory/hydra/pull/1160) ([aeneasr](https://github.com/aeneasr))
- cmd: Resolve broken wildcard cors [\#1159](https://github.com/ory/hydra/pull/1159) ([aeneasr](https://github.com/aeneasr))
- cmd: Resolve panic in migration handler [\#1151](https://github.com/ory/hydra/pull/1151) ([aeneasr](https://github.com/aeneasr))
- consent: Only fetch latest consent state  [\#1124](https://github.com/ory/hydra/pull/1124) ([aeneasr](https://github.com/aeneasr))
- server: Instantiate PKCE after oidc [\#1123](https://github.com/ory/hydra/pull/1123) ([aeneasr](https://github.com/aeneasr))
- cli: Improve migrate error messages  [\#1080](https://github.com/ory/hydra/pull/1080) ([aeneasr](https://github.com/aeneasr))
- cmd: Fix flaky port finder [\#1076](https://github.com/ory/hydra/pull/1076) ([aeneasr](https://github.com/aeneasr))

**Closed issues:**

- Resolve regression issues related to foreign keys [\#1177](https://github.com/ory/hydra/issues/1177)
- DELETE `/oauth2/auth/sessions/login/{user}` returns 404 [\#1168](https://github.com/ory/hydra/issues/1168)
- How to authenticate with POST /clients endpoint [\#1148](https://github.com/ory/hydra/issues/1148)
- Implementation of user idel time sout  [\#1146](https://github.com/ory/hydra/issues/1146)
- Move SQL migrations to files and improve test pipeline [\#1144](https://github.com/ory/hydra/issues/1144)
- cmd: Show error hint in oauth2 error view [\#1143](https://github.com/ory/hydra/issues/1143)
- Login time deteriorates over time [\#1119](https://github.com/ory/hydra/issues/1119)
- why hydra-login-consent-go didn't work, is there will have login provider and consent provider with golang? [\#1117](https://github.com/ory/hydra/issues/1117)
- Intro Blog source code is unreadable [\#1111](https://github.com/ory/hydra/issues/1111)
- consent: ignores extra claims for id and access token [\#1106](https://github.com/ory/hydra/issues/1106)
- Invalid\_request while generate the Access token in own OAuth 2.0 server [\#1104](https://github.com/ory/hydra/issues/1104)
- Invalid\_request while generate the Access token in own OAuth 2.0 server [\#1103](https://github.com/ory/hydra/issues/1103)
- Document query parameters for /oauth2/auth [\#1100](https://github.com/ory/hydra/issues/1100)
- PHP SDK is not PSR-4 compliant [\#1099](https://github.com/ory/hydra/issues/1099)
- CHALLENGE\_TOKEN\_LIFESPAN unused [\#1097](https://github.com/ory/hydra/issues/1097)
- Improve follow-up on numerous ORY repos [\#1093](https://github.com/ory/hydra/issues/1093)
- Run your own OAuth 2.0 Server : " Client authentication failed "  [\#1091](https://github.com/ory/hydra/issues/1091)
- govet cmd/tooken\_user.go: the cancel function returned by context.WithTimeout should be called [\#1090](https://github.com/ory/hydra/issues/1090)
- Enhancement: specify lifespan for refresh\_token [\#1088](https://github.com/ory/hydra/issues/1088)
- Add at\_hash claim to id\_token in code flow. [\#1085](https://github.com/ory/hydra/issues/1085)
- Disable https://api.segment.io POST request [\#1083](https://github.com/ory/hydra/issues/1083)
- Move internal dependencies to ory/x [\#1081](https://github.com/ory/hydra/issues/1081)
- Support Kubernetes Secrets [\#1079](https://github.com/ory/hydra/issues/1079)
- Silent token refresh fails with "The Authorization Server requires End-User consent" [\#1068](https://github.com/ory/hydra/issues/1068)
- Invalid login\_challenge [\#1065](https://github.com/ory/hydra/issues/1065)
- Problem with import path for go-resty and go1.11 modules [\#1063](https://github.com/ory/hydra/issues/1063)
- sql: Add auto-increment PKs [\#1059](https://github.com/ory/hydra/issues/1059)
- Feature: admin endpoint for deleting expired tokens [\#1058](https://github.com/ory/hydra/issues/1058)
- consent: Send error response if consent or login challenge is expired or invalid [\#1056](https://github.com/ory/hydra/issues/1056)
- consent: Add original request URL to login and consent request payloads [\#1055](https://github.com/ory/hydra/issues/1055)
- Fix flaky random-port generator [\#1054](https://github.com/ory/hydra/issues/1054)
- Fix flaky pseudo-random test [\#1053](https://github.com/ory/hydra/issues/1053)
- API doc: GET /userinfo works but not documented [\#1049](https://github.com/ory/hydra/issues/1049)
- go SDK userInfo response does not support extra claims [\#1048](https://github.com/ory/hydra/issues/1048)
- Issuer url is allways fallowed by / even when defined without [\#1041](https://github.com/ory/hydra/issues/1041)
- missing end\_session\_endpoint from .well-known doc [\#1040](https://github.com/ory/hydra/issues/1040)
- oryd/hydra:v1.0.0-beta.9 clients api return 404 [\#1036](https://github.com/ory/hydra/issues/1036)
- DELETE login/{user} and  DELETE consent/{user} can not redirect to Login page [\#1035](https://github.com/ory/hydra/issues/1035)
- remember in  requests/login/{challenge}/accept api cause get same subject always [\#1034](https://github.com/ory/hydra/issues/1034)
- Out of Band OAuth2 Authorization [\#1033](https://github.com/ory/hydra/issues/1033)
- \[Cleanup\] CORS Settings [\#1028](https://github.com/ory/hydra/issues/1028)
- Key rotation leads to "Could not fetch private signing key for OpenID Connect" [\#1026](https://github.com/ory/hydra/issues/1026)

**Merged pull requests:**

- More e2e tests [\#1184](https://github.com/ory/hydra/pull/1184) ([aeneasr](https://github.com/aeneasr))
- fix migrate sql command  at upgrading guide [\#1183](https://github.com/ory/hydra/pull/1183) ([sawadashota](https://github.com/sawadashota))
- rc.1 release preparations [\#1182](https://github.com/ory/hydra/pull/1182) ([aeneasr](https://github.com/aeneasr))
- e2e: Improve e2e test pipeline [\#1180](https://github.com/ory/hydra/pull/1180) ([aeneasr](https://github.com/aeneasr))
- docs: Auto-generate appendix [\#1174](https://github.com/ory/hydra/pull/1174) ([aeneasr](https://github.com/aeneasr))
- vendor: Upgrade to fosite 0.28.0 [\#1172](https://github.com/ory/hydra/pull/1172) ([aeneasr](https://github.com/aeneasr))
- ci: Generate benchmarks in docus format [\#1170](https://github.com/ory/hydra/pull/1170) ([aeneasr](https://github.com/aeneasr))
- ci: Update release pipeline for new versioning [\#1169](https://github.com/ory/hydra/pull/1169) ([aeneasr](https://github.com/aeneasr))
- oauth2: Make client registration endpoint configurable [\#1167](https://github.com/ory/hydra/pull/1167) ([aeneasr](https://github.com/aeneasr))
- sdk: Update swagger endpoint definition [\#1166](https://github.com/ory/hydra/pull/1166) ([aeneasr](https://github.com/aeneasr))
- sql: Add missing indices [\#1157](https://github.com/ory/hydra/pull/1157) ([aeneasr](https://github.com/aeneasr))
- cmd: Add ability to specify consent and login lifespan [\#1155](https://github.com/ory/hydra/pull/1155) ([aeneasr](https://github.com/aeneasr))
- cmd: Add https option to token user command [\#1150](https://github.com/ory/hydra/pull/1150) ([aeneasr](https://github.com/aeneasr))
- cmd: Improve token user error handling [\#1149](https://github.com/ory/hydra/pull/1149) ([aeneasr](https://github.com/aeneasr))
- Minor bug fix in JWK sql migrations test case [\#1136](https://github.com/ory/hydra/pull/1136) ([jacor84](https://github.com/jacor84))
- tracing: remove bad tracing config from docker-compose.yml [\#1132](https://github.com/ory/hydra/pull/1132) ([aaslamin](https://github.com/aaslamin))
- cmd: Resolve issues with secret migration [\#1129](https://github.com/ory/hydra/pull/1129) ([aeneasr](https://github.com/aeneasr))
- health: Register healthx.AliveCheckPath route for frontend [\#1128](https://github.com/ory/hydra/pull/1128) ([jayme-github](https://github.com/jayme-github))
- consent: Set fetch order to descending [\#1126](https://github.com/ory/hydra/pull/1126) ([aeneasr](https://github.com/aeneasr))
- cors: add options cors middleware handler [\#1125](https://github.com/ory/hydra/pull/1125) ([JiaLiPassion](https://github.com/JiaLiPassion))
- ci: Check vet and fix vet errors [\#1122](https://github.com/ory/hydra/pull/1122) ([aeneasr](https://github.com/aeneasr))
- jwks: cors for wellknown endpoints [\#1118](https://github.com/ory/hydra/pull/1118) ([JiaLiPassion](https://github.com/JiaLiPassion))
- oauth2: wellknown should use corsMiddleware [\#1116](https://github.com/ory/hydra/pull/1116) ([JiaLiPassion](https://github.com/JiaLiPassion))
- tracing: add support for tracing db interactions [\#1115](https://github.com/ory/hydra/pull/1115) ([aaslamin](https://github.com/aaslamin))
- build: Improve build pipeline [\#1114](https://github.com/ory/hydra/pull/1114) ([aeneasr](https://github.com/aeneasr))
- e2e: Check for access/id token claims [\#1113](https://github.com/ory/hydra/pull/1113) ([aeneasr](https://github.com/aeneasr))
- sdk/js: Declare opencollective as devdep [\#1109](https://github.com/ory/hydra/pull/1109) ([aeneasr](https://github.com/aeneasr))
- Fix missing LoginChallenge and LoginSessionID from GetConsentRequest [\#1105](https://github.com/ory/hydra/pull/1105) ([jcxplorer](https://github.com/jcxplorer))
- Update README - Benchmarks section [\#1102](https://github.com/ory/hydra/pull/1102) ([kishaningithub](https://github.com/kishaningithub))
- docs: Updates issue and pull request templates [\#1101](https://github.com/ory/hydra/pull/1101) ([aeneasr](https://github.com/aeneasr))
- Add error response if consent or login challenge is expired [\#1098](https://github.com/ory/hydra/pull/1098) ([k-lepa](https://github.com/k-lepa))
- docs: Updates issue and pull request templates [\#1096](https://github.com/ory/hydra/pull/1096) ([aeneasr](https://github.com/aeneasr))
- Move dependencies to ory/x [\#1095](https://github.com/ory/hydra/pull/1095) ([aeneasr](https://github.com/aeneasr))
- docs: Updates issue and pull request templates [\#1094](https://github.com/ory/hydra/pull/1094) ([aeneasr](https://github.com/aeneasr))
- Add schema changes introduced to UPGRADE.md [\#1082](https://github.com/ory/hydra/pull/1082) ([aaslamin](https://github.com/aaslamin))
- sql: Add auto-increment PKs [\#1078](https://github.com/ory/hydra/pull/1078) ([aeneasr](https://github.com/aeneasr))
- tracing: use context aware database methods [\#1071](https://github.com/ory/hydra/pull/1071) ([aaslamin](https://github.com/aaslamin))
- Add missing indices to resolve \#1067 [\#1069](https://github.com/ory/hydra/pull/1069) ([aaslamin](https://github.com/aaslamin))
- change go-resty import path for gopkg.in/resty.v1 [\#1064](https://github.com/ory/hydra/pull/1064) ([pierredavidbelanger](https://github.com/pierredavidbelanger))
- fosite: bump to version 0.24.0 with associated code changes [\#1062](https://github.com/ory/hydra/pull/1062) ([someone1](https://github.com/someone1))
- Bump fosite version to 0.23.0 + New tracing instrumented Hasher [\#1052](https://github.com/ory/hydra/pull/1052) ([aaslamin](https://github.com/aaslamin))
- consent: migrate to test helpers \[closes \#1043\] [\#1051](https://github.com/ory/hydra/pull/1051) ([someone1](https://github.com/someone1))
- Fix swagger [\#1045](https://github.com/ory/hydra/pull/1045) ([pierredavidbelanger](https://github.com/pierredavidbelanger))
- client: fix test to pass non-nil context [\#1044](https://github.com/ory/hydra/pull/1044) ([someone1](https://github.com/someone1))
- Bump fosite version and integrate breaking changes [\#1042](https://github.com/ory/hydra/pull/1042) ([aaslamin](https://github.com/aaslamin))
- two littles things that bugs me when I compile or run tests [\#1039](https://github.com/ory/hydra/pull/1039) ([pierredavidbelanger](https://github.com/pierredavidbelanger))
- cmd: Do not echo secrets if explicitly set [\#1038](https://github.com/ory/hydra/pull/1038) ([aeneasr](https://github.com/aeneasr))
- propagate context through to the sql store [\#1030](https://github.com/ory/hydra/pull/1030) ([aaslamin](https://github.com/aaslamin))
- consent: Add SessionsPath const [\#1027](https://github.com/ory/hydra/pull/1027) ([someone1](https://github.com/someone1))
- Use latest version of sqlcon [\#1024](https://github.com/ory/hydra/pull/1024) ([davidjwilkins](https://github.com/davidjwilkins))
- cmd/server: Export Handler bootstrap functions \(\#973\) [\#1023](https://github.com/ory/hydra/pull/1023) ([someone1](https://github.com/someone1))
- Add support for distributed tracing [\#1019](https://github.com/ory/hydra/pull/1019) ([aaslamin](https://github.com/aaslamin))

## [v1.0.0-beta.9](https://github.com/ory/hydra/tree/v1.0.0-beta.9) (2018-09-01)
[Full Changelog](https://github.com/ory/hydra/compare/v1.0.0-beta.8...v1.0.0-beta.9)

**Implemented enhancements:**

- Duplicate entry error for second consent request [\#1007](https://github.com/ory/hydra/issues/1007)
- consent: Share session state between login and consent [\#1003](https://github.com/ory/hydra/issues/1003)
- cmd: Print version when booting up [\#987](https://github.com/ory/hydra/issues/987)
- client: client specific CORS settings [\#957](https://github.com/ory/hydra/issues/957)
- cmd: Add cli helper for importing and exporting environments \(clients, policies, keys\) [\#699](https://github.com/ory/hydra/issues/699)
- sql: jsonb support for postgres [\#516](https://github.com/ory/hydra/issues/516)
- client: filter oauth2 clients by field through REST API [\#505](https://github.com/ory/hydra/issues/505)
- cmd: Allow SYSTEM\_SECRET key rotation [\#73](https://github.com/ory/hydra/issues/73)
- consent: Forward session and login information [\#1013](https://github.com/ory/hydra/pull/1013) ([aeneasr](https://github.com/aeneasr))
- jwk: Add ability to rotate SYSTEM\_SECRET  [\#1012](https://github.com/ory/hydra/pull/1012) ([aeneasr](https://github.com/aeneasr))
- vendor: Upgrade sqlcon to 0.0.6 [\#1008](https://github.com/ory/hydra/pull/1008) ([aeneasr](https://github.com/aeneasr))
- cmd: Use viper for cors detection [\#998](https://github.com/ory/hydra/pull/998) ([aeneasr](https://github.com/aeneasr))
- cmd: Disable CORS by default [\#997](https://github.com/ory/hydra/pull/997) ([aeneasr](https://github.com/aeneasr))
- cmd: Add version to banner [\#995](https://github.com/ory/hydra/pull/995) ([aeneasr](https://github.com/aeneasr))
- sdk: Add new methods to SDK interface [\#994](https://github.com/ory/hydra/pull/994) ([aeneasr](https://github.com/aeneasr))

**Fixed bugs:**

- Client creation gives incorrect error message [\#1016](https://github.com/ory/hydra/issues/1016)
- oauth2: id\_token\_hint should work with expired ID tokens [\#1014](https://github.com/ory/hydra/issues/1014)
- cors: Don't automatically auto-allow CORS [\#996](https://github.com/ory/hydra/issues/996)
- Use ID\_TOKEN\_LIFESPAN when doing refresh [\#985](https://github.com/ory/hydra/issues/985)
- MySQL/MariDB broken on default Debian installations [\#377](https://github.com/ory/hydra/issues/377)
- cmd: Clarify HYDRA\_ADMIN\_URL in missing endpoint message [\#1018](https://github.com/ory/hydra/pull/1018) ([aeneasr](https://github.com/aeneasr))
- oauth2: Accept expired JWTs as id\_token\_hint [\#1017](https://github.com/ory/hydra/pull/1017) ([aeneasr](https://github.com/aeneasr))
- cmd: Disable CORS by default [\#997](https://github.com/ory/hydra/pull/997) ([aeneasr](https://github.com/aeneasr))
- consent: Populate consent session with default values [\#989](https://github.com/ory/hydra/pull/989) ([aeneasr](https://github.com/aeneasr))

**Closed issues:**

- cmd: Replace cors fork with upstream [\#1010](https://github.com/ory/hydra/issues/1010)
- Auth State mismatch. URL Double Encoding [\#1005](https://github.com/ory/hydra/issues/1005)
- Can not remember consent because no user interaction was required with resp\['skip'\] false [\#999](https://github.com/ory/hydra/issues/999)
- invalid if condition about SubjectTypesSupport [\#992](https://github.com/ory/hydra/issues/992)
- sdk: add oauthapi functions to golang interface [\#991](https://github.com/ory/hydra/issues/991)
- After redirecting from consent -- runtime error: invalid memory address or nil pointer dereference [\#988](https://github.com/ory/hydra/issues/988)

**Merged pull requests:**

- docker: Update compose definitions [\#1020](https://github.com/ory/hydra/pull/1020) ([aeneasr](https://github.com/aeneasr))
- config: Fix use of uninitialized logger [\#1015](https://github.com/ory/hydra/pull/1015) ([vHanda](https://github.com/vHanda))
- cmd: Replace aeneasr/cors with rs/cors [\#1011](https://github.com/ory/hydra/pull/1011) ([aeneasr](https://github.com/aeneasr))
- oauth2: Enable client specific CORS settings [\#1009](https://github.com/ory/hydra/pull/1009) ([aeneasr](https://github.com/aeneasr))
- oauth2: Resolve broken expiry when refreshing id token [\#1002](https://github.com/ory/hydra/pull/1002) ([aeneasr](https://github.com/aeneasr))
- Delete Procfile [\#1001](https://github.com/ory/hydra/pull/1001) ([MOZGIII](https://github.com/MOZGIII))
- Fix serve all cmd in docker files [\#1000](https://github.com/ory/hydra/pull/1000) ([condemil](https://github.com/condemil))
- cmd: Public subject type should cause public id alg [\#993](https://github.com/ory/hydra/pull/993) ([aeneasr](https://github.com/aeneasr))
- config: disable plugin backend through 'noplugin' tag  [\#986](https://github.com/ory/hydra/pull/986) ([glerchundi](https://github.com/glerchundi))

## [v1.0.0-beta.8](https://github.com/ory/hydra/tree/v1.0.0-beta.8) (2018-08-10)
[Full Changelog](https://github.com/ory/hydra/compare/v1.0.0-beta.7...v1.0.0-beta.8)

**Implemented enhancements:**

- Allow logging out and deleting a single session cookie [\#970](https://github.com/ory/hydra/issues/970)
- vendor: Upgrade to MySQL 1.4 driver [\#965](https://github.com/ory/hydra/issues/965)
- oauth2: abstract oauth2/handler JWT Strategies [\#960](https://github.com/ory/hydra/issues/960)
- consent: expose  a list of all clients authorized by a user [\#953](https://github.com/ory/hydra/issues/953)
- oauth2: Support for Pairwise Subject Identifier Type [\#950](https://github.com/ory/hydra/issues/950)
- \[Enhancement/Proposal\] Update Plugin System [\#949](https://github.com/ory/hydra/issues/949)
- The JWK api should be able to export .pem [\#175](https://github.com/ory/hydra/issues/175)
- cmd: Add flags for new client fields in create [\#939](https://github.com/ory/hydra/issues/939)
- client: Deprecate the `public` flag [\#938](https://github.com/ory/hydra/issues/938)
- client: Clarify error message regarding client auth method [\#936](https://github.com/ory/hydra/issues/936)
- cmd: Add option to specify new oidc parameters in client [\#935](https://github.com/ory/hydra/issues/935)
- consent: Obtain previously selected scopes [\#902](https://github.com/ory/hydra/issues/902)
- oauth2: allow issuing of JWT access tokens [\#248](https://github.com/ory/hydra/issues/248)
- oauth2: Add scope to introspection test suite [\#941](https://github.com/ory/hydra/pull/941) ([aeneasr](https://github.com/aeneasr))
-  consent: Add logout api endpoint [\#984](https://github.com/ory/hydra/pull/984) ([aeneasr](https://github.com/aeneasr))
- sdk: Upgrade superagent to 3.7.0 [\#983](https://github.com/ory/hydra/pull/983) ([aeneasr](https://github.com/aeneasr))
- vendor: Upgrade to latest sqlcon [\#975](https://github.com/ory/hydra/pull/975) ([aeneasr](https://github.com/aeneasr))
- oauth2: Refactor JWT strategy [\#972](https://github.com/ory/hydra/pull/972) ([someone1](https://github.com/someone1))
- oauth2: Removes authorization from introspection [\#969](https://github.com/ory/hydra/pull/969) ([aeneasr](https://github.com/aeneasr))
- oauth2: Support for Pairwise Subject Identifier Type [\#966](https://github.com/ory/hydra/pull/966) ([aeneasr](https://github.com/aeneasr))
- cmd: Introduce public and administrative ports [\#963](https://github.com/ory/hydra/pull/963) ([aeneasr](https://github.com/aeneasr))
- oauth2: Adds JWT Access Token strategy  [\#947](https://github.com/ory/hydra/pull/947) ([aeneasr](https://github.com/aeneasr))
- oauth2: Improve token endpoint authentication error message [\#942](https://github.com/ory/hydra/pull/942) ([aeneasr](https://github.com/aeneasr))

**Fixed bugs:**

- client: Improve error messages from managers [\#976](https://github.com/ory/hydra/issues/976)
- consent: Duplicate row error should return a better error message [\#880](https://github.com/ory/hydra/issues/880)
- oauth2: error\_hint, error\_debug are not shared when redirect fails [\#974](https://github.com/ory/hydra/issues/974)
- oauth2: Introspect response is empty when `active` is false. [\#964](https://github.com/ory/hydra/issues/964)
- consent: MemoryManager should return `errNoPreviousConsentFound` when no previous consent was found [\#959](https://github.com/ory/hydra/issues/959)
- consent: Auth session should check for `pkg.ErrNotFound`, not `sql.ErrNoRows` [\#944](https://github.com/ory/hydra/issues/944)
- sdk: Add AdminURL and PublicURL to configuration [\#968](https://github.com/ory/hydra/pull/968) ([aeneasr](https://github.com/aeneasr))
- cmd: Introduce public and administrative ports [\#963](https://github.com/ory/hydra/pull/963) ([aeneasr](https://github.com/aeneasr))
- consent: Properly identify revoked login sessions [\#945](https://github.com/ory/hydra/pull/945) ([aeneasr](https://github.com/aeneasr))

**Closed issues:**

- Refresh token and access token share same lifetime [\#955](https://github.com/ory/hydra/issues/955)
- Id\_token\_hint doesn't work as expected [\#951](https://github.com/ory/hydra/issues/951)
- consent: Check if helper rejects unknown JSON fields [\#940](https://github.com/ory/hydra/issues/940)
- Unable to specify a custom claim to hydra [\#937](https://github.com/ory/hydra/issues/937)
- \[HTTP API\] get /version returns empty [\#934](https://github.com/ory/hydra/issues/934)
- docs: Add limitations section [\#839](https://github.com/ory/hydra/issues/839)
- Expose administrative APIs at a different port \(e.g. 4445\) [\#904](https://github.com/ory/hydra/issues/904)

**Merged pull requests:**

- client: Improve memory manager error messages [\#978](https://github.com/ory/hydra/pull/978) ([aeneasr](https://github.com/aeneasr))
- consent: Add ListUserConsentSessions to OAuth2API interface [\#977](https://github.com/ory/hydra/pull/977) ([clausdenk](https://github.com/clausdenk))
- docker: Update .dockerignore [\#967](https://github.com/ory/hydra/pull/967) ([aeneasr](https://github.com/aeneasr))
- cli: fix reporting of epected vs. received status codes [\#961](https://github.com/ory/hydra/pull/961) ([rjw57](https://github.com/rjw57))
- all: Introduce database backend interface and update plugin system an… [\#956](https://github.com/ory/hydra/pull/956) ([someone1](https://github.com/someone1))
- Add api endpoint to list all authorized clients by user [\#954](https://github.com/ory/hydra/pull/954) ([kingjan1999](https://github.com/kingjan1999))
- Use spdx expression for license in package.json [\#952](https://github.com/ory/hydra/pull/952) ([kingjan1999](https://github.com/kingjan1999))
- Improve client API compatibility with oidc dynamic discovery [\#943](https://github.com/ory/hydra/pull/943) ([aeneasr](https://github.com/aeneasr))
- oauth2: Share error details with redirect fallback [\#982](https://github.com/ory/hydra/pull/982) ([aeneasr](https://github.com/aeneasr))
- cli: Print "active:false" when token is inactive [\#981](https://github.com/ory/hydra/pull/981) ([aeneasr](https://github.com/aeneasr))
- consent: Return proper error when no consent was found [\#980](https://github.com/ory/hydra/pull/980) ([aeneasr](https://github.com/aeneasr))
- vendor: Upgrade sqlcon to 0.0.5 [\#979](https://github.com/ory/hydra/pull/979) ([aeneasr](https://github.com/aeneasr))

## [v1.0.0-beta.7](https://github.com/ory/hydra/tree/v1.0.0-beta.7) (2018-07-16)
[Full Changelog](https://github.com/ory/hydra/compare/v1.0.0-beta.6...v1.0.0-beta.7)

**Implemented enhancements:**

- Panic when calling oauth2/auth/sessions/consent/{user} or oauth2/auth/sessions/consent/{user}/{client} [\#928](https://github.com/ory/hydra/issues/928)
- client: Improve handling of legacy `id` field [\#927](https://github.com/ory/hydra/pull/927) ([aeneasr](https://github.com/aeneasr))

**Fixed bugs:**

- Panic when calling oauth2/auth/sessions/consent/{user} or oauth2/auth/sessions/consent/{user}/{client} [\#928](https://github.com/ory/hydra/issues/928)
- jwk: Auto-remove old keys when upgrading from \< beta.7 [\#925](https://github.com/ory/hydra/pull/925) ([aeneasr](https://github.com/aeneasr))

**Closed issues:**

- migration 0.11.10 \> 1.0 : did you forget to run hydra migrate sql" or forget to set the SYSTEM\_SECRET [\#926](https://github.com/ory/hydra/issues/926)
- ClientID property is ignored when creating a new OAuth2 Client [\#924](https://github.com/ory/hydra/issues/924)
- The CSRF value from the token does not match the CSRF value from the data store [\#923](https://github.com/ory/hydra/issues/923)
- Which version is stable? [\#922](https://github.com/ory/hydra/issues/922)
- JSON Web Key Store default keys broken after upgrading to beta.6 [\#921](https://github.com/ory/hydra/issues/921)

**Merged pull requests:**

- Document that ORY Hydra is OpenID Certified [\#933](https://github.com/ory/hydra/pull/933) ([aeneasr](https://github.com/aeneasr))
- cmd: Show error when loading x509 cert fails [\#932](https://github.com/ory/hydra/pull/932) ([aeneasr](https://github.com/aeneasr))
- Allow cookie without max age [\#930](https://github.com/ory/hydra/pull/930) ([BastianHofmann](https://github.com/BastianHofmann))
- cmd: Check dependencies are defined before instantiation [\#929](https://github.com/ory/hydra/pull/929) ([aeneasr](https://github.com/aeneasr))
- README: fix docker linux link [\#920](https://github.com/ory/hydra/pull/920) ([philips](https://github.com/philips))

## [v1.0.0-beta.6](https://github.com/ory/hydra/tree/v1.0.0-beta.6) (2018-07-11)
[Full Changelog](https://github.com/ory/hydra/compare/v1.0.0-beta.5...v1.0.0-beta.6)

**Implemented enhancements:**

- consent: Add endpoint to revoke authentication and consent sessions [\#856](https://github.com/ory/hydra/issues/856)
- jwk: improve JWK tests [\#588](https://github.com/ory/hydra/issues/588)
- cli/clients: allow to import multiple clients with one file [\#388](https://github.com/ory/hydra/issues/388)
- oauth2: allow token revocation without knowing the token \(i.e. per user\) [\#304](https://github.com/ory/hydra/issues/304)
- cmd: CLI should be able to import PEM keys to JWK store [\#98](https://github.com/ory/hydra/issues/98)

**Fixed bugs:**

- migration 0.9.x -\> 1.0: sector\_identifier\_uri contains null values  [\#918](https://github.com/ory/hydra/issues/918)

**Closed issues:**

- Hydra version 0.11.13-alpine break cli [\#917](https://github.com/ory/hydra/issues/917)
- health: Check if and why the health endpoint returns a HTTPS response [\#879](https://github.com/ory/hydra/issues/879)
- docs: disallow secrets from docs/tutorials in production mode [\#573](https://github.com/ory/hydra/issues/573)

**Merged pull requests:**

- client: Fix sql migration step for oidc [\#919](https://github.com/ory/hydra/pull/919) ([aeneasr](https://github.com/aeneasr))
- cmd: Allows import of PEM/DER/JSON encoded keys  [\#916](https://github.com/ory/hydra/pull/916) ([aeneasr](https://github.com/aeneasr))

## [v1.0.0-beta.5](https://github.com/ory/hydra/tree/v1.0.0-beta.5) (2018-07-07)
[Full Changelog](https://github.com/ory/hydra/compare/v0.11.14...v1.0.0-beta.5)

**Implemented enhancements:**

- client: Improve and DRY validation in handler [\#909](https://github.com/ory/hydra/issues/909)
- cmd/server: Die when system secret is in wrong format [\#817](https://github.com/ory/hydra/issues/817)
- OpenID Connect Certification [\#689](https://github.com/ory/hydra/issues/689)

**Fixed bugs:**

- Public and private key pair fetched from store does not match [\#912](https://github.com/ory/hydra/issues/912)
- 500 error returned on GET /clients/{id} when client doesn't exist [\#903](https://github.com/ory/hydra/issues/903)
- metrics: Properly handle metrics log messages [\#833](https://github.com/ory/hydra/issues/833)

**Closed issues:**

- go get return error [\#913](https://github.com/ory/hydra/issues/913)
- Can't create clients using the CLI [\#911](https://github.com/ory/hydra/issues/911)
- is hydra can build on window ? [\#910](https://github.com/ory/hydra/issues/910)
- Let's improve the docs! [\#385](https://github.com/ory/hydra/issues/385)
- Add benchmarks to documentation [\#161](https://github.com/ory/hydra/issues/161)

**Merged pull requests:**

- consent: Adds ability to revoke consent and login sessions [\#915](https://github.com/ory/hydra/pull/915) ([aeneasr](https://github.com/aeneasr))
- jwk: Tests for simple equality in JWT strategy  [\#914](https://github.com/ory/hydra/pull/914) ([aeneasr](https://github.com/aeneasr))
- Adds OpenID Connect Dynamic Client Registration [\#908](https://github.com/ory/hydra/pull/908) ([aeneasr](https://github.com/aeneasr))
- docs: Adds link to examples repository [\#907](https://github.com/ory/hydra/pull/907) ([aeneasr](https://github.com/aeneasr))
- docs: Removes obsolete issue template [\#906](https://github.com/ory/hydra/pull/906) ([aeneasr](https://github.com/aeneasr))

## [v0.11.14](https://github.com/ory/hydra/tree/v0.11.14) (2018-06-15)
[Full Changelog](https://github.com/ory/hydra/compare/v1.0.0-beta.4...v0.11.14)

**Fixed bugs:**

- Missing commits between v0.11.10 and v0.11.12 [\#894](https://github.com/ory/hydra/issues/894)

## [v1.0.0-beta.4](https://github.com/ory/hydra/tree/v1.0.0-beta.4) (2018-06-13)
[Full Changelog](https://github.com/ory/hydra/compare/v1.0.0-beta.3...v1.0.0-beta.4)

## [v1.0.0-beta.3](https://github.com/ory/hydra/tree/v1.0.0-beta.3) (2018-06-13)
[Full Changelog](https://github.com/ory/hydra/compare/v1.0.0-beta.2...v1.0.0-beta.3)

**Implemented enhancements:**

- cmd: Allows reading database from env in migrate sql [\#898](https://github.com/ory/hydra/pull/898) ([aeneasr](https://github.com/aeneasr))

**Fixed bugs:**

- oidc\_context empty [\#900](https://github.com/ory/hydra/issues/900)
-  consent: Propagates oidc\_context to consent request [\#901](https://github.com/ory/hydra/pull/901) ([aeneasr](https://github.com/aeneasr))

**Closed issues:**

- cmd: Add flag to allow reading database url in migration command from env [\#896](https://github.com/ory/hydra/issues/896)

**Merged pull requests:**

- ci: Stops benchmark result commit & pushes [\#905](https://github.com/ory/hydra/pull/905) ([aeneasr](https://github.com/aeneasr))
- docs: Adds CI benchmarks  [\#897](https://github.com/ory/hydra/pull/897) ([aeneasr](https://github.com/aeneasr))
- all: Moves to metrics-middleware [\#895](https://github.com/ory/hydra/pull/895) ([aeneasr](https://github.com/aeneasr))

## [v1.0.0-beta.2](https://github.com/ory/hydra/tree/v1.0.0-beta.2) (2018-05-29)
[Full Changelog](https://github.com/ory/hydra/compare/v1.0.0-beta.1...v1.0.0-beta.2)

**Closed issues:**

- 1.0.0-alpha.1 Release Notes [\#885](https://github.com/ory/hydra/issues/885)

**Merged pull requests:**

- ci: Improves build toolchain [\#893](https://github.com/ory/hydra/pull/893) ([aeneasr](https://github.com/aeneasr))

## [v1.0.0-beta.1](https://github.com/ory/hydra/tree/v1.0.0-beta.1) (2018-05-29)
[Full Changelog](https://github.com/ory/hydra/compare/v0.11.12...v1.0.0-beta.1)

**Implemented enhancements:**

- oauth2: Revoke tokens when performing refreshing grant [\#889](https://github.com/ory/hydra/issues/889)
- docs: Explicitly document in upgrade guide that hydra is no longer protected by default [\#888](https://github.com/ory/hydra/issues/888)
- Extend status page to check dependencies. [\#887](https://github.com/ory/hydra/issues/887)
- oauth2: Revoke previous and future access tokens when revoking a token [\#884](https://github.com/ory/hydra/issues/884)
- consent: Investigate if prompt=none should be allowed with implicit flows [\#866](https://github.com/ory/hydra/issues/866)
- consent: Implement login\_hint capabilities [\#860](https://github.com/ory/hydra/issues/860)
- consent: Always remove session if rememberLogin=false [\#859](https://github.com/ory/hydra/issues/859)
- consent: Resolve broken time out [\#852](https://github.com/ory/hydra/issues/852)
- oauth2: Support max\_age [\#851](https://github.com/ory/hydra/issues/851)
- consent: Include id\_token\_hint in oidc context [\#850](https://github.com/ory/hydra/issues/850)
- health: Document prometheus endpoint [\#844](https://github.com/ory/hydra/issues/844)
- config: Deprecate `ClusterURL`, `ClientID`, `ClientSecret` [\#841](https://github.com/ory/hydra/issues/841)
- oauth2: Return token type on token introspection [\#831](https://github.com/ory/hydra/issues/831)
- oauth2: Support id\_token\_hint at authorization endpoint [\#826](https://github.com/ory/hydra/issues/826)
- consent app: Restart consent flow [\#809](https://github.com/ory/hydra/issues/809)
- oauth2: Allow multiple audience claims on ID token [\#790](https://github.com/ory/hydra/issues/790)
- client: Add field `client\_secret\_expires\_at` to create [\#778](https://github.com/ory/hydra/issues/778)
- all: All JSON output/input should be using `\_` instead of camelCase [\#777](https://github.com/ory/hydra/issues/777)
- oauth2: Reject authorization requests for invalid scopes before redirecting to consent endpoint [\#776](https://github.com/ory/hydra/issues/776)
- oauth2: Improving the consent flow design [\#772](https://github.com/ory/hydra/issues/772)
- oauth2: Expire consent request on successful consent interaction [\#771](https://github.com/ory/hydra/issues/771)
- health: Add ability to retrieve version \(protected endpoint\) [\#743](https://github.com/ory/hydra/issues/743)
- Deprecate `hydra policies create -f` [\#708](https://github.com/ory/hydra/issues/708)
- Disallow unknown JSON fields [\#707](https://github.com/ory/hydra/issues/707)
- oauth2: Remember authentication and application authorization [\#697](https://github.com/ory/hydra/issues/697)
- oauth2: Revoke access and refresh tokens when authorization code is used twice [\#693](https://github.com/ory/hydra/issues/693)
- oauth2: Require consent for OAuth 2.0 public clients [\#692](https://github.com/ory/hydra/issues/692)
- oauth2: Reintroduce audience claim [\#687](https://github.com/ory/hydra/issues/687)
- policy: evaluate wildcard matching strategy [\#580](https://github.com/ory/hydra/issues/580)
- installer: homebrew recipe for macOS users [\#572](https://github.com/ory/hydra/issues/572)
- Warden group metadata [\#387](https://github.com/ory/hydra/issues/387)
- policy: search policies by subject and resource [\#362](https://github.com/ory/hydra/issues/362)
- warden: check against multiple policies [\#264](https://github.com/ory/hydra/issues/264)
- core: add warden context everywhere [\#238](https://github.com/ory/hydra/issues/238)
- better and more e2e tests [\#192](https://github.com/ory/hydra/issues/192)
- Health and test improvements [\#891](https://github.com/ory/hydra/pull/891) ([aeneasr](https://github.com/aeneasr))
- Resolves various issues related to OAuth2 [\#890](https://github.com/ory/hydra/pull/890) ([aeneasr](https://github.com/aeneasr))
- Improve oidc conformity [\#876](https://github.com/ory/hydra/pull/876) ([aeneasr](https://github.com/aeneasr))
- Improves compatibility with OIDC Conformity Tests [\#873](https://github.com/ory/hydra/pull/873) ([aeneasr](https://github.com/aeneasr))
- sdk: Remove the need for OAuth2 credentials [\#869](https://github.com/ory/hydra/pull/869) ([aeneasr](https://github.com/aeneasr))
- Minor improvements [\#868](https://github.com/ory/hydra/pull/868) ([aeneasr](https://github.com/aeneasr))
- consent: Always bust auth session if remember is false [\#864](https://github.com/ory/hydra/pull/864) ([aeneasr](https://github.com/aeneasr))
- oauth2: Returns token type on introspection [\#832](https://github.com/ory/hydra/pull/832) ([aeneasr](https://github.com/aeneasr))

**Fixed bugs:**

- Incorrect CORS-related env vars parsing [\#886](https://github.com/ory/hydra/issues/886)
- consent: Remove the client secret from consent/login response [\#878](https://github.com/ory/hydra/issues/878)
- oauth2: ID Token must be returned in both authorize and token response in hybrid flows with response type `code` [\#875](https://github.com/ory/hydra/issues/875)
- consent: On first prompt=none after authentication, times mismatch [\#874](https://github.com/ory/hydra/issues/874)
- oauth2: Reject requests without nonce unless using the code flow [\#867](https://github.com/ory/hydra/issues/867)
- oauth2: max\_age fails if max\_age=1 [\#862](https://github.com/ory/hydra/issues/862)
- oauth2: Figure out why MySQL tests are flaky on CI [\#861](https://github.com/ory/hydra/issues/861)
- oauth2: Resolve broken prompt parameter [\#843](https://github.com/ory/hydra/issues/843)
- oauth2: Duplicate requests to /oauth2/token cause 500 [\#828](https://github.com/ory/hydra/issues/828)
- consent app: Restart consent flow [\#809](https://github.com/ory/hydra/issues/809)
- Hydra connect fails when the client secret contains "%" [\#631](https://github.com/ory/hydra/issues/631)
- Health and test improvements [\#891](https://github.com/ory/hydra/pull/891) ([aeneasr](https://github.com/aeneasr))
- Resolves various issues related to OAuth2 [\#890](https://github.com/ory/hydra/pull/890) ([aeneasr](https://github.com/aeneasr))
- Improves OpenID Connect Conformity [\#882](https://github.com/ory/hydra/pull/882) ([aeneasr](https://github.com/aeneasr))
- Improve oidc conformity [\#876](https://github.com/ory/hydra/pull/876) ([aeneasr](https://github.com/aeneasr))
- cmd: Adds jwt strategy and fixes nil pointer exception [\#865](https://github.com/ory/hydra/pull/865) ([aeneasr](https://github.com/aeneasr))

**Closed issues:**

- consent: Authentication session cookie invalidation scenarios [\#855](https://github.com/ory/hydra/issues/855)
- consent: Investigate if failure during consent should cause session to be revoked [\#854](https://github.com/ory/hydra/issues/854)
- Please support Type Definition \(d.ts\) for typescript.  [\#848](https://github.com/ory/hydra/issues/848)
- security: add HttpOnly cookie flag [\#847](https://github.com/ory/hydra/issues/847)
- cmd: Deprecate `hydra connect` and replace with per-command flags and environment variables [\#840](https://github.com/ory/hydra/issues/840)
- REST API /clients limit & offset bug [\#838](https://github.com/ory/hydra/issues/838)
- Allow configuring consent URL per client [\#837](https://github.com/ory/hydra/issues/837)
- Duplicate client creation results in 500 [\#835](https://github.com/ory/hydra/issues/835)
- Error 1406: Data too long for column 'subject' at row 1 [\#829](https://github.com/ory/hydra/issues/829)
- Does warden groups work with internal Hydra APIs? [\#823](https://github.com/ory/hydra/issues/823)
- Hydra sdk error hydra.introspectOauth2Token is not a function [\#822](https://github.com/ory/hydra/issues/822)
- Improve the lint percentage [\#818](https://github.com/ory/hydra/issues/818)
- docs: Refactor examples / tutorials [\#810](https://github.com/ory/hydra/issues/810)
- Moving the access control engine to Oathkeeper [\#807](https://github.com/ory/hydra/issues/807)
- Can you build an identity provider with hydra or not? [\#789](https://github.com/ory/hydra/issues/789)
- docker: Add image capable of loading policies/clients/jwks from an init.d directory [\#760](https://github.com/ory/hydra/issues/760)
- Add PUT method for /warden/groups/:id [\#745](https://github.com/ory/hydra/issues/745)
- Document that the install guide is different from the 5 minute guide [\#718](https://github.com/ory/hydra/issues/718)
- Prometheus metrics [\#669](https://github.com/ory/hydra/issues/669)
- docs: Port numbers from docker compose and the lengthy tutorial do not match [\#653](https://github.com/ory/hydra/issues/653)
- docs: add subject + id mocks in the policy section of the swagger specs for each endpoint [\#614](https://github.com/ory/hydra/issues/614)
- docs: /warden/allowed do not fully specify security parameters [\#565](https://github.com/ory/hydra/issues/565)
- docs: explain oauth2 better [\#356](https://github.com/ory/hydra/issues/356)
- docs: have a "running hydra in production" section [\#354](https://github.com/ory/hydra/issues/354)
- docs: clarify that the consent app is responsible for implementing full OIDC [\#353](https://github.com/ory/hydra/issues/353)
- docs: add auth0 seminar to docs [\#347](https://github.com/ory/hydra/issues/347)
- docs: add bug bounty section to readme [\#84](https://github.com/ory/hydra/issues/84)
- docs: add passport.js real-world example [\#83](https://github.com/ory/hydra/issues/83)

**Merged pull requests:**

- vendor: Upgrades fosite dependency [\#892](https://github.com/ory/hydra/pull/892) ([aeneasr](https://github.com/aeneasr))
- Minor consent improvements [\#881](https://github.com/ory/hydra/pull/881) ([aeneasr](https://github.com/aeneasr))
- oauth2: Ignores JTI in userinfo [\#877](https://github.com/ory/hydra/pull/877) ([aeneasr](https://github.com/aeneasr))
- oauth2: Rejects requests without nonce in implicit/hybrid [\#872](https://github.com/ory/hydra/pull/872) ([aeneasr](https://github.com/aeneasr))
- Improves health endpoints and cleans up code [\#871](https://github.com/ory/hydra/pull/871) ([aeneasr](https://github.com/aeneasr))
- Client secret expires [\#870](https://github.com/ory/hydra/pull/870) ([zepatrik](https://github.com/zepatrik))
- Fix mysql timing bug [\#863](https://github.com/ory/hydra/pull/863) ([aeneasr](https://github.com/aeneasr))
- consent: Removes stray fmt.Print [\#858](https://github.com/ory/hydra/pull/858) ([aeneasr](https://github.com/aeneasr))
- Improves consent flow [\#857](https://github.com/ory/hydra/pull/857) ([aeneasr](https://github.com/aeneasr))
- Resolves issues with auth\_time [\#853](https://github.com/ory/hydra/pull/853) ([aeneasr](https://github.com/aeneasr))
- add /health/version endpoint [\#845](https://github.com/ory/hydra/pull/845) ([zepatrik](https://github.com/zepatrik))
- Deprecate connect [\#842](https://github.com/ory/hydra/pull/842) ([aeneasr](https://github.com/aeneasr))
- Move policy merged [\#830](https://github.com/ory/hydra/pull/830) ([aeneasr](https://github.com/aeneasr))
- \[Prometheus\] Add new prometheus metrics and metrics endpoint [\#827](https://github.com/ory/hydra/pull/827) ([dolbik](https://github.com/dolbik))
- 1.0.x [\#825](https://github.com/ory/hydra/pull/825) ([aeneasr](https://github.com/aeneasr))
- Merge from 0.11.x [\#824](https://github.com/ory/hydra/pull/824) ([aeneasr](https://github.com/aeneasr))

## [v0.11.12](https://github.com/ory/hydra/tree/v0.11.12) (2018-04-08)
[Full Changelog](https://github.com/ory/hydra/compare/v0.11.10...v0.11.12)

**Fixed bugs:**

- sdk: PHP sdk missing from releases [\#781](https://github.com/ory/hydra/issues/781)

**Closed issues:**

- Special characters in redirect url [\#819](https://github.com/ory/hydra/issues/819)
- "Could not fetch signing key for OpenID Connect" [\#816](https://github.com/ory/hydra/issues/816)

**Merged pull requests:**

- Resolves dep and tests issues [\#821](https://github.com/ory/hydra/pull/821) ([aeneasr](https://github.com/aeneasr))
- oauth2: Resolves client secrets from potentially leaking to the database in cleartext [\#820](https://github.com/ory/hydra/pull/820) ([aeneasr](https://github.com/aeneasr))
- Activating Open Collective [\#805](https://github.com/ory/hydra/pull/805) ([monkeywithacupcake](https://github.com/monkeywithacupcake))
- metrics: Improves naming of traits  [\#804](https://github.com/ory/hydra/pull/804) ([aeneasr](https://github.com/aeneasr))
- 0.11 [\#796](https://github.com/ory/hydra/pull/796) ([aeneasr](https://github.com/aeneasr))

## [v0.11.10](https://github.com/ory/hydra/tree/v0.11.10) (2018-03-19)
[Full Changelog](https://github.com/ory/hydra/compare/v0.11.9...v0.11.10)

**Closed issues:**

- docs: Link to php sdk README is wrong [\#811](https://github.com/ory/hydra/issues/811)

**Merged pull requests:**

- Minor code cleanup [\#815](https://github.com/ory/hydra/pull/815) ([euank](https://github.com/euank))
- docs: Resolves broken swagger definitions [\#812](https://github.com/ory/hydra/pull/812) ([aeneasr](https://github.com/aeneasr))
- docs: Updates banner in readme [\#808](https://github.com/ory/hydra/pull/808) ([aeneasr](https://github.com/aeneasr))
- Update links to discord and readme [\#806](https://github.com/ory/hydra/pull/806) ([aeneasr](https://github.com/aeneasr))

## [v0.11.9](https://github.com/ory/hydra/tree/v0.11.9) (2018-03-10)
[Full Changelog](https://github.com/ory/hydra/compare/v0.11.7...v0.11.9)

**Implemented enhancements:**

- telemetry: Add version and build info as custom dimensions [\#802](https://github.com/ory/hydra/issues/802)
- docs: Adds redirects for broken guide links [\#798](https://github.com/ory/hydra/pull/798) ([aeneasr](https://github.com/aeneasr))

**Fixed bugs:**

- id\_token not returned after request at the /oauth2/token endpoint using the refresh\_token [\#794](https://github.com/ory/hydra/issues/794)
- docker: Build time always return time.Now\(\) [\#792](https://github.com/ory/hydra/issues/792)
- cmd: Resolves an issue with broken build time display [\#799](https://github.com/ory/hydra/pull/799) ([aeneasr](https://github.com/aeneasr))
- cmd: Adds OpenID Connect refresh handler [\#797](https://github.com/ory/hydra/pull/797) ([aeneasr](https://github.com/aeneasr))

**Closed issues:**

- docs: document difference between scopes and policies [\#590](https://github.com/ory/hydra/issues/590)

**Merged pull requests:**

- metrics: Improves naming of traits [\#803](https://github.com/ory/hydra/pull/803) ([aeneasr](https://github.com/aeneasr))
- docs: Resolves broken images and build [\#801](https://github.com/ory/hydra/pull/801) ([aeneasr](https://github.com/aeneasr))
- docs: Moves documentation to new repository. [\#800](https://github.com/ory/hydra/pull/800) ([aeneasr](https://github.com/aeneasr))
- all: Updates license headers [\#793](https://github.com/ory/hydra/pull/793) ([aeneasr](https://github.com/aeneasr))

## [v0.11.7](https://github.com/ory/hydra/tree/v0.11.7) (2018-03-03)
[Full Changelog](https://github.com/ory/hydra/compare/v0.11.6...v0.11.7)

**Implemented enhancements:**

- make --skip-newsletter the default [\#779](https://github.com/ory/hydra/issues/779)
- group: Add pagination to group management [\#741](https://github.com/ory/hydra/issues/741)
- jwk: Add pagination to jwk lists [\#740](https://github.com/ory/hydra/issues/740)
- client: Add pagination to client list [\#739](https://github.com/ory/hydra/issues/739)
- ConsentRequest should use time.Now\(\).UTC\(\) for ExpiresAt. [\#679](https://github.com/ory/hydra/issues/679)
- sdk: add python sdk [\#639](https://github.com/ory/hydra/issues/639)
- Importing a client should fail when an unrecognized field is found [\#357](https://github.com/ory/hydra/issues/357)
- ci: Automatically pushes docs to website [\#784](https://github.com/ory/hydra/pull/784) ([aeneasr](https://github.com/aeneasr))
- oauth2: Forces UTC in consent strategy [\#775](https://github.com/ory/hydra/pull/775) ([aeneasr](https://github.com/aeneasr))
- client: Introduces pagination to client management [\#774](https://github.com/ory/hydra/pull/774) ([aeneasr](https://github.com/aeneasr))

**Fixed bugs:**

- oauth2: Remove exp and iat from ID token header [\#787](https://github.com/ory/hydra/issues/787)
- Don't push to coveralls in CI when PR comes from fork [\#782](https://github.com/ory/hydra/issues/782)
- policy: List tests do not care about offset/limit - fix that [\#746](https://github.com/ory/hydra/issues/746)

**Closed issues:**

- A way to skip the consent screen for certain clients \(first party\) [\#791](https://github.com/ory/hydra/issues/791)
- Where's the tutorial? [\#788](https://github.com/ory/hydra/issues/788)
- Feature Request: oauth2/token endpoint json payload option [\#786](https://github.com/ory/hydra/issues/786)
- docs: Deprecate recovering root access section [\#756](https://github.com/ory/hydra/issues/756)
- oauth2: Document how to make the well known endpoint public [\#688](https://github.com/ory/hydra/issues/688)
- oauth2: replace redirect uri exact match with protocol/host/path match [\#257](https://github.com/ory/hydra/issues/257)

**Merged pull requests:**

- docs: Adds automatic summary and toc generation [\#785](https://github.com/ory/hydra/pull/785) ([aeneasr](https://github.com/aeneasr))
- Remove coveralls token from circleci config [\#783](https://github.com/ory/hydra/pull/783) ([zepatrik](https://github.com/zepatrik))
- Update newsletter text [\#780](https://github.com/ory/hydra/pull/780) ([zepatrik](https://github.com/zepatrik))
- Minor improvements to the gitbook guide [\#773](https://github.com/ory/hydra/pull/773) ([aeneasr](https://github.com/aeneasr))

## [v0.11.6](https://github.com/ory/hydra/tree/v0.11.6) (2018-02-07)
[Full Changelog](https://github.com/ory/hydra/compare/v0.11.4...v0.11.6)

**Implemented enhancements:**

- server: Add default policy for well-known/jwks.json [\#761](https://github.com/ory/hydra/issues/761)
- cmd: Add newsletter info and sign up [\#755](https://github.com/ory/hydra/issues/755)
- metrics: Improve metrics endpoint [\#742](https://github.com/ory/hydra/issues/742)
- oauth2: Add ability to purge old access tokens [\#738](https://github.com/ory/hydra/issues/738)
- jwk: refactor jwk id generation [\#589](https://github.com/ory/hydra/issues/589)
- oauth2: Adds support for PKCE \(IETF RFC7636\)  [\#769](https://github.com/ory/hydra/pull/769) ([aeneasr](https://github.com/aeneasr))
- Forces unique JWK IDs and allows anonymous access to ./well-known/jwks.json [\#762](https://github.com/ory/hydra/pull/762) ([aeneasr](https://github.com/aeneasr))

**Fixed bugs:**

- Do not show client secret when client is public in CLI [\#737](https://github.com/ory/hydra/issues/737)
- oauth2: Client secret error message should be shown on creation [\#725](https://github.com/ory/hydra/issues/725)
- sdk: Resolves composer license complaint [\#763](https://github.com/ory/hydra/pull/763) ([aeneasr](https://github.com/aeneasr))

**Closed issues:**

- docker-compose encountered errors [\#758](https://github.com/ory/hydra/issues/758)
- AWS Lambda Support? [\#749](https://github.com/ory/hydra/issues/749)
- cmd/client: Ask for security newsletter sign up when using client side CLI [\#747](https://github.com/ory/hydra/issues/747)
- oauth2: Add PKCE support [\#744](https://github.com/ory/hydra/issues/744)

**Merged pull requests:**

- Gen php sdk [\#814](https://github.com/ory/hydra/pull/814) ([pnicolcev-tulipretail](https://github.com/pnicolcev-tulipretail))
- oauth2: Resolves possible session fixation attack [\#770](https://github.com/ory/hydra/pull/770) ([aeneasr](https://github.com/aeneasr))
- docs: Fix dead link to example policy [\#767](https://github.com/ory/hydra/pull/767) ([gr-eg](https://github.com/gr-eg))
- Purge tokens [\#766](https://github.com/ory/hydra/pull/766) ([aeneasr](https://github.com/aeneasr))
- client: do not show/send secret when client is public [\#765](https://github.com/ory/hydra/pull/765) ([zepatrik](https://github.com/zepatrik))
- fix \#725 [\#764](https://github.com/ory/hydra/pull/764) ([zepatrik](https://github.com/zepatrik))
- Cmd newsletter signup [\#759](https://github.com/ory/hydra/pull/759) ([aeneasr](https://github.com/aeneasr))
- sdk: Generate php sdk and point php autoloader to lib folder [\#736](https://github.com/ory/hydra/pull/736) ([pnicolcev-tulipretail](https://github.com/pnicolcev-tulipretail))

## [v0.11.4](https://github.com/ory/hydra/tree/v0.11.4) (2018-01-23)
[Full Changelog](https://github.com/ory/hydra/compare/v0.11.3...v0.11.4)

## [v0.11.3](https://github.com/ory/hydra/tree/v0.11.3) (2018-01-23)
[Full Changelog](https://github.com/ory/hydra/compare/v0.11.2...v0.11.3)

**Implemented enhancements:**

- Improve telemetry module [\#752](https://github.com/ory/hydra/pull/752) ([aeneasr](https://github.com/aeneasr))

**Closed issues:**

- possible consent session id attack? [\#753](https://github.com/ory/hydra/issues/753)

## [v0.11.2](https://github.com/ory/hydra/tree/v0.11.2) (2018-01-22)
[Full Changelog](https://github.com/ory/hydra/compare/v0.11.1...v0.11.2)

**Fixed bugs:**

- client: Returns 404 only when policy allows getting a client [\#751](https://github.com/ory/hydra/pull/751) ([aeneasr](https://github.com/aeneasr))

**Merged pull requests:**

- oauth2: Protects consent flow against session fixation [\#754](https://github.com/ory/hydra/pull/754) ([aeneasr](https://github.com/aeneasr))

## [v0.11.1](https://github.com/ory/hydra/tree/v0.11.1) (2018-01-18)
[Full Changelog](https://github.com/ory/hydra/compare/v0.11.0...v0.11.1)

**Implemented enhancements:**

- groups: Add ability to list all groups, not just by member [\#729](https://github.com/ory/hydra/issues/729)

**Fixed bugs:**

-  Resolves issues with pagination [\#750](https://github.com/ory/hydra/pull/750) ([aeneasr](https://github.com/aeneasr))

**Closed issues:**

- Timezone Issue with new consent flow in 0.10? [\#735](https://github.com/ory/hydra/issues/735)
- policies: change effect type from string to boolean [\#666](https://github.com/ory/hydra/issues/666)
- cmd: `hydra connect --url` should work with and without trailing slash [\#650](https://github.com/ory/hydra/issues/650)

**Merged pull requests:**

- add a save way to get the ClusterURL and append to it [\#748](https://github.com/ory/hydra/pull/748) ([zepatrik](https://github.com/zepatrik))

## [v0.11.0](https://github.com/ory/hydra/tree/v0.11.0) (2018-01-08)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.10...v0.11.0)

**Implemented enhancements:**

- group: List groups without owner [\#732](https://github.com/ory/hydra/issues/732)
- Add an alias for offline scope called offline\_access [\#722](https://github.com/ory/hydra/issues/722)
- oauth2: Print debug message to logs and evaluate transmitting it to clients too [\#715](https://github.com/ory/hydra/issues/715)
- groups: Add ability to list all groups, not just by member [\#734](https://github.com/ory/hydra/pull/734) ([aeneasr](https://github.com/aeneasr))
- sdk: Adds php registry dummy [\#733](https://github.com/ory/hydra/pull/733) ([aeneasr](https://github.com/aeneasr))
- oauth2: Prints debug message to logs and evaluate transmitting it to clients too [\#727](https://github.com/ory/hydra/pull/727) ([aeneasr](https://github.com/aeneasr))
- vendor: Adds offline\_access scope alias [\#724](https://github.com/ory/hydra/pull/724) ([aeneasr](https://github.com/aeneasr))

**Fixed bugs:**

- health: Should not require x-forwarded-proto [\#726](https://github.com/ory/hydra/issues/726)
- health: Stop requiring x-forwarded-proto [\#731](https://github.com/ory/hydra/pull/731) ([aeneasr](https://github.com/aeneasr))

**Closed issues:**

- variable part in the subject and resource in ladon policy to be filled by request [\#730](https://github.com/ory/hydra/issues/730)
- Trailing slash redirect strips directories from path [\#723](https://github.com/ory/hydra/issues/723)
- Resolve broken docker-compose tutorial guide [\#717](https://github.com/ory/hydra/issues/717)
- Document external dependencies [\#716](https://github.com/ory/hydra/issues/716)

**Merged pull requests:**

- docs: Adds documentation on third-party deps [\#728](https://github.com/ory/hydra/pull/728) ([aeneasr](https://github.com/aeneasr))

## [v0.10.10](https://github.com/ory/hydra/tree/v0.10.10) (2017-12-16)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.9...v0.10.10)

**Implemented enhancements:**

- Make scopes in `hydra token client` command configurable [\#711](https://github.com/ory/hydra/issues/711)
- cmd: Makes scopes in token command configurable [\#712](https://github.com/ory/hydra/pull/712) ([aeneasr](https://github.com/aeneasr))
- cmd: Adds a dedicated command for importing policies [\#709](https://github.com/ory/hydra/pull/709) ([aeneasr](https://github.com/aeneasr))

**Fixed bugs:**

- Misleading error message when using the SDK [\#686](https://github.com/ory/hydra/issues/686)
- sdk/go: Resolves incorrect error message [\#713](https://github.com/ory/hydra/pull/713) ([aeneasr](https://github.com/aeneasr))

**Closed issues:**

- Docker readme, in case it is lost [\#719](https://github.com/ory/hydra/issues/719)
- Keep track of version and build hash [\#706](https://github.com/ory/hydra/issues/706)
- Scope is documented as hydra.groups but should by hydra.warden.groups [\#702](https://github.com/ory/hydra/issues/702)
- Rename `hydra policies create -f` to `hydra policies import` [\#701](https://github.com/ory/hydra/issues/701)

**Merged pull requests:**

- docs: Resolves issue with broken 5-minute tutorial [\#721](https://github.com/ory/hydra/pull/721) ([aeneasr](https://github.com/aeneasr))
- Improves userinfo endpoint [\#714](https://github.com/ory/hydra/pull/714) ([aeneasr](https://github.com/aeneasr))
- groups: Corrects group scope documentation [\#710](https://github.com/ory/hydra/pull/710) ([aeneasr](https://github.com/aeneasr))

## [v0.10.9](https://github.com/ory/hydra/tree/v0.10.9) (2017-12-13)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.8...v0.10.9)

**Implemented enhancements:**

- Reintroduce alpine based image with shell [\#703](https://github.com/ory/hydra/issues/703)

**Merged pull requests:**

- pkg: Fixes returning nil instead of empty array in split [\#705](https://github.com/ory/hydra/pull/705) ([aeneasr](https://github.com/aeneasr))

## [v0.10.8](https://github.com/ory/hydra/tree/v0.10.8) (2017-12-12)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.7...v0.10.8)

**Implemented enhancements:**

- oauth2: Add token\_endpoint\_auth\_methods\_supported to openid-configuration [\#695](https://github.com/ory/hydra/issues/695)

**Closed issues:**

- docs: Add introspect bc to upgrade [\#698](https://github.com/ory/hydra/issues/698)

**Merged pull requests:**

- Reintroduces alpine based docker image [\#704](https://github.com/ory/hydra/pull/704) ([aeneasr](https://github.com/aeneasr))

## [v0.10.7](https://github.com/ory/hydra/tree/v0.10.7) (2017-12-09)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.6...v0.10.7)

## [v0.10.6](https://github.com/ory/hydra/tree/v0.10.6) (2017-12-09)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.5...v0.10.6)

**Closed issues:**

- oauth2: Write test for userinfo endpoint without token and test for 401 [\#691](https://github.com/ory/hydra/issues/691)

**Merged pull requests:**

- Improves OpenID Connect conformity [\#694](https://github.com/ory/hydra/pull/694) ([aeneasr](https://github.com/aeneasr))

## [v0.10.5](https://github.com/ory/hydra/tree/v0.10.5) (2017-12-09)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.4...v0.10.5)

**Closed issues:**

- oauth2: Support userinfo endpoint [\#652](https://github.com/ory/hydra/issues/652)

## [v0.10.4](https://github.com/ory/hydra/tree/v0.10.4) (2017-12-09)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.3...v0.10.4)

**Merged pull requests:**

- oauth2: Adds basic userinfo endpoint [\#690](https://github.com/ory/hydra/pull/690) ([aeneasr](https://github.com/aeneasr))

## [v0.10.3](https://github.com/ory/hydra/tree/v0.10.3) (2017-12-08)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.2...v0.10.3)

## [v0.10.2](https://github.com/ory/hydra/tree/v0.10.2) (2017-12-08)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.1...v0.10.2)

## [v0.10.1](https://github.com/ory/hydra/tree/v0.10.1) (2017-12-08)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.0...v0.10.1)

**Implemented enhancements:**

- Open source policy naming guidelines [\#680](https://github.com/ory/hydra/issues/680)

**Closed issues:**

- docs: docker --link should be replaced by networks [\#555](https://github.com/ory/hydra/issues/555)

## [v0.10.0](https://github.com/ory/hydra/tree/v0.10.0) (2017-12-08)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.0-alpha.21...v0.10.0)

**Implemented enhancements:**

- docs: Improve release and breaking changes management [\#675](https://github.com/ory/hydra/issues/675)
- oauth2: Make sub explicit in the database [\#658](https://github.com/ory/hydra/issues/658)
- oauth2: Add access control to token introspection endpoint [\#655](https://github.com/ory/hydra/issues/655)
- all: make policy resource and action names configurable [\#640](https://github.com/ory/hydra/issues/640)
- Subject field [\#674](https://github.com/ory/hydra/pull/674) ([aeneasr](https://github.com/aeneasr))
- Add changelog [\#673](https://github.com/ory/hydra/pull/673) ([aeneasr](https://github.com/aeneasr))

**Fixed bugs:**

- oauth2: Token revokation should check client id before revoking tokens [\#676](https://github.com/ory/hydra/issues/676)
- cli/policies: removing a policy subject adds the subject Instead [\#662](https://github.com/ory/hydra/issues/662)
- jwk: Rename ES521 key generation algorithm to ES512 [\#651](https://github.com/ory/hydra/issues/651)
- oauth2: Fixes clients being able to revoke any token [\#677](https://github.com/ory/hydra/pull/677) ([aeneasr](https://github.com/aeneasr))

**Closed issues:**

- Json logging [\#670](https://github.com/ory/hydra/issues/670)
- swagger: scope pattern requires a space [\#661](https://github.com/ory/hydra/issues/661)
- docs: Add list of undisclosed adopters with requests ranges to readme [\#659](https://github.com/ory/hydra/issues/659)

**Merged pull requests:**

- Update release notes and prepare 0.10.0 [\#685](https://github.com/ory/hydra/pull/685) ([aeneasr](https://github.com/aeneasr))
- docs: Adds multi-tenant best practices [\#684](https://github.com/ory/hydra/pull/684) ([aeneasr](https://github.com/aeneasr))
- ci: Resolves code climate issues [\#683](https://github.com/ory/hydra/pull/683) ([aeneasr](https://github.com/aeneasr))
- pkg: Adds test for LogError [\#682](https://github.com/ory/hydra/pull/682) ([aeneasr](https://github.com/aeneasr))
- docs: Adds ACP best practices [\#681](https://github.com/ory/hydra/pull/681) ([aeneasr](https://github.com/aeneasr))
- oauth2: Requires firewall check for introspecting access tokens [\#678](https://github.com/ory/hydra/pull/678) ([aeneasr](https://github.com/aeneasr))
- Makes policy resource names prefixes configurable [\#672](https://github.com/ory/hydra/pull/672) ([aeneasr](https://github.com/aeneasr))
- docs: Adds consent state machine [\#671](https://github.com/ory/hydra/pull/671) ([aeneasr](https://github.com/aeneasr))
- docs: Make space optional in scope regex \(\#661\) [\#668](https://github.com/ory/hydra/pull/668) ([pnicolcev-tulipretail](https://github.com/pnicolcev-tulipretail))
- Various minor fixes [\#667](https://github.com/ory/hydra/pull/667) ([aeneasr](https://github.com/aeneasr))
- telemetry: Update telemetry identification [\#654](https://github.com/ory/hydra/pull/654) ([aeneasr](https://github.com/aeneasr))

## [v0.10.0-alpha.21](https://github.com/ory/hydra/tree/v0.10.0-alpha.21) (2017-11-27)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.0-alpha.20...v0.10.0-alpha.21)

**Closed issues:**

- Add support for CORS [\#506](https://github.com/ory/hydra/issues/506)

**Merged pull requests:**

- cli: Fix hydra cli adding policy subjects on subject remove [\#665](https://github.com/ory/hydra/pull/665) ([jamesnicolas](https://github.com/jamesnicolas))

## [v0.10.0-alpha.20](https://github.com/ory/hydra/tree/v0.10.0-alpha.20) (2017-11-26)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.0-alpha.19...v0.10.0-alpha.20)

**Merged pull requests:**

- cmd: Added cors support to host process [\#664](https://github.com/ory/hydra/pull/664) ([aeneasr](https://github.com/aeneasr))

## [v0.10.0-alpha.19](https://github.com/ory/hydra/tree/v0.10.0-alpha.19) (2017-11-26)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.0-alpha.18...v0.10.0-alpha.19)

**Closed issues:**

- Working with flask-oidc [\#660](https://github.com/ory/hydra/issues/660)
- Multi stage build process removes the ability to shell into hydra container [\#657](https://github.com/ory/hydra/issues/657)
- Support ES256 JWK Algo [\#627](https://github.com/ory/hydra/issues/627)
- oauth2/introspect: skip omitempty in active flag [\#607](https://github.com/ory/hydra/issues/607)
- oauth2: provide CWT token generation [\#577](https://github.com/ory/hydra/issues/577)

**Merged pull requests:**

- vendor: Upgraded ladon and dockertest versions [\#663](https://github.com/ory/hydra/pull/663) ([aeneasr](https://github.com/aeneasr))
- pkg: Make low entropy RSA key generation explicit in function name [\#656](https://github.com/ory/hydra/pull/656) ([aeneasr](https://github.com/aeneasr))
- docs: Update hydra versions [\#649](https://github.com/ory/hydra/pull/649) ([aeneasr](https://github.com/aeneasr))

## [v0.10.0-alpha.18](https://github.com/ory/hydra/tree/v0.10.0-alpha.18) (2017-11-06)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.0-alpha.17...v0.10.0-alpha.18)

## [v0.10.0-alpha.17](https://github.com/ory/hydra/tree/v0.10.0-alpha.17) (2017-11-06)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.0-alpha.16...v0.10.0-alpha.17)

## [v0.10.0-alpha.16](https://github.com/ory/hydra/tree/v0.10.0-alpha.16) (2017-11-06)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.0-alpha.15...v0.10.0-alpha.16)

**Merged pull requests:**

- Fix static build [\#648](https://github.com/ory/hydra/pull/648) ([aeneasr](https://github.com/aeneasr))

## [v0.10.0-alpha.15](https://github.com/ory/hydra/tree/v0.10.0-alpha.15) (2017-11-06)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.0-alpha.14...v0.10.0-alpha.15)

**Merged pull requests:**

- docker: Make hydra executable [\#647](https://github.com/ory/hydra/pull/647) ([aeneasr](https://github.com/aeneasr))

## [v0.10.0-alpha.14](https://github.com/ory/hydra/tree/v0.10.0-alpha.14) (2017-11-06)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.0-alpha.13...v0.10.0-alpha.14)

**Fixed bugs:**

- sql/postgres: wherever limit/offset is used, include ORDER BY clause [\#619](https://github.com/ory/hydra/issues/619)
- oauth2: fix racy memory consent manager with RW mutex [\#600](https://github.com/ory/hydra/issues/600)

**Merged pull requests:**

- Fix racy behaviour in oauth2 memory managers [\#646](https://github.com/ory/hydra/pull/646) ([aeneasr](https://github.com/aeneasr))

## [v0.10.0-alpha.13](https://github.com/ory/hydra/tree/v0.10.0-alpha.13) (2017-11-06)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.0-alpha.11...v0.10.0-alpha.13)

**Implemented enhancements:**

- Would it make sense to build hydra statically [\#374](https://github.com/ory/hydra/issues/374)

**Merged pull requests:**

- docker: Stop building from source in docker image [\#645](https://github.com/ory/hydra/pull/645) ([aeneasr](https://github.com/aeneasr))

## [v0.10.0-alpha.11](https://github.com/ory/hydra/tree/v0.10.0-alpha.11) (2017-11-06)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.0-alpha.12...v0.10.0-alpha.11)

## [v0.10.0-alpha.12](https://github.com/ory/hydra/tree/v0.10.0-alpha.12) (2017-11-06)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.0-alpha.10...v0.10.0-alpha.12)

**Closed issues:**

- Add license header to all source files [\#643](https://github.com/ory/hydra/issues/643)
- warden: remove obsolete http manager [\#616](https://github.com/ory/hydra/issues/616)

**Merged pull requests:**

- Add license header to all source files [\#644](https://github.com/ory/hydra/pull/644) ([aeneasr](https://github.com/aeneasr))
- cmd: require url-encoding of root client id and secret [\#641](https://github.com/ory/hydra/pull/641) ([aeneasr](https://github.com/aeneasr))
- fix health link in docs [\#637](https://github.com/ory/hydra/pull/637) ([DallanQ](https://github.com/DallanQ))

## [v0.10.0-alpha.10](https://github.com/ory/hydra/tree/v0.10.0-alpha.10) (2017-10-26)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.0-alpha.9...v0.10.0-alpha.10)

**Implemented enhancements:**

- jwk: use cryptopasta library [\#629](https://github.com/ory/hydra/issues/629)
- Feature Request: ability to list all groups [\#594](https://github.com/ory/hydra/issues/594)

**Closed issues:**

- jwk: add es256 generator to jwk handler in master [\#634](https://github.com/ory/hydra/issues/634)
- groups: add ability to list all groups to master branch [\#633](https://github.com/ory/hydra/issues/633)
- travis: run genswag and gensdk before npm publish [\#610](https://github.com/ory/hydra/issues/610)

## [v0.10.0-alpha.9](https://github.com/ory/hydra/tree/v0.10.0-alpha.9) (2017-10-25)
[Full Changelog](https://github.com/ory/hydra/compare/v0.9.16...v0.10.0-alpha.9)

**Closed issues:**

- docs: followed the installation guide and was unable to get a successful consent [\#623](https://github.com/ory/hydra/issues/623)
- tests: run manager tests in parallel [\#617](https://github.com/ory/hydra/issues/617)

**Merged pull requests:**

- Changes from zvelo [\#636](https://github.com/ory/hydra/pull/636) ([aeneasr](https://github.com/aeneasr))
- Dep, JWK and groups [\#635](https://github.com/ory/hydra/pull/635) ([aeneasr](https://github.com/aeneasr))
- tests: run database tests in parallel [\#632](https://github.com/ory/hydra/pull/632) ([aeneasr](https://github.com/aeneasr))
- Use recommendations made from cryptopasta repository [\#630](https://github.com/ory/hydra/pull/630) ([aeneasr](https://github.com/aeneasr))
- Support ES256 JWK Algo [\#628](https://github.com/ory/hydra/pull/628) ([joshuarubin](https://github.com/joshuarubin))

## [v0.9.16](https://github.com/ory/hydra/tree/v0.9.16) (2017-10-23)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.0-alpha.8...v0.9.16)

**Closed issues:**

- docs: adding policy to consent app doesn't work as resource using \<.\*\> [\#621](https://github.com/ory/hydra/issues/621)
- documentation vague regarding returned client\_secret [\#620](https://github.com/ory/hydra/issues/620)

**Merged pull requests:**

- updated links to apiary as the old ones didn't work [\#626](https://github.com/ory/hydra/pull/626) ([abusaidm](https://github.com/abusaidm))
- docs: updated hydra version in the tutorial to v0.10.0-alpha.8 and consent app to v0.10.0-alpha.9 [\#625](https://github.com/ory/hydra/pull/625) ([abusaidm](https://github.com/abusaidm))
- docs: fixed spelling and wording [\#624](https://github.com/ory/hydra/pull/624) ([abusaidm](https://github.com/abusaidm))
- docs: fix bash command and version used in tutorial [\#622](https://github.com/ory/hydra/pull/622) ([abusaidm](https://github.com/abusaidm))
- add ability to list all groups [\#612](https://github.com/ory/hydra/pull/612) ([joshuarubin](https://github.com/joshuarubin))

## [v0.10.0-alpha.8](https://github.com/ory/hydra/tree/v0.10.0-alpha.8) (2017-10-18)
[Full Changelog](https://github.com/ory/hydra/compare/v0.9.15...v0.10.0-alpha.8)

**Closed issues:**

- docs: SDK for Go is actually for Node, fix this typo [\#615](https://github.com/ory/hydra/issues/615)
- server.injectConsentManager doesn't use ConsentRequestSQLManager even if \*config.SQLConnection exists [\#613](https://github.com/ory/hydra/issues/613)

**Merged pull requests:**

- cmd/server: SQLConnection should load SQLRequestManager [\#618](https://github.com/ory/hydra/pull/618) ([aeneasr](https://github.com/aeneasr))
- Clean up helpers and increase test coverage [\#611](https://github.com/ory/hydra/pull/611) ([aeneasr](https://github.com/aeneasr))
- sdk: format js sdk and remove mock tests [\#609](https://github.com/ory/hydra/pull/609) ([aeneasr](https://github.com/aeneasr))

## [v0.9.15](https://github.com/ory/hydra/tree/v0.9.15) (2017-10-11)
[Full Changelog](https://github.com/ory/hydra/compare/v0.9.14...v0.9.15)

**Merged pull requests:**

- Support dep [\#606](https://github.com/ory/hydra/pull/606) ([joshuarubin](https://github.com/joshuarubin))

## [v0.9.14](https://github.com/ory/hydra/tree/v0.9.14) (2017-10-06)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.0-alpha.7...v0.9.14)

## [v0.10.0-alpha.7](https://github.com/ory/hydra/tree/v0.10.0-alpha.7) (2017-10-06)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.0-alpha.6...v0.10.0-alpha.7)

## [v0.10.0-alpha.6](https://github.com/ory/hydra/tree/v0.10.0-alpha.6) (2017-10-05)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.0-alpha.5...v0.10.0-alpha.6)

## [v0.10.0-alpha.5](https://github.com/ory/hydra/tree/v0.10.0-alpha.5) (2017-10-05)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.0-alpha.4...v0.10.0-alpha.5)

## [v0.10.0-alpha.4](https://github.com/ory/hydra/tree/v0.10.0-alpha.4) (2017-10-05)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.0-alpha.3...v0.10.0-alpha.4)

**Merged pull requests:**

- travis: move deploy scripts to its own file [\#604](https://github.com/ory/hydra/pull/604) ([aeneasr](https://github.com/aeneasr))
- tests: skip cpu intense jwk generation in short mode [\#603](https://github.com/ory/hydra/pull/603) ([aeneasr](https://github.com/aeneasr))

## [v0.10.0-alpha.3](https://github.com/ory/hydra/tree/v0.10.0-alpha.3) (2017-10-05)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.0-alpha.2...v0.10.0-alpha.3)

## [v0.10.0-alpha.2](https://github.com/ory/hydra/tree/v0.10.0-alpha.2) (2017-10-05)
[Full Changelog](https://github.com/ory/hydra/compare/v0.10.0-alpha.1...v0.10.0-alpha.2)

**Implemented enhancements:**

- all: refactor http client endpoint logic [\#584](https://github.com/ory/hydra/issues/584)
- oauth2: refresh openid connect id token via refresh\_token grant [\#556](https://github.com/ory/hydra/issues/556)
- oauth2: change scope semantics to wildcard [\#550](https://github.com/ory/hydra/issues/550)
- warden: need endpoint that just introspects tokens [\#539](https://github.com/ory/hydra/issues/539)
- sdk: client libraries for all languages [\#249](https://github.com/ory/hydra/issues/249)
- core: enable usage statistics reporting [\#230](https://github.com/ory/hydra/issues/230)
- core: introduce a way to test for bc breaks in datastore [\#193](https://github.com/ory/hydra/issues/193)

**Merged pull requests:**

- travis: resolve deployment issues [\#602](https://github.com/ory/hydra/pull/602) ([aeneasr](https://github.com/aeneasr))
- warden: remove deprecated http manager [\#601](https://github.com/ory/hydra/pull/601) ([aeneasr](https://github.com/aeneasr))
- docs: fix sdk links [\#599](https://github.com/ory/hydra/pull/599) ([aeneasr](https://github.com/aeneasr))
- travis: re-add goveralls [\#598](https://github.com/ory/hydra/pull/598) ([aeneasr](https://github.com/aeneasr))

## [v0.10.0-alpha.1](https://github.com/ory/hydra/tree/v0.10.0-alpha.1) (2017-10-05)
[Full Changelog](https://github.com/ory/hydra/compare/v0.9.13...v0.10.0-alpha.1)

**Implemented enhancements:**

- oauth2: write test for handling consent deny [\#597](https://github.com/ory/hydra/issues/597)
- group: add warden tests [\#591](https://github.com/ory/hydra/issues/591)
- health: remove TLS restriction on health endpoint when termination is set [\#586](https://github.com/ory/hydra/issues/586)

**Fixed bugs:**

- cmd: `policies delete` says `Connection \<id\> deleted` instead of `Policy \<id\> deleted` [\#583](https://github.com/ory/hydra/issues/583)

**Closed issues:**

- oauth2: change meaning of audience claim [\#595](https://github.com/ory/hydra/issues/595)
- sdk/go: write interfaces for APIs & responses [\#593](https://github.com/ory/hydra/issues/593)

**Merged pull requests:**

- travis: fix binary building [\#596](https://github.com/ory/hydra/pull/596) ([aeneasr](https://github.com/aeneasr))
- cmd/cli: typo Connection -\> Policy [\#592](https://github.com/ory/hydra/pull/592) ([ljagiello](https://github.com/ljagiello))
- sdk: switch to swagger codegen sdk [\#585](https://github.com/ory/hydra/pull/585) ([aeneasr](https://github.com/aeneasr))
- 0.10.0 [\#557](https://github.com/ory/hydra/pull/557) ([aeneasr](https://github.com/aeneasr))

## [v0.9.13](https://github.com/ory/hydra/tree/v0.9.13) (2017-09-26)
[Full Changelog](https://github.com/ory/hydra/compare/v0.9.12...v0.9.13)

**Implemented enhancements:**

- RFC: Refactor consent flow [\#578](https://github.com/ory/hydra/issues/578)
- oauth2: remove scope parameter from introspection request [\#551](https://github.com/ory/hydra/issues/551)
- "Subject claim can not be empty" error when trying to retrieve ID Token [\#460](https://github.com/ory/hydra/issues/460)

**Fixed bugs:**

- cmd: `token user` no longer uses cluster url [\#581](https://github.com/ory/hydra/issues/581)
- warden: do not use refresh tokens as proof of authorization [\#549](https://github.com/ory/hydra/issues/549)
- Fix import path for logrus [\#477](https://github.com/ory/hydra/issues/477)

**Closed issues:**

- Support for RFC 7636 [\#576](https://github.com/ory/hydra/issues/576)
- `authorization` header in `/oauth2/token` endpoint is case sensitive [\#575](https://github.com/ory/hydra/issues/575)
- DATABASE\_URL=memory go run main.go host Error [\#571](https://github.com/ory/hydra/issues/571)
- error on mismatch uris [\#569](https://github.com/ory/hydra/issues/569)
- Relation "hydra\_jwk" does not exist [\#568](https://github.com/ory/hydra/issues/568)
- Freemium Crap [\#567](https://github.com/ory/hydra/issues/567)
- Warden API docs do not talk about access\_token [\#564](https://github.com/ory/hydra/issues/564)
- When the client is run through a container, it should pick up configuration from environment [\#563](https://github.com/ory/hydra/issues/563)
- Docker hub documentation showing up as HTML [\#562](https://github.com/ory/hydra/issues/562)
- Allow people to configure the Hydra service using a config file. [\#561](https://github.com/ory/hydra/issues/561)
- Error on go get the project [\#560](https://github.com/ory/hydra/issues/560)
- Open a Patreon account [\#558](https://github.com/ory/hydra/issues/558)
- GET /client/:id broken on master [\#538](https://github.com/ory/hydra/issues/538)

**Merged pull requests:**

- health: disable TLS restriction for health check [\#587](https://github.com/ory/hydra/pull/587) ([aeneasr](https://github.com/aeneasr))
- cmd: `token user` should use clusterurl instead of empty string [\#582](https://github.com/ory/hydra/pull/582) ([aeneasr](https://github.com/aeneasr))
- vendor: update various dependencies [\#579](https://github.com/ory/hydra/pull/579) ([aeneasr](https://github.com/aeneasr))
- Update to ladon 0.8.2 [\#570](https://github.com/ory/hydra/pull/570) ([olivierdeckers](https://github.com/olivierdeckers))
- install.md: port typo [\#566](https://github.com/ory/hydra/pull/566) ([rnback](https://github.com/rnback))
- oauth2: give meaningful hint when subject claim is empty [\#554](https://github.com/ory/hydra/pull/554) ([aeneasr](https://github.com/aeneasr))

## [v0.9.12](https://github.com/ory/hydra/tree/v0.9.12) (2017-07-06)
[Full Changelog](https://github.com/ory/hydra/compare/v0.9.11...v0.9.12)

**Implemented enhancements:**

- oauth2: use wildcards for scope strategy [\#552](https://github.com/ory/hydra/issues/552)

**Merged pull requests:**

- warden: refresh tokens are no longer proof of authZ [\#553](https://github.com/ory/hydra/pull/553) ([aeneasr](https://github.com/aeneasr))
- README.md: hydra container doesn't include bash [\#548](https://github.com/ory/hydra/pull/548) ([srenatus](https://github.com/srenatus))
- docs: fix typo in tutorial [\#547](https://github.com/ory/hydra/pull/547) ([aeneasr](https://github.com/aeneasr))
- cmd/token/user: fix auth and token-url mixup [\#546](https://github.com/ory/hydra/pull/546) ([aeneasr](https://github.com/aeneasr))
- docs: update docs [\#545](https://github.com/ory/hydra/pull/545) ([aeneasr](https://github.com/aeneasr))

## [v0.9.11](https://github.com/ory/hydra/tree/v0.9.11) (2017-06-30)
[Full Changelog](https://github.com/ory/hydra/compare/v0.9.10...v0.9.11)

**Merged pull requests:**

- docs: add step-by-step installation guide [\#544](https://github.com/ory/hydra/pull/544) ([aeneasr](https://github.com/aeneasr))
- docs: add product teasers [\#543](https://github.com/ory/hydra/pull/543) ([aeneasr](https://github.com/aeneasr))

## [v0.9.10](https://github.com/ory/hydra/tree/v0.9.10) (2017-06-29)
[Full Changelog](https://github.com/ory/hydra/compare/v0.9.9...v0.9.10)

**Implemented enhancements:**

- cmd/host: move status info from health endpoint to another one and protect it [\#532](https://github.com/ory/hydra/issues/532)

**Fixed bugs:**

- Decode Basic Auth Credentials [\#536](https://github.com/ory/hydra/issues/536)

**Closed issues:**

- Cannot try tutorial install, not existing dependencies [\#541](https://github.com/ory/hydra/issues/541)
- \[docker-compose\] ERROR: for postgresd expected string or buffer [\#540](https://github.com/ory/hydra/issues/540)

**Merged pull requests:**

- vendor: update fosite to remove forced nonce [\#542](https://github.com/ory/hydra/pull/542) ([aeneasr](https://github.com/aeneasr))
- oauth2: form-urldecode authorization basic header [\#537](https://github.com/ory/hydra/pull/537) ([aeneasr](https://github.com/aeneasr))
- \[DOC\] Update "Build from source" section to actual state [\#534](https://github.com/ory/hydra/pull/534) ([dolbik](https://github.com/dolbik))
- cmd/host: move status info to dedicated endpoint [\#533](https://github.com/ory/hydra/pull/533) ([aeneasr](https://github.com/aeneasr))

## [v0.9.9](https://github.com/ory/hydra/tree/v0.9.9) (2017-06-17)
[Full Changelog](https://github.com/ory/hydra/compare/v0.9.8...v0.9.9)

**Fixed bugs:**

- cmd/policy/create: not exiting on error [\#527](https://github.com/ory/hydra/issues/527)

**Merged pull requests:**

- cmd: add test for get handler [\#531](https://github.com/ory/hydra/pull/531) ([aeneasr](https://github.com/aeneasr))
- cmd/policy/create: exit on error - closes \#527 [\#530](https://github.com/ory/hydra/pull/530) ([aeneasr](https://github.com/aeneasr))

## [v0.9.8](https://github.com/ory/hydra/tree/v0.9.8) (2017-06-17)
[Full Changelog](https://github.com/ory/hydra/compare/v0.9.7...v0.9.8)

**Fixed bugs:**

- Updating policies may cause loss of policy data [\#503](https://github.com/ory/hydra/issues/503)

**Closed issues:**

- oauth2: investigate panic [\#512](https://github.com/ory/hydra/issues/512)

**Merged pull requests:**

- oauth2: resolve panic with nested at\_ext and id\_ext [\#529](https://github.com/ory/hydra/pull/529) ([aeneasr](https://github.com/aeneasr))
- vendor: update to ladon 0.8.0 - closes \#503 [\#528](https://github.com/ory/hydra/pull/528) ([aeneasr](https://github.com/aeneasr))

## [v0.9.7](https://github.com/ory/hydra/tree/v0.9.7) (2017-06-16)
[Full Changelog](https://github.com/ory/hydra/compare/v0.9.6...v0.9.7)

**Closed issues:**

- Fatal error when running docker container [\#525](https://github.com/ory/hydra/issues/525)

**Merged pull requests:**

- cmd/server: supply admin client policy with id [\#526](https://github.com/ory/hydra/pull/526) ([aeneasr](https://github.com/aeneasr))

## [v0.9.6](https://github.com/ory/hydra/tree/v0.9.6) (2017-06-15)
[Full Changelog](https://github.com/ory/hydra/compare/v0.9.5...v0.9.6)

**Merged pull requests:**

- Db plugin connector [\#524](https://github.com/ory/hydra/pull/524) ([aeneasr](https://github.com/aeneasr))

## [v0.9.5](https://github.com/ory/hydra/tree/v0.9.5) (2017-06-15)
[Full Changelog](https://github.com/ory/hydra/compare/v0.9.4...v0.9.5)

**Merged pull requests:**

- vendor: upgrade ladon to 0.7.7 [\#523](https://github.com/ory/hydra/pull/523) ([aeneasr](https://github.com/aeneasr))

## [v0.9.4](https://github.com/ory/hydra/tree/v0.9.4) (2017-06-14)
[Full Changelog](https://github.com/ory/hydra/compare/v0.9.3...v0.9.4)

**Merged pull requests:**

- cmd: resolve issuer test issue [\#522](https://github.com/ory/hydra/pull/522) ([aeneasr](https://github.com/aeneasr))
- all: improve test exports [\#521](https://github.com/ory/hydra/pull/521) ([aeneasr](https://github.com/aeneasr))
- docs: start writing faq from gitter [\#504](https://github.com/ory/hydra/pull/504) ([aeneasr](https://github.com/aeneasr))

## [v0.9.3](https://github.com/ory/hydra/tree/v0.9.3) (2017-06-14)
[Full Changelog](https://github.com/ory/hydra/compare/v0.9.2...v0.9.3)

**Closed issues:**

- Generating Client ID/Secret in \>= 0.8.0 [\#517](https://github.com/ory/hydra/issues/517)
- Could not gracefully run server [\#513](https://github.com/ory/hydra/issues/513)
- authorize\_code without password [\#511](https://github.com/ory/hydra/issues/511)

**Merged pull requests:**

- metrics: resolve potential data race [\#520](https://github.com/ory/hydra/pull/520) ([aeneasr](https://github.com/aeneasr))
- Fix warden docs [\#519](https://github.com/ory/hydra/pull/519) ([aeneasr](https://github.com/aeneasr))
- all: export test helpers [\#518](https://github.com/ory/hydra/pull/518) ([aeneasr](https://github.com/aeneasr))
- oauth2: add tests for refresh token grant [\#515](https://github.com/ory/hydra/pull/515) ([aeneasr](https://github.com/aeneasr))
- oauth2: use issuer-prefixed auth URL in challenge redirect [\#509](https://github.com/ory/hydra/pull/509) ([wyattanderson](https://github.com/wyattanderson))
- cmd: resolve failing test [\#501](https://github.com/ory/hydra/pull/501) ([aeneasr](https://github.com/aeneasr))

## [v0.9.2](https://github.com/ory/hydra/tree/v0.9.2) (2017-06-13)
[Full Changelog](https://github.com/ory/hydra/compare/v0.9.1...v0.9.2)

**Merged pull requests:**

- cmd/server: print full error message on http startup [\#514](https://github.com/ory/hydra/pull/514) ([aeneasr](https://github.com/aeneasr))

## [v0.9.1](https://github.com/ory/hydra/tree/v0.9.1) (2017-06-12)
[Full Changelog](https://github.com/ory/hydra/compare/v0.9.0...v0.9.1)

**Merged pull requests:**

- client: export tests [\#510](https://github.com/ory/hydra/pull/510) ([aeneasr](https://github.com/aeneasr))
- metrics: improve metrics [\#508](https://github.com/ory/hydra/pull/508) ([aeneasr](https://github.com/aeneasr))
- cmd: add auto migration image [\#502](https://github.com/ory/hydra/pull/502) ([aeneasr](https://github.com/aeneasr))

## [v0.9.0](https://github.com/ory/hydra/tree/v0.9.0) (2017-06-07)
[Full Changelog](https://github.com/ory/hydra/compare/v0.8.7...v0.9.0)

**Implemented enhancements:**

- cmd/cli: add flag for X-Forwarded-Proto for faking https termination [\#349](https://github.com/ory/hydra/issues/349)
- metrics: add metrics and telemetry package [\#500](https://github.com/ory/hydra/pull/500) ([aeneasr](https://github.com/aeneasr))

**Fixed bugs:**

- warden/group: investigate missing transaction rollback in group manager [\#462](https://github.com/ory/hydra/issues/462)
- policies: validate conditions and return error instead of silently dropping them [\#350](https://github.com/ory/hydra/issues/350)

**Closed issues:**

- Headers should be case-insensitive [\#496](https://github.com/ory/hydra/issues/496)
- docs: add FAQ on missing migrate in docker image [\#484](https://github.com/ory/hydra/issues/484)
- docs: include oauth2 example [\#358](https://github.com/ory/hydra/issues/358)
- warden: allow scopes in policies [\#330](https://github.com/ory/hydra/issues/330)

**Merged pull requests:**

- sdk: add simple example of hydra sdk [\#499](https://github.com/ory/hydra/pull/499) ([aeneasr](https://github.com/aeneasr))
- docs: add FAQ on missing migrate in docker image [\#498](https://github.com/ory/hydra/pull/498) ([aeneasr](https://github.com/aeneasr))
- vendor: upgrade to ladon 0.7.4 - closes \#350 [\#497](https://github.com/ory/hydra/pull/497) ([aeneasr](https://github.com/aeneasr))
- docs: add scopes to oauth2 [\#495](https://github.com/ory/hydra/pull/495) ([aeneasr](https://github.com/aeneasr))
- warden/group: add rollback to transactions [\#494](https://github.com/ory/hydra/pull/494) ([aeneasr](https://github.com/aeneasr))

## [v0.8.7](https://github.com/ory/hydra/tree/v0.8.7) (2017-06-05)
[Full Changelog](https://github.com/ory/hydra/compare/v0.8.6...v0.8.7)

**Implemented enhancements:**

- oauth2: add possibility for denying consent requests [\#400](https://github.com/ory/hydra/issues/400)
- oauth2: allow redirection to client if consent was denied [\#371](https://github.com/ory/hydra/issues/371)

**Fixed bugs:**

- Introspection endpoint responds with 401 on invalid payload token [\#457](https://github.com/ory/hydra/issues/457)

**Closed issues:**

- Allow configuration of `DB\_HOST`, `DB\_PASS`, `DB\_USER`, `DB\_NAME` separately. [\#480](https://github.com/ory/hydra/issues/480)

**Merged pull requests:**

- all: implement --fake-tls-termination flag [\#493](https://github.com/ory/hydra/pull/493) ([aeneasr](https://github.com/aeneasr))
- oauth2/introspect\>: resolve 401 on invalid token [\#492](https://github.com/ory/hydra/pull/492) ([aeneasr](https://github.com/aeneasr))
- client/manager\_sql: return an empty slice if string is empty [\#491](https://github.com/ory/hydra/pull/491) ([faxal](https://github.com/faxal))

## [v0.8.6](https://github.com/ory/hydra/tree/v0.8.6) (2017-06-05)
[Full Changelog](https://github.com/ory/hydra/compare/v0.8.5...v0.8.6)

**Implemented enhancements:**

- Assign clients different consent urls  [\#378](https://github.com/ory/hydra/issues/378)

**Fixed bugs:**

- Creating policies via the CLI does not populate the 'description' field [\#472](https://github.com/ory/hydra/issues/472)
- Missing "iss" field from /oauth2/introspect response [\#399](https://github.com/ory/hydra/issues/399)
- client: getting a non-existing client raises 500 instead of 404 [\#348](https://github.com/ory/hydra/issues/348)

**Closed issues:**

- Libraries version problem, build break. [\#481](https://github.com/ory/hydra/issues/481)
- oauth2: update to latest fosite which removed implicit storage [\#468](https://github.com/ory/hydra/issues/468)
- Unable to set Public flag to false [\#463](https://github.com/ory/hydra/issues/463)
- oauth2: allow client specific token TTLs [\#428](https://github.com/ory/hydra/issues/428)
- docs: hint at health check [\#355](https://github.com/ory/hydra/issues/355)
- Hydra URLs mounted to a subpath [\#352](https://github.com/ory/hydra/issues/352)
- oidc: hydra as federated user auth for AWS Console/API [\#315](https://github.com/ory/hydra/issues/315)
- jwk: when retrieving a key, stray request missing a subject 403 [\#271](https://github.com/ory/hydra/issues/271)

**Merged pull requests:**

- oauth2/introspect: send issuer in introspection [\#490](https://github.com/ory/hydra/pull/490) ([aeneasr](https://github.com/aeneasr))
- oauth2: allow redirection to client if consent was denied [\#489](https://github.com/ory/hydra/pull/489) ([aeneasr](https://github.com/aeneasr))
- docs: add health check to swagger and resolve swagger issues [\#488](https://github.com/ory/hydra/pull/488) ([aeneasr](https://github.com/aeneasr))
- jwk/handler: nest ac check and resolve stray log message [\#487](https://github.com/ory/hydra/pull/487) ([aeneasr](https://github.com/aeneasr))
- pkg/errors: make ErrNotFound return a status code [\#486](https://github.com/ory/hydra/pull/486) ([aeneasr](https://github.com/aeneasr))
- cmd/policies: description is a string field, not slice [\#485](https://github.com/ory/hydra/pull/485) ([aeneasr](https://github.com/aeneasr))
- Vendor update [\#483](https://github.com/ory/hydra/pull/483) ([aeneasr](https://github.com/aeneasr))
- vendor: update to latest versions [\#482](https://github.com/ory/hydra/pull/482) ([aeneasr](https://github.com/aeneasr))
- client/manager: remove merging of stored and updated client [\#478](https://github.com/ory/hydra/pull/478) ([faxal](https://github.com/faxal))
- Fix Swagger for Warden Groups [\#476](https://github.com/ory/hydra/pull/476) ([pbarker](https://github.com/pbarker))

## [v0.8.5](https://github.com/ory/hydra/tree/v0.8.5) (2017-06-01)
[Full Changelog](https://github.com/ory/hydra/compare/v0.8.4...v0.8.5)

**Fixed bugs:**

- max\_conns and max\_conn\_lifetime breaks db.Ping [\#464](https://github.com/ory/hydra/issues/464)
- cmd/server: resolve gorilla session mem leak - closes \#461 [\#475](https://github.com/ory/hydra/pull/475) ([aeneasr](https://github.com/aeneasr))

**Closed issues:**

- Container is not Running [\#474](https://github.com/ory/hydra/issues/474)
- Random periodic crashes  [\#461](https://github.com/ory/hydra/issues/461)

**Merged pull requests:**

- fix spelling of challenge [\#471](https://github.com/ory/hydra/pull/471) ([sstarcher](https://github.com/sstarcher))
- oauth2: remove unused implicit grant storage [\#469](https://github.com/ory/hydra/pull/469) ([aeneasr](https://github.com/aeneasr))

## [v0.8.4](https://github.com/ory/hydra/tree/v0.8.4) (2017-05-24)
[Full Changelog](https://github.com/ory/hydra/compare/v0.8.3...v0.8.4)

**Closed issues:**

- Kubernetes Helm chart [\#430](https://github.com/ory/hydra/issues/430)

**Merged pull requests:**

- config: connect to cleaned DSN [\#470](https://github.com/ory/hydra/pull/470) ([aeneasr](https://github.com/aeneasr))
- docs: hint to kubernetes helm chart - see \#430 [\#467](https://github.com/ory/hydra/pull/467) ([aeneasr](https://github.com/aeneasr))
- Improve documentation [\#466](https://github.com/ory/hydra/pull/466) ([aeneasr](https://github.com/aeneasr))

## [v0.8.3](https://github.com/ory/hydra/tree/v0.8.3) (2017-05-23)
[Full Changelog](https://github.com/ory/hydra/compare/v0.8.2...v0.8.3)

**Implemented enhancements:**

- http: harden http server for public net [\#334](https://github.com/ory/hydra/issues/334)

**Fixed bugs:**

- config: remove sql control parameters from dsn before connecting [\#465](https://github.com/ory/hydra/pull/465) ([aeneasr](https://github.com/aeneasr))

**Closed issues:**

- Listing policies not working with database [\#458](https://github.com/ory/hydra/issues/458)
- go install github.com/ory/hydra Fails to compile [\#456](https://github.com/ory/hydra/issues/456)
- Challenge claims redirect http instead of https [\#455](https://github.com/ory/hydra/issues/455)
- core/store: document aes gcm nonce limitation [\#76](https://github.com/ory/hydra/issues/76)

**Merged pull requests:**

- Policy Fix [\#459](https://github.com/ory/hydra/pull/459) ([pbarker](https://github.com/pbarker))

## [v0.8.2](https://github.com/ory/hydra/tree/v0.8.2) (2017-05-10)
[Full Changelog](https://github.com/ory/hydra/compare/v0.8.1...v0.8.2)

**Implemented enhancements:**

- Missing `kid` parameter in ID token header [\#433](https://github.com/ory/hydra/issues/433)
- no /.well-known/openid-configuration endpoint implementation [\#379](https://github.com/ory/hydra/issues/379)

**Merged pull requests:**

- Add Key Id to Header [\#454](https://github.com/ory/hydra/pull/454) ([pbarker](https://github.com/pbarker))
- cmd: improve error message for when database tables are missing [\#453](https://github.com/ory/hydra/pull/453) ([aeneasr](https://github.com/aeneasr))
- Wellknown [\#427](https://github.com/ory/hydra/pull/427) ([pbarker](https://github.com/pbarker))

## [v0.8.1](https://github.com/ory/hydra/tree/v0.8.1) (2017-05-08)
[Full Changelog](https://github.com/ory/hydra/compare/v0.8.0...v0.8.1)

**Implemented enhancements:**

- cmd: database migrations should not be run automatically but have a cmd instead [\#444](https://github.com/ory/hydra/issues/444)
- all: move herodot to ory/herodot [\#436](https://github.com/ory/hydra/issues/436)

**Fixed bugs:**

- cmd: token client fails in ci sometimes [\#443](https://github.com/ory/hydra/issues/443)

**Closed issues:**

- all: deprecating rethinkdb and redis support [\#425](https://github.com/ory/hydra/issues/425)
- oauth2: consent anti-csrf token should be forcefully removed [\#367](https://github.com/ory/hydra/issues/367)

## [v0.8.0](https://github.com/ory/hydra/tree/v0.8.0) (2017-05-07)
[Full Changelog](https://github.com/ory/hydra/compare/v0.7.13...v0.8.0)

**Closed issues:**

- Refresh token doesn't work [\#449](https://github.com/ory/hydra/issues/449)

**Merged pull requests:**

- ✏️  minor grammar typo [\#452](https://github.com/ory/hydra/pull/452) ([therebelrobot](https://github.com/therebelrobot))
- Add example about securing the consent app [\#450](https://github.com/ory/hydra/pull/450) ([matteosuppo](https://github.com/matteosuppo))
- Allow setting SkipTLSVerify Option value [\#448](https://github.com/ory/hydra/pull/448) ([faxal](https://github.com/faxal))
- 0.8.0: Towards production friendliness [\#445](https://github.com/ory/hydra/pull/445) ([aeneasr](https://github.com/aeneasr))

## [v0.7.13](https://github.com/ory/hydra/tree/v0.7.13) (2017-05-03)
[Full Changelog](https://github.com/ory/hydra/compare/v0.7.12...v0.7.13)

**Implemented enhancements:**

- ui: implement a basic management interface with react for oauth2 client, jwk, social connections and others [\#215](https://github.com/ory/hydra/issues/215)

**Fixed bugs:**

- herodot: resolve issue with infinite loop caused by certain error chain [\#441](https://github.com/ory/hydra/issues/441)
- "Could not fetch signing key for OpenID Connect" [\#439](https://github.com/ory/hydra/issues/439)
- vendor: upgrade fosite to resolve regression issue [\#446](https://github.com/ory/hydra/pull/446) ([aeneasr](https://github.com/aeneasr))

**Closed issues:**

- Peculiar EOF instead of response from the introspect endpoint. [\#368](https://github.com/ory/hydra/issues/368)

**Merged pull requests:**

- Add Auth0 to sponsor section [\#435](https://github.com/ory/hydra/pull/435) ([aeneasr](https://github.com/aeneasr))

## [v0.7.12](https://github.com/ory/hydra/tree/v0.7.12) (2017-04-30)
[Full Changelog](https://github.com/ory/hydra/compare/v0.7.11...v0.7.12)

**Fixed bugs:**

- herodot: resolve issue with infinite loop caused by certain error chain [\#442](https://github.com/ory/hydra/pull/442) ([aeneasr](https://github.com/aeneasr))

**Closed issues:**

- Freeze dependencies  [\#437](https://github.com/ory/hydra/issues/437)

## [v0.7.11](https://github.com/ory/hydra/tree/v0.7.11) (2017-04-28)
[Full Changelog](https://github.com/ory/hydra/compare/v0.7.10...v0.7.11)

**Closed issues:**

- Mismatch between library versions [\#434](https://github.com/ory/hydra/issues/434)
- Data Passthrough to IDP [\#431](https://github.com/ory/hydra/issues/431)
- Api protection [\#429](https://github.com/ory/hydra/issues/429)
- Gitter.im or irc channel [\#426](https://github.com/ory/hydra/issues/426)
- Outdated fosite [\#424](https://github.com/ory/hydra/issues/424)
- oauth2: resource owner password credentials proposal [\#214](https://github.com/ory/hydra/issues/214)

**Merged pull requests:**

- vendor: resolve issues with glide lock file [\#438](https://github.com/ory/hydra/pull/438) ([aeneasr](https://github.com/aeneasr))

## [v0.7.10](https://github.com/ory/hydra/tree/v0.7.10) (2017-04-14)
[Full Changelog](https://github.com/ory/hydra/compare/v0.7.9...v0.7.10)

**Closed issues:**

- Build instructions from Readme fail [\#420](https://github.com/ory/hydra/issues/420)
- API error \(500\) during tests [\#419](https://github.com/ory/hydra/issues/419)
- Uname in session [\#418](https://github.com/ory/hydra/issues/418)
- Resource owner password credentials grant [\#417](https://github.com/ory/hydra/issues/417)
- ory vs ory-am [\#414](https://github.com/ory/hydra/issues/414)
- Cockroachdb support [\#413](https://github.com/ory/hydra/issues/413)
- Small doc error [\#411](https://github.com/ory/hydra/issues/411)
- Rest API documentation not working [\#410](https://github.com/ory/hydra/issues/410)

**Merged pull requests:**

- Remove uname references from docs [\#423](https://github.com/ory/hydra/pull/423) ([matteosuppo](https://github.com/matteosuppo))
- vendor: update common and ladon dependencies [\#422](https://github.com/ory/hydra/pull/422) ([aeneasr](https://github.com/aeneasr))
- docs: resolve broken build instructions in readme - closes \#420 [\#421](https://github.com/ory/hydra/pull/421) ([aeneasr](https://github.com/aeneasr))
- Dropping brackets in Create Client example [\#415](https://github.com/ory/hydra/pull/415) ([pbarker](https://github.com/pbarker))
- Update bash command in tutorial [\#412](https://github.com/ory/hydra/pull/412) ([pbarker](https://github.com/pbarker))
- Update README.md [\#409](https://github.com/ory/hydra/pull/409) ([joelpickup](https://github.com/joelpickup))
- docs: changes apiary url to current version [\#407](https://github.com/ory/hydra/pull/407) ([aeneasr](https://github.com/aeneasr))

## [v0.7.9](https://github.com/ory/hydra/tree/v0.7.9) (2017-04-02)
[Full Changelog](https://github.com/ory/hydra/compare/v0.7.8...v0.7.9)

**Closed issues:**

- Flow Using Curl help \(token auth\) [\#405](https://github.com/ory/hydra/issues/405)
- Add support to mongodb [\#401](https://github.com/ory/hydra/issues/401)

**Merged pull requests:**

- Updated ladon version in glide.lock [\#404](https://github.com/ory/hydra/pull/404) ([ericalandouglas](https://github.com/ericalandouglas))
- oauth2: fix typo [\#403](https://github.com/ory/hydra/pull/403) ([maximesong](https://github.com/maximesong))

## [v0.7.8](https://github.com/ory/hydra/tree/v0.7.8) (2017-03-24)
[Full Changelog](https://github.com/ory/hydra/compare/v0.7.7...v0.7.8)

**Implemented enhancements:**

- sdk: add consent helper [\#397](https://github.com/ory/hydra/issues/397)
- Transition Dockerfile to Alpine Linux [\#393](https://github.com/ory/hydra/issues/393)
- redirect\_uri domains are case-sensitive [\#380](https://github.com/ory/hydra/issues/380)
- Per-client consent URLs [\#351](https://github.com/ory/hydra/issues/351)
- sdk: add consent helper - closes \#397 [\#398](https://github.com/ory/hydra/pull/398) ([aeneasr](https://github.com/aeneasr))
- docs: add example policy for consent app signing [\#389](https://github.com/ory/hydra/pull/389) ([aeneasr](https://github.com/aeneasr))

**Fixed bugs:**

- cli handler\_groups type error? [\#383](https://github.com/ory/hydra/issues/383)

**Closed issues:**

- oauth2: token introspection fails on HTTP without dangerous-force-http [\#395](https://github.com/ory/hydra/issues/395)
- Create User based on access token provided by Social Provider [\#394](https://github.com/ory/hydra/issues/394)
- investigate why import from json fails [\#390](https://github.com/ory/hydra/issues/390)
- gitter link doesn't work [\#386](https://github.com/ory/hydra/issues/386)
- Possible security bug in warden/group package [\#382](https://github.com/ory/hydra/issues/382)
- relation "hydra\_client" does not exist \(postgres\) [\#381](https://github.com/ory/hydra/issues/381)
- Native login support [\#375](https://github.com/ory/hydra/issues/375)
- Request denied by default [\#373](https://github.com/ory/hydra/issues/373)

**Merged pull requests:**

- docker: reduce docker image size [\#396](https://github.com/ory/hydra/pull/396) ([aeneasr](https://github.com/aeneasr))
- Added information about auth code exchange to oauth2 docs [\#392](https://github.com/ory/hydra/pull/392) ([therebelrobot](https://github.com/therebelrobot))
- Small typo. [\#391](https://github.com/ory/hydra/pull/391) ([darron](https://github.com/darron))
- all: resolve ci issues and improve readme [\#384](https://github.com/ory/hydra/pull/384) ([aeneasr](https://github.com/aeneasr))

## [v0.7.7](https://github.com/ory/hydra/tree/v0.7.7) (2017-02-11)
[Full Changelog](https://github.com/ory/hydra/compare/v0.7.4...v0.7.7)

## [v0.7.4](https://github.com/ory/hydra/tree/v0.7.4) (2017-02-11)
[Full Changelog](https://github.com/ory/hydra/compare/v0.7.5...v0.7.4)

## [v0.7.5](https://github.com/ory/hydra/tree/v0.7.5) (2017-02-11)
[Full Changelog](https://github.com/ory/hydra/compare/v0.7.6...v0.7.5)

## [v0.7.6](https://github.com/ory/hydra/tree/v0.7.6) (2017-02-11)
[Full Changelog](https://github.com/ory/hydra/compare/v0.7.3...v0.7.6)

**Implemented enhancements:**

- sql: limit maximum open connections, document timeout options through DSN [\#359](https://github.com/ory/hydra/issues/359)

**Fixed bugs:**

- oauth2: invalid consent response causes panic [\#369](https://github.com/ory/hydra/issues/369)
- oauth2: resolve issue with cookie store [\#376](https://github.com/ory/hydra/pull/376) ([aeneasr](https://github.com/aeneasr))

**Closed issues:**

- Can hydra be easily integrated \(embedded\) into any golang http application? [\#372](https://github.com/ory/hydra/issues/372)

**Merged pull requests:**

- oauth2: invalid consent response causes panic - closes  \#369 [\#370](https://github.com/ory/hydra/pull/370) ([aeneasr](https://github.com/aeneasr))
- Resolve issues with SQL maximum open connections [\#360](https://github.com/ory/hydra/pull/360) ([aeneasr](https://github.com/aeneasr))

## [v0.7.3](https://github.com/ory/hydra/tree/v0.7.3) (2017-01-22)
[Full Changelog](https://github.com/ory/hydra/compare/v0.7.2...v0.7.3)

**Fixed bugs:**

- policy: investigate potential sql connection leak - closes \#363 [\#365](https://github.com/ory/hydra/pull/365) ([aeneasr](https://github.com/aeneasr))

**Closed issues:**

- Have Hydra store usernames linked to tokens [\#364](https://github.com/ory/hydra/issues/364)
- policy: investigate potential sql connection leak [\#363](https://github.com/ory/hydra/issues/363)
- crypto/bcrypt: hashedPassword is not the hash of the given password [\#346](https://github.com/ory/hydra/issues/346)

**Merged pull requests:**

- Update fosite\_store\_redis.go [\#361](https://github.com/ory/hydra/pull/361) ([itsjamie](https://github.com/itsjamie))

## [v0.7.2](https://github.com/ory/hydra/tree/v0.7.2) (2017-01-02)
[Full Changelog](https://github.com/ory/hydra/compare/v0.7.1...v0.7.2)

**Fixed bugs:**

- Problems with the authorization code flow [\#342](https://github.com/ory/hydra/issues/342)
- sql: deleting policies does not delete associated records with mysql driver [\#326](https://github.com/ory/hydra/issues/326)
- vendor: update to fosite 0.6.11 - closes \#338 [\#343](https://github.com/ory/hydra/pull/343) ([aeneasr](https://github.com/aeneasr))

**Closed issues:**

- oidc: at\_hash / c\_hash mismatch [\#338](https://github.com/ory/hydra/issues/338)
- oidc: SCIM compliance [\#320](https://github.com/ory/hydra/issues/320)

**Merged pull requests:**

- vendor: update to fosite 0.6.12 [\#344](https://github.com/ory/hydra/pull/344) ([aeneasr](https://github.com/aeneasr))

## [v0.7.1](https://github.com/ory/hydra/tree/v0.7.1) (2016-12-30)
[Full Changelog](https://github.com/ory/hydra/compare/v0.7.0...v0.7.1)

## [v0.7.0](https://github.com/ory/hydra/tree/v0.7.0) (2016-12-30)
[Full Changelog](https://github.com/ory/hydra/compare/v0.6.10...v0.7.0)

**Implemented enhancements:**

- Implement RemoveSubjectFromPolicy and RemoveResourceFromPolicy [\#336](https://github.com/ory/hydra/issues/336)
- policy: provide rest endpoint for policy updates [\#305](https://github.com/ory/hydra/issues/305)
- 0.7.0: SQL Migrate, Groups, Hardening [\#329](https://github.com/ory/hydra/pull/329) ([aeneasr](https://github.com/aeneasr))

**Fixed bugs:**

- 0.7.0: SQL Migrate, Groups, Hardening [\#329](https://github.com/ory/hydra/pull/329) ([aeneasr](https://github.com/aeneasr))

**Closed issues:**

- Replace \# with ? in authentication response [\#337](https://github.com/ory/hydra/issues/337)

## [v0.6.10](https://github.com/ory/hydra/tree/v0.6.10) (2016-12-26)
[Full Changelog](https://github.com/ory/hydra/compare/v0.6.9...v0.6.10)

**Implemented enhancements:**

- oauth2/consent: force jti echo in consent response [\#322](https://github.com/ory/hydra/issues/322)
- include a migration routine for databases [\#194](https://github.com/ory/hydra/issues/194)
- warden: add group management and group based policy checks [\#68](https://github.com/ory/hydra/issues/68)
- Improve http-based warden/introspection error responses [\#335](https://github.com/ory/hydra/pull/335) ([aeneasr](https://github.com/aeneasr))

## [v0.6.9](https://github.com/ory/hydra/tree/v0.6.9) (2016-12-20)
[Full Changelog](https://github.com/ory/hydra/compare/v0.6.8...v0.6.9)

**Implemented enhancements:**

- cmd: add configuration options for `hydra token user` [\#327](https://github.com/ory/hydra/issues/327)
- core: add api key flow [\#234](https://github.com/ory/hydra/issues/234)

**Fixed bugs:**

-  openid: support response\_type=code id\_token - closes \#332 [\#333](https://github.com/ory/hydra/pull/333) ([aeneasr](https://github.com/aeneasr))

**Closed issues:**

- openid: support response\_type=code id\_token [\#332](https://github.com/ory/hydra/issues/332)
- Apparent failure on load with ECDSA key [\#328](https://github.com/ory/hydra/issues/328)
- Why hydra github homepage crash when I visit \( while scrolling down\) [\#323](https://github.com/ory/hydra/issues/323)
- JsonWebTokenError: jwt must be provided [\#321](https://github.com/ory/hydra/issues/321)
- write tests for cmd helpers [\#186](https://github.com/ory/hydra/issues/186)

**Merged pull requests:**

- cmd: replace newline in HTTP\_TLS [\#331](https://github.com/ory/hydra/pull/331) ([ewilde](https://github.com/ewilde))
- Log fixes [\#324](https://github.com/ory/hydra/pull/324) ([johnwu96822](https://github.com/johnwu96822))

## [v0.6.8](https://github.com/ory/hydra/tree/v0.6.8) (2016-12-06)
[Full Changelog](https://github.com/ory/hydra/compare/v0.6.7...v0.6.8)

**Implemented enhancements:**

- oauth2: http introspector should return well known error [\#319](https://github.com/ory/hydra/pull/319) ([aeneasr](https://github.com/aeneasr))

## [v0.6.7](https://github.com/ory/hydra/tree/v0.6.7) (2016-12-04)
[Full Changelog](https://github.com/ory/hydra/compare/v0.6.6...v0.6.7)

**Merged pull requests:**

- all: improve cli and oauth2 error reporting [\#318](https://github.com/ory/hydra/pull/318) ([aeneasr](https://github.com/aeneasr))

## [v0.6.6](https://github.com/ory/hydra/tree/v0.6.6) (2016-12-04)
[Full Changelog](https://github.com/ory/hydra/compare/v0.6.5...v0.6.6)

**Implemented enhancements:**

- core: Redis backend [\#306](https://github.com/ory/hydra/issues/306)

**Closed issues:**

- oauth2: aud parameter does not allow arrays [\#314](https://github.com/ory/hydra/issues/314)

**Merged pull requests:**

- add missing work in docs/oauth2.md [\#317](https://github.com/ory/hydra/pull/317) ([bbigras](https://github.com/bbigras))
- docker: --name should be before the image's name [\#316](https://github.com/ory/hydra/pull/316) ([bbigras](https://github.com/bbigras))

## [v0.6.5](https://github.com/ory/hydra/tree/v0.6.5) (2016-11-28)
[Full Changelog](https://github.com/ory/hydra/compare/v0.6.4...v0.6.5)

**Implemented enhancements:**

- store/redis: redis backend for hydra [\#313](https://github.com/ory/hydra/pull/313) ([115100](https://github.com/115100))

## [v0.6.4](https://github.com/ory/hydra/tree/v0.6.4) (2016-11-22)
[Full Changelog](https://github.com/ory/hydra/compare/v0.6.3...v0.6.4)

**Implemented enhancements:**

- oauth2/revocation: token revocation fails silently with sql store [\#312](https://github.com/ory/hydra/pull/312) ([aeneasr](https://github.com/aeneasr))

**Fixed bugs:**

- oauth2/revocation: token revocation fails silently with sql store [\#311](https://github.com/ory/hydra/issues/311)
- oauth2/revocation: token revocation fails silently with sql store [\#312](https://github.com/ory/hydra/pull/312) ([aeneasr](https://github.com/aeneasr))

**Closed issues:**

- docs: clean up TokenValid leftovers [\#310](https://github.com/ory/hydra/issues/310)

## [v0.6.3](https://github.com/ory/hydra/tree/v0.6.3) (2016-11-17)
[Full Changelog](https://github.com/ory/hydra/compare/v0.6.2...v0.6.3)

**Implemented enhancements:**

- Rejection reason code to /warden/token/allowed [\#308](https://github.com/ory/hydra/issues/308)

**Fixed bugs:**

- oauth2: resolve issues with token introspection on user tokens [\#309](https://github.com/ory/hydra/pull/309) ([aeneasr](https://github.com/aeneasr))

## [v0.6.2](https://github.com/ory/hydra/tree/v0.6.2) (2016-11-05)
[Full Changelog](https://github.com/ory/hydra/compare/v0.6.1...v0.6.2)

**Implemented enhancements:**

- github: comply with Go license terms [\#300](https://github.com/ory/hydra/issues/300)

**Merged pull requests:**

- Fix client SQL manager missing client\_name [\#303](https://github.com/ory/hydra/pull/303) ([johnwu96822](https://github.com/johnwu96822))

## [v0.6.1](https://github.com/ory/hydra/tree/v0.6.1) (2016-10-26)
[Full Changelog](https://github.com/ory/hydra/compare/v0.6.0...v0.6.1)

**Fixed bugs:**

- MySQL DB not creating on start – JSON column types only supported from MySQL 5.7 and onwards [\#299](https://github.com/ory/hydra/issues/299)
- 0.6.1 [\#301](https://github.com/ory/hydra/pull/301) ([aeneasr](https://github.com/aeneasr))

**Merged pull requests:**

- Fix some minor typos and the broken tutorial links [\#298](https://github.com/ory/hydra/pull/298) ([justinclift](https://github.com/justinclift))

## [v0.6.0](https://github.com/ory/hydra/tree/v0.6.0) (2016-10-25)
[Full Changelog](https://github.com/ory/hydra/compare/v0.5.8...v0.6.0)

**Implemented enhancements:**

- Make it possible for travis-ci to build forked repos [\#295](https://github.com/ory/hydra/issues/295)
- core: add sql support [\#292](https://github.com/ory/hydra/issues/292)
- travis: execute gox build only when new commit is a new tag [\#285](https://github.com/ory/hydra/issues/285)
- cmd: prettify the `hydra token user` output [\#281](https://github.com/ory/hydra/issues/281)
- warden: make it clear that ladon.Request.Subject is not required or break bc and remove it [\#270](https://github.com/ory/hydra/issues/270)
- connections: remove connections API [\#265](https://github.com/ory/hydra/issues/265)
- consider signing up for Core Infrastructure Initiative badge [\#246](https://github.com/ory/hydra/issues/246)
- oauth2: token revocation endpoint [\#233](https://github.com/ory/hydra/issues/233)
- oauth2/rethinkdb: clear expired access tokens from memory [\#228](https://github.com/ory/hydra/issues/228)
- 0.6.0 [\#293](https://github.com/ory/hydra/pull/293) ([aeneasr](https://github.com/aeneasr))

**Fixed bugs:**

- all: coverage report is missing covered lines of nested packages [\#296](https://github.com/ory/hydra/issues/296)
- oauth2/introspect: make endpoint rfc7662 compatible [\#289](https://github.com/ory/hydra/issues/289)
- rethink: figure out how to deal with unreliable changefeed [\#269](https://github.com/ory/hydra/issues/269)
- oauth2: requests waste a lot of time in fosite storer `requestFromRDB\(\)` routine [\#260](https://github.com/ory/hydra/issues/260)
- 0.6.0 [\#293](https://github.com/ory/hydra/pull/293) ([aeneasr](https://github.com/aeneasr))

**Closed issues:**

- docs: fix typo in consent.md [\#294](https://github.com/ory/hydra/issues/294)
- docs/apiary: add at\_ext note to warden endpoints [\#287](https://github.com/ory/hydra/issues/287)
- core/storage: with rethinkdb being closed, what is our path forward? [\#286](https://github.com/ory/hydra/issues/286)
- docs: warden resource names are wrong on apiary [\#268](https://github.com/ory/hydra/issues/268)
- Request for Comment: Fair Source License / Business Source License [\#227](https://github.com/ory/hydra/issues/227)
- core: \(health\) monitoring endpoint [\#216](https://github.com/ory/hydra/issues/216)
- add much simpler identity provider and oauth2 consumer example [\#172](https://github.com/ory/hydra/issues/172)
- 2fa: add two factor authentication helper API [\#69](https://github.com/ory/hydra/issues/69)

**Merged pull requests:**

- cmd: fix typo in host command help text [\#291](https://github.com/ory/hydra/pull/291) ([faxal](https://github.com/faxal))
- travis: Only gox build on tags and go1.7 [\#288](https://github.com/ory/hydra/pull/288) ([emilva](https://github.com/emilva))
- docs: improve introduction [\#267](https://github.com/ory/hydra/pull/267) ([aeneasr](https://github.com/aeneasr))

## [v0.5.8](https://github.com/ory/hydra/tree/v0.5.8) (2016-10-06)
[Full Changelog](https://github.com/ory/hydra/compare/v0.5.7...v0.5.8)

**Fixed bugs:**

- oauth2: refresh token does not migrate session object to new token [\#283](https://github.com/ory/hydra/issues/283)
- oauth2: refresh token does not migrate session object to new token [\#284](https://github.com/ory/hydra/pull/284) ([aeneasr](https://github.com/aeneasr))

## [v0.5.7](https://github.com/ory/hydra/tree/v0.5.7) (2016-10-04)
[Full Changelog](https://github.com/ory/hydra/compare/v0.5.6...v0.5.7)

**Implemented enhancements:**

- jwk: add use parameter to generated JWKs [\#279](https://github.com/ory/hydra/issues/279)
- jwk: add use parameter to generated JWKs - closes \#279 [\#280](https://github.com/ory/hydra/pull/280) ([aeneasr](https://github.com/aeneasr))

## [v0.5.6](https://github.com/ory/hydra/tree/v0.5.6) (2016-10-03)
[Full Changelog](https://github.com/ory/hydra/compare/v0.5.5...v0.5.6)

**Implemented enhancements:**

- oauth2: scopes should be separated by %20 and not +, to ensure javascript compatibility [\#278](https://github.com/ory/hydra/pull/278) ([aeneasr](https://github.com/aeneasr))

**Fixed bugs:**

- cmd: hydra help host profiling typo [\#274](https://github.com/ory/hydra/issues/274)
- cmd: hydra help host typos [\#272](https://github.com/ory/hydra/issues/272)

**Closed issues:**

- Scopes should be separated by %20 and not +, to ensure javascript compatibility [\#277](https://github.com/ory/hydra/issues/277)

**Merged pull requests:**

- cmd: fix \#272 typos in the host command controls [\#276](https://github.com/ory/hydra/pull/276) ([cixtor](https://github.com/cixtor))
- Fix \#274 - replace HYDRA\_PROFILING with PROFILING [\#275](https://github.com/ory/hydra/pull/275) ([otremblay](https://github.com/otremblay))

## [v0.5.5](https://github.com/ory/hydra/tree/v0.5.5) (2016-09-29)
[Full Changelog](https://github.com/ory/hydra/compare/v0.5.4...v0.5.5)

## [v0.5.4](https://github.com/ory/hydra/tree/v0.5.4) (2016-09-29)
[Full Changelog](https://github.com/ory/hydra/compare/v0.5.3...v0.5.4)

## [v0.5.3](https://github.com/ory/hydra/tree/v0.5.3) (2016-09-29)
[Full Changelog](https://github.com/ory/hydra/compare/v0.5.2...v0.5.3)

**Implemented enhancements:**

- docker: add http-only dockerfile and upgrade to go 1.7 base image [\#273](https://github.com/ory/hydra/pull/273) ([aeneasr](https://github.com/aeneasr))

**Fixed bugs:**

- investigate if and why slow rethinkdb connection causes client root to be recreated [\#191](https://github.com/ory/hydra/issues/191)

**Closed issues:**

- Consider extract Go SDK package into separate repository [\#266](https://github.com/ory/hydra/issues/266)
- Showcase: How and where are you using Hydra? [\#115](https://github.com/ory/hydra/issues/115)

## [v0.5.2](https://github.com/ory/hydra/tree/v0.5.2) (2016-09-23)
[Full Changelog](https://github.com/ory/hydra/compare/v0.5.0...v0.5.2)

## [v0.5.0](https://github.com/ory/hydra/tree/v0.5.0) (2016-09-22)
[Full Changelog](https://github.com/ory/hydra/compare/v0.5.1...v0.5.0)

## [v0.5.1](https://github.com/ory/hydra/tree/v0.5.1) (2016-09-22)
[Full Changelog](https://github.com/ory/hydra/compare/v0.4.2-alpha.4...v0.5.1)

**Implemented enhancements:**

- oauth2: include original request query parameters in the consent challenge [\#256](https://github.com/ory/hydra/issues/256)
- Need a better health check for a load balancer [\#251](https://github.com/ory/hydra/issues/251)
- client: add ability to update client [\#250](https://github.com/ory/hydra/issues/250)
- oauth2: allow access token validation for public clients [\#245](https://github.com/ory/hydra/issues/245)
- all: improve error messages regarding token validation [\#244](https://github.com/ory/hydra/issues/244)
- all: resolve naming inconsistencies in jwk set names used in hydra [\#239](https://github.com/ory/hydra/issues/239)
- sdk: resolve naming inconsistencies [\#226](https://github.com/ory/hydra/issues/226)
- oidc: support kid hint in header [\#222](https://github.com/ory/hydra/issues/222)
- 0.5.0-errors [\#263](https://github.com/ory/hydra/pull/263) ([aeneasr](https://github.com/aeneasr))
- 0.5.0 [\#243](https://github.com/ory/hydra/pull/243) ([aeneasr](https://github.com/aeneasr))

**Fixed bugs:**

- When invalid/expired token is used for /warden/allowed endpoint, status 500 is returned [\#262](https://github.com/ory/hydra/issues/262)
- docs: fix images in readme [\#261](https://github.com/ory/hydra/issues/261)
- Bad HTML encoding of the scope parameter [\#259](https://github.com/ory/hydra/issues/259)
- docs: images are broken [\#258](https://github.com/ory/hydra/issues/258)
- oauth2: id token hashes are not base64 url encoded [\#255](https://github.com/ory/hydra/issues/255)
- oauth2: state parameter is missing when response\_type=id\_token [\#254](https://github.com/ory/hydra/issues/254)
- jwk: anonymous request can't read public keys [\#253](https://github.com/ory/hydra/issues/253)
- travis: ld flags are wrong [\#242](https://github.com/ory/hydra/issues/242)
- cmd: hydra token user should show id token in browser [\#224](https://github.com/ory/hydra/issues/224)
- oidc: hybrid flow using `token+code+id\_token` returns multiple tokens of the same type [\#223](https://github.com/ory/hydra/issues/223)
- hydra clients import doesn't print client's secret [\#221](https://github.com/ory/hydra/issues/221)
- 0.5.0-errors [\#263](https://github.com/ory/hydra/pull/263) ([aeneasr](https://github.com/aeneasr))
- 0.5.0 [\#243](https://github.com/ory/hydra/pull/243) ([aeneasr](https://github.com/aeneasr))

**Closed issues:**

- core: document hard-wired JWK sets [\#247](https://github.com/ory/hydra/issues/247)
- managing client definitions [\#197](https://github.com/ory/hydra/issues/197)

**Merged pull requests:**

- docs: add notes on operational considerations [\#252](https://github.com/ory/hydra/pull/252) ([boyvinall](https://github.com/boyvinall))

## [v0.4.2-alpha.4](https://github.com/ory/hydra/tree/v0.4.2-alpha.4) (2016-09-03)
[Full Changelog](https://github.com/ory/hydra/compare/v0.4.2...v0.4.2-alpha.4)

## [v0.4.2](https://github.com/ory/hydra/tree/v0.4.2) (2016-09-03)
[Full Changelog](https://github.com/ory/hydra/compare/v0.4.3...v0.4.2)

## [v0.4.3](https://github.com/ory/hydra/tree/v0.4.3) (2016-09-03)
[Full Changelog](https://github.com/ory/hydra/compare/v0.4.2-alpha.3...v0.4.3)

## [v0.4.2-alpha.3](https://github.com/ory/hydra/tree/v0.4.2-alpha.3) (2016-09-02)
[Full Changelog](https://github.com/ory/hydra/compare/v0.4.2-alpha.2...v0.4.2-alpha.3)

## [v0.4.2-alpha.2](https://github.com/ory/hydra/tree/v0.4.2-alpha.2) (2016-09-01)
[Full Changelog](https://github.com/ory/hydra/compare/v0.4.2-alpha.1...v0.4.2-alpha.2)

## [v0.4.2-alpha.1](https://github.com/ory/hydra/tree/v0.4.2-alpha.1) (2016-09-01)
[Full Changelog](https://github.com/ory/hydra/compare/0.4.2-alpha...v0.4.2-alpha.1)

## [0.4.2-alpha](https://github.com/ory/hydra/tree/0.4.2-alpha) (2016-09-01)
[Full Changelog](https://github.com/ory/hydra/compare/v0.4.1...0.4.2-alpha)

**Implemented enhancements:**

- Add version option to Hydra's CLI [\#218](https://github.com/ory/hydra/issues/218)
- autobuild [\#240](https://github.com/ory/hydra/pull/240) ([aeneasr](https://github.com/aeneasr))
- Update jwt-go and resolve warden regression issue [\#232](https://github.com/ory/hydra/pull/232) ([aeneasr](https://github.com/aeneasr))

**Fixed bugs:**

- warden: firewal.Audience overridden with requesting clients subject [\#236](https://github.com/ory/hydra/pull/236) ([faxal](https://github.com/faxal))
- Update jwt-go and resolve warden regression issue [\#232](https://github.com/ory/hydra/pull/232) ([aeneasr](https://github.com/aeneasr))

**Closed issues:**

- how to use hydra without "--dangerous-auto-logon"? [\#241](https://github.com/ory/hydra/issues/241)
- warden: firewal.Audience overridden with requesting clients subject  [\#237](https://github.com/ory/hydra/issues/237)
- Vendor: Upgrade to jwt-go 3.0.0 [\#229](https://github.com/ory/hydra/issues/229)
- docs: warden sdk example is misleading [\#225](https://github.com/ory/hydra/issues/225)
- Typo in the apiary documentation [\#220](https://github.com/ory/hydra/issues/220)
- Importing clients with the CLI doesn't work [\#219](https://github.com/ory/hydra/issues/219)
- doc: add "what is hydra not?" section to readme [\#217](https://github.com/ory/hydra/issues/217)
- figure out a process to autobuild releases [\#210](https://github.com/ory/hydra/issues/210)

**Merged pull requests:**

- fix broken link for tutorial in README.md [\#213](https://github.com/ory/hydra/pull/213) ([allan-simon](https://github.com/allan-simon))

## [v0.4.1](https://github.com/ory/hydra/tree/v0.4.1) (2016-08-18)
[Full Changelog](https://github.com/ory/hydra/compare/v0.4.0...v0.4.1)

**Fixed bugs:**

- error bad request when running tutorial [\#211](https://github.com/ory/hydra/issues/211)
- cmd: resolve issue with token user flow [\#212](https://github.com/ory/hydra/pull/212) ([aeneasr](https://github.com/aeneasr))

## [v0.4.0](https://github.com/ory/hydra/tree/v0.4.0) (2016-08-17)
[Full Changelog](https://github.com/ory/hydra/compare/v0.3.1...v0.4.0)

**Implemented enhancements:**

- all: move docs from gitbook to github [\#204](https://github.com/ory/hydra/issues/204)
- 0.4.0 [\#203](https://github.com/ory/hydra/pull/203) ([aeneasr](https://github.com/aeneasr))

**Fixed bugs:**

- 0.4.0-prefix [\#209](https://github.com/ory/hydra/pull/209) ([aeneasr](https://github.com/aeneasr))
- 0.4.0 [\#203](https://github.com/ory/hydra/pull/203) ([aeneasr](https://github.com/aeneasr))

**Closed issues:**

- docs/guide: warden docs are outdated [\#206](https://github.com/ory/hydra/issues/206)
- fix sdk examples in readme [\#196](https://github.com/ory/hydra/issues/196)
- add tests for clients import [\#163](https://github.com/ory/hydra/issues/163)
- remove go get -t ./... from travis [\#71](https://github.com/ory/hydra/issues/71)

## [v0.3.1](https://github.com/ory/hydra/tree/v0.3.1) (2016-08-17)
[Full Changelog](https://github.com/ory/hydra/compare/v0.3.0...v0.3.1)

**Implemented enhancements:**

- oauth2: introspection should return custom session values [\#205](https://github.com/ory/hydra/issues/205)
- warden: move IntrospectToken from warden sdk to oauth2 [\#201](https://github.com/ory/hydra/issues/201)
- warden: rename InspectToken to IntrospectToken [\#200](https://github.com/ory/hydra/issues/200)

**Fixed bugs:**

- AccessTokens get overridden during startup of hydra [\#207](https://github.com/ory/hydra/issues/207)
- warden: IntrospectToken always throws an error on Hydra logs [\#199](https://github.com/ory/hydra/issues/199)
- resolve issue with at extra data [\#198](https://github.com/ory/hydra/issues/198)
- Fix 207 [\#208](https://github.com/ory/hydra/pull/208) ([aeneasr](https://github.com/aeneasr))

## [v0.3.0](https://github.com/ory/hydra/tree/v0.3.0) (2016-08-09)
[Full Changelog](https://github.com/ory/hydra/compare/v0.2.0...v0.3.0)

**Implemented enhancements:**

- 0.3.0 [\#195](https://github.com/ory/hydra/pull/195) ([aeneasr](https://github.com/aeneasr))

**Fixed bugs:**

- 0.3.0 [\#195](https://github.com/ory/hydra/pull/195) ([aeneasr](https://github.com/aeneasr))

## [v0.2.0](https://github.com/ory/hydra/tree/v0.2.0) (2016-08-09)
[Full Changelog](https://github.com/ory/hydra/compare/0.1-beta.4...v0.2.0)

**Implemented enhancements:**

- warden sdk should not make distinction between token and request [\#190](https://github.com/ory/hydra/issues/190)
- core scope should not be mandatory [\#189](https://github.com/ory/hydra/issues/189)
- id token claims should be set by consent challenge `id\_token` claim [\#188](https://github.com/ory/hydra/issues/188)
- provide default consent endpoint in hydra [\#185](https://github.com/ory/hydra/issues/185)
- make bcrypt cost configurable [\#184](https://github.com/ory/hydra/issues/184)
- make lifespans configurable [\#183](https://github.com/ory/hydra/issues/183)
- improve env to config [\#182](https://github.com/ory/hydra/issues/182)
- add memory profiling and cpu profiling [\#179](https://github.com/ory/hydra/issues/179)
- add basic http request logging [\#178](https://github.com/ory/hydra/issues/178)
- support edge tls termination [\#177](https://github.com/ory/hydra/issues/177)
- Make client HTTPManager not compatible with fosite.Storage [\#173](https://github.com/ory/hydra/issues/173)
- clean up stale branches [\#171](https://github.com/ory/hydra/issues/171)
- improve hydra connect dialogue [\#170](https://github.com/ory/hydra/issues/170)
- investigate if token creation can be speeded up [\#168](https://github.com/ory/hydra/issues/168)
- consent: allow proxying of id token claims [\#167](https://github.com/ory/hydra/issues/167)
- warden: rename authorized / allowed endpoints to something more meaningful [\#162](https://github.com/ory/hydra/issues/162)
- warden: rename `assertion` to `token` [\#158](https://github.com/ory/hydra/issues/158)
- Implement strict mode for warden [\#156](https://github.com/ory/hydra/issues/156)
- Implement token introspection endpoint [\#155](https://github.com/ory/hydra/issues/155)
- Don't log database credentials [\#147](https://github.com/ory/hydra/issues/147)
- OpenID Connect Session Management  [\#143](https://github.com/ory/hydra/issues/143)
- \[Feature request\] Import clients on startup [\#140](https://github.com/ory/hydra/issues/140)
- Warden for anonymous users [\#139](https://github.com/ory/hydra/issues/139)
- oauth2/consent: id token expiry should be configurable [\#127](https://github.com/ory/hydra/issues/127)
- warden: endpoint should only require valid client, not policy based access control [\#121](https://github.com/ory/hydra/issues/121)
- Improve error message of wrong system secret [\#104](https://github.com/ory/hydra/issues/104)
- warden: rename authorized / allowed endpoints to something more meaningful [\#187](https://github.com/ory/hydra/pull/187) ([aeneasr](https://github.com/aeneasr))
- 0.2.0 [\#165](https://github.com/ory/hydra/pull/165) ([aeneasr](https://github.com/aeneasr))
- all: add test cases for methods returning slices or maps of entities [\#152](https://github.com/ory/hydra/pull/152) ([aeneasr](https://github.com/aeneasr))
- Resolve rethinkdb connection when idle [\#148](https://github.com/ory/hydra/pull/148) ([aeneasr](https://github.com/aeneasr))
- all: resolve issues with the sdk and cli [\#142](https://github.com/ory/hydra/pull/142) ([aeneasr](https://github.com/aeneasr))
- cli: add token validation [\#134](https://github.com/ory/hydra/pull/134) ([aeneasr](https://github.com/aeneasr))
- Add wrapper library for HTTP Managers [\#130](https://github.com/ory/hydra/pull/130) ([faxal](https://github.com/faxal))

**Fixed bugs:**

- investigate runtime panic on warden allowed [\#181](https://github.com/ory/hydra/issues/181)
- oauth2 implicit flow should allow custom protocols [\#180](https://github.com/ory/hydra/issues/180)
- support edge tls termination [\#177](https://github.com/ory/hydra/issues/177)
- Token generation should be always consistent, not eventually consistent [\#176](https://github.com/ory/hydra/issues/176)
- consent: allow proxying of id token claims [\#167](https://github.com/ory/hydra/issues/167)
- config: do not store database config in hydra config [\#164](https://github.com/ory/hydra/issues/164)
- OAuth2 token endpoint does not allow GET method but reads query parameters [\#160](https://github.com/ory/hydra/issues/160)
- OAuth2 token endpoint should be able to handle simple form encoded requests [\#159](https://github.com/ory/hydra/issues/159)
- --dry option does not work correctly [\#157](https://github.com/ory/hydra/issues/157)
- client.GetClients\(\) returns invalid information [\#150](https://github.com/ory/hydra/issues/150)
- RethinkDB connection dies after a certain amount of inactive time [\#146](https://github.com/ory/hydra/issues/146)
- Fails to startup when a SSO connection is added. [\#141](https://github.com/ory/hydra/issues/141)
- id\_token: at\_hash / c\_hash is null [\#129](https://github.com/ory/hydra/issues/129)
- oauth2: some scopes are included twice [\#126](https://github.com/ory/hydra/issues/126)
- warden: iat / exp values are not being set [\#125](https://github.com/ory/hydra/issues/125)
- investigate missing scopes issue [\#124](https://github.com/ory/hydra/issues/124)
- rethinkdb: resolve an issue where missing refresh tokens cause duplicate key error [\#122](https://github.com/ory/hydra/issues/122)
- 0.2.0 [\#165](https://github.com/ory/hydra/pull/165) ([aeneasr](https://github.com/aeneasr))
- ensure client endpoint is initialised for CLI "clients import" command [\#149](https://github.com/ory/hydra/pull/149) ([boyvinall](https://github.com/boyvinall))
- Resolve rethinkdb connection when idle [\#148](https://github.com/ory/hydra/pull/148) ([aeneasr](https://github.com/aeneasr))
- all: resolve issues with the sdk and cli [\#142](https://github.com/ory/hydra/pull/142) ([aeneasr](https://github.com/aeneasr))
- Resolve warden issues [\#128](https://github.com/ory/hydra/pull/128) ([aeneasr](https://github.com/aeneasr))
- Various bugfixes [\#123](https://github.com/ory/hydra/pull/123) ([aeneasr](https://github.com/aeneasr))

**Closed issues:**

- Error trying to create a token via curl [\#174](https://github.com/ory/hydra/issues/174)
- gorethink: could not decode type \[\]uint8 into Go value of type string [\#169](https://github.com/ory/hydra/issues/169)
- document warden interface sdk [\#166](https://github.com/ory/hydra/issues/166)
- Document what OpenID Connect is and how to use it [\#154](https://github.com/ory/hydra/issues/154)
- Warden endpoints [\#137](https://github.com/ory/hydra/issues/137)
- Environment variables naming scheme [\#136](https://github.com/ory/hydra/issues/136)
- Implicit Flow redirect\_uri  does not match [\#133](https://github.com/ory/hydra/issues/133)
- hydra 2FA on cloud providers [\#132](https://github.com/ory/hydra/issues/132)
- Document HTTP client libraries for go [\#101](https://github.com/ory/hydra/issues/101)
- Document error redirect to identity provider [\#96](https://github.com/ory/hydra/issues/96)
- use dropbox example to explain oauth2 [\#95](https://github.com/ory/hydra/issues/95)

**Merged pull requests:**

- client: fix client.GetClients\(\) for multiple clients [\#151](https://github.com/ory/hydra/pull/151) ([boyvinall](https://github.com/boyvinall))
- readme: Fix table of contents links [\#145](https://github.com/ory/hydra/pull/145) ([smithrobs](https://github.com/smithrobs))
- doc: Minor grammar/spelling fixes for README [\#144](https://github.com/ory/hydra/pull/144) ([smithrobs](https://github.com/smithrobs))
- Add some precisions to installation [\#131](https://github.com/ory/hydra/pull/131) ([yageek](https://github.com/yageek))

## [0.1-beta.4](https://github.com/ory/hydra/tree/0.1-beta.4) (2016-06-26)
[Full Changelog](https://github.com/ory/hydra/compare/0.1-beta.3...0.1-beta.4)

**Implemented enhancements:**

- Connect to rethinkdb over SSL with self-signed certificate [\#114](https://github.com/ory/hydra/issues/114)

**Fixed bugs:**

- clients endpoint returns client secret base64 encoded [\#119](https://github.com/ory/hydra/issues/119)
- firewall 403s on warden endpoints [\#118](https://github.com/ory/hydra/issues/118)
- Client secrets should not be hashed when POSTing [\#113](https://github.com/ory/hydra/issues/113)
- Resolve issues with warden and client api [\#120](https://github.com/ory/hydra/pull/120) ([aeneasr](https://github.com/aeneasr))
- client: return client secret on POST and remove it from GET [\#117](https://github.com/ory/hydra/pull/117) ([aeneasr](https://github.com/aeneasr))

**Merged pull requests:**

- Connect to rethinkdb with a custom certificate [\#116](https://github.com/ory/hydra/pull/116) ([matteosuppo](https://github.com/matteosuppo))
- dist: fix typos in exemplary policies [\#112](https://github.com/ory/hydra/pull/112) ([aeneasr](https://github.com/aeneasr))

## [0.1-beta.3](https://github.com/ory/hydra/tree/0.1-beta.3) (2016-06-20)
[Full Changelog](https://github.com/ory/hydra/compare/0.1-beta.2...0.1-beta.3)

**Implemented enhancements:**

- docker: remove wait time on boot and use restart unless-stopped option [\#105](https://github.com/ory/hydra/pull/105) ([aeneasr](https://github.com/aeneasr))

**Fixed bugs:**

- Warden handlers are not mounted [\#109](https://github.com/ory/hydra/issues/109)

**Closed issues:**

- Installation fails [\#108](https://github.com/ory/hydra/issues/108)
- Exchange token from browser client [\#107](https://github.com/ory/hydra/issues/107)
- Temporary Client not working [\#106](https://github.com/ory/hydra/issues/106)
- Could not fetch initial state with docker-compose [\#103](https://github.com/ory/hydra/issues/103)

**Merged pull requests:**

- all: update jwt-go to versioned package and update dependencies [\#111](https://github.com/ory/hydra/pull/111) ([aeneasr](https://github.com/aeneasr))
- Mount warden handler [\#110](https://github.com/ory/hydra/pull/110) ([faxal](https://github.com/faxal))

## [0.1-beta.2](https://github.com/ory/hydra/tree/0.1-beta.2) (2016-06-14)
[Full Changelog](https://github.com/ory/hydra/compare/0.1-beta1...0.1-beta.2)

**Implemented enhancements:**

- CLI should have `-dry` option to show what the HTTP request looks like [\#99](https://github.com/ory/hydra/issues/99)
- Add offline scope for refresh tokens [\#97](https://github.com/ory/hydra/issues/97)
- extend jwk cert store [\#92](https://github.com/ory/hydra/issues/92)
- Creating clients with predefined credentials [\#91](https://github.com/ory/hydra/issues/91)
- Passing key and certificate to hydra [\#88](https://github.com/ory/hydra/issues/88)
- AES-GCM key should be sha256\(secret\)\[:32\] [\#86](https://github.com/ory/hydra/issues/86)
- Update GoRethink imports [\#78](https://github.com/ory/hydra/issues/78)
- link exemplary policies in the docs [\#75](https://github.com/ory/hydra/issues/75)
- support SAML in addition to OAuth2 [\#29](https://github.com/ory/hydra/issues/29)
- 0.1-beta2 [\#90](https://github.com/ory/hydra/pull/90) ([aeneasr](https://github.com/aeneasr))
- vendor: switch to versioned gorethink api [\#81](https://github.com/ory/hydra/pull/81) ([aeneasr](https://github.com/aeneasr))

**Fixed bugs:**

- fix issue where tls certificate is regenerated on boot [\#93](https://github.com/ory/hydra/issues/93)
- typo: singing instead of signing [\#89](https://github.com/ory/hydra/issues/89)
- 404 in the gitbook   [\#85](https://github.com/ory/hydra/issues/85)
- Update GoRethink imports [\#78](https://github.com/ory/hydra/issues/78)
- client: resolved that secrets can not be set when using http or cli [\#102](https://github.com/ory/hydra/pull/102) ([aeneasr](https://github.com/aeneasr))

**Closed issues:**

- document security architecture [\#82](https://github.com/ory/hydra/issues/82)
- go install fails [\#77](https://github.com/ory/hydra/issues/77)
- Security audit based on rfc6819 [\#42](https://github.com/ory/hydra/issues/42)

**Merged pull requests:**

- Fix typo of weather [\#100](https://github.com/ory/hydra/pull/100) ([smurfpandey](https://github.com/smurfpandey))
- readme: add security section [\#87](https://github.com/ory/hydra/pull/87) ([aeneasr](https://github.com/aeneasr))
- Fix idiom in README [\#79](https://github.com/ory/hydra/pull/79) ([neuhaus](https://github.com/neuhaus))

## [0.1-beta1](https://github.com/ory/hydra/tree/0.1-beta1) (2016-05-29)
**Implemented enhancements:**

- client rest endpoint: rename `name` to `client\_name` [\#72](https://github.com/ory/hydra/issues/72)
- allow using not self-signed TLS certificates [\#70](https://github.com/ory/hydra/issues/70)
- Implement OpenID Connect Dynamic Client Registration 1.0 [\#65](https://github.com/ory/hydra/issues/65)
- Implement default identity provider using postgres [\#63](https://github.com/ory/hydra/issues/63)
- Implement generic connectors [\#61](https://github.com/ory/hydra/issues/61)
- Replace osin with ory-am/fosite [\#46](https://github.com/ory/hydra/issues/46)
- Remove dockertest dependency from handlers [\#43](https://github.com/ory/hydra/issues/43)
- adding RethinkDB as a Store [\#39](https://github.com/ory/hydra/issues/39)
- Add more IdPs [\#33](https://github.com/ory/hydra/issues/33)
- Make JWT as access tokens optional and replace with a custom strategy [\#32](https://github.com/ory/hydra/issues/32)
- support for ldap for user storage [\#28](https://github.com/ory/hydra/issues/28)
- Migrate from mux to httprouter [\#14](https://github.com/ory/hydra/issues/14)
- Decompositioning, implement Fosite [\#62](https://github.com/ory/hydra/pull/62) ([aeneasr](https://github.com/aeneasr))

**Fixed bugs:**

- spec: /jwk/:set/:kid must return array  [\#74](https://github.com/ory/hydra/issues/74)
- client rest endpoint: rename `name` to `client\\_name` [\#72](https://github.com/ory/hydra/issues/72)
- Too many open files probably caused by http client [\#47](https://github.com/ory/hydra/issues/47)

**Closed issues:**

- Add Dockerfile for autobuild [\#60](https://github.com/ory/hydra/issues/60)
- CLI refactor and initial account set up [\#59](https://github.com/ory/hydra/issues/59)
- ory-am ssl cert invalid [\#58](https://github.com/ory/hydra/issues/58)
- Granted Endpoint Proposal: Performant access decisions for resource providers using REST [\#48](https://github.com/ory/hydra/issues/48)
- Security "audit" pre-analysis \(based on rfc6749\) [\#41](https://github.com/ory/hydra/issues/41)
- wrong repo [\#40](https://github.com/ory/hydra/issues/40)
- Rename providers to connectors [\#38](https://github.com/ory/hydra/issues/38)
- Are there standards for connecting to third party providers [\#37](https://github.com/ory/hydra/issues/37)
- Add support for scopes [\#36](https://github.com/ory/hydra/issues/36)
- Readme: Accounts CLI Usage [\#31](https://github.com/ory/hydra/issues/31)
- Continue using JWT as access tokens? [\#22](https://github.com/ory/hydra/issues/22)
- remove refresh token claims [\#21](https://github.com/ory/hydra/issues/21)
- godeps should only be commited on release [\#19](https://github.com/ory/hydra/issues/19)
- refactor POST workflow  [\#13](https://github.com/ory/hydra/issues/13)
- JWT assertions [\#5](https://github.com/ory/hydra/issues/5)
- Check JWT Algorithm [\#3](https://github.com/ory/hydra/issues/3)

**Merged pull requests:**

- Remove go get of govet in .travis.yml [\#67](https://github.com/ory/hydra/pull/67) ([sbani](https://github.com/sbani))
- Hydra is now using Go 1.6 vendoring and is deployable to heroku [\#56](https://github.com/ory/hydra/pull/56) ([aeneasr](https://github.com/aeneasr))
- Heroku [\#55](https://github.com/ory/hydra/pull/55) ([aeneasr](https://github.com/aeneasr))
- Update README.md [\#54](https://github.com/ory/hydra/pull/54) ([leetal](https://github.com/leetal))
- RethinkDB [\#53](https://github.com/ory/hydra/pull/53) ([leetal](https://github.com/leetal))
- handler.go:300: no formatting directive in Sprintf call [\#52](https://github.com/ory/hydra/pull/52) ([QuentinPerez](https://github.com/QuentinPerez))
- providers: added microsoft and improved existing providers [\#51](https://github.com/ory/hydra/pull/51) ([aeneasr](https://github.com/aeneasr))
- oauth: added google provider [\#50](https://github.com/ory/hydra/pull/50) ([aeneasr](https://github.com/aeneasr))
- handle multiple return values from gopass [\#49](https://github.com/ory/hydra/pull/49) ([timothyknight](https://github.com/timothyknight))
- doc: create MAINTAINERS [\#45](https://github.com/ory/hydra/pull/45) ([aeneasr](https://github.com/aeneasr))
- docs: create CONTRIBUTING.md [\#44](https://github.com/ory/hydra/pull/44) ([aeneasr](https://github.com/aeneasr))
- update accounts CLI Usage [\#34](https://github.com/ory/hydra/pull/34) ([akhedrane](https://github.com/akhedrane))
- Add a Gitter chat badge to README.md [\#30](https://github.com/ory/hydra/pull/30) ([gitter-badger](https://github.com/gitter-badger))
- Extra arguments [\#27](https://github.com/ory/hydra/pull/27) ([QuentinPerez](https://github.com/QuentinPerez))
- all: oauth and guard endpoints now accept basic auth instead of token… [\#26](https://github.com/ory/hydra/pull/26) ([aeneasr](https://github.com/aeneasr))
- account: refactor, more endpoints and tests [\#25](https://github.com/ory/hydra/pull/25) ([aeneasr](https://github.com/aeneasr))
- all: username instead of email, token revocation, introspect spec ali… [\#24](https://github.com/ory/hydra/pull/24) ([aeneasr](https://github.com/aeneasr))
- Tutorial [\#23](https://github.com/ory/hydra/pull/23) ([aeneasr](https://github.com/aeneasr))
- Unstaged [\#20](https://github.com/ory/hydra/pull/20) ([aeneasr](https://github.com/aeneasr))
- client: now tries to refresh when token is invalid [\#18](https://github.com/ory/hydra/pull/18) ([aeneasr](https://github.com/aeneasr))
- client: added possibility to skip CA check [\#17](https://github.com/ory/hydra/pull/17) ([aeneasr](https://github.com/aeneasr))
- cli: fixed default TLS and JWT filepaths [\#16](https://github.com/ory/hydra/pull/16) ([aeneasr](https://github.com/aeneasr))
- Policy changes and more tests [\#15](https://github.com/ory/hydra/pull/15) ([aeneasr](https://github.com/aeneasr))
- unstaged [\#12](https://github.com/ory/hydra/pull/12) ([aeneasr](https://github.com/aeneasr))
- Ladon api update & policy http endpoint [\#11](https://github.com/ory/hydra/pull/11) ([aeneasr](https://github.com/aeneasr))
- Improved CLI `client create` and provider workflow. [\#10](https://github.com/ory/hydra/pull/10) ([aeneasr](https://github.com/aeneasr))
- cli [\#9](https://github.com/ory/hydra/pull/9) ([aeneasr](https://github.com/aeneasr))
- all: increased test coverage [\#8](https://github.com/ory/hydra/pull/8) ([aeneasr](https://github.com/aeneasr))
- Handlers and cleanup [\#7](https://github.com/ory/hydra/pull/7) ([aeneasr](https://github.com/aeneasr))
- Single Sign On [\#6](https://github.com/ory/hydra/pull/6) ([aeneasr](https://github.com/aeneasr))
- tests: increased coverage [\#4](https://github.com/ory/hydra/pull/4) ([aeneasr](https://github.com/aeneasr))
- Implemented jwt, middleware, test coverage and handlers. [\#2](https://github.com/ory/hydra/pull/2) ([aeneasr](https://github.com/aeneasr))
- Refactor [\#1](https://github.com/ory/hydra/pull/1) ([aeneasr](https://github.com/aeneasr))



\* *This Change Log was automatically generated by [github_changelog_generator](https://github.com/skywinder/Github-Changelog-Generator)*