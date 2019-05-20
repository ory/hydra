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
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.util.ArrayList;
import java.util.List;

/**
 * It includes links to several endpoints (e.g. /oauth2/token) and exposes information on supported signature algorithms among others.
 */
@ApiModel(description = "It includes links to several endpoints (e.g. /oauth2/token) and exposes information on supported signature algorithms among others.")
@javax.annotation.Generated(value = "io.swagger.codegen.languages.JavaClientCodegen", date = "2019-05-17T10:32:12.084+02:00")
public class WellKnown {
  @JsonProperty("authorization_endpoint")
  private String authorizationEndpoint = null;

  @JsonProperty("backchannel_logout_session_supported")
  private Boolean backchannelLogoutSessionSupported = null;

  @JsonProperty("backchannel_logout_supported")
  private Boolean backchannelLogoutSupported = null;

  @JsonProperty("claims_parameter_supported")
  private Boolean claimsParameterSupported = null;

  @JsonProperty("claims_supported")
  private List<String> claimsSupported = null;

  @JsonProperty("end_session_endpoint")
  private String endSessionEndpoint = null;

  @JsonProperty("frontchannel_logout_session_supported")
  private Boolean frontchannelLogoutSessionSupported = null;

  @JsonProperty("frontchannel_logout_supported")
  private Boolean frontchannelLogoutSupported = null;

  @JsonProperty("grant_types_supported")
  private List<String> grantTypesSupported = null;

  @JsonProperty("id_token_signing_alg_values_supported")
  private List<String> idTokenSigningAlgValuesSupported = new ArrayList<String>();

  @JsonProperty("issuer")
  private String issuer = null;

  @JsonProperty("jwks_uri")
  private String jwksUri = null;

  @JsonProperty("registration_endpoint")
  private String registrationEndpoint = null;

  @JsonProperty("request_parameter_supported")
  private Boolean requestParameterSupported = null;

  @JsonProperty("request_uri_parameter_supported")
  private Boolean requestUriParameterSupported = null;

  @JsonProperty("require_request_uri_registration")
  private Boolean requireRequestUriRegistration = null;

  @JsonProperty("response_modes_supported")
  private List<String> responseModesSupported = null;

  @JsonProperty("response_types_supported")
  private List<String> responseTypesSupported = new ArrayList<String>();

  @JsonProperty("revocation_endpoint")
  private String revocationEndpoint = null;

  @JsonProperty("scopes_supported")
  private List<String> scopesSupported = null;

  @JsonProperty("subject_types_supported")
  private List<String> subjectTypesSupported = new ArrayList<String>();

  @JsonProperty("token_endpoint")
  private String tokenEndpoint = null;

  @JsonProperty("token_endpoint_auth_methods_supported")
  private List<String> tokenEndpointAuthMethodsSupported = null;

  @JsonProperty("userinfo_endpoint")
  private String userinfoEndpoint = null;

  @JsonProperty("userinfo_signing_alg_values_supported")
  private List<String> userinfoSigningAlgValuesSupported = null;

  public WellKnown authorizationEndpoint(String authorizationEndpoint) {
    this.authorizationEndpoint = authorizationEndpoint;
    return this;
  }

   /**
   * URL of the OP&#39;s OAuth 2.0 Authorization Endpoint.
   * @return authorizationEndpoint
  **/
  @ApiModelProperty(example = "https://playground.ory.sh/ory-hydra/public/oauth2/auth", required = true, value = "URL of the OP's OAuth 2.0 Authorization Endpoint.")
  public String getAuthorizationEndpoint() {
    return authorizationEndpoint;
  }

  public void setAuthorizationEndpoint(String authorizationEndpoint) {
    this.authorizationEndpoint = authorizationEndpoint;
  }

  public WellKnown backchannelLogoutSessionSupported(Boolean backchannelLogoutSessionSupported) {
    this.backchannelLogoutSessionSupported = backchannelLogoutSessionSupported;
    return this;
  }

