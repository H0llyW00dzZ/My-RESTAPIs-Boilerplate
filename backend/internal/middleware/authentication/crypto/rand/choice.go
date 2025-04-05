// Copyright (c) 2025 H0llyW00dzZ All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package rand

import (
	"crypto/rand"
	"errors"
	"math/big"
)

// ErrEmptyChoices is returned by the Choice function when the provided slice is empty.
// This error indicates that there are no elements to select from, ensuring that the
// function handles empty input gracefully and predictably.
var ErrEmptyChoices = errors.New("crypto/rand: choices slice is empty")

// Choice selects a random element from the given slice using [crypto/rand] for [secure randomness].
//
//   - See [generics] for more information on using generics in Go.
//   - See [secure randomness] for details on [cryptographically] secure random number generation.
//
// [generics]: https://go.dev/doc/tutorial/generics
// [secure randomness]: https://en.wikipedia.org/wiki/Cryptographically_secure_pseudorandom_number_generator
// [cryptographically]: https://en.wikipedia.org/wiki/Cryptography
func Choice[T any](choices []T) (T, error) {
	var zero T
	// Should not crash when the slice is empty.
	if len(choices) == 0 {
		return zero, ErrEmptyChoices
	}

	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(choices))))
	if err != nil {
		return zero, err
	}
	return choices[nBig.Int64()], nil
}
