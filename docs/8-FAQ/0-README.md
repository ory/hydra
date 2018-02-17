# FAQ

This file keeps track of questions and discussions from Gitter and general help with various issues.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [How can I control SQL connection limits?](#how-can-i-control-sql-connection-limits)
- [Why is the Resource Owner Password Credentials grant not supported?](#why-is-the-resource-owner-password-credentials-grant-not-supported)
- [Should I use OAuth2 tokens for authentication?](#should-i-use-oauth2-tokens-for-authentication)
- [How to deal with mobile apps?](#how-to-deal-with-mobile-apps)
- [How should I run migrations?](#how-should-i-run-migrations)
- [What does the installation process look like?](#what-does-the-installation-process-look-like)
- [What does a migration process look like?](#what-does-a-migration-process-look-like)
- [How can I do this in docker?](#how-can-i-do-this-in-docker)
- [Can I set the log level to warn, error, debug, ...?](#can-i-set-the-log-level-to-warn-error-debug-)
- [How can I import TLS certificates?](#how-can-i-import-tls-certificates)
- [Is there an HTTP API Documentation?](#is-there-an-http-api-documentation)
- [How can I disable HTTPS for testing?](#how-can-i-disable-https-for-testing)
- [MySQL gives `unsupported Scan, storing driver.Value type []uint8 into type *time.Time`](#mysql-gives-unsupported-scan-storing-drivervalue-type-uint8-into-type-timetime)
- [The docker image exits immediately](#the-docker-image-exits-immediately)
- [Insufficient Entropy](#insufficient-entropy)
- [I get compile errors!](#i-get-compile-errors)
- [Is JWT supported?](#is-jwt-supported)
- [Refreshing tokens](#refreshing-tokens)
- [Revoking tokens & log out](#revoking-tokens-&-log-out)
- [Operational Considerations](#operational-considerations)
  - [Managing Client/Policy Definitions](#managing-clientpolicy-definitions)
- [Recovering root client access](#recovering-root-client-access)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Consent App Client

> Alexander Alimovs @aalimovs 13:25  
but this is not advised, and ideally you split consent app (identity provider) from the client app, correct?
I'm grasping OAuth 2.0 trying to get this work and it starts to make more and more sense, almost there

> Aeneas @arekkas 13:28  
yes so your consent app is like the source of truth, the globe that knows everything

> Aeneas @arekkas 13:28  
because it has access to the cryptographic keys that allow it to issue the consent response, which in turn can be used to say: hey, this guy here is peter and he's the freaking super admin
so hydra will be like: ok cool, welcome peter, here are your keys, don't do anything stupid with that
so if you use the consent app client everywhere, it might leak, or some one who is not peter may say: oh look, I can look up the cryptographic keys for the consent repsonse, let me use that and impersonate peter
so now some random internet guy is the superadmin of your system
thus, it is extremely important to only use the consent app client in the consent app, and nowhere else
in the error case above you are requesting the scope "hydra" which has not been granted to the client (see the scopes column)

## How can I control SQL connection limits?

You can configure SQL connection limits by appending parameters `max_conns`, `max_idle_conns`, or `max_conn_lifetime`
to the DSN: `postgres://foo:bar@host:port/database?max_conns=12`.

## Why is the Resource Owner Password Credentials grant not supported?

The following is a copy of the original [comment on GitHub](https://github.com/ory/hydra/pull/297#issuecomment-294282671):

I took a long time for this issue, primarily because I felt very uncomfortable implementing it. The ROCP grant is something from the "dark ages" of OAuth2 and there are suitable replacements for mobile clients, such as public oauth2 clients, which are supported by Hydra: https://tools.ietf.org/html/draft-ietf-oauth-native-apps-09

The OAuth2 Thread Model explicitly states that the ROPC grant is commonly used in legacy/migration scenarios, and

>   This grant type has higher
   risk because it maintains the UID/password anti-pattern.
   Additionally, because the user does not have control over the
   authorization process, clients using this grant type are not limited   by scope but instead have potentially the same capabilities as the
   user themselves.  As there is no authorization step, the ability to
   offer token revocation is bypassed.

> Because passwords are often used for more than 1 service, this
   anti-pattern may also put at risk whatever else is accessible with
   the supplied credential.  Additionally, any easily derived equivalent
   (e.g., joe@example.com and joe@example.net) might easily allow
   someone to guess that the same password can be used elsewhere.

>    Impact: The resource server can only differentiate scope based on the
   access token being associated with a particular client.  The client
   could also acquire long-lived tokens and pass them up to an
   attacker's web service for further abuse.  The client, eavesdroppers,
   or endpoints could eavesdrop the user id and password.

>    o  Except for migration reasons, minimize use of this grant type.

- [source](https://tools.ietf.org/html/rfc6819#section-4.4.3)

Thus, I decided to not implement the ROPC grant in Hydra. Over time, I will add documentation how to deal with mobile scenarios and similar.

## Should I use OAuth2 tokens for authentication?

OAuth2 tokens are like money. It allows you to buy stuff, but the cashier does not really care if the money is
yours or if you stole it, as long as it's valid money. Depending on what you understand as authentication, this is a yes and no answer:

* **Yes:** You can use access tokens to find out which user ("subject") is performing an action in a resource provider (blog article service, shopping basket, ...).
Coming back to the money example: *You*, the subject, receives a cappuccino from the vendor (resource provider) in exchange for money (access token).
* **No:** Never use access tokens for logging people in, for example `http://myapp.com/login?access_token=...`.
Coming back to the money example: The police officer ("authentication server") will not accept money ("access token") as a proof of identity ("it's really you"). Unless he is corrupt ("vulnerable"), of course.

In the second example ("authentication server"), you must use OpenID Connect ID Tokens.

## How to deal with mobile apps?

Authors of apps running on client side (native apps, single page apps, hybrid apps) have for some
time relied on the Resource Owner Password Credentials grant. This is highly discouraged by the IETF, and replaced
with recommendations in [OAuth 2.0 for Native Apps](https://tools.ietf.org/html/draft-ietf-oauth-native-apps-03).

To keep things short, it allows you to perform the normal `authorize_code` flows without supplying a password. Hydra
allows this by setting the public flag, for example:

```sh
hydra clients create \
    --id my-id \
    --is-public \
    -r code,id_token \
    -g authorization_code,refresh_token \
    -a offline,openid \
    -c https://mydomain/callback
```

## How should I run migrations?

Since ORY Hydra 0.8.0, migrations are no longer run automatically on boot. This is required in production environments,
because:

1. Although SQL migrations are tested, migrating schemas can cause data loss and should only be done consciously with
prior back ups.
2. Running a production system with a user that has right such as ALTER TABLE is a security anti-pattern.

Thus, to initialize the database schemas, it is required to run `hydra migrate sql driver://user:password@host:port/db` before running
`hydra host`.

## What does the installation process look like?

1. Run `hydra migrate sql ...` on a host close to the database (e.g. a virtual machine with access to the SQL instance).

## What does a migration process look like?

1. Make sure a database update is required by checking the release notes.
2. Make a back up of the database.
3. Run the migration script on a host close to the database (e.g. a virtual machine with access to the SQL instance).
Schemas are usually backwards compatible, so instances running previous versions of ORY Hydra should keep working fine.
If backwards compatibility is not given, this will be addressed in the patch notes.
4. Upgrade all ORY Hydra instances.

## How can I do this in docker?

Many deployments of ORY Hydra use Docker. Although several options are available, we advise to extend the ORY Hydra Docker
image

**Dockerfile**
```
FROM oryd/hydra:tag

ENTRYPOINT /go/bin/hydra migrate sql $DATABASE_URL
```

and run it in your infrastructure once.

Additionally, *but not recommended*, it is possible to override the entry point of the ORY Hydra Docker image using CLI flag
`--entrypoint "hydra migrate sql $DATABASE_URL; hydra host"` or with `entrypoint: hydra migrate sql $DATABASE_URL; hydra host`
set in your docker compose config.

## Can I set the log level to warn, error, debug, ...?

Yes, you can do so by setting the environment variable `LOG_LEVEL=<level>`. There are various levels supported:

* debug
* info
* warn
* error
* fatal
* panic

## How can I import TLS certificates?

You can import TLS certificates when running `hydra host`. This can be done by setting the following environment variables:

**Read from file**
- `HTTPS_TLS_CERT_PATH`: The path to the TLS certificate (pem encoded).
- `HTTPS_TLS_KEY_PATH`: The path to the TLS private key (pem encoded).

**Embedded**
- `HTTPS_TLS_CERT`: A pem encoded TLS certificate passed as string. Can be used instead of TLS_CERT_PATH.
- `HTTPS_TLS_KEY`: A pem encoded TLS key passed as string. Can be used instead of TLS_KEY_PATH.

Or by specifying the following flags:

```
--https-tls-cert-path string   Path to the certificate file for HTTP/2 over TLS (https). You can set HTTPS_TLS_KEY_PATH or HTTPS_TLS_KEY instead.
--https-tls-key-path string    Path to the key file for HTTP/2 over TLS (https). You can set HTTPS_TLS_KEY_PATH or HTTPS_TLS_KEY instead.
```

## Is there an HTTP API Documentation?

Yes, it is available at [Apiary](http://docs.hydra13.apiary.io/).

## How can I disable HTTPS for testing?

You can do so by running `hydra host --dangerous-force-http`.

## MySQL gives `unsupported Scan, storing driver.Value type []uint8 into type *time.Time`

> did a quick test to get mysql running, but run into migrate sql issue - seems mysql related
An error occurred while running the migrations: Could not apply ladon SQL migrations: Could not migrate sql schema, applied 0 migrations: sql: Scan error on column index 0: unsupported Scan, storing driver.Value type []uint8 into type *time.Time
is this a known bug ? or any specific mysql version which is required (running 5.7) ?

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

> Hey there , I am getting this error when I try request an access token "The request used a security parameter (e.g., anti-replay, anti-csrf) with insufficient entropy (minimum of 8 characters)"

> Kareem Diaa @kimooz Jun 07 16:41  
Hey there , I am getting this error when I try request an access token "The request used a security parameter (e.g., anti-replay, anti-csrf) with insufficient entropy (minimum of 8 characters)"

> Aeneas @arekkas Jun 07 16:41  
@kimooz make sure state and nonce are set in your auth code url (http://hydra/oauth2/auth?client_id=...&nonce=THIS_NEEDS_TO_BE_SET&state=THIS_ALSO_NEEDS_TO_BE_SET

## I get compile errors!

> I would try deleting the vendor dir and glide’s files and try glide init again or clear Glide’s global cache.

> follow the steps in the readme https://github.com/ory/hydra#building-from-source

## Is JWT supported?

> Mufid @mufid 03:29  
> Could Hydra's Access Token be a JWT? So that my resource server does not need to call Introspection API for each request.

> Mufid @mufid 03:39  
Yes, the access token looks like JWT, but i am unable to decode it. Here is my example token form Hydra: LpxuGoqWy7lYp9N0Cea8mEGR6IHhyr37jxZXRHqSjRM.nU-jMnAJ7dUKQPjWF4QBEL9OQWVU8zj_ElhrT-FQrWw (JWT Tokens should have 2 dots (3 segments), so this is not a valid JWT)

> Mufid @mufid 03:56  
*form --> from, typo, sorry.
> Aeneas @arekkas 11:50  
@mufid JWT is not supported at the moment, we might add it, but not as part of the hydra community edition

## Refreshing tokens

> Kareem Diaa @kimooz 15:48  
One last question  if you don't mind
from your experience do you think that saving the user access token in a session and validating it from the client on ever refresh does that make sense or not?
using the introspect endpoint

> Aeneas @arekkas 15:51  
nah, simply write your http calls in a way that if a 401 or 403 occurrs, the token is refreshed
that's the easiest
and cleanest

## Revoking tokens & log out

> Kareem Diaa @kimooz 15:41  
Thanks @arekkas. I had two other questions:  
1\. Is there a way to revoke all access tokens for a certain user("log out user")  
2\. How can I inform the consent app that this user logged out?  

> Aeneas @arekkas 15:42  
1\. no this isn't supported currently  
2\. you can't because log out and revoking access tokens are two things  
and it would require an additional api or something, which makes the consent app harder to write and integrate

> Kareem Diaa @kimooz 15:43  
So can you suggest a workaround?
I want implement single sign off

> Aeneas @arekkas 15:44  
the user has the access and refresh token right
in his browser or somewhere

> Kareem Diaa @kimooz 15:44  
yah

> Aeneas @arekkas 15:44  
ok so why not make a request to /oauth2/revoke
and pass that refresh token
(you will probably need a proxy with a client id and secret for that to be possible, but you get the point)

> Kareem Diaa @kimooz 15:46  
yah but the moment he refreshes, the client will hit on hydra and then consent where it will find that this user is already logged in
and will return a new token although he should have logged out
ohh so you mean have two requests one for hydra to revoke and one for consent to log out correct?

> Aeneas @arekkas 15:47  
yes

## Operational Considerations

**This section might be outdated.**

This section is intended to give operations folk some useful guidelines on
how to best manage a hydra deployment.

### Managing Client/Policy Definitions

It is useful for JSON files for client and policy definitions to be persisted outside
of the hydra database itself. There are several reasons for this:

- client secrets are stored in the database as a bcrypt hash, so you will not be
  able to read back the secret.  If, at some later point, you need to update the
  client definition then you need to delete the client definition from the database
  and recreate it. So as to ensure that existing web/mobile/other apps are able to
  continue operation, you will need to recreate the same client ID/secret - for which
  you should use the `hydra clients import` command.

A good storage platform for these files is Hashicorp [Vault](https://www.vaultproject.io). The full setup of a
Vault deployment is beyond the scope of this document, please refer to the Vault website
for this.

With LDAP-based authentication configured, reading a client definition from
the Vault "secrets" backend is as simple as:

````
$ vault auth -method=ldap username=john.doe
$ vault read -field=myclient secrets/hydra/clients | hydra clients import /dev/stdin
````

The first command above will prompt you for your LDAP password and then create
a file `~/.vault-token`. Note that you may wish to consider using token authentication
instead of LDAP, but the above should be simple for ops folk to perform.

If you need to update a client later, you can delete the client from hydra using:

````
$ hydra clients delete "<clientid>"
````
and then re-import as above.

The same procedure can be followed for importing/updating policy definitions.

## Recovering root client access

If you somehow manage to lose admin access to your Hydra system, you can regain this
by making use of Hydra's temporary root client creation - which is triggered when
hydra is unable to find any client definitions upon startup. Due to the ID given to
policy used for temporary root clients, you may need to also delete configured
policies. To do so, make a (sql) back up of existing clients and policies,
then empty tables `hydra_clients` and `hydra_policies`, and:

- Restart Hydra
- Re-import your client/policy definitions, as described above
- Delete your new temporary root client
- Ensure that any Hydra clients which have read keys from hydra are refreshed, possibly
  involving a simple restart to effect a timely update
