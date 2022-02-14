# LoginRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Challenge** | **string** | ID is the identifier (\&quot;login challenge\&quot;) of the login request. It is used to identify the session. | 
**Client** | [**OAuth2Client**](OAuth2Client.md) |  | 
**OidcContext** | Pointer to [**OpenIDConnectContext**](OpenIDConnectContext.md) |  | [optional] 
**RequestUrl** | **string** | RequestURL is the original OAuth 2.0 Authorization URL requested by the OAuth 2.0 client. It is the URL which initiates the OAuth 2.0 Authorization Code or OAuth 2.0 Implicit flow. This URL is typically not needed, but might come in handy if you want to deal with additional request parameters. | 
**RequestedAccessTokenAudience** | **[]string** |  | 
**RequestedScope** | **[]string** |  | 
**SessionId** | Pointer to **string** | SessionID is the login session ID. If the user-agent reuses a login session (via cookie / remember flag) this ID will remain the same. If the user-agent did not have an existing authentication session (e.g. remember is false) this will be a new random value. This value is used as the \&quot;sid\&quot; parameter in the ID Token and in OIDC Front-/Back- channel logout. It&#39;s value can generally be used to associate consecutive login requests by a certain user. | [optional] 
**Skip** | **bool** | Skip, if true, implies that the client has requested the same scopes from the same user previously. If true, you can skip asking the user to grant the requested scopes, and simply forward the user to the redirect URL.  This feature allows you to update / set session information. | 
**Subject** | **string** | Subject is the user ID of the end-user that authenticated. Now, that end user needs to grant or deny the scope requested by the OAuth 2.0 client. If this value is set and &#x60;skip&#x60; is true, you MUST include this subject type when accepting the login request, or the request will fail. | 

## Methods

### NewLoginRequest

`func NewLoginRequest(challenge string, client OAuth2Client, requestUrl string, requestedAccessTokenAudience []string, requestedScope []string, skip bool, subject string, ) *LoginRequest`

NewLoginRequest instantiates a new LoginRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewLoginRequestWithDefaults

`func NewLoginRequestWithDefaults() *LoginRequest`

NewLoginRequestWithDefaults instantiates a new LoginRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetChallenge

`func (o *LoginRequest) GetChallenge() string`

GetChallenge returns the Challenge field if non-nil, zero value otherwise.

### GetChallengeOk

`func (o *LoginRequest) GetChallengeOk() (*string, bool)`

GetChallengeOk returns a tuple with the Challenge field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChallenge

`func (o *LoginRequest) SetChallenge(v string)`

SetChallenge sets Challenge field to given value.


### GetClient

`func (o *LoginRequest) GetClient() OAuth2Client`

GetClient returns the Client field if non-nil, zero value otherwise.

### GetClientOk

`func (o *LoginRequest) GetClientOk() (*OAuth2Client, bool)`

GetClientOk returns a tuple with the Client field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClient

`func (o *LoginRequest) SetClient(v OAuth2Client)`

SetClient sets Client field to given value.


### GetOidcContext

`func (o *LoginRequest) GetOidcContext() OpenIDConnectContext`

GetOidcContext returns the OidcContext field if non-nil, zero value otherwise.

### GetOidcContextOk

`func (o *LoginRequest) GetOidcContextOk() (*OpenIDConnectContext, bool)`

GetOidcContextOk returns a tuple with the OidcContext field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOidcContext

`func (o *LoginRequest) SetOidcContext(v OpenIDConnectContext)`

SetOidcContext sets OidcContext field to given value.

### HasOidcContext

`func (o *LoginRequest) HasOidcContext() bool`

HasOidcContext returns a boolean if a field has been set.

### GetRequestUrl

`func (o *LoginRequest) GetRequestUrl() string`

GetRequestUrl returns the RequestUrl field if non-nil, zero value otherwise.

### GetRequestUrlOk

`func (o *LoginRequest) GetRequestUrlOk() (*string, bool)`

GetRequestUrlOk returns a tuple with the RequestUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestUrl

`func (o *LoginRequest) SetRequestUrl(v string)`

SetRequestUrl sets RequestUrl field to given value.


### GetRequestedAccessTokenAudience

`func (o *LoginRequest) GetRequestedAccessTokenAudience() []string`

GetRequestedAccessTokenAudience returns the RequestedAccessTokenAudience field if non-nil, zero value otherwise.

### GetRequestedAccessTokenAudienceOk

`func (o *LoginRequest) GetRequestedAccessTokenAudienceOk() (*[]string, bool)`

GetRequestedAccessTokenAudienceOk returns a tuple with the RequestedAccessTokenAudience field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestedAccessTokenAudience

`func (o *LoginRequest) SetRequestedAccessTokenAudience(v []string)`

SetRequestedAccessTokenAudience sets RequestedAccessTokenAudience field to given value.


### GetRequestedScope

`func (o *LoginRequest) GetRequestedScope() []string`

GetRequestedScope returns the RequestedScope field if non-nil, zero value otherwise.

### GetRequestedScopeOk

`func (o *LoginRequest) GetRequestedScopeOk() (*[]string, bool)`

GetRequestedScopeOk returns a tuple with the RequestedScope field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestedScope

`func (o *LoginRequest) SetRequestedScope(v []string)`

SetRequestedScope sets RequestedScope field to given value.


### GetSessionId

`func (o *LoginRequest) GetSessionId() string`

GetSessionId returns the SessionId field if non-nil, zero value otherwise.

### GetSessionIdOk

`func (o *LoginRequest) GetSessionIdOk() (*string, bool)`

GetSessionIdOk returns a tuple with the SessionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSessionId

`func (o *LoginRequest) SetSessionId(v string)`

SetSessionId sets SessionId field to given value.

### HasSessionId

`func (o *LoginRequest) HasSessionId() bool`

HasSessionId returns a boolean if a field has been set.

### GetSkip

`func (o *LoginRequest) GetSkip() bool`

GetSkip returns the Skip field if non-nil, zero value otherwise.

### GetSkipOk

`func (o *LoginRequest) GetSkipOk() (*bool, bool)`

GetSkipOk returns a tuple with the Skip field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSkip

`func (o *LoginRequest) SetSkip(v bool)`

SetSkip sets Skip field to given value.


### GetSubject

`func (o *LoginRequest) GetSubject() string`

GetSubject returns the Subject field if non-nil, zero value otherwise.

### GetSubjectOk

`func (o *LoginRequest) GetSubjectOk() (*string, bool)`

GetSubjectOk returns a tuple with the Subject field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubject

`func (o *LoginRequest) SetSubject(v string)`

SetSubject sets Subject field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


