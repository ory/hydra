# hydra-client-resttemplate

## Requirements

Building the API client library requires [Maven](https://maven.apache.org/) to be installed.

## Installation

To install the API client library to your local Maven repository, simply execute:

```shell
mvn install
```

To deploy it to a remote Maven repository instead, configure the settings of the repository and execute:

```shell
mvn deploy
```

Refer to the [official documentation](https://maven.apache.org/plugins/maven-deploy-plugin/usage.html) for more information.

### Maven users

Add this dependency to your project's POM:

```xml
<dependency>
    <groupId>com.github.ory</groupId>
    <artifactId>hydra-client-resttemplate</artifactId>
    <version>1.0.0</version>
    <scope>compile</scope>
</dependency>
```

### Gradle users

Add this dependency to your project's build file:

```groovy
compile "com.github.ory:hydra-client-resttemplate:1.0.0"
```

### Others

At first generate the JAR by executing:

    mvn package

Then manually install the following JARs:

* target/hydra-client-resttemplate-1.0.0.jar
* target/lib/*.jar

## Getting Started

Please follow the [installation](#installation) instruction and execute the following Java code:

```java

import com.github.ory.hydra.*;
import com.github.ory.hydra.auth.*;
import com.github.ory.hydra.model.*;
import com.github.ory.hydra.api.AdminApi;

import java.io.File;
import java.util.*;

public class AdminApiExample {

    public static void main(String[] args) {
        
        AdminApi apiInstance = new AdminApi();
        String challenge = "challenge_example"; // String | 
        AcceptConsentRequest body = new AcceptConsentRequest(); // AcceptConsentRequest | 
        try {
            CompletedRequest result = apiInstance.acceptConsentRequest(challenge, body);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling AdminApi#acceptConsentRequest");
            e.printStackTrace();
        }
    }
}

```

## Documentation for API Endpoints

All URIs are relative to *http://localhost*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*AdminApi* | [**acceptConsentRequest**](docs/AdminApi.md#acceptConsentRequest) | **PUT** /oauth2/auth/requests/consent/{challenge}/accept | Accept an consent request
*AdminApi* | [**acceptLoginRequest**](docs/AdminApi.md#acceptLoginRequest) | **PUT** /oauth2/auth/requests/login/{challenge}/accept | Accept an login request
*AdminApi* | [**createJsonWebKeySet**](docs/AdminApi.md#createJsonWebKeySet) | **POST** /keys/{set} | Generate a new JSON Web Key
*AdminApi* | [**createOAuth2Client**](docs/AdminApi.md#createOAuth2Client) | **POST** /clients | Create an OAuth 2.0 client
*AdminApi* | [**deleteJsonWebKey**](docs/AdminApi.md#deleteJsonWebKey) | **DELETE** /keys/{set}/{kid} | Delete a JSON Web Key
*AdminApi* | [**deleteJsonWebKeySet**](docs/AdminApi.md#deleteJsonWebKeySet) | **DELETE** /keys/{set} | Delete a JSON Web Key Set
*AdminApi* | [**deleteOAuth2Client**](docs/AdminApi.md#deleteOAuth2Client) | **DELETE** /clients/{id} | Deletes an OAuth 2.0 Client
*AdminApi* | [**flushInactiveOAuth2Tokens**](docs/AdminApi.md#flushInactiveOAuth2Tokens) | **POST** /oauth2/flush | Flush Expired OAuth2 Access Tokens
*AdminApi* | [**getConsentRequest**](docs/AdminApi.md#getConsentRequest) | **GET** /oauth2/auth/requests/consent/{challenge} | Get consent request information
*AdminApi* | [**getJsonWebKey**](docs/AdminApi.md#getJsonWebKey) | **GET** /keys/{set}/{kid} | Fetch a JSON Web Key
*AdminApi* | [**getJsonWebKeySet**](docs/AdminApi.md#getJsonWebKeySet) | **GET** /keys/{set} | Retrieve a JSON Web Key Set
*AdminApi* | [**getLoginRequest**](docs/AdminApi.md#getLoginRequest) | **GET** /oauth2/auth/requests/login/{challenge} | Get an login request
*AdminApi* | [**getOAuth2Client**](docs/AdminApi.md#getOAuth2Client) | **GET** /clients/{id} | Get an OAuth 2.0 Client.
*AdminApi* | [**introspectOAuth2Token**](docs/AdminApi.md#introspectOAuth2Token) | **POST** /oauth2/introspect | Introspect OAuth2 tokens
*AdminApi* | [**listOAuth2Clients**](docs/AdminApi.md#listOAuth2Clients) | **GET** /clients | List OAuth 2.0 Clients
*AdminApi* | [**listUserConsentSessions**](docs/AdminApi.md#listUserConsentSessions) | **GET** /oauth2/auth/sessions/consent/{user} | Lists all consent sessions of a user
*AdminApi* | [**rejectConsentRequest**](docs/AdminApi.md#rejectConsentRequest) | **PUT** /oauth2/auth/requests/consent/{challenge}/reject | Reject an consent request
*AdminApi* | [**rejectLoginRequest**](docs/AdminApi.md#rejectLoginRequest) | **PUT** /oauth2/auth/requests/login/{challenge}/reject | Reject a login request
*AdminApi* | [**revokeAllUserConsentSessions**](docs/AdminApi.md#revokeAllUserConsentSessions) | **DELETE** /oauth2/auth/sessions/consent/{user} | Revokes all previous consent sessions of a user
*AdminApi* | [**revokeAuthenticationSession**](docs/AdminApi.md#revokeAuthenticationSession) | **DELETE** /oauth2/auth/sessions/login/{user} | Invalidates a user&#39;s authentication session
*AdminApi* | [**revokeUserClientConsentSessions**](docs/AdminApi.md#revokeUserClientConsentSessions) | **DELETE** /oauth2/auth/sessions/consent/{user}/{client} | Revokes consent sessions of a user for a specific OAuth 2.0 Client
*AdminApi* | [**revokeUserLoginCookie**](docs/AdminApi.md#revokeUserLoginCookie) | **GET** /oauth2/auth/sessions/login/revoke | Logs user out by deleting the session cookie
*AdminApi* | [**updateJsonWebKey**](docs/AdminApi.md#updateJsonWebKey) | **PUT** /keys/{set}/{kid} | Update a JSON Web Key
*AdminApi* | [**updateJsonWebKeySet**](docs/AdminApi.md#updateJsonWebKeySet) | **PUT** /keys/{set} | Update a JSON Web Key Set
*AdminApi* | [**updateOAuth2Client**](docs/AdminApi.md#updateOAuth2Client) | **PUT** /clients/{id} | Update an OAuth 2.0 Client
*HealthApi* | [**isInstanceAlive**](docs/HealthApi.md#isInstanceAlive) | **GET** /health/alive | Check alive status
*HealthApi* | [**isInstanceReady**](docs/HealthApi.md#isInstanceReady) | **GET** /health/ready | Check readiness status
*PublicApi* | [**discoverOpenIDConfiguration**](docs/PublicApi.md#discoverOpenIDConfiguration) | **GET** /.well-known/openid-configuration | OpenID Connect Discovery
*PublicApi* | [**oauthAuth**](docs/PublicApi.md#oauthAuth) | **GET** /oauth2/auth | The OAuth 2.0 authorize endpoint
*PublicApi* | [**oauthToken**](docs/PublicApi.md#oauthToken) | **POST** /oauth2/token | The OAuth 2.0 token endpoint
*PublicApi* | [**revokeOAuth2Token**](docs/PublicApi.md#revokeOAuth2Token) | **POST** /oauth2/revoke | Revoke OAuth2 tokens
*PublicApi* | [**userinfo**](docs/PublicApi.md#userinfo) | **GET** /userinfo | OpenID Connect Userinfo
*PublicApi* | [**wellKnown**](docs/PublicApi.md#wellKnown) | **GET** /.well-known/jwks.json | JSON Web Keys Discovery
*VersionApi* | [**getVersion**](docs/VersionApi.md#getVersion) | **GET** /version | Get service version


## Documentation for Models

 - [AcceptConsentRequest](docs/AcceptConsentRequest.md)
 - [AcceptLoginRequest](docs/AcceptLoginRequest.md)
 - [AuthenticationSession](docs/AuthenticationSession.md)
 - [CompletedRequest](docs/CompletedRequest.md)
 - [ConsentRequest](docs/ConsentRequest.md)
 - [ConsentRequestSession](docs/ConsentRequestSession.md)
 - [EmptyResponse](docs/EmptyResponse.md)
 - [FlushInactiveOAuth2TokensRequest](docs/FlushInactiveOAuth2TokensRequest.md)
 - [GenericError](docs/GenericError.md)
 - [HealthNotReadyStatus](docs/HealthNotReadyStatus.md)
 - [HealthStatus](docs/HealthStatus.md)
 - [JSONWebKey](docs/JSONWebKey.md)
 - [JSONWebKeySet](docs/JSONWebKeySet.md)
 - [JsonWebKeySetGeneratorRequest](docs/JsonWebKeySetGeneratorRequest.md)
 - [LoginRequest](docs/LoginRequest.md)
 - [OAuth2Client](docs/OAuth2Client.md)
 - [OAuth2TokenIntrospection](docs/OAuth2TokenIntrospection.md)
 - [OauthTokenResponse](docs/OauthTokenResponse.md)
 - [OpenIDConnectContext](docs/OpenIDConnectContext.md)
 - [PreviousConsentSession](docs/PreviousConsentSession.md)
 - [RejectRequest](docs/RejectRequest.md)
 - [SwaggerFlushInactiveAccessTokens](docs/SwaggerFlushInactiveAccessTokens.md)
 - [SwaggerJsonWebKeyQuery](docs/SwaggerJsonWebKeyQuery.md)
 - [SwaggerJwkCreateSet](docs/SwaggerJwkCreateSet.md)
 - [SwaggerJwkSetQuery](docs/SwaggerJwkSetQuery.md)
 - [SwaggerJwkUpdateSet](docs/SwaggerJwkUpdateSet.md)
 - [SwaggerJwkUpdateSetKey](docs/SwaggerJwkUpdateSetKey.md)
 - [SwaggerOAuthIntrospectionRequest](docs/SwaggerOAuthIntrospectionRequest.md)
 - [SwaggerRevokeOAuth2TokenParameters](docs/SwaggerRevokeOAuth2TokenParameters.md)
 - [UserinfoResponse](docs/UserinfoResponse.md)
 - [Version](docs/Version.md)
 - [WellKnown](docs/WellKnown.md)


## Documentation for Authorization

Authentication schemes defined for the API:
### basic

- **Type**: HTTP basic authentication

### oauth2

- **Type**: OAuth
- **Flow**: accessCode
- **Authorization URL**: /oauth2/auth
- **Scopes**: 
  - offline: A scope required when requesting refresh tokens
  - openid: Request an OpenID Connect ID Token


## Recommendation

It's recommended to create an instance of `ApiClient` per thread in a multithreaded environment to avoid any potential issues.

## Author



