// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package randx

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunePatterns(t *testing.T) {
	for k, v := range []struct {
		runes       []rune
		shouldMatch string
	}{
		{Alpha, "[a-zA-Z]{52}"},
		{AlphaLower, "[a-z]{26}"},
		{AlphaUpper, "[A-Z]{26}"},
		{AlphaUpperVowels, "[AEIOUY]{6}"},
		{AlphaUpperNoVowels, "[^AEIOUY]{20}"},
		{AlphaNum, "[a-zA-Z0-9]{62}"},
		{AlphaLowerNum, "[a-z0-9]{36}"},
		{AlphaUpperNum, "[A-Z0-9]{36}"},
		{Numeric, "[0-9]{10}"},
	} {
		valid, err := regexp.Match(v.shouldMatch, []byte(string(v.runes)))
		assert.Nil(t, err, "Case %d", k)
		assert.True(t, valid, "Case %d", k)
	}
}

func TestRuneSequenceMatchesPattern(t *testing.T) {
	for k, v := range []struct {
		runes       []rune
		shouldMatch string
		length      int
	}{
		{Alpha, "[a-zA-Z]+", 25},
		{AlphaLower, "[a-z]+", 46},
		{AlphaUpper, "[A-Z]+", 21},
		{AlphaUpperVowels, "[AEIOUY]+", 12},
		{AlphaUpperNoVowels, "[^AEIOUY]+", 42},
		{AlphaNum, "[a-zA-Z0-9]+", 123},
		{AlphaLowerNum, "[a-z0-9]+", 41},
		{AlphaUpperNum, "[A-Z0-9]+", 94914},
		{Numeric, "[0-9]+", 94914},
	} {
		seq, err := RuneSequence(v.length, v.runes)
		assert.Nil(t, err, "case %d", k)
		assert.Equal(t, v.length, len(seq), "case %d", k)

		valid, err := regexp.Match(v.shouldMatch, []byte(string(seq)))
		assert.Nil(t, err, "case %d", k)
		assert.True(t, valid, "case %d\nrunes %s\nresult %s", k, v.runes, string(seq))
	}
}

func TestRuneSequenceIsPseudoUnique(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	times := 100
	runes := []rune("ab")
	length := 32
	s := make(map[string]bool)

	for i := 0; i < times; i++ {
		k, err := RuneSequence(length, runes)
		assert.Nil(t, err)
		ks := string(k)
		_, ok := s[ks]
		assert.False(t, ok)
		if ok {
			return
		}
		s[ks] = true
	}
}

func BenchmarkTestInt64(b *testing.B) {
	length := 25
	pattern := []rune("abcdefghijklmnopqrstuvwxyz")
	for i := 0; i < b.N; i++ {
		RuneSequence(length, pattern)
	}
}