   /**
   * Boolean value specifying whether the OP can pass a sid (session ID) Claim in the Logout Token to identify the RP session with the OP. If supported, the sid Claim is also included in ID Tokens issued by the OP
   * @return backchannelLogoutSessionSupported
  **/
  @ApiModelProperty(value = "Boolean value specifying whether the OP can pass a sid (session ID) Claim in the Logout Token to identify the RP session with the OP. If supported, the sid Claim is also included in ID Tokens issued by the OP")
  public Boolean getBackchannelLogoutSessionSupported() {
    return backchannelLogoutSessionSupported;
  }

  public void setBackchannelLogoutSessionSupported(Boolean backchannelLogoutSessionSupported) {
    this.backchannelLogoutSessionSupported = backchannelLogoutSessionSupported;
  }

  public WellKnown backchannelLogoutSupported(Boolean backchannelLogoutSupported) {
    this.backchannelLogoutSupported = backchannelLogoutSupported;
    return this;
  }

   /**
   * Boolean value specifying whether the OP supports back-channel logout, with true indicating support.
   * @return backchannelLogoutSupported
  **/
  @ApiModelProperty(value = "Boolean value specifying whether the OP supports back-channel logout, with true indicating support.")
  public Boolean getBackchannelLogoutSupported() {
    return backchannelLogoutSupported;
  }

  public void setBackchannelLogoutSupported(Boolean backchannelLogoutSupported) {
    this.backchannelLogoutSupported = backchannelLogoutSupported;
  }

  public WellKnown claimsParameterSupported(Boolean claimsParameterSupported) {
    this.claimsParameterSupported = claimsParameterSupported;
    return this;
  }

   /**
   * Boolean value specifying whether the OP supports use of the claims parameter, with true indicating support.
   * @return claimsParameterSupported
  **/
  @ApiModelProperty(value = "Boolean value specifying whether the OP supports use of the claims parameter, with true indicating support.")
  public Boolean getClaimsParameterSupported() {
    return claimsParameterSupported;
  }

  public void setClaimsParameterSupported(Boolean claimsParameterSupported) {
    this.claimsParameterSupported = claimsParameterSupported;
  }

  public WellKnown claimsSupported(List<String> claimsSupported) {
    this.claimsSupported = claimsSupported;
    return this;
  }

  public WellKnown addClaimsSupportedItem(String claimsSupportedItem) {
    if (this.claimsSupported == null) {
      this.claimsSupported = new ArrayList<String>();
    }
    this.claimsSupported.add(claimsSupportedItem);
    return this;
  }

   /**
   * JSON array containing a list of the Claim Names of the Claims that the OpenID Provider MAY be able to supply values for. Note that for privacy or other reasons, this might not be an exhaustive list.
   * @return claimsSupported
  **/
  @ApiModelProperty(value = "JSON array containing a list of the Claim Names of the Claims that the OpenID Provider MAY be able to supply values for. Note that for privacy or other reasons, this might not be an exhaustive list.")
  public List<String> getClaimsSupported() {
    return claimsSupported;
  }

  public void setClaimsSupported(List<String> claimsSupported) {
    this.claimsSupported = claimsSupported;
  }

  public WellKnown endSessionEndpoint(String endSessionEndpoint) {
    this.endSessionEndpoint = endSessionEndpoint;
    return this;
  }

   /**
   * URL at the OP to which an RP can perform a redirect to request that the End-User be logged out at the OP.
   * @return endSessionEndpoint
  **/
  @ApiModelProperty(value = "URL at the OP to which an RP can perform a redirect to request that the End-User be logged out at the OP.")
  public String getEndSessionEndpoint() {
    return endSessionEndpoint;
  }

  public void setEndSessionEndpoint(String endSessionEndpoint) {
    this.endSessionEndpoint = endSessionEndpoint;
  }

  public WellKnown frontchannelLogoutSessionSupported(Boolean frontchannelLogoutSessionSupported) {
    this.frontchannelLogoutSessionSupported = frontchannelLogoutSessionSupported;
    return this;
  }

