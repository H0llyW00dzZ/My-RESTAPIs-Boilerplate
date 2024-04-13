// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package health

import (
	"github.com/gofiber/fiber/v2"

	"h0llyw00dz-template/backend/internal/database"
	log "h0llyw00dz-template/backend/internal/logger"
)

// DBHandler is a Fiber handler that checks the health of the database and Redis.
// It logs the user activity and the health status of MySQL and Redis.
// The detailed health statistics are returned as JSON.
func DBHandler(db database.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Log the user activity
		log.LogUserActivity(c, "viewed the health of the database and Redis")

		// Get the health status from the database service
		health := db.Health()

		// Log the MySQL health status
		mysqlStatus, ok := health["mysql_status"]
		if !ok {
			// If the MySQL status key is missing, log an error
			log.LogError("MySQL health check did not return a status")
		} else {
			if mysqlStatus == "up" {
				log.LogInfof("MySQL Status: %s, Stats: Open Connections: %s, In Use: %s, Idle: %s, Wait Count: %s, Wait Duration: %s",
					database.MsgDBItsHealthy, health["mysql_open_connections"], health["mysql_in_use"], health["mysql_idle"],
					health["mysql_wait_count"], health["mysql_wait_duration"])
			} else {
				mysqlErrorMessage, _ := health["mysql_error"]
				log.LogErrorf("MySQL Error: %v", mysqlErrorMessage)
			}
		}

		// Log the Redis health status
		redisStatus, ok := health["redis_status"]
		if !ok {
			// If the Redis status key is missing, log an error
			log.LogError("Redis health check did not return a status")
		} else {
			if redisStatus == "up" {
				log.LogInfof("Redis Status: %s, Stats: Version: %s, Mode: %s, Connected Clients: %s, Used Memory: %s, Peak Used Memory: %s, Uptime: %s seconds",
					health["redis_message"], health["redis_version"], health["redis_mode"], health["redis_connected_clients"],
					health["redis_used_memory"], health["redis_used_memory_peak"], health["redis_uptime_in_seconds"])
			} else {
				redisErrorMessage, _ := health["redis_error"]
				log.LogErrorf("Redis Error: %v", redisErrorMessage)
			}
		}

		// Return the detailed health statistics as JSON
		return c.JSON(health)
	}
}
