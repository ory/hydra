// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package keysetpagination

import (
	"cmp"
	"reflect"

	"github.com/jmoiron/sqlx/reflectx"
)

type (
	Paginator struct {
		token, defaultToken        PageToken
		size, defaultSize, maxSize int
		isLast                     bool
	}
	Option func(*Paginator)
)

const (
	DefaultSize    = 100
	DefaultMaxSize = 500
)

var dbStructTagMapper = reflectx.NewMapper("db")

func (p *Paginator) DefaultToken() PageToken { return p.defaultToken }
func (p *Paginator) IsLast() bool            { return p.isLast }

func (p *Paginator) PageToken() PageToken {
	if p.token.cols != nil {
		return p.token
	}
	return p.defaultToken
}

func (p *Paginator) Size() int {
	defaultSize := cmp.Or(p.defaultSize, DefaultSize)
	maxSize := cmp.Or(p.maxSize, DefaultMaxSize)

	size := p.size
	if size <= 0 {
		size = defaultSize
	}
	if size > maxSize {
		size = maxSize
	}

	return size
}

func (p *Paginator) ToOptions() []Option {
	opts := make([]Option, 0, 6)
	if p.token.cols != nil {
		opts = append(opts, WithToken(p.token))
	}
	if p.defaultToken.cols != nil {
		opts = append(opts, WithDefaultToken(p.defaultToken))
	}
	if p.size > 0 {
		opts = append(opts, WithSize(p.size))
	}
	if p.defaultSize != DefaultSize {
		opts = append(opts, WithDefaultSize(p.defaultSize))
	}
	if p.maxSize != DefaultMaxSize {
		opts = append(opts, WithMaxSize(p.maxSize))
	}
	if p.isLast {
		opts = append(opts, withIsLast(p.isLast))
	}
	return opts
}

// Result removes the last item (if applicable) and returns the paginator for the next page.
func Result[I any](items []I, p *Paginator) ([]I, *Paginator) {
	return ResultFunc(items, p, func(last I, colName string) any {
		lastItemVal := reflect.ValueOf(last)
		return dbStructTagMapper.FieldByName(lastItemVal, colName).Interface()
	})
}

// ResultFunc removes the last item (if applicable) and returns the paginator for the next page.
// The extractor function is used to extract the column values from the last item.
func ResultFunc[I any](items []I, p *Paginator, extractor func(last I, colName string) any) ([]I, *Paginator) {
	if len(items) <= p.Size() {
		return items, &Paginator{
			isLast: true,

			defaultToken: p.defaultToken,
			size:         p.size,
			defaultSize:  p.defaultSize,
			maxSize:      p.maxSize,
		}
	}

	items = items[:p.Size()]
	lastItem := items[len(items)-1]

	currentCols := p.PageToken().Columns()
	newCols := make([]Column, len(currentCols))
	for i, col := range currentCols {
		newCols[i] = Column{
			Name:  col.Name,
			Order: col.Order,
			Value: extractor(lastItem, col.Name),
		}
	}

	return items, &Paginator{
		token:        NewPageToken(newCols...),
		defaultToken: p.defaultToken,
		size:         p.size,
		defaultSize:  p.defaultSize,
		maxSize:      p.maxSize,
	}
}

func WithSize(size int) Option {
	return func(p *Paginator) { p.size = size }
}
func WithDefaultSize(size int) Option {
	return func(p *Paginator) { p.defaultSize = size }
}
func WithMaxSize(size int) Option {
	return func(p *Paginator) { p.maxSize = size }
}
func WithToken(t PageToken) Option {
	return func(p *Paginator) { p.token = t }
}
func WithDefaultToken(t PageToken) Option {
	return func(p *Paginator) { p.defaultToken = t }
}
func withIsLast(isLast bool) Option {
	return func(p *Paginator) { p.isLast = isLast }
}

func NewPaginator(modifiers ...Option) *Paginator {
	p := &Paginator{
		// these can still be overridden by the modifiers, but they should never be unset
		maxSize:     DefaultMaxSize,
		defaultSize: DefaultSize,
	}
	for _, f := range modifiers {
		f(p)
	}
	return p
}
