# OpenID Connect 1.0

**OpenID Connect 1.0** is a simple identity layer on top of the OAuth 2.0 protocol.
It allows Clients to verify the identity of the End-User based on the authentication performed
by an Authorization Server, as well as to obtain basic profile information about the End-User in an
interoperable and REST-like manner.

OpenID Connect allows clients of all types, including Web-based, mobile, and JavaScript clients,
to request and receive information about authenticated sessions and end-users. The specification
suite is extensible, allowing participants to use optional features such as encryption of identity data,
discovery of OpenID Providers, and session management, when it makes sense for them.

There are different work flows for OpenID Connect 1.0, we recommend checking out the OpenID Connect sandbox at
[openidconnect.net](https://openidconnect.net/).

In a nutshell, add `openid` to the OAuth2 scope when making an OAuth2 Authorize Code request.
You will receive an `id_token` alongside the `access_token` when making the code exchange.


Hydra uses the [JSON Web Key Manager](https://ory-am.gitbooks.io/hydra/content/key_manager.html) to retrieve the
key pair `hydra.openid.id-token` for signing ID tokens. You can use that endpoint to retrieve the public key for verification.
