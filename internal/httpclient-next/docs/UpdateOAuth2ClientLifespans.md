# UpdateOAuth2ClientLifespans

## Properties

| Name                                           | Type                  | Description                                                       | Notes      |
| ---------------------------------------------- | --------------------- | ----------------------------------------------------------------- | ---------- |
| **AuthorizationCodeGrantAccessTokenLifespan**  | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] |
| **AuthorizationCodeGrantIdTokenLifespan**      | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] |
| **AuthorizationCodeGrantRefreshTokenLifespan** | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] |
| **ClientCredentialsGrantAccessTokenLifespan**  | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] |
| **ImplicitGrantAccessTokenLifespan**           | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] |
| **ImplicitGrantIdTokenLifespan**               | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] |
| **JwtBearerGrantAccessTokenLifespan**          | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] |
| **PasswordGrantAccessTokenLifespan**           | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] |
| **PasswordGrantRefreshTokenLifespan**          | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] |
| **RefreshTokenGrantAccessTokenLifespan**       | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] |
| **RefreshTokenGrantIdTokenLifespan**           | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] |
| **RefreshTokenGrantRefreshTokenLifespan**      | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] |

## Methods

### NewUpdateOAuth2ClientLifespans

`func NewUpdateOAuth2ClientLifespans() *UpdateOAuth2ClientLifespans`

NewUpdateOAuth2ClientLifespans instantiates a new UpdateOAuth2ClientLifespans
object This constructor will assign default values to properties that have it
defined, and makes sure properties required by API are set, but the set of
arguments will change when the set of required properties is changed

### NewUpdateOAuth2ClientLifespansWithDefaults

`func NewUpdateOAuth2ClientLifespansWithDefaults() *UpdateOAuth2ClientLifespans`

NewUpdateOAuth2ClientLifespansWithDefaults instantiates a new
UpdateOAuth2ClientLifespans object This constructor will only assign default
values to properties that have it defined, but it doesn't guarantee that
properties required by API are set

### GetAuthorizationCodeGrantAccessTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) GetAuthorizationCodeGrantAccessTokenLifespan() string`

GetAuthorizationCodeGrantAccessTokenLifespan returns the
AuthorizationCodeGrantAccessTokenLifespan field if non-nil, zero value
otherwise.

### GetAuthorizationCodeGrantAccessTokenLifespanOk

`func (o *UpdateOAuth2ClientLifespans) GetAuthorizationCodeGrantAccessTokenLifespanOk() (*string, bool)`

GetAuthorizationCodeGrantAccessTokenLifespanOk returns a tuple with the
AuthorizationCodeGrantAccessTokenLifespan field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetAuthorizationCodeGrantAccessTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) SetAuthorizationCodeGrantAccessTokenLifespan(v string)`

SetAuthorizationCodeGrantAccessTokenLifespan sets
AuthorizationCodeGrantAccessTokenLifespan field to given value.

### HasAuthorizationCodeGrantAccessTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) HasAuthorizationCodeGrantAccessTokenLifespan() bool`

HasAuthorizationCodeGrantAccessTokenLifespan returns a boolean if a field has
been set.

### GetAuthorizationCodeGrantIdTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) GetAuthorizationCodeGrantIdTokenLifespan() string`

GetAuthorizationCodeGrantIdTokenLifespan returns the
AuthorizationCodeGrantIdTokenLifespan field if non-nil, zero value otherwise.

### GetAuthorizationCodeGrantIdTokenLifespanOk

`func (o *UpdateOAuth2ClientLifespans) GetAuthorizationCodeGrantIdTokenLifespanOk() (*string, bool)`

GetAuthorizationCodeGrantIdTokenLifespanOk returns a tuple with the
AuthorizationCodeGrantIdTokenLifespan field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetAuthorizationCodeGrantIdTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) SetAuthorizationCodeGrantIdTokenLifespan(v string)`

SetAuthorizationCodeGrantIdTokenLifespan sets
AuthorizationCodeGrantIdTokenLifespan field to given value.

### HasAuthorizationCodeGrantIdTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) HasAuthorizationCodeGrantIdTokenLifespan() bool`

HasAuthorizationCodeGrantIdTokenLifespan returns a boolean if a field has been
set.

### GetAuthorizationCodeGrantRefreshTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) GetAuthorizationCodeGrantRefreshTokenLifespan() string`

GetAuthorizationCodeGrantRefreshTokenLifespan returns the
AuthorizationCodeGrantRefreshTokenLifespan field if non-nil, zero value
otherwise.

### GetAuthorizationCodeGrantRefreshTokenLifespanOk

`func (o *UpdateOAuth2ClientLifespans) GetAuthorizationCodeGrantRefreshTokenLifespanOk() (*string, bool)`

