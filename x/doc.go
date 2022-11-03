// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

// ORY Hydra
//
// Welcome to the ORY Hydra HTTP API documentation. You will find documentation for all HTTP APIs here.
//
//	Schemes: http, https
//	Host:
//	BasePath: /
//	Version: latest
//
//	Consumes:
//	- application/json
//	- application/x-www-form-urlencoded
//
//	Produces:
//	- application/json
//
//	SecurityDefinitions:
//	oauth2:
//	    type: oauth2
//	    authorizationUrl: https://hydra.demo.ory.sh/oauth2/auth
//	    tokenUrl: https://hydra.demo.ory.sh/oauth2/token
//	    flow: accessCode
//	    scopes:
//	      offline: "A scope required when requesting refresh tokens (alias for `offline_access`)"
//	      offline_access: "A scope required when requesting refresh tokens"
//	      openid: "Request an OpenID Connect ID Token"
//	basic:
//	    type: basic
//	bearer:
//	    type: basic
//
//	Extensions:
//	---
//	x-request-id: string
//	x-forwarded-proto: string
//	---
//
// swagger:meta
package x
