# RejectOAuth2Request

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Error** | Pointer to **string** | The error should follow the OAuth2 error format (e.g. &#x60;invalid_request&#x60;, &#x60;login_required&#x60;).  Defaults to &#x60;request_denied&#x60;. | [optional] 
**ErrorDebug** | Pointer to **string** | Debug contains information to help resolve the problem as a developer. Usually not exposed to the public but only in the server logs. | [optional] 
**ErrorDescription** | Pointer to **string** | Description of the error in a human readable format. | [optional] 
**ErrorHint** | Pointer to **string** | Hint to help resolve the error. | [optional] 
**StatusCode** | Pointer to **int64** | Represents the HTTP status code of the error (e.g. 401 or 403)  Defaults to 400 | [optional] 

## Methods

### NewRejectOAuth2Request

`func NewRejectOAuth2Request() *RejectOAuth2Request`

NewRejectOAuth2Request instantiates a new RejectOAuth2Request object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRejectOAuth2RequestWithDefaults

`func NewRejectOAuth2RequestWithDefaults() *RejectOAuth2Request`

NewRejectOAuth2RequestWithDefaults instantiates a new RejectOAuth2Request object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetError

`func (o *RejectOAuth2Request) GetError() string`

GetError returns the Error field if non-nil, zero value otherwise.

### GetErrorOk

`func (o *RejectOAuth2Request) GetErrorOk() (*string, bool)`

GetErrorOk returns a tuple with the Error field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetError

`func (o *RejectOAuth2Request) SetError(v string)`

SetError sets Error field to given value.

### HasError

`func (o *RejectOAuth2Request) HasError() bool`

HasError returns a boolean if a field has been set.

### GetErrorDebug

`func (o *RejectOAuth2Request) GetErrorDebug() string`

GetErrorDebug returns the ErrorDebug field if non-nil, zero value otherwise.

### GetErrorDebugOk

`func (o *RejectOAuth2Request) GetErrorDebugOk() (*string, bool)`

GetErrorDebugOk returns a tuple with the ErrorDebug field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrorDebug

`func (o *RejectOAuth2Request) SetErrorDebug(v string)`

SetErrorDebug sets ErrorDebug field to given value.

### HasErrorDebug

`func (o *RejectOAuth2Request) HasErrorDebug() bool`

HasErrorDebug returns a boolean if a field has been set.

### GetErrorDescription

`func (o *RejectOAuth2Request) GetErrorDescription() string`

GetErrorDescription returns the ErrorDescription field if non-nil, zero value otherwise.

### GetErrorDescriptionOk

`func (o *RejectOAuth2Request) GetErrorDescriptionOk() (*string, bool)`

GetErrorDescriptionOk returns a tuple with the ErrorDescription field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrorDescription

`func (o *RejectOAuth2Request) SetErrorDescription(v string)`

SetErrorDescription sets ErrorDescription field to given value.

### HasErrorDescription

`func (o *RejectOAuth2Request) HasErrorDescription() bool`

HasErrorDescription returns a boolean if a field has been set.

### GetErrorHint

`func (o *RejectOAuth2Request) GetErrorHint() string`

GetErrorHint returns the ErrorHint field if non-nil, zero value otherwise.

### GetErrorHintOk

`func (o *RejectOAuth2Request) GetErrorHintOk() (*string, bool)`

GetErrorHintOk returns a tuple with the ErrorHint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrorHint

`func (o *RejectOAuth2Request) SetErrorHint(v string)`

SetErrorHint sets ErrorHint field to given value.

### HasErrorHint

`func (o *RejectOAuth2Request) HasErrorHint() bool`

HasErrorHint returns a boolean if a field has been set.

### GetStatusCode

`func (o *RejectOAuth2Request) GetStatusCode() int64`

GetStatusCode returns the StatusCode field if non-nil, zero value otherwise.

### GetStatusCodeOk

`func (o *RejectOAuth2Request) GetStatusCodeOk() (*int64, bool)`

GetStatusCodeOk returns a tuple with the StatusCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatusCode

`func (o *RejectOAuth2Request) SetStatusCode(v int64)`

SetStatusCode sets StatusCode field to given value.

### HasStatusCode

`func (o *RejectOAuth2Request) HasStatusCode() bool`

HasStatusCode returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


