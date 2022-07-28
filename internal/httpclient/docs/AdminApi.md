# \AdminApi

All URIs are relative to _http://localhost_

| Method                                                                     | HTTP request                          | Description |
| -------------------------------------------------------------------------- | ------------------------------------- | ----------- |
| [**UpdateOAuth2ClientLifespans**](AdminApi.md#UpdateOAuth2ClientLifespans) | **Put** /admin/clients/{id}/lifespans |

## UpdateOAuth2ClientLifespans

> OAuth2Client UpdateOAuth2ClientLifespans(ctx,
> id).UpdateOAuth2ClientLifespans(updateOAuth2ClientLifespans).Execute()

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
    updateOAuth2ClientLifespans := *openapiclient.NewUpdateOAuth2ClientLifespans() // UpdateOAuth2ClientLifespans |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.AdminApi.UpdateOAuth2ClientLifespans(context.Background(), id).UpdateOAuth2ClientLifespans(updateOAuth2ClientLifespans).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `AdminApi.UpdateOAuth2ClientLifespans``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UpdateOAuth2ClientLifespans`: OAuth2Client
    fmt.Fprintf(os.Stdout, "Response from `AdminApi.UpdateOAuth2ClientLifespans`: %v\n", resp)
}
```

### Path Parameters

| Name    | Type                | Description                                                                 | Notes |
| ------- | ------------------- | --------------------------------------------------------------------------- | ----- |
| **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc. |
| **id**  | **string**          | The id of the OAuth 2.0 Client.                                             |

### Other Parameters

Other parameters are passed through a pointer to a
apiUpdateOAuth2ClientLifespansRequest struct via the builder pattern

| Name | Type | Description | Notes |
| ---- | ---- | ----------- | ----- |

**updateOAuth2ClientLifespans** |
[**UpdateOAuth2ClientLifespans**](UpdateOAuth2ClientLifespans.md) | |

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
