// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jsonschemax

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/ory/jsonschema/v3"
)

type (
	byName       []Path
	PathEnhancer interface {
		EnhancePath(Path) map[string]interface{}
	}
	TypeHint int
)

func (s byName) Len() int           { return len(s) }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byName) Less(i, j int) bool { return s[i].Name < s[j].Name }

const (
	String TypeHint = iota + 1
	Float
	Int
	Bool
	JSON
	Nil

	BoolSlice
	StringSlice
	IntSlice
	FloatSlice
)

// Path represents a JSON Schema Path.
type Path struct {
	// Title of the path.
	Title string

	// Description of the path.
	Description string

	// Examples of the path.
	Examples []interface{}

	// Name is the JSON path name.
	Name string

	// Default is the default value of that path.
	Default interface{}

	// Type is a prototype (e.g. float64(0)) of the path type.
	Type interface{}

	TypeHint

	// Format is the format of the path if defined
	Format string

	// Pattern is the pattern of the path if defined
	Pattern *regexp.Regexp

	// Enum are the allowed enum values
	Enum []interface{}

	// first element in slice is constant value. note: slice is used to capture nil constant.
	Constant []interface{}

	// ReadOnly is whether the value is readonly
	ReadOnly bool

	// -1 if not specified
	MinLength int
	MaxLength int

	// Required if set indicates this field is required.
	Required bool

	Minimum *big.Float
	Maximum *big.Float

	MultipleOf *big.Float

	CustomProperties map[string]interface{}
}

// ListPathsBytes works like ListPathsWithRecursion but prepares the JSON Schema itself.
func ListPathsBytes(ctx context.Context, raw json.RawMessage, maxRecursion int16) ([]Path, error) {
	compiler := jsonschema.NewCompiler()
	compiler.ExtractAnnotations = true
	id := fmt.Sprintf("%x.json", sha256.Sum256(raw))
	if err := compiler.AddResource(id, bytes.NewReader(raw)); err != nil {
		return nil, err
	}
	compiler.ExtractAnnotations = true
	return runPathsFromCompiler(ctx, id, compiler, maxRecursion, false)
}

// ListPathsWithRecursion will follow circular references until maxRecursion is reached, without
// returning an error.
func ListPathsWithRecursion(ctx context.Context, ref string, compiler *jsonschema.Compiler, maxRecursion uint8) ([]Path, error) {
	return runPathsFromCompiler(ctx, ref, compiler, int16(maxRecursion), false)
}

// ListPaths lists all paths of a JSON Schema. Will return an error
// if circular references are found.
func ListPaths(ctx context.Context, ref string, compiler *jsonschema.Compiler) ([]Path, error) {
	return runPathsFromCompiler(ctx, ref, compiler, -1, false)
}

// ListPathsWithArraysIncluded lists all paths of a JSON Schema. Will return an error
// if circular references are found.
// Includes arrays with `#`.
func ListPathsWithArraysIncluded(ctx context.Context, ref string, compiler *jsonschema.Compiler) ([]Path, error) {
	return runPathsFromCompiler(ctx, ref, compiler, -1, true)
}

// ListPathsWithInitializedSchema loads the paths from the schema without compiling it.
//
// You MUST ensure that the compiler was using `ExtractAnnotations = true`.
func ListPathsWithInitializedSchema(schema *jsonschema.Schema) ([]Path, error) {
	return runPaths(schema, -1, false)
}

// ListPathsWithInitializedSchemaAndArraysIncluded loads the paths from the schema without compiling it.
//
// You MUST ensure that the compiler was using `ExtractAnnotations = true`.
// Includes arrays with `#`.
func ListPathsWithInitializedSchemaAndArraysIncluded(schema *jsonschema.Schema) ([]Path, error) {
	return runPaths(schema, -1, true)
}

