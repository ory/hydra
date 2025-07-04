// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package pointerx

// Ptr returns the input value's pointer.
func Ptr[T any](v T) *T {
	return &v
}

// Deref returns the input values de-referenced value, or zero value if nil.
func Deref[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

// String returns the input value's pointer.
// Deprecated: use Ptr instead.
func String(s string) *string {
	return &s
}

// StringR is the reverse to String.
// Deprecated: use Deref instead.
func StringR(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Int returns the input value's pointer.
// Deprecated: use Ptr instead.
func Int(s int) *int {
	return &s
}

// IntR is the reverse to Int.
// Deprecated: use Deref instead.
func IntR(s *int) int {
	if s == nil {
		return int(0)
	}
	return *s
}

// Int32 returns the input value's pointer.
// Deprecated: use Ptr instead.
func Int32(s int32) *int32 {
	return &s
}

// Int32R is the reverse to Int32.
// Deprecated: use Deref instead.
func Int32R(s *int32) int32 {
	if s == nil {
		return int32(0)
	}
	return *s
}

// Int64 returns the input value's pointer.
// Deprecated: use Ptr instead.
func Int64(s int64) *int64 {
	return &s
}

// Int64R is the reverse to Int64.
// Deprecated: use Deref instead.
func Int64R(s *int64) int64 {
	if s == nil {
		return int64(0)
	}
	return *s
}

// Float32 returns the input value's pointer.
// Deprecated: use Ptr instead.
func Float32(s float32) *float32 {
	return &s
}

// Float32R is the reverse to Float32.
// Deprecated: use Deref instead.
func Float32R(s *float32) float32 {
	if s == nil {
		return float32(0)
	}
	return *s
}

// Float64 returns the input value's pointer.
// Deprecated: use Ptr instead.
func Float64(s float64) *float64 {
	return &s
}

// Float64R is the reverse to Float64.
// Deprecated: use Deref instead.
func Float64R(s *float64) float64 {
	if s == nil {
		return float64(0)
	}
	return *s
}

// Bool returns the input value's pointer.
// Deprecated: use Ptr instead.
func Bool(s bool) *bool {
	return &s
}

// BoolR is the reverse to Bool.
// Deprecated: use Deref instead.
func BoolR(s *bool) bool {
	if s == nil {
		return false
	}
	return *s
}
