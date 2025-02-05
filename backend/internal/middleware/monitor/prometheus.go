// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package monitor

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
)

// PrometheusConfig defines the config for the Prometheus middleware.
type PrometheusConfig struct {
	// ServiceName is the name of the service being monitored.
	// If provided, it will be used as the prefix for the metrics.
	// Optional. Default is an empty string.
	ServiceName string

	// Namespace is the namespace for the Prometheus metrics.
	// Optional. Default is "senior_golang".
	Namespace string

	// Subsystem is the subsystem for the Prometheus metrics.
	// Optional. Default is an empty string.
	Subsystem string

	// Labels are additional labels to be added to the Prometheus metrics.
	// Optional. Default is an empty map.
	Labels map[string]string

	// SkipPaths is a list of paths that should be skipped by the Prometheus middleware.
	// Optional. Default is an empty slice.
	SkipPaths []string

	// MetricsPath is the path where the Prometheus metrics will be exposed.
	// Optional. Default is "/metrics".
	MetricsPath string

	// Next is a function that defines a custom logic for skipping the Prometheus middleware.
	// The middleware will be skipped if this function returns true.
	// Optional. Default is nil.
	Next func(c *fiber.Ctx) bool
}

// DefaultPrometheusConfig represents the default configuration options for the Prometheus middleware.
var DefaultPrometheusConfig = PrometheusConfig{
	Namespace:   "senior_golang",
	MetricsPath: "/metrics",
}

// NewPrometheus creates a new Prometheus middleware with optional custom configuration options.
//
// Warning: Do not use this in production, especially on Kubernetes, as it does not store metrics on disk or utilize a fiber storage mechanism (https://docs.gofiber.io/storage/).
// The implementation of github.com/ansrivas/fiberprometheus/v2 lacks a storage mechanism,
// which means metrics will be stored directly in memory.
func NewPrometheus(config ...PrometheusConfig) fiber.Handler {
	cfg := DefaultPrometheusConfig

	if len(config) > 0 {
		cfg = config[0]
		if cfg.Namespace == "" {
			cfg.Namespace = DefaultPrometheusConfig.Namespace
		}
		if cfg.MetricsPath == "" {
			cfg.MetricsPath = DefaultPrometheusConfig.MetricsPath
		}
	}

	var prometheus *fiberprometheus.FiberPrometheus

	if cfg.ServiceName != "" {
		prometheus = fiberprometheus.NewWith(cfg.ServiceName, cfg.Namespace, cfg.Subsystem)
	} else {
		prometheus = fiberprometheus.NewWithLabels(cfg.Labels, cfg.Namespace, cfg.Subsystem)
	}

	if len(cfg.SkipPaths) > 0 {
		prometheus.SetSkipPaths(cfg.SkipPaths)
	}

	return func(c *fiber.Ctx) error {
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		prometheus.RegisterAt(c.App(), cfg.MetricsPath)
		return prometheus.Middleware(c)
	}
}
