# OAuth2Client

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AccessTokenStrategy** | Pointer to **string** | OAuth 2.0 Access Token Strategy  AccessTokenStrategy is the strategy used to generate access tokens. Valid options are &#x60;jwt&#x60; and &#x60;opaque&#x60;. &#x60;jwt&#x60; is a bad idea, see https://www.ory.sh/docs/oauth2-oidc/jwt-access-token Setting the strategy here overrides the global setting in &#x60;strategies.access_token&#x60;. | [optional] 
**AllowedCorsOrigins** | Pointer to **[]string** |  | [optional] 
**Audience** | Pointer to **[]string** |  | [optional] 
**AuthorizationCodeGrantAccessTokenLifespan** | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] 
**AuthorizationCodeGrantIdTokenLifespan** | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] 
**AuthorizationCodeGrantRefreshTokenLifespan** | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] 
**BackchannelLogoutSessionRequired** | Pointer to **bool** | OpenID Connect Back-Channel Logout Session Required  Boolean value specifying whether the RP requires that a sid (session ID) Claim be included in the Logout Token to identify the RP session with the OP when the backchannel_logout_uri is used. If omitted, the default value is false. | [optional] 
**BackchannelLogoutUri** | Pointer to **string** | OpenID Connect Back-Channel Logout URI  RP URL that will cause the RP to log itself out when sent a Logout Token by the OP. | [optional] 
**ClientCredentialsGrantAccessTokenLifespan** | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] 
**ClientId** | Pointer to **string** | OAuth 2.0 Client ID  The ID is immutable. If no ID is provided, a UUID4 will be generated. | [optional] 
**ClientName** | Pointer to **string** | OAuth 2.0 Client Name  The human-readable name of the client to be presented to the end-user during authorization. | [optional] 
**ClientSecret** | Pointer to **string** | OAuth 2.0 Client Secret  The secret will be included in the create request as cleartext, and then never again. The secret is kept in hashed format and is not recoverable once lost. | [optional] 
**ClientSecretExpiresAt** | Pointer to **int64** | OAuth 2.0 Client Secret Expires At  The field is currently not supported and its value is always 0. | [optional] 
**ClientUri** | Pointer to **string** | OAuth 2.0 Client URI  ClientURI is a URL string of a web page providing information about the client. If present, the server SHOULD display this URL to the end-user in a clickable fashion. | [optional] 
**Contacts** | Pointer to **[]string** |  | [optional] 
**CreatedAt** | Pointer to **time.Time** | OAuth 2.0 Client Creation Date  CreatedAt returns the timestamp of the client&#39;s creation. | [optional] 
**DeviceAuthorizationGrantAccessTokenLifespan** | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] 
**DeviceAuthorizationGrantIdTokenLifespan** | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] 
**DeviceAuthorizationGrantRefreshTokenLifespan** | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] 
**FrontchannelLogoutSessionRequired** | Pointer to **bool** | OpenID Connect Front-Channel Logout Session Required  Boolean value specifying whether the RP requires that iss (issuer) and sid (session ID) query parameters be included to identify the RP session with the OP when the frontchannel_logout_uri is used. If omitted, the default value is false. | [optional] 
**FrontchannelLogoutUri** | Pointer to **string** | OpenID Connect Front-Channel Logout URI  RP URL that will cause the RP to log itself out when rendered in an iframe by the OP. An iss (issuer) query parameter and a sid (session ID) query parameter MAY be included by the OP to enable the RP to validate the request and to determine which of the potentially multiple sessions is to be logged out; if either is included, both MUST be. | [optional] 
**GrantTypes** | Pointer to **[]string** |  | [optional] 
**ImplicitGrantAccessTokenLifespan** | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] 
**ImplicitGrantIdTokenLifespan** | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] 
**Jwks** | Pointer to [**JsonWebKeySet**](JsonWebKeySet.md) |  | [optional] 
**JwksUri** | Pointer to **string** | OAuth 2.0 Client JSON Web Key Set URL  URL for the Client&#39;s JSON Web Key Set [JWK] document. If the Client signs requests to the Server, it contains the signing key(s) the Server uses to validate signatures from the Client. The JWK Set MAY also contain the Client&#39;s encryption keys(s), which are used by the Server to encrypt responses to the Client. When both signing and encryption keys are made available, a use (Key Use) parameter value is REQUIRED for all keys in the referenced JWK Set to indicate each key&#39;s intended usage. Although some algorithms allow the same key to be used for both signatures and encryption, doing so is NOT RECOMMENDED, as it is less secure. The JWK x5c parameter MAY be used to provide X.509 representations of keys provided. When used, the bare key values MUST still be present and MUST match those in the certificate. | [optional] 
**JwtBearerGrantAccessTokenLifespan** | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] 
**LogoUri** | Pointer to **string** | OAuth 2.0 Client Logo URI  A URL string referencing the client&#39;s logo. | [optional] 
**Metadata** | Pointer to **interface{}** |  | [optional] 
**Owner** | Pointer to **string** | OAuth 2.0 Client Owner  Owner is a string identifying the owner of the OAuth 2.0 Client. | [optional] 
**PolicyUri** | Pointer to **string** | OAuth 2.0 Client Policy URI  PolicyURI is a URL string that points to a human-readable privacy policy document that describes how the deployment organization collects, uses, retains, and discloses personal data. | [optional] 
**PostLogoutRedirectUris** | Pointer to **[]string** |  | [optional] 
**RedirectUris** | Pointer to **[]string** |  | [optional] 
**RefreshTokenGrantAccessTokenLifespan** | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] 
**RefreshTokenGrantIdTokenLifespan** | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] 
**RefreshTokenGrantRefreshTokenLifespan** | Pointer to **string** | Specify a time duration in milliseconds, seconds, minutes, hours. | [optional] 
**RegistrationAccessToken** | Pointer to **string** | OpenID Connect Dynamic Client Registration Access Token  RegistrationAccessToken can be used to update, get, or delete the OAuth2 Client. It is sent when creating a client using Dynamic Client Registration. | [optional] 
**RegistrationClientUri** | Pointer to **string** | OpenID Connect Dynamic Client Registration URL  RegistrationClientURI is the URL used to update, get, or delete the OAuth2 Client. | [optional] 
**RequestObjectSigningAlg** | Pointer to **string** | OpenID Connect Request Object Signing Algorithm  JWS [JWS] alg algorithm [JWA] that MUST be used for signing Request Objects sent to the OP. All Request Objects from this Client MUST be rejected, if not signed with this algorithm. | [optional] 
**RequestUris** | Pointer to **[]string** |  | [optional] 
**ResponseTypes** | Pointer to **[]string** |  | [optional] 
**Scope** | Pointer to **string** | OAuth 2.0 Client Scope  Scope is a string containing a space-separated list of scope values (as described in Section 3.3 of OAuth 2.0 [RFC6749]) that the client can use when requesting access tokens. | [optional] 
**SectorIdentifierUri** | Pointer to **string** | OpenID Connect Sector Identifier URI  URL using the https scheme to be used in calculating Pseudonymous Identifiers by the OP. The URL references a file with a single JSON array of redirect_uri values. | [optional] 
**SkipConsent** | Pointer to **bool** | SkipConsent skips the consent screen for this client. This field can only be set from the admin API. | [optional] 
**SkipLogoutConsent** | Pointer to **bool** | SkipLogoutConsent skips the logout consent screen for this client. This field can only be set from the admin API. | [optional] 
**SubjectType** | Pointer to **string** | OpenID Connect Subject Type  The &#x60;subject_types_supported&#x60; Discovery parameter contains a list of the supported subject_type values for this server. Valid types include &#x60;pairwise&#x60; and &#x60;public&#x60;. | [optional] 
**TokenEndpointAuthMethod** | Pointer to **string** | OAuth 2.0 Token Endpoint Authentication Method  Requested Client Authentication method for the Token Endpoint. The options are:  &#x60;client_secret_basic&#x60;: (default) Send &#x60;client_id&#x60; and &#x60;client_secret&#x60; as &#x60;application/x-www-form-urlencoded&#x60; encoded in the HTTP Authorization header. &#x60;client_secret_post&#x60;: Send &#x60;client_id&#x60; and &#x60;client_secret&#x60; as &#x60;application/x-www-form-urlencoded&#x60; in the HTTP body. &#x60;private_key_jwt&#x60;: Use JSON Web Tokens to authenticate the client. &#x60;none&#x60;: Used for public clients (native apps, mobile apps) which can not have secrets. | [optional] [default to "client_secret_basic"]
**TokenEndpointAuthSigningAlg** | Pointer to **string** | OAuth 2.0 Token Endpoint Signing Algorithm  Requested Client Authentication signing algorithm for the Token Endpoint. | [optional] 
**TosUri** | Pointer to **string** | OAuth 2.0 Client Terms of Service URI  A URL string pointing to a human-readable terms of service document for the client that describes a contractual relationship between the end-user and the client that the end-user accepts when authorizing the client. | [optional] 
**UpdatedAt** | Pointer to **time.Time** | OAuth 2.0 Client Last Update Date  UpdatedAt returns the timestamp of the last update. | [optional] 
**UserinfoSignedResponseAlg** | Pointer to **string** | OpenID Connect Request Userinfo Signed Response Algorithm  JWS alg algorithm [JWA] REQUIRED for signing UserInfo Responses. If this is specified, the response will be JWT [JWT] serialized, and signed using JWS. The default, if omitted, is for the UserInfo Response to return the Claims as a UTF-8 encoded JSON object using the application/json content-type. | [optional] 

