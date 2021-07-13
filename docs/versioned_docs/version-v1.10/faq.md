---
id: faq
title: FAQ
---

This file keeps track of questions and discussions from Gitter and general help
with various issues.

**Please be aware that some things noted here might be highly outdated. If you
find an outdated section, please create a PR which removes or updates it!**

<!-- toc -->

## How can I store custom data ?

From [Github](https://github.com/ory/hydra/issues/2333)

> How can I store custom data in the client that can be retrieved along with the
> login request (using the challenge code)?
> I would like send in some metadata to the /oauth2/auth endpoint that can be
> retrievable by the "challenge code" that Hydra creates.

You can do this using the
[`request url`](https://github.com/ory/hydra/blob/27dc147a37c1cb98b7beb5428c389545958dd122/consent/types.go#L505).

Both Login and Consent Requests have a
[`request url`](https://github.com/ory/hydra/blob/27dc147a37c1cb98b7beb5428c389545958dd122/consent/types.go#L505)
which contains the original OAuth2 Authorize URL (e.g.
`/oauth2/auth?client_id=...&your_custom_param=....`)

```
RequestURL is the original OAuth 2.0 Authorization URL requested by the OAuth 2.0 client.
It is the URL which initiates the OAuth 2.0 Authorization Code or OAuth 2.0 Implicit flow.
This URL is typically not needed, but might come in handy if you want to deal with additional request parameters.
```

## How can I control SQL connection limits?

You can configure SQL connection limits by appending parameters `max_conns`,
`max_idle_conns`, `max_conn_lifetime` or `max_conn_idle_time` to the DSN:
`postgres://foo:bar@host:port/database?max_conns=12`.

## Why is the Resource Owner Password Credentials grant not supported?

The following is a copy of the original
[comment on GitHub](https://github.com/ory/hydra/pull/297#issuecomment-294282671):

I took a long time for this issue, primarily because I felt very uncomfortable
implementing it. The ROCP grant is something from the "dark ages" of OAuth2 and
there are suitable replacements for mobile clients, such as public oauth2
clients, which are supported by Hydra:
https://tools.ietf.org/html/draft-ietf-oauth-native-apps-09

The OAuth2 Thread Model explicitly states that the ROPC grant is commonly used
in legacy/migration scenarios, and

> This grant type has higher risk because it maintains the UID/password
> anti-pattern. Additionally, because the user does not have control over the
> authorization process, clients using this grant type are not limited by scope
> but instead have potentially the same capabilities as the user themselves. As
> there is no authorization step, the ability to offer token revocation is
> bypassed.

> Because passwords are often used for more than 1 service, this anti-pattern
> may also put at risk whatever else is accessible with the supplied credential.
> Additionally, any easily derived equivalent (e.g., joe@example.com and
> joe@example.net) might easily allow someone to guess that the same password
> can be used elsewhere.

> Impact: The resource server can only differentiate scope based on the access
> token being associated with a particular client. The client could also acquire
> long-lived tokens and pass them up to an attacker's web service for further
> abuse. The client, eavesdroppers, or endpoints could eavesdrop the user id and
> password.

> o Except for migration reasons, minimize use of this grant type.

- [source](https://tools.ietf.org/html/rfc6819#section-4.4.3)

Thus, I decided to not implement the ROPC grant in Hydra. Over time, I will add
documentation how to deal with mobile scenarios and similar.

## Should I use OAuth2 tokens for authentication?

OAuth2 tokens are like money. It allows you to buy stuff, but the cashier does
not really care if the money is yours or if you stole it, as long as it's valid
money. Depending on what you understand as authentication, this is a yes and no
answer:

- **Yes:** You can use access tokens to find out which user ("subject") is
  performing an action in a resource provider (blog article service, shopping
  basket, ...). Coming back to the money example: _You_, the subject, receives a
  cappuccino from the vendor (resource provider) in exchange for money (access
  token).
- **No:** Never use access tokens for logging people in, for example
  `http://myapp.com/login?access_token=...`. Coming back to the money example:
  The police officer ("authentication server") will not accept money ("access
  token") as a proof of identity ("it's really you"). Unless he is corrupt
  ("vulnerable"), of course.

In the second example ("authentication server"), you must use OpenID Connect ID
Tokens.

## How to deal with mobile apps?

Read [this article](https://www.ory.sh/oauth2-for-mobile-app-spa-browser).

## How should I run migrations?

Since ORY Hydra 0.8.0, migrations are no longer run automatically on boot. This
is required in production environments, because:

1. Although SQL migrations are tested, migrating schemas can cause data loss and
   should only be done consciously with prior back ups.
2. Running a production system with a user that has right such as ALTER TABLE is
   a security anti-pattern.

Thus, to initialize the database schemas, it is required to run
`hydra migrate sql driver://user:password@host:port/db` before running
`hydra host`.

## What does the installation process look like?

1. Run `hydra migrate sql ...` on a host close to the database (e.g. a virtual
   machine with access to the SQL instance).

## What does a migration process look like?

1. Make sure a database update is required by checking the release notes.
2. Make a back up of the database.
3. Run the migration script on a host close to the database (e.g. a virtual
   machine with access to the SQL instance). Schemas are usually backwards
   compatible, so instances running previous versions of ORY Hydra should keep
   working fine. If backwards compatibility is not given, this will be addressed
   in the patch notes.
4. Upgrade all ORY Hydra instances.

## How can I do this in docker?

Many deployments of ORY Hydra use Docker. Although several options are
available, we advise to extend the ORY Hydra Docker image

**Dockerfile**

```
FROM oryd/hydra:<tag>

ENTRYPOINT hydra migrate sql --yes $DATABASE_URL
```

and run it in your infrastructure once.

Additionally, _but not recommended_, it is possible to override the entry point
of the ORY Hydra Docker image using CLI flag
`--entrypoint "hydra migrate sql --yes $DATABASE_URL; hydra host"` or with
`entrypoint: hydra migrate sql $DATABASE_URL; hydra host` set in your docker
compose config.

## Can I set the log level to warn, error, debug, ...?

Yes, you can do so by setting the environment variable `LOG_LEVEL=<level>`.
There are various levels supported:

- debug
- info
- warn
- error
- fatal
- panic

## How can I import TLS certificates?

You can import TLS certificates when running `hydra host`. This can be done by
setting the following environment variables:

**Read from file**

- `HTTPS_TLS_CERT_PATH`: The path to the TLS certificate (pem encoded).
- `HTTPS_TLS_KEY_PATH`: The path to the TLS private key (pem encoded).

**Embedded**

- `HTTPS_TLS_CERT`: A pem encoded TLS certificate passed as string. Can be used
  instead of TLS_CERT_PATH.
- `HTTPS_TLS_KEY`: A pem encoded TLS key passed as string. Can be used instead
  of TLS_KEY_PATH.

Or by specifying the following flags:

```
--https-tls-cert-path string   Path to the certificate file for HTTP/2 over TLS (https). You can set HTTPS_TLS_KEY_PATH or HTTPS_TLS_KEY instead.
--https-tls-key-path string    Path to the key file for HTTP/2 over TLS (https). You can set HTTPS_TLS_KEY_PATH or HTTPS_TLS_KEY instead.
```

## Is there an HTTP API Documentation?

[Yes](https://www.ory.sh/docs/hydra/reference/api.mdx). bro

## How can I disable HTTPS for testing?

You can do so by running `hydra host --dangerous-force-http`.

## MySQL gives `unsupported Scan, storing driver.Value type []uint8 into type *time.Time`

> did a quick test to get mysql running, but run into migrate sql issue - seems
> mysql related An error occurred while running the migrations: Could not apply
> ladon SQL migrations: Could not migrate sql schema, applied 0 migrations: sql:
> Scan error on column index 0: unsupported Scan, storing driver.Value type
> []uint8 into type \*time.Time is this a known bug ? or any specific mysql
> version which is required (running 5.7) ?

```
$ hydra help host
...
   - MySQL: If DATABASE_URL is a DSN starting with mysql:// MySQL will be used as storage backend.
        Example: DATABASE_URL=mysql://user:password@tcp(host:123)/database?parseTime=true

        Be aware that the ?parseTime=true parameter is mandatory, or timestamps will not work.
...
```

## The docker image exits immediately

Check the logs using `docker logs <container-id>`.

## Insufficient Entropy

> Hey there , I am getting this error when I try request an access token "The
> request used a security parameter (e.g., anti-replay, anti-csrf) with
> insufficient entropy (minimum of 8 characters)"

> Kareem Diaa @kimooz Jun 07 16:41 Hey there , I am getting this error when I
> try request an access token "The request used a security parameter (e.g.,
> anti-replay, anti-csrf) with insufficient entropy (minimum of 8 characters)"

> Aeneas @arekkas Jun 07 16:41 @kimooz make sure state and nonce are set in your
> auth code url
> (http://hydra/oauth2/auth?client_id=...&nonce=THIS_NEEDS_TO_BE_SET&state=THIS_ALSO_NEEDS_TO_BE_SET

## I get compile errors!

> I would try deleting the vendor dir and glide’s files and try glide init again
> or clear Glide’s global cache.

> follow the steps in the readme
> https://github.com/ory/hydra#building-from-source

## Refreshing tokens

> Kareem Diaa @kimooz 15:48 One last question if you don't mind from your
> experience do you think that saving the user access token in a session and
> validating it from the client on ever refresh does that make sense or not?
> using the introspect endpoint

> Aeneas @arekkas 15:51 nah, simply write your http calls in a way that if a 401
> or 403 occurs, the token is refreshed that's the easiest and cleanest

## Hydra won't remember my login

Hydra only keeps track of the session if you set the `remember_me` parameter in
the `accept login request`. Otherwise the login challenge will always instruct
you to show the login UI.

## Where can I get documentation on running multiple instances of Hydra?

> @killa#7335 20190103 01:06 My company is thinking about adopting Hydra. There
> is a performance benchmark including data at
> https://www.ory.sh/docs/next/performance/hydra and some information on how to
> scale hydra. Where can I get documentation on running multiple instances of
> Hydra? Source: Discord/ory/general

Hydra scales according to 12 factor principles. Just add another instance with
the same config. Please check the documentation section for 12 factor principles
for more information: https://www.ory.sh/docs/ecosystem/cloud-native. There is
also some information on collecting statistics in the section on
[prometheus](https://github.com/prometheus) in the
[five minute tutorial](./5min-tutorial.mdx).

## Is it possible to disable/enable certain flows in Hydra?

> Is it possible to leave only OpenID authorization code grant flow, while
> disabling all the rest of the flows in Hydra? Is it configurable that way?

Yes - you can configure that in your client. It has things like grant_types
etc - there you can basically whitelist the flows you need.

## How can i test if my 4445 is running properly?

> I am using ory hydra for authentication.
> I get a 404 when im trying to create a client on my private EC2 hydra task.
> Is there a way to test if my 4445 is running properly?

You can check /health/alive, to see if it's alive. and /health/ready, to see if
it's also in ready state (meaning db connectivity works).
