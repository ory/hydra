# OAuth 2.0 Token Exchange (RFC 8693)

Hydra supports [RFC 8693 OAuth 2.0 Token Exchange](https://www.rfc-editor.org/rfc/rfc8693.html), which allows a client to exchange a security token (subject token) for a new access token, optionally for a different audience or with delegation (actor token).

## When to use

- A resource server receives an access token and needs a new token to call a backend service (e.g. with a different audience or scopes).
- Impersonation or delegation scenarios where an actor token is provided along with the subject token.

## Request

`POST /oauth2/token` with `application/x-www-form-urlencoded`:

| Parameter | Required | Description |
|-----------|----------|-------------|
| grant_type | Yes | `urn:ietf:params:oauth:grant-type:token-exchange` |
| subject_token | Yes | The token to exchange (e.g. an access token or JWT). |
| subject_token_type | Yes | Type of subject_token, e.g. `urn:ietf:params:oauth:token-type:access_token` or `urn:ietf:params:oauth:token-type:jwt`. |
| resource | No | URI of the target service where the new token will be used. |
| audience | No | Logical name of the target service (space-delimited or repeated). |
| scope | No | Requested scope for the issued token. |
| requested_token_type | No | Desired type of the issued token. |
| actor_token | No | For delegation: token representing the acting party. |
| actor_token_type | No | Required when actor_token is present. |

Client authentication (e.g. HTTP Basic with client_id and client_secret) is required unless configured otherwise. The client must have `urn:ietf:params:oauth:grant-type:token-exchange` in its `grant_types`.

## Response

On success (200 OK), the response includes:

| Field | Description |
|-------|-------------|
| access_token | The issued access token. |
| issued_token_type | **Required** by RFC 8693; e.g. `urn:ietf:params:oauth:token-type:access_token`. |
| token_type | Usually `Bearer`. |
| expires_in | Lifetime of the access token in seconds. |
| scope | Optional; scope of the issued token. |

## Supported subject token types

- **urn:ietf:params:oauth:token-type:access_token**: Opaque or JWT access token issued by this Hydra. Validated via introspection (storage lookup).
- **urn:ietf:params:oauth:token-type:jwt**: JWT; when the JWT is an access token issued by this server, it is validated the same way as access_token type.

## Example

```http
POST /oauth2/token HTTP/1.1
Host: hydra.example.com
Content-Type: application/x-www-form-urlencoded
Authorization: Basic <client_credentials>

grant_type=urn%3Aietf%3Aparams%3Aoauth%3Agrant-type%3Atoken-exchange
&subject_token=<access_token>
&subject_token_type=urn%3Aietf%3Aparams%3Aoauth%3Atoken-type%3Aaccess_token
&audience=https%3A%2F%2Fbackend.example.com%2Fapi
```

Successful response:

```json
{
  "access_token": "eyJ...",
  "issued_token_type": "urn:ietf:params:oauth:token-type:access_token",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

## Configuration

- **Grant type**: Add `urn:ietf:params:oauth:grant-type:token-exchange` to the client's `grant_types` (e.g. via Admin API or OIDC Dynamic Client Registration).
- **Skip client auth** (not recommended): Set `GrantTypeTokenExchangeCanSkipClientAuth: true` in Fosite config to allow unauthenticated token exchange requests.
- **OIDC Discovery**: `/.well-known/openid-configuration` includes `grant_types_supported` with token exchange when supported.

## References

- [RFC 8693 - OAuth 2.0 Token Exchange](https://www.rfc-editor.org/rfc/rfc8693.html)
