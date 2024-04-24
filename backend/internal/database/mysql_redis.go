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
	"github.com/redis/go-redis/v9"

	_ "github.com/go-sql-driver/mysql" // MySQL driver is used for connecting to MySQL databases.
	_ "github.com/joho/godotenv/autoload"
)

// Service defines the interface for database operations that can be performed.
type Service interface {
	// Health checks the health of the database connection.
	Health(filter string) map[string]string

	// Close terminates the database connection.
	Close() error

	// Exec executes a SQL query with the provided arguments and returns the result.
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	// ExecWithoutRow executes a query that doesn't return any rows, such as
	// CREATE, ALTER, DROP, INSERT, UPDATE, or DELETE statements.
	// It's useful for initializing database schemas, migrations, or any other
	// queries that don't require retrieving rows.
	//
	// Example Usage:
	//
	//	ctx := context.Background()
	//	query := "CREATE TABLE users (id INT, name VARCHAR(255))"
	//	err := db.ExecWithoutRow(ctx, query)
	//	if err != nil {
	//	    // Handle the error
	//	}
	ExecWithoutRow(ctx context.Context, query string, args ...interface{}) error

	// EnsureTransactionClosure is a deferred function to handle transaction rollback or commit.
	// It can be used in goroutines along with an interval, such as in a scheduler.
	//
	// Example Usage:
	//
	//	func schedulerTask(interval time.Duration) {
	//	    for {
	//	        ctx := context.Background()
	//	        tx, err := db.BeginTx(ctx, nil)
	//	        if err != nil {
	//	            log.LogErrorf("Failed to start transaction: %v", err)
	//	            continue
	//	        }
	//	        defer db.EnsureTransactionClosure(tx, &err)
	//
	//	        // Perform database operations within the transaction
	//	        // ...
	//
	//	        time.Sleep(interval)
	//	    }
	//	}
	//
	//	go schedulerTask(1 * time.Minute)
	//
	// In the example above, EnsureTransactionClosure is used within a goroutine that runs
	// a scheduler task. The task is executed at a specified interval.
	//
	// The goroutine starts a new transaction using BeginTx in each iteration. EnsureTransactionClosure
	// is deferred immediately after starting the transaction to handle the transaction closure.
	//
	// If an error occurs during the transaction or if a panic is encountered, EnsureTransactionClosure
	// will rollback the transaction. If no errors occur, it will commit the transaction.
	//
	// The function also logs any errors that occur during the rollback or commit process.
	//
	// Note: Make sure to handle errors appropriately and adjust the interval as needed
	// based on your specific requirements.
	EnsureTransactionClosure(tx *sql.Tx, err *error)

	// BeginTx starts a new database transaction with the specified options.
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)

	// QueryRow executes a query that is expected to return at most one row and scans that row into the provided destination.
	QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row

	// FiberStorage returns the [fiber.Storage] interface for storage middleware.
	FiberStorage() fiber.Storage

	// ScanAndDel uses the Redis SCAN command to iterate over a set of keys and delete them.
	// It's particularly useful for deleting keys with a common pattern.
	//
	// Example Usage:
	//
	//	if err := db.ScanAndDel("gopher_key:*"); err != nil {
	//		Log.LogErrorf("Failed to clear gopher keys cache: %v", err)
	//		return err
	//	}
	ScanAndDel(pattern string) error

	// PrepareInsertStatement prepares a SQL insert statement for the transaction.
	PrepareInsertStatement(ctx context.Context, tx *sql.Tx, query string) (*sql.Stmt, error)

	// ScanKeys returns a slice of keys for a given pattern starting at the cursor.
	ScanKeys(ctx context.Context, pattern string, cursor uint64) ([]string, uint64, error)

	// DeleteKeys deletes a slice of keys from Redis and returns the updated count.
	DeleteKeys(ctx context.Context, keys []string, totalDeleted int) (int, error)
}

// service is a concrete implementation of the Service interface.
type service struct {
	db          *sql.DB
	rdb         fiber.Storage
	redisClient *redis.Client
}

