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
}

// service is a concrete implementation of the Service interface.
type service struct {
	db          *sql.DB
	rdb         fiber.Storage
	redisClient *redis.Client
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

	// Create a new Redis client for health checks
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisAddress, redisPort),
		Password: redisPassword,
		DB:       redisDB,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	})

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
		rdb:         redisStorage, // Use the Redis storage for rate limiting or any that needed in middleware
		redisClient: redisClient,
	}

	return dbInstance
}

// Close closes the database connection.
func (s *service) Close() error {
	log.LogInfof(MsgDBDisconnected, dbname)
	return s.db.Close()
}

// Health checks the health of the database and Redis connections.
// It returns a map with keys indicating various health statistics.
func (s *service) Health(filter string) map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// TODO: Improve by using a "Map of Functions" to reduce complexity (caused by if-else statements) when handling multiple databases in the future.
	if filter == "" || filter == "mysql" {
		stats = s.checkMySQLHealth(ctx, stats)
	}

	if filter == "" || filter == "redis" {
		stats = s.checkRedisHealth(stats)
	}

	return stats
}

// checkMySQLHealth checks the health of the MySQL database and adds the relevant statistics to the stats map.
func (s *service) checkMySQLHealth(ctx context.Context, stats map[string]string) map[string]string {
	// Ping the MySQL database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["mysql_status"] = "down"
		stats["mysql_error"] = fmt.Sprintf(ErrDBDown, err)
		stats["mysql_message"] = fmt.Sprintf("%v", err)
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
func (s *service) checkRedisHealth(stats map[string]string) map[string]string {
	// Ping the Redis server
	// Note: The Redis client must use the method "context.Background()" because
	// the context with timeout may cause unexpected behavior during health checks.
	pong, err := s.redisClient.Ping(context.Background()).Result()
	if err != nil {
		stats["redis_status"] = "down"
		stats["redis_error"] = fmt.Sprintf("Redis is down: %v", err)
		stats["redis_message"] = fmt.Sprintf("%v", err)
	} else {
		// Redis is up
		stats["redis_status"] = "up"
		stats["redis_message"] = "Redis is healthy"
		stats["redis_ping_response"] = pong

		// Get Redis server information
		info, err := s.redisClient.Info(context.Background()).Result()
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

			// Evaluate Redis stats to provide a health message
			stats = s.evaluateRedisStats(redisInfo, stats)
		}
	}

	return stats
}

// evaluateRedisStats evaluates the Redis server statistics and updates the stats map with the appropriate health message.
func (s *service) evaluateRedisStats(redisInfo, stats map[string]string) map[string]string {
	connectedClients, _ := strconv.Atoi(redisInfo["connected_clients"])
	if connectedClients > 40 { // Assuming 50 is the max for this example
		stats["redis_message"] = MsgRedisHighConnectedClients
	}

	usedMemory, _ := strconv.ParseInt(redisInfo["used_memory"], 10, 64)
	if usedMemory > 1024*1024*1024 { // 1 GB
		stats["redis_message"] = MsgRedisHighMemoryUsage
	}

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
	var cursor uint64
	var n int
	for {
		// Use the SCAN command to iterate over the keys.
		var keys []string
		var err error
		keys, cursor, err = s.redisClient.Scan(ctx, cursor, pattern, 0).Result()
		if err != nil {
			log.LogErrorf("Error retrieving keys from Redis: %v", err)
			return err
		}
		n += len(keys)

		// Use the DEL command to delete the keys.
		if len(keys) > 0 {
			_, err = s.redisClient.Del(ctx, keys...).Result()
			if err != nil {
				log.LogErrorf("Error deleting keys from Redis: %v", err)
				return err
			}
		}

		// If the cursor returned by SCAN is 0, we have iterated through all keys.
		if cursor == 0 {
			break
		}
	}
	log.LogInfof("Deleted %d keys with pattern: %s", n, pattern)
	return nil
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
