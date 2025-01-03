// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package server

import (
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"h0llyw00dz-template/backend/internal/database"
	log "h0llyw00dz-template/backend/internal/logger"
	"h0llyw00dz-template/backend/internal/middleware"
	"h0llyw00dz-template/env"

	"github.com/gofiber/fiber/v2"
)

var (
	// apiSubdomain is the subdomain for the API endpoints.
	// It is set using the API_SUB_DOMAIN environment variable.
	// Example: set API_SUB_DOMAIN=api.localhost:8080 for local development.
	apiSubdomain = os.Getenv(env.APISUBDOMAIN)

	// disableHTTPInsecure determines if the server should avoid listening on port 80 for HTTP.
	// This is controlled by the DISABLE_PORT_HTTPINSECURE environment variable.
	// If set to "true", the server will not set up an HTTP listener on port 80 otherwise it will listening (e.g., keep empty the env).
	disableHTTPInsecure = env.GetEnv(env.DISABLEDEFAULTPORTHTTP, "") == "true"
)

// Server defines the interface for a server that can be started, shut down, and clean up its database.
type Server interface {
	Start(addr, monitorPath string, tlsConfig *tls.Config, streamListener net.Listener)
	Shutdown(ctx context.Context, shutdownTimeout time.Duration) error
	CleanupDB() error
	Mount(prefix string, app any)
	MountPath(path string, handler any)
	SubmitToCTLog(cert *x509.Certificate, privateKey crypto.PrivateKey, ctLog CTLog, httpRequestMaker *HTTPRequestMaker) error
}

// FiberServer implements the Server interface for a Fiber application.
type FiberServer struct {
	App        *fiber.App
	db         database.Service
	httpServer *http.Server
}

// NewFiberServer returns a new FiberServer with the given Fiber app, application name, and monitor path.
// It also initializes the database and registers routes.
func NewFiberServer(app *fiber.App, appName, monitorPath string) *FiberServer {
	// Note: The database.New() function and the database.service function that takes the database as a parameter are safe from multiple calls (e.g., 10,000 calls from different parts of the codebase)
	// because they follow the singleton pattern. Without the singleton pattern, it would be unsafe
	// as it would create multiple database connections, leading to potential resource exhaustion.
	db := database.New()
	s := &FiberServer{
		App: app,
		db:  db,
	}
	middleware.RegisterRoutes(app, appName, monitorPath, db)
	return s
}

