package handler

import (
	"encoding/json"
	"io/ioutil"

	"fmt"

	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/ory-am/ladon/policy"
	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/pborman/uuid"
)

type Policy struct {
	Ctx Context
}

func (c *Policy) Import(ctx *cli.Context) error {
	// Start the database backend
	if err := c.Ctx.Start(); err != nil {
		return fmt.Errorf("Could not start context: %s", err)
	}

	var policies []policy.DefaultPolicy
	for _, arg := range ctx.Args() {
		data, err := ioutil.ReadFile(arg)
		if err != nil {
			return fmt.Errorf("Could not read file: %s", err)
		}

		if err := json.Unmarshal(data, &policies); err != nil {
			return fmt.Errorf(`Could not unmarshall file %s to JSON: %s`, arg, err)
		}

		for _, pol := range policies {
			if pol.ID == "" {
				pol.ID = uuid.New()
			}

			if err := c.Ctx.GetPolicies().Create(&pol); err != nil {
				return fmt.Errorf(`Could not create policy: %s`, err)
			}
			fmt.Printf("Successfully created policy %s.\n", pol.ID)
		}
	}
	return nil
}
