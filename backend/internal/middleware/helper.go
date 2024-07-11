// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package middleware

import (
	"crypto/sha256"
	"fmt"
	log "h0llyw00dz-template/backend/internal/logger"
	"h0llyw00dz-template/backend/pkg/restapis/helper"
	"hash/fnv"
	"time"

	validator "github.com/H0llyW00dzZ/FiberValidator"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/redirect"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/rewrite"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
)

// generateGoogleUUIDFromIP generates a deterministic UUID based on the provided IP address.
func generateGoogleUUIDFromIP(ipAddress string) string {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(ipAddress)).String()
}

// generateSignatureFromIP generates a signature based on the provided IP address.
func generateSignatureFromIP(ipAddress string) string {
	// Generate a UUID based on the IP address
	uuid := uuid.NewSHA1(uuid.NameSpaceURL, []byte(ipAddress))

	// Generate a signature by taking the first 8 characters of the UUID
	signature := uuid.String()[:8]

	return signature
}

// hashForSignature creates a hash from the IP and User-Agent to use in generating a UUID.
func hashForSignature(toHash string) string {
	h := fnv.New64a()
	h.Write([]byte(toHash))
	return fmt.Sprintf("%x", h.Sum64())
}

// ratelimiterMsg is a custom handler function for the rate limiter middleware.
// It returns a closure that logs a message indicating that a visitor has been rate limited
// and sends an error response with a "Too Many Requests" status code and an appropriate error message.
//
// Parameters:
//
//	customMessage: A custom message to be logged when a visitor is rate limited.
//
// Returns:
//
//	A closure that takes a Fiber context and returns an error.
//	The closure logs the custom message and sends an error response indicating that the rate limit has been reached.
func ratelimiterMsg(customMessage string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		log.LogUserActivity(c, customMessage)
		// Note: Custom messages for the "Too Many Requests" HTTP response will not be implemented due to lazy 🤪.
		// The default error message provided by Fiber will be used instead.
		return helper.SendErrorResponse(c, fiber.StatusTooManyRequests, fiber.ErrTooManyRequests.Message)
	}
}

// WithKeyGenerator is an option function that sets a custom key generator for various Fiber middlewares.
//
// It supports the following middleware configurations:
//
//	*cache.Config: Sets the key generator for the cache middleware.
//	*limiter.Config: Sets the key generator for the rate limiter middleware.
//
// The key generator is a function that takes a *fiber.Ctx as a parameter and returns a string key.
// It is used to generate a unique key for each request, which is used for caching or rate limiting.
//
// Example usage:
//
//	// Define a custom key generator function
//	func customKeyGenerator(c *fiber.Ctx) string {
//	    // Custom logic to generate a unique key based on the request
//	    return c.Path() + c.IP()
//	}
//
//	// Use the WithKeyGenerator option function to set the key generator for the cache middleware
//	cacheMiddleware := NewCacheMiddleware(WithKeyGenerator(customKeyGenerator))
//
//	// Use the WithKeyGenerator option function to set the key generator for the rate limiter middleware
//	rateLimiterMiddleware := NewRateLimiter(WithKeyGenerator(customKeyGenerator))
//
// Note:
//   - If an unsupported middleware configuration is passed to WithKeyGenerator, it will panic with an error message.
//   - Additional key generator support for other middlewares will be added based on future requirements.
func WithKeyGenerator(keyGenerator func(*fiber.Ctx) string) any {
	return func(config any) {
		// Note: This a better switch-statement, it doesn't matter if there is so many switch (e.g, 1 billion switch case)
		switch cfg := config.(type) {
		case *cache.Config:
			cfg.KeyGenerator = keyGenerator
			// TODO: Implement a custom key generator for any sensitive data such as API keys or OAuth tokens,
			// since the default rate limiter key in Fiber is based on c.IP()
		case *limiter.Config:
			cfg.KeyGenerator = keyGenerator
		default:
			panic(fmt.Sprintf("unsupported config type: %T", config))
		}
	}
}

// WithCacheControl is an option function for NewCacheMiddleware that enables or disables the Cache-Control header.
func WithCacheControl(cacheControl bool) func(*cache.Config) {
	return func(config *cache.Config) {
		config.CacheControl = cacheControl
	}
}

// WithCacheExpiration is an option function for NewCacheMiddleware that sets the expiration time for cached entries.
func WithCacheExpiration(expiration time.Duration) func(*cache.Config) {
	return func(config *cache.Config) {
		config.Expiration = expiration
	}
}

// WithCacheHeader is an option function for NewCacheMiddleware that sets the cache status header.
func WithCacheHeader(header string) func(*cache.Config) {
	return func(config *cache.Config) {
		config.CacheHeader = header
	}
}

