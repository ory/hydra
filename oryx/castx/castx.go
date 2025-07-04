// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package castx

import (
	"encoding/csv"
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/cast"
)

// ToFloatSlice casts an interface to a []float64 type.
func ToFloatSlice(i interface{}) []float64 {
	f, _ := ToFloatSliceE(i)
	return f
}

// ToFloatSliceE casts an interface to a []float64 type.
func ToFloatSliceE(i interface{}) ([]float64, error) {
	if i == nil {
		return []float64{}, fmt.Errorf("unable to cast %#v of type %T to []float64", i, i)
	}

	switch v := i.(type) {
	case []float64:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]float64, s.Len())
		for j := range a {
			val, err := cast.ToFloat64E(s.Index(j).Interface())
			if err != nil {
				return []float64{}, fmt.Errorf("unable to cast %#v of type %T to []float64", i, i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return []float64{}, fmt.Errorf("unable to cast %#v of type %T to []float64", i, i)
	}
}

// ToStringSlice casts an interface to a []string type and respects comma-separated values.
func ToStringSlice(i interface{}) []string {
	s, _ := ToStringSliceE(i)
	return s
}

// ToStringSliceE casts an interface to a []string type and respects comma-separated values.
func ToStringSliceE(i interface{}) ([]string, error) {
	switch s := i.(type) {
	case string:
		return parseCSV(s)
	}

	return cast.ToStringSliceE(i)
}

func parseCSV(v string) ([]string, error) {
	return csv.NewReader(strings.NewReader(v)).Read()
}
