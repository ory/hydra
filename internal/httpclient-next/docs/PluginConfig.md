# PluginConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Args** | [**PluginConfigArgs**](PluginConfigArgs.md) |  | 
**Description** | **string** | description | 
**DockerVersion** | Pointer to **string** | Docker Version used to create the plugin | [optional] 
**Documentation** | **string** | documentation | 
**Entrypoint** | **[]string** | entrypoint | 
**Env** | [**[]PluginEnv**](PluginEnv.md) | env | 
**Interface** | [**PluginConfigInterface**](PluginConfigInterface.md) |  | 
**IpcHost** | **bool** | ipc host | 
**Linux** | [**PluginConfigLinux**](PluginConfigLinux.md) |  | 
**Mounts** | [**[]PluginMount**](PluginMount.md) | mounts | 
**Network** | [**PluginConfigNetwork**](PluginConfigNetwork.md) |  | 
**PidHost** | **bool** | pid host | 
**PropagatedMount** | **string** | propagated mount | 
**User** | Pointer to [**PluginConfigUser**](PluginConfigUser.md) |  | [optional] 
**WorkDir** | **string** | work dir | 
**Rootfs** | Pointer to [**PluginConfigRootfs**](PluginConfigRootfs.md) |  | [optional] 

## Methods

### NewPluginConfig

`func NewPluginConfig(args PluginConfigArgs, description string, documentation string, entrypoint []string, env []PluginEnv, interface_ PluginConfigInterface, ipcHost bool, linux PluginConfigLinux, mounts []PluginMount, network PluginConfigNetwork, pidHost bool, propagatedMount string, workDir string, ) *PluginConfig`

NewPluginConfig instantiates a new PluginConfig object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPluginConfigWithDefaults

`func NewPluginConfigWithDefaults() *PluginConfig`

NewPluginConfigWithDefaults instantiates a new PluginConfig object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetArgs

`func (o *PluginConfig) GetArgs() PluginConfigArgs`

GetArgs returns the Args field if non-nil, zero value otherwise.

### GetArgsOk

`func (o *PluginConfig) GetArgsOk() (*PluginConfigArgs, bool)`

GetArgsOk returns a tuple with the Args field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetArgs

`func (o *PluginConfig) SetArgs(v PluginConfigArgs)`

SetArgs sets Args field to given value.


### GetDescription

`func (o *PluginConfig) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *PluginConfig) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *PluginConfig) SetDescription(v string)`

SetDescription sets Description field to given value.


### GetDockerVersion

`func (o *PluginConfig) GetDockerVersion() string`

GetDockerVersion returns the DockerVersion field if non-nil, zero value otherwise.

### GetDockerVersionOk

`func (o *PluginConfig) GetDockerVersionOk() (*string, bool)`

GetDockerVersionOk returns a tuple with the DockerVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDockerVersion

`func (o *PluginConfig) SetDockerVersion(v string)`

SetDockerVersion sets DockerVersion field to given value.

### HasDockerVersion

`func (o *PluginConfig) HasDockerVersion() bool`

HasDockerVersion returns a boolean if a field has been set.

### GetDocumentation

`func (o *PluginConfig) GetDocumentation() string`

GetDocumentation returns the Documentation field if non-nil, zero value otherwise.

### GetDocumentationOk

`func (o *PluginConfig) GetDocumentationOk() (*string, bool)`

GetDocumentationOk returns a tuple with the Documentation field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDocumentation

`func (o *PluginConfig) SetDocumentation(v string)`

SetDocumentation sets Documentation field to given value.


### GetEntrypoint

`func (o *PluginConfig) GetEntrypoint() []string`

GetEntrypoint returns the Entrypoint field if non-nil, zero value otherwise.

### GetEntrypointOk

`func (o *PluginConfig) GetEntrypointOk() (*[]string, bool)`

GetEntrypointOk returns a tuple with the Entrypoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEntrypoint

`func (o *PluginConfig) SetEntrypoint(v []string)`

SetEntrypoint sets Entrypoint field to given value.


### GetEnv

`func (o *PluginConfig) GetEnv() []PluginEnv`

GetEnv returns the Env field if non-nil, zero value otherwise.

### GetEnvOk

`func (o *PluginConfig) GetEnvOk() (*[]PluginEnv, bool)`

GetEnvOk returns a tuple with the Env field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnv

`func (o *PluginConfig) SetEnv(v []PluginEnv)`

SetEnv sets Env field to given value.


### GetInterface

`func (o *PluginConfig) GetInterface() PluginConfigInterface`

GetInterface returns the Interface field if non-nil, zero value otherwise.

### GetInterfaceOk

