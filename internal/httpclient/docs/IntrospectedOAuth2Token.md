# IntrospectedOAuth2Token

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Active** | **bool** | Active is a boolean indicator of whether or not the presented token is currently active.  The specifics of a token&#39;s \&quot;active\&quot; state will vary depending on the implementation of the authorization server and the information it keeps about its tokens, but a \&quot;true\&quot; value return for the \&quot;active\&quot; property will generally indicate that a given token has been issued by this authorization server, has not been revoked by the resource owner, and is within its given time window of validity (e.g., after its issuance time and before its expiration time). | 
**Aud** | Pointer to **[]string** | Audience contains a list of the token&#39;s intended audiences. | [optional] 
**ClientId** | Pointer to **string** | ID is aclient identifier for the OAuth 2.0 client that requested this token. | [optional] 
**Exp** | Pointer to **int64** | Expires at is an integer timestamp, measured in the number of seconds since January 1 1970 UTC, indicating when this token will expire. | [optional] 
**Ext** | Pointer to **map[string]interface{}** | Extra is arbitrary data set by the session. | [optional] 
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

### NewIntrospectedOAuth2Token

`func NewIntrospectedOAuth2Token(active bool, ) *IntrospectedOAuth2Token`

NewIntrospectedOAuth2Token instantiates a new IntrospectedOAuth2Token object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewIntrospectedOAuth2TokenWithDefaults

`func NewIntrospectedOAuth2TokenWithDefaults() *IntrospectedOAuth2Token`

NewIntrospectedOAuth2TokenWithDefaults instantiates a new IntrospectedOAuth2Token object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetActive

`func (o *IntrospectedOAuth2Token) GetActive() bool`

GetActive returns the Active field if non-nil, zero value otherwise.

### GetActiveOk

`func (o *IntrospectedOAuth2Token) GetActiveOk() (*bool, bool)`

GetActiveOk returns a tuple with the Active field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActive

`func (o *IntrospectedOAuth2Token) SetActive(v bool)`

SetActive sets Active field to given value.


### GetAud

`func (o *IntrospectedOAuth2Token) GetAud() []string`

GetAud returns the Aud field if non-nil, zero value otherwise.

### GetAudOk

`func (o *IntrospectedOAuth2Token) GetAudOk() (*[]string, bool)`

GetAudOk returns a tuple with the Aud field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAud

`func (o *IntrospectedOAuth2Token) SetAud(v []string)`

SetAud sets Aud field to given value.

### HasAud

`func (o *IntrospectedOAuth2Token) HasAud() bool`

HasAud returns a boolean if a field has been set.

### GetClientId

`func (o *IntrospectedOAuth2Token) GetClientId() string`

GetClientId returns the ClientId field if non-nil, zero value otherwise.

### GetClientIdOk

`func (o *IntrospectedOAuth2Token) GetClientIdOk() (*string, bool)`

GetClientIdOk returns a tuple with the ClientId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClientId

`func (o *IntrospectedOAuth2Token) SetClientId(v string)`

SetClientId sets ClientId field to given value.

### HasClientId

`func (o *IntrospectedOAuth2Token) HasClientId() bool`

HasClientId returns a boolean if a field has been set.

### GetExp

`func (o *IntrospectedOAuth2Token) GetExp() int64`

GetExp returns the Exp field if non-nil, zero value otherwise.

### GetExpOk

`func (o *IntrospectedOAuth2Token) GetExpOk() (*int64, bool)`

GetExpOk returns a tuple with the Exp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExp

`func (o *IntrospectedOAuth2Token) SetExp(v int64)`

SetExp sets Exp field to given value.

### HasExp

`func (o *IntrospectedOAuth2Token) HasExp() bool`

HasExp returns a boolean if a field has been set.

### GetExt

`func (o *IntrospectedOAuth2Token) GetExt() map[string]interface{}`

GetExt returns the Ext field if non-nil, zero value otherwise.

### GetExtOk

`func (o *IntrospectedOAuth2Token) GetExtOk() (*map[string]interface{}, bool)`

GetExtOk returns a tuple with the Ext field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExt

`func (o *IntrospectedOAuth2Token) SetExt(v map[string]interface{})`

SetExt sets Ext field to given value.

### HasExt

`func (o *IntrospectedOAuth2Token) HasExt() bool`

HasExt returns a boolean if a field has been set.

### GetIat

`func (o *IntrospectedOAuth2Token) GetIat() int64`

GetIat returns the Iat field if non-nil, zero value otherwise.

### GetIatOk

`func (o *IntrospectedOAuth2Token) GetIatOk() (*int64, bool)`

GetIatOk returns a tuple with the Iat field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIat

`func (o *IntrospectedOAuth2Token) SetIat(v int64)`

SetIat sets Iat field to given value.

### HasIat

