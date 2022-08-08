# Headers

## Properties

| Name      | Type                                  | Description | Notes      |
| --------- | ------------------------------------- | ----------- | ---------- |
| **Extra** | Pointer to **map[string]interface{}** |             | [optional] |

## Methods

### NewHeaders

`func NewHeaders() *Headers`

NewHeaders instantiates a new Headers object This constructor will assign
default values to properties that have it defined, and makes sure properties
required by API are set, but the set of arguments will change when the set of
required properties is changed

### NewHeadersWithDefaults

`func NewHeadersWithDefaults() *Headers`

NewHeadersWithDefaults instantiates a new Headers object This constructor will
only assign default values to properties that have it defined, but it doesn't
guarantee that properties required by API are set

### GetExtra

`func (o *Headers) GetExtra() map[string]interface{}`

GetExtra returns the Extra field if non-nil, zero value otherwise.

### GetExtraOk

`func (o *Headers) GetExtraOk() (*map[string]interface{}, bool)`

GetExtraOk returns a tuple with the Extra field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetExtra

`func (o *Headers) SetExtra(v map[string]interface{})`

SetExtra sets Extra field to given value.

### HasExtra

`func (o *Headers) HasExtra() bool`

HasExtra returns a boolean if a field has been set.

[[Back to Model list]](../README.md#documentation-for-models)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to README]](../README.md)
