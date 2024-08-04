// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package keyidentifier

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

// GetKeyFunc generates a unique key for each request and returns a function that retrieves the key from the context.
func (k *KeyIdentifier) GetKeyFunc() func(*fiber.Ctx) string {
	return func(c *fiber.Ctx) string {
		// Generate a random UUID
		id := utils.UUIDv4()

		// Sign the UUID using ECDSA
		if k.config.PrivateKey != nil && k.config.SignedContextKey != nil {
			signature, err := k.signUUID(id)
			if err != nil {
				panic(fmt.Errorf("failed to sign UUID: %v", err))
			}

			// Store the Signature for future use
			c.Locals(k.config.SignedContextKey, signature)
		}

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
