package fosite

import (
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/context"
)

const MinParameterEntropy = 8

type TokenType string

const (
	AccessToken   TokenType = "access_token"
	RefreshToken  TokenType = "refresh_token"
	AuthorizeCode TokenType = "authorize_code"
	IDToken       TokenType = "id_token"
)

// OAuth2Provider is an interface that enables you to write OAuth2 handlers with only a few lines of code.
// Check fosite.Fosite for an implementation of this interface.
type OAuth2Provider interface {
	// NewAuthorizeRequest returns an AuthorizeRequest.
	//
	// The following specs must be considered in any implementation of this method:
	// * https://tools.ietf.org/html/rfc6749#section-3.1
	//	 Extension response types MAY contain a space-delimited (%x20) list of
	//	 values, where the order of values does not matter (e.g., response
	//	 type "a b" is the same as "b a").  The meaning of such composite
	//	 response types is defined by their respective specifications.
	// * https://tools.ietf.org/html/rfc6749#section-3.1.2
	//   The redirection endpoint URI MUST be an absolute URI as defined by
	//   [RFC3986] Section 4.3.  The endpoint URI MAY include an
	//   "application/x-www-form-urlencoded" formatted (per Appendix B) query
	//   component ([RFC3986] Section 3.4), which MUST be retained when adding
	//   additional query parameters.  The endpoint URI MUST NOT include a
	//   fragment component.
	// * https://tools.ietf.org/html/rfc6749#section-3.1.2.2 (everything MUST be implemented)
	NewAuthorizeRequest(ctx context.Context, req *http.Request) (AuthorizeRequester, error)

	// NewAuthorizeResponse iterates through all response type handlers and returns their result or
	// ErrUnsupportedResponseType if none of the handler's were able to handle it.
	//
	// The following specs must be considered in any implementation of this method:
	// * https://tools.ietf.org/html/rfc6749#section-3.1.1
	//	 Extension response types MAY contain a space-delimited (%x20) list of
	//	 values, where the order of values does not matter (e.g., response
	//	 type "a b" is the same as "b a").  The meaning of such composite
	//	 response types is defined by their respective specifications.
	//	 If an authorization request is missing the "response_type" parameter,
	//	 or if the response type is not understood, the authorization server
	//	 MUST return an error response as described in Section 4.1.2.1.
	NewAuthorizeResponse(ctx context.Context, req *http.Request, requester AuthorizeRequester, session Session) (AuthorizeResponder, error)

	// WriteAuthorizeError returns the error codes to the redirection endpoint or shows the error to the user, if no valid
	// redirect uri was given. Implements rfc6749#section-4.1.2.1
	//
	// The following specs must be considered in any implementation of this method:
	// * https://tools.ietf.org/html/rfc6749#section-3.1.2
	//   The redirection endpoint URI MUST be an absolute URI as defined by
	//   [RFC3986] Section 4.3.  The endpoint URI MAY include an
	//   "application/x-www-form-urlencoded" formatted (per Appendix B) query
	//   component ([RFC3986] Section 3.4), which MUST be retained when adding
	//   additional query parameters.  The endpoint URI MUST NOT include a
	//   fragment component.
	// * https://tools.ietf.org/html/rfc6749#section-4.1.2.1 (everything)
	// * https://tools.ietf.org/html/rfc6749#section-3.1.2.2 (everything MUST be implemented)
	WriteAuthorizeError(rw http.ResponseWriter, requester AuthorizeRequester, err error)

	// WriteAuthorizeResponse persists the AuthorizeSession in the store and redirects the user agent to the provided
	// redirect url or returns an error if storage failed.
	//
	// The following specs must be considered in any implementation of this method:
	// * https://tools.ietf.org/html/rfc6749#rfc6749#section-4.1.2.1
	//   After completing its interaction with the resource owner, the
	//   authorization server directs the resource owner's user-agent back to
	//   the client.  The authorization server redirects the user-agent to the
	//   client's redirection endpoint previously established with the
	//   authorization server during the client registration process or when
	//   making the authorization request.
	// * https://tools.ietf.org/html/rfc6749#section-3.1.2.2 (everything MUST be implemented)
	WriteAuthorizeResponse(rw http.ResponseWriter, requester AuthorizeRequester, responder AuthorizeResponder)

	// NewAccessRequest creates a new access request object and validates
	// various parameters.
	//
	// The following specs must be considered in any implementation of this method:
	// * https://tools.ietf.org/html/rfc6749#section-3.2 (everything)
	// * https://tools.ietf.org/html/rfc6749#section-3.2.1 (everything)
	//
	// Furthermore the registered handlers should implement their specs accordingly.
	NewAccessRequest(ctx context.Context, req *http.Request, session Session) (AccessRequester, error)

	// NewAccessResponse creates a new access response and validates that access_token and token_type are set.
	//
	// The following specs must be considered in any implementation of this method:
	// https://tools.ietf.org/html/rfc6749#section-5.1
	NewAccessResponse(ctx context.Context, req *http.Request, requester AccessRequester) (AccessResponder, error)

	// WriteAccessError writes an access request error response.
	//
	// The following specs must be considered in any implementation of this method:
	// * https://tools.ietf.org/html/rfc6749#section-5.2 (everything)
	WriteAccessError(rw http.ResponseWriter, requester AccessRequester, err error)

	// WriteAccessResponse writes the access response.
	//
	// The following specs must be considered in any implementation of this method:
	// https://tools.ietf.org/html/rfc6749#section-5.1
	WriteAccessResponse(rw http.ResponseWriter, requester AccessRequester, responder AccessResponder)

	// NewRevocationRequest handles incoming token revocation requests and validates various parameters.
	//
	// The following specs must be considered in any implementation of this method:
	// https://tools.ietf.org/html/rfc7009#section-2.1
	NewRevocationRequest(ctx context.Context, r *http.Request) error

	// WriteRevocationResponse writes the revoke response.
	//
	// The following specs must be considered in any implementation of this method:
	// https://tools.ietf.org/html/rfc7009#section-2.2
	WriteRevocationResponse(rw http.ResponseWriter, err error)

	// IntrospectToken returns token metadata, if the token is valid. Tokens generated by the authorization endpoint,
	// such as the authorization code, can not be introspected.
	IntrospectToken(ctx context.Context, token string, tokenType TokenType, session Session, scope ...string) (AccessRequester, error)

	// NewIntrospectionRequest initiates token introspection as defined in
	// https://tools.ietf.org/search/rfc7662#section-2.1
	NewIntrospectionRequest(ctx context.Context, r *http.Request, session Session) (IntrospectionResponder, error)

	// WriteIntrospectionError responds with an error if token introspection failed as defined in
	// https://tools.ietf.org/search/rfc7662#section-2.3
	WriteIntrospectionError(rw http.ResponseWriter, err error)

	// WriteIntrospectionResponse responds with token metadata discovered by token introspection as defined in
	// https://tools.ietf.org/search/rfc7662#section-2.2
	WriteIntrospectionResponse(rw http.ResponseWriter, r IntrospectionResponder)
}