func runPathsFromCompiler(ctx context.Context, ref string, compiler *jsonschema.Compiler, maxRecursion int16, includeArrays bool) ([]Path, error) {
	if compiler == nil {
		compiler = jsonschema.NewCompiler()
	}

	compiler.ExtractAnnotations = true

	schema, err := compiler.Compile(ctx, ref)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return runPaths(schema, maxRecursion, includeArrays)
}

func runPaths(schema *jsonschema.Schema, maxRecursion int16, includeArrays bool) ([]Path, error) {
	pointers := map[string]bool{}
	paths, err := listPaths(schema, nil, nil, pointers, 0, maxRecursion, includeArrays)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sort.Stable(paths)
	return makeUnique(paths)
}

func makeUnique(in byName) (byName, error) {
	cache := make(map[string]Path)
	for _, p := range in {
		vc, ok := cache[p.Name]
		if !ok {
			cache[p.Name] = p
			continue
		}

		if fmt.Sprintf("%T", p.Type) != fmt.Sprintf("%T", p.Type) {
			return nil, errors.Errorf("multiple types %+v are not supported for path: %s", []interface{}{p.Type, vc.Type}, p.Name)
		}

		if vc.Default == nil {
			cache[p.Name] = p
		}
	}

	k := 0
	out := make([]Path, len(cache))
	for _, v := range cache {
		out[k] = v
		k++
	}

	paths := byName(out)
	sort.Sort(paths)
	return paths, nil
}

func appendPointer(in map[string]bool, pointer *jsonschema.Schema) map[string]bool {
	out := make(map[string]bool)
	for k, v := range in {
		out[k] = v
	}
	out[fmt.Sprintf("%p", pointer)] = true
	return out
}

