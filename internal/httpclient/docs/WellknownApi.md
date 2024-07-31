# \WellknownAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DiscoverJsonWebKeys**](WellknownAPI.md#DiscoverJsonWebKeys) | **Get** /.well-known/jwks.json | Discover Well-Known JSON Web Keys



## DiscoverJsonWebKeys

> JsonWebKeySet DiscoverJsonWebKeys(ctx).Execute()

Discover Well-Known JSON Web Keys



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
	resp, r, err := apiClient.WellknownAPI.DiscoverJsonWebKeys(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WellknownAPI.DiscoverJsonWebKeys``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DiscoverJsonWebKeys`: JsonWebKeySet
	fmt.Fprintf(os.Stdout, "Response from `WellknownAPI.DiscoverJsonWebKeys`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiDiscoverJsonWebKeysRequest struct via the builder pattern


### Return type

[**JsonWebKeySet**](JsonWebKeySet.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

