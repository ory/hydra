package x

import (
	"encoding/json"

	jsonpatch "github.com/evanphx/json-patch"
)

func ApplyJSONPatch(p json.RawMessage, object interface{}) error {
	patch, err := jsonpatch.DecodePatch(p)
	if err != nil {
		return err
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
