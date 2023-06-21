// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ory/hydra/v2/cmd/cliclient"

	"github.com/pkg/errors"

	"github.com/ory/graceful"

	"github.com/julienschmidt/httprouter"
	"github.com/spf13/cobra"
	"github.com/toqueteos/webbrowser"
	"golang.org/x/oauth2"

	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
	"github.com/ory/x/randx"
	"github.com/ory/x/tlsx"
	"github.com/ory/x/urlx"
)

var tokenUserWelcome = template.Must(template.New("").Parse(`<html>
<body>
<h1>Welcome to the exemplary OAuth 2.0 Consumer!</h1>
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
    <li>Expires in: <code>{{ .Expiry }}</code></li>
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
		Short:   "An exemplary OAuth 2.0 Client performing the OAuth 2.0 Authorize Code Flow",
		Long: `Starts an exemplary web server that acts as an OAuth 2.0 Client performing the Authorize Code Flow.
This command will help you to see if Ory Hydra has been configured properly.

This command must not be used for anything else than manual testing or demo purposes. The server will terminate on error
and success, unless if the --no-shutdown flag is provided.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, endpoint, err := cliclient.NewClient(cmd)
			if err != nil {
				return err
			}

			endpoint = cliclient.GetOAuth2URLOverride(cmd, endpoint)

			ctx := context.WithValue(cmd.Context(), oauth2.HTTPClient, client)
			isSSL := flagx.MustGetBool(cmd, "https")
			port := flagx.MustGetInt(cmd, "port")
			scopes := flagx.MustGetStringSlice(cmd, "scope")
			prompt := flagx.MustGetStringSlice(cmd, "prompt")
			maxAge := flagx.MustGetInt(cmd, "max-age")
			redirectUrl := flagx.MustGetString(cmd, "redirect")
			audience := flagx.MustGetStringSlice(cmd, "audience")
			noShutdown := flagx.MustGetBool(cmd, "no-shutdown")

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

			if err != nil {
				return err
			}
			conf := oauth2.Config{
				ClientID:     clientID,
				ClientSecret: clientSecret,
				Endpoint: oauth2.Endpoint{
					TokenURL: urlx.AppendPaths(endpoint, "/oauth2/token").String(),
					AuthURL:  urlx.AppendPaths(endpoint, "/oauth2/auth").String(),
				},
				RedirectURL: redirectUrl,
				Scopes:      scopes,
			}

			var generateAuthCodeURL = func() (string, []rune) {
				state, err := randx.RuneSequence(24, randx.AlphaLower)
				cmdx.Must(err, "Could not generate random state: %s", err)

				nonce, err := randx.RuneSequence(24, randx.AlphaLower)
				cmdx.Must(err, "Could not generate random state: %s", err)

				authCodeURL := conf.AuthCodeURL(
					string(state),
					oauth2.SetAuthURLParam("audience", strings.Join(audience, "+")),
					oauth2.SetAuthURLParam("nonce", string(nonce)),
					oauth2.SetAuthURLParam("prompt", strings.Join(prompt, "+")),
					oauth2.SetAuthURLParam("max_age", strconv.Itoa(maxAge)),
				)
				return authCodeURL, state
			}
			authCodeURL, state := generateAuthCodeURL()

			if !flagx.MustGetBool(cmd, "no-open") {
				_ = webbrowser.Open(serverLocation) // ignore errors
			}

			_, _ = fmt.Fprintln(os.Stderr, "Setting up home route on "+serverLocation)
			_, _ = fmt.Fprintln(os.Stderr, "Setting up callback listener on "+serverLocation+"callback")
			_, _ = fmt.Fprintln(os.Stderr, "Press ctrl + c on Linux / Windows or cmd + c on OSX to end the process.")
			_, _ = fmt.Fprintf(os.Stderr, "If your browser does not open automatically, navigate to:\n\n\t%s\n\n", serverLocation)

			r := httprouter.New()
			var tlsc *tls.Config
			if isSSL {
				key, err := rsa.GenerateKey(rand.Reader, 2048)
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "Unable to generate RSA key pair: %s", err)
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
			var shutdown = func() {
				time.Sleep(time.Second * 1)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()
				_ = server.Shutdown(ctx)
			}
			var onDone = func() {
				if !noShutdown {
					go shutdown()
				} else {
					// regenerate because we don't want to shutdown and we don't want to reuse nonce & state
					authCodeURL, state = generateAuthCodeURL()
				}
			}

			r.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
				_ = tokenUserWelcome.Execute(w, &struct{ URL string }{URL: authCodeURL})
			})

			r.GET("/perform-flow", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
				http.Redirect(w, r, authCodeURL, http.StatusFound)
			})

			type ed struct {
				Name        string
				Description string
				Hint        string
				Debug       string
			}

			r.GET("/callback", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
				if len(r.URL.Query().Get("error")) > 0 {
					_, _ = fmt.Fprintf(os.Stderr, "Got error: %s\n", r.URL.Query().Get("error_description"))

					w.WriteHeader(http.StatusInternalServerError)
					_ = tokenUserError.Execute(w, &ed{
						Name:        r.URL.Query().Get("error"),
						Description: r.URL.Query().Get("error_description"),
						Hint:        r.URL.Query().Get("error_hint"),
						Debug:       r.URL.Query().Get("error_debug"),
					})

					onDone()
					return
				}

				if r.URL.Query().Get("state") != string(state) {
					_, _ = fmt.Fprintf(os.Stderr, "States do not match. Expected %s, got %s\n", string(state), r.URL.Query().Get("state"))

					w.WriteHeader(http.StatusInternalServerError)
					_ = tokenUserError.Execute(w, &ed{
						Name:        "States do not match",
						Description: "Expected state " + string(state) + " but got " + r.URL.Query().Get("state"),
					})
					onDone()
					return
				}

				code := r.URL.Query().Get("code")
				token, err := conf.Exchange(ctx, code)
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "Unable to exchange code for token: %s\n", err)

					w.WriteHeader(http.StatusInternalServerError)
					_ = tokenUserError.Execute(w, &ed{
						Name: err.Error(),
					})
					onDone()
					return
				}

				cmdx.PrintRow(cmd, outputOAuth2Token(*token))
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
					BackURL:           serverLocation,
					DisplayBackButton: noShutdown,
				})
				onDone()
			})

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
	cmd.Flags().Int("max-age", 0, "Set the OpenID Connect max_age parameter")
	cmd.Flags().Bool("no-shutdown", false, "Do not terminate on success/error. State and nonce will be regenerated when auth flow has completed (either due to an error or success).")

	cmd.Flags().String("client-id", os.Getenv("OAUTH2_CLIENT_ID"), "Use the provided OAuth 2.0 Client ID, defaults to environment variable OAUTH2_CLIENT_ID")
	cmd.Flags().String("client-secret", os.Getenv("OAUTH2_CLIENT_SECRET"), "Use the provided OAuth 2.0 Client Secret, defaults to environment variable OAUTH2_CLIENT_SECRET")

	cmd.Flags().String("redirect", "", "Force a redirect url")
	cmd.Flags().StringSlice("audience", []string{}, "Request a specific OAuth 2.0 Access Token Audience")
	cmd.Flags().String("auth-url", "", "Usually it is enough to specify the `endpoint` flag, but if you want to force the authorization url, use this flag")
	cmd.Flags().String("token-url", "", "Usually it is enough to specify the `endpoint` flag, but if you want to force the token url, use this flag")
	cmd.Flags().Bool("https", false, "Sets up HTTPS for the endpoint using a self-signed certificate which is re-generated every time you start this command")

	return cmd
}
