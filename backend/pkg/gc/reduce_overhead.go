// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package gc

import (
	"github.com/valyala/bytebufferpool"
)

// BufferPool is used for efficient memory reuse in [I/O operations].
//
// Example usage for replacing [I/O operations] like ReadAll/ReadFull with Fiber's custom JSON encoder/decoder:
//
//	// Get a buffer from the pool
//	buf := gc.BufferPool.Get()
//
//	defer func() {
//		buf.Reset()            // Reset the buffer to prevent data leaks
//		gc.BufferPool.Put(buf) // Return the buffer to the pool for reuse
//	}()
//
//	if _, err := buf.ReadFrom(resp.Body); err != nil {
//		return "", fmt.Errorf("error reading response body: %w", err)
//	}
//
//	// Use the decoder from the Fiber app configuration
//	if err := j.c.App().Config().JSONDecoder(buf.Bytes(), &JsonStructPointer); err != nil {
//		return "", fmt.Errorf("error decoding response: %w", err)
//	}
//
// Example usage for rendering HTMX + TEMPL components:
//
//	buf := gc.BufferPool.Get()
//
//	// Use defer to guarantee buffer cleanup (reset and return to the pool)
//	// even if an error occurs during rendering.
//	defer func() {
//		buf.Reset()            // Reset the buffer to prevent data leaks.
//		gc.BufferPool.Put(buf) // Return the buffer to the pool for reuse.
//	}()
//
//	// Render the HTMX component into the byte buffer.
//	if err := component.Render(c.Context(), buf); err != nil {
//		// Handle any rendering errors by returning an internal server error page.
//		return v.renderErrorPage(c, fiber.StatusInternalServerError, "Error rendering component: %v", err)
//	}
//
//	// Convert the byte buffer to a string.
//	renderedHTML := buf.String()
//
//	// Set the appropriate response headers for HTMX.
//	c.Set("HX-Trigger", "update")
//	c.Set("Content-Type", "text/html")
//
//	// Send the rendered HTML as the response.
//	return c.SendString(renderedHTML)
//
// Example usage for efficient file reading:
//
//	// Get a buffer from the pool
//	buf := gc.BufferPool.Get()
//
//	defer func() {
//		buf.Reset()            // Reset the buffer to prevent data leaks
//		gc.BufferPool.Put(buf) // Return the buffer to the pool for reuse
//	}()
//
//	// Open the file for reading
//	file, err := os.Open("example.txt")
//	if err != nil {
//		return "", fmt.Errorf("error opening file: %w", err)
//	}
//	defer file.Close()
//
//	// Read the file contents into the buffer
//	if _, err := buf.ReadFrom(file); err != nil {
//		return "", fmt.Errorf("error reading file: %w", err)
//	}
//
//	// Process the file contents from the buffer
//	processFileContents(buf.Bytes())
//
// Example usage for handling HTTP requests and responses using the standard library [net/http]:
//
//	http.HandleFunc("/example", func(w http.ResponseWriter, r *http.Request) {
//		// Get a buffer from the pool
//		buf := gc.BufferPool.Get()
//
//		defer func() {
//			buf.Reset()            // Reset the buffer to prevent data leaks
//			gc.BufferPool.Put(buf) // Return the buffer to the pool for reuse
//		}()
//
//		// Read request body into the buffer
//		if _, err := buf.ReadFrom(r.Body); err != nil {
//			http.Error(w, "Error reading request body", http.StatusInternalServerError)
//			return
//		}
//
//		// Process the request data
//		processedData := processData(buf.Bytes())
//
//		// Set response headers
//		w.Header().Set("Content-Type", "text/plain")
//
//		// Write the processed data as the response
//		if _, err := w.Write(processedData); err != nil {
//			fmt.Printf("Error writing response: %v\n", err)
//		}
//	})
//
//	http.ListenAndServe(":8080", nil)
//
// Example usage for handling HTTP requests and responses using the standard library [net/http] with [gzip compression]:
//
//	// Server
//	http.HandleFunc("/example", func(w http.ResponseWriter, r *http.Request) {
//		// Get a buffer from the pool
//		buf := gc.BufferPool.Get()
//
//		defer func() {
//			buf.Reset()            // Reset the buffer to prevent data leaks
//			gc.BufferPool.Put(buf) // Return the buffer to the pool for reuse
//		}()
//
//		// Read request body into the buffer
//		if _, err := buf.ReadFrom(r.Body); err != nil {
//			http.Error(w, "Error reading request body", http.StatusInternalServerError)
//			return
//		}
//
//		// Process the request data
//		processedData := processData(buf.Bytes())
//
//		// Set response headers for gzip
//		w.Header().Set("Content-Encoding", "gzip")
//		w.Header().Set("Content-Type", "text/plain")
//
//		// Write the processed data as a gzip-compressed response
//		gzipWriter := gzip.NewWriter(w)
//		defer gzipWriter.Close()
//
//		if _, err := gzipWriter.Write(processedData); err != nil {
//			fmt.Printf("Error writing response: %v\n", err)
//		}
//	})
//
//	http.ListenAndServe(":8080", nil)
//
// Example usage for an HTTP client sending and receiving [gzip-compressed] data:
//
//	func sendGzipRequest(url string, data []byte) ([]byte, error) {
//		// Get a buffer from the pool
//		buf := gc.BufferPool.Get()
//		defer func() {
//			buf.Reset()            // Reset the buffer to prevent data leaks
//			gc.BufferPool.Put(buf) // Return the buffer to the pool for reuse
//		}()
//
//		// Compress the request data
//		gzipWriter := gzip.NewWriter(buf)
//		if _, err := gzipWriter.Write(data); err != nil {
//			return nil, fmt.Errorf("error compressing request data: %w", err)
//		}
//		gzipWriter.Close()
//
//		// Create a new HTTP request
//		req, err := http.NewRequest(http.MethodPost, url, buf)
//		if err != nil {
//			return nil, fmt.Errorf("error creating request: %w", err)
//		}
//		req.Header.Set("Content-Encoding", "gzip")
//
//		// Perform the request
//		client := &http.Client{}
//		resp, err := client.Do(req)
//		if err != nil {
//			return nil, fmt.Errorf("error sending request: %w", err)
//		}
//		defer resp.Body.Close()
//
//		// Read the response
//		respBuf := gc.BufferPool.Get()
//		defer func() {
//			buf.Reset()            	   // Reset the buffer to prevent data leaks
//			gc.BufferPool.Put(respBuf) // Return the buffer to the pool for reuse
//		}()
//
//		switch resp.Header.Get("Content-Encoding") {
//		case "gzip":
//			gzipReader, err := gzip.NewReader(resp.Body)
//			if err != nil {
//				return nil, fmt.Errorf("error creating gzip reader: %w", err)
//			}
//			defer gzipReader.Close()
//			if _, err := respBuf.ReadFrom(gzipReader); err != nil {
//				return nil, fmt.Errorf("error reading compressed response: %w", err)
//			}
//		default:
//			if _, err := respBuf.ReadFrom(resp.Body); err != nil {
//				return nil, fmt.Errorf("error reading response: %w", err)
//			}
//		}
//
//		return respBuf.Bytes(), nil
//	}
//
// Note: These examples demonstrate various [I/O operations], such as JSON responses,
// rendering HTML components, reading files, handling HTTP requests, and utilizing gzip compression.
// Efficient memory usage is achieved by leveraging a [buffer pool], which is especially
// beneficial in high-concurrency environments. For example, using 8 cores while keeping
// memory usage under 100MiB maintains high CPU efficiency with low memory consumption.
//
// [gzip compression]: https://pkg.go.dev/compress/gzip
// [gzip-compressed]: https://pkg.go.dev/compress/gzip
// [buffer pool]: https://pkg.go.dev/github.com/valyala/bytebufferpool
// [I/O operations]: https://en.wikipedia.org/wiki/Input/output
var BufferPool = bytebufferpool.Pool{}