// IntrospectionResponse is the response object that will be returned when token introspection was successful,
// for example when the client is allowed to perform token introspection. Refer to
// https://tools.ietf.org/search/rfc7662#section-2.2 for more details.
type IntrospectionResponder interface {
	// IsActive returns true if the introspected token is active and false otherwise.
	IsActive() bool

	// AccessRequester returns nil when IsActive() is false and the original access request object otherwise.
	GetAccessRequester() AccessRequester
}

// Requester is an abstract interface for handling requests in Fosite.
type Requester interface {
	// GetID returns a unique identifier.
	GetID() string

	// GetRequestedAt returns the time the request was created.
	GetRequestedAt() (requestedAt time.Time)

	// GetClient returns the requests client.
	GetClient() (client Client)

	// GetRequestedScopes returns the request's scopes.
	GetRequestedScopes() (scopes Arguments)

	// SetRequestedScopes sets the request's scopes.
	SetRequestedScopes(scopes Arguments)

	// AppendRequestedScope appends a scope to the request.
	AppendRequestedScope(scope string)

	// GetGrantScopes returns all granted scopes.
	GetGrantedScopes() (grantedScopes Arguments)

	// GrantScope marks a request's scope as granted.
	GrantScope(scope string)

	// GetSession returns a pointer to the request's session or nil if none is set.
	GetSession() (session Session)

	// GetSession sets the request's session pointer.
	SetSession(session Session)

	// GetRequestForm returns the request's form input.
	GetRequestForm() url.Values

	Merge(requester Requester)
}

// AccessRequester is a token endpoint's request context.
type AccessRequester interface {
	// GetGrantType returns the requests grant type.
	GetGrantTypes() (grantTypes Arguments)

	Requester
}

// AuthorizeRequester is an authorize endpoint's request context.
type AuthorizeRequester interface {
	// GetResponseTypes returns the requested response types
	GetResponseTypes() (responseTypes Arguments)

	// SetResponseTypeHandled marks a response_type (e.g. token or code) as handled indicating that the response type
	// is supported.
	SetResponseTypeHandled(responseType string)

	// DidHandleAllResponseTypes returns if all requested response types have been handled correctly
	DidHandleAllResponseTypes() (didHandle bool)

	// GetRedirectURI returns the requested redirect URI
	GetRedirectURI() (redirectURL *url.URL)

	// IsRedirectURIValid returns false if the redirect is not rfc-conform (i.e. missing client, not on white list,
	// or malformed)
	IsRedirectURIValid() (isValid bool)

	// GetState returns the request's state.
	GetState() (state string)

	Requester
}

// AccessResponder is a token endpoint's response.
type AccessResponder interface {
	// SetExtra sets a key value pair for the access response.
	SetExtra(key string, value interface{})

	// GetExtra returns a key's value.
	GetExtra(key string) interface{}

	SetExpiresIn(time.Duration)

	SetScopes(scopes Arguments)

	// SetAccessToken sets the responses mandatory access token.
	SetAccessToken(token string)

	// SetTokenType set's the responses mandatory token type
	SetTokenType(tokenType string)

	// SetAccessToken returns the responses access token.
	GetAccessToken() (token string)

	// GetTokenType returns the responses token type.
	GetTokenType() (token string)

	// ToMap converts the response to a map.
	ToMap() map[string]interface{}
}

// AuthorizeResponder is an authorization endpoint's response.
type AuthorizeResponder interface {
	// GetCode returns the response's authorize code if set.
	GetCode() string

	// GetHeader returns the response's header
	GetHeader() (header http.Header)

	// AddHeader adds an header key value pair to the response
	AddHeader(key, value string)

	// GetQuery returns the response's query
	GetQuery() (query url.Values)

	// AddQuery adds an url query key value pair to the response
	AddQuery(key, value string)

	// GetHeader returns the response's url fragments
	GetFragment() (fragment url.Values)

	// AddHeader adds a key value pair to the response's url fragment
	AddFragment(key, value string)
}
