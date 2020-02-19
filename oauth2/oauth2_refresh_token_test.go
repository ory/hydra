package oauth2_test

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ory/fosite"
	hc "github.com/ory/hydra/client"
	"github.com/ory/hydra/internal"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/x"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreateRefreshTokenSessionBug reproduces the bug raised in https://github.com/ory/hydra/issues/1719
// Once this bug is fixed, this test should start failing. It currently only deals with Postgres as that was
// what the issue was based on due to the default isolation level used by the storage engine.
func TestCreateRefreshTokenSessionBug(t *testing.T) {
	if testing.Short() {
		return
	}

	// number of workers that will concurrently hit the 'CreateRefreshTokenSession' method using the same refresh token.
	// don't set this value to be too high as it will result in connection failures to the DB instance. The test is designed such that
	// it will retry in the event we get unlucky and a transaction completes successfully prior to other requests getting past the
	// first read.
	workers := 10
	postgresRegistry := internal.NewRegistrySQL(internal.NewConfigurationWithDefaults(), connectToPG(t))
	x.CleanSQL(t, postgresRegistry.DB())
	_, err := postgresRegistry.CreateSchemas("postgres")
	require.NoError(t, err)

	token := "234c678fed33c1d2025537ae464a1ebf7d23fc4a"
	tokenSignature := "4c7c7e8b3a77ad0c3ec846a21653c48b45dbfa31"
	testClient := hc.Client{
		ClientID:      uuid.New(),
		Secret:        "secret",
		ResponseTypes: []string{"id_token", "code", "token"},
		GrantTypes:    []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"},
		Scope:         "hydra offline openid",
		Audience:      []string{"https://api.ory.sh/"},
	}

	request := &fosite.AccessRequest{
		GrantTypes: []string{
			"refresh_token",
		},
		Request: fosite.Request{
			ID: uuid.New(),
			Client: &hc.Client{
				ClientID: testClient.ClientID,
			},
			RequestedScope: []string{"offline"},
			GrantedScope:   []string{"offline"},
			Session:        oauth2.NewSession(""),
			Form: url.Values{
				"refresh_token": []string{fmt.Sprintf("%s.%s", token, tokenSignature)},
			},
		},
	}

	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	require.NoError(t, postgresRegistry.OAuth2Storage().(clientCreator).CreateClient(ctx, &testClient))
	assert.NoError(t, postgresRegistry.OAuth2Storage().CreateRefreshTokenSession(ctx, tokenSignature, request))
	_, err = postgresRegistry.OAuth2Storage().GetRefreshTokenSession(ctx, tokenSignature, nil)
	assert.NoError(t, err)
	provider := postgresRegistry.OAuth2Provider()

	var mutex sync.RWMutex
	successCh := make(chan struct{})

	for {
		barrier := make(chan struct{})
		workerCtx, stopWorkers := context.WithCancel(ctx)
		for i := 0; i < workers; i++ {
			go func() {
				<-barrier

				_, err := provider.NewAccessResponse(ctx, request)
				mutex.Lock()
				defer mutex.Unlock()

				select {
				case <-workerCtx.Done():
					return
				default:
					switch err := errors.Cause(err).(type) {
					case *fosite.RFC6749Error:

						if strings.Contains(err.Debug, "pq: duplicate key value violates unique constraint") {
							stopWorkers()
							successCh <- struct{}{}
							return
						}

						if strings.Contains(err.Debug, "not_found") {
							// too late, a goroutine finished the transaction before any other goroutine got passed the first read
							// at this point, the test will be unfruitful so let's try to add another refresh session so they can
							// race again!
							postgresRegistry.OAuth2Storage().CreateRefreshTokenSession(ctx, tokenSignature, request)
							stopWorkers()
						}
					}
				}
			}()
		}

		// let the race begin!
		close(barrier)

		// keep going until either the test timesout or we reproduce the bug - whichever comes first!
		select {
		case <-ctx.Done():
			t.Errorf("failed to reproduce bug https://github.com/ory/hydra/issues/1719")
			return
		case <-successCh:
			return
		}
	}
}
