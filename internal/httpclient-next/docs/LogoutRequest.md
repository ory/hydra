# LogoutRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Challenge** | Pointer to **string** | Challenge is the identifier (\&quot;logout challenge\&quot;) of the logout authentication request. It is used to identify the session. | [optional] 
**Client** | Pointer to [**OAuth2Client**](OAuth2Client.md) |  | [optional] 
**RequestUrl** | Pointer to **string** | RequestURL is the original Logout URL requested. | [optional] 
**RpInitiated** | Pointer to **bool** | RPInitiated is set to true if the request was initiated by a Relying Party (RP), also known as an OAuth 2.0 Client. | [optional] 
**Sid** | Pointer to **string** | SessionID is the login session ID that was requested to log out. | [optional] 
**Subject** | Pointer to **string** | Subject is the user for whom the logout was request. | [optional] 

## Methods

### NewLogoutRequest

`func NewLogoutRequest() *LogoutRequest`

NewLogoutRequest instantiates a new LogoutRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewLogoutRequestWithDefaults

`func NewLogoutRequestWithDefaults() *LogoutRequest`

NewLogoutRequestWithDefaults instantiates a new LogoutRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetChallenge

`func (o *LogoutRequest) GetChallenge() string`

GetChallenge returns the Challenge field if non-nil, zero value otherwise.

### GetChallengeOk

`func (o *LogoutRequest) GetChallengeOk() (*string, bool)`

GetChallengeOk returns a tuple with the Challenge field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChallenge

`func (o *LogoutRequest) SetChallenge(v string)`

SetChallenge sets Challenge field to given value.

### HasChallenge

`func (o *LogoutRequest) HasChallenge() bool`

HasChallenge returns a boolean if a field has been set.

### GetClient

`func (o *LogoutRequest) GetClient() OAuth2Client`

GetClient returns the Client field if non-nil, zero value otherwise.

### GetClientOk

`func (o *LogoutRequest) GetClientOk() (*OAuth2Client, bool)`

GetClientOk returns a tuple with the Client field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClient

`func (o *LogoutRequest) SetClient(v OAuth2Client)`

SetClient sets Client field to given value.

### HasClient

`func (o *LogoutRequest) HasClient() bool`

HasClient returns a boolean if a field has been set.

### GetRequestUrl

`func (o *LogoutRequest) GetRequestUrl() string`

GetRequestUrl returns the RequestUrl field if non-nil, zero value otherwise.

### GetRequestUrlOk

`func (o *LogoutRequest) GetRequestUrlOk() (*string, bool)`

GetRequestUrlOk returns a tuple with the RequestUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestUrl

`func (o *LogoutRequest) SetRequestUrl(v string)`

SetRequestUrl sets RequestUrl field to given value.

### HasRequestUrl

`func (o *LogoutRequest) HasRequestUrl() bool`

HasRequestUrl returns a boolean if a field has been set.

### GetRpInitiated

`func (o *LogoutRequest) GetRpInitiated() bool`

GetRpInitiated returns the RpInitiated field if non-nil, zero value otherwise.

### GetRpInitiatedOk

`func (o *LogoutRequest) GetRpInitiatedOk() (*bool, bool)`

GetRpInitiatedOk returns a tuple with the RpInitiated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRpInitiated

`func (o *LogoutRequest) SetRpInitiated(v bool)`

SetRpInitiated sets RpInitiated field to given value.

### HasRpInitiated

`func (o *LogoutRequest) HasRpInitiated() bool`

HasRpInitiated returns a boolean if a field has been set.

### GetSid

`func (o *LogoutRequest) GetSid() string`

GetSid returns the Sid field if non-nil, zero value otherwise.

### GetSidOk

`func (o *LogoutRequest) GetSidOk() (*string, bool)`

GetSidOk returns a tuple with the Sid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSid

`func (o *LogoutRequest) SetSid(v string)`

SetSid sets Sid field to given value.

### HasSid

`func (o *LogoutRequest) HasSid() bool`

HasSid returns a boolean if a field has been set.

### GetSubject

`func (o *LogoutRequest) GetSubject() string`

GetSubject returns the Subject field if non-nil, zero value otherwise.

### GetSubjectOk

`func (o *LogoutRequest) GetSubjectOk() (*string, bool)`

GetSubjectOk returns a tuple with the Subject field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubject

`func (o *LogoutRequest) SetSubject(v string)`

SetSubject sets Subject field to given value.

### HasSubject

`func (o *LogoutRequest) HasSubject() bool`

HasSubject returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


