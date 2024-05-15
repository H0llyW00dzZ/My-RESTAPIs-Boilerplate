// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package database

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	log "h0llyw00dz-template/backend/internal/logger"
	"runtime"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
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
	ConnMaxIdleTime       time.Duration
	ConnMaxLifetime       time.Duration
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

// Calculate the maximum number of connections based on the number of CPUs
//
// Note: This redis configuration is better for starters and provides stability.
//
// Tested on Node Spec:
//
//	2x vCPU
//
//	4x ~ 8x Compute
//
//	1 GB RAM
//
// Tested Env Configuration:
//
//	RDB_POOL_TIMEOUT - 5m
//
//	REDIS_MAXCONN_IDLE_TIME - 30m
//
//	REDIS_MAXCONN_LIFE_TIME - 1h
var maxConnections = 2 * runtime.NumCPU()

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
		PoolTimeout:           config.PoolTimeout,           // PoolTimeout should already be a time.Duration
		PoolSize:              config.PoolSize,              // adding back this for default.
		ContextTimeoutEnabled: config.ContextTimeoutEnabled, // adding back this for default.
		MinIdleConns:          config.PoolSize / 4,          // Set minimum idle connections to 25% of the pool size
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
	// Note: Implementing statistics similar to those in Redis isn't feasible due to connection limitations.
	// Even attempting to set it to unlimited will inevitably lead to a bottleneck, regardless of server specs (e.g., even on a high-spec or baremetal server).
	// So, it's best to maintain the current configuration since Redis will handle this aspect.
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

// parseRedisConfig parses the Redis configuration from environment variables and returns a RedisClientConfig struct.
// It handles parsing errors and returns an error if any of the configurations are invalid.
func parseRedisConfig() (RedisClientConfig, error) {
	// Parse the Redis database index from the environment variable.
	redisDB, err := strconv.Atoi(redisDatabase)
	if err != nil {
		return RedisClientConfig{}, fmt.Errorf("invalid Redis database index: %v", err)
	}

	// Parse Redis port from the environment variable
	redisPortInt, err := strconv.Atoi(redisPort)
	if err != nil {
		return RedisClientConfig{}, fmt.Errorf("invalid Redis port: %v", err)
	}

	// Parse pool timeout from the environment variable
	poolTimeout, err := time.ParseDuration(redisPoolTimeout)
	if err != nil {
		return RedisClientConfig{}, fmt.Errorf("invalid Redis pool timeout value: %v", err)
	}

	// Parse connection max life time from the environment variable
	redisConnMaxLifetime, err := time.ParseDuration(redisConnMaxLifetime)
	if err != nil {
		return RedisClientConfig{}, fmt.Errorf("invalid Redis connection max life time value: %v", err)
	}

	// Parse connection max idle time from the environment variable
	redisConnMaxIdleTime, err := time.ParseDuration(redisConnMaxIdleTime)
	if err != nil {
		return RedisClientConfig{}, fmt.Errorf("invalid Redis connection max idle time value: %v", err)
	}

	// Return the RedisClientConfig struct with the parsed configurations
	return RedisClientConfig{
		Address:               redisAddress,
		Port:                  redisPortInt,
		Password:              redisPassword,
		Database:              redisDB,
		PoolTimeout:           poolTimeout,
		PoolSize:              maxConnections,
		ContextTimeoutEnabled: true,
		ConnMaxIdleTime:       redisConnMaxIdleTime,
		ConnMaxLifetime:       redisConnMaxLifetime,
	}, nil
}

// initializeRedisClient initializes the Redis client using the provided Redis configuration.
// It parses the configuration from environment variables and returns a new Redis client instance.
func initializeRedisClient() (*redis.Client, error) {
	// Parse the Redis database index from the environment variable.
	redisDB, err := strconv.Atoi(redisDatabase)
	if err != nil {
		return nil, fmt.Errorf("invalid Redis database index: %v", err)
	}

	// Parse Redis port from the environment variable
	redisPortInt, err := strconv.Atoi(redisPort)
	if err != nil {
		return nil, fmt.Errorf("invalid Redis port: %v", err)
	}

	// Parse pool timeout from the environment variable
	poolTimeout, err := time.ParseDuration(redisPoolTimeout)
	if err != nil {
		return nil, fmt.Errorf("invalid Redis pool timeout value: %v", err)
	}

	// Parse connection max life time from the environment variable
	redisConnMaxLifetime, err := time.ParseDuration(redisConnMaxLifetime)
	if err != nil {
		return nil, fmt.Errorf("invalid Redis connection max life time value: %v", err)
	}

	// Parse connection max idle time from the environment variable
	redisConnMaxIdleTime, err := time.ParseDuration(redisConnMaxIdleTime)
	if err != nil {
		return nil, fmt.Errorf("invalid Redis connection max idle time value: %v", err)
	}

	// Prepare Redis client configuration
	redisClientConfig := RedisClientConfig{
		Address:               redisAddress,
		Port:                  redisPortInt,
		Password:              redisPassword,
		Database:              redisDB,
		PoolTimeout:           poolTimeout,
		ContextTimeoutEnabled: true,
		PoolSize:              maxConnections,
		ConnMaxLifetime:       redisConnMaxLifetime,
		ConnMaxIdleTime:       redisConnMaxIdleTime,
	}

	// Initialize and return the Redis client using the provided configuration
	return InitializeRedisClient(redisClientConfig), nil
}

// initializeRedisStorage initializes the Redis storage for Fiber using the provided Redis configuration.
// It parses the configuration from environment variables and returns a new Redis storage instance.
func initializeRedisStorage() (fiber.Storage, error) {
	// Parse the Redis database index from the environment variable.
	redisDB, err := strconv.Atoi(redisDatabase)
	if err != nil {
		return nil, fmt.Errorf("invalid Redis database index: %v", err)
	}

	// Parse Redis port from the environment variable
	redisPortInt, err := strconv.Atoi(redisPort)
	if err != nil {
		return nil, fmt.Errorf("invalid Redis port: %v", err)
	}

	// Prepare Fiber Redis storage configuration
	fiberRedisConfig := FiberRedisClientConfig{
		Address:  redisAddress,
		Port:     redisPortInt,
		Password: redisPassword,
		Database: redisDB,
		PoolSize: maxConnections,
		Reset:    false,
	}

	// Initialize and return the Redis storage using the provided configuration
	return InitializeRedisStorage(fiberRedisConfig), nil
}

// initializeMySQLDB initializes the MySQL database using the provided MySQL configuration.
// It prepares the configuration from environment variables and returns a new database connection.
func initializeMySQLDB() (*sql.DB, error) {
	// Prepare MySQL configuration
	mysqlConfig := MySQLConfig{
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
		Database: dbname,
	}

	// Initialize and return the MySQL database connection using the provided configuration
	return InitializeMySQLDB(mysqlConfig)
}

// model represents the Bubble Tea model for the spinner.
type model struct {
	spinner  spinner.Model
	quitting bool
}

// Init initializes the model.
func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update updates the model based on the received message.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.QuitMsg:
		return m, tea.Quit
	}
	return m, nil
}

// View renders the spinner.
func (m model) View() string {
	return fmt.Sprintf("\n   %s Initializing database...\n\n", m.spinner.View())
}
