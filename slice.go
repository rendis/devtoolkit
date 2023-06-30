package devtoolkit

// Contains checks if a slice contains an item. Item must be comparable.
func Contains[T comparable](slice []T, item T) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ContainsWithPredicate checks if a slice contains an item. Use predicate to compare items.
func ContainsWithPredicate[T any](slice []T, item T, predicate func(T, T) bool) bool {
	for _, s := range slice {
		if predicate(s, item) {
			return true
		}
	}
	return false
}

// IndexOf returns the index of the first instance of item in slice, or -1 if item is not present in slice.
func IndexOf[T comparable](slice []T, item T) int {
	for i, s := range slice {
		if s == item {
			return i
		}
	}
	return -1
}

// IndexOfWithPredicate returns the index of the first instance of item in slice, or -1 if item is not present in slice.
// Use predicate to compare items.
func IndexOfWithPredicate[T any](slice []T, item T, predicate func(T, T) bool) int {
	for i, s := range slice {
		if predicate(s, item) {
			return i
		}
	}
	return -1
}

// LastIndexOf returns the index of the last instance of item in slice, or -1 if item is not present in slice.
func LastIndexOf[T comparable](slice []T, item T) int {
	for i := len(slice) - 1; i >= 0; i-- {
		if slice[i] == item {
			return i
		}
	}
	return -1
}

// LastIndexOfWithPredicate returns the index of the last instance of item in slice, or -1 if item is not present in slice.
// Use predicate to compare items.
func LastIndexOfWithPredicate[T any](slice []T, item T, predicate func(T, T) bool) int {
	for i := len(slice) - 1; i >= 0; i-- {
		if predicate(slice[i], item) {
			return i
		}
	}
	return -1
}

// Remove removes the first instance of item from slice, if present.
// Returns true if item was removed, false otherwise.
func Remove[T comparable](slice []T, item T) bool {
	for i, s := range slice {
		if s == item {
			slice = append(slice[:i], slice[i+1:]...)
			return true
		}
	}
	return false
}

// RemoveWithPredicate removes the first instance of item from slice, if present.
// Use predicate to compare items.
// Returns true if item was removed, false otherwise.
func RemoveWithPredicate[T any](slice []T, item T, predicate func(T, T) bool) bool {
	for i, s := range slice {
		if predicate(s, item) {
			slice = append(slice[:i], slice[i+1:]...)
			return true
		}
	}
	return false
}

// RemoveAll removes all instances of item from slice, if present.
// Returns true if item was removed, false otherwise.
func RemoveAll[T comparable](slice []T, item T) bool {
	var removed bool
	for i, s := range slice {
		if s == item {
			slice = append(slice[:i], slice[i+1:]...)
			removed = true
		}
	}
	return removed
}

// RemoveAllWithPredicate removes all instances of item from slice, if present.
// Use predicate to compare items.
// Returns true if item was removed, false otherwise.
func RemoveAllWithPredicate[T any](slice []T, item T, predicate func(T, T) bool) bool {
	var removed bool
	for i, s := range slice {
		if predicate(s, item) {
			slice = append(slice[:i], slice[i+1:]...)
			removed = true
		}
	}
	return removed
}

// RemoveAt removes the item at the given index from slice.
// Returns true if item was removed, false otherwise.
func RemoveAt[T any](slice []T, index int) bool {
	if index < 0 || index >= len(slice) {
		return false
	}
	slice = append(slice[:index], slice[index+1:]...)
	return true
}

// RemoveRange removes the items in the given range from slice.
// Returns true if items were removed, false otherwise.
func RemoveRange[T any](slice []T, start, end int) bool {
	if start < 0 || end < 0 || start >= len(slice) || end >= len(slice) || start > end {
		return false
	}
	slice = append(slice[:start], slice[end+1:]...)
	return true
}

// RemoveIf removes all items from slice for which predicate returns true.
// Returns true if items were removed, false otherwise.
func RemoveIf[T any](slice []T, predicate func(T) bool) bool {
	var removed bool
	for i := 0; i < len(slice); i++ {
		if predicate(slice[i]) {
			slice = append(slice[:i], slice[i+1:]...)
			removed = true
		}
	}
	return removed
}

// Filter returns a new slice containing all items from slice for which predicate returns true.
func Filter[T any](slice []T, predicate func(T) bool) []T {
	var filtered []T
	for _, s := range slice {
		if predicate(s) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

// FilterNot returns a new slice containing all items from slice for which predicate returns false.
func FilterNot[T any](slice []T, predicate func(T) bool) []T {
	var filtered []T
	for _, s := range slice {
		if !predicate(s) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

// Map returns a new slice containing the results of applying the given mapper function to each item in slice.
func Map[T, R any](slice []T, mapper func(T) R) []R {
	var mapped []R
	for _, s := range slice {
		mapped = append(mapped, mapper(s))
	}
	return mapped
}

// RemoveDuplicates removes all duplicate items from slice.
// Returns true if items were removed, false otherwise.
func RemoveDuplicates[T comparable](slice []T) bool {
	var counter = make(map[T]bool)
	var removed bool
	for i := 0; i < len(slice); i++ {
		if counter[slice[i]] {
			slice = append(slice[:i], slice[i+1:]...)
			removed = true
		} else {
			counter[slice[i]] = true
		}
	}
	return removed
}

// Reverse reverses the order of items in slice.
func Reverse[T any](slice []T) {
	for i := 0; i < len(slice)/2; i++ {
		slice[i], slice[len(slice)-i-1] = slice[len(slice)-i-1], slice[i]
	}
}

// Difference returns a new slice containing all items from slice that are not present in other.
func Difference[T comparable](slice, other []T) []T {
	var set = make(map[T]bool)
	for _, s := range other {
		set[s] = true
	}

	var diff []T
	for _, s := range slice {
		if !set[s] {
			diff = append(diff, s)
		}
	}
	return diff
}

// Intersection returns a new slice containing all items from slice that are also present in other.
func Intersection[T comparable](slice, other []T) []T {
	var set = make(map[T]bool)
	for _, s := range other {
		set[s] = true
	}

	var inter []T
	for _, s := range slice {
		if set[s] {
			inter = append(inter, s)
		}
	}
	return inter
}

// Union returns a new slice containing all items from slice and other.
func Union[T comparable](slice, other []T) []T {
	var set = make(map[T]bool)
	for _, s := range slice {
		set[s] = true
	}
	for _, s := range other {
		set[s] = true
	}

	var union []T
	for s := range set {
		union = append(union, s)
	}
	return union
}