   /**
   * Boolean value specifying whether the OP can pass iss (issuer) and sid (session ID) query parameters to identify the RP session with the OP when the frontchannel_logout_uri is used. If supported, the sid Claim is also included in ID Tokens issued by the OP.
   * @return frontchannelLogoutSessionSupported
  **/
  @ApiModelProperty(value = "Boolean value specifying whether the OP can pass iss (issuer) and sid (session ID) query parameters to identify the RP session with the OP when the frontchannel_logout_uri is used. If supported, the sid Claim is also included in ID Tokens issued by the OP.")
  public Boolean getFrontchannelLogoutSessionSupported() {
    return frontchannelLogoutSessionSupported;
  }

  public void setFrontchannelLogoutSessionSupported(Boolean frontchannelLogoutSessionSupported) {
    this.frontchannelLogoutSessionSupported = frontchannelLogoutSessionSupported;
  }

  public WellKnown frontchannelLogoutSupported(Boolean frontchannelLogoutSupported) {
    this.frontchannelLogoutSupported = frontchannelLogoutSupported;
    return this;
  }

   /**
   * Boolean value specifying whether the OP supports HTTP-based logout, with true indicating support.
   * @return frontchannelLogoutSupported
  **/
  @ApiModelProperty(value = "Boolean value specifying whether the OP supports HTTP-based logout, with true indicating support.")
  public Boolean getFrontchannelLogoutSupported() {
    return frontchannelLogoutSupported;
  }

  public void setFrontchannelLogoutSupported(Boolean frontchannelLogoutSupported) {
    this.frontchannelLogoutSupported = frontchannelLogoutSupported;
  }

  public WellKnown grantTypesSupported(List<String> grantTypesSupported) {
    this.grantTypesSupported = grantTypesSupported;
    return this;
  }

  public WellKnown addGrantTypesSupportedItem(String grantTypesSupportedItem) {
    if (this.grantTypesSupported == null) {
      this.grantTypesSupported = new ArrayList<String>();
    }
    this.grantTypesSupported.add(grantTypesSupportedItem);
    return this;
  }

   /**
   * JSON array containing a list of the OAuth 2.0 Grant Type values that this OP supports.
   * @return grantTypesSupported
  **/
  @ApiModelProperty(value = "JSON array containing a list of the OAuth 2.0 Grant Type values that this OP supports.")
  public List<String> getGrantTypesSupported() {
    return grantTypesSupported;
  }

  public void setGrantTypesSupported(List<String> grantTypesSupported) {
    this.grantTypesSupported = grantTypesSupported;
  }

  public WellKnown idTokenSigningAlgValuesSupported(List<String> idTokenSigningAlgValuesSupported) {
    this.idTokenSigningAlgValuesSupported = idTokenSigningAlgValuesSupported;
    return this;
  }

  public WellKnown addIdTokenSigningAlgValuesSupportedItem(String idTokenSigningAlgValuesSupportedItem) {
    this.idTokenSigningAlgValuesSupported.add(idTokenSigningAlgValuesSupportedItem);
    return this;
  }

   /**
   * JSON array containing a list of the JWS signing algorithms (alg values) supported by the OP for the ID Token to encode the Claims in a JWT.
   * @return idTokenSigningAlgValuesSupported
  **/
  @ApiModelProperty(required = true, value = "JSON array containing a list of the JWS signing algorithms (alg values) supported by the OP for the ID Token to encode the Claims in a JWT.")
  public List<String> getIdTokenSigningAlgValuesSupported() {
    return idTokenSigningAlgValuesSupported;
  }

  public void setIdTokenSigningAlgValuesSupported(List<String> idTokenSigningAlgValuesSupported) {
    this.idTokenSigningAlgValuesSupported = idTokenSigningAlgValuesSupported;
  }

  public WellKnown issuer(String issuer) {
    this.issuer = issuer;
    return this;
  }

