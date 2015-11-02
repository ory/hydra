package handler

import (
	"fmt"
	"github.com/RangelReale/osin"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/ory-am/common/rand/sequence"
	"github.com/pborman/uuid"
)

type Client struct {
	Ctx *Context
}

func (c *Client) Create(ctx *cli.Context) {
	seq, err := sequence.RuneSequence(10, []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"))
	if err != nil {
		log.Fatalf("err")
	}

	client := &osin.DefaultClient{
		Id:          uuid.New(),
		Secret:      string(seq),
		RedirectUri: "",
		UserData:    "",
	}

	if err := c.Ctx.Osins.CreateClient(client); err != nil {
		log.Fatalf("%s", err)
	}

	fmt.Printf(`Created client "%s" with secret "%s".`+"\n", client.Id, client.Secret)

	if ctx.Bool("as-superuser") {
		if err := c.Ctx.Policies.Create(superUserPolicy(client.Id)); err != nil {
			log.Fatalf("%s", err)
		}
		fmt.Printf(`Granted superuser privileges to client "%s".`+"\n", client.Id)
	}
}
