// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package csp

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"unique"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

// Config represents the configuration options for the CSP middleware.
type Config struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// RandomnessGenerator is a function that generates a random string for the CSP nonce.
	// It takes a string parameter as input, which can be used as a seed or additional data for randomness generation.
	// The function should return the generated random string.
	//
	// Optional. Default: A default randomness generator function provided by the CSP middleware.
	RandomnessGenerator func(string) string

	// ContextKey is the key used to store and retrieve the CSP data from the request context.
	// It should be a unique and descriptive key to avoid conflicts with other middleware or application data.
	//
	// Optional. Default: "csp_random"
	ContextKey any

	// CSPValueGenerator is a function that generates the CSP header value based on the provided randomness and custom values.
	// It takes the generated randomness and a map of custom values as input and returns the desired CSP header value.
	//
	// Optional. Default: A default CSP value generator function provided by the CSP middleware.
	CSPValueGenerator func(string, map[string]string) string

	// IPHeader is the header name used to retrieve the client IP address.
	// If not provided, it will use the default "X-Real-IP" header.
	//
	// Optional. Default: "X-Real-IP"
	IPHeader string
}

// DefaultConfig returns the default configuration for the CSP middleware.
func DefaultConfig() Config {
	return Config{
		RandomnessGenerator: defaultRandomnessGenerator,
		ContextKey:          "csp_random",
		Next:                nil,
		CSPValueGenerator:   defaultCSPValueGenerator,
		IPHeader:            "X-Real-IP",
	}
}

// defaultRandomnessGenerator generates randomness using SHA256 (digest) of the client IP.
func defaultRandomnessGenerator(clientIP string) string {
	hash := sha256.Sum256([]byte(clientIP))
	return hex.EncodeToString(hash[:])
}

// defaultCSPValueGenerator generates the default CSP header value using the provided randomness and custom values.
// It sets the 'script-src' directives with a nonce value based on the randomness.
func defaultCSPValueGenerator(randomness string, customValues map[string]string) string {
	var cspBuilder strings.Builder

	// Add script-src directive
	cspBuilder.WriteString(fmt.Sprintf("script-src 'nonce-%s'", randomness))

	return cspBuilder.String()
}

// getClientIP retrieves the client IP address from the specified header or the remote address.
//
// It handles cases where the header contains multiple IP addresses separated by commas.
//
// Important: The real client IP address must be the first one in the list. Other IP addresses in the list are typically from proxies or load balancers.
// If the real client IP address is not the first one, it indicates that other routers/ingresses are not following best practices (bad practices) for IP address forwarding.
func getClientIP(c *fiber.Ctx, ipHeader string) []string {
	clientIP := unique.Make(c.Get(ipHeader))
	if clientIP.Value() == "" {
		return []string{c.IP()}
	}

	var validIPs []string

	// Split the header value by comma to get multiple IP addresses
	ipList := strings.Split(clientIP.Value(), ",")

	// Iterate over the IP addresses and store the valid ones
	for _, ip := range ipList {
		ip = strings.TrimSpace(ip) // Trim leading/trailing whitespace

		// Check if the IP address is a valid IPv4 address
		if utils.IsIPv4(ip) {
			validIPs = append(validIPs, ip)
			continue
		}

		// Check if the IP address is a valid IPv6 address
		if utils.IsIPv6(ip) {
			validIPs = append(validIPs, ip)
		}
	}

	// If the IP address is not valid, return [c.IP] anyway to prevent Header Spoofing.
	// This will use the Private IP or Real Client IP Address, which could be random (depending on the server configuration),
	// making it difficult to guess for bypass or any potential vulnerable purposes.
	if len(validIPs) == 0 {
		return []string{c.IP()}
	}

	return validIPs
}
