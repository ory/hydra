// Copyright (c) 2013 Théo Crevon
//
// See the file LICENSE for copying permission.

/*
Package reflections provides high level abstractions above the
reflect library.

Reflect library is very low-level and as can be quite complex when it comes to do simple things like accessing a structure field value, a field tag...

The purpose of reflections package is to make developers life easier when it comes to introspect structures at runtime.
It's API is freely inspired from python language (getattr, setattr, hasattr...) and provides a simplified access to structure fields and tags.
*/
package reflections

import (
	"errors"
	"fmt"
	"reflect"
)

// GetField returns the value of the provided obj field. obj can whether
// be a structure or pointer to structure.
func GetField(obj interface{}, name string) (interface{}, error) {
	if !hasValidType(obj, []reflect.Kind{reflect.Struct, reflect.Ptr}) {
		return nil, errors.New("Cannot use GetField on a non-struct interface")
	}

	objValue := reflectValue(obj)
	field := objValue.FieldByName(name)
	if !field.IsValid() {
		return nil, fmt.Errorf("No such field: %s in obj", name)
	}

	return field.Interface(), nil
}

// GetFieldKind returns the kind of the provided obj field. obj can whether
// be a structure or pointer to structure.
func GetFieldKind(obj interface{}, name string) (reflect.Kind, error) {
	if !hasValidType(obj, []reflect.Kind{reflect.Struct, reflect.Ptr}) {
		return reflect.Invalid, errors.New("Cannot use GetField on a non-struct interface")
	}

	objValue := reflectValue(obj)
	field := objValue.FieldByName(name)

	if !field.IsValid() {
		return reflect.Invalid, fmt.Errorf("No such field: %s in obj", name)
	}

	return field.Type().Kind(), nil
}

// GetFieldTag returns the provided obj field tag value. obj can whether
// be a structure or pointer to structure.
func GetFieldTag(obj interface{}, fieldName, tagKey string) (string, error) {
	if !hasValidType(obj, []reflect.Kind{reflect.Struct, reflect.Ptr}) {
		return "", errors.New("Cannot use GetField on a non-struct interface")
	}

	objValue := reflectValue(obj)
	objType := objValue.Type()

	field, ok := objType.FieldByName(fieldName)
	if !ok {
		return "", fmt.Errorf("No such field: %s in obj", fieldName)
	}

	if !isExportableField(field) {
		return "", errors.New("Cannot GetFieldTag on a non-exported struct field")
	}

	return field.Tag.Get(tagKey), nil
}

// SetField sets the provided obj field with provided value. obj param has
// to be a pointer to a struct, otherwise it will soundly fail. Provided
// value type should match with the struct field you're trying to set.
func SetField(obj interface{}, name string, value interface{}) error {
	// Fetch the field reflect.Value
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	// If obj field value is not settable an error is thrown
	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		invalidTypeError := errors.New("Provided value type didn't match obj field type")
		return invalidTypeError
	}

	structFieldValue.Set(val)
	return nil
}

// HasField checks if the provided field name is part of a struct. obj can whether
// be a structure or pointer to structure.
func HasField(obj interface{}, name string) (bool, error) {
	if !hasValidType(obj, []reflect.Kind{reflect.Struct, reflect.Ptr}) {
		return false, errors.New("Cannot use GetField on a non-struct interface")
	}

	objValue := reflectValue(obj)
	objType := objValue.Type()
	field, ok := objType.FieldByName(name)
	if !ok || !isExportableField(field) {
		return false, nil
	}

	return true, nil
}

// Fields returns the struct fields names list. obj can whether
// be a structure or pointer to structure.
func Fields(obj interface{}) ([]string, error) {
	return fields(obj, false)
}

// FieldsDeep returns "flattened" fields (fields from anonymous
// inner structs are treated as normal fields)
func FieldsDeep(obj interface{}) ([]string, error) {
	return fields(obj, true)
}

