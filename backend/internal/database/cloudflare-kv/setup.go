// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package cfkv

import (
	"h0llyw00dz-template/env"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/cloudflarekv"
	_ "github.com/joho/godotenv/autoload" // auto load env
)

// FiberCloudflareKVClientConfig defines the settings needed for Fiber Cloudflare-KV client initialization.
type FiberCloudflareKVClientConfig struct {
	Key         string // The Cloudflare Auth Token (Required).
	Email       string // The Cloudflare Email (Required).
	AccountID   string // The Cloudflare Account ID (Required).
	NamespaceID string // The Cloudflare NameSpace ID (Required).
	Reset       bool   // Reset it used in fiber storage operation for each time server restarting it will reset (remove) entity data (Default: True).
}

var (
	authKey       = os.Getenv(env.CFKVKEY)
	emailCF       = os.Getenv(env.CFKVEMAIL)
	accIDCF       = os.Getenv(env.CFKVACCID)
	nameSpaceIDCF = os.Getenv(env.CFKVNAMESPACEID)
)

// InitializeCfkvStorage (Alternative Redis) initializes and returns a new Fiber Cloudflare KV storage instance
// for use with Fiber middlewares such as rate limiting, caching which it suitable with network load balancer and cheap.
//
// Recommended Usage: MYSQL -> Cloudflare-KV
//
// Note: This must implement "vice versa" method, for example when the data such as username not stored in cloudflare kv storage
// then fetch it from Mysql -> stored in this cloudflare kv storage with expiration
func (config *FiberCloudflareKVClientConfig) InitializeCfkvStorage() (fiber.Storage, error) {

	storage := cloudflarekv.New(cloudflarekv.Config{
		Key:         config.Key,
		Email:       config.Email,
		AccountID:   config.AccountID,
		NamespaceID: config.NamespaceID,
		Reset:       config.Reset,
	})

	return storage, nil

}

// InitCfkvStorage initializes the Cloudflare KV storage for Fiber using the provided configuration.
func InitCfkvStorage() (fiber.Storage, error) {
	// Prepare Fiber Cloudflare KV storage configuration
	fiberStorageConfig := &FiberCloudflareKVClientConfig{
		Key:         authKey,
		Email:       emailCF,
		AccountID:   accIDCF,
		NamespaceID: nameSpaceIDCF,
		Reset:       true, // due it only used for caching from original data where it stored in mysql, so it set true
	}

	// Initialize and return the Cloudflare KV storage using the provided configuration
	return fiberStorageConfig.InitializeCfkvStorage()
}
