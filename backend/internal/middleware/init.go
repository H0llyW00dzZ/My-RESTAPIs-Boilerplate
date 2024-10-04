// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package middleware

import "unique"

// Note: This is just a test of a new package introduced in Go 1.23.
var uniqueHostnames []unique.Handle[string]

func init() {
	uniqueHostnames = make([]unique.Handle[string], 0)
}
