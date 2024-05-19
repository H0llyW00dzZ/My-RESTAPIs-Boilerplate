// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package server

import (
	"h0llyw00dz-template/backend/internal/server/k8s/metrics"

	"github.com/gofiber/fiber/v2"
)

// MountPrometheusMiddlewareHandler returns a Fiber handler that sets up the Prometheus middleware.
func MountPrometheusMiddlewareHandler(path, serviceName string, options ...interface{}) fiber.Handler {
	prometheus := metrics.NewPrometheusMiddleware(serviceName, options...)
	return func(c *fiber.Ctx) error {
		prometheus.RegisterAt(c.App(), path)
		c.App().Use(prometheus.Middleware)
		return c.Next()
	}
}