// WithCacheExpirationGenerator is an option function for NewCacheMiddleware that sets a custom expiration generator function.
func WithCacheExpirationGenerator(generator func(*fiber.Ctx, *cache.Config) time.Duration) func(*cache.Config) {
	return func(config *cache.Config) {
		config.ExpirationGenerator = generator
	}
}

// WithCacheStoreResponseHeaders is an option function for NewCacheMiddleware that enables or disables storing response headers.
func WithCacheStoreResponseHeaders(storeHeaders bool) func(*cache.Config) {
	return func(config *cache.Config) {
		config.StoreResponseHeaders = storeHeaders
	}
}

// WithCacheMaxBytes is an option function for NewCacheMiddleware that sets the maximum number of bytes to store in cache.
func WithCacheMaxBytes(maxBytes uint) func(*cache.Config) {
	return func(config *cache.Config) {
		config.MaxBytes = maxBytes
	}
}

// WithCacheMethods is an option function for NewCacheMiddleware that sets the HTTP methods to cache.
func WithCacheMethods(methods []string) func(*cache.Config) {
	return func(config *cache.Config) {
		config.Methods = methods
	}
}

// WithValidator is an option function for NewKeyAuthMiddleware that sets a custom validator.
func WithValidator(validator func(*fiber.Ctx, string) (bool, error)) func(*keyauth.Config) {
	return func(config *keyauth.Config) {
		config.Validator = validator
	}
}

// WithErrorHandler is an option function for NewKeyAuthMiddleware that sets a custom error handler.
func WithErrorHandler(errorHandler func(*fiber.Ctx, error) error) func(*keyauth.Config) {
	return func(config *keyauth.Config) {
		config.ErrorHandler = errorHandler
	}
}

// WithSuccessHandler is an option function for NewKeyAuthMiddleware that sets a custom success handler.
func WithSuccessHandler(successHandler func(*fiber.Ctx) error) func(*keyauth.Config) {
	return func(config *keyauth.Config) {
		config.SuccessHandler = successHandler
	}
}

// WithKeyLookup is an option function for NewKeyAuthMiddleware that sets a custom key lookup.
func WithKeyLookup(keyLookup string) func(*keyauth.Config) {
	return func(config *keyauth.Config) {
		config.KeyLookup = keyLookup
	}
}

// WithContextKey is an option function for NewKeyAuthMiddleware that sets a custom context key.
func WithContextKey(contextKey string) func(*keyauth.Config) {
	return func(config *keyauth.Config) {
		config.ContextKey = contextKey
	}
}

// WithMax is an option function for NewRateLimiter that sets the maximum number of requests.
func WithMax(max int) func(*limiter.Config) {
	return func(config *limiter.Config) {
		config.Max = max
	}
}

// WithExpiration is an option function for NewRateLimiter that sets the expiration time.
func WithExpiration(expiration time.Duration) func(*limiter.Config) {
	return func(config *limiter.Config) {
		config.Expiration = expiration
	}
}

// WithLimitReached is an option function for NewRateLimiter that sets a custom limit reached handler.
func WithLimitReached(limitReached func(*fiber.Ctx) error) func(*limiter.Config) {
	return func(config *limiter.Config) {
		config.LimitReached = limitReached
	}
}

// WithKey is an option function for NewEncryptedCookieMiddleware that sets the encryption key.
func WithKey(key string) func(*encryptcookie.Config) {
	return func(config *encryptcookie.Config) {
		config.Key = key
	}
}

// WithEncryptor is an option function for NewEncryptedCookieMiddleware that sets a custom encryptor function.
func WithEncryptor(encryptor func(decryptedString, key string) (string, error)) func(*encryptcookie.Config) {
	return func(config *encryptcookie.Config) {
		config.Encryptor = encryptor
	}
}

// WithDecryptor is an option function for NewEncryptedCookieMiddleware that sets a custom decryptor function.
func WithDecryptor(decryptor func(encryptedString, key string) (string, error)) func(*encryptcookie.Config) {
	return func(config *encryptcookie.Config) {
		config.Decryptor = decryptor
	}
}

// useNonNilMiddleware registers non-nil middlewares to the fiber.Router.
func useNonNilMiddleware(g fiber.Router, middlewares ...fiber.Handler) {
	for _, middleware := range middlewares {
		if middleware != nil {
			g.Use(middleware)
		}
	}
}

// appendNonNilHandler appends non-nil handlers to the handlers slice.
func appendNonNilHandler(handlers []fiber.Handler, handlerFuncs ...fiber.Handler) []fiber.Handler {
	for _, handlerFunc := range handlerFuncs {
		if handlerFunc != nil {
			handlers = append(handlers, handlerFunc)
		}
	}
	return handlers
}

