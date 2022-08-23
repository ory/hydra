# OAuth2ConsentSession

## Properties

| Name                      | Type                                                                             | Description | Notes      |
| ------------------------- | -------------------------------------------------------------------------------- | ----------- | ---------- |
| **AllowedTopLevelClaims** | Pointer to **[]string**                                                          |             | [optional] |
| **ClientId**              | Pointer to **string**                                                            |             | [optional] |
| **ConsentChallenge**      | Pointer to **string**                                                            |             | [optional] |
| **ExcludeNotBeforeClaim** | Pointer to **bool**                                                              |             | [optional] |
| **ExpiresAt**             | Pointer to [**OAuth2ConsentSessionExpiresAt**](OAuth2ConsentSessionExpiresAt.md) |             | [optional] |
| **Extra**                 | Pointer to **map[string]interface{}**                                            |             | [optional] |
| **Headers**               | Pointer to [**Headers**](Headers.md)                                             |             | [optional] |
| **IdTokenClaims**         | Pointer to [**IDTokenClaims**](IDTokenClaims.md)                                 |             | [optional] |
| **Kid**                   | Pointer to **string**                                                            |             | [optional] |
| **Subject**               | Pointer to **string**                                                            |             | [optional] |
| **Username**              | Pointer to **string**                                                            |             | [optional] |

## Methods

### NewOAuth2ConsentSession

`func NewOAuth2ConsentSession() *OAuth2ConsentSession`

NewOAuth2ConsentSession instantiates a new OAuth2ConsentSession object This
constructor will assign default values to properties that have it defined, and
makes sure properties required by API are set, but the set of arguments will
change when the set of required properties is changed

### NewOAuth2ConsentSessionWithDefaults

`func NewOAuth2ConsentSessionWithDefaults() *OAuth2ConsentSession`

NewOAuth2ConsentSessionWithDefaults instantiates a new OAuth2ConsentSession
object This constructor will only assign default values to properties that have
it defined, but it doesn't guarantee that properties required by API are set

### GetAllowedTopLevelClaims

`func (o *OAuth2ConsentSession) GetAllowedTopLevelClaims() []string`

GetAllowedTopLevelClaims returns the AllowedTopLevelClaims field if non-nil,
zero value otherwise.

### GetAllowedTopLevelClaimsOk

`func (o *OAuth2ConsentSession) GetAllowedTopLevelClaimsOk() (*[]string, bool)`

GetAllowedTopLevelClaimsOk returns a tuple with the AllowedTopLevelClaims field
if it's non-nil, zero value otherwise and a boolean to check if the value has
been set.

### SetAllowedTopLevelClaims

`func (o *OAuth2ConsentSession) SetAllowedTopLevelClaims(v []string)`

SetAllowedTopLevelClaims sets AllowedTopLevelClaims field to given value.

### HasAllowedTopLevelClaims

`func (o *OAuth2ConsentSession) HasAllowedTopLevelClaims() bool`

HasAllowedTopLevelClaims returns a boolean if a field has been set.

### GetClientId

`func (o *OAuth2ConsentSession) GetClientId() string`

GetClientId returns the ClientId field if non-nil, zero value otherwise.

### GetClientIdOk

`func (o *OAuth2ConsentSession) GetClientIdOk() (*string, bool)`

GetClientIdOk returns a tuple with the ClientId field if it's non-nil, zero
value otherwise and a boolean to check if the value has been set.

### SetClientId

`func (o *OAuth2ConsentSession) SetClientId(v string)`

SetClientId sets ClientId field to given value.

### HasClientId

`func (o *OAuth2ConsentSession) HasClientId() bool`

HasClientId returns a boolean if a field has been set.

### GetConsentChallenge

`func (o *OAuth2ConsentSession) GetConsentChallenge() string`

GetConsentChallenge returns the ConsentChallenge field if non-nil, zero value
otherwise.

### GetConsentChallengeOk

`func (o *OAuth2ConsentSession) GetConsentChallengeOk() (*string, bool)`

GetConsentChallengeOk returns a tuple with the ConsentChallenge field if it's
non-nil, zero value otherwise and a boolean to check if the value has been set.

### SetConsentChallenge

`func (o *OAuth2ConsentSession) SetConsentChallenge(v string)`

SetConsentChallenge sets ConsentChallenge field to given value.

### HasConsentChallenge

`func (o *OAuth2ConsentSession) HasConsentChallenge() bool`

HasConsentChallenge returns a boolean if a field has been set.

### GetExcludeNotBeforeClaim

`func (o *OAuth2ConsentSession) GetExcludeNotBeforeClaim() bool`

GetExcludeNotBeforeClaim returns the ExcludeNotBeforeClaim field if non-nil,
zero value otherwise.

### GetExcludeNotBeforeClaimOk

`func (o *OAuth2ConsentSession) GetExcludeNotBeforeClaimOk() (*bool, bool)`

GetExcludeNotBeforeClaimOk returns a tuple with the ExcludeNotBeforeClaim field
if it's non-nil, zero value otherwise and a boolean to check if the value has
been set.

