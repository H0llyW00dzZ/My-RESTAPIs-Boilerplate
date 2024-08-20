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
// Resource Usage:
//
// It's important to consider the resource usage of the Prometheus middleware, especially in terms of memory consumption.
// Collecting and exposing metrics can potentially consume a significant amount of memory, depending on the number of requests
// and the granularity of the metrics being collected. This is the reason why I rarely use other services
// in Kubernetes environments (cloud ecosystem) (e.g., control plane/node) that can consume a significant amount of memory, as they are not commonly used
// by the applications (this repo). Then When reviewing the bills, it becomes evident that these services, which are mostly unused by the apps (this repo),
// contribute to the high expenses.
//
// The Prometheus middleware uses the Prometheus client library to collect and expose metrics. Each metric consumes memory to
// store its value, labels, and other metadata. As the number of unique metric combinations increases, the memory usage can grow
// accordingly.
//
// To mitigate excessive memory consumption, consider the following strategies:
//
//   - Use appropriate metric types: Choose the appropriate Prometheus metric types based on the requirements. For example, use
//     Gauges for values that can go up and down, Counters for monotonically increasing values, and Histograms or Summaries for
//     measuring distributions of values.
//
//   - Limit the cardinality of labels: Be cautious when adding labels to metrics, as each unique combination of label values
//     creates a new time series, consuming additional memory. Avoid using high-cardinality labels or unbounded label values.
//
//   - Monitor and adjust scrape intervals: Configure appropriate scrape intervals for Prometheus to collect metrics from the
//     application. Longer scrape intervals can help reduce the frequency of metric collection and conserve memory, but they may
//     also impact the granularity of the collected data.
//
//   - Use Prometheus best practices: Follow Prometheus best practices for metric naming, labeling, and instrumentation to ensure
//     efficient memory usage. Avoid creating excessive numbers of metrics or using overly complex metric names.
//
//   - Profile and optimize: Regularly profile the application's memory usage and identify any memory leaks or excessive memory
//     consumption. Optimize the code and configurations to minimize memory overhead.
//
//   - Scale and distribute: If the memory usage becomes a bottleneck, consider scaling the application horizontally by distributing
//     the workload across multiple instances. This can help spread the memory consumption across different nodes.
//
// By being mindful of the memory usage and applying appropriate strategies, the resource consumption of the Prometheus middleware
// can be effectively managed while still benefiting from its powerful monitoring capabilities.
// It's worth noting that in this repository, without the Prometheus middleware (i.e., the original version),
// the average memory usage ranges from 10 MiB to 50 MiB. So, while there may be instances of bottlenecks or out-of-memory (OOM) issues,
// they are not caused by this repository (better blame cloud ecosystem).
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
//
// Additionally, the Prometheus middleware can be enhanced with an authentication mechanism to restrict access
// to the metrics endpoint. This can be achieved by implementing a custom authentication middleware or using
// existing authentication libraries compatible with Fiber.
//
// Recommended for authentication:
//   - https://docs.gofiber.io/api/middleware/keyauth + Database (bound keys mechanism (e.g., Tokens, API Keys, or other keys built with cryptography)
//     to a database instead of hardcoding them, allowing for easy key generation and management)
//
// Recommended for authentication (Advanced):
//   - https://docs.gofiber.io/api/middleware/keyauth + HSM (bound keys mechanism (e.g., Tokens, API Keys, or other keys built with cryptography)
//     to a HSM instead of hardcoding them, allowing for easy key generation and management)
//
// By implementing authentication, you can ensure that only authorized users or systems can access the Prometheus
// metrics endpoint, enhancing the security of your application.
package monitor
