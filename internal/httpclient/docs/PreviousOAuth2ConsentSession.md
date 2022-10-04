# PreviousOAuth2ConsentSession

## Properties

| Name                         | Type                                                                                     | Description                                                                                                                                                              | Notes      |
| ---------------------------- | ---------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---------- |
| **ConsentRequest**           | Pointer to [**OAuth2ConsentRequest**](OAuth2ConsentRequest.md)                           |                                                                                                                                                                          | [optional] |
| **GrantAccessTokenAudience** | Pointer to **[]string**                                                                  |                                                                                                                                                                          | [optional] |
| **GrantScope**               | Pointer to **[]string**                                                                  |                                                                                                                                                                          | [optional] |
| **HandledAt**                | Pointer to **time.Time**                                                                 |                                                                                                                                                                          | [optional] |
| **Remember**                 | Pointer to **bool**                                                                      | Remember, if set to true, tells ORY Hydra to remember this consent authorization and reuse it if the same client asks the same user for the same, or a subset of, scope. | [optional] |
| **RememberFor**              | Pointer to **int64**                                                                     | RememberFor sets how long the consent authorization should be remembered for in seconds. If set to &#x60;0&#x60;, the authorization will be remembered indefinitely.     | [optional] |
| **Session**                  | Pointer to [**AcceptOAuth2ConsentRequestSession**](AcceptOAuth2ConsentRequestSession.md) |                                                                                                                                                                          | [optional] |

## Methods

### NewPreviousOAuth2ConsentSession

`func NewPreviousOAuth2ConsentSession() *PreviousOAuth2ConsentSession`

NewPreviousOAuth2ConsentSession instantiates a new PreviousOAuth2ConsentSession
object This constructor will assign default values to properties that have it
defined, and makes sure properties required by API are set, but the set of
arguments will change when the set of required properties is changed

### NewPreviousOAuth2ConsentSessionWithDefaults

`func NewPreviousOAuth2ConsentSessionWithDefaults() *PreviousOAuth2ConsentSession`

NewPreviousOAuth2ConsentSessionWithDefaults instantiates a new
PreviousOAuth2ConsentSession object This constructor will only assign default
values to properties that have it defined, but it doesn't guarantee that
properties required by API are set

### GetConsentRequest

`func (o *PreviousOAuth2ConsentSession) GetConsentRequest() OAuth2ConsentRequest`

GetConsentRequest returns the ConsentRequest field if non-nil, zero value
otherwise.

### GetConsentRequestOk

`func (o *PreviousOAuth2ConsentSession) GetConsentRequestOk() (*OAuth2ConsentRequest, bool)`

GetConsentRequestOk returns a tuple with the ConsentRequest field if it's
non-nil, zero value otherwise and a boolean to check if the value has been set.

### SetConsentRequest

`func (o *PreviousOAuth2ConsentSession) SetConsentRequest(v OAuth2ConsentRequest)`

SetConsentRequest sets ConsentRequest field to given value.

### HasConsentRequest

`func (o *PreviousOAuth2ConsentSession) HasConsentRequest() bool`

HasConsentRequest returns a boolean if a field has been set.

### GetGrantAccessTokenAudience

`func (o *PreviousOAuth2ConsentSession) GetGrantAccessTokenAudience() []string`

GetGrantAccessTokenAudience returns the GrantAccessTokenAudience field if
non-nil, zero value otherwise.

### GetGrantAccessTokenAudienceOk

`func (o *PreviousOAuth2ConsentSession) GetGrantAccessTokenAudienceOk() (*[]string, bool)`

GetGrantAccessTokenAudienceOk returns a tuple with the GrantAccessTokenAudience
field if it's non-nil, zero value otherwise and a boolean to check if the value
has been set.

### SetGrantAccessTokenAudience

`func (o *PreviousOAuth2ConsentSession) SetGrantAccessTokenAudience(v []string)`

SetGrantAccessTokenAudience sets GrantAccessTokenAudience field to given value.

### HasGrantAccessTokenAudience