## Methods

### NewOAuth2Client

`func NewOAuth2Client() *OAuth2Client`

NewOAuth2Client instantiates a new OAuth2Client object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOAuth2ClientWithDefaults

`func NewOAuth2ClientWithDefaults() *OAuth2Client`

NewOAuth2ClientWithDefaults instantiates a new OAuth2Client object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAccessTokenStrategy

`func (o *OAuth2Client) GetAccessTokenStrategy() string`

GetAccessTokenStrategy returns the AccessTokenStrategy field if non-nil, zero value otherwise.

### GetAccessTokenStrategyOk

`func (o *OAuth2Client) GetAccessTokenStrategyOk() (*string, bool)`

GetAccessTokenStrategyOk returns a tuple with the AccessTokenStrategy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccessTokenStrategy

`func (o *OAuth2Client) SetAccessTokenStrategy(v string)`

SetAccessTokenStrategy sets AccessTokenStrategy field to given value.

### HasAccessTokenStrategy

`func (o *OAuth2Client) HasAccessTokenStrategy() bool`

HasAccessTokenStrategy returns a boolean if a field has been set.

### GetAllowedCorsOrigins

`func (o *OAuth2Client) GetAllowedCorsOrigins() []string`

GetAllowedCorsOrigins returns the AllowedCorsOrigins field if non-nil, zero value otherwise.

### GetAllowedCorsOriginsOk

`func (o *OAuth2Client) GetAllowedCorsOriginsOk() (*[]string, bool)`

GetAllowedCorsOriginsOk returns a tuple with the AllowedCorsOrigins field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllowedCorsOrigins

`func (o *OAuth2Client) SetAllowedCorsOrigins(v []string)`

SetAllowedCorsOrigins sets AllowedCorsOrigins field to given value.

### HasAllowedCorsOrigins

`func (o *OAuth2Client) HasAllowedCorsOrigins() bool`

HasAllowedCorsOrigins returns a boolean if a field has been set.

### GetAudience

`func (o *OAuth2Client) GetAudience() []string`

GetAudience returns the Audience field if non-nil, zero value otherwise.

### GetAudienceOk

`func (o *OAuth2Client) GetAudienceOk() (*[]string, bool)`

GetAudienceOk returns a tuple with the Audience field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAudience

`func (o *OAuth2Client) SetAudience(v []string)`

SetAudience sets Audience field to given value.

### HasAudience

`func (o *OAuth2Client) HasAudience() bool`

HasAudience returns a boolean if a field has been set.

### GetAuthorizationCodeGrantAccessTokenLifespan

`func (o *OAuth2Client) GetAuthorizationCodeGrantAccessTokenLifespan() string`

GetAuthorizationCodeGrantAccessTokenLifespan returns the AuthorizationCodeGrantAccessTokenLifespan field if non-nil, zero value otherwise.

### GetAuthorizationCodeGrantAccessTokenLifespanOk

`func (o *OAuth2Client) GetAuthorizationCodeGrantAccessTokenLifespanOk() (*string, bool)`

GetAuthorizationCodeGrantAccessTokenLifespanOk returns a tuple with the AuthorizationCodeGrantAccessTokenLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthorizationCodeGrantAccessTokenLifespan

