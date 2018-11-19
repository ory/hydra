
# JSONWebKey

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**alg** | **String** | The \&quot;alg\&quot; (algorithm) parameter identifies the algorithm intended for use with the key.  The values used should either be registered in the IANA \&quot;JSON Web Signature and Encryption Algorithms\&quot; registry established by [JWA] or be a value that contains a Collision- Resistant Name. | 
**crv** | **String** |  |  [optional]
**d** | **String** |  |  [optional]
**dp** | **String** |  |  [optional]
**dq** | **String** |  |  [optional]
**e** | **String** |  |  [optional]
**k** | **String** |  |  [optional]
**kid** | **String** | The \&quot;kid\&quot; (key ID) parameter is used to match a specific key.  This is used, for instance, to choose among a set of keys within a JWK Set during key rollover.  The structure of the \&quot;kid\&quot; value is unspecified.  When \&quot;kid\&quot; values are used within a JWK Set, different keys within the JWK Set SHOULD use distinct \&quot;kid\&quot; values.  (One example in which different keys might use the same \&quot;kid\&quot; value is if they have different \&quot;kty\&quot; (key type) values but are considered to be equivalent alternatives by the application using them.)  The \&quot;kid\&quot; value is a case-sensitive string. | 
**kty** | **String** | The \&quot;kty\&quot; (key type) parameter identifies the cryptographic algorithm family used with the key, such as \&quot;RSA\&quot; or \&quot;EC\&quot;. \&quot;kty\&quot; values should either be registered in the IANA \&quot;JSON Web Key Types\&quot; registry established by [JWA] or be a value that contains a Collision- Resistant Name.  The \&quot;kty\&quot; value is a case-sensitive string. | 
**n** | **String** |  |  [optional]
**p** | **String** |  |  [optional]
**q** | **String** |  |  [optional]
**qi** | **String** |  |  [optional]
**use** | **String** | Use (\&quot;public key use\&quot;) identifies the intended use of the public key. The \&quot;use\&quot; parameter is employed to indicate whether a public key is used for encrypting data or verifying the signature on data. Values are commonly \&quot;sig\&quot; (signature) or \&quot;enc\&quot; (encryption). | 
**x** | **String** |  |  [optional]
**x5c** | **List&lt;String&gt;** | The \&quot;x5c\&quot; (X.509 certificate chain) parameter contains a chain of one or more PKIX certificates [RFC5280].  The certificate chain is represented as a JSON array of certificate value strings.  Each string in the array is a base64-encoded (Section 4 of [RFC4648] -- not base64url-encoded) DER [ITU.X690.1994] PKIX certificate value. The PKIX certificate containing the key value MUST be the first certificate. |  [optional]
**y** | **String** |  |  [optional]



