// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package middleware

import (
	"os"
	"runtime/debug"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"h0llyw00dz-template/backend/internal/database"
	log "h0llyw00dz-template/backend/internal/logger"
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/keyidentifier"
	"h0llyw00dz-template/backend/internal/middleware/router/domain"
	"h0llyw00dz-template/backend/pkg/mime"
	"h0llyw00dz-template/env"
	htmx "h0llyw00dz-template/frontend/htmx/error_page_handler"

	_ "github.com/joho/godotenv/autoload" // godot autoload env
)

var (
	// apiSubdomain is the subdomain for the API endpoints.
	// It is set using the API_SUB_DOMAIN environment variable.
	// Example: set API_SUB_DOMAIN=api.localhost:8080 for local development.
	apiSubdomain = os.Getenv(env.APISUBDOMAIN)

	// frontendDomain is the domain for the frontend application.
	// It is set using the DOMAIN environment variable.
	// Example: set DOMAIN=localhost:8080 for local development.
	frontendDomain = os.Getenv(env.DOMAIN)

	// readTimeoutStr is the read timeout duration for the server.
	// It is set using the READTIMEOUT environment variable.
	// The value should be a valid duration string (e.g., "30s", "1m").
	readTimeoutStr = os.Getenv(env.READTIMEOUT)

	// writeTimeoutStr is the write timeout duration for the server.
	// It is set using the WRITETIMEOUT environment variable.
	// The value should be a valid duration string (e.g., "30s", "1m").
	writeTimeoutStr = os.Getenv(env.WRITETIMEOUT)

	// readTimeout is the parsed read timeout duration.
	// It is obtained by parsing the readTimeoutStr using time.ParseDuration.
	// If parsing fails, the default value of 0 is used.
	readTimeout, _ = time.ParseDuration(readTimeoutStr)

	// writeTimeout is the parsed write timeout duration.
	// It is obtained by parsing the writeTimeoutStr using time.ParseDuration.
	// If parsing fails, the default value of 0 is used.
	writeTimeout, _ = time.ParseDuration(writeTimeoutStr)
)

