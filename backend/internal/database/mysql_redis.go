// Copyright (c) 2024 H0llyW00dz & Melkeydev (go-blueprint author) All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// Note: The database package here is not covered by tests and won't have tests implemented for it,
// as it is not worth testing the database that requires authentication. (literally stupid testing that requires authentication unlike mock)

package database

import (
	"context"
	"database/sql"
	"fmt"
	log "h0llyw00dz-template/backend/internal/logger"
	"h0llyw00dz-template/backend/internal/middleware/authentication/crypto/bcrypt"
	"h0llyw00dz-template/env"
	"io"
	"math"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	_ "github.com/go-sql-driver/mysql" // MySQL driver is used for connecting to MySQL databases.
	// This package automatically loads environment variables from a .env file.
	//
	// Note: This may trigger a false positive in any secret scanners LMAO hahaha, However,
	// it is not an actual security issue in this case because this method is lightweight and efficient compared to using cryptographic techniques,
	// which can be expensive in terms of memory usage (potentially adding 100MB+ overhead) just for handling environment variables.
	_ "github.com/joho/godotenv/autoload"
)

// Service defines the interface for database operations that can be performed.
type Service interface {
	// Health checks the health of the database connection.
	Health(filter string) map[string]string

	// Close terminates the database connection.
	Close() error

	// Exec executes a SQL query with the provided arguments and returns the result.
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)

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
	ExecWithoutRow(ctx context.Context, query string, args ...any) error

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
	QueryRow(ctx context.Context, query string, args ...any) *sql.Row

	// FiberStorage returns the [fiber.Storage] interface for storage middleware.
	//
	// Note: This uses Redis. For other storage options, it is recommended to implement a similar approach.
	// More information can be found at https://docs.gofiber.io/storage/.
	// It is possible to use, for example, 5 different databases/storage options in this repository at a large scale.
	//
	// Additionally, due to the current implementation of this boilerplate, which supports compatibility with HPA,
	// using 5 storage options (only fiber storage) might shift it from HPA (Horizontal Pod Autoscaler) to HNA (Horizontal Node Autoscaler).
	// If you're not familiar with HNA, it functions as an auto-pilot.
	//
	// Warning:
	// 	- Switching to a stateful configuration is not recommended (bad).
	FiberStorage() fiber.Storage

	// ScanAndDel uses the Redis SCAN command to iterate over a set of keys and delete them.
	// It's particularly useful for deleting keys with a common pattern.
	//
	// Example Usage:
	//   // With Slice
	//   sliceKey := []string{"gopher_key:*", "another_gopher_key:*"}
	//
	//   if err := db.ScanAndDel(ctx, sliceKey...); err != nil {
	//       log.LogErrorf("Failed to clear keys cache: %v", err)
	//       return err
	//   }
	//
	// Note: This method is capable of handling a large number of keys (millions, billions, or more) efficiently,
	// depending on your infrastructure. It dynamically adjusts the scan count based on performance to ensure
	// the operation remains efficient and avoids overwhelming the Redis server.
	ScanAndDel(ctx context.Context, patterns []string) error

	// PrepareInsertStatement prepares a SQL insert statement for the transaction.
	PrepareInsertStatement(ctx context.Context, tx *sql.Tx, query string) (*sql.Stmt, error)

	// ScanKeys returns a slice of keys for a given pattern starting at the cursor.
	ScanKeys(ctx context.Context, pattern string, cursor uint64, count int64) ([]string, uint64, error)

	// DeleteKeys deletes a slice of keys from Redis and returns the updated count.
	DeleteKeys(ctx context.Context, keys []string, totalDeleted int) (int, error)

	// RestartRedisConnection safely closes the existing connection to Redis and establishes a new one.
	RestartRedisConnection() error

	// RestartMySQLConnection safely restarts the MySQL connection.
	RestartMySQLConnection() error

	// AuthUser returns the ServiceAuth interface for managing user authentication-related database operations.
	//
	// Note: This can be used to connect code to other functions once database.New() is called. The reason for implementing it like this
	// is that it's a better approach for maintaining a large codebase, as it uses the singleton pattern for database management.
	Auth() ServiceAuth

	// SetKeysAtPipeline efficiently sets multiple key-value pairs in Redis/Valkey with a specified TTL (Time To Live) using pipelining.
	SetKeysAtPipeline(ctx context.Context, keyValues map[string]any, ttl time.Duration) error

	// GetKeysAtPipeline efficiently retrieves the values for the given keys from Redis/Valkey using pipelining.
	GetKeysAtPipeline(ctx context.Context, keys []string) (map[string]any, error)

	// BackupTables creates a backup of specified tables in the database.
	//
	// Compatibility:
	// - Ensure the network is stable and the MySQL server is properly configured in real-world scenarios.
	// - Race conditions are minimized by default, as seen in the function implementation.
	// - When performing backups for imports, ensure Unicode characters are handled correctly.
	//   MySQL Workbench might produce errors like "\xF0\x9F\x87\xAE\xF0\x9F..." during import.
	// - Set batchSize based on calculation: if a one table has 100K rows, set it to 10K is sufficient.
	BackupTables(tablesToBackup []string, batchSize int) error

	// BackupTablesConcurrently is an interface method for creating backups of specified tables concurrently 🛰️.
	//
	// Note: This method also safely writes to a single file and has been tested with over 70MiB of data.
	//
	// Compatibility:
	// - Ensure the network is stable and the MySQL server is properly configured in real-world scenarios.
	// - Race conditions are minimized by default, as seen in the function implementation.
	// - When performing backups for imports, ensure Unicode characters are handled correctly.
	//   MySQL Workbench might produce errors like "\xF0\x9F\x87\xAE\xF0\x9F..." during import.
	// - Set batchSize based on calculation: if a one table has 100K rows, set it to 10K is sufficient.
	BackupTablesConcurrently(tablesToBackup []string, o io.Writer, batchSize int) error

	// BackupTablesWithGPG creates a backup of specified tables in the database and encrypts it using a PGP public key.
	// This function generates a .sql file containing the SQL statements needed to recreate
	// the database schema and insert all the current data for each specified table.
	// After generating the backup file, it encrypts the file using the provided PGP public key.
	//
	// Parameters:
	//   - tablesToBackup: A slice of strings containing the names of the tables to back up.
	//   - publicKey: A string representing the PGP public key used to encrypt the backup file.
	//
	// Returns:
	//   - error: An error is returned if any part of the backup or encryption process fails.
	//
	// Process:
	//   1. Validate the table names to ensure they are correct.
	//   2. Create a backup file with a timestamped name to store the SQL dump.
	//   3. Write SQL headers and dump the schema and data for each specified table into the file.
	//   4. Close the file to ensure all data is written and the file handle is released.
	//   5. Encrypt the backup file using the provided PGP public key, creating an encrypted file with a .gpg extension.
	//   6. Log the successful completion of the backup and encryption process.
	//   7. Remove the unencrypted backup file to ensure data security.
	//
	// Notes:
	//   - The function uses a context with a timeout to ensure that the backup process does not run indefinitely.
	//   - It is important to handle file closures and removals carefully to avoid resource leaks and ensure security.
	//   - The encryption process uses the proton library for PGP encryption.
	//
	// Example Usage:
	// if err := db.BackupTablesWithGPG([]string{"users", "orders"}, publicKey); err != nil {
	// 	log.LogErrorf("Backup failed: %v", err)
	// }
	//
	// Note: For decryption, this should work with any GPG frontend (e.g., https://github.com/saturneric/GpgFrontend).
	//
	// For example, decryption (tested on my strong GPG):
	//
	// Decrypt Operation - Success
	//
	// General State:
	//
	// 	- MIME: false
	// 	- Message Integrity Protection: true
	// 	- Symmetric Encryption Algorithm: AES256.CFB
	// 	- German Encryption Standards: false
	//
	// Recipient(s):
	//
	// Recipient [1]: h0llyw00dzz@pm.me<h0llyw00dzz@pm.me>
	// 	- Key ID: AB68BB42A56C9894
	// 	- Public Key Algo: ECDH
	// 	- Status: Success
	//
	// Compatibility:
	// - Ensure the network is stable and the MySQL server is properly configured in real-world scenarios.
	// - Race conditions are minimized by default, as seen in the function implementation.
	// - When performing backups for imports, ensure Unicode characters are handled correctly.
	//   MySQL Workbench might produce errors like "\xF0\x9F\x87\xAE\xF0\x9F..." during import.
	// - Set batchSize based on calculation: if a one table has 100K rows, set it to 10K is sufficient.
	BackupTablesWithGPG(tablesToBackup []string, publicGPGKey []string, batchSize int) error

	// PingDB checks the connectivity of both the MySQL database and the Redis instance.
	//
	// Note: This is effective for health probes (e.g., liveness/readiness) on Kubernetes with HPA.
	// It is also recommended to implement a new listener mechanism for health probes that can only be accessed
	// within the cluster itself and not exposed publicly. In GKE, this can effectively speed up pod replacement.
	PingDB(ctx context.Context) bool

	// GetKeysJSONAtPipeline retrieves multiple objects from Redis using JSON.GET with a custom decoder.
	// It uses pipelining to fetch multiple keys in a single network call, improving efficiency.
	//
	// Note: Ensure your Redis instance has the RedisJSON module enabled.
	// For more efficiency with simple string values, consider using [GetKeysAtPipeline].
	//
	// TODO: Switch to generics [T] or implement this specifically for generics [T] ?
	// This change will elevate it to a top-level function if you know how to handle it.
	GetKeysJSONAtPipeline(ctx context.Context, objects []any, decoder JSONDecoder, keyFunc KeyFunc, path ...string) ([]any, error)

	// SetKeysJSONAtPipeline stores multiple objects in Redis using JSON.SET with a custom encoder and key extractor.
	// It pipelines the commands to reduce network latency, making it efficient for bulk operations.
	//
	// Note: Ensure your Redis instance has the RedisJSON module enabled.
	// For more efficiency with simple string values, consider using [SetKeysAtPipeline].
	//
	// TODO: Switch to generics [T] or implement this specifically for generics [T] ?
	// This change will elevate it to a top-level function if you know how to handle it.
	SetKeysJSONAtPipeline(ctx context.Context, objects []any, encoder JSONEncoder, keyFunc KeyFunc, path ...string) error

	// GetRawJSONAtPipeline retrieves multiple JSON objects from Redis without decoding them.
	// It uses pipelining to fetch multiple keys in a single network call, improving efficiency.
	// The function returns a map of keys to their corresponding raw JSON byte slices.
	//
	// Note: Ensure your Redis instance has the RedisJSON module enabled.
	// For more efficiency with simple string values, consider using [GetKeysAtPipeline].
	//
	// TODO: Switch to generics [T] or implement this specifically for generics [T] ?
	// This change will elevate it to a top-level function if you know how to handle it.
	GetRawJSONAtPipeline(ctx context.Context, objects []any, keyFunc KeyFunc, path ...string) (map[string][]byte, error)

	// DelKeysJSONAtPipeline removes JSON objects from Redis using the JSON.DEL command.
	// It utilizes pipelining to efficiently delete multiple keys in a single network call.
	//
	// TODO: Switch to generics [T] or implement this specifically for generics [T] ?
	// This change will elevate it to a top-level function if you know how to handle it.
	DelKeysJSONAtPipeline(ctx context.Context, objects []any, keyFunc KeyFunc, path ...string) error

	// StreamRows executes a given query and streams the rows, allowing for efficient iteration over large datasets.
	//
	// Note: The caller is responsible for ensuring SQL injection protection when using StreamRows.
	// This method is suitable for syncing MySQL x Redis if implemented correctly (easy 🤪).
	StreamRows(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

// service is a concrete implementation of the Service interface.
type service struct {
	db          *sql.DB
	rdb         fiber.Storage
	redisClient *redis.Client
	mu          sync.Mutex // a mutex to guard connection restarts or any that needed
	auth        ServiceAuth
	bcrypt      *bcrypt.Hash
	initRedis   *RedisClientConfig
	initMysql   *MySQLConfig
}

// dbConfig holds the environment variables for the database connection.
//
// Note: Regarding this Using environment variables in global variables, if you think this high risk you are fucking stupid as developer or security.
var (
	dbname               = os.Getenv(env.DBDATABASE)
	password             = os.Getenv(env.DBPASSWORD)
	username             = os.Getenv(env.DBUSERNAME)
	port                 = os.Getenv(env.DBPORT)
	host                 = os.Getenv(env.DBHOST)
	redisAddress         = os.Getenv(env.RDBADDRESS)
	redisPort            = os.Getenv(env.RDBPORT)
	redisPassword        = os.Getenv(env.RDBPASSWORD)
	redisDatabase        = os.Getenv(env.RDBDATABASE)
	redisPoolTimeout     = os.Getenv(env.RDBPOOLTIMEOUT)
	redisConnMaxIdleTime = os.Getenv(env.RDBMAXCONNLIFEIDLE)
	redisConnMaxLifetime = os.Getenv(env.RDBMAXCONNLIFETIME)
	mysqltlsCAs          = os.Getenv(env.MYSQLCERTTLS)
	redistlsCAs          = os.Getenv(env.REDISCERTTLS)
	dbInstance           *service
	initOnce             sync.Once
)

// New creates a new instance of the Service interface.
// It opens a connection to the MySQL database using the environment variables
// and sets up the connection pool configuration.
//
// Note: For better connection establishment, it's recommended to put this in the "func init()"
// so that it will initialize before the "func main()" runs. This is because the connection will be
// shared across the entire codebase (sharing is caring), even if the codebase grows to billions of lines of Go code (e.g, Senior Golang).
//
// Additionally, this improvement supports scaling up to 100+ pods. For even larger scalability (HPA),
// consider multiple deployments to support up to 1,000 pods, ideally with multi-architecture (AMD64 x ARM64) support.
// The recommended ratio is 10 deployments for every 1,000 pods (e.g., 50 deployments for 5,000 pods, 100 deployments for 10,000 pods).
// Note that this calculation is for HPA (stateless) and is based on deployment ratios for scalability. For stateful setups, predictability based on demand is not feasible hahaha.
//
// It also improves connection stability for MySQL, reducing occasional drops. For Redis, it enhances latency due to the use of a pool of goroutines.
func New() Service {
	initOnce.Do(func() {
		// Initialize the Redis client
		redisClient, err := initializeRedisClient()
		if err != nil {
			// This will catch connection errors such as timeouts and parsing errors from the "strconv" package.
			log.LogFatal("Failed to initialize Redis client:", err)
		}

		// Initialize Redis storage for Fiber
		redisStorage, err := initializeRedisStorage()
		if err != nil {
			// This will catch connection errors such as timeouts and parsing errors from the "strconv" package.
			log.LogFatal("Failed to initialize Redis storage:", err)
		}

		// Initialize the MySQL database
		db, err := initializeMySQLDB()
		if err != nil {
			// This will not be a connection error, but a DSN parse error or
			// another initialization error.
			log.LogFatal("Failed to initialize MySQL database:", err)
		}

		// Initialize bcrypt
		// Note: This operation should be inexpensive as it uses a pointer,
		// and the garbage collector will be happy handling memory efficiently. 🤪
		bchash, err := bcrypt.New()

		if err != nil {
			// This will not be a connection error, but a Cost error or
			// another initialization error.
			log.LogFatal("Failed to initialize bcrypt:", err)
		}

		// Create the service instance
		dbInstance = &service{
			db:          db,
			rdb:         redisStorage, // use redisStorage for rate limiting or any other needs in middleware
			redisClient: redisClient,
			// Note: This method is safe, even with a large number of service instances (e.g., 1 billion) due to the singleton pattern.
			// Also, MySQL should be used as the primary database, while Redis should be used for caching.
			// Here's an example data flow:
			// 1. For read operations: service -> Redis (if not found in Redis) -> get from main database -> put back in Redis -> repeat.
			// 2. For insert/update operations: service -> main database -> repeat.
			// Redis will handle caching for read operations, while write operations will directly interact with the main database.
			// Also note that these example data flows are highly stable, and the reason for this logic is that traditional SQL databases (e.g., MySQL) have limited open connections,
			// unlike NoSQL databases (e.g., Redis), which are capable of up to 10K connections with basically no limits.
			// So Redis is perfect for connection pooling because the most important factor for interacting with it is the connection itself.
			auth: NewServiceAuth(db, redisStorage, bchash),
		}
	})

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
//
// Note: For MySQL, it will always show 1 connection because it's not an actual pool of goroutines;
// it's a pool of pods k8s. However, for Redis, it will show multiple connections due to goroutine usage.
func (s *service) Health(filter string) map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultCtxTimeout)
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
	if err := s.db.PingContext(ctx); err != nil {
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

	// Note: It is important to monitor the WaitCount (known as wait events) when using only a
	// single MySQL database or multiple MySQL databases without other database types (e.g., Redis).
	// If a high WaitCount is detected, it could indicate a potential bottleneck that may lead to a Denial of Service (DoS) scenario, which should not happen.
	if dbStats.WaitCount > 500 { // 500 is enough, since it using redis & mysql (multiple database)
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
			stats["redis_message"] = fmt.Sprintf(MsgRedisFailedToRetrieveInfo, err)
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
			// TODO: Implement a helper function to extract and format numerical values from health stats.
			// This will eliminate the need for repeated `strconv.FormatUint` calls, ensuring consistency
			// across the frontend REST API responses.
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
			// Note: This fixes a potential underflow issue that may occur in certain rare cases.
			// The problem only occurs occasionally.
			// Also, note that math.Max is used now because it might prevent warnings in other Go linters, such as the golangci-lint.
			//
			// TODO: Remove this due to its unpredictability, especially in Redis that has Autopilot and Kafka integration.
			activeConns := uint64(math.Max(float64(poolStats.TotalConns-poolStats.IdleConns), 0))
			stats["redis_active_connections"] = strconv.FormatUint(activeConns, 10)

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
	highConnectionThreshold := int(float64(poolSize) * 0.8)

	// Check if connected clients exceed the high connection threshold
	if connectedClients > highConnectionThreshold {
		stats["redis_message"] = MsgRedisHighConnectedClients
	}

	// Check for stale connections and append a warning if they exceed a certain percentage threshold.
	//
	// Note: This approach is more dynamic instead of being explicitly static, which is suitable for Redis with Autopilot
	// (an automated managed database for high performance (zer0 overhead) that is capable of handling 1 billion hits
	// even with low resources (e.g., 1GB Redis). Also note that it is not a cluster or sentinel setup) and Kafka integration.
	totalHitsConns := int(poolStats.Hits)
	if totalHitsConns > 0 {
		stalePercentage := float64(poolStats.StaleConns) / float64(totalHitsConns) * 100
		if stalePercentage > 70 { // Adjust this threshold as needed
			stats["redis_message"] = fmt.Sprintf(MsgRedisHasStaleConnections, poolStats.StaleConns, stalePercentage)
		}
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

	// Check the number of idle connections
	idleConns := int(poolStats.IdleConns)
	// Determine a high idle connection threshold, let's say 70% of the pool size
	highIdleConnectionThreshold := int(float64(poolSize) * 0.7)
	if idleConns > highIdleConnectionThreshold {
		stats["redis_message"] = MsgRedisHighIdleConnections
	}

	// Check the pool utilization
	poolUtilization := float64(poolStats.TotalConns-poolStats.IdleConns) / float64(poolSize) * 100
	// Determine a high pool utilization threshold, let's say 90%
	highPoolUtilizationThreshold := 90.0
	if poolUtilization > highPoolUtilizationThreshold {
		stats["redis_message"] = MsgRedisHighPoolUtilization
	}

	// this possible reached 100%++ (e.g, 150%) due wrong configuration.
	bottleneckThreshold := 100.0
	if poolUtilization > bottleneckThreshold {
		stats["redis_message"] = MsgRedisHighPoolBottleneck
	}

	return stats
}

// Exec executes a SQL query with the provided arguments.
func (s *service) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}

// ExecWithoutRow executes a query without returning any rows.
//
// Note: This method is different from "Exec". Unlike "Exec", it doesn't return "sql.Result".
// This method is better suited for initializing database schemas or running migrations before the app starts.
func (s *service) ExecWithoutRow(ctx context.Context, query string, args ...any) error {
	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
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
func (s *service) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return s.db.QueryRowContext(ctx, query, args...)
}

// FiberStorage returns the [fiber.Storage] interface for fiber storage middleware.
//
// Note: For structured databases, it is recommended to use a relational database and a stable package
// such as the standard library's [database/sql] package with a MySQL driver. This is because when using
// FiberStorage, you must implement the storage functionality from scratch to make it structurable,
// such as storing JSON data in Redis (See Key-Auth + Session Middleware Logic). Also Note that, this doesn't matter for FiberStorage in other Go code
// (e.g., auth.go) once the connection to the database has been established, as it will be shared across
// the codebase (sharing is caring).
func (s *service) FiberStorage() fiber.Storage {
	return s.rdb
}

// ScanAndDel uses the Redis SCAN command to iterate over a set of keys and delete them.
// It accepts one or more key patterns and deletes keys matching any of the patterns.
func (s *service) ScanAndDel(ctx context.Context, patterns []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Use a context with a timeout to avoid hanging indefinitely
	// Note: This should fix an issue where the function hangs when using RedisClientConfig with "ContextTimeoutEnabled" set to true.
	ctx, cancel := context.WithTimeout(ctx, DefaultCtxTimeout)
	defer cancel()

	var totalDeleted int
	// Note: The count is dynamically adjusted based on performance.
	// This is sufficient, unlike Redis Insight, which can encounter blank screens (on windows hahaha) or OOM errors (on k8s hahaha) with many key-value pairs,
	// such as millions/billions.
	baseCount := 1     // Initial count for retrieving keys
	maxCount := 100000 // Maximum limit for the count; if insufficient, it can be increased to 1 million

	for _, pattern := range patterns {
		var cursor uint64
		count := baseCount
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				// Continue scanning keys as the context is still valid.
				// Do not return here.
			}

			start := time.Now()

			// Use ScanKeys to retrieve keys
			keys, nextCursor, err := s.ScanKeys(ctx, pattern, cursor, int64(count))
			if err != nil {
				return err
			}

			// Calculate time taken for this iteration
			duration := time.Since(start)

			// Adjust count based on duration
			if duration < 50*time.Millisecond && count < maxCount {
				count *= 2 // Double the count if operation is fast
			} else if duration > 200*time.Millisecond && count > baseCount {
				count /= 2 // Halve the count if operation is slow
			}

			// Delete keys if found
			if len(keys) > 0 {
				deleted, err := s.DeleteKeys(ctx, keys, totalDeleted)
				if err != nil {
					return err
				}
				totalDeleted += deleted
			}

			cursor = nextCursor

			if cursor == 0 {
				break
			}
		}
	}

	if totalDeleted > 0 {
		log.LogInfof("Deleted %d keys with patterns: %v", totalDeleted, patterns)
	} else {
		// TODO: Define and implement custom error types, such as 'KeyNotFoundError', to provide
		// more granular error information when no keys are found for deletion.
		// This enhancement follows best practices for error handling by allowing more specific error
		// responses and the potential for error handling strategies based on error types.
		log.LogInfof("No keys found with patterns: %v", patterns)
	}

	return nil
}

