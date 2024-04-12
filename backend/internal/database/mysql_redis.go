// Copyright (c) 2024 H0llyW00dz & Melkeydev (go-blueprint author) All rights reserved.
//
// License: BSD 3-Clause License, MIT License

package database

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	log "h0llyw00dz-template/backend/internal/logger"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	redisStorage "github.com/gofiber/storage/redis/v3" // Alias the import to avoid conflict

	_ "github.com/go-sql-driver/mysql" // MySQL driver is used for connecting to MySQL databases.
	_ "github.com/joho/godotenv/autoload"
)

// Service defines the interface for database operations that can be performed.
type Service interface {
	// Health checks the health of the database connection.
	Health() map[string]string

	// Close terminates the database connection.
	Close() error

	// Exec executes a SQL query with the provided arguments and returns the result.
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	// BeginTx starts a new database transaction with the specified options.
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)

	// QueryRow executes a query that is expected to return at most one row and scans that row into the provided destination.
	QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row

	// FiberStorage returns the [fiber.Storage] interface for storage middleware.
	FiberStorage() fiber.Storage
}

// service is a concrete implementation of the Service interface.
type service struct {
	db          *sql.DB
	ratelimiter fiber.Storage
}

// dbConfig holds the environment variables for the database connection.
var (
	dbname        = os.Getenv(EnvMYSQLDBName)
	password      = os.Getenv(EnvMYSQLDBPassword)
	username      = os.Getenv(EnvMYSQLDBUsername)
	port          = os.Getenv(EnvMYSQLDBPort)
	host          = os.Getenv(EnvMYSQLDBHost)
	redisAddress  = os.Getenv(EnvRedisDBHost)
	redisPort     = os.Getenv(EnvRedisDBPort)
	redisPassword = os.Getenv(EnvRedisDBPassword)
	redisDatabase = os.Getenv(EnvRedisDBName)
	dbInstance    *service
)

// New creates a new instance of the Service interface.
// It opens a connection to the MySQL database using the environment variables
// and sets up the connection pool configuration.
func New() Service {
	if dbInstance != nil {
		return dbInstance
	}

	// Parse the Redis database index from the environment variable.
	redisDB, err := strconv.Atoi(redisDatabase)
	if err != nil {
		log.LogFatal("Invalid Redis database index:", err)
	}
	portDB, err := strconv.Atoi(redisPort)
	if err != nil {
		log.LogFatal("Invalid Redis database index:", err)
	}

	// Opening a driver typically will not attempt to connect to the database.
	db, err := sql.Open(dbMYSQL, fmt.Sprintf(MySQLConnect, username, password, host, port, dbname))
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.LogFatal(err)
	}

	// Log the successful database connection
	log.LogInfof(MsgDBConnected, dbname)

	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)

	// Initialize the Redis storage for rate limiting or any that needed in middleware
	redisStorage := redisStorage.New(redisStorage.Config{
		Host:     redisAddress,
		Port:     portDB,
		Password: redisPassword,
		Database: redisDB,
		Reset:    false, // Set to true if you want to clear the storage upon connection
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
		PoolSize: 10 * runtime.GOMAXPROCS(10), // Adjust the pool size as needed
	})

	dbInstance = &service{
		db:          db,
		ratelimiter: redisStorage, // Use the Redis storage for rate limiting or any that needed in middleware
	}

	return dbInstance
}

// Close closes the database connection.
func (s *service) Close() error {
	log.LogInfof(MsgDBDisconnected, dbname)
	return s.db.Close()
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf(ErrDBDown, err)
		stats["message"] = fmt.Sprintf("%v", err)
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = MsgDBItsHealthy

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 40 is the max for this example
		stats["message"] = MsgDBHeavyLoad
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = MsgDBHighWaitEvents
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = MsgDBManyIdleConnections
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = MsgDBManyMaxLifetimeClosures
	}

	return stats
}

// Exec executes a SQL query with the provided arguments.
func (s *service) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}

func (s *service) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return s.db.BeginTx(ctx, opts)
}

// QueryRow executes a query that is expected to return at most one row.
func (s *service) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return s.db.QueryRowContext(ctx, query, args...)
}

// FiberStorage returns the [fiber.Storage] interface for fiber storage middleware.
func (s *service) FiberStorage() fiber.Storage {
	return s.ratelimiter
}
