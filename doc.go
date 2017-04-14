// Package hydra is a server implementation of the OAuth 2.0 authorization framework and the OpenID Connect Core 1.0. Existing OAuth2 implementations usually ship as libraries or SDKs such as node-oauth2-server or fosite, or as fully featured identity solutions with user management and user interfaces, such as Dex.
//
// Implementing and using OAuth2 without understanding the whole specification is challenging and prone to errors, even when SDKs are being used. The primary goal of Hydra is to make OAuth 2.0 and OpenID Connect 1.0 better accessible.
//
// Hydra implements the flows described in OAuth2 and OpenID Connect 1.0 without forcing you to use a "Hydra User Management" or some template engine or a predefined front-end. Instead it relies on HTTP redirection and cryptographic methods to verify user consent allowing you to use Hydra with any authentication endpoint, be it authboss, auth0.com or your proprietary PHP authentication.
//
// The official repository is located at https://github.com/ory-am/hydra
package main
