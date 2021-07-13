---
id: security-architecture
title: Security Architecture
---

Hydra is built with tough security in mind.

<!-- toc -->

## OAuth 2.0 Security Overview

Hydra is an implementation of the security-first Fosite OAuth 2.0 SDK
([https://github.com/ory/fosite](https://github.com/ory/fosite)). Fosite
respects the
[OAuth 2.0 Threat Model and Security Considerations](https://tools.ietf.org/html/rfc6819#section-5.1.5.3)
by the IETF, specifically:

- No Cleartext Storage of Credentials
- Encryption of Credentials
- Use Short Expiration Time
- Limit Number of Usages or One-Time Usage
- Bind Token to Client id
- Automatic Revocation of Derived Tokens If Abuse Is Detected
- Binding of Refresh Token to "client_id"
- Refresh Token Rotation
- Revocation of Refresh Tokens
- Validate Pre-Registered "redirect_uri"
- Binding of Authorization "code" to "client_id"
- Binding of Authorization "code" to "redirect_uri"
- Opaque access tokens
- Opaque refresh tokens
- Ensure Confidentiality of Requests
- Use of Asymmetric Cryptography
- Enforcing random states: Without a random-looking state or OpenID Connect
  nonce the request will fail.

Additionally these safeguards are implemented:

- Advanced Token Validation: Tokens are laid out as
  &lt;key&gt;.&lt;signature&gt; where &lt;signature&gt; is created using
  HMAC-SHA256 using a global secret.

### Advanced Token Validation (Datastore Security)

For a OAuth2 access token, refresh token or authorize code to be valid, one
requires both the key and the signature (formatted as
&lt;key&gt;.&lt;signature&gt;). Only the signature is stored in the datastore
(SQL), thus a compromised datastore will not allow an attacker to gain access to
any valid authorize codes, access tokens, or refresh tokens.

Because HMAC-SHA256 is used, the System Secret is required to create valid
key-signature pairs, rendering an attacker unable to inject new codes or tokens
into a compromised datastore.

## Cryptography

Hydra uses different cryptographic methods, this is an overview of all of them.

### AES-GCM

AES-GCM is used to encrypt JWKs at rest using a key size of 256 bit which
exceeds requirements by Lenstra, ECRYPT II, NIST, ANSSI, and BSI, see
[https://www.keylength.com/en/compare/](https://www.keylength.com/en/compare/).

GCM (Galois/Counter Mode) is an authenticated encryption algorithm designed to
provide both data authenticity (integrity) and confidentiality. GCM uses a nonce
(“IV”) that has an upper limit of 2^32 nonces. If more nonces are used, there is
risk of repeats. This means that you risk collisions when storing more than 2^32
documents authenticated with GCM. Because AES-GCM is only used to encrypt data
at rest, this is might only impose a problem if

1. more than 2^32 documents are stored using AES-GCM
2. an attacker gains access to the datastore where &gt; 2^32 documents are
   stored
3. the attacker is able to exploit repeats, for example by authenticating
   malicious documents

### RS256

RSASSA-PKCS1-v1_5 using SHA-256 (RS256) is used to sign JWTs. It’s use is
recommended by the JWA specification, see
[https://www.rfc-editor.org/rfc/rfc7518.txt](https://www.rfc-editor.org/rfc/rfc7518.txt)

The RSA Key size is 4096 bit long, exceeding the minimum requirement of 2048 bit
by
[https://www.rfc-editor.org/rfc/rfc7518.txt](https://www.rfc-editor.org/rfc/rfc7518.txt).

Recommendations from NIST, ANSSI, IAD-NSA, BSI, Lenstra and others vary between
1300 and 2048 bit key lengths for asymmetric cryptography based on discrete
logarithms (RSA). 4096 exceeds all recommendations for 2017 from all
authorities, see
[https://www.keylength.com/en/compare/](https://www.keylength.com/en/compare/).

### HMAC-SHA256

HMAC (FIPS 198) with SHA256 (FIPS 180-4) is used to sign access tokens,
authorize codes and refresh tokens. SHA-2 (with 256 bit) is encouraged by NIST,
see
[http://csrc.nist.gov/groups/ST/hash/policy.html](http://csrc.nist.gov/groups/ST/hash/policy.html)

### BCrypt

BCrypt is used to hash client credentials at rest. It is not officially
recommended by NIST as it is not based on hashing primitives such as SHA-2, but
rather on Blowfish. However, BCrypt is much stronger than any other (salted)
hashing method for passwords, has wide adoption and is an official golang/x
library.

I recommend reading this thread on Security Stack Exchange on BCrypt, SCrypt and
PBKDF2:
[https://security.stackexchange.com/questions/4781/do-any-security-experts-recommend-bcrypt-for-password-storage](https://security.stackexchange.com/questions/4781/do-any-security-experts-recommend-bcrypt-for-password-storage)

Be aware that BCrypt causes very high CPU loads, depending on the Workload
Factor. We strongly advise reducing the number of requests that use Basic
Authorization.

## How does Access Control work with Hydra?

Hydra supports two concepts of authorization. One is called Token Introspection
which is a standard by the IETF (
[https://tools.ietf.org/html/rfc7662](https://tools.ietf.org/html/rfc7662) ). It
is primarily targeted at third-party services and is usually used by a
programmatic API. It can be used by first-party services too. The sole purpose
of this endpoint is to check whether an access token is valid or not.
