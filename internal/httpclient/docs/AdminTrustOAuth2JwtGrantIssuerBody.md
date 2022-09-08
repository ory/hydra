# AdminTrustOAuth2JwtGrantIssuerBody

## Properties

| Name                | Type                            | Description                                                                                                                                            | Notes      |
| ------------------- | ------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------ | ---------- |
| **AllowAnySubject** | Pointer to **bool**             | The \&quot;allow_any_subject\&quot; indicates that the issuer is allowed to have any principal as the subject of the JWT.                              | [optional] |
| **ExpiresAt**       | **time.Time**                   | The \&quot;expires_at\&quot; indicates, when grant will expire, so we will reject assertion from \&quot;issuer\&quot; targeting \&quot;subject\&quot;. |
| **Issuer**          | **string**                      | The \&quot;issuer\&quot; identifies the principal that issued the JWT assertion (same as \&quot;iss\&quot; claim in JWT).                              |
| **Jwk**             | [**JsonWebKey**](JsonWebKey.md) |                                                                                                                                                        |
| **Scope**           | **[]string**                    | The \&quot;scope\&quot; contains list of scope values (as described in Section 3.3 of OAuth 2.0 [RFC6749])                                             |
| **Subject**         | Pointer to **string**           | The \&quot;subject\&quot; identifies the principal that is the subject of the JWT.                                                                     | [optional] |

## Methods

### NewAdminTrustOAuth2JwtGrantIssuerBody

`func NewAdminTrustOAuth2JwtGrantIssuerBody(expiresAt time.Time, issuer string, jwk JsonWebKey, scope []string, ) *AdminTrustOAuth2JwtGrantIssuerBody`

NewAdminTrustOAuth2JwtGrantIssuerBody instantiates a new
AdminTrustOAuth2JwtGrantIssuerBody object This constructor will assign default
values to properties that have it defined, and makes sure properties required by
API are set, but the set of arguments will change when the set of required
properties is changed

### NewAdminTrustOAuth2JwtGrantIssuerBodyWithDefaults

`func NewAdminTrustOAuth2JwtGrantIssuerBodyWithDefaults() *AdminTrustOAuth2JwtGrantIssuerBody`

NewAdminTrustOAuth2JwtGrantIssuerBodyWithDefaults instantiates a new
AdminTrustOAuth2JwtGrantIssuerBody object This constructor will only assign
default values to properties that have it defined, but it doesn't guarantee that
properties required by API are set

### GetAllowAnySubject

`func (o *AdminTrustOAuth2JwtGrantIssuerBody) GetAllowAnySubject() bool`

GetAllowAnySubject returns the AllowAnySubject field if non-nil, zero value
otherwise.

### GetAllowAnySubjectOk

`func (o *AdminTrustOAuth2JwtGrantIssuerBody) GetAllowAnySubjectOk() (*bool, bool)`

GetAllowAnySubjectOk returns a tuple with the AllowAnySubject field if it's
non-nil, zero value otherwise and a boolean to check if the value has been set.

### SetAllowAnySubject

`func (o *AdminTrustOAuth2JwtGrantIssuerBody) SetAllowAnySubject(v bool)`

SetAllowAnySubject sets AllowAnySubject field to given value.

### HasAllowAnySubject

`func (o *AdminTrustOAuth2JwtGrantIssuerBody) HasAllowAnySubject() bool`

HasAllowAnySubject returns a boolean if a field has been set.

### GetExpiresAt

`func (o *AdminTrustOAuth2JwtGrantIssuerBody) GetExpiresAt() time.Time`

GetExpiresAt returns the ExpiresAt field if non-nil, zero value otherwise.

### GetExpiresAtOk

`func (o *AdminTrustOAuth2JwtGrantIssuerBody) GetExpiresAtOk() (*time.Time, bool)`

GetExpiresAtOk returns a tuple with the ExpiresAt field if it's non-nil, zero
value otherwise and a boolean to check if the value has been set.

### SetExpiresAt

`func (o *AdminTrustOAuth2JwtGrantIssuerBody) SetExpiresAt(v time.Time)`

SetExpiresAt sets ExpiresAt field to given value.

### GetIssuer

`func (o *AdminTrustOAuth2JwtGrantIssuerBody) GetIssuer() string`

GetIssuer returns the Issuer field if non-nil, zero value otherwise.

### GetIssuerOk

`func (o *AdminTrustOAuth2JwtGrantIssuerBody) GetIssuerOk() (*string, bool)`

GetIssuerOk returns a tuple with the Issuer field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetIssuer

`func (o *AdminTrustOAuth2JwtGrantIssuerBody) SetIssuer(v string)`

SetIssuer sets Issuer field to given value.

### GetJwk

`func (o *AdminTrustOAuth2JwtGrantIssuerBody) GetJwk() JsonWebKey`

GetJwk returns the Jwk field if non-nil, zero value otherwise.

### GetJwkOk

`func (o *AdminTrustOAuth2JwtGrantIssuerBody) GetJwkOk() (*JsonWebKey, bool)`

GetJwkOk returns a tuple with the Jwk field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetJwk

`func (o *AdminTrustOAuth2JwtGrantIssuerBody) SetJwk(v JsonWebKey)`

SetJwk sets Jwk field to given value.

### GetScope

`func (o *AdminTrustOAuth2JwtGrantIssuerBody) GetScope() []string`

GetScope returns the Scope field if non-nil, zero value otherwise.

### GetScopeOk

`func (o *AdminTrustOAuth2JwtGrantIssuerBody) GetScopeOk() (*[]string, bool)`

GetScopeOk returns a tuple with the Scope field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetScope

`func (o *AdminTrustOAuth2JwtGrantIssuerBody) SetScope(v []string)`

SetScope sets Scope field to given value.

### GetSubject

`func (o *AdminTrustOAuth2JwtGrantIssuerBody) GetSubject() string`

GetSubject returns the Subject field if non-nil, zero value otherwise.

### GetSubjectOk

`func (o *AdminTrustOAuth2JwtGrantIssuerBody) GetSubjectOk() (*string, bool)`

GetSubjectOk returns a tuple with the Subject field if it's non-nil, zero value
otherwise and a boolean to check if the value has been set.

### SetSubject

`func (o *AdminTrustOAuth2JwtGrantIssuerBody) SetSubject(v string)`

SetSubject sets Subject field to given value.

### HasSubject

`func (o *AdminTrustOAuth2JwtGrantIssuerBody) HasSubject() bool`

HasSubject returns a boolean if a field has been set.

[[Back to Model list]](../README.md#documentation-for-models)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to README]](../README.md)
