# FlushInactiveOAuth2TokensRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NotAfter** | Pointer to **time.Time** | NotAfter sets after which point tokens should not be flushed. This is useful when you want to keep a history of recently issued tokens for auditing. | [optional] 

## Methods

### NewFlushInactiveOAuth2TokensRequest

`func NewFlushInactiveOAuth2TokensRequest() *FlushInactiveOAuth2TokensRequest`

NewFlushInactiveOAuth2TokensRequest instantiates a new FlushInactiveOAuth2TokensRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewFlushInactiveOAuth2TokensRequestWithDefaults

`func NewFlushInactiveOAuth2TokensRequestWithDefaults() *FlushInactiveOAuth2TokensRequest`

NewFlushInactiveOAuth2TokensRequestWithDefaults instantiates a new FlushInactiveOAuth2TokensRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNotAfter

`func (o *FlushInactiveOAuth2TokensRequest) GetNotAfter() time.Time`

GetNotAfter returns the NotAfter field if non-nil, zero value otherwise.

### GetNotAfterOk

`func (o *FlushInactiveOAuth2TokensRequest) GetNotAfterOk() (*time.Time, bool)`

GetNotAfterOk returns a tuple with the NotAfter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNotAfter

`func (o *FlushInactiveOAuth2TokensRequest) SetNotAfter(v time.Time)`

SetNotAfter sets NotAfter field to given value.

### HasNotAfter

`func (o *FlushInactiveOAuth2TokensRequest) HasNotAfter() bool`

HasNotAfter returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