`func (o *OAuth2Client) SetAuthorizationCodeGrantAccessTokenLifespan(v string)`

SetAuthorizationCodeGrantAccessTokenLifespan sets AuthorizationCodeGrantAccessTokenLifespan field to given value.

### HasAuthorizationCodeGrantAccessTokenLifespan

`func (o *OAuth2Client) HasAuthorizationCodeGrantAccessTokenLifespan() bool`

HasAuthorizationCodeGrantAccessTokenLifespan returns a boolean if a field has been set.

### GetAuthorizationCodeGrantIdTokenLifespan

`func (o *OAuth2Client) GetAuthorizationCodeGrantIdTokenLifespan() string`

GetAuthorizationCodeGrantIdTokenLifespan returns the AuthorizationCodeGrantIdTokenLifespan field if non-nil, zero value otherwise.

### GetAuthorizationCodeGrantIdTokenLifespanOk

`func (o *OAuth2Client) GetAuthorizationCodeGrantIdTokenLifespanOk() (*string, bool)`

GetAuthorizationCodeGrantIdTokenLifespanOk returns a tuple with the AuthorizationCodeGrantIdTokenLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthorizationCodeGrantIdTokenLifespan

`func (o *OAuth2Client) SetAuthorizationCodeGrantIdTokenLifespan(v string)`

SetAuthorizationCodeGrantIdTokenLifespan sets AuthorizationCodeGrantIdTokenLifespan field to given value.

### HasAuthorizationCodeGrantIdTokenLifespan

`func (o *OAuth2Client) HasAuthorizationCodeGrantIdTokenLifespan() bool`

HasAuthorizationCodeGrantIdTokenLifespan returns a boolean if a field has been set.

### GetAuthorizationCodeGrantRefreshTokenLifespan

`func (o *OAuth2Client) GetAuthorizationCodeGrantRefreshTokenLifespan() string`

GetAuthorizationCodeGrantRefreshTokenLifespan returns the AuthorizationCodeGrantRefreshTokenLifespan field if non-nil, zero value otherwise.

### GetAuthorizationCodeGrantRefreshTokenLifespanOk

`func (o *OAuth2Client) GetAuthorizationCodeGrantRefreshTokenLifespanOk() (*string, bool)`

GetAuthorizationCodeGrantRefreshTokenLifespanOk returns a tuple with the AuthorizationCodeGrantRefreshTokenLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthorizationCodeGrantRefreshTokenLifespan

`func (o *OAuth2Client) SetAuthorizationCodeGrantRefreshTokenLifespan(v string)`

SetAuthorizationCodeGrantRefreshTokenLifespan sets AuthorizationCodeGrantRefreshTokenLifespan field to given value.

### HasAuthorizationCodeGrantRefreshTokenLifespan

`func (o *OAuth2Client) HasAuthorizationCodeGrantRefreshTokenLifespan() bool`

HasAuthorizationCodeGrantRefreshTokenLifespan returns a boolean if a field has been set.

### GetBackchannelLogoutSessionRequired

`func (o *OAuth2Client) GetBackchannelLogoutSessionRequired() bool`

GetBackchannelLogoutSessionRequired returns the BackchannelLogoutSessionRequired field if non-nil, zero value otherwise.

### GetBackchannelLogoutSessionRequiredOk

`func (o *OAuth2Client) GetBackchannelLogoutSessionRequiredOk() (*bool, bool)`

GetBackchannelLogoutSessionRequiredOk returns a tuple with the BackchannelLogoutSessionRequired field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBackchannelLogoutSessionRequired

`func (o *OAuth2Client) SetBackchannelLogoutSessionRequired(v bool)`

SetBackchannelLogoutSessionRequired sets BackchannelLogoutSessionRequired field to given value.

### HasBackchannelLogoutSessionRequired

`func (o *OAuth2Client) HasBackchannelLogoutSessionRequired() bool`

HasBackchannelLogoutSessionRequired returns a boolean if a field has been set.

### GetBackchannelLogoutUri

`func (o *OAuth2Client) GetBackchannelLogoutUri() string`

GetBackchannelLogoutUri returns the BackchannelLogoutUri field if non-nil, zero value otherwise.

### GetBackchannelLogoutUriOk

`func (o *OAuth2Client) GetBackchannelLogoutUriOk() (*string, bool)`

GetBackchannelLogoutUriOk returns a tuple with the BackchannelLogoutUri field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBackchannelLogoutUri

`func (o *OAuth2Client) SetBackchannelLogoutUri(v string)`

SetBackchannelLogoutUri sets BackchannelLogoutUri field to given value.

### HasBackchannelLogoutUri

`func (o *OAuth2Client) HasBackchannelLogoutUri() bool`

HasBackchannelLogoutUri returns a boolean if a field has been set.

### GetClientCredentialsGrantAccessTokenLifespan

`func (o *OAuth2Client) GetClientCredentialsGrantAccessTokenLifespan() string`

GetClientCredentialsGrantAccessTokenLifespan returns the ClientCredentialsGrantAccessTokenLifespan field if non-nil, zero value otherwise.

### GetClientCredentialsGrantAccessTokenLifespanOk

`func (o *OAuth2Client) GetClientCredentialsGrantAccessTokenLifespanOk() (*string, bool)`

GetClientCredentialsGrantAccessTokenLifespanOk returns a tuple with the ClientCredentialsGrantAccessTokenLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClientCredentialsGrantAccessTokenLifespan

`func (o *OAuth2Client) SetClientCredentialsGrantAccessTokenLifespan(v string)`

SetClientCredentialsGrantAccessTokenLifespan sets ClientCredentialsGrantAccessTokenLifespan field to given value.

### HasClientCredentialsGrantAccessTokenLifespan

`func (o *OAuth2Client) HasClientCredentialsGrantAccessTokenLifespan() bool`

HasClientCredentialsGrantAccessTokenLifespan returns a boolean if a field has been set.

### GetClientId

`func (o *OAuth2Client) GetClientId() string`

GetClientId returns the ClientId field if non-nil, zero value otherwise.

### GetClientIdOk

`func (o *OAuth2Client) GetClientIdOk() (*string, bool)`

GetClientIdOk returns a tuple with the ClientId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClientId

`func (o *OAuth2Client) SetClientId(v string)`

SetClientId sets ClientId field to given value.

### HasClientId

