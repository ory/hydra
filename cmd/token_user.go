package cmd

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/common/rand/sequence"
	"github.com/ory-am/hydra/pkg"
	"github.com/spf13/cobra"
	"github.com/toqueteos/webbrowser"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"gopkg.in/tylerb/graceful.v1"
)

// tokenUserCmd represents the token command
var tokenUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Generate an OAuth2 token using the code flow",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		if ok, _ := cmd.Flags().GetBool("skip-tls-verify"); ok {
			fmt.Println("Warning: Skipping TLS Certificate Verification.")
			ctx = context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}})
		}

		scopes, _ := cmd.Flags().GetStringSlice("scopes")
		conf := oauth2.Config{
			ClientID:     c.ClientID,
			ClientSecret: c.ClientSecret,
			Endpoint: oauth2.Endpoint{
				TokenURL: pkg.JoinURLStrings(c.ClusterURL, "/oauth2/token"),
				AuthURL:  pkg.JoinURLStrings(c.ClusterURL, "/oauth2/auth"),
			},
			RedirectURL: "http://localhost:4445/callback",
			Scopes:      scopes,
		}

		state, err := sequence.RuneSequence(24, []rune("abcdefghijklmnopqrstuvwxyz"))
		pkg.Must(err, "Could not generate random state: %s", err)

		nonce, err := sequence.RuneSequence(24, []rune("abcdefghijklmnopqrstuvwxyz"))
		pkg.Must(err, "Could not generate random state: %s", err)

		location := conf.AuthCodeURL(string(state)) + "&nonce=" + string(nonce)

		if ok, _ := cmd.Flags().GetBool("no-open"); !ok {
			webbrowser.Open(location)
		}

		fmt.Println("Setting up callback listener on http://localhost:4445/callback")
		fmt.Println("Press ctrl + c on Linux / Windows or cmd + c on OSX to end the process.")
		fmt.Printf("If your browser does not open automatically, navigate to:\n\n\t%s\n\n", location)

		srv := &graceful.Server{
			Timeout: 2 * time.Second,
			Server:  &http.Server{Addr: ":4445"},
		}
		r := httprouter.New()
		r.GET("/callback", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			defer srv.Stop(time.Second)

			if r.URL.Query().Get("error") != "" {
				message := fmt.Sprintf("Got error: %s", r.URL.Query().Get("error_description"))
				fmt.Println(message)

				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(message))
				return
			}

			if r.URL.Query().Get("state") != string(state) {
				message := fmt.Sprintf("States do not match. Expected %s, got %s", string(state), r.URL.Query().Get("state"))
				fmt.Println(message)

				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(message))
				return
			}

			code := r.URL.Query().Get("code")
			token, err := conf.Exchange(ctx, code)
			pkg.Must(err, "Could not exchange code for token: %s", err)

			fmt.Printf("Access Token:\n\t%s\n", token.AccessToken)
			fmt.Printf("Refresh Token:\n\t%s\n\n", token.RefreshToken)
			fmt.Printf("Expires in:\n\t%s\n\n", token.Expiry)

			w.Write([]byte(fmt.Sprintf(`
<html><head></head><body>
<ul>
	<li>Access Token: <code>%s</code></li>
	<li>Refresh Token: <code>%s</code></li>
	<li>Expires in: <code>%s</code></li>
`, token.AccessToken, token.RefreshToken, token.Expiry)))

			idt := token.Extra("id_token")
			if idt != nil {
				w.Write([]byte(fmt.Sprintf(`<li>ID Token: <code>%s</code></li>`, idt)))
				fmt.Printf("ID Token:\n\t%s\n\n", idt)
			}
			w.Write([]byte("</ul></body></html>"))
		})
		srv.Server.Handler = r
		srv.ListenAndServe()
	},
}

func init() {
	tokenCmd.AddCommand(tokenUserCmd)
	tokenUserCmd.Flags().Bool("no-open", false, "Do not open the browser window automatically")
	tokenUserCmd.Flags().StringSlice("scopes", []string{"hydra", "offline", "openid"}, "Ask for specific scopes")
}
