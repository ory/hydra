# IsReady503Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Errors** | Pointer to **map[string]string** | Errors contains a list of errors that caused the not ready status. | [optional] 

## Methods

### NewIsReady503Response

`func NewIsReady503Response() *IsReady503Response`

NewIsReady503Response instantiates a new IsReady503Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewIsReady503ResponseWithDefaults

`func NewIsReady503ResponseWithDefaults() *IsReady503Response`

NewIsReady503ResponseWithDefaults instantiates a new IsReady503Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetErrors

`func (o *IsReady503Response) GetErrors() map[string]string`

GetErrors returns the Errors field if non-nil, zero value otherwise.

### GetErrorsOk

`func (o *IsReady503Response) GetErrorsOk() (*map[string]string, bool)`

GetErrorsOk returns a tuple with the Errors field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrors

`func (o *IsReady503Response) SetErrors(v map[string]string)`

SetErrors sets Errors field to given value.

### HasErrors

`func (o *IsReady503Response) HasErrors() bool`

HasErrors returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


