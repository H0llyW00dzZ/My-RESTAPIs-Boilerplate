// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	log "h0llyw00dz-template/backend/internal/logger"
	"h0llyw00dz-template/backend/internal/middleware/csp"
	"h0llyw00dz-template/backend/internal/middleware/monitor"
	"h0llyw00dz-template/backend/internal/middleware/restime"

	validator "github.com/H0llyW00dzZ/FiberValidator"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/earlydata"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/redirect"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/rewrite"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/middleware/skip"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/google/uuid"
)

// NewCacheMiddleware creates a new cache middleware with optional custom configuration options.
//
// The cache middleware is built on top of the Fiber cache middleware and provides additional customization options.
// It allows you to specify a custom key generator function to generate cache keys based on the request context,
// as well as a custom cache skipper function to determine whether to skip caching for specific requests.
//
// Parameters:
//
//	options: Optional configuration options that can be used to customize the cache middleware.
//	         Available options include:
//	         - WithExpiration(expiration time.Duration): Sets the expiration time for cached entries.
//	         - WithCacheControl(cacheControl bool): Enables or disables the Cache-Control header.
//	         - WithKeyGenerator(keyGenerator func(*fiber.Ctx) string): Sets a custom key generator function.
//	         - WithNext(cacheSkipper func(*fiber.Ctx) bool): Sets a custom cache skipper function.
//	         - WithStorage(storage fiber.Storage): Sets the storage backend for the cache middleware.
//	         - WithStoreResponseHeaders(storeResponseHeaders bool): Enables or disables storing additional response headers.
//	         - WithMaxBytes(maxBytes uint): Sets the maximum number of bytes of response bodies to store in cache.
//	         - WithMethods(methods []string): Specifies the HTTP methods to cache.
//	         The options are passed as any and are type-asserted within the function.
//
// Returns:
//
//	A Fiber handler function representing the configured cache middleware.
//
// Example usage:
//
//	// Create a cache middleware with default options
//	cacheMiddleware := NewCacheMiddleware()
//
//	// Create a cache middleware with custom options
//	cacheMiddleware := NewCacheMiddleware(
//	    WithExpiration(time.Minute * 5),
//	    WithCacheControl(true),
//	    WithKeyGenerator(customKeyGenerator),
//	    WithNext(customCacheSkipper),
//	    WithStorage(customStorage),
//	    WithStoreResponseHeaders(true),
//	    WithMaxBytes(1024 * 1024),
//	    WithMethods([]string{fiber.MethodGet, fiber.MethodPost}),
//	)
func NewCacheMiddleware(options ...any) fiber.Handler {
	// Create a new cache middleware configuration.
	config := cache.Config{}

	// Apply any additional options to the cache configuration.
	for _, option := range options {
		if optFunc, ok := option.(func(*cache.Config)); ok {
			optFunc(&config)
		}
	}

	// Create the cache middleware with the configured options.
	cacheMiddleware := cache.New(config)

	// Return a custom middleware that conditionally applies the cache middleware.
	// Note: This safely integrates with the context (e.g., fiber ctx or std library ctx).
	return func(c *fiber.Ctx) error {
		// Check if caching should be skipped for the current request path.
		if config.Next != nil && config.Next(c) {
			// Caching is skipped, so don't generate a cache key and proceed to the next middleware.
			return c.Next()
		}

		// Caching is not skipped, so apply the cache middleware.
		return cacheMiddleware(c)
	}
}