`func (o *PluginConfig) GetInterfaceOk() (*PluginConfigInterface, bool)`

GetInterfaceOk returns a tuple with the Interface field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInterface

`func (o *PluginConfig) SetInterface(v PluginConfigInterface)`

SetInterface sets Interface field to given value.


### GetIpcHost

`func (o *PluginConfig) GetIpcHost() bool`

GetIpcHost returns the IpcHost field if non-nil, zero value otherwise.

### GetIpcHostOk

`func (o *PluginConfig) GetIpcHostOk() (*bool, bool)`

GetIpcHostOk returns a tuple with the IpcHost field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpcHost

`func (o *PluginConfig) SetIpcHost(v bool)`

SetIpcHost sets IpcHost field to given value.


### GetLinux

`func (o *PluginConfig) GetLinux() PluginConfigLinux`

GetLinux returns the Linux field if non-nil, zero value otherwise.

### GetLinuxOk

`func (o *PluginConfig) GetLinuxOk() (*PluginConfigLinux, bool)`

GetLinuxOk returns a tuple with the Linux field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLinux

`func (o *PluginConfig) SetLinux(v PluginConfigLinux)`

SetLinux sets Linux field to given value.


### GetMounts

`func (o *PluginConfig) GetMounts() []PluginMount`

GetMounts returns the Mounts field if non-nil, zero value otherwise.

### GetMountsOk

`func (o *PluginConfig) GetMountsOk() (*[]PluginMount, bool)`

GetMountsOk returns a tuple with the Mounts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMounts

`func (o *PluginConfig) SetMounts(v []PluginMount)`

SetMounts sets Mounts field to given value.


### GetNetwork

`func (o *PluginConfig) GetNetwork() PluginConfigNetwork`

GetNetwork returns the Network field if non-nil, zero value otherwise.

### GetNetworkOk

`func (o *PluginConfig) GetNetworkOk() (*PluginConfigNetwork, bool)`

GetNetworkOk returns a tuple with the Network field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetwork

`func (o *PluginConfig) SetNetwork(v PluginConfigNetwork)`

SetNetwork sets Network field to given value.


### GetPidHost

`func (o *PluginConfig) GetPidHost() bool`

GetPidHost returns the PidHost field if non-nil, zero value otherwise.

### GetPidHostOk

`func (o *PluginConfig) GetPidHostOk() (*bool, bool)`

GetPidHostOk returns a tuple with the PidHost field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPidHost

`func (o *PluginConfig) SetPidHost(v bool)`

SetPidHost sets PidHost field to given value.


### GetPropagatedMount

`func (o *PluginConfig) GetPropagatedMount() string`

GetPropagatedMount returns the PropagatedMount field if non-nil, zero value otherwise.

### GetPropagatedMountOk

`func (o *PluginConfig) GetPropagatedMountOk() (*string, bool)`

GetPropagatedMountOk returns a tuple with the PropagatedMount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPropagatedMount

`func (o *PluginConfig) SetPropagatedMount(v string)`

SetPropagatedMount sets PropagatedMount field to given value.


### GetUser

`func (o *PluginConfig) GetUser() PluginConfigUser`

GetUser returns the User field if non-nil, zero value otherwise.

### GetUserOk

`func (o *PluginConfig) GetUserOk() (*PluginConfigUser, bool)`

GetUserOk returns a tuple with the User field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUser

`func (o *PluginConfig) SetUser(v PluginConfigUser)`

SetUser sets User field to given value.

### HasUser

`func (o *PluginConfig) HasUser() bool`

HasUser returns a boolean if a field has been set.

### GetWorkDir

`func (o *PluginConfig) GetWorkDir() string`

GetWorkDir returns the WorkDir field if non-nil, zero value otherwise.

### GetWorkDirOk

`func (o *PluginConfig) GetWorkDirOk() (*string, bool)`

GetWorkDirOk returns a tuple with the WorkDir field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkDir

`func (o *PluginConfig) SetWorkDir(v string)`

SetWorkDir sets WorkDir field to given value.


### GetRootfs

`func (o *PluginConfig) GetRootfs() PluginConfigRootfs`

GetRootfs returns the Rootfs field if non-nil, zero value otherwise.

### GetRootfsOk

`func (o *PluginConfig) GetRootfsOk() (*PluginConfigRootfs, bool)`

GetRootfsOk returns a tuple with the Rootfs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRootfs

`func (o *PluginConfig) SetRootfs(v PluginConfigRootfs)`

SetRootfs sets Rootfs field to given value.

### HasRootfs

`func (o *PluginConfig) HasRootfs() bool`

HasRootfs returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


