# Volume

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CreatedAt** | Pointer to **string** | Date/Time the volume was created. | [optional] 
**Driver** | **string** | Name of the volume driver used by the volume. | 
**Labels** | **map[string]string** | User-defined key/value metadata. | 
**Mountpoint** | **string** | Mount path of the volume on the host. | 
**Name** | **string** | Name of the volume. | 
**Options** | **map[string]string** | The driver specific options used when creating the volume. | 
**Scope** | **string** | The level at which the volume exists. Either &#x60;global&#x60; for cluster-wide, or &#x60;local&#x60; for machine level. | 
**Status** | Pointer to **map[string]interface{}** | Low-level details about the volume, provided by the volume driver. Details are returned as a map with key/value pairs: &#x60;{\&quot;key\&quot;:\&quot;value\&quot;,\&quot;key2\&quot;:\&quot;value2\&quot;}&#x60;.  The &#x60;Status&#x60; field is optional, and is omitted if the volume driver does not support this feature. | [optional] 
**UsageData** | Pointer to [**VolumeUsageData**](VolumeUsageData.md) |  | [optional] 

## Methods

### NewVolume

`func NewVolume(driver string, labels map[string]string, mountpoint string, name string, options map[string]string, scope string, ) *Volume`

NewVolume instantiates a new Volume object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVolumeWithDefaults

`func NewVolumeWithDefaults() *Volume`

NewVolumeWithDefaults instantiates a new Volume object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCreatedAt

`func (o *Volume) GetCreatedAt() string`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *Volume) GetCreatedAtOk() (*string, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *Volume) SetCreatedAt(v string)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *Volume) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetDriver

`func (o *Volume) GetDriver() string`

GetDriver returns the Driver field if non-nil, zero value otherwise.

### GetDriverOk

`func (o *Volume) GetDriverOk() (*string, bool)`

GetDriverOk returns a tuple with the Driver field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDriver

`func (o *Volume) SetDriver(v string)`

SetDriver sets Driver field to given value.


### GetLabels

`func (o *Volume) GetLabels() map[string]string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *Volume) GetLabelsOk() (*map[string]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *Volume) SetLabels(v map[string]string)`

SetLabels sets Labels field to given value.


### GetMountpoint

`func (o *Volume) GetMountpoint() string`

GetMountpoint returns the Mountpoint field if non-nil, zero value otherwise.

### GetMountpointOk

`func (o *Volume) GetMountpointOk() (*string, bool)`

GetMountpointOk returns a tuple with the Mountpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMountpoint

`func (o *Volume) SetMountpoint(v string)`

SetMountpoint sets Mountpoint field to given value.


### GetName

`func (o *Volume) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Volume) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Volume) SetName(v string)`

SetName sets Name field to given value.


### GetOptions

`func (o *Volume) GetOptions() map[string]string`

GetOptions returns the Options field if non-nil, zero value otherwise.

### GetOptionsOk

`func (o *Volume) GetOptionsOk() (*map[string]string, bool)`

GetOptionsOk returns a tuple with the Options field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOptions

`func (o *Volume) SetOptions(v map[string]string)`

SetOptions sets Options field to given value.


### GetScope

`func (o *Volume) GetScope() string`

GetScope returns the Scope field if non-nil, zero value otherwise.

### GetScopeOk

`func (o *Volume) GetScopeOk() (*string, bool)`

GetScopeOk returns a tuple with the Scope field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScope

`func (o *Volume) SetScope(v string)`

SetScope sets Scope field to given value.


### GetStatus

`func (o *Volume) GetStatus() map[string]interface{}`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *Volume) GetStatusOk() (*map[string]interface{}, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *Volume) SetStatus(v map[string]interface{})`

SetStatus sets Status field to given value.

### HasStatus

`func (o *Volume) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetUsageData

`func (o *Volume) GetUsageData() VolumeUsageData`

GetUsageData returns the UsageData field if non-nil, zero value otherwise.

### GetUsageDataOk

`func (o *Volume) GetUsageDataOk() (*VolumeUsageData, bool)`

GetUsageDataOk returns a tuple with the UsageData field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsageData

`func (o *Volume) SetUsageData(v VolumeUsageData)`

SetUsageData sets UsageData field to given value.

### HasUsageData

`func (o *Volume) HasUsageData() bool`

HasUsageData returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


