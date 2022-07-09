# TrustedOAuth2JwtGrantIssuer

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AllowAnySubject** | Pointer to **bool** | The \&quot;allow_any_subject\&quot; indicates that the issuer is allowed to have any principal as the subject of the JWT. | [optional] 
**CreatedAt** | Pointer to **time.Time** | The \&quot;created_at\&quot; indicates, when grant was created. | [optional] 
**ExpiresAt** | Pointer to **time.Time** | The \&quot;expires_at\&quot; indicates, when grant will expire, so we will reject assertion from \&quot;issuer\&quot; targeting \&quot;subject\&quot;. | [optional] 
**Id** | Pointer to **string** |  | [optional] 
**Issuer** | Pointer to **string** | The \&quot;issuer\&quot; identifies the principal that issued the JWT assertion (same as \&quot;iss\&quot; claim in JWT). | [optional] 
**PublicKey** | Pointer to [**TrustedOAuth2JwtGrantJsonWebKey**](TrustedOAuth2JwtGrantJsonWebKey.md) |  | [optional] 
**Scope** | Pointer to **[]string** | The \&quot;scope\&quot; contains list of scope values (as described in Section 3.3 of OAuth 2.0 [RFC6749]) | [optional] 
**Subject** | Pointer to **string** | The \&quot;subject\&quot; identifies the principal that is the subject of the JWT. | [optional] 

## Methods

### NewTrustedOAuth2JwtGrantIssuer

`func NewTrustedOAuth2JwtGrantIssuer() *TrustedOAuth2JwtGrantIssuer`

NewTrustedOAuth2JwtGrantIssuer instantiates a new TrustedOAuth2JwtGrantIssuer object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTrustedOAuth2JwtGrantIssuerWithDefaults

`func NewTrustedOAuth2JwtGrantIssuerWithDefaults() *TrustedOAuth2JwtGrantIssuer`

NewTrustedOAuth2JwtGrantIssuerWithDefaults instantiates a new TrustedOAuth2JwtGrantIssuer object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAllowAnySubject

`func (o *TrustedOAuth2JwtGrantIssuer) GetAllowAnySubject() bool`

GetAllowAnySubject returns the AllowAnySubject field if non-nil, zero value otherwise.

### GetAllowAnySubjectOk

`func (o *TrustedOAuth2JwtGrantIssuer) GetAllowAnySubjectOk() (*bool, bool)`

GetAllowAnySubjectOk returns a tuple with the AllowAnySubject field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllowAnySubject

`func (o *TrustedOAuth2JwtGrantIssuer) SetAllowAnySubject(v bool)`

SetAllowAnySubject sets AllowAnySubject field to given value.

### HasAllowAnySubject

`func (o *TrustedOAuth2JwtGrantIssuer) HasAllowAnySubject() bool`

HasAllowAnySubject returns a boolean if a field has been set.

### GetCreatedAt

`func (o *TrustedOAuth2JwtGrantIssuer) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *TrustedOAuth2JwtGrantIssuer) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *TrustedOAuth2JwtGrantIssuer) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *TrustedOAuth2JwtGrantIssuer) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetExpiresAt

`func (o *TrustedOAuth2JwtGrantIssuer) GetExpiresAt() time.Time`

GetExpiresAt returns the ExpiresAt field if non-nil, zero value otherwise.

### GetExpiresAtOk

`func (o *TrustedOAuth2JwtGrantIssuer) GetExpiresAtOk() (*time.Time, bool)`

GetExpiresAtOk returns a tuple with the ExpiresAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiresAt

`func (o *TrustedOAuth2JwtGrantIssuer) SetExpiresAt(v time.Time)`

SetExpiresAt sets ExpiresAt field to given value.

### HasExpiresAt

`func (o *TrustedOAuth2JwtGrantIssuer) HasExpiresAt() bool`

HasExpiresAt returns a boolean if a field has been set.

### GetId

`func (o *TrustedOAuth2JwtGrantIssuer) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *TrustedOAuth2JwtGrantIssuer) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *TrustedOAuth2JwtGrantIssuer) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *TrustedOAuth2JwtGrantIssuer) HasId() bool`

HasId returns a boolean if a field has been set.

### GetIssuer

`func (o *TrustedOAuth2JwtGrantIssuer) GetIssuer() string`

GetIssuer returns the Issuer field if non-nil, zero value otherwise.

### GetIssuerOk

`func (o *TrustedOAuth2JwtGrantIssuer) GetIssuerOk() (*string, bool)`

GetIssuerOk returns a tuple with the Issuer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIssuer

`func (o *TrustedOAuth2JwtGrantIssuer) SetIssuer(v string)`

SetIssuer sets Issuer field to given value.

### HasIssuer

`func (o *TrustedOAuth2JwtGrantIssuer) HasIssuer() bool`

HasIssuer returns a boolean if a field has been set.

### GetPublicKey

`func (o *TrustedOAuth2JwtGrantIssuer) GetPublicKey() TrustedOAuth2JwtGrantJsonWebKey`

GetPublicKey returns the PublicKey field if non-nil, zero value otherwise.

### GetPublicKeyOk

`func (o *TrustedOAuth2JwtGrantIssuer) GetPublicKeyOk() (*TrustedOAuth2JwtGrantJsonWebKey, bool)`

GetPublicKeyOk returns a tuple with the PublicKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPublicKey

`func (o *TrustedOAuth2JwtGrantIssuer) SetPublicKey(v TrustedOAuth2JwtGrantJsonWebKey)`

SetPublicKey sets PublicKey field to given value.

### HasPublicKey

`func (o *TrustedOAuth2JwtGrantIssuer) HasPublicKey() bool`

HasPublicKey returns a boolean if a field has been set.

### GetScope

`func (o *TrustedOAuth2JwtGrantIssuer) GetScope() []string`

GetScope returns the Scope field if non-nil, zero value otherwise.

### GetScopeOk

`func (o *TrustedOAuth2JwtGrantIssuer) GetScopeOk() (*[]string, bool)`

GetScopeOk returns a tuple with the Scope field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScope

`func (o *TrustedOAuth2JwtGrantIssuer) SetScope(v []string)`

SetScope sets Scope field to given value.

### HasScope

`func (o *TrustedOAuth2JwtGrantIssuer) HasScope() bool`

HasScope returns a boolean if a field has been set.

### GetSubject

`func (o *TrustedOAuth2JwtGrantIssuer) GetSubject() string`

GetSubject returns the Subject field if non-nil, zero value otherwise.

### GetSubjectOk

`func (o *TrustedOAuth2JwtGrantIssuer) GetSubjectOk() (*string, bool)`

GetSubjectOk returns a tuple with the Subject field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubject

`func (o *TrustedOAuth2JwtGrantIssuer) SetSubject(v string)`

SetSubject sets Subject field to given value.

### HasSubject

`func (o *TrustedOAuth2JwtGrantIssuer) HasSubject() bool`

HasSubject returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


