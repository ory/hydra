# OidcConfiguration

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AuthorizationEndpoint** | **string** | OAuth 2.0 Authorization Endpoint URL | 
**BackchannelLogoutSessionSupported** | Pointer to **bool** | OpenID Connect Back-Channel Logout Session Required  Boolean value specifying whether the OP can pass a sid (session ID) Claim in the Logout Token to identify the RP session with the OP. If supported, the sid Claim is also included in ID Tokens issued by the OP | [optional] 
**BackchannelLogoutSupported** | Pointer to **bool** | OpenID Connect Back-Channel Logout Supported  Boolean value specifying whether the OP supports back-channel logout, with true indicating support. | [optional] 
**ClaimsParameterSupported** | Pointer to **bool** | OpenID Connect Claims Parameter Parameter Supported  Boolean value specifying whether the OP supports use of the claims parameter, with true indicating support. | [optional] 
**ClaimsSupported** | Pointer to **[]string** | OpenID Connect Supported Claims  JSON array containing a list of the Claim Names of the Claims that the OpenID Provider MAY be able to supply values for. Note that for privacy or other reasons, this might not be an exhaustive list. | [optional] 
**CodeChallengeMethodsSupported** | Pointer to **[]string** | OAuth 2.0 PKCE Supported Code Challenge Methods  JSON array containing a list of Proof Key for Code Exchange (PKCE) [RFC7636] code challenge methods supported by this authorization server. | [optional] 
**CredentialsEndpointDraft00** | Pointer to **string** | OpenID Connect Verifiable Credentials Endpoint  Contains the URL of the Verifiable Credentials Endpoint. | [optional] 
**CredentialsSupportedDraft00** | Pointer to [**[]CredentialSupportedDraft00**](CredentialSupportedDraft00.md) | OpenID Connect Verifiable Credentials Supported  JSON array containing a list of the Verifiable Credentials supported by this authorization server. | [optional] 
**DeviceAuthorizationEndpoint** | **string** | OAuth 2.0 Device Authorization Endpoint URL | 
**EndSessionEndpoint** | Pointer to **string** | OpenID Connect End-Session Endpoint  URL at the OP to which an RP can perform a redirect to request that the End-User be logged out at the OP. | [optional] 
**FrontchannelLogoutSessionSupported** | Pointer to **bool** | OpenID Connect Front-Channel Logout Session Required  Boolean value specifying whether the OP can pass iss (issuer) and sid (session ID) query parameters to identify the RP session with the OP when the frontchannel_logout_uri is used. If supported, the sid Claim is also included in ID Tokens issued by the OP. | [optional] 
**FrontchannelLogoutSupported** | Pointer to **bool** | OpenID Connect Front-Channel Logout Supported  Boolean value specifying whether the OP supports HTTP-based logout, with true indicating support. | [optional] 
**GrantTypesSupported** | Pointer to **[]string** | OAuth 2.0 Supported Grant Types  JSON array containing a list of the OAuth 2.0 Grant Type values that this OP supports. | [optional] 
**IdTokenSignedResponseAlg** | **[]string** | OpenID Connect Default ID Token Signing Algorithms  Algorithm used to sign OpenID Connect ID Tokens. | 
**IdTokenSigningAlgValuesSupported** | **[]string** | OpenID Connect Supported ID Token Signing Algorithms  JSON array containing a list of the JWS signing algorithms (alg values) supported by the OP for the ID Token to encode the Claims in a JWT. | 
**Issuer** | **string** | OpenID Connect Issuer URL  An URL using the https scheme with no query or fragment component that the OP asserts as its IssuerURL Identifier. If IssuerURL discovery is supported , this value MUST be identical to the issuer value returned by WebFinger. This also MUST be identical to the iss Claim value in ID Tokens issued from this IssuerURL. | 
**JwksUri** | **string** | OpenID Connect Well-Known JSON Web Keys URL  URL of the OP&#39;s JSON Web Key Set [JWK] document. This contains the signing key(s) the RP uses to validate signatures from the OP. The JWK Set MAY also contain the Server&#39;s encryption key(s), which are used by RPs to encrypt requests to the Server. When both signing and encryption keys are made available, a use (Key Use) parameter value is REQUIRED for all keys in the referenced JWK Set to indicate each key&#39;s intended usage. Although some algorithms allow the same key to be used for both signatures and encryption, doing so is NOT RECOMMENDED, as it is less secure. The JWK x5c parameter MAY be used to provide X.509 representations of keys provided. When used, the bare key values MUST still be present and MUST match those in the certificate. | 
**RegistrationEndpoint** | Pointer to **string** | OpenID Connect Dynamic Client Registration Endpoint URL | [optional] 
**RequestObjectSigningAlgValuesSupported** | Pointer to **[]string** | OpenID Connect Supported Request Object Signing Algorithms  JSON array containing a list of the JWS signing algorithms (alg values) supported by the OP for Request Objects, which are described in Section 6.1 of OpenID Connect Core 1.0 [OpenID.Core]. These algorithms are used both when the Request Object is passed by value (using the request parameter) and when it is passed by reference (using the request_uri parameter). | [optional] 
**RequestParameterSupported** | Pointer to **bool** | OpenID Connect Request Parameter Supported  Boolean value specifying whether the OP supports use of the request parameter, with true indicating support. | [optional] 
**RequestUriParameterSupported** | Pointer to **bool** | OpenID Connect Request URI Parameter Supported  Boolean value specifying whether the OP supports use of the request_uri parameter, with true indicating support. | [optional] 
**RequireRequestUriRegistration** | Pointer to **bool** | OpenID Connect Requires Request URI Registration  Boolean value specifying whether the OP requires any request_uri values used to be pre-registered using the request_uris registration parameter. | [optional] 
**ResponseModesSupported** | Pointer to **[]string** | OAuth 2.0 Supported Response Modes  JSON array containing a list of the OAuth 2.0 response_mode values that this OP supports. | [optional] 
**ResponseTypesSupported** | **[]string** | OAuth 2.0 Supported Response Types  JSON array containing a list of the OAuth 2.0 response_type values that this OP supports. Dynamic OpenID Providers MUST support the code, id_token, and the token id_token Response Type values. | 
**RevocationEndpoint** | Pointer to **string** | OAuth 2.0 Token Revocation URL  URL of the authorization server&#39;s OAuth 2.0 revocation endpoint. | [optional] 
**ScopesSupported** | Pointer to **[]string** | OAuth 2.0 Supported Scope Values  JSON array containing a list of the OAuth 2.0 [RFC6749] scope values that this server supports. The server MUST support the openid scope value. Servers MAY choose not to advertise some supported scope values even when this parameter is used | [optional] 
**SubjectTypesSupported** | **[]string** | OpenID Connect Supported Subject Types  JSON array containing a list of the Subject Identifier types that this OP supports. Valid types include pairwise and public. | 
**TokenEndpoint** | **string** | OAuth 2.0 Token Endpoint URL | 
**TokenEndpointAuthMethodsSupported** | Pointer to **[]string** | OAuth 2.0 Supported Client Authentication Methods  JSON array containing a list of Client Authentication methods supported by this Token Endpoint. The options are client_secret_post, client_secret_basic, client_secret_jwt, and private_key_jwt, as described in Section 9 of OpenID Connect Core 1.0 | [optional] 
**UserinfoEndpoint** | Pointer to **string** | OpenID Connect Userinfo URL  URL of the OP&#39;s UserInfo Endpoint. | [optional] 
**UserinfoSignedResponseAlg** | **[]string** | OpenID Connect User Userinfo Signing Algorithm  Algorithm used to sign OpenID Connect Userinfo Responses. | 
**UserinfoSigningAlgValuesSupported** | Pointer to **[]string** | OpenID Connect Supported Userinfo Signing Algorithm  JSON array containing a list of the JWS [JWS] signing algorithms (alg values) [JWA] supported by the UserInfo Endpoint to encode the Claims in a JWT [JWT]. | [optional] 

