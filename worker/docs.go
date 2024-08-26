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
// Also note that for more efficiency, as this worker mostly consumes CPU, it is recommended to use AMD CPUs (get good get AMD) for the server specification,
// as they perform better than Intel CPUs for this use case.
package worker
