---
id: limitations
title: Limitations
---

ORY Hydra tries to solve all of OAuth 2.0 and OpenID Connect uses. There are,
however, some limitations.

## MySQL <= 5.6 / MariaDB

ORY Hydra has issues with MySQL <= 5.6 (but not MySQL 5.7+) and certain MariaDB
versions. Read more about this [here](https://github.com/ory/hydra/issues/377).
Our recommendation is to use MySQL 5.7+ or PostgreSQL.

## OAuth 2.0 Client Secret Length

OAuth 2.0 Client Secrets are hashed using BCrypt. BCrypt has, by design, an
maximum password length. The Golang BCrypt library has a maximum password length
of 73 bytes. Any password longer will be "truncated":

```
$ hydra clients create --id long-secret \
	--secret 525348e77144a9cee9a7471a8b67c50ea85b9e3eb377a3c1a3a23db88f9150eefe76e6a339fdbc62b817595f53d72549d9ebe36438f8c2619846b963e9f43a94 \
	--endpoint http://localhost:4445 \
	--token-endpoint-auth-method client_secret_post \
	--grant-types client_credentials

$ hydra token client --client-id long-secret \
	--client-secret 525348e77144a9cee9a7471a8b67c50ea85b9e3eb377a3c1a3a23db88f9150eefe76e6a3 \
	--endpoint http://localhost:4444
```

For more information on this topic we recommend reading:

- https://security.stackexchange.com/questions/39849/does-bcrypt-have-a-maximum-password-length
- https://security.stackexchange.com/questions/6623/pre-hash-password-before-applying-bcrypt-to-avoid-restricting-password-length

## Resource Owner Password Credentials Grant Type (ROCP)

ORY Hydra does not and will not implement the Resource Owner Password
Credentials Grant Type. Read on for context.

### Overview

This grant type allows OAuth 2.0 Clients to exchange user credentials (username,
password) for an access token.

**Request:**

```
POST /oauth2/token HTTP/1.1
Host: server.example.com
Authorization: Basic czZCaGRSa3F0MzpnWDFmQmF0M2JW
Content-Type: application/x-www-form-urlencoded

grant_type=password&username=johndoe&password=A3ddj3w
```

**Response:**

```
HTTP/1.1 200 OK
Content-Type: application/json;charset=UTF-8
Cache-Control: no-store
Pragma: no-cache

{
  "access_token":"2YotnFZFEjr1zCsicMWpAA",
  "token_type":"example",
  "expires_in":3600,
  "refresh_token":"tGzv3JOkF0XG5Qx2TlKWIA",
  "example_parameter":"example_value"
}
```

You might think that this is the perfect grant type for your first-party
application. This grant type is most commonly used in mobile authentication for
first-party apps. If you plan on doing this, stop right now and read
[this blog article](https://www.ory.sh/oauth2-for-mobile-app-spa-browser).

### Legacy & Bad Security

The ROCP grant type is discouraged by developers, professionals, and the IETF
itself. It was originally added because big legacy corporations (not dropping
any names, but they are part of the IETF consortium) did not want to migrate
their authentication infrastructure to the modern web but instead do what
they've been doing all along "but OAuth 2.0" and for systems that want to
upgrade from OAuth (1.0) to OAuth 2.0.

There are a ton of good reasons why this is a bad flow, they are summarized in
[this excellent blog article as well](https://www.scottbrady91.com/OAuth/Why-the-Resource-Owner-Password-Credentials-Grant-Type-is-not-Authentication-nor-Suitable-for-Modern-Applications).

### What about Auth0, Okta, ...?

Auth0, Okta, Stormpath started early with OAuth 2.0 SaaS and adopted the ROPC
grant too. They since deprecated these old flows but still have them active as
existing apps rely on them.
