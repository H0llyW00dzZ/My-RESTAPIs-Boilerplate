// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package server

import (
	"h0llyw00dz-template/backend/internal/server/k8s/metrics"

	"github.com/gofiber/fiber/v2"
)

// MountPrometheusMiddlewareHandler returns a Fiber handler that sets up the Prometheus middleware.
//
// This function provides a flexible way to mount the Prometheus middleware as a Fiber handler.
// It allows you to specify the path where the Prometheus metrics will be exposed and the service name
// for labeling the metrics. You can also provide additional options for configuring the middleware.
//
// The returned handler function mounts the Prometheus middleware at the specified path and adds it to
// the application's middleware stack. It can be used in scenarios where you want to conditionally mount
// the Prometheus middleware based on certain conditions or routes.
//
// Note: This approach is different from "RegisterRoutesPrometheus" as it returns a Fiber handler
// instead of directly registering the middleware on the FiberServer instance.
//
// Parameters:
//
//	path: The path where the Prometheus metrics will be exposed.
//	serviceName: The name of the service for labeling the metrics.
//	options: Optional parameters for configuring the Prometheus middleware.
//
// Returns:
//
//	A Fiber handler function that mounts the Prometheus middleware.
func MountPrometheusMiddlewareHandler(path, serviceName string, options ...any) fiber.Handler {
	prometheus := metrics.NewPrometheusMiddleware(serviceName, options...)
	return func(c *fiber.Ctx) error {
		prometheus.RegisterAt(c.App(), path)
		c.App().Use(prometheus.Middleware)
		return c.Next()
	}
}
