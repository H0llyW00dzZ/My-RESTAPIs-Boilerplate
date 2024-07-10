// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package env

// Application Configuration
const (
	APPNAME     = "APP_NAME"     // The name of the application (default: "Gopher").
	PORT        = "PORT"         // The port number on which the server will listen (default: "8080").
	MONITORPATH = "MONITOR_PATH" // The path for the server monitoring endpoint (default: "/monitor").
	TIMEFORMAT  = "TIME_FORMAT"  // The format for logging timestamps (default: "unix").
	//     Available options:
	//   - "unix": Unix timestamp format (e.g., [1713355079]).
	//   - "default": Default timestamp format (e.g., 2024/04/17 15:04:05).
)

// App/Server Timeout Configuration
const (
	READTIMEOUT     = "READ_TIMEOUT"     // The maximum duration for reading the entire request, including the body (default: "5s").
	WRITETIMEOUT    = "WRITE_TIMEOUT"    // The maximum duration before timing out writes of the response (default: "5s").
	SHUTDOWNTIMEOUT = "SHUTDOWN_TIMEOUT" // The maximum duration to wait for active connections to finish during server shutdown (default: "5s").
)

// MySQL Database Configuration
const (
	DBHOST     = "DB_HOST"     // The hostname or IP address of the MySQL database server (required).
	DBPORT     = "DB_PORT"     // The port number on which the MySQL database server is listening (required).
	DBDATABASE = "DB_DATABASE" // The name of the MySQL database to connect to (required).
	DBUSERNAME = "DB_USERNAME" // The username for authenticating with the MySQL database (required).
	DBPASSWORD = "DB_PASSWORD" // The password for authenticating with the MySQL database (required).
)

// Redis Database Configuration
const (
	RDBADDRESS         = "RDB_ADDRESS"             // The address of the Redis server (required).
	RDBPORT            = "RDB_PORT"                // The port number on which the Redis server is listening (required).
	RDBPASSWORD        = "RDB_PASSWORD"            // The password for authenticating with the Redis server (required).
	RDBDATABASE        = "RDB_DATABASE"            // The Redis database number to use (required).
	RDBPOOLTIMEOUT     = "RDB_POOL_TIMEOUT"        // The maximum amount of time to wait for a connection from the Redis connection pool (required).
	RDBMAXCONNLIFEIDLE = "REDIS_MAXCONN_IDLE_TIME" // The maximum amount of time a Redis connection can remain idle in the connection pool (required).
	RDBMAXCONNLIFETIME = "REDIS_MAXCONN_LIFE_TIME" // The maximum lifetime of a Redis connection in the connection pool (required).
)

// TLS Configuration
const (
	// This environment variable is used to specify additional root CA / subs CA certificates that should be trusted by the application.
	EXTRACERTTLS = "EXTRA_CERTS_TLS" // Base64-encoded root CA / subs CA certificates for establishing secure connections database (required).
)

// Site Middleware Configuration (Optional since it boilerplate and must rewrite a "DomainRouter" in RegisterRoutes (see backend/internal/middleware/routes.go))
const (
	DOMAIN       = "DOMAIN"
	APISUBDOMAIN = "API_SUB_DOMAIN"
)
