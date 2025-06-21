// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// Note: The struct types defined here are partially copied from the htmx error handler and are subject to change due to the extensive struct building required for the front end.

package htmxauth

import (
	"github.com/gofiber/fiber/v2"
)

// viewData is a structure that contains the data for rendering a view.
// It embeds a pointer to the views structure to provide access to title and headers.
type viewData struct {
	*views // Embed a pointer to the views structure for access to title and headers.
}

// views represents the data that will be passed to the view template.
// It embeds a pointer to the index structure.
type views struct {
	*index
}

type index struct {
	title      string // The title of the page.
	cfheader   string // A list of Content-Security-Policy headers (e.g, CF-Ray-ID).
	xRequestID string // X-Request-ID Header
	cspRandom  string // Content-Security-Policy
	timeYears  string // Store the year as an string
	appName    string // The Fiber application name
	httpStatus string // HTTP Status Code as an string
}

// IndexHandler returns a Fiber handler that handles the rendering of the index page.
func IndexHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Prepare the view data
		vd := &viewData{
			views: &views{
				index: &index{
					title:   "Login",
					appName: c.App().Config().AppName,
				},
			},
		}

		// Render the index page using the Base template
		component := Base(*vd)
		return vd.renderAndSend(c, fiber.StatusOK, component)
	}
}
