// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package htmx

import (
	"h0llyw00dz-template/backend/pkg/restapis/helper"
	"os"

	"github.com/gofiber/fiber/v2"
)

const (
	// InternalServerError is the standard error message for a 500 Internal Server Error.
	InternalServerError = "500 Internal Server Error"

	// PageNotFound is the standard error message for a 404 Page Not Found error.
	PageNotFound = "404 Page Not Found"

	// PageForbidden is the standard error message for a 403 Forbidden error.
	PageForbidden = "403 Forbidden"
)

// Define Cloudflare formats.
//
// Note: Other header can be found here
// https://developers.cloudflare.com/fundamentals/reference/http-request-headers/#:~:text=%E2%80%8B%E2%80%8B%20CF-Connecting-IP,to%20Restoring%20original%20visitor%20IPs.
const (
	CloudflareRayIDHeader = "cf-ray"
)

// NewErrorHandler is a middleware that handles errors for all routes (dynamic).
//
// This middleware intercepts any errors that occur during route handling
// and provides a custom error response.
func NewErrorHandler(c *fiber.Ctx) error {
	vd := &viewData{
		views: &views{},
	}
	cloudflareRayID := c.Get(CloudflareRayIDHeader)
	if cloudflareRayID != "" {
		vd.cfheader = cloudflareRayID
	}

	// Call the next route handler and catch any errors
	if err := c.Next(); err != nil {
		isAPI := c.Hostname() == os.Getenv("API_SUB_DOMAIN")
		return errorHandler(c, err, vd, isAPI)
	}

	// No errors, continue with the next middleware
	return nil
}

// errorHandler is a general error handler function.
// It takes a Fiber context, an error, and a viewData struct.
// You can add additional parameters as needed (e.g., a flag for API vs. frontend).
func errorHandler(c *fiber.Ctx, err error, vd *viewData, isAPI bool) error {
	if e, ok := err.(*fiber.Error); ok {
		// Handle specific error codes based on context
		if isAPI {
			return handleAPIError(c, e)
		}
		return handleFrontendError(c, e, vd)
	}
	return handleGenericError(c, err, vd, isAPI)
}

// handleAPIError handles errors for REST API routes.
func handleAPIError(c *fiber.Ctx, e *fiber.Error) error {
	// Customize API error handling, e.g., JSON responses with appropriate status codes.
	switch e.Code {
	case fiber.StatusNotFound, fiber.StatusForbidden:
		// Return a JSON response for 404 or 403 errors in versioned APIs
		return helper.SendErrorResponse(c, e.Code, e.Message)
	default:
		return helper.SendErrorResponse(c, fiber.StatusInternalServerError, "Internal Server Error")
	}
}

// handleFrontendError handles errors for frontend routes.
func handleFrontendError(c *fiber.Ctx, e *fiber.Error, vd *viewData) error {
	switch e.Code {
	case fiber.StatusNotFound:
		// Render the 404 error page for frontend routes
		vd.title = PageNotFound + " - " + c.App().Config().AppName
		return vd.PageNotFoundHandler(c)
	case fiber.StatusForbidden:
		// Render the 403 error page for frontend routes
		vd.title = PageForbidden + " - " + c.App().Config().AppName
		return vd.PageForbidden403Handler(c)
	default:
		vd.title = InternalServerError + " - " + c.App().Config().AppName
		// Fallback to the general error page for other errors in frontend routes
		return vd.Page500InternalServerHandler(c, e)
	}
}

// handleGenericError handles non-fiber.Error errors.
func handleGenericError(c *fiber.Ctx, err error, vd *viewData, isAPI bool) error {
	if isAPI {
		return helper.SendErrorResponse(c, fiber.StatusInternalServerError, "Internal Server Error")
	}

	vd.title = InternalServerError + " - " + c.App().Config().AppName
	return vd.Page500InternalServerHandler(c, err)
}
