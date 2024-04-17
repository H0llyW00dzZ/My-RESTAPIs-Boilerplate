// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package health

import "strconv"

// bytesToMB converts bytes to megabytes (MB)
func bytesToMB(bytesStr string) float64 {
	bytes, _ := strconv.ParseFloat(bytesStr, 64)
	return bytes / (1024 * 1024)
}