GetAuthorizationCodeGrantRefreshTokenLifespanOk returns a tuple with the
AuthorizationCodeGrantRefreshTokenLifespan field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetAuthorizationCodeGrantRefreshTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) SetAuthorizationCodeGrantRefreshTokenLifespan(v string)`

SetAuthorizationCodeGrantRefreshTokenLifespan sets
AuthorizationCodeGrantRefreshTokenLifespan field to given value.

### HasAuthorizationCodeGrantRefreshTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) HasAuthorizationCodeGrantRefreshTokenLifespan() bool`

HasAuthorizationCodeGrantRefreshTokenLifespan returns a boolean if a field has
been set.

### GetClientCredentialsGrantAccessTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) GetClientCredentialsGrantAccessTokenLifespan() string`

GetClientCredentialsGrantAccessTokenLifespan returns the
ClientCredentialsGrantAccessTokenLifespan field if non-nil, zero value
otherwise.

### GetClientCredentialsGrantAccessTokenLifespanOk

`func (o *UpdateOAuth2ClientLifespans) GetClientCredentialsGrantAccessTokenLifespanOk() (*string, bool)`

GetClientCredentialsGrantAccessTokenLifespanOk returns a tuple with the
ClientCredentialsGrantAccessTokenLifespan field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetClientCredentialsGrantAccessTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) SetClientCredentialsGrantAccessTokenLifespan(v string)`

SetClientCredentialsGrantAccessTokenLifespan sets
ClientCredentialsGrantAccessTokenLifespan field to given value.

### HasClientCredentialsGrantAccessTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) HasClientCredentialsGrantAccessTokenLifespan() bool`

HasClientCredentialsGrantAccessTokenLifespan returns a boolean if a field has
been set.

### GetImplicitGrantAccessTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) GetImplicitGrantAccessTokenLifespan() string`

GetImplicitGrantAccessTokenLifespan returns the ImplicitGrantAccessTokenLifespan
field if non-nil, zero value otherwise.

### GetImplicitGrantAccessTokenLifespanOk

`func (o *UpdateOAuth2ClientLifespans) GetImplicitGrantAccessTokenLifespanOk() (*string, bool)`

GetImplicitGrantAccessTokenLifespanOk returns a tuple with the
ImplicitGrantAccessTokenLifespan field if it's non-nil, zero value otherwise and
a boolean to check if the value has been set.

### SetImplicitGrantAccessTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) SetImplicitGrantAccessTokenLifespan(v string)`

SetImplicitGrantAccessTokenLifespan sets ImplicitGrantAccessTokenLifespan field
to given value.

### HasImplicitGrantAccessTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) HasImplicitGrantAccessTokenLifespan() bool`

HasImplicitGrantAccessTokenLifespan returns a boolean if a field has been set.

### GetImplicitGrantIdTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) GetImplicitGrantIdTokenLifespan() string`

GetImplicitGrantIdTokenLifespan returns the ImplicitGrantIdTokenLifespan field
if non-nil, zero value otherwise.

### GetImplicitGrantIdTokenLifespanOk

`func (o *UpdateOAuth2ClientLifespans) GetImplicitGrantIdTokenLifespanOk() (*string, bool)`

GetImplicitGrantIdTokenLifespanOk returns a tuple with the
ImplicitGrantIdTokenLifespan field if it's non-nil, zero value otherwise and a
boolean to check if the value has been set.

### SetImplicitGrantIdTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) SetImplicitGrantIdTokenLifespan(v string)`

SetImplicitGrantIdTokenLifespan sets ImplicitGrantIdTokenLifespan field to given
value.

### HasImplicitGrantIdTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) HasImplicitGrantIdTokenLifespan() bool`

HasImplicitGrantIdTokenLifespan returns a boolean if a field has been set.

### GetJwtBearerGrantAccessTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) GetJwtBearerGrantAccessTokenLifespan() string`

GetJwtBearerGrantAccessTokenLifespan returns the
JwtBearerGrantAccessTokenLifespan field if non-nil, zero value otherwise.

### GetJwtBearerGrantAccessTokenLifespanOk

`func (o *UpdateOAuth2ClientLifespans) GetJwtBearerGrantAccessTokenLifespanOk() (*string, bool)`

GetJwtBearerGrantAccessTokenLifespanOk returns a tuple with the
JwtBearerGrantAccessTokenLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetJwtBearerGrantAccessTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) SetJwtBearerGrantAccessTokenLifespan(v string)`

SetJwtBearerGrantAccessTokenLifespan sets JwtBearerGrantAccessTokenLifespan
field to given value.

### HasJwtBearerGrantAccessTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) HasJwtBearerGrantAccessTokenLifespan() bool`

HasJwtBearerGrantAccessTokenLifespan returns a boolean if a field has been set.

### GetPasswordGrantAccessTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) GetPasswordGrantAccessTokenLifespan() string`

GetPasswordGrantAccessTokenLifespan returns the PasswordGrantAccessTokenLifespan
field if non-nil, zero value otherwise.

### GetPasswordGrantAccessTokenLifespanOk

`func (o *UpdateOAuth2ClientLifespans) GetPasswordGrantAccessTokenLifespanOk() (*string, bool)`

