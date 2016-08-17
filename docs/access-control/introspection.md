# OAuth2 Token Introspection

OAuth2 Token Introspection is an [IETF](https://tools.ietf.org/html/rfc7662) standard.
It defines a method for a protected resource to query
an OAuth 2.0 authorization server to determine the active state of an
OAuth 2.0 token and to determine meta-information about this token.
OAuth 2.0 deployments can use this method to convey information about
the authorization context of the token from the authorization server
to the protected resource.

In order to make a successful Token Introspection request, the audience of the access token you are introspecting
*must* match the subject of the access token you are using to access the introspection endpoint.

The Token Introspection endpoint is documented in more detail [here](http://docs.hdyra.apiary.io/#reference/oauth2/oauth2-token-introspection).