# VerifyUserCodeRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Client** | Pointer to [**OAuth2Client**](OAuth2Client.md) |  | [optional] 
**DeviceCodeRequestId** | Pointer to **string** |  | [optional] 
**RequestUrl** | Pointer to **string** | RequestURL is the original Device Authorization URL requested. | [optional] 
**RequestedAccessTokenAudience** | Pointer to **[]string** |  | [optional] 
**RequestedScope** | Pointer to **[]string** |  | [optional] 

## Methods

### NewVerifyUserCodeRequest

`func NewVerifyUserCodeRequest() *VerifyUserCodeRequest`

NewVerifyUserCodeRequest instantiates a new VerifyUserCodeRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVerifyUserCodeRequestWithDefaults

`func NewVerifyUserCodeRequestWithDefaults() *VerifyUserCodeRequest`

NewVerifyUserCodeRequestWithDefaults instantiates a new VerifyUserCodeRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetClient

`func (o *VerifyUserCodeRequest) GetClient() OAuth2Client`

GetClient returns the Client field if non-nil, zero value otherwise.

### GetClientOk

`func (o *VerifyUserCodeRequest) GetClientOk() (*OAuth2Client, bool)`

GetClientOk returns a tuple with the Client field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClient

`func (o *VerifyUserCodeRequest) SetClient(v OAuth2Client)`

SetClient sets Client field to given value.

### HasClient

`func (o *VerifyUserCodeRequest) HasClient() bool`

HasClient returns a boolean if a field has been set.

### GetDeviceCodeRequestId

`func (o *VerifyUserCodeRequest) GetDeviceCodeRequestId() string`

GetDeviceCodeRequestId returns the DeviceCodeRequestId field if non-nil, zero value otherwise.

### GetDeviceCodeRequestIdOk

`func (o *VerifyUserCodeRequest) GetDeviceCodeRequestIdOk() (*string, bool)`

GetDeviceCodeRequestIdOk returns a tuple with the DeviceCodeRequestId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceCodeRequestId

`func (o *VerifyUserCodeRequest) SetDeviceCodeRequestId(v string)`

SetDeviceCodeRequestId sets DeviceCodeRequestId field to given value.

### HasDeviceCodeRequestId

`func (o *VerifyUserCodeRequest) HasDeviceCodeRequestId() bool`

HasDeviceCodeRequestId returns a boolean if a field has been set.

### GetRequestUrl

`func (o *VerifyUserCodeRequest) GetRequestUrl() string`

GetRequestUrl returns the RequestUrl field if non-nil, zero value otherwise.

### GetRequestUrlOk

`func (o *VerifyUserCodeRequest) GetRequestUrlOk() (*string, bool)`

GetRequestUrlOk returns a tuple with the RequestUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestUrl

`func (o *VerifyUserCodeRequest) SetRequestUrl(v string)`

SetRequestUrl sets RequestUrl field to given value.

### HasRequestUrl

`func (o *VerifyUserCodeRequest) HasRequestUrl() bool`

HasRequestUrl returns a boolean if a field has been set.

### GetRequestedAccessTokenAudience

`func (o *VerifyUserCodeRequest) GetRequestedAccessTokenAudience() []string`

GetRequestedAccessTokenAudience returns the RequestedAccessTokenAudience field if non-nil, zero value otherwise.

### GetRequestedAccessTokenAudienceOk

`func (o *VerifyUserCodeRequest) GetRequestedAccessTokenAudienceOk() (*[]string, bool)`

GetRequestedAccessTokenAudienceOk returns a tuple with the RequestedAccessTokenAudience field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestedAccessTokenAudience

`func (o *VerifyUserCodeRequest) SetRequestedAccessTokenAudience(v []string)`

SetRequestedAccessTokenAudience sets RequestedAccessTokenAudience field to given value.

### HasRequestedAccessTokenAudience

`func (o *VerifyUserCodeRequest) HasRequestedAccessTokenAudience() bool`

HasRequestedAccessTokenAudience returns a boolean if a field has been set.

### GetRequestedScope

`func (o *VerifyUserCodeRequest) GetRequestedScope() []string`

GetRequestedScope returns the RequestedScope field if non-nil, zero value otherwise.

### GetRequestedScopeOk

`func (o *VerifyUserCodeRequest) GetRequestedScopeOk() (*[]string, bool)`

GetRequestedScopeOk returns a tuple with the RequestedScope field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestedScope

`func (o *VerifyUserCodeRequest) SetRequestedScope(v []string)`

SetRequestedScope sets RequestedScope field to given value.

### HasRequestedScope

`func (o *VerifyUserCodeRequest) HasRequestedScope() bool`

HasRequestedScope returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


