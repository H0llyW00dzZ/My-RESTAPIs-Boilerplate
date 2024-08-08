// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package htmxauth

import (
	htmx "h0llyw00dz-template/frontend/htmx/error_page_handler"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	bpool "github.com/valyala/bytebufferpool"
)

// renderAndSend renders the given HTMX component using the provided Fiber context
// and sends the rendered HTML content as an HTTP response with the specified status code.
func (vd *viewData) renderAndSend(c *fiber.Ctx, statusCode int, component templ.Component) error {
	// Get a buffer from the pool for efficient string building.
	buf := bpool.Get()

	// Use defer to guarantee buffer cleanup (reset and return to the pool)
	// even if an error occurs during rendering.
	defer func() {
		buf.Reset()    // Reset the buffer to prevent data leaks.
		bpool.Put(buf) // Return the buffer to the pool for reuse.
	}()

	// Render the HTMX component into the byte buffer.
	if err := component.Render(c.Context(), buf); err != nil {
		// Handle any rendering errors by returning an internal server error page.
		return vd.renderErrorPage(c, err)
	}

	// Send the response
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	return c.Status(statusCode).SendString(buf.String())
}

// renderErrorPage logs an error and renders a generic error page using the htmx.NewErrorHandler middleware.
func (vd *viewData) renderErrorPage(c *fiber.Ctx, err error) error {
	// Set the error in the Fiber context
	c.Context().SetUserValue("error", err)

	// Use the htmx.NewErrorHandler middleware to handle the error and render the appropriate error page.
	return htmx.NewErrorHandler(c)
}
