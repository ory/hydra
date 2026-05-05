// Package region defines the Ory Network region enumeration. Lives in x/ so
// kratos-oss and other OSS layers can use it without importing cloudlib.
package region

import (
	"database/sql/driver"

	"github.com/ory/herodot"
	"github.com/pkg/errors"
)

// Region is an Ory Network region. Specific regions map to a single CRDB
// region; super-regions (EU, Asia, US, Global) group specific regions that
// share a data-residency boundary.
//
// swagger:enum Region
type Region string

const (
	EUCentral     Region = "eu-central"
	AsiaNorthEast Region = "asia-northeast"
	USEast        Region = "us-east"
	USWest        Region = "us-west"

	EU     Region = "eu"
	Asia   Region = "asia"
	US     Region = "us"
	Global Region = "global"
)

// All returns every valid region in stable order: specific regions first,
// then super-regions.
func All() []Region {
	return []Region{EUCentral, AsiaNorthEast, USEast, USWest, EU, Asia, US, Global}
}

func (r Region) String() string {
	return string(r)
}

// Valid reports whether r is a known region.
func (r Region) Valid() bool {
	switch r {
	case EUCentral, AsiaNorthEast, USEast, USWest, EU, Asia, US, Global:
		return true
	}
	return false
}

// IsSuperRegion reports whether r is a super-region.
func (r Region) IsSuperRegion() bool {
	switch r {
	case EU, Asia, US, Global:
		return true
	}
	return false
}

// Contains reports whether r contains other:
//   - Global contains every valid region.
//   - EU contains EUCentral; Asia contains AsiaNorthEast; US contains USEast and USWest.
//   - Every region contains itself.
func (r Region) Contains(other Region) bool {
	if !r.Valid() || !other.Valid() {
		return false
	}
	if r == Global || r == other {
		return true
	}
	switch r {
	case EU:
		return other == EUCentral
	case Asia:
		return other == AsiaNorthEast
	case US:
		return other == USEast || other == USWest
	}
	return false
}

// Scan implements sql.Scanner. NULL and empty scan to the zero value;
// validate via Region.Valid if "unset" must be rejected.
func (r *Region) Scan(src any) error {
	switch s := src.(type) {
	case nil:
		*r = ""
	case string:
		*r = Region(s)
	case []byte:
		*r = Region(s)
	default:
		return errors.Errorf("cannot scan %T into region.Region", src)
	}
	return nil
}

// Value implements driver.Valuer. The empty Region writes as "".
func (r Region) Value() (driver.Value, error) {
	return string(r), nil
}

// IsEqual compares two nullable *Region pointers (both nil = equal).
func IsEqual(a, b *Region) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// ErrInvalid is wrapped when a submitted region is not a known value.
// Use NewErrInvalid for a herodot 400; use errors.Is for chain checks.
var ErrInvalid = errors.New("the provided region is not a valid Ory region")

// ErrNotAllowed is wrapped when a valid region is outside the project's
// home_region constraint.
var ErrNotAllowed = errors.New("the provided region is not allowed by this project's home region")

// NewErrInvalid returns a fresh herodot 400 wrapping ErrInvalid.
func NewErrInvalid() error {
	return errors.WithStack(
		herodot.ErrBadRequest().
			WithReason(ErrInvalid.Error()).
			WithDebug(`region must be one of eu-central, asia-northeast, us-east, us-west, eu, asia, us, global`).
			WithWrap(ErrInvalid),
	)
}

// NewErrNotAllowed returns a fresh herodot 400 wrapping ErrNotAllowed.
func NewErrNotAllowed() error {
	return errors.WithStack(
		herodot.ErrBadRequest().
			WithReason(ErrNotAllowed.Error()).
			WithWrap(ErrNotAllowed),
	)
}
