// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package helper

import "github.com/gofiber/fiber/v2"

// ErrorResponse represents the structure of an error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// SendErrorResponse sends an error response with the specified status code and error message.
func SendErrorResponse(c *fiber.Ctx, statusCode int, errorMessage string) error {
	return c.Status(statusCode).JSON(ErrorResponse{
		Error: errorMessage,
	})
}

// ErrorHandler is the error handling middleware that runs after other middleware.
func ErrorHandler(c *fiber.Ctx) error {
	// Call the next route handler and catch any errors
	err := c.Next()

	// If a crash/panics occurs, return a generic error response
	if err != nil {
		// Note: This error is used to handle crash/panics because other errors are already handled independently.
		return fiber.NewError(fiber.StatusInternalServerError)
	}

	// No errors, continue with the next middleware
	return nil
}
