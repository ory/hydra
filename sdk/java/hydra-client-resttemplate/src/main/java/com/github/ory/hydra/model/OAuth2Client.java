/*
 * ORY Hydra
 * Welcome to the ORY Hydra HTTP API documentation. You will find documentation for all HTTP APIs here.
 *
 * OpenAPI spec version: latest
 * 
 *
 * NOTE: This class is auto generated by the swagger code generator program.
 * https://github.com/swagger-api/swagger-codegen.git
 * Do not edit the class manually.
 */


package com.github.ory.hydra.model;

import java.util.Objects;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonValue;
import com.github.ory.hydra.model.JSONWebKeySet;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.util.ArrayList;
import java.util.List;
import org.joda.time.DateTime;

/**
 * OAuth2Client
 */
@javax.annotation.Generated(value = "io.swagger.codegen.languages.JavaClientCodegen", date = "2019-08-19T20:15:39.753+02:00")
public class OAuth2Client {
  @JsonProperty("allowed_cors_origins")
  private List<String> allowedCorsOrigins = null;

  @JsonProperty("audience")
  private List<String> audience = null;

  @JsonProperty("backchannel_logout_session_required")
  private Boolean backchannelLogoutSessionRequired = null;

  @JsonProperty("backchannel_logout_uri")
  private String backchannelLogoutUri = null;

  @JsonProperty("client_id")
  private String clientId = null;

  @JsonProperty("client_name")
  private String clientName = null;

  @JsonProperty("client_secret")
  private String clientSecret = null;

  @JsonProperty("client_secret_expires_at")
  private Long clientSecretExpiresAt = null;

  @JsonProperty("client_uri")
  private String clientUri = null;

  @JsonProperty("contacts")
  private List<String> contacts = null;

  @JsonProperty("created_at")
  private DateTime createdAt = null;

  @JsonProperty("frontchannel_logout_session_required")
  private Boolean frontchannelLogoutSessionRequired = null;

  @JsonProperty("frontchannel_logout_uri")
  private String frontchannelLogoutUri = null;

  @JsonProperty("grant_types")
  private List<String> grantTypes = null;

  @JsonProperty("jwks")
  private JSONWebKeySet jwks = null;

  @JsonProperty("jwks_uri")
  private String jwksUri = null;

  @JsonProperty("logo_uri")
  private String logoUri = null;

  @JsonProperty("owner")
  private String owner = null;

  @JsonProperty("policy_uri")
  private String policyUri = null;

  @JsonProperty("post_logout_redirect_uris")
  private List<String> postLogoutRedirectUris = null;

  @JsonProperty("redirect_uris")
  private List<String> redirectUris = null;

  @JsonProperty("request_object_signing_alg")
  private String requestObjectSigningAlg = null;

  @JsonProperty("request_uris")
  private List<String> requestUris = null;

  @JsonProperty("response_types")
  private List<String> responseTypes = null;

  @JsonProperty("scope")
  private String scope = null;

  @JsonProperty("sector_identifier_uri")
  private String sectorIdentifierUri = null;

  @JsonProperty("subject_type")
  private String subjectType = null;

  @JsonProperty("token_endpoint_auth_method")
  private String tokenEndpointAuthMethod = null;

  @JsonProperty("tos_uri")
  private String tosUri = null;

  @JsonProperty("updated_at")
  private DateTime updatedAt = null;

  @JsonProperty("userinfo_signed_response_alg")
  private String userinfoSignedResponseAlg = null;

  public OAuth2Client allowedCorsOrigins(List<String> allowedCorsOrigins) {
    this.allowedCorsOrigins = allowedCorsOrigins;
    return this;
  }

  public OAuth2Client addAllowedCorsOriginsItem(String allowedCorsOriginsItem) {
    if (this.allowedCorsOrigins == null) {
      this.allowedCorsOrigins = new ArrayList<String>();
    }
    this.allowedCorsOrigins.add(allowedCorsOriginsItem);
    return this;
  }

   /**
   * AllowedCORSOrigins are one or more URLs (scheme://host[:port]) which are allowed to make CORS requests to the /oauth/token endpoint. If this array is empty, the sever&#39;s CORS origin configuration (&#x60;CORS_ALLOWED_ORIGINS&#x60;) will be used instead. If this array is set, the allowed origins are appended to the server&#39;s CORS origin configuration. Be aware that environment variable &#x60;CORS_ENABLED&#x60; MUST be set to &#x60;true&#x60; for this to work.
   * @return allowedCorsOrigins
  **/
  @ApiModelProperty(value = "AllowedCORSOrigins are one or more URLs (scheme://host[:port]) which are allowed to make CORS requests to the /oauth/token endpoint. If this array is empty, the sever's CORS origin configuration (`CORS_ALLOWED_ORIGINS`) will be used instead. If this array is set, the allowed origins are appended to the server's CORS origin configuration. Be aware that environment variable `CORS_ENABLED` MUST be set to `true` for this to work.")
  public List<String> getAllowedCorsOrigins() {
    return allowedCorsOrigins;
  }