// dbConfig holds the environment variables for the database connection.
var (
	dbname           = os.Getenv(EnvMYSQLDBName)
	password         = os.Getenv(EnvMYSQLDBPassword)
	username         = os.Getenv(EnvMYSQLDBUsername)
	port             = os.Getenv(EnvMYSQLDBPort)
	host             = os.Getenv(EnvMYSQLDBHost)
	redisAddress     = os.Getenv(EnvRedisDBHost)
	redisPort        = os.Getenv(EnvRedisDBPort)
	redisPassword    = os.Getenv(EnvRedisDBPassword)
	redisDatabase    = os.Getenv(EnvRedisDBName)
	redisPoolTimeout = os.Getenv(EnvRedisDBPoolTimeout)
	dbInstance       *service
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

	// Parse redis port from the environment variable
	portDB, err := strconv.Atoi(redisPort)
	if err != nil {
		log.LogFatal("Invalid Redis database index:", err)
	}

	// Parse pool timeout from the environment variable
	poolTimeout, err := time.ParseDuration(redisPoolTimeout)
	if err != nil {
		log.LogFatal("Invalid Redis pool timeout value:", err)
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

	// Set MySQL connection pool parameters.
	// TODO: Refine the MySQL setup and statistics tracking to align with the enhancements previously implemented for Redis.
	db.SetConnMaxLifetime(0) // Connections are not closed due to being idle too long.
	db.SetMaxIdleConns(50)   // Maximum number of connections in the idle connection pool.
	db.SetMaxOpenConns(50)   // Maximum number of open connections to the database.

	// Note: This configuration is better for starters and provides stability.
	// Tested on Node Spec:
	// 2x vCPU
	// 4x ~ 8x Compute
	// 1 GB RAM
	maxConnections := 2 * runtime.NumCPU()

	// Create a new Redis client for health checks or for any other needs in middleware that do not involve using Fiber's storage.
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisAddress, redisPort),
		Password: redisPassword,
		DB:       redisDB,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
		PoolSize:              maxConnections,
		PoolTimeout:           poolTimeout,
		ContextTimeoutEnabled: true,
	})

	// Initialize Redis storage for rate limiting or for any other needs in middleware.
	redisStorage := redisStorage.New(redisStorage.Config{
		Host:     redisAddress,
		Port:     portDB,
		Password: redisPassword,
		Database: redisDB,
		Reset:    false, // Set to true to clear the storage upon establishing a connection.
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
		PoolSize: maxConnections, // Adjust the pool size as necessary.
	})

	dbInstance = &service{
		db:          db,
		rdb:         redisStorage, // Use the Redis storage for rate limiting or any that needed in middleware
		redisClient: redisClient,
	}

	return dbInstance
}

// Close closes the database connection and the Redis client.
func (s *service) Close() error {
	// Close the Redis client connection
	if err := s.redisClient.Close(); err != nil {
		log.LogErrorf("Error closing Redis client: %v", err)
		// Don't return yet because we also need to close the SQL database connection.
	}

	// Log information about closing the Redis connection
	log.LogInfo("Redis connection closed.")

	// Close the SQL database connection
	if err := s.db.Close(); err != nil {
		log.LogErrorf("Error closing database connection: %v", err)
		return err
	}

	// Log information about closing the database connection
	log.LogInfof(MsgDBDisconnected, dbname)

	return nil
}

// Health checks the health of the database and Redis connections.
// It returns a map with keys indicating various health statistics.
func (s *service) Health(filter string) map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// TODO: Improve by using a "Map of Functions" to reduce complexity (caused by if-else statements) when handling multiple databases in the future.
	if filter == "" || filter == "mysql" {
		stats = s.checkMySQLHealth(ctx, stats)
	}

	if filter == "" || filter == "redis" {
		stats = s.checkRedisHealth(ctx, stats)
	}

	return stats
}

