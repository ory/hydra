# VerifiableCredentialPrimingResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CNonce** | Pointer to **string** |  | [optional] 
**CNonceExpiresIn** | Pointer to **int64** |  | [optional] 
**Error** | Pointer to **string** |  | [optional] 
**ErrorDebug** | Pointer to **string** |  | [optional] 
**ErrorDescription** | Pointer to **string** |  | [optional] 
**ErrorHint** | Pointer to **string** |  | [optional] 
**Format** | Pointer to **string** |  | [optional] 
**StatusCode** | Pointer to **int64** |  | [optional] 

## Methods

### NewVerifiableCredentialPrimingResponse

`func NewVerifiableCredentialPrimingResponse() *VerifiableCredentialPrimingResponse`

NewVerifiableCredentialPrimingResponse instantiates a new VerifiableCredentialPrimingResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVerifiableCredentialPrimingResponseWithDefaults

`func NewVerifiableCredentialPrimingResponseWithDefaults() *VerifiableCredentialPrimingResponse`

NewVerifiableCredentialPrimingResponseWithDefaults instantiates a new VerifiableCredentialPrimingResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCNonce

`func (o *VerifiableCredentialPrimingResponse) GetCNonce() string`

GetCNonce returns the CNonce field if non-nil, zero value otherwise.

### GetCNonceOk

`func (o *VerifiableCredentialPrimingResponse) GetCNonceOk() (*string, bool)`

GetCNonceOk returns a tuple with the CNonce field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCNonce

`func (o *VerifiableCredentialPrimingResponse) SetCNonce(v string)`

SetCNonce sets CNonce field to given value.

### HasCNonce

`func (o *VerifiableCredentialPrimingResponse) HasCNonce() bool`

HasCNonce returns a boolean if a field has been set.

### GetCNonceExpiresIn

`func (o *VerifiableCredentialPrimingResponse) GetCNonceExpiresIn() int64`

GetCNonceExpiresIn returns the CNonceExpiresIn field if non-nil, zero value otherwise.

### GetCNonceExpiresInOk

`func (o *VerifiableCredentialPrimingResponse) GetCNonceExpiresInOk() (*int64, bool)`

GetCNonceExpiresInOk returns a tuple with the CNonceExpiresIn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCNonceExpiresIn

`func (o *VerifiableCredentialPrimingResponse) SetCNonceExpiresIn(v int64)`

SetCNonceExpiresIn sets CNonceExpiresIn field to given value.

### HasCNonceExpiresIn

`func (o *VerifiableCredentialPrimingResponse) HasCNonceExpiresIn() bool`

HasCNonceExpiresIn returns a boolean if a field has been set.

### GetError

`func (o *VerifiableCredentialPrimingResponse) GetError() string`

GetError returns the Error field if non-nil, zero value otherwise.

### GetErrorOk

`func (o *VerifiableCredentialPrimingResponse) GetErrorOk() (*string, bool)`

GetErrorOk returns a tuple with the Error field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetError

`func (o *VerifiableCredentialPrimingResponse) SetError(v string)`

SetError sets Error field to given value.

### HasError

`func (o *VerifiableCredentialPrimingResponse) HasError() bool`

HasError returns a boolean if a field has been set.

### GetErrorDebug

`func (o *VerifiableCredentialPrimingResponse) GetErrorDebug() string`

GetErrorDebug returns the ErrorDebug field if non-nil, zero value otherwise.

### GetErrorDebugOk

`func (o *VerifiableCredentialPrimingResponse) GetErrorDebugOk() (*string, bool)`

GetErrorDebugOk returns a tuple with the ErrorDebug field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrorDebug

`func (o *VerifiableCredentialPrimingResponse) SetErrorDebug(v string)`

SetErrorDebug sets ErrorDebug field to given value.

### HasErrorDebug

`func (o *VerifiableCredentialPrimingResponse) HasErrorDebug() bool`

HasErrorDebug returns a boolean if a field has been set.

### GetErrorDescription

`func (o *VerifiableCredentialPrimingResponse) GetErrorDescription() string`

GetErrorDescription returns the ErrorDescription field if non-nil, zero value otherwise.

### GetErrorDescriptionOk

`func (o *VerifiableCredentialPrimingResponse) GetErrorDescriptionOk() (*string, bool)`

GetErrorDescriptionOk returns a tuple with the ErrorDescription field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrorDescription

`func (o *VerifiableCredentialPrimingResponse) SetErrorDescription(v string)`

SetErrorDescription sets ErrorDescription field to given value.

### HasErrorDescription

`func (o *VerifiableCredentialPrimingResponse) HasErrorDescription() bool`

HasErrorDescription returns a boolean if a field has been set.

### GetErrorHint

`func (o *VerifiableCredentialPrimingResponse) GetErrorHint() string`

GetErrorHint returns the ErrorHint field if non-nil, zero value otherwise.

### GetErrorHintOk

`func (o *VerifiableCredentialPrimingResponse) GetErrorHintOk() (*string, bool)`

GetErrorHintOk returns a tuple with the ErrorHint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrorHint

`func (o *VerifiableCredentialPrimingResponse) SetErrorHint(v string)`

SetErrorHint sets ErrorHint field to given value.

### HasErrorHint

`func (o *VerifiableCredentialPrimingResponse) HasErrorHint() bool`

HasErrorHint returns a boolean if a field has been set.

### GetFormat

`func (o *VerifiableCredentialPrimingResponse) GetFormat() string`

GetFormat returns the Format field if non-nil, zero value otherwise.

### GetFormatOk

`func (o *VerifiableCredentialPrimingResponse) GetFormatOk() (*string, bool)`

GetFormatOk returns a tuple with the Format field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFormat

`func (o *VerifiableCredentialPrimingResponse) SetFormat(v string)`

SetFormat sets Format field to given value.

### HasFormat

`func (o *VerifiableCredentialPrimingResponse) HasFormat() bool`

HasFormat returns a boolean if a field has been set.

### GetStatusCode

`func (o *VerifiableCredentialPrimingResponse) GetStatusCode() int64`

GetStatusCode returns the StatusCode field if non-nil, zero value otherwise.

### GetStatusCodeOk

`func (o *VerifiableCredentialPrimingResponse) GetStatusCodeOk() (*int64, bool)`

GetStatusCodeOk returns a tuple with the StatusCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatusCode

`func (o *VerifiableCredentialPrimingResponse) SetStatusCode(v int64)`

SetStatusCode sets StatusCode field to given value.

### HasStatusCode

`func (o *VerifiableCredentialPrimingResponse) HasStatusCode() bool`

HasStatusCode returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


