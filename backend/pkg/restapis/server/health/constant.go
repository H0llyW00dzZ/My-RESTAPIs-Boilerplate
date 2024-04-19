// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package health

import (
	"sync"
	"time"
)

var (
	// validFiltersKeyPrefix is the prefix used for the key to store the valid filters for database health checks in Redis.
	// It is a descriptive prefix that clearly indicates the purpose of the cached data.
	// Using a descriptive prefix improves clarity, avoids naming conflicts, and enhances maintainability.
	validFiltersKeyPrefix = "Database health check valid filters TrackID:"

	// validFiltersKey is the unique key used to store the valid filters in Redis.
	// It is generated once and reused for subsequent cache operations.
	validFiltersKey string

	// validFiltersKeyOnce is a sync.Once instance used to ensure that the validFiltersKey is generated only once.
	validFiltersKeyOnce sync.Once

	// cacheExpiration is the duration for which the valid filters are cached in Redis.
	// After this duration, the cached data will be automatically expired and removed from Redis.
	// Adjust this value based on your application's requirements and the frequency of updates to the valid filters.
	cacheExpiration = 1 * time.Hour
)