// WithAllowOrigins sets the allowed origins for CORS requests.
func WithAllowOrigins(origins string) func(*cors.Config) {
	return func(config *cors.Config) {
		config.AllowOrigins = origins
	}
}

// WithAllowMethods sets the allowed HTTP methods for CORS requests.
func WithAllowMethods(methods string) func(*cors.Config) {
	return func(config *cors.Config) {
		config.AllowMethods = methods
	}
}

// WithAllowHeaders sets the allowed headers for CORS requests.
func WithAllowHeaders(headers string) func(*cors.Config) {
	return func(config *cors.Config) {
		config.AllowHeaders = headers
	}
}

// WithExposeHeaders sets the headers that should be exposed to the client.
func WithExposeHeaders(headers string) func(*cors.Config) {
	return func(config *cors.Config) {
		config.ExposeHeaders = headers
	}
}

// WithAllowCredentials sets whether credentials are allowed for CORS requests.
func WithAllowCredentials(allow bool) func(*cors.Config) {
	return func(config *cors.Config) {
		config.AllowCredentials = allow
	}
}

// WithMaxAge sets the maximum age (in seconds) for preflight requests.
func WithMaxAge(maxAge int) func(*cors.Config) {
	return func(config *cors.Config) {
		config.MaxAge = maxAge
	}
}

// WithAllowOriginsFunc sets a custom function to determine the allowed origins for CORS requests.
//
// Example Usage:
//
//	func isAllowedOrigin(origin string) bool {
//		allowedOrigins := []string{
//			"https://example.com",
//			"https://api.example.com",
//		}
//		for _, allowedOrigin := range allowedOrigins {
//			if origin == allowedOrigin {
//				return true
//			}
//		}
//			return false
//		}
//
//	corsMiddleware := NewCORSMiddleware(
//	WithAllowOriginsFunc(isAllowedOrigin),
//	// Other options...
//
//	)
func WithAllowOriginsFunc(allowOriginsFunc func(string) bool) func(*cors.Config) {
	return func(config *cors.Config) {
		config.AllowOriginsFunc = allowOriginsFunc
	}
}

// WithRedirectStatusCode sets the HTTP status code for the redirect response.
func WithRedirectStatusCode(statusCode int) func(*redirect.Config) {
	return func(config *redirect.Config) {
		config.StatusCode = statusCode
	}
}

// WithSessionExpiration is an option function for NewSessionMiddleware that sets the session expiration time.
func WithSessionExpiration(expiration time.Duration) func(*session.Config) {
	return func(config *session.Config) {
		config.Expiration = expiration
	}
}

// WithSessionStorage is an option function for NewSessionMiddleware that sets the session storage backend.
func WithSessionStorage(storage fiber.Storage) func(*session.Config) {
	return func(config *session.Config) {
		config.Storage = storage
	}
}

// WithSessionKeyLookup is an option function for NewSessionMiddleware that sets the session key lookup.
func WithSessionKeyLookup(keyLookup string) func(*session.Config) {
	return func(config *session.Config) {
		config.KeyLookup = keyLookup
	}
}

// WithSessionCookieDomain is an option function for NewSessionMiddleware that sets the session cookie domain.
func WithSessionCookieDomain(cookieDomain string) func(*session.Config) {
	return func(config *session.Config) {
		config.CookieDomain = cookieDomain
	}
}

// WithSessionCookiePath is an option function for NewSessionMiddleware that sets the session cookie path.
func WithSessionCookiePath(cookiePath string) func(*session.Config) {
	return func(config *session.Config) {
		config.CookiePath = cookiePath
	}
}

// WithSessionCookieSecure is an option function for NewSessionMiddleware that sets the session cookie secure flag.
func WithSessionCookieSecure(cookieSecure bool) func(*session.Config) {
	return func(config *session.Config) {
		config.CookieSecure = cookieSecure
	}
}

// WithSessionCookieHTTPOnly is an option function for NewSessionMiddleware that sets the session cookie HTTP only flag.
func WithSessionCookieHTTPOnly(cookieHTTPOnly bool) func(*session.Config) {
	return func(config *session.Config) {
		config.CookieHTTPOnly = cookieHTTPOnly
	}
}

// WithSessionCookieSameSite is an option function for NewSessionMiddleware that sets the session cookie SameSite attribute.
func WithSessionCookieSameSite(cookieSameSite string) func(*session.Config) {
	return func(config *session.Config) {
		config.CookieSameSite = cookieSameSite
	}
}

