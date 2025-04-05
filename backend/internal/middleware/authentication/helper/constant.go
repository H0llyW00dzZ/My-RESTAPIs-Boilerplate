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
//   - Identifier: The unique identifier associated with the API key.
//   - APIKey: The actual API key value.
//   - Status: The status of the API key (e.g., "active", "expired").
//   - Authorization: The authorization data of the API key.
//
// Note: This structure is only for Redis/Valkey, as it is used solely for caching + better performance.
// for relational database (MySQL) marked as TODO.
type APIKeyData struct {
	Identifier    string            `json:"identifier,omitempty"`
	APIKey        string            `json:"apikey"`
	Status        string            `json:"status"`
	Authorization AuthorizationData `json:"authorization"`
}

// KeyAuthSessData is a type alias for map[string]any.
// It represents a map of key-value pairs used to store session data related to key authentication.
// The keys are strings, and the values can be of any type (any).
// This type alias provides a convenient way to work with session data in a flexible manner.
//
// Note: This is currently unused because key-auth and session middleware logic are bound to Fiber storage.
// However, for other cache handlers (without being bound to Fiber storage), this can be useful.
type KeyAuthSessData map[string]any

// AuthorizationData represents the authorization data of an API key.
// It includes the following fields:
//   - AuthTime: The time of the last authorization.
//   - ExpiredTime: The expiration time of the API key.
//   - Signature: The signature data associated with the API key, which can be of any type (e.g., ECDSA signature that can be used for enhance security purpose or other purpose).
//
// Note: Current format, and may subject to changed:
//
//	{
//		"authorization": {
//		  "time": "2024-08-01T21:24:28.4352685Z",
//		  "apikey_expired_time": "2024-10-26T21:18:17Z",
//		  "signature": "..."
//		}
//	}
type AuthorizationData struct {
	AuthTime time.Time `json:"time"`
	// Note: This expiration time is retrieved from the relational database (MySQL).
	// The performance speed might be somewhat slow (taking an average of 1s response time in the frontend) during the first query due to the relational database (always slow).
	// However, when it hits Redis/Valkey and is released into cookies with encryption, the speed can be faster (possibly 0ms ~ 1ms response time).
	ExpiredTime time.Time `json:"apikey_expired_time"`
	Signature   any       `json:"signature,omitempty"`
}