## Methods

### NewOidcConfiguration

`func NewOidcConfiguration(authorizationEndpoint string, deviceAuthorizationEndpoint string, idTokenSignedResponseAlg []string, idTokenSigningAlgValuesSupported []string, issuer string, jwksUri string, responseTypesSupported []string, subjectTypesSupported []string, tokenEndpoint string, userinfoSignedResponseAlg []string, ) *OidcConfiguration`

NewOidcConfiguration instantiates a new OidcConfiguration object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOidcConfigurationWithDefaults

`func NewOidcConfigurationWithDefaults() *OidcConfiguration`

NewOidcConfigurationWithDefaults instantiates a new OidcConfiguration object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAuthorizationEndpoint

`func (o *OidcConfiguration) GetAuthorizationEndpoint() string`

GetAuthorizationEndpoint returns the AuthorizationEndpoint field if non-nil, zero value otherwise.

### GetAuthorizationEndpointOk

`func (o *OidcConfiguration) GetAuthorizationEndpointOk() (*string, bool)`

GetAuthorizationEndpointOk returns a tuple with the AuthorizationEndpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthorizationEndpoint

`func (o *OidcConfiguration) SetAuthorizationEndpoint(v string)`

SetAuthorizationEndpoint sets AuthorizationEndpoint field to given value.


