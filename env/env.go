// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package env

// Application Configuration
const (
	APPNAME     = "APP_NAME"     // The name of the application (default: "Gopher").
	PORT        = "PORT"         // The port number on which the server will listen (default: "8080").
	MONITORPATH = "MONITOR_PATH" // The path for the server monitoring endpoint (default: "/monitor").
	// TRUSTEDPROXIES is used to specify a list of IP addresses or CIDR ranges considered as trusted proxies.
	// Note: For environments using ingress controllers like NGINX or services with a load balancer,
	// the trusted proxies should include the IP addresses of these load balancers (not a pods ip).
	// It's recommended to also include any other relevant load balancers (e.g., MySQL load balancers).
	// If there are multiple IPs, specify each IP with a /32 suffix for precision.
	// The configuration may vary based on the cloud provider. If the provider supports CIDR pooling for Kubernetes,
	// it can be beneficial. CIDR pooling (e.g., CIDR pooling for node) allows securing the network by restricting external access and
	// optimizing network performance (e.g., consistent latency within regions).
	//
	// Additionally, if using pod IPs, ensure the CIDR ranges are accurate and predictable for the pool size, such as for 100 pods.
	// Accurate and predictable CIDR ranges can optimize network performance as well (e.g., consistently low latency, like in my cluster, the Black Pearl ☸️ ⛵).
	// Default: "0.0.0.0/0", allowing all IPs.
	TRUSTEDPROXIES = "TRUSTED_PROXIES" // A list of IP addresses or CIDR ranges that are considered trusted proxies (default: "0.0.0.0/0").
	TIMEFORMAT     = "TIME_FORMAT"     // The format for logging timestamps (default: "unix").
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

// App/Server Capacity Configuration
const (
	SIZEBODYLIMIT = "SIZE_BODY_LIMIT" // The maximum size limit for the request body (default: "4MiB" (equivalent to 4 * 1024 * 1024 bytes)).
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
	//
	// Best Practice: Use Private CAs, and ensure that the leaf CA adds the IP into the SANs (Subject Alternative Names).
	// This can be easily achieved by implementing a tool like PKIX for CSR generation in Go, instead of using Public CAs (e.g., Trusted CAs).
	// Then it can be easily bound to the network infrastructure (e.g., load balancers, etc.).
	MYSQLCERTTLS = "MYSQL_CERTS_TLS" // Base64-encoded root CA / subs CA certificates for establishing secure connections MySQL database (required).
	REDISCERTTLS = "REDIS_CERTS_TLS" // Base64-encoded root CA / subs CA certificates for establishing secure connections Redis database (required).
	// Note: This Path & File Name TLS secrets are supported and securely managed by Kubernetes as long as the certificate issued implementation is correct.
	// If the implementation is incorrect, this ship ⛵ BlackPearl, will be shrinking.
	SERVERCERTTLS = "TLS_CERT_FILE"
	SERVERKEYTLS  = "TLS_KEY_FILE"
	// DISABLEDEFAULTPORTHTTP determines whether the application should disable listening on the default HTTP port (port 80).
	// Set this environment variable to "true" to prevent the server from accepting insecure HTTP connections.
	// This enhances security by ensuring that only secure connections (HTTPS/TLS) are accepted.
	DISABLEDEFAULTPORTHTTP = "DISABLE_PORT_HTTPINSECURE"
	// Note: This for mTLS (Optional)
	SERVERCATLS = "TLS_CA_FILE"
	ENABLEMTLS  = "ENABLE_MTLS"
)

// Site Middleware Configuration (Optional since it boilerplate and must rewrite a "DomainRouter" in RegisterRoutes (see backend/internal/middleware/routes.go))
const (
	DOMAIN       = "DOMAIN"
	APISUBDOMAIN = "API_SUB_DOMAIN"
)

// Cloudflare-KV Storage Configuration (Optional since it alternative Redis that with better network plus cheap which is suitable for load balancer)
const (
	CFKVKEY         = "CF_KV_KEY"          // The Cloudflare Auth Token.
	CFKVEMAIL       = "CF_KV_EMAIL"        // The Cloudflare Email.
	CFKVACCID       = "CF_KV_ACC_ID"       // The Cloudflare Account ID.
	CFKVNAMESPACEID = "CF_KV_NAMESPACE_ID" // The Cloudflare NameSpace ID.
)
