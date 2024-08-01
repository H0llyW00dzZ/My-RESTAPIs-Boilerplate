// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package helper

import "time"

// CacheExpiredTTL is the time-to-live (TTL) duration for expired items in the cache.
// It determines how long an expired item should remain in the cache before being evicted.
// In this case, expired items will be kept in the cache for 5 minutes.
const (
	CacheExpiredTTL = 5 * time.Minute
)

// APIKeyStatus represents the status of an API key.
type APIKeyStatus int

const (
	// APIKeyUnknown represents an unknown API key status.
	APIKeyUnknown APIKeyStatus = iota

	// APIKeyActive represents an active API key status.
	APIKeyActive

	// APIKeyExpired represents an expired API key status.
	APIKeyExpired
)

// String returns the string representation of the APIKeyStatus.
func (s APIKeyStatus) String() string {
	switch s {
	case APIKeyActive:
		return "active"
	case APIKeyExpired:
		return "expired"
	default:
		return "unknown"
	}
}

const (
	prefixKey = "key_auth_session:"
)

// APIKeyData represents the structure of the API key data stored in the cache.
// It includes the following fields:
// 	- Identifier: The unique identifier associated with the API key.
// 	- APIKey: The actual API key value.
// 	- Status: The status of the API key (e.g., "active", "expired").
//
// Note: This structure is only for Redis, as it is used solely for caching + better performance.
// for relational database (MySQL) marked as TODO.
type APIKeyData struct {
	Identifier string `json:"identifier"`
	APIKey     string `json:"apikey"`
	Status     string `json:"status"`
}

// KeyAuthSessData is a type alias for map[string]any.
// It represents a map of key-value pairs used to store session data related to key authentication.
// The keys are strings, and the values can be of any type (any).
// This type alias provides a convenient way to work with session data in a flexible manner.
type KeyAuthSessData map[string]any
