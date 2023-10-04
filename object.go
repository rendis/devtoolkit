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
