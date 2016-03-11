// Package rethinkdb is a ladon storage backend for rethinkDB.
package rethinkdb

import (
	"encoding/json"
	"fmt"
	"time"

	rdb "github.com/dancannon/gorethink"

	"github.com/go-errors/errors"
	"github.com/ory-am/common/compiler"
	"github.com/ory-am/common/pkg"
	"github.com/ory-am/ladon/policy"
)

const policyTableName = "ladon_policy"

type rethinkPolicy struct {
	ID                string                 `gorethink:"id"`
	Description       string                 `gorethink:"description"`
	Effect            string                 `gorethink:"effect"`
	CreatedAt         int64                  `gorethink:"created_at"`
	Conditions        []byte                 `gorethink:"conditions"`
	PolicySubjects    []linkedPolicyResource `gorethink:"ladon_policy_subjects"`
	PolicyPermissions []linkedPolicyResource `gorethink:"ladon_policy_permissions"`
	PolicyResources   []linkedPolicyResource `gorethink:"ladon_policy_resources"`
}

type linkedPolicyResource struct {
	Compiled string `gorethink:"compiled"`
	Template string `gorethink:"template"`
}

type Store struct {
	session *rdb.Session
}

func New(session *rdb.Session) *Store {
	return &Store{session}
}

func (s *Store) CreateTables() error {
	exists, err := s.tableExists(policyTableName)
	if err == nil && !exists {
		_, err := rdb.TableCreate(policyTableName).RunWrite(s.session)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

// TableExists check if table(s) exists in database
func (s *Store) tableExists(table string) (bool, error) {

	res, err := rdb.TableList().Run(s.session)
	if err != nil {
		return false, err
	}
	defer res.Close()

	if res.IsNil() {
		return false, nil
	}

	var tableDB string
	for res.Next(&tableDB) {
		if table == tableDB {
			return true, nil
		}
	}

	return false, nil
}

func (s *Store) Create(policy policy.Policy) (err error) {
	conditions := []byte("[]")
	if policy.GetConditions() != nil {
		cs := policy.GetConditions()
		conditions, err = json.Marshal(&cs)
		if err != nil {
			return err
		}
	}

	policySubjects, err := createLink(policy, policy.GetSubjects())
	if err != nil {
		return err
	}
	policyPermissions, err := createLink(policy, policy.GetPermissions())
	if err != nil {
		return err
	}
	policyResources, err := createLink(policy, policy.GetResources())
	if err != nil {
		return err
	}

	dbPolicy := rethinkPolicy{
		ID:                policy.GetID(),
		Description:       policy.GetDescription(),
		Effect:            policy.GetEffect(),
		CreatedAt:         int64(time.Now().Unix()),
		Conditions:        conditions,
		PolicySubjects:    policySubjects,
		PolicyPermissions: policyPermissions,
		PolicyResources:   policyResources,
	}

	res, err := rdb.Table(policyTableName).Insert(dbPolicy).RunWrite(s.session)

	if err != nil {
		return err
	} else if res.Errors > 0 {
		return errors.New(res.FirstError)
	}

	return nil
}

func (s *Store) Get(id string) (policy.Policy, error) {
	// Query policy
	result, err := rdb.Table(policyTableName).Get(id).Run(s.session)

	if err != nil {
		return nil, err
	} else if result.IsNil() {
		return nil, pkg.ErrNotFound
	}

	defer result.Close()

	var p rethinkPolicy
	err = result.One(&p)
	if err != nil {
		return nil, err
	}

	orgPolicy := policy.DefaultPolicy{
		ID:          p.ID,
		Description: p.Description,
		Effect:      p.Effect,
		Permissions: getLinked(p.PolicyPermissions),
		Subjects:    getLinked(p.PolicySubjects),
		Resources:   getLinked(p.PolicyResources),
	}

	if err := json.Unmarshal(p.Conditions, &orgPolicy.Conditions); err != nil {
		return nil, err
	}

	return &orgPolicy, nil
}

func (s *Store) Delete(id string) error {
	if _, err := rdb.Table(policyTableName).Get(id).Delete().RunWrite(s.session); err != nil {
		return err
	}
	return nil
}

func (s *Store) FindPoliciesForSubject(subject string) (policies []policy.Policy, err error) {
	// Query all appliccable policies for subject
	res, err := rdb.Table(policyTableName).Filter(func(policy rdb.Term) rdb.Term {
		return policy.Field("ladon_policy_subjects").Contains(func(policy_subject rdb.Term) rdb.Term {
			return rdb.Expr(subject).Match(policy_subject.Field("compiled"))
		}).Or(policy.Field("ladon_policy_subjects").IsEmpty())
	}).Run(s.session)

	if err != nil {
		return nil, err
	} else if res.IsNil() {
		return nil, pkg.ErrNotFound
	}

	defer res.Close()

	var p []rethinkPolicy
	err = res.All(&p)
	if err != nil {
		return nil, err
	}

	for _, tp := range p {
		tempPolicy := policy.DefaultPolicy{
			ID:          tp.ID,
			Description: tp.Description,
			Effect:      tp.Effect,
			Permissions: getLinked(tp.PolicyPermissions),
			Subjects:    getLinked(tp.PolicySubjects),
			Resources:   getLinked(tp.PolicyResources),
		}

		if err := json.Unmarshal(tp.Conditions, &tempPolicy.Conditions); err != nil {
			return nil, err
		}
		policies = append(policies, &tempPolicy)
	}

	return policies, nil
}

func getLinked(resourceData []linkedPolicyResource) []string {
	templates := make([]string, len(resourceData))

	for i, data := range resourceData {
		templates[i] = data.Template
	}

	return templates
}

func createLink(p policy.Policy, templates []string) ([]linkedPolicyResource, error) {
	resSlice := make([]linkedPolicyResource, len(templates))
	for i, template := range templates {
		reg, err := compiler.CompileRegex(template, p.GetStartDelimiter(), p.GetEndDelimiter())

		if err != nil {
			return nil, err
		}

		policyData := linkedPolicyResource{
			Compiled: reg.String(),
			Template: template,
		}

		resSlice[i] = policyData
	}
	return resSlice, nil
}
