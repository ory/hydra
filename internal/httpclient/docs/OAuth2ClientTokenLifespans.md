# OAuth2ClientTokenLifespans

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AuthorizationCodeGrantAccessTokenLifespan** | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] 
**AuthorizationCodeGrantIdTokenLifespan** | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] 
**AuthorizationCodeGrantRefreshTokenLifespan** | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] 
**ClientCredentialsGrantAccessTokenLifespan** | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] 
**DeviceAuthorizationGrantAccessTokenLifespan** | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] 
**DeviceAuthorizationGrantIdTokenLifespan** | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] 
**DeviceAuthorizationGrantRefreshTokenLifespan** | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] 
**ImplicitGrantAccessTokenLifespan** | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] 
**ImplicitGrantIdTokenLifespan** | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] 
**JwtBearerGrantAccessTokenLifespan** | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] 
**RefreshTokenGrantAccessTokenLifespan** | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] 
**RefreshTokenGrantIdTokenLifespan** | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] 
**RefreshTokenGrantRefreshTokenLifespan** | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] 

## Methods

### NewOAuth2ClientTokenLifespans

`func NewOAuth2ClientTokenLifespans() *OAuth2ClientTokenLifespans`

NewOAuth2ClientTokenLifespans instantiates a new OAuth2ClientTokenLifespans object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOAuth2ClientTokenLifespansWithDefaults

`func NewOAuth2ClientTokenLifespansWithDefaults() *OAuth2ClientTokenLifespans`

NewOAuth2ClientTokenLifespansWithDefaults instantiates a new OAuth2ClientTokenLifespans object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAuthorizationCodeGrantAccessTokenLifespan

`func (o *OAuth2ClientTokenLifespans) GetAuthorizationCodeGrantAccessTokenLifespan() string`

GetAuthorizationCodeGrantAccessTokenLifespan returns the AuthorizationCodeGrantAccessTokenLifespan field if non-nil, zero value otherwise.

### GetAuthorizationCodeGrantAccessTokenLifespanOk

`func (o *OAuth2ClientTokenLifespans) GetAuthorizationCodeGrantAccessTokenLifespanOk() (*string, bool)`

GetAuthorizationCodeGrantAccessTokenLifespanOk returns a tuple with the AuthorizationCodeGrantAccessTokenLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthorizationCodeGrantAccessTokenLifespan

`func (o *OAuth2ClientTokenLifespans) SetAuthorizationCodeGrantAccessTokenLifespan(v string)`

SetAuthorizationCodeGrantAccessTokenLifespan sets AuthorizationCodeGrantAccessTokenLifespan field to given value.

### HasAuthorizationCodeGrantAccessTokenLifespan

`func (o *OAuth2ClientTokenLifespans) HasAuthorizationCodeGrantAccessTokenLifespan() bool`

HasAuthorizationCodeGrantAccessTokenLifespan returns a boolean if a field has been set.

### GetAuthorizationCodeGrantIdTokenLifespan

`func (o *OAuth2ClientTokenLifespans) GetAuthorizationCodeGrantIdTokenLifespan() string`

GetAuthorizationCodeGrantIdTokenLifespan returns the AuthorizationCodeGrantIdTokenLifespan field if non-nil, zero value otherwise.

### GetAuthorizationCodeGrantIdTokenLifespanOk

`func (o *OAuth2ClientTokenLifespans) GetAuthorizationCodeGrantIdTokenLifespanOk() (*string, bool)`

GetAuthorizationCodeGrantIdTokenLifespanOk returns a tuple with the AuthorizationCodeGrantIdTokenLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthorizationCodeGrantIdTokenLifespan

`func (o *OAuth2ClientTokenLifespans) SetAuthorizationCodeGrantIdTokenLifespan(v string)`

SetAuthorizationCodeGrantIdTokenLifespan sets AuthorizationCodeGrantIdTokenLifespan field to given value.

### HasAuthorizationCodeGrantIdTokenLifespan

`func (o *OAuth2ClientTokenLifespans) HasAuthorizationCodeGrantIdTokenLifespan() bool`

HasAuthorizationCodeGrantIdTokenLifespan returns a boolean if a field has been set.

### GetAuthorizationCodeGrantRefreshTokenLifespan

`func (o *OAuth2ClientTokenLifespans) GetAuthorizationCodeGrantRefreshTokenLifespan() string`

GetAuthorizationCodeGrantRefreshTokenLifespan returns the AuthorizationCodeGrantRefreshTokenLifespan field if non-nil, zero value otherwise.