  public void setAllowedCorsOrigins(List<String> allowedCorsOrigins) {
    this.allowedCorsOrigins = allowedCorsOrigins;
  }

  public OAuth2Client audience(List<String> audience) {
    this.audience = audience;
    return this;
  }

  public OAuth2Client addAudienceItem(String audienceItem) {
    if (this.audience == null) {
      this.audience = new ArrayList<String>();
    }
    this.audience.add(audienceItem);
    return this;
  }

   /**
   * Audience is a whitelist defining the audiences this client is allowed to request tokens for. An audience limits the applicability of an OAuth 2.0 Access Token to, for example, certain API endpoints. The value is a list of URLs. URLs MUST NOT contain whitespaces.
   * @return audience
  **/
  @ApiModelProperty(value = "Audience is a whitelist defining the audiences this client is allowed to request tokens for. An audience limits the applicability of an OAuth 2.0 Access Token to, for example, certain API endpoints. The value is a list of URLs. URLs MUST NOT contain whitespaces.")
  public List<String> getAudience() {
    return audience;
  }

  public void setAudience(List<String> audience) {
    this.audience = audience;
  }

  public OAuth2Client backchannelLogoutSessionRequired(Boolean backchannelLogoutSessionRequired) {
    this.backchannelLogoutSessionRequired = backchannelLogoutSessionRequired;
    return this;
  }

   /**
   * Boolean value specifying whether the RP requires that a sid (session ID) Claim be included in the Logout Token to identify the RP session with the OP when the backchannel_logout_uri is used. If omitted, the default value is false.
   * @return backchannelLogoutSessionRequired
  **/
  @ApiModelProperty(value = "Boolean value specifying whether the RP requires that a sid (session ID) Claim be included in the Logout Token to identify the RP session with the OP when the backchannel_logout_uri is used. If omitted, the default value is false.")
  public Boolean getBackchannelLogoutSessionRequired() {
    return backchannelLogoutSessionRequired;
  }

  public void setBackchannelLogoutSessionRequired(Boolean backchannelLogoutSessionRequired) {
    this.backchannelLogoutSessionRequired = backchannelLogoutSessionRequired;
  }

  public OAuth2Client backchannelLogoutUri(String backchannelLogoutUri) {
    this.backchannelLogoutUri = backchannelLogoutUri;
    return this;
  }

   /**
   * RP URL that will cause the RP to log itself out when sent a Logout Token by the OP.
   * @return backchannelLogoutUri
  **/
  @ApiModelProperty(value = "RP URL that will cause the RP to log itself out when sent a Logout Token by the OP.")
  public String getBackchannelLogoutUri() {
    return backchannelLogoutUri;
  }

  public void setBackchannelLogoutUri(String backchannelLogoutUri) {
    this.backchannelLogoutUri = backchannelLogoutUri;
  }

  public OAuth2Client clientId(String clientId) {
    this.clientId = clientId;
    return this;
  }

   /**
   * ClientID  is the id for this client.
   * @return clientId
  **/
  @ApiModelProperty(value = "ClientID  is the id for this client.")
  public String getClientId() {
    return clientId;
  }

  public void setClientId(String clientId) {
    this.clientId = clientId;
  }

  public OAuth2Client clientName(String clientName) {
    this.clientName = clientName;
    return this;
  }

   /**
   * Name is the human-readable string name of the client to be presented to the end-user during authorization.
   * @return clientName
  **/
  @ApiModelProperty(value = "Name is the human-readable string name of the client to be presented to the end-user during authorization.")
  public String getClientName() {
    return clientName;
  }

  public void setClientName(String clientName) {
    this.clientName = clientName;
  }

  public OAuth2Client clientSecret(String clientSecret) {
    this.clientSecret = clientSecret;
    return this;
  }

