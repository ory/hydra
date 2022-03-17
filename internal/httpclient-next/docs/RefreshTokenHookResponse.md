# RefreshTokenHookResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Session** | Pointer to [**ConsentRequestSession**](ConsentRequestSession.md) |  | [optional] 

## Methods

### NewRefreshTokenHookResponse

`func NewRefreshTokenHookResponse() *RefreshTokenHookResponse`

NewRefreshTokenHookResponse instantiates a new RefreshTokenHookResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRefreshTokenHookResponseWithDefaults

`func NewRefreshTokenHookResponseWithDefaults() *RefreshTokenHookResponse`

NewRefreshTokenHookResponseWithDefaults instantiates a new RefreshTokenHookResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSession

`func (o *RefreshTokenHookResponse) GetSession() ConsentRequestSession`

GetSession returns the Session field if non-nil, zero value otherwise.

### GetSessionOk

`func (o *RefreshTokenHookResponse) GetSessionOk() (*ConsentRequestSession, bool)`

GetSessionOk returns a tuple with the Session field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSession

`func (o *RefreshTokenHookResponse) SetSession(v ConsentRequestSession)`

SetSession sets Session field to given value.

### HasSession

`func (o *RefreshTokenHookResponse) HasSession() bool`

HasSession returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


