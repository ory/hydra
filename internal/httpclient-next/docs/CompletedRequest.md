# CompletedRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**RedirectTo** | **string** | RedirectURL is the URL which you should redirect the user to once the authentication process is completed. | 

## Methods

### NewCompletedRequest

`func NewCompletedRequest(redirectTo string, ) *CompletedRequest`

NewCompletedRequest instantiates a new CompletedRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCompletedRequestWithDefaults

`func NewCompletedRequestWithDefaults() *CompletedRequest`

NewCompletedRequestWithDefaults instantiates a new CompletedRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRedirectTo

`func (o *CompletedRequest) GetRedirectTo() string`

GetRedirectTo returns the RedirectTo field if non-nil, zero value otherwise.

### GetRedirectToOk

`func (o *CompletedRequest) GetRedirectToOk() (*string, bool)`

GetRedirectToOk returns a tuple with the RedirectTo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRedirectTo

`func (o *CompletedRequest) SetRedirectTo(v string)`

SetRedirectTo sets RedirectTo field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


