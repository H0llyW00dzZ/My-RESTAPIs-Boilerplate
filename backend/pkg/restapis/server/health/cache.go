// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package health

import (
	log "h0llyw00dz-template/backend/internal/logger"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// generateValidFiltersKey generates a unique key for storing the valid filters in Redis.
// It combines the validFiltersKeyPrefix with a random UUID to ensure uniqueness.
// The key is generated only once and reused for subsequent cache operations.
func generateValidFiltersKey() string {
	validFiltersKeyOnce.Do(func() {
		randomUUID := uuid.New().String()
		validFiltersKey = validFiltersKeyPrefix + randomUUID
	})
	return validFiltersKey
}

// retrieveValidFiltersFromCache attempts to retrieve valid filters from cache.
func retrieveValidFiltersFromCache(storage fiber.Storage) bool {
	validFiltersKey := generateValidFiltersKey()
	return retrieveFromCache(storage, validFiltersKey, &validFiltersSlice)
}

// storeValidFiltersInCache stores the valid filters in cache.
func storeValidFiltersInCache(storage fiber.Storage) {
	validFiltersKey := generateValidFiltersKey()
	storeInCache(storage, validFiltersKey, validFiltersSlice, cacheExpiration)
}

// retrieveFromCache attempts to retrieve data from cache based on the provided key.
func retrieveFromCache(storage fiber.Storage, key string, data interface{}) bool {
	cacheData, err := storage.Get(key)
	if err != nil {
		if err == fiber.ErrNotFound {
			log.LogInfof("Cache data not found for key: %s", key)
		} else {
			log.LogErrorf("Failed to retrieve cache data for key: %s, error: %v", key, err)
		}
		return false
	}

	if len(cacheData) == 0 {
		log.LogInfof("Cache data is empty for key: %s", key)
		return false
	}

	if err := sonic.Unmarshal(cacheData, data); err != nil {
		log.LogErrorf("Failed to unmarshal cache data for key: %s, error: %v", key, err)
		// Clear the invalid cached data from Redis
		if err := storage.Delete(key); err != nil {
			log.LogErrorf("Failed to delete invalid cached data for key: %s, error: %v", key, err)
		}
		return false
	}

	log.LogInfof("Cache data retrieved for key: %s", key)
	return true
}

// storeInCache stores the provided data in cache with the specified key and expiration.
func storeInCache(storage fiber.Storage, key string, data interface{}, expiration time.Duration) {
	cacheData, err := sonic.Marshal(data)
	if err != nil {
		log.LogErrorf("Failed to marshal cache data for key: %s, error: %v", key, err)
		return
	}

	if err := storage.Set(key, cacheData, expiration); err != nil {
		log.LogErrorf("Failed to store cache data for key: %s, error: %v", key, err)
		return
	}

	log.LogInfof("Cache data stored for key: %s with expiration: %v", key, expiration)
}
