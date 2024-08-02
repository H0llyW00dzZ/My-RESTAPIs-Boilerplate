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
	return c.Status(statusCode).JSON(ErrorResponse{
		Code:  statusCode,
		Error: errorMessage,
	}, MIMEApplicationProblemJSONCharsetUTF8)
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
