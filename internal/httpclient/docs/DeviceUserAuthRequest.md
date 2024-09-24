# DeviceUserAuthRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Challenge** | **string** | ID is the identifier (\&quot;device challenge\&quot;) of the device grant request. It is used to identify the session. | 
**Client** | Pointer to [**OAuth2Client**](OAuth2Client.md) |  | [optional] 
**HandledAt** | Pointer to **time.Time** |  | [optional] 
**RequestUrl** | Pointer to **string** | RequestURL is the original Device Authorization URL requested. | [optional] 
**RequestedAccessTokenAudience** | Pointer to **[]string** |  | [optional] 
**RequestedScope** | Pointer to **[]string** |  | [optional] 

## Methods

### NewDeviceUserAuthRequest

`func NewDeviceUserAuthRequest(challenge string, ) *DeviceUserAuthRequest`

NewDeviceUserAuthRequest instantiates a new DeviceUserAuthRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDeviceUserAuthRequestWithDefaults

`func NewDeviceUserAuthRequestWithDefaults() *DeviceUserAuthRequest`

NewDeviceUserAuthRequestWithDefaults instantiates a new DeviceUserAuthRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetChallenge

`func (o *DeviceUserAuthRequest) GetChallenge() string`

GetChallenge returns the Challenge field if non-nil, zero value otherwise.

### GetChallengeOk

`func (o *DeviceUserAuthRequest) GetChallengeOk() (*string, bool)`

GetChallengeOk returns a tuple with the Challenge field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChallenge

`func (o *DeviceUserAuthRequest) SetChallenge(v string)`

SetChallenge sets Challenge field to given value.


### GetClient

`func (o *DeviceUserAuthRequest) GetClient() OAuth2Client`

GetClient returns the Client field if non-nil, zero value otherwise.

### GetClientOk

`func (o *DeviceUserAuthRequest) GetClientOk() (*OAuth2Client, bool)`

GetClientOk returns a tuple with the Client field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClient

`func (o *DeviceUserAuthRequest) SetClient(v OAuth2Client)`

SetClient sets Client field to given value.

### HasClient

`func (o *DeviceUserAuthRequest) HasClient() bool`

HasClient returns a boolean if a field has been set.

### GetHandledAt

`func (o *DeviceUserAuthRequest) GetHandledAt() time.Time`

GetHandledAt returns the HandledAt field if non-nil, zero value otherwise.

### GetHandledAtOk

`func (o *DeviceUserAuthRequest) GetHandledAtOk() (*time.Time, bool)`

GetHandledAtOk returns a tuple with the HandledAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHandledAt

`func (o *DeviceUserAuthRequest) SetHandledAt(v time.Time)`

SetHandledAt sets HandledAt field to given value.

### HasHandledAt

`func (o *DeviceUserAuthRequest) HasHandledAt() bool`

HasHandledAt returns a boolean if a field has been set.

### GetRequestUrl

`func (o *DeviceUserAuthRequest) GetRequestUrl() string`

GetRequestUrl returns the RequestUrl field if non-nil, zero value otherwise.

### GetRequestUrlOk

`func (o *DeviceUserAuthRequest) GetRequestUrlOk() (*string, bool)`

GetRequestUrlOk returns a tuple with the RequestUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestUrl

`func (o *DeviceUserAuthRequest) SetRequestUrl(v string)`

SetRequestUrl sets RequestUrl field to given value.

### HasRequestUrl

`func (o *DeviceUserAuthRequest) HasRequestUrl() bool`

HasRequestUrl returns a boolean if a field has been set.

### GetRequestedAccessTokenAudience

`func (o *DeviceUserAuthRequest) GetRequestedAccessTokenAudience() []string`

GetRequestedAccessTokenAudience returns the RequestedAccessTokenAudience field if non-nil, zero value otherwise.

### GetRequestedAccessTokenAudienceOk

`func (o *DeviceUserAuthRequest) GetRequestedAccessTokenAudienceOk() (*[]string, bool)`

GetRequestedAccessTokenAudienceOk returns a tuple with the RequestedAccessTokenAudience field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestedAccessTokenAudience

`func (o *DeviceUserAuthRequest) SetRequestedAccessTokenAudience(v []string)`

SetRequestedAccessTokenAudience sets RequestedAccessTokenAudience field to given value.

### HasRequestedAccessTokenAudience

`func (o *DeviceUserAuthRequest) HasRequestedAccessTokenAudience() bool`

HasRequestedAccessTokenAudience returns a boolean if a field has been set.

### GetRequestedScope

`func (o *DeviceUserAuthRequest) GetRequestedScope() []string`

GetRequestedScope returns the RequestedScope field if non-nil, zero value otherwise.

### GetRequestedScopeOk

`func (o *DeviceUserAuthRequest) GetRequestedScopeOk() (*[]string, bool)`

GetRequestedScopeOk returns a tuple with the RequestedScope field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestedScope

`func (o *DeviceUserAuthRequest) SetRequestedScope(v []string)`

SetRequestedScope sets RequestedScope field to given value.

### HasRequestedScope

`func (o *DeviceUserAuthRequest) HasRequestedScope() bool`

HasRequestedScope returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


