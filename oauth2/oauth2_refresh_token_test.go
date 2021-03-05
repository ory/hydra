package oauth2_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/fosite"
	hc "github.com/ory/hydra/client"
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/x/dbal"
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

	token := "234c678fed33c1d2025537ae464a1ebf7d23fc4a"
	tokenSignature := "4c7c7e8b3a77ad0c3ec846a21653c48b45dbfa31"
	testClient := hc.Client{
		OutfacingID:   uuid.New(),
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
				OutfacingID: testClient.OutfacingID,
			},
			RequestedScope: []string{"offline"},
			GrantedScope:   []string{"offline"},
			Session:        oauth2.NewSession(""),
			Form: url.Values{
				"refresh_token": []string{fmt.Sprintf("%s.%s", token, tokenSignature)},
			},
		},
	}

	setupRegistries(t)

	for dbName, dbRegistry := range registries {
		if dbName == "memory" {
			// todo check why sqlite fails with "no such table: hydra_oauth2_refresh \n sqlite create"
			// should be fine though as nobody should use sqlite in production
			continue
		}
		ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
		require.NoError(t, dbRegistry.OAuth2Storage().(clientCreator).CreateClient(ctx, &testClient))
		require.NoError(t, dbRegistry.OAuth2Storage().CreateRefreshTokenSession(ctx, tokenSignature, request))
		_, err := dbRegistry.OAuth2Storage().GetRefreshTokenSession(ctx, tokenSignature, nil)
		require.NoError(t, err)
		provider := dbRegistry.OAuth2Provider()
		storageVersion := dbVersion(t, ctx, dbRegistry)

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
						time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
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
					if e := (&fosite.RFC6749Error{}); errors.As(err, &e) {
						switch e.ErrorField {

						// change logic below when the refresh handler starts returning 'fosite.ErrInvalidRequest' for other reasons.
						// as of now, this error is only returned due to concurrent transactions competing to refresh using the same token.

						case fosite.ErrInvalidRequest.ErrorField, fosite.ErrServerError.ErrorField:
							// the error description copy is defined by RFC 6749 and should not be different regardless of
							// the underlying transactional aware storage backend used by hydra
							assert.Contains(t, []string{fosite.ErrInvalidRequest.DescriptionField, fosite.ErrServerError.DescriptionField}, e.DescriptionField)
							// the database error debug copy will be different depending on the underlying database used
							switch dbName {
							case dbal.DriverMySQL:
							case dbal.DriverPostgreSQL, dbal.DriverCockroachDB:
								var matched bool
								for _, errSubstr := range []string{
									// both postgreSQL & cockroachDB return error code 40001 for consistency errors as a result of
									// using the REPEATABLE_READ isolation level
									"SQLSTATE 40001",
									// possible if one worker starts the transaction AFTER another worker has successfully
									// refreshed the token and committed the transaction
									"not_found",
									// postgres: duplicate key value violates unique constraint "hydra_oauth2_access_request_id_idx": Unable to insert or update resource because a resource with that value exists already: The request could not be completed due to concurrent access
									"duplicate key",
									// cockroach: restart transaction: TransactionRetryWithProtoRefreshError: TransactionRetryError: retry txn (RETRY_WRITE_TOO_OLD - WriteTooOld flag converted to WriteTooOldError): "sql txn" meta={id=7f069400 key=/Table/62/2/"02a55d6e-509b-4d7a-8458-5828b2f831a1"/0 pri=0.00598277 epo=0 ts=1600955431.566576173,2 min=1600955431.566576173,0 seq=6} lock=true stat=PENDING rts=1600955431.566576173,2 wto=false max=1600955431.566576173,0: Unable to serialize access due to a concurrent update in another session: The request could not be completed due to concurrent access
									"RETRY_WRITE_TOO_OLD",
									// postgres: pq: deadlock detected
									"deadlock detected",
									// postgres: pq: could not serialize access due to concurrent update: Unable to serialize access due to a concurrent update in another session: The request could not be completed due to concurrent access
									"concurrent update",
									// cockroach: this happens when there is an error with the storage
									"RETRY_WRITE_TOO_OLD",
									// refresh token reuse detection
									"token_inactive",
								} {
									if strings.Contains(e.DebugField, errSubstr) {
										matched = true
										break
									}
								}

								assert.True(t, matched, "received an unexpected kind of `%s`\n"+
									"DB version: %s\n"+
									"Error description: %s\n"+
									"Error debug: %s\n"+
									"Error hint: %s\n"+
									"Raw error: %T %+v\n"+
									"Raw cause: %T %+v",
									e.ErrorField,
									storageVersion,
									e.DescriptionField,
									e.DebugField,
									e.HintField,
									err, err,
									e, e)
							}
						default:
							// unfortunately, MySQL does not offer the same behaviour under the "REPEATABLE_READ" isolation
							// level so we have to relax this assertion just for MySQL for the time being as server_errors
							// resembling the following can be returned:
							//
							//    Error 1213: Deadlock found when trying to get lock; try restarting transaction
							if dbName != dbal.DriverMySQL {
								t.Errorf("an unexpected RFC6749 error with the name %q was returned.\n"+
									"Hint: has the refresh token error handling changed in fosite? If so, you need to add further "+
									"assertions here to cover the additional errors that are being returned by the handler.\n"+
									"DB version: %s\n"+
									"Error description: %s\n"+
									"Error debug: %s\n"+
									"Error hint: %s\n"+
									"Raw error: %+v",
									e.ErrorField,
									storageVersion,
									e.DescriptionField,
									e.DebugField,
									e.HintField,
									err)
							}
						}
					} else {
						t.Errorf("expected underlying error to be of type '*fosite.RFC6749Error', but it was "+
							"actually of type %T: %+v - DB version: %s", err, err, storageVersion)
					}
				} else {
					successCount++
				}
			}

			// IMPORTANT - skip consistency check for MySQL :(
			//
			// different DBMS's provide different consistency guarantees when using the "REPEATABLE_READ" isolation level
			// Currently, MySQL's implementation of "REPEATABLE_READ" makes it possible for multiple concurrent requests
			// to successfully utilize the same refresh token. Therefore, we skip the assertion below.
			//
			// TODO: this needs to be addressed by making it possible to use different isolation levels for various authorization
			//       flows depending on the underlying hydra storage backend. For example, if using MySQL, hydra should force
			//       the transaction isolation level to be "Serializable" when a request to the token handler is received.

			switch dbName {
			case dbal.DriverMySQL:
			case dbal.DriverPostgreSQL, dbal.DriverCockroachDB:
				require.Equal(t, 1, successCount, "CRITICAL: in test iteration %d, %d out of %d workers "+
					"were able to use the refresh token. Exactly ONE was expected to be have been successful.",
					run,
					successCount,
					workers)
			}

			// reset state for the next test iteration
			assert.NoError(t, dbRegistry.OAuth2Storage().DeleteRefreshTokenSession(ctx, tokenSignature))
			assert.NoError(t, dbRegistry.OAuth2Storage().CreateRefreshTokenSession(ctx, tokenSignature, request))
		}
	}
}

type version struct {
	Version string `db:"version"`
}

func dbVersion(t *testing.T, ctx context.Context, registry driver.Registry) string {
	var v version

	versionFunc := "version()"
	c := registry.Persister().Connection(ctx)
	if c.Dialect.Name() == "sqlite3" {
		versionFunc = "sqlite_version()"
	}
	/* #nosec G201 - versionFunc is an enum */
	require.NoError(t, registry.Persister().Connection(ctx).RawQuery(fmt.Sprintf("select %s as version", versionFunc)).First(&v))

	return v.Version
}