`func (o *OAuth2Client) HasClientId() bool`

HasClientId returns a boolean if a field has been set.

### GetClientName

`func (o *OAuth2Client) GetClientName() string`

GetClientName returns the ClientName field if non-nil, zero value otherwise.

### GetClientNameOk

`func (o *OAuth2Client) GetClientNameOk() (*string, bool)`

GetClientNameOk returns a tuple with the ClientName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClientName

`func (o *OAuth2Client) SetClientName(v string)`

SetClientName sets ClientName field to given value.

### HasClientName

`func (o *OAuth2Client) HasClientName() bool`

HasClientName returns a boolean if a field has been set.

### GetClientSecret

`func (o *OAuth2Client) GetClientSecret() string`

GetClientSecret returns the ClientSecret field if non-nil, zero value otherwise.

### GetClientSecretOk

`func (o *OAuth2Client) GetClientSecretOk() (*string, bool)`

GetClientSecretOk returns a tuple with the ClientSecret field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClientSecret

`func (o *OAuth2Client) SetClientSecret(v string)`

SetClientSecret sets ClientSecret field to given value.

### HasClientSecret

`func (o *OAuth2Client) HasClientSecret() bool`

HasClientSecret returns a boolean if a field has been set.

### GetClientSecretExpiresAt

`func (o *OAuth2Client) GetClientSecretExpiresAt() int64`

GetClientSecretExpiresAt returns the ClientSecretExpiresAt field if non-nil, zero value otherwise.

### GetClientSecretExpiresAtOk

`func (o *OAuth2Client) GetClientSecretExpiresAtOk() (*int64, bool)`

GetClientSecretExpiresAtOk returns a tuple with the ClientSecretExpiresAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClientSecretExpiresAt

`func (o *OAuth2Client) SetClientSecretExpiresAt(v int64)`

SetClientSecretExpiresAt sets ClientSecretExpiresAt field to given value.

### HasClientSecretExpiresAt

`func (o *OAuth2Client) HasClientSecretExpiresAt() bool`

HasClientSecretExpiresAt returns a boolean if a field has been set.

### GetClientUri

`func (o *OAuth2Client) GetClientUri() string`

GetClientUri returns the ClientUri field if non-nil, zero value otherwise.

### GetClientUriOk

`func (o *OAuth2Client) GetClientUriOk() (*string, bool)`

GetClientUriOk returns a tuple with the ClientUri field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClientUri

`func (o *OAuth2Client) SetClientUri(v string)`

SetClientUri sets ClientUri field to given value.

### HasClientUri

`func (o *OAuth2Client) HasClientUri() bool`

HasClientUri returns a boolean if a field has been set.

### GetContacts

`func (o *OAuth2Client) GetContacts() []string`

GetContacts returns the Contacts field if non-nil, zero value otherwise.

### GetContactsOk

`func (o *OAuth2Client) GetContactsOk() (*[]string, bool)`

GetContactsOk returns a tuple with the Contacts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContacts

`func (o *OAuth2Client) SetContacts(v []string)`

SetContacts sets Contacts field to given value.

### HasContacts

`func (o *OAuth2Client) HasContacts() bool`

HasContacts returns a boolean if a field has been set.

### GetCreatedAt

