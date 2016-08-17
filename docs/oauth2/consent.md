# Consent Flow

Hydra does not ship user authentication. This is something you will have to solve yourself. Usually when you are looking
at using OAuth2 for your app, you already have user authentication anyways.

In abstract, a consent flow looks like this:

![](../dist/images/consent.png)

1. A *client* application (app in browser in laptop) requests an access token from a resource owner: `https://hydra.myapp.com/oauth2/auth?client_id=c3b49cf0-88e4-4faa-9489-28d5b8957858&response_type=code&scope=core+hydra&state=vboeidlizlxrywkwlsgeggff&nonce=tedgziijemvninkuotcuuiof`.
2. Hydra generates a consent challenge and forwards the *user agent* (browser in laptop) to the *consent endpoint*: `https://login.myapp.com/?challenge=eyJhbGciOiJSUzI1N...`.
3. The *consent endpoint* verifies the resource owner's identity (e.g. cookie, username/password login form, ...). The consent challenge is then decoded and the information extracted. It is used to show the consent screen: `Do you want to grant _my cool app_ access to all your private data? [Yes] [No]`
4. When consent is given, the *consent endpoint* generates a consent token and redirects the user agent (browser in laptop) back to hydra: `https://hydra.myapp.com/oauth2/auth?client_id=c3b49cf0-88e4-4faa-9489-28d5b8957858&response_type=code&scope=core+hydra&state=vboeidlizlxrywkwlsgeggff&nonce=tedgziijemvninkuotcuuiof&consent=eyJhbGciOiJSU...`.
5. Hydra validates the consent token and issues the access token to the *user agent*.

## Detailed Example

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
* **aud:** The client id that initiated the request. You can fetch client data using the [OAuth2 Client API](http://docs.hdyra.apiary.io/#reference/oauth2/manage-the-oauth2-client-collection).
* **exp:** The challenge's expiry date. Consent endpoints must not accept challenges that have expired.
* **redir:** Where the consent endpoint should redirect the user agent to, once consent is given.

Hydra signs the consent token with a key called consent.challenge.
The public key can be looked up via the [Key Manager](https://ory-am.gitbooks.io/hydra/content/jwk.html):

```
https://192.168.99.100:4444/keys/consent.challenge/public
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
a signed consent token to the URL:

```
https://192.168.99.100:4444/oauth2/auth?client_id=c3b49cf0-88e4-4faa-9489-28d5b8957858&response_type=code&scope=core+hydra&state=myhnxqmwzhteycweovblswwj&nonce=gmmsuvtsbtoeumeohesztshg&consent=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJjM2I0OWNmMC04OGU0LTRmYWEtOTQ4OS0yOGQ1Yjg5NTc4NTgiLCJleHAiOjE0NjQ1MTUwOTksInNjcCI6WyJjb3JlIiwiaHlkcmEiXSwic3ViIjoiam9obi5kb2VAbWUuY29tIiwiaWF0IjoxNDY0NTExNTE1fQ.tX5TKdP9hHCgPbqBzKIYMjJVwqOdxf5ACScmQ6t20Qteo8AYEfavGwq8KxRF1Oz_otcQDdZY--jcl1caom0yT2eTvj1d9E2Hs7eXmYuW_xF9pTpmDwJnrcOlONFKsNZN97n41qprzMrsX5ez0T5AcopGwpPMxKhwGDSXq9CQgQU
```

The consent token is a RSA-SHA 256 (RS256) signed [JSON Web Token](https://tools.ietf.org/html/rfc7519)
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
* **aud:** The client id that initiated the OAuth2 request. You can fetch client data using the [OAuth2 Client API](http://docs.hdyra.apiary.io/#reference/oauth2/manage-the-oauth2-client-collection).
* **exp:** The expiry date of this token. Use very short lifespans (< 5 min).
* **iat:** The tokens issuance time.
* **id_ext:** If set, pass this extra data to the id token *(optional)*
* **at_ext:** If set, pass this extra data to the access token session. You can retrieve the data by using the warden endpoints *(optional)*.

Hydra validates the consent token with consent-app's public key. The public key must be stored in the (https://ory-am.gitbooks.io/hydra/content/key_manager.html) at `https://localhost:4444/keys/consent.endpoint/public`

If you want, you can use the Key Manager to store and retrieve private keys as well. When Hydra boots for the first time, a private/public `consent.endpoint` keypair is created. You can that keypair to sign consent tokens. The private key is available at `https://localhost:4444/keys/asymmetric/consent.endpoint/private`.
