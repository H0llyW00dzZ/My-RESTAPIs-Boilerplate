// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// Note: The database package here is not covered by tests and won't have tests implemented for it,
// as it is not worth testing the database that requires authentication. (literally stupid testing that requires authentication unlike mock)

package database

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-sql-driver/mysql"
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

// Note: This is balanced for VPA or HPA. However, for HPA, the maximum CPU request should be "555m" without a limit in the deployment.
// HPA doesn't require a large CPU request, and it should be used without a limit. For memory in HPA, it depends on the usage.
// For example, if Fiber storage are used with Redis, it can be implemented to optimize data storage, such as using Fiber cache or rate limiter.
// This can achieve low memory usage because the storage in Fiber cache or rate limiter doesn't allocate memory again when data is already stored in Redis.
// The data is then sent to the client (without allocating memory again), and it can be connected with MySQL when MySQL is streaming.
// In this case, MySQL data can be streamed to Fiber storage, which can then put the data into Redis as a stream as well.
//
// Additionally, for VPA, the Fiber framework is fully stable even for cases where low memory usage is not achieved by combining MySQL streaming and Redis for optimized data storage without allocating memory again.
// The average memory usage for VPA in cases where low memory usage is not achieved by combining MySQL streaming and Redis is not static (stuck at around 50MiB++).
// It will go up and down depending on how the garbage collector is recycling. If immutable is set to false, it can go down to 10MiB or 5MiB when idle.
var defaultFiberMaxConnections = 5 * runtime.GOMAXPROCS(0)

// InitializeRedisClient initializes and returns a new Redis client.
func (config *RedisClientConfig) InitializeRedisClient() (*redis.Client, error) {
	rootCAs, err := loadRedisRootCA()
	if err != nil {
		return nil, err
	}

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
			ClientCAs:  rootCAs,
			MaxVersion: tls.VersionTLS13,
			MinVersion: tls.VersionTLS13,
			// Note: Explicitly setting CurvePreferences is disabled by default to ensure future compatibility with X25519MLKEM768 or SecP256r1MLKEM768.
		},
		PoolTimeout:           config.PoolTimeout,           // PoolTimeout should already be a time.Duration
		PoolSize:              config.PoolSize,              // adding back this for default.
		ContextTimeoutEnabled: config.ContextTimeoutEnabled, // adding back this for default.
		ConnMaxIdleTime:       config.ConnMaxIdleTime,       // Required ENV = REDIS_MAXCONN_IDLE_TIME
		ConnMaxLifetime:       config.ConnMaxLifetime,       // Required ENV = REDIS_MAXCONN_LIFE_TIME
		MinIdleConns:          config.PoolSize / 4,          // Set minimum idle connections to 25% of the pool size
	})
	return client, nil
}