// ScanKeys returns a slice of keys for a given pattern starting at the cursor.
func (s *service) ScanKeys(ctx context.Context, pattern string, cursor uint64, count int64) ([]string, uint64, error) {
	return s.redisClient.Scan(ctx, cursor, pattern, count).Result()
}

// DeleteKeys deletes a slice of keys from Redis and returns the updated count.
func (s *service) DeleteKeys(ctx context.Context, keys []string, totalDeleted int) (int, error) {
	if _, err := s.redisClient.Del(ctx, keys...).Result(); err != nil {
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

// RestartRedisConnection safely closes the existing connection to Redis and establishes a new one.
func (s *service) RestartRedisConnection() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Close the existing Redis client connection.
	if err := s.redisClient.Close(); err != nil {
		log.LogErrorf("Error closing Redis client: %v", err)
		return err
	}

	// Reinitialize the Redis client.
	redisClient, err := s.initRedis.InitializeRedisClient()
	if err != nil {
		log.LogErrorf("Error initializing Redis client: %v", err)
		return err
	}
	s.redisClient = redisClient

	// Log the reconnection
	log.LogInfo("Redis connection has been restarted.")

	return nil
}

// RestartMySQLConnection safely closes the existing MySQL connection and establishes a new one.
func (s *service) RestartMySQLConnection() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Close the existing MySQL database connection.
	if err := s.db.Close(); err != nil {
		log.LogErrorf("Error closing MySQL database connection: %v", err)
		return err
	}

	// Reinitialize the MySQL database connection.
	var err error
	s.db, err = s.initMysql.InitializeMySQLDB()
	if err != nil {
		log.LogErrorf("Error reinitializing MySQL database connection: %v", err)
		return err
	}

	// Log the reconnection.
	log.LogInfo("MySQL connection has been restarted.")

	return nil
}

