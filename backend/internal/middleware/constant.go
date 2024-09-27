// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

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

const (
	loggerFormat     = "${time} [${blue}${appName}${reset}] [${green}INFO${reset}] | [${protocol} - ${hostName} - ${status}] | ${latency} - ${unixTime} | [${ip} - ${locals:visitor_uuid}] | ${method} | ${path} | ${error}\n"
	loggerFormatTime = "2006/01/02 15:04:05"
)
