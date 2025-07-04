package hasherx_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/ory/x/hasherx"
	"github.com/ory/x/randx"
)

func TestPBKDF2Performance(t *testing.T) {
	for _, iters := range []uint32{
		100, 1000, 10000, 25000, 100000, 1000000,
	} {
		t.Run(fmt.Sprintf("%d", iters), func(t *testing.T) {
			runPBKDF2(t, iters, 100)
		})
	}
}

func runPBKDF2(t *testing.T, iterations uint32, hashCount uint32) {
	c := gomock.NewController(t)
	t.Cleanup(c.Finish)
	reg := NewMockPBKDF2Configurator(c)
	reg.EXPECT().HasherPBKDF2Config(gomock.Any()).Return(&hasherx.PBKDF2Config{
		Algorithm:  "sha256",
		Iterations: iterations,
		SaltLength: 32,
		KeyLength:  32,
	}).AnyTimes()

	pw := randx.MustString(16, randx.AlphaLower)
	hasher := hasherx.NewHasherPBKDF2(reg)
	ctx := context.Background()

	var err error
	start := time.Now()
	for i := uint32(0); i < hashCount; i++ {
		if _, err = hasher.Generate(ctx, []byte(pw)); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	}
	end := time.Now()
	diff := end.Sub(start).Round(time.Millisecond)
	t.Logf("%d hashes in %s with %d iterations, %dms per hash", hashCount, diff, iterations, diff.Milliseconds()/int64(hashCount))
}
