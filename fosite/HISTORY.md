**THIS DOCUMENT HAS MOVED**

This file is no longer being updated and kept for historical reasons. Please
check the [CHANGELOG](CHANGELOG.md) instead!

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [0.28.0](#0280)
- [0.27.0](#0270)
  - [Conceptual Changes](#conceptual-changes)
  - [API Changes](#api-changes)
- [0.26.0](#0260)
- [0.24.0](#0240)
  - [Breaking change(s)](#breaking-changes)
    - [`fosite/handler/oauth2.JWTStrategy`](#fositehandleroauth2jwtstrategy)
    - [`OpenIDConnectRequestValidator.ValidatePrompt`](#openidconnectrequestvalidatorvalidateprompt)
- [0.23.0](#0230)
  - [Breaking change(s)](#breaking-changes-1)
    - [`Hasher`](#hasher)
- [0.22.0](#0220)
  - [Breaking change(s)](#breaking-changes-2)
    - [`JWTStrategy`](#jwtstrategy)
- [0.21.0](#0210)
  - [Changes to parsing of OAuth 2.0 Client `response_types`](#changes-to-parsing-of-oauth-20-client-response_types)
  - [`openid.DefaultStrategy` field name changed](#openiddefaultstrategy-field-name-changed)
  - [`oauth2.RS256JWTStrategy` was renamed and field name changed](#oauth2rs256jwtstrategy-was-renamed-and-field-name-changed)
  - [Adds `private_key_jwt` client authentication method](#adds-private_key_jwt-client-authentication-method)
  - [Response Type `id_token` no longer required for authorize_code flow](#response-type-id_token-no-longer-required-for-authorize_code-flow)
- [0.20.0](#0200)
- [Breaking Changes](#breaking-changes)
  - [JWT Claims](#jwt-claims)
  - [`AuthorizeCodeStorage`](#authorizecodestorage)
- [0.19.0](#0190)
- [0.18.0](#0180)
- [0.17.0](#0170)
- [0.16.0](#0160)
- [0.15.0](#0150)
- [0.14.0](#0140)
- [0.13.0](#0130)
  - [Breaking changes](#breaking-changes)
- [0.12.0](#0120)
  - [Breaking changes](#breaking-changes-1)
    - [Improved cryptographic methods](#improved-cryptographic-methods)
- [0.11.0](#0110)
  - [Non-breaking changes](#non-breaking-changes)
    - [Storage adapter](#storage-adapter)
    - [Reducing use of gomock](#reducing-use-of-gomock)
  - [Breaking Changes](#breaking-changes-1)
    - [`fosite/handler/oauth2.AuthorizeCodeGrantStorage` was removed](#fositehandleroauth2authorizecodegrantstorage-was-removed)
    - [`fosite/handler/oauth2.RefreshTokenGrantStorage` was removed](#fositehandleroauth2refreshtokengrantstorage-was-removed)
    - [`fosite/handler/oauth2.AuthorizeCodeGrantStorage` was removed](#fositehandleroauth2authorizecodegrantstorage-was-removed-1)
    - [WildcardScopeStrategy](#wildcardscopestrategy)
    - [Refresh tokens and authorize codes are no longer JWTs](#refresh-tokens-and-authorize-codes-are-no-longer-jwts)
    - [Delete access tokens when persisting refresh session](#delete-access-tokens-when-persisting-refresh-session)
- [0.10.0](#0100)
- [0.9.0](#090)
- [0.8.0](#080)
  - [Breaking changes](#breaking-changes-2)
    - [`ClientManager`](#clientmanager)
    - [`OAuth2Provider`](#oauth2provider)
- [0.7.0](#070)
- [0.6.0](#060)
- [0.5.0](#050)
- [0.4.0](#040)
- [0.3.0](#030)
- [0.2.0](#020)
- [0.1.0](#010)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 0.28.0

This version (re-)introduces refresh token lifespans. Per default, this feature
is enabled and set to 30 days. If a refresh token has not been used within 30
days, it will expire.

To disable refresh token lifespans (previous behaviour), set
`compose.Config.RefreshTokenLifespan = -1`.

## 0.27.0

This PR adds the ability to specify a target audience for OAuth 2.0 Access
Tokens.

### Conceptual Changes

From now on, `scope` and `audience` will be checked against the client's
whitelisted scope and audience on every refresh token exchange. This prevents
clients, which no longer are allowed to request a certain audience or scope, to
keep using those values with existing refresh tokens.

### API Changes

```go
type fosite.Client interface {
+	// GetAudience returns the allowed audience(s) for this client.
+	GetAudience() Arguments
}
```

```go
type fosite.Request struct {
-   Scopes         Argument
+   RequestedScope Argument

-   GrantedScopes  Argument
+   GrantedScope   Argument
}
```

```go
type fosite.Requester interface {
+	// GetRequestedAudience returns the requested audiences for this request.
+	GetRequestedAudience() (audience Arguments)

+	// SetRequestedAudience sets the requested audienc.
+	SetRequestedAudience(audience Arguments)

+	// GetGrantedAudience returns all granted scopes.
+	GetGrantedAudience() (grantedAudience Arguments)

+	// GrantAudience marks a request's audience as granted.
+	GrantAudience(audience string)
}
```

```go
type fosite/token/jwt.JWTClaimsContainer interface {
-	// With returns a copy of itself with expiresAt and scope set to the given values.
-	With(expiry time.Time, scope, audience []string) JWTClaimsContainer

+	// With returns a copy of itself with expiresAt, scope, audience set to the given values.
+	With(expiry time.Time, scope, audience []string) JWTClaimsContainer
}
```

## 0.26.0

This release makes it easier to define custom JWT Containers for access tokens
when using the JWT strategy. To do that, the following signatures have changed:

```go
// github.com/ory/fosite/handler/oauth2
type JWTSessionContainer interface {
	// GetJWTClaims returns the claims.
-	GetJWTClaims() *jwt.JWTClaims
+	GetJWTClaims() jwt.JWTClaimsContainer

	// GetJWTHeader returns the header.
	GetJWTHeader() *jwt.Headers

	fosite.Session
}
```

```go
+ type JWTClaimsContainer interface {
+	// With returns a copy of itself with expiresAt and scope set to the given values.
+	With(expiry time.Time, scope []string) JWTClaimsContainer
+
+	// WithDefaults returns a copy of itself with issuedAt and issuer set to the given default values. If those
+	// values are already set in the claims, they will not be updated.
+	WithDefaults(iat time.Time, issuer string) JWTClaimsContainer
+
+	// ToMapClaims returns the claims as a github.com/dgrijalva/jwt-go.MapClaims type.
+	ToMapClaims() jwt.MapClaims
+ }
```

All default session implementations have been updated to reflect this change. If
you define custom session, this patch will affect you.

## 0.24.0

This release addresses areas where the go context was missing or not propagated
down the call path properly.

### Breaking change(s)

#### `fosite/handler/oauth2.JWTStrategy`

The
[`fosite/handler/oauth2.JWTStrategy`](https://github.com/ory/fosite/blob/master/handler/oauth2/strategy.go)
interface changed as a context parameter was added to its method signature:

```go
type JWTStrategy interface {
-	Validate(tokenType fosite.TokenType, token string) (requester fosite.Requester, err error)
+	Validate(ctx context.Context, tokenType fosite.TokenType, token string) (requester fosite.Requester, err error)
}
```

#### `OpenIDConnectRequestValidator.ValidatePrompt`

The
[`OpenIDConnectRequestValidator.ValidatePrompt`](https://github.com/ory/fosite/blob/master/handler/openid/validator.go)
method signature was updated to take a go context as its first parameter:

```go
-	func (v *OpenIDConnectRequestValidator) ValidatePrompt(req fosite.AuthorizeRequester) error {
+	func (v *OpenIDConnectRequestValidator) ValidatePrompt(ctx context.Context, req fosite.AuthorizeRequester) error {
```

## 0.23.0

This releases addresses inconsistencies in some of the public interfaces by
passing in the go context to their signatures.

### Breaking change(s)

#### `Hasher`

The [`Hasher`](https://github.com/ory/fosite/blob/master/hash.go) interface
changed as a context parameter was added to its method signatures:

```go
type Hasher interface {
-	Compare(hash, data []byte) error
+	Compare(ctx context.Context, hash, data []byte) error
-	Hash(data []byte) ([]byte, error)
+	Hash(ctx context.Context, data []byte) ([]byte, error)
}
```

## 0.22.0

This releases addresses inconsistencies in some of the public interfaces by
passing in the go context to their signatures.

### Breaking change(s)

#### `JWTStrategy`

The [`JWTStrategy`](https://github.com/ory/fosite/blob/master/token/jwt/jwt.go)
interface changed as a context parameter was added to its method signatures:

```go
type JWTStrategy interface {
-	Generate(claims jwt.Claims, header Mapper) (string, string, error)
+	Generate(ctx context.Context, claims jwt.Claims, header Mapper) (string, string, error)
-	Validate(token string) (string, error)
+	Validate(ctx context.Context, token string) (string, error)
-	GetSignature(token string) (string, error)
+	GetSignature(ctx context.Context, token string) (string, error)
-	Hash(in []byte) ([]byte, error)
+	Hash(ctx context.Context, in []byte) ([]byte, error)
-	Decode(token string) (*jwt.Token, error)
+	Decode(ctx context.Context, token string) (*jwt.Token, error)
	GetSigningMethodLength() int
}
```

## 0.21.0

This release improves compatibility with the OpenID Connect Dynamic Client
Registration 1.0 specification.

### Changes to parsing of OAuth 2.0 Client `response_types`

Previously, when response types such as `code token id_token` were requested
(OpenID Connect Hybrid Flow) it was enough for the client to have
`response_types=["code", "token", "id_token"]`. This is however incompatible
with the OpenID Connect Dynamic Client Registration 1.0 spec which dictates that
the `response_types` have to match exactly.

Assuming you are requesting `&response_types=code+token+id_token`, your client
should have `response_types=["code token id_token"]`, if other response types
are required (e.g. `&response_types=code`, `&response_types=token`) they too
must be included: `response_types=["code", "token", "code token id_token"]`.

### `openid.DefaultStrategy` field name changed

Field `RS256JWTStrategy` was renamed to `JWTStrategy` and now relies on an
interface instead of a concrete struct.

### `oauth2.RS256JWTStrategy` was renamed and field name changed

The strategy `oauth2.RS256JWTStrategy` was renamed to
`oauth2.DefaultJWTStrategy` and now accepts an interface that implements
`jwt.JWTStrategy` instead of directly relying on `jwt.RS256JWTStrategy`. For
this reason, the field `RS256JWTStrategy` was renamed to `JWTStrategy`

### Adds `private_key_jwt` client authentication method

This patch adds the ability to perform the
[`private_key_jwt` client authentication method](http://openid.net/specs/openid-connect-core-1_0.html#ClientAuthentication)
defined in the OpenID Connect specification. Please note that method
`client_secret_jwt` is not supported because of the BCrypt hashing strategy.

For this strategy to work, you must set the `TokenURL` field of the
`compose.Config` object to the authorization server's Token URL.

If you would like to support this authentication method, your `Client`
implementation must also implement `fosite.DefaultOpenIDConnectClient` and then,
for example, `GetTokenEndpointAuthMethod()` should return `private_key_jwt`.

### Response Type `id_token` no longer required for authorize_code flow

The `authorize_code`
[does not require](https://openid.net/specs/openid-connect-registration-1_0.html#ClientMetadata)
the `id_token` response type to be available when performing the OpenID Connect
flow:

> grant_types OPTIONAL. JSON array containing a list of the OAuth 2.0 Grant
> Types that the Client is declaring that it will restrict itself to using. The
> Grant Type values used by OpenID Connect are:
>
>          authorization_code: The Authorization Code Grant Type described in OAuth 2.0 Section 4.1.
>          implicit: The Implicit Grant Type described in OAuth 2.0 Section 4.2.
>          refresh_token: The Refresh Token Grant Type described in OAuth 2.0 Section 6.
>
>      The following table lists the correspondence between response_type values that the Client will use and grant_type values that MUST be included in the registered grant_types list:
>
>          code: authorization_code
>          id_token: implicit
>          token id_token: implicit
>          code id_token: authorization_code, implicit
>          code token: authorization_code, implicit
>          code token id_token: authorization_code, implicit
>
>      If omitted, the default is that the Client will use only the authorization_code Grant Type.

Before this patch, the `id_token` response type was required whenever an ID
Token was requested. This patch changes that.

## 0.20.0

This release implements an OAuth 2.0 Best Practice with regards to revoking
already issued access and refresh tokens if an authorization code is used more
than one time.

## Breaking Changes

### JWT Claims

- `github.com/ory/fosite/token/jwt.JWTClaims.Audience` is no longer a `string`,
  but a string slice `[]string`.
- `github.com/ory/fosite/handler/openid.IDTokenClaims` is no longer a `string`,
  but a string slice `[]string`.

### `AuthorizeCodeStorage`

This improves security as, in the event of an authorization code being leaked,
all associated tokens are revoked. To implement this feature, a breaking change
had to be introduced. The
`github.com/ory/fosite/handler/oauth2.AuthorizeCodeStorage` interface changed as
follows:

- `DeleteAuthorizeCodeSession(ctx context.Context, code string) (err error)` has
  been removed from the interface and is no longer used by this library.
- `InvalidateAuthorizeCodeSession(ctx context.Context, code string) (err error)`
  has been introduced.
- The error `github.com/ory/fosite/handler/oauth2.ErrInvalidatedAuthorizeCode`
  has been added.

The following documentation sheds light on how you should update your storage
adapter:

```
// ErrInvalidatedAuthorizeCode is an error indicating that an authorization code has been
// used previously.
var ErrInvalidatedAuthorizeCode = errors.New("Authorization code has ben invalidated")

// AuthorizeCodeStorage handles storage requests related to authorization codes.
type AuthorizeCodeStorage interface {
	// GetAuthorizeCodeSession stores the authorization request for a given authorization code.
	CreateAuthorizeCodeSession(ctx context.Context, code string, request fosite.Requester) (err error)

	// GetAuthorizeCodeSession hydrates the session based on the given code and returns the authorization request.
	// If the authorization code has been invalidated with `InvalidateAuthorizeCodeSession`, this
	// method should return the ErrInvalidatedAuthorizeCode error.
	//
	// Make sure to also return the fosite.Requester value when returning the ErrInvalidatedAuthorizeCode error!
	GetAuthorizeCodeSession(ctx context.Context, code string, session fosite.Session) (request fosite.Requester, err error)

	// InvalidateAuthorizeCodeSession is called when an authorize code is being used. The state of the authorization
	// code should be set to invalid and consecutive requests to GetAuthorizeCodeSession should return the
	// ErrInvalidatedAuthorizeCode error.
	InvalidateAuthorizeCodeSession(ctx context.Context, code string) (err error)
}
```

## 0.19.0

This release improves the OpenID Connect vaildation strategy which now properly
handles `prompt`, `max_age`, and `id_token_hint` at the `/oauth2/auth` endpoint
instead of the `/oauth2/token` endpoint.

To achieve this, the `OpenIDConnectRequestValidator` has been modified and now
requires a `jwt.JWTStrategy` (implemented by, for example
`jwt.RS256JWTStrategy`).

The compose package has been updated accordingly. You should not expect any
major breaking changes from this release.

## 0.18.0

This release allows the introspection handler to return the token type (e.g.
`access_token`, `refresh_token`) of the introspected token. To achieve that,
some breaking API changes have been introduced:

- `OAuth2.IntrospectToken(ctx context.Context, token string, tokenType TokenType, session Session, scope ...string) (AccessRequester, error)`
  is now
  `OAuth2.IntrospectToken(ctx context.Context, token string, tokenType TokenType, session Session, scope ...string) (TokenType, AccessRequester, error)`.
- `TokenIntrospector.IntrospectToken(ctx context.Context, token string, tokenType TokenType, accessRequest AccessRequester, scopes []string) (error)`
  is now
  `TokenIntrospector.IntrospectToken(ctx context.Context, token string, tokenType TokenType, accessRequest AccessRequester, scopes []string) (TokenType, error)`.

This patch also resolves a misconfigured json key in the `IntrospectionResponse`
struct. `AccessRequester AccessRequester json:",extra"` is now properly declared
as `AccessRequester AccessRequester json:"extra"`.

## 0.17.0

This release resolves a security issue (reported by
[platform.sh](https://www.platform.sh)) related to potential storage
implementations. This library used to pass all of the request body from both
authorize and token endpoints to the storage adapters. As some of these values
are needed in consecutive requests, some storage adapters chose to drop the full
body to the database.

This implied that confidential parameters, such as the `client_secret` which can
be passed in the request body since version 0.15.0, were stored as key/value
pairs in plaintext in the database. While most client secrets are generated
programmatically (as opposed to set by the user), it's a considerable security
issue nonetheless.

The issue has been resolved by sanitizing the request body and only including
those values truly required by their respective handlers. This lead to two
breaking changes in the API:

1. The `fosite.Requester` interface has a new method
   `Sanitize(allowedParameters []string) Requester` which returns a sanitized
   clone of the method receiver. If you do not use your own `fosite.Requester`
   implementation, this won't affect you.
2. If you use the PKCE handler, you will have to add three new methods to your
   storage implementation. The methods to be added work exactly like, for
   example `CreateAuthorizeCodeSession`. A reference implementation can be found
   in [./storage/memory.go](./storage/memory.go). The method signatures are as
   follows:

```go
type PKCERequestStorage interface {
	GetPKCERequestSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error)
	CreatePKCERequestSession(ctx context.Context, signature string, requester fosite.Requester) error
	DeletePKCERequestSession(ctx context.Context, signature string) error
}
```

We encourage you to upgrade to this release and check your storage
implementations and potentially remove old data.

We would like to thank [platform.sh](https://www.platform.sh) for sponsoring the
development of a patch that resolves this issue.

## 0.16.0

This patch introduces `SendDebugMessagesToClients` to the Fosite struct which
enables/disables sending debug information to clients. Debug information may
contain sensitive information as it forwards error messages from, for example,
storage implementations. For this reason, `RevealDebugPayloads` defaults to
false. Keep in mind that the information may be very helpful when specific OAuth
2.0 requests fail and we generally recommend displaying debug information.

Additionally, error keys for JSON changed which caused a new minor version,
speicifically
[`statusCode` was changed to `status_code`](https://github.com/ory/fosite/pull/242/files#diff-dd25e0e0a594c3f3592c1c717039b85eR221).

## 0.15.0

This release focuses on improving compatibility with OpenID Connect
Certification and better error context.

- Error handling is improved by explicitly adding debug information (e.g. "Token
  invalid because it was not found in the database") to the error object.
  Previously, the original error was prepended which caused weird formatting
  issues.
- Allows client credentials in POST body at the `/oauth2/token` endpoint. Please
  note that this method is not recommended to be used, unless the client making
  the request is unable to use HTTP Basic Authorization.
- Allows public clients (without secret) to access the `/oauth2/token` endpoint
  which was previously only possible by adding an arbitrary secret.

This release has no breaking changes to the external API but due to the nature
of the changes, it is released as a new major version.

## 0.14.0

Improves error contexts. A breaking code changes to the public API was reverted
with 0.14.1.

## 0.13.0

### Breaking changes

`glide` was replaced with `dep`.

## 0.12.0

### Breaking changes

#### Improved cryptographic methods

- The minimum required secret length used to generate signatures of access
  tokens has increased from 16 to 32 byte.
- The algorithm used to generate access tokens using the HMAC-SHA strategy has
  changed from HMAC-SHA256 to HMAC-SHA512.

## 0.11.0

### Non-breaking changes

#### Storage adapter

To simplify the storage adapter logic, and also reduce the likelihoods of bugs
within the storage adapter, the interface was greatly simplified. Specifically,
these two methods have been removed:

- `PersistRefreshTokenGrantSession(ctx context.Context, requestRefreshSignature, accessSignature, refreshSignature string, request fosite.Requester) error`
- `PersistAuthorizeCodeGrantSession(ctx context.Context, authorizeCode, accessSignature, refreshSignature string, request fosite.Requester) error`

For this change, you don't need to do anything. You can however simply delete
those two methods from your store.

#### Reducing use of gomock

In the long term, fosite should remove all gomocks and instead test against the
internal implementations. This will increase iterations per line during tests
and reduce annoying mock updates.

### Breaking Changes

#### `fosite/handler/oauth2.AuthorizeCodeGrantStorage` was removed

`AuthorizeCodeGrantStorage` was used specifically in the composer. Refactor
references to `AuthorizeCodeGrantStorage` with `CoreStorage`.

#### `fosite/handler/oauth2.RefreshTokenGrantStorage` was removed

`RefreshTokenGrantStorage` was used specifically in the composer. Refactor
references to `RefreshTokenGrantStorage` with `CoreStorage`.

#### `fosite/handler/oauth2.AuthorizeCodeGrantStorage` was removed

`AuthorizeCodeGrantStorage` was used specifically in the composer. Refactor
references to `AuthorizeCodeGrantStorage` with `CoreStorage`.

#### WildcardScopeStrategy

A new [scope strategy](https://github.com/ory/fosite/pull/187) was introduced
called `WildcardScopeStrategy`. This strategy is now the default when using the
composer. To set the HierarchicScopeStrategy strategy, do:

```
import "github.com/ory/fosite/compose"

var config = &compose.Config{
    ScopeStrategy: fosite.HierarchicScopeStrategy,
}
```

#### Refresh tokens and authorize codes are no longer JWTs

Using JWTs for refresh tokens and authorize codes did not make sense:

1. Refresh tokens are long-living credentials, JWTs require an expiry date.
2. Refresh tokens are never validated client-side, only server-side. Thus access
   to the store is available.
3. Authorize codes are never validated client-side, only server-side.

Also, one compose method changed due to this:

```go
package compose

// ..

- func NewOAuth2JWTStrategy(key *rsa.PrivateKey) *oauth2.RS256JWTStrategy
+ func NewOAuth2JWTStrategy(key *rsa.PrivateKey, strategy *oauth2.HMACSHAStrategy) *oauth2.RS256JWTStrategy
```

#### Delete access tokens when persisting refresh session

Please delete access tokens in your store when you persist a refresh session.
This increases security. Here is an example of how to do that using only
existing methods:

```go
func (s *MemoryStore) PersistRefreshTokenGrantSession(ctx context.Context, originalRefreshSignature, accessSignature, refreshSignature string, request fosite.Requester) error {
	if ts, err := s.GetRefreshTokenSession(ctx, originalRefreshSignature, nil); err != nil {
		return err
	} else if err := s.RevokeAccessToken(ctx, ts.GetID()); err != nil {
		return err
	} else if err := s.RevokeRefreshToken(ctx, ts.GetID()); err != nil {
 		return err
 	} else if err := s.CreateAccessTokenSession(ctx, accessSignature, request); err != nil {
 		return err
 	} else if err := s.CreateRefreshTokenSession(ctx, refreshSignature, request); err != nil {
 		return err
 	}

 	return nil
}
```

## 0.10.0

It is no longer possible to introspect authorize codes, and passing scopes to
the introspector now also checks refresh token scopes.

## 0.9.0

This patch adds the ability to pass a custom hasher to `compose.Compose`, which
is a breaking change. You can pass nil for the fosite default hasher:

```
package compose

-func Compose(config *Config, storage interface{}, strategy interface{}, factories ...Factory) fosite.OAuth2Provider {
+func Compose(config *Config, storage interface{}, strategy interface{}, hasher fosite.Hasher, factories ...Factory) fosite.OAuth2Provider {
```

## 0.8.0

This patch addresses some inconsistencies in the public interfaces. Also
remaining references to the old repository location at `ory-am/fosite` where
updated to `ory/fosite`.

### Breaking changes

#### `ClientManager`

The
[`ClientManager`](https://github.com/ory/fosite/blob/master/client_manager.go)
interface changed, as a context parameter was added:

```go
type ClientManager interface {
	// GetClient loads the client by its ID or returns an error
  	// if the client does not exist or another error occurred.
-	GetClient(id string) (Client, error)
+	GetClient(ctx context.Context, id string) (Client, error)
}
```

#### `OAuth2Provider`

The [OAuth2Provider](https://github.com/ory/fosite/blob/master/oauth2.go)
interface changed, as the need for passing down `*http.Request` was removed.
This is justifiable because `NewAuthorizeRequest` and `NewAccessRequest` already
contain `*http.Request`.

The public api of those two methods changed:

```go
-	NewAuthorizeResponse(ctx context.Context, req *http.Request, requester AuthorizeRequester, session Session) (AuthorizeResponder, error)
+	NewAuthorizeResponse(ctx context.Context, requester AuthorizeRequester, session Session) (AuthorizeResponder, error)


-	NewAccessResponse(ctx context.Context, req *http.Request, requester AccessRequester) (AccessResponder, error)
+	NewAccessResponse(ctx context.Context, requester AccessRequester) (AccessResponder, error)
```

## 0.7.0

Breaking changes:

- Replaced `"golang.org/x/net/context"` with `"context"`.
- Move the repo from `github.com/ory-am/fosite` to `github.com/ory/fosite`

## 0.6.0

A bug related to refresh tokens was found. To mitigate it, a `Clone()` method
has been introduced to the `fosite.Session` interface. If you use a custom
session object, this will be a breaking change. Fosite's default sessions have
been upgraded and no additional work should be required. If you use your own
session struct, we encourage using package `gob/encoding` to deep-copy it in
`Clone()`.

## 0.5.0

Breaking changes:

- `compose.OpenIDConnectExplicit` is now `compose.OpenIDConnectExplicitFactory`
- `compose.OpenIDConnectImplicit` is now `compose.OpenIDConnectImplicitFactory`
- `compose.OpenIDConnectHybrid` is now `compose.OpenIDConnectHybridFactory`
- The token introspection handler is no longer added automatically by
  `compose.OAuth2*`. Add `compose.OAuth2TokenIntrospectionFactory` to your
  composer if you need token introspection.
- Session refactor:
  - The HMACSessionContainer was removed and replaced by `fosite.Session` /
    `fosite.DefaultSession`. All sessions must now implement this signature. The
    new session interface allows for better expiration time handling.
  - The OpenID `DefaultSession` signature changed as well, it is now
    implementing the `fosite.Session` interface

## 0.4.0

Breaking changes:

- `./fosite-example` is now a separate repository:
  https://github.com/ory-am/fosite-example
- `github.com/ory-am/fosite/fosite-example/pkg.Store` is now
  `github.com/ory-am/fosite/storage.MemoryStore`
- `fosite.Client` has now a new method called `IsPublic()` which can be used to
  identify public clients who do not own a client secret
- All grant types except the client_credentials grant now allow public clients.
  public clients are usually mobile apps and single page apps.
- `TokenValidator` is now `TokenIntrospector`, `TokenValidationHandlers` is now
  `TokenIntrospectionHandlers`.
- `TokenValidator.ValidateToken` is now `TokenIntrospector.IntrospectToken`
- `fosite.OAuth2Provider.NewIntrospectionRequest()` has been added
- `fosite.OAuth2Provider.WriteIntrospectionError()` has been added
- `fosite.OAuth2Provider.WriteIntrospectionResponse()` has been added

## 0.3.0

- Updated jwt-go from 2.7.0 to 3.0.0

## 0.2.0

Breaking changes:

- Token validation refactored: `ValidateRequestAuthorization` is now `Validate`
  and does not require a http request but instead a token and a token hint. A
  token can be anything, including authorization codes, refresh tokens, id
  tokens, ...
- Remove mandatory scope: The mandatory scope (`fosite`) has been removed as it
  has proven impractical.
- Allowed OAuth2 Client scopes are now being set with `scope` instead of
  `granted_scopes` when using the DefaultClient.
- There is now a scope matching strategy that can be replaced.
- OAuth2 Client scopes are now checked on every grant type.
- Handler subpackages such as `core/client` or `oidc/explicit` have been merged
  and moved one level up
- `handler/oidc` is now `handler/openid`
- `handler/core` is now `handler/oauth2`

## 0.1.0

Initial release
