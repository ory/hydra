# JSONWebKeySet

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Keys** | Pointer to [**[]JSONWebKey**](JSONWebKey.md) | The value of the \&quot;keys\&quot; parameter is an array of JWK values.  By default, the order of the JWK values within the array does not imply an order of preference among them, although applications of JWK Sets can choose to assign a meaning to the order for their purposes, if desired. | [optional] 

## Methods

### NewJSONWebKeySet

`func NewJSONWebKeySet() *JSONWebKeySet`

NewJSONWebKeySet instantiates a new JSONWebKeySet object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewJSONWebKeySetWithDefaults

`func NewJSONWebKeySetWithDefaults() *JSONWebKeySet`

NewJSONWebKeySetWithDefaults instantiates a new JSONWebKeySet object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKeys

`func (o *JSONWebKeySet) GetKeys() []JSONWebKey`

GetKeys returns the Keys field if non-nil, zero value otherwise.

### GetKeysOk

`func (o *JSONWebKeySet) GetKeysOk() (*[]JSONWebKey, bool)`

GetKeysOk returns a tuple with the Keys field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeys

`func (o *JSONWebKeySet) SetKeys(v []JSONWebKey)`

SetKeys sets Keys field to given value.

### HasKeys

`func (o *JSONWebKeySet) HasKeys() bool`

HasKeys returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


