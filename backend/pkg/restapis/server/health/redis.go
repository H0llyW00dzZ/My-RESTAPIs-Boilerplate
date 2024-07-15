// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package health

import (
	"fmt"
	log "h0llyw00dz-template/backend/internal/logger"
	"h0llyw00dz-template/backend/pkg/restapis/helper"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// MemoryUsage represents memory usage in both megabytes and gigabytes.
type MemoryUsage struct {
	MB string `json:"mb,omitempty"`
	GB string `json:"gb,omitempty"`
}

// ConnectionFigures represents the numerical statistics associated with Redis connections.
type ConnectionFigures struct {
	Hits       string `json:"hits,omitempty"`
	Misses     string `json:"misses,omitempty"`
	Timeouts   string `json:"timeouts,omitempty"`
	Total      string `json:"total,omitempty"`
	Stale      string `json:"stale,omitempty"`
	Idle       string `json:"idle,omitempty"`
	Active     string `json:"active,omitempty"`
	Percentage string `json:"percentage,omitempty"`
}

// PoolingStats now contains a ConnectionFigures struct, representing a part of the pooling stats.
type PoolingStats struct {
	Figures       ConnectionFigures `json:"figures,omitempty"`
	ObservedTotal string            `json:"observed_total,omitempty"`
}

// MemoryStats represents the memory usage statistics.
type MemoryStats struct {
	Used       MemoryUsage `json:"used,omitempty"`
	Peak       MemoryUsage `json:"peak,omitempty"`
	Free       MemoryUsage `json:"free,omitempty"`
	Percentage string      `json:"percentage,omitempty"`
}

// RedisStats groups the statistics related to Redis.
type RedisStats struct {
	Version          string       `json:"version,omitempty"`
	Mode             string       `json:"mode,omitempty"`
	ConnectedClients string       `json:"connected_clients,omitempty"`
	Memory           MemoryStats  `json:"memory,omitempty"`
	Uptime           []any        `json:"uptime,omitempty"`
	Pooling          PoolingStats `json:"pooling,omitempty"`
}

// UptimeFields represents the uptime fields in a structured format.
type UptimeFields struct {
	Day    string `json:"day"`
	Hour   string `json:"hour"`
	Minute string `json:"minute"`
	Second string `json:"second"`
}

// RedisHealth represents the health statistics for Redis.
//
// Demo: https://api-beta.btz.pm/v1/health/db?filter=redis (Better REST Formatting)
type RedisHealth struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Stats   *RedisStats `json:"stats,omitempty"`
}

// createRedisHealthResponse creates a RedisHealth struct from the provided health statistics.
func (r *RedisHealth) createRedisHealthResponse(health map[string]string) *RedisHealth {
	r.Status = health["redis_status"]
	r.Message = health["redis_message"]
	r.Error = health["redis_error"]

	// Only populate the Stats field if Redis is up and running
	if health["redis_status"] == "up" {
		// Convert used memory and peak used memory to megabytes (MB) and gigabytes (GB)
		// Note: The gigabyte value will show a nonzero value starting from 0.01GB, which is approximately 10MB.
		usedMemoryMB, usedMemoryGB := bytesToMBGB(health["redis_used_memory"])
		peakMemoryMB, peakMemoryGB := bytesToMBGB(health["redis_used_memory_peak"])
		freeMemoryMB, freeMemoryGB := bytesToMBGB(health["redis_max_memory"])

		// Calculate the memory usage percentage
		usedMemory := helper.ParseInt64Value(health["redis_used_memory"], 10, 64)
		maxMemory := helper.ParseInt64Value(health["redis_max_memory"], 10, 64)
		memoryUsage := calculateMemoryUsage(usedMemory, maxMemory)

		// Format the uptime
		_, uptime := formatUptime(health["redis_uptime_in_seconds"])

		// Parse numerical values from the health stats for calculation
		hits := helper.ParseNumericalValue(health["redis_hits_connections"], 10, 64)
		misses := helper.ParseNumericalValue(health["redis_misses_connections"], 10, 64)
		timeouts := helper.ParseNumericalValue(health["redis_timeouts_connections"], 10, 64)
		activeConns := helper.ParseNumericalValue(health["redis_active_connections"], 10, 64)
		StaleConns := helper.ParseNumericalValue(health["redis_stale_connections"], 10, 64)
		idleConns := helper.ParseNumericalValue(health["redis_idle_connections"], 10, 64)

		// Calculate the observed total connections
		observedTotalConns := hits + misses + timeouts + StaleConns + activeConns + idleConns
		observedTotal := strconv.FormatUint(observedTotalConns, 10)

		// Create PoolingStats from the health statistics
		poolingStats := PoolingStats{
			Figures: ConnectionFigures{
				Total:      health["redis_total_connections"],
				Stale:      health["redis_stale_connections"],
				Idle:       health["redis_idle_connections"],
				Active:     health["redis_active_connections"],
				Hits:       health["redis_hits_connections"],
				Misses:     health["redis_misses_connections"],
				Timeouts:   health["redis_timeouts_connections"],
				Percentage: health["redis_pool_size_percentage"],
			},
			ObservedTotal: observedTotal,
		}

		r.Stats = &RedisStats{
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
			Uptime:  uptime,
			Pooling: poolingStats,
		}
	}

	return r
}

