# AcceptLoginRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Acr** | Pointer to **string** | ACR sets the Authentication AuthorizationContext Class Reference value for this authentication session. You can use it to express that, for example, a user authenticated using two factor authentication. | [optional] 
**Amr** | Pointer to **[]string** |  | [optional] 
**Context** | Pointer to **map[string]interface{}** |  | [optional] 
**ForceSubjectIdentifier** | Pointer to **string** | ForceSubjectIdentifier forces the \&quot;pairwise\&quot; user ID of the end-user that authenticated. The \&quot;pairwise\&quot; user ID refers to the (Pairwise Identifier Algorithm)[http://openid.net/specs/openid-connect-core-1_0.html#PairwiseAlg] of the OpenID Connect specification. It allows you to set an obfuscated subject (\&quot;user\&quot;) identifier that is unique to the client.  Please note that this changes the user ID on endpoint /userinfo and sub claim of the ID Token. It does not change the sub claim in the OAuth 2.0 Introspection.  Per default, ORY Hydra handles this value with its own algorithm. In case you want to set this yourself you can use this field. Please note that setting this field has no effect if &#x60;pairwise&#x60; is not configured in ORY Hydra or the OAuth 2.0 Client does not expect a pairwise identifier (set via &#x60;subject_type&#x60; key in the client&#39;s configuration).  Please also be aware that ORY Hydra is unable to properly compute this value during authentication. This implies that you have to compute this value on every authentication process (probably depending on the client ID or some other unique value).  If you fail to compute the proper value, then authentication processes which have id_token_hint set might fail. | [optional] 
**Remember** | Pointer to **bool** | Remember, if set to true, tells ORY Hydra to remember this user by telling the user agent (browser) to store a cookie with authentication data. If the same user performs another OAuth 2.0 Authorization Request, he/she will not be asked to log in again. | [optional] 
**RememberFor** | Pointer to **int64** | RememberFor sets how long the authentication should be remembered for in seconds. If set to &#x60;0&#x60;, the authorization will be remembered for the duration of the browser session (using a session cookie). | [optional] 
**Subject** | **string** | Subject is the user ID of the end-user that authenticated. | 

## Methods

### NewAcceptLoginRequest

`func NewAcceptLoginRequest(subject string, ) *AcceptLoginRequest`

NewAcceptLoginRequest instantiates a new AcceptLoginRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAcceptLoginRequestWithDefaults

`func NewAcceptLoginRequestWithDefaults() *AcceptLoginRequest`

NewAcceptLoginRequestWithDefaults instantiates a new AcceptLoginRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAcr

`func (o *AcceptLoginRequest) GetAcr() string`

GetAcr returns the Acr field if non-nil, zero value otherwise.

### GetAcrOk

`func (o *AcceptLoginRequest) GetAcrOk() (*string, bool)`

GetAcrOk returns a tuple with the Acr field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAcr

`func (o *AcceptLoginRequest) SetAcr(v string)`

SetAcr sets Acr field to given value.

### HasAcr

`func (o *AcceptLoginRequest) HasAcr() bool`

HasAcr returns a boolean if a field has been set.

### GetAmr

`func (o *AcceptLoginRequest) GetAmr() []string`

GetAmr returns the Amr field if non-nil, zero value otherwise.

### GetAmrOk

`func (o *AcceptLoginRequest) GetAmrOk() (*[]string, bool)`

GetAmrOk returns a tuple with the Amr field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAmr

`func (o *AcceptLoginRequest) SetAmr(v []string)`

SetAmr sets Amr field to given value.

### HasAmr

`func (o *AcceptLoginRequest) HasAmr() bool`

HasAmr returns a boolean if a field has been set.

### GetContext

`func (o *AcceptLoginRequest) GetContext() map[string]interface{}`

GetContext returns the Context field if non-nil, zero value otherwise.

### GetContextOk

`func (o *AcceptLoginRequest) GetContextOk() (*map[string]interface{}, bool)`

GetContextOk returns a tuple with the Context field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContext

`func (o *AcceptLoginRequest) SetContext(v map[string]interface{})`

SetContext sets Context field to given value.

### HasContext

`func (o *AcceptLoginRequest) HasContext() bool`

HasContext returns a boolean if a field has been set.

### GetForceSubjectIdentifier

`func (o *AcceptLoginRequest) GetForceSubjectIdentifier() string`

GetForceSubjectIdentifier returns the ForceSubjectIdentifier field if non-nil, zero value otherwise.

### GetForceSubjectIdentifierOk

`func (o *AcceptLoginRequest) GetForceSubjectIdentifierOk() (*string, bool)`

GetForceSubjectIdentifierOk returns a tuple with the ForceSubjectIdentifier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetForceSubjectIdentifier

`func (o *AcceptLoginRequest) SetForceSubjectIdentifier(v string)`

SetForceSubjectIdentifier sets ForceSubjectIdentifier field to given value.

### HasForceSubjectIdentifier

`func (o *AcceptLoginRequest) HasForceSubjectIdentifier() bool`

HasForceSubjectIdentifier returns a boolean if a field has been set.

### GetRemember

`func (o *AcceptLoginRequest) GetRemember() bool`

GetRemember returns the Remember field if non-nil, zero value otherwise.

### GetRememberOk

`func (o *AcceptLoginRequest) GetRememberOk() (*bool, bool)`

GetRememberOk returns a tuple with the Remember field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRemember

`func (o *AcceptLoginRequest) SetRemember(v bool)`

SetRemember sets Remember field to given value.

### HasRemember

`func (o *AcceptLoginRequest) HasRemember() bool`

HasRemember returns a boolean if a field has been set.

### GetRememberFor

`func (o *AcceptLoginRequest) GetRememberFor() int64`

GetRememberFor returns the RememberFor field if non-nil, zero value otherwise.

### GetRememberForOk

`func (o *AcceptLoginRequest) GetRememberForOk() (*int64, bool)`

GetRememberForOk returns a tuple with the RememberFor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRememberFor

`func (o *AcceptLoginRequest) SetRememberFor(v int64)`

SetRememberFor sets RememberFor field to given value.

### HasRememberFor

`func (o *AcceptLoginRequest) HasRememberFor() bool`

HasRememberFor returns a boolean if a field has been set.

### GetSubject

`func (o *AcceptLoginRequest) GetSubject() string`

GetSubject returns the Subject field if non-nil, zero value otherwise.

### GetSubjectOk

`func (o *AcceptLoginRequest) GetSubjectOk() (*string, bool)`

GetSubjectOk returns a tuple with the Subject field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubject

`func (o *AcceptLoginRequest) SetSubject(v string)`

SetSubject sets Subject field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


