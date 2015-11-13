package main

import (
	"github.com/codegangsta/cli"
	. "github.com/ory-am/hydra/cli/hydra-host/handler"
	"os"
	//"github.com/ory-am/hydra/cli/hydra-host/templates"
	//"fmt"
)

func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "hydra-host"
	app.Usage = `Dragons guard your resources.`

	ctx := new(Context)
	cl := &Client{Ctx: ctx}
	u := &User{Ctx: ctx}
	co := &Core{Ctx: ctx}
	app.Commands = []cli.Command{
		{
			Name:  "client",
			Usage: "Client actions.",
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
			Usage: "User actions.",
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
			Usage:  "Start the host service.",
			Action: co.Start,
		},
		/*{
			Name:  "policy",
			Usage: "Policy actions.",
			Subcommands: []cli.Command{
				{
					Name:   "grant",
					ArgsUsage: "<template>",
					Usage:  `Grant grants various policy templates to subjects.`,
					Action: cl.Create,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "s, subject",
							Usage: "Set the subject's id.",
						},
						cli.BoolFlag{
							Name:  "only-if-owner",
							Usage: "Only allow access if it the subject is also the owner of the resource.",
						},
					},
					BashComplete: func(c *cli.Context) {
						if len(c.Args()) > 0 {
							return
						}
						for _, t := range templates.Templates {
							fmt.Println(t)
						}
					},
				},
			},
		},*/
	}
	app.Run(os.Args)
}
