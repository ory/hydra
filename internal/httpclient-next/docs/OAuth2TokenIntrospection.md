# OAuth2TokenIntrospection

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Active** | **bool** | Active is a boolean indicator of whether or not the presented token is currently active.  The specifics of a token&#39;s \&quot;active\&quot; state will vary depending on the implementation of the authorization server and the information it keeps about its tokens, but a \&quot;true\&quot; value return for the \&quot;active\&quot; property will generally indicate that a given token has been issued by this authorization server, has not been revoked by the resource owner, and is within its given time window of validity (e.g., after its issuance time and before its expiration time). | 
**Aud** | Pointer to **[]string** | Audience contains a list of the token&#39;s intended audiences. | [optional] 
**ClientId** | Pointer to **string** | ID is aclient identifier for the OAuth 2.0 client that requested this token. | [optional] 
**Exp** | Pointer to **int64** | Expires at is an integer timestamp, measured in the number of seconds since January 1 1970 UTC, indicating when this token will expire. | [optional] 
**Ext** | Pointer to **map[string]map[string]interface{}** | Extra is arbitrary data set by the session. | [optional] 
**Iat** | Pointer to **int64** | Issued at is an integer timestamp, measured in the number of seconds since January 1 1970 UTC, indicating when this token was originally issued. | [optional] 
**Iss** | Pointer to **string** | IssuerURL is a string representing the issuer of this token | [optional] 
**Nbf** | Pointer to **int64** | NotBefore is an integer timestamp, measured in the number of seconds since January 1 1970 UTC, indicating when this token is not to be used before. | [optional] 
**ObfuscatedSubject** | Pointer to **string** | ObfuscatedSubject is set when the subject identifier algorithm was set to \&quot;pairwise\&quot; during authorization. It is the &#x60;sub&#x60; value of the ID Token that was issued. | [optional] 
**Scope** | Pointer to **string** | Scope is a JSON string containing a space-separated list of scopes associated with this token. | [optional] 
**Sub** | Pointer to **string** | Subject of the token, as defined in JWT [RFC7519]. Usually a machine-readable identifier of the resource owner who authorized this token. | [optional] 
**TokenType** | Pointer to **string** | TokenType is the introspected token&#39;s type, typically &#x60;Bearer&#x60;. | [optional] 
**TokenUse** | Pointer to **string** | TokenUse is the introspected token&#39;s use, for example &#x60;access_token&#x60; or &#x60;refresh_token&#x60;. | [optional] 
**Username** | Pointer to **string** | Username is a human-readable identifier for the resource owner who authorized this token. | [optional] 

## Methods

### NewOAuth2TokenIntrospection

`func NewOAuth2TokenIntrospection(active bool, ) *OAuth2TokenIntrospection`

NewOAuth2TokenIntrospection instantiates a new OAuth2TokenIntrospection object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOAuth2TokenIntrospectionWithDefaults

`func NewOAuth2TokenIntrospectionWithDefaults() *OAuth2TokenIntrospection`

NewOAuth2TokenIntrospectionWithDefaults instantiates a new OAuth2TokenIntrospection object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetActive

`func (o *OAuth2TokenIntrospection) GetActive() bool`

GetActive returns the Active field if non-nil, zero value otherwise.

### GetActiveOk

`func (o *OAuth2TokenIntrospection) GetActiveOk() (*bool, bool)`

GetActiveOk returns a tuple with the Active field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActive

`func (o *OAuth2TokenIntrospection) SetActive(v bool)`

SetActive sets Active field to given value.


### GetAud

`func (o *OAuth2TokenIntrospection) GetAud() []string`

GetAud returns the Aud field if non-nil, zero value otherwise.

### GetAudOk

`func (o *OAuth2TokenIntrospection) GetAudOk() (*[]string, bool)`

GetAudOk returns a tuple with the Aud field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAud

`func (o *OAuth2TokenIntrospection) SetAud(v []string)`

SetAud sets Aud field to given value.

### HasAud

`func (o *OAuth2TokenIntrospection) HasAud() bool`

HasAud returns a boolean if a field has been set.

### GetClientId

`func (o *OAuth2TokenIntrospection) GetClientId() string`

GetClientId returns the ClientId field if non-nil, zero value otherwise.

### GetClientIdOk

`func (o *OAuth2TokenIntrospection) GetClientIdOk() (*string, bool)`

GetClientIdOk returns a tuple with the ClientId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClientId

`func (o *OAuth2TokenIntrospection) SetClientId(v string)`

SetClientId sets ClientId field to given value.

### HasClientId

`func (o *OAuth2TokenIntrospection) HasClientId() bool`

HasClientId returns a boolean if a field has been set.

### GetExp

`func (o *OAuth2TokenIntrospection) GetExp() int64`

GetExp returns the Exp field if non-nil, zero value otherwise.

### GetExpOk

`func (o *OAuth2TokenIntrospection) GetExpOk() (*int64, bool)`

GetExpOk returns a tuple with the Exp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExp

`func (o *OAuth2TokenIntrospection) SetExp(v int64)`

SetExp sets Exp field to given value.

### HasExp

`func (o *OAuth2TokenIntrospection) HasExp() bool`

HasExp returns a boolean if a field has been set.

### GetExt

`func (o *OAuth2TokenIntrospection) GetExt() map[string]map[string]interface{}`

GetExt returns the Ext field if non-nil, zero value otherwise.

### GetExtOk

`func (o *OAuth2TokenIntrospection) GetExtOk() (*map[string]map[string]interface{}, bool)`

