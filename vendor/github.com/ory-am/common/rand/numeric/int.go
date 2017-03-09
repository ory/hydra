package numeric

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"io"
)

var (
	rander = rand.Reader // random function
	r      = make([]byte, 8)
)

// Int64 creates a random 64 bit integer using crypto.rand
func Int64() (i int64) {
	randomBits(r)
	buf := bytes.NewBuffer(r)
	binary.Read(buf, binary.LittleEndian, &i)
	return i
}

// UInt64 creates a random 64 bit unsigned integer using crypto.rand
func UInt64() (i uint64) {
	randomBits(r)
	buf := bytes.NewBuffer(r)
	binary.Read(buf, binary.LittleEndian, &i)
	return i
}

// Int32 creates a random 32 bit integer using crypto.rand
func Int32() (i int32) {
	randomBits(r)
	buf := bytes.NewBuffer(r)
	binary.Read(buf, binary.LittleEndian, &i)
	return i
}

// UInt32 creates a random 32 bit unsigned integer using crypto.rand
func UInt32() (i uint32) {
	randomBits(r)
	buf := bytes.NewBuffer(r)
	binary.Read(buf, binary.LittleEndian, &i)
	return i
}

// randomBits completely fills slice b with random data.
func randomBits(b []byte) {
	if _, err := io.ReadFull(rander, b); err != nil {
		panic(err.Error()) // rand should never fail
	}
}