// NewRateLimiter creates a new rate limiter middleware with optional custom configuration options.
// It retrieves the storage interface from the provided options and configures the
// rate limiter middleware accordingly.
//
// The rate limiter middleware is built on top of the Fiber rate limiter middleware and provides
// additional customization options. It allows you to specify the maximum number of requests allowed
// within a given time period, the expiration time for the rate limit, and a custom message to log
// when the rate limit is reached.
//
// Parameters:
//
//	options: Optional configuration options that can be used to customize the rate limiter middleware.
//	         Available options include:
//	         - WithMax(max int): Sets the maximum number of requests allowed within the time period.
//	         - WithExpiration(expiration time.Duration): Sets the expiration time for the rate limit.
//	         - WithLimitReached(handler fiber.Handler): Sets a custom handler to execute when the rate limit is reached.
//	         - WithStorage(storage fiber.Storage): Sets the storage backend for the rate limiter middleware.
//	         The options are passed as any and are type-asserted within the function.
//
// Returns:
//
//	A Fiber handler function representing the configured rate limiter middleware.
//
// Example usage:
//
//	// Create a rate limiter middleware with default options
//	rateLimiter := NewRateLimiter()
//
//	// Create a rate limiter middleware with a custom maximum number of requests
//	rateLimiter := NewRateLimiter(WithMax(100))
//
//	// Create a rate limiter middleware with a custom expiration time
//	rateLimiter := NewRateLimiter(WithExpiration(time.Minute))
//
//	// Create a rate limiter middleware with a custom limit reached handler
//	rateLimiter := NewRateLimiter(WithLimitReached(customLimitReachedHandler))
//
//	// Create a rate limiter middleware with a custom storage backend
//	rateLimiter := NewRateLimiter(WithStorage(customStorage))
//
//	// Create a rate limiter middleware with multiple custom options
//	rateLimiter := NewRateLimiter(WithMax(100), WithExpiration(time.Minute), WithLimitReached(customLimitReachedHandler), WithStorage(customStorage))
func NewRateLimiter(options ...any) fiber.Handler {
	// Create a new rate limiter middleware configuration.
	config := limiter.Config{}

	// Apply any additional options to the rate limiter configuration.
	for _, option := range options {
		if optFunc, ok := option.(func(*limiter.Config)); ok {
			optFunc(&config)
		}
	}

	// Create the rate limiter middleware with the configured options.
	rateLimiterMiddleware := limiter.New(config)

	// Return the rate limiter middleware.
	return rateLimiterMiddleware
}

// NewCORSMiddleware creates a new CORS middleware with optional custom configuration options.
//
// Example Usage:
//
//	corsMiddleware := NewCORSMiddleware(
//	WithAllowOrigins("https://example.com, https://api.example.com"),
//	WithAllowMethods("GET, POST, PUT, DELETE"),
//	WithAllowHeaders("Content-Type, Authorization"),
//	WithExposeHeaders("Content-Length"),
//	WithAllowCredentials(true),
//	WithMaxAge(3600),
//
// )
func NewCORSMiddleware(options ...any) fiber.Handler {
	// Note: In the Fiber framework v3, this CORS middleware configuration provides better security and low overhead.
	// For example, it allows blocking internal IPs by setting `AllowPrivateNetwork` to false (read more: https://docs.gofiber.io/api/middleware/cors).
	// Create a new CORS middleware configuration with default values
	config := cors.Config{}

	// Apply any additional options to the CORS configuration
	for _, option := range options {
		if optFunc, ok := option.(func(*cors.Config)); ok {
			optFunc(&config)
		}
	}

	// Create the CORS middleware with the configured options.
	corsMiddleware := cors.New(config)

	// Return the CORS middleware.
	return corsMiddleware
}

// NewETagMiddleware creates a new ETag middleware with the default and optional custom configuration options.
// It generates strong ETags for response caching and validation.
func NewETagMiddleware(options ...any) fiber.Handler {
	// Create a new ETag middleware configuration.
	config := etag.Config{}

	// Apply any additional options to the ETag configuration
	for _, option := range options {
		if optFunc, ok := option.(func(*etag.Config)); ok {
			optFunc(&config)
		}
	}

	// Create the ETag middleware with the configured options
	etagMiddleware := etag.New(config)

	// Return the ETag middleware
	return etagMiddleware
}

// NewFaviconMiddleware creates a new favicon middleware to serve a favicon file.
// It takes the file path of the favicon and the URL path where the favicon will be served.
func NewFaviconMiddleware(options ...any) fiber.Handler {
	// Create a new favicon middleware configuration with default values
	config := favicon.Config{}

	// Apply any additional options to the favicon configuration
	for _, option := range options {
		if optFunc, ok := option.(func(*favicon.Config)); ok {
			optFunc(&config)
		}
	}

	// Create the favicon middleware with the configured options
	faviconMiddleware := favicon.New(config)

	// Return the favicon middleware
	return faviconMiddleware
}

// NewPprofMiddleware creates a new pprof middleware with a custom path.
// It allows easy access to the pprof profiling tools and logs user activity.
func NewPprofMiddleware(path, pprofMessage string) fiber.Handler {
	// Example Usage: app.Use(NewPprofMiddleware("/pprof", "Accessed pprof profiling tools"))
	// Note: This safely integrates with the context (e.g., fiber ctx or std library ctx).
	return func(c *fiber.Ctx) error {
		log.LogUserActivity(c, pprofMessage)
		return pprof.New(pprof.Config{
			Prefix: path,
		})(c)
	}
}

