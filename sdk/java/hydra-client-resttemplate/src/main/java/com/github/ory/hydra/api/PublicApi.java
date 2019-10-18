package com.github.ory.hydra.api;

import com.github.ory.hydra.ApiClient;

import com.github.ory.hydra.model.GenericError;
import com.github.ory.hydra.model.HealthNotReadyStatus;
import com.github.ory.hydra.model.HealthStatus;
import com.github.ory.hydra.model.JSONWebKeySet;
import com.github.ory.hydra.model.Oauth2TokenResponse;
import com.github.ory.hydra.model.UserinfoResponse;
import com.github.ory.hydra.model.WellKnown;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import org.springframework.util.LinkedMultiValueMap;
import org.springframework.util.MultiValueMap;
import org.springframework.web.client.RestClientException;
import org.springframework.web.client.HttpClientErrorException;
import org.springframework.web.util.UriComponentsBuilder;
import org.springframework.core.ParameterizedTypeReference;
import org.springframework.core.io.FileSystemResource;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpMethod;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;

@javax.annotation.Generated(value = "io.swagger.codegen.languages.JavaClientCodegen", date = "2019-10-08T15:22:06.432-04:00")
@Component("com.github.ory.hydra.api.PublicApi")
public class PublicApi {
    private ApiClient apiClient;

    public PublicApi() {
        this(new ApiClient());
    }

    @Autowired
    public PublicApi(ApiClient apiClient) {
        this.apiClient = apiClient;
    }

    public ApiClient getApiClient() {
        return apiClient;
    }

    public void setApiClient(ApiClient apiClient) {
        this.apiClient = apiClient;
    }

