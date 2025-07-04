// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package stringsx

import (
	"fmt"
	"slices"
	"strings"
)

type (
	RegisteredCases struct {
		cases  []string
		actual string
	}
	errUnknownCase struct {
		*RegisteredCases
	}
	RegisteredPrefixes struct {
		prefixes []string
		actual   string
	}
	errUnknownPrefix struct {
		*RegisteredPrefixes
	}
)

var (
	ErrUnknownCase   = errUnknownCase{}
	ErrUnknownPrefix = errUnknownPrefix{}
)

func SwitchExact(actual string) *RegisteredCases {
	return &RegisteredCases{
		actual: actual,
	}
}

func SwitchPrefix(actual string) *RegisteredPrefixes {
	return &RegisteredPrefixes{
		actual: actual,
	}
}

func (r *RegisteredCases) AddCase(cases ...string) bool {
	r.cases = append(r.cases, cases...)
	return slices.Contains(cases, r.actual)
}

func (r *RegisteredPrefixes) HasPrefix(prefixes ...string) bool {
	r.prefixes = append(r.prefixes, prefixes...)
	return slices.ContainsFunc(prefixes, func(s string) bool {
		return strings.HasPrefix(r.actual, s)
	})
}

func (r *RegisteredCases) String() string {
	return "[" + strings.Join(r.cases, ", ") + "]"
}

func (r *RegisteredPrefixes) String() string {
	return "[" + strings.Join(r.prefixes, ", ") + "]"
}

func (r *RegisteredCases) ToUnknownCaseErr() error {
	return errUnknownCase{r}
}

func (r *RegisteredPrefixes) ToUnknownPrefixErr() error {
	return errUnknownPrefix{r}
}

func (e errUnknownCase) Error() string {
	return fmt.Sprintf("expected one of %s but got %s", e.String(), e.actual)
}

func (e errUnknownCase) Is(err error) bool {
	_, ok := err.(errUnknownCase)
	return ok
}

func (e errUnknownPrefix) Error() string {
	return fmt.Sprintf("expected %s to have one of the prefixes %s", e.actual, e.String())
}

func (e errUnknownPrefix) Is(err error) bool {
	_, ok := err.(errUnknownPrefix)
	return ok
}
