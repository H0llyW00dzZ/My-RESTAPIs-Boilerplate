// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package health

import (
	"github.com/gofiber/fiber/v2"

	"h0llyw00dz-template/backend/internal/database"
	log "h0llyw00dz-template/backend/internal/logger"
)

// Response represents the structured response for the health statistics.
type Response struct {
	MySQLHealth *MySQLHealth `json:"mysql_health,omitempty"`
	RedisHealth *RedisHealth `json:"redis_health,omitempty"`
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
	Status           string              `json:"status"`
	Message          string              `json:"message"`
	Error            string              `json:"error,omitempty"`
	Version          string              `json:"version,omitempty"`
	Mode             string              `json:"mode,omitempty"`
	ConnectedClients string              `json:"connected_clients,omitempty"`
	UsedMemory       MemoryUsage         `json:"used_memory,omitempty"`
	PeakUsedMemory   MemoryUsage         `json:"peak_used_memory,omitempty"`
	UptimeStats      string              `json:"uptime_stats,omitempty"`
	Uptime           []map[string]string `json:"uptime,omitempty"`
}

// MemoryUsage represents memory usage in both megabytes and gigabytes.
type MemoryUsage struct {
	MB string `json:"mb,omitempty"`
	GB string `json:"gb,omitempty"`
}

// DBHandler is a Fiber handler that checks the health of the database and Redis.
// It logs the user activity and the health status of MySQL and Redis.
// The detailed health statistics are returned as a structured JSON response.
func DBHandler(db database.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the filter parameter from the query string
		filter := c.Query("filter")

		// Check if the filter is valid
		if !isValidFilter(filter) {
			// TODO: Deal with log errors. Typically, I wouldn't tackle this for StatusBadRequest or StatusNotFound. 🤷‍♂️ 🤪
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid filter parameter. Allowed values: mysql, redis",
			})
		}

		// Log the user activity based on the filter
		logUserActivity(c, filter)

		// Get the health status from the database service
		health := db.Health(filter)

		// Create the response struct using the createHealthResponse function
		response := createHealthResponse(health, filter)

		// Log the health status based on the filter
		logHealthStatus(response, filter)

		// Return the structured health statistics as JSON
		// Note: The "c.JSON" method uses the sonic package (related to main configuration) for JSON encoding and decoding,
		// which is one of the reasons why the Fiber framework is considered the best framework in 2024.
		// "You don't need to repeat yourself for JSON encoding/decoding (e.g., using the standard library or other JSON encoder/decoder)."
		return c.JSON(response)
	}
}

// logUserActivity logs the user activity based on the filter.
func logUserActivity(c *fiber.Ctx, filter string) {
	switch filter {
	case "mysql":
		log.LogUserActivity(c, "viewed the health of the MySQL database")
	case "redis":
		log.LogUserActivity(c, "viewed the health of Redis")
	default:
		log.LogUserActivity(c, "viewed the health of the database and Redis")
	}
}

// createHealthResponse creates a Response struct from the provided health statistics.
func createHealthResponse(health map[string]string, filter string) Response {
	// Note: By structuring the code this way, it is easily maintainable for customization, etc.
	response := Response{}

	if filter == "" || filter == "mysql" {
		response.MySQLHealth = createMySQLHealthResponse(health)
	}

	if filter == "" || filter == "redis" {
		response.RedisHealth = createRedisHealthResponse(health)
	}

	return response
}