   /**
   * Secret is the client&#39;s secret. The secret will be included in the create request as cleartext, and then never again. The secret is stored using BCrypt so it is impossible to recover it. Tell your users that they need to write the secret down as it will not be made available again.
   * @return clientSecret
  **/
  @ApiModelProperty(value = "Secret is the client's secret. The secret will be included in the create request as cleartext, and then never again. The secret is stored using BCrypt so it is impossible to recover it. Tell your users that they need to write the secret down as it will not be made available again.")
  public String getClientSecret() {
    return clientSecret;
  }

  public void setClientSecret(String clientSecret) {
    this.clientSecret = clientSecret;
  }

  public OAuth2Client clientSecretExpiresAt(Long clientSecretExpiresAt) {
    this.clientSecretExpiresAt = clientSecretExpiresAt;
    return this;
  }

   /**
   * SecretExpiresAt is an integer holding the time at which the client secret will expire or 0 if it will not expire. The time is represented as the number of seconds from 1970-01-01T00:00:00Z as measured in UTC until the date/time of expiration.  This feature is currently not supported and it&#39;s value will always be set to 0.
   * @return clientSecretExpiresAt
  **/
  @ApiModelProperty(value = "SecretExpiresAt is an integer holding the time at which the client secret will expire or 0 if it will not expire. The time is represented as the number of seconds from 1970-01-01T00:00:00Z as measured in UTC until the date/time of expiration.  This feature is currently not supported and it's value will always be set to 0.")
  public Long getClientSecretExpiresAt() {
    return clientSecretExpiresAt;
  }

  public void setClientSecretExpiresAt(Long clientSecretExpiresAt) {
    this.clientSecretExpiresAt = clientSecretExpiresAt;
  }

  public OAuth2Client clientUri(String clientUri) {
    this.clientUri = clientUri;
    return this;
  }

   /**
   * ClientURI is an URL string of a web page providing information about the client. If present, the server SHOULD display this URL to the end-user in a clickable fashion.
   * @return clientUri
  **/
  @ApiModelProperty(value = "ClientURI is an URL string of a web page providing information about the client. If present, the server SHOULD display this URL to the end-user in a clickable fashion.")
  public String getClientUri() {
    return clientUri;
  }

  public void setClientUri(String clientUri) {
    this.clientUri = clientUri;
  }

  public OAuth2Client contacts(List<String> contacts) {
    this.contacts = contacts;
    return this;
  }

  public OAuth2Client addContactsItem(String contactsItem) {
    if (this.contacts == null) {
      this.contacts = new ArrayList<String>();
    }
    this.contacts.add(contactsItem);
    return this;
  }

   /**
   * Contacts is a array of strings representing ways to contact people responsible for this client, typically email addresses.
   * @return contacts
  **/
  @ApiModelProperty(value = "Contacts is a array of strings representing ways to contact people responsible for this client, typically email addresses.")
  public List<String> getContacts() {
    return contacts;
  }

  public void setContacts(List<String> contacts) {
    this.contacts = contacts;
  }

  public OAuth2Client createdAt(DateTime createdAt) {
    this.createdAt = createdAt;
    return this;
  }

   /**
   * CreatedAt returns the timestamp of the client&#39;s creation.
   * @return createdAt
  **/
  @ApiModelProperty(value = "CreatedAt returns the timestamp of the client's creation.")
  public DateTime getCreatedAt() {
    return createdAt;
  }

  public void setCreatedAt(DateTime createdAt) {
    this.createdAt = createdAt;
  }

  public OAuth2Client frontchannelLogoutSessionRequired(Boolean frontchannelLogoutSessionRequired) {
    this.frontchannelLogoutSessionRequired = frontchannelLogoutSessionRequired;
    return this;
  }

   /**
   * Boolean value specifying whether the RP requires that iss (issuer) and sid (session ID) query parameters be included to identify the RP session with the OP when the frontchannel_logout_uri is used. If omitted, the default value is false.
   * @return frontchannelLogoutSessionRequired
  **/
  @ApiModelProperty(value = "Boolean value specifying whether the RP requires that iss (issuer) and sid (session ID) query parameters be included to identify the RP session with the OP when the frontchannel_logout_uri is used. If omitted, the default value is false.")
  public Boolean getFrontchannelLogoutSessionRequired() {
    return frontchannelLogoutSessionRequired;
  }

  public void setFrontchannelLogoutSessionRequired(Boolean frontchannelLogoutSessionRequired) {
    this.frontchannelLogoutSessionRequired = frontchannelLogoutSessionRequired;
  }

  public OAuth2Client frontchannelLogoutUri(String frontchannelLogoutUri) {
    this.frontchannelLogoutUri = frontchannelLogoutUri;
    return this;
  }

