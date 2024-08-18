// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// Package monitor provides a Prometheus middleware for Fiber applications.
//
// The Prometheus middleware is built on top of the Fiber framework and collects and exposes
// metrics about the HTTP requests handled by the Fiber application. It integrates seamlessly
// with the Prometheus monitoring system to provide insights into the performance and behavior
// of the application, without relying on any specific cloud environment.
//
// Configuration:
//
// The Prometheus middleware can be configured using the PrometheusConfig struct.
// The available configuration options are:
//   - ServiceName (string): The name of the service being monitored. If provided, it will be used as the prefix for the metric names.
//   - Namespace (string): The namespace for the Prometheus metrics. Default is "senior_golang".
//   - Subsystem (string): The subsystem for the Prometheus metrics. Default is an empty string.
//   - Labels (map[string]string): Additional labels to be added to the Prometheus metrics.
//   - SkipPaths ([]string): A list of paths to be skipped from being monitored by the Prometheus middleware.
//   - MetricsPath (string): The path at which the Prometheus metrics will be exposed. Default is "/metrics".
//   - Next (func(c *fiber.Ctx) bool): A custom function to determine whether the Prometheus middleware should be skipped for a particular request.
//   - CacheKey (string): The cache key used for caching the Prometheus metrics. Default is an empty string.
//
// Default Configuration:
//
// The [DefaultPrometheusConfig] variable provides default values for the Prometheus middleware configuration.
// By default, the Namespace is set to "senior_golang" and the MetricsPath is set to "/metrics".
//
// Example:
//
//	app := fiber.New()
//
//	prometheusConfig := monitor.PrometheusConfig{
//	    ServiceName: "my-service",
//	    Namespace:   "my-namespace",
//	    Subsystem:   "my-subsystem",
//	    Labels: map[string]string{
//	        "env": "production",
//	    },
//	    SkipPaths: []string{"/health", "/metrics"},
//	    MetricsPath: "/my-metrics",
//	    Next: func(c *fiber.Ctx) bool {
//	        return c.Path() == "/skip"
//	    },
//	    CacheKey: "my-cache-key",
//	}
//
//	app.Use(monitor.NewPrometheus(prometheusConfig))
//
// In this example, the Prometheus middleware is configured with a custom service name, namespace, subsystem, labels,
// skip paths, metrics path, a custom Next function, and a cache key. The middleware will collect metrics for all requests except
// those with paths "/health", "/metrics", and "/skip".
//
// The collected metrics will be exposed at the "/my-metrics" endpoint and can be scraped by Prometheus for monitoring purposes.
// The metrics will be cached using the specified cache key "my-cache-key".
//
// The Prometheus middleware is self-contained and does not require any specific cloud environment setup. It can be used in
// any environment where the Fiber application is deployed, making it highly portable and flexible.
//
// Security:
//
// While the Prometheus middleware is highly portable and flexible (unlike relying on any specific cloud environment,
// for example, via localhost), allowing it to be used in any environment where the Fiber application is deployed,
// it focuses on doing one thing and doing it well:
// collecting and exposing metrics for Prometheus.
//
// To ensure security and isolation, the Prometheus middleware can be easily configured to use
// TLS (Transport Layer Security) for encrypted communication. By enabling TLS, the middleware
// can securely expose the metrics endpoint, preventing unauthorized access and protecting sensitive
// information.
package monitor
