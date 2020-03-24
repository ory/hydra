package oauth2_test

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/ory/fosite"
	hc "github.com/ory/hydra/client"
	"github.com/ory/hydra/internal"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/x"
	"github.com/ory/x/errorsx"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreateRefreshTokenSessionStress is a sanity test to verify the fix for https://github.com/ory/hydra/issues/1719 &
// https://github.com/ory/hydra/issues/1735.
// It currently only deals with Postgres as that was what the issue was based on due to the default isolation level used
// by the storage engine.
func TestCreateRefreshTokenSessionStress(t *testing.T) {
	if testing.Short() {
		return
	}

	// number of iterations this test will make to ensure everything is working as expected. This test is aiming to
	// prove correct behaviour when the handler is getting hit with the same refresh token in concurrent requests. Given
	// that problems that may occur in this scenario are "racey" in nature, it is important to run this test several times
	// so to minimize the probability were we pass due to sheer luck.
	testRuns := 5
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
			RequestedAt: time.Now(),
			ID:          uuid.New(),
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
	require.NoError(t, postgresRegistry.OAuth2Storage().CreateRefreshTokenSession(ctx, tokenSignature, request))
	_, err = postgresRegistry.OAuth2Storage().GetRefreshTokenSession(ctx, tokenSignature, nil)
	require.NoError(t, err)
	provider := postgresRegistry.OAuth2Provider()

	var wg sync.WaitGroup
	for run := 0; run < testRuns; run++ {
		barrier := make(chan struct{})
		errorsCh := make(chan error, workers)

		go func() {
			for w := 0; w < workers; w++ {
				wg.Add(1)
				go func(run, worker int) {
					defer wg.Done()
					ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
					// all workers will block here until the for loop above has launched all the worker go-routines
					// this is to ensure we fire all the workers off at the same
					<-barrier
					_, err := provider.NewAccessResponse(ctx, request)
					errorsCh <- err
				}(run, w)
			}

			// wait until all workers have completed their work
			wg.Wait()
			close(errorsCh)
		}()

		// let the race begin!
		// all worker go-routines will now attempt to hit the "NewAccessResponse" method
		close(barrier)

		// process worker results

		// successCount is the number of workers that were able to call "NewAccessResponse" without receiving an error.
		// if the successCount at the end of a test run is bigger than one, it means that multiple access/refresh tokens
		// were issued using the same refresh token! - https://knowyourmeme.com/memes/scared-hamster
		var successCount int
		for err := range errorsCh {
			if err != nil {
				switch cause := errorsx.Cause(err).(type) {
				case *fosite.RFC6749Error:
					switch cause.Name {
					case fosite.ErrInvalidRequest.Name:
						assert.Equal(t, fosite.ErrInvalidRequest.Description, cause.Description)
						assert.Contains(t, cause.Debug, "could not serialize access due to concurrent update (SQLSTATE 40001)")
					default:
						t.Errorf("an unexpected RFC6749 error with the name %q was returned.\n"+
							"Hint: has the refresh token error handling changed in fosite? If so, you need to add further "+
							"assertions here to cover the additional errors that are being returned by the handler.\n\n"+
							"Error description: %s\n"+
							"Error debug: %s\n"+
							"Error hint: %s\n"+
							"Raw error: %+v",
							cause.Name,
							cause.Description,
							cause.Debug,
							cause.Hint,
							err)
					}
				default:
					t.Errorf("expected underlying error to be of type '*fosite.RFC6749Error', but it was "+
						"actually of type %T: %+v", err, err)
				}
			} else {
				successCount++
			}
		}

		assert.Equal(t, 1, successCount, "CRITICAL: in test iteration %d, %d out of %d workers "+
			"were able to use the refresh token. Exactly ONE was expected to be have been successful.",
			run,
			successCount,
			workers)

		// reset state for the next test iteration
		_ = postgresRegistry.OAuth2Storage().RevokeRefreshToken(ctx, request.ID)
		_ = postgresRegistry.OAuth2Storage().CreateRefreshTokenSession(ctx, tokenSignature, request)
	}
}
