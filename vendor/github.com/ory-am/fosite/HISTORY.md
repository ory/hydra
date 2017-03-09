This is a list of breaking changes. As long as `1.0.0` is not released, breaking changes will be addressed as minor version
bumps (`0.1.0` -> `0.2.0`).

## 0.6.0

A bug related to refresh tokens was found. To mitigate it, a `Clone()` method has been introduced to the `fosite.Session` interface.
If you use a custom session object, this will be a breaking change. Fosite's default sessions have been upgraded and no additional
work should be required. If you use your own session struct, we encourage using package `gob/encoding` to deep-copy it in `Clone()`.

## 0.5.0

Breaking changes:

* `compose.OpenIDConnectExplicit` is now `compose.OpenIDConnectExplicitFactory`
* `compose.OpenIDConnectImplicit` is now `compose.OpenIDConnectImplicitFactory`
* `compose.OpenIDConnectHybrid` is now `compose.OpenIDConnectHybridFactory`
* The token introspection handler is no longer added automatically by `compose.OAuth2*`. Add `compose.OAuth2TokenIntrospectionFactory`
to your composer if you need token introspection.
* Session refactor:
  * The HMACSessionContainer was removed and replaced by `fosite.Session` / `fosite.DefaultSession`. All sessions
  must now implement this signature. The new session interface allows for better expiration time handling.
  * The OpenID `DefaultSession` signature changed as well, it is now implementing the `fosite.Session` interface

## 0.4.0

Breaking changes:

* `./fosite-example` is now a separate repository: https://github.com/ory-am/fosite-example
* `github.com/ory-am/fosite/fosite-example/pkg.Store` is now `github.com/ory-am/fosite/storage.MemoryStore`
* `fosite.Client` has now a new method called `IsPublic()` which can be used to identify public clients who do not own a client secret
* All grant types except the client_credentials grant now allow public clients. public clients are usually mobile apps and single page apps.
* `TokenValidator` is now `TokenIntrospector`, `TokenValidationHandlers` is now `TokenIntrospectionHandlers`.
* `TokenValidator.ValidateToken` is now `TokenIntrospector.IntrospectToken`
* `fosite.OAuth2Provider.NewIntrospectionRequest()` has been added
* `fosite.OAuth2Provider.WriteIntrospectionError()` has been added
* `fosite.OAuth2Provider.WriteIntrospectionResponse()` has been added

## 0.3.0

* Updated jwt-go from 2.7.0 to 3.0.0

## 0.2.0

Breaking changes:

* Token validation refactored: `ValidateRequestAuthorization` is now `Validate` and does not require a http request
but instead a token and a token hint. A token can be anything, including authorization codes, refresh tokens,
id tokens, ...
* Remove mandatory scope: The mandatory scope (`fosite`) has been removed as it has proven impractical.
* Allowed OAuth2 Client scopes are now being set with `scope` instead of `granted_scopes` when using the DefaultClient.
* There is now a scope matching strategy that can be replaced.
* OAuth2 Client scopes are now checked on every grant type.
* Handler subpackages such as `core/client` or `oidc/explicit` have been merged and moved one level up
* `handler/oidc` is now `handler/openid`
* `handler/core` is now `handler/oauth2`

## 0.1.0

Initial release