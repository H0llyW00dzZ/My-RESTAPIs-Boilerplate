// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package htmxlogin

import "github.com/gofiber/fiber/v2"

// SignupHandler returns a Fiber handler that handles the signup form submission.
func SignupHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		vd := &viewData{}
		// TODO: Implement signup logic here
		// ...

		// Render a response using the renderAndSend method
		component := Base(*vd)
		return vd.renderAndSend(c, fiber.StatusOK, component)
	}
}
