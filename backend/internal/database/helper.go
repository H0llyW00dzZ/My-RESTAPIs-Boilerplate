// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package database

import (
	log "h0llyw00dz-template/backend/internal/logger"
	"strings"
	"time"
)

// parseDateAdded parses the date_added field from a byte slice.
//
// Note: This helper is useful for MySQL, for example, when you need to convert a time from MySQL into JSON or plain text.
func parseDateAdded(dateAddedBytes []uint8) (time.Time, error) {
	const layout = "2006-01-02 15:04:05" // Define the layout constant
	dateAdded, err := time.Parse(layout, string(dateAddedBytes))
	if err != nil {
		log.LogErrorf("Error parsing date_added: %v", err)
		return time.Time{}, err
	}
	return dateAdded, nil
}

// parseRedisInfo parses the Redis info response and returns a map of key-value pairs.
func parseRedisInfo(info string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(info, "\r\n")
	for _, line := range lines {
		if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			result[key] = value
		}
	}
	return result
}
