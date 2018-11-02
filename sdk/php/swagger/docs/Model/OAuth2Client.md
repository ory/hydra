# OAuth2Client

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**allowed_cors_origins** | **string[]** | AllowedCORSOrigins are one or more URLs (scheme://host[:port]) which are allowed to make CORS requests to the /oauth/token endpoint. If this array is empty, the sever&#39;s CORS origin configuration (&#x60;CORS_ALLOWED_ORIGINS&#x60;) will be used instead. If this array is set, the allowed origins are appended to the server&#39;s CORS origin configuration. Be aware that environment variable &#x60;CORS_ENABLED&#x60; MUST be set to &#x60;true&#x60; for this to work. | [optional] 
**audience** | **string[]** | Audience is a whitelist defining the audiences this client is allowed to request tokens for. An audience limits the applicability of an OAuth 2.0 Access Token to, for example, certain API endpoints. The value is a list of URLs. URLs MUST NOT contain whitespaces. | [optional] 
**client_id** | **string** | ClientID  is the id for this client. | [optional] 
**client_name** | **string** | Name is the human-readable string name of the client to be presented to the end-user during authorization. | [optional] 
**client_secret** | **string** | Secret is the client&#39;s secret. The secret will be included in the create request as cleartext, and then never again. The secret is stored using BCrypt so it is impossible to recover it. Tell your users that they need to write the secret down as it will not be made available again. | [optional] 
**client_secret_expires_at** | **int** | SecretExpiresAt is an integer holding the time at which the client secret will expire or 0 if it will not expire. The time is represented as the number of seconds from 1970-01-01T00:00:00Z as measured in UTC until the date/time of expiration. | [optional] 
**client_uri** | **string** | ClientURI is an URL string of a web page providing information about the client. If present, the server SHOULD display this URL to the end-user in a clickable fashion. | [optional] 
**contacts** | **string[]** | Contacts is a array of strings representing ways to contact people responsible for this client, typically email addresses. | [optional] 
**grant_types** | **string[]** | GrantTypes is an array of grant types the client is allowed to use. | [optional] 
**jwks** | [**\HydraSDK\Model\JSONWebKeySet**](JSONWebKeySet.md) |  | [optional] 
**jwks_uri** | **string** | URL for the Client&#39;s JSON Web Key Set [JWK] document. If the Client signs requests to the Server, it contains the signing key(s) the Server uses to validate signatures from the Client. The JWK Set MAY also contain the Client&#39;s encryption keys(s), which are used by the Server to encrypt responses to the Client. When both signing and encryption keys are made available, a use (Key Use) parameter value is REQUIRED for all keys in the referenced JWK Set to indicate each key&#39;s intended usage. Although some algorithms allow the same key to be used for both signatures and encryption, doing so is NOT RECOMMENDED, as it is less secure. The JWK x5c parameter MAY be used to provide X.509 representations of keys provided. When used, the bare key values MUST still be present and MUST match those in the certificate. | [optional] 
**logo_uri** | **string** | LogoURI is an URL string that references a logo for the client. | [optional] 
**owner** | **string** | Owner is a string identifying the owner of the OAuth 2.0 Client. | [optional] 
**policy_uri** | **string** | PolicyURI is a URL string that points to a human-readable privacy policy document that describes how the deployment organization collects, uses, retains, and discloses personal data. | [optional] 
**redirect_uris** | **string[]** | RedirectURIs is an array of allowed redirect urls for the client, for example http://mydomain/oauth/callback . | [optional] 
**request_object_signing_alg** | **string** | JWS [JWS] alg algorithm [JWA] that MUST be used for signing Request Objects sent to the OP. All Request Objects from this Client MUST be rejected, if not signed with this algorithm. | [optional] 
**request_uris** | **string[]** | Array of request_uri values that are pre-registered by the RP for use at the OP. Servers MAY cache the contents of the files referenced by these URIs and not retrieve them at the time they are used in a request. OPs can require that request_uri values used be pre-registered with the require_request_uri_registration discovery parameter. | [optional] 
**response_types** | **string[]** | ResponseTypes is an array of the OAuth 2.0 response type strings that the client can use at the authorization endpoint. | [optional] 
**scope** | **string** | Scope is a string containing a space-separated list of scope values (as described in Section 3.3 of OAuth 2.0 [RFC6749]) that the client can use when requesting access tokens. | [optional] 
**sector_identifier_uri** | **string** | URL using the https scheme to be used in calculating Pseudonymous Identifiers by the OP. The URL references a file with a single JSON array of redirect_uri values. | [optional] 
**subject_type** | **string** | SubjectType requested for responses to this Client. The subject_types_supported Discovery parameter contains a list of the supported subject_type values for this server. Valid types include &#x60;pairwise&#x60; and &#x60;public&#x60;. | [optional] 
**token_endpoint_auth_method** | **string** | Requested Client Authentication method for the Token Endpoint. The options are client_secret_post, client_secret_basic, private_key_jwt, and none. | [optional] 
**tos_uri** | **string** | TermsOfServiceURI is a URL string that points to a human-readable terms of service document for the client that describes a contractual relationship between the end-user and the client that the end-user accepts when authorizing the client. | [optional] 
**userinfo_signed_response_alg** | **string** | JWS alg algorithm [JWA] REQUIRED for signing UserInfo Responses. If this is specified, the response will be JWT [JWT] serialized, and signed using JWS. The default, if omitted, is for the UserInfo Response to return the Claims as a UTF-8 encoded JSON object using the application/json content-type. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


