// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package database

import "time"

// Error and message constants for database-related operations.
const (
	ErrDBDown                    = "db down: %v"
	MsgDBItsHealthy              = "It's healthy"
	MsgDBDisconnected            = "Disconnected from database: %s"
	MsgDBHeavyLoad               = "The database is experiencing heavy load."
	MsgDBHighWaitEvents          = "The database has a high number of wait events, indicating potential bottlenecks."
	MsgDBManyIdleConnections     = "Many idle connections are being closed, consider revising the connection pool settings."
	MsgDBManyMaxLifetimeClosures = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	MsgDBNotAccessible           = "The database is not accessible, please check the connection and credentials."
)

// Constants for MySQL connection string and response messages.
//
// Note: Ignore any false positives reported by code scanners, such as secret scanners or other tools, that are not 100% accurate LMAO hahaha.
// For example, "Hardcoded Credentials" might be flagged incorrectly in this case.
const (
	MySQLConnect = "%s:%s@tcp(%s:%s)/%s"
	dbMYSQL      = "mysql"
)

// Message constants for Redis-related operations.
const (
	// MsgRedisHighConnectedClients indicates that Redis has a high number of connected clients.
	MsgRedisHighConnectedClients = "Redis has a high number of connected clients"

	// MsgRedisHighMemoryUsage indicates that Redis is using a significant amount of memory.
	MsgRedisHighMemoryUsage = "Redis is using a significant amount of memory"

	// MsgRedisRecentlyRestarted indicates that Redis has been recently restarted.
	MsgRedisRecentlyRestarted = "Redis has been recently restarted"

	// MsgRedisHasStaleConnections indicates the number and percentage of stale connections in Redis.
	MsgRedisHasStaleConnections = "Redis has %d stale (%.2f%% High) connections."

	// MsgRedisFailedToRetrieveInfo indicates a failure to retrieve Redis information.
	MsgRedisFailedToRetrieveInfo = "Failed to retrieve Redis info: %v"

	// MsgRedisHighIdleConnections indicates that Redis has a high number of idle connections.
	MsgRedisHighIdleConnections = "Redis has a high number of idle connections"

	// MsgRedisHighPoolUtilization indicates that the Redis connection pool utilization is high.
	MsgRedisHighPoolUtilization = "Redis connection pool utilization is high"

	// MsgRedisHighPoolBottleneck indicates a critical bottleneck in the Redis connection pool.
	MsgRedisHighPoolBottleneck = "CRITICAL: Redis connection pool bottlenecked! Utilization exceeds 100%. Possible misconfiguration!"
)

// Default context timeouts for operations.
const (
	// DefaultCtxTimeout is the default timeout for context operations, set to 5 minutes.
	DefaultCtxTimeout = 5 * time.Minute

	// DefaultBackupCtxTimeout is the default timeout for backup context operations, set to 30 minutes.
	DefaultBackupCtxTimeout = 30 * time.Minute
)

// Miscellaneous constants.
const (
	// nullObject represents a null value as a string.
	nullObject = "NULL"
)
