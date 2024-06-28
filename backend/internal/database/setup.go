// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

// Note: The database package here is not covered by tests and won't have tests implemented for it,
// as it is not worth testing the database that requires authentication. (literally stupid testing that requires authentication unlike mock)

package database

import (
	"crypto/tls"
	"database/sql"
	"fmt"
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
// Note: This Redis configuration is better for starters and provides stability.
//
// Tested on Node Spec:
//   - 2x vCPU
//   - 4x ~ 8x Compute
//   - 1 GB RAM
//
// Redis RAM Spec:
//   - 2 GB total
//   - 1 GB for master known as "primary node" or "master node"
//   - 1 GB for slave known as "replica node" or "slave node" (automated synchronization replica)
//
// Tested Env Configuration:
//   - RDB_POOL_TIMEOUT: 5m
//   - REDIS_MAXCONN_IDLE_TIME: 30m
//   - REDIS_MAXCONN_LIFE_TIME: 1h
var maxConnections = 2 * runtime.NumCPU()

// InitializeRedisClient initializes and returns a new Redis client.
func (config *RedisClientConfig) InitializeRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Address, config.Port),
		Password: config.Password,
		DB:       config.Database,
		// Note: TLSConfig is optional, but it is recommended for better security, so it's advisable to use it.
		// Also note that for non-Kubernetes environments, it is recommended to use TLS. For certificates, packages from https://pkg.go.dev/golang.org/x/crypto@v0.24.0/acme or Caddy can be used.
		// Personally, I don't use this because I am running on Kubernetes with another secure connection method (e.g., bound pods/node ports).
		// For Mutual TLS or whatever it is, see: https://redis.io/docs/latest/operate/rc/security/database-security/tls-ssl/. However,
		// the requirement for Mutual TLS or whatever it is depends on how the cloud provider sets it up.
		// For example, in some cloud providers, Mutual TLS or whatever it is may not be needed, and only the following settings are required.
		TLSConfig: &tls.Config{
			// Explicitly set the maximum and minimum TLS versions to 1.3 this server anyways.
			// However Go's standard TLS 1.3 implementation is broken because it keeps forcing the use of the AES-GCM cipher suite.
			MaxVersion: tls.VersionTLS13,
			MinVersion: tls.VersionTLS13,
			CurvePreferences: []tls.CurveID{
				tls.X25519,
				tls.CurveP256,
				tls.CurveP384,
				tls.CurveP521,
			},
		},
		PoolTimeout:           config.PoolTimeout,           // PoolTimeout should already be a time.Duration
		PoolSize:              config.PoolSize,              // adding back this for default.
		ContextTimeoutEnabled: config.ContextTimeoutEnabled, // adding back this for default.
		MinIdleConns:          config.PoolSize / 4,          // Set minimum idle connections to 25% of the pool size
	})
	return client
}

// InitializeMySQLDB initializes and returns a new MySQL database client.
func (config *MySQLConfig) InitializeMySQLDB() (*sql.DB, error) {
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
	return db, nil
}

// InitializeRedisStorage initializes and returns a new Redis storage instance
// for use with Fiber middlewares such as rate limiting.
func (config *FiberRedisClientConfig) InitializeRedisStorage() fiber.Storage {
	storage := redisStorage.New(redisStorage.Config{
		Host:     config.Address,
		Port:     config.Port,
		Password: config.Password,
		Database: config.Database,
		Reset:    config.Reset,
		// Note: TLSConfig is optional, but it is recommended for better security, so it's advisable to use it.
		// Also note that for non-Kubernetes environments, it is recommended to use TLS. For certificates, packages from https://pkg.go.dev/golang.org/x/crypto@v0.24.0/acme or Caddy can be used.
		// Personally, I don't use this because I am running on Kubernetes with another secure connection method (e.g., bound pods/node ports).
		// For Mutual TLS or whatever it is, see: https://redis.io/docs/latest/operate/rc/security/database-security/tls-ssl/. However,
		// the requirement for Mutual TLS or whatever it is depends on how the cloud provider sets it up.
		// For example, in some cloud providers, Mutual TLS or whatever it is may not be needed, and only the following settings are required.
		TLSConfig: &tls.Config{
			// Explicitly set the maximum and minimum TLS versions to 1.3 this server anyways.
			// However Go's standard TLS 1.3 implementation is broken because it keeps forcing the use of the AES-GCM cipher suite.
			MaxVersion: tls.VersionTLS13,
			MinVersion: tls.VersionTLS13,
			CurvePreferences: []tls.CurveID{
				tls.X25519,
				tls.CurveP256,
				tls.CurveP384,
				tls.CurveP521,
			},
		},
		PoolSize: config.PoolSize, // Adjust the pool size as necessary.
	})
	return storage
}

