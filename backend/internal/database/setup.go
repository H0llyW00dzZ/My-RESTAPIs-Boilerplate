// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package database

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	log "h0llyw00dz-template/backend/internal/logger"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver is used for connecting to MySQL databases.
	"github.com/gofiber/fiber/v2"
	redisStorage "github.com/gofiber/storage/redis/v3" // Alias the import to avoid conflict
	"github.com/redis/go-redis/v9"
)

// RedisClientConfig defines the settings needed for Redis client initialization.
type RedisClientConfig struct {
	Address               string
	Port                  int
	Password              string
	Database              int
	PoolTimeout           time.Duration
	ContextTimeoutEnabled bool
	PoolSize              int
}

// FiberRedisClientConfig defines the settings needed for Fiber Redis client initialization.
type FiberRedisClientConfig struct {
	Address  string
	Port     int
	Password string
	Database int
	Reset    bool
	PoolSize int
}

// MySQLConfig defines the settings needed for MySQL client initialization.
type MySQLConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

// InitializeRedisClient initializes and returns a new Redis client.
func InitializeRedisClient(config RedisClientConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Address, config.Port),
		Password: config.Password,
		DB:       config.Database,
		// Note: TLSConfig is optional, but it is recommended for better security, so it's advisable to use it.
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
		PoolTimeout: config.PoolTimeout, // PoolTimeout should already be a time.Duration
	})
	return client
}

// InitializeMySQLDB initializes and returns a new MySQL database client.
func InitializeMySQLDB(config MySQLConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(MySQLConnect, config.Username, config.Password, config.Host, config.Port, config.Database)
	db, err := sql.Open(dbMYSQL, dsn)
	if err != nil {
		return nil, err
	}
	// Set MySQL connection pool parameters.
	// TODO: Refine the MySQL setup and statistics tracking to align with the enhancements previously implemented for Redis.
	db.SetConnMaxLifetime(0) // Connections are not closed due to being idle too long.
	db.SetMaxIdleConns(50)   // Maximum number of connections in the idle connection pool.
	db.SetMaxOpenConns(50)   // Maximum number of open connections to the database.
	// Log the successful database connection
	log.LogInfof(MsgDBConnected, dbname)
	return db, nil
}

// InitializeRedisStorage initializes and returns a new Redis storage instance
// for use with Fiber middlewares such as rate limiting.
func InitializeRedisStorage(config FiberRedisClientConfig) fiber.Storage {
	storage := redisStorage.New(redisStorage.Config{
		Host:     config.Address,
		Port:     config.Port,
		Password: config.Password,
		Database: config.Database,
		Reset:    config.Reset,
		// Note: TLSConfig is optional, but it is recommended for better security, so it's advisable to use it.
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
		PoolSize: config.PoolSize, // Adjust the pool size as necessary.
	})
	return storage
}
