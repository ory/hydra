# DeviceGrantRequest

## Properties

| Name                             | Type                                | Description                                                                                                            | Notes      |
| -------------------------------- | ----------------------------------- | ---------------------------------------------------------------------------------------------------------------------- | ---------- |
| **Challenge**                    | **string**                          | ID is the identifier (\&quot;device challenge\&quot;) of the device grant request. It is used to identify the session. |
| **Client**                       | [**OAuth2Client**](OAuth2Client.md) |                                                                                                                        |
| **HandledAt**                    | Pointer to **time.Time**            |                                                                                                                        | [optional] |
| **RequestUrl**                   | Pointer to **string**               | RequestURL is the original Device Grant URL requested.                                                                 | [optional] |
| **RequestedAccessTokenAudience** | **[]string**                        |                                                                                                                        |
| **RequestedScope**               | **[]string**                        |                                                                                                                        |

## Methods

### NewDeviceGrantRequest

`func NewDeviceGrantRequest(challenge string, client OAuth2Client, requestedAccessTokenAudience []string, requestedScope []string, ) *DeviceGrantRequest`

NewDeviceGrantRequest instantiates a new DeviceGrantRequest object This
constructor will assign default values to properties that have it defined, and
makes sure properties required by API are set, but the set of arguments will
change when the set of required properties is changed

### NewDeviceGrantRequestWithDefaults

`func NewDeviceGrantRequestWithDefaults() *DeviceGrantRequest`

NewDeviceGrantRequestWithDefaults instantiates a new DeviceGrantRequest object
This constructor will only assign default values to properties that have it
defined, but it doesn't guarantee that properties required by API are set

### GetChallenge

`func (o *DeviceGrantRequest) GetChallenge() string`

GetChallenge returns the Challenge field if non-nil, zero value otherwise.

### GetChallengeOk

`func (o *DeviceGrantRequest) GetChallengeOk() (*string, bool)`

GetChallengeOk returns a tuple with the Challenge field if it's non-nil, zero
value otherwise and a boolean to check if the value has been set.

### SetChallenge

`func (o *DeviceGrantRequest) SetChallenge(v string)`

SetChallenge sets Challenge field to given value.

### GetClient

`func (o *DeviceGrantRequest) GetClient() OAuth2Client`

GetClient returns the Client field if non-nil, zero value otherwise.

### GetClientOk

`func (o *DeviceGrantRequest) GetClientOk() (*OAuth2Client, bool)`

GetClientOk returns a tuple with the Client field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetClient

`func (o *DeviceGrantRequest) SetClient(v OAuth2Client)`

SetClient sets Client field to given value.

### GetHandledAt

`func (o *DeviceGrantRequest) GetHandledAt() time.Time`

GetHandledAt returns the HandledAt field if non-nil, zero value otherwise.

### GetHandledAtOk

`func (o *DeviceGrantRequest) GetHandledAtOk() (*time.Time, bool)`

GetHandledAtOk returns a tuple with the HandledAt field if it's non-nil, zero
value otherwise and a boolean to check if the value has been set.

### SetHandledAt

`func (o *DeviceGrantRequest) SetHandledAt(v time.Time)`

SetHandledAt sets HandledAt field to given value.

### HasHandledAt

`func (o *DeviceGrantRequest) HasHandledAt() bool`

HasHandledAt returns a boolean if a field has been set.

### GetRequestUrl

`func (o *DeviceGrantRequest) GetRequestUrl() string`

GetRequestUrl returns the RequestUrl field if non-nil, zero value otherwise.

### GetRequestUrlOk

`func (o *DeviceGrantRequest) GetRequestUrlOk() (*string, bool)`

GetRequestUrlOk returns a tuple with the RequestUrl field if it's non-nil, zero
value otherwise and a boolean to check if the value has been set.

### SetRequestUrl

`func (o *DeviceGrantRequest) SetRequestUrl(v string)`

SetRequestUrl sets RequestUrl field to given value.

### HasRequestUrl

`func (o *DeviceGrantRequest) HasRequestUrl() bool`

HasRequestUrl returns a boolean if a field has been set.

### GetRequestedAccessTokenAudience

`func (o *DeviceGrantRequest) GetRequestedAccessTokenAudience() []string`

GetRequestedAccessTokenAudience returns the RequestedAccessTokenAudience field
if non-nil, zero value otherwise.

### GetRequestedAccessTokenAudienceOk

`func (o *DeviceGrantRequest) GetRequestedAccessTokenAudienceOk() (*[]string, bool)`

GetRequestedAccessTokenAudienceOk returns a tuple with the
RequestedAccessTokenAudience field if it's non-nil, zero value otherwise and a
boolean to check if the value has been set.

### SetRequestedAccessTokenAudience

`func (o *DeviceGrantRequest) SetRequestedAccessTokenAudience(v []string)`

SetRequestedAccessTokenAudience sets RequestedAccessTokenAudience field to given
value.

### GetRequestedScope

`func (o *DeviceGrantRequest) GetRequestedScope() []string`

GetRequestedScope returns the RequestedScope field if non-nil, zero value
otherwise.

### GetRequestedScopeOk

`func (o *DeviceGrantRequest) GetRequestedScopeOk() (*[]string, bool)`

GetRequestedScopeOk returns a tuple with the RequestedScope field if it's
non-nil, zero value otherwise and a boolean to check if the value has been set.

### SetRequestedScope

`func (o *DeviceGrantRequest) SetRequestedScope(v []string)`

SetRequestedScope sets RequestedScope field to given value.

[[Back to Model list]](../README.md#documentation-for-models)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to README]](../README.md)
