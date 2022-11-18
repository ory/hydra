# OAuth2RedirectTo

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**RedirectTo** | **string** | RedirectURL is the URL which you should redirect the user&#39;s browser to once the authentication process is completed. | 

## Methods

### NewOAuth2RedirectTo

`func NewOAuth2RedirectTo(redirectTo string, ) *OAuth2RedirectTo`

NewOAuth2RedirectTo instantiates a new OAuth2RedirectTo object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOAuth2RedirectToWithDefaults

`func NewOAuth2RedirectToWithDefaults() *OAuth2RedirectTo`

NewOAuth2RedirectToWithDefaults instantiates a new OAuth2RedirectTo object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRedirectTo

`func (o *OAuth2RedirectTo) GetRedirectTo() string`

GetRedirectTo returns the RedirectTo field if non-nil, zero value otherwise.

### GetRedirectToOk

`func (o *OAuth2RedirectTo) GetRedirectToOk() (*string, bool)`

GetRedirectToOk returns a tuple with the RedirectTo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRedirectTo

`func (o *OAuth2RedirectTo) SetRedirectTo(v string)`

SetRedirectTo sets RedirectTo field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