// NewIPBasedUUIDMiddleware creates a new middleware that generates a deterministic UUID based on the client's IP address
// and attaches it to the Fiber context for reusability.
func NewIPBasedUUIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the client's IP address from the Fiber context
		ipAddress := c.IP()

		// Generate a deterministic UUID based on the IP address
		uuid := generateGoogleUUIDFromIP(ipAddress)

		// Attach the generated UUID to the Fiber context
		c.Locals("ip_based_uuid", uuid)

		// Continue to the next middleware or handler
		return c.Next()
	}
}

// NewSignatureMiddleware creates a new middleware that generates a signature based on the client's IP address
// and attaches it to the Fiber context for security purposes.
func NewSignatureMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the client's IP address from the Fiber context
		ipAddress := c.IP()

		// Generate a signature based on the IP address
		signature := generateSignatureFromIP(ipAddress)

		// Attach the generated signature to the Fiber context
		c.Locals("signature", signature)

		// Continue to the next middleware or handler
		return c.Next()
	}
}

// CustomKeyGenerator generates a custom cache key based on the request and logs the visitor activity.
func CustomKeyGenerator(c *fiber.Ctx) string {
	// Get client's IP and User-Agent
	clientIP := c.IP()
	userAgent := c.Get(fiber.HeaderUserAgent)

	// Create a string to hash
	toHash := fmt.Sprintf("%s-%s", clientIP, userAgent)

	// Create a fnv hash and write our string to it
	signature := hashForSignature(toHash)

	// Generate a UUID based on the hash
	signatureUUID := uuid.NewSHA1(uuid.NameSpaceOID, []byte(signature))

	// Log visitor activity with the signature for the frontend
	// Note: Rename "generated" to "initiated" because This cache is used for fiber storage operations (e.g., get, set, delete, reset, close).
	log.LogVisitorf("Frontend cache initiated for visitor activity: IP [%s], User-Agent [%s], Signature [%s], UUID [%s]", clientIP, userAgent, signature, signatureUUID.String())

	// Generate a custom cache key with the hashed signature and UUID for fiber storage operations (e.g., get, set, delete, reset, close).
	return fmt.Sprintf(utils.CopyString("cache_front_end:%s:%s:%s"), signature, signatureUUID.String(), c.Path())
}

// CustomCacheSkipper is a function that determines whether to skip caching for a given request path.
// It returns true if the request path starts with any of the specified prefixes.
func CustomCacheSkipper(prefixes ...string) func(*fiber.Ctx) bool {
	// Note: This safely integrates with the context (e.g., fiber ctx or std library ctx).
	return func(c *fiber.Ctx) bool {
		for _, prefix := range prefixes {
			if strings.HasPrefix(c.Path(), prefix) {
				return true
			}
		}
		return false
	}
}

// NewKeyAuthMiddleware creates a new key authentication middleware with the provided configuration.
//
// WARNING: Do not try to modify this by integrating it with the context (e.g., fiber ctx or std library ctx).
// Doing so may lead to high vulnerability if not handled correctly for each handler. It's better to keep it as is.
// For example (for advanced Go developers only), if you want to modify this to integrate it with the context (e.g., fiber ctx or std library ctx),
// each handler must have this function:
//
//	// Retrieve the authenticated API key from the request context
//
// apiKey, ok := c.Locals("token").(string)
//
//	if !ok {
//		log.LogUserActivity(c, "Invalid API key")
//		return helper.SendErrorResponse(c, fiber.StatusUnauthorized, "Missing or invalid API key")
//	}
//
// TODO: Implement a custom "Next" function that can skip authorization for admin/security roles,
// as they utilize another highly secure authentication mechanism with zero vulnerabilities and exploits 💀.
func NewKeyAuthMiddleware(options ...any) fiber.Handler {
	// Create a new key authentication middleware configuration.
	config := keyauth.Config{}

	// Apply any additional options to the key authentication configuration.
	for _, option := range options {
		if optFunc, ok := option.(func(*keyauth.Config)); ok {
			optFunc(&config)
		}
	}

	// Create the key authentication middleware with the configured options.
	keyAuthMiddleware := keyauth.New(config)

	// Return the key authentication middleware.
	return keyAuthMiddleware
}

