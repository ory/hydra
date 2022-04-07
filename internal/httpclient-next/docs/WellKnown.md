# WellKnown

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AuthorizationEndpoint** | **string** | URL of the OP&#39;s OAuth 2.0 Authorization Endpoint. | 
**BackchannelLogoutSessionSupported** | Pointer to **bool** | Boolean value specifying whether the OP can pass a sid (session ID) Claim in the Logout Token to identify the RP session with the OP. If supported, the sid Claim is also included in ID Tokens issued by the OP | [optional] 
**BackchannelLogoutSupported** | Pointer to **bool** | Boolean value specifying whether the OP supports back-channel logout, with true indicating support. | [optional] 
**ClaimsParameterSupported** | Pointer to **bool** | Boolean value specifying whether the OP supports use of the claims parameter, with true indicating support. | [optional] 
**ClaimsSupported** | Pointer to **[]string** | JSON array containing a list of the Claim Names of the Claims that the OpenID Provider MAY be able to supply values for. Note that for privacy or other reasons, this might not be an exhaustive list. | [optional] 
**CodeChallengeMethodsSupported** | Pointer to **[]string** | JSON array containing a list of Proof Key for Code Exchange (PKCE) [RFC7636] code challenge methods supported by this authorization server. | [optional] 
**EndSessionEndpoint** | Pointer to **string** | URL at the OP to which an RP can perform a redirect to request that the End-User be logged out at the OP. | [optional] 
**FrontchannelLogoutSessionSupported** | Pointer to **bool** | Boolean value specifying whether the OP can pass iss (issuer) and sid (session ID) query parameters to identify the RP session with the OP when the frontchannel_logout_uri is used. If supported, the sid Claim is also included in ID Tokens issued by the OP. | [optional] 
**FrontchannelLogoutSupported** | Pointer to **bool** | Boolean value specifying whether the OP supports HTTP-based logout, with true indicating support. | [optional] 
**GrantTypesSupported** | Pointer to **[]string** | JSON array containing a list of the OAuth 2.0 Grant Type values that this OP supports. | [optional] 
**IdTokenSigningAlgValuesSupported** | **[]string** | JSON array containing a list of the JWS signing algorithms (alg values) supported by the OP for the ID Token to encode the Claims in a JWT. | 
**Issuer** | **string** | URL using the https scheme with no query or fragment component that the OP asserts as its IssuerURL Identifier. If IssuerURL discovery is supported , this value MUST be identical to the issuer value returned by WebFinger. This also MUST be identical to the iss Claim value in ID Tokens issued from this IssuerURL. | 
**JwksUri** | **string** | URL of the OP&#39;s JSON Web Key Set [JWK] document. This contains the signing key(s) the RP uses to validate signatures from the OP. The JWK Set MAY also contain the Server&#39;s encryption key(s), which are used by RPs to encrypt requests to the Server. When both signing and encryption keys are made available, a use (Key Use) parameter value is REQUIRED for all keys in the referenced JWK Set to indicate each key&#39;s intended usage. Although some algorithms allow the same key to be used for both signatures and encryption, doing so is NOT RECOMMENDED, as it is less secure. The JWK x5c parameter MAY be used to provide X.509 representations of keys provided. When used, the bare key values MUST still be present and MUST match those in the certificate. | 
**RegistrationEndpoint** | Pointer to **string** | URL of the OP&#39;s Dynamic Client Registration Endpoint. | [optional] 
**RequestObjectSigningAlgValuesSupported** | Pointer to **[]string** | JSON array containing a list of the JWS signing algorithms (alg values) supported by the OP for Request Objects, which are described in Section 6.1 of OpenID Connect Core 1.0 [OpenID.Core]. These algorithms are used both when the Request Object is passed by value (using the request parameter) and when it is passed by reference (using the request_uri parameter). | [optional] 
**RequestParameterSupported** | Pointer to **bool** | Boolean value specifying whether the OP supports use of the request parameter, with true indicating support. | [optional] 
**RequestUriParameterSupported** | Pointer to **bool** | Boolean value specifying whether the OP supports use of the request_uri parameter, with true indicating support. | [optional] 
**RequireRequestUriRegistration** | Pointer to **bool** | Boolean value specifying whether the OP requires any request_uri values used to be pre-registered using the request_uris registration parameter. | [optional] 
**ResponseModesSupported** | Pointer to **[]string** | JSON array containing a list of the OAuth 2.0 response_mode values that this OP supports. | [optional] 
**ResponseTypesSupported** | **[]string** | JSON array containing a list of the OAuth 2.0 response_type values that this OP supports. Dynamic OpenID Providers MUST support the code, id_token, and the token id_token Response Type values. | 
**RevocationEndpoint** | Pointer to **string** | URL of the authorization server&#39;s OAuth 2.0 revocation endpoint. | [optional] 
**ScopesSupported** | Pointer to **[]string** | SON array containing a list of the OAuth 2.0 [RFC6749] scope values that this server supports. The server MUST support the openid scope value. Servers MAY choose not to advertise some supported scope values even when this parameter is used | [optional] 
**SubjectTypesSupported** | **[]string** | JSON array containing a list of the Subject Identifier types that this OP supports. Valid types include pairwise and public. | 
**TokenEndpoint** | **string** | URL of the OP&#39;s OAuth 2.0 Token Endpoint | 
**TokenEndpointAuthMethodsSupported** | Pointer to **[]string** | JSON array containing a list of Client Authentication methods supported by this Token Endpoint. The options are client_secret_post, client_secret_basic, client_secret_jwt, and private_key_jwt, as described in Section 9 of OpenID Connect Core 1.0 | [optional] 
**UserinfoEndpoint** | Pointer to **string** | URL of the OP&#39;s UserInfo Endpoint. | [optional] 
**UserinfoSigningAlgValuesSupported** | Pointer to **[]string** | JSON array containing a list of the JWS [JWS] signing algorithms (alg values) [JWA] supported by the UserInfo Endpoint to encode the Claims in a JWT [JWT]. | [optional] 

