# RFC6749ErrorJson

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Error** | Pointer to **string** |  | [optional] 
**ErrorDebug** | Pointer to **string** |  | [optional] 
**ErrorDescription** | Pointer to **string** |  | [optional] 
**ErrorHint** | Pointer to **string** |  | [optional] 
**StatusCode** | Pointer to **int64** |  | [optional] 

## Methods

### NewRFC6749ErrorJson

`func NewRFC6749ErrorJson() *RFC6749ErrorJson`

NewRFC6749ErrorJson instantiates a new RFC6749ErrorJson object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRFC6749ErrorJsonWithDefaults

`func NewRFC6749ErrorJsonWithDefaults() *RFC6749ErrorJson`

NewRFC6749ErrorJsonWithDefaults instantiates a new RFC6749ErrorJson object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetError

`func (o *RFC6749ErrorJson) GetError() string`

GetError returns the Error field if non-nil, zero value otherwise.

### GetErrorOk

`func (o *RFC6749ErrorJson) GetErrorOk() (*string, bool)`

GetErrorOk returns a tuple with the Error field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetError

`func (o *RFC6749ErrorJson) SetError(v string)`

SetError sets Error field to given value.

### HasError

`func (o *RFC6749ErrorJson) HasError() bool`

HasError returns a boolean if a field has been set.

### GetErrorDebug

`func (o *RFC6749ErrorJson) GetErrorDebug() string`

GetErrorDebug returns the ErrorDebug field if non-nil, zero value otherwise.

### GetErrorDebugOk

`func (o *RFC6749ErrorJson) GetErrorDebugOk() (*string, bool)`

GetErrorDebugOk returns a tuple with the ErrorDebug field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrorDebug

`func (o *RFC6749ErrorJson) SetErrorDebug(v string)`

SetErrorDebug sets ErrorDebug field to given value.

### HasErrorDebug

`func (o *RFC6749ErrorJson) HasErrorDebug() bool`

HasErrorDebug returns a boolean if a field has been set.

### GetErrorDescription

`func (o *RFC6749ErrorJson) GetErrorDescription() string`

GetErrorDescription returns the ErrorDescription field if non-nil, zero value otherwise.

### GetErrorDescriptionOk

`func (o *RFC6749ErrorJson) GetErrorDescriptionOk() (*string, bool)`

GetErrorDescriptionOk returns a tuple with the ErrorDescription field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrorDescription

`func (o *RFC6749ErrorJson) SetErrorDescription(v string)`

SetErrorDescription sets ErrorDescription field to given value.

### HasErrorDescription

`func (o *RFC6749ErrorJson) HasErrorDescription() bool`

HasErrorDescription returns a boolean if a field has been set.

### GetErrorHint

`func (o *RFC6749ErrorJson) GetErrorHint() string`

GetErrorHint returns the ErrorHint field if non-nil, zero value otherwise.

### GetErrorHintOk

`func (o *RFC6749ErrorJson) GetErrorHintOk() (*string, bool)`

GetErrorHintOk returns a tuple with the ErrorHint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrorHint

`func (o *RFC6749ErrorJson) SetErrorHint(v string)`

SetErrorHint sets ErrorHint field to given value.

### HasErrorHint

`func (o *RFC6749ErrorJson) HasErrorHint() bool`

HasErrorHint returns a boolean if a field has been set.

### GetStatusCode

`func (o *RFC6749ErrorJson) GetStatusCode() int64`

GetStatusCode returns the StatusCode field if non-nil, zero value otherwise.

### GetStatusCodeOk

`func (o *RFC6749ErrorJson) GetStatusCodeOk() (*int64, bool)`

GetStatusCodeOk returns a tuple with the StatusCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatusCode

`func (o *RFC6749ErrorJson) SetStatusCode(v int64)`

SetStatusCode sets StatusCode field to given value.

### HasStatusCode

`func (o *RFC6749ErrorJson) HasStatusCode() bool`

HasStatusCode returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


