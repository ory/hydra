# GenericError

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Code** | Pointer to **int64** | The status code | [optional] 
**Debug** | Pointer to **string** | Debug information  This field is often not exposed to protect against leaking sensitive information. | [optional] 
**Details** | Pointer to **interface{}** | Further error details | [optional] 
**Id** | Pointer to **string** | The error ID  Useful when trying to identify various errors in application logic. | [optional] 
**Message** | **string** | Error message  The error&#39;s message. | 
**Reason** | Pointer to **string** | A human-readable reason for the error | [optional] 
**Request** | Pointer to **string** | The request ID  The request ID is often exposed internally in order to trace errors across service architectures. This is often a UUID. | [optional] 
**Status** | Pointer to **string** | The status description | [optional] 

## Methods

### NewGenericError

`func NewGenericError(message string, ) *GenericError`

NewGenericError instantiates a new GenericError object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGenericErrorWithDefaults

`func NewGenericErrorWithDefaults() *GenericError`

NewGenericErrorWithDefaults instantiates a new GenericError object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCode

`func (o *GenericError) GetCode() int64`

GetCode returns the Code field if non-nil, zero value otherwise.

### GetCodeOk

`func (o *GenericError) GetCodeOk() (*int64, bool)`

GetCodeOk returns a tuple with the Code field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCode

`func (o *GenericError) SetCode(v int64)`

SetCode sets Code field to given value.

### HasCode

`func (o *GenericError) HasCode() bool`

HasCode returns a boolean if a field has been set.

### GetDebug

`func (o *GenericError) GetDebug() string`

GetDebug returns the Debug field if non-nil, zero value otherwise.

### GetDebugOk

`func (o *GenericError) GetDebugOk() (*string, bool)`

GetDebugOk returns a tuple with the Debug field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDebug

`func (o *GenericError) SetDebug(v string)`

SetDebug sets Debug field to given value.

### HasDebug

`func (o *GenericError) HasDebug() bool`

HasDebug returns a boolean if a field has been set.

### GetDetails

`func (o *GenericError) GetDetails() interface{}`

GetDetails returns the Details field if non-nil, zero value otherwise.

### GetDetailsOk

`func (o *GenericError) GetDetailsOk() (*interface{}, bool)`

GetDetailsOk returns a tuple with the Details field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDetails

`func (o *GenericError) SetDetails(v interface{})`

SetDetails sets Details field to given value.

### HasDetails

`func (o *GenericError) HasDetails() bool`

HasDetails returns a boolean if a field has been set.

### SetDetailsNil

`func (o *GenericError) SetDetailsNil(b bool)`

 SetDetailsNil sets the value for Details to be an explicit nil

### UnsetDetails
`func (o *GenericError) UnsetDetails()`

UnsetDetails ensures that no value is present for Details, not even an explicit nil
### GetId

`func (o *GenericError) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *GenericError) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *GenericError) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *GenericError) HasId() bool`

HasId returns a boolean if a field has been set.

### GetMessage

`func (o *GenericError) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *GenericError) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *GenericError) SetMessage(v string)`

SetMessage sets Message field to given value.


### GetReason

`func (o *GenericError) GetReason() string`

GetReason returns the Reason field if non-nil, zero value otherwise.

### GetReasonOk

`func (o *GenericError) GetReasonOk() (*string, bool)`

GetReasonOk returns a tuple with the Reason field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReason

`func (o *GenericError) SetReason(v string)`

SetReason sets Reason field to given value.

### HasReason

`func (o *GenericError) HasReason() bool`

HasReason returns a boolean if a field has been set.

### GetRequest

`func (o *GenericError) GetRequest() string`

GetRequest returns the Request field if non-nil, zero value otherwise.

### GetRequestOk

`func (o *GenericError) GetRequestOk() (*string, bool)`

GetRequestOk returns a tuple with the Request field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequest

`func (o *GenericError) SetRequest(v string)`

SetRequest sets Request field to given value.

### HasRequest

`func (o *GenericError) HasRequest() bool`

HasRequest returns a boolean if a field has been set.

### GetStatus

`func (o *GenericError) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *GenericError) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *GenericError) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *GenericError) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


