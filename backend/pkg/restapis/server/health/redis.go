// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package health

import (
	"fmt"
	log "h0llyw00dz-template/backend/internal/logger"
	"strconv"
)

// MemoryUsage represents memory usage in both megabytes and gigabytes.
type MemoryUsage struct {
	MB string `json:"mb,omitempty"`
	GB string `json:"gb,omitempty"`
}

// PoolingStats represents the statistics of the current connection pooling state.
type PoolingStats struct {
	Total              string `json:"total,omitempty"`
	Idle               string `json:"idle,omitempty"`
	Active             string `json:"active,omitempty"`
	PoolSizePercentage string `json:"percentage,omitempty"`
}

// MemoryStats represents the memory usage statistics.
type MemoryStats struct {
	Used       MemoryUsage `json:"used,omitempty"`
	Peak       MemoryUsage `json:"peak,omitempty"`
	Percentage string      `json:"percentage,omitempty"`
	Free       MemoryUsage `json:"free,omitempty"`
}

// RedisStats groups the statistics related to Redis.
type RedisStats struct {
	Version          string              `json:"version,omitempty"`
	Mode             string              `json:"mode,omitempty"`
	ConnectedClients string              `json:"connected_clients,omitempty"`
	Memory           MemoryStats         `json:"memory,omitempty"`
	UptimeStats      string              `json:"uptime_stats,omitempty"`
	Uptime           []map[string]string `json:"uptime,omitempty"`
	Pooling          PoolingStats        `json:"pooling,omitempty"`
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
	redisHealth := &RedisHealth{
		Status:  health["redis_status"],
		Message: health["redis_message"],
		Error:   health["redis_error"],
	}

	// Only populate the Stats field if Redis is up and running
	if health["redis_status"] == "up" {
		// Convert used memory and peak used memory to megabytes (MB) and gigabytes (GB)
		// Note: The gigabyte value will show a nonzero value starting from 0.01GB, which is approximately 10MB.
		usedMemoryMB, usedMemoryGB := bytesToMBGB(health["redis_used_memory"])
		peakMemoryMB, peakMemoryGB := bytesToMBGB(health["redis_used_memory_peak"])
		freeMemoryMB, freeMemoryGB := bytesToMBGB(health["redis_max_memory"])

		// Calculate the memory usage percentage
		usedMemory, _ := strconv.ParseInt(health["redis_used_memory"], 10, 64)
		maxMemory, _ := strconv.ParseInt(health["redis_max_memory"], 10, 64)
		memoryUsage := calculateMemoryUsage(usedMemory, maxMemory)

		// Format the uptime
		uptimeStats, uptime := formatUptime(health["redis_uptime_in_seconds"])

		// Create PoolingStats from the health statistics
		poolingStats := PoolingStats{
			Total:  health["redis_total_connections"],
			Idle:   health["redis_idle_connections"],
			Active: health["redis_active_connections"],
		}

		redisHealth.Stats = RedisStats{
			Version:          health["redis_version"],
			Mode:             health["redis_mode"],
			ConnectedClients: health["redis_connected_clients"],
			Memory: MemoryStats{
				Used: MemoryUsage{
					// Better formatting it should be raw "%.2f"
					MB: fmt.Sprintf("%.2f", usedMemoryMB),
					GB: fmt.Sprintf("%.2f", usedMemoryGB),
				},
				Peak: MemoryUsage{
					// Better formatting it should be raw "%.2f"
					MB: fmt.Sprintf("%.2f", peakMemoryMB),
					GB: fmt.Sprintf("%.2f", peakMemoryGB),
				},
				Percentage: memoryUsage,
				Free: MemoryUsage{
					// Better formatting it should be raw "%.2f"
					MB: fmt.Sprintf("%.2f", freeMemoryMB),
					GB: fmt.Sprintf("%.2f", freeMemoryGB),
				},
			},
			UptimeStats: uptimeStats,
			Uptime:      uptime,
			Pooling:     poolingStats,
		}
	}

	return redisHealth
}

// logRedisHealthStatus logs the Redis health status.
func logRedisHealthStatus(response Response) {
	if response.RedisHealth.Status == "up" {
		log.LogInfof("Redis Status: %s, Stats: Version: %s, Mode: %s, Used Memory: %s MB (%s GB), Peak Memory: %s MB (%s GB), Memory Usage: %s, Free Memory: %s MB (%s GB), Uptime: %s, Total Connections: %s, Active Connections: %s, Idle Connections: %s, Pool Size Usage: %s",
			response.RedisHealth.Message, response.RedisHealth.Stats.Version, response.RedisHealth.Stats.Mode,
			response.RedisHealth.Stats.Memory.Used.MB, response.RedisHealth.Stats.Memory.Used.GB,
			response.RedisHealth.Stats.Memory.Peak.MB, response.RedisHealth.Stats.Memory.Peak.GB,
			response.RedisHealth.Stats.Memory.Percentage,
			response.RedisHealth.Stats.Memory.Free.MB, response.RedisHealth.Stats.Memory.Free.GB,
			response.RedisHealth.Stats.UptimeStats,
			response.RedisHealth.Stats.Pooling.Total, response.RedisHealth.Stats.Pooling.Active, response.RedisHealth.Stats.Pooling.Idle,
			response.RedisHealth.Stats.Pooling.PoolSizePercentage)
	} else {
		log.LogErrorf("Redis Error: %v", response.RedisHealth.Error)
	}
}
