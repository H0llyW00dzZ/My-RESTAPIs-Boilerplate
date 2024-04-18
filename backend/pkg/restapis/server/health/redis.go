// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package health

import "fmt"

// createRedisHealthResponse creates a RedisHealth struct from the provided health statistics.
func createRedisHealthResponse(health map[string]string) *RedisHealth {
	// Convert used memory and peak used memory to megabytes (MB) and gigabytes (GB)
	// Note: gigabytes will be showing 0.00GB if under 100MB usage
	usedMemoryMB, usedMemoryGB := bytesToMBGB(health["redis_used_memory"])
	peakUsedMemoryMB, peakUsedMemoryGB := bytesToMBGB(health["redis_used_memory_peak"])
	// Format the uptime
	uptimeStats, uptime := formatUptime(health["redis_uptime_in_seconds"])

	return &RedisHealth{
		Status:           health["redis_status"],
		Message:          health["redis_message"],
		Error:            health["redis_error"],
		Version:          health["redis_version"],
		Mode:             health["redis_mode"],
		ConnectedClients: health["redis_connected_clients"],
		// Better formatting it should be raw "%.2f"
		UsedMemory: MemoryUsage{
			MB: fmt.Sprintf("%.2f", usedMemoryMB),
			GB: fmt.Sprintf("%.2f", usedMemoryGB),
		},
		// Better formatting it should be raw "%.2f"
		PeakUsedMemory: MemoryUsage{
			MB: fmt.Sprintf("%.2f", peakUsedMemoryMB),
			GB: fmt.Sprintf("%.2f", peakUsedMemoryGB),
		},
		UptimeStats: uptimeStats,
		Uptime:      uptime,
	}
}