### GetBackchannelLogoutSessionSupported

`func (o *OidcConfiguration) GetBackchannelLogoutSessionSupported() bool`

GetBackchannelLogoutSessionSupported returns the BackchannelLogoutSessionSupported field if non-nil, zero value otherwise.

### GetBackchannelLogoutSessionSupportedOk

`func (o *OidcConfiguration) GetBackchannelLogoutSessionSupportedOk() (*bool, bool)`

GetBackchannelLogoutSessionSupportedOk returns a tuple with the BackchannelLogoutSessionSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBackchannelLogoutSessionSupported

`func (o *OidcConfiguration) SetBackchannelLogoutSessionSupported(v bool)`

SetBackchannelLogoutSessionSupported sets BackchannelLogoutSessionSupported field to given value.

### HasBackchannelLogoutSessionSupported

`func (o *OidcConfiguration) HasBackchannelLogoutSessionSupported() bool`

HasBackchannelLogoutSessionSupported returns a boolean if a field has been set.

### GetBackchannelLogoutSupported

`func (o *OidcConfiguration) GetBackchannelLogoutSupported() bool`

GetBackchannelLogoutSupported returns the BackchannelLogoutSupported field if non-nil, zero value otherwise.

### GetBackchannelLogoutSupportedOk

`func (o *OidcConfiguration) GetBackchannelLogoutSupportedOk() (*bool, bool)`

GetBackchannelLogoutSupportedOk returns a tuple with the BackchannelLogoutSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBackchannelLogoutSupported

`func (o *OidcConfiguration) SetBackchannelLogoutSupported(v bool)`

SetBackchannelLogoutSupported sets BackchannelLogoutSupported field to given value.

### HasBackchannelLogoutSupported

`func (o *OidcConfiguration) HasBackchannelLogoutSupported() bool`

HasBackchannelLogoutSupported returns a boolean if a field has been set.

### GetClaimsParameterSupported

`func (o *OidcConfiguration) GetClaimsParameterSupported() bool`

GetClaimsParameterSupported returns the ClaimsParameterSupported field if non-nil, zero value otherwise.

### GetClaimsParameterSupportedOk

`func (o *OidcConfiguration) GetClaimsParameterSupportedOk() (*bool, bool)`

GetClaimsParameterSupportedOk returns a tuple with the ClaimsParameterSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClaimsParameterSupported

`func (o *OidcConfiguration) SetClaimsParameterSupported(v bool)`

SetClaimsParameterSupported sets ClaimsParameterSupported field to given value.

### HasClaimsParameterSupported

`func (o *OidcConfiguration) HasClaimsParameterSupported() bool`

HasClaimsParameterSupported returns a boolean if a field has been set.

### GetClaimsSupported

`func (o *OidcConfiguration) GetClaimsSupported() []string`

GetClaimsSupported returns the ClaimsSupported field if non-nil, zero value otherwise.

### GetClaimsSupportedOk

`func (o *OidcConfiguration) GetClaimsSupportedOk() (*[]string, bool)`

GetClaimsSupportedOk returns a tuple with the ClaimsSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClaimsSupported

`func (o *OidcConfiguration) SetClaimsSupported(v []string)`

SetClaimsSupported sets ClaimsSupported field to given value.

### HasClaimsSupported

`func (o *OidcConfiguration) HasClaimsSupported() bool`

HasClaimsSupported returns a boolean if a field has been set.

### GetCodeChallengeMethodsSupported

`func (o *OidcConfiguration) GetCodeChallengeMethodsSupported() []string`

