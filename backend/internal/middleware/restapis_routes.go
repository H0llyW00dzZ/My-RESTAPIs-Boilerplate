// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package middleware

import (
	"h0llyw00dz-template/backend/internal/database"
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/hybrid"
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/keyidentifier"
	healthz "h0llyw00dz-template/backend/pkg/restapis/server/health"
	htmx "h0llyw00dz-template/frontend/htmx/error_page_handler"
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
	CompressJSON              fiber.Handler
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
	CompressJSON              fiber.Handler
	Prometheus                fiber.Handler
}

// APIConfig represents the configuration parameters for the API routes.
// It contains the necessary dependencies and settings for registering and handling API routes.
//
// TODO: Refactor this struct (Currently unused) to follow the Single Responsibility Principle (SRP).
// Consider splitting it into smaller, more focused structs or configurations.
type APIConfig struct {
	V1          fiber.Router
	V2          fiber.Router
	AppName     string
	MonitorPath string
	DB          database.Service
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
	// Note: This is just an example that can be integrated with other Fiber middleware.
	// If needed to store it in storage, use a prefix for group keys and call "GetKeyFunc".
	genReqID := keyidentifier.New(keyidentifier.Config{
		Prefix: "",
	})
	// Generate Request ID
	xRequestID := NewRequestIDMiddleware(
		WithRequestIDHeaderContextKey("rest_apis_visitor_uuid"),
		WithRequestIDGenerator(genReqID.GetKey),
	)
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

	//serverAPIs(v1, db, rateLimiterRESTAPIs)
	// Note: To test the Prometheus middleware, make a request to any URL in http://api.localhost:8080/ (restAPIS Router),
	// then visit http://api.localhost:8080/v1/server/metrics to see how it works.
	newPrometheus := NewPrometheusMiddleware(
		WithPrometheusServiceName("senior_golang"),
		WithPrometheusNamespace("restapis"),
		WithPrometheusSubsystem("http"),
		WithPrometheusLabels(map[string]string{
			"environment": "production",
		}),
		WithPrometheusSkipPaths([]string{
			"/health",
			"/v1/server/metrics",
			// Note: This should work because Prometheus can consume a lot of memory (it's surprising lmao since this repo is built on top of Fiber,
			// which is a zero-allocation framework built on top of fasthttp), which seems to be an issue with how Prometheus works and its implementation.
			// Ideally, the data should be stored in storage (e.g., disk) instead of memory. If this doesn't work, then use Next.
			"/favicon.ico",
		}),
		WithPrometheusMetricsPaths("/v1/server/metrics"),
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
	server := v1.Group("/server") // '/v1/server' prefix

	// Create the root group and redirect middleware
	// Note: This is a method similar to nginx/apache .htaccess, if you're familiar with it.
	// This is just an example where it would redirect from api.localhost:8080/v1/ to api.localhost:8080.
	// In this root API group, it is possible to set the index root path `/` (e.g., to host the Swagger UI documentation).
	// Also, note that this method won't conflict with another path that already has a handler (e.g., api.localhost:8080/v1/server/health/db).
	rootGroup := APIGroup{
		// When other Fiber middleware mechanisms are applied here, they will be applied across all REST API routes.
		// For example, if the "encryptcookie" middleware is applied here, it will encrypt any cookies sent in the REST APIs.
		// This is a tip to stack middleware mechanisms instead of applying them one by one.
		Prefix:                    "/",
		EncryptedCookieMiddleware: encryptcookie,
		Prometheus:                newPrometheus,
		RequestID:                 xRequestID,
		Routes:                    []APIRoute{},
	}

	// Note: This method is also called a "higher-order function",
	// similar to another configuration where "...any" is used.
	// This is one of the reasons why I like Go. For example, when I'm lazy to implement something from scratch,
	// I can just use a package that is already stable then build on top of it using higher-order functions.
	redirectMiddleware := NewRedirectMiddleware(
		// Note: This is a tip for manipulating bot scanners (bad crawls) that attempt to access sensitive directories like credentials or configs.
		// For example, if http://api.localhost:8080/v1/server/health/db is a registered endpoint, when a request is made to http://api.localhost:8080/v1/server/health/,
		// it will be redirected to http://api.localhost:8080/.
		WithRedirectRules(map[string]string{
			"v1":                "/",
			"v1/server/health":  "/",
			"v1/server/health/": "/",
			"v1/server":         "/",
			"v1/server/":        "/",
		}),
		WithRedirectStatusCode(fiber.StatusMovedPermanently),
	)

	methods := []string{
		fiber.MethodHead,
		fiber.MethodGet,
		fiber.MethodPost,
		fiber.MethodPut,
		fiber.MethodDelete,
		fiber.MethodPatch,
		fiber.MethodConnect,
		fiber.MethodOptions,
		fiber.MethodTrace,
	}

	for _, method := range methods {
		rootGroup.Routes = append(rootGroup.Routes, APIRoute{
			Path:    "*",
			Method:  method,
			Handler: redirectMiddleware,
		})
	}

	// Register the root group
	registerGroup(api, rootGroup)

	// Apply the rate limiter middleware directly to the REST API routes
	// Note: This method is called "higher-order function" which is better than (if-else statement which is bad)
	gopherStorage := db.FiberStorage()
	rateLimiterRESTAPIs := NewRateLimiter(
		WithRateLimiterStorage(gopherStorage),
		WithMax(maxRequestRESTAPIsRateLimiter),
		WithExpiration(maxExpirationRESTAPIsRateLimiter),
		WithLimitReached(ratelimiterMsg(MsgRESTAPIsVisitorGotRateLimited)),
	)
	server.Get("/health/db", rateLimiterRESTAPIs, encryptcookie, healthz.DBHandler(db))
	// This is for the Prometheus Handler. When an authentication mechanism is implemented, simply add the handler for the authentication mechanism here.
	// For example: server.Get("/metrics", rateLimiterRESTAPIs, keyAuth)
	//
	// Demo:
	//  - The currently unavailable service was stopped, as I was lazy to continue due to it being hosted on Heroku.
	//    The latency was unstable (always high) for handling two databases (Redis and MySQL), and memory leaks always occurred
	//    because Heroku stores metrics in memory (actually the same as this Prometheus Middleware), unlike on Kubernetes.
	server.Get("/metrics", rateLimiterRESTAPIs)

	// Register server APIs routes
	// Custom error handling for versioned APIs
	api.Use(htmx.NewErrorHandler)
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
//
// Note: This currently unused, might will removed it later.
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
					Handler: healthz.DBHandler(db),
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
		group.CompressJSON,
		group.Prometheus,
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
