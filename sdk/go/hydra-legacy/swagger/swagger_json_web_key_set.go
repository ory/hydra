/*
 * ORY Hydra
 *
 * Welcome to the ORY Hydra HTTP API documentation. You will find documentation for all HTTP APIs here.
 *
 * OpenAPI spec version: latest
 *
 * Generated by: https://github.com/swagger-api/swagger-codegen.git
 */

package swagger

// SwaggerJSONWebKeySet swagger JSON web key set
type SwaggerJsonWebKeySet struct {

	// The value of the \"keys\" parameter is an array of JWK values.  By default, the order of the JWK values within the array does not imply an order of preference among them, although applications of JWK Sets can choose to assign a meaning to the order for their purposes, if desired.
	Keys []SwaggerJsonWebKey `json:"keys,omitempty"`
}
