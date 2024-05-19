// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package server

import "h0llyw00dz-template/backend/internal/server/k8s/metrics"

// RegisterRoutesPrometheus registers the Prometheus middleware at the specified path.
//
// This method is a convenience function defined on the FiberServer struct. It directly registers
// the Prometheus middleware at the specified path and adds it to the application's middleware stack.
//
// It is a straightforward approach if you want to register the Prometheus middleware globally for all routes.
// The method assumes the existence of a Fiber application instance and is called directly on the FiberServer instance.
//
// Parameters:
//    path: The path where the Prometheus metrics will be exposed.
//    serviceName: The name of the service for labeling the metrics.
//    options: Optional parameters for configuring the Prometheus middleware.
func (s *FiberServer) RegisterRoutesPrometheus(path string, serviceName string, options ...interface{}) {
	prometheus := metrics.NewPrometheusMiddleware(serviceName, options...)
	prometheus.RegisterAt(s.app, path)
	s.app.Use(prometheus.Middleware)
}
