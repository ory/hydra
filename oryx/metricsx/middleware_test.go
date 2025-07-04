// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package metricsx

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnonymizePath(t *testing.T) {
	m := &Service{
		o: &Options{WhitelistedPaths: []string{"/keys"}},
	}

	assert.Equal(t, "/keys", m.anonymizePath("/keys/1234/sub-path"))
	assert.Equal(t, "/keys", m.anonymizePath("/keys/1234"))
	assert.Equal(t, "/keys", m.anonymizePath("/keys"))
	assert.Equal(t, "/", m.anonymizePath("/not-keys"))
}

func TestAnonymizeQuery(t *testing.T) {
	m := &Service{}

	assert.EqualValues(t, "foo=2ec879270efe890972d975251e9d454f4af49df1f07b4317fd5b6ae90de4c774&foo=1864a573566eba1b9ddab79d8f4bab5a39c938918a21b80a64ae1c9c12fa9aa2&foo2=186084f6bd8e222bedade9439d6ae69ed274b954eeebe9b54fd5f47e54dd7675&foo2=1ee7158281cc3b5a27de4c337e07987e8677f5f687a4671ca369b79c653d379d", m.anonymizeQuery(url.Values{
		"foo":  []string{"bar", "baz"},
		"foo2": []string{"bar2", "baz2"},
	}, "somesupersaltysalt"))
	assert.EqualValues(t, "", m.anonymizeQuery(url.Values{
		"foo": []string{},
	}, "somesupersaltysalt"))
	assert.EqualValues(t, "foo=", m.anonymizeQuery(url.Values{
		"foo": []string{""},
	}, "somesupersaltysalt"))
	assert.EqualValues(t, "", m.anonymizeQuery(url.Values{}, "somesupersaltysalt"))
}
