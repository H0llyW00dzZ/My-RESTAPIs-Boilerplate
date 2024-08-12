// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package helper

import "github.com/gofiber/fiber/v2"

// ErrorResponse represents the structure of an error response.
type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

const (

	// MIMEApplicationProblemJSON represents the MIME type for problem+json (RFC 7807).
	MIMEApplicationProblemJSON = "application/problem+json"

	// MIMEApplicationProblemJSONCharsetUTF8 represents the MIME type for problem+json with UTF-8 charset (RFC 7807 Enhancement).
	MIMEApplicationProblemJSONCharsetUTF8 = "application/problem+json; charset=utf-8"
)

// SendErrorResponse sends an error response with the specified status code and error message.
func SendErrorResponse(c *fiber.Ctx, statusCode int, errorMessage string) error {
	// Check if the response body contains non-ASCII characters
	if !isASCII(errorMessage) {
		// If non-ASCII characters are present, use the MIME type with charset
		return c.Status(statusCode).JSON(ErrorResponse{
			Code:  statusCode,
			Error: errorMessage,
		}, MIMEApplicationProblemJSONCharsetUTF8)
	}

	// If the response body contains only ASCII characters, use the MIME type without charset
	return c.Status(statusCode).JSON(ErrorResponse{
		Code:  statusCode,
		Error: errorMessage,
	}, MIMEApplicationProblemJSON)
}

// ErrorHandler is the error handling middleware that runs after other middleware.
func ErrorHandler(c *fiber.Ctx) error {
	// Call the next route handler and catch any errors
	err := c.Next()

	// If a crash/panics occurs, return a generic error response
	if err != nil {
		// Note: This error is used to handle crash/panics because other errors are already handled independently.
		return SendErrorResponse(c, fiber.StatusInternalServerError, fiber.ErrInternalServerError.Message)
	}

	// No errors, continue with the next middleware
	return nil
}

// isASCII checks if a string contains only ASCII characters.
func isASCII(s string) bool {
	for _, c := range s {
		if c > 127 {
			return false
		}
	}
	return true
}
