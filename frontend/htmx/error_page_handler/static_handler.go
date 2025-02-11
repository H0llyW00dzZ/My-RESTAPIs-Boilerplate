// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package htmx

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// NewStaticHandleVersionedAPIError handles errors for versioned static REST API routes.
func NewStaticHandleVersionedAPIError(c *fiber.Ctx, e *fiber.Error) error {
	timeYearNow := time.Now().Year()
	vd := &viewData{
		views: &views{},
	}

	// Convert the integer year to a string
	vd.timeYears = strconv.Itoa(timeYearNow)
	// Get Application name
	vd.appName = c.App().Config().AppName

	if cloudflareRayID := c.Get(CloudflareRayIDHeader); cloudflareRayID != "" {
		vd.cfheader = cloudflareRayID
		// Get xRequestID Where it was generated.
	} else if xRequestID := c.Locals(XRequestID); xRequestID != nil {
		vd.xRequestID = xRequestID.(string)
	}

	// Get cspRandom Where it was generated.
	if cspRandom := c.Locals(cspRandom); cspRandom != nil {
		vd.cspRandom = cspRandom.(string)
	}

	// Convert the integer httpStatusCode to a string
	if e.Code != fiber.StatusOK {
		vd.httpStatus = strconv.Itoa(e.Code)
		c.Locals("error", e.Error())
	}

	return handleError(c, e, vd)
}

// NewStaticHandleFrontendError handles errors for static frontend routes.
func NewStaticHandleFrontendError(c *fiber.Ctx, e *fiber.Error) error {
	timeYearNow := time.Now().Year()
	vd := &viewData{
		views: &views{},
	}

	// Convert the integer year to a string
	vd.timeYears = strconv.Itoa(timeYearNow)
	// Get Application name
	vd.appName = c.App().Config().AppName

	if cloudflareRayID := c.Get(CloudflareRayIDHeader); cloudflareRayID != "" {
		vd.cfheader = cloudflareRayID
		// Get xRequestID Where it was generated.
	} else if xRequestID := c.Locals(XRequestID); xRequestID != nil {
		vd.xRequestID = xRequestID.(string)
	}

	// Get cspRandom Where it was generated.
	if cspRandom := c.Locals(cspRandom); cspRandom != nil {
		vd.cspRandom = cspRandom.(string)
	}

	// Convert the integer httpStatusCode to a string
	if e.Code != fiber.StatusOK {
		vd.httpStatus = strconv.Itoa(e.Code)
		c.Locals("error", e.Error())
	}
	return handleError(c, e, vd)
}
