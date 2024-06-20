// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"h0llyw00dz-template/backend/internal/database"
	log "h0llyw00dz-template/backend/internal/logger"
	"h0llyw00dz-template/backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// Server defines the interface for a server that can be started, shut down, and clean up its database.
type Server interface {
	Start(addr, monitorPath string)
	Shutdown(ctx context.Context) error
	CleanupDB() error
	Mount(prefix string, app any)
	MountPath(path string, handler any)
}

// FiberServer implements the Server interface for a Fiber application.
type FiberServer struct {
	app *fiber.App
	db  database.Service
}

// NewFiberServer returns a new FiberServer with the given Fiber app, application name, and monitor path.
// It also initializes the database and registers routes.
func NewFiberServer(app *fiber.App, appName, monitorPath string) *FiberServer {
	// Note: The database.New() function and the database.service function that takes the database as a parameter are safe from multiple calls (e.g., 10,000 calls from different parts of the codebase)
	// because they follow the singleton pattern. Without the singleton pattern, it would be unsafe
	// as it would create multiple database connections, leading to potential resource exhaustion.
	db := database.New()
	s := &FiberServer{
		app: app,
		db:  db,
	}
	middleware.RegisterRoutes(app, appName, monitorPath, db)
	return s
}

// Start runs the Fiber server in a separate goroutine to listen for incoming requests.
func (s *FiberServer) Start(addr, monitorPath string) {
	// TODO: Implement environment mode. For example, when the environment is set to "dev" or "local", it will switch to "Listen".
	// Otherwise, it will force a change from Listen to ListenTLS to use TLS 1.3 protocols.
	// For the certificate, if used at the enterprise or government level, it should be issued to the organization named "Boring TLS" hahaha.
	go func() {
		log.LogInfof(MsgServerStart, addr)
		if err := s.app.Listen(addr); err != nil {
			log.LogErrorf(ErrorHTTPListenAndServe, err)
		}
	}()
}

// Shutdown gracefully stops the Fiber server using the provided context.
func (s *FiberServer) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}

// CleanupDB closes the database connection and Redis client, then logs the outcome.
func (s *FiberServer) CleanupDB() error {
	var err error

	// If the database service is present, close it which will close both the SQL db and Redis connections
	if s.db != nil {
		err = s.db.Close()
		if err != nil {
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
		s.app.Mount(prefix, v)
	case func(router fiber.Router):
		group := s.app.Group(prefix)
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
		s.app.Get(path, v)
	case func(router fiber.Router):
		group := s.app.Group(path)
		v(group)
	default:
		panic(fmt.Errorf("unknown type for mounting path: %T", v))
	}
}

// StartServer initializes and starts the server, then waits for a shutdown signal.
// It manages the lifecycle of the server, including graceful shutdown.
func StartServer(server Server, addr, monitorPath string, shutdownTimeout time.Duration) {
	startServerAsync(server, addr, monitorPath)
	waitForShutdownSignal(shutdownTimeout, server)
}

// startServerAsync initiates the server's start process in a non-blocking manner.
func startServerAsync(server Server, addr, monitorPath string) {
	server.Start(addr, monitorPath)
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
		if err := server.Shutdown(ctx); err != nil {
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