func fields(obj interface{}, deep bool) ([]string, error) {
	if !hasValidType(obj, []reflect.Kind{reflect.Struct, reflect.Ptr}) {
		return nil, errors.New("Cannot use GetField on a non-struct interface")
	}

	objValue := reflectValue(obj)
	objType := objValue.Type()
	fieldsCount := objType.NumField()

	var allFields []string
	for i := 0; i < fieldsCount; i++ {
		field := objType.Field(i)
		if isExportableField(field) {
			if deep && field.Anonymous {
				fieldValue := objValue.Field(i)
				subFields, err := fields(fieldValue.Interface(), deep)
				if err != nil {
					return nil, fmt.Errorf("Cannot get fields in %s: %s", field.Name, err.Error())
				}
				allFields = append(allFields, subFields...)
			} else {
				allFields = append(allFields, field.Name)
			}
		}
	}

	return allFields, nil
}

// Items returns the field - value struct pairs as a map. obj can whether
// be a structure or pointer to structure.
func Items(obj interface{}) (map[string]interface{}, error) {
	return items(obj, false)
}

// FieldsDeep returns "flattened" items (fields from anonymous
// inner structs are treated as normal fields)
func ItemsDeep(obj interface{}) (map[string]interface{}, error) {
	return items(obj, true)
}

func items(obj interface{}, deep bool) (map[string]interface{}, error) {
	if !hasValidType(obj, []reflect.Kind{reflect.Struct, reflect.Ptr}) {
		return nil, errors.New("Cannot use GetField on a non-struct interface")
	}

	objValue := reflectValue(obj)
	objType := objValue.Type()
	fieldsCount := objType.NumField()

	allItems := make(map[string]interface{})

	for i := 0; i < fieldsCount; i++ {
		field := objType.Field(i)
		fieldValue := objValue.Field(i)
		if isExportableField(field) {
			if deep && field.Anonymous {
				if m, err := items(fieldValue.Interface(), deep); err == nil {
					for k, v := range m {
						allItems[k] = v
					}
				} else {
					return nil, fmt.Errorf("Cannot get items in %s: %s", field.Name, err.Error())
				}
			} else {
				allItems[field.Name] = fieldValue.Interface()
			}
		}
	}

	return allItems, nil
}

// Tags lists the struct tag fields. obj can whether
// be a structure or pointer to structure.
func Tags(obj interface{}, key string) (map[string]string, error) {
	return tags(obj, key, false)
}

// FieldsDeep returns "flattened" tags (fields from anonymous
// inner structs are treated as normal fields)
func TagsDeep(obj interface{}, key string) (map[string]string, error) {
	return tags(obj, key, true)
}

func tags(obj interface{}, key string, deep bool) (map[string]string, error) {
	if !hasValidType(obj, []reflect.Kind{reflect.Struct, reflect.Ptr}) {
		return nil, errors.New("Cannot use GetField on a non-struct interface")
	}

	objValue := reflectValue(obj)
	objType := objValue.Type()
	fieldsCount := objType.NumField()

	allTags := make(map[string]string)

	for i := 0; i < fieldsCount; i++ {
		structField := objType.Field(i)
		if isExportableField(structField) {
			if deep && structField.Anonymous {
				fieldValue := objValue.Field(i)
				if m, err := tags(fieldValue.Interface(), key, deep); err == nil {
					for k, v := range m {
						allTags[k] = v
					}
				} else {
					return nil, fmt.Errorf("Cannot get items in %s: %s", structField.Name, err.Error())
				}
			} else {
				allTags[structField.Name] = structField.Tag.Get(key)
			}
		}
	}

	return allTags, nil
}

func reflectValue(obj interface{}) reflect.Value {
	var val reflect.Value

	if reflect.TypeOf(obj).Kind() == reflect.Ptr {
		val = reflect.ValueOf(obj).Elem()
	} else {
		val = reflect.ValueOf(obj)
	}

	return val
}

func isExportableField(field reflect.StructField) bool {
	// PkgPath is empty for exported fields.
	return field.PkgPath == ""
}

func hasValidType(obj interface{}, types []reflect.Kind) bool {
	for _, t := range types {
		if reflect.TypeOf(obj).Kind() == t {
			return true
		}
	}

	return false
}

func isStruct(obj interface{}) bool {
	return reflect.TypeOf(obj).Kind() == reflect.Struct
}

func isPointer(obj interface{}) bool {
	return reflect.TypeOf(obj).Kind() == reflect.Ptr
}