### GetAuthorizationCodeGrantRefreshTokenLifespanOk

`func (o *OAuth2ClientTokenLifespans) GetAuthorizationCodeGrantRefreshTokenLifespanOk() (*string, bool)`

GetAuthorizationCodeGrantRefreshTokenLifespanOk returns a tuple with the AuthorizationCodeGrantRefreshTokenLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthorizationCodeGrantRefreshTokenLifespan

`func (o *OAuth2ClientTokenLifespans) SetAuthorizationCodeGrantRefreshTokenLifespan(v string)`

SetAuthorizationCodeGrantRefreshTokenLifespan sets AuthorizationCodeGrantRefreshTokenLifespan field to given value.

### HasAuthorizationCodeGrantRefreshTokenLifespan

`func (o *OAuth2ClientTokenLifespans) HasAuthorizationCodeGrantRefreshTokenLifespan() bool`

HasAuthorizationCodeGrantRefreshTokenLifespan returns a boolean if a field has been set.

### GetClientCredentialsGrantAccessTokenLifespan

`func (o *OAuth2ClientTokenLifespans) GetClientCredentialsGrantAccessTokenLifespan() string`

GetClientCredentialsGrantAccessTokenLifespan returns the ClientCredentialsGrantAccessTokenLifespan field if non-nil, zero value otherwise.

### GetClientCredentialsGrantAccessTokenLifespanOk

`func (o *OAuth2ClientTokenLifespans) GetClientCredentialsGrantAccessTokenLifespanOk() (*string, bool)`

GetClientCredentialsGrantAccessTokenLifespanOk returns a tuple with the ClientCredentialsGrantAccessTokenLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClientCredentialsGrantAccessTokenLifespan

`func (o *OAuth2ClientTokenLifespans) SetClientCredentialsGrantAccessTokenLifespan(v string)`

SetClientCredentialsGrantAccessTokenLifespan sets ClientCredentialsGrantAccessTokenLifespan field to given value.

### HasClientCredentialsGrantAccessTokenLifespan

`func (o *OAuth2ClientTokenLifespans) HasClientCredentialsGrantAccessTokenLifespan() bool`

HasClientCredentialsGrantAccessTokenLifespan returns a boolean if a field has been set.

### GetDeviceAuthorizationGrantAccessTokenLifespan

`func (o *OAuth2ClientTokenLifespans) GetDeviceAuthorizationGrantAccessTokenLifespan() string`

GetDeviceAuthorizationGrantAccessTokenLifespan returns the DeviceAuthorizationGrantAccessTokenLifespan field if non-nil, zero value otherwise.

### GetDeviceAuthorizationGrantAccessTokenLifespanOk

`func (o *OAuth2ClientTokenLifespans) GetDeviceAuthorizationGrantAccessTokenLifespanOk() (*string, bool)`

GetDeviceAuthorizationGrantAccessTokenLifespanOk returns a tuple with the DeviceAuthorizationGrantAccessTokenLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceAuthorizationGrantAccessTokenLifespan

`func (o *OAuth2ClientTokenLifespans) SetDeviceAuthorizationGrantAccessTokenLifespan(v string)`

SetDeviceAuthorizationGrantAccessTokenLifespan sets DeviceAuthorizationGrantAccessTokenLifespan field to given value.

### HasDeviceAuthorizationGrantAccessTokenLifespan

`func (o *OAuth2ClientTokenLifespans) HasDeviceAuthorizationGrantAccessTokenLifespan() bool`

HasDeviceAuthorizationGrantAccessTokenLifespan returns a boolean if a field has been set.

### GetDeviceAuthorizationGrantIdTokenLifespan

`func (o *OAuth2ClientTokenLifespans) GetDeviceAuthorizationGrantIdTokenLifespan() string`

GetDeviceAuthorizationGrantIdTokenLifespan returns the DeviceAuthorizationGrantIdTokenLifespan field if non-nil, zero value otherwise.

### GetDeviceAuthorizationGrantIdTokenLifespanOk

`func (o *OAuth2ClientTokenLifespans) GetDeviceAuthorizationGrantIdTokenLifespanOk() (*string, bool)`

GetDeviceAuthorizationGrantIdTokenLifespanOk returns a tuple with the DeviceAuthorizationGrantIdTokenLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceAuthorizationGrantIdTokenLifespan

`func (o *OAuth2ClientTokenLifespans) SetDeviceAuthorizationGrantIdTokenLifespan(v string)`