`func (o *IntrospectedOAuth2Token) HasIat() bool`

HasIat returns a boolean if a field has been set.

### GetIss

`func (o *IntrospectedOAuth2Token) GetIss() string`

GetIss returns the Iss field if non-nil, zero value otherwise.

### GetIssOk

`func (o *IntrospectedOAuth2Token) GetIssOk() (*string, bool)`

GetIssOk returns a tuple with the Iss field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIss

`func (o *IntrospectedOAuth2Token) SetIss(v string)`

SetIss sets Iss field to given value.

### HasIss

`func (o *IntrospectedOAuth2Token) HasIss() bool`

HasIss returns a boolean if a field has been set.

### GetNbf

`func (o *IntrospectedOAuth2Token) GetNbf() int64`

GetNbf returns the Nbf field if non-nil, zero value otherwise.

### GetNbfOk

`func (o *IntrospectedOAuth2Token) GetNbfOk() (*int64, bool)`

GetNbfOk returns a tuple with the Nbf field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNbf

`func (o *IntrospectedOAuth2Token) SetNbf(v int64)`

SetNbf sets Nbf field to given value.

### HasNbf

`func (o *IntrospectedOAuth2Token) HasNbf() bool`

HasNbf returns a boolean if a field has been set.

### GetObfuscatedSubject

`func (o *IntrospectedOAuth2Token) GetObfuscatedSubject() string`

GetObfuscatedSubject returns the ObfuscatedSubject field if non-nil, zero value otherwise.

### GetObfuscatedSubjectOk

`func (o *IntrospectedOAuth2Token) GetObfuscatedSubjectOk() (*string, bool)`

GetObfuscatedSubjectOk returns a tuple with the ObfuscatedSubject field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObfuscatedSubject

`func (o *IntrospectedOAuth2Token) SetObfuscatedSubject(v string)`

SetObfuscatedSubject sets ObfuscatedSubject field to given value.

### HasObfuscatedSubject

`func (o *IntrospectedOAuth2Token) HasObfuscatedSubject() bool`

HasObfuscatedSubject returns a boolean if a field has been set.

### GetScope

`func (o *IntrospectedOAuth2Token) GetScope() string`

GetScope returns the Scope field if non-nil, zero value otherwise.

### GetScopeOk

`func (o *IntrospectedOAuth2Token) GetScopeOk() (*string, bool)`

GetScopeOk returns a tuple with the Scope field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScope

`func (o *IntrospectedOAuth2Token) SetScope(v string)`

SetScope sets Scope field to given value.

### HasScope

`func (o *IntrospectedOAuth2Token) HasScope() bool`

HasScope returns a boolean if a field has been set.

### GetSub

`func (o *IntrospectedOAuth2Token) GetSub() string`

GetSub returns the Sub field if non-nil, zero value otherwise.

### GetSubOk

`func (o *IntrospectedOAuth2Token) GetSubOk() (*string, bool)`

GetSubOk returns a tuple with the Sub field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSub

`func (o *IntrospectedOAuth2Token) SetSub(v string)`

SetSub sets Sub field to given value.

### HasSub

`func (o *IntrospectedOAuth2Token) HasSub() bool`

HasSub returns a boolean if a field has been set.

### GetTokenType

`func (o *IntrospectedOAuth2Token) GetTokenType() string`

GetTokenType returns the TokenType field if non-nil, zero value otherwise.

### GetTokenTypeOk

`func (o *IntrospectedOAuth2Token) GetTokenTypeOk() (*string, bool)`

GetTokenTypeOk returns a tuple with the TokenType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTokenType

`func (o *IntrospectedOAuth2Token) SetTokenType(v string)`

SetTokenType sets TokenType field to given value.

### HasTokenType

`func (o *IntrospectedOAuth2Token) HasTokenType() bool`

HasTokenType returns a boolean if a field has been set.

### GetTokenUse

`func (o *IntrospectedOAuth2Token) GetTokenUse() string`

GetTokenUse returns the TokenUse field if non-nil, zero value otherwise.

### GetTokenUseOk

`func (o *IntrospectedOAuth2Token) GetTokenUseOk() (*string, bool)`

GetTokenUseOk returns a tuple with the TokenUse field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTokenUse

`func (o *IntrospectedOAuth2Token) SetTokenUse(v string)`

SetTokenUse sets TokenUse field to given value.

### HasTokenUse

`func (o *IntrospectedOAuth2Token) HasTokenUse() bool`

HasTokenUse returns a boolean if a field has been set.

### GetUsername

`func (o *IntrospectedOAuth2Token) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *IntrospectedOAuth2Token) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *IntrospectedOAuth2Token) SetUsername(v string)`

SetUsername sets Username field to given value.

### HasUsername

`func (o *IntrospectedOAuth2Token) HasUsername() bool`

HasUsername returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


