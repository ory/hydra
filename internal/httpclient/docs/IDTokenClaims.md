# IDTokenClaims

## Properties

| Name         | Type                                  | Description | Notes      |
| ------------ | ------------------------------------- | ----------- | ---------- |
| **Acr**      | Pointer to **string**                 |             | [optional] |
| **Amr**      | Pointer to **[]string**               |             | [optional] |
| **AtHash**   | Pointer to **string**                 |             | [optional] |
| **Aud**      | Pointer to **[]string**               |             | [optional] |
| **AuthTime** | Pointer to **time.Time**              |             | [optional] |
| **CHash**    | Pointer to **string**                 |             | [optional] |
| **Exp**      | Pointer to **time.Time**              |             | [optional] |
| **Ext**      | Pointer to **map[string]interface{}** |             | [optional] |
| **Iat**      | Pointer to **time.Time**              |             | [optional] |
| **Iss**      | Pointer to **string**                 |             | [optional] |
| **Jti**      | Pointer to **string**                 |             | [optional] |
| **Nonce**    | Pointer to **string**                 |             | [optional] |
| **Rat**      | Pointer to **time.Time**              |             | [optional] |
| **Sub**      | Pointer to **string**                 |             | [optional] |

## Methods

### NewIDTokenClaims

`func NewIDTokenClaims() *IDTokenClaims`

NewIDTokenClaims instantiates a new IDTokenClaims object This constructor will
assign default values to properties that have it defined, and makes sure
properties required by API are set, but the set of arguments will change when
the set of required properties is changed

### NewIDTokenClaimsWithDefaults

`func NewIDTokenClaimsWithDefaults() *IDTokenClaims`

NewIDTokenClaimsWithDefaults instantiates a new IDTokenClaims object This
constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAcr

`func (o *IDTokenClaims) GetAcr() string`

GetAcr returns the Acr field if non-nil, zero value otherwise.

### GetAcrOk

`func (o *IDTokenClaims) GetAcrOk() (*string, bool)`

GetAcrOk returns a tuple with the Acr field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetAcr

`func (o *IDTokenClaims) SetAcr(v string)`

SetAcr sets Acr field to given value.

### HasAcr

`func (o *IDTokenClaims) HasAcr() bool`

HasAcr returns a boolean if a field has been set.

### GetAmr

`func (o *IDTokenClaims) GetAmr() []string`

GetAmr returns the Amr field if non-nil, zero value otherwise.

### GetAmrOk

`func (o *IDTokenClaims) GetAmrOk() (*[]string, bool)`

GetAmrOk returns a tuple with the Amr field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetAmr

`func (o *IDTokenClaims) SetAmr(v []string)`

SetAmr sets Amr field to given value.

### HasAmr

`func (o *IDTokenClaims) HasAmr() bool`

HasAmr returns a boolean if a field has been set.

### GetAtHash

`func (o *IDTokenClaims) GetAtHash() string`

GetAtHash returns the AtHash field if non-nil, zero value otherwise.

### GetAtHashOk

`func (o *IDTokenClaims) GetAtHashOk() (*string, bool)`

GetAtHashOk returns a tuple with the AtHash field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetAtHash

`func (o *IDTokenClaims) SetAtHash(v string)`

SetAtHash sets AtHash field to given value.

### HasAtHash

`func (o *IDTokenClaims) HasAtHash() bool`

HasAtHash returns a boolean if a field has been set.

### GetAud

`func (o *IDTokenClaims) GetAud() []string`

GetAud returns the Aud field if non-nil, zero value otherwise.

### GetAudOk

`func (o *IDTokenClaims) GetAudOk() (*[]string, bool)`

GetAudOk returns a tuple with the Aud field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetAud

`func (o *IDTokenClaims) SetAud(v []string)`

SetAud sets Aud field to given value.

### HasAud

`func (o *IDTokenClaims) HasAud() bool`

HasAud returns a boolean if a field has been set.

### GetAuthTime

`func (o *IDTokenClaims) GetAuthTime() time.Time`

GetAuthTime returns the AuthTime field if non-nil, zero value otherwise.

### GetAuthTimeOk

`func (o *IDTokenClaims) GetAuthTimeOk() (*time.Time, bool)`

GetAuthTimeOk returns a tuple with the AuthTime field if it's non-nil, zero
value otherwise and a boolean to check if the value has been set.

### SetAuthTime

`func (o *IDTokenClaims) SetAuthTime(v time.Time)`

SetAuthTime sets AuthTime field to given value.

### HasAuthTime

`func (o *IDTokenClaims) HasAuthTime() bool`

HasAuthTime returns a boolean if a field has been set.

### GetCHash

`func (o *IDTokenClaims) GetCHash() string`

