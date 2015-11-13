package handler

/*
import (
	"fmt"
	"github.com/RangelReale/osin"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/ory-am/common/rand/sequence"
)

type Permission struct {
	Ctx *Context
}

func (c *Permission) Grant(ctx *cli.Context) {
	subject := ctx.String("subject")
	if subject == "" {
		log.Fatalf("Please provide a subject.")
	}

	onlyIfOwner := ctx.Bool("only-if-owner")
	template := ctx.String("template")

	secret := ctx.String("secret")
	if secret == "" {
		if seq, err := sequence.RuneSequence(10, sequence.AlphaNum); err != nil {
			log.Fatalf("err")
		} else {
			secret = string(seq)
		}
	}

	client := &osin.DefaultClient{
		Id:          subject,
		Secret:      secret,
		RedirectUri: ctx.String("redirect-url"),
		UserData:    "",
	}

	c.Ctx.Start()
	if err := c.Ctx.Osins.CreateClient(client); err != nil {
		log.Fatalf("%s", err)
	}

	fmt.Printf(`Created client "%s" with secret "%s" and redirect url "%s".` + "\n", client.Id, client.Secret, client.RedirectUri)

	if ctx.Bool("as-superuser") {
		if err := c.Ctx.Policies.Create(superUserPolicy(client.Id)); err != nil {
			log.Fatalf("%s", err)
		}
		fmt.Printf(`Granted superuser privileges to client "%s".` + "\n", client.Id)
	}
}*/
