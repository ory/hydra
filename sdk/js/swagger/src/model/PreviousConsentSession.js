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
    define(['ApiClient', 'model/ConsentRequest', 'model/ConsentRequestSessionData'], factory);
  } else if (typeof module === 'object' && module.exports) {
    // CommonJS-like environments that support module.exports, like Node.
    module.exports = factory(require('../ApiClient'), require('./ConsentRequest'), require('./ConsentRequestSessionData'));
  } else {
    // Browser globals (root is window)
    if (!root.OryHydra) {
      root.OryHydra = {};
    }
    root.OryHydra.PreviousConsentSession = factory(root.OryHydra.ApiClient, root.OryHydra.ConsentRequest, root.OryHydra.ConsentRequestSessionData);
  }
}(this, function(ApiClient, ConsentRequest, ConsentRequestSessionData) {
  'use strict';




  /**
   * The PreviousConsentSession model module.
   * @module model/PreviousConsentSession
   * @version latest
   */

  /**
   * Constructs a new <code>PreviousConsentSession</code>.
   * PreviousConsentSession The response used to return handled consent requests same as HandledAuthenticationRequest, just with consent_request exposed as json
   * @alias module:model/PreviousConsentSession
   * @class
   */
  var exports = function() {
    var _this = this;







  };

  /**
   * Constructs a <code>PreviousConsentSession</code> from a plain JavaScript object, optionally creating a new instance.
   * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
   * @param {Object} data The plain JavaScript object bearing properties of interest.
   * @param {module:model/PreviousConsentSession} obj Optional instance to populate.
   * @return {module:model/PreviousConsentSession} The populated <code>PreviousConsentSession</code> instance.
   */
  exports.constructFromObject = function(data, obj) {
    if (data) {
      obj = obj || new exports();

      if (data.hasOwnProperty('consent_request')) {
        obj['consent_request'] = ConsentRequest.constructFromObject(data['consent_request']);
      }
      if (data.hasOwnProperty('grant_access_token_audience')) {
        obj['grant_access_token_audience'] = ApiClient.convertToType(data['grant_access_token_audience'], ['String']);
      }
      if (data.hasOwnProperty('grant_scope')) {
        obj['grant_scope'] = ApiClient.convertToType(data['grant_scope'], ['String']);
      }
      if (data.hasOwnProperty('remember')) {
        obj['remember'] = ApiClient.convertToType(data['remember'], 'Boolean');
      }
      if (data.hasOwnProperty('remember_for')) {
        obj['remember_for'] = ApiClient.convertToType(data['remember_for'], 'Number');
      }
      if (data.hasOwnProperty('session')) {
        obj['session'] = ConsentRequestSessionData.constructFromObject(data['session']);
      }
    }
    return obj;
  }

  /**
   * @member {module:model/ConsentRequest} consent_request
   */
  exports.prototype['consent_request'] = undefined;
  /**
   * GrantedAudience sets the audience the user authorized the client to use. Should be a subset of `requested_access_token_audience`.
   * @member {Array.<String>} grant_access_token_audience
   */
  exports.prototype['grant_access_token_audience'] = undefined;
  /**
   * GrantScope sets the scope the user authorized the client to use. Should be a subset of `requested_scope`
   * @member {Array.<String>} grant_scope
   */
  exports.prototype['grant_scope'] = undefined;
  /**
   * Remember, if set to true, tells ORY Hydra to remember this consent authorization and reuse it if the same client asks the same user for the same, or a subset of, scope.
   * @member {Boolean} remember
   */
  exports.prototype['remember'] = undefined;
  /**
   * RememberFor sets how long the consent authorization should be remembered for in seconds. If set to `0`, the authorization will be remembered indefinitely.
   * @member {Number} remember_for
   */
  exports.prototype['remember_for'] = undefined;
  /**
   * @member {module:model/ConsentRequestSessionData} session
   */
  exports.prototype['session'] = undefined;



  return exports;
}));


