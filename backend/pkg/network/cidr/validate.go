// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package cidr

import (
	"fmt"
	"h0llyw00dz-template/env"
	"net"
	"strings"
)

// ValidateAndParseIPs checks if the given IPs or CIDR ranges are valid from environment variable
func ValidateAndParseIPs(envVar string, defaultIPS string) ([]string, error) {
	var trustedProxies []string

	// Get the IPs from the environment variable
	ips := env.GetEnv(envVar, defaultIPS)

	// If the default value is used, skip validation
	if ips == defaultIPS {
		return []string{ips}, nil
	}

	// Split the IPs by comma
	ipList := strings.Split(ips, ",")

	for _, ip := range ipList {
		ip = strings.TrimSpace(ip)
		// Check if it's a valid IP or CIDR
		if net.ParseIP(ip) != nil || isValidCIDR(ip) {
			trustedProxies = append(trustedProxies, ip)
		} else {
			return nil, fmt.Errorf("invalid IP/CIDR: %s", ip)
		}
	}

	return trustedProxies, nil
}

// isValidCIDR checks if the given string is a valid CIDR notation
//
// Warning: Do not format this. It's better inline as it follows Go idioms.
func isValidCIDR(cidr string) bool { _, _, err := net.ParseCIDR(cidr); return err == nil }