// RegisterRoutes sets up the API routing for the application.
// It organizes routes into versioned groups for better API version management.
//
// Note: There are now 3 routers: restapis, frontend, and wildcard handler (503) (wildcard handler (503) known as root).
// They operate independently. Also note that as the codebase grows, the routing structure
// may become a binary tree (see https://en.wikipedia.org/wiki/Binary_tree), which is considered one of the best art in Go programming.
func RegisterRoutes(app *fiber.App, appName, monitorPath string, db database.Service) {
	// Hosts
	hosts := map[string]*fiber.App{}

	// restAPIs subdomain
	api := fiber.New(fiber.Config{
		ServerHeader: appName,
		AppName:      appName,
		// Note: Using the sonic JSON encoder/decoder provides better performance and is more memory-efficient
		// since Fiber is designed for zero allocation memory usage.
		JSONEncoder:      sonic.Marshal,
		JSONDecoder:      sonic.Unmarshal,
		CaseSensitive:    true,
		StrictRouting:    true,
		DisableKeepalive: false,
		ReadTimeout:      readTimeout,
		WriteTimeout:     writeTimeout,
		// Note: It's important to set Prefork to false because if it's enabled and running in Kubernetes,
		// it may get killed by an Out-of-Memory (OOM) error due to a conflict with the Horizontal Pod Autoscaler (HPA).
		Prefork: false,
		// Which is suitable for streaming AI Response.
		StreamRequestBody: true,
		// When running behind an ingress controller/proxy, disable "EnableIPValidation"
		// because the ingress controller/proxy will forward the real IP anyway from the header, which is already valid.
		EnableIPValidation:      false,
		EnableTrustedProxyCheck: true,
		// By default, it is set to 0.0.0.0/0 for local development; however, it can be bound to an ingress controller/proxy.
		// This can be a private IP range (e.g., 10.0.0.0/8).
		TrustedProxies: []string{"0.0.0.0/0"},
		// Trust X-Forwarded-For headers; additionally, this can be customized if using an ingress controller/proxy, especially Ingress Nginx.
		ProxyHeader: "X-Forwarded-For", // Fix where * (wildcard header) doesn't work in some kubernetes ingress eco-system
	})
	registerRESTAPIsRoutes(api, db)
	hosts[apiSubdomain] = api

	// Frontend domain
	frontend := fiber.New(fiber.Config{
		ServerHeader: appName,
		AppName:      appName,
		// Note: Using the sonic JSON encoder/decoder provides better performance and is more memory-efficient
		// since Fiber is designed for zero allocation memory usage.
		JSONEncoder:      sonic.Marshal,
		JSONDecoder:      sonic.Unmarshal,
		CaseSensitive:    true,
		StrictRouting:    true,
		DisableKeepalive: false,
		ReadTimeout:      readTimeout,
		WriteTimeout:     writeTimeout,
		// Note: It's important to set Prefork to false because if it's enabled and running in Kubernetes,
		// it may get killed by an Out-of-Memory (OOM) error due to a conflict with the Horizontal Pod Autoscaler (HPA).
		Prefork: false,
		// Which is suitable for streaming AI Response.
		StreamRequestBody: true,
		// When running behind an ingress controller/proxy, disable "EnableIPValidation"
		// because the ingress controller/proxy will forward the real IP anyway from the header, which is already valid.
		EnableIPValidation:      false,
		EnableTrustedProxyCheck: true,
		// By default, it is set to 0.0.0.0/0 for local development; however, it can be bound to an ingress controller/proxy.
		// This can be a private IP range (e.g., 10.0.0.0/8).
		TrustedProxies: []string{"0.0.0.0/0"},
		// Trust X-Forwarded-For headers; additionally, this can be customized if using an ingress controller/proxy, especially Ingress Nginx.
		ProxyHeader: "X-Forwarded-For", // Fix where * (wildcard header) doesn't work in some kubernetes ingress eco-system
	})
	registerStaticFrontendRoutes(frontend, appName, db)
	hosts[frontendDomain] = frontend

	// Apply the combined middlewares
	registerRouteConfigMiddleware(app, db)
	registerRootRouter(app)

	// Configure the DomainRouter middleware
	domainRouter := domain.New(domain.Config{
		Hosts: hosts,
	})

	// Apply the subdomain & domain routing middleware
	app.Use(domainRouter)
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
	// Note: This is just an example that can be integrated with other Fiber middleware.
	// If needed to store it in storage, use a prefix for group keys and call "GetKeyFunc".
	genReqID := keyidentifier.New(keyidentifier.Config{
		Prefix: "",
	})
	xRequestID := NewRequestIDMiddleware(
		WithRequestIDHeaderContextKey("visitor_uuid"),
		WithRequestIDGenerator(genReqID.GetKey),
	)
	// Note: This is just an example. It should work with SHA-256 for the key, however it may not properly bind to a UUID.
	cacheKeyGen := keyidentifier.New(keyidentifier.Config{
		Prefix: "go_frontend:",
	})
	// Speed depends of database connection as well.
	gopherstorage := db.FiberStorage()
	// stack Skipper
	contentTypeSkip := CustomNextContentType(
		// Note: Its important to disabling cache for this MIME
		fiber.MIMETextHTML,
		fiber.MIMETextHTMLCharsetUTF8,
		fiber.MIMEApplicationJSON,
		fiber.MIMEApplicationJSONCharsetUTF8,
		mime.ApplicationProblemJSON,
		mime.ApplicationProblemJSONCharsetUTF8,
		mime.TextEventStream,
		// This is temporary because it only registers 2 routers (currently).
		// When there are 3 or more routers, it will be structured like this in the demo:
		// - TLSv1.3 & HTTP/3 (QUIC): https://btz.pm (frontend currently disabled because I don't have any ideas for building the front-end, so it will return to the wildcard (see fiber.NewError in DomainRouter))
		// - TLSv1.3 & mTLSv1.3: https://api-beta.btz.pm (REST APIs)
		//
		// Also, note that the demo might be rare because having a single domain that can handle different protocols
		// to do one thing and do it well in the same host and repository is uncommon; however, it is secure.
		// fiber.MIMETextPlain,
		// fiber.MIMETextPlainCharsetUTF8,
		//
		// Note: It's important to disable caching for this MIME type, which is particularly suitable when using Grafana,
		// especially when playing HTMX MinesweeperX through Grafana plugins while monitoring.
		// Also note that while caching is disabled for Prometheus, it will become real-time because the Prometheus MIME type basically streams to serve HTTP.
		mime.PrometheusMetrics,
	)
	// Note: It's important to skip caching for redirect status codes, which can enhance security (e.g., for auth mechanisms).
	// If redirect status codes are cached, it can lead to security issues (e.g., new CVEs, exploits such as cache poisoning) because when redirect status codes are cached (hit),
	// they store only the header with 0 content (the reason why it appears blank in the browser). This can lead to security issues, especially for auth mechanisms, due to the information stored in the header with 0 content.
	statusCodeSkip := CustomNextStatusCode(
		fiber.StatusMovedPermanently,
		fiber.StatusPermanentRedirect,
		fiber.StatusTemporaryRedirect,
	)
	cacheMiddleware := NewCacheMiddleware(
		WithCacheStorage(gopherstorage),
		WithCacheKeyGenerator(cacheKeyGen.GenerateCacheKey),
		WithCacheExpiration(1*time.Hour),
		// Note: It is recommended to set "WithCacheControl" to true. If "WithCacheControl" is false, it will use server-side caching (in this repo),
		// which can waste memory resources because it stores everything on the server-side (in this repo) instead of the client-side.
		// When "WithCacheControl" is set to true, it can be combined with eTag. Additionally, by setting up caching in this way,
		// basically it creating own CDN (Content Delivery Network) solution.
		WithCacheControl(true),
		WithCacheNext(
			// Note: This actually work lmao.
			// Also, note that if it doesn't work, the browser would display a blank page
			// because it hits the cache, not an unreachable cache. If the cache is unreachable, it will redirect that mean works.
			CustomNextStack(map[string]func(*fiber.Ctx) bool{
				"contentTypeSkip": contentTypeSkip,
				"statusCodeSkip":  statusCodeSkip,
			}),
		),
		WithCacheHeader("X-Go-Frontend"),
	)

	// Create a custom middleware to set the CSP header
	cspMiddleware := NewCSPHeaderGenerator()

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

	etagMiddleware := NewETagMiddleware()
	// Apply the recover middleware
	app.Use(xRequestID, etagMiddleware, cspMiddleware, cacheMiddleware, htmx.NewErrorHandler, recoverMiddleware)
}

// registerRootRouter sets up the root router for the application.
// It registers static file serving and applies the favicon middleware.
func registerRootRouter(app *fiber.App) {
	// Register static file serving
	app.Static("/styles/", "./frontend/public/assets", fiber.Static{
		// This "ByteRange" Enhance QUIC
		ByteRange: true,
		Compress:  true,
		// optional
	})

	// Favicon setup
	// Note: This is just an example
	favicon := NewFaviconMiddleware(
		WithFaviconFile("./frontend/public/assets/images/favicon.ico"),
		WithFaviconURL("/favicon.ico"),
	)

	app.Use(favicon)
}
