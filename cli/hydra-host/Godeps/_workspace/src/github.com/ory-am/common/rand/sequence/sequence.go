package sequence

import (
    "crypto/rand"
    "io"
    "math"
    "math/big"
)

var rander = rand.Reader // random function

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

// randomBits completely fills slice b with random data.
func randomBits(b []byte) {
    if _, err := io.ReadFull(rander, b); err != nil {
        panic(err.Error()) // rand should never fail
    }
}

func runesCap(r []rune, l int) int64 {
    rs := make(map[rune]bool)
    p := float64(0)
    for _, v := range r {
        if _, ok := rs[v]; !ok {
            p++
        }
    }

    return int64(math.Pow(p, float64(l)) - 1)
}