// CleanupExpiredSessions is a goroutine that periodically cleans up expired sessions from the session store.
// It takes a session store and a cleanup interval as parameters.
//
// The function starts a ticker that triggers the cleanup process at the specified interval.
// On each tick, it calls the Reset method of the session store to remove all sessions.
// If an error occurs during the reset process, it is logged using the Logger.
//
// Note: The Reset method removes all sessions from the store, not just the expired ones.
// If more fine-grained control over removing only expired sessions is needed,
// additional implementation may be required based on the specific storage backend being used.
//
// The cleanup goroutine runs indefinitely until the ticker is stopped.
// It is typically started as a separate goroutine within the NewSessionMiddleware function.
//
// TODO: Leverage this "goroutine scheduler task".
func CleanupExpiredSessions(store *session.Store, interval time.Duration) {
	// Create a new ticker that triggers the cleanup process at the specified interval.
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Run the cleanup process indefinitely until the ticker is stopped.
	for range ticker.C {
		// Reset the session store to remove all sessions.
		err := store.Reset()
		if err != nil {
			// Log any errors that occur during the reset process.
			log.LogErrorf("Failed to reset session store: %v", err)
		}
	}
}

// WithStorage is an option function that sets the storage backend for various Fiber middlewares.
//
// It supports the following middleware configurations:
//
//	*cache.Config: Sets the storage backend for the cache middleware.
//	*session.Config: Sets the storage backend for the session middleware.
//	*limiter.Config: Sets the storage backend for the rate limiter middleware.
//	*csrf.Config: Sets the storage backend for the CSRF middleware.
//	*idempotency.Config: Sets the storage backend for the idempotency middleware.
//
// The storage backend must implement the fiber.Storage interface.
//
// Example usage:
//
//	// Create a custom storage backend
//	storage := myCustomStorage{}
//
//	// Use the WithStorage option function to set the storage backend for the cache middleware
//	cacheMiddleware := NewCacheMiddleware(WithStorage(storage))
//
//	// Use the WithStorage option function to set the storage backend for the session middleware
//	sessionMiddleware := NewSessionMiddleware(WithStorage(storage))
//
//	// Use the WithStorage option function to set the storage backend for the rate limiter middleware
//	rateLimiterMiddleware := NewRateLimiter(WithStorage(storage))
//
//	// Use the WithStorage option function to set the storage backend for the CSRF middleware
//	csrfMiddleware := NewCSRFMiddleware(WithStorage(storage))
//
//	// Use the WithStorage option function to set the storage backend for the idempotency middleware
//	idempotencyMiddleware := NewIdempotencyMiddleware(WithStorage(storage))
//
// Note:
//   - If an unsupported middleware configuration is passed to WithStorage, it will panic with an error message.
//   - Additional storage support for other middlewares will be implemented in the future as needed.
func WithStorage(storage fiber.Storage) any {
	return func(config any) {
		switch cfg := config.(type) {
		case *cache.Config:
			cfg.Storage = storage
		case *session.Config:
			cfg.Storage = storage
		case *limiter.Config:
			cfg.Storage = storage
		case *csrf.Config:
			cfg.Storage = storage
		case *idempotency.Config:
			cfg.Storage = storage
		default:
			panic(fmt.Sprintf("unsupported config type: %T", config))
		}
	}
}

// WithCSRFKeyLookup is an option function for NewCSRFMiddleware that sets the key lookup for the CSRF token.
func WithCSRFKeyLookup(keyLookup string) func(*csrf.Config) {
	return func(config *csrf.Config) {
		config.KeyLookup = keyLookup
	}
}

// WithCSRFCookieName is an option function for NewCSRFMiddleware that sets the name of the CSRF cookie.
func WithCSRFCookieName(cookieName string) func(*csrf.Config) {
	return func(config *csrf.Config) {
		config.CookieName = cookieName
	}
}

// WithCSRFCookieDomain is an option function for NewCSRFMiddleware that sets the domain of the CSRF cookie.
func WithCSRFCookieDomain(cookieDomain string) func(*csrf.Config) {
	return func(config *csrf.Config) {
		config.CookieDomain = cookieDomain
	}
}

// WithCSRFCookiePath is an option function for NewCSRFMiddleware that sets the path of the CSRF cookie.
func WithCSRFCookiePath(cookiePath string) func(*csrf.Config) {
	return func(config *csrf.Config) {
		config.CookiePath = cookiePath
	}
}

// WithCSRFCookieSecure is an option function for NewCSRFMiddleware that sets the secure flag of the CSRF cookie.
func WithCSRFCookieSecure(cookieSecure bool) func(*csrf.Config) {
	return func(config *csrf.Config) {
		config.CookieSecure = cookieSecure
	}
}