// NewEncryptedCookieMiddleware creates a new encrypted cookie middleware with optional custom configuration options.
//
// Note: This middleware can be integrated with authentication cryptography techniques
// that use double encryption and decryption, such as the Last Enhance technique.
//
// WARNING: When using this middleware with custom cryptography that has already been implemented,
// make sure to use different keys for AES-GCM and ChaCha20-Poly1305 encryption for another.
// The current implementation of the Fiber encrypted cookie middleware only supports a single key,
// which is likely a mistake and a limitation.
// To enhance professional security, it's recommended to use separate keys in this function (e.g., create new keys specifically for cookies).
//
// TODO: Enhance this to integrate it with a perfect hybrid cryptosystem 🛡️🔐
// when the Fiber encrypted cookie supports multiple keys (as currently, it only supports one key)
func NewEncryptedCookieMiddleware(options ...any) fiber.Handler {
	// Create a new encrypted cookie middleware configuration.
	config := encryptcookie.Config{}

	// Apply any additional options to the encrypted cookie configuration.
	for _, option := range options {
		if optFunc, ok := option.(func(*encryptcookie.Config)); ok {
			optFunc(&config)
		}
	}

	// Create the encrypted cookie middleware with the configured options.
	encryptedCookieMiddleware := encryptcookie.New(config)

	// Return the encrypted cookie middleware.
	return encryptedCookieMiddleware
}

// NewRedirectMiddleware creates a new redirect middleware with optional custom configuration options.
//
// Example Usage:
//
//	redirectMiddleware := NewRedirectMiddleware(
//	  WithRedirectRules(map[string]string{
//	    "/old":   "/new",
//	    "/old/*": "/new/$1",
//	  }),
//	  WithRedirectStatusCode(fiber.StatusMovedPermanently),
//	)
func NewRedirectMiddleware(options ...any) fiber.Handler {
	// Create a new redirect configuration with default values
	config := redirect.Config{}

	// Apply any additional options to the redirect configuration.
	for _, option := range options {
		if optFunc, ok := option.(func(*redirect.Config)); ok {
			optFunc(&config)
		}
	}

	// Create the redirect middleware with the configured options
	redirectMiddleware := redirect.New(config)

	// Return the redirect middleware.
	return redirectMiddleware
}

// NewSessionMiddleware creates a new session middleware with optional custom configuration options.
//
// Note: When using this session with session storage,
// it is recommended to use a database that can handle high connections,
// for example, Redis is recommended because it can handle 10K++ connections,
// which is perfect for pooling without bottlenecks, and it's essentially unlimited connection.
func NewSessionMiddleware(options ...any) fiber.Handler {
	// Create a new session middleware configuration.
	config := session.Config{}

	// Default cleanup interval of 10 minutes.
	cleanupInterval := 10 * time.Minute

	// Default context key for storing the session.
	// Example Usage:
	// sessionMiddleware := NewSessionMiddleware(
	// 	WithSessionExpiration(time.Hour),
	// 	WithSessionStorage(customStorage),
	// 	"customSessionKey",
	// )
	contextKey := "session"

	// Apply any additional options to the session configuration.
	for _, option := range options {
		switch opt := option.(type) {
		case func(*session.Config):
			opt(&config)
		case time.Duration:
			cleanupInterval = opt
		case string:
			contextKey = opt
		}
	}

	// Create the session store with the configured options.
	store := session.New(config)

	// Start the cleanup goroutine for expired sessions.
	go CleanupExpiredSessions(store, cleanupInterval)

	// Return the session middleware function.
	return func(c *fiber.Ctx) error {
		// Get the session from the context.
		sess, err := store.Get(c)
		if err != nil {
			return err
		}

		// Save the session in the context using the custom context key for further usage.
		c.Locals(contextKey, sess)

		// Continue to the next middleware or handler.
		return c.Next()
	}
}

// NewCSRFMiddleware creates a new CSRF middleware with optional custom configuration options.
func NewCSRFMiddleware(options ...any) fiber.Handler {
	// Create a new CSRF middleware configuration.
	config := csrf.Config{}

	// Apply any additional options to the CSRF configuration.
	for _, option := range options {
		if optFunc, ok := option.(func(*csrf.Config)); ok {
			optFunc(&config)
		}
	}

	// Create the CSRF middleware with the configured options.
	csrfMiddleware := csrf.New(config)

	// Return the CSRF middleware.
	return csrfMiddleware
}

