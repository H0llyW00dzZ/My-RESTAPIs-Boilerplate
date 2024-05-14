// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package middleware

import (
	"h0llyw00dz-template/backend/internal/database"
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/hybrid"
	"h0llyw00dz-template/backend/pkg/restapis/server/health"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/google/uuid"
)

// APIRoute represents a single API route, containing the path, HTTP method,
// handler function, and an optional rate limiter.
type APIRoute struct {
	Path                      string
	Method                    string
	Handler                   fiber.Handler
	RateLimiter               fiber.Handler
	KeyAuth                   fiber.Handler
	RequestID                 fiber.Handler
	EncryptedCookieMiddleware fiber.Handler
}

// APIGroup represents a group of API routes under a common prefix.
// It also allows for a group-wide rate limiter.
type APIGroup struct {
	Prefix                    string
	Routes                    []APIRoute
	RateLimiter               fiber.Handler
	KeyAuth                   fiber.Handler
	RequestID                 fiber.Handler
	EncryptedCookieMiddleware fiber.Handler
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
	// Example encrypt cookie
	// Note: This is suitable with session middleware logic.
	encryptcookie := NewEncryptedCookieMiddleware(
		WithKey(encryptcookie.GenerateKey()),
		WithEncryptor(func(value, key string) (string, error) {
			return hybrid.EncryptCookie(value, key, "hex")
		}),
		WithDecryptor(func(encodedCookie, key string) (string, error) {
			return hybrid.DecryptCookie(encodedCookie, key, "hex")
		}),
	)

	v1 := api.Group("/v1", func(c *fiber.Ctx) error { // '/v1/' prefix
		c.Set("Version", "v1")
		// Set Cookie for group "v1" only if it doesn't exist
		// Note: This fix where a cookie keep generating from server-side into client-side.
		if c.Cookies("GhoperCookie") == "" {
			c.Cookie(&fiber.Cookie{
				// This should be safe against cookie poisoning, MITM, etc, even without a hash function,
				// because it would require 99999999999 cpu to attack this encryptcookie.
				Name:  "GhoperCookie",
				Value: uuid.NewSHA1(uuid.NameSpaceURL, []byte(c.IP())).String(),
			})
		}
		return c.Next()
	})

	// Create the root group and redirect middleware
	// Note: This is a method similar to nginx/apache .htaccess, if you're familiar with it.
	// This is just an example where it would redirect from api.localhost:8080/v1/ to api.localhost:8080.
	// In this root API group, it is possible to set the index root path `/` (e.g., to host the Swagger UI documentation).
	// Also, note that this method won't conflict with another path that already has a handler (e.g., api.localhost:8080/v1/server/health/db).
	rootGroup := APIGroup{
		Prefix:                    "/",
		EncryptedCookieMiddleware: encryptcookie,
		Routes:                    []APIRoute{},
	}

	// Note: This method is also called a "higher-order function",
	// similar to another configuration where "...interface{}" is used.
	// This is one of the reasons why I like Go. For example, when I'm lazy to implement something from scratch,
	// I can just use a package that is already stable then build on top of it using higher-order functions.
	redirectMiddleware := NewRedirectMiddleware(
		WithRules(map[string]string{
			"v1": "/",
		}),
		WithRedirectStatusCode(fiber.StatusMovedPermanently),
	)

	rootGroup.Routes = append(rootGroup.Routes, APIRoute{
		Path:    "v*",
		Method:  fiber.MethodGet,
		Handler: redirectMiddleware,
	})

	// Register the root group
	registerGroup(api, rootGroup)

	// Apply the rate limiter middleware directly to the REST API routes
	// Note: This method is called "higher-order function" which is better than (if-else statement which is bad)
	gopherStorage := db.FiberStorage()
	rateLimiterRESTAPIs := NewRateLimiter(
		WithStorage(gopherStorage),
		WithMax(maxRequestRESTAPIsRateLimiter),
		WithExpiration(maxExpirationRESTAPIsRateLimiter),
		WithLimitReached(ratelimiterMsg(MsgRESTAPIsVisitorGotRateLimited)),
	)

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
	// as I've had over 600 handlers across around 351 files.
	apiGroups := []APIGroup{
		{ // Note: Example https://localhost:8080/v1/server/health/db
			Prefix:      "/server/health",
			RateLimiter: rateLimiterRESTAPIs, // This is an optional example.
			Routes: []APIRoute{
				{
					Path: "/db",
					// Note: This approach allows defining multiple HTTP methods (e.g., GET, POST, PUT, DELETE) for a single handler & path.
					Method: strings.Join([]string{
						fiber.MethodGet,
					}, ","),
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

	registerGroupMiddlewares(g, group)
	registerGroupRoutes(g, group)
}

// registerGroupMiddlewares registers the middlewares for an API group.
func registerGroupMiddlewares(g fiber.Router, group APIGroup) {

	// Note: This approach uses a "higher-order function" called useNonNilMiddleware.
	// Also Note that Higher-order functions are powerful especially for "Cryptography Technique" and can handle multiple functions as arguments.
	// They provide a more concise and expressive way to work with functions compared to
	// using multiple if-else statements or switch cases.
	useNonNilMiddleware(
		g,
		group.RateLimiter,
		group.KeyAuth,
		group.RequestID,
		group.EncryptedCookieMiddleware,
	)
}

// registerGroupRoutes registers the routes for an API group.
func registerGroupRoutes(g fiber.Router, group APIGroup) {
	for _, route := range group.Routes {
		registerRoute(g, route)
	}
}

// registerRoute registers a single API route with its middlewares and handler.
func registerRoute(g fiber.Router, route APIRoute) {
	handlers := getRouteHandlers(route)
	g.Add(route.Method, route.Path, handlers...)
}

// getRouteHandlers returns the handlers for an API route.
func getRouteHandlers(route APIRoute) []fiber.Handler {
	handlers := make([]fiber.Handler, 0, 5)

	// Note: This approach uses a "higher-order function" called appendNonNilHandler.
	// Also Note that Higher-order functions are powerful especially for "Cryptography Technique" and can handle multiple functions as arguments.
	// They provide a more concise and expressive way to work with functions compared to
	// using multiple if-else statements or switch cases.
	handlers = appendNonNilHandler(
		handlers,
		route.RateLimiter,
		route.KeyAuth,
		route.RequestID,
		route.EncryptedCookieMiddleware,
	)

	handlers = append(handlers, route.Handler)

	return handlers
}