## Methods

### NewWellKnown

`func NewWellKnown(authorizationEndpoint string, idTokenSigningAlgValuesSupported []string, issuer string, jwksUri string, responseTypesSupported []string, subjectTypesSupported []string, tokenEndpoint string, ) *WellKnown`

NewWellKnown instantiates a new WellKnown object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWellKnownWithDefaults

`func NewWellKnownWithDefaults() *WellKnown`

NewWellKnownWithDefaults instantiates a new WellKnown object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAuthorizationEndpoint

`func (o *WellKnown) GetAuthorizationEndpoint() string`

GetAuthorizationEndpoint returns the AuthorizationEndpoint field if non-nil, zero value otherwise.

### GetAuthorizationEndpointOk

`func (o *WellKnown) GetAuthorizationEndpointOk() (*string, bool)`

GetAuthorizationEndpointOk returns a tuple with the AuthorizationEndpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthorizationEndpoint

`func (o *WellKnown) SetAuthorizationEndpoint(v string)`

SetAuthorizationEndpoint sets AuthorizationEndpoint field to given value.


### GetBackchannelLogoutSessionSupported

`func (o *WellKnown) GetBackchannelLogoutSessionSupported() bool`

GetBackchannelLogoutSessionSupported returns the BackchannelLogoutSessionSupported field if non-nil, zero value otherwise.

### GetBackchannelLogoutSessionSupportedOk

`func (o *WellKnown) GetBackchannelLogoutSessionSupportedOk() (*bool, bool)`

GetBackchannelLogoutSessionSupportedOk returns a tuple with the BackchannelLogoutSessionSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBackchannelLogoutSessionSupported

`func (o *WellKnown) SetBackchannelLogoutSessionSupported(v bool)`

