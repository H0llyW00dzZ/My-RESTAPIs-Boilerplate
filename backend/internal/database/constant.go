// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package database

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
	MySQLConnect       = "%s:%s@tcp(%s:%s)/%s"
	dbMYSQL            = "mysql"
	EnvMYSQLDBName     = "DB_DATABASE"
	EnvMYSQLDBPassword = "DB_PASSWORD"
	EnvMYSQLDBUsername = "DB_USERNAME"
	EnvMYSQLDBPort     = "DB_PORT"
	EnvMYSQLDBHost     = "DB_HOST"
)

// Constants for Redis NoSQL name and environment variable names.
//
// Note: Ignore any false positives reported by code scanners, such as secret scanners or other tools, that are not 100% accurate LMAO hahaha.
// For example, "Hardcoded Credentials" might be flagged incorrectly in this case.
const (
	EnvRedisDBName            = "RDB_DATABASE"
	EnvRedisDBPassword        = "RDB_PASSWORD"
	EnvRedisDBPort            = "RDB_PORT"
	EnvRedisDBHost            = "RDB_ADDRESS"
	EnvRedisDBPoolTimeout     = "RDB_POOL_TIMEOUT"
	EnvRedisDBConnMaxIdleTime = "REDIS_MAXCONN_IDLE_TIME"
	EnvRedisDBConnMaxLifeTime = "REDIS_MAXCONN_LIFE_TIME"
)

// Message constants for Redis-related operations.
const (
	MsgRedisHighConnectedClients = "Redis has a high number of connected clients"
	MsgRedisHighMemoryUsage      = "Redis is using a significant amount of memory"
	MsgRedisRecentlyRestarted    = "Redis has been recently restarted"
	MsgRedisHasStaleConnections  = "Redis has %d stale connections."
	MsgRedisFailedToRetrieveInfo = "Failed to retrieve Redis info: %v"
	MsgRedisHighIdleConnections  = "Redis has a high number of idle connections"
	MsgRedisHighPoolUtilization  = "Redis connection pool utilization is high"
)
