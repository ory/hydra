/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

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

	"github.com/ory/graceful"

	"github.com/ory/hydra/cmd/cli"

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
<a href="{{ .BackURL }}">Back to Welcome Page</a>
{{ end }}
</body>
</html>`))

// cmd represents the token command
func NewTokenUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "An exemplary OAuth 2.0 Client performing the OAuth 2.0 Authorize Code Flow",
		Long: `Starts an exemplary web server that acts as an OAuth 2.0 Client performing the Authorize Code Flow.
This command will help you to see if Ory Hydra has been configured properly.

This command must not be used for anything else than manual testing or demo purposes. The server will terminate on error
and success, unless if the --no-shutdown flag is provided.`,
		Run: func(cmd *cobra.Command, args []string) {
			/* #nosec G402 - we want to support dev environments, hence tls trickery */
			ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: flagx.MustGetBool(cmd, "skip-tls-verify")},
			}})

			isSSL := flagx.MustGetBool(cmd, "https")
			port := flagx.MustGetInt(cmd, "port")
			scopes := flagx.MustGetStringSlice(cmd, "scope")
			prompt := flagx.MustGetStringSlice(cmd, "prompt")
			maxAge := flagx.MustGetInt(cmd, "max-age")
			redirectUrl := flagx.MustGetString(cmd, "redirect")
			backend := flagx.MustGetString(cmd, "token-url")
			frontend := flagx.MustGetString(cmd, "auth-url")
			audience := flagx.MustGetStringSlice(cmd, "audience")
			noShutdown := flagx.MustGetBool(cmd, "no-shutdown")

			clientID := flagx.MustGetString(cmd, "client-id")
			if clientID == "" {
				fmt.Print(cmd.UsageString())
				fmt.Println("Please provide a Client ID using --client-id flag, or OAUTH2_CLIENT_ID environment variable.")
				return
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

			if backend == "" {
				backend = urlx.AppendPaths(cli.RemoteURI(cmd), "/oauth2/token").String()
			}
			if frontend == "" {
				frontend = urlx.AppendPaths(cli.RemoteURI(cmd), "/oauth2/auth").String()
			}

			conf := oauth2.Config{
				ClientID:     clientID,
				ClientSecret: clientSecret,
				Endpoint: oauth2.Endpoint{
					TokenURL: backend,
					AuthURL:  frontend,
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

			fmt.Println("Setting up home route on " + serverLocation)
			fmt.Println("Setting up callback listener on " + serverLocation + "callback")
			fmt.Println("Press ctrl + c on Linux / Windows or cmd + c on OSX to end the process.")
			fmt.Printf("If your browser does not open automatically, navigate to:\n\n\t%s\n\n", serverLocation)

			r := httprouter.New()
			var tlsc *tls.Config
			if isSSL {
				key, err := rsa.GenerateKey(rand.Reader, 2048)
				cmdx.Must(err, "Unable to generate RSA key pair: %s", err)
				cert, err := tlsx.CreateSelfSignedTLSCertificate(key)
				cmdx.Must(err, "Unable to generate self-signed TLS Certificate: %s", err)
				// #nosec G402 - This is a false positive because we use graceful.WithDefaults which sets the correct TLS settings.
				tlsc = &tls.Config{Certificates: []tls.Certificate{*cert}}
			}

			server := graceful.WithDefaults(&http.Server{Addr: fmt.Sprintf(":%d", port), Handler: r, TLSConfig: tlsc})
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

			type ed struct {
				Name        string
				Description string
				Hint        string
				Debug       string
			}

			r.GET("/callback", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
				if len(r.URL.Query().Get("error")) > 0 {
					fmt.Printf("Got error: %s\n", r.URL.Query().Get("error_description"))

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
					fmt.Printf("States do not match. Expected %s, got %s\n", string(state), r.URL.Query().Get("state"))

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
					fmt.Printf("Unable to exchange code for token: %s\n", err)

					w.WriteHeader(http.StatusInternalServerError)
					_ = tokenUserError.Execute(w, &ed{
						Name: err.Error(),
					})
					onDone()
					return
				}

				idt := token.Extra("id_token")
				fmt.Printf("Access Token:\n\t%s\n", token.AccessToken)
				fmt.Printf("Refresh Token:\n\t%s\n", token.RefreshToken)
				fmt.Printf("Expires in:\n\t%s\n", token.Expiry.Format(time.RFC1123))
				fmt.Printf("ID Token:\n\t%v\n\n", idt)

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
					IDToken:           fmt.Sprintf("%v", idt),
					BackURL:           serverLocation,
					DisplayBackButton: noShutdown,
				})
				onDone()
			})
			var err error
			if isSSL {
				err = server.ListenAndServeTLS("", "")
			} else {
				err = server.ListenAndServe()
			}
			cmdx.Must(err, "%s", err)
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
	cmd.Flags().String("endpoint", os.Getenv("HYDRA_URL"), "Set the URL where Ory Hydra is hosted, defaults to environment variable HYDRA_URL")
	cmd.Flags().Bool("https", false, "Sets up HTTPS for the endpoint using a self-signed certificate which is re-generated every time you start this command")

	return cmd
}
