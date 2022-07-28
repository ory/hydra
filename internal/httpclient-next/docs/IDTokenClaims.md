# IDTokenClaims

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AccessTokenHash** | Pointer to **string** |  | [optional] 
**Audience** | Pointer to **[]string** |  | [optional] 
**AuthTime** | Pointer to **time.Time** |  | [optional] 
**AuthenticationContextClassReference** | Pointer to **string** |  | [optional] 
**AuthenticationMethodsReferences** | Pointer to **[]string** |  | [optional] 
**CodeHash** | Pointer to **string** |  | [optional] 
**ExpiresAt** | Pointer to **time.Time** |  | [optional] 
**Extra** | Pointer to **map[string]map[string]interface{}** |  | [optional] 
**IssuedAt** | Pointer to **time.Time** |  | [optional] 
**Issuer** | Pointer to **string** |  | [optional] 
**JTI** | Pointer to **string** |  | [optional] 
**Nonce** | Pointer to **string** |  | [optional] 
**RequestedAt** | Pointer to **time.Time** |  | [optional] 
**Subject** | Pointer to **string** |  | [optional] 

## Methods

### NewIDTokenClaims

`func NewIDTokenClaims() *IDTokenClaims`

NewIDTokenClaims instantiates a new IDTokenClaims object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewIDTokenClaimsWithDefaults

`func NewIDTokenClaimsWithDefaults() *IDTokenClaims`

NewIDTokenClaimsWithDefaults instantiates a new IDTokenClaims object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAccessTokenHash

`func (o *IDTokenClaims) GetAccessTokenHash() string`

GetAccessTokenHash returns the AccessTokenHash field if non-nil, zero value otherwise.

### GetAccessTokenHashOk

`func (o *IDTokenClaims) GetAccessTokenHashOk() (*string, bool)`

GetAccessTokenHashOk returns a tuple with the AccessTokenHash field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccessTokenHash

`func (o *IDTokenClaims) SetAccessTokenHash(v string)`

SetAccessTokenHash sets AccessTokenHash field to given value.

### HasAccessTokenHash

`func (o *IDTokenClaims) HasAccessTokenHash() bool`

HasAccessTokenHash returns a boolean if a field has been set.

### GetAudience

`func (o *IDTokenClaims) GetAudience() []string`

GetAudience returns the Audience field if non-nil, zero value otherwise.

### GetAudienceOk

`func (o *IDTokenClaims) GetAudienceOk() (*[]string, bool)`

GetAudienceOk returns a tuple with the Audience field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAudience

`func (o *IDTokenClaims) SetAudience(v []string)`

SetAudience sets Audience field to given value.

### HasAudience

`func (o *IDTokenClaims) HasAudience() bool`

HasAudience returns a boolean if a field has been set.

### GetAuthTime

`func (o *IDTokenClaims) GetAuthTime() time.Time`

GetAuthTime returns the AuthTime field if non-nil, zero value otherwise.

### GetAuthTimeOk

`func (o *IDTokenClaims) GetAuthTimeOk() (*time.Time, bool)`

GetAuthTimeOk returns a tuple with the AuthTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthTime

`func (o *IDTokenClaims) SetAuthTime(v time.Time)`

SetAuthTime sets AuthTime field to given value.

### HasAuthTime

`func (o *IDTokenClaims) HasAuthTime() bool`

HasAuthTime returns a boolean if a field has been set.

### GetAuthenticationContextClassReference

`func (o *IDTokenClaims) GetAuthenticationContextClassReference() string`

GetAuthenticationContextClassReference returns the AuthenticationContextClassReference field if non-nil, zero value otherwise.

### GetAuthenticationContextClassReferenceOk

`func (o *IDTokenClaims) GetAuthenticationContextClassReferenceOk() (*string, bool)`

GetAuthenticationContextClassReferenceOk returns a tuple with the AuthenticationContextClassReference field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthenticationContextClassReference

`func (o *IDTokenClaims) SetAuthenticationContextClassReference(v string)`

SetAuthenticationContextClassReference sets AuthenticationContextClassReference field to given value.

### HasAuthenticationContextClassReference

`func (o *IDTokenClaims) HasAuthenticationContextClassReference() bool`

HasAuthenticationContextClassReference returns a boolean if a field has been set.

### GetAuthenticationMethodsReferences

`func (o *IDTokenClaims) GetAuthenticationMethodsReferences() []string`

GetAuthenticationMethodsReferences returns the AuthenticationMethodsReferences field if non-nil, zero value otherwise.

### GetAuthenticationMethodsReferencesOk

`func (o *IDTokenClaims) GetAuthenticationMethodsReferencesOk() (*[]string, bool)`

GetAuthenticationMethodsReferencesOk returns a tuple with the AuthenticationMethodsReferences field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthenticationMethodsReferences

`func (o *IDTokenClaims) SetAuthenticationMethodsReferences(v []string)`

SetAuthenticationMethodsReferences sets AuthenticationMethodsReferences field to given value.

### HasAuthenticationMethodsReferences

`func (o *IDTokenClaims) HasAuthenticationMethodsReferences() bool`

HasAuthenticationMethodsReferences returns a boolean if a field has been set.

### GetCodeHash

`func (o *IDTokenClaims) GetCodeHash() string`

