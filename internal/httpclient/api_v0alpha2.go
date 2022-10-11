/*
Ory Hydra API

Documentation for all of Ory Hydra's APIs.

API version:
Contact: hi@ory.sh
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// V0alpha2ApiService V0alpha2Api service
type V0alpha2ApiService service

type ApiAdminDeleteOAuth2TokenRequest struct {
	ctx        context.Context
	ApiService *V0alpha2ApiService
	clientId   *string
}

func (r ApiAdminDeleteOAuth2TokenRequest) ClientId(clientId string) ApiAdminDeleteOAuth2TokenRequest {
	r.clientId = &clientId
	return r
}

func (r ApiAdminDeleteOAuth2TokenRequest) Execute() (*http.Response, error) {
	return r.ApiService.AdminDeleteOAuth2TokenExecute(r)
}

/*
AdminDeleteOAuth2Token Delete OAuth2 Access Tokens from a Client

This endpoint deletes OAuth2 access tokens issued for a client from the database

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiAdminDeleteOAuth2TokenRequest
*/
func (a *V0alpha2ApiService) AdminDeleteOAuth2Token(ctx context.Context) ApiAdminDeleteOAuth2TokenRequest {
	return ApiAdminDeleteOAuth2TokenRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// Execute executes the request
func (a *V0alpha2ApiService) AdminDeleteOAuth2TokenExecute(r ApiAdminDeleteOAuth2TokenRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   interface{}
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "V0alpha2ApiService.AdminDeleteOAuth2Token")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/admin/oauth2/tokens"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.clientId == nil {
		return nil, reportError("clientId is required and must be specified")
	}

	localVarQueryParams.Add("client_id", parameterToString(*r.clientId, ""))
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		var v ErrorOAuth2
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarHTTPResponse, newErr
		}
		newErr.model = v
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}

type ApiAdminIntrospectOAuth2TokenRequest struct {
	ctx        context.Context
	ApiService *V0alpha2ApiService
	token      *string
	scope      *string
}

// The string value of the token. For access tokens, this is the \\\&quot;access_token\\\&quot; value returned from the token endpoint defined in OAuth 2.0. For refresh tokens, this is the \\\&quot;refresh_token\\\&quot; value returned.
func (r ApiAdminIntrospectOAuth2TokenRequest) Token(token string) ApiAdminIntrospectOAuth2TokenRequest {
	r.token = &token
	return r
}

// An optional, space separated list of required scopes. If the access token was not granted one of the scopes, the result of active will be false.
func (r ApiAdminIntrospectOAuth2TokenRequest) Scope(scope string) ApiAdminIntrospectOAuth2TokenRequest {
	r.scope = &scope
	return r
}

func (r ApiAdminIntrospectOAuth2TokenRequest) Execute() (*IntrospectedOAuth2Token, *http.Response, error) {
	return r.ApiService.AdminIntrospectOAuth2TokenExecute(r)
}

