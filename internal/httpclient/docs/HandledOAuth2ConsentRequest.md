# HandledOAuth2ConsentRequest

## Properties

| Name           | Type       | Description                                                                                | Notes |
| -------------- | ---------- | ------------------------------------------------------------------------------------------ | ----- |
| **RedirectTo** | **string** | Original request URL to which you should redirect the user if request was already handled. |

## Methods

### NewHandledOAuth2ConsentRequest

`func NewHandledOAuth2ConsentRequest(redirectTo string, ) *HandledOAuth2ConsentRequest`

NewHandledOAuth2ConsentRequest instantiates a new HandledOAuth2ConsentRequest
object This constructor will assign default values to properties that have it
defined, and makes sure properties required by API are set, but the set of
arguments will change when the set of required properties is changed

### NewHandledOAuth2ConsentRequestWithDefaults

`func NewHandledOAuth2ConsentRequestWithDefaults() *HandledOAuth2ConsentRequest`

NewHandledOAuth2ConsentRequestWithDefaults instantiates a new
HandledOAuth2ConsentRequest object This constructor will only assign default
values to properties that have it defined, but it doesn't guarantee that
properties required by API are set

### GetRedirectTo

`func (o *HandledOAuth2ConsentRequest) GetRedirectTo() string`

GetRedirectTo returns the RedirectTo field if non-nil, zero value otherwise.

### GetRedirectToOk

`func (o *HandledOAuth2ConsentRequest) GetRedirectToOk() (*string, bool)`

GetRedirectToOk returns a tuple with the RedirectTo field if it's non-nil, zero
value otherwise and a boolean to check if the value has been set.

### SetRedirectTo

`func (o *HandledOAuth2ConsentRequest) SetRedirectTo(v string)`

SetRedirectTo sets RedirectTo field to given value.

[[Back to Model list]](../README.md#documentation-for-models)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to README]](../README.md)