   /**
   * URL using the https scheme with no query or fragment component that the OP asserts as its IssuerURL Identifier. If IssuerURL discovery is supported , this value MUST be identical to the issuer value returned by WebFinger. This also MUST be identical to the iss Claim value in ID Tokens issued from this IssuerURL.
   * @return issuer
  **/
  @ApiModelProperty(example = "https://playground.ory.sh/ory-hydra/public/", required = true, value = "URL using the https scheme with no query or fragment component that the OP asserts as its IssuerURL Identifier. If IssuerURL discovery is supported , this value MUST be identical to the issuer value returned by WebFinger. This also MUST be identical to the iss Claim value in ID Tokens issued from this IssuerURL.")
  public String getIssuer() {
    return issuer;
  }

  public void setIssuer(String issuer) {
    this.issuer = issuer;
  }

  public WellKnown jwksUri(String jwksUri) {
    this.jwksUri = jwksUri;
    return this;
  }

   /**
   * URL of the OP&#39;s JSON Web Key Set [JWK] document. This contains the signing key(s) the RP uses to validate signatures from the OP. The JWK Set MAY also contain the Server&#39;s encryption key(s), which are used by RPs to encrypt requests to the Server. When both signing and encryption keys are made available, a use (Key Use) parameter value is REQUIRED for all keys in the referenced JWK Set to indicate each key&#39;s intended usage. Although some algorithms allow the same key to be used for both signatures and encryption, doing so is NOT RECOMMENDED, as it is less secure. The JWK x5c parameter MAY be used to provide X.509 representations of keys provided. When used, the bare key values MUST still be present and MUST match those in the certificate.
   * @return jwksUri
  **/
  @ApiModelProperty(example = "https://playground.ory.sh/ory-hydra/public/.well-known/jwks.json", required = true, value = "URL of the OP's JSON Web Key Set [JWK] document. This contains the signing key(s) the RP uses to validate signatures from the OP. The JWK Set MAY also contain the Server's encryption key(s), which are used by RPs to encrypt requests to the Server. When both signing and encryption keys are made available, a use (Key Use) parameter value is REQUIRED for all keys in the referenced JWK Set to indicate each key's intended usage. Although some algorithms allow the same key to be used for both signatures and encryption, doing so is NOT RECOMMENDED, as it is less secure. The JWK x5c parameter MAY be used to provide X.509 representations of keys provided. When used, the bare key values MUST still be present and MUST match those in the certificate.")
  public String getJwksUri() {
    return jwksUri;
  }

  public void setJwksUri(String jwksUri) {
    this.jwksUri = jwksUri;
  }

  public WellKnown registrationEndpoint(String registrationEndpoint) {
    this.registrationEndpoint = registrationEndpoint;
    return this;
  }

   /**
   * URL of the OP&#39;s Dynamic Client Registration Endpoint.
   * @return registrationEndpoint
  **/
  @ApiModelProperty(example = "https://playground.ory.sh/ory-hydra/admin/client", value = "URL of the OP's Dynamic Client Registration Endpoint.")
  public String getRegistrationEndpoint() {
    return registrationEndpoint;
  }

  public void setRegistrationEndpoint(String registrationEndpoint) {
    this.registrationEndpoint = registrationEndpoint;
  }

  public WellKnown requestParameterSupported(Boolean requestParameterSupported) {
    this.requestParameterSupported = requestParameterSupported;
    return this;
  }

   /**
   * Boolean value specifying whether the OP supports use of the request parameter, with true indicating support.
   * @return requestParameterSupported
  **/
  @ApiModelProperty(value = "Boolean value specifying whether the OP supports use of the request parameter, with true indicating support.")
  public Boolean getRequestParameterSupported() {
    return requestParameterSupported;
  }

  public void setRequestParameterSupported(Boolean requestParameterSupported) {
    this.requestParameterSupported = requestParameterSupported;
  }

  public WellKnown requestUriParameterSupported(Boolean requestUriParameterSupported) {
    this.requestUriParameterSupported = requestUriParameterSupported;
    return this;
  }

   /**
   * Boolean value specifying whether the OP supports use of the request_uri parameter, with true indicating support.
   * @return requestUriParameterSupported
  **/
  @ApiModelProperty(value = "Boolean value specifying whether the OP supports use of the request_uri parameter, with true indicating support.")
  public Boolean getRequestUriParameterSupported() {
    return requestUriParameterSupported;
  }

