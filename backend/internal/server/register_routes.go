// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package server

import "h0llyw00dz-template/backend/internal/server/k8s/metrics"

// RegisterRoutesPrometheus registers the Prometheus middleware at the specified path.
func (s *FiberServer) RegisterRoutesPrometheus(path string, serviceName string, options ...interface{}) {
	prometheus := metrics.NewPrometheusMiddleware(serviceName, options...)
	prometheus.RegisterAt(s.app, path)
	s.app.Use(prometheus.Middleware)
}
