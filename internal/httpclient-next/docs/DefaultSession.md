# DefaultSession

## Properties

| Name              | Type                                                | Description | Notes      |
| ----------------- | --------------------------------------------------- | ----------- | ---------- |
| **ExpiresAt**     | Pointer to [**map[string]time.Time**](time.Time.md) |             | [optional] |
| **Headers**       | Pointer to [**Headers**](Headers.md)                |             | [optional] |
| **IdTokenClaims** | Pointer to [**IDTokenClaims**](IDTokenClaims.md)    |             | [optional] |
| **Subject**       | Pointer to **string**                               |             | [optional] |
| **Username**      | Pointer to **string**                               |             | [optional] |

## Methods

### NewDefaultSession

`func NewDefaultSession() *DefaultSession`

NewDefaultSession instantiates a new DefaultSession object This constructor will
assign default values to properties that have it defined, and makes sure
properties required by API are set, but the set of arguments will change when
the set of required properties is changed

### NewDefaultSessionWithDefaults

`func NewDefaultSessionWithDefaults() *DefaultSession`

NewDefaultSessionWithDefaults instantiates a new DefaultSession object This
constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetExpiresAt

`func (o *DefaultSession) GetExpiresAt() map[string]time.Time`

GetExpiresAt returns the ExpiresAt field if non-nil, zero value otherwise.

### GetExpiresAtOk

`func (o *DefaultSession) GetExpiresAtOk() (*map[string]time.Time, bool)`

GetExpiresAtOk returns a tuple with the ExpiresAt field if it's non-nil, zero
value otherwise and a boolean to check if the value has been set.

### SetExpiresAt

`func (o *DefaultSession) SetExpiresAt(v map[string]time.Time)`

SetExpiresAt sets ExpiresAt field to given value.

### HasExpiresAt

`func (o *DefaultSession) HasExpiresAt() bool`

HasExpiresAt returns a boolean if a field has been set.

### GetHeaders

`func (o *DefaultSession) GetHeaders() Headers`

GetHeaders returns the Headers field if non-nil, zero value otherwise.

### GetHeadersOk

`func (o *DefaultSession) GetHeadersOk() (*Headers, bool)`

GetHeadersOk returns a tuple with the Headers field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetHeaders

`func (o *DefaultSession) SetHeaders(v Headers)`

SetHeaders sets Headers field to given value.

### HasHeaders

`func (o *DefaultSession) HasHeaders() bool`

HasHeaders returns a boolean if a field has been set.

### GetIdTokenClaims

`func (o *DefaultSession) GetIdTokenClaims() IDTokenClaims`

GetIdTokenClaims returns the IdTokenClaims field if non-nil, zero value
otherwise.

### GetIdTokenClaimsOk

`func (o *DefaultSession) GetIdTokenClaimsOk() (*IDTokenClaims, bool)`

GetIdTokenClaimsOk returns a tuple with the IdTokenClaims field if it's non-nil,
zero value otherwise and a boolean to check if the value has been set.

### SetIdTokenClaims

`func (o *DefaultSession) SetIdTokenClaims(v IDTokenClaims)`

SetIdTokenClaims sets IdTokenClaims field to given value.

### HasIdTokenClaims

`func (o *DefaultSession) HasIdTokenClaims() bool`

HasIdTokenClaims returns a boolean if a field has been set.

### GetSubject

`func (o *DefaultSession) GetSubject() string`

GetSubject returns the Subject field if non-nil, zero value otherwise.

### GetSubjectOk

`func (o *DefaultSession) GetSubjectOk() (*string, bool)`

GetSubjectOk returns a tuple with the Subject field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetSubject

`func (o *DefaultSession) SetSubject(v string)`

SetSubject sets Subject field to given value.

### HasSubject

`func (o *DefaultSession) HasSubject() bool`

HasSubject returns a boolean if a field has been set.

### GetUsername

`func (o *DefaultSession) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *DefaultSession) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero
value otherwise and a boolean to check if the value has been set.

### SetUsername

`func (o *DefaultSession) SetUsername(v string)`

SetUsername sets Username field to given value.

### HasUsername

`func (o *DefaultSession) HasUsername() bool`

HasUsername returns a boolean if a field has been set.

[[Back to Model list]](../README.md#documentation-for-models)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to README]](../README.md)
