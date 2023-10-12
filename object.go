package devtoolkit

import (
	"encoding/json"
	"reflect"
)

// ToPtr returns a pointer to the given value.
func ToPtr[T any](t T) *T {
	return &t
}

// IsZero returns true if the given value is the zero value for its type.
func IsZero(t any) bool {
	return t == nil || reflect.DeepEqual(t, reflect.Zero(reflect.TypeOf(t)).Interface())
}

// StructToMap converts a struct to a map[string]any.
func StructToMap(t any) (map[string]any, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	var mapResult map[string]interface{}
	err = json.Unmarshal(data, &mapResult)
	if err != nil {
		return nil, err
	}
	return mapResult, nil
}

// MapToStruct converts a map[string]any to a struct.
func MapToStruct[T any](m map[string]any) (*T, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	var structResult T
	err = json.Unmarshal(data, &structResult)
	if err != nil {
		return nil, err
	}
	return &structResult, nil
}

// CastToPointer casts the given value to a pointer of the given type.
// v must be a pointer.
// if v not a pointer, returns false.
// if v is nil, returns false.
// if v is a pointer but not of the given type, returns false.
// if v is a pointer of the given type, returns true.
func CastToPointer[T any](v any) (*T, bool) {
	vType := reflect.TypeOf(v)

	// v must not be nil
	if vType == nil {
		return nil, false
	}

	// v must be a pointer
	if vType.Kind() != reflect.Ptr {
		return nil, false
	}

	// cast v to *T
	resp, ok := v.(*T)
	if !ok {
		return nil, false
	}

	return resp, true
}
