# CredentialSupportedDraft00

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CryptographicBindingMethodsSupported** | Pointer to **[]string** | OpenID Connect Verifiable Credentials Cryptographic Binding Methods Supported  Contains a list of cryptographic binding methods supported for signing the proof. | [optional] 
**CryptographicSuitesSupported** | Pointer to **[]string** | OpenID Connect Verifiable Credentials Cryptographic Suites Supported  Contains a list of cryptographic suites methods supported for signing the proof. | [optional] 
**Format** | Pointer to **string** | OpenID Connect Verifiable Credentials Format  Contains the format that is supported by this authorization server. | [optional] 
**Types** | Pointer to **[]string** | OpenID Connect Verifiable Credentials Types  Contains the types of verifiable credentials supported. | [optional] 

## Methods

### NewCredentialSupportedDraft00

`func NewCredentialSupportedDraft00() *CredentialSupportedDraft00`

NewCredentialSupportedDraft00 instantiates a new CredentialSupportedDraft00 object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCredentialSupportedDraft00WithDefaults

`func NewCredentialSupportedDraft00WithDefaults() *CredentialSupportedDraft00`

NewCredentialSupportedDraft00WithDefaults instantiates a new CredentialSupportedDraft00 object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCryptographicBindingMethodsSupported

`func (o *CredentialSupportedDraft00) GetCryptographicBindingMethodsSupported() []string`

GetCryptographicBindingMethodsSupported returns the CryptographicBindingMethodsSupported field if non-nil, zero value otherwise.

### GetCryptographicBindingMethodsSupportedOk

`func (o *CredentialSupportedDraft00) GetCryptographicBindingMethodsSupportedOk() (*[]string, bool)`

GetCryptographicBindingMethodsSupportedOk returns a tuple with the CryptographicBindingMethodsSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCryptographicBindingMethodsSupported

`func (o *CredentialSupportedDraft00) SetCryptographicBindingMethodsSupported(v []string)`

SetCryptographicBindingMethodsSupported sets CryptographicBindingMethodsSupported field to given value.

### HasCryptographicBindingMethodsSupported

`func (o *CredentialSupportedDraft00) HasCryptographicBindingMethodsSupported() bool`

HasCryptographicBindingMethodsSupported returns a boolean if a field has been set.

### GetCryptographicSuitesSupported

`func (o *CredentialSupportedDraft00) GetCryptographicSuitesSupported() []string`

GetCryptographicSuitesSupported returns the CryptographicSuitesSupported field if non-nil, zero value otherwise.

### GetCryptographicSuitesSupportedOk

`func (o *CredentialSupportedDraft00) GetCryptographicSuitesSupportedOk() (*[]string, bool)`

GetCryptographicSuitesSupportedOk returns a tuple with the CryptographicSuitesSupported field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCryptographicSuitesSupported

`func (o *CredentialSupportedDraft00) SetCryptographicSuitesSupported(v []string)`

SetCryptographicSuitesSupported sets CryptographicSuitesSupported field to given value.

### HasCryptographicSuitesSupported

`func (o *CredentialSupportedDraft00) HasCryptographicSuitesSupported() bool`

HasCryptographicSuitesSupported returns a boolean if a field has been set.

### GetFormat

`func (o *CredentialSupportedDraft00) GetFormat() string`

GetFormat returns the Format field if non-nil, zero value otherwise.

### GetFormatOk

`func (o *CredentialSupportedDraft00) GetFormatOk() (*string, bool)`

GetFormatOk returns a tuple with the Format field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFormat

`func (o *CredentialSupportedDraft00) SetFormat(v string)`

SetFormat sets Format field to given value.

### HasFormat

`func (o *CredentialSupportedDraft00) HasFormat() bool`

HasFormat returns a boolean if a field has been set.

### GetTypes

`func (o *CredentialSupportedDraft00) GetTypes() []string`

GetTypes returns the Types field if non-nil, zero value otherwise.

### GetTypesOk

`func (o *CredentialSupportedDraft00) GetTypesOk() (*[]string, bool)`

GetTypesOk returns a tuple with the Types field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTypes

`func (o *CredentialSupportedDraft00) SetTypes(v []string)`

SetTypes sets Types field to given value.

### HasTypes

`func (o *CredentialSupportedDraft00) HasTypes() bool`

HasTypes returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


