// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jsonx

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/gobwas/glob"
	"github.com/pkg/errors"

	"github.com/ory/x/pointerx"
)

var opAllowList = map[string]struct{}{
	"add":     {},
	"remove":  {},
	"replace": {},
}

func isUnsupported(op jsonpatch.Operation) bool {
	_, ok := opAllowList[op.Kind()]

	return !ok
}

func isElementAccess(path string) bool {
	if path == "" {
		return false
	}
	elements := strings.Split(path, "/")
	lastElement := elements[len(elements)-1:][0]
	if lastElement == "-" {
		return true
	}
	if _, err := strconv.Atoi(lastElement); err == nil {
		return true
	}

	return false
}

// ApplyJSONPatch applies a JSON patch to an object and returns the modified
// object. The original object is not modified. It returns an error if the patch
// is invalid or if the patch includes paths that are denied. denyPaths is a
// list of path globs (interpreted with [glob.Compile] that are not allowed to
// be patched.
func ApplyJSONPatch[T any](p json.RawMessage, object T, denyPaths ...string) (result T, err error) {
	patch, err := jsonpatch.DecodePatch(p)
	if err != nil {
		return result, errors.WithStack(err)
	}

	denyPattern := fmt.Sprintf("{%s}", strings.ToLower(strings.Join(denyPaths, ",")))
	matcher, err := glob.Compile(denyPattern, '/')
	if err != nil {
		return result, errors.WithStack(err)
	}

	for _, op := range patch {
		// Some operations are buggy, see https://github.com/evanphx/json-patch/pull/158
		if isUnsupported(op) {
			return result, errors.Errorf("unsupported operation: %s", op.Kind())
		}
		path, err := op.Path()
		if err != nil {
			return result, errors.Errorf("error parsing patch operations: %v", err)
		}
		if matcher.Match(strings.ToLower(path)) {
			return result, errors.Errorf("patch includes denied path: %s", path)
		}

		// JSON patch officially rejects replacing paths that don't exist, but we want to be more tolerant.
		// Therefore, we will ensure that all paths that we want to replace exist in the original document.
		if op.Kind() == "replace" && !isElementAccess(path) {
			op["op"] = pointerx.Ptr(json.RawMessage(`"add"`))
		}
	}

	original, err := json.Marshal(object)
	if err != nil {
		return result, errors.WithStack(err)
	}

	options := jsonpatch.NewApplyOptions()
	options.EnsurePathExistsOnAdd = true

	modified, err := patch.ApplyWithOptions(original, options)
	if err != nil {
		return result, errors.WithStack(err)
	}

	err = json.Unmarshal(modified, &result)
	return result, errors.WithStack(err)
}