// NewHelmetMiddleware creates a new Helmet middleware with optional custom configuration options.
//
// Example Usage:
//
//	helmetMiddleware := NewHelmetMiddleware(
//	WithXSSProtection("0"),
//	WithContentTypeNosniff("nosniff"),
//	WithXFrameOptions("SAMEORIGIN"),
//	WithReferrerPolicy("no-referrer"),
//	WithCrossOriginEmbedderPolicy("require-corp"),
//	WithCrossOriginOpenerPolicy("same-origin"),
//	WithCrossOriginResourcePolicy("same-origin"),
//	WithOriginAgentCluster("?1"),
//	WithXDNSPrefetchControl("off"),
//	WithXDownloadOptions("noopen"),
//	WithXPermittedCrossDomain("none"),
//	)
//
// Note: This suitable for frontend.
func NewHelmetMiddleware(options ...any) fiber.Handler {
	// Create a new Helmet middleware configuration with default values
	config := helmet.Config{}

	// Apply the provided options to the Helmet configuration
	for _, option := range options {
		if optFunc, ok := option.(func(*helmet.Config)); ok {
			optFunc(&config)
		}
	}

	// Create the Helmet middleware with the configured options
	helmetMiddleware := helmet.New(config)

	// Return the Helmet middleware.
	return helmetMiddleware
}

// NewBasicAuthMiddleware creates a new basic authentication middleware with optional custom configuration options.
//
// TODO: Consider customizing this middleware to support alternative authentication methods like OAuth, cryptocurrency-based authentication,
// or Single Sign-On (SSO) by modifying the username/password logic to handle session tokens or other authentication mechanisms (NOTE: NO JWT and their base standards).
func NewBasicAuthMiddleware(options ...any) fiber.Handler {
	// Create a new basic authentication middleware configuration.
	config := basicauth.Config{}

	// Apply any additional options to the basic authentication configuration.
	for _, option := range options {
		if optFunc, ok := option.(func(*basicauth.Config)); ok {
			optFunc(&config)
		}
	}

	// Create the basic authentication middleware with the configured options.
	basicAuthMiddleware := basicauth.New(config)

	// Return the basic authentication middleware.
	return basicAuthMiddleware
}

// NewSwaggerMiddleware creates a new Swagger middleware with optional custom configuration options.
func NewSwaggerMiddleware(options ...any) fiber.Handler {
	// Create a new Swagger middleware configuration.
	config := swagger.Config{}

	// Apply any additional options to the Swagger configuration.
	for _, option := range options {
		if optFunc, ok := option.(func(*swagger.Config)); ok {
			optFunc(&config)
		}
	}

	// Create the Swagger middleware with the configured options.
	swaggerMiddleware := swagger.New(config)

	// Return the Swagger middleware.
	return swaggerMiddleware
}

// NewIdempotencyMiddleware creates a new idempotency middleware with optional custom configuration options.
//
// Ref: https://datatracker.ietf.org/doc/html/draft-ietf-httpapi-idempotency-key-header-02
func NewIdempotencyMiddleware(options ...any) fiber.Handler {
	// Create a new idempotency middleware configuration.
	config := idempotency.Config{}

	// Apply any additional options to the idempotency configuration.
	for _, option := range options {
		if optFunc, ok := option.(func(*idempotency.Config)); ok {
			optFunc(&config)
		}
	}

	// Create the idempotency middleware with the configured options.
	idempotencyMiddleware := idempotency.New(config)

	// Return the idempotency middleware.
	return idempotencyMiddleware
}

// NewRewriteMiddleware creates a new Rewrite middleware with optional custom configuration options.
func NewRewriteMiddleware(options ...any) fiber.Handler {
	// Create a new Rewrite middleware configuration
	config := rewrite.Config{}

	// Apply any additional options to the Rewrite configuration
	for _, option := range options {
		if optFunc, ok := option.(func(*rewrite.Config)); ok {
			optFunc(&config)
		}
	}

	// Create the Rewrite middleware with the configured options
	rewriteMiddleware := rewrite.New(config)

	// Return the Rewrite middleware
	return rewriteMiddleware
}

// NewHTTPHandlerMiddleware creates a new middleware that adapts an http.Handler to a Fiber handler.
func NewHTTPHandlerMiddleware(handler http.Handler) fiber.Handler {
	return adaptor.HTTPHandler(handler)
}

