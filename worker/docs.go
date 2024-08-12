// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// Package worker implement goroutine worker just like human being, and it pretty usefull for large go application.
//
// Important: Be cautious when implementing worker goroutines that in jobs.
// Improper implementation can lead to resource exhaustion (e.g., consuming too much memory, smiliar memory leak).
package worker
