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
<p>This example requests an OAuth 2.0 Access, Refresh, and OpenID Connect ID Token from the OAuth 2.0 Server (ORY
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
<html>
<head></head>
<body>
<ul>
    <li>Access Token: <code>{{ .AccessToken }}</code></li>
    <li>Refresh Token: <code>{{ .RefreshToken }}</code></li>
    <li>Expires in: <code>{{ .Expiry }}</code></li>
    <li>ID Token: <code>{{ .IDToken }}</code></li>
</ul>
</body>
</html>`))

// tokenUserCmd represents the token command
var tokenUserCmd = &cobra.Command{
	Use:   "user",
	Short: "An exemplary OAuth 2.0 Client performing the OAuth 2.0 Authorize Code Flow",
	Long: `Starts an exemplary web server that acts as an OAuth 2.0 Client performing the Authorize Code Flow.
This command will help you to see if ORY Hydra has been configured properly.

This command must not be used for anything else than manual testing or demo purposes. The server will terminate on error
and success.`,
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

		clientID := flagx.MustGetString(cmd, "client-id")
		clientSecret := flagx.MustGetString(cmd, "client-secret")
		if clientID == "" || clientSecret == "" {
			fmt.Print(cmd.UsageString())
			fmt.Println("Please provide a Client ID and Client Secret using flags --client-id and --client-secret, or environment variables OAUTH2_CLIENT_ID and OAUTH2_CLIENT_SECRET.")
			return
		}

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
			tlsc = &tls.Config{Certificates: []tls.Certificate{*cert}}
		}

		server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: r, TLSConfig: tlsc}
		var shutdown = func() {
			time.Sleep(time.Second * 1)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			_ = server.Shutdown(ctx)
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

				go shutdown()
				return
			}

			if r.URL.Query().Get("state") != string(state) {
				fmt.Printf("States do not match. Expected %s, got %s\n", string(state), r.URL.Query().Get("state"))

				w.WriteHeader(http.StatusInternalServerError)
				_ = tokenUserError.Execute(w, &ed{
					Name:        "States do not match",
					Description: "Expected state " + string(state) + " but got " + r.URL.Query().Get("state"),
				})
				go shutdown()
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
				go shutdown()
				return
			}

			idt := token.Extra("id_token")
			fmt.Printf("Access Token:\n\t%s\n", token.AccessToken)
			fmt.Printf("Refresh Token:\n\t%s\n", token.RefreshToken)
			fmt.Printf("Expires in:\n\t%s\n", token.Expiry.Format(time.RFC1123))
			fmt.Printf("ID Token:\n\t%v\n\n", idt)

			_ = tokenUserResult.Execute(w, struct {
				AccessToken  string
				RefreshToken string
				Expiry       string
				IDToken      string
			}{
				AccessToken:  token.AccessToken,
				RefreshToken: token.RefreshToken,
				Expiry:       token.Expiry.Format(time.RFC1123),
				IDToken:      fmt.Sprintf("%v", idt),
			})

			go shutdown()
		})

		if isSSL {
			err = server.ListenAndServeTLS("", "")
		} else {
			err = server.ListenAndServe()
		}
		cmdx.Must(err, "%s", err)
	},
}

func init() {
	tokenCmd.AddCommand(tokenUserCmd)
	tokenUserCmd.Flags().Bool("no-open", false, "Do not open the browser window automatically")
	tokenUserCmd.Flags().IntP("port", "p", 4446, "The port on which the server should run")
	tokenUserCmd.Flags().StringSlice("scope", []string{"offline", "openid"}, "Request OAuth2 scope")
	tokenUserCmd.Flags().StringSlice("prompt", []string{}, "Set the OpenID Connect prompt parameter")
	tokenUserCmd.Flags().Int("max-age", 0, "Set the OpenID Connect max_age parameter")

	tokenUserCmd.Flags().String("client-id", os.Getenv("OAUTH2_CLIENT_ID"), "Use the provided OAuth 2.0 Client ID, defaults to environment variable OAUTH2_CLIENT_ID")
	tokenUserCmd.Flags().String("client-secret", os.Getenv("OAUTH2_CLIENT_SECRET"), "Use the provided OAuth 2.0 Client Secret, defaults to environment variable OAUTH2_CLIENT_SECRET")

	tokenUserCmd.Flags().String("redirect", "", "Force a redirect url")
	tokenUserCmd.Flags().StringSlice("audience", []string{}, "Request a specific OAuth 2.0 Access Token Audience")
	tokenUserCmd.Flags().String("auth-url", "", "Usually it is enough to specify the `endpoint` flag, but if you want to force the authorization url, use this flag")
	tokenUserCmd.Flags().String("token-url", "", "Usually it is enough to specify the `endpoint` flag, but if you want to force the token url, use this flag")
	tokenUserCmd.Flags().String("endpoint", os.Getenv("HYDRA_URL"), "Set the URL where ORY Hydra is hosted, defaults to environment variable HYDRA_URL")
	tokenUserCmd.Flags().Bool("https", false, "Sets up HTTPS for the endpoint using a self-signed certificate which is re-generated every time you start this command")
}
