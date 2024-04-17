// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package health

import (
	"fmt"
	"strconv"
)

// bytesToMB converts bytes to megabytes (MB)
func bytesToMB(bytesStr string) float64 {
	// Note: Error handling is omitted here as it would unnecessarily complicate the code
	// for a simple conversion/reformatting operation.
	bytes, _ := strconv.ParseFloat(bytesStr, 64)
	return bytes / (1024 * 1024)
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