// WithCSRFCookieHTTPOnly is an option function for NewCSRFMiddleware that sets the HTTP only flag of the CSRF cookie.
func WithCSRFCookieHTTPOnly(cookieHTTPOnly bool) func(*csrf.Config) {
	return func(config *csrf.Config) {
		config.CookieHTTPOnly = cookieHTTPOnly
	}
}

// WithCSRFCookieSameSite is an option function for NewCSRFMiddleware that sets the SameSite attribute of the CSRF cookie.
func WithCSRFCookieSameSite(cookieSameSite string) func(*csrf.Config) {
	return func(config *csrf.Config) {
		config.CookieSameSite = cookieSameSite
	}
}

// WithCSRFCookieSessionOnly is an option function for NewCSRFMiddleware that sets the session-only flag of the CSRF cookie.
func WithCSRFCookieSessionOnly(cookieSessionOnly bool) func(*csrf.Config) {
	return func(config *csrf.Config) {
		config.CookieSessionOnly = cookieSessionOnly
	}
}

// WithCSRFExpiration is an option function for NewCSRFMiddleware that sets the expiration time of the CSRF token.
func WithCSRFExpiration(expiration time.Duration) func(*csrf.Config) {
	return func(config *csrf.Config) {
		config.Expiration = expiration
	}
}

// WithCSRFSingleUseToken is an option function for NewCSRFMiddleware that sets the single-use token flag.
func WithCSRFSingleUseToken(singleUseToken bool) func(*csrf.Config) {
	return func(config *csrf.Config) {
		config.SingleUseToken = singleUseToken
	}
}

// WithCSRFSession is an option function for NewCSRFMiddleware that sets the session store for the CSRF middleware.
func WithCSRFSession(session *session.Store) func(*csrf.Config) {
	return func(config *csrf.Config) {
		config.Session = session
	}
}

// WithCSRFSessionKey is an option function for NewCSRFMiddleware that sets the session key for storing the CSRF token.
func WithCSRFSessionKey(sessionKey string) func(*csrf.Config) {
	return func(config *csrf.Config) {
		config.SessionKey = sessionKey
	}
}

// WithCSRFContextKey is an option function for NewCSRFMiddleware that sets the context key for storing the CSRF token.
func WithCSRFContextKey(contextKey any) func(*csrf.Config) {
	return func(config *csrf.Config) {
		config.ContextKey = contextKey
	}
}

// WithCSRFKeyGenerator is an option function for NewCSRFMiddleware that sets the key generator function for the CSRF token.
func WithCSRFKeyGenerator(keyGenerator func() string) func(*csrf.Config) {
	return func(config *csrf.Config) {
		config.KeyGenerator = keyGenerator
	}
}

// WithCSRFErrorHandler is an option function for NewCSRFMiddleware that sets the error handler for the CSRF middleware.
func WithCSRFErrorHandler(errorHandler fiber.ErrorHandler) func(*csrf.Config) {
	return func(config *csrf.Config) {
		config.ErrorHandler = errorHandler
	}
}

// WithCSRFExtractor is an option function for NewCSRFMiddleware that sets the extractor function for retrieving the CSRF token.
func WithCSRFExtractor(extractor func(*fiber.Ctx) (string, error)) func(*csrf.Config) {
	return func(config *csrf.Config) {
		config.Extractor = extractor
	}
}

// WithCSRFHandlerContextKey is an option function for NewCSRFMiddleware that sets the context key for storing the CSRF handler.
func WithCSRFHandlerContextKey(handlerContextKey any) func(*csrf.Config) {
	return func(config *csrf.Config) {
		config.HandlerContextKey = handlerContextKey
	}
}

// WithXSSProtection is an option function for NewHelmetMiddleware that sets the X-XSS-Protection header.
func WithXSSProtection(protection string) func(*helmet.Config) {
	return func(config *helmet.Config) {
		config.XSSProtection = protection
	}
}

// WithContentTypeNosniff is an option function for NewHelmetMiddleware that sets the X-Content-Type-Options header.
func WithContentTypeNosniff(nosniff string) func(*helmet.Config) {
	return func(config *helmet.Config) {
		config.ContentTypeNosniff = nosniff
	}
}

// WithXFrameOptions is an option function for NewHelmetMiddleware that sets the X-Frame-Options header.
func WithXFrameOptions(options string) func(*helmet.Config) {
	return func(config *helmet.Config) {
		config.XFrameOptions = options
	}
}

// WithReferrerPolicy is an option function for NewHelmetMiddleware that sets the Referrer-Policy header.
func WithReferrerPolicy(policy string) func(*helmet.Config) {
	return func(config *helmet.Config) {
		config.ReferrerPolicy = policy
	}
}

// WithCrossOriginEmbedderPolicy is an option function for NewHelmetMiddleware that sets the Cross-Origin-Embedder-Policy header.
func WithCrossOriginEmbedderPolicy(policy string) func(*helmet.Config) {
	return func(config *helmet.Config) {
		config.CrossOriginEmbedderPolicy = policy
	}
}