GetCodeChallengeMethodsSupported returns the CodeChallengeMethodsSupported field if non-nil, zero value otherwise.

### GetCodeChallengeMethodsSupportedOk

`func (o *OidcConfiguration) GetCodeChallengeMethodsSupportedOk() (*[]string, bool)`

GetCodeChallengeMethodsSupportedOk returns a tuple with the CodeChallengeMethodsSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCodeChallengeMethodsSupported

`func (o *OidcConfiguration) SetCodeChallengeMethodsSupported(v []string)`

SetCodeChallengeMethodsSupported sets CodeChallengeMethodsSupported field to given value.

### HasCodeChallengeMethodsSupported

`func (o *OidcConfiguration) HasCodeChallengeMethodsSupported() bool`

HasCodeChallengeMethodsSupported returns a boolean if a field has been set.

### GetCredentialsEndpointDraft00

`func (o *OidcConfiguration) GetCredentialsEndpointDraft00() string`

GetCredentialsEndpointDraft00 returns the CredentialsEndpointDraft00 field if non-nil, zero value otherwise.

### GetCredentialsEndpointDraft00Ok

`func (o *OidcConfiguration) GetCredentialsEndpointDraft00Ok() (*string, bool)`

GetCredentialsEndpointDraft00Ok returns a tuple with the CredentialsEndpointDraft00 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCredentialsEndpointDraft00

`func (o *OidcConfiguration) SetCredentialsEndpointDraft00(v string)`

SetCredentialsEndpointDraft00 sets CredentialsEndpointDraft00 field to given value.

### HasCredentialsEndpointDraft00

`func (o *OidcConfiguration) HasCredentialsEndpointDraft00() bool`

HasCredentialsEndpointDraft00 returns a boolean if a field has been set.

### GetCredentialsSupportedDraft00

`func (o *OidcConfiguration) GetCredentialsSupportedDraft00() []CredentialSupportedDraft00`

GetCredentialsSupportedDraft00 returns the CredentialsSupportedDraft00 field if non-nil, zero value otherwise.

### GetCredentialsSupportedDraft00Ok

`func (o *OidcConfiguration) GetCredentialsSupportedDraft00Ok() (*[]CredentialSupportedDraft00, bool)`

GetCredentialsSupportedDraft00Ok returns a tuple with the CredentialsSupportedDraft00 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCredentialsSupportedDraft00

`func (o *OidcConfiguration) SetCredentialsSupportedDraft00(v []CredentialSupportedDraft00)`

SetCredentialsSupportedDraft00 sets CredentialsSupportedDraft00 field to given value.

### HasCredentialsSupportedDraft00

`func (o *OidcConfiguration) HasCredentialsSupportedDraft00() bool`

HasCredentialsSupportedDraft00 returns a boolean if a field has been set.

### GetDeviceAuthorizationEndpoint

`func (o *OidcConfiguration) GetDeviceAuthorizationEndpoint() string`

GetDeviceAuthorizationEndpoint returns the DeviceAuthorizationEndpoint field if non-nil, zero value otherwise.

### GetDeviceAuthorizationEndpointOk

`func (o *OidcConfiguration) GetDeviceAuthorizationEndpointOk() (*string, bool)`

GetDeviceAuthorizationEndpointOk returns a tuple with the DeviceAuthorizationEndpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceAuthorizationEndpoint

`func (o *OidcConfiguration) SetDeviceAuthorizationEndpoint(v string)`

SetDeviceAuthorizationEndpoint sets DeviceAuthorizationEndpoint field to given value.


### GetEndSessionEndpoint

`func (o *OidcConfiguration) GetEndSessionEndpoint() string`

GetEndSessionEndpoint returns the EndSessionEndpoint field if non-nil, zero value otherwise.

### GetEndSessionEndpointOk

`func (o *OidcConfiguration) GetEndSessionEndpointOk() (*string, bool)`

GetEndSessionEndpointOk returns a tuple with the EndSessionEndpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndSessionEndpoint

`func (o *OidcConfiguration) SetEndSessionEndpoint(v string)`

SetEndSessionEndpoint sets EndSessionEndpoint field to given value.

### HasEndSessionEndpoint

`func (o *OidcConfiguration) HasEndSessionEndpoint() bool`

HasEndSessionEndpoint returns a boolean if a field has been set.

