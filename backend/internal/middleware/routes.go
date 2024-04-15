// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package middleware

import (
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"h0llyw00dz-template/backend/internal/database"
	log "h0llyw00dz-template/backend/internal/logger"
)

// RegisterRoutes sets up the API routing for the application.
// It organizes routes into versioned groups for better API version management (Currently unimplemented for this boilerplate).
func RegisterRoutes(app *fiber.App, appName, monitorPath string, db database.Service) {
	// Apply the combined middlewares
	registerRouteConfigMiddleware(app)
	// Register the Static Frontend Routes
	registerStaticFrontendRoutes(app, appName, db)
}

// registerRouteConfigMiddleware applies middleware configurations to the Fiber application.
// It sets up the necessary middleware such as recovery, logging, and custom error handling for manipulate panics (Currently unimplemented for this boilerplate).
func registerRouteConfigMiddleware(app *fiber.App) {

	// Favicon front end setup
	// Note: this just an example
	favicon := NewFaviconMiddleware("./frontend/public/assets/images/favicon.ico", "/styles/images/favicon.ico")

	// Recovery middleware setup
	recoverMiddleware := recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			// Log the panic and stack trace
			log.LogUserActivity(c, "attempted to panic occurred")
			log.LogCrashf(MsgPanicOccurred, e)
			log.LogCrashf(MsgStackTrace, debug.Stack())
		},
	})

	// Apply the recover middleware
	app.Use(favicon, recoverMiddleware)
}
