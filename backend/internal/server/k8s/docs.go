// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

// Package k8s provides functionality for integrating with Kubernetes.
//
// This package contains subpackages and modules that are specific to Kubernetes-related
// features and configurations. It aims to simplify the integration of the application
// with Kubernetes by providing utilities, handlers, and middleware.
//
// Subpackages:
//
//   - metrics: Contains code related to metrics and monitoring in the context of Kubernetes.
//     It provides middleware and handlers for exposing metrics to Prometheus.
//
// Example usage of the metrics subpackage:
//
//	// Create a new Prometheus middleware with a service name
//	prometheusMiddleware := metrics.NewPrometheusMiddleware("my-service")
//
//	// Create a new Prometheus middleware with a service name and namespace
//	prometheusMiddleware := metrics.NewPrometheusMiddleware("my-service", "my-namespace")
//
//	// Create a new Prometheus middleware with a service name, namespace, and subsystem
//	prometheusMiddleware := metrics.NewPrometheusMiddleware("my-service", "my-namespace", "my-subsystem")
//
//	// Create a new Prometheus middleware with a service name and custom labels
//	prometheusMiddleware := metrics.NewPrometheusMiddleware("my-service", map[string]string{
//		"custom_label1": "custom_value1",
//		"custom_label2": "custom_value2",
//	})
//
//	// Create a new Prometheus middleware with a service name, namespace, subsystem, and custom labels
//	prometheusMiddleware := metrics.NewPrometheusMiddleware("my-service", "my-namespace", "my-subsystem", map[string]string{
//		"custom_label1": "custom_value1",
//		"custom_label2": "custom_value2",
//	})
//
//	// Create a new Prometheus middleware with a custom Prometheus registry
//	customRegistry := prometheus.NewRegistry()
//	prometheusMiddleware := metrics.NewPrometheusMiddleware("my-service", customRegistry)
//
//	// Register the Prometheus middleware at a specific path
//	prometheusMiddleware.RegisterAt(app, "/metrics")
//
//	// Use the Prometheus middleware
//	app.Use(prometheusMiddleware.Middleware)
//
// The metrics subpackage provides the NewPrometheusMiddleware function for creating a new
// Prometheus middleware with optional custom configuration options. It allows creating a middleware
// with a service name, namespace, subsystem, custom labels, and a custom Prometheus registry.
//
// As more subpackages or modules are added to the k8s package, this documentation can be extended
// to provide an overview of their functionality and usage examples.
//
// Note: The k8s package and its subpackages are intended to be used in conjunction with
// a Kubernetes environment. Ensure that the necessary Kubernetes dependencies and
// configurations are set up in the project.
package k8s
