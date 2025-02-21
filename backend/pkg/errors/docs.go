// Copyright (c) 2025 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// Package errors provides simple error wrapping functionality
// using the built-in errors and fmt packages.
//
// This package offers two main functions:
// - Wrap: Wraps an error with a simple message.
// - WrapF: Wraps an error with a formatted message.
//
// These functions serve as a replacement for github.com/pkg/errors
// by utilizing Go's native error wrapping capabilities introduced in Go 1.13.
//
// Reasons to use this package:
// - Leverage Go's built-in error wrapping and inspection.
// - Avoid external dependencies that are no longer maintained, keeping the codebase lightweight.
// - Easily wrap multiple errors to provide additional context.
// - Enhance reusability by standardizing error handling across projects.
//
// Example usage:
//
//	err := stdErrors.New("original error")
//	wrappedErr := errors.Wrap(err, "additional context")
//	wrappedErrF := errors.WrapF(err, "context with value: %d", 42)
//
// These functions are particularly useful for adding context to errors
// in a consistent and idiomatic manner, making the code easier to debug
// and maintain.
package errors
