// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package helper

import (
	"h0llyw00dz-template/backend/pkg/mime"

	"github.com/gofiber/fiber/v2"
)

// ErrorResponse represents the structure of an error response.
type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

// SendErrorResponse sends an error response with the specified status code and error message.
func SendErrorResponse(c *fiber.Ctx, statusCode int, errorMessage string) error {
	// Check if the response body contains non-ASCII characters
	if !isASCII(errorMessage) {
		// If non-ASCII characters are present, use the MIME type with charset
		return c.Status(statusCode).JSON(ErrorResponse{
			Code:  statusCode,
			Error: errorMessage,
		}, mime.ApplicationProblemJSONCharsetUTF8)
	}

	// If the response body contains only ASCII characters, use the MIME type without charset
	return c.Status(statusCode).JSON(ErrorResponse{
		Code:  statusCode,
		Error: errorMessage,
	}, mime.ApplicationProblemJSON)
}

// ErrorHandler is the error handling middleware that runs after other middleware.
//
// TODO: Deprecate/Remove This Function - it will be replaced by [htmx.NewErrorHandler] when the [middleware.RegisterRoutes] is reorganized.
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
