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
	//	w.Header().Set("Content-Type", MIMETextEventStream)
	//	w.Write([]byte("event: update\ndata: {\"message\": \"Hello, world!\"}\n\n"))
	//	w.Flush()
	//
	// See the Server-Sent Events specification for more details:
	// https://html.spec.whatwg.org/multipage/server-sent-events.html
	//
	// Note: This suitable for AI
	TextEventStream = "text/event-stream"
)