SetBackchannelLogoutSessionSupported sets BackchannelLogoutSessionSupported field to given value.

### HasBackchannelLogoutSessionSupported

`func (o *WellKnown) HasBackchannelLogoutSessionSupported() bool`

HasBackchannelLogoutSessionSupported returns a boolean if a field has been set.

### GetBackchannelLogoutSupported

`func (o *WellKnown) GetBackchannelLogoutSupported() bool`

GetBackchannelLogoutSupported returns the BackchannelLogoutSupported field if non-nil, zero value otherwise.

### GetBackchannelLogoutSupportedOk

`func (o *WellKnown) GetBackchannelLogoutSupportedOk() (*bool, bool)`

GetBackchannelLogoutSupportedOk returns a tuple with the BackchannelLogoutSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBackchannelLogoutSupported

`func (o *WellKnown) SetBackchannelLogoutSupported(v bool)`

SetBackchannelLogoutSupported sets BackchannelLogoutSupported field to given value.

### HasBackchannelLogoutSupported

`func (o *WellKnown) HasBackchannelLogoutSupported() bool`

HasBackchannelLogoutSupported returns a boolean if a field has been set.

### GetClaimsParameterSupported

`func (o *WellKnown) GetClaimsParameterSupported() bool`

GetClaimsParameterSupported returns the ClaimsParameterSupported field if non-nil, zero value otherwise.

### GetClaimsParameterSupportedOk

`func (o *WellKnown) GetClaimsParameterSupportedOk() (*bool, bool)`

GetClaimsParameterSupportedOk returns a tuple with the ClaimsParameterSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClaimsParameterSupported

`func (o *WellKnown) SetClaimsParameterSupported(v bool)`

SetClaimsParameterSupported sets ClaimsParameterSupported field to given value.

### HasClaimsParameterSupported

`func (o *WellKnown) HasClaimsParameterSupported() bool`

HasClaimsParameterSupported returns a boolean if a field has been set.

### GetClaimsSupported

`func (o *WellKnown) GetClaimsSupported() []string`

GetClaimsSupported returns the ClaimsSupported field if non-nil, zero value otherwise.

### GetClaimsSupportedOk

`func (o *WellKnown) GetClaimsSupportedOk() (*[]string, bool)`

GetClaimsSupportedOk returns a tuple with the ClaimsSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClaimsSupported

`func (o *WellKnown) SetClaimsSupported(v []string)`

SetClaimsSupported sets ClaimsSupported field to given value.

### HasClaimsSupported

`func (o *WellKnown) HasClaimsSupported() bool`

HasClaimsSupported returns a boolean if a field has been set.

### GetCodeChallengeMethodsSupported

`func (o *WellKnown) GetCodeChallengeMethodsSupported() []string`

GetCodeChallengeMethodsSupported returns the CodeChallengeMethodsSupported field if non-nil, zero value otherwise.

### GetCodeChallengeMethodsSupportedOk

`func (o *WellKnown) GetCodeChallengeMethodsSupportedOk() (*[]string, bool)`

GetCodeChallengeMethodsSupportedOk returns a tuple with the CodeChallengeMethodsSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCodeChallengeMethodsSupported

`func (o *WellKnown) SetCodeChallengeMethodsSupported(v []string)`

SetCodeChallengeMethodsSupported sets CodeChallengeMethodsSupported field to given value.

### HasCodeChallengeMethodsSupported

`func (o *WellKnown) HasCodeChallengeMethodsSupported() bool`

HasCodeChallengeMethodsSupported returns a boolean if a field has been set.

### GetEndSessionEndpoint

`func (o *WellKnown) GetEndSessionEndpoint() string`

GetEndSessionEndpoint returns the EndSessionEndpoint field if non-nil, zero value otherwise.

### GetEndSessionEndpointOk

`func (o *WellKnown) GetEndSessionEndpointOk() (*string, bool)`

