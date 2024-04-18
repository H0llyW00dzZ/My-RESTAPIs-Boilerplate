// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package health

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