### GetFrontchannelLogoutSessionSupported

`func (o *OidcConfiguration) GetFrontchannelLogoutSessionSupported() bool`

GetFrontchannelLogoutSessionSupported returns the FrontchannelLogoutSessionSupported field if non-nil, zero value otherwise.

### GetFrontchannelLogoutSessionSupportedOk

`func (o *OidcConfiguration) GetFrontchannelLogoutSessionSupportedOk() (*bool, bool)`

GetFrontchannelLogoutSessionSupportedOk returns a tuple with the FrontchannelLogoutSessionSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFrontchannelLogoutSessionSupported

`func (o *OidcConfiguration) SetFrontchannelLogoutSessionSupported(v bool)`

SetFrontchannelLogoutSessionSupported sets FrontchannelLogoutSessionSupported field to given value.

### HasFrontchannelLogoutSessionSupported

`func (o *OidcConfiguration) HasFrontchannelLogoutSessionSupported() bool`

HasFrontchannelLogoutSessionSupported returns a boolean if a field has been set.

### GetFrontchannelLogoutSupported

`func (o *OidcConfiguration) GetFrontchannelLogoutSupported() bool`

GetFrontchannelLogoutSupported returns the FrontchannelLogoutSupported field if non-nil, zero value otherwise.

### GetFrontchannelLogoutSupportedOk

`func (o *OidcConfiguration) GetFrontchannelLogoutSupportedOk() (*bool, bool)`

GetFrontchannelLogoutSupportedOk returns a tuple with the FrontchannelLogoutSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFrontchannelLogoutSupported

`func (o *OidcConfiguration) SetFrontchannelLogoutSupported(v bool)`

SetFrontchannelLogoutSupported sets FrontchannelLogoutSupported field to given value.

### HasFrontchannelLogoutSupported

`func (o *OidcConfiguration) HasFrontchannelLogoutSupported() bool`

HasFrontchannelLogoutSupported returns a boolean if a field has been set.

### GetGrantTypesSupported

`func (o *OidcConfiguration) GetGrantTypesSupported() []string`

GetGrantTypesSupported returns the GrantTypesSupported field if non-nil, zero value otherwise.

### GetGrantTypesSupportedOk

`func (o *OidcConfiguration) GetGrantTypesSupportedOk() (*[]string, bool)`

GetGrantTypesSupportedOk returns a tuple with the GrantTypesSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGrantTypesSupported

`func (o *OidcConfiguration) SetGrantTypesSupported(v []string)`

SetGrantTypesSupported sets GrantTypesSupported field to given value.

### HasGrantTypesSupported

`func (o *OidcConfiguration) HasGrantTypesSupported() bool`

HasGrantTypesSupported returns a boolean if a field has been set.

### GetIdTokenSignedResponseAlg

`func (o *OidcConfiguration) GetIdTokenSignedResponseAlg() []string`

GetIdTokenSignedResponseAlg returns the IdTokenSignedResponseAlg field if non-nil, zero value otherwise.

### GetIdTokenSignedResponseAlgOk

`func (o *OidcConfiguration) GetIdTokenSignedResponseAlgOk() (*[]string, bool)`

GetIdTokenSignedResponseAlgOk returns a tuple with the IdTokenSignedResponseAlg field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIdTokenSignedResponseAlg

`func (o *OidcConfiguration) SetIdTokenSignedResponseAlg(v []string)`

SetIdTokenSignedResponseAlg sets IdTokenSignedResponseAlg field to given value.


### GetIdTokenSigningAlgValuesSupported

`func (o *OidcConfiguration) GetIdTokenSigningAlgValuesSupported() []string`

GetIdTokenSigningAlgValuesSupported returns the IdTokenSigningAlgValuesSupported field if non-nil, zero value otherwise.

### GetIdTokenSigningAlgValuesSupportedOk

`func (o *OidcConfiguration) GetIdTokenSigningAlgValuesSupportedOk() (*[]string, bool)`

GetIdTokenSigningAlgValuesSupportedOk returns a tuple with the IdTokenSigningAlgValuesSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIdTokenSigningAlgValuesSupported

`func (o *OidcConfiguration) SetIdTokenSigningAlgValuesSupported(v []string)`

SetIdTokenSigningAlgValuesSupported sets IdTokenSigningAlgValuesSupported field to given value.


