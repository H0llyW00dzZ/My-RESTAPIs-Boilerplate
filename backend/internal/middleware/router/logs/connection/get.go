// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package connectionlogger

import (
	"fmt"
	"sync/atomic"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// Tracks the number of active connections.
var activeConnections int64

// GetActiveConnections returns the current active connection count.
//
// This effectively counts any active connection (e.g., keep-alive). However, it currently doesn't support Prometheus for creating custom metrics for HPA.
func GetActiveConnections(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
	// Get the current count
	if count := atomic.LoadInt64(&activeConnections); count > 1 {
		return output.WriteString(fmt.Sprintf("%d Active Connections", count))
	} else if count == 1 {
		return output.WriteString("1 Active Connections")
	}

	return output.WriteString("0 Active Connections")
}
