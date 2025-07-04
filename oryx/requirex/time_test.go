// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package requirex

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type MockT struct {
	Failed bool
}

func (t *MockT) FailNow() {
	t.Failed = true
}

func (t *MockT) Errorf(format string, args ...interface{}) {
	_, _ = format, args
}

func TestEqualDurationAndTime(t *testing.T) {
	type args struct {
		expected  time.Duration
		actual    time.Duration
		precision time.Duration
	}
	tests := []struct {
		name string
		ok   bool
		args args
	}{
		{ok: true, name: "zero precision", args: args{expected: time.Nanosecond, actual: time.Nanosecond}},
		{ok: true, name: "small precision", args: args{expected: time.Nanosecond, actual: time.Nanosecond, precision: time.Nanosecond}},
		{ok: true, name: "large precision", args: args{expected: time.Nanosecond, actual: time.Nanosecond, precision: time.Hour}},
		{ok: false, name: "not within duration", args: args{expected: 12 * time.Second, actual: 13 * time.Second, precision: time.Nanosecond}},
		{ok: false, name: "not within duration negative value", args: args{expected: -12 * time.Second, actual: 13 * time.Second, precision: 20 * time.Second}},
		{ok: true, name: "within duration", args: args{expected: 12 * time.Second, actual: 13 * time.Second, precision: time.Second + time.Nanosecond}},
		{ok: true, name: "within duration negative value", args: args{expected: -12 * time.Second, actual: 13 * time.Second, precision: 30 * time.Second}},
		{ok: true, name: "exactly one precision apart", args: args{expected: 12 * time.Second, actual: 13 * time.Second, precision: time.Second}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("test equal duration", func(t *testing.T) {
				mt := MockT{}
				EqualDuration(&mt, tt.args.expected, tt.args.actual, tt.args.precision)
				require.Equal(t, !tt.ok, mt.Failed)

				mt = MockT{}
				EqualDuration(&mt, tt.args.actual, tt.args.expected, tt.args.precision)
				require.Equal(t, !tt.ok, mt.Failed)
			})

			t.Run("test equal time", func(t *testing.T) {
				rt := time.Now()
				mt := MockT{}
				EqualTime(&mt, rt.Add(tt.args.expected), rt.Add(tt.args.actual), tt.args.precision)
				require.Equal(t, !tt.ok, mt.Failed)

				mt = MockT{}
				EqualTime(&mt, rt.Add(tt.args.actual), rt.Add(tt.args.expected), tt.args.precision)
				require.Equal(t, !tt.ok, mt.Failed)

				rt = time.Time{}
				mt = MockT{}
				EqualTime(&mt, rt.Add(-tt.args.actual), rt.Add(-tt.args.expected), tt.args.precision)
				require.Equal(t, !tt.ok, mt.Failed)

			})
		})
	}
}