### GetIssuer

`func (o *OidcConfiguration) GetIssuer() string`

GetIssuer returns the Issuer field if non-nil, zero value otherwise.

### GetIssuerOk

`func (o *OidcConfiguration) GetIssuerOk() (*string, bool)`

GetIssuerOk returns a tuple with the Issuer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIssuer

`func (o *OidcConfiguration) SetIssuer(v string)`

SetIssuer sets Issuer field to given value.


### GetJwksUri

`func (o *OidcConfiguration) GetJwksUri() string`

GetJwksUri returns the JwksUri field if non-nil, zero value otherwise.

### GetJwksUriOk

`func (o *OidcConfiguration) GetJwksUriOk() (*string, bool)`

GetJwksUriOk returns a tuple with the JwksUri field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetJwksUri

`func (o *OidcConfiguration) SetJwksUri(v string)`

SetJwksUri sets JwksUri field to given value.


### GetRegistrationEndpoint

`func (o *OidcConfiguration) GetRegistrationEndpoint() string`

GetRegistrationEndpoint returns the RegistrationEndpoint field if non-nil, zero value otherwise.

### GetRegistrationEndpointOk

`func (o *OidcConfiguration) GetRegistrationEndpointOk() (*string, bool)`

GetRegistrationEndpointOk returns a tuple with the RegistrationEndpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegistrationEndpoint

`func (o *OidcConfiguration) SetRegistrationEndpoint(v string)`

SetRegistrationEndpoint sets RegistrationEndpoint field to given value.

### HasRegistrationEndpoint

`func (o *OidcConfiguration) HasRegistrationEndpoint() bool`

HasRegistrationEndpoint returns a boolean if a field has been set.

### GetRequestObjectSigningAlgValuesSupported

`func (o *OidcConfiguration) GetRequestObjectSigningAlgValuesSupported() []string`

GetRequestObjectSigningAlgValuesSupported returns the RequestObjectSigningAlgValuesSupported field if non-nil, zero value otherwise.

### GetRequestObjectSigningAlgValuesSupportedOk

`func (o *OidcConfiguration) GetRequestObjectSigningAlgValuesSupportedOk() (*[]string, bool)`

GetRequestObjectSigningAlgValuesSupportedOk returns a tuple with the RequestObjectSigningAlgValuesSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestObjectSigningAlgValuesSupported

`func (o *OidcConfiguration) SetRequestObjectSigningAlgValuesSupported(v []string)`

SetRequestObjectSigningAlgValuesSupported sets RequestObjectSigningAlgValuesSupported field to given value.

### HasRequestObjectSigningAlgValuesSupported

`func (o *OidcConfiguration) HasRequestObjectSigningAlgValuesSupported() bool`

HasRequestObjectSigningAlgValuesSupported returns a boolean if a field has been set.

### GetRequestParameterSupported

`func (o *OidcConfiguration) GetRequestParameterSupported() bool`

GetRequestParameterSupported returns the RequestParameterSupported field if non-nil, zero value otherwise.

### GetRequestParameterSupportedOk

`func (o *OidcConfiguration) GetRequestParameterSupportedOk() (*bool, bool)`

GetRequestParameterSupportedOk returns a tuple with the RequestParameterSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestParameterSupported

`func (o *OidcConfiguration) SetRequestParameterSupported(v bool)`

SetRequestParameterSupported sets RequestParameterSupported field to given value.

### HasRequestParameterSupported

`func (o *OidcConfiguration) HasRequestParameterSupported() bool`

HasRequestParameterSupported returns a boolean if a field has been set.

### GetRequestUriParameterSupported

`func (o *OidcConfiguration) GetRequestUriParameterSupported() bool`

GetRequestUriParameterSupported returns the RequestUriParameterSupported field if non-nil, zero value otherwise.

### GetRequestUriParameterSupportedOk

`func (o *OidcConfiguration) GetRequestUriParameterSupportedOk() (*bool, bool)`

GetRequestUriParameterSupportedOk returns a tuple with the RequestUriParameterSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestUriParameterSupported

`func (o *OidcConfiguration) SetRequestUriParameterSupported(v bool)`

SetRequestUriParameterSupported sets RequestUriParameterSupported field to given value.

### HasRequestUriParameterSupported

