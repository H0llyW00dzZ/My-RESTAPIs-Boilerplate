// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package health

import (
	log "h0llyw00dz-template/backend/internal/logger"
	"h0llyw00dz-template/backend/pkg/restapis/helper"

	"github.com/gofiber/fiber/v2"
)

// ConnectionStats represents the statistics of the current database connection state.
type ConnectionStats struct {
	Open      string `json:"open,omitempty"`
	InUse     string `json:"in_use,omitempty"`
	Idle      string `json:"idle,omitempty"`
	WaitCount string `json:"wait_count,omitempty"`
	Duration  string `json:"duration,omitempty"` // Renamed for clarity
}

// MySQLHealth represents the health statistics for MySQL.
type MySQLHealth struct {
	Status  string           `json:"status"`
	Message string           `json:"message,omitempty"`
	Error   string           `json:"error,omitempty"`
	Stats   *ConnectionStats `json:"stats,omitempty"`
}

// createMySQLHealthResponse creates a MySQLHealth struct from the provided health statistics.
func createMySQLHealthResponse(health map[string]string) *MySQLHealth {
	mysqlHealth := &MySQLHealth{
		Status:  health["mysql_status"],
		Message: health["mysql_message"],
		Error:   health["mysql_error"],
	}

	// Only populate the Stats field if MySQL is up and running
	if health["mysql_status"] == "up" {
		mysqlHealth.Stats = &ConnectionStats{
			Open:      health["mysql_open_connections"],
			InUse:     health["mysql_in_use"],
			Idle:      health["mysql_idle"],
			WaitCount: health["mysql_wait_count"],
			Duration:  health["mysql_wait_duration"],
		}
	}

	return mysqlHealth
}

// logMySQLHealthStatus logs the MySQL health status and sends an error response if MySQL is down.
func logMySQLHealthStatus(c *fiber.Ctx, response Response) error {
	// Extract mysqlHealth from the response
	mysqlHealth := response.MySQLHealth

	if mysqlHealth != nil && mysqlHealth.Status == "up" {
		// Log general MySQL status
		// TODO: Improve this by using charm.sh TUI components for a better and more modern experience (not the ancient experience).
		log.LogInfof("MySQL Status: %s, Stats: Open Connections: %s, In Use: %s, Idle: %s, Wait Count: %s, Duration: %s",
			response.MySQLHealth.Message,
			response.MySQLHealth.Stats.Open,
			response.MySQLHealth.Stats.InUse,
			response.MySQLHealth.Stats.Idle,
			response.MySQLHealth.Stats.WaitCount,
			response.MySQLHealth.Stats.Duration)
	} else {
		// Log the error if MySQL is not up or if mysqlHealth is nil
		log.LogErrorf("MySQL Error: %v", mysqlHealth.Error)

		// Send an error response
		// Note: This is dynamic and it's not possible to set the "errorCode" because it depends on internal/database/mysql_redis.go,
		// so it only works to set the HTTP status code as ServiceUnavailable.
		return helper.SendErrorResponse(c, fiber.StatusServiceUnavailable, mysqlHealth.Error)
	}

	return nil
}