// AuthUser returns the ServiceAuth interface for managing user authentication-related database operations.
func (s *service) Auth() ServiceAuth {
	return s.auth
}

// SetKeysAtPipeline reduces the latency cost associated with round-trip time (RTT) by batching multiple commands (e.g, 1 billion commands that save cost money $$$) into a single network request.
// This method is particularly useful for bulk-insert scenarios where performance is critical.
//
// Example Usage:
//
//	data := map[string]any{
//	    "name":      "Gopher",
//	    "isStudent": true,
//	}
//
// Type assert to their respective types (type-safe manner):
//
//	name := data["name"].(string)
//	isStudent := data["isStudent"].(bool)
//
// Note: On my rack-server machine, it has a 0-ms (backend) and 20-ms (frontend for visitor/client in the SEA region) response time because we are neighbors, so ¯\_(ツ)_/¯
//
// Important:
//   - For better performance, avoid using Redis/Valkey JSON commands.
//     String commands are sufficient as they are immutable and can easily enhance performance by utilizing other JSON encoders/decoders (Most Important Go, unlike other language) for the value key. so ¯\_(ツ)_/¯
func (s *service) SetKeysAtPipeline(ctx context.Context, keyValues map[string]any, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Use a context with a timeout to avoid hanging indefinitely
	//
	// This is compatible/supported with Fiber's context (c.Context()), but it's recommended to use context.Background() if you're familiar with handling contexts.
	// By default, this explicitly uses "context.WithTimeout".
	ctx, cancel := context.WithTimeout(ctx, DefaultCtxTimeout)
	defer cancel()

	// Initialize a new pipeline.
	pipe := s.redisClient.Pipeline()

	// Iterate over the provided key-value pairs, queuing up each one in the pipeline.
	for key, value := range keyValues {
		// Queue a SET command for the key-value pair with the given TTL.
		pipe.Set(ctx, key, value, ttl)
	}

	// Execute the pipelined commands in a single round-trip to the Redis server.
	if _, err := pipe.Exec(ctx); err != nil {
		// Return an enhanced error if the pipelining fails, wrapping the original error for more context.
		return fmt.Errorf("Pipeline execution failed: %w", err)
	}

	// signify successful execution.
	return nil
}