// WithCrossOriginOpenerPolicy is an option function for NewHelmetMiddleware that sets the Cross-Origin-Opener-Policy header.
func WithCrossOriginOpenerPolicy(policy string) func(*helmet.Config) {
	return func(config *helmet.Config) {
		config.CrossOriginOpenerPolicy = policy
	}
}

// WithCrossOriginResourcePolicy is an option function for NewHelmetMiddleware that sets the Cross-Origin-Resource-Policy header.
func WithCrossOriginResourcePolicy(policy string) func(*helmet.Config) {
	return func(config *helmet.Config) {
		config.CrossOriginResourcePolicy = policy
	}
}

// WithOriginAgentCluster is an option function for NewHelmetMiddleware that sets the Origin-Agent-Cluster header.
func WithOriginAgentCluster(value string) func(*helmet.Config) {
	return func(config *helmet.Config) {
		config.OriginAgentCluster = value
	}
}

// WithXDNSPrefetchControl is an option function for NewHelmetMiddleware that sets the X-DNS-Prefetch-Control header.
func WithXDNSPrefetchControl(control string) func(*helmet.Config) {
	return func(config *helmet.Config) {
		config.XDNSPrefetchControl = control
	}
}

// WithXDownloadOptions is an option function for NewHelmetMiddleware that sets the X-Download-Options header.
func WithXDownloadOptions(options string) func(*helmet.Config) {
	return func(config *helmet.Config) {
		config.XDownloadOptions = options
	}
}

// WithXPermittedCrossDomain is an option function for NewHelmetMiddleware that sets the X-Permitted-Cross-Domain-Policies header.
func WithXPermittedCrossDomain(policies string) func(*helmet.Config) {
	return func(config *helmet.Config) {
		config.XPermittedCrossDomain = policies
	}
}

// WithUsers is an option function for NewBasicAuthMiddleware that sets the authorized users.
func WithUsers(users map[string]string) func(*basicauth.Config) {
	return func(config *basicauth.Config) {
		config.Users = users
	}
}

// WithRealm is an option function for NewBasicAuthMiddleware that sets the authentication realm.
func WithRealm(realm string) func(*basicauth.Config) {
	return func(config *basicauth.Config) {
		config.Realm = realm
	}
}

// WithAuthorizer is an option function for NewBasicAuthMiddleware that sets the authorizer function.
func WithAuthorizer(authorizer func(string, string) bool) func(*basicauth.Config) {
	return func(config *basicauth.Config) {
		config.Authorizer = authorizer
	}
}

// WithUnauthorized is an option function for NewBasicAuthMiddleware that sets the unauthorized handler.
func WithUnauthorized(unauthorized fiber.Handler) func(*basicauth.Config) {
	return func(config *basicauth.Config) {
		config.Unauthorized = unauthorized
	}
}

// WithContextUsername is an option function for NewBasicAuthMiddleware that sets the context key for the authenticated username.
func WithContextUsername(contextUsername string) func(*basicauth.Config) {
	return func(config *basicauth.Config) {
		config.ContextUsername = contextUsername
	}
}

// WithContextPassword is an option function for NewBasicAuthMiddleware that sets the context key for the authenticated password.
func WithContextPassword(contextPassword string) func(*basicauth.Config) {
	return func(config *basicauth.Config) {
		config.ContextPassword = contextPassword
	}
}

// WithSwaggerBasePath is an option function for NewSwaggerMiddleware that sets the base path for the UI path.
func WithSwaggerBasePath(basePath string) func(*swagger.Config) {
	return func(config *swagger.Config) {
		config.BasePath = basePath
	}
}

// WithSwaggerFilePath is an option function for NewSwaggerMiddleware that sets the file path for the swagger.json or swagger.yaml file.
func WithSwaggerFilePath(filePath string) func(*swagger.Config) {
	return func(config *swagger.Config) {
		config.FilePath = filePath
	}
}

// WithSwaggerPath is an option function for NewSwaggerMiddleware that sets the path that combines with BasePath for the full UI path.
func WithSwaggerPath(path string) func(*swagger.Config) {
	return func(config *swagger.Config) {
		config.Path = path
	}
}

// WithSwaggerTitle is an option function for NewSwaggerMiddleware that sets the title for the documentation site.
func WithSwaggerTitle(title string) func(*swagger.Config) {
	return func(config *swagger.Config) {
		config.Title = title
	}
}

