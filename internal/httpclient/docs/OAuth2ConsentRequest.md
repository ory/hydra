# OAuth2ConsentRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Acr** | Pointer to **string** | ACR represents the Authentication AuthorizationContext Class Reference value for this authentication session. You can use it to express that, for example, a user authenticated using two factor authentication. | [optional] 
**Amr** | Pointer to **[]string** |  | [optional] 
**Challenge** | **string** | Challenge is used to retrieve/accept/deny the consent request. | 
**Client** | Pointer to [**OAuth2Client**](OAuth2Client.md) |  | [optional] 
**ConsentRequestId** | Pointer to **string** | ConsentRequestID is the ID of the consent request. | [optional] 
**Context** | Pointer to **interface{}** |  | [optional] 
**LoginChallenge** | Pointer to **string** | LoginChallenge is the login challenge this consent challenge belongs to. It can be used to associate a login and consent request in the login &amp; consent app. | [optional] 
**LoginSessionId** | Pointer to **string** | LoginSessionID is the login session ID. If the user-agent reuses a login session (via cookie / remember flag) this ID will remain the same. If the user-agent did not have an existing authentication session (e.g. remember is false) this will be a new random value. This value is used as the \&quot;sid\&quot; parameter in the ID Token and in OIDC Front-/Back- channel logout. It&#39;s value can generally be used to associate consecutive login requests by a certain user. | [optional] 
**OidcContext** | Pointer to [**OAuth2ConsentRequestOpenIDConnectContext**](OAuth2ConsentRequestOpenIDConnectContext.md) |  | [optional] 
**RequestUrl** | Pointer to **string** | RequestURL is the original OAuth 2.0 Authorization URL requested by the OAuth 2.0 client. It is the URL which initiates the OAuth 2.0 Authorization Code or OAuth 2.0 Implicit flow. This URL is typically not needed, but might come in handy if you want to deal with additional request parameters. | [optional] 
**RequestedAccessTokenAudience** | Pointer to **[]string** |  | [optional] 
**RequestedScope** | Pointer to **[]string** |  | [optional] 
**Skip** | Pointer to **bool** | Skip, if true, implies that the client has requested the same scopes from the same user previously. If true, you must not ask the user to grant the requested scopes. You must however either allow or deny the consent request using the usual API call. | [optional] 
**Subject** | Pointer to **string** | Subject is the user ID of the end-user that authenticated. Now, that end user needs to grant or deny the scope requested by the OAuth 2.0 client. | [optional] 

## Methods

### NewOAuth2ConsentRequest

`func NewOAuth2ConsentRequest(challenge string, ) *OAuth2ConsentRequest`

NewOAuth2ConsentRequest instantiates a new OAuth2ConsentRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOAuth2ConsentRequestWithDefaults

`func NewOAuth2ConsentRequestWithDefaults() *OAuth2ConsentRequest`

NewOAuth2ConsentRequestWithDefaults instantiates a new OAuth2ConsentRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAcr

`func (o *OAuth2ConsentRequest) GetAcr() string`

GetAcr returns the Acr field if non-nil, zero value otherwise.

### GetAcrOk

`func (o *OAuth2ConsentRequest) GetAcrOk() (*string, bool)`

GetAcrOk returns a tuple with the Acr field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAcr

`func (o *OAuth2ConsentRequest) SetAcr(v string)`

SetAcr sets Acr field to given value.

### HasAcr

`func (o *OAuth2ConsentRequest) HasAcr() bool`

HasAcr returns a boolean if a field has been set.

### GetAmr

`func (o *OAuth2ConsentRequest) GetAmr() []string`

GetAmr returns the Amr field if non-nil, zero value otherwise.

### GetAmrOk

`func (o *OAuth2ConsentRequest) GetAmrOk() (*[]string, bool)`

GetAmrOk returns a tuple with the Amr field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAmr

`func (o *OAuth2ConsentRequest) SetAmr(v []string)`

SetAmr sets Amr field to given value.

### HasAmr

`func (o *OAuth2ConsentRequest) HasAmr() bool`

HasAmr returns a boolean if a field has been set.

### GetChallenge

`func (o *OAuth2ConsentRequest) GetChallenge() string`