GetCodeHash returns the CodeHash field if non-nil, zero value otherwise.

### GetCodeHashOk

`func (o *IDTokenClaims) GetCodeHashOk() (*string, bool)`

GetCodeHashOk returns a tuple with the CodeHash field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCodeHash

`func (o *IDTokenClaims) SetCodeHash(v string)`

SetCodeHash sets CodeHash field to given value.

### HasCodeHash

`func (o *IDTokenClaims) HasCodeHash() bool`

HasCodeHash returns a boolean if a field has been set.

### GetExpiresAt

`func (o *IDTokenClaims) GetExpiresAt() time.Time`

GetExpiresAt returns the ExpiresAt field if non-nil, zero value otherwise.

### GetExpiresAtOk

`func (o *IDTokenClaims) GetExpiresAtOk() (*time.Time, bool)`

GetExpiresAtOk returns a tuple with the ExpiresAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiresAt

`func (o *IDTokenClaims) SetExpiresAt(v time.Time)`

SetExpiresAt sets ExpiresAt field to given value.

### HasExpiresAt

`func (o *IDTokenClaims) HasExpiresAt() bool`

HasExpiresAt returns a boolean if a field has been set.

### GetExtra

`func (o *IDTokenClaims) GetExtra() map[string]map[string]interface{}`

GetExtra returns the Extra field if non-nil, zero value otherwise.

### GetExtraOk

`func (o *IDTokenClaims) GetExtraOk() (*map[string]map[string]interface{}, bool)`

GetExtraOk returns a tuple with the Extra field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExtra

`func (o *IDTokenClaims) SetExtra(v map[string]map[string]interface{})`

SetExtra sets Extra field to given value.

### HasExtra

`func (o *IDTokenClaims) HasExtra() bool`

HasExtra returns a boolean if a field has been set.

### GetIssuedAt

`func (o *IDTokenClaims) GetIssuedAt() time.Time`

GetIssuedAt returns the IssuedAt field if non-nil, zero value otherwise.

### GetIssuedAtOk

`func (o *IDTokenClaims) GetIssuedAtOk() (*time.Time, bool)`

GetIssuedAtOk returns a tuple with the IssuedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIssuedAt

`func (o *IDTokenClaims) SetIssuedAt(v time.Time)`

SetIssuedAt sets IssuedAt field to given value.

### HasIssuedAt

`func (o *IDTokenClaims) HasIssuedAt() bool`

HasIssuedAt returns a boolean if a field has been set.

### GetIssuer

`func (o *IDTokenClaims) GetIssuer() string`

GetIssuer returns the Issuer field if non-nil, zero value otherwise.

### GetIssuerOk

`func (o *IDTokenClaims) GetIssuerOk() (*string, bool)`

GetIssuerOk returns a tuple with the Issuer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIssuer

`func (o *IDTokenClaims) SetIssuer(v string)`

SetIssuer sets Issuer field to given value.

### HasIssuer

`func (o *IDTokenClaims) HasIssuer() bool`

HasIssuer returns a boolean if a field has been set.

### GetJTI

`func (o *IDTokenClaims) GetJTI() string`

GetJTI returns the JTI field if non-nil, zero value otherwise.

### GetJTIOk

`func (o *IDTokenClaims) GetJTIOk() (*string, bool)`

GetJTIOk returns a tuple with the JTI field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetJTI

`func (o *IDTokenClaims) SetJTI(v string)`

SetJTI sets JTI field to given value.

### HasJTI

`func (o *IDTokenClaims) HasJTI() bool`

HasJTI returns a boolean if a field has been set.

### GetNonce

`func (o *IDTokenClaims) GetNonce() string`

GetNonce returns the Nonce field if non-nil, zero value otherwise.

### GetNonceOk

`func (o *IDTokenClaims) GetNonceOk() (*string, bool)`

GetNonceOk returns a tuple with the Nonce field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNonce

`func (o *IDTokenClaims) SetNonce(v string)`

SetNonce sets Nonce field to given value.

### HasNonce

`func (o *IDTokenClaims) HasNonce() bool`

HasNonce returns a boolean if a field has been set.

### GetRequestedAt

`func (o *IDTokenClaims) GetRequestedAt() time.Time`

GetRequestedAt returns the RequestedAt field if non-nil, zero value otherwise.

### GetRequestedAtOk

`func (o *IDTokenClaims) GetRequestedAtOk() (*time.Time, bool)`

GetRequestedAtOk returns a tuple with the RequestedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestedAt

`func (o *IDTokenClaims) SetRequestedAt(v time.Time)`

SetRequestedAt sets RequestedAt field to given value.

### HasRequestedAt

`func (o *IDTokenClaims) HasRequestedAt() bool`

HasRequestedAt returns a boolean if a field has been set.

### GetSubject

`func (o *IDTokenClaims) GetSubject() string`

GetSubject returns the Subject field if non-nil, zero value otherwise.

### GetSubjectOk

`func (o *IDTokenClaims) GetSubjectOk() (*string, bool)`

GetSubjectOk returns a tuple with the Subject field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubject

`func (o *IDTokenClaims) SetSubject(v string)`

SetSubject sets Subject field to given value.

### HasSubject

`func (o *IDTokenClaims) HasSubject() bool`

HasSubject returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


