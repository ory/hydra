# \HealthApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**IsInstanceAlive**](HealthApi.md#IsInstanceAlive) | **Get** /health/alive | Check the Alive Status
[**IsInstanceReady**](HealthApi.md#IsInstanceReady) | **Get** /health/ready | Check the Readiness Status


# **IsInstanceAlive**
> HealthStatus IsInstanceAlive()

Check the Alive Status

This endpoint returns a 200 status code when the HTTP server is up running. This status does currently not include checks whether the database connection is working. This endpoint does not require the `X-Forwarded-Proto` header when TLS termination is set.  Be aware that if you are running multiple nodes of ORY Hydra, the health status will never refer to the cluster state, only to a single instance.


### Parameters
This endpoint does not need any parameter.

### Return type

[**HealthStatus**](healthStatus.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json, application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **IsInstanceReady**
> HealthStatus IsInstanceReady()

Check the Readiness Status

This endpoint returns a 200 status code when the HTTP server is up running and the environment dependencies (e.g. the database) are responsive as well.  This status does currently not include checks whether the database connection is working. This endpoint does not require the `X-Forwarded-Proto` header when TLS termination is set.  Be aware that if you are running multiple nodes of ORY Hydra, the health status will never refer to the cluster state, only to a single instance.


### Parameters
This endpoint does not need any parameter.

### Return type

[**HealthStatus**](healthStatus.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json, application/x-www-form-urlencoded
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

