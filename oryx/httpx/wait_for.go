// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package httpx

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

// WaitForEndpoint waits for the endpoint to be available.
func WaitForEndpoint(ctx context.Context, endpoint string, opts ...retry.Option) error {
	return WaitForEndpointWithClient(ctx, http.DefaultClient, endpoint, opts...)
}

// WaitForEndpointWithClient waits for the endpoint to be available while using the given http.Client.
func WaitForEndpointWithClient(ctx context.Context, client *http.Client, endpoint string, opts ...retry.Option) error {
	return retry.Do(func() error {
		req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
		if err != nil {
			return err
		}

		res, err := client.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		if gjson.GetBytes(body, "status").String() != "ok" {
			return errors.Errorf("status is not yet ok: %s", body)
		}

		return nil
	},
		append([]retry.Option{
			retry.DelayType(retry.BackOffDelay),
			retry.Delay(time.Second),
			retry.MaxDelay(time.Second * 2),
			retry.Attempts(20),
		}, opts...)...)
}
