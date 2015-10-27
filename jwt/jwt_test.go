package jwt

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestLoadCertificate(t *testing.T) {
	for _, c := range TestCertificates {
		out, err := LoadCertificate(c[0])
		assert.Nil(t, err)
		assert.Equal(t, c[1], string(out))
	}
	_, err := LoadCertificate("")
	assert.NotNil(t, err)
	_, err = LoadCertificate("foobar")
	assert.NotNil(t, err)
}

func TestSignRejectsAlgAndTypHeader(t *testing.T) {
	j := New([]byte(TestCertificates[0][1]), []byte(TestCertificates[1][1]))
	for _, c := range []map[string]interface{}{
		map[string]interface{}{"alg": "foo"},
		map[string]interface{}{"typ": "foo"},
		map[string]interface{}{"typ": "foo", "alg": "foo"},
	} {
		_, err := j.SignToken(map[string]interface{}{}, c)
		assert.NotNil(t, err)
	}
}

func TestSignAndVerify(t *testing.T) {
	type test struct {
		private []byte
		public  []byte
		header  map[string]interface{}
		claims  map[string]interface{}
		valid   bool
		signOk  bool
	}

	cases := []test{
		test{
			[]byte(""),
			[]byte(TestCertificates[1][1]),
			map[string]interface{}{"foo": "bar"},
			map[string]interface{}{"nbf": time.Now().Add(time.Hour).Unix()},
			false,
			false,
		},
		test{
			[]byte(TestCertificates[0][1]),
			[]byte(""),
			map[string]interface{}{"foo": "bar"},
			map[string]interface{}{"nbf": time.Now().Add(time.Hour).Unix()},
			false,
			true,
		},
		test{
			[]byte(TestCertificates[0][1]),
			[]byte(TestCertificates[1][1]),
			map[string]interface{}{"foo": "bar"},
			map[string]interface{}{"nbf": time.Now().Add(time.Hour).Unix()},
			false,
			true,
		},
		test{
			[]byte(TestCertificates[0][1]),
			[]byte(TestCertificates[1][1]),
			map[string]interface{}{"foo": "bar"},
			map[string]interface{}{"exp": time.Now().Add(-time.Hour).Unix()},
			false,
			true,
		},
		test{
			[]byte(TestCertificates[0][1]),
			[]byte(TestCertificates[1][1]),
			map[string]interface{}{"foo": "bar"},
			map[string]interface{}{
				"nbf": time.Now().Add(-time.Hour).Unix(),
				"exp": time.Now().Add(time.Hour).Unix(),
			},
			true,
			true,
		},
		test{
			[]byte(TestCertificates[0][1]),
			[]byte(TestCertificates[1][1]),
			map[string]interface{}{"foo": "bar"},
			map[string]interface{}{
				"nbf": time.Now().Add(-time.Hour).Unix(),
			},
			true,
			true,
		},
		test{
			[]byte(TestCertificates[0][1]),
			[]byte(TestCertificates[1][1]),
			map[string]interface{}{"foo": "bar"},
			map[string]interface{}{
				"exp": time.Now().Add(time.Hour).Unix(),
			},
			true,
			true,
		},
		test{
			[]byte(TestCertificates[0][1]),
			[]byte(TestCertificates[1][1]),
			map[string]interface{}{"foo": "bar"},
			map[string]interface{}{},
			true,
			true,
		},
	}

	for i, c := range cases {
		j := New(c.private, c.public)
		data, err := j.SignToken(c.claims, c.header)
		if c.signOk {
			require.Nil(t, err, "Case %d", i)
		} else {
			require.NotNil(t, err, "Case %d", i)
		}
		tok, err := j.VerifyToken([]byte(data))
		if c.valid {
			require.Nil(t, err, "Case %d", i)
			require.Equal(t, c.valid, tok.Valid, "Case %d", i)
		} else {
			require.NotNil(t, err, "Case %d", i)
		}
	}
}
