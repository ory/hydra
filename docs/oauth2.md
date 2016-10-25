# OAuth 2.0 & OpenID Connect

If you are new to OAuth2, please read the [Introduction to OAuth 2.0 and OpenID Connect](README.md#introduction-to-oauth-20-and-openid-connect)
first.

## Overview

This section defines a glossary, provides additional information on OpenID Connect and introduces OAuth 2.0 Clients.

### Glossary

1. **The resource owner** is the user who authorizes an application to access their account. The application's access to
the user's account is limited to the "scope" of the authorization granted (e.g. read or write access).
2. **Authorization Server (Hydra)** verifies the identity of the user and issues access tokens to the *client application*.
3. **Client** is the *application* that wants to access the user's account. Before it may do so, it must be authorized
by the user.
4. **Identity Provider** contains a log in user interface and a database of all your users. To integrate Hydra,
you must modify the Identity Provider. It mus be able to generate consent tokens and ask for the user's consent.
5. **User Agent** is usually the resource owner's browser.
6. **Consent App** is an app (e.g. NodeJS) that is able to receive consent challenges and create consent tokens.
It must verify the identity of the user that is giving the consent. This can be achieved using Cookie Auth,
HTTP Basic Auth, Login HTML Form, or any other mean of authentication. Upon authentication, the user must be asked
if he consents to allowing the client access to his resources.

Examples:
1. Peter wants to give MyPhotoBook access to his Dropbox. Peter is the resource owner.
2. The Authorization Server (Hydra) is responsible for managing the access request fom MyPhotoBook. Hydra handles
the communication between the resource owner, the consent endpoint and the client. Hydra is the authorization server.
In this case, Dropbox would be the one who uses Hydra.
3. MyPhotoBook is the client and was issued an id and a password by Hydra. MyPhotoBook uses these credentials
to talk with Hydra.
4. Dropbox has a database and a frontend that allow their users to log in, using their username and password.
This is what an Identity Provider does.
5. The User Agent is Peter's FireFox.
6. The Consent App is a frontend app that asks the user if he is willing to give MyPhotoBook access to his pictures stored
on Dropbox. It is responsible to tell Hydra if the user accepted or rejected the request by MyPhotoBook. The Consent App
uses the Identity Provider to authenticate peter, for example by using cookies or presenting a user/password login view.

### OpenID Connect 1.0

If you are new to OpenID Connect, please read the [Introduction to OAuth 2.0 and OpenID Connect](README.md#introduction-to-oauth-20-and-openid-connect)
first. 

Hydra uses the [JSON Web Key Manager](./jwk.md) to retrieve the
key pair `hydra.openid.id-token` for signing ID tokens. You can use that endpoint to retrieve the public key for verification,
has Hydra is not supporting OpenID Connect Discovery yet.

### OAuth 2.0 Clients

You can manage *OAuth 2.0 clients* using the cli or the HTTP REST API.

* **CLI:** `hydra clients -h`
* **REST:** Read the [API Docs](http://docs.hdyra.apiary.io/#reference/oauth2-clients)

## Consent App Flow

Hydra does not include user authentication and things like lost password, user registration or user activation.
The consent app flow is used to let Hydra identify who resource owner is. In abstract, the consent flow looks like this:

![](../images/consent.png)

1. A *client* application (app in browser in laptop) requests an access token from a resource owner:
`https://hydra.myapp.com/oauth2/auth?client_id=c3b49cf0-88e4-4faa-9489-28d5b8957858&response_type=code&scope=core+hydra&state=vboeidlizlxrywkwlsgeggff&nonce=tedgziijemvninkuotcuuiof`.
2. Hydra generates a consent challenge and forwards the *user agent* (browser in laptop) to the *consent endpoint*:
`https://login.myapp.com/?challenge=eyJhbGciOiJSUzI1N...`.
3. The *consent endpoint* verifies the resource owner's identity (e.g. cookie, username/password login form, ...).
The consent challenge is then decoded and the information extracted. It is used to show the consent screen: `Do you want to grant _my cool app_ access to all your private data? [Yes] [No]`
4. When consent is given, the *consent endpoint* generates a consent response token and redirects the user
agent (browser in laptop) back to hydra:
`https://hydra.myapp.com/oauth2/auth?client_id=c3b49cf0-88e4-4faa-9489-28d5b8957858&response_type=code&scope=core+hydra&state=vboeidlizlxrywkwlsgeggff&nonce=tedgziijemvninkuotcuuiof&consent=eyJhbGciOiJSU...`.
5. Hydra validates the consent response token and issues the access token to the *user agent*.

### Consent App Flow Example

In this section we assume that hydra runs on `https://192.168.99.100:4444` and our
consent app on `https:/192.168.99.100:3000`.

All user-based OAuth2 requests, including the OpenID Connect workflow, begin at the `/oauth2/auth` endpoint.
For example: `https://192.168.99.100:4444/oauth2/auth?client_id=c3b49cf0-88e4-4faa-9489-28d5b8957858&response_type=code&scope=core+hydra&state=wewuphkgywhtldsmainefkyx&nonce=uqfjjzftqpjccdvxltaposri`

If the request includes a valid redirect uri and a valid client id, hydra redirects the user to then consent url.
The consent url can be set using the `CONSENT_URL` environment variable.

Let's set the `CONSENT_URL` to `https:/192.168.99.100:3000/consent`, where a NodeJS application
is running (the *consent app*). Next, Hydra appends a consent challenge to the consent url and redirects the user to it.
For example: `http://192.168.99.100:3000/consent/?challenge=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJjM2I0OWNmMC04OGU0LTRmYWEtOTQ4OS0yOGQ1Yjg5NTc4NTgiLCJleHAiOjE0NjQ1MTU0ODIsImp0aSI6IjNmYWRlN2NjLTdlYTItNGViMi05MGI1LWY5OTUwNTI4MzgyOSIsInJlZGlyIjoiaHR0cHM6Ly8xOTIuMTY4Ljk5LjEwMDo0NDQ0L29hdXRoMi9hdXRoP2NsaWVudF9pZD1jM2I0OWNmMC04OGU0LTRmYWEtOTQ4OS0yOGQ1Yjg5NTc4NThcdTAwMjZyZXNwb25zZV90eXBlPWNvZGVcdTAwMjZzY29wZT1jb3JlK2h5ZHJhXHUwMDI2c3RhdGU9d2V3dXBoa2d5d2h0bGRzbWFpbmVma3l4XHUwMDI2bm9uY2U9dXFmamp6ZnRxcGpjY2R2eGx0YXBvc3JpIiwic2NwIjpbImNvcmUiLCJoeWRyYSJdfQ.KpLBotIEE4izVSAjLOeCCfm_wYZ7UWSCA81akr6Ci1yycKs8e_bhBYdSThy8JW3bAvofNcZ0v48ov9KxZVegWm8GuNbBEcNvKeiyW_8PiJXWE92YsMv-tDIL3VFPOp0469FmDLsSg5ohsFj5S89FzykNYfVxLPBAFcAS_JElWbo`

The consent challenge is a signed RSA-SHA 256 (RS256) [JSON Web Token](https://tools.ietf.org/html/rfc7519) and contains
the following claims:


```
eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJjM2I0OWNmMC04OGU0LTRmYWEtOTQ4OS0yOGQ1Yjg5NTc4NTgiLCJleHAiOjE0NjQ1MTUwOTksImp0aSI6IjY0YzRmNzllLWUwMTYtNDViOC04YzBlLWQ5NmM2NzFjMWU4YSIsInJlZGlyIjoiaHR0cHM6Ly8xOTIuMTY4Ljk5LjEwMDo0NDQ0L29hdXRoMi9hdXRoP2NsaWVudF9pZD1jM2I0OWNmMC04OGU0LTRmYWEtOTQ4OS0yOGQ1Yjg5NTc4NThcdTAwMjZyZXNwb25zZV90eXBlPWNvZGVcdTAwMjZzY29wZT1jb3JlK2h5ZHJhXHUwMDI2c3RhdGU9bXlobnhxbXd6aHRleWN3ZW92Ymxzd3dqXHUwMDI2bm9uY2U9Z21tc3V2dHNidG9ldW1lb2hlc3p0c2hnIiwic2NwIjpbImNvcmUiLCJoeWRyYSJdfQ.v4K1-AuT5Uwu1DRNvdf7SwjjPT8KO97thRYa3pDWzjBLyjkCNvgp0P5V0oA3XqRutoFpYx4AtQyz0bY7n3XcPE7ZQ2nBWTBnZ04GzWbxcJNFhBvgc_jiQBECebdxN29kgxHoU0frtVDcz6Uur468nBa9D_BDBpN-KgEBsI5Hjhc

{
  "aud": "c3b49cf0-88e4-4faa-9489-28d5b8957858",
  "exp": 1464515099,
  "jti": "64c4f79e-e016-45b8-8c0e-d96c671c1e8a",
  "redir": "https://192.168.99.100:4444/oauth2/auth?client_id=c3b49cf0-88e4-4faa-9489-28d5b8957858&response_type=code&scope=core+hydra&state=myhnxqmwzhteycweovblswwj&nonce=gmmsuvtsbtoeumeohesztshg",
  "scp": [
    "core",
    "hydra"
  ]
}
```

The challenge claims are:
* **jti:** A unique id.
* **scp:** The requested scopes, e.g. `["blog.readall", "blog.writeall"]`
* **aud:** The client id that initiated the request. You can fetch client data using the [OAuth2 Client API](http://docs.hdyra.apiary.io/#reference/oauth2/manage-the-oauth2-client-collection).
* **exp:** The challenge's expiry date. Consent endpoints must not accept challenges that have expired.
* **redir:** Where the consent endpoint should redirect the user agent to, once consent is given.

Hydra signs the consent response token with a key called `hydra.consent.challenge`.
The public key can be looked up via the [Key Manager](https://ory-am.gitbooks.io/hydra/content/jwk.html):

```
https://192.168.99.100:4444/keys/hydra.consent.challenge/public
```

Next, the consent-app must check if the user is authenticated. This can be done by e.g. using a session cookie.
If the user is not authenticate, he must be challenged to provide valid credentials through e.g. a HTML form.
The consent-app could use LDAP, MySQL, RethinkDB or any other backend to store and verify the credentials.

Upon user authentication, the consent-app must ask for the user's consent. This could look like:

> _That super useful service app_ would like to:
> * Know who you are
> * View your extended profile info
> * Get read access to all your cloud pictures
> 
> [Deny] - [Allow]

If the user clicks *Allow*, the consent-app redirects him back to the *redir* claim value. The consent-app appends
a signed consent response token to the URL:

```
https://192.168.99.100:4444/oauth2/auth?client_id=c3b49cf0-88e4-4faa-9489-28d5b8957858&response_type=code&scope=core+hydra&state=myhnxqmwzhteycweovblswwj&nonce=gmmsuvtsbtoeumeohesztshg&consent=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJjM2I0OWNmMC04OGU0LTRmYWEtOTQ4OS0yOGQ1Yjg5NTc4NTgiLCJleHAiOjE0NjQ1MTUwOTksInNjcCI6WyJjb3JlIiwiaHlkcmEiXSwic3ViIjoiam9obi5kb2VAbWUuY29tIiwiaWF0IjoxNDY0NTExNTE1fQ.tX5TKdP9hHCgPbqBzKIYMjJVwqOdxf5ACScmQ6t20Qteo8AYEfavGwq8KxRF1Oz_otcQDdZY--jcl1caom0yT2eTvj1d9E2Hs7eXmYuW_xF9pTpmDwJnrcOlONFKsNZN97n41qprzMrsX5ez0T5AcopGwpPMxKhwGDSXq9CQgQU
```

The consent response token is a RSA-SHA 256 (RS256) signed [JSON Web Token](https://tools.ietf.org/html/rfc7519)
that contains the following claims:

```
eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJjM2I0OWNmMC04OGU0LTRmYWEtOTQ4OS0yOGQ1Yjg5NTc4NTgiLCJleHAiOjE0NjQ1MTUwOTksInNjcCI6WyJjb3JlIiwiaHlkcmEiXSwic3ViIjoiam9obi5kb2VAbWUuY29tIiwiaWF0IjoxNDY0NTExNTE1fQ.tX5TKdP9hHCgPbqBzKIYMjJVwqOdxf5ACScmQ6t20Qteo8AYEfavGwq8KxRF1Oz_otcQDdZY--jcl1caom0yT2eTvj1d9E2Hs7eXmYuW_xF9pTpmDwJnrcOlONFKsNZN97n41qprzMrsX5ez0T5AcopGwpPMxKhwGDSXq9CQgQU

{
  "aud": "c3b49cf0-88e4-4faa-9489-28d5b8957858",
  "exp": 1464515099,
  "scp": [
    "core",
    "hydra"
  ],
  "sub": "john.doe@me.com",
  "iat": 1464511515,
  "id_ext": { "foo": "bar" },
  "at_ext": { "baz": true }
}
```

The consent claims are:
* **jti:** A unique id.
* **scp:** The scopes the user opted in to *grant* access to, e.g. only `["blog.readall"]`.
* **aud:** The client id that initiated the OAuth2 request. You can fetch
client data using the [OAuth2 Client API](http://docs.hdyra.apiary.io/#reference/oauth2/manage-the-oauth2-client-collection).
* **exp:** The expiry date of this token. Use very short lifespans (< 5 min).
* **iat:** The tokens issuance time.
* **id_ext:** If set, pass this extra data to the id token. This data is not available at OAuth2 Token Introspection
 nor at the warden endpoints. *(optional)*
* **at_ext:** If set, pass this extra data to the access token session. You can retrieve the data
by using OAuth2 Token Introspection or the warden endpoints. *(optional)*

Hydra validates the consent response token with consent-app's public key. The public
key must be stored in the [JSON Web Key Manager](./jwk.md)
at `https://localhost:4444/keys/hydra.consent.response/public`

If you want, you can use the Key Manager to store and retrieve private keys as well. When Hydra boots for the first time,
a private/public `hydra.consent.response` keypair is created.
You can that keypair to sign consent response tokens. The private key is available at
`https://localhost:4444/keys/asymmetric/hydra.consent.response/private`.

### Error Handling during Consent App Flow

Hydra follows the OAuth 2.0 error response specifications. Some errors however must be handled by the consent app.
In the case of such an error, the user agent will be redirected to the consent app
endpoint and an `error` and `error_description` query parameter will be appended to the URL.

# OAuth2 Token Introspection

OAuth2 Token Introspection is an [IETF](https://tools.ietf.org/html/rfc7662) standard.
It defines a method for a protected resource to query
an OAuth 2.0 authorization server to determine the active state of an
OAuth 2.0 token and to determine meta-information about this token.
OAuth 2.0 deployments can use this method to convey information about
the authorization context of the token from the authorization server
to the protected resource.

The Token Introspection endpoint is documented in the
[API Docs](http://docs.hdyra.apiary.io/#reference/oauth2/oauth2-token-introspection).