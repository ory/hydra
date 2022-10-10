# ErrorBody

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

### NewErrorBody

`func NewErrorBody(message string, ) *ErrorBody`

NewErrorBody instantiates a new ErrorBody object This constructor will assign
default values to properties that have it defined, and makes sure properties
required by API are set, but the set of arguments will change when the set of
required properties is changed

### NewErrorBodyWithDefaults

`func NewErrorBodyWithDefaults() *ErrorBody`

NewErrorBodyWithDefaults instantiates a new ErrorBody object This constructor
will only assign default values to properties that have it defined, but it
doesn't guarantee that properties required by API are set

### GetCode

`func (o *ErrorBody) GetCode() int64`

GetCode returns the Code field if non-nil, zero value otherwise.

### GetCodeOk

`func (o *ErrorBody) GetCodeOk() (*int64, bool)`

GetCodeOk returns a tuple with the Code field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetCode

`func (o *ErrorBody) SetCode(v int64)`

SetCode sets Code field to given value.

### HasCode

`func (o *ErrorBody) HasCode() bool`

HasCode returns a boolean if a field has been set.

### GetDebug

`func (o *ErrorBody) GetDebug() string`

GetDebug returns the Debug field if non-nil, zero value otherwise.

### GetDebugOk

`func (o *ErrorBody) GetDebugOk() (*string, bool)`

GetDebugOk returns a tuple with the Debug field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetDebug

`func (o *ErrorBody) SetDebug(v string)`

SetDebug sets Debug field to given value.

### HasDebug

`func (o *ErrorBody) HasDebug() bool`

HasDebug returns a boolean if a field has been set.

### GetDetails

`func (o *ErrorBody) GetDetails() map[string]interface{}`

GetDetails returns the Details field if non-nil, zero value otherwise.

### GetDetailsOk

`func (o *ErrorBody) GetDetailsOk() (*map[string]interface{}, bool)`

GetDetailsOk returns a tuple with the Details field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetDetails

`func (o *ErrorBody) SetDetails(v map[string]interface{})`

SetDetails sets Details field to given value.

### HasDetails

`func (o *ErrorBody) HasDetails() bool`

HasDetails returns a boolean if a field has been set.

### GetId

`func (o *ErrorBody) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ErrorBody) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ErrorBody) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *ErrorBody) HasId() bool`

HasId returns a boolean if a field has been set.

### GetMessage

`func (o *ErrorBody) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *ErrorBody) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetMessage

`func (o *ErrorBody) SetMessage(v string)`

SetMessage sets Message field to given value.

### GetReason

`func (o *ErrorBody) GetReason() string`

GetReason returns the Reason field if non-nil, zero value otherwise.

### GetReasonOk

`func (o *ErrorBody) GetReasonOk() (*string, bool)`

GetReasonOk returns a tuple with the Reason field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetReason

`func (o *ErrorBody) SetReason(v string)`

SetReason sets Reason field to given value.

### HasReason

`func (o *ErrorBody) HasReason() bool`

HasReason returns a boolean if a field has been set.

### GetRequest

`func (o *ErrorBody) GetRequest() string`

GetRequest returns the Request field if non-nil, zero value otherwise.

### GetRequestOk

`func (o *ErrorBody) GetRequestOk() (*string, bool)`

GetRequestOk returns a tuple with the Request field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetRequest

`func (o *ErrorBody) SetRequest(v string)`

SetRequest sets Request field to given value.

### HasRequest

`func (o *ErrorBody) HasRequest() bool`

HasRequest returns a boolean if a field has been set.

### GetStatus

`func (o *ErrorBody) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ErrorBody) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetStatus

`func (o *ErrorBody) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *ErrorBody) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

[[Back to Model list]](../README.md#documentation-for-models)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to README]](../README.md)
