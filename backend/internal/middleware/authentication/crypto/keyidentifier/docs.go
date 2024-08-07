// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// Package keyidentifier provides a utility for generating unique keys that can be used for database operations such as Redis/Valkey or key-value stores.
//
// Compatibility:
//
//   - This package is not directly compatible with Fiber middleware that expects a key generator to provide a direct value (e.g., Fiber Rate Limiter where the IP is used as the key, not the value).
//     It generates keys independently of the [fiber.Ctx] and does not modify the context.
//     Also note that using a direct value (e.g., IP address, other sensitive) as the key, as in the case of Fiber Rate Limiter, can lead to a security compromise (possibly a new CVE) because it exposes sensitive information in logs/commander panel.
//
//   - The key exchange mechanism for ECDSA is not supported (and won't be implemented even though implementing key exchange is easy) because it is not used for external services such as TLS or other external services.
package keyidentifier
