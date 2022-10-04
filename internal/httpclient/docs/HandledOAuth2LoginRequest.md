# HandledOAuth2LoginRequest

## Properties

| Name           | Type       | Description                                                                                | Notes |
| -------------- | ---------- | ------------------------------------------------------------------------------------------ | ----- |
| **RedirectTo** | **string** | Original request URL to which you should redirect the user if request was already handled. |

## Methods

### NewHandledOAuth2LoginRequest

`func NewHandledOAuth2LoginRequest(redirectTo string, ) *HandledOAuth2LoginRequest`

NewHandledOAuth2LoginRequest instantiates a new HandledOAuth2LoginRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments will
change when the set of required properties is changed

### NewHandledOAuth2LoginRequestWithDefaults

`func NewHandledOAuth2LoginRequestWithDefaults() *HandledOAuth2LoginRequest`

NewHandledOAuth2LoginRequestWithDefaults instantiates a new
HandledOAuth2LoginRequest object This constructor will only assign default
values to properties that have it defined, but it doesn't guarantee that
properties required by API are set

### GetRedirectTo

`func (o *HandledOAuth2LoginRequest) GetRedirectTo() string`

GetRedirectTo returns the RedirectTo field if non-nil, zero value otherwise.

### GetRedirectToOk

`func (o *HandledOAuth2LoginRequest) GetRedirectToOk() (*string, bool)`

GetRedirectToOk returns a tuple with the RedirectTo field if it's non-nil, zero
value otherwise and a boolean to check if the value has been set.

### SetRedirectTo

`func (o *HandledOAuth2LoginRequest) SetRedirectTo(v string)`

SetRedirectTo sets RedirectTo field to given value.

[[Back to Model list]](../README.md#documentation-for-models)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to README]](../README.md)
