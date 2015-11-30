package handler

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/ory-am/ladon/policy"
	"github.com/pborman/uuid"
)

type Policy struct {
	Ctx *Context
}

type DefaultPolicy struct {
	Description string                    `json:"description"`
	Subjects    []string                  `json:"subjects"`
	Effect      string                    `json:"effect"`
	Resources   []string                  `json:"resources"`
	Permissions []string                  `json:"permissions"`
	Conditions  []policy.DefaultCondition `json:"conditions"`
}

func (c *Policy) Import(ctx *cli.Context) {
	var policies []DefaultPolicy
	c.Ctx.Start()
	for _, arg := range ctx.Args() {
		data, err := ioutil.ReadFile(arg)
		if err != nil {
			log.Fatalf("Could not read file: %s", err)
		}

		if err := json.Unmarshal(data, &policies); err != nil {
			log.Fatalf(`Could not unmarshall file %s to JSON: %s`, arg, err)
		}

		for _, pol := range policies {
			conditions := make([]policy.Condition, len(pol.Conditions))
			for k, c := range pol.Conditions {
				conditions[k] = &c
			}

			pp := &policy.DefaultPolicy{
				ID:          uuid.New(),
				Description: pol.Description,
				Subjects:    pol.Subjects,
				Effect:      pol.Effect,
				Resources:   pol.Resources,
				Permissions: pol.Permissions,
				Conditions:  conditions,
			}
			if err := c.Ctx.Policies.Create(pp); err != nil {
				log.Fatalf(`Could not create policy: %s`, err)
			}

			log.Printf("Successfully created policy %s.", pp.ID)
		}
	}
}
