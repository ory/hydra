// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package configx

import (
	stdjson "encoding/json"
	"testing"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
)

func TestKoanfMergeArray(t *testing.T) {
	k := koanf.NewWithConf(koanf.Conf{Delim: Delimiter, StrictMerge: true})
	if err := k.Load(rawbytes.Provider([]byte(`{"foo":[{"id":"bar"}]}`)), json.Parser()); err != nil {
		t.Fatal(err)
	}

	if err := k.Load(rawbytes.Provider([]byte(`{"foo":[{"key":"baz"},{"baz":"bar"}]}`)), json.Parser(), koanf.WithMergeFunc(MergeAllTypes)); err != nil {
		t.Fatal(err)
	}

	expected := `{"foo":[{"id":"bar","key":"baz"},{"baz":"bar"}]}`
	out, _ := stdjson.Marshal(k.All())
	if string(out) != expected {
		t.Fatalf("Expected %s but got: %s", expected, out)
	}
}
