// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package health

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"h0llyw00dz-template/backend/internal/database"
	log "h0llyw00dz-template/backend/internal/logger"
)

// Response represents the structured response for the health statistics.
type Response struct {
	MySQLHealth MySQLHealth `json:"mysql_health"`
	RedisHealth RedisHealth `json:"redis_health"`
}

// MySQLHealth represents the health statistics for MySQL.
type MySQLHealth struct {
	Status          string `json:"status"`
	Message         string `json:"message"`
	Error           string `json:"error,omitempty"`
	OpenConnections string `json:"open_connections,omitempty"`
	InUse           string `json:"in_use,omitempty"`
	Idle            string `json:"idle,omitempty"`
	WaitCount       string `json:"wait_count,omitempty"`
	WaitDuration    string `json:"wait_duration,omitempty"`
}

// RedisHealth represents the health statistics for Redis.
type RedisHealth struct {
	Status           string `json:"status"`
	Message          string `json:"message"`
	Error            string `json:"error,omitempty"`
	Version          string `json:"version,omitempty"`
	Mode             string `json:"mode,omitempty"`
	ConnectedClients string `json:"connected_clients,omitempty"`
	UsedMemory       string `json:"used_memory,omitempty"`
	PeakUsedMemory   string `json:"peak_used_memory,omitempty"`
	Uptime           string `json:"uptime,omitempty"`
}

// DBHandler is a Fiber handler that checks the health of the database and Redis.
// It logs the user activity and the health status of MySQL and Redis.
// The detailed health statistics are returned as a structured JSON response.
func DBHandler(db database.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Log the user activity
		log.LogUserActivity(c, "viewed the health of the database and Redis")

		// Get the health status from the database service
		health := db.Health()

		// Create the response struct using the createHealthResponse function
		response := createHealthResponse(health)

		// Log the MySQL health status
		if response.MySQLHealth.Status == "up" {
			log.LogInfof("MySQL Status: %s, Stats: Open Connections: %s, In Use: %s, Idle: %s, Wait Count: %s, Wait Duration: %s",
				response.MySQLHealth.Message, response.MySQLHealth.OpenConnections, response.MySQLHealth.InUse,
				response.MySQLHealth.Idle, response.MySQLHealth.WaitCount, response.MySQLHealth.WaitDuration)
		} else {
			// If the MySQL status key is missing, log an error
			log.LogErrorf("MySQL Error: %v", response.MySQLHealth.Error)
		}

		// Log the Redis health status
		if response.RedisHealth.Status == "up" {
			log.LogInfof("Redis Status: %s, Stats: Version: %s, Mode: %s, Connected Clients: %s, Used Memory: %s, Peak Used Memory: %s, Uptime: %s",
				response.RedisHealth.Message, response.RedisHealth.Version, response.RedisHealth.Mode,
				response.RedisHealth.ConnectedClients, response.RedisHealth.UsedMemory, response.RedisHealth.PeakUsedMemory,
				response.RedisHealth.Uptime)
		} else {
			// If the Redis status key is missing, log an error
			log.LogErrorf("Redis Error: %v", response.RedisHealth.Error)
		}

		// Return the structured health statistics as JSON
		// Note: The "c.JSON" method uses the sonic package (related to main configuration) for JSON encoding and decoding,
		// which is one of the reasons why the Fiber framework is considered the best framework in 2024.
		// "You don't need to repeat yourself for JSON encoding/decoding (e.g., using the standard library or other JSON encoder/decoder)."
		return c.JSON(response)
	}
}

// createHealthResponse creates a Response struct from the provided health statistics.
func createHealthResponse(health map[string]string) Response {
	// Convert used memory and peak used memory to megabytes (MB)
	usedMemoryMB := bytesToMB(health["redis_used_memory"])
	peakUsedMemoryMB := bytesToMB(health["redis_used_memory_peak"])
	// Format the uptime
	formattedUptime := formatUptime(health["redis_uptime_in_seconds"])

	// Note: By structuring the code this way, it is easily maintainable for customization,etc.
	// Also note that, this method no need to use pointer to match into struct,
	// as pointers are typically recommended for database interfaces
	// (e.g., implementing database interfaces for viewing table data and values in JSON format, such as in real-time database systems).
	return Response{
		MySQLHealth: MySQLHealth{
			Status:          health["mysql_status"],
			Message:         health["mysql_message"],
			Error:           health["mysql_error"],
			OpenConnections: health["mysql_open_connections"],
			InUse:           health["mysql_in_use"],
			Idle:            health["mysql_idle"],
			WaitCount:       health["mysql_wait_count"],
			WaitDuration:    health["mysql_wait_duration"],
		},
		RedisHealth: RedisHealth{
			Status:           health["redis_status"],
			Message:          health["redis_message"],
			Error:            health["redis_error"],
			Version:          health["redis_version"],
			Mode:             health["redis_mode"],
			ConnectedClients: health["redis_connected_clients"],
			UsedMemory:       fmt.Sprintf("%.2f MB", usedMemoryMB),
			PeakUsedMemory:   fmt.Sprintf("%.2f MB", peakUsedMemoryMB),
			Uptime:           formattedUptime,
		},
	}
}