### SetExcludeNotBeforeClaim

`func (o *OAuth2ConsentSession) SetExcludeNotBeforeClaim(v bool)`

SetExcludeNotBeforeClaim sets ExcludeNotBeforeClaim field to given value.

### HasExcludeNotBeforeClaim

`func (o *OAuth2ConsentSession) HasExcludeNotBeforeClaim() bool`

HasExcludeNotBeforeClaim returns a boolean if a field has been set.

### GetExpiresAt

`func (o *OAuth2ConsentSession) GetExpiresAt() OAuth2ConsentSessionExpiresAt`

GetExpiresAt returns the ExpiresAt field if non-nil, zero value otherwise.

### GetExpiresAtOk

`func (o *OAuth2ConsentSession) GetExpiresAtOk() (*OAuth2ConsentSessionExpiresAt, bool)`

GetExpiresAtOk returns a tuple with the ExpiresAt field if it's non-nil, zero
value otherwise and a boolean to check if the value has been set.

### SetExpiresAt

`func (o *OAuth2ConsentSession) SetExpiresAt(v OAuth2ConsentSessionExpiresAt)`

SetExpiresAt sets ExpiresAt field to given value.

### HasExpiresAt

`func (o *OAuth2ConsentSession) HasExpiresAt() bool`

HasExpiresAt returns a boolean if a field has been set.

### GetExtra

`func (o *OAuth2ConsentSession) GetExtra() map[string]interface{}`

GetExtra returns the Extra field if non-nil, zero value otherwise.

### GetExtraOk

`func (o *OAuth2ConsentSession) GetExtraOk() (*map[string]interface{}, bool)`

GetExtraOk returns a tuple with the Extra field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetExtra

`func (o *OAuth2ConsentSession) SetExtra(v map[string]interface{})`

SetExtra sets Extra field to given value.

### HasExtra

`func (o *OAuth2ConsentSession) HasExtra() bool`

HasExtra returns a boolean if a field has been set.

### GetHeaders

`func (o *OAuth2ConsentSession) GetHeaders() Headers`

GetHeaders returns the Headers field if non-nil, zero value otherwise.

### GetHeadersOk

`func (o *OAuth2ConsentSession) GetHeadersOk() (*Headers, bool)`

GetHeadersOk returns a tuple with the Headers field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetHeaders

`func (o *OAuth2ConsentSession) SetHeaders(v Headers)`

SetHeaders sets Headers field to given value.

### HasHeaders

`func (o *OAuth2ConsentSession) HasHeaders() bool`

HasHeaders returns a boolean if a field has been set.

### GetIdTokenClaims

`func (o *OAuth2ConsentSession) GetIdTokenClaims() IDTokenClaims`

GetIdTokenClaims returns the IdTokenClaims field if non-nil, zero value
otherwise.

### GetIdTokenClaimsOk

`func (o *OAuth2ConsentSession) GetIdTokenClaimsOk() (*IDTokenClaims, bool)`

GetIdTokenClaimsOk returns a tuple with the IdTokenClaims field if it's non-nil,
zero value otherwise and a boolean to check if the value has been set.

### SetIdTokenClaims

`func (o *OAuth2ConsentSession) SetIdTokenClaims(v IDTokenClaims)`

SetIdTokenClaims sets IdTokenClaims field to given value.

### HasIdTokenClaims

`func (o *OAuth2ConsentSession) HasIdTokenClaims() bool`

HasIdTokenClaims returns a boolean if a field has been set.

### GetKid

`func (o *OAuth2ConsentSession) GetKid() string`

GetKid returns the Kid field if non-nil, zero value otherwise.

### GetKidOk

`func (o *OAuth2ConsentSession) GetKidOk() (*string, bool)`

GetKidOk returns a tuple with the Kid field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetKid

`func (o *OAuth2ConsentSession) SetKid(v string)`

SetKid sets Kid field to given value.

### HasKid

`func (o *OAuth2ConsentSession) HasKid() bool`

HasKid returns a boolean if a field has been set.

### GetSubject

`func (o *OAuth2ConsentSession) GetSubject() string`

GetSubject returns the Subject field if non-nil, zero value otherwise.

### GetSubjectOk

`func (o *OAuth2ConsentSession) GetSubjectOk() (*string, bool)`

GetSubjectOk returns a tuple with the Subject field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetSubject

`func (o *OAuth2ConsentSession) SetSubject(v string)`

SetSubject sets Subject field to given value.

### HasSubject

`func (o *OAuth2ConsentSession) HasSubject() bool`

HasSubject returns a boolean if a field has been set.

### GetUsername

`func (o *OAuth2ConsentSession) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *OAuth2ConsentSession) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero
value otherwise and a boolean to check if the value has been set.

### SetUsername

`func (o *OAuth2ConsentSession) SetUsername(v string)`

SetUsername sets Username field to given value.

### HasUsername

`func (o *OAuth2ConsentSession) HasUsername() bool`

HasUsername returns a boolean if a field has been set.

[[Back to Model list]](../README.md#documentation-for-models)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to README]](../README.md)
