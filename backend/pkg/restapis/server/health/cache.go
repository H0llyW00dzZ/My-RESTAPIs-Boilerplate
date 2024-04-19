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
// It attempts to retrieve an existing key for the given IP address from the cache.
// If no existing key is found or an error occurs, a new key is generated and saved.
func generateValidFiltersKey(storage fiber.Storage, ipAddress string) (string, error) {
	var err error
	validFiltersKeyOnce.Do(func() {
		// Attempt to retrieve the existing key for the IP address
		existingKey, errRetrieve := retrieveIPKeyMapping(storage, ipAddress)
		if errRetrieve == nil && existingKey != "" {
			validFiltersKey = existingKey
		} else {
			// If no existing key or error, generate a new one
			randomUUID := uuid.New().String()
			validFiltersKey = validFiltersKeyPrefix + randomUUID
			// Save the new IP-key mapping
			err = saveIPKeyMapping(storage, ipAddress, validFiltersKey)
		}
	})
	return validFiltersKey, err
}

// retrieveValidFiltersFromCache attempts to retrieve valid filters from cache.
// It generates the valid filters key based on the provided IP address and
// retrieves the cached data using the generated key.
// Returns a boolean indicating whether the retrieval was successful and any error encountered.
func retrieveValidFiltersFromCache(storage fiber.Storage, ipAddress string) (bool, error) {
	validFiltersKey, err := generateValidFiltersKey(storage, ipAddress)
	if err != nil {
		log.LogErrorf("Failed to generate valid filters key: %v", err)
		return false, err
	}
	return retrieveFromCache(storage, validFiltersKey, &validFiltersSlice), nil
}

// storeValidFiltersInCache stores the valid filters in cache.
// It generates the valid filters key based on the provided IP address and
// stores the valid filters data in the cache using the generated key.
// Returns any error encountered during the process.
func storeValidFiltersInCache(storage fiber.Storage, ipAddress string) error {
	validFiltersKey, err := generateValidFiltersKey(storage, ipAddress)
	if err != nil {
		log.LogErrorf("Failed to generate valid filters key: %v", err)
		return err
	}
	storeInCache(storage, validFiltersKey, validFiltersSlice, cacheExpiration)
	return nil
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

// saveIPKeyMapping stores the mapping of IP address to the valid filters key.
func saveIPKeyMapping(storage fiber.Storage, ipAddress string, key string) error {
	ipKey := ipToKeyPrefix + ipAddress
	return storage.Set(ipKey, []byte(key), cacheExpiration) // Convert string to []byte
}

// retrieveIPKeyMapping retrieves the mapping of IP address to the valid filters key.
func retrieveIPKeyMapping(storage fiber.Storage, ipAddress string) (string, error) {
	ipKey := ipToKeyPrefix + ipAddress
	data, err := storage.Get(ipKey)
	if err != nil {
		return "", err
	}
	return string(data), nil // Convert []byte to string
}
