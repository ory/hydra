# PluginConfigLinux

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AllowAllDevices** | **bool** | allow all devices | 
**Capabilities** | **[]string** | capabilities | 
**Devices** | [**[]PluginDevice**](PluginDevice.md) | devices | 

## Methods

### NewPluginConfigLinux

`func NewPluginConfigLinux(allowAllDevices bool, capabilities []string, devices []PluginDevice, ) *PluginConfigLinux`

NewPluginConfigLinux instantiates a new PluginConfigLinux object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPluginConfigLinuxWithDefaults

`func NewPluginConfigLinuxWithDefaults() *PluginConfigLinux`

NewPluginConfigLinuxWithDefaults instantiates a new PluginConfigLinux object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAllowAllDevices

`func (o *PluginConfigLinux) GetAllowAllDevices() bool`

GetAllowAllDevices returns the AllowAllDevices field if non-nil, zero value otherwise.

### GetAllowAllDevicesOk

`func (o *PluginConfigLinux) GetAllowAllDevicesOk() (*bool, bool)`

GetAllowAllDevicesOk returns a tuple with the AllowAllDevices field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllowAllDevices

`func (o *PluginConfigLinux) SetAllowAllDevices(v bool)`

SetAllowAllDevices sets AllowAllDevices field to given value.


### GetCapabilities

`func (o *PluginConfigLinux) GetCapabilities() []string`

GetCapabilities returns the Capabilities field if non-nil, zero value otherwise.

### GetCapabilitiesOk

`func (o *PluginConfigLinux) GetCapabilitiesOk() (*[]string, bool)`

GetCapabilitiesOk returns a tuple with the Capabilities field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCapabilities

`func (o *PluginConfigLinux) SetCapabilities(v []string)`

SetCapabilities sets Capabilities field to given value.


### GetDevices

`func (o *PluginConfigLinux) GetDevices() []PluginDevice`

GetDevices returns the Devices field if non-nil, zero value otherwise.

### GetDevicesOk

`func (o *PluginConfigLinux) GetDevicesOk() (*[]PluginDevice, bool)`

GetDevicesOk returns a tuple with the Devices field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDevices

`func (o *PluginConfigLinux) SetDevices(v []PluginDevice)`

SetDevices sets Devices field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