   /**
   * RP URL that will cause the RP to log itself out when rendered in an iframe by the OP. An iss (issuer) query parameter and a sid (session ID) query parameter MAY be included by the OP to enable the RP to validate the request and to determine which of the potentially multiple sessions is to be logged out; if either is included, both MUST be.
   * @return frontchannelLogoutUri
  **/
  @ApiModelProperty(value = "RP URL that will cause the RP to log itself out when rendered in an iframe by the OP. An iss (issuer) query parameter and a sid (session ID) query parameter MAY be included by the OP to enable the RP to validate the request and to determine which of the potentially multiple sessions is to be logged out; if either is included, both MUST be.")
  public String getFrontchannelLogoutUri() {
    return frontchannelLogoutUri;
  }

  public void setFrontchannelLogoutUri(String frontchannelLogoutUri) {
    this.frontchannelLogoutUri = frontchannelLogoutUri;
  }

  public OAuth2Client grantTypes(List<String> grantTypes) {
    this.grantTypes = grantTypes;
    return this;
  }

  public OAuth2Client addGrantTypesItem(String grantTypesItem) {
    if (this.grantTypes == null) {
      this.grantTypes = new ArrayList<String>();
    }
    this.grantTypes.add(grantTypesItem);
    return this;
  }

   /**
   * GrantTypes is an array of grant types the client is allowed to use.
   * @return grantTypes
  **/
  @ApiModelProperty(value = "GrantTypes is an array of grant types the client is allowed to use.")
  public List<String> getGrantTypes() {
    return grantTypes;
  }

  public void setGrantTypes(List<String> grantTypes) {
    this.grantTypes = grantTypes;
  }

  public OAuth2Client jwks(JSONWebKeySet jwks) {
    this.jwks = jwks;
    return this;
  }

   /**
   * Get jwks
   * @return jwks
  **/
  @ApiModelProperty(value = "")
  public JSONWebKeySet getJwks() {
    return jwks;
  }

  public void setJwks(JSONWebKeySet jwks) {
    this.jwks = jwks;
  }

  public OAuth2Client jwksUri(String jwksUri) {
    this.jwksUri = jwksUri;
    return this;
  }

   /**
   * URL for the Client&#39;s JSON Web Key Set [JWK] document. If the Client signs requests to the Server, it contains the signing key(s) the Server uses to validate signatures from the Client. The JWK Set MAY also contain the Client&#39;s encryption keys(s), which are used by the Server to encrypt responses to the Client. When both signing and encryption keys are made available, a use (Key Use) parameter value is REQUIRED for all keys in the referenced JWK Set to indicate each key&#39;s intended usage. Although some algorithms allow the same key to be used for both signatures and encryption, doing so is NOT RECOMMENDED, as it is less secure. The JWK x5c parameter MAY be used to provide X.509 representations of keys provided. When used, the bare key values MUST still be present and MUST match those in the certificate.
   * @return jwksUri
  **/
  @ApiModelProperty(value = "URL for the Client's JSON Web Key Set [JWK] document. If the Client signs requests to the Server, it contains the signing key(s) the Server uses to validate signatures from the Client. The JWK Set MAY also contain the Client's encryption keys(s), which are used by the Server to encrypt responses to the Client. When both signing and encryption keys are made available, a use (Key Use) parameter value is REQUIRED for all keys in the referenced JWK Set to indicate each key's intended usage. Although some algorithms allow the same key to be used for both signatures and encryption, doing so is NOT RECOMMENDED, as it is less secure. The JWK x5c parameter MAY be used to provide X.509 representations of keys provided. When used, the bare key values MUST still be present and MUST match those in the certificate.")
  public String getJwksUri() {
    return jwksUri;
  }

  public void setJwksUri(String jwksUri) {
    this.jwksUri = jwksUri;
  }

  public OAuth2Client logoUri(String logoUri) {
    this.logoUri = logoUri;
    return this;
  }

   /**
   * LogoURI is an URL string that references a logo for the client.
   * @return logoUri
  **/
  @ApiModelProperty(value = "LogoURI is an URL string that references a logo for the client.")
  public String getLogoUri() {
    return logoUri;
  }

  public void setLogoUri(String logoUri) {
    this.logoUri = logoUri;
  }

  public OAuth2Client owner(String owner) {
    this.owner = owner;
    return this;
  }

   /**
   * Owner is a string identifying the owner of the OAuth 2.0 Client.
   * @return owner
  **/
  @ApiModelProperty(value = "Owner is a string identifying the owner of the OAuth 2.0 Client.")
  public String getOwner() {
    return owner;
  }

