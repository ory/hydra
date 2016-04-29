package policy

import (
	"net/http"
	"net/url"

	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/ladon"
)

type HTTPManager struct {
	Endpoint *url.URL

	Client *http.Client
}

// Create persists the policy.
func (m *HTTPManager) Create(policy ladon.Policy) error {
	var r = pkg.NewSuperAgent(m.Endpoint.String())
	r.Client = m.Client
	if err := r.POST(policy); err != nil {
		return nil
	}

	return nil
}

// Get retrieves a policy.
func (m *HTTPManager) Get(id string) (ladon.Policy, error) {
	var policy ladon.DefaultPolicy
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, id).String())
	r.Client = m.Client
	if err := r.GET(&policy); err != nil {
		return nil, err
	}

	return &policy, nil

}

// Delete removes a policy.
func (m *HTTPManager) Delete(id string) error {
	var r = pkg.NewSuperAgent(pkg.JoinURL(m.Endpoint, id).String())
	r.Client = m.Client
	if err := r.DELETE(); err != nil {
		return err
	}
	return nil
}

// Finds all policies associated with the subject.
func (m *HTTPManager) FindPoliciesForSubject(subject string) (ladon.Policies, error) {
	var policies []*ladon.DefaultPolicy
	var r = pkg.NewSuperAgent(m.Endpoint.String() + "?subject=" + subject)
	r.Client = m.Client
	if err := r.GET(&policies); err != nil {
		return nil, err
	}

	return func(ps []*ladon.DefaultPolicy) (r ladon.Policies) {
		r = make(ladon.Policies, len(ps))
		for k, p := range ps {
			r[k] = ladon.Policy(p)
		}
		return r
	}(policies), nil
}