GetPasswordGrantAccessTokenLifespanOk returns a tuple with the
PasswordGrantAccessTokenLifespan field if it's non-nil, zero value otherwise and
a boolean to check if the value has been set.

### SetPasswordGrantAccessTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) SetPasswordGrantAccessTokenLifespan(v string)`

SetPasswordGrantAccessTokenLifespan sets PasswordGrantAccessTokenLifespan field
to given value.

### HasPasswordGrantAccessTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) HasPasswordGrantAccessTokenLifespan() bool`

HasPasswordGrantAccessTokenLifespan returns a boolean if a field has been set.

### GetPasswordGrantRefreshTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) GetPasswordGrantRefreshTokenLifespan() string`

GetPasswordGrantRefreshTokenLifespan returns the
PasswordGrantRefreshTokenLifespan field if non-nil, zero value otherwise.

### GetPasswordGrantRefreshTokenLifespanOk

`func (o *UpdateOAuth2ClientLifespans) GetPasswordGrantRefreshTokenLifespanOk() (*string, bool)`

GetPasswordGrantRefreshTokenLifespanOk returns a tuple with the
PasswordGrantRefreshTokenLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPasswordGrantRefreshTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) SetPasswordGrantRefreshTokenLifespan(v string)`

SetPasswordGrantRefreshTokenLifespan sets PasswordGrantRefreshTokenLifespan
field to given value.

### HasPasswordGrantRefreshTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) HasPasswordGrantRefreshTokenLifespan() bool`

HasPasswordGrantRefreshTokenLifespan returns a boolean if a field has been set.

### GetRefreshTokenGrantAccessTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) GetRefreshTokenGrantAccessTokenLifespan() string`

GetRefreshTokenGrantAccessTokenLifespan returns the
RefreshTokenGrantAccessTokenLifespan field if non-nil, zero value otherwise.

### GetRefreshTokenGrantAccessTokenLifespanOk

`func (o *UpdateOAuth2ClientLifespans) GetRefreshTokenGrantAccessTokenLifespanOk() (*string, bool)`

GetRefreshTokenGrantAccessTokenLifespanOk returns a tuple with the
RefreshTokenGrantAccessTokenLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRefreshTokenGrantAccessTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) SetRefreshTokenGrantAccessTokenLifespan(v string)`

SetRefreshTokenGrantAccessTokenLifespan sets
RefreshTokenGrantAccessTokenLifespan field to given value.

### HasRefreshTokenGrantAccessTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) HasRefreshTokenGrantAccessTokenLifespan() bool`

HasRefreshTokenGrantAccessTokenLifespan returns a boolean if a field has been
set.

### GetRefreshTokenGrantIdTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) GetRefreshTokenGrantIdTokenLifespan() string`

GetRefreshTokenGrantIdTokenLifespan returns the RefreshTokenGrantIdTokenLifespan
field if non-nil, zero value otherwise.

### GetRefreshTokenGrantIdTokenLifespanOk

`func (o *UpdateOAuth2ClientLifespans) GetRefreshTokenGrantIdTokenLifespanOk() (*string, bool)`

GetRefreshTokenGrantIdTokenLifespanOk returns a tuple with the
RefreshTokenGrantIdTokenLifespan field if it's non-nil, zero value otherwise and
a boolean to check if the value has been set.

### SetRefreshTokenGrantIdTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) SetRefreshTokenGrantIdTokenLifespan(v string)`

SetRefreshTokenGrantIdTokenLifespan sets RefreshTokenGrantIdTokenLifespan field
to given value.

### HasRefreshTokenGrantIdTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) HasRefreshTokenGrantIdTokenLifespan() bool`

HasRefreshTokenGrantIdTokenLifespan returns a boolean if a field has been set.

### GetRefreshTokenGrantRefreshTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) GetRefreshTokenGrantRefreshTokenLifespan() string`

GetRefreshTokenGrantRefreshTokenLifespan returns the
RefreshTokenGrantRefreshTokenLifespan field if non-nil, zero value otherwise.

### GetRefreshTokenGrantRefreshTokenLifespanOk

`func (o *UpdateOAuth2ClientLifespans) GetRefreshTokenGrantRefreshTokenLifespanOk() (*string, bool)`

GetRefreshTokenGrantRefreshTokenLifespanOk returns a tuple with the
RefreshTokenGrantRefreshTokenLifespan field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetRefreshTokenGrantRefreshTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) SetRefreshTokenGrantRefreshTokenLifespan(v string)`

SetRefreshTokenGrantRefreshTokenLifespan sets
RefreshTokenGrantRefreshTokenLifespan field to given value.

### HasRefreshTokenGrantRefreshTokenLifespan

`func (o *UpdateOAuth2ClientLifespans) HasRefreshTokenGrantRefreshTokenLifespan() bool`

HasRefreshTokenGrantRefreshTokenLifespan returns a boolean if a field has been
set.

[[Back to Model list]](../README.md#documentation-for-models)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to README]](../README.md)
