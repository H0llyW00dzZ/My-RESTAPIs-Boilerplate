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
// Average usage:
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
// Handling GraphQL With Fiber Framework + Anti Memory Leaks/Wasted:
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
// Also note that for more efficiency, as this worker mostly consumes CPU, it is recommended to use AMD CPUs (get good get AMD) for the server specification,
// as they perform better than Intel CPUs for this use case.
//
// Important: While using this worker, do not use Prometheus middleware or any metrics that are directly bound to this repository,
// because it can lead to excessive memory consumption (possibly memory leaks) due to the improper implementation of metrics (wrong implementation regarding metrics).
package worker
