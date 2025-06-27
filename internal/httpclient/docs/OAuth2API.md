# \OAuth2API

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AcceptOAuth2ConsentRequest**](OAuth2API.md#AcceptOAuth2ConsentRequest) | **Put** /admin/oauth2/auth/requests/consent/accept | Accept OAuth 2.0 Consent Request
[**AcceptOAuth2LoginRequest**](OAuth2API.md#AcceptOAuth2LoginRequest) | **Put** /admin/oauth2/auth/requests/login/accept | Accept OAuth 2.0 Login Request
[**AcceptOAuth2LogoutRequest**](OAuth2API.md#AcceptOAuth2LogoutRequest) | **Put** /admin/oauth2/auth/requests/logout/accept | Accept OAuth 2.0 Session Logout Request
[**AcceptUserCodeRequest**](OAuth2API.md#AcceptUserCodeRequest) | **Put** /admin/oauth2/auth/requests/device/accept | Accepts a device grant user_code request
[**CreateOAuth2Client**](OAuth2API.md#CreateOAuth2Client) | **Post** /admin/clients | Create OAuth 2.0 Client
[**DeleteOAuth2Client**](OAuth2API.md#DeleteOAuth2Client) | **Delete** /admin/clients/{id} | Delete OAuth 2.0 Client
[**DeleteOAuth2Token**](OAuth2API.md#DeleteOAuth2Token) | **Delete** /admin/oauth2/tokens | Delete OAuth 2.0 Access Tokens from specific OAuth 2.0 Client
[**DeleteTrustedOAuth2JwtGrantIssuer**](OAuth2API.md#DeleteTrustedOAuth2JwtGrantIssuer) | **Delete** /admin/trust/grants/jwt-bearer/issuers/{id} | Delete Trusted OAuth2 JWT Bearer Grant Type Issuer
[**GetOAuth2Client**](OAuth2API.md#GetOAuth2Client) | **Get** /admin/clients/{id} | Get an OAuth 2.0 Client
[**GetOAuth2ConsentRequest**](OAuth2API.md#GetOAuth2ConsentRequest) | **Get** /admin/oauth2/auth/requests/consent | Get OAuth 2.0 Consent Request
[**GetOAuth2LoginRequest**](OAuth2API.md#GetOAuth2LoginRequest) | **Get** /admin/oauth2/auth/requests/login | Get OAuth 2.0 Login Request
[**GetOAuth2LogoutRequest**](OAuth2API.md#GetOAuth2LogoutRequest) | **Get** /admin/oauth2/auth/requests/logout | Get OAuth 2.0 Session Logout Request
[**GetTrustedOAuth2JwtGrantIssuer**](OAuth2API.md#GetTrustedOAuth2JwtGrantIssuer) | **Get** /admin/trust/grants/jwt-bearer/issuers/{id} | Get Trusted OAuth2 JWT Bearer Grant Type Issuer
[**IntrospectOAuth2Token**](OAuth2API.md#IntrospectOAuth2Token) | **Post** /admin/oauth2/introspect | Introspect OAuth2 Access and Refresh Tokens
[**ListOAuth2Clients**](OAuth2API.md#ListOAuth2Clients) | **Get** /admin/clients | List OAuth 2.0 Clients
[**ListOAuth2ConsentSessions**](OAuth2API.md#ListOAuth2ConsentSessions) | **Get** /admin/oauth2/auth/sessions/consent | List OAuth 2.0 Consent Sessions of a Subject
[**ListTrustedOAuth2JwtGrantIssuers**](OAuth2API.md#ListTrustedOAuth2JwtGrantIssuers) | **Get** /admin/trust/grants/jwt-bearer/issuers | List Trusted OAuth2 JWT Bearer Grant Type Issuers
[**OAuth2Authorize**](OAuth2API.md#OAuth2Authorize) | **Get** /oauth2/auth | OAuth 2.0 Authorize Endpoint
[**OAuth2DeviceFlow**](OAuth2API.md#OAuth2DeviceFlow) | **Post** /oauth2/device/auth | The OAuth 2.0 Device Authorize Endpoint
[**Oauth2TokenExchange**](OAuth2API.md#Oauth2TokenExchange) | **Post** /oauth2/token | The OAuth 2.0 Token Endpoint
[**PatchOAuth2Client**](OAuth2API.md#PatchOAuth2Client) | **Patch** /admin/clients/{id} | Patch OAuth 2.0 Client
[**PerformOAuth2DeviceVerificationFlow**](OAuth2API.md#PerformOAuth2DeviceVerificationFlow) | **Get** /oauth2/device/verify | OAuth 2.0 Device Verification Endpoint
[**RejectOAuth2ConsentRequest**](OAuth2API.md#RejectOAuth2ConsentRequest) | **Put** /admin/oauth2/auth/requests/consent/reject | Reject OAuth 2.0 Consent Request
[**RejectOAuth2LoginRequest**](OAuth2API.md#RejectOAuth2LoginRequest) | **Put** /admin/oauth2/auth/requests/login/reject | Reject OAuth 2.0 Login Request
[**RejectOAuth2LogoutRequest**](OAuth2API.md#RejectOAuth2LogoutRequest) | **Put** /admin/oauth2/auth/requests/logout/reject | Reject OAuth 2.0 Session Logout Request
[**RevokeOAuth2ConsentSessions**](OAuth2API.md#RevokeOAuth2ConsentSessions) | **Delete** /admin/oauth2/auth/sessions/consent | Revoke OAuth 2.0 Consent Sessions of a Subject
[**RevokeOAuth2LoginSessions**](OAuth2API.md#RevokeOAuth2LoginSessions) | **Delete** /admin/oauth2/auth/sessions/login | Revokes OAuth 2.0 Login Sessions by either a Subject or a SessionID
[**RevokeOAuth2Token**](OAuth2API.md#RevokeOAuth2Token) | **Post** /oauth2/revoke | Revoke OAuth 2.0 Access or Refresh Token
[**SetOAuth2Client**](OAuth2API.md#SetOAuth2Client) | **Put** /admin/clients/{id} | Set OAuth 2.0 Client
[**SetOAuth2ClientLifespans**](OAuth2API.md#SetOAuth2ClientLifespans) | **Put** /admin/clients/{id}/lifespans | Set OAuth2 Client Token Lifespans
[**TrustOAuth2JwtGrantIssuer**](OAuth2API.md#TrustOAuth2JwtGrantIssuer) | **Post** /admin/trust/grants/jwt-bearer/issuers | Trust OAuth2 JWT Bearer Grant Type Issuer



## AcceptOAuth2ConsentRequest

> OAuth2RedirectTo AcceptOAuth2ConsentRequest(ctx).ConsentChallenge(consentChallenge).AcceptOAuth2ConsentRequest(acceptOAuth2ConsentRequest).Execute()

Accept OAuth 2.0 Consent Request



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
	consentChallenge := "consentChallenge_example" // string | OAuth 2.0 Consent Request Challenge
	acceptOAuth2ConsentRequest := *openapiclient.NewAcceptOAuth2ConsentRequest() // AcceptOAuth2ConsentRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OAuth2API.AcceptOAuth2ConsentRequest(context.Background()).ConsentChallenge(consentChallenge).AcceptOAuth2ConsentRequest(acceptOAuth2ConsentRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.AcceptOAuth2ConsentRequest``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AcceptOAuth2ConsentRequest`: OAuth2RedirectTo
	fmt.Fprintf(os.Stdout, "Response from `OAuth2API.AcceptOAuth2ConsentRequest`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAcceptOAuth2ConsentRequestRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **consentChallenge** | **string** | OAuth 2.0 Consent Request Challenge | 
 **acceptOAuth2ConsentRequest** | [**AcceptOAuth2ConsentRequest**](AcceptOAuth2ConsentRequest.md) |  | 

### Return type

[**OAuth2RedirectTo**](OAuth2RedirectTo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AcceptOAuth2LoginRequest

> OAuth2RedirectTo AcceptOAuth2LoginRequest(ctx).LoginChallenge(loginChallenge).AcceptOAuth2LoginRequest(acceptOAuth2LoginRequest).Execute()

Accept OAuth 2.0 Login Request



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
	loginChallenge := "loginChallenge_example" // string | OAuth 2.0 Login Request Challenge
	acceptOAuth2LoginRequest := *openapiclient.NewAcceptOAuth2LoginRequest("Subject_example") // AcceptOAuth2LoginRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OAuth2API.AcceptOAuth2LoginRequest(context.Background()).LoginChallenge(loginChallenge).AcceptOAuth2LoginRequest(acceptOAuth2LoginRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.AcceptOAuth2LoginRequest``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AcceptOAuth2LoginRequest`: OAuth2RedirectTo
	fmt.Fprintf(os.Stdout, "Response from `OAuth2API.AcceptOAuth2LoginRequest`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAcceptOAuth2LoginRequestRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **loginChallenge** | **string** | OAuth 2.0 Login Request Challenge | 
 **acceptOAuth2LoginRequest** | [**AcceptOAuth2LoginRequest**](AcceptOAuth2LoginRequest.md) |  | 

### Return type

[**OAuth2RedirectTo**](OAuth2RedirectTo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AcceptOAuth2LogoutRequest

> OAuth2RedirectTo AcceptOAuth2LogoutRequest(ctx).LogoutChallenge(logoutChallenge).Execute()

Accept OAuth 2.0 Session Logout Request



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
	logoutChallenge := "logoutChallenge_example" // string | OAuth 2.0 Logout Request Challenge

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OAuth2API.AcceptOAuth2LogoutRequest(context.Background()).LogoutChallenge(logoutChallenge).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.AcceptOAuth2LogoutRequest``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AcceptOAuth2LogoutRequest`: OAuth2RedirectTo
	fmt.Fprintf(os.Stdout, "Response from `OAuth2API.AcceptOAuth2LogoutRequest`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAcceptOAuth2LogoutRequestRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **logoutChallenge** | **string** | OAuth 2.0 Logout Request Challenge | 

### Return type

[**OAuth2RedirectTo**](OAuth2RedirectTo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AcceptUserCodeRequest

> OAuth2RedirectTo AcceptUserCodeRequest(ctx).DeviceChallenge(deviceChallenge).AcceptDeviceUserCodeRequest(acceptDeviceUserCodeRequest).Execute()

Accepts a device grant user_code request



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
	deviceChallenge := "deviceChallenge_example" // string | 
	acceptDeviceUserCodeRequest := *openapiclient.NewAcceptDeviceUserCodeRequest() // AcceptDeviceUserCodeRequest |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OAuth2API.AcceptUserCodeRequest(context.Background()).DeviceChallenge(deviceChallenge).AcceptDeviceUserCodeRequest(acceptDeviceUserCodeRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.AcceptUserCodeRequest``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AcceptUserCodeRequest`: OAuth2RedirectTo
	fmt.Fprintf(os.Stdout, "Response from `OAuth2API.AcceptUserCodeRequest`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAcceptUserCodeRequestRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deviceChallenge** | **string** |  | 
 **acceptDeviceUserCodeRequest** | [**AcceptDeviceUserCodeRequest**](AcceptDeviceUserCodeRequest.md) |  | 

### Return type

[**OAuth2RedirectTo**](OAuth2RedirectTo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateOAuth2Client

> OAuth2Client CreateOAuth2Client(ctx).OAuth2Client(oAuth2Client).Execute()

Create OAuth 2.0 Client



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
	oAuth2Client := *openapiclient.NewOAuth2Client() // OAuth2Client | OAuth 2.0 Client Request Body

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OAuth2API.CreateOAuth2Client(context.Background()).OAuth2Client(oAuth2Client).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.CreateOAuth2Client``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateOAuth2Client`: OAuth2Client
	fmt.Fprintf(os.Stdout, "Response from `OAuth2API.CreateOAuth2Client`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateOAuth2ClientRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **oAuth2Client** | [**OAuth2Client**](OAuth2Client.md) | OAuth 2.0 Client Request Body | 

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


## DeleteOAuth2Client

> DeleteOAuth2Client(ctx, id).Execute()

Delete OAuth 2.0 Client



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
	r, err := apiClient.OAuth2API.DeleteOAuth2Client(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.DeleteOAuth2Client``: %v\n", err)
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

Other parameters are passed through a pointer to a apiDeleteOAuth2ClientRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteOAuth2Token

> DeleteOAuth2Token(ctx).ClientId(clientId).Execute()

Delete OAuth 2.0 Access Tokens from specific OAuth 2.0 Client



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
	clientId := "clientId_example" // string | OAuth 2.0 Client ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.OAuth2API.DeleteOAuth2Token(context.Background()).ClientId(clientId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.DeleteOAuth2Token``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDeleteOAuth2TokenRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **clientId** | **string** | OAuth 2.0 Client ID | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteTrustedOAuth2JwtGrantIssuer

> DeleteTrustedOAuth2JwtGrantIssuer(ctx, id).Execute()

Delete Trusted OAuth2 JWT Bearer Grant Type Issuer



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
	id := "id_example" // string | The id of the desired grant

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.OAuth2API.DeleteTrustedOAuth2JwtGrantIssuer(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.DeleteTrustedOAuth2JwtGrantIssuer``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | The id of the desired grant | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteTrustedOAuth2JwtGrantIssuerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetOAuth2Client

> OAuth2Client GetOAuth2Client(ctx, id).Execute()

Get an OAuth 2.0 Client



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
	resp, r, err := apiClient.OAuth2API.GetOAuth2Client(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.GetOAuth2Client``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetOAuth2Client`: OAuth2Client
	fmt.Fprintf(os.Stdout, "Response from `OAuth2API.GetOAuth2Client`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | The id of the OAuth 2.0 Client. | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetOAuth2ClientRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**OAuth2Client**](OAuth2Client.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetOAuth2ConsentRequest

> OAuth2ConsentRequest GetOAuth2ConsentRequest(ctx).ConsentChallenge(consentChallenge).Execute()

Get OAuth 2.0 Consent Request



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
	consentChallenge := "consentChallenge_example" // string | OAuth 2.0 Consent Request Challenge

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OAuth2API.GetOAuth2ConsentRequest(context.Background()).ConsentChallenge(consentChallenge).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.GetOAuth2ConsentRequest``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetOAuth2ConsentRequest`: OAuth2ConsentRequest
	fmt.Fprintf(os.Stdout, "Response from `OAuth2API.GetOAuth2ConsentRequest`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetOAuth2ConsentRequestRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **consentChallenge** | **string** | OAuth 2.0 Consent Request Challenge | 

### Return type

[**OAuth2ConsentRequest**](OAuth2ConsentRequest.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetOAuth2LoginRequest

> OAuth2LoginRequest GetOAuth2LoginRequest(ctx).LoginChallenge(loginChallenge).Execute()

Get OAuth 2.0 Login Request



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
	loginChallenge := "loginChallenge_example" // string | OAuth 2.0 Login Request Challenge

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OAuth2API.GetOAuth2LoginRequest(context.Background()).LoginChallenge(loginChallenge).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.GetOAuth2LoginRequest``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetOAuth2LoginRequest`: OAuth2LoginRequest
	fmt.Fprintf(os.Stdout, "Response from `OAuth2API.GetOAuth2LoginRequest`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetOAuth2LoginRequestRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **loginChallenge** | **string** | OAuth 2.0 Login Request Challenge | 

### Return type

[**OAuth2LoginRequest**](OAuth2LoginRequest.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetOAuth2LogoutRequest

> OAuth2LogoutRequest GetOAuth2LogoutRequest(ctx).LogoutChallenge(logoutChallenge).Execute()

Get OAuth 2.0 Session Logout Request



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
	logoutChallenge := "logoutChallenge_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OAuth2API.GetOAuth2LogoutRequest(context.Background()).LogoutChallenge(logoutChallenge).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.GetOAuth2LogoutRequest``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetOAuth2LogoutRequest`: OAuth2LogoutRequest
	fmt.Fprintf(os.Stdout, "Response from `OAuth2API.GetOAuth2LogoutRequest`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetOAuth2LogoutRequestRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **logoutChallenge** | **string** |  | 

### Return type

[**OAuth2LogoutRequest**](OAuth2LogoutRequest.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetTrustedOAuth2JwtGrantIssuer

> TrustedOAuth2JwtGrantIssuer GetTrustedOAuth2JwtGrantIssuer(ctx, id).Execute()

Get Trusted OAuth2 JWT Bearer Grant Type Issuer



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
	id := "id_example" // string | The id of the desired grant

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OAuth2API.GetTrustedOAuth2JwtGrantIssuer(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.GetTrustedOAuth2JwtGrantIssuer``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetTrustedOAuth2JwtGrantIssuer`: TrustedOAuth2JwtGrantIssuer
	fmt.Fprintf(os.Stdout, "Response from `OAuth2API.GetTrustedOAuth2JwtGrantIssuer`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | The id of the desired grant | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetTrustedOAuth2JwtGrantIssuerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**TrustedOAuth2JwtGrantIssuer**](TrustedOAuth2JwtGrantIssuer.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## IntrospectOAuth2Token

> IntrospectedOAuth2Token IntrospectOAuth2Token(ctx).Token(token).Scope(scope).Execute()

Introspect OAuth2 Access and Refresh Tokens



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
	token := "token_example" // string | The string value of the token. For access tokens, this is the \\\"access_token\\\" value returned from the token endpoint defined in OAuth 2.0. For refresh tokens, this is the \\\"refresh_token\\\" value returned.
	scope := "scope_example" // string | An optional, space separated list of required scopes. If the access token was not granted one of the scopes, the result of active will be false. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OAuth2API.IntrospectOAuth2Token(context.Background()).Token(token).Scope(scope).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.IntrospectOAuth2Token``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `IntrospectOAuth2Token`: IntrospectedOAuth2Token
	fmt.Fprintf(os.Stdout, "Response from `OAuth2API.IntrospectOAuth2Token`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiIntrospectOAuth2TokenRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **token** | **string** | The string value of the token. For access tokens, this is the \\\&quot;access_token\\\&quot; value returned from the token endpoint defined in OAuth 2.0. For refresh tokens, this is the \\\&quot;refresh_token\\\&quot; value returned. | 
 **scope** | **string** | An optional, space separated list of required scopes. If the access token was not granted one of the scopes, the result of active will be false. | 

### Return type

[**IntrospectedOAuth2Token**](IntrospectedOAuth2Token.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/x-www-form-urlencoded
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListOAuth2Clients

> []OAuth2Client ListOAuth2Clients(ctx).PageSize(pageSize).PageToken(pageToken).ClientName(clientName).Owner(owner).Execute()

List OAuth 2.0 Clients



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
	pageSize := int64(789) // int64 | Items per Page  This is the number of items per page to return. For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination). (optional) (default to 250)
	pageToken := "pageToken_example" // string | Next Page Token  The next page token. For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination). (optional)
	clientName := "clientName_example" // string | The name of the clients to filter by. (optional)
	owner := "owner_example" // string | The owner of the clients to filter by. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OAuth2API.ListOAuth2Clients(context.Background()).PageSize(pageSize).PageToken(pageToken).ClientName(clientName).Owner(owner).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.ListOAuth2Clients``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListOAuth2Clients`: []OAuth2Client
	fmt.Fprintf(os.Stdout, "Response from `OAuth2API.ListOAuth2Clients`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListOAuth2ClientsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pageSize** | **int64** | Items per Page  This is the number of items per page to return. For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination). | [default to 250]
 **pageToken** | **string** | Next Page Token  The next page token. For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination). | 
 **clientName** | **string** | The name of the clients to filter by. | 
 **owner** | **string** | The owner of the clients to filter by. | 

### Return type

[**[]OAuth2Client**](OAuth2Client.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListOAuth2ConsentSessions

> []OAuth2ConsentSession ListOAuth2ConsentSessions(ctx).Subject(subject).PageSize(pageSize).PageToken(pageToken).LoginSessionId(loginSessionId).Execute()

List OAuth 2.0 Consent Sessions of a Subject



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
	subject := "subject_example" // string | The subject to list the consent sessions for.
	pageSize := int64(789) // int64 | Items per Page  This is the number of items per page to return. For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination). (optional) (default to 250)
	pageToken := "pageToken_example" // string | Next Page Token  The next page token. For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination). (optional) (default to "1")
	loginSessionId := "loginSessionId_example" // string | The login session id to list the consent sessions for. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OAuth2API.ListOAuth2ConsentSessions(context.Background()).Subject(subject).PageSize(pageSize).PageToken(pageToken).LoginSessionId(loginSessionId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.ListOAuth2ConsentSessions``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListOAuth2ConsentSessions`: []OAuth2ConsentSession
	fmt.Fprintf(os.Stdout, "Response from `OAuth2API.ListOAuth2ConsentSessions`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListOAuth2ConsentSessionsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subject** | **string** | The subject to list the consent sessions for. | 
 **pageSize** | **int64** | Items per Page  This is the number of items per page to return. For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination). | [default to 250]
 **pageToken** | **string** | Next Page Token  The next page token. For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination). | [default to &quot;1&quot;]
 **loginSessionId** | **string** | The login session id to list the consent sessions for. | 

### Return type

[**[]OAuth2ConsentSession**](OAuth2ConsentSession.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListTrustedOAuth2JwtGrantIssuers

> []TrustedOAuth2JwtGrantIssuer ListTrustedOAuth2JwtGrantIssuers(ctx).PageSize(pageSize).PageToken(pageToken).Issuer(issuer).Execute()

List Trusted OAuth2 JWT Bearer Grant Type Issuers



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
	pageSize := int64(789) // int64 | Items per Page  This is the number of items per page to return. For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination). (optional) (default to 250)
	pageToken := "pageToken_example" // string | Next Page Token  The next page token. For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination). (optional)
	issuer := "issuer_example" // string | If optional \"issuer\" is supplied, only jwt-bearer grants with this issuer will be returned. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OAuth2API.ListTrustedOAuth2JwtGrantIssuers(context.Background()).PageSize(pageSize).PageToken(pageToken).Issuer(issuer).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.ListTrustedOAuth2JwtGrantIssuers``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListTrustedOAuth2JwtGrantIssuers`: []TrustedOAuth2JwtGrantIssuer
	fmt.Fprintf(os.Stdout, "Response from `OAuth2API.ListTrustedOAuth2JwtGrantIssuers`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListTrustedOAuth2JwtGrantIssuersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pageSize** | **int64** | Items per Page  This is the number of items per page to return. For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination). | [default to 250]
 **pageToken** | **string** | Next Page Token  The next page token. For details on pagination please head over to the [pagination documentation](https://www.ory.sh/docs/ecosystem/api-design#pagination). | 
 **issuer** | **string** | If optional \&quot;issuer\&quot; is supplied, only jwt-bearer grants with this issuer will be returned. | 

### Return type

[**[]TrustedOAuth2JwtGrantIssuer**](TrustedOAuth2JwtGrantIssuer.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## OAuth2Authorize

> ErrorOAuth2 OAuth2Authorize(ctx).Execute()

OAuth 2.0 Authorize Endpoint



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
	resp, r, err := apiClient.OAuth2API.OAuth2Authorize(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.OAuth2Authorize``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `OAuth2Authorize`: ErrorOAuth2
	fmt.Fprintf(os.Stdout, "Response from `OAuth2API.OAuth2Authorize`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiOAuth2AuthorizeRequest struct via the builder pattern


### Return type

[**ErrorOAuth2**](ErrorOAuth2.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## OAuth2DeviceFlow

> DeviceAuthorization OAuth2DeviceFlow(ctx).Execute()

The OAuth 2.0 Device Authorize Endpoint



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
	resp, r, err := apiClient.OAuth2API.OAuth2DeviceFlow(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.OAuth2DeviceFlow``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `OAuth2DeviceFlow`: DeviceAuthorization
	fmt.Fprintf(os.Stdout, "Response from `OAuth2API.OAuth2DeviceFlow`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiOAuth2DeviceFlowRequest struct via the builder pattern


### Return type

[**DeviceAuthorization**](DeviceAuthorization.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## Oauth2TokenExchange

> OAuth2TokenExchange Oauth2TokenExchange(ctx).GrantType(grantType).ClientId(clientId).Code(code).RedirectUri(redirectUri).RefreshToken(refreshToken).Execute()

The OAuth 2.0 Token Endpoint



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
	grantType := "grantType_example" // string | 
	clientId := "clientId_example" // string |  (optional)
	code := "code_example" // string |  (optional)
	redirectUri := "redirectUri_example" // string |  (optional)
	refreshToken := "refreshToken_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OAuth2API.Oauth2TokenExchange(context.Background()).GrantType(grantType).ClientId(clientId).Code(code).RedirectUri(redirectUri).RefreshToken(refreshToken).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.Oauth2TokenExchange``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `Oauth2TokenExchange`: OAuth2TokenExchange
	fmt.Fprintf(os.Stdout, "Response from `OAuth2API.Oauth2TokenExchange`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiOauth2TokenExchangeRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **grantType** | **string** |  | 
 **clientId** | **string** |  | 
 **code** | **string** |  | 
 **redirectUri** | **string** |  | 
 **refreshToken** | **string** |  | 

### Return type

[**OAuth2TokenExchange**](OAuth2TokenExchange.md)

### Authorization

[basic](../README.md#basic), [oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: application/x-www-form-urlencoded
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PatchOAuth2Client

> OAuth2Client PatchOAuth2Client(ctx, id).JsonPatch(jsonPatch).Execute()

Patch OAuth 2.0 Client



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
	jsonPatch := []openapiclient.JsonPatch{*openapiclient.NewJsonPatch("replace", "/name")} // []JsonPatch | OAuth 2.0 Client JSON Patch Body

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OAuth2API.PatchOAuth2Client(context.Background(), id).JsonPatch(jsonPatch).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.PatchOAuth2Client``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PatchOAuth2Client`: OAuth2Client
	fmt.Fprintf(os.Stdout, "Response from `OAuth2API.PatchOAuth2Client`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | The id of the OAuth 2.0 Client. | 

### Other Parameters

Other parameters are passed through a pointer to a apiPatchOAuth2ClientRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **jsonPatch** | [**[]JsonPatch**](JsonPatch.md) | OAuth 2.0 Client JSON Patch Body | 

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


## PerformOAuth2DeviceVerificationFlow

> ErrorOAuth2 PerformOAuth2DeviceVerificationFlow(ctx).Execute()

OAuth 2.0 Device Verification Endpoint



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
	resp, r, err := apiClient.OAuth2API.PerformOAuth2DeviceVerificationFlow(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.PerformOAuth2DeviceVerificationFlow``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PerformOAuth2DeviceVerificationFlow`: ErrorOAuth2
	fmt.Fprintf(os.Stdout, "Response from `OAuth2API.PerformOAuth2DeviceVerificationFlow`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiPerformOAuth2DeviceVerificationFlowRequest struct via the builder pattern


### Return type

[**ErrorOAuth2**](ErrorOAuth2.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RejectOAuth2ConsentRequest

> OAuth2RedirectTo RejectOAuth2ConsentRequest(ctx).ConsentChallenge(consentChallenge).RejectOAuth2Request(rejectOAuth2Request).Execute()

Reject OAuth 2.0 Consent Request



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
	consentChallenge := "consentChallenge_example" // string | OAuth 2.0 Consent Request Challenge
	rejectOAuth2Request := *openapiclient.NewRejectOAuth2Request() // RejectOAuth2Request |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OAuth2API.RejectOAuth2ConsentRequest(context.Background()).ConsentChallenge(consentChallenge).RejectOAuth2Request(rejectOAuth2Request).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.RejectOAuth2ConsentRequest``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RejectOAuth2ConsentRequest`: OAuth2RedirectTo
	fmt.Fprintf(os.Stdout, "Response from `OAuth2API.RejectOAuth2ConsentRequest`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRejectOAuth2ConsentRequestRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **consentChallenge** | **string** | OAuth 2.0 Consent Request Challenge | 
 **rejectOAuth2Request** | [**RejectOAuth2Request**](RejectOAuth2Request.md) |  | 

### Return type

[**OAuth2RedirectTo**](OAuth2RedirectTo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RejectOAuth2LoginRequest

> OAuth2RedirectTo RejectOAuth2LoginRequest(ctx).LoginChallenge(loginChallenge).RejectOAuth2Request(rejectOAuth2Request).Execute()

Reject OAuth 2.0 Login Request



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
	loginChallenge := "loginChallenge_example" // string | OAuth 2.0 Login Request Challenge
	rejectOAuth2Request := *openapiclient.NewRejectOAuth2Request() // RejectOAuth2Request |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OAuth2API.RejectOAuth2LoginRequest(context.Background()).LoginChallenge(loginChallenge).RejectOAuth2Request(rejectOAuth2Request).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.RejectOAuth2LoginRequest``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RejectOAuth2LoginRequest`: OAuth2RedirectTo
	fmt.Fprintf(os.Stdout, "Response from `OAuth2API.RejectOAuth2LoginRequest`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRejectOAuth2LoginRequestRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **loginChallenge** | **string** | OAuth 2.0 Login Request Challenge | 
 **rejectOAuth2Request** | [**RejectOAuth2Request**](RejectOAuth2Request.md) |  | 

### Return type

[**OAuth2RedirectTo**](OAuth2RedirectTo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RejectOAuth2LogoutRequest

> RejectOAuth2LogoutRequest(ctx).LogoutChallenge(logoutChallenge).Execute()

Reject OAuth 2.0 Session Logout Request



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
	logoutChallenge := "logoutChallenge_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.OAuth2API.RejectOAuth2LogoutRequest(context.Background()).LogoutChallenge(logoutChallenge).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.RejectOAuth2LogoutRequest``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRejectOAuth2LogoutRequestRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **logoutChallenge** | **string** |  | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RevokeOAuth2ConsentSessions

> RevokeOAuth2ConsentSessions(ctx).Subject(subject).Client(client).ConsentRequestId(consentRequestId).All(all).Execute()

Revoke OAuth 2.0 Consent Sessions of a Subject



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
	subject := "subject_example" // string | OAuth 2.0 Consent Subject  The subject whose consent sessions should be deleted. (optional)
	client := "client_example" // string | OAuth 2.0 Client ID  If set, deletes only those consent sessions that have been granted to the specified OAuth 2.0 Client ID. (optional)
	consentRequestId := "consentRequestId_example" // string | Consent Request ID  If set, revoke all token chains derived from this particular consent request ID. (optional)
	all := true // bool | Revoke All Consent Sessions  If set to `true` deletes all consent sessions by the Subject that have been granted. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.OAuth2API.RevokeOAuth2ConsentSessions(context.Background()).Subject(subject).Client(client).ConsentRequestId(consentRequestId).All(all).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.RevokeOAuth2ConsentSessions``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRevokeOAuth2ConsentSessionsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subject** | **string** | OAuth 2.0 Consent Subject  The subject whose consent sessions should be deleted. | 
 **client** | **string** | OAuth 2.0 Client ID  If set, deletes only those consent sessions that have been granted to the specified OAuth 2.0 Client ID. | 
 **consentRequestId** | **string** | Consent Request ID  If set, revoke all token chains derived from this particular consent request ID. | 
 **all** | **bool** | Revoke All Consent Sessions  If set to &#x60;true&#x60; deletes all consent sessions by the Subject that have been granted. | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RevokeOAuth2LoginSessions

> RevokeOAuth2LoginSessions(ctx).Subject(subject).Sid(sid).Execute()

Revokes OAuth 2.0 Login Sessions by either a Subject or a SessionID



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
	subject := "subject_example" // string | OAuth 2.0 Subject  The subject to revoke authentication sessions for. (optional)
	sid := "sid_example" // string | Login Session ID  The login session to revoke. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.OAuth2API.RevokeOAuth2LoginSessions(context.Background()).Subject(subject).Sid(sid).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.RevokeOAuth2LoginSessions``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRevokeOAuth2LoginSessionsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **subject** | **string** | OAuth 2.0 Subject  The subject to revoke authentication sessions for. | 
 **sid** | **string** | Login Session ID  The login session to revoke. | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RevokeOAuth2Token

> RevokeOAuth2Token(ctx).Token(token).ClientId(clientId).ClientSecret(clientSecret).Execute()

Revoke OAuth 2.0 Access or Refresh Token



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
	token := "token_example" // string | 
	clientId := "clientId_example" // string |  (optional)
	clientSecret := "clientSecret_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.OAuth2API.RevokeOAuth2Token(context.Background()).Token(token).ClientId(clientId).ClientSecret(clientSecret).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.RevokeOAuth2Token``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRevokeOAuth2TokenRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **token** | **string** |  | 
 **clientId** | **string** |  | 
 **clientSecret** | **string** |  | 

### Return type

 (empty response body)

### Authorization

[basic](../README.md#basic), [oauth2](../README.md#oauth2)

### HTTP request headers

- **Content-Type**: application/x-www-form-urlencoded
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SetOAuth2Client

> OAuth2Client SetOAuth2Client(ctx, id).OAuth2Client(oAuth2Client).Execute()

Set OAuth 2.0 Client



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
	resp, r, err := apiClient.OAuth2API.SetOAuth2Client(context.Background(), id).OAuth2Client(oAuth2Client).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.SetOAuth2Client``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `SetOAuth2Client`: OAuth2Client
	fmt.Fprintf(os.Stdout, "Response from `OAuth2API.SetOAuth2Client`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | OAuth 2.0 Client ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiSetOAuth2ClientRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **oAuth2Client** | [**OAuth2Client**](OAuth2Client.md) | OAuth 2.0 Client Request Body | 

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


## SetOAuth2ClientLifespans

> OAuth2Client SetOAuth2ClientLifespans(ctx, id).OAuth2ClientTokenLifespans(oAuth2ClientTokenLifespans).Execute()

Set OAuth2 Client Token Lifespans



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
	oAuth2ClientTokenLifespans := *openapiclient.NewOAuth2ClientTokenLifespans() // OAuth2ClientTokenLifespans |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OAuth2API.SetOAuth2ClientLifespans(context.Background(), id).OAuth2ClientTokenLifespans(oAuth2ClientTokenLifespans).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.SetOAuth2ClientLifespans``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `SetOAuth2ClientLifespans`: OAuth2Client
	fmt.Fprintf(os.Stdout, "Response from `OAuth2API.SetOAuth2ClientLifespans`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | OAuth 2.0 Client ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiSetOAuth2ClientLifespansRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **oAuth2ClientTokenLifespans** | [**OAuth2ClientTokenLifespans**](OAuth2ClientTokenLifespans.md) |  | 

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


## TrustOAuth2JwtGrantIssuer

> TrustedOAuth2JwtGrantIssuer TrustOAuth2JwtGrantIssuer(ctx).TrustOAuth2JwtGrantIssuer(trustOAuth2JwtGrantIssuer).Execute()

Trust OAuth2 JWT Bearer Grant Type Issuer



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
    "time"
	openapiclient "github.com/ory/hydra-client-go/v2"
)

func main() {
	trustOAuth2JwtGrantIssuer := *openapiclient.NewTrustOAuth2JwtGrantIssuer(time.Now(), "https://jwt-idp.example.com", *openapiclient.NewJsonWebKey("RS256", "1603dfe0af8f4596", "RSA", "sig"), []string{"Scope_example"}) // TrustOAuth2JwtGrantIssuer |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OAuth2API.TrustOAuth2JwtGrantIssuer(context.Background()).TrustOAuth2JwtGrantIssuer(trustOAuth2JwtGrantIssuer).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.TrustOAuth2JwtGrantIssuer``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `TrustOAuth2JwtGrantIssuer`: TrustedOAuth2JwtGrantIssuer
	fmt.Fprintf(os.Stdout, "Response from `OAuth2API.TrustOAuth2JwtGrantIssuer`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiTrustOAuth2JwtGrantIssuerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **trustOAuth2JwtGrantIssuer** | [**TrustOAuth2JwtGrantIssuer**](TrustOAuth2JwtGrantIssuer.md) |  | 

### Return type

[**TrustedOAuth2JwtGrantIssuer**](TrustedOAuth2JwtGrantIssuer.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