  public void setRequestUriParameterSupported(Boolean requestUriParameterSupported) {
    this.requestUriParameterSupported = requestUriParameterSupported;
  }

  public WellKnown requireRequestUriRegistration(Boolean requireRequestUriRegistration) {
    this.requireRequestUriRegistration = requireRequestUriRegistration;
    return this;
  }

   /**
   * Boolean value specifying whether the OP requires any request_uri values used to be pre-registered using the request_uris registration parameter.
   * @return requireRequestUriRegistration
  **/
  @ApiModelProperty(value = "Boolean value specifying whether the OP requires any request_uri values used to be pre-registered using the request_uris registration parameter.")
  public Boolean getRequireRequestUriRegistration() {
    return requireRequestUriRegistration;
  }

  public void setRequireRequestUriRegistration(Boolean requireRequestUriRegistration) {
    this.requireRequestUriRegistration = requireRequestUriRegistration;
  }

  public WellKnown responseModesSupported(List<String> responseModesSupported) {
    this.responseModesSupported = responseModesSupported;
    return this;
  }

  public WellKnown addResponseModesSupportedItem(String responseModesSupportedItem) {
    if (this.responseModesSupported == null) {
      this.responseModesSupported = new ArrayList<String>();
    }
    this.responseModesSupported.add(responseModesSupportedItem);
    return this;
  }

   /**
   * JSON array containing a list of the OAuth 2.0 response_mode values that this OP supports.
   * @return responseModesSupported
  **/
  @ApiModelProperty(value = "JSON array containing a list of the OAuth 2.0 response_mode values that this OP supports.")
  public List<String> getResponseModesSupported() {
    return responseModesSupported;
  }

  public void setResponseModesSupported(List<String> responseModesSupported) {
    this.responseModesSupported = responseModesSupported;
  }

  public WellKnown responseTypesSupported(List<String> responseTypesSupported) {
    this.responseTypesSupported = responseTypesSupported;
    return this;
  }

  public WellKnown addResponseTypesSupportedItem(String responseTypesSupportedItem) {
    this.responseTypesSupported.add(responseTypesSupportedItem);
    return this;
  }

   /**
   * JSON array containing a list of the OAuth 2.0 response_type values that this OP supports. Dynamic OpenID Providers MUST support the code, id_token, and the token id_token Response Type values.
   * @return responseTypesSupported
  **/
  @ApiModelProperty(required = true, value = "JSON array containing a list of the OAuth 2.0 response_type values that this OP supports. Dynamic OpenID Providers MUST support the code, id_token, and the token id_token Response Type values.")
  public List<String> getResponseTypesSupported() {
    return responseTypesSupported;
  }

  public void setResponseTypesSupported(List<String> responseTypesSupported) {
    this.responseTypesSupported = responseTypesSupported;
  }

  public WellKnown revocationEndpoint(String revocationEndpoint) {
    this.revocationEndpoint = revocationEndpoint;
    return this;
  }

   /**
   * URL of the authorization server&#39;s OAuth 2.0 revocation endpoint.
   * @return revocationEndpoint
  **/
  @ApiModelProperty(value = "URL of the authorization server's OAuth 2.0 revocation endpoint.")
  public String getRevocationEndpoint() {
    return revocationEndpoint;
  }

  public void setRevocationEndpoint(String revocationEndpoint) {
    this.revocationEndpoint = revocationEndpoint;
  }

  public WellKnown scopesSupported(List<String> scopesSupported) {
    this.scopesSupported = scopesSupported;
    return this;
  }

  public WellKnown addScopesSupportedItem(String scopesSupportedItem) {
    if (this.scopesSupported == null) {
      this.scopesSupported = new ArrayList<String>();
    }
    this.scopesSupported.add(scopesSupportedItem);
    return this;
  }

