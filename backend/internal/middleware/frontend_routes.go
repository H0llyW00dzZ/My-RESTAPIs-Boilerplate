// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package middleware

import (
	"h0llyw00dz-template/backend/internal/database"
	"time"

	"github.com/gofiber/fiber/v2"
)

// registerStaticFrontendRoutes sets up the frontend routing for a given Fiber app.
// It registers custom application routes and static file serving.
//
// Parameters:
//
//	app: The Fiber app on which to register the routes.
func registerStaticFrontendRoutes(app *fiber.App, _ string, _ database.Service) {
	// Note: This is an example, and the unused parameters are assigned to "_" to avoid compiler errors.
	// Register static file serving
	app.Static("/styles/", "./frontend/public/assets", fiber.Static{
		Compress: true,
		// Note: It's important to disable this when using middleware cache to avoid confusion,
		// as caching is already handled by the middleware cache.
		CacheDuration: -1 * time.Microsecond,
	})
}
