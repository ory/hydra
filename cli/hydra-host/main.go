package main

import (
	"os"

	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/codegangsta/cli"
	. "github.com/ory-am/hydra/cli/hydra-host/handler"
	//"github.com/ory-am/hydra/cli/hydra-host/templates"
	"fmt"
	"time"
)

func errorWrapper(f func(ctx *cli.Context) error) func(ctx *cli.Context) {
	return func(ctx *cli.Context) {
		if err := f(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred: %s\n", err)
		}
	}
}

var (
	ctx      = new(DefaultSystemContext)
	cl       = &Client{Ctx: ctx.GetSystemContext()}
	u        = &Account{Ctx: ctx.GetSystemContext()}
	co       = &Core{Ctx: ctx.GetSystemContext()}
	pl       = &Policy{Ctx: ctx.GetSystemContext()}
	Commands = []cli.Command{
		{
			Name:  "client",
			Usage: "Client actions",
			Subcommands: []cli.Command{
				{
					Name:   "create",
					Usage:  `Create a new client`,
					Action: errorWrapper(cl.Create),
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "i, id",
							Usage: "Set client's id",
						},
						cli.StringFlag{
							Name:  "s, secret",
							Usage: "The client's secret",
						},
						cli.StringFlag{
							Name:  "r, redirect-url",
							Usage: `A list of allowed redirect URLs: https://foobar.com/callback|https://bazbar.com/cb|http://localhost:3000/authcb`,
						},
						cli.BoolFlag{
							Name:  "as-superuser",
							Usage: "Grant superuser privileges to the client",
						},
					},
				},
			},
		},
		{
			Name:  "account",
			Usage: "Account actions",
			Subcommands: []cli.Command{
				{
					Name:      "create",
					Usage:     "Create a new account",
					ArgsUsage: "<username>",
					Action:    errorWrapper(u.Create),
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "password",
							Usage: "The account's password",
						},
						cli.BoolFlag{
							Name:  "as-superuser",
							Usage: "Grant superuser privileges to the account",
						},
					},
				},
			},
		},
		{
			Name:   "start",
			Usage:  "Start the host service",
			Action: errorWrapper(co.Start),
		},
		{
			Name:  "jwt",
			Usage: "JWT actions",
			Subcommands: []cli.Command{
				{
					Name:   "generate-keypair",
					Usage:  "Create a JWT PEM keypair.\n\n   You can use these files by providing the environment variables JWT_PRIVATE_KEY_PATH and JWT_PUBLIC_KEY_PATH",
					Action: errorWrapper(CreatePublicPrivatePEMFiles),
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "s, private-file-path",
							Value: "rs256-private.pem",
							Usage: "Where to save the private key PEM file",
						},
						cli.StringFlag{
							Name:  "p, public-file-path",
							Value: "rs256-public.pem",
							Usage: "Where to save the private key PEM file",
						},
					},
				},
			},
		},
		{
			Name:  "tls",
			Usage: "TLS actions",
			Subcommands: []cli.Command{
				{
					Name:   "generate-dummy-certificate",
					Usage:  "Create a dummy TLS certificate and private key.\n\n   You can use these files (in development!) by providing the environment variables TLS_CERT_PATH and TLS_KEY_PATH",
					Action: errorWrapper(CreateDummyTLSCert),
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "c, certificate-file-path",
							Value: "tls-cert.pem",
							Usage: "Where to save the private key PEM file",
						},
						cli.StringFlag{
							Name:  "k, key-file-path",
							Value: "tls-key.pem",
							Usage: "Where to save the private key PEM file",
						},
						cli.StringFlag{
							Name:  "u, host",
							Usage: "Comma-separated hostnames and IPs to generate a certificate for",
						},
						cli.StringFlag{
							Name:  "sd, start-date",
							Usage: "Creation date formatted as Jan 1 15:04:05 2011",
						},
						cli.DurationFlag{
							Name:  "d, duration",
							Value: 365 * 24 * time.Hour,
							Usage: "Duration that certificate is valid for",
						},
						cli.BoolFlag{
							Name:  "ca",
							Usage: "whether this cert should be its own Certificate Authority",
						},
						cli.IntFlag{
							Name:  "rb, rsa-bits",
							Value: 2048,
							Usage: "Size of RSA key to generate. Ignored if --ecdsa-curve is set",
						},
						cli.StringFlag{
							Name:  "ec, ecdsa-curve",
							Usage: "ECDSA curve to use to generate a key. Valid values are P224, P256, P384, P521",
						},
					},
				},
			},
		},
		{
			Name:  "policy",
			Usage: "Policy actions",
			Subcommands: []cli.Command{
				{
					Name:      "import",
					ArgsUsage: "<policies1.json> <policies2.json> <policies3.json>",
					Usage:     `Import a json file which defines an array of policies`,
					Action:    errorWrapper(pl.Import),
					Flags:     []cli.Flag{},
				},
			},
		},
	}
)

func main() {
	NewApp().Run(os.Args)
}

func NewApp() *cli.App {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "hydra-host"
	app.Usage = `Dragons guard your resources`
	app.Commands = Commands
	return app
}
