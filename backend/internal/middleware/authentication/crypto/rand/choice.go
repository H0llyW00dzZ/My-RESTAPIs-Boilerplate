// Copyright (c) 2025 H0llyW00dzZ All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package rand

import (
	"crypto/rand"
	"math/big"
)

// Choice selects a random element from the given slice using [crypto/rand] for [secure randomness].
//
//   - See [generics] for more information on using generics in Go.
//   - See [secure randomness] for details on cryptographically secure random number generation.
//
// [generics]: https://go.dev/doc/tutorial/generics
// [secure randomness]: https://en.wikipedia.org/wiki/Cryptographically_secure_pseudorandom_number_generator
func Choice[T any](choices []T) (T, error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(choices))))
	if err != nil {
		var zero T
		return zero, err
	}
	return choices[nBig.Int64()], nil
}
