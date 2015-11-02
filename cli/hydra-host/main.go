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
			Usage: "client actions",
			Subcommands: []cli.Command{
				{
					Name:   "create",
					Usage:  "create a new client",
					Action: cl.Create,
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "as-superuser",
							Usage: "grant superuser privileges to the client",
						},
					},
				},
			},
		},
		{
			Name:  "user",
			Usage: "user actions",
			Subcommands: []cli.Command{
				{
					Name:      "create",
					Usage:     "create a new user",
					ArgsUsage: "<email>",
					Action:    u.Create,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "password",
							Usage: "the user's password",
						},
						cli.BoolFlag{
							Name:  "as-superuser",
							Usage: "grant superuser privileges to the user",
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