`func (o *OidcConfiguration) HasRequestUriParameterSupported() bool`

HasRequestUriParameterSupported returns a boolean if a field has been set.

### GetRequireRequestUriRegistration

`func (o *OidcConfiguration) GetRequireRequestUriRegistration() bool`

GetRequireRequestUriRegistration returns the RequireRequestUriRegistration field if non-nil, zero value otherwise.

### GetRequireRequestUriRegistrationOk

`func (o *OidcConfiguration) GetRequireRequestUriRegistrationOk() (*bool, bool)`

GetRequireRequestUriRegistrationOk returns a tuple with the RequireRequestUriRegistration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequireRequestUriRegistration

`func (o *OidcConfiguration) SetRequireRequestUriRegistration(v bool)`

SetRequireRequestUriRegistration sets RequireRequestUriRegistration field to given value.

### HasRequireRequestUriRegistration

`func (o *OidcConfiguration) HasRequireRequestUriRegistration() bool`

HasRequireRequestUriRegistration returns a boolean if a field has been set.

### GetResponseModesSupported

`func (o *OidcConfiguration) GetResponseModesSupported() []string`

GetResponseModesSupported returns the ResponseModesSupported field if non-nil, zero value otherwise.

### GetResponseModesSupportedOk

`func (o *OidcConfiguration) GetResponseModesSupportedOk() (*[]string, bool)`

GetResponseModesSupportedOk returns a tuple with the ResponseModesSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResponseModesSupported

`func (o *OidcConfiguration) SetResponseModesSupported(v []string)`

SetResponseModesSupported sets ResponseModesSupported field to given value.

### HasResponseModesSupported

`func (o *OidcConfiguration) HasResponseModesSupported() bool`

HasResponseModesSupported returns a boolean if a field has been set.

### GetResponseTypesSupported

`func (o *OidcConfiguration) GetResponseTypesSupported() []string`

GetResponseTypesSupported returns the ResponseTypesSupported field if non-nil, zero value otherwise.

### GetResponseTypesSupportedOk

`func (o *OidcConfiguration) GetResponseTypesSupportedOk() (*[]string, bool)`

GetResponseTypesSupportedOk returns a tuple with the ResponseTypesSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResponseTypesSupported

`func (o *OidcConfiguration) SetResponseTypesSupported(v []string)`

SetResponseTypesSupported sets ResponseTypesSupported field to given value.


### GetRevocationEndpoint

`func (o *OidcConfiguration) GetRevocationEndpoint() string`

GetRevocationEndpoint returns the RevocationEndpoint field if non-nil, zero value otherwise.

### GetRevocationEndpointOk

`func (o *OidcConfiguration) GetRevocationEndpointOk() (*string, bool)`

GetRevocationEndpointOk returns a tuple with the RevocationEndpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRevocationEndpoint

`func (o *OidcConfiguration) SetRevocationEndpoint(v string)`

SetRevocationEndpoint sets RevocationEndpoint field to given value.

### HasRevocationEndpoint

`func (o *OidcConfiguration) HasRevocationEndpoint() bool`

HasRevocationEndpoint returns a boolean if a field has been set.

### GetScopesSupported

`func (o *OidcConfiguration) GetScopesSupported() []string`

GetScopesSupported returns the ScopesSupported field if non-nil, zero value otherwise.

### GetScopesSupportedOk

`func (o *OidcConfiguration) GetScopesSupportedOk() (*[]string, bool)`

GetScopesSupportedOk returns a tuple with the ScopesSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScopesSupported

`func (o *OidcConfiguration) SetScopesSupported(v []string)`

SetScopesSupported sets ScopesSupported field to given value.

### HasScopesSupported

`func (o *OidcConfiguration) HasScopesSupported() bool`

HasScopesSupported returns a boolean if a field has been set.

### GetSubjectTypesSupported

`func (o *OidcConfiguration) GetSubjectTypesSupported() []string`

GetSubjectTypesSupported returns the SubjectTypesSupported field if non-nil, zero value otherwise.

### GetSubjectTypesSupportedOk

`func (o *OidcConfiguration) GetSubjectTypesSupportedOk() (*[]string, bool)`

GetSubjectTypesSupportedOk returns a tuple with the SubjectTypesSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubjectTypesSupported

`func (o *OidcConfiguration) SetSubjectTypesSupported(v []string)`

