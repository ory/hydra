# AcceptOAuth2ConsentRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Context** | Pointer to **interface{}** |  | [optional] 
**GrantAccessTokenAudience** | Pointer to **[]string** |  | [optional] 
**GrantScope** | Pointer to **[]string** |  | [optional] 
**Remember** | Pointer to **bool** | Remember, if set to true, tells ORY Hydra to remember this consent authorization and reuse it if the same client asks the same user for the same, or a subset of, scope. | [optional] 
**RememberFor** | Pointer to **int64** | RememberFor sets how long the consent authorization should be remembered for in seconds. If set to &#x60;0&#x60;, the authorization will be remembered indefinitely. | [optional] 
**Session** | Pointer to [**AcceptOAuth2ConsentRequestSession**](AcceptOAuth2ConsentRequestSession.md) |  | [optional] 

## Methods

### NewAcceptOAuth2ConsentRequest

`func NewAcceptOAuth2ConsentRequest() *AcceptOAuth2ConsentRequest`

NewAcceptOAuth2ConsentRequest instantiates a new AcceptOAuth2ConsentRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAcceptOAuth2ConsentRequestWithDefaults

`func NewAcceptOAuth2ConsentRequestWithDefaults() *AcceptOAuth2ConsentRequest`

NewAcceptOAuth2ConsentRequestWithDefaults instantiates a new AcceptOAuth2ConsentRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetContext

`func (o *AcceptOAuth2ConsentRequest) GetContext() interface{}`

GetContext returns the Context field if non-nil, zero value otherwise.

### GetContextOk

`func (o *AcceptOAuth2ConsentRequest) GetContextOk() (*interface{}, bool)`

GetContextOk returns a tuple with the Context field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContext

`func (o *AcceptOAuth2ConsentRequest) SetContext(v interface{})`

SetContext sets Context field to given value.

### HasContext

`func (o *AcceptOAuth2ConsentRequest) HasContext() bool`

HasContext returns a boolean if a field has been set.

### SetContextNil

`func (o *AcceptOAuth2ConsentRequest) SetContextNil(b bool)`

 SetContextNil sets the value for Context to be an explicit nil

### UnsetContext
`func (o *AcceptOAuth2ConsentRequest) UnsetContext()`

UnsetContext ensures that no value is present for Context, not even an explicit nil
### GetGrantAccessTokenAudience

`func (o *AcceptOAuth2ConsentRequest) GetGrantAccessTokenAudience() []string`

GetGrantAccessTokenAudience returns the GrantAccessTokenAudience field if non-nil, zero value otherwise.

### GetGrantAccessTokenAudienceOk

`func (o *AcceptOAuth2ConsentRequest) GetGrantAccessTokenAudienceOk() (*[]string, bool)`

GetGrantAccessTokenAudienceOk returns a tuple with the GrantAccessTokenAudience field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGrantAccessTokenAudience

`func (o *AcceptOAuth2ConsentRequest) SetGrantAccessTokenAudience(v []string)`

SetGrantAccessTokenAudience sets GrantAccessTokenAudience field to given value.

### HasGrantAccessTokenAudience

`func (o *AcceptOAuth2ConsentRequest) HasGrantAccessTokenAudience() bool`

HasGrantAccessTokenAudience returns a boolean if a field has been set.

### GetGrantScope

`func (o *AcceptOAuth2ConsentRequest) GetGrantScope() []string`

GetGrantScope returns the GrantScope field if non-nil, zero value otherwise.

### GetGrantScopeOk

`func (o *AcceptOAuth2ConsentRequest) GetGrantScopeOk() (*[]string, bool)`

GetGrantScopeOk returns a tuple with the GrantScope field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGrantScope

`func (o *AcceptOAuth2ConsentRequest) SetGrantScope(v []string)`

SetGrantScope sets GrantScope field to given value.

### HasGrantScope

`func (o *AcceptOAuth2ConsentRequest) HasGrantScope() bool`

HasGrantScope returns a boolean if a field has been set.

### GetRemember

`func (o *AcceptOAuth2ConsentRequest) GetRemember() bool`

GetRemember returns the Remember field if non-nil, zero value otherwise.

### GetRememberOk

`func (o *AcceptOAuth2ConsentRequest) GetRememberOk() (*bool, bool)`

GetRememberOk returns a tuple with the Remember field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRemember

`func (o *AcceptOAuth2ConsentRequest) SetRemember(v bool)`

SetRemember sets Remember field to given value.

### HasRemember

`func (o *AcceptOAuth2ConsentRequest) HasRemember() bool`

HasRemember returns a boolean if a field has been set.

### GetRememberFor

`func (o *AcceptOAuth2ConsentRequest) GetRememberFor() int64`

GetRememberFor returns the RememberFor field if non-nil, zero value otherwise.

### GetRememberForOk

`func (o *AcceptOAuth2ConsentRequest) GetRememberForOk() (*int64, bool)`

GetRememberForOk returns a tuple with the RememberFor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRememberFor

`func (o *AcceptOAuth2ConsentRequest) SetRememberFor(v int64)`

SetRememberFor sets RememberFor field to given value.

### HasRememberFor

`func (o *AcceptOAuth2ConsentRequest) HasRememberFor() bool`

HasRememberFor returns a boolean if a field has been set.

### GetSession

`func (o *AcceptOAuth2ConsentRequest) GetSession() AcceptOAuth2ConsentRequestSession`

GetSession returns the Session field if non-nil, zero value otherwise.

### GetSessionOk

`func (o *AcceptOAuth2ConsentRequest) GetSessionOk() (*AcceptOAuth2ConsentRequestSession, bool)`

GetSessionOk returns a tuple with the Session field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSession

`func (o *AcceptOAuth2ConsentRequest) SetSession(v AcceptOAuth2ConsentRequestSession)`

SetSession sets Session field to given value.

### HasSession

`func (o *AcceptOAuth2ConsentRequest) HasSession() bool`

HasSession returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