// WithSwaggerCacheAge is an option function for NewSwaggerMiddleware that sets the max-age for the Cache-Control header in seconds.
func WithSwaggerCacheAge(cacheAge int) func(*swagger.Config) {
	return func(config *swagger.Config) {
		config.CacheAge = cacheAge
	}
}

// WithNext is an option function that sets the Next function to skip the middleware when returned true.
//
// It supports the following middleware configurations:
//
//	*cache.Config: Sets the Next function for the cache middleware.
//	*cors.Config: Sets the Next function for the CORS middleware.
//	*csrf.Config: Sets the Next function for the CSRF middleware.
//	*redirect.Config: Sets the Next function for the redirect middleware.
//	*swagger.Config: Sets the Next function for the Swagger middleware.
//	*validator.Config: Sets the Next function for the Validator middleware.
//
// The Next function takes a *fiber.Ctx as a parameter and returns a boolean value.
// If the Next function returns true, the middleware will be skipped for the current request.
//
// Example usage:
//
//	// Define a custom Next function
//	func customNext(c *fiber.Ctx) bool {
//	    // Custom logic to determine whether to skip the middleware
//	    // Return true to skip the middleware, false otherwise
//	    return c.Path() == "/skip"
//	}
//
//	// Use the WithNext option function to set the Next function for the cache middleware
//	cacheMiddleware := NewCacheMiddleware(WithNext(customNext))
//
//	// Use the WithNext option function to set the Next function for the CORS middleware
//	corsMiddleware := NewCORSMiddleware(WithNext(customNext))
//
//	// Use the WithNext option function to set the Next function for the CSRF middleware
//	csrfMiddleware := NewCSRFMiddleware(WithNext(customNext))
//
//	// Use the WithNext option function to set the Next function for the redirect middleware
//	redirectMiddleware := NewRedirectMiddleware(WithNext(customNext))
//
//	// Use the WithNext option function to set the Next function for the Swagger middleware
//	swaggerMiddleware := NewSwaggerMiddleware(WithNext(customNext))
//
//	// Use the WithNext option function to set the Next function for the Validator middleware
//	validatorMiddleware := NewValidatorMiddleware(WithNext(customNext))
//
// Note:
//   - If an unsupported middleware configuration is passed to WithNext, it will panic with an error message.
//   - Additional "Next" functionality for other middlewares will be added based on future requirements.
func WithNext(next func(c *fiber.Ctx) bool) any {
	return func(config any) {
		switch cfg := config.(type) {
		case *cache.Config:
			cfg.Next = next
		case *cors.Config:
			cfg.Next = next
		case *csrf.Config:
			cfg.Next = next
		case *redirect.Config:
			cfg.Next = next
		case *swagger.Config:
			cfg.Next = next
		case *validator.Config:
			cfg.Next = next
		default:
			panic(fmt.Sprintf("unsupported config type: %T", config))
		}
	}
}

// WithETagWeak is an option function for NewETagMiddleware that sets the weak ETag flag.
func WithETagWeak(weak bool) func(*etag.Config) {
	return func(config *etag.Config) {
		config.Weak = weak
	}
}

// WithFaviconFile is an option function for NewFaviconMiddleware that sets the file path of the favicon.
func WithFaviconFile(file string) func(*favicon.Config) {
	return func(config *favicon.Config) {
		config.File = file
	}
}

// WithFaviconURL is an option function for NewFaviconMiddleware that sets the URL path where the favicon will be served.
func WithFaviconURL(url string) func(*favicon.Config) {
	return func(config *favicon.Config) {
		config.URL = url
	}
}

// WithIdempotencyLifetime is an option function for NewIdempotencyMiddleware that sets the maximum lifetime of an idempotency key.
func WithIdempotencyLifetime(lifetime time.Duration) func(*idempotency.Config) {
	return func(config *idempotency.Config) {
		config.Lifetime = lifetime
	}
}

// WithIdempotencyKeyHeader is an option function for NewIdempotencyMiddleware that sets the name of the header that contains the idempotency key.
func WithIdempotencyKeyHeader(keyHeader string) func(*idempotency.Config) {
	return func(config *idempotency.Config) {
		config.KeyHeader = keyHeader
	}
}

// WithIdempotencyKeyHeaderValidate is an option function for NewIdempotencyMiddleware that sets the function to validate the syntax of the idempotency header.
func WithIdempotencyKeyHeaderValidate(keyHeaderValidate func(string) error) func(*idempotency.Config) {
	return func(config *idempotency.Config) {
		config.KeyHeaderValidate = keyHeaderValidate
	}
}

// WithIdempotencyKeepResponseHeaders is an option function for NewIdempotencyMiddleware that sets the list of headers that should be kept from the original response.
func WithIdempotencyKeepResponseHeaders(keepResponseHeaders []string) func(*idempotency.Config) {
	return func(config *idempotency.Config) {
		config.KeepResponseHeaders = keepResponseHeaders
	}
}

