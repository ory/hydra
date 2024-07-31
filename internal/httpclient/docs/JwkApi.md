# \JwkAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateJsonWebKeySet**](JwkAPI.md#CreateJsonWebKeySet) | **Post** /admin/keys/{set} | Create JSON Web Key
[**DeleteJsonWebKey**](JwkAPI.md#DeleteJsonWebKey) | **Delete** /admin/keys/{set}/{kid} | Delete JSON Web Key
[**DeleteJsonWebKeySet**](JwkAPI.md#DeleteJsonWebKeySet) | **Delete** /admin/keys/{set} | Delete JSON Web Key Set
[**GetJsonWebKey**](JwkAPI.md#GetJsonWebKey) | **Get** /admin/keys/{set}/{kid} | Get JSON Web Key
[**GetJsonWebKeySet**](JwkAPI.md#GetJsonWebKeySet) | **Get** /admin/keys/{set} | Retrieve a JSON Web Key Set
[**SetJsonWebKey**](JwkAPI.md#SetJsonWebKey) | **Put** /admin/keys/{set}/{kid} | Set JSON Web Key
[**SetJsonWebKeySet**](JwkAPI.md#SetJsonWebKeySet) | **Put** /admin/keys/{set} | Update a JSON Web Key Set



## CreateJsonWebKeySet

> JsonWebKeySet CreateJsonWebKeySet(ctx, set).CreateJsonWebKeySet(createJsonWebKeySet).Execute()

Create JSON Web Key



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
	set := "set_example" // string | The JSON Web Key Set ID
	createJsonWebKeySet := *openapiclient.NewCreateJsonWebKeySet("Alg_example", "Kid_example", "Use_example") // CreateJsonWebKeySet | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.JwkAPI.CreateJsonWebKeySet(context.Background(), set).CreateJsonWebKeySet(createJsonWebKeySet).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `JwkAPI.CreateJsonWebKeySet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateJsonWebKeySet`: JsonWebKeySet
	fmt.Fprintf(os.Stdout, "Response from `JwkAPI.CreateJsonWebKeySet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**set** | **string** | The JSON Web Key Set ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateJsonWebKeySetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **createJsonWebKeySet** | [**CreateJsonWebKeySet**](CreateJsonWebKeySet.md) |  | 

### Return type

[**JsonWebKeySet**](JsonWebKeySet.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteJsonWebKey

> DeleteJsonWebKey(ctx, set, kid).Execute()

Delete JSON Web Key



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
	set := "set_example" // string | The JSON Web Key Set
	kid := "kid_example" // string | The JSON Web Key ID (kid)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.JwkAPI.DeleteJsonWebKey(context.Background(), set, kid).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `JwkAPI.DeleteJsonWebKey``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**set** | **string** | The JSON Web Key Set | 
**kid** | **string** | The JSON Web Key ID (kid) | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteJsonWebKeyRequest struct via the builder pattern


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


## DeleteJsonWebKeySet

> DeleteJsonWebKeySet(ctx, set).Execute()

Delete JSON Web Key Set



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
	set := "set_example" // string | The JSON Web Key Set

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.JwkAPI.DeleteJsonWebKeySet(context.Background(), set).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `JwkAPI.DeleteJsonWebKeySet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**set** | **string** | The JSON Web Key Set | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteJsonWebKeySetRequest struct via the builder pattern


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


## GetJsonWebKey

> JsonWebKeySet GetJsonWebKey(ctx, set, kid).Execute()

Get JSON Web Key



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
	set := "set_example" // string | JSON Web Key Set ID
	kid := "kid_example" // string | JSON Web Key ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.JwkAPI.GetJsonWebKey(context.Background(), set, kid).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `JwkAPI.GetJsonWebKey``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetJsonWebKey`: JsonWebKeySet
	fmt.Fprintf(os.Stdout, "Response from `JwkAPI.GetJsonWebKey`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**set** | **string** | JSON Web Key Set ID | 
**kid** | **string** | JSON Web Key ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetJsonWebKeyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



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


## GetJsonWebKeySet

> JsonWebKeySet GetJsonWebKeySet(ctx, set).Execute()

Retrieve a JSON Web Key Set



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
	set := "set_example" // string | JSON Web Key Set ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.JwkAPI.GetJsonWebKeySet(context.Background(), set).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `JwkAPI.GetJsonWebKeySet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetJsonWebKeySet`: JsonWebKeySet
	fmt.Fprintf(os.Stdout, "Response from `JwkAPI.GetJsonWebKeySet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**set** | **string** | JSON Web Key Set ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetJsonWebKeySetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


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


## SetJsonWebKey

> JsonWebKey SetJsonWebKey(ctx, set, kid).JsonWebKey(jsonWebKey).Execute()

Set JSON Web Key



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
	set := "set_example" // string | The JSON Web Key Set ID
	kid := "kid_example" // string | JSON Web Key ID
	jsonWebKey := *openapiclient.NewJsonWebKey("RS256", "1603dfe0af8f4596", "RSA", "sig") // JsonWebKey |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.JwkAPI.SetJsonWebKey(context.Background(), set, kid).JsonWebKey(jsonWebKey).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `JwkAPI.SetJsonWebKey``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `SetJsonWebKey`: JsonWebKey
	fmt.Fprintf(os.Stdout, "Response from `JwkAPI.SetJsonWebKey`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**set** | **string** | The JSON Web Key Set ID | 
**kid** | **string** | JSON Web Key ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiSetJsonWebKeyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **jsonWebKey** | [**JsonWebKey**](JsonWebKey.md) |  | 

### Return type

[**JsonWebKey**](JsonWebKey.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SetJsonWebKeySet

> JsonWebKeySet SetJsonWebKeySet(ctx, set).JsonWebKeySet(jsonWebKeySet).Execute()

Update a JSON Web Key Set



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
	set := "set_example" // string | The JSON Web Key Set ID
	jsonWebKeySet := *openapiclient.NewJsonWebKeySet() // JsonWebKeySet |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.JwkAPI.SetJsonWebKeySet(context.Background(), set).JsonWebKeySet(jsonWebKeySet).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `JwkAPI.SetJsonWebKeySet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `SetJsonWebKeySet`: JsonWebKeySet
	fmt.Fprintf(os.Stdout, "Response from `JwkAPI.SetJsonWebKeySet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**set** | **string** | The JSON Web Key Set ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiSetJsonWebKeySetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **jsonWebKeySet** | [**JsonWebKeySet**](JsonWebKeySet.md) |  | 

### Return type

[**JsonWebKeySet**](JsonWebKeySet.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

