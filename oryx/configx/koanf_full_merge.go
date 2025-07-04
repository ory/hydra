// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package configx

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/tidwall/sjson"

	"github.com/ory/x/jsonx"
)

func MergeAllTypes(src, dst map[string]interface{}) error {
	rawSrc, err := json.Marshal(src)
	if err != nil {
		return errors.WithStack(err)
	}

	dstSrc, err := json.Marshal(dst)
	if err != nil {
		return errors.WithStack(err)
	}

	keys := jsonx.Flatten(rawSrc)
	for key, value := range keys {
		dstSrc, err = sjson.SetBytes(dstSrc, key, value)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return errors.WithStack(json.Unmarshal(dstSrc, &dst))
}