SetSubjectTypesSupported sets SubjectTypesSupported field to given value.


### GetTokenEndpoint

`func (o *OidcConfiguration) GetTokenEndpoint() string`

GetTokenEndpoint returns the TokenEndpoint field if non-nil, zero value otherwise.

### GetTokenEndpointOk

`func (o *OidcConfiguration) GetTokenEndpointOk() (*string, bool)`

GetTokenEndpointOk returns a tuple with the TokenEndpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTokenEndpoint

`func (o *OidcConfiguration) SetTokenEndpoint(v string)`

SetTokenEndpoint sets TokenEndpoint field to given value.


### GetTokenEndpointAuthMethodsSupported

`func (o *OidcConfiguration) GetTokenEndpointAuthMethodsSupported() []string`

GetTokenEndpointAuthMethodsSupported returns the TokenEndpointAuthMethodsSupported field if non-nil, zero value otherwise.

### GetTokenEndpointAuthMethodsSupportedOk

`func (o *OidcConfiguration) GetTokenEndpointAuthMethodsSupportedOk() (*[]string, bool)`

GetTokenEndpointAuthMethodsSupportedOk returns a tuple with the TokenEndpointAuthMethodsSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTokenEndpointAuthMethodsSupported

`func (o *OidcConfiguration) SetTokenEndpointAuthMethodsSupported(v []string)`

SetTokenEndpointAuthMethodsSupported sets TokenEndpointAuthMethodsSupported field to given value.

### HasTokenEndpointAuthMethodsSupported

`func (o *OidcConfiguration) HasTokenEndpointAuthMethodsSupported() bool`

HasTokenEndpointAuthMethodsSupported returns a boolean if a field has been set.

### GetUserinfoEndpoint

`func (o *OidcConfiguration) GetUserinfoEndpoint() string`

GetUserinfoEndpoint returns the UserinfoEndpoint field if non-nil, zero value otherwise.

### GetUserinfoEndpointOk

`func (o *OidcConfiguration) GetUserinfoEndpointOk() (*string, bool)`

GetUserinfoEndpointOk returns a tuple with the UserinfoEndpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserinfoEndpoint

`func (o *OidcConfiguration) SetUserinfoEndpoint(v string)`

SetUserinfoEndpoint sets UserinfoEndpoint field to given value.

### HasUserinfoEndpoint

`func (o *OidcConfiguration) HasUserinfoEndpoint() bool`

HasUserinfoEndpoint returns a boolean if a field has been set.

### GetUserinfoSignedResponseAlg

`func (o *OidcConfiguration) GetUserinfoSignedResponseAlg() []string`

GetUserinfoSignedResponseAlg returns the UserinfoSignedResponseAlg field if non-nil, zero value otherwise.

### GetUserinfoSignedResponseAlgOk

`func (o *OidcConfiguration) GetUserinfoSignedResponseAlgOk() (*[]string, bool)`

GetUserinfoSignedResponseAlgOk returns a tuple with the UserinfoSignedResponseAlg field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserinfoSignedResponseAlg

`func (o *OidcConfiguration) SetUserinfoSignedResponseAlg(v []string)`

SetUserinfoSignedResponseAlg sets UserinfoSignedResponseAlg field to given value.


### GetUserinfoSigningAlgValuesSupported

`func (o *OidcConfiguration) GetUserinfoSigningAlgValuesSupported() []string`

GetUserinfoSigningAlgValuesSupported returns the UserinfoSigningAlgValuesSupported field if non-nil, zero value otherwise.

### GetUserinfoSigningAlgValuesSupportedOk

`func (o *OidcConfiguration) GetUserinfoSigningAlgValuesSupportedOk() (*[]string, bool)`

GetUserinfoSigningAlgValuesSupportedOk returns a tuple with the UserinfoSigningAlgValuesSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserinfoSigningAlgValuesSupported

`func (o *OidcConfiguration) SetUserinfoSigningAlgValuesSupported(v []string)`

SetUserinfoSigningAlgValuesSupported sets UserinfoSigningAlgValuesSupported field to given value.

### HasUserinfoSigningAlgValuesSupported

`func (o *OidcConfiguration) HasUserinfoSigningAlgValuesSupported() bool`

HasUserinfoSigningAlgValuesSupported returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


