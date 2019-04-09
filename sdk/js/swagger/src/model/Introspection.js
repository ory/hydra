/**
 * ORY Hydra
 * Welcome to the ORY Hydra HTTP API documentation. You will find documentation for all HTTP APIs here.
 *
 * OpenAPI spec version: latest
 *
 * NOTE: This class is auto generated by the swagger code generator program.
 * https://github.com/swagger-api/swagger-codegen.git
 *
 * Swagger Codegen version: 2.2.3
 *
 * Do not edit the class manually.
 *
 */

(function(root, factory) {
  if (typeof define === 'function' && define.amd) {
    // AMD. Register as an anonymous module.
    define(['ApiClient'], factory);
  } else if (typeof module === 'object' && module.exports) {
    // CommonJS-like environments that support module.exports, like Node.
    module.exports = factory(require('../ApiClient'));
  } else {
    // Browser globals (root is window)
    if (!root.OryHydra) {
      root.OryHydra = {};
    }
    root.OryHydra.Introspection = factory(root.OryHydra.ApiClient);
  }
}(this, function(ApiClient) {
  'use strict';




  /**
   * The Introspection model module.
   * @module model/Introspection
   * @version latest
   */

  /**
   * Constructs a new <code>Introspection</code>.
   * https://tools.ietf.org/html/rfc7662
   * @alias module:model/Introspection
   * @class
   * @param active {Boolean} Active is a boolean indicator of whether or not the presented token is currently active.  The specifics of a token's \"active\" state will vary depending on the implementation of the authorization server and the information it keeps about its tokens, but a \"true\" value return for the \"active\" property will generally indicate that a given token has been issued by this authorization server, has not been revoked by the resource owner, and is within its given time window of validity (e.g., after its issuance time and before its expiration time).
   */
  var exports = function(active) {
    var _this = this;

    _this['active'] = active;












  };

  /**
   * Constructs a <code>Introspection</code> from a plain JavaScript object, optionally creating a new instance.
   * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
   * @param {Object} data The plain JavaScript object bearing properties of interest.
   * @param {module:model/Introspection} obj Optional instance to populate.
   * @return {module:model/Introspection} The populated <code>Introspection</code> instance.
   */
  exports.constructFromObject = function(data, obj) {
    if (data) {
      obj = obj || new exports();

      if (data.hasOwnProperty('active')) {
        obj['active'] = ApiClient.convertToType(data['active'], 'Boolean');
      }
      if (data.hasOwnProperty('aud')) {
        obj['aud'] = ApiClient.convertToType(data['aud'], ['String']);
      }
      if (data.hasOwnProperty('client_id')) {
        obj['client_id'] = ApiClient.convertToType(data['client_id'], 'String');
      }
      if (data.hasOwnProperty('exp')) {
        obj['exp'] = ApiClient.convertToType(data['exp'], 'Number');
      }
      if (data.hasOwnProperty('ext')) {
        obj['ext'] = ApiClient.convertToType(data['ext'], {'String': Object});
      }
      if (data.hasOwnProperty('iat')) {
        obj['iat'] = ApiClient.convertToType(data['iat'], 'Number');
      }
      if (data.hasOwnProperty('iss')) {
        obj['iss'] = ApiClient.convertToType(data['iss'], 'String');
      }
      if (data.hasOwnProperty('nbf')) {
        obj['nbf'] = ApiClient.convertToType(data['nbf'], 'Number');
      }
      if (data.hasOwnProperty('obfuscated_subject')) {
        obj['obfuscated_subject'] = ApiClient.convertToType(data['obfuscated_subject'], 'String');
      }
      if (data.hasOwnProperty('scope')) {
        obj['scope'] = ApiClient.convertToType(data['scope'], 'String');
      }
      if (data.hasOwnProperty('sub')) {
        obj['sub'] = ApiClient.convertToType(data['sub'], 'String');
      }
      if (data.hasOwnProperty('token_type')) {
        obj['token_type'] = ApiClient.convertToType(data['token_type'], 'String');
      }
      if (data.hasOwnProperty('username')) {
        obj['username'] = ApiClient.convertToType(data['username'], 'String');
      }
    }
    return obj;
  }

  /**
   * Active is a boolean indicator of whether or not the presented token is currently active.  The specifics of a token's \"active\" state will vary depending on the implementation of the authorization server and the information it keeps about its tokens, but a \"true\" value return for the \"active\" property will generally indicate that a given token has been issued by this authorization server, has not been revoked by the resource owner, and is within its given time window of validity (e.g., after its issuance time and before its expiration time).
   * @member {Boolean} active
   */
  exports.prototype['active'] = undefined;
  /**
   * Audience contains a list of the token's intended audiences.
   * @member {Array.<String>} aud
   */
  exports.prototype['aud'] = undefined;
  /**
   * ClientID is aclient identifier for the OAuth 2.0 client that requested this token.
   * @member {String} client_id
   */
  exports.prototype['client_id'] = undefined;
  /**
   * Expires at is an integer timestamp, measured in the number of seconds since January 1 1970 UTC, indicating when this token will expire.
   * @member {Number} exp
   */
  exports.prototype['exp'] = undefined;
  /**
   * Extra is arbitrary data set by the session.
   * @member {Object.<String, Object>} ext
   */
  exports.prototype['ext'] = undefined;
  /**
   * Issued at is an integer timestamp, measured in the number of seconds since January 1 1970 UTC, indicating when this token was originally issued.
   * @member {Number} iat
   */
  exports.prototype['iat'] = undefined;
  /**
   * IssuerURL is a string representing the issuer of this token
   * @member {String} iss
   */
  exports.prototype['iss'] = undefined;
  /**
   * NotBefore is an integer timestamp, measured in the number of seconds since January 1 1970 UTC, indicating when this token is not to be used before.
   * @member {Number} nbf
   */
  exports.prototype['nbf'] = undefined;
  /**
   * ObfuscatedSubject is set when the subject identifier algorithm was set to \"pairwise\" during authorization. It is the `sub` value of the ID Token that was issued.
   * @member {String} obfuscated_subject
   */
  exports.prototype['obfuscated_subject'] = undefined;
  /**
   * Scope is a JSON string containing a space-separated list of scopes associated with this token.
   * @member {String} scope
   */
  exports.prototype['scope'] = undefined;
  /**
   * Subject of the token, as defined in JWT [RFC7519]. Usually a machine-readable identifier of the resource owner who authorized this token.
   * @member {String} sub
   */
  exports.prototype['sub'] = undefined;
  /**
   * TokenType is the introspected token's type, for example `access_token` or `refresh_token`.
   * @member {String} token_type
   */
  exports.prototype['token_type'] = undefined;
  /**
   * Username is a human-readable identifier for the resource owner who authorized this token.
   * @member {String} username
   */
  exports.prototype['username'] = undefined;



  return exports;
}));


