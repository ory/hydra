# \OidcAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateOidcDynamicClient**](OidcAPI.md#CreateOidcDynamicClient) | **Post** /oauth2/register | Register OAuth2 Client using OpenID Dynamic Client Registration
[**CreateVerifiableCredential**](OidcAPI.md#CreateVerifiableCredential) | **Post** /credentials | Issues a Verifiable Credential
[**DeleteOidcDynamicClient**](OidcAPI.md#DeleteOidcDynamicClient) | **Delete** /oauth2/register/{id} | Delete OAuth 2.0 Client using the OpenID Dynamic Client Registration Management Protocol
[**DiscoverOidcConfiguration**](OidcAPI.md#DiscoverOidcConfiguration) | **Get** /.well-known/openid-configuration | OpenID Connect Discovery
[**GetOidcDynamicClient**](OidcAPI.md#GetOidcDynamicClient) | **Get** /oauth2/register/{id} | Get OAuth2 Client using OpenID Dynamic Client Registration
[**GetOidcUserInfo**](OidcAPI.md#GetOidcUserInfo) | **Get** /userinfo | OpenID Connect Userinfo
[**RevokeOidcSession**](OidcAPI.md#RevokeOidcSession) | **Get** /oauth2/sessions/logout | OpenID Connect Front- and Back-channel Enabled Logout
[**SetOidcDynamicClient**](OidcAPI.md#SetOidcDynamicClient) | **Put** /oauth2/register/{id} | Set OAuth2 Client using OpenID Dynamic Client Registration



## CreateOidcDynamicClient

> OAuth2Client CreateOidcDynamicClient(ctx).OAuth2Client(oAuth2Client).Execute()

Register OAuth2 Client using OpenID Dynamic Client Registration



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/ory/hydra-client-go/v2"
)

func main() {
	oAuth2Client := *openapiclient.NewOAuth2Client() // OAuth2Client | Dynamic Client Registration Request Body

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OidcAPI.CreateOidcDynamicClient(context.Background()).OAuth2Client(oAuth2Client).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OidcAPI.CreateOidcDynamicClient``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateOidcDynamicClient`: OAuth2Client
	fmt.Fprintf(os.Stdout, "Response from `OidcAPI.CreateOidcDynamicClient`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateOidcDynamicClientRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **oAuth2Client** | [**OAuth2Client**](OAuth2Client.md) | Dynamic Client Registration Request Body | 

### Return type

[**OAuth2Client**](OAuth2Client.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateVerifiableCredential

> VerifiableCredentialResponse CreateVerifiableCredential(ctx).CreateVerifiableCredentialRequestBody(createVerifiableCredentialRequestBody).Execute()

Issues a Verifiable Credential



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/ory/hydra-client-go/v2"
)

func main() {
	createVerifiableCredentialRequestBody := *openapiclient.NewCreateVerifiableCredentialRequestBody() // CreateVerifiableCredentialRequestBody |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OidcAPI.CreateVerifiableCredential(context.Background()).CreateVerifiableCredentialRequestBody(createVerifiableCredentialRequestBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OidcAPI.CreateVerifiableCredential``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateVerifiableCredential`: VerifiableCredentialResponse
	fmt.Fprintf(os.Stdout, "Response from `OidcAPI.CreateVerifiableCredential`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateVerifiableCredentialRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **createVerifiableCredentialRequestBody** | [**CreateVerifiableCredentialRequestBody**](CreateVerifiableCredentialRequestBody.md) |  | 

### Return type

[**VerifiableCredentialResponse**](VerifiableCredentialResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteOidcDynamicClient

> DeleteOidcDynamicClient(ctx, id).Execute()

Delete OAuth 2.0 Client using the OpenID Dynamic Client Registration Management Protocol



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/ory/hydra-client-go/v2"
)

func main() {
	id := "id_example" // string | The id of the OAuth 2.0 Client.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.OidcAPI.DeleteOidcDynamicClient(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OidcAPI.DeleteOidcDynamicClient``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | The id of the OAuth 2.0 Client. | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteOidcDynamicClientRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DiscoverOidcConfiguration

> OidcConfiguration DiscoverOidcConfiguration(ctx).Execute()

OpenID Connect Discovery



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/ory/hydra-client-go/v2"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OidcAPI.DiscoverOidcConfiguration(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OidcAPI.DiscoverOidcConfiguration``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DiscoverOidcConfiguration`: OidcConfiguration
	fmt.Fprintf(os.Stdout, "Response from `OidcAPI.DiscoverOidcConfiguration`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiDiscoverOidcConfigurationRequest struct via the builder pattern


### Return type

[**OidcConfiguration**](OidcConfiguration.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetOidcDynamicClient

> OAuth2Client GetOidcDynamicClient(ctx, id).Execute()

Get OAuth2 Client using OpenID Dynamic Client Registration



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/ory/hydra-client-go/v2"
)

func main() {
	id := "id_example" // string | The id of the OAuth 2.0 Client.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OidcAPI.GetOidcDynamicClient(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OidcAPI.GetOidcDynamicClient``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetOidcDynamicClient`: OAuth2Client
	fmt.Fprintf(os.Stdout, "Response from `OidcAPI.GetOidcDynamicClient`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | The id of the OAuth 2.0 Client. | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetOidcDynamicClientRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**OAuth2Client**](OAuth2Client.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetOidcUserInfo

> OidcUserInfo GetOidcUserInfo(ctx).Execute()

OpenID Connect Userinfo



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/ory/hydra-client-go/v2"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OidcAPI.GetOidcUserInfo(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OidcAPI.GetOidcUserInfo``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetOidcUserInfo`: OidcUserInfo
	fmt.Fprintf(os.Stdout, "Response from `OidcAPI.GetOidcUserInfo`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetOidcUserInfoRequest struct via the builder pattern


### Return type

[**OidcUserInfo**](OidcUserInfo.md)

### Authorization

[oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RevokeOidcSession

> RevokeOidcSession(ctx).Execute()

OpenID Connect Front- and Back-channel Enabled Logout



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/ory/hydra-client-go/v2"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.OidcAPI.RevokeOidcSession(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OidcAPI.RevokeOidcSession``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiRevokeOidcSessionRequest struct via the builder pattern


### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SetOidcDynamicClient

> OAuth2Client SetOidcDynamicClient(ctx, id).OAuth2Client(oAuth2Client).Execute()

Set OAuth2 Client using OpenID Dynamic Client Registration



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/ory/hydra-client-go/v2"
)

func main() {
	id := "id_example" // string | OAuth 2.0 Client ID
	oAuth2Client := *openapiclient.NewOAuth2Client() // OAuth2Client | OAuth 2.0 Client Request Body

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OidcAPI.SetOidcDynamicClient(context.Background(), id).OAuth2Client(oAuth2Client).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OidcAPI.SetOidcDynamicClient``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `SetOidcDynamicClient`: OAuth2Client
	fmt.Fprintf(os.Stdout, "Response from `OidcAPI.SetOidcDynamicClient`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | OAuth 2.0 Client ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiSetOidcDynamicClientRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **oAuth2Client** | [**OAuth2Client**](OAuth2Client.md) | OAuth 2.0 Client Request Body | 

### Return type

[**OAuth2Client**](OAuth2Client.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

