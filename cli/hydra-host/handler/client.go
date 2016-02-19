package handler

import (
	"fmt"

	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/RangelReale/osin"
	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/ory-am/common/rand/sequence"
	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/pborman/uuid"
)

type Client struct {
	Ctx Context
}

func (c *Client) Create(ctx *cli.Context) error {
	// Start the database backend
	if err := c.Ctx.Start(); err != nil {
		return fmt.Errorf("Could not start context: %s", err)
	}

	id := ctx.String("id")
	if id == "" {
		id = uuid.New()
	}

	secret := ctx.String("secret")
	if secret == "" {
		if seq, err := sequence.RuneSequence(10, sequence.AlphaNum); err != nil {
			return fmt.Errorf("Could not create rune sequence because %s", err)
		} else {
			secret = string(seq)
		}
	}

	client := &osin.DefaultClient{
		Id:          id,
		Secret:      secret,
		RedirectUri: ctx.String("redirect-url"),
		UserData:    "",
	}

	if err := c.Ctx.GetOsins().CreateClient(client); err != nil {
		return fmt.Errorf("Could not create client because %s", err)
	}
	fmt.Printf(`Created client "%s" with secret "%s" and redirect url "%s".`+"\n", client.Id, client.Secret, client.RedirectUri)

	if ctx.Bool("as-superuser") {
		if err := c.Ctx.GetPolicies().Create(superUserPolicy(client.Id)); err != nil {
			return fmt.Errorf("Could not create policy for client because %s", err)
		}
		fmt.Printf(`Granted superuser privileges to client "%s".`+"\n", client.Id)
	}
	return nil
}
