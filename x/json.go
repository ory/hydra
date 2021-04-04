package x

import (
	"encoding/json"
	"fmt"

	jsonpatch "github.com/evanphx/json-patch"
)

func ApplyJSONPatch(p json.RawMessage, object interface{}, denyPaths ...string) error {
	patch, err := jsonpatch.DecodePatch(p)
	if err != nil {
		return err
	}

	denySet := make(map[string]struct{})
	for _, path := range denyPaths {
		denySet[path] = struct{}{}
	}

	for _, op := range patch {
		path, err := op.Path()
		if err != nil {
			return fmt.Errorf("error parsing patch operations: %v", err)
		}
		if _, ok := denySet[path]; ok {
			return fmt.Errorf("patch includes denied path: %s", path)
		}
	}

	original, err := json.Marshal(object)
	if err != nil {
		return err
	}

	modified, err := patch.Apply(original)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(modified, object); err != nil {
		return err
	}
	return nil
}
