// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package gc

import "net/http"

// ReadResponseBody reads the response body into a string and returns it.
// It uses the bufferPool to minimize memory allocations.
func ReadResponseBody(resp *http.Response) (string, error) {
	// Get a buffer from the pool.
	buf := BufferPool.Get()

	// Use defer to guarantee buffer cleanup (reset and return to the pool)
	// even if an error occurs during reading the response body.
	defer func() {
		buf.Reset()         // Reset the buffer to prevent data leaks.
		BufferPool.Put(buf) // Return the buffer to the pool for reuse.
	}()

	// Copy the response body into the buffer.
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return "", err
	}

	// Convert the buffer's contents into a string.
	return buf.String(), nil
}