// checkMySQLHealth checks the health of the MySQL database and adds the relevant statistics to the stats map.
func (s *service) checkMySQLHealth(ctx context.Context, stats map[string]string) map[string]string {
	// Ping the MySQL database
	err := s.db.PingContext(ctx)
	if err != nil {
		// Note: While using `log.Fatal` is an option, it is not recommended for this REST API.
		// These APIs are designed for large-scale applications with complex infrastructure rather than
		// small systems reliant on a single database. Using `log.Fatal` can prematurely terminate
		// the service, which is undesirable in a distributed and resilient application environment.
		stats["mysql_status"] = "down"
		stats["mysql_error"] = fmt.Sprintf(ErrDBDown, err)
	} else {
		// MySQL is up, add more statistics
		stats["mysql_status"] = "up"
		stats["mysql_message"] = MsgDBItsHealthy

		// Get MySQL database stats (like open connections, in use, idle, etc.)
		dbStats := s.db.Stats()
		stats["mysql_open_connections"] = strconv.Itoa(dbStats.OpenConnections)
		stats["mysql_in_use"] = strconv.Itoa(dbStats.InUse)
		stats["mysql_idle"] = strconv.Itoa(dbStats.Idle)
		stats["mysql_wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
		stats["mysql_wait_duration"] = dbStats.WaitDuration.String()
		stats["mysql_max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
		stats["mysql_max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

		// Evaluate MySQL stats to provide a health message
		stats = s.evaluateMySQLStats(dbStats, stats)
	}

	return stats
}

// evaluateMySQLStats evaluates the MySQL database statistics and updates the stats map with the appropriate health message.
func (s *service) evaluateMySQLStats(dbStats sql.DBStats, stats map[string]string) map[string]string {
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["mysql_message"] = MsgDBHeavyLoad
	}

	if dbStats.WaitCount > 1000 {
		stats["mysql_message"] = MsgDBHighWaitEvents
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["mysql_message"] = MsgDBManyIdleConnections
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["mysql_message"] = MsgDBManyMaxLifetimeClosures
	}

	return stats
}

// checkRedisHealth checks the health of the Redis server and adds the relevant statistics to the stats map.
func (s *service) checkRedisHealth(ctx context.Context, stats map[string]string) map[string]string {
	// Ping the Redis server
	pong, err := s.redisClient.Ping(ctx).Result()
	if err != nil {
		// Note: While using `log.Fatal` is an option, it is not recommended for this REST API.
		// These APIs are designed for large-scale applications with complex infrastructure rather than
		// small systems reliant on a single database. Using `log.Fatal` can prematurely terminate
		// the service, which is undesirable in a distributed and resilient application environment.
		stats["redis_status"] = "down"
		stats["redis_error"] = fmt.Sprintf(ErrDBDown, err)
	} else {
		// Redis is up
		stats["redis_status"] = "up"
		stats["redis_message"] = MsgDBItsHealthy
		stats["redis_ping_response"] = pong

		// Get Redis server information
		info, err := s.redisClient.Info(ctx).Result()
		if err != nil {
			stats["redis_info_error"] = fmt.Sprintf("Failed to retrieve Redis info: %v", err)
		} else {
			// Parse the Redis info response
			redisInfo := parseRedisInfo(info)
			stats["redis_version"] = redisInfo["redis_version"]
			stats["redis_mode"] = redisInfo["redis_mode"]
			stats["redis_connected_clients"] = redisInfo["connected_clients"]
			stats["redis_used_memory"] = redisInfo["used_memory"]
			stats["redis_used_memory_peak"] = redisInfo["used_memory_peak"]
			stats["redis_uptime_in_seconds"] = redisInfo["uptime_in_seconds"]

			// Get the pool stats of the Redis client
			poolStats := s.redisClient.PoolStats()

			// Extract the number of hits (free times) connections in the pool
			stats["redis_hits_connections"] = strconv.FormatUint(uint64(poolStats.Hits), 10)

			// Extract the number of misses (not found) connections in the pool
			stats["redis_misses_connections"] = strconv.FormatUint(uint64(poolStats.Misses), 10)

			// Extract the number of timeouts (wait a timeouts) connections in the pool
			stats["redis_timeouts_connections"] = strconv.FormatUint(uint64(poolStats.Timeouts), 10)

			// Extract the total number of connections in the pool
			stats["redis_total_connections"] = strconv.FormatUint(uint64(poolStats.TotalConns), 10)

			// Extract the number of idle connections in the pool
			stats["redis_idle_connections"] = strconv.FormatUint(uint64(poolStats.IdleConns), 10)

			// Extract the number of stale connections in the pool
			stats["redis_stale_connections"] = strconv.FormatUint(uint64(poolStats.StaleConns), 10)

			// Extract the number of active connections (TotalConns - IdleConns gives us the ActiveConns)
			activeConns := poolStats.TotalConns - poolStats.IdleConns
			stats["redis_active_connections"] = strconv.FormatUint(uint64(activeConns), 10)

			// Get the used memory of the Redis server in bytes
			stats["redis_max_memory"] = redisInfo["maxmemory"] // Raw max memory in bytes

			// Get the pool size percentage
			poolSize := s.redisClient.Options().PoolSize
			connectedClients, _ := strconv.Atoi(redisInfo["connected_clients"])
			poolSizePercentage := float64(connectedClients) / float64(poolSize) * 100
			stats["redis_pool_size_percentage"] = fmt.Sprintf("%.2f%%", poolSizePercentage)

			// Evaluate Redis stats to provide a health message
			stats = s.evaluateRedisStats(redisInfo, stats)
		}
	}

	return stats
}

// evaluateRedisStats evaluates the Redis server statistics and updates the stats map with the appropriate health message.
func (s *service) evaluateRedisStats(redisInfo, stats map[string]string) map[string]string {
	// Retrieve the pool size from the Redis client configuration
	poolSize := s.redisClient.Options().PoolSize

	// Get the pool stats of the Redis client
	poolStats := s.redisClient.PoolStats()

	// Check the number of connected clients
	connectedClients, _ := strconv.Atoi(redisInfo["connected_clients"])

	// Determine a high connection threshold, let's say 80% of the pool size because 20% is must be free (genius thinking 🤪)
	highConnectionThreshold := float64(poolSize) * 0.8

	// Check if connected clients exceed the high connection threshold
	if float64(connectedClients) > highConnectionThreshold {
		stats["redis_message"] = MsgRedisHighConnectedClients
	}

	// Check for any stale connections and set a warning if any are found
	// Note: This can sometimes happen, especially with Unix clients.
	// It might be less common when connecting from Linux (client) to Linux (server), as opposed to Unix (client) to Linux (server).
	if poolStats.StaleConns > 0 {
		staleConns := uint64(poolStats.StaleConns)
		stats["redis_message"] = fmt.Sprintf(MsgRedisHasStaleConnections, staleConns)
	}

	// Check if used memory is close to the maximum memory
	// Note: this is now more dynamic instead of hardcoded static 🤪
	usedMemory, _ := strconv.ParseInt(redisInfo["used_memory"], 10, 64)
	maxMemory, _ := strconv.ParseInt(redisInfo["maxmemory"], 10, 64)
	if maxMemory > 0 {
		// Calculate the percentage of used memory
		usedMemoryPercentage := float64(usedMemory) / float64(maxMemory) * 100
		// If used memory is greater than or equal to 90% of the maximum memory,
		// set the redis_health_message to indicate high memory usage
		if usedMemoryPercentage >= 90 {
			stats["redis_message"] = MsgRedisHighMemoryUsage
		}
	}

	// Check the uptime of the Redis server
	uptimeInSeconds, _ := strconv.ParseInt(redisInfo["uptime_in_seconds"], 10, 64)
	if uptimeInSeconds < 3600 { // 1 hour
		stats["redis_message"] = MsgRedisRecentlyRestarted
	}

	return stats
}

// Exec executes a SQL query with the provided arguments.
func (s *service) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}

// ExecWithoutRow executes a query without returning any rows.
//
// Note: This method is different from "Exec". Unlike "Exec", it doesn't return "sql.Result".
// This method is better suited for initializing database schemas or running migrations before the app starts.
func (s *service) ExecWithoutRow(ctx context.Context, query string, args ...interface{}) error {
	_, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		log.LogErrorf("Error executing query: %v", err)
		return err
	}
	return nil
}

