// Copyright © 2026 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package paginationplanner

import (
	"github.com/pkg/errors"

	keysetpagination "github.com/ory/x/pagination/keysetpagination_v2"
)

// PaginationPlanner chooses the best [PaginationPlan] for a queriedColumns.
// Plan is eligible when its required constraint set is satisfied (and an optional Condition matches).
// If no plan matches, the FallbackPlan is used.
type PaginationPlanner struct {
	Plans        []PaginationPlan
	FallbackPlan PaginationPlan
}

func NewPaginationPlanner(fallbackPlan PaginationPlan, plans []PaginationPlan) (*PaginationPlanner, error) {
	if len(fallbackPlan.DefaultPageToken.Columns()) == 0 {
		return nil, errors.New("plan must define at least one PageTokenColumn")
	}

	for i := range plans {
		plan := &plans[i]
		if len(plan.ApplicableQueries) == 0 {
			return nil, errors.New("plan must define at least one ApplicableQueries")
		}
		if len(plan.DefaultPageToken.Columns()) == 0 {
			return nil, errors.New("plan must define a DefaultPageToken")
		}
		plan.populateInternals()
	}

	return &PaginationPlanner{
		Plans:        plans,
		FallbackPlan: fallbackPlan,
	}, nil
}

// GetPaginator selects the first eligible plan for the given queriedColumns constraints.
// Eligibility requires an exact ColumnSet match and (if set) Condition match.
// If none match, the FallbackPlan is used for building the Paginator.
func (p *PaginationPlanner) GetPaginator(q Query, pageOpts ...keysetpagination.Option) *keysetpagination.Paginator {
	if len(q) == 0 {
		return keysetpagination.NewPaginator(append(pageOpts, keysetpagination.WithDefaultToken(p.FallbackPlan.DefaultPageToken))...)
	}

	plan := p.pickPlan(q)

	origCols := plan.DefaultPageToken.Columns()

	// Make a fresh copy so we don't mutate the original backing array.
	cols := make([]keysetpagination.Column, len(origCols))
	copy(cols, origCols)

	for i, col := range cols {
		if val, ok := q.colByName(col.Name); ok && val.IsConstrained() {
			cols[i].HasConstraint = true
		}
	}

	defaultToken := keysetpagination.WithDefaultToken(keysetpagination.NewPageToken(cols...))
	return keysetpagination.NewPaginator(append(pageOpts, defaultToken)...)
}

func (p *PaginationPlanner) pickPlan(q Query) PaginationPlan {
	constrainedCols := q.constrainedCols()

	for _, plan := range p.Plans {
		if _, ok := plan.applicableQueries[constrainedCols]; !ok {
			continue
		}
		if plan.Condition != nil && !plan.Condition(q) {
			continue
		}
		return plan
	}

	return p.FallbackPlan
}

type PaginationPlan struct {
	Name string

	// DefaultPageToken defines the columns that the pagination will be made of
	DefaultPageToken keysetpagination.PageToken

	// ApplicableQueries defines the exact sets of columns that this pagination plan is applicable for.
	ApplicableQueries [][]Column
	applicableQueries map[uint]struct{}

	// Condition further restricts the Plan eligibility for given queriedColumns.
	// This can be used for plan-specific logic, e.g. verifying partial-index predicates match.
	Condition func(q Query) bool
}

func (pp *PaginationPlan) GetPaginator(pageOpts ...keysetpagination.Option) *keysetpagination.Paginator {
	defaultToken := keysetpagination.WithDefaultToken(pp.DefaultPageToken)
	return keysetpagination.NewPaginator(append(pageOpts, defaultToken)...)
}

func (pp *PaginationPlan) populateInternals() {
	pp.applicableQueries = make(map[uint]struct{}, len(pp.ApplicableQueries))
	for _, cols := range pp.ApplicableQueries {
		var colSet uint
		for _, col := range cols {
			colSet |= col.bit
		}
		pp.applicableQueries[colSet] = struct{}{}
	}
}

type Table struct {
	nextBit   uint
	usedNames map[string]struct{}
}

func NewTable() *Table {
	return &Table{
		nextBit:   1,
		usedNames: make(map[string]struct{}),
	}
}

func (t *Table) NewColumn(name string) Column {
	defer func() { t.nextBit = t.nextBit << 1 }()
	if _, exists := t.usedNames[name]; exists {
		panic("column name already used: " + name)
	}
	t.usedNames[name] = struct{}{}
	return Column{
		name: name,
		bit:  t.nextBit,
	}
}

type Column struct {
	bit  uint
	name string
}

func (c Column) Name() string {
	return c.name
}

type ColumnConstraint uint8

const (
	colUnconstrained    ColumnConstraint = iota
	colConstraintEq                      // col = ?
	colConstraintIsNull                  // col IS NULL
)

func (cs ColumnConstraint) IsConstrained() bool {
	return cs == colConstraintEq || cs == colConstraintIsNull
}

func (cs ColumnConstraint) IsNull() bool {
	return cs == colConstraintIsNull
}

// Query is used to track which column constraints are set.
type Query map[Column]ColumnConstraint

func NewQuery() Query { return make(Query) }

func (qc Query) SetEq(cols ...Column) Query {
	return qc.set(colConstraintEq, cols...)
}

func (qc Query) SetIsNull(cols ...Column) Query {
	return qc.set(colConstraintIsNull, cols...)
}

func (qc Query) set(c ColumnConstraint, cols ...Column) Query {
	for _, col := range cols {
		qc[col] = c
	}
	return qc
}

func (qc Query) colByName(name string) (ColumnConstraint, bool) {
	for col, constraint := range qc {
		if col.name == name {
			return constraint, true
		}
	}
	return colUnconstrained, false
}

func (qc Query) constrainedCols() uint {
	var cols uint
	for col, state := range qc {
		if !state.IsConstrained() {
			continue
		}
		cols |= col.bit
	}
	return cols
}
