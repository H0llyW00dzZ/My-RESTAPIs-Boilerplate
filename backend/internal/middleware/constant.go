// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package middleware

import "time"

const (
	// MsgPanicOccurred is the log message indicating that a panic occurred within the application.
	MsgPanicOccurred = "Panic occurred: %v"
	// MsgStackTrace is the log message containing the stack trace of a panic.
	MsgStackTrace = "Stack trace:\n%s"
	// MsgRESTAPIsVisitorGotRateLimited is the log message indicating that a visitor to the REST APIs
	// has been rate limited.
	MsgRESTAPIsVisitorGotRateLimited = "[REST APIs] visitor got rate limited"
)

const (

	// maxRequestRESTAPIsRateLimiter is the maximum number of requests a client can make to
	// REST API endpoints within the time window specified by maxExpirationRESTAPIsRateLimiter.
	maxRequestRESTAPIsRateLimiter = 10

	// maxExpirationRESTAPIsRateLimiter is the time window for counting REST API requests
	// towards the rate limit. After this duration, the request count for each client is reset.
	maxExpirationRESTAPIsRateLimiter = 1 * time.Minute
)
