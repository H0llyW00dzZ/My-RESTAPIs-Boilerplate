// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package htmx

import "github.com/gofiber/fiber/v2"

// NewStaticHandleVersionedAPIError handles errors for versioned static REST API routes.
func NewStaticHandleVersionedAPIError(c *fiber.Ctx, e *fiber.Error) error {
	vd := &viewData{
		views: &views{},
	}
	cloudflareRayID := c.Get(CloudflareRayIDHeader)
	if cloudflareRayID != "" {
		vd.cfheader = cloudflareRayID
	}
	return handleError(c, e, vd)
}

// NewStaticHandleFrontendError handles errors for static frontend routes.
func NewStaticHandleFrontendError(c *fiber.Ctx, e *fiber.Error) error {
	vd := &viewData{
		views: &views{},
	}
	cloudflareRayID := c.Get(CloudflareRayIDHeader)
	if cloudflareRayID != "" {
		vd.cfheader = cloudflareRayID
	}
	return handleError(c, e, vd)
}
