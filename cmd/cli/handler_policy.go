package cli

import (
	"fmt"
	"os"

	"github.com/ory/hydra/config"
	"github.com/ory/hydra/pkg"
	"github.com/ory/hydra/policy"
	"github.com/ory/ladon"
	"github.com/spf13/cobra"
	"github.com/square/go-jose/json"
)

type PolicyHandler struct {
	Config *config.Config
}

func (h *PolicyHandler) newJwkManager(cmd *cobra.Command) *policy.HTTPManager {
	dry, _ := cmd.Flags().GetBool("dry")
	term, _ := cmd.Flags().GetBool("fake-tls-termination")

	return &policy.HTTPManager{
		Dry:                dry,
		Endpoint:           h.Config.Resolve("/policies"),
		Client:             h.Config.OAuth2Client(cmd),
		FakeTLSTermination: term,
	}
}

func newPolicyHandler(c *config.Config) *PolicyHandler {
	return &PolicyHandler{
		Config: c,
	}
}

func (h *PolicyHandler) CreatePolicy(cmd *cobra.Command, args []string) {
	m := h.newJwkManager(cmd)
	files, _ := cmd.Flags().GetStringSlice("files")
	if len(files) > 0 {
		for _, path := range files {
			reader, err := os.Open(path)
			pkg.Must(err, "Could not open file %s: %s", path, err)
			var p ladon.DefaultPolicy
			err = json.NewDecoder(reader).Decode(&p)
			pkg.Must(err, "Could not parse JSON: %s", err)
			err = m.Create(&p)
			pkg.Must(err, "Could not create policy: %s", err)
			fmt.Printf("Imported policy %s from %s.\n", p.ID, path)
		}
		return
	}

	id, _ := cmd.Flags().GetString("id")
	description, _ := cmd.Flags().GetString("description")
	subjects, _ := cmd.Flags().GetStringSlice("subjects")
	resources, _ := cmd.Flags().GetStringSlice("resources")
	actions, _ := cmd.Flags().GetStringSlice("actions")
	isAllow, _ := cmd.Flags().GetBool("allow")
	if len(subjects) == 0 || len(resources) == 0 || len(actions) == 0 {
		fmt.Println(cmd.UsageString())
		fmt.Println("")
		fmt.Println("Got empty subject, resource or action list")
		return
	}

	effect := ladon.DenyAccess
	if isAllow {
		effect = ladon.AllowAccess
	}

	p := &ladon.DefaultPolicy{
		ID:          id,
		Description: description,
		Subjects:    subjects,
		Resources:   resources,
		Actions:     actions,
		Effect:      effect,
	}
	err := m.Create(p)
	if m.Dry {
		fmt.Printf("%s\n", err)
		return
	}
	pkg.Must(err, "Could not create policy: %s", err)
	fmt.Printf("Created policy %s.\n", p.ID)

}

func (h *PolicyHandler) AddResourceToPolicy(cmd *cobra.Command, args []string) {
	m := h.newJwkManager(cmd)
	if len(args) < 2 {
		fmt.Print(cmd.UsageString())
		return
	}

	pp, err := m.Get(args[0])
	pkg.Must(err, "Could not get policy: %s", err)

	p := pp.(*ladon.DefaultPolicy)
	p.Resources = append(p.Resources, args[1:]...)

	err = m.Update(p)
	if m.Dry {
		fmt.Printf("%s\n", err)
		return
	}
	pkg.Must(err, "Could not update policy: %s", err)
	fmt.Printf("Added resources to policy %s", p.ID)
}

func (h *PolicyHandler) RemoveResourceFromPolicy(cmd *cobra.Command, args []string) {
	m := h.newJwkManager(cmd)
	if len(args) < 2 {
		fmt.Print(cmd.UsageString())
		return
	}

	pp, err := m.Get(args[0])
	pkg.Must(err, "Could not get policy: %s", err)

	p := pp.(*ladon.DefaultPolicy)
	resources := []string{}
	for _, r := range p.Resources {
		var filter bool
		for _, a := range args[1:] {
			if r == a {
				filter = true
			}
		}
		if !filter {
			resources = append(resources, r)
		}
	}
	p.Resources = resources

	err = m.Update(p)
	if m.Dry {
		fmt.Printf("%s\n", err)
		return
	}
	pkg.Must(err, "Could not update policy: %s", err)
	fmt.Printf("Removed resources from policy %s", p.ID)
}

