// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package helper

import (
	"h0llyw00dz-template/backend/internal/database"
	log "h0llyw00dz-template/backend/internal/logger"
	"time"

	"github.com/bytedance/sonic"
)

// GetAPIKeyStatusFromCache retrieves the API key status from the Redis cache.
func GetAPIKeyStatusFromCache(db database.ServiceAuth, identifier, key string) (APIKeyStatus, error) {
	cacheKey := prefixKey + identifier
	cachedData, err := db.FiberStorage().Get(cacheKey)
	if err != nil {
		log.LogErrorf("Failed to get API key status from cache: %v", err)
		// Returning APIKeyUnknown is better than returning APIKeyExpired.
		// If the key is not found in the cache, it will return APIKeyUnknown without an error.
		// Otherwise, it will return APIKeyUnknown with an error if an error occurs.
		return APIKeyUnknown, err
	}

	if len(cachedData) == 0 {
		// Return APIKeyUnknown without an error if the key is not found in the cache
		return APIKeyUnknown, nil
	}

	// Note: Custom JSON encoder/decoder configuration, similar to what Fiber currently supports,
	// currently unavailable in this enhancement due to its focus on better performance.
	var data APIKeyData
	if err = sonic.Unmarshal([]byte(cachedData), &data); err != nil {
		log.LogErrorf("Failed to unmarshal API key data from cache: %v", err)
		return APIKeyUnknown, err
	}

	switch data.Status {
	case APIKeyActive.String():
		return APIKeyActive, nil
	case APIKeyExpired.String():
		return APIKeyExpired, nil
	default:
		return APIKeyUnknown, nil
	}
}

// UpdateCacheWithExpiredStatus updates the Redis cache with the expired status.
func UpdateCacheWithExpiredStatus(db database.ServiceAuth, identifier, key string) {
	cacheKey := prefixKey + identifier
	data := APIKeyData{
		Identifier: identifier,
		APIKey:     key,
		Status:     APIKeyExpired.String(),
	}

	// Note: Custom JSON encoder/decoder configuration, similar to what Fiber currently supports,
	// currently unavailable in this enhancement due to its focus on better performance.
	jsonData, err := sonic.Marshal(data)
	if err != nil {
		log.LogErrorf("Failed to marshal API key data for cache: %v", err)
		return
	}

	// Note: This should be set 5 minute as minimum, because it will covered by rate limiter.
	if err := db.FiberStorage().Set(cacheKey, jsonData, CacheExpiredTTL); err != nil {
		log.LogErrorf("Failed to update Redis cache for expired API key: %v", err)
	}
}

// UpdateCacheWithActiveStatus updates the Redis cache with the active status and expiration time.
func UpdateCacheWithActiveStatus(db database.ServiceAuth, identifier, key string, expirationDate time.Time) {
	cacheKey := prefixKey + identifier
	data := APIKeyData{
		Identifier: identifier,
		APIKey:     key,
		Status:     APIKeyActive.String(),
	}

	// Note: Custom JSON encoder/decoder configuration, similar to what Fiber currently supports,
	// currently unavailable in this enhancement due to its focus on better performance.
	jsonData, err := sonic.Marshal(data)
	if err != nil {
		log.LogErrorf("Failed to marshal API key data for cache: %v", err)
		return
	}

	// Note: This should be set 5 minute as minimum, because it will covered by rate limiter.
	if err := db.FiberStorage().Set(cacheKey, jsonData, CacheExpiredTTL); err != nil {
		log.LogErrorf("Failed to update Redis cache for active API key: %v", err)
	}
}
