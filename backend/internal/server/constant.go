// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package server

const (
	// MsgServerStart is the log message indicating that the server is starting.
	MsgServerStart = "Starting server on %s"
	// MsgServerShutdown is the log message indicating that the server is shutting down.
	MsgServerShutdown = "Shutting down server... reason: %v"
	// MsgErrorDuringShutdown is the log message for an error encountered during server shutdown.
	MsgErrorDuringShutdown = "Error during server shutdown: %v"
	// MsgDatabaseCleanupFailed is the log message for a failure in database cleanup.
	MsgDatabaseCleanupFailed = "Database cleanup failed: %v"
	// MsgServerShutdownCompleted is the log message indicating that the server shutdown process has completed.
	MsgServerShutdownCompleted = "Server shutdown completed."
	// MsgServerShutdownExceedTimeout is the log message indicating that the server shutdown process exceeded the specified timeout.
	MsgServerShutdownExceedTimeout = "Server shutdown exceeded the timeout."
	// MsgDatabaseConnectionClosed is the log message indicating that the database connection has been closed gracefully.
	MsgDatabaseConnectionClosed = "Database connection closed."
	// ErrorHTTPListenAndServe is the error message format for an HTTP server ListenAndServe error.
	ErrorHTTPListenAndServe = "HTTP server ListenAndServe: %v"
	// ErrorClosingTheDatabase is the error message format for an error encountered when closing the database connection.
	ErrorClosingTheDatabase = "Error closing the database: %v"
)