// parseRedisConfig parses the Redis configuration from environment variables and returns a RedisClientConfig struct.
// It handles parsing errors and returns an error if any of the configurations are invalid.
func parseRedisConfig() (*RedisClientConfig, error) {
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

	// Return the RedisClientConfig struct with the parsed configurations
	return &RedisClientConfig{
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
	// Parse the Redis configuration from environment variables
	redisConfig, err := parseRedisConfig()
	if err != nil {
		return nil, err
	}

	// Initialize and return the Redis client using the provided configuration
	return redisConfig.InitializeRedisClient(), nil
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
	fiberRedisConfig := &FiberRedisClientConfig{
		Address:  redisAddress,
		Port:     redisPortInt,
		Password: redisPassword,
		Database: redisDB,
		PoolSize: maxConnections,
		Reset:    false,
	}

	// Initialize and return the Redis storage using the provided configuration
	return fiberRedisConfig.InitializeRedisStorage(), nil
}

// initializeMySQLDB initializes the MySQL database using the provided MySQL configuration.
// It prepares the configuration from environment variables and returns a new database connection.
func initializeMySQLDB() (*sql.DB, error) {
	// Prepare MySQL configuration
	mysqlConfig := &MySQLConfig{
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
		Database: dbname,
	}

	// Initialize and return the MySQL database connection using the provided configuration
	return mysqlConfig.InitializeMySQLDB()
}

// model represents the Bubble Tea model for the spinners.
type model struct {
	dotSpinner    spinner.Model
	meterSpinner  spinner.Model
	pointsSpinner spinner.Model
	progress      float64
	quitting      bool
	done          bool
}

// Init initializes the model.
func (m model) Init() tea.Cmd {
	return tea.Batch(m.dotSpinner.Tick, m.meterSpinner.Tick, m.pointsSpinner.Tick)
}

// Update updates the model based on the received message.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmds []tea.Cmd
		dotSpinner, cmd := m.dotSpinner.Update(msg)
		cmds = append(cmds, cmd)
		meterSpinner, cmd := m.meterSpinner.Update(msg)
		cmds = append(cmds, cmd)
		pointsSpinner, cmd := m.pointsSpinner.Update(msg)
		cmds = append(cmds, cmd)

		// Update the progress value
		m.progress += 0.1
		if m.progress > 1.0 {
			m.progress = 0.0
		}

		return model{
			dotSpinner:    dotSpinner,
			meterSpinner:  meterSpinner,
			pointsSpinner: pointsSpinner,
			progress:      m.progress,
			quitting:      m.quitting,
		}, tea.Batch(cmds...)
	case tea.QuitMsg:
		return m, tea.Quit
	}
	return m, nil
}

// View renders the spinners.
func (m model) View() string {
	// Apply the color style to the spinner frames
	styledDotSpinner := m.dotSpinner.Style.Render(m.dotSpinner.View())
	styledMeterSpinner := m.meterSpinner.Style.Render(m.meterSpinner.View())
	styledPointsSpinner := m.pointsSpinner.Style.Render(m.pointsSpinner.View())

	// Note: This looks better now.
	// TODO: Handle initialization failure scenarios, such as connection timeouts, since this initialization is only connecting to the database.
	if m.done {
		return fmt.Sprintf("\r   ✓ Database initialization completed   \n\n")
	}
	return fmt.Sprintf("\r\n   %s Initializing database%s   %s Progress%s", styledDotSpinner, styledPointsSpinner, styledMeterSpinner, styledPointsSpinner)
}
