# OAuth2ConsentSession

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ConsentRequest** | Pointer to [**OAuth2ConsentRequest**](OAuth2ConsentRequest.md) |  | [optional] 
**Context** | Pointer to **interface{}** |  | [optional] 
**ExpiresAt** | Pointer to [**OAuth2ConsentSessionExpiresAt**](OAuth2ConsentSessionExpiresAt.md) |  | [optional] 
**GrantAccessTokenAudience** | Pointer to **[]string** |  | [optional] 
**GrantScope** | Pointer to **[]string** |  | [optional] 
**HandledAt** | Pointer to **time.Time** |  | [optional] 
**Remember** | Pointer to **bool** | Remember Consent  Remember, if set to true, tells ORY Hydra to remember this consent authorization and reuse it if the same client asks the same user for the same, or a subset of, scope. | [optional] 
**RememberFor** | Pointer to **int64** | Remember Consent For  RememberFor sets how long the consent authorization should be remembered for in seconds. If set to &#x60;0&#x60;, the authorization will be remembered indefinitely. | [optional] 
**Session** | Pointer to [**AcceptOAuth2ConsentRequestSession**](AcceptOAuth2ConsentRequestSession.md) |  | [optional] 

## Methods

### NewOAuth2ConsentSession

`func NewOAuth2ConsentSession() *OAuth2ConsentSession`

NewOAuth2ConsentSession instantiates a new OAuth2ConsentSession object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOAuth2ConsentSessionWithDefaults

`func NewOAuth2ConsentSessionWithDefaults() *OAuth2ConsentSession`

NewOAuth2ConsentSessionWithDefaults instantiates a new OAuth2ConsentSession object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetConsentRequest

`func (o *OAuth2ConsentSession) GetConsentRequest() OAuth2ConsentRequest`

GetConsentRequest returns the ConsentRequest field if non-nil, zero value otherwise.

### GetConsentRequestOk

`func (o *OAuth2ConsentSession) GetConsentRequestOk() (*OAuth2ConsentRequest, bool)`

GetConsentRequestOk returns a tuple with the ConsentRequest field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConsentRequest

`func (o *OAuth2ConsentSession) SetConsentRequest(v OAuth2ConsentRequest)`

SetConsentRequest sets ConsentRequest field to given value.

### HasConsentRequest

`func (o *OAuth2ConsentSession) HasConsentRequest() bool`

HasConsentRequest returns a boolean if a field has been set.

### GetContext

`func (o *OAuth2ConsentSession) GetContext() interface{}`

GetContext returns the Context field if non-nil, zero value otherwise.

### GetContextOk

`func (o *OAuth2ConsentSession) GetContextOk() (*interface{}, bool)`

GetContextOk returns a tuple with the Context field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContext

`func (o *OAuth2ConsentSession) SetContext(v interface{})`

SetContext sets Context field to given value.

### HasContext

`func (o *OAuth2ConsentSession) HasContext() bool`

HasContext returns a boolean if a field has been set.

### SetContextNil

`func (o *OAuth2ConsentSession) SetContextNil(b bool)`

 SetContextNil sets the value for Context to be an explicit nil

### UnsetContext
`func (o *OAuth2ConsentSession) UnsetContext()`

UnsetContext ensures that no value is present for Context, not even an explicit nil
### GetExpiresAt

`func (o *OAuth2ConsentSession) GetExpiresAt() OAuth2ConsentSessionExpiresAt`

GetExpiresAt returns the ExpiresAt field if non-nil, zero value otherwise.

### GetExpiresAtOk

`func (o *OAuth2ConsentSession) GetExpiresAtOk() (*OAuth2ConsentSessionExpiresAt, bool)`

GetExpiresAtOk returns a tuple with the ExpiresAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiresAt

`func (o *OAuth2ConsentSession) SetExpiresAt(v OAuth2ConsentSessionExpiresAt)`

SetExpiresAt sets ExpiresAt field to given value.

### HasExpiresAt

`func (o *OAuth2ConsentSession) HasExpiresAt() bool`

HasExpiresAt returns a boolean if a field has been set.

