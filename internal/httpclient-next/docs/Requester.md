# Requester

## Properties

| Name                | Type                    | Description                                                               | Notes      |
| ------------------- | ----------------------- | ------------------------------------------------------------------------- | ---------- |
| **ClientId**        | Pointer to **string**   | ClientID is the identifier of the OAuth 2.0 client.                       | [optional] |
| **GrantTypes**      | Pointer to **[]string** | GrantTypes is the requests grant types.                                   | [optional] |
| **GrantedAudience** | Pointer to **[]string** | GrantedAudience is the list of audiences granted to the OAuth 2.0 client. | [optional] |
| **GrantedScopes**   | Pointer to **[]string** | GrantedScopes is the list of scopes granted to the OAuth 2.0 client.      | [optional] |

## Methods

### NewRequester

`func NewRequester() *Requester`

NewRequester instantiates a new Requester object This constructor will assign
default values to properties that have it defined, and makes sure properties
required by API are set, but the set of arguments will change when the set of
required properties is changed

### NewRequesterWithDefaults

`func NewRequesterWithDefaults() *Requester`

NewRequesterWithDefaults instantiates a new Requester object This constructor
will only assign default values to properties that have it defined, but it
doesn't guarantee that properties required by API are set

### GetClientId

`func (o *Requester) GetClientId() string`

GetClientId returns the ClientId field if non-nil, zero value otherwise.

### GetClientIdOk

`func (o *Requester) GetClientIdOk() (*string, bool)`

GetClientIdOk returns a tuple with the ClientId field if it's non-nil, zero
value otherwise and a boolean to check if the value has been set.

### SetClientId

`func (o *Requester) SetClientId(v string)`

SetClientId sets ClientId field to given value.

### HasClientId

`func (o *Requester) HasClientId() bool`

HasClientId returns a boolean if a field has been set.

### GetGrantTypes

`func (o *Requester) GetGrantTypes() []string`

GetGrantTypes returns the GrantTypes field if non-nil, zero value otherwise.

### GetGrantTypesOk

`func (o *Requester) GetGrantTypesOk() (*[]string, bool)`

GetGrantTypesOk returns a tuple with the GrantTypes field if it's non-nil, zero
value otherwise and a boolean to check if the value has been set.

### SetGrantTypes

`func (o *Requester) SetGrantTypes(v []string)`

SetGrantTypes sets GrantTypes field to given value.

### HasGrantTypes

`func (o *Requester) HasGrantTypes() bool`

HasGrantTypes returns a boolean if a field has been set.

### GetGrantedAudience

`func (o *Requester) GetGrantedAudience() []string`

GetGrantedAudience returns the GrantedAudience field if non-nil, zero value
otherwise.

### GetGrantedAudienceOk

`func (o *Requester) GetGrantedAudienceOk() (*[]string, bool)`

GetGrantedAudienceOk returns a tuple with the GrantedAudience field if it's
non-nil, zero value otherwise and a boolean to check if the value has been set.

### SetGrantedAudience

`func (o *Requester) SetGrantedAudience(v []string)`

SetGrantedAudience sets GrantedAudience field to given value.

### HasGrantedAudience

`func (o *Requester) HasGrantedAudience() bool`

HasGrantedAudience returns a boolean if a field has been set.

### GetGrantedScopes

`func (o *Requester) GetGrantedScopes() []string`

GetGrantedScopes returns the GrantedScopes field if non-nil, zero value
otherwise.

### GetGrantedScopesOk

`func (o *Requester) GetGrantedScopesOk() (*[]string, bool)`

GetGrantedScopesOk returns a tuple with the GrantedScopes field if it's non-nil,
zero value otherwise and a boolean to check if the value has been set.

### SetGrantedScopes

`func (o *Requester) SetGrantedScopes(v []string)`

SetGrantedScopes sets GrantedScopes field to given value.

### HasGrantedScopes

`func (o *Requester) HasGrantedScopes() bool`

HasGrantedScopes returns a boolean if a field has been set.

[[Back to Model list]](../README.md#documentation-for-models)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to README]](../README.md)
