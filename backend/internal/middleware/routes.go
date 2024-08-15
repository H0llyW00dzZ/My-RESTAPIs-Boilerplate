// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package middleware

import (
	"os"
	"runtime/debug"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"h0llyw00dz-template/backend/internal/database"
	log "h0llyw00dz-template/backend/internal/logger"
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/keyidentifier"
	"h0llyw00dz-template/backend/pkg/mime"
	"h0llyw00dz-template/backend/pkg/restapis/helper"
	"h0llyw00dz-template/env"
	htmx "h0llyw00dz-template/frontend/htmx/error_page_handler"

	_ "github.com/joho/godotenv/autoload" // godot autoload env
)

var (
	apiSubdomain = os.Getenv(env.APISUBDOMAIN) // set API_SUB_DOMAIN=api.localhost:8080 (depends of port available) for local development
)

// Note: This method works well Docs: https://github.com/gofiber/fiber/issues/750
// Also note that There is no limit to this feature. For example, you can add a billion domains or subdomains.
// Another note: When running this in a container with Kubernetes, make sure to have a configuration for allow internal IPs (e.g., 10.0.0.0/24).
// Because this method creates an additional internal IP for handling routes (e.g., 10.0.0.1 for REST APIs, then 10.0.0.2 for the frontend).
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
// It organizes routes into versioned groups for better API version management.
func RegisterRoutes(app *fiber.App, appName, monitorPath string, db database.Service) {
	// Note: This is just an example that can be integrated with other Fiber middleware.
	// If needed to store it in storage, use a prefix for group keys and call "GetKeyFunc".
	genReqID := keyidentifier.New(keyidentifier.Config{
		Prefix: "",
	})
	// Generate Request ID
	//
	// Note: This just example and "visitor_uuid" contextkey can be used for c.Locals
	// Previously I've been done implement this X-Request-ID bound into hash from TLS 1.3 with Private Protocols Cryptography (not open source) not UUID.
	xRequestID := NewRequestIDMiddleware(
		WithRequestIDHeaderContextKey("visitor_uuid"),
		WithRequestIDGenerator(genReqID.GetKey),
	)

	// Create a custom middleware to set the CSP header
	cspMiddleware := NewCSPHeaderGenerator()

	// Hosts
	// TODO: Reorganize this.
	// When this is reorganized, it will create 3 Routers (3 domains):
	//   1. Static Frontend:
	//      - HTMX
	//      - TEMPL
	//   2. Hostname:
	//      - Root Middleware Handler: When applying the middleware mechanism from registerRouteConfigMiddleware,
	//        it will be applied across the frontend and REST APIs.
	//      - Wildcard StatusServiceUnavailable Handler:
	//        Demo: https://api-beta.btz.pm speed might be slow at first due to the firewall implementation in Go that uses MySQL.
	//   3. REST APIs:
	//      - The Services
	// For TLS (e.g., in Ingress), it requires a Wildcard certificate instead of issuing 3 separate certificates.
	// The 3-router implementation is based on my previous work that has been done before.
	// It offers a better design (all-in-one) and is easier to maintain, even with 500+ files.
	hosts := map[string]*Host{}
	// Apply the combined middlewares
	registerRouteConfigMiddleware(app, db)
	// API subdomain
	api := fiber.New()
	// Register the REST APIs Routes
	registerRESTAPIsRoutes(app, db)
	// Note: This is just an example. In production, replace `api.localhost:8080` with a specific domain/subdomain, such as api.example.com.
	// Similarly, for the frontend, specify a domain like `hosts["example.com"] = &Host{frontend}`.
	// Additionally, instead of hard-coding the domain or subdomain,
	// it is possible to integrate it with environment variables or other configurations (e.g, YAML).
	hosts[apiSubdomain] = &Host{api}
	// Register the Static Frontend Routes
	registerStaticFrontendRoutes(app, appName, db)
	// Apply the subdomain routing middleware
	//
	// Note: "htmx.NewErrorHandler" will apply to localhost:8080 by default.
	// For "api.localhost:8080" to function correctly, REST API routes must be implemented.
	// Additionally, define environment variables for "DOMAIN" and "API_SUB_DOMAIN" to enable multi-site support (up to 1 billion domains).
	app.Use(xRequestID, cspMiddleware, htmx.NewErrorHandler, DomainRouter(hosts)) // When "htmx.NewErrorHandler" Applied, Generic Error (E.g, Crash/Panic will render "Internal Server Error" as JSON due It use recoverMiddleware)
}