   /**
   * SON array containing a list of the OAuth 2.0 [RFC6749] scope values that this server supports. The server MUST support the openid scope value. Servers MAY choose not to advertise some supported scope values even when this parameter is used
   * @return scopesSupported
  **/
  @ApiModelProperty(value = "SON array containing a list of the OAuth 2.0 [RFC6749] scope values that this server supports. The server MUST support the openid scope value. Servers MAY choose not to advertise some supported scope values even when this parameter is used")
  public List<String> getScopesSupported() {
    return scopesSupported;
  }

  public void setScopesSupported(List<String> scopesSupported) {
    this.scopesSupported = scopesSupported;
  }

  public WellKnown subjectTypesSupported(List<String> subjectTypesSupported) {
    this.subjectTypesSupported = subjectTypesSupported;
    return this;
  }

  public WellKnown addSubjectTypesSupportedItem(String subjectTypesSupportedItem) {
    this.subjectTypesSupported.add(subjectTypesSupportedItem);
    return this;
  }

   /**
   * JSON array containing a list of the Subject Identifier types that this OP supports. Valid types include pairwise and public.
   * @return subjectTypesSupported
  **/
  @ApiModelProperty(example = "\"public, pairwise\"", required = true, value = "JSON array containing a list of the Subject Identifier types that this OP supports. Valid types include pairwise and public.")
  public List<String> getSubjectTypesSupported() {
    return subjectTypesSupported;
  }

  public void setSubjectTypesSupported(List<String> subjectTypesSupported) {
    this.subjectTypesSupported = subjectTypesSupported;
  }

  public WellKnown tokenEndpoint(String tokenEndpoint) {
    this.tokenEndpoint = tokenEndpoint;
    return this;
  }

   /**
   * URL of the OP&#39;s OAuth 2.0 Token Endpoint
   * @return tokenEndpoint
  **/
  @ApiModelProperty(example = "https://playground.ory.sh/ory-hydra/public/oauth2/token", required = true, value = "URL of the OP's OAuth 2.0 Token Endpoint")
  public String getTokenEndpoint() {
    return tokenEndpoint;
  }

  public void setTokenEndpoint(String tokenEndpoint) {
    this.tokenEndpoint = tokenEndpoint;
  }

  public WellKnown tokenEndpointAuthMethodsSupported(List<String> tokenEndpointAuthMethodsSupported) {
    this.tokenEndpointAuthMethodsSupported = tokenEndpointAuthMethodsSupported;
    return this;
  }

  public WellKnown addTokenEndpointAuthMethodsSupportedItem(String tokenEndpointAuthMethodsSupportedItem) {
    if (this.tokenEndpointAuthMethodsSupported == null) {
      this.tokenEndpointAuthMethodsSupported = new ArrayList<String>();
    }
    this.tokenEndpointAuthMethodsSupported.add(tokenEndpointAuthMethodsSupportedItem);
    return this;
  }

   /**
   * JSON array containing a list of Client Authentication methods supported by this Token Endpoint. The options are client_secret_post, client_secret_basic, client_secret_jwt, and private_key_jwt, as described in Section 9 of OpenID Connect Core 1.0
   * @return tokenEndpointAuthMethodsSupported
  **/
  @ApiModelProperty(value = "JSON array containing a list of Client Authentication methods supported by this Token Endpoint. The options are client_secret_post, client_secret_basic, client_secret_jwt, and private_key_jwt, as described in Section 9 of OpenID Connect Core 1.0")
  public List<String> getTokenEndpointAuthMethodsSupported() {
    return tokenEndpointAuthMethodsSupported;
  }

  public void setTokenEndpointAuthMethodsSupported(List<String> tokenEndpointAuthMethodsSupported) {
    this.tokenEndpointAuthMethodsSupported = tokenEndpointAuthMethodsSupported;
  }

  public WellKnown userinfoEndpoint(String userinfoEndpoint) {
    this.userinfoEndpoint = userinfoEndpoint;
    return this;
  }

   /**
   * URL of the OP&#39;s UserInfo Endpoint.
   * @return userinfoEndpoint
  **/
  @ApiModelProperty(value = "URL of the OP's UserInfo Endpoint.")
  public String getUserinfoEndpoint() {
    return userinfoEndpoint;
  }

