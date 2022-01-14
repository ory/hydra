---
id: advanced
title: Advanced OAuth2 and OpenID Connect Flows
---

## OAuth 2.0

### Audience

There are two types of audience concept in the context of OAuth 2.0 and OpenID
Connect:

1. OAuth 2.0: Access and Refresh Tokens are "internal-facing". The `aud` claim
   of an OAuth 2.0 Access and Refresh token defines at which _endpoints_ the
   token can be used.
2. OpenID Connect: The ID Token is "external-facing". The `aud` claim of an
   OpenID Connect ID Token defines which _clients_ should accept it.

While modifying the audience of an ID Token is not desirable, specifying the
audience of an OAuth 2.0 Access Token is. This is not defined as an IETF
Standard but is considered good practice in certain environments.

For this reason, Hydra allows you to control the aud claim of the access token.
To do so, you must specify the intended audiences in the OAuth 2.0 Client's
metadata on a per-client basis:

```
{
    "client_id": "...",
    "audience": ["https://api.my-cloud.com/user", "https://some-tenant.my-cloud.com/"]
}
```

The audience is a list of case-sensitive URLs. **URLs must not contain
whitespaces**.

#### OAuth 2.0 Authorization Code, Implicit, Hybrid Flows

When performing an OAuth 2.0 authorize code, implicit, or hybrid flow, you can
request audiences at the `/oauth2/auth` endpoint
`https://my-hydra.com/oauth2/auth?client_id=...&scope=...&audience=https%3A%2F%2Fapi.my-cloud.com%2Fuser+https%3A%2F%2Fsome-tenant.my-cloud.com%2F`
which requests audiences `https://api.my-cloud.com/user` and
`https://some-tenant.my-cloud.com/`.

The `audience` query parameter may contain multiple strings separated by a
url-encoded space (`+` or `%20`). The audience values themselves must also be
url encoded. The values will be validated against the whitelisted audiences
defined in the OAuth 2.0 Client:

- An OAuth 2.0 Client with the allowed audience `https://api.my-cloud/user` is
  allowed to request audience values `https://api.my-cloud/user`
  `https://api.my-cloud/user/1234` but not `https://api.my-cloud/not-user` nor
  `https://something-else/`.

The requested audience from the query parameter is then part of the login and
consent request payload as field `requested_access_token_audience`. You can then
alter the audience using `grant_audience.access_token` when accepting the
consent request:

```
hydra.acceptConsentRequest(challenge, {
  // Ory Hydra checks if requested audiences are allowed by the client, so we can simply echo this.
  grant_audience: {
    access_token: response.requested_access_token_audience,
    // or, for example:
    // access_token: ["https://api.my-cloud/not-user"]
  },

  // ... remember: false
  // ...
})
```

When introspecting the OAuth 2.0 Access Token, the response payload will include
the audience:

```
{
  "active": true,
  // ...
  "audience": ["https://api.my-cloud/user", "https://api.my-cloud/user/1234"]
}
```

#### OAuth 2.0 Client Credentials Grant

When performing the client credentials grant, the audience parameter from the
POST body of the `/oauth2/token` is decoded and validated according to the same
rules of the previous section, except for the login and consent part which does
not exist for this flow.

### JSON Web Tokens

Ory Hydra issues opaque OAuth 2.0 Access Tokens per default for the following
reasons:

1. **OAuth 2.0 Access Tokens represent internal state but are public
   knowledge:** An Access Token often contains internal data (e.g. session data)
   or other sensitive data (e.g. user roles and permissions) and is sometimes
   used as a means of transporting system-relevant information in a stateless
   manner. Therefore, making these tokens transparent (by using JSON Web Tokens
   as Access Tokens) comes with risk of exposing this information accidentally,
   and with the downside of not storing this information in the OAuth 2.0 Access
   Token at all.
2. **JSON Web Tokens can not hold secrets:** Unless encrypted, JSON Web Tokens
   can be read by everyone, including 3rd Parties. Therefore, they can not keep
   secrets. This point is similar to (1), but it is important to stress this.
3. **Access Tokens as JSON Web Tokens can not be revoked:** Well, you can revoke
   them, but they will be considered valid until the "expiry" of the token is
   reached. Unless, of course, you have a blacklist or check with Hydra if the
   token was revoked, which however defeats the purpose of using JSON Web Tokens
   in the first place.
4. **Certain OpenID Connect features will not work** when using JSON Web Tokens
   as Access Tokens, such as the pairwise subject identifier algorithm.
