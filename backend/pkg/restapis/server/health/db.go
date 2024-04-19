// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package health

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"h0llyw00dz-template/backend/internal/database"
)

// Response represents the structured response for the health statistics.
type Response struct {
	MySQLHealth *MySQLHealth `json:"mysql_health,omitempty"`
	RedisHealth *RedisHealth `json:"redis_health,omitempty"`
}

// DBHandler is a Fiber handler that checks the health of the database and Redis.
// It logs the user activity and the health status of MySQL and Redis.
// The detailed health statistics are returned as a structured JSON response.
func DBHandler(db database.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the IP address from the request context
		ipAddress := c.IP()

		// Initialize the valid filters slice using the fiber.Storage and IP address
		initValidFiltersSlice(db.FiberStorage(), ipAddress)
		// Get the filter parameter from the query string
		filter := c.Query("filter")

		// Check if the filter is valid
		if !isValidFilter(filter) {
			// TODO: Deal with log errors. Typically, I wouldn't tackle this for StatusBadRequest or StatusNotFound. 🤷‍♂️ 🤪
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid filter parameter. Allowed values: %s", strings.Join(getValidFilters(), ", ")),
			})
		}

		// Log the user activity based on the filter
		logUserActivity(c, filter)

		// Get the health status from the database service
		health := db.Health(filter)

		// Create the response struct using the createHealthResponse function
		response := createHealthResponse(health, filter)

		// Log the health status based on the filter
		logHealthStatus(response, filter)

		// Return the structured health statistics as JSON
		// Note: The "c.JSON" method uses the sonic package (related to main configuration) for JSON encoding and decoding,
		// which is one of the reasons why the Fiber framework is considered the best framework in 2024.
		// "You don't need to repeat yourself for JSON encoding/decoding (e.g., using the standard library or other JSON encoder/decoder)."
		return c.JSON(response)
	}
}

// getValidFilters returns a slice of valid filter values.
func getValidFilters() []string {
	// Note: this kind of optimization is part of a strategy called "premature optimization"
	return validFiltersSlice
}

// createHealthResponse creates a Response struct from the provided health statistics.
func createHealthResponse(health map[string]string, filter string) Response {
	// Note: By structuring the code this way, it is easily maintainable for customization, etc.
	// Define a map of filter-specific response creation functions
	responseCreators := map[string]func(map[string]string) interface{}{
		"mysql": func(h map[string]string) interface{} {
			return createMySQLHealthResponse(h)
		},
		"redis": func(h map[string]string) interface{} {
			return createRedisHealthResponse(h)
		},
	}

	response := Response{}

	// Check if the filter is empty or exists in the responseCreators map
	if filter == "" {
		// If the filter is empty, create responses for all available filters
		for _, creator := range responseCreators {
			applyHealthResponse(&response, creator(health))
		}
	} else if creator, ok := responseCreators[filter]; ok {
		// If the filter exists in the responseCreators map, create the corresponding response
		applyHealthResponse(&response, creator(health))
	}

	return response
}

// applyHealthResponse applies the health response to the Response struct based on the type of response.
func applyHealthResponse(response *Response, healthResponse interface{}) {
	switch r := healthResponse.(type) {
	case *MySQLHealth:
		response.MySQLHealth = r
	case *RedisHealth:
		response.RedisHealth = r
	}
}
