# JsonWebKeyApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createJsonWebKeySet**](JsonWebKeyApi.md#createJsonWebKeySet) | **POST** /keys/{set} | Generate a new JSON Web Key
[**deleteJsonWebKey**](JsonWebKeyApi.md#deleteJsonWebKey) | **DELETE** /keys/{set}/{kid} | Delete a JSON Web Key
[**deleteJsonWebKeySet**](JsonWebKeyApi.md#deleteJsonWebKeySet) | **DELETE** /keys/{set} | Delete a JSON Web Key Set
[**getJsonWebKey**](JsonWebKeyApi.md#getJsonWebKey) | **GET** /keys/{set}/{kid} | Retrieve a JSON Web Key
[**getJsonWebKeySet**](JsonWebKeyApi.md#getJsonWebKeySet) | **GET** /keys/{set} | Retrieve a JSON Web Key Set
[**updateJsonWebKey**](JsonWebKeyApi.md#updateJsonWebKey) | **PUT** /keys/{set}/{kid} | Update a JSON Web Key
[**updateJsonWebKeySet**](JsonWebKeyApi.md#updateJsonWebKeySet) | **PUT** /keys/{set} | Update a JSON Web Key Set


<a name="createJsonWebKeySet"></a>
# **createJsonWebKeySet**
> JSONWebKeySet createJsonWebKeySet(set, body)

Generate a new JSON Web Key

This endpoint is capable of generating JSON Web Key Sets for you. There a different strategies available, such as symmetric cryptographic keys (HS256, HS512) and asymetric cryptographic keys (RS256, ECDSA). If the specified JSON Web Key Set does not exist, it will be created.  A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.JsonWebKeyApi;


JsonWebKeyApi apiInstance = new JsonWebKeyApi();
String set = "set_example"; // String | The set
JsonWebKeySetGeneratorRequest body = new JsonWebKeySetGeneratorRequest(); // JsonWebKeySetGeneratorRequest | 
try {
    JSONWebKeySet result = apiInstance.createJsonWebKeySet(set, body);
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling JsonWebKeyApi#createJsonWebKeySet");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **set** | **String**| The set |
 **body** | [**JsonWebKeySetGeneratorRequest**](JsonWebKeySetGeneratorRequest.md)|  | [optional]

### Return type

[**JSONWebKeySet**](JSONWebKeySet.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="deleteJsonWebKey"></a>
# **deleteJsonWebKey**
> deleteJsonWebKey(kid, set)

Delete a JSON Web Key

Use this endpoint to delete a single JSON Web Key.  A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.JsonWebKeyApi;


JsonWebKeyApi apiInstance = new JsonWebKeyApi();
String kid = "kid_example"; // String | The kid of the desired key
String set = "set_example"; // String | The set
try {
    apiInstance.deleteJsonWebKey(kid, set);
} catch (ApiException e) {
    System.err.println("Exception when calling JsonWebKeyApi#deleteJsonWebKey");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **kid** | **String**| The kid of the desired key |
 **set** | **String**| The set |

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="deleteJsonWebKeySet"></a>
# **deleteJsonWebKeySet**
> deleteJsonWebKeySet(set)

Delete a JSON Web Key Set

Use this endpoint to delete a complete JSON Web Key Set and all the keys in that set.  A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.JsonWebKeyApi;


JsonWebKeyApi apiInstance = new JsonWebKeyApi();
String set = "set_example"; // String | The set
try {
    apiInstance.deleteJsonWebKeySet(set);
} catch (ApiException e) {
    System.err.println("Exception when calling JsonWebKeyApi#deleteJsonWebKeySet");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **set** | **String**| The set |

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getJsonWebKey"></a>
# **getJsonWebKey**
> JSONWebKeySet getJsonWebKey(kid, set)

Retrieve a JSON Web Key

This endpoint can be used to retrieve JWKs stored in ORY Hydra.  A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.JsonWebKeyApi;


JsonWebKeyApi apiInstance = new JsonWebKeyApi();
String kid = "kid_example"; // String | The kid of the desired key
String set = "set_example"; // String | The set
try {
    JSONWebKeySet result = apiInstance.getJsonWebKey(kid, set);
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling JsonWebKeyApi#getJsonWebKey");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **kid** | **String**| The kid of the desired key |
 **set** | **String**| The set |

### Return type

[**JSONWebKeySet**](JSONWebKeySet.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getJsonWebKeySet"></a>
# **getJsonWebKeySet**
> JSONWebKeySet getJsonWebKeySet(set)

Retrieve a JSON Web Key Set

This endpoint can be used to retrieve JWK Sets stored in ORY Hydra.  A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.JsonWebKeyApi;


JsonWebKeyApi apiInstance = new JsonWebKeyApi();
String set = "set_example"; // String | The set
try {
    JSONWebKeySet result = apiInstance.getJsonWebKeySet(set);
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling JsonWebKeyApi#getJsonWebKeySet");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **set** | **String**| The set |

### Return type

[**JSONWebKeySet**](JSONWebKeySet.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="updateJsonWebKey"></a>
# **updateJsonWebKey**
> JSONWebKey updateJsonWebKey(kid, set, body)

Update a JSON Web Key

Use this method if you do not want to let Hydra generate the JWKs for you, but instead save your own.  A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.JsonWebKeyApi;


JsonWebKeyApi apiInstance = new JsonWebKeyApi();
String kid = "kid_example"; // String | The kid of the desired key
String set = "set_example"; // String | The set
JSONWebKey body = new JSONWebKey(); // JSONWebKey | 
try {
    JSONWebKey result = apiInstance.updateJsonWebKey(kid, set, body);
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling JsonWebKeyApi#updateJsonWebKey");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **kid** | **String**| The kid of the desired key |
 **set** | **String**| The set |
 **body** | [**JSONWebKey**](JSONWebKey.md)|  | [optional]

### Return type

[**JSONWebKey**](JSONWebKey.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="updateJsonWebKeySet"></a>
# **updateJsonWebKeySet**
> JSONWebKeySet updateJsonWebKeySet(set, body)

Update a JSON Web Key Set

Use this method if you do not want to let Hydra generate the JWKs for you, but instead save your own.  A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.

### Example
```java
// Import classes:
//import com.github.ory.hydra.ApiException;
//import com.github.ory.hydra.api.JsonWebKeyApi;


JsonWebKeyApi apiInstance = new JsonWebKeyApi();
String set = "set_example"; // String | The set
JSONWebKeySet body = new JSONWebKeySet(); // JSONWebKeySet | 
try {
    JSONWebKeySet result = apiInstance.updateJsonWebKeySet(set, body);
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling JsonWebKeyApi#updateJsonWebKeySet");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **set** | **String**| The set |
 **body** | [**JSONWebKeySet**](JSONWebKeySet.md)|  | [optional]

### Return type

[**JSONWebKeySet**](JSONWebKeySet.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

