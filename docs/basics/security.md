## Security

> Why should I use Hydra? It's not that hard to implement two OAuth2 endpoints and there are numerous SDKs out there!

OAuth2 and OAuth2 related specifications are over 200 written pages. Implementing OAuth2 is easy, getting it right is hard.
Even if you use a secure SDK (there are numerous SDKs not secure by design in the wild), messing up the implementation
is a real threat - no matter how good you or your team is. To err is human. Now, let us take a look at security in Hydra.

Hydra uses [Fosite](https://github.com/ory-am/fosite#a-word-on-security), a secure-by-design OAuth2 SDK. Fosite implements
best practices proposed by the IETF:  
* [No Cleartext Storage of Credentials](https://tools.ietf.org/html/rfc6819#section-5.1.4.1.3)
* [Encryption of Credentials](https://tools.ietf.org/html/rfc6819#section-5.1.4.1.4)
* [Use Short Expiration Time](https://tools.ietf.org/html/rfc6819#section-5.1.5.3)
* [Limit Number of Usages or One-Time Usage](https://tools.ietf.org/html/rfc6819#section-5.1.5.4)
* [Bind Token to Client id](https://tools.ietf.org/html/rfc6819#section-5.1.5.8)
* [Automatic Revocation of Derived Tokens If Abuse Is Detected](https://tools.ietf.org/html/rfc6819#section-5.2.1.1)
* [Binding of Refresh Token to "client_id"](https://tools.ietf.org/html/rfc6819#section-5.2.2.2)
* [Refresh Token Rotation](https://tools.ietf.org/html/rfc6819#section-5.2.2.3)
* [Revocation of Refresh Tokens](https://tools.ietf.org/html/rfc6819#section-5.2.2.4)
* [Validate Pre-Registered "redirect_uri"](https://tools.ietf.org/html/rfc6819#section-5.2.3.5)
* [Binding of Authorization "code" to "client_id"](https://tools.ietf.org/html/rfc6819#section-5.2.4.4)
* [Binding of Authorization "code" to "redirect_uri"](https://tools.ietf.org/html/rfc6819#section-5.2.4.6)
* [Opaque access tokens](https://tools.ietf.org/html/rfc6749#section-1.4)
* [Opaque refresh tokens](https://tools.ietf.org/html/rfc6749#section-1.5)
* [Ensure Confidentiality of Requests](https://tools.ietf.org/html/rfc6819#section-5.1.1)
* [Use of Asymmetric Cryptography](https://tools.ietf.org/html/rfc6819#section-5.1.4.1.5)
* **Enforcing random states:** Without a random-looking state or OpenID Connect nonce the request will fail.
* **Advanced Token Validation:** Tokens are laid out as `<key>.<signature>` where `<signature>` is created using HMAC-SHA256
     and a global secret. This is what a token can look like: `/tgBeUhWlAT8tM8Bhmnx+Amf8rOYOUhrDi3pGzmjP7c=.BiV/Yhma+5moTP46anxMT6cWW8gz5R5vpC9RbpwSDdM=`
* **Enforcing scopes:** By default, you always need to include the `core` scope or Hydra will not execute the request.

Hydra uses [Ladon](https://github.com/ory-am/ladon) for policy management and access control. Ladon's API is minimalistic
and well tested.

Hydra encrypts symmetric and asymmetric keys at rest using AES-GCM 256bit.

Hydra does not store tokens, only their signatures. An attacker gaining database access is neither able to steal tokens nor
to issue new ones.

Hydra has automated unit and integration tests.

Hydra does not use hacks. We would rather rewrite the whole thing instead of introducing a hack.

APIs are uniform, well documented and secured using the warden's access control infrastructure.

Hydra is open source and can be reviewed by anyone.

Hydra is designed by a [security enthusiast](https://github.com/arekkas), who has written and participated in numerous auth* projects.
