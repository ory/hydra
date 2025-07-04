package hasherx_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/inhies/go-bytesize"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/ory/x/hasherx"
)

func mkpw(t *testing.T, length int) []byte {
	pw := make([]byte, length)
	_, err := rand.Read(pw)
	require.NoError(t, err)
	return pw
}

func TestArgonHasher(t *testing.T) {
	c := gomock.NewController(t)
	t.Cleanup(c.Finish)
	reg := NewMockArgon2Configurator(c)
	reg.EXPECT().HasherArgon2Config(gomock.Any()).Return(&hasherx.Argon2Config{
		Memory:      bytesize.KB,
		Iterations:  2,
		Parallelism: 1,
		SaltLength:  32,
		KeyLength:   32,
	}).AnyTimes()

	for k, pw := range [][]byte{
		mkpw(t, 8),
		mkpw(t, 16),
		mkpw(t, 32),
		mkpw(t, 64),
		mkpw(t, 128),
	} {
		k := k
		pw := pw
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()
			for kk, h := range []hasherx.Hasher{
				hasherx.NewHasherArgon2(reg),
			} {
				kk := kk
				h := h
				t.Run(fmt.Sprintf("hasher=%T/password=%d", h, kk), func(t *testing.T) {
					t.Parallel()
					hs, err := h.Generate(context.Background(), pw)
					require.NoError(t, err)
					assert.NotEqual(t, pw, hs)

					t.Logf("hash: %s", hs)
					require.NoError(t, hasherx.CompareArgon2id(context.Background(), pw, hs))

					mod := make([]byte, len(pw))
					copy(mod, pw)
					mod[len(pw)-1] = ^pw[len(pw)-1]
					require.Error(t, hasherx.CompareArgon2id(context.Background(), mod, hs))
				})
			}
		})
	}
}

func newBCryptRegistry(t *testing.T) *MockBCryptConfigurator {
	c := gomock.NewController(t)
	t.Cleanup(c.Finish)
	reg := NewMockBCryptConfigurator(c)
	reg.EXPECT().HasherBcryptConfig(gomock.Any()).Return(&hasherx.BCryptConfig{Cost: 4}).AnyTimes()
	return reg
}

func TestBcryptHasherGeneratesErrorWhenPasswordIsLong(t *testing.T) {
	hasher := hasherx.NewHasherBcrypt(newBCryptRegistry(t))

	password := mkpw(t, 73)
	res, err := hasher.Generate(context.Background(), password)

	assert.Error(t, err, "password is too long")
	assert.Nil(t, res)
}

func TestBcryptHasherGeneratesHash(t *testing.T) {
	for k, pw := range [][]byte{
		mkpw(t, 8),
		mkpw(t, 16),
		mkpw(t, 32),
		mkpw(t, 64),
		mkpw(t, 72),
	} {
		k := k
		pw := pw
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()
			hasher := hasherx.NewHasherBcrypt(newBCryptRegistry(t))
			hs, err := hasher.Generate(context.Background(), pw)

			assert.Nil(t, err)
			assert.True(t, hasher.Understands(hs))

			// Valid format: $2a$12$[22 character salt][31 character hash]
			assert.Equal(t, 60, len(string(hs)), "invalid bcrypt hash length")
			assert.Equal(t, "$2a$04$", string(hs)[:7], "invalid bcrypt identifier")
		})
	}
}

func TestComparatorBcryptFailsWhenPasswordIsTooLong(t *testing.T) {
	password := mkpw(t, 73)
	err := hasherx.CompareBcrypt(context.Background(), password, []byte("hash"))

	assert.Error(t, err, "password is too long")
}

func TestComparatorBcryptSuccess(t *testing.T) {
	for k, pw := range [][]byte{
		mkpw(t, 8),
		mkpw(t, 16),
		mkpw(t, 32),
		mkpw(t, 64),
		mkpw(t, 72),
	} {
		k := k
		pw := pw
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()
			hasher := hasherx.NewHasherBcrypt(newBCryptRegistry(t))

			hs, err := hasher.Generate(context.Background(), pw)

			assert.Nil(t, err)
			assert.True(t, hasher.Understands(hs))

			err = hasherx.CompareBcrypt(context.Background(), pw, hs)
			assert.Nil(t, err, "hash validation fails")
		})
	}
}

func TestComparatorBcryptFail(t *testing.T) {
	for k, pw := range [][]byte{
		mkpw(t, 8),
		mkpw(t, 16),
		mkpw(t, 32),
		mkpw(t, 64),
		mkpw(t, 72),
	} {
		k := k
		pw := pw
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()
			mod := make([]byte, len(pw))
			copy(mod, pw)
			mod[len(pw)-1] = ^pw[len(pw)-1]

			err := hasherx.CompareBcrypt(context.Background(), pw, mod)
			assert.Error(t, err)
		})
	}
}

