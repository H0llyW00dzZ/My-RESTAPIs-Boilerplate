// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

//go:build immutable
// +build immutable

package main

import (
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

// Config holds the application configuration settings
type Config struct {
	AppName         string
	Port            string
	MonitorPath     string
	TimeFormat      string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// setupFiber initializes a new Fiber application with custom configuration.
// It sets up the JSON encoder/decoder, case sensitivity, and strict routing,
// and applies the application name to the server headers.
func setupFiber(config Config) *fiber.App {
	// TODO: Implement a server startup message mechanism similar to "Fiber" ASCII art,
	// with animation (e.g., similar to a streaming/bubble tea spinner) for multiple sites or large codebases.
	// The current static "Fiber" ASCII art only shows one site when there are multiple, which isn't ideal.
	// However, animated ASCII art may not be necessary right now, as it only works properly in terminals.
	return fiber.New(fiber.Config{
		ServerHeader: config.AppName,
		AppName:      config.AppName,
		// Note: Using the sonic JSON encoder/decoder provides better performance and is more memory-efficient
		// since Fiber is designed for zero allocation memory usage.
		JSONEncoder:      sonic.Marshal,
		JSONDecoder:      sonic.Unmarshal,
		CaseSensitive:    true,
		StrictRouting:    true,
		DisableKeepalive: false,
		ReadTimeout:      config.ReadTimeout,
		WriteTimeout:     config.WriteTimeout,
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
		// Trust X-Forwarded-For headers. This can be customized if using an ingress controller or proxy, especially Ingress NGINX.
		//
		// Note: X-Forwarded-* or any * (wildcard header) from a reverse proxy don't work with Kubernetes Ingress NGINX.
		// It's better to explicitly use X-Forwarded-For or other specific headers without * (wildcard header).
		ProxyHeader: fiber.HeaderXForwardedFor, // Fix where * (wildcard header) doesn't work in some kubernetes ingress eco-system
		// This immutable setting is more efficient and cost-effective than the standard library's new package.
		// It is also safe to use in combination with the worker package for concurrency.
		Immutable: true,
	})
}
