package main

import (
	"github.com/codegangsta/cli"
	. "github.com/ory-am/hydra/cli/hydra-host/handler"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "hydra-host"
	app.Usage = `Dragons guard your resources.`

	ctx := new(Context)
	cl := &Client{Ctx: ctx}
	u := &User{Ctx: ctx}
	co := &Core{Ctx: ctx}
	app.Commands = []cli.Command{
		{
			Name:  "client",
			Usage: "Client actions",
			Subcommands: []cli.Command{
				{
					Name:   "create",
					Usage:  `Create a new client.`,
					Action: cl.Create,
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
			Name:  "user",
			Usage: "User actions",
			Subcommands: []cli.Command{
				{
					Name:      "create",
					Usage:     "Create a new user",
					ArgsUsage: "<email>",
					Action:    u.Create,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "password",
							Usage: "The user's password",
						},
						cli.BoolFlag{
							Name:  "as-superuser",
							Usage: "Grant superuser privileges to the user",
						},
					},
				},
			},
		},
		{
			Name:   "start",
			Usage:  "start hydra-host",
			Action: co.Start,
		},
	}
	app.Run(os.Args)
}