GetEndSessionEndpointOk returns a tuple with the EndSessionEndpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndSessionEndpoint

`func (o *WellKnown) SetEndSessionEndpoint(v string)`

SetEndSessionEndpoint sets EndSessionEndpoint field to given value.

### HasEndSessionEndpoint

`func (o *WellKnown) HasEndSessionEndpoint() bool`

HasEndSessionEndpoint returns a boolean if a field has been set.

### GetFrontchannelLogoutSessionSupported

`func (o *WellKnown) GetFrontchannelLogoutSessionSupported() bool`

GetFrontchannelLogoutSessionSupported returns the FrontchannelLogoutSessionSupported field if non-nil, zero value otherwise.

### GetFrontchannelLogoutSessionSupportedOk

`func (o *WellKnown) GetFrontchannelLogoutSessionSupportedOk() (*bool, bool)`

GetFrontchannelLogoutSessionSupportedOk returns a tuple with the FrontchannelLogoutSessionSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFrontchannelLogoutSessionSupported

`func (o *WellKnown) SetFrontchannelLogoutSessionSupported(v bool)`

SetFrontchannelLogoutSessionSupported sets FrontchannelLogoutSessionSupported field to given value.

### HasFrontchannelLogoutSessionSupported

`func (o *WellKnown) HasFrontchannelLogoutSessionSupported() bool`

HasFrontchannelLogoutSessionSupported returns a boolean if a field has been set.

### GetFrontchannelLogoutSupported

`func (o *WellKnown) GetFrontchannelLogoutSupported() bool`

GetFrontchannelLogoutSupported returns the FrontchannelLogoutSupported field if non-nil, zero value otherwise.

### GetFrontchannelLogoutSupportedOk

`func (o *WellKnown) GetFrontchannelLogoutSupportedOk() (*bool, bool)`

GetFrontchannelLogoutSupportedOk returns a tuple with the FrontchannelLogoutSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFrontchannelLogoutSupported

`func (o *WellKnown) SetFrontchannelLogoutSupported(v bool)`

SetFrontchannelLogoutSupported sets FrontchannelLogoutSupported field to given value.

### HasFrontchannelLogoutSupported

`func (o *WellKnown) HasFrontchannelLogoutSupported() bool`

HasFrontchannelLogoutSupported returns a boolean if a field has been set.

### GetGrantTypesSupported

`func (o *WellKnown) GetGrantTypesSupported() []string`

GetGrantTypesSupported returns the GrantTypesSupported field if non-nil, zero value otherwise.

### GetGrantTypesSupportedOk

`func (o *WellKnown) GetGrantTypesSupportedOk() (*[]string, bool)`

GetGrantTypesSupportedOk returns a tuple with the GrantTypesSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGrantTypesSupported

`func (o *WellKnown) SetGrantTypesSupported(v []string)`

SetGrantTypesSupported sets GrantTypesSupported field to given value.

### HasGrantTypesSupported

`func (o *WellKnown) HasGrantTypesSupported() bool`

HasGrantTypesSupported returns a boolean if a field has been set.

### GetIdTokenSigningAlgValuesSupported

`func (o *WellKnown) GetIdTokenSigningAlgValuesSupported() []string`

GetIdTokenSigningAlgValuesSupported returns the IdTokenSigningAlgValuesSupported field if non-nil, zero value otherwise.

### GetIdTokenSigningAlgValuesSupportedOk

`func (o *WellKnown) GetIdTokenSigningAlgValuesSupportedOk() (*[]string, bool)`

GetIdTokenSigningAlgValuesSupportedOk returns a tuple with the IdTokenSigningAlgValuesSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIdTokenSigningAlgValuesSupported

`func (o *WellKnown) SetIdTokenSigningAlgValuesSupported(v []string)`

SetIdTokenSigningAlgValuesSupported sets IdTokenSigningAlgValuesSupported field to given value.


