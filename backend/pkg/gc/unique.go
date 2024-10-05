// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package gc

import (
	"unique"

	"github.com/gofiber/fiber/v2"
)

// UniqueMakeFiberCTX is a global variable that holds a function to generate a unique handle for any value.
// It takes a [fiber.Ctx] and a value of any type as input and returns a function that, when called, returns the unique handle value.
//
// The unique handle is generated using the [unique.Make] function, which creates a globally unique identity
// for the provided value.
//
// The returned function, when called, retrieves the actual value of the unique handle using the Value method.
//
// This allows for efficient and consistent retrieval of unique handles for any value throughout the code.
var UniqueMakeFiberCTX = func(c *fiber.Ctx) func(value any) any {
	return func(value any) any {
		return unique.Make(value).Value()
	}
}

// UniqueMake is a global variable that holds a function to generate a unique handle for any value.
// It takes a value of any type as input and returns a function that, when called, returns the unique handle value.
//
// The unique handle is generated using the [unique.Make] function, which creates a globally unique identity
// for the provided value.
//
// The returned function, when called, retrieves the actual value of the unique handle using the Value method.
//
// This allows for efficient and consistent retrieval of unique handles for any value throughout the code.
var UniqueMake = func(value any) func() any {
	return func() any {
		return unique.Make(value).Value()
	}
}
