// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package oauth2

import (
	"sync"

	"github.com/ory/hydra/pkg"
	"github.com/pkg/errors"
)

type ConsentRequestMemoryManager struct {
	requests map[string]ConsentRequest
	sync.RWMutex
}

func NewConsentRequestMemoryManager() *ConsentRequestMemoryManager {
	return &ConsentRequestMemoryManager{requests: map[string]ConsentRequest{}}
}

func (m *ConsentRequestMemoryManager) PersistConsentRequest(session *ConsentRequest) error {
	m.Lock()
	defer m.Unlock()
	m.requests[session.ID] = *session
	return nil
}

func (m *ConsentRequestMemoryManager) AcceptConsentRequest(id string, payload *AcceptConsentRequestPayload) error {
	session, err := m.GetConsentRequest(id)
	if err != nil {
		return err
	}

	session.Subject = payload.Subject
	session.AccessTokenExtra = payload.AccessTokenExtra
	session.IDTokenExtra = payload.IDTokenExtra
	session.Consent = ConsentRequestAccepted
	session.GrantedScopes = payload.GrantScopes

	return m.PersistConsentRequest(session)
}

func (m *ConsentRequestMemoryManager) RejectConsentRequest(id string, payload *RejectConsentRequestPayload) error {
	session, err := m.GetConsentRequest(id)
	if err != nil {
		return err
	}

	session.Consent = ConsentRequestRejected
	session.DenyReason = payload.Reason
	return m.PersistConsentRequest(session)
}

func (m *ConsentRequestMemoryManager) GetConsentRequest(id string) (*ConsentRequest, error) {
	m.RLock()
	defer m.RUnlock()
	if session, found := m.requests[id]; !found {
		return nil, errors.Wrap(pkg.ErrNotFound, "")
	} else {
		return &session, nil
	}
}
