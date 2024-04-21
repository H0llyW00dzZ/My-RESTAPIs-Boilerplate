// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package health

import log "h0llyw00dz-template/backend/internal/logger"

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
	Message string           `json:"message"`
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

// logMySQLHealthStatus logs the MySQL health status.
func logMySQLHealthStatus(response Response) {
	if response.MySQLHealth.Status == "up" {
		log.LogInfof("MySQL Status: %s, Stats: Open Connections: %s, In Use: %s, Idle: %s, Wait Count: %s, Duration: %s",
			response.MySQLHealth.Message, response.MySQLHealth.Stats.Open, response.MySQLHealth.Stats.InUse,
			response.MySQLHealth.Stats.Idle, response.MySQLHealth.Stats.WaitCount, response.MySQLHealth.Stats.Duration)
	} else {
		log.LogErrorf("MySQL Error: %v", response.MySQLHealth.Error)
	}
}
