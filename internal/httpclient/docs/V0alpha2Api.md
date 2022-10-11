# \V0alpha2Api

All URIs are relative to _http://localhost_

| Method                                                                | HTTP request                     | Description                                                                              |
| --------------------------------------------------------------------- | -------------------------------- | ---------------------------------------------------------------------------------------- |
| [**DeleteOidcDynamicClient**](V0alpha2Api.md#DeleteOidcDynamicClient) | **Delete** /oauth2/register/{id} | Delete OAuth 2.0 Client using the OpenID Dynamic Client Registration Management Protocol |

## DeleteOidcDynamicClient

> DeleteOidcDynamicClient(ctx, id).Execute()

Delete OAuth 2.0 Client using the OpenID Dynamic Client Registration Management
Protocol

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
    resp, r, err := apiClient.V0alpha2Api.DeleteOidcDynamicClient(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V0alpha2Api.DeleteOidcDynamicClient``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters

| Name    | Type                | Description                                                                 | Notes |
| ------- | ------------------- | --------------------------------------------------------------------------- | ----- |
| **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc. |
| **id**  | **string**          | The id of the OAuth 2.0 Client.                                             |

### Other Parameters

Other parameters are passed through a pointer to a
apiDeleteOidcDynamicClientRequest struct via the builder pattern

| Name | Type | Description | Notes |
| ---- | ---- | ----------- | ----- |

### Return type

(empty response body)

### Authorization

[bearer](../README.md#bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#)
[[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)
