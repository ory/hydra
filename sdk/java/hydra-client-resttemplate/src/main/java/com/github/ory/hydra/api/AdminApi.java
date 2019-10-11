package com.github.ory.hydra.api;

import com.github.ory.hydra.ApiClient;

import com.github.ory.hydra.model.AcceptConsentRequest;
import com.github.ory.hydra.model.AcceptLoginRequest;
import com.github.ory.hydra.model.CompletedRequest;
import com.github.ory.hydra.model.ConsentRequest;
import com.github.ory.hydra.model.FlushInactiveOAuth2TokensRequest;
import com.github.ory.hydra.model.GenericError;
import com.github.ory.hydra.model.HealthStatus;
import com.github.ory.hydra.model.JSONWebKey;
import com.github.ory.hydra.model.JSONWebKeySet;
import com.github.ory.hydra.model.JsonWebKeySetGeneratorRequest;
import com.github.ory.hydra.model.LoginRequest;
import com.github.ory.hydra.model.LogoutRequest;
import com.github.ory.hydra.model.OAuth2Client;
import com.github.ory.hydra.model.OAuth2TokenIntrospection;
import com.github.ory.hydra.model.PreviousConsentSession;
import com.github.ory.hydra.model.RejectRequest;
import com.github.ory.hydra.model.Version;

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
@Component("com.github.ory.hydra.api.AdminApi")
public class AdminApi {
    private ApiClient apiClient;

    public AdminApi() {
        this(new ApiClient());
    }

    @Autowired
    public AdminApi(ApiClient apiClient) {
        this.apiClient = apiClient;
    }

    public ApiClient getApiClient() {
        return apiClient;
    }

    public void setApiClient(ApiClient apiClient) {
        this.apiClient = apiClient;
    }

