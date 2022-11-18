# CreateJsonWebKeySet

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Alg** | **string** | JSON Web Key Algorithm  The algorithm to be used for creating the key. Supports &#x60;RS256&#x60;, &#x60;ES256&#x60;, &#x60;ES512&#x60;, &#x60;HS512&#x60;, and &#x60;HS256&#x60;. | 
**Kid** | **string** | JSON Web Key ID  The Key ID of the key to be created. | 
**Use** | **string** | JSON Web Key Use  The \&quot;use\&quot; (public key use) parameter identifies the intended use of the public key. The \&quot;use\&quot; parameter is employed to indicate whether a public key is used for encrypting data or verifying the signature on data. Valid values are \&quot;enc\&quot; and \&quot;sig\&quot;. | 

## Methods

### NewCreateJsonWebKeySet

`func NewCreateJsonWebKeySet(alg string, kid string, use string, ) *CreateJsonWebKeySet`

NewCreateJsonWebKeySet instantiates a new CreateJsonWebKeySet object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateJsonWebKeySetWithDefaults

`func NewCreateJsonWebKeySetWithDefaults() *CreateJsonWebKeySet`

NewCreateJsonWebKeySetWithDefaults instantiates a new CreateJsonWebKeySet object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAlg

`func (o *CreateJsonWebKeySet) GetAlg() string`

GetAlg returns the Alg field if non-nil, zero value otherwise.

### GetAlgOk

`func (o *CreateJsonWebKeySet) GetAlgOk() (*string, bool)`

GetAlgOk returns a tuple with the Alg field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAlg

`func (o *CreateJsonWebKeySet) SetAlg(v string)`

SetAlg sets Alg field to given value.


### GetKid

`func (o *CreateJsonWebKeySet) GetKid() string`

GetKid returns the Kid field if non-nil, zero value otherwise.

### GetKidOk

`func (o *CreateJsonWebKeySet) GetKidOk() (*string, bool)`

GetKidOk returns a tuple with the Kid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKid

`func (o *CreateJsonWebKeySet) SetKid(v string)`

SetKid sets Kid field to given value.


### GetUse

`func (o *CreateJsonWebKeySet) GetUse() string`

GetUse returns the Use field if non-nil, zero value otherwise.

### GetUseOk

`func (o *CreateJsonWebKeySet) GetUseOk() (*string, bool)`

GetUseOk returns a tuple with the Use field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUse

`func (o *CreateJsonWebKeySet) SetUse(v string)`

SetUse sets Use field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