  public void setUserinfoEndpoint(String userinfoEndpoint) {
    this.userinfoEndpoint = userinfoEndpoint;
  }

  public WellKnown userinfoSigningAlgValuesSupported(List<String> userinfoSigningAlgValuesSupported) {
    this.userinfoSigningAlgValuesSupported = userinfoSigningAlgValuesSupported;
    return this;
  }

  public WellKnown addUserinfoSigningAlgValuesSupportedItem(String userinfoSigningAlgValuesSupportedItem) {
    if (this.userinfoSigningAlgValuesSupported == null) {
      this.userinfoSigningAlgValuesSupported = new ArrayList<String>();
    }
    this.userinfoSigningAlgValuesSupported.add(userinfoSigningAlgValuesSupportedItem);
    return this;
  }

   /**
   * JSON array containing a list of the JWS [JWS] signing algorithms (alg values) [JWA] supported by the UserInfo Endpoint to encode the Claims in a JWT [JWT].
   * @return userinfoSigningAlgValuesSupported
  **/
  @ApiModelProperty(value = "JSON array containing a list of the JWS [JWS] signing algorithms (alg values) [JWA] supported by the UserInfo Endpoint to encode the Claims in a JWT [JWT].")
  public List<String> getUserinfoSigningAlgValuesSupported() {
    return userinfoSigningAlgValuesSupported;
  }

  public void setUserinfoSigningAlgValuesSupported(List<String> userinfoSigningAlgValuesSupported) {
    this.userinfoSigningAlgValuesSupported = userinfoSigningAlgValuesSupported;
  }


  @Override
  public boolean equals(java.lang.Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    WellKnown wellKnown = (WellKnown) o;
    return Objects.equals(this.authorizationEndpoint, wellKnown.authorizationEndpoint) &&
        Objects.equals(this.backchannelLogoutSessionSupported, wellKnown.backchannelLogoutSessionSupported) &&
        Objects.equals(this.backchannelLogoutSupported, wellKnown.backchannelLogoutSupported) &&
        Objects.equals(this.claimsParameterSupported, wellKnown.claimsParameterSupported) &&
        Objects.equals(this.claimsSupported, wellKnown.claimsSupported) &&
        Objects.equals(this.endSessionEndpoint, wellKnown.endSessionEndpoint) &&
        Objects.equals(this.frontchannelLogoutSessionSupported, wellKnown.frontchannelLogoutSessionSupported) &&
        Objects.equals(this.frontchannelLogoutSupported, wellKnown.frontchannelLogoutSupported) &&
        Objects.equals(this.grantTypesSupported, wellKnown.grantTypesSupported) &&
        Objects.equals(this.idTokenSigningAlgValuesSupported, wellKnown.idTokenSigningAlgValuesSupported) &&
        Objects.equals(this.issuer, wellKnown.issuer) &&
        Objects.equals(this.jwksUri, wellKnown.jwksUri) &&
        Objects.equals(this.registrationEndpoint, wellKnown.registrationEndpoint) &&
        Objects.equals(this.requestParameterSupported, wellKnown.requestParameterSupported) &&
        Objects.equals(this.requestUriParameterSupported, wellKnown.requestUriParameterSupported) &&
        Objects.equals(this.requireRequestUriRegistration, wellKnown.requireRequestUriRegistration) &&
        Objects.equals(this.responseModesSupported, wellKnown.responseModesSupported) &&
        Objects.equals(this.responseTypesSupported, wellKnown.responseTypesSupported) &&
        Objects.equals(this.revocationEndpoint, wellKnown.revocationEndpoint) &&
        Objects.equals(this.scopesSupported, wellKnown.scopesSupported) &&
        Objects.equals(this.subjectTypesSupported, wellKnown.subjectTypesSupported) &&
        Objects.equals(this.tokenEndpoint, wellKnown.tokenEndpoint) &&
        Objects.equals(this.tokenEndpointAuthMethodsSupported, wellKnown.tokenEndpointAuthMethodsSupported) &&
        Objects.equals(this.userinfoEndpoint, wellKnown.userinfoEndpoint) &&
        Objects.equals(this.userinfoSigningAlgValuesSupported, wellKnown.userinfoSigningAlgValuesSupported);
  }

