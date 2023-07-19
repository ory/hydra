# VerifiableCredentialProof

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Jwt** | Pointer to **string** |  | [optional] 
**ProofType** | Pointer to **string** |  | [optional] 

## Methods

### NewVerifiableCredentialProof

`func NewVerifiableCredentialProof() *VerifiableCredentialProof`

NewVerifiableCredentialProof instantiates a new VerifiableCredentialProof object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVerifiableCredentialProofWithDefaults

`func NewVerifiableCredentialProofWithDefaults() *VerifiableCredentialProof`

NewVerifiableCredentialProofWithDefaults instantiates a new VerifiableCredentialProof object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetJwt

`func (o *VerifiableCredentialProof) GetJwt() string`

GetJwt returns the Jwt field if non-nil, zero value otherwise.

### GetJwtOk

`func (o *VerifiableCredentialProof) GetJwtOk() (*string, bool)`

GetJwtOk returns a tuple with the Jwt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetJwt

`func (o *VerifiableCredentialProof) SetJwt(v string)`

SetJwt sets Jwt field to given value.

### HasJwt

`func (o *VerifiableCredentialProof) HasJwt() bool`

HasJwt returns a boolean if a field has been set.

### GetProofType

`func (o *VerifiableCredentialProof) GetProofType() string`

GetProofType returns the ProofType field if non-nil, zero value otherwise.

### GetProofTypeOk

`func (o *VerifiableCredentialProof) GetProofTypeOk() (*string, bool)`

GetProofTypeOk returns a tuple with the ProofType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProofType

`func (o *VerifiableCredentialProof) SetProofType(v string)`

SetProofType sets ProofType field to given value.

### HasProofType

`func (o *VerifiableCredentialProof) HasProofType() bool`

HasProofType returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


