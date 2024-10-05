// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// Package worker implement goroutine worker just like human being, and it pretty usefull for large go application.
//
// Important: Be cautious when implementing worker goroutines that in jobs.
// Improper implementation can lead to resource exhaustion (e.g., consuming too much memory, smiliar memory leak).
//
// Recommended: Use this worker in Kubernetes, which is suitable for Horizontal Pod Autoscaling (HPA).
//
// # Average usage:
//
//	NAME                     CPU(cores)   MEMORY(bytes)
//	senior-golang-worker-775b64c9b5-52t2v   176m         51Mi
//	senior-golang-worker-775b64c9b5-5rgph   133m         56Mi
//	senior-golang-worker-775b64c9b5-9d8z5   177m         46Mi
//	senior-golang-worker-775b64c9b5-bvkpg   146m         42Mi
//	senior-golang-worker-775b64c9b5-lxrc7   164m         50Mi
//	senior-golang-worker-775b64c9b5-mkc5k   170m         54Mi
//	senior-golang-worker-775b64c9b5-pk9gz   175m         47Mi
//	senior-golang-worker-775b64c9b5-zrcsh   129m         48Mi
//
// Handling GraphQL With Fiber Framework + Anti Memory Leaks/Wasted (e.g., GC (Garbage Collector) Overhead caused by memory):
//
//	NAME                     CPU(cores)   MEMORY(bytes)
//	senior-golang-worker-59d75b6884-2s4d7   145m         21Mi
//	senior-golang-worker-59d75b6884-5nhjn   170m         19Mi
//	senior-golang-worker-59d75b6884-82df8   151m         24Mi
//	senior-golang-worker-59d75b6884-9g7g2   169m         20Mi
//	senior-golang-worker-59d75b6884-hn59z   201m         20Mi
//	senior-golang-worker-59d75b6884-ltfb5   179m         21Mi
//	senior-golang-worker-59d75b6884-ppw4l   164m         19Mi
//	senior-golang-worker-59d75b6884-pt67x   193m         23Mi
//	senior-golang-worker-59d75b6884-qwz5w   156m         22Mi
//	senior-golang-worker-59d75b6884-sfc95   159m         20Mi
//	senior-golang-worker-59d75b6884-tc6br   179m         20Mi
//	senior-golang-worker-59d75b6884-tq27m   183m         19Mi
//	senior-golang-worker-59d75b6884-zmbx8   122m         23Mi
//
// HPA Suitable:
//
//		Name:                                                     senior-golang-worker-hpa
//		Namespace:                                                senior-golang
//		Labels:                                                   <none>
//		Annotations:                                              <none>
//		CreationTimestamp:                                        Sun, 01 Sep 2024 02:40:18 +0700
//		Reference:                                                Deployment/senior-golang-worker
//		Metrics:                                                  ( current / target )
//	  	resource cpu on pods  (as a percentage of request):     62% (157m) / 80%
//	  	resource memory on pods  (as a percentage of request):  8% (22274363076m) / 50%
//		Min replicas:                                             1
//		Max replicas:                                             50
//		Deployment pods:                                          13 current / 13 desired
//		Conditions:
//	  Type            Status  Reason               Message
//	  ----            ------  ------               -------
//	  AbleToScale     True    ScaleDownStabilized  recent recommendations were higher than current one, applying the highest recent recommendation
//	  ScalingActive   True    ValidMetricFound     the HPA was able to successfully calculate a replica count from cpu resource utilization (percentage of request)
//	  ScalingLimited  False   DesiredWithinRange   the desired count is within the acceptable range
//		Events:
//	  Type    Reason             Age    From                       Message
//	  ----    ------             ----   ----                       -------
//	  Normal  SuccessfulRescale  35m    horizontal-pod-autoscaler  New size: 3; reason: cpu resource utilization (percentage of request) above target
//	  Normal  SuccessfulRescale  35m    horizontal-pod-autoscaler  New size: 5; reason: cpu resource utilization (percentage of request) above target
//	  Normal  SuccessfulRescale  22m    horizontal-pod-autoscaler  New size: 10; reason: cpu resource utilization (percentage of request) above target
//	  Normal  SuccessfulRescale  22m    horizontal-pod-autoscaler  New size: 20; reason: cpu resource utilization (percentage of request) above target
//	  Normal  SuccessfulRescale  16m    horizontal-pod-autoscaler  New size: 17; reason: All metrics below target
//	  Normal  SuccessfulRescale  9m51s  horizontal-pod-autoscaler  New size: 15; reason: All metrics below target
//	  Normal  SuccessfulRescale  20s    horizontal-pod-autoscaler  New size: 13; reason: All metrics below target
//
// HPA Only vCPU (Stable) With QoS: Burstable + Anti Memory Leaks/Wasted (e.g., GC (Garbage Collector) Overhead caused by memory),
// handling billions of requests concurrently for long-running processes:
//
//	Name:                                                  senior-golang-worker-hpa
//	Namespace:                                             senior-golang
//	Labels:                                                <none>
//	Annotations:                                           <none>
//	CreationTimestamp:                                     Sun, 01 Sep 2024 02:40:18 +0700
//	Reference:                                             Deployment/senior-golang-worker
//	Metrics:                                               ( current / target )
//	  resource cpu on pods  (as a percentage of request):  67% (167m) / 80%
//	Min replicas:                                          1
//	Max replicas:                                          50
//	Deployment pods:                                       12 current / 12 desired
//	Conditions:
//	  Type            Status  Reason               Message
//	  ----            ------  ------               -------
//	  AbleToScale     True    ScaleDownStabilized  recent recommendations were higher than current one, applying the highest recent recommendation
//	  ScalingActive   True    ValidMetricFound     the HPA was able to successfully calculate a replica count from cpu resource utilization (percentage of request)
//	  ScalingLimited  False   DesiredWithinRange   the desired count is within the acceptable range
//	Events:
//	  Type    Reason             Age                   From                       Message
//	  ----    ------             ----                  ----                       -------
//	  Normal  SuccessfulRescale  56m (x59 over 5d)     horizontal-pod-autoscaler  New size: 11; reason: All metrics below target
//	  Normal  SuccessfulRescale  49m (x19 over 3d3h)   horizontal-pod-autoscaler  New size: 15; reason: cpu resource utilization (percentage of request) above target
//	  Normal  SuccessfulRescale  49m (x11 over 2d14h)  horizontal-pod-autoscaler  New size: 17; reason: cpu resource utilization (percentage of request) above target
//	  Normal  SuccessfulRescale  44m (x13 over 2d14h)  horizontal-pod-autoscaler  New size: 15; reason: All metrics below target
//	  Normal  SuccessfulRescale  41m (x36 over 3d3h)   horizontal-pod-autoscaler  New size: 13; reason: All metrics below target
//	  Normal  SuccessfulRescale  36m (x56 over 5d)     horizontal-pod-autoscaler  New size: 10; reason: All metrics below target
//	  Normal  SuccessfulRescale  29m (x31 over 4d1h)   horizontal-pod-autoscaler  New size: 14; reason: cpu resource utilization (percentage of request) above target
//	  Normal  SuccessfulRescale  24m (x53 over 5d)     horizontal-pod-autoscaler  New size: 12; reason: All metrics below target
//
// Also note that for more efficiency, as this worker mostly consumes CPU, it is recommended to use AMD CPUs (get good get AMD) for the server specification,
// as they perform better than Intel CPUs for this use case.
//
// Important: While using this worker, do not use Prometheus middleware or any metrics (e.g, Heroku Go Metrics, other) that are directly bound to this repository,
// because it can lead to excessive memory consumption (possibly memory leaks) due to the improper implementation of metrics (wrong implementation regarding metrics).
//
// For example, here's how excessive memory consumption (possibly memory leaks) works in Go:
//   - When the GC (Garbage Collector) becomes overloaded (overhead), it will take a lot of time to free memory resources.
//
// For example, here's how metrics can be wrongly implemented:
//   - Metrics should not be stored in memory and then wait for collection, because when store in memory then waiting for collection, the garbage collector will become overhead as goroutines hold the metrics
//     that must be collected by an external process, caller, or whatever it is.
//
// # Compatibility:
//
//   - Due to this worker being designed similar to a semaphore, it is recommended to use it in Kubernetes batch/job services as it can be useful for the cluster
//     (e.g., implementing its own self-healing mechanism, maintaining databases, fully managed by the goroutines).
//
//   - While using this worker, do not use a mutex again in functions that will be executed/managed by goroutines, because it can degrade the performance (making it slower).
//     Instead, use channels for communication. Additionally, to ensure avoiding data races (as mutexes can degrade performance), use the "unique" package (a new package) by copying from the original value. It is safe.
//     However, while using "unique", resource memory usage will be dynamic (unpredictable, e.g., 70MiB then scaling down to 40MiB), and it is suitable for HPA due to its spreading across every pod.
//
// # Boost The Worker:
//
//   - Since this worker is concurrent, it is possible to boost the worker (goroutines). For example, if you are only using this worker for handling Fiber requests concurrently (e.g., from a http client),
//     you should also set the concurrency level in the Fiber configuration accordingly (e.g., "512 * 1024"). This can make both the worker and Fiber faster, as concurrency
//     should be handled by concurrency.
package worker
