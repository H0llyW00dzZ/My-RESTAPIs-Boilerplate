// Copyright (c) 2025 H0llyW00dzZ All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package rand

import "errors"

var (
	// ErrMapIsEmpty is returned by the Map function when the provided map is empty.
	// This error indicates that there are no elements to choose from.
	ErrMapIsEmpty = errors.New("crypto/rand: map is empty")
)

// Map selects a random value from the provided map.
// It returns the value associated with a randomly chosen key.
// If the map is empty, it returns the zero value of the map's value type and an error.
func Map[K comparable, V any](m map[K]V) (V, error) {
	if len(m) == 0 {
		var zero V
		return zero, ErrMapIsEmpty
	}

	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	key, err := Choice(keys)
	if err != nil {
		var zero V
		return zero, err
	}
	return m[key], nil
}

// MapValue selects a random value from the provided map.
// It does this by first collecting all values into a slice and then using the [Choice] function
// to randomly select one of these values.
// If the map is empty, it returns the zero value of the map's value type and an error.
func MapValue[K comparable, V any](m map[K]V) (V, error) {
	if len(m) == 0 {
		var zero V
		return zero, ErrMapIsEmpty
	}

	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}

	return Choice(values)
}
