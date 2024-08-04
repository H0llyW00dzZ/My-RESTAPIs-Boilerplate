// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package keyidentifier

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

// GetKeyFunc generates a unique key for each request and returns a function that retrieves the key from the context.
//
// TODO: Implement a Value Mechanism from [fiber.Ctx] for signing/binding a cryptographic technique to the UUID.
func (k *KeyIdentifier) GetKeyFunc() func(*fiber.Ctx) string {
	return func(c *fiber.Ctx) string {
		// Generate a random UUID
		//
		// TODO: Do we really need to improve this by using a cryptographic technique similar to how Bitcoin generates addresses?
		id := utils.UUIDv4()

		// Set the key in the context
		key := k.config.Prefix + id

		// Return the generated key
		return key
	}
}

// GetKey generates a unique key for each request and retrieves it from the context.
func (k *KeyIdentifier) GetKey() string {
	// Generate a random UUID
	//
	// TODO: Do we really need to improve this by using a cryptographic technique similar to how Bitcoin generates addresses?
	id := utils.UUIDv4()

	// Set the key in the context
	key := k.config.Prefix + id

	// Return the generated key
	return key
}
