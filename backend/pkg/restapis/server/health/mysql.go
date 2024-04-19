// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package health

import log "h0llyw00dz-template/backend/internal/logger"

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

// createMySQLHealthResponse creates a MySQLHealth struct from the provided health statistics.
func createMySQLHealthResponse(health map[string]string) *MySQLHealth {
	return &MySQLHealth{
		Status:          health["mysql_status"],
		Message:         health["mysql_message"],
		Error:           health["mysql_error"],
		OpenConnections: health["mysql_open_connections"],
		InUse:           health["mysql_in_use"],
		Idle:            health["mysql_idle"],
		WaitCount:       health["mysql_wait_count"],
		WaitDuration:    health["mysql_wait_duration"],
	}
}

// logMySQLHealthStatus logs the MySQL health status.
func logMySQLHealthStatus(response Response) {
	if response.MySQLHealth.Status == "up" {
		log.LogInfof("MySQL Status: %s, Stats: Open Connections: %s, In Use: %s, Idle: %s, Wait Count: %s, Wait Duration: %s",
			response.MySQLHealth.Message, response.MySQLHealth.OpenConnections, response.MySQLHealth.InUse,
			response.MySQLHealth.Idle, response.MySQLHealth.WaitCount, response.MySQLHealth.WaitDuration)
	} else {
		log.LogErrorf("MySQL Error: %v", response.MySQLHealth.Error)
	}
}
