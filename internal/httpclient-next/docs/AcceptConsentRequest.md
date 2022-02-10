# AcceptConsentRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**GrantAccessTokenAudience** | Pointer to **[]string** |  | [optional] 
**GrantScope** | Pointer to **[]string** |  | [optional] 
**HandledAt** | Pointer to **time.Time** |  | [optional] 
**Remember** | Pointer to **bool** | Remember, if set to true, tells ORY Hydra to remember this consent authorization and reuse it if the same client asks the same user for the same, or a subset of, scope. | [optional] 
**RememberFor** | Pointer to **int64** | RememberFor sets how long the consent authorization should be remembered for in seconds. If set to &#x60;0&#x60;, the authorization will be remembered indefinitely. | [optional] 
**Session** | Pointer to [**ConsentRequestSession**](ConsentRequestSession.md) |  | [optional] 

## Methods

### NewAcceptConsentRequest

`func NewAcceptConsentRequest() *AcceptConsentRequest`

NewAcceptConsentRequest instantiates a new AcceptConsentRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAcceptConsentRequestWithDefaults

`func NewAcceptConsentRequestWithDefaults() *AcceptConsentRequest`

NewAcceptConsentRequestWithDefaults instantiates a new AcceptConsentRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetGrantAccessTokenAudience

`func (o *AcceptConsentRequest) GetGrantAccessTokenAudience() []string`

GetGrantAccessTokenAudience returns the GrantAccessTokenAudience field if non-nil, zero value otherwise.

### GetGrantAccessTokenAudienceOk

`func (o *AcceptConsentRequest) GetGrantAccessTokenAudienceOk() (*[]string, bool)`

GetGrantAccessTokenAudienceOk returns a tuple with the GrantAccessTokenAudience field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGrantAccessTokenAudience

`func (o *AcceptConsentRequest) SetGrantAccessTokenAudience(v []string)`

SetGrantAccessTokenAudience sets GrantAccessTokenAudience field to given value.

### HasGrantAccessTokenAudience

`func (o *AcceptConsentRequest) HasGrantAccessTokenAudience() bool`

HasGrantAccessTokenAudience returns a boolean if a field has been set.

### GetGrantScope

`func (o *AcceptConsentRequest) GetGrantScope() []string`

GetGrantScope returns the GrantScope field if non-nil, zero value otherwise.

### GetGrantScopeOk

`func (o *AcceptConsentRequest) GetGrantScopeOk() (*[]string, bool)`

GetGrantScopeOk returns a tuple with the GrantScope field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGrantScope

`func (o *AcceptConsentRequest) SetGrantScope(v []string)`

SetGrantScope sets GrantScope field to given value.

### HasGrantScope

`func (o *AcceptConsentRequest) HasGrantScope() bool`

HasGrantScope returns a boolean if a field has been set.

### GetHandledAt

`func (o *AcceptConsentRequest) GetHandledAt() time.Time`

GetHandledAt returns the HandledAt field if non-nil, zero value otherwise.

### GetHandledAtOk

`func (o *AcceptConsentRequest) GetHandledAtOk() (*time.Time, bool)`

GetHandledAtOk returns a tuple with the HandledAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHandledAt

`func (o *AcceptConsentRequest) SetHandledAt(v time.Time)`

SetHandledAt sets HandledAt field to given value.

### HasHandledAt

`func (o *AcceptConsentRequest) HasHandledAt() bool`

HasHandledAt returns a boolean if a field has been set.

### GetRemember

`func (o *AcceptConsentRequest) GetRemember() bool`

GetRemember returns the Remember field if non-nil, zero value otherwise.

### GetRememberOk

`func (o *AcceptConsentRequest) GetRememberOk() (*bool, bool)`

GetRememberOk returns a tuple with the Remember field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRemember

`func (o *AcceptConsentRequest) SetRemember(v bool)`

SetRemember sets Remember field to given value.

### HasRemember

`func (o *AcceptConsentRequest) HasRemember() bool`

HasRemember returns a boolean if a field has been set.

### GetRememberFor

`func (o *AcceptConsentRequest) GetRememberFor() int64`

GetRememberFor returns the RememberFor field if non-nil, zero value otherwise.

### GetRememberForOk

`func (o *AcceptConsentRequest) GetRememberForOk() (*int64, bool)`

GetRememberForOk returns a tuple with the RememberFor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRememberFor

`func (o *AcceptConsentRequest) SetRememberFor(v int64)`

SetRememberFor sets RememberFor field to given value.

### HasRememberFor

`func (o *AcceptConsentRequest) HasRememberFor() bool`

HasRememberFor returns a boolean if a field has been set.

### GetSession

`func (o *AcceptConsentRequest) GetSession() ConsentRequestSession`

GetSession returns the Session field if non-nil, zero value otherwise.

### GetSessionOk

`func (o *AcceptConsentRequest) GetSessionOk() (*ConsentRequestSession, bool)`

GetSessionOk returns a tuple with the Session field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSession

`func (o *AcceptConsentRequest) SetSession(v ConsentRequestSession)`

SetSession sets Session field to given value.

### HasSession

`func (o *AcceptConsentRequest) HasSession() bool`

HasSession returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


