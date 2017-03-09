package fosite

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

// NewIntrospectionRequest initiates token introspection as defined in
// https://tools.ietf.org/search/rfc7662#section-2.1
//
// The protected resource calls the introspection endpoint using an HTTP
// POST [RFC7231] request with parameters sent as
// "application/x-www-form-urlencoded" data as defined in
// [W3C.REC-html5-20141028].  The protected resource sends a parameter
// representing the token along with optional parameters representing
// additional context that is known by the protected resource to aid the
// authorization server in its response.
//
// * token
// REQUIRED.  The string value of the token.  For access tokens, this
// is the "access_token" value returned from the token endpoint
// defined in OAuth 2.0 [RFC6749], Section 5.1.  For refresh tokens,
// this is the "refresh_token" value returned from the token endpoint
// as defined in OAuth 2.0 [RFC6749], Section 5.1.  Other token types
// are outside the scope of this specification.
//
// * token_type_hint
// OPTIONAL.  A hint about the type of the token submitted for
// introspection.  The protected resource MAY pass this parameter to
// help the authorization server optimize the token lookup.  If the
// server is unable to locate the token using the given hint, it MUST
// extend its search across all of its supported token types.  An
// authorization server MAY ignore this parameter, particularly if it
// is able to detect the token type automatically.  Values for this
// field are defined in the "OAuth Token Type Hints" registry defined
// in OAuth Token Revocation [RFC7009].
//
// The introspection endpoint MAY accept other OPTIONAL parameters to
// provide further context to the query.  For instance, an authorization
// server may desire to know the IP address of the client accessing the
// protected resource to determine if the correct client is likely to be
// presenting the token.  The definition of this or any other parameters
// are outside the scope of this specification, to be defined by service
// documentation or extensions to this specification.  If the
// authorization server is unable to determine the state of the token
// without additional information, it SHOULD return an introspection
// response indicating the token is not active as described in
// Section 2.2.
//
// To prevent token scanning attacks, the endpoint MUST also require
// some form of authorization to access this endpoint, such as client
// authentication as described in OAuth 2.0 [RFC6749] or a separate
// OAuth 2.0 access token such as the bearer token described in OAuth
// 2.0 Bearer Token Usage [RFC6750].  The methods of managing and
// validating these authentication credentials are out of scope of this
// specification.
//
// For example, the following shows a protected resource calling the
// token introspection endpoint to query about an OAuth 2.0 bearer
// token.  The protected resource is using a separate OAuth 2.0 bearer
// token to authorize this call.
//
// The following is a non-normative example request:
//
//	POST /introspect HTTP/1.1
//	Host: server.example.com
//	Accept: application/json
//	Content-Type: application/x-www-form-urlencoded
//	Authorization: Bearer 23410913-abewfq.123483
//
//	token=2YotnFZFEjr1zCsicMWpAA
//
// In this example, the protected resource uses a client identifier and
// client secret to authenticate itself to the introspection endpoint.
// The protected resource also sends a token type hint indicating that
// it is inquiring about an access token.
//
// The following is a non-normative example request:
//
//	POST /introspect HTTP/1.1
//	Host: server.example.com
//	Accept: application/json
//	Content-Type: application/x-www-form-urlencoded
//	Authorization: Basic czZCaGRSa3F0MzpnWDFmQmF0M2JW
//
//	token=mF_9.B5f-4.1JqM&token_type_hint=access_token
func (f *Fosite) NewIntrospectionRequest(ctx context.Context, r *http.Request, session Session) (IntrospectionResponder, error) {
	if r.Method != "POST" {
		return &IntrospectionResponse{Active: false}, errors.Wrap(ErrInvalidRequest, "HTTP method is not POST")
	} else if err := r.ParseForm(); err != nil {
		return &IntrospectionResponse{Active: false}, errors.Wrap(ErrInvalidRequest, err.Error())
	}

	token := r.PostForm.Get("token")
	tokenType := r.PostForm.Get("token_type_hint")
	scope := r.PostForm.Get("scope")
	if clientToken := AccessTokenFromRequest(r); clientToken != "" {
		if token == clientToken {
			return &IntrospectionResponse{Active: false}, errors.Wrap(ErrRequestUnauthorized, "Bearer and introspection token are identical")
		}

		if _, err := f.IntrospectToken(ctx, clientToken, AccessToken, session.Clone()); err != nil {
			return &IntrospectionResponse{Active: false}, errors.Wrap(ErrRequestUnauthorized, "HTTP Authorization header missing, malformed or credentials used are invalid")
		}
	} else {
		clientID, clientSecret, ok := r.BasicAuth()
		if !ok {
			return &IntrospectionResponse{Active: false}, errors.Wrap(ErrRequestUnauthorized, "HTTP Authorization header missing, malformed or credentials used are invalid")
		}

		client, err := f.Store.GetClient(clientID)
		if err != nil {
			return &IntrospectionResponse{Active: false}, errors.Wrap(ErrRequestUnauthorized, "HTTP Authorization header missing, malformed or credentials used are invalid")
		}

		// Enforce client authentication
		if err := f.Hasher.Compare(client.GetHashedSecret(), []byte(clientSecret)); err != nil {
			return &IntrospectionResponse{Active: false}, errors.Wrap(ErrRequestUnauthorized, "HTTP Authorization header missing, malformed or credentials used are invalid")
		}
	}

	ar, err := f.IntrospectToken(ctx, token, TokenType(tokenType), session, strings.Split(scope, " ")...)
	if err != nil {
		return &IntrospectionResponse{Active: false}, err
	}

	return &IntrospectionResponse{
		Active:          true,
		AccessRequester: ar,
	}, nil
}

type IntrospectionResponse struct {
	Active          bool            `json:"active"`
	AccessRequester AccessRequester `json:",extra"`
}

func (r *IntrospectionResponse) IsActive() bool {
	return r.Active
}

func (r *IntrospectionResponse) GetAccessRequester() AccessRequester {
	return r.AccessRequester
}
