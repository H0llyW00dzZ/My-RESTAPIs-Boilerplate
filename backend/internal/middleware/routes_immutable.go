// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

//go:build immutable
// +build immutable

package middleware

import (
	log "h0llyw00dz-template/backend/internal/logger"
	"h0llyw00dz-template/backend/pkg/convert"
	"h0llyw00dz-template/backend/pkg/network/cidr"

	"h0llyw00dz-template/backend/internal/database"
	"h0llyw00dz-template/backend/internal/middleware/router/domain"
	"h0llyw00dz-template/env"

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
	// Validate and parse trusted proxies
	trustedProxies, err := cidr.ValidateAndParseIPs(env.TRUSTEDPROXIES, "0.0.0.0/0")
	if err != nil {
		log.LogFatal(err)
	}

	// Validate and parse size body limit
	sizeBodyLimit, err := convert.ToBytes(env.GetEnv(env.SIZEBODYLIMIT, "4MiB"))
	if err != nil {
		log.LogFatal(err)
	}

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
		TrustedProxies: trustedProxies,
		// Trust X-Forwarded-For headers. This can be customized if using an ingress controller or proxy, especially Ingress NGINX.
		//
		// Note: X-Forwarded-* or any * (wildcard header) from a reverse proxy don't work with Kubernetes Ingress NGINX.
		// It's better to explicitly use X-Forwarded-For or other specific headers without * (wildcard header).
		ProxyHeader: fiber.HeaderXForwardedFor, // Fix where * (wildcard header) doesn't work in some kubernetes ingress eco-system
		// This immutable setting is more efficient and cost-effective than the standard library's new package.
		// It is also safe to use in combination with the worker package for concurrency.
		Immutable: true,
		// Note: The body limit should be adjusted based on the application's requirements.
		// For optimal performance in a concurrent environment, ensure that the body limit is set appropriately.
		// When the concurrency configuration is well-tuned and matches the Horizontal Pod Autoscaler (HPA) settings in Kubernetes,
		// it can result in a highly stable and scalable system for large-scale deployments (as demonstrated through extensive testing with multiple nodes until stability was consistently achieved).
		BodyLimit: sizeBodyLimit,
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
		TrustedProxies: trustedProxies,
		// Trust X-Forwarded-For headers. This can be customized if using an ingress controller or proxy, especially Ingress NGINX.
		//
		// Note: X-Forwarded-* or any * (wildcard header) from a reverse proxy don't work with Kubernetes Ingress NGINX.
		// It's better to explicitly use X-Forwarded-For or other specific headers without * (wildcard header).
		ProxyHeader: fiber.HeaderXForwardedFor, // Fix where * (wildcard header) doesn't work in some kubernetes ingress eco-system
		// This immutable setting is more efficient and cost-effective than the standard library's new package.
		// It is also safe to use in combination with the worker package for concurrency.
		Immutable: true,
		// Note: The body limit should be adjusted based on the application's requirements.
		// For optimal performance in a concurrent environment, ensure that the body limit is set appropriately.
		// When the concurrency configuration is well-tuned and matches the Horizontal Pod Autoscaler (HPA) settings in Kubernetes,
		// it can result in a highly stable and scalable system for large-scale deployments (as demonstrated through extensive testing with multiple nodes until stability was consistently achieved).
		BodyLimit: sizeBodyLimit,
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
