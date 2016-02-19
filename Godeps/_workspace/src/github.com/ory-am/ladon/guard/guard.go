package guard

import (
	log "github.com/ory-am/hydra/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/ory-am/common/compiler"
	. "github.com/ory-am/hydra/Godeps/_workspace/src/github.com/ory-am/ladon/guard/operator"
	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/ory-am/ladon/policy"
	"regexp"
	"strings"
)

type Guard struct {
	operators      map[string]Operator
	disableLogging bool
}

func (g *Guard) SetOperators(ops map[string]Operator) {
	g.operators = make(map[string]Operator)
	for key, op := range ops {
		g.AddOperator(key, op)
	}
}

func (g *Guard) AddOperator(name string, op Operator) {
	name = strings.ToLower(name)
	g.operators[name] = op
}

func (g *Guard) SetDefaultOperators() {
	g.SetOperators(map[string]Operator{
		"SubjectIsOwner":    SubjectIsOwner,
		"SubjectIsNotOwner": SubjectIsNotOwner,
	})
}

func (g *Guard) GetOperator(name string) (Operator, bool) {
	if len(g.operators) == 0 {
		g.SetDefaultOperators()
	}
	op, ok := g.operators[strings.ToLower(name)]
	return op, ok
}

func (g *Guard) IsGranted(resource, permission, subject string, policies []policy.Policy, ctx *Context) (allowed bool, err error) {
	// Iterate through all policies
	for _, p := range policies {
		// Does the resource match with one of the policies?
		if rm, err := Matches(p, p.GetResources(), resource); err != nil {
			if !g.disableLogging {
				log.WithFields(log.Fields{
					"resources": p.GetResources(),
					"resource":  resource,
					"error":     err,
				}).Warn("Could not match resources.")
			}
			return false, err
		} else if !rm {
			continue
		}

		// Does the action match with one of the policies?
		if pm, err := Matches(p, p.GetPermissions(), permission); err != nil {
			if !g.disableLogging {
				log.WithFields(log.Fields{
					"permissions": p.GetPermissions(),
					"permission":  permission,
					"error":       err,
				}).Warn("Could not match permissions.")
			}
			return false, err
		} else if !pm {
			continue
		}

		// Does the subject match with one of the policies?
		if sm, err := Matches(p, p.GetSubjects(), subject); err != nil {
			if !g.disableLogging {
				log.WithFields(log.Fields{
					"subjects": p.GetSubjects(),
					"subject":  subject,
					"error":    err,
				}).Warn("Could not match subjects.")
			}
			return false, err
		} else if !sm && len(p.GetSubjects()) > 0 {
			// If no match exists, but the subjects are scoped, this policy is irrelevant
			continue
		}

		if !g.PassesConditions(p, ctx, permission, resource, subject) {
			continue
		}

		// Does the policy enforce a deny policy? If yes, this overrides all allow policies -> access denied.
		if !p.HasAccess() {
			return false, nil
		}
		allowed = true
	}
	return allowed, nil
}

func (g *Guard) PassesConditions(p policy.Policy, ctx *Context, permission, resource, subject string) (passes bool) {
	var extra map[string]interface{}
	passes = len(p.GetConditions()) == 0
	for _, condition := range p.GetConditions() {
		op, ok := g.GetOperator(condition.GetOperator())
		if !ok {
			if !g.disableLogging {
				log.WithFields(log.Fields{
					"subjects": p.GetSubjects(),
					"subject":  subject,
				}).Warn("Could not check conditions.")
			}
			return false
		}

		extra = condition.GetExtra()
		extra["permission"] = permission
		extra["resource"] = resource
		extra["subject"] = subject
		if !op(extra, ctx) {
			return false
		}

		passes = true
	}
	return passes
}

func Matches(p policy.Policy, patterns []string, match string) (bool, error) {
	var reg *regexp.Regexp
	var err error
	var matches bool
	for _, h := range patterns {
		reg, err = compiler.CompileRegex(h, p.GetStartDelimiter(), p.GetEndDelimiter())
		if err != nil {
			return false, err
		}

		matches = reg.MatchString(match)
		if matches {
			return true, nil
		}
	}
	return false, nil
}
