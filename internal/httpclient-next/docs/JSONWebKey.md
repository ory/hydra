# JSONWebKey

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Alg** | **string** | The \&quot;alg\&quot; (algorithm) parameter identifies the algorithm intended for use with the key.  The values used should either be registered in the IANA \&quot;JSON Web Signature and Encryption Algorithms\&quot; registry established by [JWA] or be a value that contains a Collision- Resistant Name. | 
**Crv** | Pointer to **string** |  | [optional] 
**D** | Pointer to **string** |  | [optional] 
**Dp** | Pointer to **string** |  | [optional] 
**Dq** | Pointer to **string** |  | [optional] 
**E** | Pointer to **string** |  | [optional] 
**K** | Pointer to **string** |  | [optional] 
**Kid** | **string** | The \&quot;kid\&quot; (key ID) parameter is used to match a specific key.  This is used, for instance, to choose among a set of keys within a JWK Set during key rollover.  The structure of the \&quot;kid\&quot; value is unspecified.  When \&quot;kid\&quot; values are used within a JWK Set, different keys within the JWK Set SHOULD use distinct \&quot;kid\&quot; values.  (One example in which different keys might use the same \&quot;kid\&quot; value is if they have different \&quot;kty\&quot; (key type) values but are considered to be equivalent alternatives by the application using them.)  The \&quot;kid\&quot; value is a case-sensitive string. | 
**Kty** | **string** | The \&quot;kty\&quot; (key type) parameter identifies the cryptographic algorithm family used with the key, such as \&quot;RSA\&quot; or \&quot;EC\&quot;. \&quot;kty\&quot; values should either be registered in the IANA \&quot;JSON Web Key Types\&quot; registry established by [JWA] or be a value that contains a Collision- Resistant Name.  The \&quot;kty\&quot; value is a case-sensitive string. | 
**N** | Pointer to **string** |  | [optional] 
**P** | Pointer to **string** |  | [optional] 
**Q** | Pointer to **string** |  | [optional] 
**Qi** | Pointer to **string** |  | [optional] 
**Use** | **string** | Use (\&quot;public key use\&quot;) identifies the intended use of the public key. The \&quot;use\&quot; parameter is employed to indicate whether a public key is used for encrypting data or verifying the signature on data. Values are commonly \&quot;sig\&quot; (signature) or \&quot;enc\&quot; (encryption). | 
**X** | Pointer to **string** |  | [optional] 
**X5c** | Pointer to **[]string** | The \&quot;x5c\&quot; (X.509 certificate chain) parameter contains a chain of one or more PKIX certificates [RFC5280].  The certificate chain is represented as a JSON array of certificate value strings.  Each string in the array is a base64-encoded (Section 4 of [RFC4648] -- not base64url-encoded) DER [ITU.X690.1994] PKIX certificate value. The PKIX certificate containing the key value MUST be the first certificate. | [optional] 
**Y** | Pointer to **string** |  | [optional] 

## Methods

### NewJSONWebKey

`func NewJSONWebKey(alg string, kid string, kty string, use string, ) *JSONWebKey`

NewJSONWebKey instantiates a new JSONWebKey object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewJSONWebKeyWithDefaults

`func NewJSONWebKeyWithDefaults() *JSONWebKey`

NewJSONWebKeyWithDefaults instantiates a new JSONWebKey object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAlg

`func (o *JSONWebKey) GetAlg() string`

GetAlg returns the Alg field if non-nil, zero value otherwise.

### GetAlgOk

`func (o *JSONWebKey) GetAlgOk() (*string, bool)`

GetAlgOk returns a tuple with the Alg field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAlg

`func (o *JSONWebKey) SetAlg(v string)`

SetAlg sets Alg field to given value.


### GetCrv

`func (o *JSONWebKey) GetCrv() string`

GetCrv returns the Crv field if non-nil, zero value otherwise.

### GetCrvOk

`func (o *JSONWebKey) GetCrvOk() (*string, bool)`

GetCrvOk returns a tuple with the Crv field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCrv

`func (o *JSONWebKey) SetCrv(v string)`

SetCrv sets Crv field to given value.

### HasCrv

`func (o *JSONWebKey) HasCrv() bool`

HasCrv returns a boolean if a field has been set.

### GetD

`func (o *JSONWebKey) GetD() string`

GetD returns the D field if non-nil, zero value otherwise.

### GetDOk

`func (o *JSONWebKey) GetDOk() (*string, bool)`

GetDOk returns a tuple with the D field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetD

`func (o *JSONWebKey) SetD(v string)`

SetD sets D field to given value.

### HasD

`func (o *JSONWebKey) HasD() bool`

HasD returns a boolean if a field has been set.

### GetDp

`func (o *JSONWebKey) GetDp() string`

GetDp returns the Dp field if non-nil, zero value otherwise.

### GetDpOk

`func (o *JSONWebKey) GetDpOk() (*string, bool)`

GetDpOk returns a tuple with the Dp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDp

`func (o *JSONWebKey) SetDp(v string)`

SetDp sets Dp field to given value.

### HasDp

`func (o *JSONWebKey) HasDp() bool`

HasDp returns a boolean if a field has been set.

### GetDq

`func (o *JSONWebKey) GetDq() string`

GetDq returns the Dq field if non-nil, zero value otherwise.

### GetDqOk

`func (o *JSONWebKey) GetDqOk() (*string, bool)`

GetDqOk returns a tuple with the Dq field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDq

`func (o *JSONWebKey) SetDq(v string)`

