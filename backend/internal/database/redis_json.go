// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.
//
// Note: If you're unsure how to use this function, consider using GetKeysAtPipeline or SetKeysAtPipeline instead.

package database

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// JSONEncoder is a function type for encoding an object into JSON.
//
// Note: Currently, this implementation uses the default JSON encoder/decoder from the standard library or other libraries like Sonic.
// Encoding with indentation (e.g., MarshalIndent) is not supported, as it's generally unnecessary for this use case.
//
// TODO: Switch to generics [T] or implement this specifically for generics [T] ?
// This change will elevate it to a top-level function if you know how to handle it.
type JSONEncoder func(v any) ([]byte, error)

// JSONDecoder is a function type for decoding JSON data into an object.
//
// TODO: Switch to generics [T] or implement this specifically for generics [T] ?
// This change will elevate it to a top-level function if you know how to handle it.
type JSONDecoder func(data []byte, v any) error

// KeyFunc is a function type for extracting one or more keys from an object.
//
// TODO: Switch to generics [T] or implement this specifically for generics [T] ?
// This change will elevate it to a top-level function if you know how to handle it.
type KeyFunc func(v any) ([]string, error)

// SetKeysJSONAtPipeline stores multiple objects in Redis using JSON.SET with a custom encoder and key extractor.
// It allows specifying an optional JSON path. If no path is provided, it defaults to the root path.
func (s *service) SetKeysJSONAtPipeline(ctx context.Context, objects []any, encoder JSONEncoder, keyFunc KeyFunc, path ...string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	pipe := s.redisClient.Pipeline()

	// Use a context with a timeout to avoid hanging indefinitely
	//
	// This is compatible/supported with Fiber's context (c.Context()), but it's recommended to use context.Background() if you're familiar with handling contexts.
	// By default, this explicitly uses "context.WithTimeout".
	ctx, cancel := context.WithTimeout(ctx, DefaultCtxTimeout)
	defer cancel()

	// Default to root path if no path is provided
	jsonPath := "$"
	if len(path) > 0 && path[0] != "" {
		jsonPath = path[0]
	}

	for _, obj := range objects {
		data, err := encoder(obj)
		if err != nil {
			return fmt.Errorf("failed to encode object: %w", err)
		}

		keys, err := keyFunc(obj)
		if err != nil {
			return fmt.Errorf("failed to get keys: %w", err)
		}

		for _, key := range keys {
			pipe.Do(ctx, "JSON.SET", key, jsonPath, data)
		}
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("pipeline execution failed: %w", err)
	}

	return nil
}

// GetKeysJSONAtPipeline retrieves multiple objects from Redis using JSON.GET with a custom decoder.
// It allows specifying an optional JSON path. If no path is provided, it defaults to the root path.
func (s *service) GetKeysJSONAtPipeline(ctx context.Context, objects []any, decoder JSONDecoder, keyFunc KeyFunc, path ...string) ([]any, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	pipe := s.redisClient.Pipeline()

	// Use a context with a timeout to avoid hanging indefinitely
	//
	// This is compatible/supported with Fiber's context (c.Context()), but it's recommended to use context.Background() if you're familiar with handling contexts.
	// By default, this explicitly uses "context.WithTimeout".
	ctx, cancel := context.WithTimeout(ctx, DefaultCtxTimeout)
	defer cancel()

	// Default to root path if no path is provided
	jsonPath := "$"
	if len(path) > 0 && path[0] != "" {
		jsonPath = path[0]
	}

	var cmds []*redis.Cmd

	for _, obj := range objects {
		keys, err := keyFunc(obj)
		if err != nil {
			return nil, fmt.Errorf("failed to get keys: %w", err)
		}

		for _, key := range keys {
			cmds = append(cmds, pipe.Do(ctx, "JSON.GET", key, jsonPath))
		}
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return nil, fmt.Errorf("pipeline execution failed: %w", err)
	}

	var results []any
	for i, cmd := range cmds {
		result, err := cmd.Result()
		if err != nil {
			if err == redis.Nil {
				continue // Key not found, skip
			}
			return nil, fmt.Errorf("error getting key '%s': %w", objects[i], err)
		}

		data, ok := result.([]byte)
		if !ok {
			return nil, fmt.Errorf("unexpected type for key '%s'", objects[i])
		}

		var obj any
		if err := decoder(data, &obj); err != nil {
			return nil, fmt.Errorf("failed to decode object: %w", err)
		}
		results = append(results, obj)
	}

	return results, nil
}

// GetRawJSONAtPipeline retrieves multiple JSON objects from Redis without decoding them.
// It allows specifying an optional JSON path. If no path is provided, it defaults to the root path.
func (s *service) GetRawJSONAtPipeline(ctx context.Context, objects []any, keyFunc KeyFunc, path ...string) (map[string][]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	pipe := s.redisClient.Pipeline()

	// Use a context with a timeout to avoid hanging indefinitely
	//
	// This is compatible/supported with Fiber's context (c.Context()), but it's recommended to use context.Background() if you're familiar with handling contexts.
	// By default, this explicitly uses "context.WithTimeout".
	ctx, cancel := context.WithTimeout(ctx, DefaultCtxTimeout)
	defer cancel()

	// Default to root path if no path is provided
	jsonPath := "$"
	if len(path) > 0 && path[0] != "" {
		jsonPath = path[0]
	}

	var cmds []*redis.Cmd
	keyIndex := make(map[int]string)

	for _, obj := range objects {
		keys, err := keyFunc(obj)
		if err != nil {
			return nil, fmt.Errorf("failed to get keys: %w", err)
		}

		for _, key := range keys {
			cmdIndex := len(cmds)
			cmds = append(cmds, pipe.Do(ctx, "JSON.GET", key, jsonPath))
			keyIndex[cmdIndex] = key
		}
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return nil, fmt.Errorf("pipeline execution failed: %w", err)
	}

	results := make(map[string][]byte)
	for i, cmd := range cmds {
		result, err := cmd.Result()
		if err != nil {
			if err == redis.Nil {
				continue // Key not found, skip
			}
			return nil, fmt.Errorf("error getting key '%s': %w", keyIndex[i], err)
		}

		data, ok := result.([]byte)
		if !ok {
			return nil, fmt.Errorf("unexpected type for key '%s'", keyIndex[i])
		}

		results[keyIndex[i]] = data
	}

	return results, nil
}

// DelKeysJSONAtPipeline deletes JSON objects from Redis using JSON.DEL.
// It allows specifying an optional JSON path. If no path is provided, it defaults to the root path.
func (s *service) DelKeysJSONAtPipeline(ctx context.Context, objects []any, keyFunc KeyFunc, path ...string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	pipe := s.redisClient.Pipeline()

	// Use a context with a timeout to avoid hanging indefinitely
	//
	// This is compatible/supported with Fiber's context (c.Context()), but it's recommended to use context.Background() if you're familiar with handling contexts.
	// By default, this explicitly uses "context.WithTimeout".
	ctx, cancel := context.WithTimeout(ctx, DefaultCtxTimeout)
	defer cancel()

	jsonPath := "$"
	if len(path) > 0 && path[0] != "" {
		jsonPath = path[0]
	}

	for _, obj := range objects {
		keys, err := keyFunc(obj)
		if err != nil {
			return fmt.Errorf("failed to get keys: %w", err)
		}

		for _, key := range keys {
			pipe.Do(ctx, "JSON.DEL", key, jsonPath)
		}
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("pipeline execution failed: %w", err)
	}

	return nil
}
