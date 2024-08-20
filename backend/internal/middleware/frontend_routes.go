// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package middleware

import (
	"h0llyw00dz-template/backend/internal/database"
	htmx "h0llyw00dz-template/frontend/htmx/error_page_handler"

	"github.com/gofiber/fiber/v2"
)

// FrontendRoute represents a single frontend route, containing the path, HTTP method,
// handler function, and an optional rate limiter.
type FrontendRoute struct {
	Path                      string
	Method                    string
	Handler                   fiber.Handler
	RateLimiter               fiber.Handler
	CacheMiddleware           fiber.Handler
	ETagMiddleware            fiber.Handler
	FaviconMiddleware         fiber.Handler
	EncryptedCookieMiddleware fiber.Handler
}

// FrontendGroup represents a group of frontend routes under a common prefix.
// It also allows for a group-wide rate limiter.
type FrontendGroup struct {
	Prefix                    string
	Routes                    []FrontendRoute
	RateLimiter               fiber.Handler
	CacheMiddleware           fiber.Handler
	ETagMiddleware            fiber.Handler
	FaviconMiddleware         fiber.Handler
	EncryptedCookieMiddleware fiber.Handler
}

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
		ByteRange: true,
		Compress:  true,
		// Note: When running on K8S don't have to compress because it will handled by nginx or other controller.
		// Compress: true, // Optional
		// Note: It's important to disable this when using middleware cache to avoid confusion,
		// as caching is already handled by the middleware cache.
		// CacheDuration: -1 * time.Microsecond, // Optional
	})

	// As there is currently no handler, it will return a 404 error.
	// Example: http://localhost:8080/
	app.Use(htmx.NewErrorHandler)
}

// registerFrontendGroup adds all routes from a FrontendGroup to a specific Fiber app.
//
// Parameters:
//
//	app: The Fiber app on which to register the group's routes.
//	group: The FrontendGroup containing the routes to be registered.
func registerFrontendGroup(app *fiber.App, group FrontendGroup) {
	g := app.Group(group.Prefix)

	registerFrontendGroupMiddlewares(g, group)
	registerFrontendGroupRoutes(g, group)
}

// registerFrontendGroupMiddlewares registers the middlewares for a frontend group.
func registerFrontendGroupMiddlewares(g fiber.Router, group FrontendGroup) {
	useNonNilMiddleware(
		g,
		group.RateLimiter,
		group.CacheMiddleware,
		group.ETagMiddleware,
		group.FaviconMiddleware,
		group.EncryptedCookieMiddleware,
	)
}

// registerFrontendGroupRoutes registers the routes for a frontend group.
func registerFrontendGroupRoutes(g fiber.Router, group FrontendGroup) {
	for _, route := range group.Routes {
		registerFrontendRoute(g, route)
	}
}

// registerFrontendRoute registers a single frontend route with its middlewares and handler.
func registerFrontendRoute(g fiber.Router, route FrontendRoute) {
	handlers := getFrontendRouteHandlers(route)
	g.Add(route.Method, route.Path, handlers...)
}

// getFrontendRouteHandlers returns the handlers for a frontend route.
func getFrontendRouteHandlers(route FrontendRoute) []fiber.Handler {
	handlers := make([]fiber.Handler, 0, 6)

	handlers = appendNonNilHandler(
		handlers,
		route.RateLimiter,
		route.CacheMiddleware,
		route.ETagMiddleware,
		route.FaviconMiddleware,
		route.EncryptedCookieMiddleware,
	)

	handlers = append(handlers, route.Handler)

	return handlers
}
