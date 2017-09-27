# HydraOAuth2OpenIdConnectServer.WellKnown

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**authorizationEndpoint** | **String** | URL of the OP&#39;s OAuth 2.0 Authorization Endpoint | 
**idTokenSigningAlgValuesSupported** | **[String]** | JSON array containing a list of the JWS signing algorithms (alg values) supported by the OP for the ID Token to encode the Claims in a JWT [JWT]. The algorithm RS256 MUST be included. The value none MAY be supported, but MUST NOT be used unless the Response Type used returns no ID Token from the Authorization Endpoint (such as when using the Authorization Code Flow). | 
**issuer** | **String** | URL using the https scheme with no query or fragment component that the OP asserts as its Issuer Identifier. If Issuer discovery is supported , this value MUST be identical to the issuer value returned by WebFinger. This also MUST be identical to the iss Claim value in ID Tokens issued from this Issuer. | 
**jwksUri** | **String** | URL of the OP&#39;s JSON Web Key Set [JWK] document. This contains the signing key(s) the RP uses to validate signatures from the OP. The JWK Set MAY also contain the Server&#39;s encryption key(s), which are used by RPs to encrypt requests to the Server. When both signing and encryption keys are made available, a use (Key Use) parameter value is REQUIRED for all keys in the referenced JWK Set to indicate each key&#39;s intended usage. Although some algorithms allow the same key to be used for both signatures and encryption, doing so is NOT RECOMMENDED, as it is less secure. The JWK x5c parameter MAY be used to provide X.509 representations of keys provided. When used, the bare key values MUST still be present and MUST match those in the certificate. | 
**responseTypesSupported** | **[String]** | JSON array containing a list of the OAuth 2.0 response_type values that this OP supports. Dynamic OpenID Providers MUST support the code, id_token, and the token id_token Response Type values. | 
**subjectTypesSupported** | **[String]** | JSON array containing a list of the Subject Identifier types that this OP supports. Valid types include pairwise and public. | 
**tokenEndpoint** | **String** | URL of the OP&#39;s OAuth 2.0 Token Endpoint | 


