# OAuth2

[This introduction was taken from the Digital Ocean Blog.](https://www.digitalocean.com/community/tutorials/an-introduction-to-oauth-2)

**OAuth 2** is an authorization framework that enables applications to obtain limited access to user accounts on an HTTP service, such as Facebook, GitHub, and Hydra. It works by delegating user authentication to the service that hosts the user account, and authorizing third-party applications to access the user account. OAuth 2 provides authorization flows for web and desktop applications, and mobile devices.

![](abstract_flow.png) 

Here is a more detailed explanation of the steps in the diagram:

* The application requests authorization to access service resources from the user
* If the user authorized the request, the application receives an authorization grant
* The application requests an access token from the authorization server (API) by presenting authentication of its own identity, and the authorization grant.
* If the application identity is authenticated and the authorization grant is valid, the authorization server (API) issues an access token to the application. Authorization is complete.
* The application requests the resource from the resource server (API) and presents the access token for authentication.
* If the access token is valid, the resource server (API) serves the resource to the application.

The actual flow of this process will differ depending on the authorization grant type in use, but this is the general idea.

Read more on OAuth2 on [the Digital Ocean Blog](https://www.digitalocean.com/community/tutorials/an-introduction-to-oauth-2). We also recommend reading [API Security: Deep Dive into OAuth and OpenID Connect](http://nordicapis.com/api-security-oauth-openid-connect-depth/).

**Glossary**
* **The resource owner** is the user who authorizes an application to access their account. The application's access to the user's account is limited to the "scope" of the authorization granted (e.g. read or write access).
* **Authorization Server (Hydra)** verifies the identity of the user and issues access tokens to the *client application*.
* **Client** is the *application* that wants to access the user's account. Before it may do so, it must be authorized by the user.
* **Identity Provider** contains a log in user interface and a database of all your users. To integrate Hydra, you must modify the Identity Provider. It mus be able to generate consent tokens and ask for the user's consent.
* **User Agent** is usually the resource owner's browser.
* **Consent Endpoint** is an app (e.g. NodeJS) that is able to receive consent challenges and create consent tokens. It must verify the identity of the user that is giving the consent. This can be achieved using Cookie Auth, HTTP Basic Auth, Login HTML Form, or any other mean of authentication. Upon authentication, the user must be asked if he consents to allowing the client access to his resources.

## OAuth2 Clients

We already covered some basic OAuth2 conctepts [in the Introduction](introduction.html). You can manage *clients* using the cli or the HTTP REST API.

* **CLI:** `hydra clients -h`
* **REST:** Read the [API Docs](http://docs.hdyra.apiary.io/#reference/oauth2-clients)

## Authentication Flow

### Overview

![](hydra.png)

2. 1. A *client* application (app in browser in laptop) requests an access token from a resource owner: `https://hydra.myapp.com/oauth2/auth?client_id=c3b49cf0-88e4-4faa-9489-28d5b8957858&response_type=code&scope=core+hydra&state=vboeidlizlxrywkwlsgeggff&nonce=tedgziijemvninkuotcuuiof`.
2. Hydra generates a consent challenge and forwards the *user agent* (browser in laptop) to the *consent endpoint*: `https://login.myapp.com/?challenge=eyJhbGciOiJSUzI1N...`.
3. The *consent endpoint* verifies the resource owner's identity (e.g. cookie, username/password login form, ...). The consent challenge is then decoded and the information extracted. It is used to show the consent screen: `Do you want to grant _my cool app_ access to all your private data? [Yes] [No]`
4. When consent is given, the *consent endpoint* generates a consent token and redirects the user agent (browser in laptop) back to hydra: `https://hydra.myapp.com/oauth2/auth?client_id=c3b49cf0-88e4-4faa-9489-28d5b8957858&response_type=code&scope=core+hydra&state=vboeidlizlxrywkwlsgeggff&nonce=tedgziijemvninkuotcuuiof&consent=eyJhbGciOiJSU...`.
5. Hydra validates the consent token and issues the access token to the *user agent*.

### Example

In this section we assume that hydra runs on `https://192.168.99.100:4444` and our consent app on `https:/192.168.99.100:3000`.

All user-based OAuth2 requests, including the OpenID Connect workflow, begin at the `/oauth2/auth` endpoint. For example: `https://192.168.99.100:4444/oauth2/auth?client_id=c3b49cf0-88e4-4faa-9489-28d5b8957858&response_type=code&scope=core+hydra&state=wewuphkgywhtldsmainefkyx&nonce=uqfjjzftqpjccdvxltaposri`

If the request includes a valid redirect uri and a valid client id, hydra redirects the user to then consent url. The consent url can be set using the `$CONSENT_URL` environment variable.

Let's set the `$CONSENT_URL` to `https:/192.168.99.100:3000/consent`, where a NodeJS application is running (the *consent app*). Next, Hydra appends a consent challenge to the consent url and redirects the user to it. For example: `http://192.168.99.100:3000/?challenge=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJjM2I0OWNmMC04OGU0LTRmYWEtOTQ4OS0yOGQ1Yjg5NTc4NTgiLCJleHAiOjE0NjQ1MTU0ODIsImp0aSI6IjNmYWRlN2NjLTdlYTItNGViMi05MGI1LWY5OTUwNTI4MzgyOSIsInJlZGlyIjoiaHR0cHM6Ly8xOTIuMTY4Ljk5LjEwMDo0NDQ0L29hdXRoMi9hdXRoP2NsaWVudF9pZD1jM2I0OWNmMC04OGU0LTRmYWEtOTQ4OS0yOGQ1Yjg5NTc4NThcdTAwMjZyZXNwb25zZV90eXBlPWNvZGVcdTAwMjZzY29wZT1jb3JlK2h5ZHJhXHUwMDI2c3RhdGU9d2V3dXBoa2d5d2h0bGRzbWFpbmVma3l4XHUwMDI2bm9uY2U9dXFmamp6ZnRxcGpjY2R2eGx0YXBvc3JpIiwic2NwIjpbImNvcmUiLCJoeWRyYSJdfQ.KpLBotIEE4izVSAjLOeCCfm_wYZ7UWSCA81akr6Ci1yycKs8e_bhBYdSThy8JW3bAvofNcZ0v48ov9KxZVegWm8GuNbBEcNvKeiyW_8PiJXWE92YsMv-tDIL3VFPOp0469FmDLsSg5ohsFj5S89FzykNYfVxLPBAFcAS_JElWbo`

The consent challenge is a signed RSA-SHA 256 (RS256) [JSON Web Token](https://tools.ietf.org/html/rfc7519) and contains the following claims:


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
* **aud:** The client id that initiated the request. You can fetch client data using the [OAuth2 Client API](http://docs.hdyra.apiary.io/#reference/oauth2-clients/oauth2-client/get-an-oauth2-client).
* **exp:** The challenge's expiry date. Consent endpoints must not accept challenges that have expired.
* **redir:** Where the consent endpoint should redirect the user agent to, once consent is given.

Hydra signs the consent token with a key called consent.challenge. The public key can be looked up via the [Key Manager](https://ory-am.gitbooks.io/hydra/content/key_manager.html): `https://192.168.99.100:4444/keys/consent.challenge/public`

Next, the consent-app must check if the user is authenticated. This can be done by e.g. using a session cookie. If the user is not authenticate, he must be challenged to provide valid credentials through e.g. a HTML form. The consent-app could use LDAP, MySQL, RethinkDB or any other backend to store and verify the credentials.

Upon user authentication, the consent-app must ask for the user's consent. This could look like:

> _That super useful service app_ would like to:
> * Know who you are
> * View your extended profile info
> * Get read access to all your cloud pictures
> 
> [Deny] - [Allow]

If the user clicks *Allow*, the consent-app redirects him back to the *redir* claim value. The consent-app appends a signed consent token to the URL: `https://192.168.99.100:4444/oauth2/auth?client_id=c3b49cf0-88e4-4faa-9489-28d5b8957858&response_type=code&scope=core+hydra&state=myhnxqmwzhteycweovblswwj&nonce=gmmsuvtsbtoeumeohesztshg&consent=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJjM2I0OWNmMC04OGU0LTRmYWEtOTQ4OS0yOGQ1Yjg5NTc4NTgiLCJleHAiOjE0NjQ1MTUwOTksInNjcCI6WyJjb3JlIiwiaHlkcmEiXSwic3ViIjoiam9obi5kb2VAbWUuY29tIiwiaWF0IjoxNDY0NTExNTE1fQ.tX5TKdP9hHCgPbqBzKIYMjJVwqOdxf5ACScmQ6t20Qteo8AYEfavGwq8KxRF1Oz_otcQDdZY--jcl1caom0yT2eTvj1d9E2Hs7eXmYuW_xF9pTpmDwJnrcOlONFKsNZN97n41qprzMrsX5ez0T5AcopGwpPMxKhwGDSXq9CQgQU`.

The consent token is a RSA-SHA 256 (RS256) signed [JSON Web Token](https://tools.ietf.org/html/rfc7519) that contains the following claims:

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
* **aud:** The client id that initiated the OAuth2 request. You can fetch client data using the [OAuth2 Client API](http://docs.hdyra.apiary.io/#reference/oauth2-clients/oauth2-client/get-an-oauth2-client).
* **exp:** The expiry date of this token. Use very short lifespans (< 5 min).
* **iat:** The tokens issuance time.
* **id_ext:** If set, pass this extra data to the id token *(optional)*
* **at_ext:** If set, pass this extra data to the access token session. You can retrieve the data by using the warden endpoints *(optional)*.

Hydra validates the consent token with consent-app's public key. The public key must be stored in the (https://ory-am.gitbooks.io/hydra/content/key_manager.html) at `https://localhost:4444/keys/consent.endpoint/public`

If you want, you can use the Key Manager to store and retrieve private keys as well. When Hydra boots for the first time, a private/public `consent.endpoint` keypair is created. You can that keypair to sign consent tokens. The private key is available at `https://localhost:4444/keys/asymmetric/consent.endpoint/private`.

# OpenID Connect 1.0

**OpenID Connect 1.0** is a simple identity layer on top of the OAuth 2.0 protocol. It allows Clients to verify the identity of the End-User based on the authentication performed by an Authorization Server, as well as to obtain basic profile information about the End-User in an interoperable and REST-like manner.

OpenID Connect allows clients of all types, including Web-based, mobile, and JavaScript clients, to request and receive information about authenticated sessions and end-users. The specification suite is extensible, allowing participants to use optional features such as encryption of identity data, discovery of OpenID Providers, and session management, when it makes sense for them.

There are different work flows for OpenID Connect 1.0. We are still looking for a good introductionairy blog post. If you have one, let us know.

*In a nutshell:* Add `openid` to your scope when making an OAuth2 Authorize Code request. You will receive an `id_token` alongside an `access_token` and a `refresh_token` when making the code exchange.