func TestPbkdf2Hasher(t *testing.T) {
	for k, pw := range [][]byte{
		mkpw(t, 8),
		mkpw(t, 16),
		mkpw(t, 32),
		mkpw(t, 64),
		mkpw(t, 128),
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			t.Parallel()
			for kk, config := range []*hasherx.PBKDF2Config{
				{
					Algorithm:  "sha1",
					Iterations: 100000,
					SaltLength: 32,
					KeyLength:  32,
				},
				{
					Algorithm:  "sha224",
					Iterations: 100000,
					SaltLength: 32,
					KeyLength:  32,
				},
				{
					Algorithm:  "sha256",
					Iterations: 100000,
					SaltLength: 32,
					KeyLength:  32,
				},
				{
					Algorithm:  "sha384",
					Iterations: 100000,
					SaltLength: 32,
					KeyLength:  32,
				},
				{
					Algorithm:  "sha512",
					Iterations: 100000,
					SaltLength: 32,
					KeyLength:  32,
				},
			} {
				kk := kk
				config := config
				t.Run(fmt.Sprintf("config=%T/password=%d", config.Algorithm, kk), func(t *testing.T) {
					t.Parallel()
					c := gomock.NewController(t)
					t.Cleanup(c.Finish)
					reg := NewMockPBKDF2Configurator(c)
					reg.EXPECT().HasherPBKDF2Config(gomock.Any()).Return(config).AnyTimes()

					hasher := hasherx.NewHasherPBKDF2(reg)
					hs, err := hasher.Generate(context.Background(), pw)
					require.NoError(t, err)
					assert.NotEqual(t, pw, hs)

					t.Logf("hash: %s", hs)
					require.NoError(t, hasherx.ComparePbkdf2(context.Background(), pw, hs))

					assert.True(t, hasher.Understands(hs))

					mod := make([]byte, len(pw))
					copy(mod, pw)
					mod[len(pw)-1] = ^pw[len(pw)-1]
					require.Error(t, hasherx.ComparePbkdf2(context.Background(), mod, hs))
				})
			}
		})
	}
}

