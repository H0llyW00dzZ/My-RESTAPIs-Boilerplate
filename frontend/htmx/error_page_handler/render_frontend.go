// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

// Demo: dev.btz.pm

package htmx

import (
	"bytes"
	"h0llyw00dz-template/backend/pkg/restapis/helper"

	"github.com/gofiber/fiber/v2"
)

// views represents the data that will be passed to the view template.
type views struct {
	title      string // The title of the page.
	cfheader   string // A list of Content-Security-Policy headers (e.g, CF-Ray-ID).
	xRequestID string // X-Request-ID Header
}

// viewData is a structure that contains the data for rendering a view.
// It embeds the views structure to provide access to title and headers.
type viewData struct {
	*views // Embed the views structure for access to title and headers.
}

// handleError is a general error handler for both API and frontend routes.
func handleError(c *fiber.Ctx, e *fiber.Error, vd *viewData) error {
	switch e.Code {
	case fiber.StatusNotFound:
		vd.title = PageNotFound + " - " + c.App().Config().AppName
		return vd.PageNotFoundHandler(c)
	case fiber.StatusForbidden:
		vd.title = PageForbidden + " - " + c.App().Config().AppName
		return vd.PageForbidden403Handler(c)
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
	component := PageNotFound404(v.title, v.cfheader, v.xRequestID)

	// Note: This Optional can be used to builder string. However,
	// it is intended for low-level operations where the efficiency of using a string builder is not significant.
	buf := new(bytes.Buffer)
	if err := component.Render(c.Context(), buf); err != nil {
		return v.renderErrorPage(c, fiber.StatusInternalServerError, "Error rendering PageNotFound: %v", err)
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return c.Status(fiber.StatusNotFound).SendString(buf.String())
}

// PageForbidden403Handler renders the 403 Forbidden error page.
//
// This function takes a Fiber context and renders the 403 page.
func (v *viewData) PageForbidden403Handler(c *fiber.Ctx) error {
	component := PageForbidden403(v.title, v.cfheader, v.xRequestID)

	// Note: This Optional can be used to builder string. However,
	// it is intended for low-level operations where the efficiency of using a string builder is not significant.
	buf := new(bytes.Buffer)
	if err := component.Render(c.Context(), buf); err != nil {
		return v.renderErrorPage(c, fiber.StatusInternalServerError, "Error rendering Forbidden Page: %v", err)
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return c.Status(fiber.StatusForbidden).SendString(buf.String())
}

// Page500InternalServerHandler handles 500 Internal Server errors.
func (v *viewData) Page500InternalServerHandler(c *fiber.Ctx) error {
	component := PageInternalServerError500(v.title, v.cfheader, v.xRequestID)

	// Note: This Optional can be used to builder string. However,
	// it is intended for low-level operations where the efficiency of using a string builder is not significant.
	buf := new(bytes.Buffer)
	if err := component.Render(c.Context(), buf); err != nil {
		return v.renderErrorPage(c, fiber.StatusInternalServerError, "Error rendering Internal Server Error Page: %v", err)
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return c.Status(fiber.StatusForbidden).SendString(buf.String())
}

// GenericErrorInternalServerHandler handles Generic 500 Internal Server errors.
func (v *viewData) GenericErrorInternalServerHandler(c *fiber.Ctx, err error) error {
	// Return a JSON response with the 500 Internal Server Error status code
	return helper.SendErrorResponse(c, fiber.StatusInternalServerError, fiber.ErrInternalServerError.Message)
}