GetCHash returns the CHash field if non-nil, zero value otherwise.

### GetCHashOk

`func (o *IDTokenClaims) GetCHashOk() (*string, bool)`

GetCHashOk returns a tuple with the CHash field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetCHash

`func (o *IDTokenClaims) SetCHash(v string)`

SetCHash sets CHash field to given value.

### HasCHash

`func (o *IDTokenClaims) HasCHash() bool`

HasCHash returns a boolean if a field has been set.

### GetExp

`func (o *IDTokenClaims) GetExp() time.Time`

GetExp returns the Exp field if non-nil, zero value otherwise.

### GetExpOk

`func (o *IDTokenClaims) GetExpOk() (*time.Time, bool)`

GetExpOk returns a tuple with the Exp field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetExp

`func (o *IDTokenClaims) SetExp(v time.Time)`

SetExp sets Exp field to given value.

### HasExp

`func (o *IDTokenClaims) HasExp() bool`

HasExp returns a boolean if a field has been set.

### GetExt

`func (o *IDTokenClaims) GetExt() map[string]interface{}`

GetExt returns the Ext field if non-nil, zero value otherwise.

### GetExtOk

`func (o *IDTokenClaims) GetExtOk() (*map[string]interface{}, bool)`

GetExtOk returns a tuple with the Ext field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetExt

`func (o *IDTokenClaims) SetExt(v map[string]interface{})`

SetExt sets Ext field to given value.

### HasExt

`func (o *IDTokenClaims) HasExt() bool`

HasExt returns a boolean if a field has been set.

### GetIat

`func (o *IDTokenClaims) GetIat() time.Time`

GetIat returns the Iat field if non-nil, zero value otherwise.

### GetIatOk

`func (o *IDTokenClaims) GetIatOk() (*time.Time, bool)`

GetIatOk returns a tuple with the Iat field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetIat

`func (o *IDTokenClaims) SetIat(v time.Time)`

SetIat sets Iat field to given value.

### HasIat

`func (o *IDTokenClaims) HasIat() bool`

HasIat returns a boolean if a field has been set.

### GetIss

`func (o *IDTokenClaims) GetIss() string`

GetIss returns the Iss field if non-nil, zero value otherwise.

### GetIssOk

`func (o *IDTokenClaims) GetIssOk() (*string, bool)`

GetIssOk returns a tuple with the Iss field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetIss

`func (o *IDTokenClaims) SetIss(v string)`

SetIss sets Iss field to given value.

### HasIss

`func (o *IDTokenClaims) HasIss() bool`

HasIss returns a boolean if a field has been set.

### GetJti

`func (o *IDTokenClaims) GetJti() string`

GetJti returns the Jti field if non-nil, zero value otherwise.

### GetJtiOk

`func (o *IDTokenClaims) GetJtiOk() (*string, bool)`

GetJtiOk returns a tuple with the Jti field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetJti

`func (o *IDTokenClaims) SetJti(v string)`

SetJti sets Jti field to given value.

### HasJti

`func (o *IDTokenClaims) HasJti() bool`

HasJti returns a boolean if a field has been set.

### GetNonce

`func (o *IDTokenClaims) GetNonce() string`

GetNonce returns the Nonce field if non-nil, zero value otherwise.

### GetNonceOk

`func (o *IDTokenClaims) GetNonceOk() (*string, bool)`

GetNonceOk returns a tuple with the Nonce field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetNonce

`func (o *IDTokenClaims) SetNonce(v string)`

SetNonce sets Nonce field to given value.

### HasNonce

`func (o *IDTokenClaims) HasNonce() bool`

HasNonce returns a boolean if a field has been set.

### GetRat

`func (o *IDTokenClaims) GetRat() time.Time`

GetRat returns the Rat field if non-nil, zero value otherwise.

### GetRatOk

`func (o *IDTokenClaims) GetRatOk() (*time.Time, bool)`

GetRatOk returns a tuple with the Rat field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetRat

`func (o *IDTokenClaims) SetRat(v time.Time)`

SetRat sets Rat field to given value.

### HasRat

`func (o *IDTokenClaims) HasRat() bool`

HasRat returns a boolean if a field has been set.

### GetSub

`func (o *IDTokenClaims) GetSub() string`

GetSub returns the Sub field if non-nil, zero value otherwise.

### GetSubOk

`func (o *IDTokenClaims) GetSubOk() (*string, bool)`

GetSubOk returns a tuple with the Sub field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetSub

`func (o *IDTokenClaims) SetSub(v string)`

SetSub sets Sub field to given value.

### HasSub

`func (o *IDTokenClaims) HasSub() bool`

HasSub returns a boolean if a field has been set.

[[Back to Model list]](../README.md#documentation-for-models)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to README]](../README.md)
