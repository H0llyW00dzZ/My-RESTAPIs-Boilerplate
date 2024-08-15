// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// Demo: closed atm due lazy setup everything.

package htmx

import (
	"h0llyw00dz-template/backend/pkg/mime"
	"h0llyw00dz-template/backend/pkg/restapis/helper"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	bpool "github.com/valyala/bytebufferpool"
)

// views represents the data that will be passed to the view template.
type views struct {
	title      string // The title of the page.
	cfheader   string // A list of Content-Security-Policy headers (e.g, CF-Ray-ID).
	xRequestID string // X-Request-ID Header
	cspRandom  string // Content-Security-Policy
	timeYears  string // Store the year as an string
	appName    string // The Fiber application name
	httpStatus string // HTTP Status Code as an string
}

// viewData is a structure that contains the data for rendering a view.
// It embeds the views structure to provide access to title and headers.
type viewData struct {
	*views // Embed the views structure for access to title and headers.
}

// handleError is a general error handler for both API and frontend routes.
//
// Note: This different way, unlike static "handleFrontendError" or "handleAPIError", and this useful for multiple website (e.g, frontend, restapi, hostname)
// also this used for hostname (e.g, host.example.com) smiliar a default load balancer index page (e.g, nginx)
func handleError(c *fiber.Ctx, e *fiber.Error, vd *viewData) error {
	switch e.Code {
	case fiber.StatusNotFound:
		vd.title = PageNotFound + " - " + c.App().Config().AppName
		return vd.PageNotFoundHandler(c)
	case fiber.StatusForbidden:
		vd.title = PageForbidden + " - " + c.App().Config().AppName
		return vd.PageForbidden403Handler(c)
	case fiber.StatusServiceUnavailable:
		vd.title = PageServiceUnavailableError + " - " + c.App().Config().AppName
		return vd.PageServiceUnavailableHandler(c)
	case fiber.StatusUnauthorized:
		vd.title = fiber.ErrUnauthorized.Message + " - " + c.App().Config().AppName
		return vd.PageUnauthorizeHandler(c)
	case fiber.StatusBadRequest:
		vd.title = fiber.ErrBadRequest.Message + " - " + c.App().Config().AppName
		return vd.PageBadRequestHandler(c)
	default:
		vd.title = PageInternalServerError + " - " + c.App().Config().AppName
		return vd.Page500InternalServerHandler(c)
	}
}

// renderErrorPage logs an error and renders a generic error page.
//
// This function takes a Fiber context, a status code, a log message,
// and an optional error and sends an error response to the client.
//
// Note: The "_" parameters were previously used for logging but have been removed
// since logging is now handled at the internal package level and cannot be imported here.
func (v *viewData) renderErrorPage(c *fiber.Ctx, statusCode int, _ string, _ error) error {
	return helper.SendErrorResponse(c, statusCode, "An error occurred while rendering the page.")
}

// PageNotFoundHandler renders the 404 Not Found error page.
//
// This function takes a Fiber context and renders the 404 page.
func (v *viewData) PageNotFoundHandler(c *fiber.Ctx) error {
	component := PageNotFound404(*v) // magic pointer.
	return v.renderAndSend(c, fiber.StatusNotFound, component)
}

// PageForbidden403Handler renders the 403 Forbidden error page.
//
// This function takes a Fiber context and renders the 403 page.
func (v *viewData) PageForbidden403Handler(c *fiber.Ctx) error {
	component := PageForbidden403(*v) // magic pointer.
	return v.renderAndSend(c, fiber.StatusForbidden, component)
}

// Page500InternalServerHandler handles 500 Internal Server errors.
func (v *viewData) Page500InternalServerHandler(c *fiber.Ctx) error {
	component := PageInternalServerError500(*v) // magic pointer.
	return v.renderAndSend(c, fiber.StatusInternalServerError, component)
}

// PageServiceUnavailableHandler handles 503 Service Unavailable errors.
func (v *viewData) PageServiceUnavailableHandler(c *fiber.Ctx) error {
	component := PageServiceUnavailable(*v) // magic pointer.
	return v.renderAndSend(c, fiber.StatusServiceUnavailable, component)
}

// PageUnauthorizeHandler handles 401 Authentication required.
func (v *viewData) PageUnauthorizeHandler(c *fiber.Ctx) error {
	component := PageUnauthorize401(*v) // magic pointer.
	return v.renderAndSend(c, fiber.StatusUnauthorized, component)
}

// PageBadRequestHandler handles 400 Bad Request.
func (v *viewData) PageBadRequestHandler(c *fiber.Ctx) error {
	component := PageBadRequest400(*v) // magic pointer.
	return v.renderAndSend(c, fiber.StatusBadRequest, component)
}

// GenericErrorInternalServerHandler handles Generic 500 Internal Server errors.
func (v *viewData) GenericErrorInternalServerHandler(c *fiber.Ctx, err error) error {
	// Return a JSON response with the 500 Internal Server Error status code
	return helper.SendErrorResponse(c, fiber.StatusInternalServerError, fiber.ErrInternalServerError.Message)
}

// renderAndSend renders the given HTMX component using the provided Fiber context
// and sends the rendered HTML content as an HTTP response with the specified status code.
//
// This function utilizes valyala/bytebufferpool for efficient string building,
// ensuring reduced garbage collection overhead, however it not possible to made it
// cracked zer0-ms response (required senior only) in production, unless encapsulating
// load balancer so it might possible to made it cracked zer0-ms response (required senior only).
//
// It follows the DRY principle (Don't Repeat Yourself) by encapsulating the
// common logic for rendering and sending HTMX component responses,
// making the code cleaner and easier to maintain.
func (v *viewData) renderAndSend(c *fiber.Ctx, statusCode int, component templ.Component) error {
	// Note: This Optional can be used to builder string. However,
	// it is intended for low-level operations where the efficiency of using a string builder is not significant.
	//
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
		return v.renderErrorPage(c, fiber.StatusInternalServerError, "Error rendering component: %v", err)
	}

	// Convert the byte buffer to a string.
	renderedHTML := buf.String()

	// Set the appropriate Content-Type header based on the presence of non-ASCII characters.
	if !mime.IsASCII(renderedHTML) {
		// If non-ASCII characters are present, use the MIME type with charset
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	} else {
		// If the response body contains only ASCII characters, use the MIME type without charset
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	}

	// Send the rendered HTML content as a response with the appropriate status code.
	return c.Status(statusCode).SendString(renderedHTML)
}
