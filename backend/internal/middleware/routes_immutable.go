// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

//go:build immutable
// +build immutable

package middleware

import (
	"h0llyw00dz-template/backend/internal/database"
	"h0llyw00dz-template/backend/internal/middleware/router/domain"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes sets up the API routing for the application.
// It organizes routes into versioned groups for better API version management.
//
// Note: There are now 3 routers: restapis, frontend, and wildcard handler (503) (wildcard handler (503) known as root).
// They operate independently. Also note that as the codebase grows, the routing structure
// may become a binary tree (see https://en.wikipedia.org/wiki/Binary_tree), which is considered one of the best art in Go programming.
func RegisterRoutes(app *fiber.App, appName, monitorPath string, db database.Service) {
	// Hosts
	hosts := map[string]*fiber.App{}

	// restAPIs subdomain
	//
	// TODO: Implement a server startup message mechanism similar to "Fiber" ASCII art,
	// with animation (e.g., similar to a streaming/bubble tea spinner) for multiple sites or large codebases.
	// The current static "Fiber" ASCII art only shows one site when there are multiple, which isn't ideal.
	// However, animated ASCII art may not be necessary right now, as it only works properly in terminals.
	api := fiber.New(fiber.Config{
		ServerHeader: appName,
		AppName:      appName,
		// Note: Using the sonic JSON encoder/decoder provides better performance and is more memory-efficient
		// since Fiber is designed for zero allocation memory usage.
		JSONEncoder:      sonic.Marshal,
		JSONDecoder:      sonic.Unmarshal,
		CaseSensitive:    true,
		StrictRouting:    true,
		DisableKeepalive: false,
		ReadTimeout:      readTimeout,
		WriteTimeout:     writeTimeout,
		// Note: It's important to set Prefork to false because if it's enabled and running in Kubernetes,
		// it may get killed by an Out-of-Memory (OOM) error due to a conflict with the Horizontal Pod Autoscaler (HPA).
		Prefork: false,
		// Which is suitable for streaming AI Response.
		StreamRequestBody:       true,
		EnableIPValidation:      true,
		EnableTrustedProxyCheck: true,
		// By default, it is set to 0.0.0.0/0 for local development; however, it can be bound to an ingress controller/proxy.
		// This can be a private IP range (e.g., 10.0.0.0/8).
		TrustedProxies: []string{"0.0.0.0/0"},
		// Trust X-Forwarded-For headers; additionally, this can be customized if using an ingress controller/proxy, especially Ingress Nginx.
		ProxyHeader: fiber.HeaderXForwardedFor, // Fix where * (wildcard header) doesn't work in some kubernetes ingress eco-system
		// Using this immutable setting is more efficient + cheap than the standard library's unique (new pkg).
		Immutable: true,
	})
	registerRESTAPIsRoutes(api, db)
	hosts[apiSubdomain] = api

	// Frontend domain
	//
	// TODO: Implement a server startup message mechanism similar to "Fiber" ASCII art,
	// with animation (e.g., similar to a streaming/bubble tea spinner) for multiple sites or large codebases.
	// The current static "Fiber" ASCII art only shows one site when there are multiple, which isn't ideal.
	// However, animated ASCII art may not be necessary right now, as it only works properly in terminals.
	frontend := fiber.New(fiber.Config{
		ServerHeader: appName,
		AppName:      appName,
		// Note: Using the sonic JSON encoder/decoder provides better performance and is more memory-efficient
		// since Fiber is designed for zero allocation memory usage.
		JSONEncoder:      sonic.Marshal,
		JSONDecoder:      sonic.Unmarshal,
		CaseSensitive:    true,
		StrictRouting:    true,
		DisableKeepalive: false,
		ReadTimeout:      readTimeout,
		WriteTimeout:     writeTimeout,
		// Note: It's important to set Prefork to false because if it's enabled and running in Kubernetes,
		// it may get killed by an Out-of-Memory (OOM) error due to a conflict with the Horizontal Pod Autoscaler (HPA).
		Prefork: false,
		// Which is suitable for streaming AI Response.
		StreamRequestBody:       true,
		EnableIPValidation:      true,
		EnableTrustedProxyCheck: true,
		// By default, it is set to 0.0.0.0/0 for local development; however, it can be bound to an ingress controller/proxy.
		// This can be a private IP range (e.g., 10.0.0.0/8).
		TrustedProxies: []string{"0.0.0.0/0"},
		// Trust X-Forwarded-For headers; additionally, this can be customized if using an ingress controller/proxy, especially Ingress Nginx.
		ProxyHeader: fiber.HeaderXForwardedFor, // Fix where * (wildcard header) doesn't work in some kubernetes ingress eco-system
		// Using this immutable setting is more efficient + cheap than the standard library's unique (new pkg).
		Immutable: true,
	})
	registerStaticFrontendRoutes(frontend, appName, db)
	hosts[frontendDomain] = frontend

	// Apply the combined middlewares
	registerRouteConfigMiddleware(app, db)
	registerRootRouter(app)

	// Configure the DomainRouter middleware
	domainRouter := domain.New(domain.Config{
		Hosts: hosts,
	})

	// Apply the subdomain & domain routing middleware
	app.Use(domainRouter)
}
