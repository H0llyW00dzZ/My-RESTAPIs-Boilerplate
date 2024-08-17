// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package monitor

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
)

// PrometheusConfig represents the configuration options for the Prometheus middleware.
type PrometheusConfig struct {
	ServiceName string
	Namespace   string
	Subsystem   string
	Labels      map[string]string
	SkipPaths   []string
	MetricsPath string
	Next        func(c *fiber.Ctx) bool
}

// DefaultPrometheusConfig represents the default configuration options for the Prometheus middleware.
var DefaultPrometheusConfig = PrometheusConfig{
	Namespace:   "senior_golang",
	MetricsPath: "/metrics",
}

// NewPrometheus creates a new Prometheus middleware with optional custom configuration options.
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
