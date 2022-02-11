# PluginSettings

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Args** | **[]string** | args | 
**Devices** | [**[]PluginDevice**](PluginDevice.md) | devices | 
**Env** | **[]string** | env | 
**Mounts** | [**[]PluginMount**](PluginMount.md) | mounts | 

## Methods

### NewPluginSettings

`func NewPluginSettings(args []string, devices []PluginDevice, env []string, mounts []PluginMount, ) *PluginSettings`

NewPluginSettings instantiates a new PluginSettings object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPluginSettingsWithDefaults

`func NewPluginSettingsWithDefaults() *PluginSettings`

NewPluginSettingsWithDefaults instantiates a new PluginSettings object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetArgs

`func (o *PluginSettings) GetArgs() []string`

GetArgs returns the Args field if non-nil, zero value otherwise.

### GetArgsOk

`func (o *PluginSettings) GetArgsOk() (*[]string, bool)`

GetArgsOk returns a tuple with the Args field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetArgs

`func (o *PluginSettings) SetArgs(v []string)`

SetArgs sets Args field to given value.


### GetDevices

`func (o *PluginSettings) GetDevices() []PluginDevice`

GetDevices returns the Devices field if non-nil, zero value otherwise.

### GetDevicesOk

`func (o *PluginSettings) GetDevicesOk() (*[]PluginDevice, bool)`

GetDevicesOk returns a tuple with the Devices field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDevices

`func (o *PluginSettings) SetDevices(v []PluginDevice)`

SetDevices sets Devices field to given value.


### GetEnv

`func (o *PluginSettings) GetEnv() []string`

GetEnv returns the Env field if non-nil, zero value otherwise.

### GetEnvOk

`func (o *PluginSettings) GetEnvOk() (*[]string, bool)`

GetEnvOk returns a tuple with the Env field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnv

`func (o *PluginSettings) SetEnv(v []string)`

SetEnv sets Env field to given value.


### GetMounts

`func (o *PluginSettings) GetMounts() []PluginMount`

GetMounts returns the Mounts field if non-nil, zero value otherwise.

### GetMountsOk

`func (o *PluginSettings) GetMountsOk() (*[]PluginMount, bool)`

GetMountsOk returns a tuple with the Mounts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMounts

`func (o *PluginSettings) SetMounts(v []PluginMount)`

SetMounts sets Mounts field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


