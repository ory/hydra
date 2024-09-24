# DeviceAuthorization

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DeviceCode** | Pointer to **string** | The device verification code. | [optional] 
**ExpiresIn** | Pointer to **int64** | The lifetime in seconds of the \&quot;device_code\&quot; and \&quot;user_code\&quot;. | [optional] 
**Interval** | Pointer to **int64** | The minimum amount of time in seconds that the client SHOULD wait between polling requests to the token endpoint.  If no value is provided, clients MUST use 5 as the default. | [optional] 
**UserCode** | Pointer to **string** | The end-user verification code. | [optional] 
**VerificationUri** | Pointer to **string** | The end-user verification URI on the authorization server.  The URI should be short and easy to remember as end users will be asked to manually type it into their user agent. | [optional] 
**VerificationUriComplete** | Pointer to **string** | A verification URI that includes the \&quot;user_code\&quot; (or other information with the same function as the \&quot;user_code\&quot;), which is designed for non-textual transmission. | [optional] 

## Methods

### NewDeviceAuthorization

`func NewDeviceAuthorization() *DeviceAuthorization`

NewDeviceAuthorization instantiates a new DeviceAuthorization object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDeviceAuthorizationWithDefaults

`func NewDeviceAuthorizationWithDefaults() *DeviceAuthorization`

NewDeviceAuthorizationWithDefaults instantiates a new DeviceAuthorization object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDeviceCode

`func (o *DeviceAuthorization) GetDeviceCode() string`

GetDeviceCode returns the DeviceCode field if non-nil, zero value otherwise.

### GetDeviceCodeOk

`func (o *DeviceAuthorization) GetDeviceCodeOk() (*string, bool)`

GetDeviceCodeOk returns a tuple with the DeviceCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceCode

`func (o *DeviceAuthorization) SetDeviceCode(v string)`

SetDeviceCode sets DeviceCode field to given value.

### HasDeviceCode

`func (o *DeviceAuthorization) HasDeviceCode() bool`

HasDeviceCode returns a boolean if a field has been set.

### GetExpiresIn

`func (o *DeviceAuthorization) GetExpiresIn() int64`

GetExpiresIn returns the ExpiresIn field if non-nil, zero value otherwise.

### GetExpiresInOk

`func (o *DeviceAuthorization) GetExpiresInOk() (*int64, bool)`

GetExpiresInOk returns a tuple with the ExpiresIn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiresIn

`func (o *DeviceAuthorization) SetExpiresIn(v int64)`

SetExpiresIn sets ExpiresIn field to given value.

### HasExpiresIn

`func (o *DeviceAuthorization) HasExpiresIn() bool`

HasExpiresIn returns a boolean if a field has been set.

### GetInterval

`func (o *DeviceAuthorization) GetInterval() int64`

GetInterval returns the Interval field if non-nil, zero value otherwise.

### GetIntervalOk

`func (o *DeviceAuthorization) GetIntervalOk() (*int64, bool)`

GetIntervalOk returns a tuple with the Interval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInterval

`func (o *DeviceAuthorization) SetInterval(v int64)`

SetInterval sets Interval field to given value.

### HasInterval

`func (o *DeviceAuthorization) HasInterval() bool`

HasInterval returns a boolean if a field has been set.

### GetUserCode

`func (o *DeviceAuthorization) GetUserCode() string`

GetUserCode returns the UserCode field if non-nil, zero value otherwise.

### GetUserCodeOk

`func (o *DeviceAuthorization) GetUserCodeOk() (*string, bool)`

GetUserCodeOk returns a tuple with the UserCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserCode

`func (o *DeviceAuthorization) SetUserCode(v string)`

SetUserCode sets UserCode field to given value.

### HasUserCode

`func (o *DeviceAuthorization) HasUserCode() bool`

HasUserCode returns a boolean if a field has been set.

### GetVerificationUri

`func (o *DeviceAuthorization) GetVerificationUri() string`

GetVerificationUri returns the VerificationUri field if non-nil, zero value otherwise.

### GetVerificationUriOk

`func (o *DeviceAuthorization) GetVerificationUriOk() (*string, bool)`

GetVerificationUriOk returns a tuple with the VerificationUri field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVerificationUri

`func (o *DeviceAuthorization) SetVerificationUri(v string)`

SetVerificationUri sets VerificationUri field to given value.

### HasVerificationUri

`func (o *DeviceAuthorization) HasVerificationUri() bool`

HasVerificationUri returns a boolean if a field has been set.

### GetVerificationUriComplete

`func (o *DeviceAuthorization) GetVerificationUriComplete() string`

GetVerificationUriComplete returns the VerificationUriComplete field if non-nil, zero value otherwise.

### GetVerificationUriCompleteOk

`func (o *DeviceAuthorization) GetVerificationUriCompleteOk() (*string, bool)`

GetVerificationUriCompleteOk returns a tuple with the VerificationUriComplete field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVerificationUriComplete

`func (o *DeviceAuthorization) SetVerificationUriComplete(v string)`

SetVerificationUriComplete sets VerificationUriComplete field to given value.

### HasVerificationUriComplete

`func (o *DeviceAuthorization) HasVerificationUriComplete() bool`

HasVerificationUriComplete returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