// GetKeysAtPipeline reduces the latency cost associated with round-trip time (RTT) by batching multiple commands (e.g, 1 billion commands that save cost money $$$) into a single network request.
// This method is particularly useful for bulk-insert scenarios where performance is critical.
//
// Note:
//   - If a key does not exist in Redis/Valkey, its corresponding value in the returned map will be nil.
//
// Important:
//   - For better performance, avoid using Redis/Valkey JSON commands.
//     String commands are sufficient as they are immutable and can easily enhance performance by utilizing other JSON encoders/decoders (Most Important Go, unlike other language) for the value key. so ¯\_(ツ)_/¯
func (s *service) GetKeysAtPipeline(ctx context.Context, keys []string) (map[string]any, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Use a context with a timeout to avoid hanging indefinitely
	//
	// This is compatible/supported with Fiber's context (c.Context()), but it's recommended to use context.Background() if you're familiar with handling contexts.
	// By default, this explicitly uses "context.WithTimeout".
	ctx, cancel := context.WithTimeout(ctx, DefaultCtxTimeout)
	defer cancel()

	// Initialize a new pipeline
	pipe := s.redisClient.Pipeline()

	// Create a slice to hold the Redis commands and a map to store the results
	cmds := make([]*redis.StringCmd, len(keys))
	results := make(map[string]any)

	// Iterate over the keys and queue a Get command for each key
	for i, key := range keys {
		cmds[i] = pipe.Get(ctx, key)
	}

	// Execute the pipelined commands in a single round-trip to the Redis/Valkey server
	//
	// Note: For handling Redis/Valkey Nil, in this Pipe execution, it must be handled by the caller. For example:
	//
	// cachedData, err := db.GetKeysAtPipeline([]string{stackCachePrefix})
	// if err != nil {
	// 	if err.Error() == "Pipeline execution failed: redis: nil" {
	// 		Logger.LogInfof("Stack Data: %s not found in Cache", stackData)
	// 		return nil, nil // Not found in cache, return nil data and nil error
	// 	}
	// 	return nil, fmt.Errorf("error retrieving from cache: %w", err)
	// }
	//
	// Additionally, for more suitable + idiom go best practices error handling, implement an [errors.Is] mechanism by using a switch statement (not multiple if-else statements)
	// because this error handling will depend on the server and it's different when it's in production unlike in testing/mock testing. For Example:
	//
	// cachedData, err := db.GetKeysAtPipeline([]string{stackCachePrefix})
	// if err != nil {
	// 	switch {
	// 	case errors.Is(err, "Pipeline execution failed: redis: nil"):
	// 		Logger.LogInfof("Stack Data: %s not found in Cache", stackData)
	// 		return nil, nil // Not found in cache, return nil data and nil error
	// 	default:
	// 		return nil, fmt.Errorf("error retrieving from cache: %w", err)
	// 	}
	// }
	if _, err := pipe.Exec(ctx); err != nil {
		// Return an enhanced error if the pipelining fails, wrapping the original error for more context
		return nil, fmt.Errorf("Pipeline execution failed: %w", err)
	}

	// Iterate over the command results, checking for errors and populating the results map
	for i, cmd := range cmds {
		result, err := cmd.Result()
		if err != nil {
			if err == redis.Nil {
				// If the key does not exist, set its value to nil in the results map
				results[keys[i]] = nil
			} else {
				// Return an enhanced error if there's an error getting a key, including the key in the error message
				return nil, fmt.Errorf("Error getting key '%s': %w", keys[i], err)
			}
		} else {
			// If the key exists and there's no error, store its value in the results map
			results[keys[i]] = result
		}
	}

	// Return the populated results map and a nil error to signify successful execution
	return results, nil
}

// PingDB checks the connectivity of both the MySQL database and the Redis instance.
// It returns true if both services are reachable, otherwise it returns false.
//
// Note: Logging is optional and not strictly necessary, effective replacing pods automated in Kubernetes.
func (s *service) PingDB(ctx context.Context) bool {
	// Use the context directly without setting a new timeout,
	// as the timeout can be set by the caller.

	// Ping the MySQL database to verify connectivity.
	if err := s.db.PingContext(ctx); err != nil {
		return false
	}

	// Ping the Redis server to verify connectivity.
	if err := s.redisClient.Ping(ctx).Err(); err != nil {
		return false
	}

	// Return true if both MySQL and Redis are reachable.
	return true
}

// StreamRows executes a given query and streams the rows, allowing for efficient iteration over large datasets.
func (s *service) StreamRows(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	// Execute the query with the provided arguments
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.LogErrorf("Failed to stream rows: %v", err)
		return nil, fmt.Errorf("failed to stream rows: %w", err)
	}

	return rows, nil
}
