// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package htmxlogin

import (
	htmx "h0llyw00dz-template/frontend/htmx/error_page_handler"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler returns a Fiber handler that handles errors using the htmx.NewErrorHandler middleware.
func ErrorHandler() fiber.Handler {
	return htmx.NewErrorHandler
}
