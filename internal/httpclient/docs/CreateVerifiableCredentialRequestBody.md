# CreateVerifiableCredentialRequestBody

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Format** | Pointer to **string** |  | [optional] 
**Proof** | Pointer to [**VerifiableCredentialProof**](VerifiableCredentialProof.md) |  | [optional] 
**Types** | Pointer to **[]string** |  | [optional] 

## Methods

### NewCreateVerifiableCredentialRequestBody

`func NewCreateVerifiableCredentialRequestBody() *CreateVerifiableCredentialRequestBody`

NewCreateVerifiableCredentialRequestBody instantiates a new CreateVerifiableCredentialRequestBody object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateVerifiableCredentialRequestBodyWithDefaults

`func NewCreateVerifiableCredentialRequestBodyWithDefaults() *CreateVerifiableCredentialRequestBody`

NewCreateVerifiableCredentialRequestBodyWithDefaults instantiates a new CreateVerifiableCredentialRequestBody object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetFormat

`func (o *CreateVerifiableCredentialRequestBody) GetFormat() string`

GetFormat returns the Format field if non-nil, zero value otherwise.

### GetFormatOk

`func (o *CreateVerifiableCredentialRequestBody) GetFormatOk() (*string, bool)`

GetFormatOk returns a tuple with the Format field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFormat

`func (o *CreateVerifiableCredentialRequestBody) SetFormat(v string)`

SetFormat sets Format field to given value.

### HasFormat

`func (o *CreateVerifiableCredentialRequestBody) HasFormat() bool`

HasFormat returns a boolean if a field has been set.

### GetProof

`func (o *CreateVerifiableCredentialRequestBody) GetProof() VerifiableCredentialProof`

GetProof returns the Proof field if non-nil, zero value otherwise.

### GetProofOk

`func (o *CreateVerifiableCredentialRequestBody) GetProofOk() (*VerifiableCredentialProof, bool)`

GetProofOk returns a tuple with the Proof field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProof

`func (o *CreateVerifiableCredentialRequestBody) SetProof(v VerifiableCredentialProof)`

SetProof sets Proof field to given value.

### HasProof

`func (o *CreateVerifiableCredentialRequestBody) HasProof() bool`

HasProof returns a boolean if a field has been set.

### GetTypes

`func (o *CreateVerifiableCredentialRequestBody) GetTypes() []string`

GetTypes returns the Types field if non-nil, zero value otherwise.

### GetTypesOk

`func (o *CreateVerifiableCredentialRequestBody) GetTypesOk() (*[]string, bool)`

GetTypesOk returns a tuple with the Types field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTypes

`func (o *CreateVerifiableCredentialRequestBody) SetTypes(v []string)`

SetTypes sets Types field to given value.

### HasTypes

`func (o *CreateVerifiableCredentialRequestBody) HasTypes() bool`

HasTypes returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


