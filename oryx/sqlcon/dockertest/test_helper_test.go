// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package dockertest

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
)

type mockPool struct{ mock.Mock }

func (p *mockPool) Purge(r *dockertest.Resource) error {
	args := p.Called(r)
	return args.Error(0)
}

func (p *mockPool) Run(repository string, tag string, env []string) (*dockertest.Resource, error) {
	args := p.Called(repository, tag, env)
	return args.Get(0).(*dockertest.Resource), args.Error(1)
}

func (p *mockPool) RunWithOptions(opts *dockertest.RunOptions, hcOpts ...func(*dc.HostConfig)) (*dockertest.Resource, error) {
	args := p.Called(opts, hcOpts)
	return args.Get(0).(*dockertest.Resource), args.Error(1)
}

func setupMock(t *testing.T) *mockPool {
	m := &mockPool{}
	m.Test(t)
	pool = m
	return m
}

func TestRunTestDBs(t *testing.T) {
	tc := []struct {
		name   string
		env    string
		testFn func(t testing.TB) string
	}{
		{
			name:   "postgres",
			env:    "TEST_DATABASE_POSTGRESQL",
			testFn: RunTestPostgreSQL,
		}, {
			name:   "mysql",
			env:    "TEST_DATABASE_MYSQL",
			testFn: RunTestMySQL,
		}, {
			name:   "cockroachdb",
			env:    "TEST_DATABASE_COCKROACHDB",
			testFn: RunTestCockroachDB,
		},
	}

	for _, tt := range tc {
		t.Run("db="+tt.name, func(t *testing.T) {
			t.Run("case=from_docker", func(t *testing.T) {
				m := setupMock(t)
				t.Setenv(tt.env, "")
				resource := &dockertest.Resource{}
				m.On("Run", mock.Anything, mock.Anything, mock.Anything).Return(resource, nil)
				m.On("RunWithOptions", mock.Anything, mock.Anything).Return(resource, nil)
				m.On("Purge", resource).Return(nil)

				t.Run("in test", func(t *testing.T) { tt.testFn(t) })

				m.AssertCalled(t, "Purge", resource)
			})

			t.Run("case=from_env", func(t *testing.T) {
				m := setupMock(t)
				t.Setenv(tt.env, "conn")

				tt.testFn(t)

				m.AssertExpectations(t)
			})
		})
	}
}
