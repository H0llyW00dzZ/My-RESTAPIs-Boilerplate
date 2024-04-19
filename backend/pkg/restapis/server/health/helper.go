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
	// Define a map of filter-specific logging functions
	loggers := map[string]func(Response){
		"mysql": logMySQLHealthStatus,
		"redis": logRedisHealthStatus,
	}

	// Check if the filter is empty or exists in the loggers map
	if filter == "" {
		// If the filter is empty, log the health status for all available filters
		for _, logger := range loggers {
			logger(response)
		}
	} else if logger, ok := loggers[filter]; ok {
		// If the filter exists in the loggers map, log the corresponding health status
		logger(response)
	}
}