  @Override
  public int hashCode() {
    return Objects.hash(authorizationEndpoint, backchannelLogoutSessionSupported, backchannelLogoutSupported, claimsParameterSupported, claimsSupported, endSessionEndpoint, frontchannelLogoutSessionSupported, frontchannelLogoutSupported, grantTypesSupported, idTokenSigningAlgValuesSupported, issuer, jwksUri, registrationEndpoint, requestParameterSupported, requestUriParameterSupported, requireRequestUriRegistration, responseModesSupported, responseTypesSupported, revocationEndpoint, scopesSupported, subjectTypesSupported, tokenEndpoint, tokenEndpointAuthMethodsSupported, userinfoEndpoint, userinfoSigningAlgValuesSupported);
  }


  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class WellKnown {\n");
    
    sb.append("    authorizationEndpoint: ").append(toIndentedString(authorizationEndpoint)).append("\n");
    sb.append("    backchannelLogoutSessionSupported: ").append(toIndentedString(backchannelLogoutSessionSupported)).append("\n");
    sb.append("    backchannelLogoutSupported: ").append(toIndentedString(backchannelLogoutSupported)).append("\n");
    sb.append("    claimsParameterSupported: ").append(toIndentedString(claimsParameterSupported)).append("\n");
    sb.append("    claimsSupported: ").append(toIndentedString(claimsSupported)).append("\n");
    sb.append("    endSessionEndpoint: ").append(toIndentedString(endSessionEndpoint)).append("\n");
    sb.append("    frontchannelLogoutSessionSupported: ").append(toIndentedString(frontchannelLogoutSessionSupported)).append("\n");
    sb.append("    frontchannelLogoutSupported: ").append(toIndentedString(frontchannelLogoutSupported)).append("\n");
    sb.append("    grantTypesSupported: ").append(toIndentedString(grantTypesSupported)).append("\n");
    sb.append("    idTokenSigningAlgValuesSupported: ").append(toIndentedString(idTokenSigningAlgValuesSupported)).append("\n");
    sb.append("    issuer: ").append(toIndentedString(issuer)).append("\n");
    sb.append("    jwksUri: ").append(toIndentedString(jwksUri)).append("\n");
    sb.append("    registrationEndpoint: ").append(toIndentedString(registrationEndpoint)).append("\n");
    sb.append("    requestParameterSupported: ").append(toIndentedString(requestParameterSupported)).append("\n");
    sb.append("    requestUriParameterSupported: ").append(toIndentedString(requestUriParameterSupported)).append("\n");
    sb.append("    requireRequestUriRegistration: ").append(toIndentedString(requireRequestUriRegistration)).append("\n");
    sb.append("    responseModesSupported: ").append(toIndentedString(responseModesSupported)).append("\n");
    sb.append("    responseTypesSupported: ").append(toIndentedString(responseTypesSupported)).append("\n");
    sb.append("    revocationEndpoint: ").append(toIndentedString(revocationEndpoint)).append("\n");
    sb.append("    scopesSupported: ").append(toIndentedString(scopesSupported)).append("\n");
    sb.append("    subjectTypesSupported: ").append(toIndentedString(subjectTypesSupported)).append("\n");
    sb.append("    tokenEndpoint: ").append(toIndentedString(tokenEndpoint)).append("\n");
    sb.append("    tokenEndpointAuthMethodsSupported: ").append(toIndentedString(tokenEndpointAuthMethodsSupported)).append("\n");
    sb.append("    userinfoEndpoint: ").append(toIndentedString(userinfoEndpoint)).append("\n");
    sb.append("    userinfoSigningAlgValuesSupported: ").append(toIndentedString(userinfoSigningAlgValuesSupported)).append("\n");
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

