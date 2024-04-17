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
	bytes, _ := strconv.ParseFloat(bytesStr, 64)
	return bytes / (1024 * 1024)
}

// formatUptime converts the uptime from seconds to a more readable format (days, hours, minutes)
func formatUptime(uptimeSeconds string) (string, []map[string]string) {
	// Parse the uptime seconds string to a float64
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

	// Create the formatted uptime string
	uptimeStats := fmt.Sprintf("%d days, %d hours, %d minutes", days, hours, minutes)

	// Create an array of objects with labels and values for days, hours, and minutes
	uptime := []map[string]string{
		{"day": strconv.Itoa(days)},
		{"hour": strconv.Itoa(hours)},
		{"minute": strconv.Itoa(minutes)},
	}

	// Return the formatted uptime string and the array of uptime objects
	return uptimeStats, uptime
}
