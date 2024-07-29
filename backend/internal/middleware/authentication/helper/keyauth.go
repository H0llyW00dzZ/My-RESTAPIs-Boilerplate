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
func GetAPIKeyStatusFromCache(db database.ServiceAuth, key string) (APIKeyStatus, error) {
	cachedStatus, err := db.FiberStorage().Get(key)
	if err != nil {
		log.LogErrorf("Failed to get API key status from cache: %v", err)
		// Returning APIKeyUnknown is better than returning APIKeyExpired.
		// If the key is not found in the cache, it will return APIKeyUnknown without an error.
		// Otherwise, it will return APIKeyUnknown with an error if an error occurs.
		return APIKeyUnknown, err
	}
	status := string(cachedStatus)
	switch status {
	case APIKeyActive.String():
		return APIKeyActive, nil
	case APIKeyExpired.String():
		return APIKeyExpired, nil
	default:
		return APIKeyUnknown, nil
	}
}

// UpdateCacheWithExpiredStatus updates the Redis cache with the expired status.
func UpdateCacheWithExpiredStatus(db database.ServiceAuth, key string) {
	err := db.FiberStorage().Set(key, []byte(APIKeyExpired.String()), CacheExpiredTTL)
	if err != nil {
		log.LogErrorf("Failed to update Redis cache for expired API key: %v", err)
	}
}

// UpdateCacheWithActiveStatus updates the Redis cache with the active status and expiration time.
func UpdateCacheWithActiveStatus(db database.ServiceAuth, key string, expirationDate time.Time) {
	// Note: This should be set 5 minute as minimum, because it will covered by rate limiter.
	err := db.FiberStorage().Set(key, []byte(APIKeyActive.String()), CacheExpiredTTL)
	if err != nil {
		log.LogErrorf("Failed to update Redis cache for active API key: %v", err)
	}
}
