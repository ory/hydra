# ErrorOAuth2

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Error** | Pointer to **string** | Error | [optional] 
**ErrorDebug** | Pointer to **string** | Error Debug Information  Only available in dev mode. | [optional] 
**ErrorDescription** | Pointer to **string** | Error Description | [optional] 
**ErrorHint** | Pointer to **string** | Error Hint  Helps the user identify the error cause. | [optional] 
**StatusCode** | Pointer to **int64** | HTTP Status Code | [optional] 

## Methods

### NewErrorOAuth2

`func NewErrorOAuth2() *ErrorOAuth2`

NewErrorOAuth2 instantiates a new ErrorOAuth2 object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewErrorOAuth2WithDefaults

`func NewErrorOAuth2WithDefaults() *ErrorOAuth2`

NewErrorOAuth2WithDefaults instantiates a new ErrorOAuth2 object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetError

`func (o *ErrorOAuth2) GetError() string`

GetError returns the Error field if non-nil, zero value otherwise.

### GetErrorOk

`func (o *ErrorOAuth2) GetErrorOk() (*string, bool)`

GetErrorOk returns a tuple with the Error field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetError

`func (o *ErrorOAuth2) SetError(v string)`

SetError sets Error field to given value.

### HasError

`func (o *ErrorOAuth2) HasError() bool`

HasError returns a boolean if a field has been set.

### GetErrorDebug

`func (o *ErrorOAuth2) GetErrorDebug() string`

GetErrorDebug returns the ErrorDebug field if non-nil, zero value otherwise.

### GetErrorDebugOk

`func (o *ErrorOAuth2) GetErrorDebugOk() (*string, bool)`

GetErrorDebugOk returns a tuple with the ErrorDebug field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrorDebug

`func (o *ErrorOAuth2) SetErrorDebug(v string)`

SetErrorDebug sets ErrorDebug field to given value.

### HasErrorDebug

`func (o *ErrorOAuth2) HasErrorDebug() bool`

HasErrorDebug returns a boolean if a field has been set.

### GetErrorDescription

`func (o *ErrorOAuth2) GetErrorDescription() string`

GetErrorDescription returns the ErrorDescription field if non-nil, zero value otherwise.

### GetErrorDescriptionOk

`func (o *ErrorOAuth2) GetErrorDescriptionOk() (*string, bool)`

GetErrorDescriptionOk returns a tuple with the ErrorDescription field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrorDescription

`func (o *ErrorOAuth2) SetErrorDescription(v string)`

SetErrorDescription sets ErrorDescription field to given value.

### HasErrorDescription

`func (o *ErrorOAuth2) HasErrorDescription() bool`

HasErrorDescription returns a boolean if a field has been set.

### GetErrorHint

`func (o *ErrorOAuth2) GetErrorHint() string`

GetErrorHint returns the ErrorHint field if non-nil, zero value otherwise.

### GetErrorHintOk

`func (o *ErrorOAuth2) GetErrorHintOk() (*string, bool)`

GetErrorHintOk returns a tuple with the ErrorHint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrorHint

`func (o *ErrorOAuth2) SetErrorHint(v string)`

SetErrorHint sets ErrorHint field to given value.

### HasErrorHint

`func (o *ErrorOAuth2) HasErrorHint() bool`

HasErrorHint returns a boolean if a field has been set.

### GetStatusCode

`func (o *ErrorOAuth2) GetStatusCode() int64`

GetStatusCode returns the StatusCode field if non-nil, zero value otherwise.

### GetStatusCodeOk

`func (o *ErrorOAuth2) GetStatusCodeOk() (*int64, bool)`

GetStatusCodeOk returns a tuple with the StatusCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatusCode

`func (o *ErrorOAuth2) SetStatusCode(v int64)`

SetStatusCode sets StatusCode field to given value.

### HasStatusCode

`func (o *ErrorOAuth2) HasStatusCode() bool`

HasStatusCode returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


