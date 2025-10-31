// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/url"
	"strings"

	"github.com/ory/x/errorsx"

	"github.com/asaskevich/govalidator"
)

var DefaultFormPostTemplate = template.Must(template.New("form_post").Parse(`<html>
   <head>
      <title>Submit This Form</title>
   </head>
   <body onload="javascript:document.forms[0].submit()">
      <form method="post" action="{{ .RedirURL }}">
         {{ range $key,$value := .Parameters }}
            {{ range $parameter:= $value}}
		      <input type="hidden" name="{{$key}}" value="{{$parameter}}"/>
            {{end}}
         {{ end }}
      </form>
   </body>
</html>`))

// MatchRedirectURIWithClientRedirectURIs if the given uri is a registered redirect uri. Does not perform
// uri validation.
//
// Considered specifications
//
//   - https://tools.ietf.org/html/rfc6749#section-3.1.2.3
//     If multiple redirection URIs have been registered, if only part of
//     the redirection URI has been registered, or if no redirection URI has
//     been registered, the client MUST include a redirection URI with the
//     authorization request using the "redirect_uri" request parameter.
//
//     When a redirection URI is included in an authorization request, the
//     authorization server MUST compare and match the value received
//     against at least one of the registered redirection URIs (or URI
//     components) as defined in [RFC3986] Section 6, if any redirection
//     URIs were registered.  If the client registration included the full
//     redirection URI, the authorization server MUST compare the two URIs
//     using simple string comparison as defined in [RFC3986] Section 6.2.1.
//
// * https://tools.ietf.org/html/rfc6819#section-4.4.1.7
//   - The authorization server may also enforce the usage and validation
//     of pre-registered redirect URIs (see Section 5.2.3.5).  This will
//     allow for early recognition of authorization "code" disclosure to
//     counterfeit clients.
//   - The attacker will need to use another redirect URI for its
//     authorization process rather than the target web site because it
//     needs to intercept the flow.  So, if the authorization server
//     associates the authorization "code" with the redirect URI of a
//     particular end-user authorization and validates this redirect URI
//     with the redirect URI passed to the token's endpoint, such an
//     attack is detected (see Section 5.2.4.5).
func MatchRedirectURIWithClientRedirectURIs(rawurl string, client Client) (*url.URL, error) {
	if rawurl == "" && len(client.GetRedirectURIs()) == 1 {
		if redirectURIFromClient, err := url.Parse(client.GetRedirectURIs()[0]); err == nil && IsValidRedirectURI(redirectURIFromClient) {
			// If no redirect_uri was given and the client has exactly one valid redirect_uri registered, use that instead
			return redirectURIFromClient, nil
		}
	} else if redirectTo, ok := isMatchingRedirectURI(rawurl, client.GetRedirectURIs()); rawurl != "" && ok {
		// If a redirect_uri was given and the clients knows it (simple string comparison!)
		// return it.
		if parsed, err := url.Parse(redirectTo); err == nil && IsValidRedirectURI(parsed) {
			// If no redirect_uri was given and the client has exactly one valid redirect_uri registered, use that instead
			return parsed, nil
		}
	}

	return nil, errorsx.WithStack(ErrInvalidRequest.WithHint("The 'redirect_uri' parameter does not match any of the OAuth 2.0 Client's pre-registered redirect urls."))
}

// Match a requested  redirect URI against a pool of registered client URIs
//
// Test a given redirect URI against a pool of URIs provided by a registered client.
// If the OAuth 2.0 Client has loopback URIs registered either an IPv4 URI http://127.0.0.1 or
// an IPv6 URI http://[::1] a client is allowed to request a dynamic port and the server MUST accept
// it as a valid redirection uri.
//
// https://tools.ietf.org/html/rfc8252#section-7.3
// Native apps that are able to open a port on the loopback network
// interface without needing special permissions (typically, those on
// desktop operating systems) can use the loopback interface to receive
// the OAuth redirect.
//
// Loopback redirect URIs use the "http" scheme and are constructed with
// the loopback IP literal and whatever port the client is listening on.
func isMatchingRedirectURI(uri string, haystack []string) (string, bool) {
	requested, err := url.Parse(uri)
	if err != nil {
		return "", false
	}

	for _, b := range haystack {
		if b == uri {
			return b, true
		} else if isMatchingAsLoopback(requested, b) {
			// We have to return the requested URL here because otherwise the port might get lost (see isMatchingAsLoopback)
			// description.
			return uri, true
		}
	}
	return "", false
}