// NewHTTPHandlerFuncMiddleware creates a new middleware that adapts an http.HandlerFunc to a Fiber handler.
func NewHTTPHandlerFuncMiddleware(handlerFunc http.HandlerFunc) fiber.Handler {
	return adaptor.HTTPHandlerFunc(handlerFunc)
}

// NewHTTPMiddlewareMiddleware creates a new middleware that adapts an http.Handler middleware to a Fiber middleware.
func NewHTTPMiddlewareMiddleware(mw func(http.Handler) http.Handler) fiber.Handler {
	return adaptor.HTTPMiddleware(mw)
}

// NewFiberHandlerMiddleware creates a new http.Handler that adapts a Fiber handler.
func NewFiberHandlerMiddleware(handler fiber.Handler) http.Handler {
	return adaptor.FiberHandler(handler)
}

// NewFiberHandlerFuncMiddleware creates a new http.HandlerFunc that adapts a Fiber handler.
func NewFiberHandlerFuncMiddleware(handler fiber.Handler) http.HandlerFunc {
	return adaptor.FiberHandlerFunc(handler)
}

// NewFiberAppMiddleware creates a new http.HandlerFunc that adapts a Fiber application.
func NewFiberAppMiddleware(app *fiber.App) http.HandlerFunc {
	return adaptor.FiberApp(app)
}

// ConvertRequestMiddleware converts a Fiber context to an http.Request.
// It allows specifying a custom context key for storing the converted request.
// If no custom context key is provided, it defaults to using "http_request" as the key.
func ConvertRequestMiddleware(forServer bool, contextKey ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req, err := adaptor.ConvertRequest(c, forServer)
		if err != nil {
			return err
		}

		key := "http_request"
		if len(contextKey) > 0 {
			key = contextKey[0]
		}

		c.Locals(key, req)
		return c.Next()
	}
}

// NewValidatorMiddleware creates a new Validator middleware with the provided configuration.
//
// Note: This is still a work in progress and might not work as expected for MIME types other than JSON.
// The reason for implementing it here is to avoid having to add it again.
func NewValidatorMiddleware(options ...any) fiber.Handler {
	// Create a new Validator middleware configuration.
	config := validator.Config{}

	// Apply any additional options to the Validator configuration.
	for _, option := range options {
		if optFunc, ok := option.(func(*validator.Config)); ok {
			optFunc(&config)
		}
	}

	// Create the Validator middleware with the configured options.
	validatorMiddleware := validator.New(config)

	// Return the Validator middleware.
	return validatorMiddleware
}

// NewRequestIDMiddleware creates a new request ID middleware with optional custom configuration options.
//
// The request ID middleware generates a unique ID for each incoming request and adds it to the response headers.
// It can be used for tracking and debugging purposes.
func NewRequestIDMiddleware(options ...any) fiber.Handler {
	// Create a new request ID middleware configuration.
	config := requestid.Config{}

	// Apply any additional options to the request ID configuration.
	for _, option := range options {
		if optFunc, ok := option.(func(*requestid.Config)); ok {
			optFunc(&config)
		}
	}

	// Create the request ID middleware with the configured options.
	requestIDMiddleware := requestid.New(config)

	// Return the request ID middleware.
	return requestIDMiddleware
}

// NewHealthZCheck creates a new HealthZ Check middleware for the Fiber web framework.
//
// The HealthZ Check middleware is used to perform health checks on the application and its dependencies.
// It is particularly useful in cloud environments to ensure the application is functioning correctly and
// to monitor its health status.
func NewHealthZCheck(options ...any) fiber.Handler {
	// Create a new health check middleware configuration.
	config := healthcheck.Config{}

	// Apply any additional options to the HealthZ Check configuration.
	for _, option := range options {
		if optFunc, ok := option.(func(*healthcheck.Config)); ok {
			optFunc(&config)
		}
	}

	// Create the health check middleware with the configured options.
	healthzCheckMiddleware := healthcheck.New(config)

	// Return the health check middleware.
	return healthzCheckMiddleware
}