GetExtOk returns a tuple with the Ext field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExt

`func (o *OAuth2TokenIntrospection) SetExt(v map[string]map[string]interface{})`

SetExt sets Ext field to given value.

### HasExt

`func (o *OAuth2TokenIntrospection) HasExt() bool`

HasExt returns a boolean if a field has been set.

### GetIat

`func (o *OAuth2TokenIntrospection) GetIat() int64`

GetIat returns the Iat field if non-nil, zero value otherwise.

### GetIatOk

`func (o *OAuth2TokenIntrospection) GetIatOk() (*int64, bool)`

GetIatOk returns a tuple with the Iat field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIat

`func (o *OAuth2TokenIntrospection) SetIat(v int64)`

SetIat sets Iat field to given value.

### HasIat

`func (o *OAuth2TokenIntrospection) HasIat() bool`

HasIat returns a boolean if a field has been set.

### GetIss

`func (o *OAuth2TokenIntrospection) GetIss() string`

GetIss returns the Iss field if non-nil, zero value otherwise.

### GetIssOk

`func (o *OAuth2TokenIntrospection) GetIssOk() (*string, bool)`

GetIssOk returns a tuple with the Iss field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIss

`func (o *OAuth2TokenIntrospection) SetIss(v string)`

SetIss sets Iss field to given value.

### HasIss

`func (o *OAuth2TokenIntrospection) HasIss() bool`

HasIss returns a boolean if a field has been set.

### GetNbf

`func (o *OAuth2TokenIntrospection) GetNbf() int64`

GetNbf returns the Nbf field if non-nil, zero value otherwise.

### GetNbfOk

`func (o *OAuth2TokenIntrospection) GetNbfOk() (*int64, bool)`

GetNbfOk returns a tuple with the Nbf field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNbf

`func (o *OAuth2TokenIntrospection) SetNbf(v int64)`

SetNbf sets Nbf field to given value.

### HasNbf

`func (o *OAuth2TokenIntrospection) HasNbf() bool`

HasNbf returns a boolean if a field has been set.

### GetObfuscatedSubject

`func (o *OAuth2TokenIntrospection) GetObfuscatedSubject() string`

GetObfuscatedSubject returns the ObfuscatedSubject field if non-nil, zero value otherwise.

### GetObfuscatedSubjectOk

`func (o *OAuth2TokenIntrospection) GetObfuscatedSubjectOk() (*string, bool)`

GetObfuscatedSubjectOk returns a tuple with the ObfuscatedSubject field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObfuscatedSubject

`func (o *OAuth2TokenIntrospection) SetObfuscatedSubject(v string)`

SetObfuscatedSubject sets ObfuscatedSubject field to given value.

### HasObfuscatedSubject

`func (o *OAuth2TokenIntrospection) HasObfuscatedSubject() bool`

HasObfuscatedSubject returns a boolean if a field has been set.

### GetScope

`func (o *OAuth2TokenIntrospection) GetScope() string`

GetScope returns the Scope field if non-nil, zero value otherwise.

### GetScopeOk

`func (o *OAuth2TokenIntrospection) GetScopeOk() (*string, bool)`

GetScopeOk returns a tuple with the Scope field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScope

`func (o *OAuth2TokenIntrospection) SetScope(v string)`

SetScope sets Scope field to given value.

### HasScope

`func (o *OAuth2TokenIntrospection) HasScope() bool`

HasScope returns a boolean if a field has been set.

### GetSub

`func (o *OAuth2TokenIntrospection) GetSub() string`

GetSub returns the Sub field if non-nil, zero value otherwise.

### GetSubOk

`func (o *OAuth2TokenIntrospection) GetSubOk() (*string, bool)`

GetSubOk returns a tuple with the Sub field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSub

`func (o *OAuth2TokenIntrospection) SetSub(v string)`

SetSub sets Sub field to given value.

### HasSub

`func (o *OAuth2TokenIntrospection) HasSub() bool`

HasSub returns a boolean if a field has been set.

### GetTokenType

`func (o *OAuth2TokenIntrospection) GetTokenType() string`

GetTokenType returns the TokenType field if non-nil, zero value otherwise.

### GetTokenTypeOk

`func (o *OAuth2TokenIntrospection) GetTokenTypeOk() (*string, bool)`

GetTokenTypeOk returns a tuple with the TokenType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTokenType

`func (o *OAuth2TokenIntrospection) SetTokenType(v string)`

SetTokenType sets TokenType field to given value.

### HasTokenType

`func (o *OAuth2TokenIntrospection) HasTokenType() bool`

HasTokenType returns a boolean if a field has been set.

### GetTokenUse

`func (o *OAuth2TokenIntrospection) GetTokenUse() string`

GetTokenUse returns the TokenUse field if non-nil, zero value otherwise.

### GetTokenUseOk

`func (o *OAuth2TokenIntrospection) GetTokenUseOk() (*string, bool)`

GetTokenUseOk returns a tuple with the TokenUse field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTokenUse

`func (o *OAuth2TokenIntrospection) SetTokenUse(v string)`

SetTokenUse sets TokenUse field to given value.

### HasTokenUse

`func (o *OAuth2TokenIntrospection) HasTokenUse() bool`

HasTokenUse returns a boolean if a field has been set.

### GetUsername

`func (o *OAuth2TokenIntrospection) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *OAuth2TokenIntrospection) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *OAuth2TokenIntrospection) SetUsername(v string)`

SetUsername sets Username field to given value.

### HasUsername

`func (o *OAuth2TokenIntrospection) HasUsername() bool`

HasUsername returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