### GetGrantAccessTokenAudience

`func (o *OAuth2ConsentSession) GetGrantAccessTokenAudience() []string`

GetGrantAccessTokenAudience returns the GrantAccessTokenAudience field if non-nil, zero value otherwise.

### GetGrantAccessTokenAudienceOk

`func (o *OAuth2ConsentSession) GetGrantAccessTokenAudienceOk() (*[]string, bool)`

GetGrantAccessTokenAudienceOk returns a tuple with the GrantAccessTokenAudience field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGrantAccessTokenAudience

`func (o *OAuth2ConsentSession) SetGrantAccessTokenAudience(v []string)`

SetGrantAccessTokenAudience sets GrantAccessTokenAudience field to given value.

### HasGrantAccessTokenAudience

`func (o *OAuth2ConsentSession) HasGrantAccessTokenAudience() bool`

HasGrantAccessTokenAudience returns a boolean if a field has been set.

### GetGrantScope

`func (o *OAuth2ConsentSession) GetGrantScope() []string`

GetGrantScope returns the GrantScope field if non-nil, zero value otherwise.

### GetGrantScopeOk

`func (o *OAuth2ConsentSession) GetGrantScopeOk() (*[]string, bool)`

GetGrantScopeOk returns a tuple with the GrantScope field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGrantScope

`func (o *OAuth2ConsentSession) SetGrantScope(v []string)`

SetGrantScope sets GrantScope field to given value.

### HasGrantScope

`func (o *OAuth2ConsentSession) HasGrantScope() bool`

HasGrantScope returns a boolean if a field has been set.

### GetHandledAt

`func (o *OAuth2ConsentSession) GetHandledAt() time.Time`

GetHandledAt returns the HandledAt field if non-nil, zero value otherwise.

### GetHandledAtOk

`func (o *OAuth2ConsentSession) GetHandledAtOk() (*time.Time, bool)`

GetHandledAtOk returns a tuple with the HandledAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHandledAt

`func (o *OAuth2ConsentSession) SetHandledAt(v time.Time)`

SetHandledAt sets HandledAt field to given value.

### HasHandledAt

`func (o *OAuth2ConsentSession) HasHandledAt() bool`

HasHandledAt returns a boolean if a field has been set.

### GetRemember

`func (o *OAuth2ConsentSession) GetRemember() bool`

GetRemember returns the Remember field if non-nil, zero value otherwise.

### GetRememberOk

`func (o *OAuth2ConsentSession) GetRememberOk() (*bool, bool)`

GetRememberOk returns a tuple with the Remember field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRemember

`func (o *OAuth2ConsentSession) SetRemember(v bool)`

SetRemember sets Remember field to given value.

### HasRemember

`func (o *OAuth2ConsentSession) HasRemember() bool`

HasRemember returns a boolean if a field has been set.

### GetRememberFor

`func (o *OAuth2ConsentSession) GetRememberFor() int64`

GetRememberFor returns the RememberFor field if non-nil, zero value otherwise.

### GetRememberForOk

`func (o *OAuth2ConsentSession) GetRememberForOk() (*int64, bool)`

GetRememberForOk returns a tuple with the RememberFor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRememberFor

`func (o *OAuth2ConsentSession) SetRememberFor(v int64)`

SetRememberFor sets RememberFor field to given value.

### HasRememberFor

`func (o *OAuth2ConsentSession) HasRememberFor() bool`

HasRememberFor returns a boolean if a field has been set.

### GetSession

`func (o *OAuth2ConsentSession) GetSession() AcceptOAuth2ConsentRequestSession`

GetSession returns the Session field if non-nil, zero value otherwise.

### GetSessionOk

`func (o *OAuth2ConsentSession) GetSessionOk() (*AcceptOAuth2ConsentRequestSession, bool)`

GetSessionOk returns a tuple with the Session field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSession

`func (o *OAuth2ConsentSession) SetSession(v AcceptOAuth2ConsentRequestSession)`

SetSession sets Session field to given value.

### HasSession

`func (o *OAuth2ConsentSession) HasSession() bool`

HasSession returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