// logRedisHealthStatus logs the Redis health status and sends an error response if Redis is down.
func (r *RedisHealth) logRedisHealthStatus(c *fiber.Ctx, response Response) error {
	// Extract redisHealth from the response
	redisHealth := response.RedisHealth

	// Note: This method `Map of Function` improves data structuring for logging purposes. It provides a clear and
	// efficient way to access the Redis health information, which is crucial for maintaining
	// the integrity of the logged data.
	if redisHealth != nil && redisHealth.Status == "up" {
		// Log general Redis status
		// TODO: Improve this by using charm.sh TUI components for a better and more modern experience (not the ancient experience).
		log.LogInfof("Redis Status: %s, Stats: Version: %s, Mode: %s",
			redisHealth.Message, redisHealth.Stats.Version, redisHealth.Stats.Mode)

		// Log memory usage
		log.LogInfof("Redis Memory Usage: Used: %s MB (%s GB), Peak: %s MB (%s GB), Percentage: %s, Free: %s MB (%s GB)",
			redisHealth.Stats.Memory.Used.MB, redisHealth.Stats.Memory.Used.GB,
			redisHealth.Stats.Memory.Peak.MB, redisHealth.Stats.Memory.Peak.GB,
			redisHealth.Stats.Memory.Percentage,
			redisHealth.Stats.Memory.Free.MB, redisHealth.Stats.Memory.Free.GB)

		// Log uptime stats
		if len(redisHealth.Stats.Uptime) > 1 {
			if stats, ok := redisHealth.Stats.Uptime[1].(map[string]string); ok {
				log.LogInfof("Redis Uptime: %s, Pooling Connections: %s, Connected Clients: %s",
					stats["stats"],
					redisHealth.Stats.Pooling.ObservedTotal,
					redisHealth.Stats.ConnectedClients)
			}
		}

		// Log detailed pooling stats
		log.LogInfof("Redis Pooling Figures: Hits: %s, Misses: %s, Timeouts: %s, Total: %s, Stale: %s, Idle: %s, Active: %s, Pool Size Usage: %s, Observed Total: %s",
			redisHealth.Stats.Pooling.Figures.Hits,
			redisHealth.Stats.Pooling.Figures.Misses,
			redisHealth.Stats.Pooling.Figures.Timeouts,
			redisHealth.Stats.Pooling.Figures.Total,
			redisHealth.Stats.Pooling.Figures.Stale,
			redisHealth.Stats.Pooling.Figures.Idle,
			redisHealth.Stats.Pooling.Figures.Active,
			redisHealth.Stats.Pooling.Figures.Percentage,
			redisHealth.Stats.Pooling.ObservedTotal)
	} else {
		// Log the error if Redis is not up or if redisHealth is nil
		log.LogErrorf("Redis Error: %v", redisHealth.Error)

		// Send an error response
		// Note: This is dynamic and it's not possible to set the "errorCode" because it depends on internal/database/mysql_redis.go,
		// so it only works to set the HTTP status code as ServiceUnavailable.
		return helper.SendErrorResponse(c, fiber.StatusServiceUnavailable, redisHealth.Error)
	}

	return nil
}
