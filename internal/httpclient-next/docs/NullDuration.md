# NullDuration

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Duration** | Pointer to **int64** | A Duration represents the elapsed time between two instants as an int64 nanosecond count. The representation limits the largest representable duration to approximately 290 years. | [optional] 
**Valid** | Pointer to **bool** |  | [optional] 

## Methods

### NewNullDuration

`func NewNullDuration() *NullDuration`

NewNullDuration instantiates a new NullDuration object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNullDurationWithDefaults

`func NewNullDurationWithDefaults() *NullDuration`

NewNullDurationWithDefaults instantiates a new NullDuration object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDuration

`func (o *NullDuration) GetDuration() int64`

GetDuration returns the Duration field if non-nil, zero value otherwise.

### GetDurationOk

`func (o *NullDuration) GetDurationOk() (*int64, bool)`

GetDurationOk returns a tuple with the Duration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDuration

`func (o *NullDuration) SetDuration(v int64)`

SetDuration sets Duration field to given value.

### HasDuration

`func (o *NullDuration) HasDuration() bool`

HasDuration returns a boolean if a field has been set.

### GetValid

`func (o *NullDuration) GetValid() bool`

GetValid returns the Valid field if non-nil, zero value otherwise.

### GetValidOk

`func (o *NullDuration) GetValidOk() (*bool, bool)`

GetValidOk returns a tuple with the Valid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValid

`func (o *NullDuration) SetValid(v bool)`

SetValid sets Valid field to given value.

### HasValid

`func (o *NullDuration) HasValid() bool`

HasValid returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


