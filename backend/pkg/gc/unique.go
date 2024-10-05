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

// UniqueMakeT is a global variable that holds a function to generate a unique handle for any comparable value.
// It takes a value of any comparable type as input and returns a function that, when called, returns the unique handle.
//
// The unique handle is generated using the [unique.Make] function, which creates a globally unique identity
// for the provided value. The handle is of type [unique.Handle[T]], where T is the type of the input value.
//
// The returned function, when called, retrieves the actual value of the unique handle using the Value method.
// The Value method returns a shallow copy of the original value that produced the handle.
//
// Two handles compare equal exactly if the two values used to create the handles would have also compared equal.
// The comparison of two handles is trivial and typically much more efficient than comparing the values used to create them.
//
// This allows for efficient and consistent retrieval of unique handles for any comparable value throughout the code.
//
// Note: This is suitable for horizontal pod autoscaling (HPA) due to its dynamic resource usage.
// As the CPU usage grows, the memory usage is spread across every pod, just like waterfall.
// This allows for efficient scaling and resource utilization in a distributed environment.
var UniqueMakeT = func(value any) func() unique.Handle[any] {
	return func() unique.Handle[any] {
		return unique.Make(value)
	}
}

// UniqueMakeTFiberCTX is a global variable that holds a function to generate a unique handle for any comparable value within a Fiber context.
// It takes a [fiber.Ctx] and a value of any comparable type as input and returns a function that, when called, returns the unique handle.
//
// The unique handle is generated using the [unique.Make] function, which creates a globally unique identity
// for the provided value. The handle is of type [unique.Handle[T]], where T is the type of the input value.
//
// The returned function, when called, retrieves the actual value of the unique handle using the Value method.
// The Value method returns a shallow copy of the original value that produced the handle.
//
// Two handles compare equal exactly if the two values used to create the handles would have also compared equal.
// The comparison of two handles is trivial and typically much more efficient than comparing the values used to create them.
//
// This allows for efficient and consistent retrieval of unique handles for any comparable value throughout the code within a Fiber context.
//
// Note: This is suitable for horizontal pod autoscaling (HPA) due to its dynamic resource usage.
// As the CPU usage grows, the memory usage is spread across every pod, just like waterfall.
// This allows for efficient scaling and resource utilization in a distributed environment.
var UniqueMakeTFiberCTX = func(c *fiber.Ctx) func(value any) unique.Handle[any] {
	return func(value any) unique.Handle[any] {
		return unique.Make(value)
	}
}