  public void setOwner(String owner) {
    this.owner = owner;
  }

  public OAuth2Client policyUri(String policyUri) {
    this.policyUri = policyUri;
    return this;
  }

   /**
   * PolicyURI is a URL string that points to a human-readable privacy policy document that describes how the deployment organization collects, uses, retains, and discloses personal data.
   * @return policyUri
  **/
  @ApiModelProperty(value = "PolicyURI is a URL string that points to a human-readable privacy policy document that describes how the deployment organization collects, uses, retains, and discloses personal data.")
  public String getPolicyUri() {
    return policyUri;
  }

  public void setPolicyUri(String policyUri) {
    this.policyUri = policyUri;
  }

  public OAuth2Client postLogoutRedirectUris(List<String> postLogoutRedirectUris) {
    this.postLogoutRedirectUris = postLogoutRedirectUris;
    return this;
  }

  public OAuth2Client addPostLogoutRedirectUrisItem(String postLogoutRedirectUrisItem) {
    if (this.postLogoutRedirectUris == null) {
      this.postLogoutRedirectUris = new ArrayList<String>();
    }
    this.postLogoutRedirectUris.add(postLogoutRedirectUrisItem);
    return this;
  }

   /**
   * Array of URLs supplied by the RP to which it MAY request that the End-User&#39;s User Agent be redirected using the post_logout_redirect_uri parameter after a logout has been performed.
   * @return postLogoutRedirectUris
  **/
  @ApiModelProperty(value = "Array of URLs supplied by the RP to which it MAY request that the End-User's User Agent be redirected using the post_logout_redirect_uri parameter after a logout has been performed.")
  public List<String> getPostLogoutRedirectUris() {
    return postLogoutRedirectUris;
  }

  public void setPostLogoutRedirectUris(List<String> postLogoutRedirectUris) {
    this.postLogoutRedirectUris = postLogoutRedirectUris;
  }

  public OAuth2Client redirectUris(List<String> redirectUris) {
    this.redirectUris = redirectUris;
    return this;
  }

  public OAuth2Client addRedirectUrisItem(String redirectUrisItem) {
    if (this.redirectUris == null) {
      this.redirectUris = new ArrayList<String>();
    }
    this.redirectUris.add(redirectUrisItem);
    return this;
  }

   /**
   * RedirectURIs is an array of allowed redirect urls for the client, for example http://mydomain/oauth/callback .
   * @return redirectUris
  **/
  @ApiModelProperty(value = "RedirectURIs is an array of allowed redirect urls for the client, for example http://mydomain/oauth/callback .")
  public List<String> getRedirectUris() {
    return redirectUris;
  }

  public void setRedirectUris(List<String> redirectUris) {
    this.redirectUris = redirectUris;
  }

  public OAuth2Client requestObjectSigningAlg(String requestObjectSigningAlg) {
    this.requestObjectSigningAlg = requestObjectSigningAlg;
    return this;
  }

   /**
   * JWS [JWS] alg algorithm [JWA] that MUST be used for signing Request Objects sent to the OP. All Request Objects from this Client MUST be rejected, if not signed with this algorithm.
   * @return requestObjectSigningAlg
  **/
  @ApiModelProperty(value = "JWS [JWS] alg algorithm [JWA] that MUST be used for signing Request Objects sent to the OP. All Request Objects from this Client MUST be rejected, if not signed with this algorithm.")
  public String getRequestObjectSigningAlg() {
    return requestObjectSigningAlg;
  }

  public void setRequestObjectSigningAlg(String requestObjectSigningAlg) {
    this.requestObjectSigningAlg = requestObjectSigningAlg;
  }

  public OAuth2Client requestUris(List<String> requestUris) {
    this.requestUris = requestUris;
    return this;
  }

  public OAuth2Client addRequestUrisItem(String requestUrisItem) {
    if (this.requestUris == null) {
      this.requestUris = new ArrayList<String>();
    }
    this.requestUris.add(requestUrisItem);
    return this;
  }

   /**
   * Array of request_uri values that are pre-registered by the RP for use at the OP. Servers MAY cache the contents of the files referenced by these URIs and not retrieve them at the time they are used in a request. OPs can require that request_uri values used be pre-registered with the require_request_uri_registration discovery parameter.
   * @return requestUris
  **/
  @ApiModelProperty(value = "Array of request_uri values that are pre-registered by the RP for use at the OP. Servers MAY cache the contents of the files referenced by these URIs and not retrieve them at the time they are used in a request. OPs can require that request_uri values used be pre-registered with the require_request_uri_registration discovery parameter.")
  public List<String> getRequestUris() {
    return requestUris;
  }

