// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// Package keyidentifier provides a utility for generating unique keys that can be used for database operations such as Redis/Valkey or key-value stores.
//
// Compatibility:
//   - This package is not directly compatible with Fiber middleware that expects a key generator direct value (e.g, Fiber Rate Limiter).
//     It generates keys independently of the [fiber.Ctx] and does not modify the context.
package keyidentifier