SetDeviceAuthorizationGrantIdTokenLifespan sets DeviceAuthorizationGrantIdTokenLifespan field to given value.

### HasDeviceAuthorizationGrantIdTokenLifespan

`func (o *OAuth2ClientTokenLifespans) HasDeviceAuthorizationGrantIdTokenLifespan() bool`

HasDeviceAuthorizationGrantIdTokenLifespan returns a boolean if a field has been set.

### GetDeviceAuthorizationGrantRefreshTokenLifespan

`func (o *OAuth2ClientTokenLifespans) GetDeviceAuthorizationGrantRefreshTokenLifespan() string`

GetDeviceAuthorizationGrantRefreshTokenLifespan returns the DeviceAuthorizationGrantRefreshTokenLifespan field if non-nil, zero value otherwise.

### GetDeviceAuthorizationGrantRefreshTokenLifespanOk

`func (o *OAuth2ClientTokenLifespans) GetDeviceAuthorizationGrantRefreshTokenLifespanOk() (*string, bool)`

GetDeviceAuthorizationGrantRefreshTokenLifespanOk returns a tuple with the DeviceAuthorizationGrantRefreshTokenLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceAuthorizationGrantRefreshTokenLifespan

`func (o *OAuth2ClientTokenLifespans) SetDeviceAuthorizationGrantRefreshTokenLifespan(v string)`

SetDeviceAuthorizationGrantRefreshTokenLifespan sets DeviceAuthorizationGrantRefreshTokenLifespan field to given value.

### HasDeviceAuthorizationGrantRefreshTokenLifespan

`func (o *OAuth2ClientTokenLifespans) HasDeviceAuthorizationGrantRefreshTokenLifespan() bool`

HasDeviceAuthorizationGrantRefreshTokenLifespan returns a boolean if a field has been set.

### GetImplicitGrantAccessTokenLifespan

`func (o *OAuth2ClientTokenLifespans) GetImplicitGrantAccessTokenLifespan() string`

GetImplicitGrantAccessTokenLifespan returns the ImplicitGrantAccessTokenLifespan field if non-nil, zero value otherwise.

### GetImplicitGrantAccessTokenLifespanOk

`func (o *OAuth2ClientTokenLifespans) GetImplicitGrantAccessTokenLifespanOk() (*string, bool)`

GetImplicitGrantAccessTokenLifespanOk returns a tuple with the ImplicitGrantAccessTokenLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImplicitGrantAccessTokenLifespan

`func (o *OAuth2ClientTokenLifespans) SetImplicitGrantAccessTokenLifespan(v string)`

SetImplicitGrantAccessTokenLifespan sets ImplicitGrantAccessTokenLifespan field to given value.

### HasImplicitGrantAccessTokenLifespan

`func (o *OAuth2ClientTokenLifespans) HasImplicitGrantAccessTokenLifespan() bool`

HasImplicitGrantAccessTokenLifespan returns a boolean if a field has been set.

### GetImplicitGrantIdTokenLifespan

`func (o *OAuth2ClientTokenLifespans) GetImplicitGrantIdTokenLifespan() string`

GetImplicitGrantIdTokenLifespan returns the ImplicitGrantIdTokenLifespan field if non-nil, zero value otherwise.

### GetImplicitGrantIdTokenLifespanOk

`func (o *OAuth2ClientTokenLifespans) GetImplicitGrantIdTokenLifespanOk() (*string, bool)`

GetImplicitGrantIdTokenLifespanOk returns a tuple with the ImplicitGrantIdTokenLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImplicitGrantIdTokenLifespan

`func (o *OAuth2ClientTokenLifespans) SetImplicitGrantIdTokenLifespan(v string)`

SetImplicitGrantIdTokenLifespan sets ImplicitGrantIdTokenLifespan field to given value.

### HasImplicitGrantIdTokenLifespan

`func (o *OAuth2ClientTokenLifespans) HasImplicitGrantIdTokenLifespan() bool`

HasImplicitGrantIdTokenLifespan returns a boolean if a field has been set.

### GetJwtBearerGrantAccessTokenLifespan

`func (o *OAuth2ClientTokenLifespans) GetJwtBearerGrantAccessTokenLifespan() string`

GetJwtBearerGrantAccessTokenLifespan returns the JwtBearerGrantAccessTokenLifespan field if non-nil, zero value otherwise.

### GetJwtBearerGrantAccessTokenLifespanOk