func isMatchingAsLoopback(requested *url.URL, registeredURI string) bool {
	registered, err := url.Parse(registeredURI)
	if err != nil {
		return false
	}

	// Native apps that are able to open a port on the loopback network
	// interface without needing special permissions (typically, those on
	// desktop operating systems) can use the loopback interface to receive
	// the OAuth redirect.
	//
	// Loopback redirect URIs use the "http" scheme and are constructed with
	// the loopback IP literal and whatever port the client is listening on.
	//
	// Source: https://tools.ietf.org/html/rfc8252#section-7.3
	if requested.Scheme == "http" &&
		isLoopbackAddress(requested.Hostname()) &&
		registered.Hostname() == requested.Hostname() &&
		// The port is skipped here - see codedoc above!
		registered.Path == requested.Path &&
		registered.RawQuery == requested.RawQuery {
		return true
	}

	return false
}

// Check if address is either an IPv4 loopback or an IPv6 loopback.
func isLoopbackAddress(hostname string) bool {
	return net.ParseIP(hostname).IsLoopback()
}

// IsValidRedirectURI validates a redirect_uri as specified in:
//
// * https://tools.ietf.org/html/rfc6749#section-3.1.2
//   - The redirection endpoint URI MUST be an absolute URI as defined by [RFC3986] Section 4.3.
//   - The endpoint URI MUST NOT include a fragment component.
//   - https://tools.ietf.org/html/rfc3986#section-4.3
//     absolute-URI  = scheme ":" hier-part [ "?" query ]
//   - https://tools.ietf.org/html/rfc6819#section-5.1.1
func IsValidRedirectURI(redirectURI *url.URL) bool {
	// We need to explicitly check for a scheme
	if !govalidator.IsRequestURL(redirectURI.String()) {
		return false
	}

	if redirectURI.Fragment != "" {
		// "The endpoint URI MUST NOT include a fragment component."
		return false
	}

	return true
}

func IsRedirectURISecure(ctx context.Context, redirectURI *url.URL) bool {
	return !(redirectURI.Scheme == "http" && !IsLocalhost(redirectURI))
}

// IsRedirectURISecureStrict is stricter than IsRedirectURISecure and it does not allow custom-scheme
// URLs because they can be hijacked for native apps. Use claimed HTTPS redirects instead.
// See discussion in https://github.com/ory/hydra/v2/fosite/pull/489.
func IsRedirectURISecureStrict(ctx context.Context, redirectURI *url.URL) bool {
	return redirectURI.Scheme == "https" || (redirectURI.Scheme == "http" && IsLocalhost(redirectURI))
}

func IsLocalhost(redirectURI *url.URL) bool {
	hn := redirectURI.Hostname()
	return strings.HasSuffix(hn, ".localhost") || isLoopbackAddress(hn) || hn == "localhost"
}

func WriteAuthorizeFormPostResponse(redirectURL string, parameters url.Values, template *template.Template, rw io.Writer) {
	_ = template.Execute(rw, struct {
		RedirURL   string
		Parameters url.Values
	}{
		RedirURL:   redirectURL,
		Parameters: parameters,
	})
}

// Deprecated: Do not use.
func URLSetFragment(source *url.URL, fragment url.Values) {
	var f string
	for k, v := range fragment {
		for _, vv := range v {
			if len(f) != 0 {
				f += fmt.Sprintf("&%s=%s", k, vv)
			} else {
				f += fmt.Sprintf("%s=%s", k, vv)
			}
		}
	}
	source.Fragment = f
}

func GetPostFormHTMLTemplate(ctx context.Context, f *Fosite) *template.Template {
	if t := f.Config.GetFormPostHTMLTemplate(ctx); t != nil {
		return t
	}
	return DefaultFormPostTemplate
}