SetDq sets Dq field to given value.

### HasDq

`func (o *JSONWebKey) HasDq() bool`

HasDq returns a boolean if a field has been set.

### GetE

`func (o *JSONWebKey) GetE() string`

GetE returns the E field if non-nil, zero value otherwise.

### GetEOk

`func (o *JSONWebKey) GetEOk() (*string, bool)`

GetEOk returns a tuple with the E field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetE

`func (o *JSONWebKey) SetE(v string)`

SetE sets E field to given value.

### HasE

`func (o *JSONWebKey) HasE() bool`

HasE returns a boolean if a field has been set.

### GetK

`func (o *JSONWebKey) GetK() string`

GetK returns the K field if non-nil, zero value otherwise.

### GetKOk

`func (o *JSONWebKey) GetKOk() (*string, bool)`

GetKOk returns a tuple with the K field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetK

`func (o *JSONWebKey) SetK(v string)`

SetK sets K field to given value.

### HasK

`func (o *JSONWebKey) HasK() bool`

HasK returns a boolean if a field has been set.

### GetKid

`func (o *JSONWebKey) GetKid() string`

GetKid returns the Kid field if non-nil, zero value otherwise.

### GetKidOk

`func (o *JSONWebKey) GetKidOk() (*string, bool)`

GetKidOk returns a tuple with the Kid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKid

`func (o *JSONWebKey) SetKid(v string)`

SetKid sets Kid field to given value.


### GetKty

`func (o *JSONWebKey) GetKty() string`

GetKty returns the Kty field if non-nil, zero value otherwise.

### GetKtyOk

`func (o *JSONWebKey) GetKtyOk() (*string, bool)`

GetKtyOk returns a tuple with the Kty field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKty

`func (o *JSONWebKey) SetKty(v string)`

SetKty sets Kty field to given value.


### GetN

`func (o *JSONWebKey) GetN() string`

GetN returns the N field if non-nil, zero value otherwise.

### GetNOk

`func (o *JSONWebKey) GetNOk() (*string, bool)`

GetNOk returns a tuple with the N field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetN

`func (o *JSONWebKey) SetN(v string)`

SetN sets N field to given value.

### HasN

`func (o *JSONWebKey) HasN() bool`

HasN returns a boolean if a field has been set.

### GetP

`func (o *JSONWebKey) GetP() string`

GetP returns the P field if non-nil, zero value otherwise.

### GetPOk

`func (o *JSONWebKey) GetPOk() (*string, bool)`

GetPOk returns a tuple with the P field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetP

`func (o *JSONWebKey) SetP(v string)`

SetP sets P field to given value.

### HasP

`func (o *JSONWebKey) HasP() bool`

HasP returns a boolean if a field has been set.

### GetQ

`func (o *JSONWebKey) GetQ() string`

GetQ returns the Q field if non-nil, zero value otherwise.

### GetQOk

`func (o *JSONWebKey) GetQOk() (*string, bool)`

GetQOk returns a tuple with the Q field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQ

`func (o *JSONWebKey) SetQ(v string)`

SetQ sets Q field to given value.

### HasQ

`func (o *JSONWebKey) HasQ() bool`

HasQ returns a boolean if a field has been set.

### GetQi

`func (o *JSONWebKey) GetQi() string`

GetQi returns the Qi field if non-nil, zero value otherwise.

### GetQiOk

`func (o *JSONWebKey) GetQiOk() (*string, bool)`

GetQiOk returns a tuple with the Qi field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQi

`func (o *JSONWebKey) SetQi(v string)`

SetQi sets Qi field to given value.

### HasQi

`func (o *JSONWebKey) HasQi() bool`

HasQi returns a boolean if a field has been set.

### GetUse

`func (o *JSONWebKey) GetUse() string`

GetUse returns the Use field if non-nil, zero value otherwise.

### GetUseOk

`func (o *JSONWebKey) GetUseOk() (*string, bool)`

GetUseOk returns a tuple with the Use field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUse

`func (o *JSONWebKey) SetUse(v string)`

SetUse sets Use field to given value.


### GetX

`func (o *JSONWebKey) GetX() string`

GetX returns the X field if non-nil, zero value otherwise.

### GetXOk

`func (o *JSONWebKey) GetXOk() (*string, bool)`

GetXOk returns a tuple with the X field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetX

`func (o *JSONWebKey) SetX(v string)`

SetX sets X field to given value.

### HasX

`func (o *JSONWebKey) HasX() bool`

HasX returns a boolean if a field has been set.

### GetX5c

`func (o *JSONWebKey) GetX5c() []string`

GetX5c returns the X5c field if non-nil, zero value otherwise.

### GetX5cOk

`func (o *JSONWebKey) GetX5cOk() (*[]string, bool)`

GetX5cOk returns a tuple with the X5c field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetX5c

`func (o *JSONWebKey) SetX5c(v []string)`

SetX5c sets X5c field to given value.

### HasX5c

`func (o *JSONWebKey) HasX5c() bool`

HasX5c returns a boolean if a field has been set.

### GetY

`func (o *JSONWebKey) GetY() string`

GetY returns the Y field if non-nil, zero value otherwise.

### GetYOk

`func (o *JSONWebKey) GetYOk() (*string, bool)`

GetYOk returns a tuple with the Y field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetY

`func (o *JSONWebKey) SetY(v string)`

SetY sets Y field to given value.

### HasY

`func (o *JSONWebKey) HasY() bool`

HasY returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


