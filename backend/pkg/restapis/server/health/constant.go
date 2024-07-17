// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package health

import (
	"sync"
	"time"
)

var (
	// validFiltersKey is the unique key used to store the valid filters in Redis.
	// It is generated once and reused for subsequent cache operations.
	validFiltersKey string

	// validFiltersKeyOnce is a sync.Once instance used to ensure that the validFiltersKey is generated only once.
	validFiltersKeyOnce sync.Once

	// cacheExpiration is the duration for which the valid filters are cached in Redis.
	// After this duration, the cached data will be automatically expired and removed from Redis.
	// Adjust this value based on your application's requirements and the frequency of updates to the valid filters.
	cacheExpiration = 24 * time.Hour
)

const (
	// validFiltersKeyPrefix is the prefix used for the key to store the valid filters for database health checks in Redis.
	// It is a descriptive prefix that clearly indicates the purpose of the cached data.
	// Using a descriptive prefix improves clarity, avoids naming conflicts, and enhances maintainability.
	validFiltersKeyPrefix = "DBHealthCheckTrackID:"

	// ipToKeyPrefix is the prefix used for the key to store the IP to valid filters key mapping in Redis.
	ipToKeyPrefix = "DBHealthCheckIPKey:"
)

// Define Cloudflare formats.
const (
	CloudflareConnectingIPHeader = "Cf-Connecting-IP"
)
