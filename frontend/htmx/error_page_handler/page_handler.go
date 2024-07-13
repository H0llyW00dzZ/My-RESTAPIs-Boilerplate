// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package htmx

import (
	"h0llyw00dz-template/backend/pkg/restapis/helper"
	"h0llyw00dz-template/env"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	// PageInternalServerError is the standard error message for a 500 Internal Server Error.
	PageInternalServerError = "500 Internal Server Error"

	// PageNotFound is the standard error message for a 404 Page Not Found error.
	PageNotFound = "404 Page Not Found"

	// PageForbidden is the standard error message for a 403 Forbidden error.
	PageForbidden = "403 Forbidden"
	// PageServiceUnavailableError is the standard error message for a 503 Service Unavailable
	PageServiceUnavailableError = "503 Service Unavailable"
)

// Define Cloudflare formats.
//
// Note: Other header can be found here
// https://developers.cloudflare.com/fundamentals/reference/http-request-headers
const (
	CloudflareRayIDHeader = "cf-ray"
	XRequestID            = "visitor_uuid"
	cspRandom             = "csp_random"
)

// NewErrorHandler is a middleware that handles errors for all routes (dynamic).
//
// This middleware intercepts any errors that occur during route handling
// and provides a custom error response.
func NewErrorHandler(c *fiber.Ctx) error {
	timeYearNow := time.Now().Year()
	// Get xRequestID Where it was generated.
	xRequestID := c.Locals(XRequestID)
	vd := &viewData{
		views: &views{},
	}

	// Convert the integer year to a string
	vd.timeYears = strconv.Itoa(timeYearNow)
	// Get Application name
	vd.appName = c.App().Config().AppName

	cloudflareRayID := c.Get(CloudflareRayIDHeader)
	if cloudflareRayID != "" {
		vd.cfheader = cloudflareRayID
	} else if xRequestID != nil {
		vd.xRequestID = xRequestID.(string)
	}

	// Get cspRandom Where it was generated.
	cspRandom := c.Locals(cspRandom)
	if cspRandom != nil {
		vd.cspRandom = cspRandom.(string)
	}

	// Call the next route handler and catch any errors
	if err := c.Next(); err != nil {
		isAPI := c.Hostname() == os.Getenv(env.APISUBDOMAIN)
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
		// Convert the integer httpStatusCode to a string
		if e.Code != fiber.StatusOK {
			vd.httpStatus = strconv.Itoa(e.Code)
		}

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
	case fiber.StatusNotFound, fiber.StatusForbidden,
		fiber.StatusServiceUnavailable, fiber.StatusUnauthorized:
		// Return a JSON response for 404, 403, 503, 401 errors in versioned APIs
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
	case fiber.StatusInternalServerError:
		// Render the 500 error page for frontend routes
		vd.title = PageInternalServerError + " - " + c.App().Config().AppName
		return vd.Page500InternalServerHandler(c)
	case fiber.StatusServiceUnavailable:
		// Render the 503 error page for frontend routes
		vd.title = PageServiceUnavailableError + " - " + c.App().Config().AppName
		return vd.PageServiceUnavailableHandler(c)
	case fiber.StatusUnauthorized:
		// Render the 401 error page for frontend routes
		vd.title = fiber.ErrUnauthorized.Message + " - " + c.App().Config().AppName
		return vd.PageUnauthorizeHandler(c)
	default:
		vd.title = PageInternalServerError + " - " + c.App().Config().AppName
		// Fallback to the general error page for other errors in frontend routes
		return vd.GenericErrorInternalServerHandler(c, e)
	}
}

// handleGenericError handles non-fiber.Error errors.
func handleGenericError(c *fiber.Ctx, err error, vd *viewData, isAPI bool) error {
	if isAPI {
		return helper.SendErrorResponse(c, fiber.StatusInternalServerError, "Internal Server Error")
	}

	vd.title = PageInternalServerError + " - " + c.App().Config().AppName
	return vd.GenericErrorInternalServerHandler(c, err)
}
