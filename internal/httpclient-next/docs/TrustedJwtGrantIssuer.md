# TrustedJwtGrantIssuer

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AllowAnySubject** | Pointer to **bool** | The \&quot;allow_any_subject\&quot; indicates that the issuer is allowed to have any principal as the subject of the JWT. | [optional] 
**CreatedAt** | Pointer to **time.Time** | The \&quot;created_at\&quot; indicates, when grant was created. | [optional] 
**ExpiresAt** | Pointer to **time.Time** | The \&quot;expires_at\&quot; indicates, when grant will expire, so we will reject assertion from \&quot;issuer\&quot; targeting \&quot;subject\&quot;. | [optional] 
**Id** | Pointer to **string** |  | [optional] 
**Issuer** | Pointer to **string** | The \&quot;issuer\&quot; identifies the principal that issued the JWT assertion (same as \&quot;iss\&quot; claim in JWT). | [optional] 
**PublicKey** | Pointer to [**TrustedJsonWebKey**](TrustedJsonWebKey.md) |  | [optional] 
**Scope** | Pointer to **[]string** | The \&quot;scope\&quot; contains list of scope values (as described in Section 3.3 of OAuth 2.0 [RFC6749]) | [optional] 
**Subject** | Pointer to **string** | The \&quot;subject\&quot; identifies the principal that is the subject of the JWT. | [optional] 

## Methods

### NewTrustedJwtGrantIssuer

`func NewTrustedJwtGrantIssuer() *TrustedJwtGrantIssuer`

NewTrustedJwtGrantIssuer instantiates a new TrustedJwtGrantIssuer object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTrustedJwtGrantIssuerWithDefaults

`func NewTrustedJwtGrantIssuerWithDefaults() *TrustedJwtGrantIssuer`

NewTrustedJwtGrantIssuerWithDefaults instantiates a new TrustedJwtGrantIssuer object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAllowAnySubject

`func (o *TrustedJwtGrantIssuer) GetAllowAnySubject() bool`

GetAllowAnySubject returns the AllowAnySubject field if non-nil, zero value otherwise.

### GetAllowAnySubjectOk

`func (o *TrustedJwtGrantIssuer) GetAllowAnySubjectOk() (*bool, bool)`

GetAllowAnySubjectOk returns a tuple with the AllowAnySubject field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllowAnySubject

`func (o *TrustedJwtGrantIssuer) SetAllowAnySubject(v bool)`

SetAllowAnySubject sets AllowAnySubject field to given value.

### HasAllowAnySubject

`func (o *TrustedJwtGrantIssuer) HasAllowAnySubject() bool`

HasAllowAnySubject returns a boolean if a field has been set.

### GetCreatedAt

`func (o *TrustedJwtGrantIssuer) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *TrustedJwtGrantIssuer) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *TrustedJwtGrantIssuer) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *TrustedJwtGrantIssuer) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetExpiresAt

`func (o *TrustedJwtGrantIssuer) GetExpiresAt() time.Time`

GetExpiresAt returns the ExpiresAt field if non-nil, zero value otherwise.

### GetExpiresAtOk

`func (o *TrustedJwtGrantIssuer) GetExpiresAtOk() (*time.Time, bool)`

GetExpiresAtOk returns a tuple with the ExpiresAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiresAt

`func (o *TrustedJwtGrantIssuer) SetExpiresAt(v time.Time)`

SetExpiresAt sets ExpiresAt field to given value.

### HasExpiresAt

`func (o *TrustedJwtGrantIssuer) HasExpiresAt() bool`

HasExpiresAt returns a boolean if a field has been set.

### GetId

`func (o *TrustedJwtGrantIssuer) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *TrustedJwtGrantIssuer) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *TrustedJwtGrantIssuer) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *TrustedJwtGrantIssuer) HasId() bool`

HasId returns a boolean if a field has been set.

### GetIssuer

`func (o *TrustedJwtGrantIssuer) GetIssuer() string`

GetIssuer returns the Issuer field if non-nil, zero value otherwise.

### GetIssuerOk

`func (o *TrustedJwtGrantIssuer) GetIssuerOk() (*string, bool)`

GetIssuerOk returns a tuple with the Issuer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIssuer

`func (o *TrustedJwtGrantIssuer) SetIssuer(v string)`

SetIssuer sets Issuer field to given value.

### HasIssuer

`func (o *TrustedJwtGrantIssuer) HasIssuer() bool`

HasIssuer returns a boolean if a field has been set.

### GetPublicKey

`func (o *TrustedJwtGrantIssuer) GetPublicKey() TrustedJsonWebKey`

GetPublicKey returns the PublicKey field if non-nil, zero value otherwise.

### GetPublicKeyOk

`func (o *TrustedJwtGrantIssuer) GetPublicKeyOk() (*TrustedJsonWebKey, bool)`

GetPublicKeyOk returns a tuple with the PublicKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPublicKey

`func (o *TrustedJwtGrantIssuer) SetPublicKey(v TrustedJsonWebKey)`

SetPublicKey sets PublicKey field to given value.

### HasPublicKey

`func (o *TrustedJwtGrantIssuer) HasPublicKey() bool`

HasPublicKey returns a boolean if a field has been set.

### GetScope

`func (o *TrustedJwtGrantIssuer) GetScope() []string`

GetScope returns the Scope field if non-nil, zero value otherwise.

### GetScopeOk

`func (o *TrustedJwtGrantIssuer) GetScopeOk() (*[]string, bool)`

GetScopeOk returns a tuple with the Scope field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScope

`func (o *TrustedJwtGrantIssuer) SetScope(v []string)`

SetScope sets Scope field to given value.

### HasScope

`func (o *TrustedJwtGrantIssuer) HasScope() bool`

HasScope returns a boolean if a field has been set.

### GetSubject

`func (o *TrustedJwtGrantIssuer) GetSubject() string`

GetSubject returns the Subject field if non-nil, zero value otherwise.

### GetSubjectOk

`func (o *TrustedJwtGrantIssuer) GetSubjectOk() (*string, bool)`

GetSubjectOk returns a tuple with the Subject field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubject

`func (o *TrustedJwtGrantIssuer) SetSubject(v string)`

SetSubject sets Subject field to given value.

### HasSubject

`func (o *TrustedJwtGrantIssuer) HasSubject() bool`

HasSubject returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


