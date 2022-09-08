# AdminCreateJsonWebKeySetBody

## Properties

| Name    | Type       | Description                                                                                                                                                                                                                                                                                            | Notes |
| ------- | ---------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ----- |
| **Alg** | **string** | The algorithm to be used for creating the key. Supports \&quot;RS256\&quot;, \&quot;ES256\&quot;, \&quot;ES512\&quot;, \&quot;HS512\&quot;, and \&quot;HS256\&quot;                                                                                                                                    |
| **Kid** | **string** | The kid of the key to be created                                                                                                                                                                                                                                                                       |
| **Use** | **string** | The \&quot;use\&quot; (public key use) parameter identifies the intended use of the public key. The \&quot;use\&quot; parameter is employed to indicate whether a public key is used for encrypting data or verifying the signature on data. Valid values are \&quot;enc\&quot; and \&quot;sig\&quot;. |

## Methods

### NewAdminCreateJsonWebKeySetBody

`func NewAdminCreateJsonWebKeySetBody(alg string, kid string, use string, ) *AdminCreateJsonWebKeySetBody`

NewAdminCreateJsonWebKeySetBody instantiates a new AdminCreateJsonWebKeySetBody
object This constructor will assign default values to properties that have it
defined, and makes sure properties required by API are set, but the set of
arguments will change when the set of required properties is changed

### NewAdminCreateJsonWebKeySetBodyWithDefaults

`func NewAdminCreateJsonWebKeySetBodyWithDefaults() *AdminCreateJsonWebKeySetBody`

NewAdminCreateJsonWebKeySetBodyWithDefaults instantiates a new
AdminCreateJsonWebKeySetBody object This constructor will only assign default
values to properties that have it defined, but it doesn't guarantee that
properties required by API are set

### GetAlg

`func (o *AdminCreateJsonWebKeySetBody) GetAlg() string`

GetAlg returns the Alg field if non-nil, zero value otherwise.

### GetAlgOk

`func (o *AdminCreateJsonWebKeySetBody) GetAlgOk() (*string, bool)`

GetAlgOk returns a tuple with the Alg field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetAlg

`func (o *AdminCreateJsonWebKeySetBody) SetAlg(v string)`

SetAlg sets Alg field to given value.

### GetKid

`func (o *AdminCreateJsonWebKeySetBody) GetKid() string`

GetKid returns the Kid field if non-nil, zero value otherwise.

### GetKidOk

`func (o *AdminCreateJsonWebKeySetBody) GetKidOk() (*string, bool)`

GetKidOk returns a tuple with the Kid field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetKid

`func (o *AdminCreateJsonWebKeySetBody) SetKid(v string)`

SetKid sets Kid field to given value.

### GetUse

`func (o *AdminCreateJsonWebKeySetBody) GetUse() string`

GetUse returns the Use field if non-nil, zero value otherwise.

### GetUseOk

`func (o *AdminCreateJsonWebKeySetBody) GetUseOk() (*string, bool)`

GetUseOk returns a tuple with the Use field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetUse

`func (o *AdminCreateJsonWebKeySetBody) SetUse(v string)`

SetUse sets Use field to given value.

[[Back to Model list]](../README.md#documentation-for-models)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to README]](../README.md)
