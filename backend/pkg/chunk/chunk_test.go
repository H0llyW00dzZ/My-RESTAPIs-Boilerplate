// Copyright (c) 2025 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package chunk_test

import (
	"h0llyw00dz-template/backend/pkg/chunk"
	"testing"
)

func TestSplit(t *testing.T) {
	tests := []struct {
		name      string
		items     []int
		chunkSize int
		expected  [][]int
	}{
		{
			name:      "chunk size greater than length",
			items:     []int{1, 2, 3},
			chunkSize: 5,
			expected:  [][]int{{1, 2, 3}},
		},
		{
			name:      "chunk size is zero",
			items:     []int{1, 2, 3},
			chunkSize: 0,
			expected:  [][]int{{1, 2, 3}},
		},
		{
			name:      "chunk size is negative",
			items:     []int{1, 2, 3},
			chunkSize: -1,
			expected:  [][]int{{1, 2, 3}},
		},
		{
			name:      "exact chunks",
			items:     []int{1, 2, 3, 4},
			chunkSize: 2,
			expected:  [][]int{{1, 2}, {3, 4}},
		},
		{
			name:      "non-exact chunks",
			items:     []int{1, 2, 3, 4, 5},
			chunkSize: 2,
			expected:  [][]int{{1, 2}, {3, 4}, {5}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := chunk.Split(tt.items, tt.chunkSize)
			if !equal(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Helper function to compare slices of slices
func equal(a, b [][]int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) || !equalSlice(a[i], b[i]) {
			return false
		}
	}
	return true
}

// Helper function to compare slices
func equalSlice(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
