// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

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
	err := db.FiberStorage().Set(key, []byte("expired"), 30*time.Minute)
	if err != nil {
		log.LogErrorf("Failed to update Redis cache for expired API key: %v", err)
	}
}

// UpdateCacheWithActiveStatus updates the Redis cache with the active status and expiration time.
func UpdateCacheWithActiveStatus(db database.ServiceAuth, key string, expirationDate time.Time) {
	ttl := time.Until(expirationDate)
	err := db.FiberStorage().Set(key, []byte("active"), ttl)
	if err != nil {
		log.LogErrorf("Failed to update Redis cache for active API key: %v", err)
	}
}
