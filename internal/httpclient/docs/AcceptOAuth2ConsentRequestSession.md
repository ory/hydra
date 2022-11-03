# AcceptOAuth2ConsentRequestSession

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AccessToken** | Pointer to **interface{}** | AccessToken sets session data for the access and refresh token, as well as any future tokens issued by the refresh grant. Keep in mind that this data will be available to anyone performing OAuth 2.0 Challenge Introspection. If only your services can perform OAuth 2.0 Challenge Introspection, this is usually fine. But if third parties can access that endpoint as well, sensitive data from the session might be exposed to them. Use with care! | [optional] 
**IdToken** | Pointer to **interface{}** | IDToken sets session data for the OpenID Connect ID token. Keep in mind that the session&#39;id payloads are readable by anyone that has access to the ID Challenge. Use with care! | [optional] 

## Methods

### NewAcceptOAuth2ConsentRequestSession

`func NewAcceptOAuth2ConsentRequestSession() *AcceptOAuth2ConsentRequestSession`

NewAcceptOAuth2ConsentRequestSession instantiates a new AcceptOAuth2ConsentRequestSession object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAcceptOAuth2ConsentRequestSessionWithDefaults

`func NewAcceptOAuth2ConsentRequestSessionWithDefaults() *AcceptOAuth2ConsentRequestSession`

NewAcceptOAuth2ConsentRequestSessionWithDefaults instantiates a new AcceptOAuth2ConsentRequestSession object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAccessToken

`func (o *AcceptOAuth2ConsentRequestSession) GetAccessToken() interface{}`

GetAccessToken returns the AccessToken field if non-nil, zero value otherwise.

### GetAccessTokenOk

`func (o *AcceptOAuth2ConsentRequestSession) GetAccessTokenOk() (*interface{}, bool)`

GetAccessTokenOk returns a tuple with the AccessToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccessToken

`func (o *AcceptOAuth2ConsentRequestSession) SetAccessToken(v interface{})`

SetAccessToken sets AccessToken field to given value.

### HasAccessToken

`func (o *AcceptOAuth2ConsentRequestSession) HasAccessToken() bool`

HasAccessToken returns a boolean if a field has been set.

### SetAccessTokenNil

`func (o *AcceptOAuth2ConsentRequestSession) SetAccessTokenNil(b bool)`

 SetAccessTokenNil sets the value for AccessToken to be an explicit nil

### UnsetAccessToken
`func (o *AcceptOAuth2ConsentRequestSession) UnsetAccessToken()`

UnsetAccessToken ensures that no value is present for AccessToken, not even an explicit nil
### GetIdToken

`func (o *AcceptOAuth2ConsentRequestSession) GetIdToken() interface{}`

GetIdToken returns the IdToken field if non-nil, zero value otherwise.

### GetIdTokenOk

`func (o *AcceptOAuth2ConsentRequestSession) GetIdTokenOk() (*interface{}, bool)`

GetIdTokenOk returns a tuple with the IdToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIdToken

`func (o *AcceptOAuth2ConsentRequestSession) SetIdToken(v interface{})`

SetIdToken sets IdToken field to given value.

### HasIdToken

`func (o *AcceptOAuth2ConsentRequestSession) HasIdToken() bool`

HasIdToken returns a boolean if a field has been set.

### SetIdTokenNil

`func (o *AcceptOAuth2ConsentRequestSession) SetIdTokenNil(b bool)`

 SetIdTokenNil sets the value for IdToken to be an explicit nil

### UnsetIdToken
`func (o *AcceptOAuth2ConsentRequestSession) UnsetIdToken()`

UnsetIdToken ensures that no value is present for IdToken, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


