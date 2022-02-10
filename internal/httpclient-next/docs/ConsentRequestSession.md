# ConsentRequestSession

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AccessToken** | Pointer to **map[string]map[string]interface{}** | AccessToken sets session data for the access and refresh token, as well as any future tokens issued by the refresh grant. Keep in mind that this data will be available to anyone performing OAuth 2.0 Challenge Introspection. If only your services can perform OAuth 2.0 Challenge Introspection, this is usually fine. But if third parties can access that endpoint as well, sensitive data from the session might be exposed to them. Use with care! | [optional] 
**IdToken** | Pointer to **map[string]map[string]interface{}** | IDToken sets session data for the OpenID Connect ID token. Keep in mind that the session&#39;id payloads are readable by anyone that has access to the ID Challenge. Use with care! | [optional] 

## Methods

### NewConsentRequestSession

`func NewConsentRequestSession() *ConsentRequestSession`

NewConsentRequestSession instantiates a new ConsentRequestSession object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConsentRequestSessionWithDefaults

`func NewConsentRequestSessionWithDefaults() *ConsentRequestSession`

NewConsentRequestSessionWithDefaults instantiates a new ConsentRequestSession object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAccessToken

`func (o *ConsentRequestSession) GetAccessToken() map[string]map[string]interface{}`

GetAccessToken returns the AccessToken field if non-nil, zero value otherwise.

### GetAccessTokenOk

`func (o *ConsentRequestSession) GetAccessTokenOk() (*map[string]map[string]interface{}, bool)`

GetAccessTokenOk returns a tuple with the AccessToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccessToken

`func (o *ConsentRequestSession) SetAccessToken(v map[string]map[string]interface{})`

SetAccessToken sets AccessToken field to given value.

### HasAccessToken

`func (o *ConsentRequestSession) HasAccessToken() bool`

HasAccessToken returns a boolean if a field has been set.

### GetIdToken

`func (o *ConsentRequestSession) GetIdToken() map[string]map[string]interface{}`

GetIdToken returns the IdToken field if non-nil, zero value otherwise.

### GetIdTokenOk

`func (o *ConsentRequestSession) GetIdTokenOk() (*map[string]map[string]interface{}, bool)`

GetIdTokenOk returns a tuple with the IdToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIdToken

`func (o *ConsentRequestSession) SetIdToken(v map[string]map[string]interface{})`

SetIdToken sets IdToken field to given value.

### HasIdToken

`func (o *ConsentRequestSession) HasIdToken() bool`

HasIdToken returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


