# AcceptOAuth2LoginRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Acr** | Pointer to **string** | ACR sets the Authentication AuthorizationContext Class Reference value for this authentication session. You can use it to express that, for example, a user authenticated using two factor authentication. | [optional] 
**Amr** | Pointer to **[]string** |  | [optional] 
**Context** | Pointer to **interface{}** |  | [optional] 
**ExtendSessionLifespan** | Pointer to **bool** | Extend OAuth2 authentication session lifespan  If set to &#x60;true&#x60;, the OAuth2 authentication cookie lifespan is extended. This is for example useful if you want the user to be able to use &#x60;prompt&#x3D;none&#x60; continuously.  This value can only be set to &#x60;true&#x60; if the user has an authentication, which is the case if the &#x60;skip&#x60; value is &#x60;true&#x60;. | [optional] 
**ForceSubjectIdentifier** | Pointer to **string** | ForceSubjectIdentifier forces the \&quot;pairwise\&quot; user ID of the end-user that authenticated. The \&quot;pairwise\&quot; user ID refers to the (Pairwise Identifier Algorithm)[http://openid.net/specs/openid-connect-core-1_0.html#PairwiseAlg] of the OpenID Connect specification. It allows you to set an obfuscated subject (\&quot;user\&quot;) identifier that is unique to the client.  Please note that this changes the user ID on endpoint /userinfo and sub claim of the ID Token. It does not change the sub claim in the OAuth 2.0 Introspection.  Per default, ORY Hydra handles this value with its own algorithm. In case you want to set this yourself you can use this field. Please note that setting this field has no effect if &#x60;pairwise&#x60; is not configured in ORY Hydra or the OAuth 2.0 Client does not expect a pairwise identifier (set via &#x60;subject_type&#x60; key in the client&#39;s configuration).  Please also be aware that ORY Hydra is unable to properly compute this value during authentication. This implies that you have to compute this value on every authentication process (probably depending on the client ID or some other unique value).  If you fail to compute the proper value, then authentication processes which have id_token_hint set might fail. | [optional] 
**IdentityProviderSessionId** | Pointer to **string** | IdentityProviderSessionID is the session ID of the end-user that authenticated. If specified, we will use this value to propagate the logout. | [optional] 
**Remember** | Pointer to **bool** | Remember, if set to true, tells ORY Hydra to remember this user by telling the user agent (browser) to store a cookie with authentication data. If the same user performs another OAuth 2.0 Authorization Request, he/she will not be asked to log in again. | [optional] 
**RememberFor** | Pointer to **int64** | RememberFor sets how long the authentication should be remembered for in seconds. If set to &#x60;0&#x60;, the authorization will be remembered for the duration of the browser session (using a session cookie). | [optional] 
**Subject** | **string** | Subject is the user ID of the end-user that authenticated. | 

## Methods

### NewAcceptOAuth2LoginRequest

`func NewAcceptOAuth2LoginRequest(subject string, ) *AcceptOAuth2LoginRequest`

NewAcceptOAuth2LoginRequest instantiates a new AcceptOAuth2LoginRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAcceptOAuth2LoginRequestWithDefaults

`func NewAcceptOAuth2LoginRequestWithDefaults() *AcceptOAuth2LoginRequest`

NewAcceptOAuth2LoginRequestWithDefaults instantiates a new AcceptOAuth2LoginRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAcr

`func (o *AcceptOAuth2LoginRequest) GetAcr() string`

GetAcr returns the Acr field if non-nil, zero value otherwise.

### GetAcrOk

`func (o *AcceptOAuth2LoginRequest) GetAcrOk() (*string, bool)`

GetAcrOk returns a tuple with the Acr field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAcr

`func (o *AcceptOAuth2LoginRequest) SetAcr(v string)`

SetAcr sets Acr field to given value.

### HasAcr

`func (o *AcceptOAuth2LoginRequest) HasAcr() bool`

HasAcr returns a boolean if a field has been set.

### GetAmr

`func (o *AcceptOAuth2LoginRequest) GetAmr() []string`

GetAmr returns the Amr field if non-nil, zero value otherwise.

### GetAmrOk

`func (o *AcceptOAuth2LoginRequest) GetAmrOk() (*[]string, bool)`

GetAmrOk returns a tuple with the Amr field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAmr

`func (o *AcceptOAuth2LoginRequest) SetAmr(v []string)`

SetAmr sets Amr field to given value.

### HasAmr

`func (o *AcceptOAuth2LoginRequest) HasAmr() bool`

HasAmr returns a boolean if a field has been set.

### GetContext

`func (o *AcceptOAuth2LoginRequest) GetContext() interface{}`

GetContext returns the Context field if non-nil, zero value otherwise.

### GetContextOk

`func (o *AcceptOAuth2LoginRequest) GetContextOk() (*interface{}, bool)`

GetContextOk returns a tuple with the Context field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContext

`func (o *AcceptOAuth2LoginRequest) SetContext(v interface{})`

SetContext sets Context field to given value.

### HasContext

`func (o *AcceptOAuth2LoginRequest) HasContext() bool`

HasContext returns a boolean if a field has been set.

### SetContextNil

`func (o *AcceptOAuth2LoginRequest) SetContextNil(b bool)`

 SetContextNil sets the value for Context to be an explicit nil

### UnsetContext
`func (o *AcceptOAuth2LoginRequest) UnsetContext()`

UnsetContext ensures that no value is present for Context, not even an explicit nil
### GetExtendSessionLifespan

`func (o *AcceptOAuth2LoginRequest) GetExtendSessionLifespan() bool`

GetExtendSessionLifespan returns the ExtendSessionLifespan field if non-nil, zero value otherwise.

### GetExtendSessionLifespanOk

`func (o *AcceptOAuth2LoginRequest) GetExtendSessionLifespanOk() (*bool, bool)`

GetExtendSessionLifespanOk returns a tuple with the ExtendSessionLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExtendSessionLifespan

`func (o *AcceptOAuth2LoginRequest) SetExtendSessionLifespan(v bool)`

SetExtendSessionLifespan sets ExtendSessionLifespan field to given value.

### HasExtendSessionLifespan

`func (o *AcceptOAuth2LoginRequest) HasExtendSessionLifespan() bool`

HasExtendSessionLifespan returns a boolean if a field has been set.

### GetForceSubjectIdentifier

`func (o *AcceptOAuth2LoginRequest) GetForceSubjectIdentifier() string`

GetForceSubjectIdentifier returns the ForceSubjectIdentifier field if non-nil, zero value otherwise.

### GetForceSubjectIdentifierOk

`func (o *AcceptOAuth2LoginRequest) GetForceSubjectIdentifierOk() (*string, bool)`

GetForceSubjectIdentifierOk returns a tuple with the ForceSubjectIdentifier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetForceSubjectIdentifier

`func (o *AcceptOAuth2LoginRequest) SetForceSubjectIdentifier(v string)`

SetForceSubjectIdentifier sets ForceSubjectIdentifier field to given value.

### HasForceSubjectIdentifier

`func (o *AcceptOAuth2LoginRequest) HasForceSubjectIdentifier() bool`

HasForceSubjectIdentifier returns a boolean if a field has been set.

### GetIdentityProviderSessionId

`func (o *AcceptOAuth2LoginRequest) GetIdentityProviderSessionId() string`

GetIdentityProviderSessionId returns the IdentityProviderSessionId field if non-nil, zero value otherwise.

### GetIdentityProviderSessionIdOk

`func (o *AcceptOAuth2LoginRequest) GetIdentityProviderSessionIdOk() (*string, bool)`

GetIdentityProviderSessionIdOk returns a tuple with the IdentityProviderSessionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIdentityProviderSessionId

`func (o *AcceptOAuth2LoginRequest) SetIdentityProviderSessionId(v string)`

SetIdentityProviderSessionId sets IdentityProviderSessionId field to given value.

### HasIdentityProviderSessionId

`func (o *AcceptOAuth2LoginRequest) HasIdentityProviderSessionId() bool`

HasIdentityProviderSessionId returns a boolean if a field has been set.

### GetRemember

`func (o *AcceptOAuth2LoginRequest) GetRemember() bool`

GetRemember returns the Remember field if non-nil, zero value otherwise.

### GetRememberOk

`func (o *AcceptOAuth2LoginRequest) GetRememberOk() (*bool, bool)`

GetRememberOk returns a tuple with the Remember field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRemember

`func (o *AcceptOAuth2LoginRequest) SetRemember(v bool)`

SetRemember sets Remember field to given value.

### HasRemember

`func (o *AcceptOAuth2LoginRequest) HasRemember() bool`

HasRemember returns a boolean if a field has been set.

### GetRememberFor

`func (o *AcceptOAuth2LoginRequest) GetRememberFor() int64`

GetRememberFor returns the RememberFor field if non-nil, zero value otherwise.

### GetRememberForOk

`func (o *AcceptOAuth2LoginRequest) GetRememberForOk() (*int64, bool)`

GetRememberForOk returns a tuple with the RememberFor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRememberFor

`func (o *AcceptOAuth2LoginRequest) SetRememberFor(v int64)`

SetRememberFor sets RememberFor field to given value.

### HasRememberFor

`func (o *AcceptOAuth2LoginRequest) HasRememberFor() bool`

HasRememberFor returns a boolean if a field has been set.

### GetSubject

`func (o *AcceptOAuth2LoginRequest) GetSubject() string`

GetSubject returns the Subject field if non-nil, zero value otherwise.

### GetSubjectOk

`func (o *AcceptOAuth2LoginRequest) GetSubjectOk() (*string, bool)`

GetSubjectOk returns a tuple with the Subject field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubject

`func (o *AcceptOAuth2LoginRequest) SetSubject(v string)`

SetSubject sets Subject field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