  public void setRequestUris(List<String> requestUris) {
    this.requestUris = requestUris;
  }

  public OAuth2Client responseTypes(List<String> responseTypes) {
    this.responseTypes = responseTypes;
    return this;
  }

  public OAuth2Client addResponseTypesItem(String responseTypesItem) {
    if (this.responseTypes == null) {
      this.responseTypes = new ArrayList<String>();
    }
    this.responseTypes.add(responseTypesItem);
    return this;
  }

   /**
   * ResponseTypes is an array of the OAuth 2.0 response type strings that the client can use at the authorization endpoint.
   * @return responseTypes
  **/
  @ApiModelProperty(value = "ResponseTypes is an array of the OAuth 2.0 response type strings that the client can use at the authorization endpoint.")
  public List<String> getResponseTypes() {
    return responseTypes;
  }

  public void setResponseTypes(List<String> responseTypes) {
    this.responseTypes = responseTypes;
  }

  public OAuth2Client scope(String scope) {
    this.scope = scope;
    return this;
  }

   /**
   * Scope is a string containing a space-separated list of scope values (as described in Section 3.3 of OAuth 2.0 [RFC6749]) that the client can use when requesting access tokens.
   * @return scope
  **/
  @ApiModelProperty(value = "Scope is a string containing a space-separated list of scope values (as described in Section 3.3 of OAuth 2.0 [RFC6749]) that the client can use when requesting access tokens.")
  public String getScope() {
    return scope;
  }

  public void setScope(String scope) {
    this.scope = scope;
  }

  public OAuth2Client sectorIdentifierUri(String sectorIdentifierUri) {
    this.sectorIdentifierUri = sectorIdentifierUri;
    return this;
  }

   /**
   * URL using the https scheme to be used in calculating Pseudonymous Identifiers by the OP. The URL references a file with a single JSON array of redirect_uri values.
   * @return sectorIdentifierUri
  **/
  @ApiModelProperty(value = "URL using the https scheme to be used in calculating Pseudonymous Identifiers by the OP. The URL references a file with a single JSON array of redirect_uri values.")
  public String getSectorIdentifierUri() {
    return sectorIdentifierUri;
  }

  public void setSectorIdentifierUri(String sectorIdentifierUri) {
    this.sectorIdentifierUri = sectorIdentifierUri;
  }

  public OAuth2Client subjectType(String subjectType) {
    this.subjectType = subjectType;
    return this;
  }

   /**
   * SubjectType requested for responses to this Client. The subject_types_supported Discovery parameter contains a list of the supported subject_type values for this server. Valid types include &#x60;pairwise&#x60; and &#x60;public&#x60;.
   * @return subjectType
  **/
  @ApiModelProperty(value = "SubjectType requested for responses to this Client. The subject_types_supported Discovery parameter contains a list of the supported subject_type values for this server. Valid types include `pairwise` and `public`.")
  public String getSubjectType() {
    return subjectType;
  }

  public void setSubjectType(String subjectType) {
    this.subjectType = subjectType;
  }

  public OAuth2Client tokenEndpointAuthMethod(String tokenEndpointAuthMethod) {
    this.tokenEndpointAuthMethod = tokenEndpointAuthMethod;
    return this;
  }

   /**
   * Requested Client Authentication method for the Token Endpoint. The options are client_secret_post, client_secret_basic, private_key_jwt, and none.
   * @return tokenEndpointAuthMethod
  **/
  @ApiModelProperty(value = "Requested Client Authentication method for the Token Endpoint. The options are client_secret_post, client_secret_basic, private_key_jwt, and none.")
  public String getTokenEndpointAuthMethod() {
    return tokenEndpointAuthMethod;
  }

  public void setTokenEndpointAuthMethod(String tokenEndpointAuthMethod) {
    this.tokenEndpointAuthMethod = tokenEndpointAuthMethod;
  }

  public OAuth2Client tosUri(String tosUri) {
    this.tosUri = tosUri;
    return this;
  }

   /**
   * TermsOfServiceURI is a URL string that points to a human-readable terms of service document for the client that describes a contractual relationship between the end-user and the client that the end-user accepts when authorizing the client.
   * @return tosUri
  **/
  @ApiModelProperty(value = "TermsOfServiceURI is a URL string that points to a human-readable terms of service document for the client that describes a contractual relationship between the end-user and the client that the end-user accepts when authorizing the client.")
  public String getTosUri() {
    return tosUri;
  }

