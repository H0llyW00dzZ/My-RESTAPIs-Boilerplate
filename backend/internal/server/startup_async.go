// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"h0llyw00dz-template/backend/internal/database"
	log "h0llyw00dz-template/backend/internal/logger"

	"github.com/gofiber/fiber/v2"
)

// Server defines the interface for a server that can be started, shut down, and clean up its database.
type Server interface {
	Start(addr, monitorPath string)
	Shutdown(ctx context.Context) error
	CleanupDB() error
}

// FiberServer implements the Server interface for a Fiber application.
type FiberServer struct {
	app *fiber.App
	db  database.Service
}

// NewFiberServer returns a new FiberServer with the given Fiber app, application name, and monitor path.
// It also initializes the database and registers routes.
func NewFiberServer(app *fiber.App, appName, monitorPath string) *FiberServer {
	db := database.New()
	return &FiberServer{
		app: app,
		db:  db,
	}
}

// Start runs the Fiber server in a separate goroutine to listen for incoming requests.
func (s *FiberServer) Start(addr, monitorPath string) {
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

// CleanupDB closes the database connection and logs the outcome.
func (s *FiberServer) CleanupDB() error {
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			log.LogErrorf(ErrorClosingTheDatabase, err)
			return err
		}
		log.LogInfo(MsgDatabaseConnectionClosed)
	}
	return nil
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
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	sig := <-quit
	log.LogInfof(MsgServerShutdown, sig)

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel() // Always call cancel to release resources.

	// Start graceful shutdown in a separate goroutine.
	go func() {
		if err := server.Shutdown(ctx); err != nil {
			log.LogErrorf(MsgErrorDuringShutdown, err)
		}

		if err := server.CleanupDB(); err != nil {
			log.LogErrorf(MsgDatabaseCleanupFailed, err)
		}

		cancel() // Signal shutdown completion.
	}()

	select {
	case <-ctx.Done():
		log.LogInfo(MsgServerShutdownCompleted)
	case <-time.After(shutdownTimeout):
		log.LogError(MsgServerShutdownExceedTimeout)
		cancel() // Ensure resources are released.
	}
}