// WithIdempotencyLock is an option function for NewIdempotencyMiddleware that sets the locker for idempotency keys.
func WithIdempotencyLock(lock idempotency.Locker) func(*idempotency.Config) {
	return func(config *idempotency.Config) {
		config.Lock = lock
	}
}

// WithRules is an option function that sets the rewrite or redirect rules for the Rewrite or Redirect middleware.
//
// It supports the following middleware configurations:
//
//	*rewrite.Config: Sets the rewrite rules for the Rewrite middleware.
//	*redirect.Config: Sets the redirect rules for the Redirect middleware.
//
// The rules are defined as a map of string keys and values. For the Rewrite middleware, the keys represent the URL path
// patterns to match, and the values represent the replacement paths. Captured values can be retrieved by index using the
// $1, $2, etc. syntax. For the Redirect middleware, the keys represent the source paths, and the values represent the
// destination paths.
//
// Example usage:
//
//	// Define the rewrite rules
//	rewriteRules := map[string]string{
//	    "/old":              "/new",
//	    "/api/*":            "/$1",
//	    "/js/*":             "/public/javascripts/$1",
//	    "/users/*/orders/*": "/user/$1/order/$2",
//	}
//
//	// Use the WithRules option function to set the rewrite rules for the Rewrite middleware
//	rewriteMiddleware := NewRewriteMiddleware(WithRules(rewriteRules))
//
//	// Define the redirect rules
//	redirectRules := map[string]string{
//	    "/old":              "/new",
//	    "/api/*":            "/$1",
//	    "/js/*":             "/public/javascripts/$1",
//	    "/users/*/orders/*": "/user/$1/order/$2",
//	}
//
//	// Use the WithRules option function to set the redirect rules for the Redirect middleware
//	redirectMiddleware := NewRedirectMiddleware(WithRules(redirectRules))
//
// Note:
//   - If an unsupported middleware configuration is passed to WithRules, it will panic with an error message.
func WithRules(rules map[string]string) any {
	// Note: now, this reusable, get good get golang.
	return func(config any) {
		switch cfg := config.(type) {
		case *rewrite.Config:
			cfg.Rules = rules
		case *redirect.Config:
			cfg.Rules = rules
		default:
			panic(fmt.Sprintf("unsupported config type: %T", config))
		}
	}
}

// WithValidatorRules is an option function for NewValidatorMiddleware that sets the validation rules.
func WithValidatorRules(rules []validator.Restrictor) func(*validator.Config) {
	return func(config *validator.Config) {
		config.Rules = rules
	}
}

// WithValidatorErrorHandler is an option function for NewValidatorMiddleware that sets the error handler function.
func WithValidatorErrorHandler(errorHandler func(c *fiber.Ctx, err error) error) func(*validator.Config) {
	return func(config *validator.Config) {
		config.ErrorHandler = errorHandler
	}
}

// WithValidatorContextKey is an option function for NewValidatorMiddleware that sets the context key for storing the validation result.
func WithValidatorContextKey(contextKey string) func(*validator.Config) {
	return func(config *validator.Config) {
		config.ContextKey = contextKey
	}
}

// ptr is a helper function that takes an int value and returns a pointer to it.
// This is useful when a pointer to an int value needs to be passed, such as when setting
// the Max value in a validator.Restrictor.
//
// Example usage:
//
//	app.Use(validator.New(validator.Config{
//		Rules: []validator.Restrictor{
//			validator.RestrictNumberOnly{
//				Fields: []string{"seafood_price"},
//				Max:    ptr(100),
//			},
//		},
//	}))
func ptr(i int) *int {
	return &i
}

// WithRequestIDHeader is an option function for NewRequestIDMiddleware that sets the header name for the request ID.
func WithRequestIDHeader(header string) func(*requestid.Config) {
	return func(config *requestid.Config) {
		config.Header = header
	}
}

// WithRequestIDHeaderContextKey is an option function for NewRequestIDMiddleware that sets the header Context Key name for the request ID.
func WithRequestIDHeaderContextKey(headerContextKey string) func(*requestid.Config) {
	return func(config *requestid.Config) {
		config.ContextKey = headerContextKey
	}
}

// WithRequestIDGenerator is an option function for NewRequestIDMiddleware that sets a custom generator function for the request ID.
func WithRequestIDGenerator(generator func() string) func(*requestid.Config) {
	return func(config *requestid.Config) {
		config.Generator = generator
	}
}

func digest(clientIP string) string {
	h := sha256.New()
	h.Write([]byte(clientIP))
	digest := h.Sum(nil)

	return fmt.Sprintf("%x", digest)
}