  public void setTosUri(String tosUri) {
    this.tosUri = tosUri;
  }

  public OAuth2Client updatedAt(DateTime updatedAt) {
    this.updatedAt = updatedAt;
    return this;
  }

   /**
   * UpdatedAt returns the timestamp of the last update.
   * @return updatedAt
  **/
  @ApiModelProperty(value = "UpdatedAt returns the timestamp of the last update.")
  public DateTime getUpdatedAt() {
    return updatedAt;
  }

  public void setUpdatedAt(DateTime updatedAt) {
    this.updatedAt = updatedAt;
  }

  public OAuth2Client userinfoSignedResponseAlg(String userinfoSignedResponseAlg) {
    this.userinfoSignedResponseAlg = userinfoSignedResponseAlg;
    return this;
  }

   /**
   * JWS alg algorithm [JWA] REQUIRED for signing UserInfo Responses. If this is specified, the response will be JWT [JWT] serialized, and signed using JWS. The default, if omitted, is for the UserInfo Response to return the Claims as a UTF-8 encoded JSON object using the application/json content-type.
   * @return userinfoSignedResponseAlg
  **/
  @ApiModelProperty(value = "JWS alg algorithm [JWA] REQUIRED for signing UserInfo Responses. If this is specified, the response will be JWT [JWT] serialized, and signed using JWS. The default, if omitted, is for the UserInfo Response to return the Claims as a UTF-8 encoded JSON object using the application/json content-type.")
  public String getUserinfoSignedResponseAlg() {
    return userinfoSignedResponseAlg;
  }

  public void setUserinfoSignedResponseAlg(String userinfoSignedResponseAlg) {
    this.userinfoSignedResponseAlg = userinfoSignedResponseAlg;
  }


  @Override
  public boolean equals(java.lang.Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    OAuth2Client oAuth2Client = (OAuth2Client) o;
    return Objects.equals(this.allowedCorsOrigins, oAuth2Client.allowedCorsOrigins) &&
        Objects.equals(this.audience, oAuth2Client.audience) &&
        Objects.equals(this.backchannelLogoutSessionRequired, oAuth2Client.backchannelLogoutSessionRequired) &&
        Objects.equals(this.backchannelLogoutUri, oAuth2Client.backchannelLogoutUri) &&
        Objects.equals(this.clientId, oAuth2Client.clientId) &&
        Objects.equals(this.clientName, oAuth2Client.clientName) &&
        Objects.equals(this.clientSecret, oAuth2Client.clientSecret) &&
        Objects.equals(this.clientSecretExpiresAt, oAuth2Client.clientSecretExpiresAt) &&
        Objects.equals(this.clientUri, oAuth2Client.clientUri) &&
        Objects.equals(this.contacts, oAuth2Client.contacts) &&
        Objects.equals(this.createdAt, oAuth2Client.createdAt) &&
        Objects.equals(this.frontchannelLogoutSessionRequired, oAuth2Client.frontchannelLogoutSessionRequired) &&
        Objects.equals(this.frontchannelLogoutUri, oAuth2Client.frontchannelLogoutUri) &&
        Objects.equals(this.grantTypes, oAuth2Client.grantTypes) &&
        Objects.equals(this.jwks, oAuth2Client.jwks) &&
        Objects.equals(this.jwksUri, oAuth2Client.jwksUri) &&
        Objects.equals(this.logoUri, oAuth2Client.logoUri) &&
        Objects.equals(this.owner, oAuth2Client.owner) &&
        Objects.equals(this.policyUri, oAuth2Client.policyUri) &&
        Objects.equals(this.postLogoutRedirectUris, oAuth2Client.postLogoutRedirectUris) &&
        Objects.equals(this.redirectUris, oAuth2Client.redirectUris) &&
        Objects.equals(this.requestObjectSigningAlg, oAuth2Client.requestObjectSigningAlg) &&
        Objects.equals(this.requestUris, oAuth2Client.requestUris) &&
        Objects.equals(this.responseTypes, oAuth2Client.responseTypes) &&
        Objects.equals(this.scope, oAuth2Client.scope) &&
        Objects.equals(this.sectorIdentifierUri, oAuth2Client.sectorIdentifierUri) &&
        Objects.equals(this.subjectType, oAuth2Client.subjectType) &&
        Objects.equals(this.tokenEndpointAuthMethod, oAuth2Client.tokenEndpointAuthMethod) &&
        Objects.equals(this.tosUri, oAuth2Client.tosUri) &&
        Objects.equals(this.updatedAt, oAuth2Client.updatedAt) &&
        Objects.equals(this.userinfoSignedResponseAlg, oAuth2Client.userinfoSignedResponseAlg);
  }

