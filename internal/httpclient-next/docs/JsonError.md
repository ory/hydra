# JsonError

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Error** | Pointer to **string** | Name is the error name. | [optional] 
**ErrorDebug** | Pointer to **string** | Debug contains debug information. This is usually not available and has to be enabled. | [optional] 
**ErrorDescription** | Pointer to **string** | Description contains further information on the nature of the error. | [optional] 
**StatusCode** | Pointer to **int64** | Code represents the error status code (404, 403, 401, ...). | [optional] 

## Methods

### NewJsonError

`func NewJsonError() *JsonError`

NewJsonError instantiates a new JsonError object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewJsonErrorWithDefaults

`func NewJsonErrorWithDefaults() *JsonError`

NewJsonErrorWithDefaults instantiates a new JsonError object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetError

`func (o *JsonError) GetError() string`

GetError returns the Error field if non-nil, zero value otherwise.

### GetErrorOk

`func (o *JsonError) GetErrorOk() (*string, bool)`

GetErrorOk returns a tuple with the Error field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetError

`func (o *JsonError) SetError(v string)`

SetError sets Error field to given value.

### HasError

`func (o *JsonError) HasError() bool`

HasError returns a boolean if a field has been set.

### GetErrorDebug

`func (o *JsonError) GetErrorDebug() string`

GetErrorDebug returns the ErrorDebug field if non-nil, zero value otherwise.

### GetErrorDebugOk

`func (o *JsonError) GetErrorDebugOk() (*string, bool)`

GetErrorDebugOk returns a tuple with the ErrorDebug field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrorDebug

`func (o *JsonError) SetErrorDebug(v string)`

SetErrorDebug sets ErrorDebug field to given value.

### HasErrorDebug

`func (o *JsonError) HasErrorDebug() bool`

HasErrorDebug returns a boolean if a field has been set.

### GetErrorDescription

`func (o *JsonError) GetErrorDescription() string`

GetErrorDescription returns the ErrorDescription field if non-nil, zero value otherwise.

### GetErrorDescriptionOk

`func (o *JsonError) GetErrorDescriptionOk() (*string, bool)`

GetErrorDescriptionOk returns a tuple with the ErrorDescription field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrorDescription

`func (o *JsonError) SetErrorDescription(v string)`

SetErrorDescription sets ErrorDescription field to given value.

### HasErrorDescription

`func (o *JsonError) HasErrorDescription() bool`

HasErrorDescription returns a boolean if a field has been set.

### GetStatusCode

`func (o *JsonError) GetStatusCode() int64`

GetStatusCode returns the StatusCode field if non-nil, zero value otherwise.

### GetStatusCodeOk

`func (o *JsonError) GetStatusCodeOk() (*int64, bool)`

GetStatusCodeOk returns a tuple with the StatusCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatusCode

`func (o *JsonError) SetStatusCode(v int64)`

SetStatusCode sets StatusCode field to given value.

### HasStatusCode

`func (o *JsonError) HasStatusCode() bool`

HasStatusCode returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


