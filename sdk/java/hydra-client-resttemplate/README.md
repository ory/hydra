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
import com.github.ory.hydra.api.HealthApi;

import java.io.File;
import java.util.*;

public class HealthApiExample {

    public static void main(String[] args) {
        
        HealthApi apiInstance = new HealthApi();
        try {
            HealthStatus result = apiInstance.isInstanceAlive();
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling HealthApi#isInstanceAlive");
            e.printStackTrace();
        }
    }
}

```

## Documentation for API Endpoints

All URIs are relative to *http://localhost*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*HealthApi* | [**isInstanceAlive**](docs/HealthApi.md#isInstanceAlive) | **GET** /health/alive | Check alive status
*HealthApi* | [**isInstanceReady**](docs/HealthApi.md#isInstanceReady) | **GET** /health/ready | Check readiness status
*JsonWebKeyApi* | [**createJsonWebKeySet**](docs/JsonWebKeyApi.md#createJsonWebKeySet) | **POST** /keys/{set} | Generate a new JSON Web Key
*JsonWebKeyApi* | [**deleteJsonWebKey**](docs/JsonWebKeyApi.md#deleteJsonWebKey) | **DELETE** /keys/{set}/{kid} | Delete a JSON Web Key
*JsonWebKeyApi* | [**deleteJsonWebKeySet**](docs/JsonWebKeyApi.md#deleteJsonWebKeySet) | **DELETE** /keys/{set} | Delete a JSON Web Key Set
*JsonWebKeyApi* | [**getJsonWebKey**](docs/JsonWebKeyApi.md#getJsonWebKey) | **GET** /keys/{set}/{kid} | Retrieve a JSON Web Key
*JsonWebKeyApi* | [**getJsonWebKeySet**](docs/JsonWebKeyApi.md#getJsonWebKeySet) | **GET** /keys/{set} | Retrieve a JSON Web Key Set
*JsonWebKeyApi* | [**updateJsonWebKey**](docs/JsonWebKeyApi.md#updateJsonWebKey) | **PUT** /keys/{set}/{kid} | Update a JSON Web Key
*JsonWebKeyApi* | [**updateJsonWebKeySet**](docs/JsonWebKeyApi.md#updateJsonWebKeySet) | **PUT** /keys/{set} | Update a JSON Web Key Set
*OAuth2Api* | [**acceptConsentRequest**](docs/OAuth2Api.md#acceptConsentRequest) | **PUT** /oauth2/auth/requests/consent/{challenge}/accept | Accept an consent request
*OAuth2Api* | [**acceptLoginRequest**](docs/OAuth2Api.md#acceptLoginRequest) | **PUT** /oauth2/auth/requests/login/{challenge}/accept | Accept an login request
*OAuth2Api* | [**createOAuth2Client**](docs/OAuth2Api.md#createOAuth2Client) | **POST** /clients | Create an OAuth 2.0 client
*OAuth2Api* | [**deleteOAuth2Client**](docs/OAuth2Api.md#deleteOAuth2Client) | **DELETE** /clients/{id} | Deletes an OAuth 2.0 Client
*OAuth2Api* | [**flushInactiveOAuth2Tokens**](docs/OAuth2Api.md#flushInactiveOAuth2Tokens) | **POST** /oauth2/flush | Flush Expired OAuth2 Access Tokens
*OAuth2Api* | [**getConsentRequest**](docs/OAuth2Api.md#getConsentRequest) | **GET** /oauth2/auth/requests/consent/{challenge} | Get consent request information
*OAuth2Api* | [**getLoginRequest**](docs/OAuth2Api.md#getLoginRequest) | **GET** /oauth2/auth/requests/login/{challenge} | Get an login request
*OAuth2Api* | [**getOAuth2Client**](docs/OAuth2Api.md#getOAuth2Client) | **GET** /clients/{id} | Get an OAuth 2.0 Client.
*OAuth2Api* | [**getWellKnown**](docs/OAuth2Api.md#getWellKnown) | **GET** /.well-known/openid-configuration | Server well known configuration
*OAuth2Api* | [**introspectOAuth2Token**](docs/OAuth2Api.md#introspectOAuth2Token) | **POST** /oauth2/introspect | Introspect OAuth2 tokens
*OAuth2Api* | [**listOAuth2Clients**](docs/OAuth2Api.md#listOAuth2Clients) | **GET** /clients | List OAuth 2.0 Clients
*OAuth2Api* | [**listUserConsentSessions**](docs/OAuth2Api.md#listUserConsentSessions) | **GET** /oauth2/auth/sessions/consent/{user} | Lists all consent sessions of a user
*OAuth2Api* | [**oauthAuth**](docs/OAuth2Api.md#oauthAuth) | **GET** /oauth2/auth | The OAuth 2.0 authorize endpoint
*OAuth2Api* | [**oauthToken**](docs/OAuth2Api.md#oauthToken) | **POST** /oauth2/token | The OAuth 2.0 token endpoint
*OAuth2Api* | [**rejectConsentRequest**](docs/OAuth2Api.md#rejectConsentRequest) | **PUT** /oauth2/auth/requests/consent/{challenge}/reject | Reject an consent request
*OAuth2Api* | [**rejectLoginRequest**](docs/OAuth2Api.md#rejectLoginRequest) | **PUT** /oauth2/auth/requests/login/{challenge}/reject | Reject a login request
*OAuth2Api* | [**revokeAllUserConsentSessions**](docs/OAuth2Api.md#revokeAllUserConsentSessions) | **DELETE** /oauth2/auth/sessions/consent/{user} | Revokes all previous consent sessions of a user
*OAuth2Api* | [**revokeAuthenticationSession**](docs/OAuth2Api.md#revokeAuthenticationSession) | **DELETE** /oauth2/auth/sessions/login/{user} | Invalidates a user&#39;s authentication session
*OAuth2Api* | [**revokeOAuth2Token**](docs/OAuth2Api.md#revokeOAuth2Token) | **POST** /oauth2/revoke | Revoke OAuth2 tokens
*OAuth2Api* | [**revokeUserClientConsentSessions**](docs/OAuth2Api.md#revokeUserClientConsentSessions) | **DELETE** /oauth2/auth/sessions/consent/{user}/{client} | Revokes consent sessions of a user for a specific OAuth 2.0 Client
*OAuth2Api* | [**revokeUserLoginCookie**](docs/OAuth2Api.md#revokeUserLoginCookie) | **GET** /oauth2/auth/sessions/login/revoke | Logs user out by deleting the session cookie
*OAuth2Api* | [**updateOAuth2Client**](docs/OAuth2Api.md#updateOAuth2Client) | **PUT** /clients/{id} | Update an OAuth 2.0 Client
*OAuth2Api* | [**userinfo**](docs/OAuth2Api.md#userinfo) | **POST** /userinfo | OpenID Connect Userinfo
*OAuth2Api* | [**wellKnown**](docs/OAuth2Api.md#wellKnown) | **GET** /.well-known/jwks.json | Get Well-Known JSON Web Keys
*VersionApi* | [**getVersion**](docs/VersionApi.md#getVersion) | **GET** /version | Get service version


## Documentation for Models

 - [AcceptConsentRequest](docs/AcceptConsentRequest.md)
 - [AcceptLoginRequest](docs/AcceptLoginRequest.md)
 - [AuthenticationSession](docs/AuthenticationSession.md)
 - [CompletedRequest](docs/CompletedRequest.md)
 - [ConsentRequest](docs/ConsentRequest.md)
 - [ConsentRequestSession](docs/ConsentRequestSession.md)
 - [FlushInactiveOAuth2TokensRequest](docs/FlushInactiveOAuth2TokensRequest.md)
 - [HealthNotReadyStatus](docs/HealthNotReadyStatus.md)
 - [HealthStatus](docs/HealthStatus.md)
 - [InlineResponse401](docs/InlineResponse401.md)
 - [JSONWebKey](docs/JSONWebKey.md)
 - [JSONWebKeySet](docs/JSONWebKeySet.md)
 - [JoseWebKeySetRequest](docs/JoseWebKeySetRequest.md)
 - [JsonWebKeySetGeneratorRequest](docs/JsonWebKeySetGeneratorRequest.md)
 - [LoginRequest](docs/LoginRequest.md)
 - [OAuth2Client](docs/OAuth2Client.md)
 - [OAuth2TokenIntrospection](docs/OAuth2TokenIntrospection.md)
 - [OauthTokenResponse](docs/OauthTokenResponse.md)
 - [OpenIDConnectContext](docs/OpenIDConnectContext.md)
 - [PreviousConsentSession](docs/PreviousConsentSession.md)
 - [RawMessage](docs/RawMessage.md)
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
- **Authorization URL**: https://your-hydra-instance.com/oauth2/auth
- **Scopes**: 
  - offline: A scope required when requesting refresh tokens
  - openid: Request an OpenID Connect ID Token


## Recommendation

It's recommended to create an instance of `ApiClient` per thread in a multithreaded environment to avoid any potential issues.

## Author

hi@ory.am