GetChallenge returns the Challenge field if non-nil, zero value otherwise.

### GetChallengeOk

`func (o *OAuth2ConsentRequest) GetChallengeOk() (*string, bool)`

GetChallengeOk returns a tuple with the Challenge field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChallenge

`func (o *OAuth2ConsentRequest) SetChallenge(v string)`

SetChallenge sets Challenge field to given value.


### GetClient

`func (o *OAuth2ConsentRequest) GetClient() OAuth2Client`

GetClient returns the Client field if non-nil, zero value otherwise.

### GetClientOk

`func (o *OAuth2ConsentRequest) GetClientOk() (*OAuth2Client, bool)`

GetClientOk returns a tuple with the Client field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClient

`func (o *OAuth2ConsentRequest) SetClient(v OAuth2Client)`

SetClient sets Client field to given value.

### HasClient

`func (o *OAuth2ConsentRequest) HasClient() bool`

HasClient returns a boolean if a field has been set.

### GetConsentRequestId

`func (o *OAuth2ConsentRequest) GetConsentRequestId() string`

GetConsentRequestId returns the ConsentRequestId field if non-nil, zero value otherwise.

### GetConsentRequestIdOk

`func (o *OAuth2ConsentRequest) GetConsentRequestIdOk() (*string, bool)`

GetConsentRequestIdOk returns a tuple with the ConsentRequestId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConsentRequestId

`func (o *OAuth2ConsentRequest) SetConsentRequestId(v string)`

SetConsentRequestId sets ConsentRequestId field to given value.

### HasConsentRequestId

`func (o *OAuth2ConsentRequest) HasConsentRequestId() bool`

HasConsentRequestId returns a boolean if a field has been set.

### GetContext

`func (o *OAuth2ConsentRequest) GetContext() interface{}`

GetContext returns the Context field if non-nil, zero value otherwise.

### GetContextOk

`func (o *OAuth2ConsentRequest) GetContextOk() (*interface{}, bool)`

GetContextOk returns a tuple with the Context field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContext

`func (o *OAuth2ConsentRequest) SetContext(v interface{})`

SetContext sets Context field to given value.

### HasContext

`func (o *OAuth2ConsentRequest) HasContext() bool`

HasContext returns a boolean if a field has been set.

### SetContextNil

`func (o *OAuth2ConsentRequest) SetContextNil(b bool)`

 SetContextNil sets the value for Context to be an explicit nil

### UnsetContext
`func (o *OAuth2ConsentRequest) UnsetContext()`

UnsetContext ensures that no value is present for Context, not even an explicit nil
### GetLoginChallenge

`func (o *OAuth2ConsentRequest) GetLoginChallenge() string`

GetLoginChallenge returns the LoginChallenge field if non-nil, zero value otherwise.

### GetLoginChallengeOk

`func (o *OAuth2ConsentRequest) GetLoginChallengeOk() (*string, bool)`

GetLoginChallengeOk returns a tuple with the LoginChallenge field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLoginChallenge

`func (o *OAuth2ConsentRequest) SetLoginChallenge(v string)`

SetLoginChallenge sets LoginChallenge field to given value.

### HasLoginChallenge

`func (o *OAuth2ConsentRequest) HasLoginChallenge() bool`

HasLoginChallenge returns a boolean if a field has been set.

### GetLoginSessionId

`func (o *OAuth2ConsentRequest) GetLoginSessionId() string`

GetLoginSessionId returns the LoginSessionId field if non-nil, zero value otherwise.

### GetLoginSessionIdOk

`func (o *OAuth2ConsentRequest) GetLoginSessionIdOk() (*string, bool)`

GetLoginSessionIdOk returns a tuple with the LoginSessionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLoginSessionId

`func (o *OAuth2ConsentRequest) SetLoginSessionId(v string)`

SetLoginSessionId sets LoginSessionId field to given value.

### HasLoginSessionId

`func (o *OAuth2ConsentRequest) HasLoginSessionId() bool`

HasLoginSessionId returns a boolean if a field has been set.

### GetOidcContext

`func (o *OAuth2ConsentRequest) GetOidcContext() OAuth2ConsentRequestOpenIDConnectContext`

