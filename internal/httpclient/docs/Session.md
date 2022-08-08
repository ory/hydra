# Session

## Properties

| Name                      | Type                                                | Description | Notes      |
| ------------------------- | --------------------------------------------------- | ----------- | ---------- |
| **AllowedTopLevelClaims** | Pointer to **[]string**                             |             | [optional] |
| **ClientId**              | Pointer to **string**                               |             | [optional] |
| **ConsentChallenge**      | Pointer to **string**                               |             | [optional] |
| **ExcludeNotBeforeClaim** | Pointer to **bool**                                 |             | [optional] |
| **ExpiresAt**             | Pointer to [**map[string]time.Time**](time.Time.md) |             | [optional] |
| **Extra**                 | Pointer to **map[string]interface{}**               |             | [optional] |
| **Headers**               | Pointer to [**Headers**](Headers.md)                |             | [optional] |
| **IdTokenClaims**         | Pointer to [**IDTokenClaims**](IDTokenClaims.md)    |             | [optional] |
| **Kid**                   | Pointer to **string**                               |             | [optional] |
| **Subject**               | Pointer to **string**                               |             | [optional] |
| **Username**              | Pointer to **string**                               |             | [optional] |

## Methods

### NewSession

`func NewSession() *Session`

NewSession instantiates a new Session object This constructor will assign
default values to properties that have it defined, and makes sure properties
required by API are set, but the set of arguments will change when the set of
required properties is changed

### NewSessionWithDefaults

`func NewSessionWithDefaults() *Session`

NewSessionWithDefaults instantiates a new Session object This constructor will
only assign default values to properties that have it defined, but it doesn't
guarantee that properties required by API are set

### GetAllowedTopLevelClaims

`func (o *Session) GetAllowedTopLevelClaims() []string`

GetAllowedTopLevelClaims returns the AllowedTopLevelClaims field if non-nil,
zero value otherwise.

### GetAllowedTopLevelClaimsOk

`func (o *Session) GetAllowedTopLevelClaimsOk() (*[]string, bool)`

GetAllowedTopLevelClaimsOk returns a tuple with the AllowedTopLevelClaims field
if it's non-nil, zero value otherwise and a boolean to check if the value has
been set.

### SetAllowedTopLevelClaims

`func (o *Session) SetAllowedTopLevelClaims(v []string)`

SetAllowedTopLevelClaims sets AllowedTopLevelClaims field to given value.

### HasAllowedTopLevelClaims

`func (o *Session) HasAllowedTopLevelClaims() bool`

HasAllowedTopLevelClaims returns a boolean if a field has been set.

### GetClientId

`func (o *Session) GetClientId() string`

GetClientId returns the ClientId field if non-nil, zero value otherwise.

### GetClientIdOk

`func (o *Session) GetClientIdOk() (*string, bool)`

GetClientIdOk returns a tuple with the ClientId field if it's non-nil, zero
value otherwise and a boolean to check if the value has been set.

### SetClientId

`func (o *Session) SetClientId(v string)`

SetClientId sets ClientId field to given value.

### HasClientId

`func (o *Session) HasClientId() bool`

HasClientId returns a boolean if a field has been set.

### GetConsentChallenge

`func (o *Session) GetConsentChallenge() string`

GetConsentChallenge returns the ConsentChallenge field if non-nil, zero value
otherwise.

### GetConsentChallengeOk

`func (o *Session) GetConsentChallengeOk() (*string, bool)`

GetConsentChallengeOk returns a tuple with the ConsentChallenge field if it's
non-nil, zero value otherwise and a boolean to check if the value has been set.

### SetConsentChallenge

`func (o *Session) SetConsentChallenge(v string)`

SetConsentChallenge sets ConsentChallenge field to given value.

### HasConsentChallenge

`func (o *Session) HasConsentChallenge() bool`

HasConsentChallenge returns a boolean if a field has been set.

### GetExcludeNotBeforeClaim

`func (o *Session) GetExcludeNotBeforeClaim() bool`

GetExcludeNotBeforeClaim returns the ExcludeNotBeforeClaim field if non-nil,
zero value otherwise.

### GetExcludeNotBeforeClaimOk

`func (o *Session) GetExcludeNotBeforeClaimOk() (*bool, bool)`

GetExcludeNotBeforeClaimOk returns a tuple with the ExcludeNotBeforeClaim field
if it's non-nil, zero value otherwise and a boolean to check if the value has
been set.

