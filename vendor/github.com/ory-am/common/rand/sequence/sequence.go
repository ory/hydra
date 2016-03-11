package sequence

import (
	"crypto/rand"
	"math/big"
)

var rander = rand.Reader // random function

var (
	AlphaNum      = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	Alpha         = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	AlphaLowerNum = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	AlphaUpperNum = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	AlphaLower    = []rune("abcdefghijklmnopqrstuvwxyz")
	AlphaUpper    = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	Numeric       = []rune("0123456789")
)

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
