// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package database

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// JSONEncoder is a function type for encoding an object into JSON.
type JSONEncoder func(v any) ([]byte, error)

// JSONDecoder is a function type for decoding JSON data into an object.
type JSONDecoder func(data []byte, v any) error

// KeyFunc is a function type for extracting a key from an object.
type KeyFunc func(v any) (string, error)

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

		key, err := keyFunc(obj)
		if err != nil {
			return fmt.Errorf("failed to get key: %w", err)
		}

		pipe.Do(ctx, "JSON.SET", key, jsonPath, data)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("pipeline execution failed: %w", err)
	}

	return nil
}

// GetKeysJSONAtPipeline retrieves multiple objects from Redis using JSON.GET with a custom decoder.
// It allows specifying an optional JSON path. If no path is provided, it defaults to the root path.
func (s *service) GetKeysJSONAtPipeline(ctx context.Context, ids []string, decoder JSONDecoder, path ...string) ([]any, error) {
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

	cmds := make([]*redis.Cmd, len(ids))

	for i, id := range ids {
		cmds[i] = pipe.Do(ctx, "JSON.GET", id, jsonPath)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return nil, fmt.Errorf("pipeline execution failed: %w", err)
	}

	var objects []any
	for i, cmd := range cmds {
		result, err := cmd.Result()
		if err != nil {
			if err == redis.Nil {
				continue // Key not found, skip
			}
			return nil, fmt.Errorf("error getting key '%s': %w", ids[i], err)
		}

		data, ok := result.([]byte)
		if !ok {
			return nil, fmt.Errorf("unexpected type for key '%s'", ids[i])
		}

		var obj any
		if err := decoder(data, &obj); err != nil {
			return nil, fmt.Errorf("failed to decode object: %w", err)
		}
		objects = append(objects, obj)
	}

	return objects, nil
}

// GetRawJSONAtPipeline retrieves multiple JSON objects from Redis without decoding them.
// It allows specifying an optional JSON path. If no path is provided, it defaults to the root path.
func (s *service) GetRawJSONAtPipeline(ctx context.Context, ids []string, path ...string) (map[string][]byte, error) {
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

	cmds := make([]*redis.Cmd, len(ids))

	for i, id := range ids {
		cmds[i] = pipe.Do(ctx, "JSON.GET", id, jsonPath)
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
			return nil, fmt.Errorf("error getting key '%s': %w", ids[i], err)
		}

		data, ok := result.([]byte)
		if !ok {
			return nil, fmt.Errorf("unexpected type for key '%s'", ids[i])
		}

		results[ids[i]] = data
	}

	return results, nil
}

// DelKeysJSONAtPipeline deletes JSON objects from Redis using JSON.DEL.
// It allows specifying an optional JSON path. If no path is provided, it defaults to the root path.
func (s *service) DelKeysJSONAtPipeline(ctx context.Context, ids []string, path ...string) error {
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

	for _, id := range ids {
		pipe.Do(ctx, "JSON.DEL", id, jsonPath)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("pipeline execution failed: %w", err)
	}

	return nil
}
