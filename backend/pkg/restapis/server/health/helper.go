// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package health

import (
	"fmt"
	log "h0llyw00dz-template/backend/internal/logger"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// validFilters is a map that defines the valid filters and their corresponding log messages.
var validFilters = map[string]string{
	"":      "viewed the health of the database and Redis",
	"mysql": "viewed the health of the MySQL database",
	"redis": "viewed the health of Redis",
}

// bytesToMBGB converts bytes to megabytes (MB) and gigabytes (GB)
func bytesToMBGB(bytesStr string) (float64, float64) {
	// Note: Error handling is omitted here as it would unnecessarily complicate the code
	// for a simple conversion/reformatting operation.
	bytes, _ := strconv.ParseFloat(bytesStr, 64)
	// Note: Terabyte support is not implemented since gigabytes should be sufficient for most use cases in planet earth hahaha.
	mb := bytes / (1024 * 1024)
	gb := bytes / (1024 * 1024 * 1024)
	return mb, gb
}

// formatUptime converts the uptime from seconds to a more readable format (days, hours, minutes, seconds)
func formatUptime(uptimeSeconds string) (string, []map[string]string) {
	// Parse the uptime seconds string to a float64
	// Note: Error handling is omitted here as it would unnecessarily complicate the code
	// for a simple conversion/reformatting operation.
	seconds, _ := strconv.ParseFloat(uptimeSeconds, 64)

	// Calculate the number of days
	days := int(seconds) / (24 * 3600)
	seconds -= float64(days * 24 * 3600)

	// Calculate the number of hours
	hours := int(seconds) / 3600
	seconds -= float64(hours * 3600)

	// Calculate the number of minutes
	minutes := int(seconds) / 60
	seconds -= float64(minutes * 60)

	// Calculate the remaining seconds
	remainingSeconds := int(seconds)

	// Create the formatted uptime string
	uptimeStats := fmt.Sprintf("%d days, %d hours, %d minutes, %d seconds", days, hours, minutes, remainingSeconds)

	// Create an array of objects with labels and values for days, hours, minutes, and seconds
	uptime := []map[string]string{
		{"day": strconv.Itoa(days)},
		{"hour": strconv.Itoa(hours)},
		{"minute": strconv.Itoa(minutes)},
		{"second": strconv.Itoa(remainingSeconds)},
	}

	// Return the formatted uptime string and the array of uptime objects
	return uptimeStats, uptime
}

// isValidFilter checks if the provided filter is valid.
func isValidFilter(filter string) bool {
	_, valid := validFilters[filter]
	return valid
}

// logUserActivity logs the user activity based on the filter.
func logUserActivity(c *fiber.Ctx, filter string) {
	activity, ok := validFilters[filter]
	if ok {
		log.LogUserActivity(c, activity)
	}
}

// logHealthStatus logs the health status based on the filter.
func logHealthStatus(response Response, filter string) {
	if filter == "" || filter == "mysql" {
		// Log the MySQL health status
		if response.MySQLHealth.Status == "up" {
			log.LogInfof("MySQL Status: %s, Stats: Open Connections: %s, In Use: %s, Idle: %s, Wait Count: %s, Wait Duration: %s",
				response.MySQLHealth.Message, response.MySQLHealth.OpenConnections, response.MySQLHealth.InUse,
				response.MySQLHealth.Idle, response.MySQLHealth.WaitCount, response.MySQLHealth.WaitDuration)
		} else {
			// If the MySQL status key is missing, log an error
			log.LogErrorf("MySQL Error: %v", response.MySQLHealth.Error)
		}
	}

	if filter == "" || filter == "redis" {
		// Log the Redis health status
		if response.RedisHealth.Status == "up" {
			log.LogInfof("Redis Status: %s, Stats: Version: %s, Mode: %s, Connected Clients: %s, Used Memory: %s MB (%s GB), Peak Used Memory: %s MB (%s GB), Uptime: %s",
				response.RedisHealth.Message, response.RedisHealth.Version, response.RedisHealth.Mode,
				response.RedisHealth.ConnectedClients, response.RedisHealth.UsedMemory.MB, response.RedisHealth.UsedMemory.GB,
				response.RedisHealth.PeakUsedMemory.MB, response.RedisHealth.PeakUsedMemory.GB, response.RedisHealth.UptimeStats)
		} else {
			// If the Redis status key is missing, log an error
			log.LogErrorf("Redis Error: %v", response.RedisHealth.Error)
		}
	}
}
