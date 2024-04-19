// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package health

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
