package policy

import (
	"github.com/ory-am/ladon"
	"github.com/go-errors/errors"
	r "github.com/dancannon/gorethink"
	"sync"
	"golang.org/x/net/context"
	"github.com/ory-am/hydra/pkg"
	"encoding/json"
)

// stupid hack
type rdbSchema struct {
	ID          string     `json:"id" gorethink:"id"`
	Description string     `json:"description" gorethink:"description"`
	Subjects    []string   `json:"subjects" gorethink:"subjects"`
	Effect      string     `json:"effect" gorethink:"effect"`
	Resources   []string   `json:"resources" gorethink:"resources"`
	Actions     []string   `json:"actions" gorethink:"actions"`
	Conditions  json.RawMessage `json:"conditions" gorethink:"conditions"`
}

func rdbToPolicy(s *rdbSchema) (*ladon.DefaultPolicy, error) {
	if s == nil {
		return nil, nil
	}

	ret := &ladon.DefaultPolicy{
		ID: s.ID,
		Description: s.Description,
		Subjects: s.Subjects,
		Effect: s.Effect,
		Resources: s.Resources,
		Actions: s.Actions,
		Conditions: ladon.Conditions{},
	}

	if err := ret.Conditions.UnmarshalJSON(s.Conditions); err != nil {
		return nil, errors.New(err)
	}

	return ret, nil

}

func rdbFromPolicy(p ladon.Policy) (*rdbSchema, error) {
	cs, err := p.GetConditions().MarshalJSON()
	if err != nil {
		return nil, err
	}
	return &rdbSchema{
		ID: p.GetID(),
		Description: p.GetDescription(),
		Subjects: p.GetSubjects(),
		Effect: p.GetEffect(),
		Resources: p.GetResources(),
		Actions: p.GetActions(),
		Conditions: cs,
	}, err
}

type RethinkManager struct {
	Session  *r.Session
	Table    r.Term
	sync.RWMutex

	Policies map[string]ladon.Policy
}

func (m *RethinkManager) ColdStart() error {
	m.Policies = map[string]ladon.Policy{}
	clients, err := m.Table.Run(m.Session)
	if err != nil {
		return errors.New(err)
	}

	var policy ladon.DefaultPolicy
	m.Lock()
	defer m.Unlock()
	for clients.Next(&policy) {
		m.Policies[policy.ID] = &policy
	}

	return nil
}

func (m *RethinkManager) Create(policy ladon.Policy) error {
	if err := m.publishCreate(policy); err != nil {
		return err
	}

	return nil
}

// Get retrieves a policy.
func (m *RethinkManager) Get(id string) (ladon.Policy, error) {
	p, ok := m.Policies[id]
	if !ok {
		return nil, errors.New("Not found")
	}

	return p, nil
}

// Delete removes a policy.
func (m *RethinkManager) Delete(id string) error {
	if err := m.publishDelete(id); err != nil {
		return err
	}

	return nil
}

// Finds all policies associated with the subject.
func (m *RethinkManager) FindPoliciesForSubject(subject string) (ladon.Policies, error) {
	ps := ladon.Policies{}
	for _, p := range m.Policies {
		if ok, err := ladon.Match(p, p.GetSubjects(), subject); err != nil {
			return ladon.Policies{}, err
		} else if !ok {
			continue
		}
		ps = append(ps, p)
	}
	return ps, nil
}

func (m *RethinkManager) fetch() error {
	m.Policies = map[string]ladon.Policy{}
	policies, err := m.Table.Run(m.Session)
	if err != nil {
		return errors.New(err)
	}

	var policy ladon.DefaultPolicy
	m.Lock()
	defer m.Unlock()
	for policies.Next(&policy) {
		m.Policies[policy.ID] = &policy
	}

	return nil
}

func (m *RethinkManager) publishCreate(policy ladon.Policy) error {
	p, err := rdbFromPolicy(policy)
	if err != nil {
		return err
	}
	if _, err := m.Table.Insert(p).RunWrite(m.Session); err != nil {
		return errors.New(err)
	}
	return nil
}

func (m *RethinkManager) publishDelete(id string) error {
	if _, err := m.Table.Get(id).Delete().RunWrite(m.Session); err != nil {
		return errors.New(err)
	}
	return nil
}

func (m *RethinkManager) Watch(ctx context.Context) error {
	policies, err := m.Table.Changes().Run(m.Session)
	if err != nil {
		return errors.New(err)
	}

	go func() {
		for {
			var update = make(map[string]*rdbSchema)
			for policies.Next(&update) {
				newVal, err := rdbToPolicy(update["new_val"])
				if err != nil {
					pkg.LogError(err)
					continue
				}

				oldVal, err := rdbToPolicy(update["old_val"])
				if err != nil {
					pkg.LogError(err)
					continue
				}

				m.Lock()
				if newVal == nil && oldVal != nil {
					delete(m.Policies, oldVal.GetID())
				} else if newVal != nil && oldVal != nil {
					delete(m.Policies, oldVal.GetID())
					m.Policies[newVal.GetID()] = newVal
				} else {
					m.Policies[newVal.GetID()] = newVal
				}
				m.Unlock()
			}

			policies.Close()
			if policies.Err() != nil {
				pkg.LogError(errors.New(policies.Err()))
			}

			policies, err = m.Table.Changes().Run(m.Session)
			if err != nil {
				pkg.LogError(errors.New(policies.Err()))
			}
		}
	}()
	return nil
}