/*
AdminIntrospectOAuth2Token Introspect OAuth2 Access or Refresh Tokens

The introspection endpoint allows to check if a token (both refresh and access) is active or not. An active token
is neither expired nor revoked. If a token is active, additional information on the token will be included. You can
set additional data for a token by setting `accessTokenExtra` during the consent flow.

For more information [read this blog post](https://www.oauth.com/oauth2-servers/token-introspection-endpoint/).

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiAdminIntrospectOAuth2TokenRequest
*/
func (a *V0alpha2ApiService) AdminIntrospectOAuth2Token(ctx context.Context) ApiAdminIntrospectOAuth2TokenRequest {
	return ApiAdminIntrospectOAuth2TokenRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// Execute executes the request
//
//	@return IntrospectedOAuth2Token
func (a *V0alpha2ApiService) AdminIntrospectOAuth2TokenExecute(r ApiAdminIntrospectOAuth2TokenRequest) (*IntrospectedOAuth2Token, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *IntrospectedOAuth2Token
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "V0alpha2ApiService.AdminIntrospectOAuth2Token")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/admin/oauth2/introspect"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.token == nil {
		return localVarReturnValue, nil, reportError("token is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/x-www-form-urlencoded"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	if r.scope != nil {
		localVarFormParams.Add("scope", parameterToString(*r.scope, ""))
	}
	localVarFormParams.Add("token", parameterToString(*r.token, ""))
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		var v ErrorOAuth2
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiDeleteOidcDynamicClientRequest struct {
	ctx        context.Context
	ApiService *V0alpha2ApiService
	id         string
}

func (r ApiDeleteOidcDynamicClientRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteOidcDynamicClientExecute(r)
}

/*
DeleteOidcDynamicClient Delete OAuth 2.0 Client using the OpenID Dynamic Client Registration Management Protocol

This endpoint behaves like the administrative counterpart (`deleteOAuth2Client`) but is capable of facing the
public internet directly and can be used in self-service. It implements the OpenID Connect
Dynamic Client Registration Protocol. This feature needs to be enabled in the configuration. This endpoint
is disabled by default. It can be enabled by an administrator.

To use this endpoint, you will need to present the client's authentication credentials. If the OAuth2 Client
uses the Token Endpoint Authentication Method `client_secret_post`, you need to present the client secret in the URL query.
If it uses `client_secret_basic`, present the Client ID and the Client Secret in the Authorization header.

OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are
generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param id The id of the OAuth 2.0 Client.
	@return ApiDeleteOidcDynamicClientRequest
*/
func (a *V0alpha2ApiService) DeleteOidcDynamicClient(ctx context.Context, id string) ApiDeleteOidcDynamicClientRequest {
	return ApiDeleteOidcDynamicClientRequest{
		ApiService: a,
		ctx:        ctx,
		id:         id,
	}
}

// Execute executes the request
func (a *V0alpha2ApiService) DeleteOidcDynamicClientExecute(r ApiDeleteOidcDynamicClientRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   interface{}
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "V0alpha2ApiService.DeleteOidcDynamicClient")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauth2/register/{id}"
	localVarPath = strings.Replace(localVarPath, "{"+"id"+"}", url.PathEscape(parameterToString(r.id, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		var v GenericError
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarHTTPResponse, newErr
		}
		newErr.model = v
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}

type ApiDiscoverOidcConfigurationRequest struct {
	ctx        context.Context
	ApiService *V0alpha2ApiService
}

func (r ApiDiscoverOidcConfigurationRequest) Execute() (*OidcConfiguration, *http.Response, error) {
	return r.ApiService.DiscoverOidcConfigurationExecute(r)
}

/*
DiscoverOidcConfiguration OpenID Connect Discovery

The well known endpoint an be used to retrieve information for OpenID Connect clients. We encourage you to not roll
your own OpenID Connect client but to use an OpenID Connect client library instead. You can learn more on this
flow at https://openid.net/specs/openid-connect-discovery-1_0.html .

Popular libraries for OpenID Connect clients include oidc-client-js (JavaScript), go-oidc (Golang), and others.
For a full list of clients go here: https://openid.net/developers/certified/

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiDiscoverOidcConfigurationRequest
*/
func (a *V0alpha2ApiService) DiscoverOidcConfiguration(ctx context.Context) ApiDiscoverOidcConfigurationRequest {
	return ApiDiscoverOidcConfigurationRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// Execute executes the request
//
//	@return OidcConfiguration
func (a *V0alpha2ApiService) DiscoverOidcConfigurationExecute(r ApiDiscoverOidcConfigurationRequest) (*OidcConfiguration, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *OidcConfiguration
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "V0alpha2ApiService.DiscoverOidcConfiguration")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/.well-known/openid-configuration"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		var v ErrorOAuth2
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiGetOidcUserInfoRequest struct {
	ctx        context.Context
	ApiService *V0alpha2ApiService
}

func (r ApiGetOidcUserInfoRequest) Execute() (*OidcUserInfo, *http.Response, error) {
	return r.ApiService.GetOidcUserInfoExecute(r)
}

/*
GetOidcUserInfo OpenID Connect Userinfo

This endpoint returns the payload of the ID Token, including the idTokenExtra values, of
the provided OAuth 2.0 Access Token.

For more information please [refer to the spec](http://openid.net/specs/openid-connect-core-1_0.html#UserInfo).

In the case of authentication error, a WWW-Authenticate header might be set in the response
with more information about the error. See [the spec](https://datatracker.ietf.org/doc/html/rfc6750#section-3)
for more details about header format.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiGetOidcUserInfoRequest
*/
func (a *V0alpha2ApiService) GetOidcUserInfo(ctx context.Context) ApiGetOidcUserInfoRequest {
	return ApiGetOidcUserInfoRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// Execute executes the request
//
//	@return OidcUserInfo
func (a *V0alpha2ApiService) GetOidcUserInfoExecute(r ApiGetOidcUserInfoRequest) (*OidcUserInfo, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *OidcUserInfo
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "V0alpha2ApiService.GetOidcUserInfo")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/userinfo"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		var v ErrorOAuth2
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiPerformOAuth2AuthorizationFlowRequest struct {
	ctx        context.Context
	ApiService *V0alpha2ApiService
}

func (r ApiPerformOAuth2AuthorizationFlowRequest) Execute() (*ErrorOAuth2, *http.Response, error) {
	return r.ApiService.PerformOAuth2AuthorizationFlowExecute(r)
}

/*
PerformOAuth2AuthorizationFlow The OAuth 2.0 Authorize Endpoint

This endpoint is not documented here because you should never use your own implementation to perform OAuth2 flows.
OAuth2 is a very popular protocol and a library for your programming language will exists.

To learn more about this flow please refer to the specification: https://tools.ietf.org/html/rfc6749

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiPerformOAuth2AuthorizationFlowRequest
*/
func (a *V0alpha2ApiService) PerformOAuth2AuthorizationFlow(ctx context.Context) ApiPerformOAuth2AuthorizationFlowRequest {
	return ApiPerformOAuth2AuthorizationFlowRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// Execute executes the request
//
//	@return ErrorOAuth2
func (a *V0alpha2ApiService) PerformOAuth2AuthorizationFlowExecute(r ApiPerformOAuth2AuthorizationFlowRequest) (*ErrorOAuth2, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *ErrorOAuth2
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "V0alpha2ApiService.PerformOAuth2AuthorizationFlow")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauth2/auth"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		var v ErrorOAuth2
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiPerformOAuth2TokenFlowRequest struct {
	ctx          context.Context
	ApiService   *V0alpha2ApiService
	grantType    *string
	clientId     *string
	code         *string
	redirectUri  *string
	refreshToken *string
}

func (r ApiPerformOAuth2TokenFlowRequest) GrantType(grantType string) ApiPerformOAuth2TokenFlowRequest {
	r.grantType = &grantType
	return r
}

func (r ApiPerformOAuth2TokenFlowRequest) ClientId(clientId string) ApiPerformOAuth2TokenFlowRequest {
	r.clientId = &clientId
	return r
}

func (r ApiPerformOAuth2TokenFlowRequest) Code(code string) ApiPerformOAuth2TokenFlowRequest {
	r.code = &code
	return r
}

func (r ApiPerformOAuth2TokenFlowRequest) RedirectUri(redirectUri string) ApiPerformOAuth2TokenFlowRequest {
	r.redirectUri = &redirectUri
	return r
}

func (r ApiPerformOAuth2TokenFlowRequest) RefreshToken(refreshToken string) ApiPerformOAuth2TokenFlowRequest {
	r.refreshToken = &refreshToken
	return r
}

func (r ApiPerformOAuth2TokenFlowRequest) Execute() (*OAuth2TokenResponse, *http.Response, error) {
	return r.ApiService.PerformOAuth2TokenFlowExecute(r)
}

/*
PerformOAuth2TokenFlow The OAuth 2.0 Token Endpoint

The client makes a request to the token endpoint by sending the
following parameters using the "application/x-www-form-urlencoded" HTTP
request entity-body.

> Do not implement a client for this endpoint yourself. Use a library. There are many libraries
> available for any programming language. You can find a list of libraries here: https://oauth.net/code/
>
> Do note that Hydra SDK does not implement this endpoint properly. Use one of the libraries listed above

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiPerformOAuth2TokenFlowRequest
*/
func (a *V0alpha2ApiService) PerformOAuth2TokenFlow(ctx context.Context) ApiPerformOAuth2TokenFlowRequest {
	return ApiPerformOAuth2TokenFlowRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// Execute executes the request
//
//	@return OAuth2TokenResponse
func (a *V0alpha2ApiService) PerformOAuth2TokenFlowExecute(r ApiPerformOAuth2TokenFlowRequest) (*OAuth2TokenResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *OAuth2TokenResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "V0alpha2ApiService.PerformOAuth2TokenFlow")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauth2/token"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.grantType == nil {
		return localVarReturnValue, nil, reportError("grantType is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/x-www-form-urlencoded"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	if r.clientId != nil {
		localVarFormParams.Add("client_id", parameterToString(*r.clientId, ""))
	}
	if r.code != nil {
		localVarFormParams.Add("code", parameterToString(*r.code, ""))
	}
	localVarFormParams.Add("grant_type", parameterToString(*r.grantType, ""))
	if r.redirectUri != nil {
		localVarFormParams.Add("redirect_uri", parameterToString(*r.redirectUri, ""))
	}
	if r.refreshToken != nil {
		localVarFormParams.Add("refresh_token", parameterToString(*r.refreshToken, ""))
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		var v ErrorOAuth2
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiPerformOidcFrontOrBackChannelLogoutRequest struct {
	ctx        context.Context
	ApiService *V0alpha2ApiService
}

func (r ApiPerformOidcFrontOrBackChannelLogoutRequest) Execute() (*http.Response, error) {
	return r.ApiService.PerformOidcFrontOrBackChannelLogoutExecute(r)
}

/*
PerformOidcFrontOrBackChannelLogout OpenID Connect Front- or Back-channel Enabled Logout

This endpoint initiates and completes user logout at Ory Hydra and initiates OpenID Connect Front- / Back-channel logout:

https://openid.net/specs/openid-connect-frontchannel-1_0.html
https://openid.net/specs/openid-connect-backchannel-1_0.html

Back-channel logout is performed asynchronously and does not affect logout flow.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiPerformOidcFrontOrBackChannelLogoutRequest
*/
func (a *V0alpha2ApiService) PerformOidcFrontOrBackChannelLogout(ctx context.Context) ApiPerformOidcFrontOrBackChannelLogoutRequest {
	return ApiPerformOidcFrontOrBackChannelLogoutRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// Execute executes the request
func (a *V0alpha2ApiService) PerformOidcFrontOrBackChannelLogoutExecute(r ApiPerformOidcFrontOrBackChannelLogoutRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodGet
		localVarPostBody   interface{}
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "V0alpha2ApiService.PerformOidcFrontOrBackChannelLogout")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauth2/sessions/logout"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}

type ApiRevokeOAuth2TokenRequest struct {
	ctx        context.Context
	ApiService *V0alpha2ApiService
	token      *string
}

func (r ApiRevokeOAuth2TokenRequest) Token(token string) ApiRevokeOAuth2TokenRequest {
	r.token = &token
	return r
}

func (r ApiRevokeOAuth2TokenRequest) Execute() (*http.Response, error) {
	return r.ApiService.RevokeOAuth2TokenExecute(r)
}

/*
RevokeOAuth2Token Revoke an OAuth2 Access or Refresh Token

Revoking a token (both access and refresh) means that the tokens will be invalid. A revoked access token can no
longer be used to make access requests, and a revoked refresh token can no longer be used to refresh an access token.
Revoking a refresh token also invalidates the access token that was created with it. A token may only be revoked by
the client the token was generated for.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiRevokeOAuth2TokenRequest
*/
func (a *V0alpha2ApiService) RevokeOAuth2Token(ctx context.Context) ApiRevokeOAuth2TokenRequest {
	return ApiRevokeOAuth2TokenRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// Execute executes the request
func (a *V0alpha2ApiService) RevokeOAuth2TokenExecute(r ApiRevokeOAuth2TokenRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodPost
		localVarPostBody   interface{}
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "V0alpha2ApiService.RevokeOAuth2Token")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/oauth2/revoke"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.token == nil {
		return nil, reportError("token is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/x-www-form-urlencoded"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	localVarFormParams.Add("token", parameterToString(*r.token, ""))
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		var v ErrorOAuth2
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarHTTPResponse, newErr
		}
		newErr.model = v
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}