    /**
     * Accept an consent request
     * When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider to authenticate the subject and then tell ORY Hydra now about it. If the subject authenticated, he/she must now be asked if the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the subject&#39;s behalf.  The consent provider which handles this request and is a web app implemented and hosted by you. It shows a subject interface which asks the subject to grant or deny the client access to the requested scope (\&quot;Application my-dropbox-app wants write access to all your private files\&quot;).  The consent challenge is appended to the consent provider&#39;s URL to which the subject&#39;s user-agent (browser) is redirected to. The consent provider uses that challenge to fetch information on the OAuth2 request and then tells ORY Hydra if the subject accepted or rejected the request.  This endpoint tells ORY Hydra that the subject has authorized the OAuth 2.0 client to access resources on his/her behalf. The consent provider includes additional information, such as session data for access and ID tokens, and if the consent request should be used as basis for future requests.  The response contains a redirect URL which the consent provider should redirect the user-agent to.
     * <p><b>200</b> - completedRequest
     * <p><b>404</b> - genericError
     * <p><b>500</b> - genericError
     * @param consentChallenge The consentChallenge parameter
     * @param body The body parameter
     * @return CompletedRequest
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public CompletedRequest acceptConsentRequest(String consentChallenge, AcceptConsentRequest body) throws RestClientException {
        Object postBody = body;
        
        // verify the required parameter 'consentChallenge' is set
        if (consentChallenge == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'consentChallenge' when calling acceptConsentRequest");
        }
        
        String path = UriComponentsBuilder.fromPath("/oauth2/auth/requests/consent/accept").build().toUriString();
        
        final MultiValueMap<String, String> queryParams = new LinkedMultiValueMap<String, String>();
        final HttpHeaders headerParams = new HttpHeaders();
        final MultiValueMap<String, Object> formParams = new LinkedMultiValueMap<String, Object>();
        
        queryParams.putAll(apiClient.parameterToMultiValueMap(null, "consent_challenge", consentChallenge));

        final String[] accepts = { 
            "application/json"
        };
        final List<MediaType> accept = apiClient.selectHeaderAccept(accepts);
        final String[] contentTypes = { 
            "application/json"
        };
        final MediaType contentType = apiClient.selectHeaderContentType(contentTypes);

        String[] authNames = new String[] {  };

        ParameterizedTypeReference<CompletedRequest> returnType = new ParameterizedTypeReference<CompletedRequest>() {};
        return apiClient.invokeAPI(path, HttpMethod.PUT, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * Accept an login request
     * When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider (sometimes called \&quot;identity provider\&quot;) to authenticate the subject and then tell ORY Hydra now about it. The login provider is an web-app you write and host, and it must be able to authenticate (\&quot;show the subject a login screen\&quot;) a subject (in OAuth2 the proper name for subject is \&quot;resource owner\&quot;).  The authentication challenge is appended to the login provider URL to which the subject&#39;s user-agent (browser) is redirected to. The login provider uses that challenge to fetch information on the OAuth2 request and then accept or reject the requested authentication process.  This endpoint tells ORY Hydra that the subject has successfully authenticated and includes additional information such as the subject&#39;s ID and if ORY Hydra should remember the subject&#39;s subject agent for future authentication attempts by setting a cookie.  The response contains a redirect URL which the login provider should redirect the user-agent to.
     * <p><b>200</b> - completedRequest
     * <p><b>401</b> - genericError
     * <p><b>404</b> - genericError
     * <p><b>500</b> - genericError
     * @param loginChallenge The loginChallenge parameter
     * @param body The body parameter
     * @return CompletedRequest
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public CompletedRequest acceptLoginRequest(String loginChallenge, AcceptLoginRequest body) throws RestClientException {
        Object postBody = body;
        
        // verify the required parameter 'loginChallenge' is set
        if (loginChallenge == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'loginChallenge' when calling acceptLoginRequest");
        }
        
        String path = UriComponentsBuilder.fromPath("/oauth2/auth/requests/login/accept").build().toUriString();
        
        final MultiValueMap<String, String> queryParams = new LinkedMultiValueMap<String, String>();
        final HttpHeaders headerParams = new HttpHeaders();
        final MultiValueMap<String, Object> formParams = new LinkedMultiValueMap<String, Object>();
        
        queryParams.putAll(apiClient.parameterToMultiValueMap(null, "login_challenge", loginChallenge));

        final String[] accepts = { 
            "application/json"
        };
        final List<MediaType> accept = apiClient.selectHeaderAccept(accepts);
        final String[] contentTypes = { 
            "application/json"
        };
        final MediaType contentType = apiClient.selectHeaderContentType(contentTypes);

        String[] authNames = new String[] {  };

        ParameterizedTypeReference<CompletedRequest> returnType = new ParameterizedTypeReference<CompletedRequest>() {};
        return apiClient.invokeAPI(path, HttpMethod.PUT, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * Accept a logout request
     * When a user or an application requests ORY Hydra to log out a user, this endpoint is used to confirm that logout request. No body is required.  The response contains a redirect URL which the consent provider should redirect the user-agent to.
     * <p><b>200</b> - completedRequest
     * <p><b>404</b> - genericError
     * <p><b>500</b> - genericError
     * @param logoutChallenge The logoutChallenge parameter
     * @return CompletedRequest
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public CompletedRequest acceptLogoutRequest(String logoutChallenge) throws RestClientException {
        Object postBody = null;
        
        // verify the required parameter 'logoutChallenge' is set
        if (logoutChallenge == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'logoutChallenge' when calling acceptLogoutRequest");
        }
        
        String path = UriComponentsBuilder.fromPath("/oauth2/auth/requests/logout/accept").build().toUriString();
        
        final MultiValueMap<String, String> queryParams = new LinkedMultiValueMap<String, String>();
        final HttpHeaders headerParams = new HttpHeaders();
        final MultiValueMap<String, Object> formParams = new LinkedMultiValueMap<String, Object>();
        
        queryParams.putAll(apiClient.parameterToMultiValueMap(null, "logout_challenge", logoutChallenge));

        final String[] accepts = { 
            "application/json"
        };
        final List<MediaType> accept = apiClient.selectHeaderAccept(accepts);
        final String[] contentTypes = { 
            "application/json", "application/x-www-form-urlencoded"
        };
        final MediaType contentType = apiClient.selectHeaderContentType(contentTypes);

        String[] authNames = new String[] {  };

        ParameterizedTypeReference<CompletedRequest> returnType = new ParameterizedTypeReference<CompletedRequest>() {};
        return apiClient.invokeAPI(path, HttpMethod.PUT, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * Generate a new JSON Web Key
     * This endpoint is capable of generating JSON Web Key Sets for you. There a different strategies available, such as symmetric cryptographic keys (HS256, HS512) and asymetric cryptographic keys (RS256, ECDSA). If the specified JSON Web Key Set does not exist, it will be created.  A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.
     * <p><b>201</b> - JSONWebKeySet
     * <p><b>401</b> - genericError
     * <p><b>403</b> - genericError
     * <p><b>500</b> - genericError
     * @param set The set
     * @param body The body parameter
     * @return JSONWebKeySet
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public JSONWebKeySet createJsonWebKeySet(String set, JsonWebKeySetGeneratorRequest body) throws RestClientException {
        Object postBody = body;
        
        // verify the required parameter 'set' is set
        if (set == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'set' when calling createJsonWebKeySet");
        }
        
        // create path and map variables
        final Map<String, Object> uriVariables = new HashMap<String, Object>();
        uriVariables.put("set", set);
        String path = UriComponentsBuilder.fromPath("/keys/{set}").buildAndExpand(uriVariables).toUriString();
        
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
        return apiClient.invokeAPI(path, HttpMethod.POST, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * Create an OAuth 2.0 client
     * Create a new OAuth 2.0 client If you pass &#x60;client_secret&#x60; the secret will be used, otherwise a random secret will be generated. The secret will be returned in the response and you will not be able to retrieve it later on. Write the secret down and keep it somwhere safe.  OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.
     * <p><b>201</b> - oAuth2Client
     * <p><b>400</b> - genericError
     * <p><b>409</b> - genericError
     * <p><b>500</b> - genericError
     * @param body The body parameter
     * @return OAuth2Client
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public OAuth2Client createOAuth2Client(OAuth2Client body) throws RestClientException {
        Object postBody = body;
        
        // verify the required parameter 'body' is set
        if (body == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'body' when calling createOAuth2Client");
        }
        
        String path = UriComponentsBuilder.fromPath("/clients").build().toUriString();
        
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

        ParameterizedTypeReference<OAuth2Client> returnType = new ParameterizedTypeReference<OAuth2Client>() {};
        return apiClient.invokeAPI(path, HttpMethod.POST, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * Delete a JSON Web Key
     * Use this endpoint to delete a single JSON Web Key.  A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.
     * <p><b>204</b> - Empty responses are sent when, for example, resources are deleted. The HTTP status code for empty responses is typically 201.
     * <p><b>401</b> - genericError
     * <p><b>403</b> - genericError
     * <p><b>500</b> - genericError
     * @param kid The kid of the desired key
     * @param set The set
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public void deleteJsonWebKey(String kid, String set) throws RestClientException {
        Object postBody = null;
        
        // verify the required parameter 'kid' is set
        if (kid == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'kid' when calling deleteJsonWebKey");
        }
        
        // verify the required parameter 'set' is set
        if (set == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'set' when calling deleteJsonWebKey");
        }
        
        // create path and map variables
        final Map<String, Object> uriVariables = new HashMap<String, Object>();
        uriVariables.put("kid", kid);
        uriVariables.put("set", set);
        String path = UriComponentsBuilder.fromPath("/keys/{set}/{kid}").buildAndExpand(uriVariables).toUriString();
        
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

        ParameterizedTypeReference<Void> returnType = new ParameterizedTypeReference<Void>() {};
        apiClient.invokeAPI(path, HttpMethod.DELETE, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * Delete a JSON Web Key Set
     * Use this endpoint to delete a complete JSON Web Key Set and all the keys in that set.  A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.
     * <p><b>204</b> - Empty responses are sent when, for example, resources are deleted. The HTTP status code for empty responses is typically 201.
     * <p><b>401</b> - genericError
     * <p><b>403</b> - genericError
     * <p><b>500</b> - genericError
     * @param set The set
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public void deleteJsonWebKeySet(String set) throws RestClientException {
        Object postBody = null;
        
        // verify the required parameter 'set' is set
        if (set == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'set' when calling deleteJsonWebKeySet");
        }
        
        // create path and map variables
        final Map<String, Object> uriVariables = new HashMap<String, Object>();
        uriVariables.put("set", set);
        String path = UriComponentsBuilder.fromPath("/keys/{set}").buildAndExpand(uriVariables).toUriString();
        
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

        ParameterizedTypeReference<Void> returnType = new ParameterizedTypeReference<Void>() {};
        apiClient.invokeAPI(path, HttpMethod.DELETE, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * Deletes an OAuth 2.0 Client
     * Delete an existing OAuth 2.0 Client by its ID.  OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.
     * <p><b>204</b> - Empty responses are sent when, for example, resources are deleted. The HTTP status code for empty responses is typically 201.
     * <p><b>404</b> - genericError
     * <p><b>500</b> - genericError
     * @param id The id of the OAuth 2.0 Client.
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public void deleteOAuth2Client(String id) throws RestClientException {
        Object postBody = null;
        
        // verify the required parameter 'id' is set
        if (id == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'id' when calling deleteOAuth2Client");
        }
        
        // create path and map variables
        final Map<String, Object> uriVariables = new HashMap<String, Object>();
        uriVariables.put("id", id);
        String path = UriComponentsBuilder.fromPath("/clients/{id}").buildAndExpand(uriVariables).toUriString();
        
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

        ParameterizedTypeReference<Void> returnType = new ParameterizedTypeReference<Void>() {};
        apiClient.invokeAPI(path, HttpMethod.DELETE, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * Flush Expired OAuth2 Access Tokens
     * This endpoint flushes expired OAuth2 access tokens from the database. You can set a time after which no tokens will be not be touched, in case you want to keep recent tokens for auditing. Refresh tokens can not be flushed as they are deleted automatically when performing the refresh flow.
     * <p><b>204</b> - Empty responses are sent when, for example, resources are deleted. The HTTP status code for empty responses is typically 201.
     * <p><b>401</b> - genericError
     * <p><b>500</b> - genericError
     * @param body The body parameter
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public void flushInactiveOAuth2Tokens(FlushInactiveOAuth2TokensRequest body) throws RestClientException {
        Object postBody = body;
        
        String path = UriComponentsBuilder.fromPath("/oauth2/flush").build().toUriString();
        
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

        ParameterizedTypeReference<Void> returnType = new ParameterizedTypeReference<Void>() {};
        apiClient.invokeAPI(path, HttpMethod.POST, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * Get consent request information
     * When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider to authenticate the subject and then tell ORY Hydra now about it. If the subject authenticated, he/she must now be asked if the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the subject&#39;s behalf.  The consent provider which handles this request and is a web app implemented and hosted by you. It shows a subject interface which asks the subject to grant or deny the client access to the requested scope (\&quot;Application my-dropbox-app wants write access to all your private files\&quot;).  The consent challenge is appended to the consent provider&#39;s URL to which the subject&#39;s user-agent (browser) is redirected to. The consent provider uses that challenge to fetch information on the OAuth2 request and then tells ORY Hydra if the subject accepted or rejected the request.
     * <p><b>200</b> - consentRequest
     * <p><b>404</b> - genericError
     * <p><b>409</b> - genericError
     * <p><b>500</b> - genericError
     * @param consentChallenge The consentChallenge parameter
     * @return ConsentRequest
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public ConsentRequest getConsentRequest(String consentChallenge) throws RestClientException {
        Object postBody = null;
        
        // verify the required parameter 'consentChallenge' is set
        if (consentChallenge == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'consentChallenge' when calling getConsentRequest");
        }
        
        String path = UriComponentsBuilder.fromPath("/oauth2/auth/requests/consent").build().toUriString();
        
        final MultiValueMap<String, String> queryParams = new LinkedMultiValueMap<String, String>();
        final HttpHeaders headerParams = new HttpHeaders();
        final MultiValueMap<String, Object> formParams = new LinkedMultiValueMap<String, Object>();
        
        queryParams.putAll(apiClient.parameterToMultiValueMap(null, "consent_challenge", consentChallenge));

        final String[] accepts = { 
            "application/json"
        };
        final List<MediaType> accept = apiClient.selectHeaderAccept(accepts);
        final String[] contentTypes = { 
            "application/json"
        };
        final MediaType contentType = apiClient.selectHeaderContentType(contentTypes);

        String[] authNames = new String[] {  };

        ParameterizedTypeReference<ConsentRequest> returnType = new ParameterizedTypeReference<ConsentRequest>() {};
        return apiClient.invokeAPI(path, HttpMethod.GET, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * Fetch a JSON Web Key
     * This endpoint returns a singular JSON Web Key, identified by the set and the specific key ID (kid).
     * <p><b>200</b> - JSONWebKeySet
     * <p><b>404</b> - genericError
     * <p><b>500</b> - genericError
     * @param kid The kid of the desired key
     * @param set The set
     * @return JSONWebKeySet
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public JSONWebKeySet getJsonWebKey(String kid, String set) throws RestClientException {
        Object postBody = null;
        
        // verify the required parameter 'kid' is set
        if (kid == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'kid' when calling getJsonWebKey");
        }
        
        // verify the required parameter 'set' is set
        if (set == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'set' when calling getJsonWebKey");
        }
        
        // create path and map variables
        final Map<String, Object> uriVariables = new HashMap<String, Object>();
        uriVariables.put("kid", kid);
        uriVariables.put("set", set);
        String path = UriComponentsBuilder.fromPath("/keys/{set}/{kid}").buildAndExpand(uriVariables).toUriString();
        
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
    /**
     * Retrieve a JSON Web Key Set
     * This endpoint can be used to retrieve JWK Sets stored in ORY Hydra.  A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.
     * <p><b>200</b> - JSONWebKeySet
     * <p><b>401</b> - genericError
     * <p><b>403</b> - genericError
     * <p><b>500</b> - genericError
     * @param set The set
     * @return JSONWebKeySet
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public JSONWebKeySet getJsonWebKeySet(String set) throws RestClientException {
        Object postBody = null;
        
        // verify the required parameter 'set' is set
        if (set == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'set' when calling getJsonWebKeySet");
        }
        
        // create path and map variables
        final Map<String, Object> uriVariables = new HashMap<String, Object>();
        uriVariables.put("set", set);
        String path = UriComponentsBuilder.fromPath("/keys/{set}").buildAndExpand(uriVariables).toUriString();
        
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
    /**
     * Get an login request
     * When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider (sometimes called \&quot;identity provider\&quot;) to authenticate the subject and then tell ORY Hydra now about it. The login provider is an web-app you write and host, and it must be able to authenticate (\&quot;show the subject a login screen\&quot;) a subject (in OAuth2 the proper name for subject is \&quot;resource owner\&quot;).  The authentication challenge is appended to the login provider URL to which the subject&#39;s user-agent (browser) is redirected to. The login provider uses that challenge to fetch information on the OAuth2 request and then accept or reject the requested authentication process.
     * <p><b>200</b> - loginRequest
     * <p><b>400</b> - genericError
     * <p><b>404</b> - genericError
     * <p><b>409</b> - genericError
     * <p><b>500</b> - genericError
     * @param loginChallenge The loginChallenge parameter
     * @return LoginRequest
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public LoginRequest getLoginRequest(String loginChallenge) throws RestClientException {
        Object postBody = null;
        
        // verify the required parameter 'loginChallenge' is set
        if (loginChallenge == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'loginChallenge' when calling getLoginRequest");
        }
        
        String path = UriComponentsBuilder.fromPath("/oauth2/auth/requests/login").build().toUriString();
        
        final MultiValueMap<String, String> queryParams = new LinkedMultiValueMap<String, String>();
        final HttpHeaders headerParams = new HttpHeaders();
        final MultiValueMap<String, Object> formParams = new LinkedMultiValueMap<String, Object>();
        
        queryParams.putAll(apiClient.parameterToMultiValueMap(null, "login_challenge", loginChallenge));

        final String[] accepts = { 
            "application/json"
        };
        final List<MediaType> accept = apiClient.selectHeaderAccept(accepts);
        final String[] contentTypes = { 
            "application/json"
        };
        final MediaType contentType = apiClient.selectHeaderContentType(contentTypes);

        String[] authNames = new String[] {  };

        ParameterizedTypeReference<LoginRequest> returnType = new ParameterizedTypeReference<LoginRequest>() {};
        return apiClient.invokeAPI(path, HttpMethod.GET, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * Get a logout request
     * Use this endpoint to fetch a logout request.
     * <p><b>200</b> - logoutRequest
     * <p><b>404</b> - genericError
     * <p><b>500</b> - genericError
     * @param logoutChallenge The logoutChallenge parameter
     * @return LogoutRequest
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public LogoutRequest getLogoutRequest(String logoutChallenge) throws RestClientException {
        Object postBody = null;
        
        // verify the required parameter 'logoutChallenge' is set
        if (logoutChallenge == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'logoutChallenge' when calling getLogoutRequest");
        }
        
        String path = UriComponentsBuilder.fromPath("/oauth2/auth/requests/logout").build().toUriString();
        
        final MultiValueMap<String, String> queryParams = new LinkedMultiValueMap<String, String>();
        final HttpHeaders headerParams = new HttpHeaders();
        final MultiValueMap<String, Object> formParams = new LinkedMultiValueMap<String, Object>();
        
        queryParams.putAll(apiClient.parameterToMultiValueMap(null, "logout_challenge", logoutChallenge));

        final String[] accepts = { 
            "application/json"
        };
        final List<MediaType> accept = apiClient.selectHeaderAccept(accepts);
        final String[] contentTypes = { 
            "application/json", "application/x-www-form-urlencoded"
        };
        final MediaType contentType = apiClient.selectHeaderContentType(contentTypes);

        String[] authNames = new String[] {  };

        ParameterizedTypeReference<LogoutRequest> returnType = new ParameterizedTypeReference<LogoutRequest>() {};
        return apiClient.invokeAPI(path, HttpMethod.GET, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * Get an OAuth 2.0 Client.
     * Get an OAUth 2.0 client by its ID. This endpoint never returns passwords.  OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.
     * <p><b>200</b> - oAuth2Client
     * <p><b>404</b> - genericError
     * <p><b>500</b> - genericError
     * @param id The id of the OAuth 2.0 Client.
     * @return OAuth2Client
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public OAuth2Client getOAuth2Client(String id) throws RestClientException {
        Object postBody = null;
        
        // verify the required parameter 'id' is set
        if (id == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'id' when calling getOAuth2Client");
        }
        
        // create path and map variables
        final Map<String, Object> uriVariables = new HashMap<String, Object>();
        uriVariables.put("id", id);
        String path = UriComponentsBuilder.fromPath("/clients/{id}").buildAndExpand(uriVariables).toUriString();
        
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

        ParameterizedTypeReference<OAuth2Client> returnType = new ParameterizedTypeReference<OAuth2Client>() {};
        return apiClient.invokeAPI(path, HttpMethod.GET, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * Get service version
     * This endpoint returns the service version typically notated using semantic versioning.  If the service supports TLS Edge Termination, this endpoint does not require the &#x60;X-Forwarded-Proto&#x60; header to be set.  Be aware that if you are running multiple nodes of this service, the health status will never refer to the cluster state, only to a single instance.
     * <p><b>200</b> - version
     * @return Version
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public Version getVersion() throws RestClientException {
        Object postBody = null;
        
        String path = UriComponentsBuilder.fromPath("/version").build().toUriString();
        
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

        ParameterizedTypeReference<Version> returnType = new ParameterizedTypeReference<Version>() {};
        return apiClient.invokeAPI(path, HttpMethod.GET, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * Introspect OAuth2 tokens
     * The introspection endpoint allows to check if a token (both refresh and access) is active or not. An active token is neither expired nor revoked. If a token is active, additional information on the token will be included. You can set additional data for a token by setting &#x60;accessTokenExtra&#x60; during the consent flow.  For more information [read this blog post](https://www.oauth.com/oauth2-servers/token-introspection-endpoint/).
     * <p><b>200</b> - oAuth2TokenIntrospection
     * <p><b>401</b> - genericError
     * <p><b>500</b> - genericError
     * @param token The string value of the token. For access tokens, this is the \&quot;access_token\&quot; value returned from the token endpoint defined in OAuth 2.0. For refresh tokens, this is the \&quot;refresh_token\&quot; value returned.
     * @param scope An optional, space separated list of required scopes. If the access token was not granted one of the scopes, the result of active will be false.
     * @return OAuth2TokenIntrospection
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public OAuth2TokenIntrospection introspectOAuth2Token(String token, String scope) throws RestClientException {
        Object postBody = null;
        
        // verify the required parameter 'token' is set
        if (token == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'token' when calling introspectOAuth2Token");
        }
        
        String path = UriComponentsBuilder.fromPath("/oauth2/introspect").build().toUriString();
        
        final MultiValueMap<String, String> queryParams = new LinkedMultiValueMap<String, String>();
        final HttpHeaders headerParams = new HttpHeaders();
        final MultiValueMap<String, Object> formParams = new LinkedMultiValueMap<String, Object>();
        
        if (token != null)
            formParams.add("token", token);
        if (scope != null)
            formParams.add("scope", scope);

        final String[] accepts = { 
            "application/json"
        };
        final List<MediaType> accept = apiClient.selectHeaderAccept(accepts);
        final String[] contentTypes = { 
            "application/x-www-form-urlencoded"
        };
        final MediaType contentType = apiClient.selectHeaderContentType(contentTypes);

        String[] authNames = new String[] { "basic", "oauth2" };

        ParameterizedTypeReference<OAuth2TokenIntrospection> returnType = new ParameterizedTypeReference<OAuth2TokenIntrospection>() {};
        return apiClient.invokeAPI(path, HttpMethod.POST, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * Check alive status
     * This endpoint returns a 200 status code when the HTTP server is up running. This status does currently not include checks whether the database connection is working.  If the service supports TLS Edge Termination, this endpoint does not require the &#x60;X-Forwarded-Proto&#x60; header to be set.  Be aware that if you are running multiple nodes of this service, the health status will never refer to the cluster state, only to a single instance.
     * <p><b>200</b> - healthStatus
     * <p><b>500</b> - genericError
     * @return HealthStatus
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public HealthStatus isInstanceAlive() throws RestClientException {
        Object postBody = null;
        
        String path = UriComponentsBuilder.fromPath("/health/alive").build().toUriString();
        
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
     * List OAuth 2.0 Clients
     * This endpoint lists all clients in the database, and never returns client secrets.  OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components. The \&quot;Link\&quot; header is also included in successful responses, which contains one or more links for pagination, formatted like so: &#39;&lt;https://hydra-url/admin/clients?limit&#x3D;{limit}&amp;offset&#x3D;{offset}&gt;; rel&#x3D;\&quot;{page}\&quot;&#39;, where page is one of the following applicable pages: &#39;first&#39;, &#39;next&#39;, &#39;last&#39;, and &#39;previous&#39;. Multiple links can be included in this header, and will be separated by a comma.
     * <p><b>200</b> - A list of clients.
     * <p><b>500</b> - genericError
     * @param limit The maximum amount of policies returned.
     * @param offset The offset from where to start looking.
     * @return List&lt;OAuth2Client&gt;
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public List<OAuth2Client> listOAuth2Clients(Long limit, Long offset) throws RestClientException {
        Object postBody = null;
        
        String path = UriComponentsBuilder.fromPath("/clients").build().toUriString();
        
        final MultiValueMap<String, String> queryParams = new LinkedMultiValueMap<String, String>();
        final HttpHeaders headerParams = new HttpHeaders();
        final MultiValueMap<String, Object> formParams = new LinkedMultiValueMap<String, Object>();
        
        queryParams.putAll(apiClient.parameterToMultiValueMap(null, "limit", limit));
        queryParams.putAll(apiClient.parameterToMultiValueMap(null, "offset", offset));

        final String[] accepts = { 
            "application/json"
        };
        final List<MediaType> accept = apiClient.selectHeaderAccept(accepts);
        final String[] contentTypes = { 
            "application/json"
        };
        final MediaType contentType = apiClient.selectHeaderContentType(contentTypes);

        String[] authNames = new String[] {  };

        ParameterizedTypeReference<List<OAuth2Client>> returnType = new ParameterizedTypeReference<List<OAuth2Client>>() {};
        return apiClient.invokeAPI(path, HttpMethod.GET, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * Lists all consent sessions of a subject
     * This endpoint lists all subject&#39;s granted consent sessions, including client and granted scope. The \&quot;Link\&quot; header is also included in successful responses, which contains one or more links for pagination, formatted like so: &#39;&lt;https://hydra-url/admin/oauth2/auth/sessions/consent?subject&#x3D;{user}&amp;limit&#x3D;{limit}&amp;offset&#x3D;{offset}&gt;; rel&#x3D;\&quot;{page}\&quot;&#39;, where page is one of the following applicable pages: &#39;first&#39;, &#39;next&#39;, &#39;last&#39;, and &#39;previous&#39;. Multiple links can be included in this header, and will be separated by a comma.
     * <p><b>200</b> - A list of used consent requests.
     * <p><b>400</b> - genericError
     * <p><b>404</b> - genericError
     * <p><b>500</b> - genericError
     * @param subject The subject parameter
     * @return List&lt;PreviousConsentSession&gt;
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public List<PreviousConsentSession> listSubjectConsentSessions(String subject) throws RestClientException {
        Object postBody = null;
        
        // verify the required parameter 'subject' is set
        if (subject == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'subject' when calling listSubjectConsentSessions");
        }
        
        String path = UriComponentsBuilder.fromPath("/oauth2/auth/sessions/consent").build().toUriString();
        
        final MultiValueMap<String, String> queryParams = new LinkedMultiValueMap<String, String>();
        final HttpHeaders headerParams = new HttpHeaders();
        final MultiValueMap<String, Object> formParams = new LinkedMultiValueMap<String, Object>();
        
        queryParams.putAll(apiClient.parameterToMultiValueMap(null, "subject", subject));

        final String[] accepts = { 
            "application/json"
        };
        final List<MediaType> accept = apiClient.selectHeaderAccept(accepts);
        final String[] contentTypes = { 
            "application/json"
        };
        final MediaType contentType = apiClient.selectHeaderContentType(contentTypes);

        String[] authNames = new String[] {  };

        ParameterizedTypeReference<List<PreviousConsentSession>> returnType = new ParameterizedTypeReference<List<PreviousConsentSession>>() {};
        return apiClient.invokeAPI(path, HttpMethod.GET, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * Get snapshot metrics from the Hydra service. If you&#39;re using k8s, you can then add annotations to your deployment like so:
     * &#x60;&#x60;&#x60; metadata: annotations: prometheus.io/port: \&quot;4445\&quot; prometheus.io/path: \&quot;/metrics/prometheus\&quot; &#x60;&#x60;&#x60;
     * <p><b>200</b> - Empty responses are sent when, for example, resources are deleted. The HTTP status code for empty responses is typically 201.
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public void prometheus() throws RestClientException {
        Object postBody = null;
        
        String path = UriComponentsBuilder.fromPath("/metrics/prometheus").build().toUriString();
        
        final MultiValueMap<String, String> queryParams = new LinkedMultiValueMap<String, String>();
        final HttpHeaders headerParams = new HttpHeaders();
        final MultiValueMap<String, Object> formParams = new LinkedMultiValueMap<String, Object>();

        final String[] accepts = { 
            "plain/text"
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
     * Reject an consent request
     * When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider to authenticate the subject and then tell ORY Hydra now about it. If the subject authenticated, he/she must now be asked if the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the subject&#39;s behalf.  The consent provider which handles this request and is a web app implemented and hosted by you. It shows a subject interface which asks the subject to grant or deny the client access to the requested scope (\&quot;Application my-dropbox-app wants write access to all your private files\&quot;).  The consent challenge is appended to the consent provider&#39;s URL to which the subject&#39;s user-agent (browser) is redirected to. The consent provider uses that challenge to fetch information on the OAuth2 request and then tells ORY Hydra if the subject accepted or rejected the request.  This endpoint tells ORY Hydra that the subject has not authorized the OAuth 2.0 client to access resources on his/her behalf. The consent provider must include a reason why the consent was not granted.  The response contains a redirect URL which the consent provider should redirect the user-agent to.
     * <p><b>200</b> - completedRequest
     * <p><b>404</b> - genericError
     * <p><b>500</b> - genericError
     * @param consentChallenge The consentChallenge parameter
     * @param body The body parameter
     * @return CompletedRequest
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public CompletedRequest rejectConsentRequest(String consentChallenge, RejectRequest body) throws RestClientException {
        Object postBody = body;
        
        // verify the required parameter 'consentChallenge' is set
        if (consentChallenge == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'consentChallenge' when calling rejectConsentRequest");
        }
        
        String path = UriComponentsBuilder.fromPath("/oauth2/auth/requests/consent/reject").build().toUriString();
        
        final MultiValueMap<String, String> queryParams = new LinkedMultiValueMap<String, String>();
        final HttpHeaders headerParams = new HttpHeaders();
        final MultiValueMap<String, Object> formParams = new LinkedMultiValueMap<String, Object>();
        
        queryParams.putAll(apiClient.parameterToMultiValueMap(null, "consent_challenge", consentChallenge));

        final String[] accepts = { 
            "application/json"
        };
        final List<MediaType> accept = apiClient.selectHeaderAccept(accepts);
        final String[] contentTypes = { 
            "application/json"
        };
        final MediaType contentType = apiClient.selectHeaderContentType(contentTypes);

        String[] authNames = new String[] {  };

        ParameterizedTypeReference<CompletedRequest> returnType = new ParameterizedTypeReference<CompletedRequest>() {};
        return apiClient.invokeAPI(path, HttpMethod.PUT, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * Reject a login request
     * When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, ORY Hydra asks the login provider (sometimes called \&quot;identity provider\&quot;) to authenticate the subject and then tell ORY Hydra now about it. The login provider is an web-app you write and host, and it must be able to authenticate (\&quot;show the subject a login screen\&quot;) a subject (in OAuth2 the proper name for subject is \&quot;resource owner\&quot;).  The authentication challenge is appended to the login provider URL to which the subject&#39;s user-agent (browser) is redirected to. The login provider uses that challenge to fetch information on the OAuth2 request and then accept or reject the requested authentication process.  This endpoint tells ORY Hydra that the subject has not authenticated and includes a reason why the authentication was be denied.  The response contains a redirect URL which the login provider should redirect the user-agent to.
     * <p><b>200</b> - completedRequest
     * <p><b>401</b> - genericError
     * <p><b>404</b> - genericError
     * <p><b>500</b> - genericError
     * @param loginChallenge The loginChallenge parameter
     * @param body The body parameter
     * @return CompletedRequest
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public CompletedRequest rejectLoginRequest(String loginChallenge, RejectRequest body) throws RestClientException {
        Object postBody = body;
        
        // verify the required parameter 'loginChallenge' is set
        if (loginChallenge == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'loginChallenge' when calling rejectLoginRequest");
        }
        
        String path = UriComponentsBuilder.fromPath("/oauth2/auth/requests/login/reject").build().toUriString();
        
        final MultiValueMap<String, String> queryParams = new LinkedMultiValueMap<String, String>();
        final HttpHeaders headerParams = new HttpHeaders();
        final MultiValueMap<String, Object> formParams = new LinkedMultiValueMap<String, Object>();
        
        queryParams.putAll(apiClient.parameterToMultiValueMap(null, "login_challenge", loginChallenge));

        final String[] accepts = { 
            "application/json"
        };
        final List<MediaType> accept = apiClient.selectHeaderAccept(accepts);
        final String[] contentTypes = { 
            "application/json"
        };
        final MediaType contentType = apiClient.selectHeaderContentType(contentTypes);

        String[] authNames = new String[] {  };

        ParameterizedTypeReference<CompletedRequest> returnType = new ParameterizedTypeReference<CompletedRequest>() {};
        return apiClient.invokeAPI(path, HttpMethod.PUT, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * Reject a logout request
     * When a user or an application requests ORY Hydra to log out a user, this endpoint is used to deny that logout request. No body is required.  The response is empty as the logout provider has to chose what action to perform next.
     * <p><b>204</b> - Empty responses are sent when, for example, resources are deleted. The HTTP status code for empty responses is typically 201.
     * <p><b>404</b> - genericError
     * <p><b>500</b> - genericError
     * @param logoutChallenge The logoutChallenge parameter
     * @param body The body parameter
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public void rejectLogoutRequest(String logoutChallenge, RejectRequest body) throws RestClientException {
        Object postBody = body;
        
        // verify the required parameter 'logoutChallenge' is set
        if (logoutChallenge == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'logoutChallenge' when calling rejectLogoutRequest");
        }
        
        String path = UriComponentsBuilder.fromPath("/oauth2/auth/requests/logout/reject").build().toUriString();
        
        final MultiValueMap<String, String> queryParams = new LinkedMultiValueMap<String, String>();
        final HttpHeaders headerParams = new HttpHeaders();
        final MultiValueMap<String, Object> formParams = new LinkedMultiValueMap<String, Object>();
        
        queryParams.putAll(apiClient.parameterToMultiValueMap(null, "logout_challenge", logoutChallenge));

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
        apiClient.invokeAPI(path, HttpMethod.PUT, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * Invalidates all login sessions of a certain user Invalidates a subject&#39;s authentication session
     * This endpoint invalidates a subject&#39;s authentication session. After revoking the authentication session, the subject has to re-authenticate at ORY Hydra. This endpoint does not invalidate any tokens and does not work with OpenID Connect Front- or Back-channel logout.
     * <p><b>204</b> - Empty responses are sent when, for example, resources are deleted. The HTTP status code for empty responses is typically 201.
     * <p><b>400</b> - genericError
     * <p><b>404</b> - genericError
     * <p><b>500</b> - genericError
     * @param subject The subject parameter
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public void revokeAuthenticationSession(String subject) throws RestClientException {
        Object postBody = null;
        
        // verify the required parameter 'subject' is set
        if (subject == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'subject' when calling revokeAuthenticationSession");
        }
        
        String path = UriComponentsBuilder.fromPath("/oauth2/auth/sessions/login").build().toUriString();
        
        final MultiValueMap<String, String> queryParams = new LinkedMultiValueMap<String, String>();
        final HttpHeaders headerParams = new HttpHeaders();
        final MultiValueMap<String, Object> formParams = new LinkedMultiValueMap<String, Object>();
        
        queryParams.putAll(apiClient.parameterToMultiValueMap(null, "subject", subject));

        final String[] accepts = { 
            "application/json"
        };
        final List<MediaType> accept = apiClient.selectHeaderAccept(accepts);
        final String[] contentTypes = { 
            "application/json"
        };
        final MediaType contentType = apiClient.selectHeaderContentType(contentTypes);

        String[] authNames = new String[] {  };

        ParameterizedTypeReference<Void> returnType = new ParameterizedTypeReference<Void>() {};
        apiClient.invokeAPI(path, HttpMethod.DELETE, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * Revokes consent sessions of a subject for a specific OAuth 2.0 Client
     * This endpoint revokes a subject&#39;s granted consent sessions for a specific OAuth 2.0 Client and invalidates all associated OAuth 2.0 Access Tokens.
     * <p><b>204</b> - Empty responses are sent when, for example, resources are deleted. The HTTP status code for empty responses is typically 201.
     * <p><b>400</b> - genericError
     * <p><b>404</b> - genericError
     * <p><b>500</b> - genericError
     * @param subject The subject (Subject) who&#39;s consent sessions should be deleted.
     * @param client If set, deletes only those consent sessions by the Subject that have been granted to the specified OAuth 2.0 Client ID
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public void revokeConsentSessions(String subject, String client) throws RestClientException {
        Object postBody = null;
        
        // verify the required parameter 'subject' is set
        if (subject == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'subject' when calling revokeConsentSessions");
        }
        
        String path = UriComponentsBuilder.fromPath("/oauth2/auth/sessions/consent").build().toUriString();
        
        final MultiValueMap<String, String> queryParams = new LinkedMultiValueMap<String, String>();
        final HttpHeaders headerParams = new HttpHeaders();
        final MultiValueMap<String, Object> formParams = new LinkedMultiValueMap<String, Object>();
        
        queryParams.putAll(apiClient.parameterToMultiValueMap(null, "subject", subject));
        queryParams.putAll(apiClient.parameterToMultiValueMap(null, "client", client));

        final String[] accepts = { 
            "application/json"
        };
        final List<MediaType> accept = apiClient.selectHeaderAccept(accepts);
        final String[] contentTypes = { 
            "application/json"
        };
        final MediaType contentType = apiClient.selectHeaderContentType(contentTypes);

        String[] authNames = new String[] {  };

        ParameterizedTypeReference<Void> returnType = new ParameterizedTypeReference<Void>() {};
        apiClient.invokeAPI(path, HttpMethod.DELETE, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * Update a JSON Web Key
     * Use this method if you do not want to let Hydra generate the JWKs for you, but instead save your own.  A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.
     * <p><b>200</b> - JSONWebKey
     * <p><b>401</b> - genericError
     * <p><b>403</b> - genericError
     * <p><b>500</b> - genericError
     * @param kid The kid of the desired key
     * @param set The set
     * @param body The body parameter
     * @return JSONWebKey
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public JSONWebKey updateJsonWebKey(String kid, String set, JSONWebKey body) throws RestClientException {
        Object postBody = body;
        
        // verify the required parameter 'kid' is set
        if (kid == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'kid' when calling updateJsonWebKey");
        }
        
        // verify the required parameter 'set' is set
        if (set == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'set' when calling updateJsonWebKey");
        }
        
        // create path and map variables
        final Map<String, Object> uriVariables = new HashMap<String, Object>();
        uriVariables.put("kid", kid);
        uriVariables.put("set", set);
        String path = UriComponentsBuilder.fromPath("/keys/{set}/{kid}").buildAndExpand(uriVariables).toUriString();
        
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

        ParameterizedTypeReference<JSONWebKey> returnType = new ParameterizedTypeReference<JSONWebKey>() {};
        return apiClient.invokeAPI(path, HttpMethod.PUT, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * Update a JSON Web Key Set
     * Use this method if you do not want to let Hydra generate the JWKs for you, but instead save your own.  A JSON Web Key (JWK) is a JavaScript Object Notation (JSON) data structure that represents a cryptographic key. A JWK Set is a JSON data structure that represents a set of JWKs. A JSON Web Key is identified by its set and key id. ORY Hydra uses this functionality to store cryptographic keys used for TLS and JSON Web Tokens (such as OpenID Connect ID tokens), and allows storing user-defined keys as well.
     * <p><b>200</b> - JSONWebKeySet
     * <p><b>401</b> - genericError
     * <p><b>403</b> - genericError
     * <p><b>500</b> - genericError
     * @param set The set
     * @param body The body parameter
     * @return JSONWebKeySet
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public JSONWebKeySet updateJsonWebKeySet(String set, JSONWebKeySet body) throws RestClientException {
        Object postBody = body;
        
        // verify the required parameter 'set' is set
        if (set == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'set' when calling updateJsonWebKeySet");
        }
        
        // create path and map variables
        final Map<String, Object> uriVariables = new HashMap<String, Object>();
        uriVariables.put("set", set);
        String path = UriComponentsBuilder.fromPath("/keys/{set}").buildAndExpand(uriVariables).toUriString();
        
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
        return apiClient.invokeAPI(path, HttpMethod.PUT, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
    /**
     * Update an OAuth 2.0 Client
     * Update an existing OAuth 2.0 Client. If you pass &#x60;client_secret&#x60; the secret will be updated and returned via the API. This is the only time you will be able to retrieve the client secret, so write it down and keep it safe.  OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are generated for applications which want to consume your OAuth 2.0 or OpenID Connect capabilities. To manage ORY Hydra, you will need an OAuth 2.0 Client as well. Make sure that this endpoint is well protected and only callable by first-party components.
     * <p><b>200</b> - oAuth2Client
     * <p><b>500</b> - genericError
     * @param id The id parameter
     * @param body The body parameter
     * @return OAuth2Client
     * @throws RestClientException if an error occurs while attempting to invoke the API
     */
    public OAuth2Client updateOAuth2Client(String id, OAuth2Client body) throws RestClientException {
        Object postBody = body;
        
        // verify the required parameter 'id' is set
        if (id == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'id' when calling updateOAuth2Client");
        }
        
        // verify the required parameter 'body' is set
        if (body == null) {
            throw new HttpClientErrorException(HttpStatus.BAD_REQUEST, "Missing the required parameter 'body' when calling updateOAuth2Client");
        }
        
        // create path and map variables
        final Map<String, Object> uriVariables = new HashMap<String, Object>();
        uriVariables.put("id", id);
        String path = UriComponentsBuilder.fromPath("/clients/{id}").buildAndExpand(uriVariables).toUriString();
        
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

        ParameterizedTypeReference<OAuth2Client> returnType = new ParameterizedTypeReference<OAuth2Client>() {};
        return apiClient.invokeAPI(path, HttpMethod.PUT, queryParams, postBody, headerParams, formParams, accept, contentType, authNames, returnType);
    }
}