### GetIssuer

`func (o *WellKnown) GetIssuer() string`

GetIssuer returns the Issuer field if non-nil, zero value otherwise.

### GetIssuerOk

`func (o *WellKnown) GetIssuerOk() (*string, bool)`

GetIssuerOk returns a tuple with the Issuer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIssuer

`func (o *WellKnown) SetIssuer(v string)`

SetIssuer sets Issuer field to given value.


### GetJwksUri

`func (o *WellKnown) GetJwksUri() string`

GetJwksUri returns the JwksUri field if non-nil, zero value otherwise.

### GetJwksUriOk

`func (o *WellKnown) GetJwksUriOk() (*string, bool)`

GetJwksUriOk returns a tuple with the JwksUri field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetJwksUri

`func (o *WellKnown) SetJwksUri(v string)`

SetJwksUri sets JwksUri field to given value.


### GetRegistrationEndpoint

`func (o *WellKnown) GetRegistrationEndpoint() string`

GetRegistrationEndpoint returns the RegistrationEndpoint field if non-nil, zero value otherwise.

### GetRegistrationEndpointOk

`func (o *WellKnown) GetRegistrationEndpointOk() (*string, bool)`

GetRegistrationEndpointOk returns a tuple with the RegistrationEndpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegistrationEndpoint

`func (o *WellKnown) SetRegistrationEndpoint(v string)`

SetRegistrationEndpoint sets RegistrationEndpoint field to given value.

### HasRegistrationEndpoint

`func (o *WellKnown) HasRegistrationEndpoint() bool`

HasRegistrationEndpoint returns a boolean if a field has been set.

### GetRequestObjectSigningAlgValuesSupported

`func (o *WellKnown) GetRequestObjectSigningAlgValuesSupported() []string`

GetRequestObjectSigningAlgValuesSupported returns the RequestObjectSigningAlgValuesSupported field if non-nil, zero value otherwise.

### GetRequestObjectSigningAlgValuesSupportedOk

`func (o *WellKnown) GetRequestObjectSigningAlgValuesSupportedOk() (*[]string, bool)`

GetRequestObjectSigningAlgValuesSupportedOk returns a tuple with the RequestObjectSigningAlgValuesSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestObjectSigningAlgValuesSupported

`func (o *WellKnown) SetRequestObjectSigningAlgValuesSupported(v []string)`

SetRequestObjectSigningAlgValuesSupported sets RequestObjectSigningAlgValuesSupported field to given value.

### HasRequestObjectSigningAlgValuesSupported

`func (o *WellKnown) HasRequestObjectSigningAlgValuesSupported() bool`

HasRequestObjectSigningAlgValuesSupported returns a boolean if a field has been set.

### GetRequestParameterSupported

`func (o *WellKnown) GetRequestParameterSupported() bool`

GetRequestParameterSupported returns the RequestParameterSupported field if non-nil, zero value otherwise.

### GetRequestParameterSupportedOk

`func (o *WellKnown) GetRequestParameterSupportedOk() (*bool, bool)`

GetRequestParameterSupportedOk returns a tuple with the RequestParameterSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestParameterSupported

`func (o *WellKnown) SetRequestParameterSupported(v bool)`

SetRequestParameterSupported sets RequestParameterSupported field to given value.

### HasRequestParameterSupported

`func (o *WellKnown) HasRequestParameterSupported() bool`

HasRequestParameterSupported returns a boolean if a field has been set.

### GetRequestUriParameterSupported

`func (o *WellKnown) GetRequestUriParameterSupported() bool`

GetRequestUriParameterSupported returns the RequestUriParameterSupported field if non-nil, zero value otherwise.

### GetRequestUriParameterSupportedOk

`func (o *WellKnown) GetRequestUriParameterSupportedOk() (*bool, bool)`

GetRequestUriParameterSupportedOk returns a tuple with the RequestUriParameterSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestUriParameterSupported

`func (o *WellKnown) SetRequestUriParameterSupported(v bool)`

SetRequestUriParameterSupported sets RequestUriParameterSupported field to given value.

### HasRequestUriParameterSupported

`func (o *WellKnown) HasRequestUriParameterSupported() bool`

HasRequestUriParameterSupported returns a boolean if a field has been set.

### GetRequireRequestUriRegistration

`func (o *WellKnown) GetRequireRequestUriRegistration() bool`

GetRequireRequestUriRegistration returns the RequireRequestUriRegistration field if non-nil, zero value otherwise.

### GetRequireRequestUriRegistrationOk

`func (o *WellKnown) GetRequireRequestUriRegistrationOk() (*bool, bool)`

GetRequireRequestUriRegistrationOk returns a tuple with the RequireRequestUriRegistration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequireRequestUriRegistration

`func (o *WellKnown) SetRequireRequestUriRegistration(v bool)`

SetRequireRequestUriRegistration sets RequireRequestUriRegistration field to given value.

### HasRequireRequestUriRegistration

`func (o *WellKnown) HasRequireRequestUriRegistration() bool`

HasRequireRequestUriRegistration returns a boolean if a field has been set.

### GetResponseModesSupported

`func (o *WellKnown) GetResponseModesSupported() []string`

GetResponseModesSupported returns the ResponseModesSupported field if non-nil, zero value otherwise.

### GetResponseModesSupportedOk

`func (o *WellKnown) GetResponseModesSupportedOk() (*[]string, bool)`

GetResponseModesSupportedOk returns a tuple with the ResponseModesSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResponseModesSupported

`func (o *WellKnown) SetResponseModesSupported(v []string)`

SetResponseModesSupported sets ResponseModesSupported field to given value.

### HasResponseModesSupported

`func (o *WellKnown) HasResponseModesSupported() bool`

HasResponseModesSupported returns a boolean if a field has been set.

### GetResponseTypesSupported

`func (o *WellKnown) GetResponseTypesSupported() []string`

GetResponseTypesSupported returns the ResponseTypesSupported field if non-nil, zero value otherwise.

### GetResponseTypesSupportedOk

`func (o *WellKnown) GetResponseTypesSupportedOk() (*[]string, bool)`

GetResponseTypesSupportedOk returns a tuple with the ResponseTypesSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResponseTypesSupported

`func (o *WellKnown) SetResponseTypesSupported(v []string)`

SetResponseTypesSupported sets ResponseTypesSupported field to given value.


### GetRevocationEndpoint

`func (o *WellKnown) GetRevocationEndpoint() string`

GetRevocationEndpoint returns the RevocationEndpoint field if non-nil, zero value otherwise.

### GetRevocationEndpointOk

`func (o *WellKnown) GetRevocationEndpointOk() (*string, bool)`

GetRevocationEndpointOk returns a tuple with the RevocationEndpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRevocationEndpoint

`func (o *WellKnown) SetRevocationEndpoint(v string)`

SetRevocationEndpoint sets RevocationEndpoint field to given value.

### HasRevocationEndpoint

`func (o *WellKnown) HasRevocationEndpoint() bool`

HasRevocationEndpoint returns a boolean if a field has been set.

### GetScopesSupported

`func (o *WellKnown) GetScopesSupported() []string`

GetScopesSupported returns the ScopesSupported field if non-nil, zero value otherwise.

### GetScopesSupportedOk

`func (o *WellKnown) GetScopesSupportedOk() (*[]string, bool)`

GetScopesSupportedOk returns a tuple with the ScopesSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScopesSupported

`func (o *WellKnown) SetScopesSupported(v []string)`

SetScopesSupported sets ScopesSupported field to given value.

### HasScopesSupported

`func (o *WellKnown) HasScopesSupported() bool`

HasScopesSupported returns a boolean if a field has been set.

### GetSubjectTypesSupported

