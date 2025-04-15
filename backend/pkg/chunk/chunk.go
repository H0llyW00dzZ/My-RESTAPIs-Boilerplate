// Copyright (c) 2025 H0llyW00dzZ All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package chunk

// Split splits a slice into chunks of the specified size
func Split[T any](items []T, chunkSize int) [][]T {
	if chunkSize <= 0 {
		return [][]T{items}
	}

	var chunks [][]T
	for i := 0; i < len(items); i += chunkSize {
		end := min(i+chunkSize, len(items))
		chunks = append(chunks, items[i:end])
	}
	return chunks
}
