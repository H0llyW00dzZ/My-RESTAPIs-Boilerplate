// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package htmx

import (
	"github.com/gofiber/fiber/v2"
)

// NewStaticHandleVersionedAPIError handles errors for versioned static REST API routes.
func NewStaticHandleVersionedAPIError(c *fiber.Ctx, e *fiber.Error) error {
	// Get xRequestID Where it was generated.
	xRequestID := c.Locals(XRequestID)
	vd := &viewData{
		views: &views{},
	}
	cloudflareRayID := c.Get(CloudflareRayIDHeader)
	if cloudflareRayID != "" {
		vd.cfheader = cloudflareRayID
	} else if xRequestID != nil {
		vd.xRequestID = xRequestID.(string)

	}
	return handleError(c, e, vd)
}

// NewStaticHandleFrontendError handles errors for static frontend routes.
func NewStaticHandleFrontendError(c *fiber.Ctx, e *fiber.Error) error {
	// Get xRequestID Where it was generated.
	xRequestID := c.Locals(XRequestID)
	vd := &viewData{
		views: &views{},
	}
	cloudflareRayID := c.Get(CloudflareRayIDHeader)
	if cloudflareRayID != "" {
		vd.cfheader = cloudflareRayID
	} else if xRequestID != nil {
		vd.xRequestID = xRequestID.(string)

	}
	return handleError(c, e, vd)
}