`func (o *WellKnown) GetSubjectTypesSupported() []string`

GetSubjectTypesSupported returns the SubjectTypesSupported field if non-nil, zero value otherwise.

### GetSubjectTypesSupportedOk

`func (o *WellKnown) GetSubjectTypesSupportedOk() (*[]string, bool)`

GetSubjectTypesSupportedOk returns a tuple with the SubjectTypesSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubjectTypesSupported

`func (o *WellKnown) SetSubjectTypesSupported(v []string)`

SetSubjectTypesSupported sets SubjectTypesSupported field to given value.


### GetTokenEndpoint

`func (o *WellKnown) GetTokenEndpoint() string`

GetTokenEndpoint returns the TokenEndpoint field if non-nil, zero value otherwise.

### GetTokenEndpointOk

`func (o *WellKnown) GetTokenEndpointOk() (*string, bool)`

GetTokenEndpointOk returns a tuple with the TokenEndpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTokenEndpoint

`func (o *WellKnown) SetTokenEndpoint(v string)`

SetTokenEndpoint sets TokenEndpoint field to given value.


### GetTokenEndpointAuthMethodsSupported

`func (o *WellKnown) GetTokenEndpointAuthMethodsSupported() []string`

GetTokenEndpointAuthMethodsSupported returns the TokenEndpointAuthMethodsSupported field if non-nil, zero value otherwise.

### GetTokenEndpointAuthMethodsSupportedOk

`func (o *WellKnown) GetTokenEndpointAuthMethodsSupportedOk() (*[]string, bool)`

GetTokenEndpointAuthMethodsSupportedOk returns a tuple with the TokenEndpointAuthMethodsSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTokenEndpointAuthMethodsSupported

`func (o *WellKnown) SetTokenEndpointAuthMethodsSupported(v []string)`

SetTokenEndpointAuthMethodsSupported sets TokenEndpointAuthMethodsSupported field to given value.

### HasTokenEndpointAuthMethodsSupported

`func (o *WellKnown) HasTokenEndpointAuthMethodsSupported() bool`

HasTokenEndpointAuthMethodsSupported returns a boolean if a field has been set.

### GetUserinfoEndpoint

`func (o *WellKnown) GetUserinfoEndpoint() string`

GetUserinfoEndpoint returns the UserinfoEndpoint field if non-nil, zero value otherwise.

### GetUserinfoEndpointOk

`func (o *WellKnown) GetUserinfoEndpointOk() (*string, bool)`

GetUserinfoEndpointOk returns a tuple with the UserinfoEndpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserinfoEndpoint

`func (o *WellKnown) SetUserinfoEndpoint(v string)`

SetUserinfoEndpoint sets UserinfoEndpoint field to given value.

### HasUserinfoEndpoint

`func (o *WellKnown) HasUserinfoEndpoint() bool`

HasUserinfoEndpoint returns a boolean if a field has been set.

### GetUserinfoSigningAlgValuesSupported

`func (o *WellKnown) GetUserinfoSigningAlgValuesSupported() []string`

GetUserinfoSigningAlgValuesSupported returns the UserinfoSigningAlgValuesSupported field if non-nil, zero value otherwise.

### GetUserinfoSigningAlgValuesSupportedOk

`func (o *WellKnown) GetUserinfoSigningAlgValuesSupportedOk() (*[]string, bool)`

GetUserinfoSigningAlgValuesSupportedOk returns a tuple with the UserinfoSigningAlgValuesSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserinfoSigningAlgValuesSupported

`func (o *WellKnown) SetUserinfoSigningAlgValuesSupported(v []string)`

SetUserinfoSigningAlgValuesSupported sets UserinfoSigningAlgValuesSupported field to given value.

### HasUserinfoSigningAlgValuesSupported

`func (o *WellKnown) HasUserinfoSigningAlgValuesSupported() bool`

HasUserinfoSigningAlgValuesSupported returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