### SetExcludeNotBeforeClaim

`func (o *Session) SetExcludeNotBeforeClaim(v bool)`

SetExcludeNotBeforeClaim sets ExcludeNotBeforeClaim field to given value.

### HasExcludeNotBeforeClaim

`func (o *Session) HasExcludeNotBeforeClaim() bool`

HasExcludeNotBeforeClaim returns a boolean if a field has been set.

### GetExpiresAt

`func (o *Session) GetExpiresAt() map[string]time.Time`

GetExpiresAt returns the ExpiresAt field if non-nil, zero value otherwise.

### GetExpiresAtOk

`func (o *Session) GetExpiresAtOk() (*map[string]time.Time, bool)`

GetExpiresAtOk returns a tuple with the ExpiresAt field if it's non-nil, zero
value otherwise and a boolean to check if the value has been set.

### SetExpiresAt

`func (o *Session) SetExpiresAt(v map[string]time.Time)`

SetExpiresAt sets ExpiresAt field to given value.

### HasExpiresAt

`func (o *Session) HasExpiresAt() bool`

HasExpiresAt returns a boolean if a field has been set.

### GetExtra

`func (o *Session) GetExtra() map[string]interface{}`

GetExtra returns the Extra field if non-nil, zero value otherwise.

### GetExtraOk

`func (o *Session) GetExtraOk() (*map[string]interface{}, bool)`

GetExtraOk returns a tuple with the Extra field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetExtra

`func (o *Session) SetExtra(v map[string]interface{})`

SetExtra sets Extra field to given value.

### HasExtra

`func (o *Session) HasExtra() bool`

HasExtra returns a boolean if a field has been set.

### GetHeaders

`func (o *Session) GetHeaders() Headers`

GetHeaders returns the Headers field if non-nil, zero value otherwise.

### GetHeadersOk

`func (o *Session) GetHeadersOk() (*Headers, bool)`

GetHeadersOk returns a tuple with the Headers field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetHeaders

`func (o *Session) SetHeaders(v Headers)`

SetHeaders sets Headers field to given value.

### HasHeaders

`func (o *Session) HasHeaders() bool`

HasHeaders returns a boolean if a field has been set.

### GetIdTokenClaims

`func (o *Session) GetIdTokenClaims() IDTokenClaims`

GetIdTokenClaims returns the IdTokenClaims field if non-nil, zero value
otherwise.

### GetIdTokenClaimsOk

`func (o *Session) GetIdTokenClaimsOk() (*IDTokenClaims, bool)`

GetIdTokenClaimsOk returns a tuple with the IdTokenClaims field if it's non-nil,
zero value otherwise and a boolean to check if the value has been set.

### SetIdTokenClaims

`func (o *Session) SetIdTokenClaims(v IDTokenClaims)`

SetIdTokenClaims sets IdTokenClaims field to given value.

### HasIdTokenClaims

`func (o *Session) HasIdTokenClaims() bool`

HasIdTokenClaims returns a boolean if a field has been set.

### GetKid

`func (o *Session) GetKid() string`

GetKid returns the Kid field if non-nil, zero value otherwise.

### GetKidOk

`func (o *Session) GetKidOk() (*string, bool)`

GetKidOk returns a tuple with the Kid field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetKid

`func (o *Session) SetKid(v string)`

SetKid sets Kid field to given value.

### HasKid

`func (o *Session) HasKid() bool`

HasKid returns a boolean if a field has been set.

### GetSubject

`func (o *Session) GetSubject() string`

GetSubject returns the Subject field if non-nil, zero value otherwise.

### GetSubjectOk

`func (o *Session) GetSubjectOk() (*string, bool)`

GetSubjectOk returns a tuple with the Subject field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetSubject

`func (o *Session) SetSubject(v string)`

SetSubject sets Subject field to given value.

### HasSubject

`func (o *Session) HasSubject() bool`

HasSubject returns a boolean if a field has been set.

### GetUsername

`func (o *Session) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *Session) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero
value otherwise and a boolean to check if the value has been set.

### SetUsername

`func (o *Session) SetUsername(v string)`

SetUsername sets Username field to given value.

### HasUsername

`func (o *Session) HasUsername() bool`

HasUsername returns a boolean if a field has been set.

[[Back to Model list]](../README.md#documentation-for-models)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to README]](../README.md)
