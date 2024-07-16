// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package helper

import (
	"h0llyw00dz-template/backend/internal/database"
	log "h0llyw00dz-template/backend/internal/logger"
	"time"
)

// GetAPIKeyStatusFromCache retrieves the API key status from the Redis cache.
func GetAPIKeyStatusFromCache(db database.ServiceAuth, key string) (string, error) {
	cachedStatus, err := db.FiberStorage().Get(key)
	if err != nil {
		log.LogErrorf("Failed to get API key status from cache: %v", err)
		return "", err
	}
	return string(cachedStatus), nil
}

// UpdateCacheWithExpiredStatus updates the Redis cache with the expired status.
func UpdateCacheWithExpiredStatus(db database.ServiceAuth, key string) {
	err := db.FiberStorage().Set(key, []byte("expired"), 5*time.Minute)
	if err != nil {
		log.LogErrorf("Failed to update Redis cache for expired API key: %v", err)
	}
}

// UpdateCacheWithActiveStatus updates the Redis cache with the active status and expiration time.
func UpdateCacheWithActiveStatus(db database.ServiceAuth, key string, expirationDate time.Time) {
	// Note: This should be set 5 minute as minimum, because it will covered by rate limiter.
	ttl := 5 * time.Minute
	err := db.FiberStorage().Set(key, []byte("active"), ttl)
	if err != nil {
		log.LogErrorf("Failed to update Redis cache for active API key: %v", err)
	}
}