`func (o *PreviousOAuth2ConsentSession) HasGrantAccessTokenAudience() bool`

HasGrantAccessTokenAudience returns a boolean if a field has been set.

### GetGrantScope

`func (o *PreviousOAuth2ConsentSession) GetGrantScope() []string`

GetGrantScope returns the GrantScope field if non-nil, zero value otherwise.

### GetGrantScopeOk

`func (o *PreviousOAuth2ConsentSession) GetGrantScopeOk() (*[]string, bool)`

GetGrantScopeOk returns a tuple with the GrantScope field if it's non-nil, zero
value otherwise and a boolean to check if the value has been set.

### SetGrantScope

`func (o *PreviousOAuth2ConsentSession) SetGrantScope(v []string)`

SetGrantScope sets GrantScope field to given value.

### HasGrantScope

`func (o *PreviousOAuth2ConsentSession) HasGrantScope() bool`

HasGrantScope returns a boolean if a field has been set.

### GetHandledAt

`func (o *PreviousOAuth2ConsentSession) GetHandledAt() time.Time`

GetHandledAt returns the HandledAt field if non-nil, zero value otherwise.

### GetHandledAtOk

`func (o *PreviousOAuth2ConsentSession) GetHandledAtOk() (*time.Time, bool)`

GetHandledAtOk returns a tuple with the HandledAt field if it's non-nil, zero
value otherwise and a boolean to check if the value has been set.

### SetHandledAt

`func (o *PreviousOAuth2ConsentSession) SetHandledAt(v time.Time)`

SetHandledAt sets HandledAt field to given value.

### HasHandledAt

`func (o *PreviousOAuth2ConsentSession) HasHandledAt() bool`

HasHandledAt returns a boolean if a field has been set.

### GetRemember

`func (o *PreviousOAuth2ConsentSession) GetRemember() bool`

GetRemember returns the Remember field if non-nil, zero value otherwise.

### GetRememberOk

`func (o *PreviousOAuth2ConsentSession) GetRememberOk() (*bool, bool)`

GetRememberOk returns a tuple with the Remember field if it's non-nil, zero
value otherwise and a boolean to check if the value has been set.

### SetRemember

`func (o *PreviousOAuth2ConsentSession) SetRemember(v bool)`

SetRemember sets Remember field to given value.

### HasRemember

`func (o *PreviousOAuth2ConsentSession) HasRemember() bool`

HasRemember returns a boolean if a field has been set.

### GetRememberFor

`func (o *PreviousOAuth2ConsentSession) GetRememberFor() int64`

GetRememberFor returns the RememberFor field if non-nil, zero value otherwise.

### GetRememberForOk

`func (o *PreviousOAuth2ConsentSession) GetRememberForOk() (*int64, bool)`

GetRememberForOk returns a tuple with the RememberFor field if it's non-nil,
zero value otherwise and a boolean to check if the value has been set.

### SetRememberFor

`func (o *PreviousOAuth2ConsentSession) SetRememberFor(v int64)`

SetRememberFor sets RememberFor field to given value.

### HasRememberFor

`func (o *PreviousOAuth2ConsentSession) HasRememberFor() bool`

HasRememberFor returns a boolean if a field has been set.

### GetSession

`func (o *PreviousOAuth2ConsentSession) GetSession() AcceptOAuth2ConsentRequestSession`

GetSession returns the Session field if non-nil, zero value otherwise.

### GetSessionOk

`func (o *PreviousOAuth2ConsentSession) GetSessionOk() (*AcceptOAuth2ConsentRequestSession, bool)`

GetSessionOk returns a tuple with the Session field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetSession

`func (o *PreviousOAuth2ConsentSession) SetSession(v AcceptOAuth2ConsentRequestSession)`

SetSession sets Session field to given value.

### HasSession

`func (o *PreviousOAuth2ConsentSession) HasSession() bool`

HasSession returns a boolean if a field has been set.

[[Back to Model list]](../README.md#documentation-for-models)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to README]](../README.md)
