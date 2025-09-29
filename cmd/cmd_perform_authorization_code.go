// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/toqueteos/webbrowser"
	"golang.org/x/oauth2"

	"github.com/ory/graceful"
	openapi "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/cmd/cliclient"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
	"github.com/ory/x/pointerx"
	"github.com/ory/x/randx"
	"github.com/ory/x/tlsx"
	"github.com/ory/x/urlx"
)

var tokenUserLogin = template.Must(template.New("").Parse(`<html>
<body>
<h1>Login step</h1>
<form action="/login" method="post">
	<input type="hidden" name="ls" value="{{ .LoginChallenge }}">
	<input type="text" name="username" value="Username" required>
	<input type="checkbox" name="remember" checked>Remember login<br>
	<input type="checkbox" name="revoke-consents">Revoke previous consents<br>
	<button type="submit" name="action" value="accept">Submit</button>
	<button type="submit" name="action" value="deny">Cancel</button>
</form>
{{ if .Skip }}
	<b>user authenticated, could skip login UI.</b>
{{ else }}
	User unknown.
{{ end }}
<hr>
<h2>Complete login request</h2>
<pre>{{ .Raw }}</pre>
</body>
</html>`))

var tokenUserConsent = template.Must(template.New("").Parse(`<html>
<body>
<h1>Consent step</h1>
<form action="/consent" method="post">
	<input type="hidden" name="cs" value="{{ .ConsentChallenge }}">
	{{ if not .Audiences }}
		No token audiences requested.
	{{ else }}
		<h2>Requested audiences:</h2>
		<ul>
		{{ range .Audiences }}
			<li><input type="hidden" name="audience" value="{{ . }}">{{ . }}</li>
		{{ end }}
		</ul>
	{{ end }}
	{{ if not .Scopes }}
		No scopes requested.
	{{ else }}
		<h2>Requested scopes:</h2>
		{{ range .Scopes }}
		<input type="checkbox" name="scope" value="{{ . }}" checked>{{ . }}<br>
		{{ end }}
	{{ end }}
	<br>
	<input type="checkbox" name="remember" checked>Remember consent<br>
	<button type="submit" name="action" value="accept">Submit</button>
	<button type="submit" name="action" value="deny">Cancel</button>
</form>
{{ if .Skip }}
	<b>Consent established, could skip consent UI.</b>
{{ else }}
	No previous matching consent found, or client has requested re-consent.
{{ end }}
<hr>
<h2>Previous consents for this login session ({{ .SessionID }})</h2>
<pre>{{ .PreviousConsents }}</pre>
<hr>
<h2>Complete consent request</h2>
<pre>{{ .Raw }}</pre>
</body>
</html>`))

var tokenUserWelcome = template.Must(template.New("").Parse(`<html>
<body>
<h1>Welcome to the example OAuth 2.0 Consumer!</h1>
<p>This is an example app which emulates an OAuth 2.0 consumer application. Usually, this would be your web or mobile
    application and would use an <a href="https://oauth.net/code/">OAuth 2.0</a> or <a href="https://oauth.net/code/">OpenID
        Connect</a> library.</p>
<p>This example requests an OAuth 2.0 Access, Refresh, and OpenID Connect ID Token from the OAuth 2.0 Server (Ory
    Hydra).
    To initiate the flow, click the "Authorize Application" button.</p>
<p><a href="{{ .URL }}">Authorize application</a></p>
</body>
</html>`))

var tokenUserError = template.Must(template.New("").Parse(`<html>
<body>
<h1>An error occurred</h1>
<h2>{{ .Name }}</h2>
<p>{{ .Description }}</p>
<p>{{ .Hint }}</p>
<p>{{ .Debug }}</p>
</body>
</html>`))

var tokenUserResult = template.Must(template.New("").Parse(`<html>
<head></head>
<body>
<ul>
    <li>Access Token: <code>{{ .AccessToken }}</code></li>
    <li>Refresh Token: <code>{{ .RefreshToken }}</code></li>
    <li>Expires at: <code>{{ .Expiry }}</code></li>
    <li>ID Token: <code>{{ .IDToken }}</code></li>
</ul>
{{ if .DisplayBackButton }}
<a href="{{ .BackURL }}">Back to Welcome PageToken</a>
{{ end }}
</body>
</html>`))

func NewPerformAuthorizationCodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "authorization-code",
		Example: "{{ .CommandPath }} --client-id ... --client-secret ...",
		Short:   "Example OAuth 2.0 Client performing the OAuth 2.0 Authorize Code Flow",
		Long: `Starts an example web server that acts as an OAuth 2.0 Client performing the Authorize Code Flow.
This command will help you to see if Ory Hydra has been configured properly.

This command must not be used for anything else than manual testing or demo purposes. The server will terminate on error
and success, unless if the --no-shutdown flag is provided.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, endpoint, err := cliclient.NewClient(cmd)
			if err != nil {
				return err
			}

			endpoint = cliclient.GetOAuth2URLOverride(cmd, endpoint)

			isSSL := flagx.MustGetBool(cmd, "https")
			port := flagx.MustGetInt(cmd, "port")
			scopes := flagx.MustGetStringSlice(cmd, "scope")
			prompt := flagx.MustGetStringSlice(cmd, "prompt")
			maxAge := flagx.MustGetInt(cmd, "max-age")
			redirectUrl := flagx.MustGetString(cmd, "redirect")
			authUrl := flagx.MustGetString(cmd, "auth-url")
			tokenUrl := flagx.MustGetString(cmd, "token-url")
			audience := flagx.MustGetStringSlice(cmd, "audience")
			noShutdown := flagx.MustGetBool(cmd, "no-shutdown")
			skip := flagx.MustGetBool(cmd, "skip")
			responseMode := flagx.MustGetString(cmd, "response-mode")

			clientID := flagx.MustGetString(cmd, "client-id")
			if clientID == "" {
				_, _ = fmt.Fprint(cmd.OutOrStdout(), cmd.UsageString())
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Please provide a Client ID using --client-id flag, or OAUTH2_CLIENT_ID environment variable.")
				return cmdx.FailSilently(cmd)
			}

			clientSecret := flagx.MustGetString(cmd, "client-secret")
			proto := "http"
			if isSSL {
				proto = "https"
			}

			serverLocation := fmt.Sprintf("%s://127.0.0.1:%d/", proto, port)
			if redirectUrl == "" {
				redirectUrl = serverLocation + "callback"
			}

			if authUrl == "" {
				authUrl = urlx.AppendPaths(endpoint, "/oauth2/auth").String()
			}

			if tokenUrl == "" {
				tokenUrl = urlx.AppendPaths(endpoint, "/oauth2/token").String()
			}

			conf := oauth2.Config{
				ClientID:     clientID,
				ClientSecret: clientSecret,
				Endpoint: oauth2.Endpoint{
					AuthURL:  authUrl,
					TokenURL: tokenUrl,
				},
				RedirectURL: redirectUrl,
				Scopes:      scopes,
			}

			var generateAuthCodeURL = func() (string, string) {
				state := flagx.MustGetString(cmd, "state")
				if len(state) == 0 {
					generatedState, err := randx.RuneSequence(24, randx.AlphaLower)
					cmdx.Must(err, "Could not generate random state: %s", err)
					state = string(generatedState)
				}

				nonce, err := randx.RuneSequence(24, randx.AlphaLower)
				cmdx.Must(err, "Could not generate random state: %s", err)

				opts := []oauth2.AuthCodeOption{oauth2.SetAuthURLParam("nonce", string(nonce))}
				if len(audience) > 0 {
					opts = append(opts, oauth2.SetAuthURLParam("audience", strings.Join(audience, " ")))
				}
				if len(prompt) > 0 {
					opts = append(opts, oauth2.SetAuthURLParam("prompt", strings.Join(prompt, " ")))
				}
				if maxAge >= 0 {
					opts = append(opts, oauth2.SetAuthURLParam("max_age", strconv.Itoa(maxAge)))
				}
				if responseMode != "" {
					opts = append(opts, oauth2.SetAuthURLParam("response_mode", responseMode))
				}

				authCodeURL := conf.AuthCodeURL(state, opts...)
				return authCodeURL, state
			}
			authCodeURL, state := generateAuthCodeURL()

			r := http.NewServeMux()
			var tlsc *tls.Config
			if isSSL {
				key, err := rsa.GenerateKey(rand.Reader, 2048)
				if err != nil {
					_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Unable to generate RSA key pair: %s", err)
					return cmdx.FailSilently(cmd)
				}

				cert, err := tlsx.CreateSelfSignedTLSCertificate(key)
				cmdx.Must(err, "Unable to generate self-signed TLS Certificate: %s", err)
				// #nosec G402 - This is a false positive because we use graceful.WithDefaults which sets the correct TLS settings.
				tlsc = &tls.Config{Certificates: []tls.Certificate{*cert}}
			}

			server := graceful.WithDefaults(&http.Server{
				Addr:    fmt.Sprintf(":%d", port),
				Handler: r, TLSConfig: tlsc,
				ReadHeaderTimeout: time.Second * 5,
			})
			shutdown := func() {
				time.Sleep(time.Second * 1)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()
				_ = server.Shutdown(ctx)
			}

			r.Handle("GET /", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = tokenUserWelcome.Execute(w, &struct{ URL string }{URL: authCodeURL})
			}))

			r.Handle("GET /perform-flow", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, authCodeURL, http.StatusFound)
			}))

			rt := router{
				cl:    client,
				skip:  skip,
				cmd:   cmd,
				state: &state,
				conf:  &conf,
				onDone: func() {
					if !noShutdown {
						go shutdown()
					} else {
						// regenerate because we don't want to shutdown and we don't want to reuse nonce & state
						authCodeURL, state = generateAuthCodeURL()
					}
				},
				serverLocation: serverLocation,
				noShutdown:     noShutdown,
			}

			r.Handle("GET /login", http.HandlerFunc(rt.loginGET))
			r.Handle("POST /login", http.HandlerFunc(rt.loginPOST))
			r.Handle("GET /consent", http.HandlerFunc(rt.consentGET))
			r.Handle("POST /consent", http.HandlerFunc(rt.consentPOST))
			r.Handle("GET /callback", http.HandlerFunc(rt.callback))
			r.Handle("POST /callback", http.HandlerFunc(rt.callbackPOSTForm))

			if !flagx.MustGetBool(cmd, "no-open") {
				_ = webbrowser.Open(serverLocation) // ignore errors
			}

			_, _ = fmt.Fprintln(rt.cmd.ErrOrStderr(), "Setting up home route on "+serverLocation)
			_, _ = fmt.Fprintln(rt.cmd.ErrOrStderr(), "Setting up callback listener on "+serverLocation+"callback")
			_, _ = fmt.Fprintln(rt.cmd.ErrOrStderr(), "Press ctrl + c on Linux / Windows or cmd + c on OSX to end the process.")
			_, _ = fmt.Fprintf(rt.cmd.ErrOrStderr(), "If your browser does not open automatically, navigate to:\n\n\t%s\n\n", serverLocation)

			if isSSL {
				err = server.ListenAndServeTLS("", "")
			} else {
				err = server.ListenAndServe()
			}

			if errors.Is(err, http.ErrServerClosed) {
				return nil
			} else if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().Bool("no-open", false, "Do not open the browser window automatically")
	cmd.Flags().IntP("port", "p", 4446, "The port on which the server should run")
	cmd.Flags().StringSlice("scope", []string{"offline", "openid"}, "Request OAuth2 scope")
	cmd.Flags().StringSlice("prompt", []string{}, "Set the OpenID Connect prompt parameter")
	cmd.Flags().Int("max-age", -1, "Set the OpenID Connect max_age parameter. -1 means no max_age parameter will be used.")
	cmd.Flags().Bool("no-shutdown", false, "Do not terminate on success/error. State and nonce will be regenerated when auth flow has completed (either due to an error or success).")

	cmd.Flags().String("client-id", os.Getenv("OAUTH2_CLIENT_ID"), "Use the provided OAuth 2.0 Client ID, defaults to environment variable OAUTH2_CLIENT_ID")
	cmd.Flags().String("client-secret", os.Getenv("OAUTH2_CLIENT_SECRET"), "Use the provided OAuth 2.0 Client Secret, defaults to environment variable OAUTH2_CLIENT_SECRET")

	cmd.Flags().String("state", "", "Force a state value (insecure)")
	cmd.Flags().String("redirect", "", "Force a redirect url")
	cmd.Flags().StringSlice("audience", []string{}, "Request a specific OAuth 2.0 Access Token Audience")
	cmd.Flags().String("auth-url", "", "Usually it is enough to specify the `endpoint` flag, but if you want to force the authorization url, use this flag")
	cmd.Flags().String("token-url", "", "Usually it is enough to specify the `endpoint` flag, but if you want to force the token url, use this flag")
	cmd.Flags().Bool("https", false, "Sets up HTTPS for the endpoint using a self-signed certificate which is re-generated every time you start this command")
	cmd.Flags().Bool("skip", false, "Skip login and/or consent steps if possible. Only effective if you have configured the Login and Consent UI URLs to point to this server.")
	cmd.Flags().String("response-mode", "", "Set the response mode. Can be query (default) or form_post.")

	return cmd
}

type router struct {
	cl             *openapi.APIClient
	skip           bool
	cmd            *cobra.Command
	state          *string
	conf           *oauth2.Config
	onDone         func()
	serverLocation string
	noShutdown     bool
}

func (rt *router) loginGET(w http.ResponseWriter, r *http.Request) {
	req, raw, err := rt.cl.OAuth2API.GetOAuth2LoginRequest(r.Context()).
		LoginChallenge(r.URL.Query().Get("login_challenge")).
		Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer raw.Body.Close() //nolint:errcheck

	if rt.skip && req.GetSkip() {
		req, res, err := rt.cl.OAuth2API.AcceptOAuth2LoginRequest(r.Context()).
			LoginChallenge(req.Challenge).
			AcceptOAuth2LoginRequest(openapi.AcceptOAuth2LoginRequest{Subject: req.Subject}).
			Execute()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close() //nolint:errcheck
		http.Redirect(w, r, req.RedirectTo, http.StatusFound)
		return
	}

	pretty, err := prettyJSON(raw.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = tokenUserLogin.Execute(w, struct {
		LoginChallenge string
		Skip           bool
		SessionID      string
		Raw            string
	}{
		LoginChallenge: req.Challenge,
		Skip:           req.GetSkip(),
		SessionID:      req.GetSessionId(),
		Raw:            pretty,
	})
}

func (rt *router) loginPOST(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if r.FormValue("revoke-consents") == "on" {
		res, err := rt.cl.OAuth2API.RevokeOAuth2ConsentSessions(r.Context()).
			Subject(r.FormValue("username")).
			All(true).
			Execute()
		if err != nil {
			_, _ = fmt.Fprintln(rt.cmd.ErrOrStderr(), "Error revoking previous consents:", err)
		} else {
			_, _ = fmt.Fprintln(rt.cmd.ErrOrStderr(), "Revoked all previous consents")
		}
		defer res.Body.Close() //nolint:errcheck
	}
	switch r.FormValue("action") {
	case "accept":

		req, res, err := rt.cl.OAuth2API.AcceptOAuth2LoginRequest(r.Context()).
			LoginChallenge(r.FormValue("ls")).
			AcceptOAuth2LoginRequest(openapi.AcceptOAuth2LoginRequest{
				Subject:     r.FormValue("username"),
				Remember:    pointerx.Ptr(r.FormValue("remember") == "on"),
				RememberFor: pointerx.Int64(3600),
				Context: map[string]string{
					"context from": "login step",
				},
			}).Execute()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close() //nolint:errcheck
		http.Redirect(w, r, req.RedirectTo, http.StatusFound)

	case "deny":
		req, res, err := rt.cl.OAuth2API.RejectOAuth2LoginRequest(r.Context()).LoginChallenge(r.FormValue("ls")).Execute()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close() //nolint:errcheck
		http.Redirect(w, r, req.RedirectTo, http.StatusFound)

	default:
		http.Error(w, "Invalid action", http.StatusBadRequest)
	}
}

func (rt *router) consentGET(w http.ResponseWriter, r *http.Request) {
	req, raw, err := rt.cl.OAuth2API.GetOAuth2ConsentRequest(r.Context()).
		ConsentChallenge(r.URL.Query().Get("consent_challenge")).
		Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer raw.Body.Close() //nolint:errcheck

	if rt.skip && req.GetSkip() {
		req, res, err := rt.cl.OAuth2API.AcceptOAuth2ConsentRequest(r.Context()).
			ConsentChallenge(req.Challenge).
			AcceptOAuth2ConsentRequest(openapi.AcceptOAuth2ConsentRequest{
				GrantScope:               req.GetRequestedScope(),
				GrantAccessTokenAudience: req.GetRequestedAccessTokenAudience(),
				Remember:                 pointerx.Ptr(true),
				RememberFor:              pointerx.Int64(3600),
				Session: &openapi.AcceptOAuth2ConsentRequestSession{
					AccessToken: map[string]string{
						"foo": "bar",
					},
					IdToken: map[string]string{
						"baz": "bar",
					},
				},
			}).Execute()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close() //nolint:errcheck
		http.Redirect(w, r, req.RedirectTo, http.StatusFound)
		return
	}

	pretty, err := prettyJSON(raw.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, raw, err = rt.cl.OAuth2API.ListOAuth2ConsentSessions(r.Context()).
		Subject(req.GetSubject()).
		LoginSessionId(req.GetLoginSessionId()).
		Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer raw.Body.Close() //nolint:errcheck
	prettyPrevConsent, err := prettyJSON(raw.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = tokenUserConsent.Execute(w, struct {
		ConsentChallenge string
		Audiences        []string
		Scopes           []string
		Skip             bool
		SessionID        string
		PreviousConsents string
		Raw              string
	}{
		ConsentChallenge: req.Challenge,
		Audiences:        req.RequestedAccessTokenAudience,
		Scopes:           req.RequestedScope,
		Skip:             req.GetSkip(),
		SessionID:        req.GetLoginSessionId(),
		PreviousConsents: prettyPrevConsent,
		Raw:              pretty,
	})
}

func (rt *router) consentPOST(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	switch r.FormValue("action") {
	case "accept":
		req, res, err := rt.cl.OAuth2API.AcceptOAuth2ConsentRequest(r.Context()).
			ConsentChallenge(r.FormValue("cs")).
			AcceptOAuth2ConsentRequest(openapi.AcceptOAuth2ConsentRequest{
				GrantScope:               r.Form["scope"],
				GrantAccessTokenAudience: r.Form["audience"],
				Remember:                 pointerx.Ptr(r.FormValue("remember") == "on"),
				RememberFor:              pointerx.Int64(3600),
				Session: &openapi.AcceptOAuth2ConsentRequestSession{
					AccessToken: map[string]string{
						"foo": "bar",
					},
					IdToken: map[string]string{
						"baz": "bar",
					},
				},
			}).Execute()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close() //nolint:errcheck
		http.Redirect(w, r, req.RedirectTo, http.StatusFound)

	case "deny":
		req, res, err := rt.cl.OAuth2API.RejectOAuth2ConsentRequest(r.Context()).
			ConsentChallenge(r.FormValue("cs")).
			Execute()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close() //nolint:errcheck
		http.Redirect(w, r, req.RedirectTo, http.StatusFound)

	default:
		http.Error(w, "Invalid action", http.StatusBadRequest)
	}
}

func (rt *router) callback(w http.ResponseWriter, r *http.Request) {
	defer rt.onDone()

	if len(r.URL.Query().Get("error")) > 0 {
		_, _ = fmt.Fprintf(rt.cmd.ErrOrStderr(), "Got error: %s\n", r.URL.Query().Get("error_description"))

		w.WriteHeader(http.StatusInternalServerError)
		_ = tokenUserError.Execute(w, &ed{
			Name:        r.URL.Query().Get("error"),
			Description: r.URL.Query().Get("error_description"),
			Hint:        r.URL.Query().Get("error_hint"),
			Debug:       r.URL.Query().Get("error_debug"),
		})
		return
	}

	if r.URL.Query().Get("state") != *rt.state {
		descr := fmt.Sprintf("States do not match. Expected %q, got %q.", *rt.state, r.URL.Query().Get("state"))
		_, _ = fmt.Fprintln(rt.cmd.ErrOrStderr(), descr)

		w.WriteHeader(http.StatusInternalServerError)
		_ = tokenUserError.Execute(w, &ed{
			Name:        "States do not match",
			Description: descr,
		})
		return
	}

	code := r.URL.Query().Get("code")
	ctx := context.WithValue(rt.cmd.Context(), oauth2.HTTPClient, rt.cl.GetConfig().HTTPClient)
	token, err := rt.conf.Exchange(ctx, code)
	if err != nil {
		_, _ = fmt.Fprintf(rt.cmd.ErrOrStderr(), "Unable to exchange code for token: %s\n", err)

		w.WriteHeader(http.StatusInternalServerError)
		_ = tokenUserError.Execute(w, &ed{
			Name: err.Error(),
		})
		return
	}

	cmdx.PrintRow(rt.cmd, outputOAuth2Token(*token))
	_ = tokenUserResult.Execute(w, struct {
		AccessToken       string
		RefreshToken      string
		Expiry            string
		IDToken           string
		BackURL           string
		DisplayBackButton bool
	}{
		AccessToken:       token.AccessToken,
		RefreshToken:      token.RefreshToken,
		Expiry:            token.Expiry.Format(time.RFC1123),
		IDToken:           fmt.Sprintf("%s", token.Extra("id_token")),
		BackURL:           rt.serverLocation,
		DisplayBackButton: rt.noShutdown,
	})
}

func (rt *router) callbackPOSTForm(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	u := url.URL{
		Path:     r.URL.Path,
		RawQuery: r.PostForm.Encode(),
	}
	http.Redirect(w, r, u.String(), http.StatusFound)
}

type ed struct {
	Name        string
	Description string
	Hint        string
	Debug       string
}

func prettyJSON(r io.Reader) (string, error) {
	contentsRaw, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := json.Indent(&buf, contentsRaw, "", "  "); err != nil {
		return "", err
	}
	return buf.String(), nil
}
