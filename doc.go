// Package main Hydra OAuth2 & OpenID Connect Server
//
// Please refer to the user guide for in-depth documentation: https://ory.gitbooks.io/hydra/content/
//
// Hydra offers OAuth 2.0 and OpenID Connect Core 1.0 capabilities as a service. Hydra is different, because it works with any existing authentication infrastructure, not just LDAP or SAML. By implementing a consent app (works with any programming language) you build a bridge between Hydra and your authentication infrastructure.
// Hydra is able to securely manage JSON Web Keys, and has a sophisticated policy-based access control you can use if you want to.
// Hydra is suitable for green- (new) and brownfield (existing) projects. If you are not familiar with OAuth 2.0 and are working on a greenfield project, we recommend evaluating if OAuth 2.0 really serves your purpose. Knowledge of OAuth 2.0 is imperative in understanding what Hydra does and how it works.
//
// The official repository is located at https://github.com/ory/hydra
//
//     Schemes: http, https
//     Host:
//     BasePath: /
//     Version: Latest
//     License: Apache 2.0 https://github.com/ory/hydra/blob/master/LICENSE
//     Contact: ORY <hi@ory.am> https://www.ory.am
//
//     Consumes:
//     - application/json
//     - application/x-www-form-urlencoded
//
//     Produces:
//     - application/json
//
//     SecurityDefinitions:
//     - oauth2:
//         type: oauth2
//         authorizationUrl: /oauth2/auth
//         tokenUrl: /oauth2/token
//         in: header
//         flow: accessCode
//
//     Extensions:
//     ---
//     x-request-id: string
//     ---
//
// swagger:meta
package main
