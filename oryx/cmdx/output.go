// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmdx

import "strconv"

type (
	// OutputIder outputs an ID
	OutputIder string
	// OutputIderCollection outputs a list of IDs
	OutputIderCollection struct {
		Items []OutputIder
	}
)

func (OutputIder) Header() []string {
	return []string{"ID"}
}

func (i OutputIder) Columns() []string {
	return []string{string(i)}
}

func (i OutputIder) Interface() interface{} {
	return i
}

func (OutputIderCollection) Header() []string {
	return []string{"ID"}
}

func (c OutputIderCollection) Table() [][]string {
	rows := make([][]string, len(c.Items))
	for i, ident := range c.Items {
		rows[i] = []string{string(ident)}
	}
	return rows
}

func (c OutputIderCollection) Interface() interface{} {
	return c.Items
}

func (c OutputIderCollection) Len() int {
	return len(c.Items)
}

type PaginatedList struct {
	Collection interface {
		Table
		IDs() []string
	} `json:"-"`
	Items         []interface{} `json:"items"`
	NextPageToken string        `json:"next_page_token"`
	IsLastPage    bool          `json:"is_last_page"`
}

func (r *PaginatedList) Header() []string {
	return r.Collection.Header()
}

func (r *PaginatedList) Table() [][]string {
	return append(
		r.Collection.Table(),
		[]string{},
		[]string{"NEXT PAGE TOKEN", r.NextPageToken},
		[]string{"IS LAST PAGE", strconv.FormatBool(r.IsLastPage)},
	)
}

func (r *PaginatedList) Interface() interface{} {
	return r
}

func (r *PaginatedList) Len() int {
	return r.Collection.Len() + 3
}

func (r *PaginatedList) IDs() []string {
	return r.Collection.IDs()
}

var _ Table = (*PaginatedList)(nil)