// registerRouteConfigMiddleware applies middleware configurations to the Fiber application.
// It sets up the necessary middleware such as recovery, logging, and custom error handling for manipulating panics.
//
// Note: This is the root of the router configuration. When a Fiber middleware mechanism is applied here, it will be applied across both the frontend and the REST APIs.
// If there is a need to split the middleware configuration, it must be applied separately to the frontend and the REST APIs.
// If the root, frontend, and REST APIs configurations are still not enough, it can be implemented with own middleware configuration.
// This can lead to a complex setup, similar to the best art of binary trees (see https://en.wikipedia.org/wiki/Binary_tree).
// However, it's not actually complex; it's just the art of Go programming.
func registerRouteConfigMiddleware(app *fiber.App, db database.Service) {
	// Favicon front end setup
	// Note: this just an example
	favicon := NewFaviconMiddleware(
		WithFaviconFile("./frontend/public/assets/images/favicon.ico"),
		WithFaviconURL("/favicon.ico"),
	)
	// Note: This is just an example. It should work with SHA-256 for the key, however it may not properly bind to a UUID.
	cacheKeyGen := keyidentifier.New(keyidentifier.Config{
		Prefix: "go_frontend:",
	})
	// Speed depends of database connection as well.
	gopherstorage := db.FiberStorage()
	cacheMiddleware := NewCacheMiddleware(
		WithCacheStorage(gopherstorage),
		WithCacheKeyGenerator(cacheKeyGen.GenerateCacheKey),
		WithCacheExpiration(1*time.Hour),
		WithCacheControl(true),
		WithCacheNext(
			CustomNextContentType(
				// Note: Its important to disabling cache for this MIME
				fiber.MIMETextHTML,
				fiber.MIMETextHTMLCharsetUTF8,
				fiber.MIMEApplicationJSON,
				fiber.MIMEApplicationJSONCharsetUTF8,
				mime.ApplicationProblemJSON,
				mime.ApplicationProblemJSONCharsetUTF8,
				mime.TextEventStream,
			),
		),
		WithCacheHeader("X-Go-Frontend"),
	)

	// Recovery middleware setup
	// TODO: Move this into the server package because it should be initialized as the root before other functions.
	// This way, it can catch any panics, for example, catch any panic through the sub-package k8s/metrics.
	recoverMiddleware := recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e any) {
			// Log the panic and stack trace
			log.LogUserActivity(c, "attempted to panic occurred")
			log.LogCrashf(MsgPanicOccurred, e)
			log.LogCrashf(MsgStackTrace, debug.Stack())
		},
	})

	// Apply the recover middleware
	app.Use(helper.ErrorHandler, cacheMiddleware, recoverMiddleware, favicon)
}

// DomainRouter is a middleware function that handles subdomain or domain routing.
// It takes a map of subdomain or domain hosts and routes the request to the corresponding Fiber app.
//
// Note: This is useful for large Go applications, especially when running in Kubernetes,
// as it eliminates the need for multiple containers. It also supports integration with the Kubernetes ecosystem,
// such as pointing to CNAME/NS or manually (if not using Kubernetes).
// Also note that for TLS certificates, a wildcard/advanced certificate is required.
//
// Known Bugs:
//   - Wildcard/advanced certificates (e.g, issued by digicert, sectigo, google trust service, private ca) are not supported/compatible on Heroku.
//     Using a wildcard/advanced certificate on Heroku will cause an "SSL certificate error: There is conflicting information between the SSL connection, its certificate, and/or the included HTTP requests."
//     If using a wildcard/advanced certificate, it is recommended to deploy the application in a cloud environment such as Kubernetes, where you can easily control the ingress controller (e.g, Implement own such as universe).
//     Also note that regarding known bugs, it is not caused by this repository; it is an issue with Heroku's router.
//
// TODO: Consider moving this middleware into a separate package for better maintainability. This might involve creating a new repository.
func DomainRouter(hosts map[string]*Host) fiber.Handler {
	return func(c *fiber.Ctx) error {
		host := hosts[c.Hostname()]
		if host == nil {
			// Note: Returning a new error is a better approach instead of returning directly,
			// as it allows the error to be handled by the caller somewhere else in the codebase,
			// especially when the codebase grows larger.
			return fiber.NewError(fiber.StatusServiceUnavailable)
		}
		// Use c.Context() to pass the underlying context to the host's Fiber app.
		host.Fiber.Handler()(c.Context())
		return nil
	}
}