GetOidcContext returns the OidcContext field if non-nil, zero value otherwise.

### GetOidcContextOk

`func (o *OAuth2ConsentRequest) GetOidcContextOk() (*OAuth2ConsentRequestOpenIDConnectContext, bool)`

GetOidcContextOk returns a tuple with the OidcContext field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOidcContext

`func (o *OAuth2ConsentRequest) SetOidcContext(v OAuth2ConsentRequestOpenIDConnectContext)`

SetOidcContext sets OidcContext field to given value.

### HasOidcContext

`func (o *OAuth2ConsentRequest) HasOidcContext() bool`

HasOidcContext returns a boolean if a field has been set.

### GetRequestUrl

`func (o *OAuth2ConsentRequest) GetRequestUrl() string`

GetRequestUrl returns the RequestUrl field if non-nil, zero value otherwise.

### GetRequestUrlOk

`func (o *OAuth2ConsentRequest) GetRequestUrlOk() (*string, bool)`

GetRequestUrlOk returns a tuple with the RequestUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestUrl

`func (o *OAuth2ConsentRequest) SetRequestUrl(v string)`

SetRequestUrl sets RequestUrl field to given value.

### HasRequestUrl

`func (o *OAuth2ConsentRequest) HasRequestUrl() bool`

HasRequestUrl returns a boolean if a field has been set.

### GetRequestedAccessTokenAudience

`func (o *OAuth2ConsentRequest) GetRequestedAccessTokenAudience() []string`

GetRequestedAccessTokenAudience returns the RequestedAccessTokenAudience field if non-nil, zero value otherwise.

### GetRequestedAccessTokenAudienceOk

`func (o *OAuth2ConsentRequest) GetRequestedAccessTokenAudienceOk() (*[]string, bool)`

GetRequestedAccessTokenAudienceOk returns a tuple with the RequestedAccessTokenAudience field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestedAccessTokenAudience

`func (o *OAuth2ConsentRequest) SetRequestedAccessTokenAudience(v []string)`

SetRequestedAccessTokenAudience sets RequestedAccessTokenAudience field to given value.

### HasRequestedAccessTokenAudience

`func (o *OAuth2ConsentRequest) HasRequestedAccessTokenAudience() bool`

HasRequestedAccessTokenAudience returns a boolean if a field has been set.

### GetRequestedScope

`func (o *OAuth2ConsentRequest) GetRequestedScope() []string`

GetRequestedScope returns the RequestedScope field if non-nil, zero value otherwise.

### GetRequestedScopeOk

`func (o *OAuth2ConsentRequest) GetRequestedScopeOk() (*[]string, bool)`

GetRequestedScopeOk returns a tuple with the RequestedScope field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestedScope

`func (o *OAuth2ConsentRequest) SetRequestedScope(v []string)`

SetRequestedScope sets RequestedScope field to given value.

### HasRequestedScope

`func (o *OAuth2ConsentRequest) HasRequestedScope() bool`

HasRequestedScope returns a boolean if a field has been set.

### GetSkip

`func (o *OAuth2ConsentRequest) GetSkip() bool`

GetSkip returns the Skip field if non-nil, zero value otherwise.

### GetSkipOk

`func (o *OAuth2ConsentRequest) GetSkipOk() (*bool, bool)`

GetSkipOk returns a tuple with the Skip field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSkip

`func (o *OAuth2ConsentRequest) SetSkip(v bool)`

SetSkip sets Skip field to given value.

### HasSkip

`func (o *OAuth2ConsentRequest) HasSkip() bool`

HasSkip returns a boolean if a field has been set.

### GetSubject

`func (o *OAuth2ConsentRequest) GetSubject() string`

GetSubject returns the Subject field if non-nil, zero value otherwise.

### GetSubjectOk

`func (o *OAuth2ConsentRequest) GetSubjectOk() (*string, bool)`

GetSubjectOk returns a tuple with the Subject field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubject

`func (o *OAuth2ConsentRequest) SetSubject(v string)`

SetSubject sets Subject field to given value.

### HasSubject

`func (o *OAuth2ConsentRequest) HasSubject() bool`

HasSubject returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


