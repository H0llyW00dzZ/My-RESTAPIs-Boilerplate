// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package domain

import (
	"github.com/gofiber/fiber/v2"
)

// Note: This method works well Docs: https://github.com/gofiber/fiber/issues/750
// Also note that There is no limit to this feature. For example, you can add a billion domains or subdomains.
// Another note: When running this in a container with Kubernetes, make sure to have a configuration for allow internal IPs (e.g., 10.0.0.0/24).
// Because this method creates an additional internal IP for handling routes (e.g., 10.0.0.1 for REST APIs, then 10.0.0.2 for the frontend).
type (
	// Config represents the configuration for the DomainRouter middleware.
	Config struct {
		// Hosts is a map of subdomain or domain hosts to their corresponding Fiber application instances.
		Hosts map[string]*fiber.App
	}
)