func (h *PolicyHandler) AddSubjectToPolicy(cmd *cobra.Command, args []string) {
	m := h.newJwkManager(cmd)
	if len(args) < 2 {
		fmt.Print(cmd.UsageString())
		return
	}

	pp, err := m.Get(args[0])
	if m.Dry {
		fmt.Printf("%s\n", err)
	} else {
		pkg.Must(err, "Could not get policy: %s", err)
	}

	p := pp.(*ladon.DefaultPolicy)
	p.Subjects = append(p.Subjects, args[1:]...)

	err = m.Update(p)
	if m.Dry {
		fmt.Printf("%s\n", err)
		return
	}
	pkg.Must(err, "Could not update policy: %s", err)
	fmt.Printf("Added subjects to policy %s", p.ID)
}

func (h *PolicyHandler) RemoveSubjectFromPolicy(cmd *cobra.Command, args []string) {
	m := h.newJwkManager(cmd)
	if len(args) < 2 {
		fmt.Print(cmd.UsageString())
		return
	}

	pp, err := m.Get(args[0])
	pkg.Must(err, "Could not get policy: %s", err)

	p := pp.(*ladon.DefaultPolicy)
	subjects := []string{}
	for _, r := range p.Subjects {
		var filter bool
		for _, a := range args[1:] {
			if r == a {
				filter = true
			}
		}
		if !filter {
			subjects = append(subjects, r)
		}
	}
	p.Subjects = subjects

	err = m.Update(p)
	if m.Dry {
		fmt.Printf("%s\n", err)
		return
	}
	pkg.Must(err, "Could not update policy: %s", err)
	fmt.Printf("Removed subjects from policy %s", p.ID)
}

func (h *PolicyHandler) AddActionToPolicy(cmd *cobra.Command, args []string) {
	m := h.newJwkManager(cmd)
	if len(args) < 2 {
		fmt.Print(cmd.UsageString())
		return
	}

	pp, err := m.Get(args[0])
	pkg.Must(err, "Could not get policy: %s", err)

	p := pp.(*ladon.DefaultPolicy)
	p.Actions = append(p.Actions, args[1:]...)

	err = m.Update(p)
	if m.Dry {
		fmt.Printf("%s\n", err)
		return
	}
	pkg.Must(err, "Could not update policy: %s", err)
	fmt.Printf("Added actions to policy %s", p.ID)
}

func (h *PolicyHandler) RemoveActionFromPolicy(cmd *cobra.Command, args []string) {
	m := h.newJwkManager(cmd)
	if len(args) < 2 {
		fmt.Print(cmd.UsageString())
		return
	}

	pp, err := m.Get(args[0])
	pkg.Must(err, "Could not get policy: %s", err)

	p := pp.(*ladon.DefaultPolicy)
	actions := []string{}
	for _, r := range p.Actions {
		var filter bool
		for _, a := range args[1:] {
			if r == a {
				filter = true
			}
		}
		if !filter {
			actions = append(actions, r)
		}
	}
	p.Actions = actions

	err = m.Update(p)
	if m.Dry {
		fmt.Printf("%s\n", err)
		return
	}
	pkg.Must(err, "Could not update policy: %s", err)
	fmt.Printf("Removed actions from policy %s", p.ID)
}

func (h *PolicyHandler) GetPolicy(cmd *cobra.Command, args []string) {
	m := h.newJwkManager(cmd)
	if len(args) == 0 {
		fmt.Print(cmd.UsageString())
		return
	}

	p, err := m.Get(args[0])
	if m.Dry {
		fmt.Printf("%s\n", err)
		return
	}
	pkg.Must(err, "Could not retrieve policy: %s", err)

	out, err := json.MarshalIndent(p, "", "\t")
	pkg.Must(err, "Could not convert policy to JSON: %s", err)

	fmt.Printf("%s\n", out)
}

func (h *PolicyHandler) DeletePolicy(cmd *cobra.Command, args []string) {
	m := h.newJwkManager(cmd)
	if len(args) == 0 {
		fmt.Print(cmd.UsageString())
		return
	}

	for _, arg := range args {
		err := m.Delete(arg)
		if m.Dry {
			fmt.Printf("%s\n", err)
			continue
		}
		pkg.Must(err, "Could not delete policy: %s", err)
		fmt.Printf("Connection %s deleted.\n", arg)
	}
}
