# RefreshTokenHookRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ClientId** | Pointer to **string** | ClientID is the identifier of the OAuth 2.0 client. | [optional] 
**GrantedAudience** | Pointer to **[]string** | GrantedAudience is the list of audiences granted to the OAuth 2.0 client. | [optional] 
**GrantedScopes** | Pointer to **[]string** | GrantedScopes is the list of scopes granted to the OAuth 2.0 client. | [optional] 
**Subject** | Pointer to **string** | Subject is the identifier of the authenticated end-user. | [optional] 

## Methods

### NewRefreshTokenHookRequest

`func NewRefreshTokenHookRequest() *RefreshTokenHookRequest`

NewRefreshTokenHookRequest instantiates a new RefreshTokenHookRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRefreshTokenHookRequestWithDefaults

`func NewRefreshTokenHookRequestWithDefaults() *RefreshTokenHookRequest`

NewRefreshTokenHookRequestWithDefaults instantiates a new RefreshTokenHookRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetClientId

`func (o *RefreshTokenHookRequest) GetClientId() string`

GetClientId returns the ClientId field if non-nil, zero value otherwise.

### GetClientIdOk

`func (o *RefreshTokenHookRequest) GetClientIdOk() (*string, bool)`

GetClientIdOk returns a tuple with the ClientId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClientId

`func (o *RefreshTokenHookRequest) SetClientId(v string)`

SetClientId sets ClientId field to given value.

### HasClientId

`func (o *RefreshTokenHookRequest) HasClientId() bool`

HasClientId returns a boolean if a field has been set.

### GetGrantedAudience

`func (o *RefreshTokenHookRequest) GetGrantedAudience() []string`

GetGrantedAudience returns the GrantedAudience field if non-nil, zero value otherwise.

### GetGrantedAudienceOk

`func (o *RefreshTokenHookRequest) GetGrantedAudienceOk() (*[]string, bool)`

GetGrantedAudienceOk returns a tuple with the GrantedAudience field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGrantedAudience

`func (o *RefreshTokenHookRequest) SetGrantedAudience(v []string)`

SetGrantedAudience sets GrantedAudience field to given value.

### HasGrantedAudience

`func (o *RefreshTokenHookRequest) HasGrantedAudience() bool`

HasGrantedAudience returns a boolean if a field has been set.

### GetGrantedScopes

`func (o *RefreshTokenHookRequest) GetGrantedScopes() []string`

GetGrantedScopes returns the GrantedScopes field if non-nil, zero value otherwise.

### GetGrantedScopesOk

`func (o *RefreshTokenHookRequest) GetGrantedScopesOk() (*[]string, bool)`

GetGrantedScopesOk returns a tuple with the GrantedScopes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGrantedScopes

`func (o *RefreshTokenHookRequest) SetGrantedScopes(v []string)`

SetGrantedScopes sets GrantedScopes field to given value.

### HasGrantedScopes

`func (o *RefreshTokenHookRequest) HasGrantedScopes() bool`

HasGrantedScopes returns a boolean if a field has been set.

### GetSubject

`func (o *RefreshTokenHookRequest) GetSubject() string`

GetSubject returns the Subject field if non-nil, zero value otherwise.

### GetSubjectOk

`func (o *RefreshTokenHookRequest) GetSubjectOk() (*string, bool)`

GetSubjectOk returns a tuple with the Subject field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubject

`func (o *RefreshTokenHookRequest) SetSubject(v string)`

SetSubject sets Subject field to given value.

### HasSubject

`func (o *RefreshTokenHookRequest) HasSubject() bool`

HasSubject returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