`func (o *OAuth2ClientTokenLifespans) GetJwtBearerGrantAccessTokenLifespanOk() (*string, bool)`

GetJwtBearerGrantAccessTokenLifespanOk returns a tuple with the JwtBearerGrantAccessTokenLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetJwtBearerGrantAccessTokenLifespan

`func (o *OAuth2ClientTokenLifespans) SetJwtBearerGrantAccessTokenLifespan(v string)`

SetJwtBearerGrantAccessTokenLifespan sets JwtBearerGrantAccessTokenLifespan field to given value.

### HasJwtBearerGrantAccessTokenLifespan

`func (o *OAuth2ClientTokenLifespans) HasJwtBearerGrantAccessTokenLifespan() bool`

HasJwtBearerGrantAccessTokenLifespan returns a boolean if a field has been set.

### GetRefreshTokenGrantAccessTokenLifespan

`func (o *OAuth2ClientTokenLifespans) GetRefreshTokenGrantAccessTokenLifespan() string`

GetRefreshTokenGrantAccessTokenLifespan returns the RefreshTokenGrantAccessTokenLifespan field if non-nil, zero value otherwise.

### GetRefreshTokenGrantAccessTokenLifespanOk

`func (o *OAuth2ClientTokenLifespans) GetRefreshTokenGrantAccessTokenLifespanOk() (*string, bool)`

GetRefreshTokenGrantAccessTokenLifespanOk returns a tuple with the RefreshTokenGrantAccessTokenLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRefreshTokenGrantAccessTokenLifespan

`func (o *OAuth2ClientTokenLifespans) SetRefreshTokenGrantAccessTokenLifespan(v string)`

SetRefreshTokenGrantAccessTokenLifespan sets RefreshTokenGrantAccessTokenLifespan field to given value.

### HasRefreshTokenGrantAccessTokenLifespan

`func (o *OAuth2ClientTokenLifespans) HasRefreshTokenGrantAccessTokenLifespan() bool`

HasRefreshTokenGrantAccessTokenLifespan returns a boolean if a field has been set.

### GetRefreshTokenGrantIdTokenLifespan

`func (o *OAuth2ClientTokenLifespans) GetRefreshTokenGrantIdTokenLifespan() string`

GetRefreshTokenGrantIdTokenLifespan returns the RefreshTokenGrantIdTokenLifespan field if non-nil, zero value otherwise.

### GetRefreshTokenGrantIdTokenLifespanOk

`func (o *OAuth2ClientTokenLifespans) GetRefreshTokenGrantIdTokenLifespanOk() (*string, bool)`

GetRefreshTokenGrantIdTokenLifespanOk returns a tuple with the RefreshTokenGrantIdTokenLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRefreshTokenGrantIdTokenLifespan

`func (o *OAuth2ClientTokenLifespans) SetRefreshTokenGrantIdTokenLifespan(v string)`

SetRefreshTokenGrantIdTokenLifespan sets RefreshTokenGrantIdTokenLifespan field to given value.

### HasRefreshTokenGrantIdTokenLifespan

`func (o *OAuth2ClientTokenLifespans) HasRefreshTokenGrantIdTokenLifespan() bool`

HasRefreshTokenGrantIdTokenLifespan returns a boolean if a field has been set.

### GetRefreshTokenGrantRefreshTokenLifespan

`func (o *OAuth2ClientTokenLifespans) GetRefreshTokenGrantRefreshTokenLifespan() string`

GetRefreshTokenGrantRefreshTokenLifespan returns the RefreshTokenGrantRefreshTokenLifespan field if non-nil, zero value otherwise.

### GetRefreshTokenGrantRefreshTokenLifespanOk

`func (o *OAuth2ClientTokenLifespans) GetRefreshTokenGrantRefreshTokenLifespanOk() (*string, bool)`

GetRefreshTokenGrantRefreshTokenLifespanOk returns a tuple with the RefreshTokenGrantRefreshTokenLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRefreshTokenGrantRefreshTokenLifespan

`func (o *OAuth2ClientTokenLifespans) SetRefreshTokenGrantRefreshTokenLifespan(v string)`

SetRefreshTokenGrantRefreshTokenLifespan sets RefreshTokenGrantRefreshTokenLifespan field to given value.

### HasRefreshTokenGrantRefreshTokenLifespan

`func (o *OAuth2ClientTokenLifespans) HasRefreshTokenGrantRefreshTokenLifespan() bool`

HasRefreshTokenGrantRefreshTokenLifespan returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