func listPaths(schema *jsonschema.Schema, parent *jsonschema.Schema, parents []string, pointers map[string]bool, currentRecursion int16, maxRecursion int16, includeArrays bool) (byName, error) {
	var pathType interface{}
	var pathTypeHint TypeHint
	var paths []Path
	_, isCircular := pointers[fmt.Sprintf("%p", schema)]

	if len(schema.Constant) > 0 {
		switch schema.Constant[0].(type) {
		case float64, json.Number:
			pathType = float64(0)
			pathTypeHint = Float
		case int8, int16, int, int64:
			pathType = int64(0)
			pathTypeHint = Int
		case string:
			pathType = ""
			pathTypeHint = String
		case bool:
			pathType = false
			pathTypeHint = Bool
		default:
			pathType = schema.Constant[0]
			pathTypeHint = JSON
		}
	} else if len(schema.Types) == 1 {
		switch schema.Types[0] {
		case "null":
			pathType = nil
			pathTypeHint = Nil
		case "boolean":
			pathType = false
			pathTypeHint = Bool
		case "number":
			pathType = float64(0)
			pathTypeHint = Float
		case "integer":
			pathType = float64(0)
			pathTypeHint = Int
		case "string":
			pathType = ""
			pathTypeHint = String
		case "array":
			pathType = []interface{}{}
			if schema.Items != nil {
				var itemSchemas []*jsonschema.Schema
				switch t := schema.Items.(type) {
				case []*jsonschema.Schema:
					itemSchemas = t
				case *jsonschema.Schema:
					itemSchemas = []*jsonschema.Schema{t}
				}
				var types []string
				for _, is := range itemSchemas {
					types = append(types, is.Types...)
					if is.Ref != nil {
						types = append(types, is.Ref.Types...)
					}
				}
				types = slices.Compact(types)
				if len(types) == 1 {
					switch types[0] {
					case "boolean":
						pathType = []bool{}
						pathTypeHint = BoolSlice
					case "number":
						pathType = []float64{}
						pathTypeHint = FloatSlice
					case "integer":
						pathType = []float64{}
						pathTypeHint = IntSlice
					case "string":
						pathType = []string{}
						pathTypeHint = StringSlice
					default:
						pathType = []interface{}{}
						pathTypeHint = JSON
					}
				}
			}
		case "object":
			pathType = map[string]interface{}{}
			pathTypeHint = JSON
		}
	} else if len(schema.Types) > 2 {
		pathType = nil
		pathTypeHint = JSON
	}

	var def interface{} = schema.Default
	if v, ok := def.(json.Number); ok {
		def, _ = v.Float64()
	}

	if (pathType != nil || schema.Default != nil) && len(parents) > 0 {
		name := parents[len(parents)-1]
		var required bool
		if parent != nil {
			for _, r := range parent.Required {
				if r == name {
					required = true
					break
				}
			}
		}

		path := Path{
			Name:        strings.Join(parents, "."),
			Default:     def,
			Type:        pathType,
			TypeHint:    pathTypeHint,
			Format:      schema.Format,
			Pattern:     schema.Pattern,
			Enum:        schema.Enum,
			Constant:    schema.Constant,
			MinLength:   schema.MinLength,
			MaxLength:   schema.MaxLength,
			Minimum:     schema.Minimum,
			Maximum:     schema.Maximum,
			MultipleOf:  schema.MultipleOf,
			ReadOnly:    schema.ReadOnly,
			Title:       schema.Title,
			Description: schema.Description,
			Examples:    schema.Examples,
			Required:    required,
		}

		for _, e := range schema.Extensions {
			if enhancer, ok := e.(PathEnhancer); ok {
				path.CustomProperties = enhancer.EnhancePath(path)
			}
		}
		paths = append(paths, path)
	}

	if isCircular {
		if maxRecursion == -1 {
			return nil, errors.Errorf("detected circular dependency in schema path: %s", strings.Join(parents, "."))
		} else if currentRecursion > maxRecursion {
			return paths, nil
		}
		currentRecursion++
	}

	if schema.Ref != nil {
		path, err := listPaths(schema.Ref, schema, parents, appendPointer(pointers, schema), currentRecursion, maxRecursion, includeArrays)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path...)
	}

	if schema.Not != nil {
		path, err := listPaths(schema.Not, schema, parents, appendPointer(pointers, schema), currentRecursion, maxRecursion, includeArrays)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path...)
	}

	if schema.If != nil {
		path, err := listPaths(schema.If, schema, parents, appendPointer(pointers, schema), currentRecursion, maxRecursion, includeArrays)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path...)
	}

	if schema.Then != nil {
		path, err := listPaths(schema.Then, schema, parents, appendPointer(pointers, schema), currentRecursion, maxRecursion, includeArrays)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path...)
	}

	if schema.Else != nil {
		path, err := listPaths(schema.Else, schema, parents, appendPointer(pointers, schema), currentRecursion, maxRecursion, includeArrays)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path...)
	}

	for _, sub := range schema.AllOf {
		path, err := listPaths(sub, schema, parents, appendPointer(pointers, schema), currentRecursion, maxRecursion, includeArrays)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path...)
	}

	for _, sub := range schema.AnyOf {
		path, err := listPaths(sub, schema, parents, appendPointer(pointers, schema), currentRecursion, maxRecursion, includeArrays)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path...)
	}

	for _, sub := range schema.OneOf {
		path, err := listPaths(sub, schema, parents, appendPointer(pointers, schema), currentRecursion, maxRecursion, includeArrays)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path...)
	}

	for name, sub := range schema.Properties {
		path, err := listPaths(sub, schema, append(parents, name), appendPointer(pointers, schema), currentRecursion, maxRecursion, includeArrays)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path...)
	}

	if schema.Items != nil && includeArrays {
		switch t := schema.Items.(type) {
		case []*jsonschema.Schema:
			for _, sub := range t {
				path, err := listPaths(sub, schema, append(parents, "#"), appendPointer(pointers, schema), currentRecursion, maxRecursion, includeArrays)
				if err != nil {
					return nil, err
				}
				paths = append(paths, path...)
			}
		case *jsonschema.Schema:
			path, err := listPaths(t, schema, append(parents, "#"), appendPointer(pointers, schema), currentRecursion, maxRecursion, includeArrays)
			if err != nil {
				return nil, err
			}
			paths = append(paths, path...)
		}
	}

	return paths, nil
}
