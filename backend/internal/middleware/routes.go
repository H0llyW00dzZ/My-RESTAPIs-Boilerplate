// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"h0llyw00dz-template/backend/internal/database"
	log "h0llyw00dz-template/backend/internal/logger"
	"h0llyw00dz-template/backend/pkg/restapis/helper"
	htmx "h0llyw00dz-template/frontend/htmx/error_page_handler"
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
	// Generate Request ID
	//
	// Note: This just example and "visitor_uuid" contextkey can be used for c.Locals
	// Previously I've been done implement this X-Request-ID bound into hash from TLS 1.3 with Private Protocols Cryptography (not open source) not UUID.
	xRequestID := NewRequestIDMiddleware(
		WithRequestIDHeaderContextKey("visitor_uuid"),
	)

	// Create a custom middleware to set the CSP header
	//
	// Note: This safe and secure, not possible to spoof/bypass due it designed backend at top level,
	// even it's possible to try spoof/bypass, however browser will validate, and other Controller
	// (e.g, cloudflare for frontend, and nginx for controller load balancer) will validate as well.
	cspMiddleware := func(c *fiber.Ctx) error {
		clientIP := c.Get(log.CloudflareConnectingIPHeader)
		if clientIP == "" {
			clientIP = c.IP()
		}
		cloudflareRayID := c.Get(log.CloudflareRayIDHeader)
		if cloudflareRayID != "" {
			clientIP += " - Cloudflare detected - Ray ID: " + cloudflareRayID
		}
		countryCode := c.Get(log.CloudflareIPCountryHeader)
		if countryCode != "" {
			clientIP += ", Country: " + countryCode
		}

		// Digest a visitor
		digest := digest(clientIP)

		c.Locals("csp_random", digest)

		// Set the CSP header
		c.Set("Content-Security-Policy", fmt.Sprintf("script-src 'nonce-%s'", digest))

		// Continue to the next middleware/route handler
		return c.Next()
	}

	// Hosts
	hosts := map[string]*Host{}
	// Apply the combined middlewares
	registerRouteConfigMiddleware(app)
	// API subdomain
	api := fiber.New()
	// Register the REST APIs Routes
	registerRESTAPIsRoutes(app, db)
	// Note: This is just an example. In production, replace `api.localhost:8080` with a specific domain/subdomain, such as api.example.com.
	// Similarly, for the frontend, specify a domain like `hosts["example.com"] = &Host{frontend}`.
	// Additionally, instead of hard-coding the domain or subdomain,
	// it is possible to integrate it with environment variables or other configurations (e.g, YAML).
	hosts["api.localhost:8080"] = &Host{api}
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
func registerRouteConfigMiddleware(app *fiber.App) {
	// Favicon front end setup
	// Note: this just an example
	favicon := NewFaviconMiddleware(
		WithFaviconFile("./frontend/public/assets/images/favicon.ico"),
		WithFaviconURL("/favicon.ico"),
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
	app.Use(helper.ErrorHandler, recoverMiddleware, favicon)
}

// DomainRouter is a middleware function that handles subdomain or domain routing.
// It takes a map of subdomain or domain hosts and routes the request to the corresponding Fiber app.
//
// Note: This is useful for large Go applications, especially when running in Kubernetes,
// as it eliminates the need for multiple containers. It also supports integration with the Kubernetes ecosystem,
// such as pointing to CNAME/NS or manually (if not using Kubernetes).
// Also note that For TLS certificates, a wildcard/advanced certificate is required.
func DomainRouter(hosts map[string]*Host) fiber.Handler {
	return func(c *fiber.Ctx) error {
		host := hosts[c.Hostname()]
		if host == nil {
			// Note: Returning a new error is a better approach instead of returning directly,
			// as it allows the error to be handled by the caller somewhere else in the codebase,
			// especially when the codebase grows larger.
			return fiber.NewError(fiber.StatusNotFound)
		}
		// Use c.Context() to pass the underlying context to the host's Fiber app.
		host.Fiber.Handler()(c.Context())
		return nil
	}
}