// Start runs the Fiber server in a separate goroutine to listen for incoming requests.
func (s *FiberServer) Start(addr, monitorPath string, tlsConfig *tls.Config, streamListener net.Listener) {
	// Important: Do not modify the current implementation of the HTTPS/TLS mechanism (e.g., by removing the tlsHandler struct).
	// This implementation is similar to the default Fiber setup, with the key difference being that it can be customized to support
	// both HTTP and HTTPS/TLS protocols with custom ports and can operate over the network if you have experience with HTTP.
	// Additionally, this implementation works effectively and stably on Kubernetes with Ingress NGINX without terminating the HTTPS/TLS.
	// Modifying this (e.g., removing the tlsHandler struct) may lead to security issues (e.g., vulnerabilities, CVEs).
	tlsHandler := &fiber.TLSHandler{}
	go func() {
		// TODO: Improve the Listener by creating another Fiber app when tlsConfig and streamListener are configured. This way, it can connect to other Fiber apps (Sharing is caring).
		if tlsConfig != nil && streamListener != nil {
			// Note: This branch handles Boring TLS 1.3 scenarios where the TLS configuration is provided in "run.go".
			// However, Boring TLS 1.3 is currently unavailable.
			if err := s.App.Listener(streamListener); err != nil {
				log.LogFatalf(ErrorHTTPListenAndServe, err)
			}
		} else if tlsConfig != nil {
			// Note: This branch handles standard TLS 1.3 scenarios where the TLS configuration is provided in "run.go".
			// It Force TLS 1.3, due Fiber wrong implementation, about ListenTLS related in "ListenTLSWithCertificate"
			// it should be "if tlsConfig == nil" then load default instead of using "config := &tls.Config".
			ln, err := net.Listen(s.App.Config().Network, addr)
			if err != nil {
				log.LogFatal(err)
			}

			// Note: This TLS handler mechanism uses a pointer to tlsHandler.GetClientInfo,
			// allowing it to be bound here, for example:
			//
			// tlsConfig.GetCertificate = tlsHandler.GetClientInfo
			// However, this is optional. In Kubernetes with Ingress-nginx, you typically don't need
			// "tlsConfig.GetCertificate = tlsHandler.GetClientInfo" because you can use eBPF instead.
			//
			// For mTLS, you might set "tlsConfig.GetCertificate = tlsHandler.GetClientInfo",
			// but this is optional, as mTLS is not always necessary when TLS is supported.
			// Don't be misled into setting up mTLS if TLS is sufficient, unless you specifically need mTLS for IoT.
			// However, connecting to this Boilerplate can also be challenging for IoT itself due to its capabilities.
			tlsListener := tls.NewListener(ln, tlsConfig)
			s.App.SetTLSHandler(tlsHandler)
			// Pass the TLS listener directly to the Fiber app
			if err := s.App.Listener(tlsListener); err != nil {
				log.LogFatalf(ErrorHTTPListenAndServe, err)
			}
		} else {
			// Note: This branch handles TLS 1.2 scenarios or TLS 1.3 when run as a receiver forwarder (e.g. from nginx (Non Kubernetes), Ingress from nginx if it's running on Kubernetes)
			// due to its non-secure nature and requirement to be in internal/development mode.
			//
			// TODO: Use a custom listener in the Fiber framework (similar to SetupTLS) to disable prefork,
			// as it is not compatible with Kubernetes (could this lead to an OOM error for this)?
			if err := s.App.Listen(addr); err != nil {
				log.LogFatalf(ErrorHTTPListenAndServe, err)
			}
		}
	}()

	// Start the HTTP server for redirecting to HTTPS only if TLS is configured
	// this actually work lmao 2 goroutine listening
	//
	// Note: When using Kubernetes with Ingress NGINX, this HTTP redirect is not necessary
	// because NGINX handles traffic at the front end. It is recommended not to terminate
	// HTTPS/TLS in Ingress NGINX, as doing so can degrade performance due to concurrency
	// issues (a logical flaw, related nginx itself not a bug) and pose security risks.
	// This is because traffic is unencrypted (easily get compromised)
	// if HTTPS/TLS is terminated at the ingress and then forwarded as HTTP.
	if tlsConfig != nil && !disableHTTPInsecure {
		// TODO: Improve this that can be customize
		go func() {
			httpAddr := ":80" // Listen on port 80 for HTTP
			s.httpServer = &http.Server{
				Addr: httpAddr,
				Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					httpsPort := strings.Split(addr, ":")[1]
					portPart := ""
					if httpsPort != "443" {
						portPart = ":" + httpsPort
					}
					target := httpsURI + r.Host + portPart + r.URL.RequestURI()

					// Extract the host from the API subdomain
					apiHost := strings.Split(apiSubdomain, ":")[0]

					// Check if the request is for the API subdomain
					if apiHost != "" && strings.HasPrefix(r.Host, apiHost) {
						// Note: Use a 308 redirect for REST APIs to preserve the HTTP method and body.
						// A 301 redirect is better for SEO, as it is well-recognized by search engines for business with google, bing, other search engine hahaha.
						http.Redirect(w, r, target, http.StatusPermanentRedirect) // 308 redirect for API
					} else {
						http.Redirect(w, r, target, http.StatusMovedPermanently) // 301 redirect for others
					}
				}),
			}
			if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.LogFatalf("Error starting HTTP redirect server: %v", err)
			}
		}()
	}

}

// Shutdown gracefully stops the Fiber server using the provided timeout and context for the HTTP insecure server.
//
// Note: This enhances the handling of long-running and keep-alive HTTP/HTTPS connections.
// It remains stable with HPA, but bandwidth usage might increase due to keep-alive connections (e.g., 10 Pods with an average of 1MB/s each).
// Memory leaks could occur if the network in the ingress becomes unstable, such as bottlenecks in ingress-nginx caused by complex configurations.
// For bandwidth considerations, consider hosting in a cost-effective cloud environment (e.g., DOKS).
func (s *FiberServer) Shutdown(ctx context.Context, shutdownTimeout time.Duration) error {
	// http server (insecure) it will be first
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.LogErrorf("Error shutting down HTTP server (insecure): %v", err)
			return err
		}
	}
	if err := s.App.ShutdownWithTimeout(shutdownTimeout); err != nil {
		log.LogErrorf("Error during server shutdown: %v", err)
		return err
	}
	return nil
}