5. **There is a better solution: Use
   [Ory Oathkeeper](https://github.com/ory/oathkeeper)!** Ory Oathkeeper is a
   proxy you deploy in front of your services. It will "convert" Ory Hydra's
   opaque Access Tokens into JSON Web Tokens for your backend services. This
   allows your services to work without additional REST Calls while solving all
   previous points. **We really recommend this option if you want JWTs!**

If you are not convinced that Ory Oathkeeper is the right tool for the job, you
can still enable JSON Web Tokens in Ory Hydra by setting:

```yaml
strategies:
  access_token: jwt
```

Be aware that only access tokens are formatted as JSON Web Tokens. Refresh
tokens are not impacted by this strategy. By performing OAuth 2.0 Token
Introspection you can check if the token is still valid. If a token is revoked
or otherwise blacklisted, the OAuth 2.0 Token Introspection will return
`{ "active": false }`. This is useful when you do not want to rely only on the
token's expiry.

#### JSON Web Token Validation

You can validate JSON Web Tokens issued by Ory Hydra by pointing your `jwt`
library (e.g. [node-jwks-rsa](https://github.com/auth0/node-jwks-rsa)) to
`http://ory-hydra-public-api/.well-known/jwks.json`. All necessary keys are
available there.

#### Adding custom claims top-level to the Access Token

Assume you want to add custom claims to the access token with the following
code:

```typescript
let session: ConsentRequestSession = {
  access_token: {
    foo: 'bar'
  }
}
```

Then part of the resulting access token will look like this:

```json
{
  "ext": {
    "foo": "bar"
  }
}
```

If you instead want "foo" to be added top-level in the access token, you need to
set the configuration flag `oauth2.allowed_top_level_claims` like described in
[the reference Configuration](https://www.ory.sh/hydra/docs/reference/configuration).

Note: Any user defined allowed top level claim may not override standardized
access token claim names.

Configuring Hydra to allow "foo" as a top-level claim will result in the
following access token part (allowed claims get mirrored, for backwards
compatibility):

```json
{
  "foo": "bar",
  "ext": {
    "foo": "bar"
  }
}
```

#### Updating claims at token refresh

Hydra can be configured to retrieve updated token claims from an endpoint at
token refresh, which provides updated claims for a given subject and scopes.
This is similar to accepting consent request, where the application provides the
session data by calling Hydra Admin API.

:::note

This endpoint is called _before_ any logic in Ory Hydra is executed. If the
hook, for example, returns an error, the refresh token will remain unused!

:::note

You can configure `oauth2.refresh_token_hook` config key:

```yaml
oauth2:
  refresh_token_hook: https://my-example.app/token-refresh-hook
```

Hydra makes a `POST` request to this hook with the following payload:

```json
{
  "subject": "foo",
  "client_id": "bar",
  "granted_scopes": ["openid", "offline"],
  "granted_audience": []
}
```

Hook has to respond with `200 OK` and updated session data (i.e. "extra" claims)
for a token refresh to continue:

```json
{
  "session": {
    "access_token": {
      "foo": "bar"
    },
    "id_token": {
      "bar": "baz"
    }
  }
}
```

This will overwrite existing session data from the original consent request.

Hydra will gracefully deny refresh requests if the hook responds with
`403 Forbidden`. Any other response from the hook will fail refresh requests.

### OAuth 2.0 Client Authentication with private/public keypairs

Please head over to the
[RFC7523 Documentation](guides/oauth2-grant-type-jwt-bearer.mdx).

## OpenID Connect

### Subject Identifier Algorithms

Hydra supports two
[Subject Identifier Algorithms](http://openid.net/specs/openid-connect-core-1_0.html#SubjectIDTypes):

- `public`: This provides the same `sub` (subject) value to all Clients
  (default).
- `pairwise`: This provides a different `sub` value to each Client, so as not to
  enable Clients to correlate the End-User's activities without permission.

You can enable either one or both algorithms using the following configuration
layout:

```yaml
oidc:
  subject_identifiers:
    supported_types:
      - public
      - pairwise
```

When `pairwise` is enabled, you must also set
`oidc.subject_identifiers.pairwise.salt`. The salt is used to obfuscate the
`sub` value:

```yaml
oidc:
  subject_identifiers:
    supported_types:
      - public
      - pairwise
    pairwise:
      salt: some-salt
```

**This value should not be changed once set in production. Changing it will
cause all client applications to receive new user IDs from Ory Hydra which will
lead to serious complications with authentication on their side!**

Each OAuth 2.0 Client has a configuration field `subject_type`. The value of
that `subject_type` is either `public` or `pairwise`. If the identifier
algorithm is enabled, Ory Hydra will choose the right strategy automatically.

While Ory Hydra handles `sub` obfuscation out of the box, you may also override
this value with your own obfuscated `sub` value by setting
`force_subject_identifier` when accepting the login challenge in your user login
app.

### Using login_hint with Different Subject

When a user already logged in with a subject(e.g. user-A), and she would like to
login as another user using login_hint(e.g. login_hint=user-B), directly
accepting the latter login request in your login provider will make hydra reply:
`Subject from payload does not match subject from previous authentication`

The suggested flow is:

Check the response from
[GET login request](reference/api.mdx#get-a-login-request), if both the
`subject` and `login_hint` are NOT empty and also NOT the same user, redirect
UserAgent to `request_url` which is appended with '?prompt=login'. This will
make hydra ignore the existing authentication, and allow your login provider to
login a different subject.

For more information on `prompt=login` and other options, please check
[Authentication Request](https://openid.net/specs/openid-connect-core-1_0.html#AuthRequest).
