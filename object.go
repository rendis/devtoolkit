package devtoolkit

import "reflect"

// ToPtr returns a pointer to the given value.
func ToPtr[T any](t T) *T {
	return &t
}

// IsZero returns true if the given value is the zero value for its type.
func IsZero(t any) bool {
	return t == nil || reflect.DeepEqual(t, reflect.Zero(reflect.TypeOf(t)).Interface())
}
