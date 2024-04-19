// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package health

import "fmt"

// RedisHealth represents the health statistics for Redis.
type RedisHealth struct {
	Status           string              `json:"status"`
	Message          string              `json:"message"`
	Error            string              `json:"error,omitempty"`
	Version          string              `json:"version,omitempty"`
	Mode             string              `json:"mode,omitempty"`
	ConnectedClients string              `json:"connected_clients,omitempty"`
	UsedMemory       MemoryUsage         `json:"used_memory,omitempty"`
	PeakUsedMemory   MemoryUsage         `json:"peak_used_memory,omitempty"`
	UptimeStats      string              `json:"uptime_stats,omitempty"`
	Uptime           []map[string]string `json:"uptime,omitempty"`
}

// MemoryUsage represents memory usage in both megabytes and gigabytes.
type MemoryUsage struct {
	MB string `json:"mb,omitempty"`
	GB string `json:"gb,omitempty"`
}

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
