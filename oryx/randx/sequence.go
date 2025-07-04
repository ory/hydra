// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package randx

import (
	"crypto/rand"
	"math/big"
)

var rander = rand.Reader // random function

var (
	// AlphaNum contains runes [abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789].
	AlphaNum = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	// Alpha contains runes [abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ].
	Alpha = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	// AlphaLowerNum contains runes [abcdefghijklmnopqrstuvwxyz0123456789].
	AlphaLowerNum = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	// AlphaUpperNum contains runes [ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789].
	AlphaUpperNum = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	// AlphaLower contains runes [abcdefghijklmnopqrstuvwxyz].
	AlphaLower = []rune("abcdefghijklmnopqrstuvwxyz")
	// AlphaUpperVowels contains runes [AEIOUY].
	AlphaUpperVowels = []rune("AEIOUY")
	// AlphaUpperNoVowels contains runes [BCDFGHJKLMNPQRSTVWXZ].
	AlphaUpperNoVowels = []rune("BCDFGHJKLMNPQRSTVWXZ")
	// AlphaUpper contains runes [ABCDEFGHIJKLMNOPQRSTUVWXYZ].
	AlphaUpper = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	// Numeric contains runes [0123456789].
	Numeric = []rune("0123456789")
	// AlphaNumNoAmbiguous is equivalent to AlphaNum but without visually ambiguous characters [0Oo1IlB8S5Z2].
	AlphaNumNoAmbiguous = []rune("abcdefghijkmnpqrstuvwxyzACDEFGHJKLMNPQRTUVWXY34679")
)

// RuneSequence returns a random sequence using the defined allowed runes.
func RuneSequence(l int, allowedRunes []rune) (seq []rune, err error) {
	c := big.NewInt(int64(len(allowedRunes)))
	seq = make([]rune, l)

	for i := 0; i < l; i++ {
		r, err := rand.Int(rander, c)
		if err != nil {
			return seq, err
		}
		rn := allowedRunes[r.Uint64()]
		seq[i] = rn
	}

	return seq, nil
}

// MustString returns a random string sequence using the defined runes. Panics on error.
func MustString(l int, allowedRunes []rune) string {
	seq, err := RuneSequence(l, allowedRunes)
	if err != nil {
		panic(err)
	}
	return string(seq)
}