// EnsureTransactionClosure is a deferred function to handle transaction rollback or commit.
//
// Note: This method requires the database service to be initialized in the "func init()"
// before it can be used.
//
// Example usage in init function:
//
//	func init() {
//	    // Initialize the database service
//	    db = database.New().GopherService()
//	}
//
// For example usage, see the documentation for the EnsureTransactionClosure method in the Service Interface.
func (s *service) EnsureTransactionClosure(tx *sql.Tx, err *error) {
	if p := recover(); p != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.LogErrorf("Error rolling back transaction: %v", rollbackErr)
		}
		panic(p) // re-throw panic after Rollback
	} else if *err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.LogErrorf("Error rolling back transaction: %v", rollbackErr)
		}
	} else {
		if commitErr := tx.Commit(); commitErr != nil {
			log.LogErrorf("Error committing transaction: %v", commitErr)
			*err = commitErr // capture commit error
		}
	}
}

// BeginTx starts a new transaction.
func (s *service) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return s.db.BeginTx(ctx, opts)
}

// QueryRow executes a query that is expected to return at most one row.
func (s *service) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return s.db.QueryRowContext(ctx, query, args...)
}

// FiberStorage returns the [fiber.Storage] interface for fiber storage middleware.
func (s *service) FiberStorage() fiber.Storage {
	return s.rdb
}

