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
	if !mime.IsASCII(errorMessage) {
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
