# \OidcApi

All URIs are relative to _http://localhost_

| Method                                                            | HTTP request                  | Description                                                     |
| ----------------------------------------------------------------- | ----------------------------- | --------------------------------------------------------------- |
| [**CreateOidcDynamicClient**](OidcApi.md#CreateOidcDynamicClient) | **Post** /oauth2/register     | Register OAuth2 Client using OpenID Dynamic Client Registration |
| [**GetOidcDynamicClient**](OidcApi.md#GetOidcDynamicClient)       | **Get** /oauth2/register/{id} | Get OAuth2 Client using OpenID Dynamic Client Registration      |
| [**SetOidcDynamicClient**](OidcApi.md#SetOidcDynamicClient)       | **Put** /oauth2/register/{id} | Set OAuth2 Client using OpenID Dynamic Client Registration      |

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
    openapiclient "./openapi"
)

func main() {
    oAuth2Client := *openapiclient.NewOAuth2Client() // OAuth2Client | Dynamic Client Registration Request Body

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.OidcApi.CreateOidcDynamicClient(context.Background()).OAuth2Client(oAuth2Client).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OidcApi.CreateOidcDynamicClient``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateOidcDynamicClient`: OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `OidcApi.CreateOidcDynamicClient`: %v\n", resp)
}
```

### Path Parameters

### Other Parameters

Other parameters are passed through a pointer to a
apiCreateOidcDynamicClientRequest struct via the builder pattern

| Name             | Type                                | Description                              | Notes |
| ---------------- | ----------------------------------- | ---------------------------------------- | ----- |
| **oAuth2Client** | [**OAuth2Client**](OAuth2Client.md) | Dynamic Client Registration Request Body |

### Return type

[**OAuth2Client**](OAuth2Client.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
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
    openapiclient "./openapi"
)

func main() {
    id := "id_example" // string | The id of the OAuth 2.0 Client.

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.OidcApi.GetOidcDynamicClient(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OidcApi.GetOidcDynamicClient``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetOidcDynamicClient`: OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `OidcApi.GetOidcDynamicClient`: %v\n", resp)
}
```

### Path Parameters

| Name    | Type                | Description                                                                 | Notes |
| ------- | ------------------- | --------------------------------------------------------------------------- | ----- |
| **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc. |
| **id**  | **string**          | The id of the OAuth 2.0 Client.                                             |

### Other Parameters

Other parameters are passed through a pointer to a
apiGetOidcDynamicClientRequest struct via the builder pattern

| Name | Type | Description | Notes |
| ---- | ---- | ----------- | ----- |

### Return type

[**OAuth2Client**](OAuth2Client.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

## SetOidcDynamicClient

> OAuth2Client SetOidcDynamicClient(ctx,
> id).OAuth2Client(oAuth2Client).Execute()

Set OAuth2 Client using OpenID Dynamic Client Registration

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    id := "id_example" // string | OAuth 2.0 Client ID
    oAuth2Client := *openapiclient.NewOAuth2Client() // OAuth2Client | OAuth 2.0 Client Request Body

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.OidcApi.SetOidcDynamicClient(context.Background(), id).OAuth2Client(oAuth2Client).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `OidcApi.SetOidcDynamicClient``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `SetOidcDynamicClient`: OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `OidcApi.SetOidcDynamicClient`: %v\n", resp)
}
```

### Path Parameters

| Name    | Type                | Description                                                                 | Notes |
| ------- | ------------------- | --------------------------------------------------------------------------- | ----- |
| **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc. |
| **id**  | **string**          | OAuth 2.0 Client ID                                                         |

### Other Parameters

Other parameters are passed through a pointer to a
apiSetOidcDynamicClientRequest struct via the builder pattern

| Name | Type | Description | Notes |
| ---- | ---- | ----------- | ----- |

**oAuth2Client** | [**OAuth2Client**](OAuth2Client.md) | OAuth 2.0 Client
Request Body |

### Return type

[**OAuth2Client**](OAuth2Client.md)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)
