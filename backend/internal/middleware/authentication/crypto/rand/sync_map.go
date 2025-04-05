// Copyright (c) 2025 H0llyW00dzZ All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package rand

import (
	"errors"
	"sync"
)

var (
	// ErrSyncMapIsEmpty is returned by the SyncMap and SyncMapValue functions
	// when the provided [sync.Map] is empty. This error indicates that there are
	// no elements to choose from.
	ErrSyncMapIsEmpty = errors.New("crypto/rand: sync.Map is empty")
)

// SyncMap selects a random value from the provided [sync.Map].
// It returns the value associated with a randomly chosen key.
// If the map is empty, it returns the zero value of the map's value type and an error.
//
// The function uses a type parameter [K] for keys and [V] for values,
// ensuring type safety when retrieving values from the map.
// It leverages the [sync.Map.Range] method to iterate over the map and collect keys,
// which are then passed to the [Choice] function to select one at random.
//
// Example usage:
//
//	m := &sync.Map{}
//	m.Store("a", 1)
//	m.Store("b", 2)
//	value, err := rand.SyncMap[string, int](m)
//	if err != nil {
//	    // handle error, you poggers.
//	}
//	fmt.Println(value)
func SyncMap[K comparable, V any](m *sync.Map) (V, error) {
	var zero V
	keys := []K{}

	m.Range(func(key, _ any) bool {
		k, ok := key.(K)
		if ok {
			keys = append(keys, k)
		}
		return true
	})

	if len(keys) == 0 {
		return zero, ErrSyncMapIsEmpty
	}

	key, err := Choice(keys)
	if err != nil {
		return zero, err
	}

	value, _ := m.Load(key)
	return value.(V), nil
}

// SyncMapValue selects a random value from the provided [sync.Map].
// It collects all values into a slice and uses the [Choice] function
// to randomly select one of these values.
// If the map is empty, it returns the zero value of the map's value type and an error.
//
// This function is useful when the keys are not needed, and only a random value is desired.
// The function uses a type parameter [V] for values to ensure type safety.
//
// Example usage:
//
//	m := &sync.Map{}
//	m.Store(1, "one")
//	m.Store(2, "two")
//	value, err := rand.SyncMapValue[int, string](m)
//	if err != nil {
//	    // handle error, you poggers.
//	}
//	fmt.Println(value)
func SyncMapValue[K comparable, V any](m *sync.Map) (V, error) {
	var zero V
	values := []V{}

	m.Range(func(_, value any) bool {
		v, ok := value.(V)
		if ok {
			values = append(values, v)
		}
		return true
	})

	if len(values) == 0 {
		return zero, ErrSyncMapIsEmpty
	}

	return Choice(values)
}
