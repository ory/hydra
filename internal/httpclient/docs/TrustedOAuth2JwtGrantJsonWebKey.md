# TrustedOAuth2JwtGrantJsonWebKey

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Kid** | Pointer to **string** | The \&quot;key_id\&quot; is key unique identifier (same as kid header in jws/jwt). | [optional] 
**Set** | Pointer to **string** | The \&quot;set\&quot; is basically a name for a group(set) of keys. Will be the same as \&quot;issuer\&quot; in grant. | [optional] 

## Methods

### NewTrustedOAuth2JwtGrantJsonWebKey

`func NewTrustedOAuth2JwtGrantJsonWebKey() *TrustedOAuth2JwtGrantJsonWebKey`

NewTrustedOAuth2JwtGrantJsonWebKey instantiates a new TrustedOAuth2JwtGrantJsonWebKey object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTrustedOAuth2JwtGrantJsonWebKeyWithDefaults

`func NewTrustedOAuth2JwtGrantJsonWebKeyWithDefaults() *TrustedOAuth2JwtGrantJsonWebKey`

NewTrustedOAuth2JwtGrantJsonWebKeyWithDefaults instantiates a new TrustedOAuth2JwtGrantJsonWebKey object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKid

`func (o *TrustedOAuth2JwtGrantJsonWebKey) GetKid() string`

GetKid returns the Kid field if non-nil, zero value otherwise.

### GetKidOk

`func (o *TrustedOAuth2JwtGrantJsonWebKey) GetKidOk() (*string, bool)`

GetKidOk returns a tuple with the Kid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKid

`func (o *TrustedOAuth2JwtGrantJsonWebKey) SetKid(v string)`

SetKid sets Kid field to given value.

### HasKid

`func (o *TrustedOAuth2JwtGrantJsonWebKey) HasKid() bool`

HasKid returns a boolean if a field has been set.

### GetSet

`func (o *TrustedOAuth2JwtGrantJsonWebKey) GetSet() string`

GetSet returns the Set field if non-nil, zero value otherwise.

### GetSetOk

`func (o *TrustedOAuth2JwtGrantJsonWebKey) GetSetOk() (*string, bool)`

GetSetOk returns a tuple with the Set field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSet

`func (o *TrustedOAuth2JwtGrantJsonWebKey) SetSet(v string)`

SetSet sets Set field to given value.

### HasSet

`func (o *TrustedOAuth2JwtGrantJsonWebKey) HasSet() bool`

HasSet returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


