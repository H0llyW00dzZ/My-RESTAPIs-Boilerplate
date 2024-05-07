// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package middleware

import (
	"fmt"
	"h0llyw00dz-template/backend/internal/database"
	"strings"
	"time"

	log "h0llyw00dz-template/backend/internal/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/redirect"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/google/uuid"
)

// NewCacheMiddleware creates a new cache middleware with the specified expiration time, cache control flag,
// and optional custom configuration options. It retrieves the Redis storage interface from the provided database
// service and configures the cache middleware accordingly.
//
// The cache middleware is built on top of the Fiber cache middleware and provides additional customization options.
// It allows you to specify a custom key generator function to generate cache keys based on the request context,
// as well as a custom cache skipper function to determine whether to skip caching for specific requests.
//
// Parameters:
//
//	db: The database service instance that provides the Redis storage interface.
//	expiration: The expiration time for cached entries.
//	cacheControl: A boolean flag indicating whether to include cache control headers in the response.
//	options: Optional configuration options that can be used to customize the cache middleware.
//	         Available options include:
//	         - WithKeyGenerator(keyGenerator func(*fiber.Ctx) string): Sets a custom key generator function.
//	         - WithCacheSkipper(cacheSkipper func(*fiber.Ctx) bool): Sets a custom cache skipper function.
//	         The options are passed as interface{} and are type-asserted within the function.
//
// Returns:
//
//	A Fiber handler function representing the configured cache middleware.
//
// Example usage:
//
//	// Create a cache middleware with default options
//	cacheMiddleware := NewCacheMiddleware(db, expiration, cacheControl)
//
//	// Create a cache middleware with a custom key generator
//	cacheMiddleware := NewCacheMiddleware(db, expiration, cacheControl, WithKeyGenerator(customKeyGenerator))
//
//	// Create a cache middleware with a custom cache skipper
//	cacheMiddleware := NewCacheMiddleware(db, expiration, cacheControl, WithCacheSkipper(customCacheSkipper))
//
//	// Create a cache middleware with both custom key generator and cache skipper
//	cacheMiddleware := NewCacheMiddleware(db, expiration, cacheControl, WithKeyGenerator(customKeyGenerator), WithCacheSkipper(customCacheSkipper))
func NewCacheMiddleware(db database.Service, expiration time.Duration, cacheControl bool, options ...interface{}) fiber.Handler {
	// Retrieve the Redis storage interface from the database service.
	cacheMiddlewareService := db.FiberStorage()

	// Create a new cache middleware configuration.
	config := cache.Config{
		Expiration:   expiration,
		CacheControl: cacheControl,
		Storage:      cacheMiddlewareService,
	}

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
//	         The options are passed as interface{} and are type-asserted within the function.
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
func NewRateLimiter(options ...interface{}) fiber.Handler {
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
func NewCORSMiddleware(options ...CORSOption) fiber.Handler {
	// Note: In the Fiber framework v3, this CORS middleware configuration provides better security and low overhead.
	// For example, it allows blocking internal IPs by setting `AllowPrivateNetwork` to false (read more: https://docs.gofiber.io/api/middleware/cors).
	// Create a new CORS middleware configuration with default values
	config := cors.Config{}

	// Apply any additional options to the CORS configuration
	for _, option := range options {
		option(&config)
	}

	// Create the CORS middleware with the configured options
	return cors.New(config)
}

// NewETagMiddleware creates a new ETag middleware with the default configuration.
// It generates strong ETags for response caching and validation.
func NewETagMiddleware() fiber.Handler {
	return etag.New(etag.Config{
		Weak: false,
	})
}

// NewFaviconMiddleware creates a new favicon middleware to serve a favicon file.
// It takes the file path of the favicon and the URL path where the favicon will be served.
func NewFaviconMiddleware(filePath, urlPath string) fiber.Handler {
	return favicon.New(favicon.Config{
		File: filePath,
		URL:  urlPath,
	})
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
func NewKeyAuthMiddleware(db database.Service, options ...func(*keyauth.Config)) fiber.Handler {
	// Create a new key authentication middleware configuration.
	config := keyauth.Config{
		KeyLookup:  "header:" + fiber.HeaderAuthorization,
		AuthScheme: "Bearer",
		ContextKey: "token",
	}

	// Apply any additional options to the key authentication configuration.
	for _, option := range options {
		option(&config)
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
func NewEncryptedCookieMiddleware(options ...interface{}) fiber.Handler {
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
func NewRedirectMiddleware(options ...func(*redirect.Config)) fiber.Handler {
	// Create a new redirect configuration with default values
	config := redirect.Config{}

	// Apply the provided options to the redirect configuration
	for _, option := range options {
		option(&config)
	}

	// Create the redirect middleware with the configured options
	return redirect.New(config)
}

// NewSessionMiddleware creates a new session middleware with optional custom configuration options.
func NewSessionMiddleware(options ...interface{}) fiber.Handler {
	// Create a new session middleware configuration.
	config := session.Config{}

	// Default cleanup interval of 10 minutes.
	cleanupInterval := 10 * time.Minute

	// Apply any additional options to the session configuration.
	for _, option := range options {
		switch opt := option.(type) {
		case func(*session.Config):
			opt(&config)
		case time.Duration:
			cleanupInterval = opt
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

		// Save the session in the context for further usage.
		c.Locals("session", sess)

		// Continue to the next middleware or handler.
		return c.Next()
	}
}

// NewCSRFMiddleware creates a new CSRF middleware with optional custom configuration options.
func NewCSRFMiddleware(options ...interface{}) fiber.Handler {
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
func NewHelmetMiddleware(options ...interface{}) fiber.Handler {
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