// NewCSPHeaderGenerator creates a new Content Security Policy (CSP) header generator middleware for the Fiber web framework.
//
// The CSP header generator middleware is used to generate and set the Content Security Policy headers in the HTTP response.
// It helps to prevent cross-site scripting (XSS), clickjacking, and other code injection attacks by specifying the allowed
// sources of content and restricting the behavior of web pages.
func NewCSPHeaderGenerator(options ...any) fiber.Handler {
	// Create a new CSP configuration with default values.
	config := csp.DefaultConfig()

	// Apply any additional options to the CSP configuration.
	for _, option := range options {
		if optFunc, ok := option.(func(*csp.Config)); ok {
			optFunc(&config)
		}
	}

	// Create a new CSP header generator with the configured options.
	cspgenerator := csp.New(config)

	// Return the CSP header generator middleware.
	return cspgenerator
}

// NewSkipMiddleware creates a new middleware handler that skips the wrapped handler
// based on the provided exclude predicate function.
//
// The exclude predicate function takes a [*fiber.Ctx] as input and returns a boolean value.
// If the exclude predicate returns true, the wrapped handler will be skipped, and the request
// will be passed to the next middleware in the chain using [c.Next].
//
// If the exclude predicate returns false, the wrapped handler will be executed.
func NewSkipMiddleware(handler fiber.Handler, exclude func(c *fiber.Ctx) bool) fiber.Handler {
	// Use the [skip.New] function from the Fiber framework to create a new skip middleware.
	return skip.New(handler, exclude)
}

// NewCompressMiddleware creates a new compression middleware for the Fiber web framework.
//
// The compression middleware is used to compress the HTTP response body using the specified compression algorithm.
// It helps to reduce the size of the response, resulting in faster transfer times and reduced bandwidth usage.
//
// Note: This is suitable for enhancing the load balancer's performance when using HTTP/3.
func NewCompressMiddleware(options ...any) fiber.Handler {
	// Create a new compression configuration with default values.
	config := compress.Config{}

	// Apply any additional options to the compression configuration.
	for _, option := range options {
		if optFunc, ok := option.(func(*compress.Config)); ok {
			optFunc(&config)
		}
	}

	// Create a new compression middleware with the configured options.
	compressMiddleware := compress.New(config)

	// Return the compression middleware.
	return compressMiddleware
}

// NewResponseTime creates a new response time middleware with optional custom configuration options.
func NewResponseTime(options ...any) fiber.Handler {
	// Create a new response time configuration.
	config := restime.Config{}

	// Apply any additional options to the response time configuration.
	for _, option := range options {
		if optFunc, ok := option.(func(*restime.Config)); ok {
			optFunc(&config)
		}
	}

	// Create the response time middleware with the configured options.
	restimeMiddleware := restime.New(config)

	// Return the response time middleware.
	return restimeMiddleware
}

// NewEarlyData (Another QUIC) creates a new early data middleware with optional custom configuration options.
//
// Note: Consider carefully when using an external load balancer/ingress (e.g., nginx) due to Fiber's built-in support for
// Another QUIC, which is already safe and requires explicit TLSv1.3 for the external load balancer/ingress.
func NewEarlyData(options ...any) fiber.Handler {
	// Create a new early data configuration.
	config := earlydata.Config{}

	// Apply any additional options to the early data configuration.
	for _, option := range options {
		if optFunc, ok := option.(func(*earlydata.Config)); ok {
			optFunc(&config)
		}
	}

	// Create the early data middleware with the configured options.
	earlyDataMiddleware := earlydata.New(config)

	// Return the early data middleware.
	return earlyDataMiddleware
}

// NewPrometheusMiddleware creates a new Prometheus middleware with optional custom configuration options.
//
// Note: While this Prometheus middleware is highly flexible and portable, it is easy to enhance security by adding an authentication mechanism.
//
// Recommended for authentication:
//   - https://docs.gofiber.io/api/middleware/keyauth + Database
//
// Demo:
//   - https://api-beta.btz.pm/v1/server/metrics (Authentication required)
func NewPrometheusMiddleware(options ...any) fiber.Handler {
	// Create a new Prometheus configuration.
	config := monitor.PrometheusConfig{}

	// Apply any additional options to the Prometheus configuration.
	for _, option := range options {
		if optFunc, ok := option.(func(*monitor.PrometheusConfig)); ok {
			optFunc(&config)
		}
	}

	// Create the Prometheus middleware using the NewPrometheus function.
	prometheusMiddleware := monitor.NewPrometheus(config)

	// Return the Prometheus middleware.
	return prometheusMiddleware
}
