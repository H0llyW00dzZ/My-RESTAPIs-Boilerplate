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

// Note: This method works well Docs: https://github.com/gofiber/fiber/issues/750
// Also note that There is no limit to this feature. For example, you can add a billion domains or subdomains.
type (
	// Host represents a subdomain or domain host configuration.
	// It contains a reference to a Fiber application instance.
	Host struct {
		// Fiber is a pointer to a Fiber application instance.
		// It represents the Fiber app associated with the subdomain or domain host.
		Fiber *fiber.App
	}
)

// RegisterRoutes sets up the API routing for the application.
// It organizes routes into versioned groups for better API version management (Currently unimplemented for this boilerplate).
func RegisterRoutes(app *fiber.App, appName, monitorPath string, db database.Service) {
	// Apply the combined middlewares
	registerRouteConfigMiddleware(app)
	// Register the REST APIs Routes
	registerRESTAPIsRoutes(app, db)
	// Register the Static Frontend Routes
	registerStaticFrontendRoutes(app, appName, db)
}

// registerRouteConfigMiddleware applies middleware configurations to the Fiber application.
// It sets up the necessary middleware such as recovery, logging, and custom error handling for manipulate panics (Currently unimplemented for this boilerplate).
func registerRouteConfigMiddleware(app *fiber.App) {

	// Favicon front end setup
	// Note: this just an example
	favicon := NewFaviconMiddleware("./frontend/public/assets/images/favicon.ico", "/favicon.ico")

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
	app.Use(recoverMiddleware, favicon)
}

// DomainRouter is a middleware function that handles subdomain or domain routing.
// It takes a map of subdomain or domain hosts and routes the request to the corresponding Fiber app.
func DomainRouter(hosts map[string]*Host) fiber.Handler {
	return func(c *fiber.Ctx) error {
		host := hosts[c.Hostname()]
		if host == nil {
			return c.SendStatus(fiber.StatusNotFound)
		}
		host.Fiber.Handler()(c.Context())
		return nil
	}
}
