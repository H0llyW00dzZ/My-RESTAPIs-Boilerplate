// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package metrics

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/prometheus/client_golang/prometheus"
)

// NewPrometheusMiddleware creates a new Prometheus middleware with optional custom configuration options.
func NewPrometheusMiddleware(serviceName string, options ...any) *fiberprometheus.FiberPrometheus {
	var registry *prometheus.Registry
	var namespace, subsystem string
	var labels map[string]string

	// Extract namespace, subsystem, labels, and registry from the options.
	for _, option := range options {
		switch opt := option.(type) {
		case *prometheus.Registry:
			registry = opt
		case string:
			if namespace == "" {
				namespace = opt
			} else if subsystem == "" {
				subsystem = opt
			}
		case map[string]string:
			labels = opt
		}
	}

	// Create a new Prometheus instance based on the provided options.
	var prometheus *fiberprometheus.FiberPrometheus
	if registry != nil {
		prometheus = fiberprometheus.NewWithRegistry(registry, serviceName, namespace, subsystem, labels)
	} else if labels != nil {
		prometheus = fiberprometheus.NewWithLabels(labels, namespace, subsystem)
	} else {
		prometheus = fiberprometheus.NewWith(serviceName, namespace, subsystem)
	}

	// Return the Prometheus middleware.
	return prometheus
}
