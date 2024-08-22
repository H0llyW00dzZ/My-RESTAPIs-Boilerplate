// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package domain

import (
	"github.com/gofiber/fiber/v2"
)

// New creates a new instance of the DomainRouter middleware with the provided configuration.
//
// Note: This is useful for large Go applications, especially when running in Kubernetes,
// as it eliminates the need for multiple containers. It also supports integration with the Kubernetes ecosystem,
// such as pointing to CNAME/NS or manually (if not using Kubernetes).
// Also note that for TLS certificates, a wildcard/advanced certificate is required.
//
// Known Bugs:
//
//   - Wildcard/advanced certificates (e.g., issued by DigiCert, Sectigo, Google Trust Services, or a private CA) are not supported/compatible on Heroku.
//     Using a wildcard/advanced certificate on Heroku will cause an "SSL certificate error: There is conflicting information between the SSL connection, its certificate, and/or the included HTTP requests."
//     If using a wildcard/advanced certificate, it is recommended to deploy the application in a cloud environment such as Kubernetes, where you can easily control the ingress controller (e.g., implement your own, such as Universe).
//     Also note that regarding known bugs, it is not caused by this repository; it is an issue with Heroku's router.
//
// Example public wildcard CAs that can be used for an ingress or directly:
//
//   - https://crt.sh/?q=a8bc9093e1f4ba202fc769b8818b8a279a5f70c91bee458d29d6ad3c5ac5e88c
func New(config Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		host := config.Hosts[c.Hostname()]
		if host == nil {
			// Note: Returning a new error is a better approach instead of returning directly,
			// as it allows the error to be handled by the caller somewhere else in the codebase,
			// especially when the codebase grows larger.
			return fiber.NewError(fiber.StatusServiceUnavailable)
		}
		// Use c.Context() to pass the underlying context to the host's Fiber app.
		host.Handler()(c.Context())
		return nil
	}
}
