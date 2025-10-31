// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/ory/hydra/v2/fosite"
)

func TestMemoryStore_Authenticate(t *testing.T) {
	type fields struct {
		Users      map[string]MemoryUserRelation
		usersMutex sync.RWMutex
	}
	type args struct {
		in0    context.Context
		name   string
		secret string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "invalid_password",
			args: args{
				name:   "peter",
				secret: "invalid",
			},
			fields: fields{
				Users: map[string]MemoryUserRelation{
					"peter": {
						Username: "peter",
						Password: "secret",
					},
				},
			},
			// ResourceOwnerPasswordCredentialsGrantHandler expects ErrNotFound
			wantErr: fosite.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemoryStore{
				Users:      tt.fields.Users,
				usersMutex: tt.fields.usersMutex,
			}
			if _, err := s.Authenticate(tt.args.in0, tt.args.name, tt.args.secret); err == nil || !errors.Is(err, tt.wantErr) {
				t.Errorf("Authenticate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
