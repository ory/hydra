# ModelError

## Properties

| Name        | Type                                  | Description                                                                                                                             | Notes      |
| ----------- | ------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------- | ---------- |
| **Code**    | Pointer to **int64**                  | HTTP Status Code                                                                                                                        | [optional] |
| **Debug**   | Pointer to **string**                 | Debug Details This field is often not exposed to protect against leaking sensitive information.                                         | [optional] |
| **Details** | Pointer to **map[string]interface{}** | Additional Error Details Further error details                                                                                          | [optional] |
| **Id**      | Pointer to **string**                 | Error ID Useful when trying to identify various errors in application logic.                                                            | [optional] |
| **Message** | **string**                            | Error Message The error&#39;s message.                                                                                                  |
| **Reason**  | Pointer to **string**                 | Error Reason                                                                                                                            | [optional] |
| **Request** | Pointer to **string**                 | HTTP Request ID The request ID is often exposed internally in order to trace errors across service architectures. This is often a UUID. | [optional] |
| **Status**  | Pointer to **string**                 | HTTP Status Description                                                                                                                 | [optional] |

## Methods

### NewModelError

`func NewModelError(message string, ) *ModelError`

NewModelError instantiates a new ModelError object This constructor will assign
default values to properties that have it defined, and makes sure properties
required by API are set, but the set of arguments will change when the set of
required properties is changed

### NewModelErrorWithDefaults

`func NewModelErrorWithDefaults() *ModelError`

NewModelErrorWithDefaults instantiates a new ModelError object This constructor
will only assign default values to properties that have it defined, but it
doesn't guarantee that properties required by API are set

### GetCode

`func (o *ModelError) GetCode() int64`

GetCode returns the Code field if non-nil, zero value otherwise.

### GetCodeOk

`func (o *ModelError) GetCodeOk() (*int64, bool)`

GetCodeOk returns a tuple with the Code field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetCode

`func (o *ModelError) SetCode(v int64)`

SetCode sets Code field to given value.

### HasCode

`func (o *ModelError) HasCode() bool`

HasCode returns a boolean if a field has been set.

### GetDebug

`func (o *ModelError) GetDebug() string`

GetDebug returns the Debug field if non-nil, zero value otherwise.

### GetDebugOk

`func (o *ModelError) GetDebugOk() (*string, bool)`

GetDebugOk returns a tuple with the Debug field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetDebug

`func (o *ModelError) SetDebug(v string)`

SetDebug sets Debug field to given value.

### HasDebug

`func (o *ModelError) HasDebug() bool`

HasDebug returns a boolean if a field has been set.

### GetDetails

`func (o *ModelError) GetDetails() map[string]interface{}`

GetDetails returns the Details field if non-nil, zero value otherwise.

### GetDetailsOk

`func (o *ModelError) GetDetailsOk() (*map[string]interface{}, bool)`

GetDetailsOk returns a tuple with the Details field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetDetails

`func (o *ModelError) SetDetails(v map[string]interface{})`

SetDetails sets Details field to given value.

### HasDetails

`func (o *ModelError) HasDetails() bool`

HasDetails returns a boolean if a field has been set.

### GetId

`func (o *ModelError) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ModelError) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ModelError) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *ModelError) HasId() bool`

HasId returns a boolean if a field has been set.

### GetMessage

`func (o *ModelError) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *ModelError) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetMessage

`func (o *ModelError) SetMessage(v string)`

SetMessage sets Message field to given value.

### GetReason

`func (o *ModelError) GetReason() string`

GetReason returns the Reason field if non-nil, zero value otherwise.

### GetReasonOk

`func (o *ModelError) GetReasonOk() (*string, bool)`

GetReasonOk returns a tuple with the Reason field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetReason

`func (o *ModelError) SetReason(v string)`

SetReason sets Reason field to given value.

### HasReason

`func (o *ModelError) HasReason() bool`

HasReason returns a boolean if a field has been set.

### GetRequest

`func (o *ModelError) GetRequest() string`

GetRequest returns the Request field if non-nil, zero value otherwise.

### GetRequestOk

`func (o *ModelError) GetRequestOk() (*string, bool)`

GetRequestOk returns a tuple with the Request field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetRequest

`func (o *ModelError) SetRequest(v string)`

SetRequest sets Request field to given value.

### HasRequest

`func (o *ModelError) HasRequest() bool`

HasRequest returns a boolean if a field has been set.

### GetStatus

`func (o *ModelError) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ModelError) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetStatus

`func (o *ModelError) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *ModelError) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

[[Back to Model list]](../README.md#documentation-for-models)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to README]](../README.md)
