# OAuth2TokenResponse

## Properties

| Name             | Type                  | Description                                                                                                                                                                            | Notes      |
| ---------------- | --------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------- |
| **AccessToken**  | Pointer to **string** | The access token issued by the authorization server.                                                                                                                                   | [optional] |
| **ExpiresIn**    | Pointer to **int64**  | The lifetime in seconds of the access token. For example, the value \&quot;3600\&quot; denotes that the access token will expire in one hour from the time the response was generated. | [optional] |
| **IdToken**      | Pointer to **int64**  | To retrieve a refresh token request the id_token scope.                                                                                                                                | [optional] |
| **RefreshToken** | Pointer to **string** | The refresh token, which can be used to obtain new access tokens. To retrieve it add the scope \&quot;offline\&quot; to your access token request.                                     | [optional] |
| **Scope**        | Pointer to **int64**  | The scope of the access token                                                                                                                                                          | [optional] |
| **TokenType**    | Pointer to **string** | The type of the token issued                                                                                                                                                           | [optional] |

## Methods

### NewOAuth2TokenResponse

`func NewOAuth2TokenResponse() *OAuth2TokenResponse`

NewOAuth2TokenResponse instantiates a new OAuth2TokenResponse object This
constructor will assign default values to properties that have it defined, and
makes sure properties required by API are set, but the set of arguments will
change when the set of required properties is changed

### NewOAuth2TokenResponseWithDefaults

`func NewOAuth2TokenResponseWithDefaults() *OAuth2TokenResponse`

NewOAuth2TokenResponseWithDefaults instantiates a new OAuth2TokenResponse object
This constructor will only assign default values to properties that have it
defined, but it doesn't guarantee that properties required by API are set

### GetAccessToken

`func (o *OAuth2TokenResponse) GetAccessToken() string`

GetAccessToken returns the AccessToken field if non-nil, zero value otherwise.

### GetAccessTokenOk

`func (o *OAuth2TokenResponse) GetAccessTokenOk() (*string, bool)`

GetAccessTokenOk returns a tuple with the AccessToken field if it's non-nil,
zero value otherwise and a boolean to check if the value has been set.

### SetAccessToken

`func (o *OAuth2TokenResponse) SetAccessToken(v string)`

SetAccessToken sets AccessToken field to given value.

### HasAccessToken

`func (o *OAuth2TokenResponse) HasAccessToken() bool`

HasAccessToken returns a boolean if a field has been set.

### GetExpiresIn

`func (o *OAuth2TokenResponse) GetExpiresIn() int64`

GetExpiresIn returns the ExpiresIn field if non-nil, zero value otherwise.

### GetExpiresInOk

`func (o *OAuth2TokenResponse) GetExpiresInOk() (*int64, bool)`

GetExpiresInOk returns a tuple with the ExpiresIn field if it's non-nil, zero
value otherwise and a boolean to check if the value has been set.

### SetExpiresIn

`func (o *OAuth2TokenResponse) SetExpiresIn(v int64)`

SetExpiresIn sets ExpiresIn field to given value.

### HasExpiresIn

`func (o *OAuth2TokenResponse) HasExpiresIn() bool`

HasExpiresIn returns a boolean if a field has been set.

### GetIdToken

`func (o *OAuth2TokenResponse) GetIdToken() int64`

GetIdToken returns the IdToken field if non-nil, zero value otherwise.

### GetIdTokenOk

`func (o *OAuth2TokenResponse) GetIdTokenOk() (*int64, bool)`

GetIdTokenOk returns a tuple with the IdToken field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetIdToken

`func (o *OAuth2TokenResponse) SetIdToken(v int64)`

SetIdToken sets IdToken field to given value.

### HasIdToken

`func (o *OAuth2TokenResponse) HasIdToken() bool`

HasIdToken returns a boolean if a field has been set.

### GetRefreshToken

`func (o *OAuth2TokenResponse) GetRefreshToken() string`

GetRefreshToken returns the RefreshToken field if non-nil, zero value otherwise.

### GetRefreshTokenOk

`func (o *OAuth2TokenResponse) GetRefreshTokenOk() (*string, bool)`

GetRefreshTokenOk returns a tuple with the RefreshToken field if it's non-nil,
zero value otherwise and a boolean to check if the value has been set.

### SetRefreshToken

`func (o *OAuth2TokenResponse) SetRefreshToken(v string)`

SetRefreshToken sets RefreshToken field to given value.

### HasRefreshToken

`func (o *OAuth2TokenResponse) HasRefreshToken() bool`

HasRefreshToken returns a boolean if a field has been set.

### GetScope

`func (o *OAuth2TokenResponse) GetScope() int64`

GetScope returns the Scope field if non-nil, zero value otherwise.

### GetScopeOk

`func (o *OAuth2TokenResponse) GetScopeOk() (*int64, bool)`

GetScopeOk returns a tuple with the Scope field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetScope

`func (o *OAuth2TokenResponse) SetScope(v int64)`

SetScope sets Scope field to given value.

### HasScope

`func (o *OAuth2TokenResponse) HasScope() bool`

HasScope returns a boolean if a field has been set.

### GetTokenType

`func (o *OAuth2TokenResponse) GetTokenType() string`

GetTokenType returns the TokenType field if non-nil, zero value otherwise.

### GetTokenTypeOk

`func (o *OAuth2TokenResponse) GetTokenTypeOk() (*string, bool)`

GetTokenTypeOk returns a tuple with the TokenType field if it's non-nil, zero
value otherwise and a boolean to check if the value has been set.

### SetTokenType

`func (o *OAuth2TokenResponse) SetTokenType(v string)`

SetTokenType sets TokenType field to given value.

### HasTokenType

`func (o *OAuth2TokenResponse) HasTokenType() bool`

HasTokenType returns a boolean if a field has been set.

[[Back to Model list]](../README.md#documentation-for-models)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to README]](../README.md)
