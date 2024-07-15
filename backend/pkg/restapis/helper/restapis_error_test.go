// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package helper_test

import (
	"net/http/httptest"
	"testing"

	"h0llyw00dz-template/backend/pkg/restapis/helper"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func TestSendErrorResponse_BadRequest(t *testing.T) {
	app := fiber.New()
	app.Get("/gopher/test", func(c *fiber.Ctx) error {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request")
	})

	req := httptest.NewRequest("GET", "/gopher/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}

	var errorResponse helper.ErrorResponse
	err = sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&errorResponse)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	expectedErrorCode := fiber.StatusBadRequest
	if errorResponse.Code != expectedErrorCode {
		t.Errorf("Expected error code %d, got %d", expectedErrorCode, errorResponse.Code)
	}

	expectedErrorMessage := "Invalid request"
	if errorResponse.Error != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, errorResponse.Error)
	}
}

func TestSendErrorResponse_Unauthorized(t *testing.T) {
	app := fiber.New()
	app.Get("/gopher/test", func(c *fiber.Ctx) error {
		return helper.SendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized access")
	})

	req := httptest.NewRequest("GET", "/gopher/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
	}

	var errorResponse helper.ErrorResponse
	err = sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&errorResponse)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	expectedErrorCode := fiber.StatusUnauthorized
	if errorResponse.Code != expectedErrorCode {
		t.Errorf("Expected error code %d, got %d", expectedErrorCode, errorResponse.Code)
	}

	expectedErrorMessage := "Unauthorized access"
	if errorResponse.Error != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, errorResponse.Error)
	}
}

func TestSendErrorResponse_Forbidden(t *testing.T) {
	app := fiber.New()
	app.Get("/gopher/test", func(c *fiber.Ctx) error {
		return helper.SendErrorResponse(c, fiber.StatusForbidden, "Forbidden resource")
	})

	req := httptest.NewRequest("GET", "/gopher/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.StatusCode != fiber.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", fiber.StatusForbidden, resp.StatusCode)
	}

	var errorResponse helper.ErrorResponse
	err = sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&errorResponse)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	expectedErrorCode := fiber.StatusForbidden
	if errorResponse.Code != expectedErrorCode {
		t.Errorf("Expected error code %d, got %d", expectedErrorCode, errorResponse.Code)
	}

	expectedErrorMessage := "Forbidden resource"
	if errorResponse.Error != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, errorResponse.Error)
	}
}

func TestSendErrorResponse_NotFound(t *testing.T) {
	app := fiber.New()
	app.Get("/gopher/test", func(c *fiber.Ctx) error {
		return helper.SendErrorResponse(c, fiber.StatusNotFound, "Resource not found")
	})

	req := httptest.NewRequest("GET", "/gopher/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.StatusCode != fiber.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", fiber.StatusNotFound, resp.StatusCode)
	}

	var errorResponse helper.ErrorResponse
	err = sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&errorResponse)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	expectedErrorCode := fiber.StatusNotFound
	if errorResponse.Code != expectedErrorCode {
		t.Errorf("Expected error code %d, got %d", expectedErrorCode, errorResponse.Code)
	}

	expectedErrorMessage := "Resource not found"
	if errorResponse.Error != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, errorResponse.Error)
	}
}

func TestSendErrorResponse_Conflict(t *testing.T) {
	app := fiber.New()
	app.Get("/gopher/test", func(c *fiber.Ctx) error {
		return helper.SendErrorResponse(c, fiber.StatusConflict, "Duplicate resource")
	})

	req := httptest.NewRequest("GET", "/gopher/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.StatusCode != fiber.StatusConflict {
		t.Errorf("Expected status code %d, got %d", fiber.StatusConflict, resp.StatusCode)
	}

	var errorResponse helper.ErrorResponse
	err = sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&errorResponse)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	expectedErrorCode := fiber.StatusConflict
	if errorResponse.Code != expectedErrorCode {
		t.Errorf("Expected error code %d, got %d", expectedErrorCode, errorResponse.Code)
	}

	expectedErrorMessage := "Duplicate resource"
	if errorResponse.Error != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, errorResponse.Error)
	}
}

func TestSendErrorResponse_BadGateway(t *testing.T) {
	app := fiber.New()
	app.Get("/gopher/test", func(c *fiber.Ctx) error {
		return helper.SendErrorResponse(c, fiber.StatusBadGateway, "Bad gateway")
	})

	req := httptest.NewRequest("GET", "/gopher/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadGateway {
		t.Errorf("Expected status code %d, got %d", fiber.StatusBadGateway, resp.StatusCode)
	}

	var errorResponse helper.ErrorResponse
	err = sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&errorResponse)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	expectedErrorCode := fiber.StatusBadGateway
	if errorResponse.Code != expectedErrorCode {
		t.Errorf("Expected error code %d, got %d", expectedErrorCode, errorResponse.Code)
	}

	expectedErrorMessage := "Bad gateway"
	if errorResponse.Error != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, errorResponse.Error)
	}
}

func TestSendErrorResponse_InternalServerError(t *testing.T) {
	app := fiber.New()
	app.Get("/gopher/test", func(c *fiber.Ctx) error {
		return helper.SendErrorResponse(c, fiber.StatusInternalServerError, "Internal server error")
	})

	req := httptest.NewRequest("GET", "/gopher/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", fiber.StatusInternalServerError, resp.StatusCode)
	}

	var errorResponse helper.ErrorResponse
	err = sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&errorResponse)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	expectedErrorCode := fiber.StatusInternalServerError
	if errorResponse.Code != expectedErrorCode {
		t.Errorf("Expected error code %d, got %d", expectedErrorCode, errorResponse.Code)
	}

	expectedErrorMessage := "Internal server error"
	if errorResponse.Error != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, errorResponse.Error)
	}
}

func TestSendErrorResponse_TooManyRequests(t *testing.T) {
	app := fiber.New()
	app.Get("/gopher/test", func(c *fiber.Ctx) error {
		return helper.SendErrorResponse(c, fiber.StatusTooManyRequests, "Too many requests")
	})

	req := httptest.NewRequest("GET", "/gopher/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.StatusCode != fiber.StatusTooManyRequests {
		t.Errorf("Expected status code %d, got %d", fiber.StatusTooManyRequests, resp.StatusCode)
	}

	var errorResponse helper.ErrorResponse
	err = sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&errorResponse)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	expectedErrorCode := fiber.StatusTooManyRequests
	if errorResponse.Code != expectedErrorCode {
		t.Errorf("Expected error code %d, got %d", expectedErrorCode, errorResponse.Code)
	}

	expectedErrorMessage := "Too many requests"
	if errorResponse.Error != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, errorResponse.Error)
	}
}

func TestErrorHandler(t *testing.T) {
	app := fiber.New()

	// Register the ErrorHandler (for handling panic) & Recover middleware
	app.Use(helper.ErrorHandler, recover.New())

	// Create a test route that panics
	app.Get("/gopher/test", func(c *fiber.Ctx) error {
		panic("Test panic")
	})

	req := httptest.NewRequest("GET", "/gopher/test", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", fiber.StatusInternalServerError, resp.StatusCode)
	}

	var errorResponse helper.ErrorResponse
	err = sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&errorResponse)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	expectedErrorCode := fiber.StatusInternalServerError
	if errorResponse.Code != expectedErrorCode {
		t.Errorf("Expected error code %d, got %d", expectedErrorCode, errorResponse.Code)
	}

	expectedErrorMessage := fiber.ErrInternalServerError.Message
	if errorResponse.Error != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, errorResponse.Error)
	}
}
