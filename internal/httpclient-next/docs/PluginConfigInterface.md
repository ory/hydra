# PluginConfigInterface

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ProtocolScheme** | Pointer to **string** | Protocol to use for clients connecting to the plugin. | [optional] 
**Socket** | **string** | socket | 
**Types** | [**[]PluginInterfaceType**](PluginInterfaceType.md) | types | 

## Methods

### NewPluginConfigInterface

`func NewPluginConfigInterface(socket string, types []PluginInterfaceType, ) *PluginConfigInterface`

NewPluginConfigInterface instantiates a new PluginConfigInterface object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPluginConfigInterfaceWithDefaults

`func NewPluginConfigInterfaceWithDefaults() *PluginConfigInterface`

NewPluginConfigInterfaceWithDefaults instantiates a new PluginConfigInterface object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetProtocolScheme

`func (o *PluginConfigInterface) GetProtocolScheme() string`

GetProtocolScheme returns the ProtocolScheme field if non-nil, zero value otherwise.

### GetProtocolSchemeOk

`func (o *PluginConfigInterface) GetProtocolSchemeOk() (*string, bool)`

GetProtocolSchemeOk returns a tuple with the ProtocolScheme field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProtocolScheme

`func (o *PluginConfigInterface) SetProtocolScheme(v string)`

SetProtocolScheme sets ProtocolScheme field to given value.

### HasProtocolScheme

`func (o *PluginConfigInterface) HasProtocolScheme() bool`

HasProtocolScheme returns a boolean if a field has been set.

### GetSocket

`func (o *PluginConfigInterface) GetSocket() string`

GetSocket returns the Socket field if non-nil, zero value otherwise.

### GetSocketOk

`func (o *PluginConfigInterface) GetSocketOk() (*string, bool)`

GetSocketOk returns a tuple with the Socket field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSocket

`func (o *PluginConfigInterface) SetSocket(v string)`

SetSocket sets Socket field to given value.


### GetTypes

`func (o *PluginConfigInterface) GetTypes() []PluginInterfaceType`

GetTypes returns the Types field if non-nil, zero value otherwise.

### GetTypesOk

`func (o *PluginConfigInterface) GetTypesOk() (*[]PluginInterfaceType, bool)`

GetTypesOk returns a tuple with the Types field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTypes

`func (o *PluginConfigInterface) SetTypes(v []PluginInterfaceType)`

SetTypes sets Types field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


