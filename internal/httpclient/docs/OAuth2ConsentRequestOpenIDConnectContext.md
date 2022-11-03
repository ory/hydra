# OAuth2ConsentRequestOpenIDConnectContext

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AcrValues** | Pointer to **[]string** | ACRValues is the Authentication AuthorizationContext Class Reference requested in the OAuth 2.0 Authorization request. It is a parameter defined by OpenID Connect and expresses which level of authentication (e.g. 2FA) is required.  OpenID Connect defines it as follows: &gt; Requested Authentication AuthorizationContext Class Reference values. Space-separated string that specifies the acr values that the Authorization Server is being requested to use for processing this Authentication Request, with the values appearing in order of preference. The Authentication AuthorizationContext Class satisfied by the authentication performed is returned as the acr Claim Value, as specified in Section 2. The acr Claim is requested as a Voluntary Claim by this parameter. | [optional] 
**Display** | Pointer to **string** | Display is a string value that specifies how the Authorization Server displays the authentication and consent user interface pages to the End-User. The defined values are: page: The Authorization Server SHOULD display the authentication and consent UI consistent with a full User Agent page view. If the display parameter is not specified, this is the default display mode. popup: The Authorization Server SHOULD display the authentication and consent UI consistent with a popup User Agent window. The popup User Agent window should be of an appropriate size for a login-focused dialog and should not obscure the entire window that it is popping up over. touch: The Authorization Server SHOULD display the authentication and consent UI consistent with a device that leverages a touch interface. wap: The Authorization Server SHOULD display the authentication and consent UI consistent with a \&quot;feature phone\&quot; type display.  The Authorization Server MAY also attempt to detect the capabilities of the User Agent and present an appropriate display. | [optional] 
**IdTokenHintClaims** | Pointer to **map[string]interface{}** | IDTokenHintClaims are the claims of the ID Token previously issued by the Authorization Server being passed as a hint about the End-User&#39;s current or past authenticated session with the Client. | [optional] 
**LoginHint** | Pointer to **string** | LoginHint hints about the login identifier the End-User might use to log in (if necessary). This hint can be used by an RP if it first asks the End-User for their e-mail address (or other identifier) and then wants to pass that value as a hint to the discovered authorization service. This value MAY also be a phone number in the format specified for the phone_number Claim. The use of this parameter is optional. | [optional] 
**UiLocales** | Pointer to **[]string** | UILocales is the End-User&#39;id preferred languages and scripts for the user interface, represented as a space-separated list of BCP47 [RFC5646] language tag values, ordered by preference. For instance, the value \&quot;fr-CA fr en\&quot; represents a preference for French as spoken in Canada, then French (without a region designation), followed by English (without a region designation). An error SHOULD NOT result if some or all of the requested locales are not supported by the OpenID Provider. | [optional] 

## Methods

### NewOAuth2ConsentRequestOpenIDConnectContext

`func NewOAuth2ConsentRequestOpenIDConnectContext() *OAuth2ConsentRequestOpenIDConnectContext`

NewOAuth2ConsentRequestOpenIDConnectContext instantiates a new OAuth2ConsentRequestOpenIDConnectContext object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOAuth2ConsentRequestOpenIDConnectContextWithDefaults

`func NewOAuth2ConsentRequestOpenIDConnectContextWithDefaults() *OAuth2ConsentRequestOpenIDConnectContext`

NewOAuth2ConsentRequestOpenIDConnectContextWithDefaults instantiates a new OAuth2ConsentRequestOpenIDConnectContext object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAcrValues

`func (o *OAuth2ConsentRequestOpenIDConnectContext) GetAcrValues() []string`

GetAcrValues returns the AcrValues field if non-nil, zero value otherwise.

### GetAcrValuesOk

`func (o *OAuth2ConsentRequestOpenIDConnectContext) GetAcrValuesOk() (*[]string, bool)`

GetAcrValuesOk returns a tuple with the AcrValues field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAcrValues

`func (o *OAuth2ConsentRequestOpenIDConnectContext) SetAcrValues(v []string)`

SetAcrValues sets AcrValues field to given value.

### HasAcrValues

`func (o *OAuth2ConsentRequestOpenIDConnectContext) HasAcrValues() bool`

HasAcrValues returns a boolean if a field has been set.

### GetDisplay

`func (o *OAuth2ConsentRequestOpenIDConnectContext) GetDisplay() string`

GetDisplay returns the Display field if non-nil, zero value otherwise.

### GetDisplayOk

`func (o *OAuth2ConsentRequestOpenIDConnectContext) GetDisplayOk() (*string, bool)`

GetDisplayOk returns a tuple with the Display field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisplay

`func (o *OAuth2ConsentRequestOpenIDConnectContext) SetDisplay(v string)`

SetDisplay sets Display field to given value.

### HasDisplay

`func (o *OAuth2ConsentRequestOpenIDConnectContext) HasDisplay() bool`

HasDisplay returns a boolean if a field has been set.

### GetIdTokenHintClaims

`func (o *OAuth2ConsentRequestOpenIDConnectContext) GetIdTokenHintClaims() map[string]interface{}`

GetIdTokenHintClaims returns the IdTokenHintClaims field if non-nil, zero value otherwise.

### GetIdTokenHintClaimsOk

`func (o *OAuth2ConsentRequestOpenIDConnectContext) GetIdTokenHintClaimsOk() (*map[string]interface{}, bool)`

GetIdTokenHintClaimsOk returns a tuple with the IdTokenHintClaims field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIdTokenHintClaims

`func (o *OAuth2ConsentRequestOpenIDConnectContext) SetIdTokenHintClaims(v map[string]interface{})`

SetIdTokenHintClaims sets IdTokenHintClaims field to given value.

### HasIdTokenHintClaims

`func (o *OAuth2ConsentRequestOpenIDConnectContext) HasIdTokenHintClaims() bool`

HasIdTokenHintClaims returns a boolean if a field has been set.

### GetLoginHint

`func (o *OAuth2ConsentRequestOpenIDConnectContext) GetLoginHint() string`

GetLoginHint returns the LoginHint field if non-nil, zero value otherwise.

### GetLoginHintOk

`func (o *OAuth2ConsentRequestOpenIDConnectContext) GetLoginHintOk() (*string, bool)`

GetLoginHintOk returns a tuple with the LoginHint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLoginHint

`func (o *OAuth2ConsentRequestOpenIDConnectContext) SetLoginHint(v string)`

SetLoginHint sets LoginHint field to given value.

### HasLoginHint

`func (o *OAuth2ConsentRequestOpenIDConnectContext) HasLoginHint() bool`

HasLoginHint returns a boolean if a field has been set.

### GetUiLocales

`func (o *OAuth2ConsentRequestOpenIDConnectContext) GetUiLocales() []string`

GetUiLocales returns the UiLocales field if non-nil, zero value otherwise.

### GetUiLocalesOk

`func (o *OAuth2ConsentRequestOpenIDConnectContext) GetUiLocalesOk() (*[]string, bool)`

GetUiLocalesOk returns a tuple with the UiLocales field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUiLocales

`func (o *OAuth2ConsentRequestOpenIDConnectContext) SetUiLocales(v []string)`

SetUiLocales sets UiLocales field to given value.

### HasUiLocales

`func (o *OAuth2ConsentRequestOpenIDConnectContext) HasUiLocales() bool`

HasUiLocales returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