// CleanupDB closes the database connection and Redis client, then logs the outcome.
func (s *FiberServer) CleanupDB() error {
	var err error

	// If the database service is present, close it which will close both the SQL db and Redis connections
	if s.db != nil {
		// Use the common Go idiom style for concise error checking.
		if err = s.db.Close(); err != nil {
			log.LogErrorf("Error closing the database service: %v", err)
			// Do not return here yet, to ensure all cleanup is attempted
		} else {
			log.LogInfo("Database service connections closed.")
		}
	}

	// Return the last error encountered, if any
	return err
}

// Mount mounts a Fiber application or a group of routes onto the main application.
func (s *FiberServer) Mount(prefix string, app any) {
	// Note: It seems possible to integrate it with gRPC (protoc) for internal services,
	// but it's not really needed at the moment.
	switch v := app.(type) {
	case *fiber.App:
		s.App.Mount(prefix, v)
	case func(router fiber.Router):
		group := s.App.Group(prefix)
		v(group)
	default:
		panic(fmt.Errorf("unknown type for mounting: %T", v))
	}
}

// MountPath mounts a Fiber handler or a group of routes onto the main application at a specific path.
func (s *FiberServer) MountPath(path string, handler any) {
	// Note: It seems possible to integrate it with gRPC (protoc) for internal services,
	// but it's not really needed at the moment.
	switch v := handler.(type) {
	case fiber.Handler:
		s.App.Get(path, v)
	case func(router fiber.Router):
		group := s.App.Group(path)
		v(group)
	default:
		panic(fmt.Errorf("unknown type for mounting path: %T", v))
	}
}

// StartServer initializes and starts the server, then waits for a shutdown signal.
// It manages the lifecycle of the server, including graceful shutdown.
func StartServer(server Server, addr, monitorPath string, shutdownTimeout time.Duration, tlsConfig *tls.Config, streamListener net.Listener) {
	startServerAsync(server, addr, monitorPath, tlsConfig, streamListener)
	waitForShutdownSignal(shutdownTimeout, server)
}

// startServerAsync initiates the server's start process in a non-blocking manner.
func startServerAsync(server Server, addr, monitorPath string, tlsConfig *tls.Config, streamListener net.Listener) {
	server.Start(addr, monitorPath, tlsConfig, streamListener)
}

// waitForShutdownSignal listens for OS interrupt or SIGTERM signals to gracefully shut down the server.
// It ensures that the server attempts to shut down gracefully within the provided timeout duration.
func waitForShutdownSignal(shutdownTimeout time.Duration, server Server) {
	quit := make(chan os.Signal, 1)                    // Buffer is one to ensure the signal can be received immediately.
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM) // Notify on interrupt and SIGTERM signals.

	sig := <-quit // Block until a signal is received.
	log.LogInfof(MsgServerShutdown, sig)

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel() // Ensure resources are released.

	// Note: It is important to provide a reasonable shutdownTimeout for the goroutine closure
	// below, because if shutdownTimeout is set to 0s, the goroutine will force a shutdown
	// regardless. Therefore, a recommended timeout is 5 seconds or more to allow for graceful
	// shutdown activities to proceed.
	go func() {
		defer cancel() // Cancel the context when this goroutine completes.
		if err := server.Shutdown(ctx, shutdownTimeout); err != nil {
			// Handle shutdown error.
			log.LogErrorf(MsgErrorDuringShutdown, err)
		}
		if err := server.CleanupDB(); err != nil {
			// Handle cleanup error.
			log.LogErrorf(MsgDatabaseCleanupFailed, err)
		}
	}()

	// Block until the context is done, which occurs when cancel is called or the shutdownTimeout is exceeded.
	<-ctx.Done()
	err := ctx.Err()
	switch err {
	case context.Canceled:
		log.LogInfo(MsgServerShutdownCompleted)
	case context.DeadlineExceeded:
		log.LogError(MsgServerShutdownExceedTimeout)
	default:
		// Typically this shouldn't happen as ctx.Err() should only return nil, context.Canceled,
		// or context.DeadlineExceeded according to the current context package implementation
		// Logging the unexpected error for diagnostic purposes
		log.LogErrorf("An unexpected error occurred during shutdown: %v", err)
	}

}