// ScanAndDel uses the Redis SCAN command to iterate over a set of keys and delete them.
func (s *service) ScanAndDel(pattern string) error {
	ctx := context.Background()
	var totalDeleted int
	var scanComplete bool
	var err error

	for !scanComplete {
		var keys []string
		var cursor uint64

		keys, cursor, err = s.ScanKeys(ctx, pattern, cursor)
		if err != nil {
			log.LogErrorf("Error retrieving keys from Redis: %v", err)
			return err
		}

		// Skip deletion if no keys are found, but continue scanning if not finished.
		if len(keys) > 0 {
			totalDeleted, err = s.DeleteKeys(ctx, keys, totalDeleted)
			if err != nil {
				log.LogErrorf("Error deleting keys from Redis: %v", err)
				return err
			}
		}

		scanComplete = (cursor == 0)
	}

	if totalDeleted > 0 {
		log.LogInfof("Deleted %d keys with pattern: %s", totalDeleted, pattern)
	} else {
		// TODO: Define and implement custom error types, such as 'KeyNotFoundError', to provide
		// more granular error information when no keys are found for deletion.
		// This enhancement follows best practices for error handling by allowing more specific error
		// responses and the potential for error handling strategies based on error types.
		log.LogInfof("No keys found with pattern: %s", pattern)
	}

	return nil
}

// ScanKeys returns a slice of keys for a given pattern starting at the cursor.
func (s *service) ScanKeys(ctx context.Context, pattern string, cursor uint64) ([]string, uint64, error) {
	return s.redisClient.Scan(ctx, cursor, pattern, 0).Result()
}

// DeleteKeys deletes a slice of keys from Redis and returns the updated count.
func (s *service) DeleteKeys(ctx context.Context, keys []string, totalDeleted int) (int, error) {
	_, err := s.redisClient.Del(ctx, keys...).Result()
	if err != nil {
		return totalDeleted, err
	}
	return totalDeleted + len(keys), nil
}

// PrepareInsertStatement prepares a SQL insert statement for the transaction.
// The query parameter should be a valid SQL insert statement.
//
// Example Usage:
//
//	ctx := context.Background()
//	tx, err := db.BeginTx(ctx, nil)
//	if err != nil {
//	    log.LogErrorf("Failed to start transaction: %v", err)
//	    return err
//	}
//	defer db.EnsureTransactionClosure(tx, &err)
//
//	query := "INSERT INTO users (name, email) VALUES (?, ?)"
//	stmt, err := db.PrepareInsertStatement(ctx, tx, query)
//	if err != nil {
//	    log.LogErrorf("Failed to prepare insert statement: %v", err)
//	    return err
//	}
//	defer stmt.Close()
//
//	// Use the prepared statement to execute the insert
//	_, err = stmt.ExecContext(ctx, "Gopher", "gopher@go.dev")
//	if err != nil {
//	    log.LogErrorf("Failed to insert user: %v", err)
//	    return err
//	}
func (s *service) PrepareInsertStatement(ctx context.Context, tx *sql.Tx, query string) (*sql.Stmt, error) {
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		log.LogErrorf("Error preparing insert statement: %v", err)
		return nil, err
	}
	return stmt, nil
}
