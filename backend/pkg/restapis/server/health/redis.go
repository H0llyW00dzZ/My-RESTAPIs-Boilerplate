// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package health

import (
	"fmt"
	log "h0llyw00dz-template/backend/internal/logger"
)

// MemoryUsage represents memory usage in both megabytes and gigabytes.
type MemoryUsage struct {
	MB string `json:"mb,omitempty"`
	GB string `json:"gb,omitempty"`
}

// PoolingStats represents the statistics of the current connection pooling state.
type PoolingStats struct {
	Total  string `json:"total,omitempty"`
	Idle   string `json:"idle,omitempty"`
	Active string `json:"active,omitempty"`
}

// RedisStats groups the statistics related to Redis.
type RedisStats struct {
	Version          string              `json:"version,omitempty"`
	Mode             string              `json:"mode,omitempty"`
	ConnectedClients string              `json:"connected_clients,omitempty"`
	UsedMemory       MemoryUsage         `json:"used_memory,omitempty"`
	PeakUsedMemory   MemoryUsage         `json:"peak_used_memory,omitempty"`
	UptimeStats      string              `json:"uptime_stats,omitempty"`
	Uptime           []map[string]string `json:"uptime,omitempty"`
	Pooling          PoolingStats        `json:"pooling,omitempty"`
	ServerFreeMemory MemoryUsage         `json:"server_free_memory,omitempty"`
}

// RedisHealth represents the health statistics for Redis.
type RedisHealth struct {
	Status  string     `json:"status"`
	Message string     `json:"message"`
	Error   string     `json:"error,omitempty"`
	Stats   RedisStats `json:"stats,omitempty"`
}

// createRedisHealthResponse creates a RedisHealth struct from the provided health statistics.
func createRedisHealthResponse(health map[string]string) *RedisHealth {
	// Convert used memory and peak used memory to megabytes (MB) and gigabytes (GB)
	// Note: The gigabyte value will show a nonzero value starting from 0.01GB, which is approximately 10MB.
	usedMemoryMB, usedMemoryGB := bytesToMBGB(health["redis_used_memory"])
	peakUsedMemoryMB, peakUsedMemoryGB := bytesToMBGB(health["redis_used_memory_peak"])
	serverFreeMemoryMB, serverFreeMemoryGB := bytesToMBGB(health["redis_max_memory"])
	// Format the uptime
	uptimeStats, uptime := formatUptime(health["redis_uptime_in_seconds"])
	// Create PoolingStats from the health statistics
	poolingStats := PoolingStats{
		Total:  health["redis_total_connections"],
		Idle:   health["redis_idle_connections"],
		Active: health["redis_active_connections"],
	}

	// Create RedisStats from the health statistics
	redisStats := RedisStats{
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
		Pooling:     poolingStats,
		// Better formatting it should be raw "%.2f"
		ServerFreeMemory: MemoryUsage{
			MB: fmt.Sprintf("%.2f", serverFreeMemoryMB),
			GB: fmt.Sprintf("%.2f", serverFreeMemoryGB),
		},
	}

	return &RedisHealth{
		Status:  health["redis_status"],
		Message: health["redis_message"],
		Error:   health["redis_error"],
		Stats:   redisStats, // This now contains all the Redis related stats
	}
}

// logRedisHealthStatus logs the Redis health status.
func logRedisHealthStatus(response Response) {
	if response.RedisHealth.Status == "up" {
		log.LogInfof("Redis Status: %s, Stats: Version: %s, Mode: %s, Used Memory: %s MB (%s GB), Peak Used Memory: %s MB (%s GB), Uptime: %s, Total Connections: %s, Active Connections: %s, Idle Connections: %s, Server Free Memory: %s MB (%s GB)",
			response.RedisHealth.Message, response.RedisHealth.Stats.Version, response.RedisHealth.Stats.Mode,
			response.RedisHealth.Stats.UsedMemory.MB, response.RedisHealth.Stats.UsedMemory.GB,
			response.RedisHealth.Stats.PeakUsedMemory.MB, response.RedisHealth.Stats.PeakUsedMemory.GB, response.RedisHealth.Stats.UptimeStats,
			response.RedisHealth.Stats.Pooling.Total, response.RedisHealth.Stats.Pooling.Active, response.RedisHealth.Stats.Pooling.Idle,
			response.RedisHealth.Stats.ServerFreeMemory.MB, response.RedisHealth.Stats.ServerFreeMemory.GB)
	} else {
		log.LogErrorf("Redis Error: %v", response.RedisHealth.Error)
	}
}
