// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package reqlog

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"golang.org/x/sync/errgroup"
)

func TestExternalLatencyMiddleware(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		NewMiddleware().ServeHTTP(w, r, func(w http.ResponseWriter, r *http.Request) {
			var wg sync.WaitGroup

			wg.Add(3)
			for i := range 3 {
				ctx := r.Context()
				if i%3 == 0 {
					ctx = WithDisableExternalLatencyMeasurement(ctx)
				}
				go func() {
					defer StartMeasureExternalCall(ctx, "", "", time.Now())
					time.Sleep(100 * time.Millisecond)
					wg.Done()
				}()
			}
			wg.Wait()
			total := totalExternalLatency(r.Context())
			_ = json.NewEncoder(w).Encode(map[string]any{
				"total": total,
			})
		})
	}))
	defer ts.Close()

	bodies := make([][]byte, 100)
	eg := errgroup.Group{}
	for i := range bodies {
		eg.Go(func() error {
			res, err := http.Get(ts.URL)
			if err != nil {
				return err
			}
			defer res.Body.Close()
			bodies[i], err = io.ReadAll(res.Body)
			if err != nil {
				return err
			}
			return nil
		})
	}

	require.NoError(t, eg.Wait())

	for _, body := range bodies {
		actualTotal := gjson.GetBytes(body, "total").Int()
		assert.GreaterOrEqual(t, actualTotal, int64(200*time.Millisecond), string(body))
		assert.Less(t, actualTotal, int64(300*time.Millisecond), string(body))
	}
}
