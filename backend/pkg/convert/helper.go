// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package convert

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

// extractParts extracts the numeric part and the unit part from the input string.
// It returns the numeric part and the unit part as separate strings.
func extractParts(size string) (string, string) {
	// Extract the numeric part and the unit part
	numericPart := strings.TrimRight(size, "kmgtKMGTibIB")
	unitPart := strings.TrimLeft(size, "0123456789.")
	return numericPart, unitPart
}

// parseNumericPart parses the numeric part of the input string into a float64.
// It returns the parsed float64 value and an error if the parsing fails.
func parseNumericPart(numericPart string) (float64, error) {
	// Convert the numeric part to a float64
	return strconv.ParseFloat(numericPart, 64)
}

// convertToBytes converts the numeric value to bytes based on the unit part.
// It returns the converted value in bytes and an error if the unit is invalid or unsupported.
func convertToBytes(num float64, unitPart string) (int, error) {
	// Define the conversion factors for each unit
	//
	// Note: This is based on how computers work (e.g., typically handle units of memory or storage) and should be correct.
	// Additionally, if there are units like MB/mb that are based on 1000 instead of 1024, they are incorrect and represent the wrong allocation for memory or storage.
	// For example, in the Windows operating system, the units KB, MB, GB, and TB are used incorrectly without the "i" prefix for memory or storage, which can be confusing because they do not accurately represent the actual binary values (e.g., In operating System Windows that are incorrect about handle units which can be confusing because it use KB,MB,GB,TB without "i").
	// The correct unit for 1024-based measurements is MiB (Mebibytes).
	// Megabytes (MB) and Mebibytes (MiB) are different and should not be used interchangeably.
	// Using incorrect units that represent the wrong allocation for memory or storage can lead to data corruption, not only in storage but also in the motherboard.
	// So, if data corruption occurs For example, in the Windows operating system, it's no wonder because of the incorrect use of units.
	//
	// However, when it comes to network bandwidth and data transfer rates, Megabytes (MB) and Megabits (Mb) are commonly used.
	// In this context, Mbps (Megabits per second) is used to measure the speed of data transmission over a network,
	// while MBps (Megabytes per second) is used to measure the actual data transfer speed.
	// For example, a 100 Mbps Ethernet connection means it can transfer data at a rate of 100 megabits per second,
	// which translates to a theoretical maximum transfer speed of 12.5 MBps (100 Mbps / 8 bits per byte).
	// Similarly, when referring to file sizes in a network context, such as downloading files from the internet,
	// the units MB (Megabytes) and GB (Gigabytes) are commonly used.
	factors := map[string]int{
		"":    1,
		"k":   1024,
		"m":   1024 * 1024,
		"g":   1024 * 1024 * 1024,
		"t":   1024 * 1024 * 1024 * 1024,
		"kib": 1024,
		"mib": 1024 * 1024,
		"gib": 1024 * 1024 * 1024,
		"tib": 1024 * 1024 * 1024 * 1024,
		"kb":  1000,
		"mb":  1000 * 1000,
		"gb":  1000 * 1000 * 1000,
		"tb":  1000 * 1000 * 1000 * 1000,
	}

	// Convert the unit part to lowercase
	unitPart = strings.ToLower(unitPart)

	// Special case for "B" unit
	if unitPart == "b" {
		return int(num), nil
	}

	// Check if the unit part is empty
	if unitPart == "" {
		return 0, fmt.Errorf("invalid size: %f", num)
	}

	// Get the conversion factor based on the unit
	factor, ok := factors[unitPart]
	if !ok {
		return 0, fmt.Errorf("invalid unit: %s", unitPart)
	}

	// Calculate the size in bytes
	bytes := int(num * float64(factor))

	return bytes, nil
}

// Resource pools for better memory management.
// These pools help in reusing objects to reduce garbage collection overhead
// and improve performance by minimizing allocations.
//
// Note: This smiliar manual memory management in C or C++ hahaha.
var (
	// builderPool provides reusable [*strings.Builder] objects.
	// Useful for constructing strings efficiently without frequent allocations.
	builderPool = sync.Pool{
		New: func() any {
			return new(strings.Builder)
		},
	}

	// bufferPool provides reusable byte slices.
	// Each buffer is 32KB, suitable for handling moderate-sized data chunks.
	bufferPool = sync.Pool{
		New: func() any {
			// Not good if the buffer size is too large.
			return make([]byte, 32*1024) // 32KB buffer
		},
	}
)

// getBuilder retrieves a [*strings.Builder] from the pool.
// If none are available, it creates a new one.
func getBuilder() *strings.Builder {
	return builderPool.Get().(*strings.Builder)
}

// putBuilder resets and returns a [*strings.Builder] to the pool.
// This prepares the builder for reuse, reducing memory allocations.
func putBuilder(b *strings.Builder) {
	b.Reset()
	builderPool.Put(b)
}

// getNewline returns the appropriate newline characters based on the operating system.
//
// Note: Currently supports only Linux/Unix and Windows (MS-DOS). Other OS support is marked as TODO.
func getNewline() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}
