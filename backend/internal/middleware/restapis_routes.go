// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package middleware

import (
	"h0llyw00dz-template/backend/internal/database"
	"h0llyw00dz-template/backend/pkg/restapis/server/health"

	"github.com/gofiber/fiber/v2"
)

// APIRoute represents a single API route, containing the path, HTTP method,
// handler function, and an optional rate limiter.
type APIRoute struct {
	Path        string
	Method      string
	Handler     fiber.Handler
	RateLimiter fiber.Handler
}

// APIGroup represents a group of API routes under a common prefix.
// It also allows for a group-wide rate limiter.
type APIGroup struct {
	Prefix      string
	Routes      []APIRoute
	RateLimiter fiber.Handler
}

// registerRESTAPIsRoutes registers the REST API routes for the application.
// It creates a version group ('/v1'), applies rate limiting middleware,
// and calls the serverAPIs function to register server-related API routes.
//
// The function follows the idiomatic Go practices of using descriptive names,
// handling errors, and using concise and readable code.
//
// Parameters:
//
//	api: The Fiber router to register the routes on.
//	db: The database service to be used by the API handlers.
func registerRESTAPIsRoutes(api fiber.Router, db database.Service) {
	v1 := api.Group("/v1", func(c *fiber.Ctx) error { // '/v1/' prefix
		c.Set("Version", "v1")
		return c.Next()
	})

	// Apply the rate limiter middleware directly to the REST API routes
	rateLimiterRESTAPIs := NewRateLimiter(db, maxRequestRESTAPIsRateLimiter, maxExpirationRESTAPIsRateLimiter, MsgRESTAPIsVisitorGotRateLimited)

	// Register server APIs routes
	serverAPIs(v1, db, rateLimiterRESTAPIs)
}

// serverAPIs registers the server-related API routes.
// It defines the API groups and routes, and registers them using the registerGroup function.
//
// The function is organized in a clear and concise manner, making it easy to understand
// and maintain. It follows the idiomatic Go practices of using descriptive names,
// handling errors, and using a slice to store related data (API groups).
//
// Parameters:
//
//	v1: The Fiber router to register the routes on (version 1 in this case).
//	db: The database service to be used by the API handlers.
//	rateLimiterRESTAPIs: The rate limiter middleware to be applied to the API routes.
func serverAPIs(v1 fiber.Router, db database.Service, rateLimiterRESTAPIs fiber.Handler) {
	// Define the API groups and routes
	// Note: By refactoring like this, it allows for an unlimited number of handlers and easy maintainability,
	// as I've had over 500 handlers across around 250 files.
	apiGroups := []APIGroup{
		{ // Note: Example https://localhost:8080/v1/server/health/db
			Prefix:      "/server/health",
			RateLimiter: rateLimiterRESTAPIs, // This is an optional example.
			Routes: []APIRoute{
				{
					Path:    "/db",
					Method:  fiber.MethodGet,
					Handler: health.DBHandler(db),
				},
			},
		},
	}

	// Register the API routes for each version
	for _, group := range apiGroups {
		registerGroup(v1, group)
	}
}

// registerGroup adds all routes from an APIGroup to a specific Fiber router.
//
// Parameters:
//
//	router: The Fiber router on which to register the group's routes.
//	group: The APIGroup containing the routes to be registered.
func registerGroup(router fiber.Router, group APIGroup) {
	g := router.Group(group.Prefix)

	if group.RateLimiter != nil {
		g.Use(group.RateLimiter)
	}

	for _, route := range group.Routes {
		if route.RateLimiter != nil {
			g.Add(route.Method, route.Path, route.RateLimiter, route.Handler)
		} else {
			g.Add(route.Method, route.Path, route.Handler)
		}
	}
}