    /**
     * OpenID Connect Front-Backchannel enabled Logout
     * This endpoint initiates and completes user logout at ORY Hydra and initiates OpenID Connect Front-/Back-channel logout:  https://openid.net/specs/openid-connect-frontchannel-1_0.html https://openid.net/specs/openid-connect-backchannel-1_0.html
     * <p><b>302</b> - Empty responses are sent when, for example, resources are deleted. The HTTP status code for empty responses is typically 201.
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public void disconnectUser() throws RestClientException {
        Object postBody = null;
        
        String path = UriComponentsBuilder.fromPath("/oauth2/sessions/logout").build().toUriString();
        
        final MultiValueMap<String, String> queryParams = new LinkedMultiValueMap<String, String>();
        final HttpHeaders headerParams = new HttpHeaders();
        final MultiValueMap<String, Object> formParams = new LinkedMultiValueMap<String, Object>();

        final String[] accepts = { 
            "application/json"
        };
        final List<MediaType> accept = apiClient.selectHeaderAccept(accepts);
        final String[] contentTypes = { 
            "application/json", "application/x-www-form-urlencoded"
        };
        final MediaType contentType = apiClient.selectHeaderContentType(contentTypes);

        String[] authNames = new String[] {  };

        ParameterizedTypeReference<Void> returnType = new ParameterizedTypeReference<Void>() {};
        apiClient.invokeAPI(path, HttpMethod.GET, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * OpenID Connect Discovery
     * The well known endpoint an be used to retrieve information for OpenID Connect clients. We encourage you to not roll your own OpenID Connect client but to use an OpenID Connect client library instead. You can learn more on this flow at https://openid.net/specs/openid-connect-discovery-1_0.html .  Popular libraries for OpenID Connect clients include oidc-client-js (JavaScript), go-oidc (Golang), and others. For a full list of clients go here: https://openid.net/developers/certified/
     * <p><b>200</b> - wellKnown
     * <p><b>401</b> - genericError
     * <p><b>500</b> - genericError
     * @return WellKnown
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public WellKnown discoverOpenIDConfiguration() throws RestClientException {
        Object postBody = null;
        
        String path = UriComponentsBuilder.fromPath("/.well-known/openid-configuration").build().toUriString();
        
        final MultiValueMap<String, String> queryParams = new LinkedMultiValueMap<String, String>();
        final HttpHeaders headerParams = new HttpHeaders();
        final MultiValueMap<String, Object> formParams = new LinkedMultiValueMap<String, Object>();

        final String[] accepts = { 
            "application/json"
        };
        final List<MediaType> accept = apiClient.selectHeaderAccept(accepts);
        final String[] contentTypes = { 
            "application/json", "application/x-www-form-urlencoded"
        };
        final MediaType contentType = apiClient.selectHeaderContentType(contentTypes);

        String[] authNames = new String[] {  };

        ParameterizedTypeReference<WellKnown> returnType = new ParameterizedTypeReference<WellKnown>() {};
        return apiClient.invokeAPI(path, HttpMethod.GET, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * Check readiness status
     * This endpoint returns a 200 status code when the HTTP server is up running and the environment dependencies (e.g. the database) are responsive as well.  If the service supports TLS Edge Termination, this endpoint does not require the &#x60;X-Forwarded-Proto&#x60; header to be set.  Be aware that if you are running multiple nodes of this service, the health status will never refer to the cluster state, only to a single instance.
     * <p><b>200</b> - healthStatus
     * <p><b>503</b> - healthNotReadyStatus
     * @return HealthStatus
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public HealthStatus isInstanceReady() throws RestClientException {
        Object postBody = null;
        
        String path = UriComponentsBuilder.fromPath("/health/ready").build().toUriString();
        
        final MultiValueMap<String, String> queryParams = new LinkedMultiValueMap<String, String>();
        final HttpHeaders headerParams = new HttpHeaders();
        final MultiValueMap<String, Object> formParams = new LinkedMultiValueMap<String, Object>();

        final String[] accepts = { 
            "application/json"
        };
        final List<MediaType> accept = apiClient.selectHeaderAccept(accepts);
        final String[] contentTypes = { 
            "application/json", "application/x-www-form-urlencoded"
        };
        final MediaType contentType = apiClient.selectHeaderContentType(contentTypes);

        String[] authNames = new String[] {  };

        ParameterizedTypeReference<HealthStatus> returnType = new ParameterizedTypeReference<HealthStatus>() {};
        return apiClient.invokeAPI(path, HttpMethod.GET, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * The OAuth 2.0 token endpoint
     * The client makes a request to the token endpoint by sending the following parameters using the \&quot;application/x-www-form-urlencoded\&quot; HTTP request entity-body.  &gt; Do not implement a client for this endpoint yourself. Use a library. There are many libraries &gt; available for any programming language. You can find a list of libraries here: https://oauth.net/code/ &gt; &gt; Do not the the Hydra SDK does not implement this endpoint properly. Use one of the libraries listed above!
     * <p><b>200</b> - oauth2TokenResponse
     * <p><b>401</b> - genericError
     * <p><b>500</b> - genericError
     * @param grantType The grantType parameter
     * @param code The code parameter
     * @param refreshToken The refreshToken parameter
     * @param redirectUri The redirectUri parameter
     * @param clientId The clientId parameter
     * @return Oauth2TokenResponse
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public Oauth2TokenResponse oauth2Token(String grantType, String code, String refreshToken, String redirectUri, String clientId) throws RestClientException {
        Object postBody = null;
        
        // verify the required parameter 'grantType' is set
        if (grantType == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'grantType' when calling oauth2Token");
        }
        
        String path = UriComponentsBuilder.fromPath("/oauth2/token").build().toUriString();
        
        final MultiValueMap<String, String> queryParams = new LinkedMultiValueMap<String, String>();
        final HttpHeaders headerParams = new HttpHeaders();
        final MultiValueMap<String, Object> formParams = new LinkedMultiValueMap<String, Object>();
        
        if (grantType != null)
            formParams.add("grant_type", grantType);
        if (code != null)
            formParams.add("code", code);
        if (refreshToken != null)
            formParams.add("refresh_token", refreshToken);
        if (redirectUri != null)
            formParams.add("redirect_uri", redirectUri);
        if (clientId != null)
            formParams.add("client_id", clientId);

        final String[] accepts = { 
            "application/json"
        };
        final List<MediaType> accept = apiClient.selectHeaderAccept(accepts);
        final String[] contentTypes = { 
            "application/x-www-form-urlencoded"
        };
        final MediaType contentType = apiClient.selectHeaderContentType(contentTypes);

        String[] authNames = new String[] { "basic", "oauth2" };

        ParameterizedTypeReference<Oauth2TokenResponse> returnType = new ParameterizedTypeReference<Oauth2TokenResponse>() {};
        return apiClient.invokeAPI(path, HttpMethod.POST, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * The OAuth 2.0 authorize endpoint
     * This endpoint is not documented here because you should never use your own implementation to perform OAuth2 flows. OAuth2 is a very popular protocol and a library for your programming language will exists.  To learn more about this flow please refer to the specification: https://tools.ietf.org/html/rfc6749
     * <p><b>302</b> - Empty responses are sent when, for example, resources are deleted. The HTTP status code for empty responses is typically 201.
     * <p><b>401</b> - genericError
     * <p><b>500</b> - genericError
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public void oauthAuth() throws RestClientException {
        Object postBody = null;
        
        String path = UriComponentsBuilder.fromPath("/oauth2/auth").build().toUriString();
        
        final MultiValueMap<String, String> queryParams = new LinkedMultiValueMap<String, String>();
        final HttpHeaders headerParams = new HttpHeaders();
        final MultiValueMap<String, Object> formParams = new LinkedMultiValueMap<String, Object>();

        final String[] accepts = { 
            "application/json"
        };
        final List<MediaType> accept = apiClient.selectHeaderAccept(accepts);
        final String[] contentTypes = { 
            "application/x-www-form-urlencoded"
        };
        final MediaType contentType = apiClient.selectHeaderContentType(contentTypes);

        String[] authNames = new String[] {  };

        ParameterizedTypeReference<Void> returnType = new ParameterizedTypeReference<Void>() {};
        apiClient.invokeAPI(path, HttpMethod.GET, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * Revoke OAuth2 tokens
     * Revoking a token (both access and refresh) means that the tokens will be invalid. A revoked access token can no longer be used to make access requests, and a revoked refresh token can no longer be used to refresh an access token. Revoking a refresh token also invalidates the access token that was created with it. A token may only be revoked by the client the token was generated for.
     * <p><b>200</b> - Empty responses are sent when, for example, resources are deleted. The HTTP status code for empty responses is typically 201.
     * <p><b>401</b> - genericError
     * <p><b>500</b> - genericError
     * @param token The token parameter
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public void revokeOAuth2Token(String token) throws RestClientException {
        Object postBody = null;
        
        // verify the required parameter 'token' is set
        if (token == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'token' when calling revokeOAuth2Token");
        }
        
        String path = UriComponentsBuilder.fromPath("/oauth2/revoke").build().toUriString();
        
        final MultiValueMap<String, String> queryParams = new LinkedMultiValueMap<String, String>();
        final HttpHeaders headerParams = new HttpHeaders();
        final MultiValueMap<String, Object> formParams = new LinkedMultiValueMap<String, Object>();
        
        if (token != null)
            formParams.add("token", token);

        final String[] accepts = { 
            "application/json"
        };
        final List<MediaType> accept = apiClient.selectHeaderAccept(accepts);
        final String[] contentTypes = { 
            "application/x-www-form-urlencoded"
        };
        final MediaType contentType = apiClient.selectHeaderContentType(contentTypes);

        String[] authNames = new String[] { "basic", "oauth2" };

        ParameterizedTypeReference<Void> returnType = new ParameterizedTypeReference<Void>() {};
        apiClient.invokeAPI(path, HttpMethod.POST, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * OpenID Connect Userinfo
     * This endpoint returns the payload of the ID Token, including the idTokenExtra values, of the provided OAuth 2.0 Access Token.  For more information please [refer to the spec](http://openid.net/specs/openid-connect-core-1_0.html#UserInfo).
     * <p><b>200</b> - userinfoResponse
     * <p><b>401</b> - genericError
     * <p><b>500</b> - genericError
     * @return UserinfoResponse
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public UserinfoResponse userinfo() throws RestClientException {
        Object postBody = null;
        
        String path = UriComponentsBuilder.fromPath("/userinfo").build().toUriString();
        
        final MultiValueMap<String, String> queryParams = new LinkedMultiValueMap<String, String>();
        final HttpHeaders headerParams = new HttpHeaders();
        final MultiValueMap<String, Object> formParams = new LinkedMultiValueMap<String, Object>();

        final String[] accepts = { 
            "application/json"
        };
        final List<MediaType> accept = apiClient.selectHeaderAccept(accepts);
        final String[] contentTypes = { 
            "application/json", "application/x-www-form-urlencoded"
        };
        final MediaType contentType = apiClient.selectHeaderContentType(contentTypes);

        String[] authNames = new String[] { "oauth2" };

        ParameterizedTypeReference<UserinfoResponse> returnType = new ParameterizedTypeReference<UserinfoResponse>() {};
        return apiClient.invokeAPI(path, HttpMethod.GET, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * JSON Web Keys Discovery
     * This endpoint returns JSON Web Keys to be used as public keys for verifying OpenID Connect ID Tokens and, if enabled, OAuth 2.0 JWT Access Tokens. This endpoint can be used with client libraries like [node-jwks-rsa](https://github.com/auth0/node-jwks-rsa) among others.
     * <p><b>200</b> - JSONWebKeySet
     * <p><b>500</b> - genericError
     * @return JSONWebKeySet
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public JSONWebKeySet wellKnown() throws RestClientException {
        Object postBody = null;
        
        String path = UriComponentsBuilder.fromPath("/.well-known/jwks.json").build().toUriString();
        
        final MultiValueMap<String, String> queryParams = new LinkedMultiValueMap<String, String>();
        final HttpHeaders headerParams = new HttpHeaders();
        final MultiValueMap<String, Object> formParams = new LinkedMultiValueMap<String, Object>();

        final String[] accepts = { 
            "application/json"
        };
        final List<MediaType> accept = apiClient.selectHeaderAccept(accepts);
        final String[] contentTypes = { 
            "application/json"
        };
        final MediaType contentType = apiClient.selectHeaderContentType(contentTypes);

        String[] authNames = new String[] {  };

        ParameterizedTypeReference<JSONWebKeySet> returnType = new ParameterizedTypeReference<JSONWebKeySet>() {};
        return apiClient.invokeAPI(path, HttpMethod.GET, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
}
