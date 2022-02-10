# VolumeUsageData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**RefCount** | **int64** | The number of containers referencing this volume. This field is set to &#x60;-1&#x60; if the reference-count is not available. | 
**Size** | **int64** | Amount of disk space used by the volume (in bytes). This information is only available for volumes created with the &#x60;\&quot;local\&quot;&#x60; volume driver. For volumes created with other volume drivers, this field is set to &#x60;-1&#x60; (\&quot;not available\&quot;) | 

## Methods

### NewVolumeUsageData

`func NewVolumeUsageData(refCount int64, size int64, ) *VolumeUsageData`

NewVolumeUsageData instantiates a new VolumeUsageData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVolumeUsageDataWithDefaults

`func NewVolumeUsageDataWithDefaults() *VolumeUsageData`

NewVolumeUsageDataWithDefaults instantiates a new VolumeUsageData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRefCount

`func (o *VolumeUsageData) GetRefCount() int64`

GetRefCount returns the RefCount field if non-nil, zero value otherwise.

### GetRefCountOk

`func (o *VolumeUsageData) GetRefCountOk() (*int64, bool)`

GetRefCountOk returns a tuple with the RefCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRefCount

`func (o *VolumeUsageData) SetRefCount(v int64)`

SetRefCount sets RefCount field to given value.


### GetSize

`func (o *VolumeUsageData) GetSize() int64`

GetSize returns the Size field if non-nil, zero value otherwise.

### GetSizeOk

`func (o *VolumeUsageData) GetSizeOk() (*int64, bool)`

GetSizeOk returns a tuple with the Size field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSize

`func (o *VolumeUsageData) SetSize(v int64)`

SetSize sets Size field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