`func (o *OAuth2Client) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *OAuth2Client) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *OAuth2Client) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *OAuth2Client) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetDeviceAuthorizationGrantAccessTokenLifespan

`func (o *OAuth2Client) GetDeviceAuthorizationGrantAccessTokenLifespan() string`

GetDeviceAuthorizationGrantAccessTokenLifespan returns the DeviceAuthorizationGrantAccessTokenLifespan field if non-nil, zero value otherwise.

### GetDeviceAuthorizationGrantAccessTokenLifespanOk

`func (o *OAuth2Client) GetDeviceAuthorizationGrantAccessTokenLifespanOk() (*string, bool)`

GetDeviceAuthorizationGrantAccessTokenLifespanOk returns a tuple with the DeviceAuthorizationGrantAccessTokenLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceAuthorizationGrantAccessTokenLifespan

`func (o *OAuth2Client) SetDeviceAuthorizationGrantAccessTokenLifespan(v string)`

SetDeviceAuthorizationGrantAccessTokenLifespan sets DeviceAuthorizationGrantAccessTokenLifespan field to given value.

### HasDeviceAuthorizationGrantAccessTokenLifespan

`func (o *OAuth2Client) HasDeviceAuthorizationGrantAccessTokenLifespan() bool`

HasDeviceAuthorizationGrantAccessTokenLifespan returns a boolean if a field has been set.

### GetDeviceAuthorizationGrantIdTokenLifespan

`func (o *OAuth2Client) GetDeviceAuthorizationGrantIdTokenLifespan() string`

GetDeviceAuthorizationGrantIdTokenLifespan returns the DeviceAuthorizationGrantIdTokenLifespan field if non-nil, zero value otherwise.

### GetDeviceAuthorizationGrantIdTokenLifespanOk

`func (o *OAuth2Client) GetDeviceAuthorizationGrantIdTokenLifespanOk() (*string, bool)`

GetDeviceAuthorizationGrantIdTokenLifespanOk returns a tuple with the DeviceAuthorizationGrantIdTokenLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceAuthorizationGrantIdTokenLifespan

`func (o *OAuth2Client) SetDeviceAuthorizationGrantIdTokenLifespan(v string)`

SetDeviceAuthorizationGrantIdTokenLifespan sets DeviceAuthorizationGrantIdTokenLifespan field to given value.

### HasDeviceAuthorizationGrantIdTokenLifespan

`func (o *OAuth2Client) HasDeviceAuthorizationGrantIdTokenLifespan() bool`

HasDeviceAuthorizationGrantIdTokenLifespan returns a boolean if a field has been set.

### GetDeviceAuthorizationGrantRefreshTokenLifespan

`func (o *OAuth2Client) GetDeviceAuthorizationGrantRefreshTokenLifespan() string`

GetDeviceAuthorizationGrantRefreshTokenLifespan returns the DeviceAuthorizationGrantRefreshTokenLifespan field if non-nil, zero value otherwise.

### GetDeviceAuthorizationGrantRefreshTokenLifespanOk

`func (o *OAuth2Client) GetDeviceAuthorizationGrantRefreshTokenLifespanOk() (*string, bool)`

GetDeviceAuthorizationGrantRefreshTokenLifespanOk returns a tuple with the DeviceAuthorizationGrantRefreshTokenLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceAuthorizationGrantRefreshTokenLifespan

`func (o *OAuth2Client) SetDeviceAuthorizationGrantRefreshTokenLifespan(v string)`

SetDeviceAuthorizationGrantRefreshTokenLifespan sets DeviceAuthorizationGrantRefreshTokenLifespan field to given value.

### HasDeviceAuthorizationGrantRefreshTokenLifespan

`func (o *OAuth2Client) HasDeviceAuthorizationGrantRefreshTokenLifespan() bool`

HasDeviceAuthorizationGrantRefreshTokenLifespan returns a boolean if a field has been set.

### GetFrontchannelLogoutSessionRequired

`func (o *OAuth2Client) GetFrontchannelLogoutSessionRequired() bool`

GetFrontchannelLogoutSessionRequired returns the FrontchannelLogoutSessionRequired field if non-nil, zero value otherwise.

### GetFrontchannelLogoutSessionRequiredOk

`func (o *OAuth2Client) GetFrontchannelLogoutSessionRequiredOk() (*bool, bool)`

GetFrontchannelLogoutSessionRequiredOk returns a tuple with the FrontchannelLogoutSessionRequired field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFrontchannelLogoutSessionRequired

`func (o *OAuth2Client) SetFrontchannelLogoutSessionRequired(v bool)`

SetFrontchannelLogoutSessionRequired sets FrontchannelLogoutSessionRequired field to given value.

### HasFrontchannelLogoutSessionRequired

`func (o *OAuth2Client) HasFrontchannelLogoutSessionRequired() bool`

HasFrontchannelLogoutSessionRequired returns a boolean if a field has been set.

### GetFrontchannelLogoutUri

`func (o *OAuth2Client) GetFrontchannelLogoutUri() string`

GetFrontchannelLogoutUri returns the FrontchannelLogoutUri field if non-nil, zero value otherwise.

### GetFrontchannelLogoutUriOk

`func (o *OAuth2Client) GetFrontchannelLogoutUriOk() (*string, bool)`

GetFrontchannelLogoutUriOk returns a tuple with the FrontchannelLogoutUri field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFrontchannelLogoutUri

`func (o *OAuth2Client) SetFrontchannelLogoutUri(v string)`

SetFrontchannelLogoutUri sets FrontchannelLogoutUri field to given value.

### HasFrontchannelLogoutUri

`func (o *OAuth2Client) HasFrontchannelLogoutUri() bool`

HasFrontchannelLogoutUri returns a boolean if a field has been set.

### GetGrantTypes

`func (o *OAuth2Client) GetGrantTypes() []string`

GetGrantTypes returns the GrantTypes field if non-nil, zero value otherwise.

### GetGrantTypesOk

`func (o *OAuth2Client) GetGrantTypesOk() (*[]string, bool)`

GetGrantTypesOk returns a tuple with the GrantTypes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGrantTypes

`func (o *OAuth2Client) SetGrantTypes(v []string)`

SetGrantTypes sets GrantTypes field to given value.

### HasGrantTypes

`func (o *OAuth2Client) HasGrantTypes() bool`

HasGrantTypes returns a boolean if a field has been set.

### GetImplicitGrantAccessTokenLifespan

`func (o *OAuth2Client) GetImplicitGrantAccessTokenLifespan() string`

GetImplicitGrantAccessTokenLifespan returns the ImplicitGrantAccessTokenLifespan field if non-nil, zero value otherwise.

### GetImplicitGrantAccessTokenLifespanOk

`func (o *OAuth2Client) GetImplicitGrantAccessTokenLifespanOk() (*string, bool)`

GetImplicitGrantAccessTokenLifespanOk returns a tuple with the ImplicitGrantAccessTokenLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImplicitGrantAccessTokenLifespan

`func (o *OAuth2Client) SetImplicitGrantAccessTokenLifespan(v string)`

SetImplicitGrantAccessTokenLifespan sets ImplicitGrantAccessTokenLifespan field to given value.

### HasImplicitGrantAccessTokenLifespan

`func (o *OAuth2Client) HasImplicitGrantAccessTokenLifespan() bool`

HasImplicitGrantAccessTokenLifespan returns a boolean if a field has been set.

### GetImplicitGrantIdTokenLifespan

`func (o *OAuth2Client) GetImplicitGrantIdTokenLifespan() string`

GetImplicitGrantIdTokenLifespan returns the ImplicitGrantIdTokenLifespan field if non-nil, zero value otherwise.

### GetImplicitGrantIdTokenLifespanOk

`func (o *OAuth2Client) GetImplicitGrantIdTokenLifespanOk() (*string, bool)`

GetImplicitGrantIdTokenLifespanOk returns a tuple with the ImplicitGrantIdTokenLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImplicitGrantIdTokenLifespan

`func (o *OAuth2Client) SetImplicitGrantIdTokenLifespan(v string)`

SetImplicitGrantIdTokenLifespan sets ImplicitGrantIdTokenLifespan field to given value.

### HasImplicitGrantIdTokenLifespan

`func (o *OAuth2Client) HasImplicitGrantIdTokenLifespan() bool`

HasImplicitGrantIdTokenLifespan returns a boolean if a field has been set.

### GetJwks

`func (o *OAuth2Client) GetJwks() JsonWebKeySet`

GetJwks returns the Jwks field if non-nil, zero value otherwise.

### GetJwksOk

`func (o *OAuth2Client) GetJwksOk() (*JsonWebKeySet, bool)`

GetJwksOk returns a tuple with the Jwks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetJwks

`func (o *OAuth2Client) SetJwks(v JsonWebKeySet)`

SetJwks sets Jwks field to given value.

### HasJwks

`func (o *OAuth2Client) HasJwks() bool`

HasJwks returns a boolean if a field has been set.

### GetJwksUri

`func (o *OAuth2Client) GetJwksUri() string`

GetJwksUri returns the JwksUri field if non-nil, zero value otherwise.

### GetJwksUriOk

`func (o *OAuth2Client) GetJwksUriOk() (*string, bool)`

GetJwksUriOk returns a tuple with the JwksUri field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetJwksUri

`func (o *OAuth2Client) SetJwksUri(v string)`

SetJwksUri sets JwksUri field to given value.

### HasJwksUri

`func (o *OAuth2Client) HasJwksUri() bool`

HasJwksUri returns a boolean if a field has been set.

### GetJwtBearerGrantAccessTokenLifespan

`func (o *OAuth2Client) GetJwtBearerGrantAccessTokenLifespan() string`

GetJwtBearerGrantAccessTokenLifespan returns the JwtBearerGrantAccessTokenLifespan field if non-nil, zero value otherwise.

### GetJwtBearerGrantAccessTokenLifespanOk

`func (o *OAuth2Client) GetJwtBearerGrantAccessTokenLifespanOk() (*string, bool)`

GetJwtBearerGrantAccessTokenLifespanOk returns a tuple with the JwtBearerGrantAccessTokenLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetJwtBearerGrantAccessTokenLifespan

`func (o *OAuth2Client) SetJwtBearerGrantAccessTokenLifespan(v string)`

SetJwtBearerGrantAccessTokenLifespan sets JwtBearerGrantAccessTokenLifespan field to given value.

### HasJwtBearerGrantAccessTokenLifespan

`func (o *OAuth2Client) HasJwtBearerGrantAccessTokenLifespan() bool`

HasJwtBearerGrantAccessTokenLifespan returns a boolean if a field has been set.

### GetLogoUri

`func (o *OAuth2Client) GetLogoUri() string`

GetLogoUri returns the LogoUri field if non-nil, zero value otherwise.

### GetLogoUriOk

`func (o *OAuth2Client) GetLogoUriOk() (*string, bool)`

GetLogoUriOk returns a tuple with the LogoUri field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLogoUri

`func (o *OAuth2Client) SetLogoUri(v string)`

SetLogoUri sets LogoUri field to given value.

### HasLogoUri

`func (o *OAuth2Client) HasLogoUri() bool`

HasLogoUri returns a boolean if a field has been set.

### GetMetadata

`func (o *OAuth2Client) GetMetadata() interface{}`

GetMetadata returns the Metadata field if non-nil, zero value otherwise.

### GetMetadataOk

`func (o *OAuth2Client) GetMetadataOk() (*interface{}, bool)`

GetMetadataOk returns a tuple with the Metadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetadata

`func (o *OAuth2Client) SetMetadata(v interface{})`

SetMetadata sets Metadata field to given value.

### HasMetadata

`func (o *OAuth2Client) HasMetadata() bool`

HasMetadata returns a boolean if a field has been set.

### SetMetadataNil

`func (o *OAuth2Client) SetMetadataNil(b bool)`

 SetMetadataNil sets the value for Metadata to be an explicit nil

### UnsetMetadata
`func (o *OAuth2Client) UnsetMetadata()`

UnsetMetadata ensures that no value is present for Metadata, not even an explicit nil
### GetOwner

`func (o *OAuth2Client) GetOwner() string`

GetOwner returns the Owner field if non-nil, zero value otherwise.

### GetOwnerOk

`func (o *OAuth2Client) GetOwnerOk() (*string, bool)`

GetOwnerOk returns a tuple with the Owner field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOwner

`func (o *OAuth2Client) SetOwner(v string)`

SetOwner sets Owner field to given value.

### HasOwner

`func (o *OAuth2Client) HasOwner() bool`

HasOwner returns a boolean if a field has been set.

### GetPolicyUri

`func (o *OAuth2Client) GetPolicyUri() string`

GetPolicyUri returns the PolicyUri field if non-nil, zero value otherwise.

### GetPolicyUriOk

`func (o *OAuth2Client) GetPolicyUriOk() (*string, bool)`

GetPolicyUriOk returns a tuple with the PolicyUri field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPolicyUri

`func (o *OAuth2Client) SetPolicyUri(v string)`

SetPolicyUri sets PolicyUri field to given value.

### HasPolicyUri

`func (o *OAuth2Client) HasPolicyUri() bool`

HasPolicyUri returns a boolean if a field has been set.

### GetPostLogoutRedirectUris

`func (o *OAuth2Client) GetPostLogoutRedirectUris() []string`

GetPostLogoutRedirectUris returns the PostLogoutRedirectUris field if non-nil, zero value otherwise.

### GetPostLogoutRedirectUrisOk

`func (o *OAuth2Client) GetPostLogoutRedirectUrisOk() (*[]string, bool)`

GetPostLogoutRedirectUrisOk returns a tuple with the PostLogoutRedirectUris field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPostLogoutRedirectUris

`func (o *OAuth2Client) SetPostLogoutRedirectUris(v []string)`

SetPostLogoutRedirectUris sets PostLogoutRedirectUris field to given value.

### HasPostLogoutRedirectUris

`func (o *OAuth2Client) HasPostLogoutRedirectUris() bool`

HasPostLogoutRedirectUris returns a boolean if a field has been set.

### GetRedirectUris

`func (o *OAuth2Client) GetRedirectUris() []string`

GetRedirectUris returns the RedirectUris field if non-nil, zero value otherwise.

### GetRedirectUrisOk

`func (o *OAuth2Client) GetRedirectUrisOk() (*[]string, bool)`

GetRedirectUrisOk returns a tuple with the RedirectUris field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRedirectUris

`func (o *OAuth2Client) SetRedirectUris(v []string)`

SetRedirectUris sets RedirectUris field to given value.

### HasRedirectUris

`func (o *OAuth2Client) HasRedirectUris() bool`

HasRedirectUris returns a boolean if a field has been set.

### GetRefreshTokenGrantAccessTokenLifespan

`func (o *OAuth2Client) GetRefreshTokenGrantAccessTokenLifespan() string`

GetRefreshTokenGrantAccessTokenLifespan returns the RefreshTokenGrantAccessTokenLifespan field if non-nil, zero value otherwise.

### GetRefreshTokenGrantAccessTokenLifespanOk

`func (o *OAuth2Client) GetRefreshTokenGrantAccessTokenLifespanOk() (*string, bool)`

GetRefreshTokenGrantAccessTokenLifespanOk returns a tuple with the RefreshTokenGrantAccessTokenLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRefreshTokenGrantAccessTokenLifespan

`func (o *OAuth2Client) SetRefreshTokenGrantAccessTokenLifespan(v string)`

SetRefreshTokenGrantAccessTokenLifespan sets RefreshTokenGrantAccessTokenLifespan field to given value.

### HasRefreshTokenGrantAccessTokenLifespan

`func (o *OAuth2Client) HasRefreshTokenGrantAccessTokenLifespan() bool`

HasRefreshTokenGrantAccessTokenLifespan returns a boolean if a field has been set.

### GetRefreshTokenGrantIdTokenLifespan

`func (o *OAuth2Client) GetRefreshTokenGrantIdTokenLifespan() string`

GetRefreshTokenGrantIdTokenLifespan returns the RefreshTokenGrantIdTokenLifespan field if non-nil, zero value otherwise.

### GetRefreshTokenGrantIdTokenLifespanOk

`func (o *OAuth2Client) GetRefreshTokenGrantIdTokenLifespanOk() (*string, bool)`

GetRefreshTokenGrantIdTokenLifespanOk returns a tuple with the RefreshTokenGrantIdTokenLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRefreshTokenGrantIdTokenLifespan

`func (o *OAuth2Client) SetRefreshTokenGrantIdTokenLifespan(v string)`

SetRefreshTokenGrantIdTokenLifespan sets RefreshTokenGrantIdTokenLifespan field to given value.

### HasRefreshTokenGrantIdTokenLifespan

`func (o *OAuth2Client) HasRefreshTokenGrantIdTokenLifespan() bool`

HasRefreshTokenGrantIdTokenLifespan returns a boolean if a field has been set.

### GetRefreshTokenGrantRefreshTokenLifespan

`func (o *OAuth2Client) GetRefreshTokenGrantRefreshTokenLifespan() string`

GetRefreshTokenGrantRefreshTokenLifespan returns the RefreshTokenGrantRefreshTokenLifespan field if non-nil, zero value otherwise.

### GetRefreshTokenGrantRefreshTokenLifespanOk

`func (o *OAuth2Client) GetRefreshTokenGrantRefreshTokenLifespanOk() (*string, bool)`

GetRefreshTokenGrantRefreshTokenLifespanOk returns a tuple with the RefreshTokenGrantRefreshTokenLifespan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRefreshTokenGrantRefreshTokenLifespan

`func (o *OAuth2Client) SetRefreshTokenGrantRefreshTokenLifespan(v string)`

SetRefreshTokenGrantRefreshTokenLifespan sets RefreshTokenGrantRefreshTokenLifespan field to given value.

### HasRefreshTokenGrantRefreshTokenLifespan

`func (o *OAuth2Client) HasRefreshTokenGrantRefreshTokenLifespan() bool`

HasRefreshTokenGrantRefreshTokenLifespan returns a boolean if a field has been set.

### GetRegistrationAccessToken

`func (o *OAuth2Client) GetRegistrationAccessToken() string`

GetRegistrationAccessToken returns the RegistrationAccessToken field if non-nil, zero value otherwise.

### GetRegistrationAccessTokenOk

`func (o *OAuth2Client) GetRegistrationAccessTokenOk() (*string, bool)`

GetRegistrationAccessTokenOk returns a tuple with the RegistrationAccessToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegistrationAccessToken

`func (o *OAuth2Client) SetRegistrationAccessToken(v string)`

SetRegistrationAccessToken sets RegistrationAccessToken field to given value.

### HasRegistrationAccessToken

`func (o *OAuth2Client) HasRegistrationAccessToken() bool`

HasRegistrationAccessToken returns a boolean if a field has been set.

### GetRegistrationClientUri

`func (o *OAuth2Client) GetRegistrationClientUri() string`

GetRegistrationClientUri returns the RegistrationClientUri field if non-nil, zero value otherwise.

### GetRegistrationClientUriOk

`func (o *OAuth2Client) GetRegistrationClientUriOk() (*string, bool)`

GetRegistrationClientUriOk returns a tuple with the RegistrationClientUri field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegistrationClientUri

`func (o *OAuth2Client) SetRegistrationClientUri(v string)`

SetRegistrationClientUri sets RegistrationClientUri field to given value.

### HasRegistrationClientUri

`func (o *OAuth2Client) HasRegistrationClientUri() bool`

HasRegistrationClientUri returns a boolean if a field has been set.

### GetRequestObjectSigningAlg

`func (o *OAuth2Client) GetRequestObjectSigningAlg() string`

GetRequestObjectSigningAlg returns the RequestObjectSigningAlg field if non-nil, zero value otherwise.

### GetRequestObjectSigningAlgOk

`func (o *OAuth2Client) GetRequestObjectSigningAlgOk() (*string, bool)`

GetRequestObjectSigningAlgOk returns a tuple with the RequestObjectSigningAlg field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestObjectSigningAlg

`func (o *OAuth2Client) SetRequestObjectSigningAlg(v string)`

SetRequestObjectSigningAlg sets RequestObjectSigningAlg field to given value.

### HasRequestObjectSigningAlg

`func (o *OAuth2Client) HasRequestObjectSigningAlg() bool`

HasRequestObjectSigningAlg returns a boolean if a field has been set.

### GetRequestUris

`func (o *OAuth2Client) GetRequestUris() []string`

GetRequestUris returns the RequestUris field if non-nil, zero value otherwise.

### GetRequestUrisOk

`func (o *OAuth2Client) GetRequestUrisOk() (*[]string, bool)`

GetRequestUrisOk returns a tuple with the RequestUris field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestUris

`func (o *OAuth2Client) SetRequestUris(v []string)`

SetRequestUris sets RequestUris field to given value.

### HasRequestUris

`func (o *OAuth2Client) HasRequestUris() bool`

HasRequestUris returns a boolean if a field has been set.

### GetResponseTypes

`func (o *OAuth2Client) GetResponseTypes() []string`

GetResponseTypes returns the ResponseTypes field if non-nil, zero value otherwise.

### GetResponseTypesOk

`func (o *OAuth2Client) GetResponseTypesOk() (*[]string, bool)`

GetResponseTypesOk returns a tuple with the ResponseTypes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResponseTypes

`func (o *OAuth2Client) SetResponseTypes(v []string)`

SetResponseTypes sets ResponseTypes field to given value.

### HasResponseTypes

`func (o *OAuth2Client) HasResponseTypes() bool`

HasResponseTypes returns a boolean if a field has been set.

### GetScope

`func (o *OAuth2Client) GetScope() string`

GetScope returns the Scope field if non-nil, zero value otherwise.

### GetScopeOk

`func (o *OAuth2Client) GetScopeOk() (*string, bool)`

GetScopeOk returns a tuple with the Scope field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScope

`func (o *OAuth2Client) SetScope(v string)`

SetScope sets Scope field to given value.

### HasScope

`func (o *OAuth2Client) HasScope() bool`

HasScope returns a boolean if a field has been set.

### GetSectorIdentifierUri

`func (o *OAuth2Client) GetSectorIdentifierUri() string`

GetSectorIdentifierUri returns the SectorIdentifierUri field if non-nil, zero value otherwise.

### GetSectorIdentifierUriOk

`func (o *OAuth2Client) GetSectorIdentifierUriOk() (*string, bool)`

GetSectorIdentifierUriOk returns a tuple with the SectorIdentifierUri field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSectorIdentifierUri

`func (o *OAuth2Client) SetSectorIdentifierUri(v string)`

SetSectorIdentifierUri sets SectorIdentifierUri field to given value.

### HasSectorIdentifierUri

`func (o *OAuth2Client) HasSectorIdentifierUri() bool`

HasSectorIdentifierUri returns a boolean if a field has been set.

### GetSkipConsent

`func (o *OAuth2Client) GetSkipConsent() bool`

GetSkipConsent returns the SkipConsent field if non-nil, zero value otherwise.

### GetSkipConsentOk

`func (o *OAuth2Client) GetSkipConsentOk() (*bool, bool)`

GetSkipConsentOk returns a tuple with the SkipConsent field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSkipConsent

`func (o *OAuth2Client) SetSkipConsent(v bool)`

SetSkipConsent sets SkipConsent field to given value.

### HasSkipConsent

`func (o *OAuth2Client) HasSkipConsent() bool`

HasSkipConsent returns a boolean if a field has been set.

### GetSkipLogoutConsent

`func (o *OAuth2Client) GetSkipLogoutConsent() bool`

GetSkipLogoutConsent returns the SkipLogoutConsent field if non-nil, zero value otherwise.

### GetSkipLogoutConsentOk

`func (o *OAuth2Client) GetSkipLogoutConsentOk() (*bool, bool)`

GetSkipLogoutConsentOk returns a tuple with the SkipLogoutConsent field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSkipLogoutConsent

`func (o *OAuth2Client) SetSkipLogoutConsent(v bool)`

SetSkipLogoutConsent sets SkipLogoutConsent field to given value.

### HasSkipLogoutConsent

`func (o *OAuth2Client) HasSkipLogoutConsent() bool`

HasSkipLogoutConsent returns a boolean if a field has been set.

### GetSubjectType

`func (o *OAuth2Client) GetSubjectType() string`

GetSubjectType returns the SubjectType field if non-nil, zero value otherwise.

### GetSubjectTypeOk

`func (o *OAuth2Client) GetSubjectTypeOk() (*string, bool)`

GetSubjectTypeOk returns a tuple with the SubjectType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubjectType

`func (o *OAuth2Client) SetSubjectType(v string)`

SetSubjectType sets SubjectType field to given value.

### HasSubjectType

`func (o *OAuth2Client) HasSubjectType() bool`

HasSubjectType returns a boolean if a field has been set.

### GetTokenEndpointAuthMethod

`func (o *OAuth2Client) GetTokenEndpointAuthMethod() string`

GetTokenEndpointAuthMethod returns the TokenEndpointAuthMethod field if non-nil, zero value otherwise.

### GetTokenEndpointAuthMethodOk

`func (o *OAuth2Client) GetTokenEndpointAuthMethodOk() (*string, bool)`

GetTokenEndpointAuthMethodOk returns a tuple with the TokenEndpointAuthMethod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTokenEndpointAuthMethod

`func (o *OAuth2Client) SetTokenEndpointAuthMethod(v string)`

SetTokenEndpointAuthMethod sets TokenEndpointAuthMethod field to given value.

### HasTokenEndpointAuthMethod

`func (o *OAuth2Client) HasTokenEndpointAuthMethod() bool`

HasTokenEndpointAuthMethod returns a boolean if a field has been set.

### GetTokenEndpointAuthSigningAlg

`func (o *OAuth2Client) GetTokenEndpointAuthSigningAlg() string`

GetTokenEndpointAuthSigningAlg returns the TokenEndpointAuthSigningAlg field if non-nil, zero value otherwise.

### GetTokenEndpointAuthSigningAlgOk

`func (o *OAuth2Client) GetTokenEndpointAuthSigningAlgOk() (*string, bool)`

GetTokenEndpointAuthSigningAlgOk returns a tuple with the TokenEndpointAuthSigningAlg field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTokenEndpointAuthSigningAlg

`func (o *OAuth2Client) SetTokenEndpointAuthSigningAlg(v string)`

SetTokenEndpointAuthSigningAlg sets TokenEndpointAuthSigningAlg field to given value.

### HasTokenEndpointAuthSigningAlg

`func (o *OAuth2Client) HasTokenEndpointAuthSigningAlg() bool`

HasTokenEndpointAuthSigningAlg returns a boolean if a field has been set.

### GetTosUri

`func (o *OAuth2Client) GetTosUri() string`

GetTosUri returns the TosUri field if non-nil, zero value otherwise.

### GetTosUriOk

`func (o *OAuth2Client) GetTosUriOk() (*string, bool)`

GetTosUriOk returns a tuple with the TosUri field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTosUri

`func (o *OAuth2Client) SetTosUri(v string)`

SetTosUri sets TosUri field to given value.

### HasTosUri

`func (o *OAuth2Client) HasTosUri() bool`

HasTosUri returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *OAuth2Client) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *OAuth2Client) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *OAuth2Client) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *OAuth2Client) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.

### GetUserinfoSignedResponseAlg

`func (o *OAuth2Client) GetUserinfoSignedResponseAlg() string`

GetUserinfoSignedResponseAlg returns the UserinfoSignedResponseAlg field if non-nil, zero value otherwise.

### GetUserinfoSignedResponseAlgOk

`func (o *OAuth2Client) GetUserinfoSignedResponseAlgOk() (*string, bool)`

GetUserinfoSignedResponseAlgOk returns a tuple with the UserinfoSignedResponseAlg field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserinfoSignedResponseAlg

`func (o *OAuth2Client) SetUserinfoSignedResponseAlg(v string)`

SetUserinfoSignedResponseAlg sets UserinfoSignedResponseAlg field to given value.

### HasUserinfoSignedResponseAlg

`func (o *OAuth2Client) HasUserinfoSignedResponseAlg() bool`

HasUserinfoSignedResponseAlg returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