  @Override
  public int hashCode() {
    return Objects.hash(allowedCorsOrigins, audience, backchannelLogoutSessionRequired, backchannelLogoutUri, clientId, clientName, clientSecret, clientSecretExpiresAt, clientUri, contacts, createdAt, frontchannelLogoutSessionRequired, frontchannelLogoutUri, grantTypes, jwks, jwksUri, logoUri, owner, policyUri, postLogoutRedirectUris, redirectUris, requestObjectSigningAlg, requestUris, responseTypes, scope, sectorIdentifierUri, subjectType, tokenEndpointAuthMethod, tosUri, updatedAt, userinfoSignedResponseAlg);
  }


  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class OAuth2Client {\n");
    
    sb.append("    allowedCorsOrigins: ").append(toIndentedString(allowedCorsOrigins)).append("\n");
    sb.append("    audience: ").append(toIndentedString(audience)).append("\n");
    sb.append("    backchannelLogoutSessionRequired: ").append(toIndentedString(backchannelLogoutSessionRequired)).append("\n");
    sb.append("    backchannelLogoutUri: ").append(toIndentedString(backchannelLogoutUri)).append("\n");
    sb.append("    clientId: ").append(toIndentedString(clientId)).append("\n");
    sb.append("    clientName: ").append(toIndentedString(clientName)).append("\n");
    sb.append("    clientSecret: ").append(toIndentedString(clientSecret)).append("\n");
    sb.append("    clientSecretExpiresAt: ").append(toIndentedString(clientSecretExpiresAt)).append("\n");
    sb.append("    clientUri: ").append(toIndentedString(clientUri)).append("\n");
    sb.append("    contacts: ").append(toIndentedString(contacts)).append("\n");
    sb.append("    createdAt: ").append(toIndentedString(createdAt)).append("\n");
    sb.append("    frontchannelLogoutSessionRequired: ").append(toIndentedString(frontchannelLogoutSessionRequired)).append("\n");
    sb.append("    frontchannelLogoutUri: ").append(toIndentedString(frontchannelLogoutUri)).append("\n");
    sb.append("    grantTypes: ").append(toIndentedString(grantTypes)).append("\n");
    sb.append("    jwks: ").append(toIndentedString(jwks)).append("\n");
    sb.append("    jwksUri: ").append(toIndentedString(jwksUri)).append("\n");
    sb.append("    logoUri: ").append(toIndentedString(logoUri)).append("\n");
    sb.append("    owner: ").append(toIndentedString(owner)).append("\n");
    sb.append("    policyUri: ").append(toIndentedString(policyUri)).append("\n");
    sb.append("    postLogoutRedirectUris: ").append(toIndentedString(postLogoutRedirectUris)).append("\n");
    sb.append("    redirectUris: ").append(toIndentedString(redirectUris)).append("\n");
    sb.append("    requestObjectSigningAlg: ").append(toIndentedString(requestObjectSigningAlg)).append("\n");
    sb.append("    requestUris: ").append(toIndentedString(requestUris)).append("\n");
    sb.append("    responseTypes: ").append(toIndentedString(responseTypes)).append("\n");
    sb.append("    scope: ").append(toIndentedString(scope)).append("\n");
    sb.append("    sectorIdentifierUri: ").append(toIndentedString(sectorIdentifierUri)).append("\n");
    sb.append("    subjectType: ").append(toIndentedString(subjectType)).append("\n");
    sb.append("    tokenEndpointAuthMethod: ").append(toIndentedString(tokenEndpointAuthMethod)).append("\n");
    sb.append("    tosUri: ").append(toIndentedString(tosUri)).append("\n");
    sb.append("    updatedAt: ").append(toIndentedString(updatedAt)).append("\n");
    sb.append("    userinfoSignedResponseAlg: ").append(toIndentedString(userinfoSignedResponseAlg)).append("\n");
    sb.append("}");
    return sb.toString();
  }

  /**
   * Convert the given object to string with each line indented by 4 spaces
   * (except the first line).
   */
  private String toIndentedString(java.lang.Object o) {
    if (o == null) {
      return "null";
    }
    return o.toString().replace("\n", "\n    ");
  }
  
}

