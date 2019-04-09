# WellKnown

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AuthorizationEndpoint** | **string** | URL of the OP&#39;s OAuth 2.0 Authorization Endpoint. | [default to null]
**ClaimsParameterSupported** | **bool** | Boolean value specifying whether the OP supports use of the claims parameter, with true indicating support. | [optional] [default to null]
**ClaimsSupported** | **[]string** | JSON array containing a list of the Claim Names of the Claims that the OpenID Provider MAY be able to supply values for. Note that for privacy or other reasons, this might not be an exhaustive list. | [optional] [default to null]
**GrantTypesSupported** | **[]string** | JSON array containing a list of the OAuth 2.0 Grant Type values that this OP supports. | [optional] [default to null]
**IdTokenSigningAlgValuesSupported** | **[]string** | JSON array containing a list of the JWS signing algorithms (alg values) supported by the OP for the ID Token to encode the Claims in a JWT. | [default to null]
**Issuer** | **string** | URL using the https scheme with no query or fragment component that the OP asserts as its IssuerURL Identifier. If IssuerURL discovery is supported , this value MUST be identical to the issuer value returned by WebFinger. This also MUST be identical to the iss Claim value in ID Tokens issued from this IssuerURL. | [default to null]
**JwksUri** | **string** | URL of the OP&#39;s JSON Web Key Set [JWK] document. This contains the signing key(s) the RP uses to validate signatures from the OP. The JWK Set MAY also contain the Server&#39;s encryption key(s), which are used by RPs to encrypt requests to the Server. When both signing and encryption keys are made available, a use (Key Use) parameter value is REQUIRED for all keys in the referenced JWK Set to indicate each key&#39;s intended usage. Although some algorithms allow the same key to be used for both signatures and encryption, doing so is NOT RECOMMENDED, as it is less secure. The JWK x5c parameter MAY be used to provide X.509 representations of keys provided. When used, the bare key values MUST still be present and MUST match those in the certificate. | [default to null]
**RegistrationEndpoint** | **string** | URL of the OP&#39;s Dynamic Client Registration Endpoint. | [optional] [default to null]
**RequestParameterSupported** | **bool** | Boolean value specifying whether the OP supports use of the request parameter, with true indicating support. | [optional] [default to null]
**RequestUriParameterSupported** | **bool** | Boolean value specifying whether the OP supports use of the request_uri parameter, with true indicating support. | [optional] [default to null]
**RequireRequestUriRegistration** | **bool** | Boolean value specifying whether the OP requires any request_uri values used to be pre-registered using the request_uris registration parameter. | [optional] [default to null]
**ResponseModesSupported** | **[]string** | JSON array containing a list of the OAuth 2.0 response_mode values that this OP supports. | [optional] [default to null]
**ResponseTypesSupported** | **[]string** | JSON array containing a list of the OAuth 2.0 response_type values that this OP supports. Dynamic OpenID Providers MUST support the code, id_token, and the token id_token Response Type values. | [default to null]
**ScopesSupported** | **[]string** | SON array containing a list of the OAuth 2.0 [RFC6749] scope values that this server supports. The server MUST support the openid scope value. Servers MAY choose not to advertise some supported scope values even when this parameter is used | [optional] [default to null]
**SubjectTypesSupported** | **[]string** | JSON array containing a list of the Subject Identifier types that this OP supports. Valid types include pairwise and public. | [default to null]
**TokenEndpoint** | **string** | URL of the OP&#39;s OAuth 2.0 Token Endpoint | [default to null]
**TokenEndpointAuthMethodsSupported** | **[]string** | JSON array containing a list of Client Authentication methods supported by this Token Endpoint. The options are client_secret_post, client_secret_basic, client_secret_jwt, and private_key_jwt, as described in Section 9 of OpenID Connect Core 1.0 | [optional] [default to null]
**UserinfoEndpoint** | **string** | URL of the OP&#39;s UserInfo Endpoint. | [optional] [default to null]
**UserinfoSigningAlgValuesSupported** | **[]string** | JSON array containing a list of the JWS [JWS] signing algorithms (alg values) [JWA] supported by the UserInfo Endpoint to encode the Claims in a JWT [JWT]. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