func TestCompare(t *testing.T) {
	assert.Error(t, hasherx.Compare(context.Background(), []byte("test"), []byte("$unknown$12$o6hx.Wog/wvFSkT/Bp/6DOxCtLRTDj7lm9on9suF/WaCGNVHbkfL6")))

	assert.Nil(t, hasherx.Compare(context.Background(), []byte("test"), []byte("$2a$12$o6hx.Wog/wvFSkT/Bp/6DOxCtLRTDj7lm9on9suF/WaCGNVHbkfL6")))
	assert.Nil(t, hasherx.CompareBcrypt(context.Background(), []byte("test"), []byte("$2a$12$o6hx.Wog/wvFSkT/Bp/6DOxCtLRTDj7lm9on9suF/WaCGNVHbkfL6")))
	assert.Error(t, hasherx.Compare(context.Background(), []byte("test"), []byte("$2a$12$o6hx.Wog/wvFSkT/Bp/6DOxCtLRTDj7lm9on9suF/WaCGNVHbkfL7")))

	assert.Nil(t, hasherx.Compare(context.Background(), []byte("test"), []byte("$2a$15$GRvRO2nrpYTEuPQX6AieaOlZ4.7nMGsXpt.QWMev1zrP86JNspZbO")))
	assert.Nil(t, hasherx.CompareBcrypt(context.Background(), []byte("test"), []byte("$2a$15$GRvRO2nrpYTEuPQX6AieaOlZ4.7nMGsXpt.QWMev1zrP86JNspZbO")))
	assert.Error(t, hasherx.Compare(context.Background(), []byte("test"), []byte("$2a$15$GRvRO2nrpYTEuPQX6AieaOlZ4.7nMGsXpt.QWMev1zrP86JNspZb1")))

	assert.Nil(t, hasherx.Compare(context.Background(), []byte("test"), []byte("$argon2id$v=19$m=32,t=2,p=4$cm94YnRVOW5jZzFzcVE4bQ$MNzk5BtR2vUhrp6qQEjRNw")))
	assert.Nil(t, hasherx.CompareArgon2id(context.Background(), []byte("test"), []byte("$argon2id$v=19$m=32,t=2,p=4$cm94YnRVOW5jZzFzcVE4bQ$MNzk5BtR2vUhrp6qQEjRNw")))
	assert.Error(t, hasherx.Compare(context.Background(), []byte("test"), []byte("$argon2id$v=19$m=32,t=2,p=4$cm94YnRVOW5jZzFzcVE4bQ$MNzk5BtR2vUhrp6qQEjRN2")))

	assert.Nil(t, hasherx.Compare(context.Background(), []byte("test"), []byte("$argon2i$v=19$m=65536,t=3,p=4$kk51rW/vxIVCYn+EG4kTSg$NyT88uraJ6im6dyha/M5jhXvpqlEdlS/9fEm7ScMb8c")))
	assert.Nil(t, hasherx.CompareArgon2i(context.Background(), []byte("test"), []byte("$argon2i$v=19$m=65536,t=3,p=4$kk51rW/vxIVCYn+EG4kTSg$NyT88uraJ6im6dyha/M5jhXvpqlEdlS/9fEm7ScMb8c")))
	assert.Error(t, hasherx.Compare(context.Background(), []byte("test"), []byte("$argon2i$v=19$m=65536,t=3,p=4$pZ+27D6B0bCi0DwSmANF1w$4RNCUu4Uyu7eTIvzIdSuKz+I9idJlX/ykn6J10/W0EU")))

	assert.Nil(t, hasherx.Compare(context.Background(), []byte("test"), []byte("$argon2id$v=19$m=32,t=5,p=4$cm94YnRVOW5jZzFzcVE4bQ$fBxypOL0nP/zdPE71JtAV71i487LbX3fJI5PoTN6Lp4")))
	assert.Nil(t, hasherx.CompareArgon2id(context.Background(), []byte("test"), []byte("$argon2id$v=19$m=32,t=5,p=4$cm94YnRVOW5jZzFzcVE4bQ$fBxypOL0nP/zdPE71JtAV71i487LbX3fJI5PoTN6Lp4")))
	assert.Error(t, hasherx.Compare(context.Background(), []byte("test"), []byte("$argon2id$v=19$m=32,t=5,p=4$cm94YnRVOW5jZzFzcVE4bQ$fBxypOL0nP/zdPE71JtAV71i487LbX3fJI5PoTN6Lp5")))

	assert.Nil(t, hasherx.Compare(context.Background(), []byte("test"), []byte("$pbkdf2-sha256$i=100000,l=32$1jP+5Zxpxgtee/iPxGgOz0RfE9/KJuDElP1ley4VxXc$QJxzfvdbHYBpydCbHoFg3GJEqMFULwskiuqiJctoYpI")))
	assert.Nil(t, hasherx.ComparePbkdf2(context.Background(), []byte("test"), []byte("$pbkdf2-sha256$i=100000,l=32$1jP+5Zxpxgtee/iPxGgOz0RfE9/KJuDElP1ley4VxXc$QJxzfvdbHYBpydCbHoFg3GJEqMFULwskiuqiJctoYpI")))
	assert.Error(t, hasherx.Compare(context.Background(), []byte("test"), []byte("$pbkdf2-sha256$i=100000,l=32$1jP+5Zxpxgtee/iPxGgOz0RfE9/KJuDElP1ley4VxXc$QJxzfvdbHYBpydCbHoFg3GJEqMFULwskiuqiJctoYpp")))

	assert.Nil(t, hasherx.Compare(context.Background(), []byte("test"), []byte("$pbkdf2-sha512$i=100000,l=32$bdHBpn7OWOivJMVJypy2UqR0UnaD5prQXRZevj/05YU$+wArTfv1a+bNGO1iZrmEdVjhA+lL11wF4/IxpgYfPwc")))
	assert.Nil(t, hasherx.ComparePbkdf2(context.Background(), []byte("test"), []byte("$pbkdf2-sha512$i=100000,l=32$bdHBpn7OWOivJMVJypy2UqR0UnaD5prQXRZevj/05YU$+wArTfv1a+bNGO1iZrmEdVjhA+lL11wF4/IxpgYfPwc")))
	assert.Error(t, hasherx.Compare(context.Background(), []byte("test"), []byte("$pbkdf2-sha512$i=100000,l=32$bdHBpn7OWOivJMVJypy2UqR0UnaD5prQXRZevj/05YU$+wArTfv1a+bNGO1iZrmEdVjhA+lL11wF4/IxpgYfPww")))

	assert.Error(t, hasherx.Compare(context.Background(), []byte("test"), []byte("$pbkdf2-sha256$1jP+5Zxpxgtee/iPxGgOz0RfE9/KJuDElP1ley4VxXc$QJxzfvdbHYBpydCbHoFg3GJEqMFULwskiuqiJctoYpI")))
	assert.Error(t, hasherx.Compare(context.Background(), []byte("test"), []byte("$pbkdf2-sha256$aaaa$1jP+5Zxpxgtee/iPxGgOz0RfE9/KJuDElP1ley4VxXc$QJxzfvdbHYBpydCbHoFg3GJEqMFULwskiuqiJctoYpI")))
	assert.Error(t, hasherx.Compare(context.Background(), []byte("test"), []byte("$pbkdf2-sha256$i=100000,l=32$1jP+5Zxpxgtee/iPxGgOz0RfE9/KJuDElP1ley4VxXcc$QJxzfvdbHYBpydCbHoFg3GJEqMFULwskiuqiJctoYpI")))
	assert.Error(t, hasherx.Compare(context.Background(), []byte("test"), []byte("$pbkdf2-sha256$i=100000,l=32$1jP+5Zxpxgtee/iPxGgOz0RfE9/KJuDElP1ley4VxXc$QJxzfvdbHYBpydCbHoFg3GJEqMFULwskiuqiJctoYpII")))
	assert.Error(t, hasherx.Compare(context.Background(), []byte("test"), []byte("$pbkdf2-sha512$I=100000,l=32$bdHBpn7OWOivJMVJypy2UqR0UnaD5prQXRZevj/05YU$+wArTfv1a+bNGO1iZrmEdVjhA+lL11wF4/IxpgYfPwc")))
}
