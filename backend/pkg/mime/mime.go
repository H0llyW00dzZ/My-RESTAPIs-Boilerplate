// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package mime

const (
	// ApplicationProblemJSON represents the MIME type for problem+json (RFC 7807).
	// It is used to indicate that the response body contains a problem details object
	// serialized using the JSON format.
	ApplicationProblemJSON = "application/problem+json"

	// ApplicationProblemJSONCharsetUTF8 represents the MIME type for problem+json
	// with UTF-8 charset (RFC 7807 Enhancement).
	ApplicationProblemJSONCharsetUTF8 = "application/problem+json; charset=utf-8"
)

const (
	// TextEventStream represents the MIME type for text/event-stream.
	//
	// It is used to indicate that the response body contains a stream of server-sent events (SSE).
	// Server-sent events allow the server to push updates to the client in real-time over a
	// long-lived HTTP connection.
	//
	// The event stream consists of a series of event messages, where each message is a single
	// line of text terminated by a newline character. The message format includes fields such
	// as "event", "data", "id", and "retry" to convey event information.
	//
	// Example usage:
	//
	//	w.Header().Set("Content-Type", mime.TextEventStream)
	//	w.Write([]byte("event: update\ndata: {\"message\": \"Hello, world!\"}\n\n"))
	//	w.Flush()
	//
	// See the Server-Sent Events specification for more details:
	// https://html.spec.whatwg.org/multipage/server-sent-events.html
	//
	// Note: This MIME type is suitable for AI applications, such as chat systems.
	// It's important to note that Server-Sent Events (SSE) is different from WebSocket.
	// While WebSocket enables bidirectional communication and can be more complex to implement securely,
	// SSE provides a simpler, unidirectional communication channel from the server to the client,
	// making it a suitable choice for scenarios where real-time updates are needed without the
	// additional complexity and potential security risks associated with WebSocket.
	TextEventStream = "text/event-stream"
)

const (
	// ImageXIcon represents the MIME type for image/x-icon.
	//
	// Note: This MIME type is suitable for use with the CustomNextContentType option in the CompressMiddleware (brotli).
	// The CompressMiddleware (brotli) is typically used for compressing JSON and TextEventStream responses (Excellent Performance especially in http/3).
	// However, the "image/x-icon" MIME type can also be added to the CustomNextContentType option to enable
	// compression for ICO files when using the CompressMiddleware (brotli).
	ImageXIcon = "image/x-icon"
)