// InitializeMySQLDB initializes and returns a new MySQL database client.
//
// Example Configuration:
//
//	spec:
//
//	containers:
//	- args:
//	  - --ssl-cert=/etc/mysql/tls/db-chain.pem # (leaf,subsca,root) Issued by Subsca
//	  - --ssl-key=/etc/mysql/tls/db.key.pem # Private Key Issued by Subsca
//	  - --ssl-capath=/etc/mysql/tls/root.pem # (rootCA)
//	  - --ssl-capath=/etc/ssl/certs
//	  env: # Secrets from environment variables, other secrets are bound into Hardware Security Modules It's Encrypted.
//	  - name: MYSQL_ROOT_PASSWORD
//	    valueFrom:
//	      configMapKeyRef:
//	        key: MYSQL_ROOT_PASSWORD
//	        name: mysql-config-krmr
//	  image: mysql:latest
//	  imagePullPolicy: Always
//	  name: mysql-1
//	  resources:
//	    limits:
//	      cpu: 500m
//	      ephemeral-storage: 1Gi
//	      memory: 2Gi
//	    requests:
//	      cpu: 500m
//	      ephemeral-storage: 1Gi
//	      memory: 2Gi
//	  securityContext:
//	    capabilities:
//	      drop:
//	      - NET_RAW
//	  terminationMessagePath: /dev/termination-log
//	  terminationMessagePolicy: File
//	  volumeMounts:
//	  - mountPath: /etc/mysql/tls
//	    name: mysql-tls
//	  - mountPath: /etc/mysql/tls/ca-certs
//	    name: mysql-ca-certs
//	  - mountPath: /etc/mysql/conf.d/my.cnf
//	    name: mysql-config
//	    subPath: my.cnf
//
// Note: The Example Configuration It required Run As root because of Image.
//
// Best Practice: Remove Default CAs in the Image (Include Public (Trusted) CAs), then put 1 Private Root CAs.
//
// TODO: Consider improving this by using a pool of goroutines. However, it's not necessary right now
// because having too many connections for MySQL can lead to bottlenecks (MySQL bottlenecks). For now, the current setup
// is sufficient, as Redis will handle most of the connection pooling.
func (config *MySQLConfig) InitializeMySQLDB() (*sql.DB, error) {
	rootCAs, err := loadMySQLRootCA()
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf(MySQLConnect, config.Username, config.Password, config.Host, config.Port, config.Database)

	// Set the TLS configuration for the MySQL connection
	//
	// Note: When connecting via mTLS, set the parameter to "?tls=required".
	// Also, note that for private CAs, when trying to connect via IP, the leaf CA must add the IP into the SANs (Subject Alternative Names).
	//
	// Best Practice: Never set the parameter to "?tls=skip-verify" or disable certificate verification, as it compromises security.
	// Always ensure proper certificate verification is in place to maintain a secure connection.
	if err = mysql.RegisterTLSConfig("custom",
		&tls.Config{
			RootCAs:    rootCAs,
			MaxVersion: tls.VersionTLS13,
			MinVersion: tls.VersionTLS13,
		}); err != nil {
		return nil, err
	}

	// Set the TLS connection parameter in the DSN
	//
	// Note: When connecting via mTLS, set the parameter to "?tls=required".
	// Also, note that for private CAs, when trying to connect via IP, the leaf CA must add the IP into the SANs (Subject Alternative Names).
	//
	// Best Practice: Never set the parameter to "?tls=skip-verify" or disable certificate verification, as it compromises security.
	// Always ensure proper certificate verification is in place to maintain a secure connection.
	dsnWithTLS := fmt.Sprintf("%s?tls=custom", dsn)

	// Open a new connection with the updated DSN
	db, err := sql.Open(dbMYSQL, dsnWithTLS)
	if err != nil {
		return nil, err
	}
	// Set MySQL connection pool parameters.
	// Note: Implementing statistics similar to those in Redis isn't feasible due to connection limitations.
	// Even attempting to set it to unlimited will inevitably lead to a bottleneck, regardless of server specs (e.g., even on a high-spec or baremetal server).
	// So, it's best to maintain the current configuration since Redis will handle this aspect.
	db.SetConnMaxIdleTime(time.Minute * 3) // Connections are not closed due to being idle too long.
	// Note: This is highly scalable when running on Kubernetes, especially with Fiber, which is the best choice with HPA (Horizontal Pod Autoscaling)
	// due to its built-in zer0-allocation and can be dynamic resource usage (e.g., CPU, Memory).
	// The values for "SetMaxIdleConns" and "SetMaxOpenConns" depend on the number of Pods.
	// For example, if the maximum number of replicas in HPA is set to 99999, then "SetMaxIdleConns" and "SetMaxOpenConns" should also be set to 99999.
	// Don't forget to set the maximum connections in the MySQL container to 99999 as well.
	db.SetMaxIdleConns(50) // Maximum number of connections in the idle connection pool.
	db.SetMaxOpenConns(50) // Maximum number of open connections to the database.

	// Set a timeout for the Ping operation.
	//
	// TODO: Use an environment variable to customize the timeout (e.g., "10*time.Second").
	// For now, explicitly setting it to 10 seconds should be sufficient, as it is only used during initialization to avoid runtime confusion.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}

// InitializeRedisStorage initializes and returns a new Redis storage instance
// for use with Fiber middlewares such as rate limiting.
func (config *FiberRedisClientConfig) InitializeRedisStorage() (fiber.Storage, error) {
	rootCAs, err := loadRedisRootCA()
	if err != nil {
		return nil, err
	}

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
			ClientCAs:  rootCAs,
			MaxVersion: tls.VersionTLS13,
			MinVersion: tls.VersionTLS13,
			// Note: Explicitly setting CurvePreferences is disabled by default to ensure future compatibility with X25519MLKEM768 or SecP256r1MLKEM768.
		},
		PoolSize: config.PoolSize, // Adjust the pool size as necessary.
	})
	return storage, nil

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
		PoolSize:              defaultFiberMaxConnections,
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
	return redisConfig.InitializeRedisClient()
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
		PoolSize: defaultFiberMaxConnections,
		// TODO: When ENV (e.g, GO_APP=local) it will set to true.
		Reset: false,
	}

	// Initialize and return the Redis storage using the provided configuration
	return fiberRedisConfig.InitializeRedisStorage()
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
//
// Deprecated: This is no longer used as it is incompatible with the old method for better rendering.
type model struct {
	dotSpinner    spinner.Model
	meterSpinner  spinner.Model
	pointsSpinner spinner.Model
	progress      float64
	quitting      bool
	done          bool
}

// Init initializes the model.
//
// Deprecated: This is no longer used as it is incompatible with the old method for better rendering.
func (m model) Init() tea.Cmd {
	return tea.Batch(m.dotSpinner.Tick, m.meterSpinner.Tick, m.pointsSpinner.Tick)
}

// Update updates the model based on the received message.
//
// Deprecated: This is no longer used as it is incompatible with the old method for better rendering.
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
//
// Deprecated: This is no longer used as it is incompatible with the old method for better rendering.
func (m model) View() string {
	// Apply the color style to the spinner frames
	styledDotSpinner := m.dotSpinner.Style.Render(m.dotSpinner.View())
	styledMeterSpinner := m.meterSpinner.Style.Render(m.meterSpinner.View())
	styledPointsSpinner := m.pointsSpinner.Style.Render(m.pointsSpinner.View())

	// Note: This looks better now.
	// TODO: Handle initialization failure scenarios, such as connection timeouts, since this initialization is only connecting to the database.
	if m.done {
		return fmt.Sprintf("\r   ✓ Database initialization completed   \n")
	}
	return fmt.Sprintf("\r\n   %s Initializing database%s   %s Progress%s", styledDotSpinner, styledPointsSpinner, styledMeterSpinner, styledPointsSpinner)
}
