// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package middleware

import (
	"crypto/tls"
	"fmt"
	log "h0llyw00dz-template/backend/internal/logger"
	"h0llyw00dz-template/backend/internal/middleware/csp"
	"h0llyw00dz-template/backend/internal/middleware/monitor"
	"h0llyw00dz-template/backend/internal/middleware/restime"
	"h0llyw00dz-template/backend/pkg/restapis/helper"
	"hash/fnv"
	"time"

	validator "github.com/H0llyW00dzZ/FiberValidator"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
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
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/gofiber/fiber/v2/middleware/redirect"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/rewrite"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
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

// WithLimiterKeyGenerator is an option function that sets a custom key generator for the rate limiter middleware.
func WithLimiterKeyGenerator(keyGenerator func(*fiber.Ctx) string) func(*limiter.Config) {
	return func(config *limiter.Config) {
		config.KeyGenerator = keyGenerator
	}
}

// WithCacheKeyGenerator is an option function that sets a custom key generator for the cache middleware.
func WithCacheKeyGenerator(keyGenerator func(*fiber.Ctx) string) func(*cache.Config) {
	return func(config *cache.Config) {
		config.KeyGenerator = keyGenerator
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

// WithKeyAuthContextKey is an option function for NewKeyAuthMiddleware that sets the context key for the Key Auth middleware.
func WithKeyAuthContextKey(contextKey any) func(*keyauth.Config) {
	return func(config *keyauth.Config) {
		config.ContextKey = contextKey
	}
}

// WithCSPContextKey is an option function for NewCSPHeaderGenerator that sets the context key for the CSP Header Generator middleware.
func WithCSPContextKey(contextKey any) func(*csp.Config) {
	return func(config *csp.Config) {
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

// WithCORSNext is an option function for NewCORSMiddleware that sets the Next function to skip the CORS middleware.
func WithCORSNext(next func(*fiber.Ctx) bool) func(*cors.Config) {
	return func(config *cors.Config) {
		config.Next = next
	}
}

// WithCSRFNext is an option function for NewCSRFMiddleware that sets the Next function to skip the CSRF middleware.
func WithCSRFNext(next func(*fiber.Ctx) bool) func(*csrf.Config) {
	return func(config *csrf.Config) {
		config.Next = next
	}
}

// WithRedirectNext is an option function for NewRedirectMiddleware that sets the Next function to skip the redirect middleware.
func WithRedirectNext(next func(*fiber.Ctx) bool) func(*redirect.Config) {
	return func(config *redirect.Config) {
		config.Next = next
	}
}

// WithSwaggerNext is an option function for NewSwaggerMiddleware that sets the Next function to skip the Swagger middleware.
func WithSwaggerNext(next func(*fiber.Ctx) bool) func(*swagger.Config) {
	return func(config *swagger.Config) {
		config.Next = next
	}
}

// WithValidatorNext is an option function for NewValidatorMiddleware that sets the Next function to skip the Validator middleware.
func WithValidatorNext(next func(*fiber.Ctx) bool) func(*validator.Config) {
	return func(config *validator.Config) {
		config.Next = next
	}
}

// WithHealthCheckNext is an option function for NewHealthZCheck that sets the Next function to skip the HealthZ Check middleware.
func WithHealthCheckNext(next func(*fiber.Ctx) bool) func(*healthcheck.Config) {
	return func(config *healthcheck.Config) {
		config.Next = next
	}
}

// WithCSPNext is an option function for NewCSPHeaderGenerator that sets the Next function to skip the CSP Header Generator middleware.
func WithCSPNext(next func(*fiber.Ctx) bool) func(*csp.Config) {
	return func(config *csp.Config) {
		config.Next = next
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

// WithRewriteRules is an option function for the Rewrite middleware that sets the rewrite rules.
//
// The rewrite rules are defined as a map of string keys and values. The keys represent the URL path patterns to match,
// and the values represent the replacement paths. Captured values can be retrieved by index using the $1, $2, etc. syntax.
//
// Example usage:
//
//	rewriteRules := map[string]string{
//	    "/old":              "/new",
//	    "/api/*":            "/$1",
//	    "/js/*":             "/public/javascripts/$1",
//	    "/users/*/orders/*": "/user/$1/order/$2",
//	}
//
//	rewriteMiddleware := NewRewriteMiddleware(WithRewriteRules(rewriteRules))
func WithRewriteRules(rules map[string]string) func(*rewrite.Config) {
	return func(config *rewrite.Config) {
		config.Rules = rules
	}
}

// WithRewriteNext is an option function for the Rewrite middleware that sets the Next function.
func WithRewriteNext(next func(c *fiber.Ctx) bool) func(*rewrite.Config) {
	return func(config *rewrite.Config) {
		config.Next = next
	}
}

// WithRedirectRules is an option function for the Redirect middleware that sets the redirect rules.
//
// The redirect rules are defined as a map of string keys and values. The keys represent the source paths,
// and the values represent the destination paths.
//
// Example usage:
//
//	redirectRules := map[string]string{
//	    "/old":              "/new",
//	    "/api/*":            "/$1",
//	    "/js/*":             "/public/javascripts/$1",
//	    "/users/*/orders/*": "/user/$1/order/$2",
//	}
//
//	redirectMiddleware := NewRedirectMiddleware(WithRedirectRules(redirectRules))
func WithRedirectRules(rules map[string]string) func(*redirect.Config) {
	return func(config *redirect.Config) {
		config.Rules = rules
	}
}

// WithNextRedirect is an option function for the Redirect middleware that sets the Next function.
func WithNextRedirect(next func(c *fiber.Ctx) bool) func(*redirect.Config) {
	return func(config *redirect.Config) {
		config.Next = next
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

// WithLivenessProbe is an option function for NewHealthZCheck that sets the liveness probe function.
func WithLivenessProbe(livenessProbe healthcheck.HealthChecker) func(*healthcheck.Config) {
	return func(config *healthcheck.Config) {
		config.LivenessProbe = livenessProbe
	}
}

// WithLivenessEndpoint is an option function for NewHealthZCheck that sets the HTTP endpoint for the liveness probe.
func WithLivenessEndpoint(livenessEndpoint string) func(*healthcheck.Config) {
	return func(config *healthcheck.Config) {
		config.LivenessEndpoint = livenessEndpoint
	}
}

// WithReadinessProbe is an option function for NewHealthZCheck that sets the readiness probe function.
func WithReadinessProbe(readinessProbe healthcheck.HealthChecker) func(*healthcheck.Config) {
	return func(config *healthcheck.Config) {
		config.ReadinessProbe = readinessProbe
	}
}

// WithReadinessEndpoint is an option function for NewHealthZCheck that sets the HTTP endpoint for the readiness probe.
func WithReadinessEndpoint(readinessEndpoint string) func(*healthcheck.Config) {
	return func(config *healthcheck.Config) {
		config.ReadinessEndpoint = readinessEndpoint
	}
}

// WithRandomnessGenerator is an option function for NewCSPHeaderGenerator that sets the randomness generator function.
func WithRandomnessGenerator(customRand func(string) string) func(*csp.Config) {
	return func(config *csp.Config) {
		config.RandomnessGenerator = customRand
	}
}

// WithCSPValueGenerator is an option function for NewCSPHeaderGenerator that sets the Content-Security-Policy Value generator function.
func WithCSPValueGenerator(cspvalue func(string, map[string]string) string) func(*csp.Config) {
	return func(config *csp.Config) {
		config.CSPValueGenerator = func(randomness string, customValues map[string]string) string {
			return cspvalue(randomness, customValues)
		}
	}
}

// WithCSPIPHeader is an option function that sets the header name used to retrieve the client IP address for the CSP middleware.
func WithCSPIPHeader(ipHeader string) func(*csp.Config) {
	return func(config *csp.Config) {
		config.IPHeader = ipHeader
	}
}

// WithSessionIDGenerator is an option function for NewSessionMiddleware that sets a custom generator function for the session ID.
func WithSessionIDGenerator(generator func() string) func(*session.Config) {
	return func(config *session.Config) {
		config.KeyGenerator = generator
	}
}

// WithCacheNext is an option function for NewCacheMiddleware that sets the Next function to skip the cache middleware.
func WithCacheNext(next func(*fiber.Ctx) bool) func(*cache.Config) {
	return func(config *cache.Config) {
		config.Next = next
	}
}

// WithCompressLevel is an option function for NewCompressMiddleware that sets the compression level.
func WithCompressLevel(level compress.Level) func(*compress.Config) {
	return func(config *compress.Config) {
		config.Level = level
	}
}

// WithCompressNext is an option function for NewCompressMiddleware that sets the Next function to skip the compression middleware.
func WithCompressNext(next func(*fiber.Ctx) bool) func(*compress.Config) {
	return func(config *compress.Config) {
		config.Next = next
	}
}

// WithEarlyDataNext is an option function for NewEarlyData (Another QUIC) that sets the Next function.
func WithEarlyDataNext(next func(c *fiber.Ctx) bool) func(*earlydata.Config) {
	return func(config *earlydata.Config) {
		config.Next = next
	}
}

// WithIsEarlyData is an option function for NewEarlyData (Another QUIC) that sets the IsEarlyData function.
func WithIsEarlyData(isEarlyData func(c *fiber.Ctx) bool) func(*earlydata.Config) {
	return func(config *earlydata.Config) {
		config.IsEarlyData = isEarlyData
	}
}

// WithAllowEarlyData is an option function for NewEarlyData (Another QUIC) that sets the AllowEarlyData function.
func WithAllowEarlyData(allowEarlyData func(c *fiber.Ctx) bool) func(*earlydata.Config) {
	return func(config *earlydata.Config) {
		config.AllowEarlyData = allowEarlyData
	}
}

// WithEarlyDataError is an option function for NewEarlyData (Another QUIC) that sets the Error value.
func WithEarlyDataError(err error) func(*earlydata.Config) {
	return func(config *earlydata.Config) {
		config.Error = err
	}
}

// WithPrometheusServiceName is an option function for NewPrometheus that sets the service name.
func WithPrometheusServiceName(serviceName string) func(*monitor.PrometheusConfig) {
	return func(config *monitor.PrometheusConfig) {
		config.ServiceName = serviceName
	}
}

// WithPrometheusNamespace is an option function for NewPrometheus that sets the namespace.
func WithPrometheusNamespace(namespace string) func(*monitor.PrometheusConfig) {
	return func(config *monitor.PrometheusConfig) {
		config.Namespace = namespace
	}
}

// WithPrometheusSubsystem is an option function for NewPrometheus that sets the subsystem.
func WithPrometheusSubsystem(subsystem string) func(*monitor.PrometheusConfig) {
	return func(config *monitor.PrometheusConfig) {
		config.Subsystem = subsystem
	}
}

// WithPrometheusLabels is an option function for NewPrometheus that sets the labels.
func WithPrometheusLabels(labels map[string]string) func(*monitor.PrometheusConfig) {
	return func(config *monitor.PrometheusConfig) {
		config.Labels = labels
	}
}

// WithPrometheusSkipPaths is an option function for NewPrometheus that sets the skip paths.
func WithPrometheusSkipPaths(skipPaths []string) func(*monitor.PrometheusConfig) {
	return func(config *monitor.PrometheusConfig) {
		config.SkipPaths = skipPaths
	}
}

// WithPrometheusMetricsPaths is an option function for NewPrometheus that sets the metrics path.
func WithPrometheusMetricsPaths(paths string) func(*monitor.PrometheusConfig) {
	return func(config *monitor.PrometheusConfig) {
		config.MetricsPath = paths
	}
}

// WithPrometheusMetricsNext is an option function for NewPrometheus that sets the next function.
func WithPrometheusMetricsNext(next func(c *fiber.Ctx) bool) func(*monitor.PrometheusConfig) {
	return func(config *monitor.PrometheusConfig) {
		config.Next = next
	}
}

// WithPrometheusCacheKey is an option function for NewPrometheus that sets the cache key.
func WithPrometheusCacheKey(cacheKey string) func(*monitor.PrometheusConfig) {
	return func(config *monitor.PrometheusConfig) {
		config.CacheKey = cacheKey
	}
}

// WithRestimeHeaderName is an option function for ResponseTime that sets the header name.
func WithRestimeHeaderName(headerName string) func(*restime.Config) {
	return func(config *restime.Config) {
		config.HeaderName = headerName
	}
}

// WithRestimeNext is an option function for ResponseTime that sets the next function.
func WithRestimeNext(next func(c *fiber.Ctx) bool) func(*restime.Config) {
	return func(config *restime.Config) {
		config.Next = next
	}
}

// WithProxyingNext is an option function for NewProxying that sets the next function.
func WithProxyingNext(next func(c *fiber.Ctx) bool) func(*proxy.Config) {
	return func(config *proxy.Config) {
		config.Next = next
	}
}

// WithProxyingServers is an option function for NewProxying that sets the list of servers.
func WithProxyingServers(servers []string) func(*proxy.Config) {
	return func(config *proxy.Config) {
		config.Servers = servers
	}
}

// WithProxyingModifyRequest is an option function for NewProxying that sets the request modifier function.
func WithProxyingModifyRequest(modifyRequest fiber.Handler) func(*proxy.Config) {
	return func(config *proxy.Config) {
		config.ModifyRequest = modifyRequest
	}
}

// WithProxyingModifyResponse is an option function for NewProxying that sets the response modifier function.
func WithProxyingModifyResponse(modifyResponse fiber.Handler) func(*proxy.Config) {
	return func(config *proxy.Config) {
		config.ModifyResponse = modifyResponse
	}
}

// WithProxyingTimeout is an option function for NewProxying that sets the timeout duration.
func WithProxyingTimeout(timeout time.Duration) func(*proxy.Config) {
	return func(config *proxy.Config) {
		config.Timeout = timeout
	}
}

// WithProxyingReadBufferSize is an option function for NewProxying that sets the read buffer size.
func WithProxyingReadBufferSize(size int) func(*proxy.Config) {
	return func(config *proxy.Config) {
		config.ReadBufferSize = size
	}
}

// WithProxyingWriteBufferSize is an option function for NewProxying that sets the write buffer size.
func WithProxyingWriteBufferSize(size int) func(*proxy.Config) {
	return func(config *proxy.Config) {
		config.WriteBufferSize = size
	}
}

// WithProxyingTLSConfig is an option function for NewProxying that sets the TLS configuration.
//
// Note: For private CAs/public CAs, it is recommended to use a certificate with X509v3 Extended Key Usage
// for TLS Web Server Authentication and TLS Web Client Authentication. To make it faster,
// use ECC (Elliptic Curve Cryptography) instead of RSA, as it can save bandwidth costs.
// If bandwidth is free because this proxying is only used for internal mode, then RSA can be used.
func WithProxyingTLSConfig(tlsConfig *tls.Config) func(*proxy.Config) {
	return func(config *proxy.Config) {
		config.TlsConfig = tlsConfig
	}
}

// WithProxyingClient is an option function for NewProxying that sets the custom client.
//
// TODO: Implement a sub-helper configuration for this, as it is possible to implement an ingress mechanism through a load balancer,
// which can be useful for bare metal or creating an own network gateway in a cloud provider such as GKE.
func WithProxyingClient(client *fasthttp.LBClient) func(*proxy.Config) {
	return func(config *proxy.Config) {
		config.Client = client
	}
}
